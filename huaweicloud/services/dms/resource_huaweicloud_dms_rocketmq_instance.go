package dms

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"regexp"
	"strings"
	"time"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/jmespath/go-jmespath"

	"github.com/chnsz/golangsdk"
	"github.com/chnsz/golangsdk/openstack/common/tags"
	"github.com/chnsz/golangsdk/openstack/eps/v1/enterpriseprojects"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

type dmsError struct {
	ErrorCode string `json:"error_code"`
	ErrorMsg  string `json:"error_msg"`
}

// @API RocketMQ POST /v2/{engine}/{project_id}/instances
// @API RocketMQ PUT /v2/{project_id}/instances/{instance_id}
// @API RocketMQ GET /v2/{project_id}/instances/{instance_id}
// @API RocketMQ DELETE /v2/{project_id}/instances/{instance_id}
// @API RocketMQ POST /v2/{project_id}/rocketmq/{instance_id}/tags/action
// @API RocketMQ GET /v2/{project_id}/rocketmq/{instance_id}/tags
// @API RocketMQ POST /v2/{project_id}/instances/{instance_id}/crossvpc/modify
// @API RocketMQ GET /v2/{project_id}/instances/{instance_id}/tasks
// @API RocketMQ POST /v2/{engine}/{project_id}/instances/{instance_id}/extend
// @API EIP GET /v1/{project_id}/publicips
// @API BSS GET /v2/orders/customer-orders/details/{order_id}
// @API BSS POST /v2/orders/subscriptions/resources/autorenew/{instance_id}
// @API BSS DELETE /v2/orders/subscriptions/resources/autorenew/{instance_id}
// @API BSS POST /v2/orders/subscriptions/resources/unsubscribe
func ResourceDmsRocketMQInstance() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceDmsRocketMQInstanceCreate,
		UpdateContext: resourceDmsRocketMQInstanceUpdate,
		ReadContext:   resourceDmsRocketMQInstanceRead,
		DeleteContext: resourceDmsRocketMQInstanceDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(50 * time.Minute),
			Update: schema.DefaultTimeout(50 * time.Minute),
			Delete: schema.DefaultTimeout(15 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `Specifies the name of the DMS RocketMQ instance`,
				ValidateFunc: validation.All(
					validation.StringMatch(regexp.MustCompile(`^[A-Za-z-_0-9]*$`),
						"An instance name starts with a letter and can contain only letters,"+
							"digits, underscores (_), and hyphens (-)"),
					validation.StringLenBetween(4, 64),
				),
			},
			"engine_version": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: `Specifies the version of the RocketMQ engine.`,
			},
			"storage_space": {
				Type:         schema.TypeInt,
				Required:     true,
				Description:  `Specifies the message storage capacity, Unit: GB.`,
				ValidateFunc: validation.IntBetween(300, 3000),
			},
			"vpc_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: `Specifies the ID of a VPC`,
			},
			"subnet_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: `Specifies the ID of a subnet`,
			},
			"security_group_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `Specifies the ID of a security group`,
			},
			"availability_zones": {
				Type:        schema.TypeSet,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Required:    true,
				ForceNew:    true,
				Set:         schema.HashString,
				Description: `Specifies the list of availability zone names`,
			},
			"flavor_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `Specifies a product ID`,
			},
			"storage_spec_code": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: `Specifies the storage I/O specification`,
			},
			"description": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				Description:  `Specifies the description of the DMS RocketMQ instance.`,
				ValidateFunc: validation.StringLenBetween(0, 1024),
			},
			"ssl_enable": {
				Type:        schema.TypeBool,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: `Specifies whether the RocketMQ SASL_SSL is enabled.`,
			},
			"ipv6_enable": {
				Type:        schema.TypeBool,
				Optional:    true,
				ForceNew:    true,
				Description: `Specifies whether to support IPv6`,
			},
			"enable_publicip": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: `Specifies whether to enable public access.`,
			},
			"publicip_id": {
				Type:             schema.TypeString,
				Optional:         true,
				Computed:         true,
				DiffSuppressFunc: utils.SuppressStringSepratedByCommaDiffs,
				Description:      `Specifies the ID of the EIP bound to the instance.`,
			},
			"broker_num": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: `Specifies the broker numbers.`,
			},
			"enterprise_project_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `Specifies the enterprise project id of the instance.`,
			},
			"enable_acl": {
				Type:        schema.TypeBool,
				Optional:    true,
				Computed:    true,
				Description: `Specifies whether access control is enabled.`,
			},
			"status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Indicates the status of the DMS RocketMQ instance.`,
			},
			"type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Indicates the DMS RocketMQ instance type. Value: cluster.`,
			},
			"specification": {
				Type:     schema.TypeString,
				Computed: true,
				Description: `Indicates the instance specification. For a cluster DMS RocketMQ instance, VM specifications
  and the number of nodes are returned.`,
			},
			"maintain_begin": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Indicates the time at which the maintenance window starts. The format is HH:mm:ss.`,
			},
			"maintain_end": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Indicates the time at which the maintenance window ends. The format is HH:mm:ss.`,
			},
			"used_storage_space": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: `Indicates the used message storage space. Unit: GB.`,
			},
			"publicip_address": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Indicates the public IP address.`,
			},
			"node_num": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: `Indicates the node quantity.`,
			},
			"new_spec_billing_enable": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: `Indicates whether billing based on new specifications is enabled.`,
			},
			"namesrv_address": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Indicates the metadata address.`,
			},
			"broker_address": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Indicates the service data address.`,
			},
			"public_namesrv_address": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Indicates the public network metadata address.`,
			},
			"public_broker_address": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Indicates the public network service data address.`,
			},
			"resource_spec_code": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Indicates the resource specifications.`,
			},
			"retention_policy": {
				Type:        schema.TypeBool,
				Optional:    true,
				Computed:    true,
				Description: `Specifies whether access control is enabled.`,
				Deprecated:  "Use 'enable_acl' instead",
			},
			"tags": common.TagsSchema(),
			"cross_vpc_accesses": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"advertised_ip": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"listener_ip": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"port": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"port_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						// Typo, it is only kept in the code, will not be shown in the docs.
						"lisenter_ip": {
							Type:       schema.TypeString,
							Computed:   true,
							Deprecated: "typo in lisenter_ip, please use \"listener_ip\" instead.",
						},
					},
				},
			},
			"charging_mode": common.SchemaChargingMode(nil),
			"period_unit":   common.SchemaPeriodUnit(nil),
			"period":        common.SchemaPeriod(nil),
			"auto_renew":    common.SchemaAutoRenewUpdatable(nil),
		},
	}
}

func resourceDmsRocketMQInstanceCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	// createRocketmqInstance: create DMS rocketmq instance
	var (
		createRocketmqInstanceHttpUrl = "v2/reliability/{project_id}/instances"
		createRocketmqInstanceProduct = "dmsv2"
	)
	createRocketmqInstanceClient, err := cfg.NewServiceClient(createRocketmqInstanceProduct, region)
	if err != nil {
		return diag.Errorf("error creating DmsRocketMQInstance Client: %s", err)
	}

	createRocketmqInstancePath := createRocketmqInstanceClient.Endpoint + createRocketmqInstanceHttpUrl
	createRocketmqInstancePath = strings.ReplaceAll(createRocketmqInstancePath, "{project_id}",
		createRocketmqInstanceClient.ProjectID)

	createRocketmqInstanceOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		OkCodes: []int{
			200,
		},
	}
	var availableZones []string
	// convert the codes of the availability zone into ids
	azCodes := d.Get("availability_zones").(*schema.Set)
	availableZones, err = getAvailableZoneIDByCode(cfg, region, azCodes.List())
	if err != nil {
		return diag.FromErr(err)
	}
	createRocketmqInstanceOpt.JSONBody = utils.RemoveNil(buildCreateRocketmqInstanceBodyParams(d, cfg, availableZones))
	createRocketmqInstanceResp, err := createRocketmqInstanceClient.Request("POST", createRocketmqInstancePath,
		&createRocketmqInstanceOpt)
	if err != nil {
		return diag.Errorf("error creating DmsRocketMQInstance: %s", err)
	}
	createRocketmqInstanceRespBody, err := utils.FlattenResponse(createRocketmqInstanceResp)
	if err != nil {
		return diag.FromErr(err)
	}

	id, err := jmespath.Search("instance_id", createRocketmqInstanceRespBody)
	if err != nil {
		return diag.Errorf("error creating DmsRocketMQInstance: ID is not found in API response")
	}

	var delayTime time.Duration = 300
	if chargingMode, ok := d.GetOk("charging_mode"); ok && chargingMode == "prePaid" {
		err = waitForRocketMQOrderComplete(ctx, d, cfg, createRocketmqInstanceClient, id.(string))
		if err != nil {
			return diag.FromErr(err)
		}
		delayTime = 5
	}

	stateConf := &resource.StateChangeConf{
		Pending:      []string{"CREATING"},
		Target:       []string{"RUNNING"},
		Refresh:      rocketmqInstanceStateRefreshFunc(createRocketmqInstanceClient, id.(string)),
		Timeout:      d.Timeout(schema.TimeoutCreate),
		Delay:        delayTime * time.Second,
		PollInterval: 10 * time.Second,
	}

	_, err = stateConf.WaitForStateContext(ctx)
	if err != nil {
		return diag.Errorf("error waiting for instance (%s) to create: %s", id.(string), err)
	}

	d.SetId(id.(string))

	if _, ok := d.GetOk("cross_vpc_accesses"); ok {
		if err = updateCrossVpcAccess(ctx, createRocketmqInstanceClient, d); err != nil {
			return diag.Errorf("failed to update default advertised IP: %v", err)
		}
	}

	// set tags
	tagRaw := d.Get("tags").(map[string]interface{})
	if len(tagRaw) > 0 {
		tagList := utils.ExpandResourceTags(tagRaw)
		if tagErr := tags.Create(createRocketmqInstanceClient, "rocketmq", id.(string), tagList).ExtractErr(); tagErr != nil {
			return diag.Errorf("error setting tags of RocketMQ %s: %s", id.(string), tagErr)
		}
	}

	return resourceDmsRocketMQInstanceRead(ctx, d, meta)
}

func waitForRocketMQOrderComplete(ctx context.Context, d *schema.ResourceData, cfg *config.Config,
	client *golangsdk.ServiceClient, instanceID string) error {
	region := cfg.GetRegion(d)
	orderId, err := getRocketMQInstanceOrderId(ctx, d, client, instanceID)
	if err != nil {
		return err
	}

	if orderId == "" {
		log.Printf("[WARN] error get order id by instance ID: %s", instanceID)
		return nil
	}

	bssClient, err := cfg.BssV2Client(region)
	if err != nil {
		return fmt.Errorf("error creating BSS v2 client: %s", err)
	}
	// wait for order success
	err = common.WaitOrderComplete(ctx, bssClient, orderId, d.Timeout(schema.TimeoutCreate))
	if err != nil {
		return err
	}
	_, err = common.WaitOrderResourceComplete(ctx, bssClient, orderId, d.Timeout(schema.TimeoutCreate))
	if err != nil {
		return fmt.Errorf("error waiting for RocketMQ order resource %s complete: %s", orderId, err)
	}
	return nil
}

func getRocketMQInstanceOrderId(ctx context.Context, d *schema.ResourceData, client *golangsdk.ServiceClient,
	instanceID string) (string, error) {
	stateConf := &resource.StateChangeConf{
		Pending:      []string{"EMPTY"},
		Target:       []string{"CREATING"},
		Refresh:      rocketMQInstanceCreatingFunc(client, instanceID),
		Timeout:      d.Timeout(schema.TimeoutCreate),
		Delay:        500 * time.Millisecond,
		PollInterval: 500 * time.Millisecond,
	}
	orderId, err := stateConf.WaitForStateContext(ctx)
	if err != nil {
		return "", fmt.Errorf("error waiting for RocketMQ instance (%s) to creating: %s", instanceID, err)
	}
	return orderId.(string), nil
}

func rocketMQInstanceCreatingFunc(client *golangsdk.ServiceClient, instanceID string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		var (
			getRocketmqInstanceHttpUrl = "v2/{project_id}/instances/{instance_id}"
		)

		getRocketmqInstancePath := client.Endpoint + getRocketmqInstanceHttpUrl
		getRocketmqInstancePath = strings.ReplaceAll(getRocketmqInstancePath, "{project_id}", client.ProjectID)
		getRocketmqInstancePath = strings.ReplaceAll(getRocketmqInstancePath, "{instance_id}", fmt.Sprintf("%v", instanceID))

		getRocketmqInstanceOpt := golangsdk.RequestOpts{
			KeepResponseBody: true,
			OkCodes: []int{
				200,
			},
		}
		getRocketmqInstanceResp, err := client.Request("GET", getRocketmqInstancePath, &getRocketmqInstanceOpt)

		if err != nil {
			if errCode, ok := err.(golangsdk.ErrDefault404); ok {
				var rocketMQError dmsError
				err = json.Unmarshal(errCode.Body, &rocketMQError)
				if err != nil {
					return nil, "", fmt.Errorf("error get DmsRocketMQInstance: error format error: %s", err)
				}
				if rocketMQError.ErrorCode == "DMS.00404022" {
					return getRocketmqInstanceResp, "EMPTY", nil
				}
			}
			return nil, "", fmt.Errorf("error retrieving DmsRocketMQInstance: %s", err)
		}

		res, err := utils.FlattenResponse(getRocketmqInstanceResp)
		if err != nil {
			return nil, "", err
		}
		orderID := utils.PathSearch("order_id", res, "")
		return orderID, "CREATING", nil
	}
}

func buildCreateRocketmqInstanceBodyParams(d *schema.ResourceData, cfg *config.Config,
	availableZones []string) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"name":                  utils.ValueIgnoreEmpty(d.Get("name")),
		"enable_acl":            utils.ValueIgnoreEmpty(d.Get("enable_acl")),
		"description":           utils.ValueIgnoreEmpty(d.Get("description")),
		"engine":                "reliability",
		"engine_version":        utils.ValueIgnoreEmpty(d.Get("engine_version")),
		"storage_space":         utils.ValueIgnoreEmpty(d.Get("storage_space")),
		"vpc_id":                utils.ValueIgnoreEmpty(d.Get("vpc_id")),
		"subnet_id":             utils.ValueIgnoreEmpty(d.Get("subnet_id")),
		"security_group_id":     utils.ValueIgnoreEmpty(d.Get("security_group_id")),
		"available_zones":       availableZones,
		"product_id":            utils.ValueIgnoreEmpty(d.Get("flavor_id")),
		"ssl_enable":            utils.ValueIgnoreEmpty(d.Get("ssl_enable")),
		"storage_spec_code":     utils.ValueIgnoreEmpty(d.Get("storage_spec_code")),
		"ipv6_enable":           utils.ValueIgnoreEmpty(d.Get("ipv6_enable")),
		"enable_publicip":       utils.ValueIgnoreEmpty(d.Get("enable_publicip")),
		"publicip_id":           utils.ValueIgnoreEmpty(d.Get("publicip_id")),
		"broker_num":            utils.ValueIgnoreEmpty(d.Get("broker_num")),
		"enterprise_project_id": utils.ValueIgnoreEmpty(common.GetEnterpriseProjectID(d, cfg)),
	}
	if chargingMode, ok := d.GetOk("charging_mode"); ok && chargingMode == "prePaid" {
		bodyParams["bss_param"] = buildCreateRocketmqInstanceBodyBssParams(d)
	}
	return bodyParams
}

func buildCreateRocketmqInstanceBodyBssParams(d *schema.ResourceData) map[string]interface{} {
	var autoRenew bool
	if d.Get("auto_renew").(string) == "true" {
		autoRenew = true
	}
	isAutoPay := true
	bodyParams := map[string]interface{}{
		"charging_mode": utils.ValueIgnoreEmpty(d.Get("charging_mode")),
		"period_type":   utils.ValueIgnoreEmpty(d.Get("period_unit")),
		"period_num":    utils.ValueIgnoreEmpty(d.Get("period")),
		"is_auto_renew": &autoRenew,
		"is_auto_pay":   &isAutoPay,
	}
	return bodyParams
}

func resourceDmsRocketMQInstanceUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	instanceId := d.Id()

	updateRocketmqInstanceHasChanges := []string{
		"name",
		"description",
		"security_group_id",
		"retention_policy",
		"enable_acl",
		"enable_publicip",
		"publicip_id",
	}

	// updateRocketmqInstance: update DMS rocketmq instance
	var (
		updateRocketmqInstanceHttpUrl = "v2/{project_id}/instances/{instance_id}"
		updateRocketmqInstanceProduct = "dmsv2"
	)
	updateRocketmqInstanceClient, err := cfg.NewServiceClient(updateRocketmqInstanceProduct, region)
	if err != nil {
		return diag.Errorf("error creating DMS Client: %s", err)
	}

	updateRocketmqInstancePath := updateRocketmqInstanceClient.Endpoint + updateRocketmqInstanceHttpUrl
	updateRocketmqInstancePath = strings.ReplaceAll(updateRocketmqInstancePath, "{project_id}", updateRocketmqInstanceClient.ProjectID)
	updateRocketmqInstancePath = strings.ReplaceAll(updateRocketmqInstancePath, "{instance_id}", fmt.Sprintf("%v", instanceId))

	if d.HasChanges(updateRocketmqInstanceHasChanges...) {
		updateRocketmqInstanceOpt := golangsdk.RequestOpts{
			KeepResponseBody: true,
			OkCodes: []int{
				204,
			},
		}
		updateRocketmqInstanceOpt.JSONBody = utils.RemoveNil(buildUpdateRocketmqInstanceBodyParams(d))

		retryFunc := func() (interface{}, bool, error) {
			_, err := updateRocketmqInstanceClient.Request("PUT", updateRocketmqInstancePath, &updateRocketmqInstanceOpt)
			retry, err := handleMultiOperationsError(err)
			return nil, retry, err
		}
		_, err := common.RetryContextWithWaitForState(&common.RetryContextWithWaitForStateParam{
			Ctx:          ctx,
			RetryFunc:    retryFunc,
			WaitFunc:     rocketmqInstanceStateRefreshFunc(updateRocketmqInstanceClient, d.Id()),
			WaitTarget:   []string{"RUNNING"},
			Timeout:      d.Timeout(schema.TimeoutUpdate),
			DelayTimeout: 10 * time.Second,
			PollInterval: 10 * time.Second,
		})

		if err != nil {
			return diag.Errorf("error updating DMS RocketMQ instance: %s", err)
		}

		if d.HasChange("enable_publicip") {
			if err = waitUpdatePublicipSuccess(ctx, d, updateRocketmqInstanceClient); err != nil {
				return diag.FromErr(err)
			}
		}
	}

	if d.HasChange("cross_vpc_accesses") {
		if err = updateCrossVpcAccess(ctx, updateRocketmqInstanceClient, d); err != nil {
			return diag.Errorf("error updating DMS RocketMQ Cross-VPC access information: %s", err)
		}
	}

	if d.HasChange("auto_renew") {
		bssClient, err := cfg.BssV2Client(cfg.GetRegion(d))
		if err != nil {
			return diag.Errorf("error creating BSS V2 client: %s", err)
		}
		if err = common.UpdateAutoRenew(bssClient, d.Get("auto_renew").(string), instanceId); err != nil {
			return diag.Errorf("error updating the auto-renew of the RocketMQ instance (%s): %s", instanceId, err)
		}
	}
	// update tags
	if d.HasChange("tags") {
		tagErr := utils.UpdateResourceTags(updateRocketmqInstanceClient, d, "rocketmq", instanceId)
		if tagErr != nil {
			return diag.Errorf("error updating tags of RocketMQ:%s, err:%s", instanceId, tagErr)
		}
	}

	if d.HasChanges("flavor_id", "broker_num", "storage_space") {
		err := resizeRocketmqInstance(ctx, updateRocketmqInstanceClient, d)
		if err != nil {
			return diag.FromErr(err)
		}
	}

	if d.HasChange("enterprise_project_id") {
		migrateOpts := enterpriseprojects.MigrateResourceOpts{
			ResourceId:   instanceId,
			ResourceType: "rocketmq",
			RegionId:     region,
			ProjectId:    cfg.GetProjectID(region),
		}
		if err := common.MigrateEnterpriseProject(ctx, cfg, d, migrateOpts); err != nil {
			return diag.FromErr(err)
		}
	}

	return resourceDmsRocketMQInstanceRead(ctx, d, meta)
}

func buildUpdateRocketmqInstanceBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"description":       utils.ValueIgnoreEmpty(d.Get("description")),
		"security_group_id": utils.ValueIgnoreEmpty(d.Get("security_group_id")),
	}

	if d.HasChange("enable_acl") {
		bodyParams["enable_acl"] = utils.ValueIgnoreEmpty(d.Get("enable_acl"))
	} else if d.HasChange("retention_policy") {
		bodyParams["enable_acl"] = utils.ValueIgnoreEmpty(d.Get("retention_policy"))
	}

	if d.HasChange("name") {
		bodyParams["name"] = utils.ValueIgnoreEmpty(d.Get("name"))
	}

	if d.HasChange("enable_publicip") {
		bodyParams["enable_publicip"] = utils.ValueIgnoreEmpty(d.Get("enable_publicip"))
	}

	if d.HasChange("publicip_id") {
		bodyParams["publicip_id"] = utils.ValueIgnoreEmpty(d.Get("publicip_id"))
	}
	return bodyParams
}

func resourceDmsRocketMQInstanceRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var mErr *multierror.Error

	// getRocketmqInstance: Query DMS rocketmq instance
	var (
		getRocketmqInstanceHttpUrl = "v2/{project_id}/instances/{instance_id}"
		getRocketmqInstanceProduct = "dmsv2"
	)
	getRocketmqInstanceClient, err := cfg.NewServiceClient(getRocketmqInstanceProduct, region)
	if err != nil {
		return diag.Errorf("error creating DmsRocketMQInstance Client: %s", err)
	}

	getRocketmqInstancePath := getRocketmqInstanceClient.Endpoint + getRocketmqInstanceHttpUrl
	getRocketmqInstancePath = strings.ReplaceAll(getRocketmqInstancePath, "{project_id}",
		getRocketmqInstanceClient.ProjectID)
	getRocketmqInstancePath = strings.ReplaceAll(getRocketmqInstancePath, "{instance_id}", fmt.Sprintf("%v", d.Id()))

	getRocketmqInstanceOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		OkCodes: []int{
			200,
		},
	}
	getRocketmqInstanceResp, err := getRocketmqInstanceClient.Request("GET", getRocketmqInstancePath, &getRocketmqInstanceOpt)

	if err != nil {
		return common.CheckDeletedDiag(d, err, "error retrieving DmsRocketMQInstance")
	}

	getRocketmqInstanceRespBody, err := utils.FlattenResponse(getRocketmqInstanceResp)
	if err != nil {
		return diag.FromErr(err)
	}

	// convert the ids of the availability zone into codes
	var availableZoneCodes []string
	availableZoneIDs := utils.PathSearch("available_zones", getRocketmqInstanceRespBody, nil)
	if availableZoneIDs != nil {
		azIDs := make([]string, 0)
		for _, v := range availableZoneIDs.([]interface{}) {
			azIDs = append(azIDs, v.(string))
		}
		availableZoneCodes, err = getAvailableZoneCodeByID(cfg, region, azIDs)
		if err != nil {
			return diag.FromErr(err)
		}
	}
	crossVpcInfo := utils.PathSearch("cross_vpc_info", getRocketmqInstanceRespBody, nil)
	var crossVpcAccess []map[string]interface{}
	if crossVpcInfo != nil {
		crossVpcAccess, err = flattenCrossVpcInfo(crossVpcInfo.(string))
		if err != nil {
			return diag.FromErr(err)
		}
	}

	var chargingMode = "postPaid"
	if utils.PathSearch("charging_mode", getRocketmqInstanceRespBody, 1).(float64) == 0 {
		chargingMode = "prePaid"
	}
	epsID := "all_granted_eps"
	ipIdList, addressList, err := getPublicipInfoByAddresses(meta, region, epsID, getRocketmqInstanceRespBody)
	if err != nil {
		return diag.Errorf("error retrieving public access: %s", err)
	}

	mErr = multierror.Append(
		mErr,
		d.Set("region", region),
		d.Set("name", utils.PathSearch("name", getRocketmqInstanceRespBody, nil)),
		d.Set("status", utils.PathSearch("status", getRocketmqInstanceRespBody, nil)),
		d.Set("description", utils.PathSearch("description", getRocketmqInstanceRespBody, nil)),
		d.Set("type", utils.PathSearch("type", getRocketmqInstanceRespBody, nil)),
		d.Set("specification", utils.PathSearch("specification", getRocketmqInstanceRespBody, nil)),
		d.Set("engine_version", utils.PathSearch("engine_version", getRocketmqInstanceRespBody, nil)),
		d.Set("vpc_id", utils.PathSearch("vpc_id", getRocketmqInstanceRespBody, nil)),
		d.Set("flavor_id", utils.PathSearch("product_id", getRocketmqInstanceRespBody, nil)),
		d.Set("security_group_id", utils.PathSearch("security_group_id", getRocketmqInstanceRespBody, nil)),
		d.Set("subnet_id", utils.PathSearch("subnet_id", getRocketmqInstanceRespBody, nil)),
		d.Set("availability_zones", availableZoneCodes),
		d.Set("maintain_begin", utils.PathSearch("maintain_begin", getRocketmqInstanceRespBody, nil)),
		d.Set("maintain_end", utils.PathSearch("maintain_end", getRocketmqInstanceRespBody, nil)),
		d.Set("storage_space", utils.PathSearch("total_storage_space", getRocketmqInstanceRespBody, nil)),
		d.Set("used_storage_space", utils.PathSearch("used_storage_space", getRocketmqInstanceRespBody, nil)),
		d.Set("enable_publicip", utils.PathSearch("enable_publicip", getRocketmqInstanceRespBody, nil)),
		d.Set("publicip_id", strings.Join(ipIdList, ",")),
		d.Set("publicip_address", strings.Join(addressList, ",")),
		d.Set("ssl_enable", utils.PathSearch("ssl_enable", getRocketmqInstanceRespBody, nil)),
		d.Set("storage_spec_code", utils.PathSearch("storage_spec_code", getRocketmqInstanceRespBody, nil)),
		d.Set("ipv6_enable", utils.PathSearch("ipv6_enable", getRocketmqInstanceRespBody, nil)),
		d.Set("node_num", utils.PathSearch("node_num", getRocketmqInstanceRespBody, nil)),
		d.Set("new_spec_billing_enable", utils.PathSearch("new_spec_billing_enable", getRocketmqInstanceRespBody, nil)),
		d.Set("enable_acl", utils.PathSearch("enable_acl", getRocketmqInstanceRespBody, nil)),
		d.Set("enterprise_project_id", utils.PathSearch("enterprise_project_id", getRocketmqInstanceRespBody, nil)),
		d.Set("broker_num", utils.PathSearch("broker_num", getRocketmqInstanceRespBody, nil)),
		d.Set("namesrv_address", utils.PathSearch("namesrv_address", getRocketmqInstanceRespBody, nil)),
		d.Set("broker_address", utils.PathSearch("broker_address", getRocketmqInstanceRespBody, nil)),
		d.Set("public_namesrv_address", utils.PathSearch("public_namesrv_address", getRocketmqInstanceRespBody, nil)),
		d.Set("public_broker_address", utils.PathSearch("public_broker_address", getRocketmqInstanceRespBody, nil)),
		d.Set("resource_spec_code", utils.PathSearch("resource_spec_code", getRocketmqInstanceRespBody, nil)),
		d.Set("cross_vpc_accesses", crossVpcAccess),
		d.Set("charging_mode", chargingMode),
	)

	// fetch tags
	if resourceTags, err := tags.Get(getRocketmqInstanceClient, "rocketmq", d.Id()).Extract(); err == nil {
		tagMap := utils.TagsToMap(resourceTags.Tags)
		mErr = multierror.Append(mErr, d.Set("tags", tagMap))
	} else {
		fmt.Printf("[WARN] fetching tags of RocketMQ failed: %s", err)
	}

	return diag.FromErr(mErr.ErrorOrNil())
}

func resourceDmsRocketMQInstanceDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	// deleteRocketmqInstance: Delete DMS rocketmq instance
	var (
		deleteRocketmqInstanceHttpUrl = "v2/{project_id}/instances/{instance_id}"
		deleteRocketmqInstanceProduct = "dmsv2"
	)
	deleteRocketmqInstanceClient, err := cfg.NewServiceClient(deleteRocketmqInstanceProduct, region)
	if err != nil {
		return diag.Errorf("error creating DmsRocketMQInstance Client: %s", err)
	}

	if d.Get("charging_mode") == "prePaid" {
		if err := common.UnsubscribePrePaidResource(d, cfg, []string{d.Id()}); err != nil {
			return diag.Errorf("error unsubscribe RocketMQ instance: %s", err)
		}
	} else {
		deleteRocketmqInstancePath := deleteRocketmqInstanceClient.Endpoint + deleteRocketmqInstanceHttpUrl
		deleteRocketmqInstancePath = strings.ReplaceAll(deleteRocketmqInstancePath, "{project_id}",
			deleteRocketmqInstanceClient.ProjectID)
		deleteRocketmqInstancePath = strings.ReplaceAll(deleteRocketmqInstancePath, "{instance_id}",
			fmt.Sprintf("%v", d.Id()))

		deleteRocketmqInstanceOpt := golangsdk.RequestOpts{
			KeepResponseBody: true,
			OkCodes: []int{
				204,
			},
		}
		_, err = deleteRocketmqInstanceClient.Request("DELETE", deleteRocketmqInstancePath, &deleteRocketmqInstanceOpt)
		if err != nil {
			return diag.Errorf("error deleting DmsRocketMQInstance: %s", err)
		}
	}

	stateConf := &resource.StateChangeConf{
		Pending:      []string{"DELETING", "RUNNING", "ERROR"},
		Target:       []string{"DELETED"},
		Refresh:      rocketmqInstanceStateRefreshFunc(deleteRocketmqInstanceClient, d.Id()),
		Timeout:      d.Timeout(schema.TimeoutDelete),
		Delay:        90 * time.Second,
		PollInterval: 15 * time.Second,
	}

	_, err = stateConf.WaitForStateContext(ctx)
	if err != nil {
		return diag.Errorf("error waiting for instance (%s) to delete: %s", d.Id(), err)
	}

	d.SetId("")

	return nil
}

func rocketmqInstanceStateRefreshFunc(client *golangsdk.ServiceClient, instanceID string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		getRocketmqInstancePath := client.Endpoint + "v2/{project_id}/instances/{instance_id}"
		getRocketmqInstancePath = strings.ReplaceAll(getRocketmqInstancePath, "{project_id}", client.ProjectID)
		getRocketmqInstancePath = strings.ReplaceAll(getRocketmqInstancePath, "{instance_id}", fmt.Sprintf("%v", instanceID))
		getRocketmqInstanceOpt := golangsdk.RequestOpts{
			KeepResponseBody: true,
			OkCodes: []int{
				200,
			},
		}
		v, err := client.Request("GET", getRocketmqInstancePath, &getRocketmqInstanceOpt)
		if err != nil {
			if _, ok := err.(golangsdk.ErrDefault404); ok {
				return v, "DELETED", nil
			}
			return nil, "", err
		}
		respBody, err := utils.FlattenResponse(v)
		if err != nil {
			return nil, "", err
		}
		status := utils.PathSearch("status", respBody, "").(string)
		return respBody, status, nil
	}
}

func publicipTaskRefreshFunc(client *golangsdk.ServiceClient, instanceID, taskName string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		// getPublicipTask: get publicip task
		getPublicipTaskHttpUrl := "v2/{project_id}/instances/{instance_id}/tasks"
		getPublicipTaskPath := client.Endpoint + getPublicipTaskHttpUrl
		getPublicipTaskPath = strings.ReplaceAll(getPublicipTaskPath, "{project_id}",
			client.ProjectID)
		getPublicipTaskPath = strings.ReplaceAll(getPublicipTaskPath, "{instance_id}", instanceID)

		reqOpt := golangsdk.RequestOpts{
			KeepResponseBody: true,
		}
		getPublicipTaskResp, err := client.Request("GET", getPublicipTaskPath, &reqOpt)

		if err != nil {
			return nil, "QUERY ERROR", err
		}

		getPublicipTaskRespBody, err := utils.FlattenResponse(getPublicipTaskResp)
		if err != nil {
			return nil, "PARSE ERROR", err
		}

		task := utils.PathSearch(fmt.Sprintf("tasks|[?name=='%s']|[0]", taskName), getPublicipTaskRespBody, nil)
		if task == nil {
			return nil, "NOT FOUND", nil
		}
		status := utils.PathSearch("status", task, "").(string)
		return task, status, nil
	}
}

func waitUpdatePublicipSuccess(ctx context.Context, d *schema.ResourceData, client *golangsdk.ServiceClient) error {
	// The RocketMQ enabling or disabling publicip is done if the status of its task is SUCCESS.
	actionName := "unbindInstancePublicIp"
	if d.Get("enable_publicip").(bool) {
		actionName = "bindInstancePublicIp"
	}
	stateConf := &resource.StateChangeConf{
		Pending:      []string{"CREATED"},
		Target:       []string{"SUCCESS"},
		Refresh:      publicipTaskRefreshFunc(client, d.Id(), actionName),
		Timeout:      d.Timeout(schema.TimeoutUpdate),
		Delay:        1 * time.Second,
		PollInterval: 2 * time.Second,
	}

	_, err := stateConf.WaitForStateContext(ctx)

	if err != nil {
		return fmt.Errorf("error waiting for enabling or diabling publicip to be done: %s", err)
	}

	return nil
}

func resizeRocketmqInstance(ctx context.Context, client *golangsdk.ServiceClient, d *schema.ResourceData) error {
	if d.HasChange("flavor_id") {
		resizeBodyParams := map[string]interface{}{
			"oper_type":      "vertical",
			"new_product_id": d.Get("flavor_id"),
		}

		if err := doRocketmqInstanceResize(ctx, client, d, resizeBodyParams); err != nil {
			return err
		}
	}

	if d.HasChange("broker_num") {
		resizeBodyParams := map[string]interface{}{
			"oper_type":      "horizontal",
			"new_broker_num": d.Get("broker_num"),
		}

		// the increased publicip IDs are used for brokers which will be extended if public access is already enabled.
		if d.Get("enable_publicip").(bool) {
			// get increased publicip IDs.
			oldPublicipIds, newPublicipIds := d.GetChange("publicip_id")
			publicipIds := getExtendIpids(oldPublicipIds.(string), newPublicipIds.(string))

			resizeBodyParams["publicip_id"] = publicipIds
		}

		if err := doRocketmqInstanceResize(ctx, client, d, resizeBodyParams); err != nil {
			return err
		}
	}

	if d.HasChange("storage_space") {
		resizeBodyParams := map[string]interface{}{
			"oper_type":         "storage",
			"new_storage_space": d.Get("storage_space"),
		}

		if err := doRocketmqInstanceResize(ctx, client, d, resizeBodyParams); err != nil {
			return err
		}
	}

	return nil
}

func doRocketmqInstanceResize(ctx context.Context, client *golangsdk.ServiceClient, d *schema.ResourceData, bodyParams map[string]interface{}) error {
	instanceID := d.Id()
	resizeRocketmqInstanceHttpUrl := "v2/rocketmq/{project_id}/instances/{instance_id}/extend"
	resizeRocketmqInstancePath := client.Endpoint + resizeRocketmqInstanceHttpUrl
	resizeRocketmqInstancePath = strings.ReplaceAll(resizeRocketmqInstancePath, "{project_id}", client.ProjectID)
	resizeRocketmqInstancePath = strings.ReplaceAll(resizeRocketmqInstancePath, "{instance_id}", instanceID)

	resizeRocketmqInstanceOpt := golangsdk.RequestOpts{KeepResponseBody: true}
	resizeRocketmqInstanceOpt.JSONBody = utils.RemoveNil(bodyParams)

	retryFunc := func() (interface{}, bool, error) {
		_, err := client.Request("POST", resizeRocketmqInstancePath, &resizeRocketmqInstanceOpt)
		retry, err := handleMultiOperationsError(err)
		return nil, retry, err
	}
	_, err := common.RetryContextWithWaitForState(&common.RetryContextWithWaitForStateParam{
		Ctx:          ctx,
		RetryFunc:    retryFunc,
		WaitFunc:     rocketmqInstanceStateRefreshFunc(client, instanceID),
		WaitTarget:   []string{"RUNNING"},
		Timeout:      d.Timeout(schema.TimeoutUpdate),
		DelayTimeout: 10 * time.Second,
		PollInterval: 10 * time.Second,
	})
	if err != nil {
		return fmt.Errorf("error resizing RocketMQ instance: bodyParams: %#v, err: %s", bodyParams, err)
	}

	stateConf := &resource.StateChangeConf{
		Pending:      []string{"EXTENDING"},
		Target:       []string{"RUNNING"},
		Refresh:      rocketmqInstanceStateRefreshFunc(client, instanceID),
		Timeout:      d.Timeout(schema.TimeoutUpdate),
		Delay:        60 * time.Second,
		PollInterval: 10 * time.Second,
	}
	if _, err := stateConf.WaitForStateContext(ctx); err != nil {
		return fmt.Errorf("error waiting for instance (%s) to resize: %v", instanceID, err)
	}
	return nil
}

// get public access information filtered by addresses from public_broker_address, e.g. "121.37.221.67:10105,139.159.159.46:10106"
func getPublicipInfoByAddresses(meta interface{}, region, epsID string, resp interface{}) ([]string, []string, error) {
	publicBrokerAddress := utils.PathSearch("public_broker_address", resp, "").(string)
	publicNamesrvAddress := utils.PathSearch("public_namesrv_address", resp, "").(string)
	if publicBrokerAddress == "" || publicNamesrvAddress == "" {
		return nil, nil, nil
	}

	cfg := meta.(*config.Config)
	eipClient, err := cfg.NetworkingV1Client(region)
	if err != nil {
		return nil, nil, fmt.Errorf("error creating networking client: %s", err)
	}

	allAddressPortList := strings.Split(publicBrokerAddress+","+publicNamesrvAddress, ",")
	addressList := make([]string, len(allAddressPortList))
	for i, addressPort := range allAddressPortList {
		addressList[i] = strings.Split(addressPort, ":")[0]
	}
	publicips, err := common.GetEipsbyAddresses(eipClient, addressList, epsID)
	if err != nil {
		return nil, nil, err
	}
	ipIdList := make([]string, len(publicips))
	ipAddressList := make([]string, len(publicips))
	for i, ip := range publicips {
		ipIdList[i] = ip.ID
		ipAddressList[i] = ip.PublicAddress
	}
	return ipIdList, ipAddressList, nil
}

func getExtendIpids(old, new string) string {
	oldIpidList := strings.Split(old, ",")
	newIpidList := strings.Split(new, ",")
	extendIpidList := make([]string, 0, len(newIpidList))
	for _, ipid := range newIpidList {
		if !utils.StrSliceContains(oldIpidList, ipid) {
			extendIpidList = append(extendIpidList, ipid)
		}
	}
	extendIpids := strings.Join(extendIpidList, ",")
	return extendIpids
}
