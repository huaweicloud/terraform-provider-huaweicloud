package dli

import (
	"context"
	"fmt"
	"log"
	"math"
	"strconv"
	"strings"
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
	v3queues "github.com/chnsz/golangsdk/openstack/dli/v3/queues"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

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

	MaxInstance         = "computeEngine.maxInstance"
	MaxConcurrent       = "job.maxConcurrent"
	MaxPrefetchInstance = "computeEngine.maxPrefetchInstance"
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
// @API DLI GET /v3/{project_id}/queues/{queue_name}/properties
// @API DLI PUT /v3/{project_id}/queues/{queue_name}/properties
// @API DLI DELETE /v3/{project_id}/queues/{queue_name}/properties

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
			// When the shared mode queue is offline, these parameter are mandatory for creating the exclusive mode
			// queue and must be marked as mandatory in the document, but the ability to create a shared mode queue is
			// retained (some regions may retain the ability to manage existing resources).
			"elastic_resource_pool_name": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				Description: utils.SchemaDesc(
					`The name of the elastic resource pool to which the queue belongs.`,
					utils.SchemaDescInput{
						Required: true,
					},
				),
			},
			"resource_mode": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
				ForceNew: true,
				Description: utils.SchemaDesc(
					`The queue resource mode.`,
					utils.SchemaDescInput{
						Required: true,
					},
				),
			},

			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
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
				Type:     schema.TypeInt,
				Required: true,
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

			"feature": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{queueFeatureBasic, queueFeatureAI}, false),
			},

			"tags": common.TagsForceNewSchema(),

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
			"spark_driver": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				// API
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"max_instance": {
							Type:     schema.TypeInt,
							Optional: true,
						},
						"max_concurrent": {
							Type:     schema.TypeInt,
							Optional: true,
						},
						// The type received by these parameters in the update interface is int.
						// When it is 0, it cannot be judged whether the value exists, so it is set to string.
						"max_prefetch_instance": {
							Type:     schema.TypeString,
							Optional: true,
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

			// Deprecated parameters
			"vpc_cidr": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				Description: utils.SchemaDesc(
					`The CIDR block of the queue.`,
					utils.SchemaDescInput{
						Deprecated: true,
					},
				),
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

	v3Client, err := cfg.DliV3Client(region)
	if err != nil {
		return diag.Errorf("error creating DLI V3 client: %s", err)
	}
	if v, ok := d.GetOk("scaling_policies"); ok {
		err = updateQueueScalePolicies(v3Client, elasticResourcePoolName, queueName, v.(*schema.Set))
		if err != nil {
			return diag.FromErr(err)
		}
	}

	if _, ok := d.GetOk("spark_driver"); ok {
		maxInstance, maxConcurrent, maxPrefetchInstance := getQueueProperties(d)
		if err = updateQueueSparkDriver(v3Client, queueName, maxInstance, maxConcurrent, maxPrefetchInstance); err != nil {
			return diag.Errorf("error setting properties of the queue (%s): %s", queueName, err)
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

	v3Client, err := cfg.DliV3Client(region)
	if err != nil {
		return diag.Errorf("error creating DLI V3 client: %s", err)
	}

	if elasticResourcePoolName, ok := d.GetOk("elastic_resource_pool_name"); ok {
		policies, err := getQueueScalingPolicies(v3Client, elasticResourcePoolName.(string), queueName)
		if err != nil {
			return diag.FromErr(err)
		}

		mErr = multierror.Append(mErr, d.Set("scaling_policies", policies))
	}

	sparkDriver, err := getSparkDriverByQueueName(v3Client, queueName)
	if err != nil {
		return diag.Errorf("error getting properties of the queue (%s): %s", queueName, err)
	}
	mErr = multierror.Append(mErr, d.Set("spark_driver", sparkDriver))
	return diag.FromErr(mErr.ErrorOrNil())
}

func getSparkDriverByQueueName(client *golangsdk.ServiceClient, queueName string) ([]map[string]interface{}, error) {
	opts := v3queues.ListQueuePropertyOpts{
		QueueName: queueName,
	}
	resp, err := v3queues.ListQueueProperty(client, opts)
	if err != nil {
		return nil, err
	}

	var mErr *multierror.Error
	sparkDriver := map[string]interface{}{}
	for _, property := range resp {
		switch property.Key {
		case MaxInstance:
			sparkDriver["max_instance"], err = strconv.Atoi(property.Value)
			err = multierror.Append(mErr, err)
		case MaxConcurrent:
			sparkDriver["max_concurrent"], err = strconv.Atoi(property.Value)
			err = multierror.Append(mErr, err)
		case MaxPrefetchInstance:
			sparkDriver["max_prefetch_instance"] = property.Value
		}
	}

	if mErr.ErrorOrNil() != nil {
		return nil, err
	}

	if len(sparkDriver) == 0 {
		return nil, nil
	}

	result := make([]map[string]interface{}, 0)
	result = append(result, sparkDriver)
	return result, nil
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

	if d.HasChange("spark_driver") {
		maxInstance, maxConcurrent, maxPrefetchInstance := getQueueProperties(d)
		if err = deleteQueueProperties(v3Client, queueName, maxInstance, maxConcurrent, maxPrefetchInstance); err != nil {
			return diag.Errorf("error deleting properties of the queue (%s): %s", queueName, err)
		}

		sparkDriver, err := getSparkDriverByQueueName(v3Client, queueName)
		if err != nil {
			return diag.Errorf("error getting properties of the queue (%s): %s", queueName, err)
		}

		// Set at least one property when updating.
		if len(sparkDriver) > 0 {
			if err = updateQueueSparkDriver(v3Client, queueName, maxInstance, maxConcurrent, maxPrefetchInstance); err != nil {
				return diag.Errorf("error updating properties of the queue (%s): %s", queueName, err)
			}
		}
	}
	return resourceDliQueueRead(ctx, d, meta)
}

func getQueueProperties(d *schema.ResourceData) (maxInstance, maxConcurrent int, maxPrefetchInstance string) {
	maxInstance = d.Get("spark_driver.0.max_instance").(int)
	maxConcurrent = d.Get("spark_driver.0.max_concurrent").(int)
	maxPrefetchInstance = d.Get("spark_driver.0.max_prefetch_instance").(string)
	return
}

func buildDeleteQueueProperities(maxInstance, maxConcurrent int, maxPrefetchInstance string) []string {
	result := []string{}
	if maxInstance == 0 {
		result = append(result, MaxInstance)
	}
	if maxConcurrent == 0 {
		result = append(result, MaxConcurrent)
	}
	if maxPrefetchInstance == "" {
		result = append(result, MaxPrefetchInstance)
	}

	return result
}

func deleteQueueProperties(client *golangsdk.ServiceClient, queueName string, maxInstance, maxCon int, maxPrefetchInstance string) error {
	deleteOpts := buildDeleteQueueProperities(maxInstance, maxCon, maxPrefetchInstance)
	if len(deleteOpts) == 0 {
		return nil
	}

	resp, err := v3queues.DeleteQueueProperties(client, queueName, deleteOpts)
	if err != nil {
		return err
	}

	if !resp.IsSuccess {
		return fmt.Errorf(resp.Message)
	}
	return nil
}

func buildScaleActionParam(oldValue, newValue int) string {
	if oldValue > newValue {
		return actionScaleIn
	}
	return actionScaleOut
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

func updateQueueSparkDriver(client *golangsdk.ServiceClient, queueName string, maxInstance, maxCon int, maxPrefetchInstance string) error {
	opts := v3queues.Property{
		MaxInstance:   maxInstance,
		MaxConcurrent: maxCon,
	}
	if maxPrefetchInstance != "" {
		num, err := strconv.Atoi(maxPrefetchInstance)
		if err != nil {
			return fmt.Errorf("the string (%s) cannot be converted to number", maxPrefetchInstance)
		}
		opts.MaxPrefetchInstance = utils.Int(num)
	}
	resp, err := v3queues.UpdateQueueProperty(client, queueName, opts)
	if err != nil {
		return err
	}

	if !resp.IsSuccess {
		return fmt.Errorf(resp.Message)
	}
	return nil
}

func resourceQueueImportState(_ context.Context, d *schema.ResourceData, _ interface{}) ([]*schema.ResourceData, error) {
	var (
		mErr *multierror.Error

		importedId = d.Id()
		parts      = strings.Split(importedId, "/")
	)
	switch len(parts) {
	case 1:
		mErr = multierror.Append(mErr, d.Set("name", parts[0]))
	case 2:
		mErr = multierror.Append(mErr,
			d.Set("queue_type", parts[0]),
			d.Set("name", parts[1]),
		)
	default:
		return nil, fmt.Errorf("invalid format specified for import ID, want '<queue_type>/<name>' or '<name>', but got '%s'",
			importedId)
	}
	return []*schema.ResourceData{d}, mErr.ErrorOrNil()
}
