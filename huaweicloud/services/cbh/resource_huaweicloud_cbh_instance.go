package cbh

import (
	"context"
	"errors"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/chnsz/golangsdk"
	"github.com/chnsz/golangsdk/openstack/networking/v1/eips"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// Power action constants.
const (
	Start      = "start"
	Stop       = "stop"
	SoftReboot = "soft-reboot"
	HardReboot = "hard-reboot"
)

// @API CBH POST /v2/{project_id}/cbs/instance/{server_id}/eip/bind
// @API CBH POST /v2/{project_id}/cbs/instance/{server_id}/eip/unbind
// @API CBH POST /v2/{project_id}/cbs/instance
// @API CBH GET /v2/{project_id}/cbs/instance/list
// @API CBH PUT /v2/{project_id}/cbs/instance/password
// @API CBH PUT /v2/{project_id}/cbs/instance/{server_id}/security-groups
// @API CBH PUT /v2/{project_id}/cbs/instance
// @API CBH POST /v2/{project_id}/cbs/instance/{resource_id}/tags/action
// @API CBH GET /v2/{project_id}/cbs/instance/{resource_id}/tags
// @API CBH POST /v2/{project_id}/cbs/instance/start
// @API CBH POST /v2/{project_id}/cbs/instance/stop
// @API CBH POST /v2/{project_id}/cbs/instance/reboot
// @API CBH PUT /v2/{project_id}/cbs/instance/vpc
// @API BSS GET /v2/orders/customer-orders/details/{order_id}
// @API BSS POST /v2/orders/suscriptions/resources/query
// @API BSS POST /v2/orders/subscriptions/resources/unsubscribe
// @API BSS POST /v2/orders/subscriptions/resources/autorenew/{instance_id}
// @API BSS DELETE /v2/orders/subscriptions/resources/autorenew/{instance_id}
func ResourceCBHInstance() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceCBHInstanceCreate,
		UpdateContext: resourceCBHInstanceUpdate,
		ReadContext:   resourceCBHInstanceRead,
		DeleteContext: resourceCBHInstanceDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(60 * time.Minute),
			Update: schema.DefaultTimeout(60 * time.Minute),
			Delete: schema.DefaultTimeout(30 * time.Minute),
		},

		CustomizeDiff: config.MergeDefaultTags(),

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
				ForceNew:    true,
				Description: `Specifies the name of the CBH instance.`,
			},
			"flavor_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `Specifies the product ID of the CBH server.`,
			},
			"vpc_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `Specifies the ID of a VPC.`,
			},
			"subnet_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `Specifies the ID of a subnet.`,
			},
			"security_group_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `Specifies the IDs of the security group.`,
			},
			"availability_zone": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: `Specifies the availability zone name.`,
			},
			"password": {
				Type:        schema.TypeString,
				Required:    true,
				Sensitive:   true,
				Description: `Specifies the password for logging in to the management console.`,
			},
			"charging_mode": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
				ValidateFunc: validation.StringInSlice([]string{
					"prePaid",
					"postPaid",
				}, false),
				Description: `Specifies the charging mode of the CBH instance.`,
			},
			"period_unit": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				ValidateFunc: validation.StringInSlice([]string{
					"month", "year",
				}, false),
				Description: `Specifies the charging period unit of the instance.`,
			},
			"period": {
				Type:        schema.TypeInt,
				Optional:    true,
				ForceNew:    true,
				Description: `Specifies the charging period of the CBH instance.`,
			},
			"auto_renew": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specifies whether auto renew is enabled.`,
			},
			"subnet_address": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `Specifies the IP address of the subnet.`,
			},
			"public_ip_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specifies the ID of the elastic IP.`,
			},
			"ipv6_enable": {
				Type:        schema.TypeBool,
				Optional:    true,
				ForceNew:    true,
				Description: `Specifies whether the IPv6 network is enabled.`,
			},
			"attach_disk_size": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: `Specifies the size of the additional data disk for the CBH instance.`,
			},
			"power_action": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specifies the power action after the CBH instance is created.`,
			},
			"tags": common.TagsSchema(),
			"enterprise_project_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Specifies the enterprise project ID to which the CBH instance belongs.",
			},
			"public_ip": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `schema: Computed; The elastic IP address.`,
			},
			"private_ip": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Indicates the private IP of the instance.`,
			},
			"status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Indicates the status of the instance.`,
			},
			"version": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Indicates the current version of the instance image.`,
			},
			"data_disk_size": {
				Type:        schema.TypeFloat,
				Computed:    true,
				Description: `Indicates the data disk size of the instance.`,
			},
		},
	}
}

func resourceCBHInstanceCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg                      = meta.(*config.Config)
		region                   = cfg.GetRegion(d)
		createCbhInstanceProduct = "cbh"
		serverId                 string
	)

	client, err := cfg.NewServiceClient(createCbhInstanceProduct, region)
	if err != nil {
		return diag.Errorf("error creating CBH client: %s", err)
	}

	bssClient, err := cfg.BssV2Client(cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating BSS v2 client: %s", err)
	}

	if d.Get("charging_mode").(string) == "prePaid" {
		orderId, err := createCBHInstance(client, d, cfg)
		if err != nil {
			return diag.FromErr(err)
		}

		if err := common.WaitOrderComplete(ctx, bssClient, orderId, d.Timeout(schema.TimeoutCreate)); err != nil {
			return diag.FromErr(err)
		}

		resourceId, err := common.WaitOrderResourceComplete(ctx, bssClient, orderId, d.Timeout(schema.TimeoutCreate))
		if err != nil {
			return diag.Errorf("error waiting for CBH instance order %s complete: %s", orderId, err)
		}

		instances, err := getCBHInstanceList(client)
		if err != nil {
			return diag.FromErr(err)
		}
		expression := fmt.Sprintf("[?resource_info.resource_id == '%s']|[0].server_id", resourceId)
		serverId = utils.PathSearch(expression, instances, "").(string)
		if serverId == "" {
			return diag.Errorf("unable to find the CBH instance ID from the API response")
		}

		d.SetId(serverId)
		if err := waitingForCBHInstanceActive(ctx, client, d, d.Timeout(schema.TimeoutCreate)); err != nil {
			return diag.Errorf("error waiting for CBH instance (%s) creation to active: %s", d.Id(), err)
		}
	} else {
		instanceId, err := createCBHInstance(client, d, cfg)
		if err != nil {
			return diag.FromErr(err)
		}

		if err := waitingForCBHInstanceActiveByInstanceID(ctx, client, instanceId, d, d.Timeout(schema.TimeoutCreate)); err != nil {
			return diag.Errorf("error waiting for CBH instance (%s) creation to active: %s", d.Id(), err)
		}
	}

	// After successfully creating an instance using the first security group ID, check if it is necessary to update the
	// security group to the target value.
	securityGroupIDs := d.Get("security_group_id").(string)
	sgIDs := strings.Split(securityGroupIDs, ",")
	if len(sgIDs) > 1 {
		if err := updateSecurityGroup(client, d.Id(), sgIDs); err != nil {
			return diag.Errorf("error updating the security group after successful creation of CBH instance (%s): %s", d.Id(), err)
		}
	}

	// Create an instance in the shutdown state.
	if action, ok := d.GetOk("power_action"); ok {
		if action.(string) == Stop {
			if err = doPowerAction(ctx, client, d, d.Timeout(schema.TimeoutCreate)); err != nil {
				return diag.FromErr(err)
			}
		} else {
			log.Printf("[WARN] the power action (%s) is invalid after CBH instance created", action)
		}
	}

	return resourceCBHInstanceRead(ctx, d, meta)
}

func createCBHInstance(client *golangsdk.ServiceClient, d *schema.ResourceData, cfg *config.Config) (string, error) {
	publicIp, err := buildCreateNetworkPublicIpBodyParam(d, cfg)
	if err != nil {
		return "", fmt.Errorf("error creating CBH instance: error building network public IP body: %s", err)
	}

	createInstancePath := client.Endpoint + "v2/{project_id}/cbs/instance"
	createInstancePath = strings.ReplaceAll(createInstancePath, "{project_id}", client.ProjectID)
	createInstanceOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		JSONBody:         utils.RemoveNil(buildCreateCBHInstanceBodyParam(d, cfg.GetRegion(d), cfg.GetEnterpriseProjectID(d), publicIp)),
	}

	createInstanceResp, err := client.Request("POST", createInstancePath, &createInstanceOpt)
	if err != nil {
		return "", fmt.Errorf("error creating CBH instance: %s", err)
	}
	createInstanceRespBody, err := utils.FlattenResponse(createInstanceResp)
	if err != nil {
		return "", err
	}

	var targetId string
	if d.Get("charging_mode").(string) == "prePaid" {
		targetId = utils.PathSearch("order_id", createInstanceRespBody, "").(string)
	} else {
		targetId = utils.PathSearch("instance_id", createInstanceRespBody, "").(string)
	}

	if targetId == "" {
		return "", errors.New("unable to find the order/instance ID of the CBH instance from the API response")
	}
	return targetId, nil
}

func buildCreateCBHInstanceBodyParam(d *schema.ResourceData, region string, epsId string, publicIp interface{}) map[string]interface{} {
	bodyParam := map[string]interface{}{
		"specification":         d.Get("flavor_id"),
		"instance_name":         d.Get("name"),
		"password":              d.Get("password"),
		"region":                region,
		"availability_zone":     d.Get("availability_zone"),
		"charging_mode":         buildCreateChargingModeParam(d),
		"network":               buildCreateNetworkBodyParam(d, publicIp),
		"attach_disk_size":      utils.ValueIgnoreEmpty(d.Get("attach_disk_size")),
		"enterprise_project_id": utils.ValueIgnoreEmpty(epsId),
	}
	if d.Get("charging_mode").(string) == "prePaid" {
		bodyParam["period_type"] = buildCreatePeriodTypeParam(d)
		bodyParam["period_num"] = utils.ValueIgnoreEmpty(d.Get("period"))
		bodyParam["is_auto_renew"] = buildCreateIsAutoRenewParam(d)
		bodyParam["is_auto_pay"] = 1
	}
	// The default value of the field `ipv6_enable` is false
	if d.Get("ipv6_enable").(bool) {
		bodyParam["ipv6_enable"] = true
	}

	if tagsRaw, ok := d.GetOk("tags"); ok {
		bodyParam["tags"] = utils.ExpandResourceTags(tagsRaw.(map[string]interface{}))
	}

	return bodyParam
}

// Currently, the CBH instance only supports prePaid charging mode
func buildCreateChargingModeParam(d *schema.ResourceData) interface{} {
	if d.Get("charging_mode").(string) == "prePaid" {
		return 0
	}
	if d.Get("charging_mode").(string) == "postPaid" {
		return 1
	}
	return nil
}

// Currently, the CBH instance only supports `year` and `month` period unit in prePaid charging mode.
func buildCreatePeriodTypeParam(d *schema.ResourceData) interface{} {
	if d.Get("period_unit").(string) == "year" {
		return 3
	}

	if d.Get("period_unit").(string) == "month" {
		return 2
	}
	return nil
}

// `1` indicates automatic renewal. `0` indicates non-automatic renewal.
func buildCreateIsAutoRenewParam(d *schema.ResourceData) interface{} {
	if d.Get("auto_renew").(string) == "true" {
		return 1
	}

	if d.Get("auto_renew").(string) == "false" {
		return 0
	}
	return nil
}

func buildCreateNetworkBodyParam(d *schema.ResourceData, publicIp interface{}) interface{} {
	return map[string]interface{}{
		"vpc_id":          d.Get("vpc_id"),
		"subnet_id":       d.Get("subnet_id"),
		"public_ip":       publicIp,
		"security_groups": buildCreateNetworkSecurityGroupsBodyParam(d),
		"private_ip":      buildCreateNetworkPrivateIpBodyParam(d),
	}
}

func buildCreateNetworkPublicIpBodyParam(d *schema.ResourceData, cfg *config.Config) (interface{}, error) {
	publicIpId := d.Get("public_ip_id").(string)
	if publicIpId == "" {
		return nil, nil
	}

	publicIp, err := getPublicAddressById(d, cfg, publicIpId)
	if err != nil {
		return nil, err
	}

	params := map[string]interface{}{
		"id":         publicIpId,
		"public_eip": publicIp,
	}
	return params, nil
}

func buildCreateNetworkSecurityGroupsBodyParam(d *schema.ResourceData) interface{} {
	securityGroupIDs := d.Get("security_group_id").(string)
	sgIDList := strings.Split(securityGroupIDs, ",")
	if len(sgIDList) == 0 {
		return nil
	}
	// When creating an instance, if multiple security group IDs are specified,
	// prioritize using the first one for creation.
	return []map[string]interface{}{
		{
			"id": sgIDList[0],
		},
	}
}

func buildCreateNetworkPrivateIpBodyParam(d *schema.ResourceData) interface{} {
	if v, ok := d.GetOk("subnet_address"); ok {
		return map[string]interface{}{
			"ip": v,
		}
	}
	return nil
}

func getCBHInstanceList(client *golangsdk.ServiceClient) ([]interface{}, error) {
	getCbhInstancesPath := client.Endpoint + "v2/{project_id}/cbs/instance/list"
	getCbhInstancesPath = strings.ReplaceAll(getCbhInstancesPath, "{project_id}", client.ProjectID)
	getCbhInstancesOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	getCbhInstancesResp, err := client.Request("GET", getCbhInstancesPath, &getCbhInstancesOpt)
	if err != nil {
		return nil, fmt.Errorf("error retrieving CBH instance list: %s", err)
	}

	getCbhInstancesRespBody, err := utils.FlattenResponse(getCbhInstancesResp)
	if err != nil {
		return nil, err
	}
	instances := utils.PathSearch("instance", getCbhInstancesRespBody, make([]interface{}, 0)).([]interface{})
	return instances, nil
}

func convertPowerActionValue(action string) string {
	switch action {
	case SoftReboot:
		return "SOFT"
	case HardReboot:
		return "HARD"
	default:
		return action
	}
}

// doPowerAction is a method for CBH instance power doing startup, shutdown and reboot actions.
func doPowerAction(ctx context.Context, client *golangsdk.ServiceClient, d *schema.ResourceData, timeout time.Duration) error {
	doActionInstancePath := client.Endpoint + "v2/{project_id}/cbs/instance/{action}"
	doActionInstancePath = strings.ReplaceAll(doActionInstancePath, "{project_id}", client.ProjectID)
	var (
		id          = d.Id()
		action      = d.Get("power_action").(string)
		jsonBodyMap = map[string]interface{}{
			"server_id": id,
		}
	)

	if action == SoftReboot || action == HardReboot {
		doActionInstancePath = strings.ReplaceAll(doActionInstancePath, "{action}", "reboot")
		jsonBodyMap["reboot_type"] = convertPowerActionValue(action)
	} else {
		doActionInstancePath = strings.ReplaceAll(doActionInstancePath, "{action}", action)
	}

	doActionInstanceOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		JSONBody:         jsonBodyMap,
	}
	_, err := client.Request("POST", doActionInstancePath, &doActionInstanceOpt)
	if err != nil {
		return fmt.Errorf("doing power action (%s) for CBH instance (%s) failed: %s", action, id, err)
	}

	if err := waitingForCBHInstanceTaskCompleted(ctx, client, d, timeout); err != nil {
		return fmt.Errorf("error waiting for CBH instance (%s) doing action (%s) task to complete: %s", id, action, err)
	}

	if action == Start || action == SoftReboot || action == HardReboot {
		if err := waitingForCBHInstanceActive(ctx, client, d, timeout); err != nil {
			return fmt.Errorf("error waiting for CBH instance (%s) doing action (%s) to active: %s", id, action, err)
		}
	}

	if action == Stop {
		if err := waitingForCBHInstanceShutoff(ctx, client, d, timeout); err != nil {
			return fmt.Errorf("error waiting for CBH instance (%s) doing action (%s) to shutoff: %s", id, action, err)
		}
	}

	return nil
}

func waitingForCBHInstanceShutoff(ctx context.Context, client *golangsdk.ServiceClient, d *schema.ResourceData,
	timeout time.Duration) error {
	expression := fmt.Sprintf("[?server_id == '%s']|[0]", d.Id())
	unexpectedStatus := []string{"DELETING", "DELETED", "ERROR", "FROZEN"}
	stateConf := &resource.StateChangeConf{
		Pending: []string{"PENDING"},
		Target:  []string{"COMPLETED"},
		Refresh: func() (interface{}, string, error) {
			instances, err := getCBHInstanceList(client)
			if err != nil {
				return nil, "ERROR", err
			}

			instance := utils.PathSearch(expression, instances, nil)
			if instance == nil {
				return nil, "ERROR", golangsdk.ErrDefault404{}
			}

			status := utils.PathSearch("status_info.status", instance, "").(string)
			if status == "SHUTOFF" {
				return instance, "COMPLETED", nil
			}

			if status == "" {
				return instance, "ERROR", fmt.Errorf("status is not found in list API response")
			}

			if utils.StrSliceContains(unexpectedStatus, status) {
				return instance, status, nil
			}

			return instance, "PENDING", nil
		},
		Timeout:      timeout,
		Delay:        10 * time.Second,
		PollInterval: 10 * time.Second,
	}

	_, err := stateConf.WaitForStateContext(ctx)
	return err
}

func waitingForCBHInstanceActiveByInstanceID(ctx context.Context, client *golangsdk.ServiceClient, instanceId string, d *schema.ResourceData,
	timeout time.Duration) error {
	expression := fmt.Sprintf("[?instance_id == '%s']|[0]", instanceId)
	unexpectedStatus := []string{"SHUTOFF", "DELETING", "DELETED", "ERROR", "FROZEN"}
	stateConf := &resource.StateChangeConf{
		Pending: []string{"PENDING"},
		Target:  []string{"COMPLETED"},
		Refresh: func() (interface{}, string, error) {
			instances, err := getCBHInstanceList(client)
			if err != nil {
				return nil, "ERROR", err
			}
			instance := utils.PathSearch(expression, instances, nil)
			if instance == nil {
				return nil, "ERROR", golangsdk.ErrDefault404{}
			}

			status := utils.PathSearch("status_info.status", instance, "").(string)
			if status == "ACTIVE" {
				serverId := utils.PathSearch("server_id", instance, "").(string)
				d.SetId(serverId)
				return instance, "COMPLETED", nil
			}

			if status == "" {
				return instance, "ERROR", errors.New("status is not found in list API response")
			}

			if utils.StrSliceContains(unexpectedStatus, status) {
				return instance, status, nil
			}
			return instance, "PENDING", nil
		},
		Timeout:      timeout,
		Delay:        10 * time.Second,
		PollInterval: 10 * time.Second,
	}

	_, err := stateConf.WaitForStateContext(ctx)
	return err
}

func waitingForCBHInstanceActive(ctx context.Context, client *golangsdk.ServiceClient, d *schema.ResourceData,
	timeout time.Duration) error {
	expression := fmt.Sprintf("[?server_id == '%s']|[0]", d.Id())
	unexpectedStatus := []string{"SHUTOFF", "DELETING", "DELETED", "ERROR", "FROZEN"}
	stateConf := &resource.StateChangeConf{
		Pending: []string{"PENDING"},
		Target:  []string{"COMPLETED"},
		Refresh: func() (interface{}, string, error) {
			instances, err := getCBHInstanceList(client)
			if err != nil {
				return nil, "ERROR", err
			}
			instance := utils.PathSearch(expression, instances, nil)
			if instance == nil {
				return nil, "ERROR", golangsdk.ErrDefault404{}
			}

			status := utils.PathSearch("status_info.status", instance, "").(string)
			if status == "ACTIVE" {
				return instance, "COMPLETED", nil
			}

			if status == "" {
				return instance, "ERROR", errors.New("status is not found in list API response")
			}

			if utils.StrSliceContains(unexpectedStatus, status) {
				return instance, status, nil
			}
			return instance, "PENDING", nil
		},
		Timeout:      timeout,
		Delay:        10 * time.Second,
		PollInterval: 10 * time.Second,
	}

	_, err := stateConf.WaitForStateContext(ctx)
	return err
}

func waitingForCBHInstanceTaskCompleted(ctx context.Context, client *golangsdk.ServiceClient, d *schema.ResourceData,
	timeout time.Duration) error {
	expression := fmt.Sprintf("[?server_id == '%s']|[0]", d.Id())
	unexpectedTaskStatus := []string{"delete_wait", "frozen", "unfrozen", "updating", "configuring-ha",
		"data-migrating", "rollback", "traffic-switchover"}
	stateConf := &resource.StateChangeConf{
		Pending: []string{"PENDING"},
		Target:  []string{"COMPLETED"},
		Refresh: func() (interface{}, string, error) {
			instances, err := getCBHInstanceList(client)
			if err != nil {
				return nil, "ERROR", err
			}
			instance := utils.PathSearch(expression, instances, nil)
			if instance == nil {
				return nil, "ERROR", golangsdk.ErrDefault404{}
			}

			taskStatus := utils.PathSearch("status_info.task_status", instance, "").(string)
			if taskStatus == "NO_TASK" {
				return instance, "COMPLETED", nil
			}

			if taskStatus == "" {
				return instance, "ERROR", fmt.Errorf("the cbh instnace task_status is not found in list API response")
			}

			if utils.StrSliceContains(unexpectedTaskStatus, taskStatus) {
				return instance, taskStatus, nil
			}

			return instance, "PENDING", nil
		},
		Timeout:      timeout,
		Delay:        10 * time.Second,
		PollInterval: 10 * time.Second,
	}

	_, err := stateConf.WaitForStateContext(ctx)
	return err
}

func resourceCBHInstanceUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg                      = meta.(*config.Config)
		region                   = cfg.GetRegion(d)
		updateCBHInstanceProduct = "cbh"
		ID                       = d.Id()
	)

	client, err := cfg.NewServiceClient(updateCBHInstanceProduct, region)
	if err != nil {
		return diag.Errorf("error creating CBH client: %s", err)
	}

	// Update instance can only be performed when the instance status is active.
	// Therefore, if `power_action` is start, it must be done in the first step of the update function.
	action := d.Get("power_action").(string)
	if d.HasChanges("power_action") && action == Start {
		if err = doPowerAction(ctx, client, d, d.Timeout(schema.TimeoutUpdate)); err != nil {
			return diag.FromErr(err)
		}
	}

	// Update `vpc_id`, `subnet_id`, and `subnet_address` using the same API.
	if d.HasChanges("vpc_id", "subnet_id", "subnet_address") {
		if err := updateVpc(client, d); err != nil {
			return diag.Errorf("error updating the vpc of the CBH instance (%s): %s", ID, err)
		}

		if err := waitingForCBHInstanceActive(ctx, client, d, d.Timeout(schema.TimeoutUpdate)); err != nil {
			return diag.Errorf("error waiting for CBH instance (%s) update vpc complete: %s", ID, err)
		}
	}

	if d.HasChanges("security_group_id") {
		securityGroupIDs := d.Get("security_group_id").(string)
		sgIDs := strings.Split(securityGroupIDs, ",")
		if err := updateSecurityGroup(client, ID, sgIDs); err != nil {
			return diag.Errorf("error updating the security group of the CBH instance (%s): %s", ID, err)
		}
	}

	if d.HasChanges("public_ip_id") {
		if err := updatePublicIpId(client, d); err != nil {
			return diag.FromErr(err)
		}
	}

	if d.HasChanges("password") {
		if err = updatePassword(client, d); err != nil {
			return diag.FromErr(err)
		}
	}

	// Due to API limitations, update `flavor_id` and `attach_disk_size` must call the API separately.
	if d.HasChanges("flavor_id") {
		orderId, err := updateFlavorId(client, ID, d.Get("flavor_id").(string))
		if err != nil {
			return diag.FromErr(err)
		}

		bssClient, err := cfg.BssV2Client(cfg.GetRegion(d))
		if err != nil {
			return diag.Errorf("error creating BSS v2 client: %s", err)
		}

		if err := common.WaitOrderComplete(ctx, bssClient, orderId, d.Timeout(schema.TimeoutUpdate)); err != nil {
			return diag.FromErr(err)
		}

		if err := waitingForCBHInstanceActive(ctx, client, d, d.Timeout(schema.TimeoutUpdate)); err != nil {
			return diag.Errorf("error waiting for CBH instance (%s) update to active: %s", ID, err)
		}

		// After updating the `flavor_id` and `attach_disk_size`, the instance will automatically restart;
		// At this point, the instance is restarting, but the value of the `status` attribute is already **active**,
		// so we need to continue waiting for the instance's `task_status` attribute to become **NO_TASK**
		// before we can proceed.
		if err := waitingForCBHInstanceTaskCompleted(ctx, client, d, d.Timeout(schema.TimeoutUpdate)); err != nil {
			return diag.Errorf("error waiting for CBH instance (%s) update flavor to complete: %s", ID, err)
		}
	}

	if d.HasChanges("attach_disk_size") {
		orderId, err := updateAttachDiskSize(client, ID, int32(d.Get("attach_disk_size").(int)))
		if err != nil {
			return diag.FromErr(err)
		}

		bssClient, err := cfg.BssV2Client(cfg.GetRegion(d))
		if err != nil {
			return diag.Errorf("error creating BSS v2 client: %s", err)
		}

		if err := common.WaitOrderComplete(ctx, bssClient, orderId, d.Timeout(schema.TimeoutUpdate)); err != nil {
			return diag.FromErr(err)
		}

		if err := waitingForCBHInstanceActive(ctx, client, d, d.Timeout(schema.TimeoutUpdate)); err != nil {
			return diag.Errorf("error waiting for CBH instance (%s) update to active: %s", ID, err)
		}

		if err := waitingForCBHInstanceTaskCompleted(ctx, client, d, d.Timeout(schema.TimeoutUpdate)); err != nil {
			return diag.Errorf("error waiting for CBH instance (%s) update additional disk to complete: %s", ID, err)
		}
	}

	if d.HasChange("auto_renew") {
		bssClient, err := cfg.BssV2Client(region)
		if err != nil {
			return diag.Errorf("error creating BSS V2 client: %s", err)
		}

		resourceId, err := getInstanceResourceIdById(client, ID)
		if err != nil {
			return diag.FromErr(err)
		}

		if resourceId == "" {
			return diag.Errorf("error updating the auto-renew of the CBH instance (%s): "+
				"resource ID is not found in list API response", ID)
		}

		if err = common.UpdateAutoRenew(bssClient, d.Get("auto_renew").(string), resourceId); err != nil {
			return diag.Errorf("error updating the auto-renew of the CBH instance (%s): %s", ID, err)
		}
	}

	if d.HasChange("enterprise_project_id") {
		resourceId, err := getInstanceResourceIdById(client, ID)
		if err != nil {
			return diag.FromErr(err)
		}

		if resourceId == "" {
			return diag.Errorf("error updating the enterprise project ID of the CBH instance (%s): "+
				"resource ID is not found in list API response", ID)
		}

		migrateOpts := config.MigrateResourceOpts{
			ResourceId:   resourceId,
			ResourceType: "cbh",
			RegionId:     region,
			ProjectId:    client.ProjectID,
		}
		if err := cfg.MigrateEnterpriseProject(ctx, d, migrateOpts); err != nil {
			return diag.FromErr(err)
		}
	}

	if d.HasChange("tags") {
		resourceId, err := getInstanceResourceIdById(client, ID)
		if err != nil {
			return diag.FromErr(err)
		}

		if resourceId == "" {
			return diag.Errorf("error updating tags of the CBH instance (%s): "+
				"resource ID is not found in list API response", ID)
		}

		oRaw, nRaw := d.GetChange("tags")
		oMap := oRaw.(map[string]interface{})
		nMap := nRaw.(map[string]interface{})

		if len(oMap) > 0 {
			if err = doActionInstanceTags(resourceId, "delete", client, oMap); err != nil {
				return diag.FromErr(err)
			}
		}

		if len(nMap) > 0 {
			if err := doActionInstanceTags(resourceId, "create", client, nMap); err != nil {
				return diag.FromErr(err)
			}
		}
	}

	// The stop and reboot operations should be performed after the instance updated other parameters.
	// Therefore, they must be performed at the end of the update function.
	if d.HasChange("power_action") && (action == Stop || action == SoftReboot || action == HardReboot) {
		if err = doPowerAction(ctx, client, d, d.Timeout(schema.TimeoutUpdate)); err != nil {
			return diag.FromErr(err)
		}
	}

	return resourceCBHInstanceRead(ctx, d, meta)
}

func doActionInstanceTags(resourceId, action string, client *golangsdk.ServiceClient, tagsMap map[string]interface{}) error {
	doActionTagsHttpUrl := "v2/{project_id}/cbs/instance/{resource_id}/tags/action"
	doActionTagsPath := client.Endpoint + doActionTagsHttpUrl
	doActionTagsPath = strings.ReplaceAll(doActionTagsPath, "{project_id}", client.ProjectID)
	doActionTagsPath = strings.ReplaceAll(doActionTagsPath, "{resource_id}", resourceId)
	doActionTagsOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		OkCodes: []int{
			204,
		},
	}
	doActionTagsOpt.JSONBody = map[string]interface{}{
		"action": action,
		"tags":   utils.ExpandResourceTags(tagsMap),
	}

	_, err := client.Request("POST", doActionTagsPath, &doActionTagsOpt)
	if err != nil {
		return fmt.Errorf("error updating (action: %s) tags of the CBH instance: %s", action, err)
	}

	return nil
}

func getInstanceResourceIdById(client *golangsdk.ServiceClient, instanceId string) (string, error) {
	instances, err := getCBHInstanceList(client)
	if err != nil {
		return "", err
	}
	expression := fmt.Sprintf("[?server_id == '%s']|[0].resource_info.resource_id", instanceId)
	resourceId := utils.PathSearch(expression, instances, "").(string)

	return resourceId, nil
}

func buildUpdateVpcNetWorkBodyParam(d *schema.ResourceData) interface{} {
	return map[string]interface{}{
		"vpc_id":    d.Get("vpc_id"),
		"subnet_id": d.Get("subnet_id"),
		// When updating VPC, security group ID is a required parameter, but it will not take effect in reality.
		"security_groups": []map[string]interface{}{
			{"id": ""},
		},
		"private_ip": buildCreateNetworkPrivateIpBodyParam(d),
	}
}

func updateVpc(client *golangsdk.ServiceClient, d *schema.ResourceData) error {
	updateVpcPath := client.Endpoint + "v2/{project_id}/cbs/instance/vpc"
	updateVpcPath = strings.ReplaceAll(updateVpcPath, "{project_id}", client.ProjectID)
	updateVpcOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		JSONBody: map[string]interface{}{
			"server_id": d.Id(),
			"network":   buildUpdateVpcNetWorkBodyParam(d),
		},
	}

	_, err := client.Request("PUT", updateVpcPath, &updateVpcOpt)

	return err
}

func updateFlavorId(client *golangsdk.ServiceClient, resourceId, flavorId string) (string, error) {
	updateFlavorIdPath := client.Endpoint + "v2/{project_id}/cbs/instance"
	updateFlavorIdPath = strings.ReplaceAll(updateFlavorIdPath, "{project_id}", client.ProjectID)
	updateFlavorIdOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		JSONBody: map[string]interface{}{
			"server_id":              resourceId,
			"new_resource_spec_code": flavorId,
			"is_auto_pay":            1,
		},
	}
	updateInstanceResp, err := client.Request("PUT", updateFlavorIdPath, &updateFlavorIdOpt)
	if err != nil {
		return "", fmt.Errorf("error updating CBH instance flavor: %s", err)
	}

	updateInstanceRespBody, err := utils.FlattenResponse(updateInstanceResp)
	if err != nil {
		return "", err
	}

	orderId := utils.PathSearch("order_id", updateInstanceRespBody, "").(string)
	if orderId == "" {
		return "", fmt.Errorf("unable to find the order ID of the CBH HA instance flavor from the API response")
	}

	return orderId, nil
}

func updateAttachDiskSize(client *golangsdk.ServiceClient, resourceId string, attachDiskSize int32) (string, error) {
	updateAttachDiskSizePath := client.Endpoint + "v2/{project_id}/cbs/instance"
	updateAttachDiskSizePath = strings.ReplaceAll(updateAttachDiskSizePath, "{project_id}", client.ProjectID)
	updateAttachDiskSizeOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		JSONBody: map[string]interface{}{
			"server_id":        resourceId,
			"attach_disk_size": attachDiskSize,
			"is_auto_pay":      1,
		},
	}
	updateInstanceResp, err := client.Request("PUT", updateAttachDiskSizePath, &updateAttachDiskSizeOpt)
	if err != nil {
		return "", fmt.Errorf("error updating CBH instance additional disk: %s", err)
	}

	updateInstanceRespBody, err := utils.FlattenResponse(updateInstanceResp)
	if err != nil {
		return "", err
	}

	orderId := utils.PathSearch("order_id", updateInstanceRespBody, "").(string)
	if orderId == "" {
		return "", fmt.Errorf("unable to find the order ID of the CBH HA instance additional disk from the API response")
	}

	return orderId, nil
}

func updateSecurityGroup(client *golangsdk.ServiceClient, resourceId string, sgIDs []string) error {
	updateSecurityGroupPath := client.Endpoint + "v2/{project_id}/cbs/instance/{server_id}/security-groups"
	updateSecurityGroupPath = strings.ReplaceAll(updateSecurityGroupPath, "{project_id}", client.ProjectID)
	updateSecurityGroupPath = strings.ReplaceAll(updateSecurityGroupPath, "{server_id}", resourceId)
	updateSecurityGroupOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		JSONBody: map[string]interface{}{
			"security_groups": sgIDs,
		},
	}
	_, err := client.Request("PUT", updateSecurityGroupPath, &updateSecurityGroupOpt)

	return err
}

func updatePublicIpId(client *golangsdk.ServiceClient, d *schema.ResourceData) error {
	oPublicIpIdRaw, nPublicIpIdRaw := d.GetChange("public_ip_id")
	oPublicIpId := strings.TrimSpace(oPublicIpIdRaw.(string))
	nPublicIpId := strings.TrimSpace(nPublicIpIdRaw.(string))

	if len(oPublicIpId) > 0 {
		if err := unbindEip(client, d, oPublicIpId); err != nil {
			return err
		}
	}
	if len(nPublicIpId) > 0 {
		if err := bindEip(client, d, nPublicIpId); err != nil {
			// if bind new eip fail, then bind the old eip to CBH instance
			if len(oPublicIpId) > 0 {
				if err := bindEip(client, d, oPublicIpId); err != nil {
					log.Printf("[WARN] error bind old EIP: %s", err)
				}
			}
			return err
		}
	}
	return nil
}

func unbindEip(client *golangsdk.ServiceClient, d *schema.ResourceData, publicIpId string) error {
	unbindPath := client.Endpoint + "v2/{project_id}/cbs/instance/{server_id}/eip/unbind"
	unbindPath = strings.ReplaceAll(unbindPath, "{project_id}", client.ProjectID)
	unbindPath = strings.ReplaceAll(unbindPath, "{server_id}", d.Id())
	unbindEipOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		JSONBody: map[string]interface{}{
			"publicip_id": publicIpId,
		},
	}
	_, err := client.Request("POST", unbindPath, &unbindEipOpt)
	if err != nil {
		return fmt.Errorf("error unbind EIP (%s) from CBH instance (%s): %s", publicIpId, d.Id(), err)
	}
	return nil
}

func bindEip(client *golangsdk.ServiceClient, d *schema.ResourceData, publicIpId string) error {
	bindPath := client.Endpoint + "v2/{project_id}/cbs/instance/{server_id}/eip/bind"
	bindPath = strings.ReplaceAll(bindPath, "{project_id}", client.ProjectID)
	bindPath = strings.ReplaceAll(bindPath, "{server_id}", d.Id())
	unbindEipOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		JSONBody: map[string]interface{}{
			"publicip_id": publicIpId,
		},
	}
	_, err := client.Request("POST", bindPath, &unbindEipOpt)
	if err != nil {
		return fmt.Errorf("error bind EIP (%s) to CBH instance (%s): %s", publicIpId, d.Id(), err)
	}
	return nil
}

func updatePassword(client *golangsdk.ServiceClient, d *schema.ResourceData) error {
	updatePasswordPath := client.Endpoint + "v2/{project_id}/cbs/instance/password"
	updatePasswordPath = strings.ReplaceAll(updatePasswordPath, "{project_id}", client.ProjectID)
	updatePasswordOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		JSONBody:         buildUpdateAdminPasswordParams(d),
	}
	_, err := client.Request("PUT", updatePasswordPath, &updatePasswordOpt)
	if err != nil {
		return fmt.Errorf("error update CBH instance (%s) password: %s", d.Id(), err)
	}
	return nil
}

func buildUpdateAdminPasswordParams(d *schema.ResourceData) map[string]interface{} {
	return map[string]interface{}{
		"new_password": d.Get("password"),
		"server_id":    d.Id(),
	}
}

func resourceCBHInstanceRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		mErr                  *multierror.Error
		cfg                   = meta.(*config.Config)
		region                = cfg.GetRegion(d)
		getCbhInstanceProduct = "cbh"
	)

	client, err := cfg.NewServiceClient(getCbhInstanceProduct, region)
	if err != nil {
		return diag.Errorf("error creating CBH client: %s", err)
	}

	instances, err := getCBHInstanceList(client)
	if err != nil {
		return diag.FromErr(err)
	}

	expression := fmt.Sprintf("[?server_id == '%s']|[0]", d.Id())
	instance := utils.PathSearch(expression, instances, nil)
	if instance == nil {
		return common.CheckDeletedDiag(d, golangsdk.ErrDefault404{}, "")
	}

	// When EIP is not configured, the query interface field will return a space string.
	publicIpId := strings.TrimSpace(utils.PathSearch("network.public_id", instance, "").(string))

	resourceId := utils.PathSearch("resource_info.resource_id", instance, "").(string)
	tags, err := getInstanceTags(resourceId, client)
	if err != nil {
		return diag.FromErr(err)
	}

	mErr = multierror.Append(
		mErr,
		d.Set("region", region),
		d.Set("public_ip_id", publicIpId),
		d.Set("public_ip", utils.PathSearch("network.public_ip", instance, nil)),
		d.Set("name", utils.PathSearch("name", instance, nil)),
		d.Set("private_ip", utils.PathSearch("network.private_ip", instance, nil)),
		d.Set("subnet_address", utils.PathSearch("network.private_ip", instance, nil)),
		d.Set("status", utils.PathSearch("status_info.status", instance, nil)),
		d.Set("vpc_id", utils.PathSearch("network.vpc_id", instance, nil)),
		d.Set("subnet_id", utils.PathSearch("network.subnet_id", instance, nil)),
		d.Set("security_group_id", utils.PathSearch("network.security_group_id", instance, nil)),
		d.Set("flavor_id", utils.PathSearch("resource_info.specification", instance, nil)),
		d.Set("data_disk_size", utils.PathSearch("resource_info.data_disk_size", instance, float64(0)).(float64)),
		d.Set("availability_zone", utils.PathSearch("az_info.zone", instance, nil)),
		d.Set("version", utils.PathSearch("bastion_version", instance, nil)),
		d.Set("enterprise_project_id", utils.PathSearch("enterprise_project_id", instance, nil)),
		d.Set("tags", tags),
	)
	return diag.FromErr(mErr.ErrorOrNil())
}

func getInstanceTags(resourceId string, client *golangsdk.ServiceClient) (
	map[string]interface{}, error) {
	getTagsHttpUrl := "v2/{project_id}/cbs/instance/{resource_id}/tags"
	getTagsPath := client.Endpoint + getTagsHttpUrl
	getTagsPath = strings.ReplaceAll(getTagsPath, "{project_id}", client.ProjectID)
	getTagsPath = strings.ReplaceAll(getTagsPath, "{resource_id}", resourceId)
	getTagsOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	getTagsResp, err := client.Request("GET", getTagsPath, &getTagsOpt)
	if err != nil {
		return nil, fmt.Errorf("error retrieving tags of the CBH instance: %s", err)
	}

	getTagsRespBody, err := utils.FlattenResponse(getTagsResp)
	if err != nil {
		return nil, err
	}

	tags := utils.PathSearch("tags", getTagsRespBody, make([]interface{}, 0)).([]interface{})
	result := make(map[string]interface{})
	for _, val := range tags {
		valMap := val.(map[string]interface{})
		result[valMap["key"].(string)] = valMap["value"]
	}

	return result, nil
}

func getPublicAddressById(d *schema.ResourceData, cfg *config.Config, publicIpId string) (string, error) {
	region := cfg.GetRegion(d)
	networkingClient, err := cfg.NetworkingV1Client(region)
	if err != nil {
		return "", fmt.Errorf("error creating VPC v1 client: %s", err)
	}
	publicIp, err := eips.Get(networkingClient, publicIpId).Extract()
	if err != nil {
		return "", err
	}
	return publicIp.PublicAddress, nil
}

func resourceCBHInstanceDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg                    = meta.(*config.Config)
		region                 = cfg.GetRegion(d)
		getCbhInstancesProduct = "cbh"
	)

	client, err := cfg.NewServiceClient(getCbhInstancesProduct, region)
	if err != nil {
		return diag.Errorf("error creating CBH client: %s", err)
	}

	instances, err := getCBHInstanceList(client)
	if err != nil {
		return diag.FromErr(err)
	}

	expression := fmt.Sprintf("[?server_id == '%s']|[0]", d.Id())
	instance := utils.PathSearch(expression, instances, nil)
	if instance == nil {
		// Before deleting the CBH instance, it is necessary to first call the query API to obtain the resource_id of
		// the instance. If the instance cannot be found, then execute the logic of checkDeleted.
		return common.CheckDeletedDiag(d, golangsdk.ErrDefault404{}, "error deleting the CBH instance")
	}

	resourceId := utils.PathSearch("resource_info.resource_id", instance, "").(string)
	if resourceId == "" {
		return diag.Errorf("error deleting the CBH instance (%s): resource ID is not found in list API response", d.Id())
	}

	if err = common.UnsubscribePrePaidResource(d, cfg, []string{resourceId}); err != nil {
		return diag.Errorf("error unsubscribe CBH instance (%s): %s", d.Id(), err)
	}

	if err := waitingForCBHInstanceDeleted(ctx, client, d, d.Timeout(schema.TimeoutDelete)); err != nil {
		return diag.Errorf("error waiting for CBH instance (%s) deleted: %s", d.Id(), err)
	}
	return nil
}

func waitingForCBHInstanceDeleted(ctx context.Context, client *golangsdk.ServiceClient, d *schema.ResourceData,
	timeout time.Duration) error {
	expression := fmt.Sprintf("[?server_id == '%s']|[0]", d.Id())
	stateConf := &resource.StateChangeConf{
		Pending: []string{"PENDING"},
		Target:  []string{"COMPLETED"},
		Refresh: func() (interface{}, string, error) {
			instances, err := getCBHInstanceList(client)
			if err != nil {
				return nil, "ERROR", err
			}
			instance := utils.PathSearch(expression, instances, nil)
			if instance == nil {
				m := map[string]string{"code": "COMPLETED"}
				return m, "COMPLETED", nil
			}

			status := utils.PathSearch("status_info.status", instance, "").(string)
			if status == "DELETED" {
				return instance, "COMPLETED", nil
			}
			return instance, "PENDING", nil
		},
		Timeout:      timeout,
		Delay:        30 * time.Second,
		PollInterval: 10 * time.Second,
	}

	_, err := stateConf.WaitForStateContext(ctx)
	return err
}
