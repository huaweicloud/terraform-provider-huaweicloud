package dli

import (
	"context"
	"fmt"
	"log"
	"math"
	"regexp"
	"time"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/chnsz/golangsdk"
	"github.com/chnsz/golangsdk/openstack/common/tags"
	"github.com/chnsz/golangsdk/openstack/dli/v1/queues"
	"github.com/chnsz/golangsdk/openstack/dli/v3/elasticresourcepool"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

var regexp4Name = regexp.MustCompile(`^[a-z0-9_]{1,128}$`)

const (
	CU16                  = 16
	CU64                  = 64
	CU256                 = 256
	resourceModeShared    = 0
	resourceModeExclusive = 1

	QueueTypeSQL         = "sql"
	QueueTypeGeneral     = "general"
	queueFeatureBasic    = "basic"
	queueFeatureAI       = "ai"
	queuePlatformX86     = "x86_64"
	queuePlatformAARCH64 = "aarch64"

	actionRestart  = "restart"
	actionScaleOut = "scale_out"
	actionScaleIn  = "scale_in"
)

// @API DLI POST /v1.0/{project_id}/queues
// @API DLI GET /v1.0/{project_id}/queues/{queue_name}
// @API DLI GET /v1.0/{project_id}/queues
// @API DLI PUT /v1.0/{project_id}/queues/{queue_name}/action
// @API DLI PUT /v1.0/{project_id}/queues/{queue_name}
// @API DLI DELETE /v1.0/{project_id}/queues/{queue_name}
// @API DLI POST /v3/{project_id}/elastic-resource-pools/{elastic_resource_pool_name}/queues
// @API DLI GET /v3/{project_id}/elastic-resource-pools/{elastic_resource_pool_name}/queues
// @API DLI PUT /v3/{project_id}/elastic-resource-pools/{elastic_resource_pool_name}/queues/{queue_name}
func ResourceDliQueue() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceDliQueueCreate,
		ReadContext:   resourceDliQueueRead,
		UpdateContext: resourceDliQueueUpdate,
		DeleteContext: resourceDliQueueDelete,

		Importer: &schema.ResourceImporter{
			StateContext: resourceQueueImportState,
		},

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},

			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
				ValidateFunc: validation.StringMatch(regexp4Name,
					"only contain digits, lower letters, and underscores (_)"),
			},

			"queue_type": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				Default:      QueueTypeSQL,
				ValidateFunc: validation.StringInSlice([]string{QueueTypeSQL, QueueTypeGeneral}, false),
			},

			"description": {
				Type:     schema.TypeString,
				Computed: true,
				Optional: true,
				ForceNew: true,
			},

			"cu_count": {
				Type:         schema.TypeInt,
				Required:     true,
				ValidateFunc: validCuCount,
			},

			"enterprise_project_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Computed: true,
			},

			"platform": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				Default:      queuePlatformX86,
				ValidateFunc: validation.StringInSlice([]string{queuePlatformX86, queuePlatformAARCH64}, false),
			},

			"resource_mode": {
				Type:         schema.TypeInt,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.IntInSlice([]int{resourceModeShared, resourceModeExclusive}),
			},

			"feature": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{queueFeatureBasic, queueFeatureAI}, false),
			},

			"tags": {
				Type: schema.TypeMap,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Optional: true,
				ForceNew: true,
			},

			"vpc_cidr": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"elastic_resource_pool_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"scaling_policies": {
				Type:         schema.TypeSet,
				Optional:     true,
				Computed:     true,
				RequiredWith: []string{"elastic_resource_pool_name"},
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"priority": {
							Type:     schema.TypeInt,
							Required: true,
						},
						"impact_start_time": {
							Type:     schema.TypeString,
							Required: true,
						},
						"impact_stop_time": {
							Type:     schema.TypeString,
							Required: true,
						},
						"min_cu": {
							Type:     schema.TypeInt,
							Required: true,
						},
						"max_cu": {
							Type:     schema.TypeInt,
							Required: true,
						},
					},
				},
			},
			"create_time": {
				Type:     schema.TypeInt,
				Computed: true,
			},

			"management_subnet_cidr": {
				Type:       schema.TypeString,
				Optional:   true,
				ForceNew:   true,
				Deprecated: "management_subnet_cidr is Deprecated",
			},

			"subnet_cidr": {
				Type:       schema.TypeString,
				Optional:   true,
				ForceNew:   true,
				Deprecated: "subnet_cidr is Deprecated",
			},
		},

		Timeouts: &schema.ResourceTimeout{
			Update: schema.DefaultTimeout(45 * time.Minute),
		},
	}
}

func resourceDliQueueCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	dliClient, err := cfg.DliV1Client(region)
	if err != nil {
		return diag.Errorf("creating dli client failed: %s", err)
	}

	queueName := d.Get("name").(string)

	log.Printf("[DEBUG] create dli queues queueName: %s", queueName)
	elasticResourcePoolName := d.Get("elastic_resource_pool_name").(string)
	createOpts := queues.CreateOpts{
		QueueName:               queueName,
		QueueType:               d.Get("queue_type").(string),
		Description:             d.Get("description").(string),
		CuCount:                 d.Get("cu_count").(int),
		EnterpriseProjectId:     cfg.GetEnterpriseProjectID(d),
		Platform:                d.Get("platform").(string),
		ResourceMode:            d.Get("resource_mode").(int),
		Feature:                 d.Get("feature").(string),
		Tags:                    assembleTagsFromRecource("tags", d),
		ElasticResourcePoolName: elasticResourcePoolName,
	}

	log.Printf("[DEBUG] create dli queues using parameters: %+v", createOpts)
	createResult := queues.Create(dliClient, createOpts)
	if createResult.Err != nil {
		return diag.Errorf("create dli queues failed: %s", createResult.Err)
	}

	// The resource ID (queue name) at this time is only used as a mark the resource, and the value will be refreshed
	// in the READ method.
	d.SetId(queueName)

	// This is a workaround to avoid issue: the queue is assigning, which is not available
	time.Sleep(4 * time.Minute) // lintignore:R018

	if v, ok := d.GetOk("vpc_cidr"); ok {
		err = updateVpcCidrOfQueue(dliClient, queueName, v.(string))
		if err != nil {
			return diag.Errorf("update cidr failed when creating dli queues: %s", err)
		}
	}

	if v, ok := d.GetOk("scaling_policies"); ok {
		v3Client, err := cfg.DliV3Client(region)
		if err != nil {
			return diag.Errorf("error creating DLI V3 client: %s", err)
		}
		err = updateQueueScalePolicies(v3Client, elasticResourcePoolName, queueName, v.(*schema.Set))
		if err != nil {
			return diag.FromErr(err)
		}
	}

	return resourceDliQueueRead(ctx, d, meta)
}

func assembleTagsFromRecource(key string, d *schema.ResourceData) []tags.ResourceTag {
	if v, ok := d.GetOk(key); ok {
		tagRaw := v.(map[string]interface{})
		taglist := utils.ExpandResourceTags(tagRaw)
		return taglist
	}
	return nil
}

func resourceDliQueueRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	client, err := cfg.DliV1Client(region)
	if err != nil {
		return diag.Errorf("error creating DliV1Client, err=%s", err)
	}

	queueName := d.Get("name").(string)

	queryOpts := queues.ListOpts{
		QueueType: d.Get("queue_type").(string),
	}

	log.Printf("[DEBUG] query dli queues using parameters: %+v", queryOpts)

	queryAllResult := queues.List(client, queryOpts)
	if queryAllResult.Err != nil {
		return diag.Errorf("query queues failed: %s", queryAllResult.Err)
	}

	// filter by queue_name
	queueDetail, err := filterByQueueName(queryAllResult.Body, queueName)
	if err != nil {
		return common.CheckDeletedDiag(d, err, "DLI queue")
	}
	d.SetId(queueDetail.ResourceId)

	log.Printf("[DEBUG]The detail of queue from SDK:%+v", queueDetail)

	mErr := multierror.Append(
		d.Set("region", region),
		d.Set("name", queueDetail.QueueName),
		d.Set("queue_type", queueDetail.QueueType),
		d.Set("description", queueDetail.Description),
		d.Set("cu_count", queueDetail.CuCount),
		d.Set("enterprise_project_id", utils.StringIgnoreEmpty(queueDetail.EnterpriseProjectId)),
		d.Set("platform", queueDetail.Platform),
		d.Set("resource_mode", queueDetail.ResourceMode),
		d.Set("feature", queueDetail.Feature),
		d.Set("create_time", queueDetail.CreateTime),
		d.Set("vpc_cidr", queueDetail.CidrInVpc),
		d.Set("elastic_resource_pool_name", queueDetail.ElasticResourcePoolName),
	)

	if elasticResourcePoolName, ok := d.GetOk("elastic_resource_pool_name"); ok {
		v3Client, err := cfg.DliV3Client(region)
		if err != nil {
			return diag.Errorf("error creating DLI V3 client: %s", err)
		}

		policies, err := getQueueScalingPolicies(v3Client, elasticResourcePoolName.(string), queueName)
		if err != nil {
			return diag.FromErr(err)
		}

		mErr = multierror.Append(mErr, d.Set("scaling_policies", policies))
	}

	return diag.FromErr(mErr.ErrorOrNil())
}

func filterByQueueName(body interface{}, queueName string) (r *queues.Queue, err error) {
	if queueList, ok := body.(*queues.ListResult); ok {
		log.Printf("[DEBUG]The list of queue from SDK:%+v", queueList)
		if !queueList.IsSuccess {
			return nil, fmt.Errorf("unable to query the queue list: %s", queueList.Message)
		}

		for _, v := range queueList.Queues {
			if v.QueueName == queueName {
				return &v, nil
			}
		}
		return nil, golangsdk.ErrDefault404{}
	}

	return nil, fmt.Errorf("sdk-client response type is wrong, expect type:*queues.ListResult,acutal Type:%T",
		body)
}

func getQueueScalingPoliciesByName(client *golangsdk.ServiceClient, poolName, queueName string) ([]elasticresourcepool.QueueScalingPolicy, error) {
	opts := elasticresourcepool.ListElasticResourcePoolQueuesOpts{
		ElasticResourcePoolName: poolName,
		QueueName:               queueName,
		// The API document states that the default value of limit is 100, but if offset is specified without limit, paging will not take effect.
		Limit: 100,
	}
	queueList, err := elasticresourcepool.ListElasticResourcePoolQueues(client, opts)
	if err != nil {
		return nil, err
	}

	// The `queue_name` parameter supports fuzzy search.
	for _, queue := range queueList {
		if queue.QueueName == queueName {
			return queue.QueueScalingPolicies, nil
		}
	}

	return nil, fmt.Errorf("unable to find the queue (%s) under the elastic resource pool (%s)",
		queueName, poolName)
}

func getQueueScalingPolicies(client *golangsdk.ServiceClient, elasticResourcePoolName, queueName string) ([]interface{}, error) {
	policies, err := getQueueScalingPoliciesByName(client, elasticResourcePoolName, queueName)
	if err != nil {
		return nil, err
	}

	result := make([]interface{}, len(policies))
	for i, policy := range policies {
		result[i] = map[string]interface{}{
			"priority":          policy.Priority,
			"impact_start_time": policy.ImpactStartTime,
			"impact_stop_time":  policy.ImpactStopTime,
			"min_cu":            policy.MinCu,
			"max_cu":            policy.MaxCu,
		}
	}

	return result, nil
}

func resourceDliQueueDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	client, err := cfg.DliV1Client(region)
	if err != nil {
		return diag.Errorf("error creating DliV1Client, err=%s", err)
	}

	queueName := d.Get("name").(string)
	log.Printf("[DEBUG] Deleting dli Queue %q", queueName)

	result := queues.Delete(client, queueName)
	if result.Err != nil {
		return diag.Errorf("error deleting dli Queue %q, err=%s", queueName, result.Err)
	}

	return nil
}

// support cu_count scaling
func resourceDliQueueUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	client, err := cfg.DliV1Client(region)
	if err != nil {
		return diag.Errorf("error creating DliV1Client: %s", err)
	}

	queueName := d.Get("name").(string)
	opt := queues.ActionOpts{
		QueueName: queueName,
	}
	if d.HasChange("cu_count") {
		oldValue, newValue := d.GetChange("cu_count")
		cuChange := newValue.(int) - oldValue.(int)

		opt.CuCount = int(math.Abs(float64(cuChange)))
		opt.Action = buildScaleActionParam(oldValue.(int), newValue.(int))

		log.Printf("[DEBUG]DLI queue Update Option: %#v", opt)
		result := queues.ScaleOrRestart(client, opt)
		if result.Err != nil {
			return diag.Errorf("update dli queues failed, queueName=%s, error:%s", queueName, result.Err)
		}

		updateStateConf := &resource.StateChangeConf{
			Pending: []string{fmt.Sprintf("%d", oldValue)},
			Target:  []string{fmt.Sprintf("%d", newValue)},
			Refresh: func() (interface{}, string, error) {
				getResult := queues.Get(client, queueName)
				queueDetail := getResult.Body.(*queues.Queue4Get)
				return getResult, fmt.Sprintf("%d", queueDetail.CuCount), nil
			},
			Timeout:      d.Timeout(schema.TimeoutUpdate),
			Delay:        30 * time.Second,
			PollInterval: 20 * time.Second,
		}
		_, err = updateStateConf.WaitForStateContext(ctx)
		if err != nil {
			return diag.Errorf("error waiting for dli.queue (%s) to be scale: %s", queueName, err)
		}
	}

	if d.HasChange("vpc_cidr") {
		cidr := d.Get("vpc_cidr").(string)
		err = updateVpcCidrOfQueue(client, queueName, cidr)
		if err != nil {
			return diag.Errorf("update cidr failed when updating dli queues: %s", err)
		}
	}

	v3Client, err := cfg.DliV3Client(region)
	if err != nil {
		return diag.Errorf("error creating DLI V3 client: %s", err)
	}
	if d.HasChange("elastic_resource_pool_name") {
		oldVal, newVal := d.GetChange("elastic_resource_pool_name")
		if oldVal != "" {
			return diag.Errorf("error the queue has been associate with an Elastic resopurce pool")
		}

		associateQueueToElasticResourcePoolOpts := elasticresourcepool.AssociateQueueOpts{
			ElasticResourcePoolName: newVal.(string),
			QueueName:               queueName,
		}
		err = associateQueueToElasticResourcePool(v3Client, associateQueueToElasticResourcePoolOpts)
		if err != nil {
			return diag.Errorf("error associate queue to elastic resopurce pool: %s", err)
		}
	}

	if d.HasChange("scaling_policies") {
		elasticResourcePoolName := d.Get("elastic_resource_pool_name").(string)
		queuePolicies := d.Get("scaling_policies").(*schema.Set)
		if err := updateQueueScalePolicies(v3Client, elasticResourcePoolName, queueName, queuePolicies); err != nil {
			return diag.FromErr(err)
		}
	}

	return resourceDliQueueRead(ctx, d, meta)
}

func buildScaleActionParam(oldValue, newValue int) string {
	if oldValue > newValue {
		return actionScaleIn
	}
	return actionScaleOut
}

func validCuCount(val interface{}, key string) (warns []string, errs []error) {
	diviNum := 16
	warns, errs = validation.IntAtLeast(diviNum)(val, key)
	if len(errs) > 0 {
		return warns, errs
	}
	return validation.IntDivisibleBy(diviNum)(val, key)
}

func updateVpcCidrOfQueue(client *golangsdk.ServiceClient, queueName, cidr string) error {
	resp, err := queues.UpdateCidr(client, queueName, queues.UpdateCidrOpts{Cidr: cidr})
	if err != nil {
		return err
	}
	if !resp.IsSuccess {
		return fmt.Errorf("unable to update the VPC CIDR: %s", resp.Message)
	}
	return err
}

func buildQueueScalingPolicies(policies *schema.Set) []elasticresourcepool.QueueScalingPolicy {
	if policies.Len() == 0 {
		return nil
	}

	result := make([]elasticresourcepool.QueueScalingPolicy, policies.Len())
	for i, val := range policies.List() {
		policy := val.(map[string]interface{})
		result[i] = elasticresourcepool.QueueScalingPolicy{
			Priority:        policy["priority"].(int),
			ImpactStartTime: policy["impact_start_time"].(string),
			ImpactStopTime:  policy["impact_stop_time"].(string),
			MinCu:           policy["min_cu"].(int),
			MaxCu:           policy["max_cu"].(int),
		}
	}
	return result
}

func updateQueueScalePolicies(client *golangsdk.ServiceClient, elasticResourcePoolName, queueName string, policys *schema.Set) error {
	opts := elasticresourcepool.UpdateQueuePolicyOpts{
		ElasticResourcePoolName: elasticResourcePoolName,
		QueueName:               queueName,
		QueueScalingPolicies:    buildQueueScalingPolicies(policys),
	}
	resp, err := elasticresourcepool.UpdateElasticResourcePoolQueuePolicy(client, opts)
	if err != nil {
		return fmt.Errorf("error updating scaling policies of the queue (%s) associated with the elastic resource pool (%s): %s",
			elasticResourcePoolName, queueName, err)
	}

	if !resp.IsSuccess {
		return fmt.Errorf("unable to update scaling policies to the queue (%s): %s", queueName, resp.Message)
	}
	return err
}

func associateQueueToElasticResourcePool(client *golangsdk.ServiceClient, opts elasticresourcepool.AssociateQueueOpts) error {
	resp, err := elasticresourcepool.AssociateQueue(client, opts)
	if err != nil {
		return err
	}
	if !resp.IsSuccess {
		return fmt.Errorf("unable to associate the queue to the elastic resource pool: %s", resp.Message)
	}
	return err
}

func resourceQueueImportState(_ context.Context, d *schema.ResourceData, _ interface{}) ([]*schema.ResourceData, error) {
	err := d.Set("name", d.Id())
	if err != nil {
		return []*schema.ResourceData{d}, fmt.Errorf("error saving resource name of the DLI queue: %s", err)
	}
	return []*schema.ResourceData{d}, nil
}
