package dsc

import (
	"context"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/cbc"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

var adgInstanceNonUpdatableParams = []string{
	"name",
	"specification",
	"vpc_id",
	"subnet_id",
	"security_group_id",
	"availability_zone",
	"charge_mode",
	"period_unit",
	"period",
	"deploy_mode",
	"mode",
	"type",
	"outside_ins_config",
	"outside_ins_config.*.master_node_ip",
	"outside_ins_config.*.slave_node_ip",
	"outside_ins_config.*.virtual_ip",
}

// Due to API issues, currently only the prepaid mode is supported. When the API is fixed, the postpaid mode will be
// expanded.
// @API DSC POST /v1/{project_id}/instances
// @API DSC GET /v1/{project_id}/instances/{instance_id}
// @API DSC POST /v1/{project_id}/instances/{instance_id}/reset-password
// @API DSC POST /v1/{project_id}/instances/{instance_id}/bind-eip
// @API DSC POST /v1/{project_id}/instances/{instance_id}/unbind-eip
func ResourceDscAdgInstance() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceDscAdgInstanceCreate,
		ReadContext:   resourceDscAdgInstanceRead,
		UpdateContext: resourceDscAdgInstanceUpdate,
		DeleteContext: resourceDscAdgInstanceDelete,

		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(30 * time.Minute),
			Delete: schema.DefaultTimeout(30 * time.Minute),
		},

		CustomizeDiff: config.FlexibleForceNew(adgInstanceNonUpdatableParams),

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
				Description: `Specifies the name of the ADG instance.`,
			},
			"type": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `Specifies the type of the ADG instance.`,
			},
			"specification": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `Specifies the specification of the ADG instance.`,
			},
			"vpc_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `Specifies the VPC ID to which the ADG instance belongs.`,
			},
			"subnet_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `Specifies the subnet ID to which the ADG instance belongs.`,
			},
			"security_group_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `Specifies the security group ID of the ADG instance.`,
			},
			"availability_zone": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `Specifies the availability zone of the ADG instance.`,
			},
			"charge_mode": {
				Type:     schema.TypeString,
				Required: true,
				ValidateFunc: validation.StringInSlice([]string{
					"prePaid",
				}, false),
				Description: `Specifies the charging mode of the ADG instance. Currently, only **prePaid** is supported.`,
			},
			"period_unit": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringInSlice([]string{"month", "year"}, false),
				Description:  `Specifies the charging period unit of the ADG instance.`,
			},
			"period": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: `Specifies the charging period of the ADG instance.`,
			},
			"deploy_mode": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `Specifies the deploy mode of the ADG instance.`,
			},
			"mode": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `Specifies the mode of the ADG instance.`,
			},
			"password": {
				Type:        schema.TypeString,
				Optional:    true,
				Sensitive:   true,
				Description: `Specifies the password of the ADG instance.`,
			},
			"admin_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specifies the administrator name. Valid values are sysadmin, secadmin, audadmin.`,
			},
			"publicip_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specifies the public IP ID to bind to the ADG instance.`,
			},
			"outside_ins_config": {
				Type:        schema.TypeList,
				Optional:    true,
				MaxItems:    1,
				Elem:        adgInstanceOutsideInsConfigSchema(),
				Description: `Specifies the cloud outside instance configuration.`,
			},
			"auto_renew": common.SchemaAutoRenewUpdatable(nil),
			"enable_force_new": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"true", "false"}, false),
				Description:  utils.SchemaDesc("", utils.SchemaDescInput{Internal: true}),
			},
			"status": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: `The status of the ADG instance.`,
			},
			"public_ip": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The public IP address of the ADG instance.`,
			},
			"virtual_ip": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The virtual IP address of the ADG instance.`,
			},
			"version": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The version of the ADG instance.`,
			},
			"create_time": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: `The creation time of the ADG instance.`,
			},
			"started_time": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: `The start time of the ADG instance.`,
			},
			"fail_reason": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The failure reason of the ADG instance.`,
			},
			"nodes": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: `The node information of the ADG instance.`,
				Elem:        adgInstanceNodesSchema(),
			},
		},
	}
}

func adgInstanceOutsideInsConfigSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"master_node_ip": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specifies the master node IP address.`,
			},
			"slave_node_ip": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specifies the slave node IP address.`,
			},
			"virtual_ip": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specifies the virtual IP address.`,
			},
		},
	}
}

func adgInstanceNodesSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The node ID.`,
			},
			"name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The node name.`,
			},
			"availability_zone": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The availability zone of the node.`,
			},
			"private_ip": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The private IP address of the node.`,
			},
			"role": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The role of the node.`,
			},
			"status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The status of the node.`,
			},
			"vm_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The VM ID of the node.`,
			},
			"error_reason": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The error reason of the node.`,
			},
		},
	}
}

func resourceDscAdgInstanceCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	client, err := cfg.NewServiceClient("dsc", region)
	if err != nil {
		return diag.Errorf("error creating DSC client: %s", err)
	}

	httpUrl := "v1/{project_id}/instances"
	createPath := client.Endpoint + httpUrl
	createPath = strings.ReplaceAll(createPath, "{project_id}", client.ProjectID)

	createOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"Content-Type": "application/json"},
		JSONBody:         utils.RemoveNil(buildCreateAdgInstanceBodyParams(d, cfg)),
	}

	createResp, err := client.Request("POST", createPath, &createOpt)
	if err != nil {
		return diag.Errorf("error creating DSC ADG instance: %s", err)
	}

	createRespBody, err := utils.FlattenResponse(createResp)
	if err != nil {
		return diag.FromErr(err)
	}

	orderId := utils.PathSearch("order_id", createRespBody, "").(string)
	if orderId == "" {
		return diag.Errorf("error creating DSC ADG instance: order_id is not found in API response")
	}

	bssClient, err := cfg.BssV2Client(region)
	if err != nil {
		return diag.Errorf("error creating BSS v2 client: %s", err)
	}
	if err = common.WaitOrderComplete(ctx, bssClient, orderId, d.Timeout(schema.TimeoutCreate)); err != nil {
		return diag.FromErr(err)
	}
	resourceId, err := common.WaitOrderResourceComplete(ctx, bssClient, orderId, d.Timeout(schema.TimeoutCreate))
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(resourceId)

	if adminName, ok := d.GetOk("admin_name"); ok && adminName.(string) != "" {
		if password, ok := d.GetOk("password"); ok && password.(string) != "" {
			if err = resetAdgInstancePassword(client, d.Id(), d); err != nil {
				return diag.Errorf("error resetting password for DSC ADG instance (%s): %s", d.Id(), err)
			}
		}
	}

	return resourceDscAdgInstanceRead(ctx, d, meta)
}

func buildCreateAdgInstanceBodyParams(d *schema.ResourceData, cfg *config.Config) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"availability_zone":  d.Get("availability_zone"),
		"deploy_mode":        utils.ValueIgnoreEmpty(d.Get("deploy_mode")),
		"mode":               utils.ValueIgnoreEmpty(d.Get("mode")),
		"name":               d.Get("name"),
		"password":           utils.ValueIgnoreEmpty(d.Get("password")),
		"publicip_id":        utils.ValueIgnoreEmpty(d.Get("publicip_id")),
		"region":             cfg.GetRegion(d),
		"security_group_id":  d.Get("security_group_id"),
		"specification":      d.Get("specification"),
		"subnet_id":          d.Get("subnet_id"),
		"type":               d.Get("type"),
		"vpc_id":             d.Get("vpc_id"),
		"charge_info":        buildAdgInstanceChargeInfoBodyParams(d),
		"outside_ins_config": buildAdgInstanceOutsideInsConfigBodyParams(d.Get("outside_ins_config").([]interface{})),
	}

	return bodyParams
}

func buildAdgInstanceOutsideInsConfigBodyParams(outsideInsConfig []interface{}) map[string]interface{} {
	if len(outsideInsConfig) == 0 || outsideInsConfig[0] == nil {
		return nil
	}

	insConfig := outsideInsConfig[0].(map[string]interface{})
	return map[string]interface{}{
		"master_node_ip": utils.ValueIgnoreEmpty(insConfig["master_node_ip"]),
		"slave_node_ip":  utils.ValueIgnoreEmpty(insConfig["slave_node_ip"]),
		"virtual_ip":     utils.ValueIgnoreEmpty(insConfig["virtual_ip"]),
	}
}

func buildAdgInstanceChargeInfoBodyParams(d *schema.ResourceData) map[string]interface{} {
	chargeInfo := map[string]interface{}{
		"charge_mode":   "prePaid",
		"is_auto_pay":   true,
		"is_auto_renew": parseAutoRenewValue(d.Get("auto_renew").(string)),
		"period_type":   d.Get("period_unit").(string),
		"period_num":    d.Get("period"),
	}

	return chargeInfo
}

func parseAutoRenewValue(autoRenew string) bool {
	result, err := strconv.ParseBool(autoRenew)
	if err != nil {
		log.Printf("[WARN] unable to convert auto_renew to bool value")
		return false
	}
	return result
}

func GetAdgInstanceById(client *golangsdk.ServiceClient, instanceId string) (interface{}, error) {
	httpUrl := "v1/{project_id}/instances/{instance_id}"
	getPath := client.Endpoint + httpUrl
	getPath = strings.ReplaceAll(getPath, "{project_id}", client.ProjectID)
	getPath = strings.ReplaceAll(getPath, "{instance_id}", instanceId)

	getOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	resp, err := client.Request("GET", getPath, &getOpt)
	if err != nil {
		return nil, common.ConvertExpected400ErrInto404Err(err, "error_code", "dsc.40000059")
	}

	return utils.FlattenResponse(resp)
}

func resourceDscAdgInstanceRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	client, err := cfg.NewServiceClient("dsc", region)
	if err != nil {
		return diag.Errorf("error creating DSC client: %s", err)
	}

	instance, err := GetAdgInstanceById(client, d.Id())
	if err != nil {
		return common.CheckDeletedDiag(d, err, "error retrieving DSC ADG instance")
	}

	mErr := multierror.Append(
		d.Set("region", region),
		d.Set("name", utils.PathSearch("name", instance, nil)),
		d.Set("specification", utils.PathSearch("specification", instance, nil)),
		d.Set("vpc_id", utils.PathSearch("vpc_id", instance, nil)),
		d.Set("subnet_id", utils.PathSearch("subnet_id", instance, nil)),
		d.Set("security_group_id", utils.PathSearch("security_group_id", instance, nil)),
		d.Set("availability_zone", utils.PathSearch("availability_zone", instance, nil)),
		d.Set("deploy_mode", utils.PathSearch("deploy_mode", instance, nil)),
		d.Set("mode", utils.PathSearch("mode", instance, nil)),
		d.Set("type", utils.PathSearch("type", instance, nil)),
		d.Set("publicip_id", utils.PathSearch("publicip_id", instance, nil)),
		d.Set("status", utils.PathSearch("status", instance, nil)),
		d.Set("public_ip", utils.PathSearch("public_ip", instance, nil)),
		d.Set("virtual_ip", utils.PathSearch("virtual_ip", instance, nil)),
		d.Set("version", utils.PathSearch("version", instance, nil)),
		d.Set("fail_reason", utils.PathSearch("fail_reason", instance, nil)),
		d.Set("nodes", flattenAdgInstanceNodes(instance)),
		d.Set("charge_mode", utils.PathSearch("charge_mode", instance, nil)),
		d.Set("create_time", utils.PathSearch("create_time", instance, nil)),
		d.Set("started_time", utils.PathSearch("started_time", instance, nil)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenAdgInstanceNodes(resp interface{}) []map[string]interface{} {
	if resp == nil {
		return nil
	}

	nodesRaw := utils.PathSearch("nodes", resp, make([]interface{}, 0)).([]interface{})
	if len(nodesRaw) == 0 {
		return nil
	}

	result := make([]map[string]interface{}, len(nodesRaw))
	for i, node := range nodesRaw {
		result[i] = map[string]interface{}{
			"id":                utils.PathSearch("id", node, nil),
			"name":              utils.PathSearch("name", node, nil),
			"availability_zone": utils.PathSearch("availability_zone", node, nil),
			"private_ip":        utils.PathSearch("private_ip", node, nil),
			"role":              utils.PathSearch("role", node, nil),
			"status":            utils.PathSearch("status", node, nil),
			"vm_id":             utils.PathSearch("vm_id", node, nil),
			"error_reason":      utils.PathSearch("error_reason", node, nil),
		}
	}
	return result
}

func resourceDscAdgInstanceUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	client, err := cfg.NewServiceClient("dsc", region)
	if err != nil {
		return diag.Errorf("error creating DSC client: %s", err)
	}

	instanceId := d.Id()

	// reset password
	if d.HasChanges("password", "admin_name") {
		if err = resetAdgInstancePassword(client, instanceId, d); err != nil {
			return diag.Errorf("error resetting password for DSC ADG instance (%s): %s", instanceId, err)
		}
	}

	// bind/unbind EIP
	if d.HasChange("publicip_id") {
		oldPublicipId, newPublicipId := d.GetChange("publicip_id")
		// unbind old EIP first
		if oldPublicipId.(string) != "" {
			if err = unbindAdgInstanceEip(client, instanceId, oldPublicipId.(string)); err != nil {
				return diag.Errorf("error unbinding EIP from DSC ADG instance (%s): %s", instanceId, err)
			}
		}
		// bind new EIP
		if newPublicipId.(string) != "" {
			if err = bindAdgInstanceEip(client, instanceId, newPublicipId.(string)); err != nil {
				return diag.Errorf("error binding EIP to DSC ADG instance (%s): %s", instanceId, err)
			}
		}
	}

	// update auto_renew
	if d.HasChange("auto_renew") {
		bssClient, err := cfg.BssV2Client(region)
		if err != nil {
			return diag.Errorf("error creating BSS V2 client: %s", err)
		}
		if err = cbc.UpdateAutoRenew(bssClient, d.Get("auto_renew").(string), instanceId); err != nil {
			return diag.Errorf("error updating the auto-renew of the DSC ADG instance (%s): %s", instanceId, err)
		}
	}

	return resourceDscAdgInstanceRead(ctx, d, meta)
}

func resetAdgInstancePassword(client *golangsdk.ServiceClient, instanceId string, d *schema.ResourceData) error {
	httpUrl := "v1/{project_id}/instances/{instance_id}/reset-password"
	resetPath := client.Endpoint + httpUrl
	resetPath = strings.ReplaceAll(resetPath, "{project_id}", client.ProjectID)
	resetPath = strings.ReplaceAll(resetPath, "{instance_id}", instanceId)

	resetOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"Content-Type": "application/json"},
		JSONBody: map[string]interface{}{
			"admin_name": d.Get("admin_name").(string),
			"admin_pw":   d.Get("password").(string),
		},
	}

	_, err := client.Request("POST", resetPath, &resetOpt)
	return err
}

func bindAdgInstanceEip(client *golangsdk.ServiceClient, instanceId, eipId string) error {
	httpUrl := "v1/{project_id}/instances/{instance_id}/bind-eip"
	bindPath := client.Endpoint + httpUrl
	bindPath = strings.ReplaceAll(bindPath, "{project_id}", client.ProjectID)
	bindPath = strings.ReplaceAll(bindPath, "{instance_id}", instanceId)

	bindOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"Content-Type": "application/json"},
		JSONBody: map[string]interface{}{
			"eip_id": eipId,
		},
	}

	_, err := client.Request("POST", bindPath, &bindOpt)
	return err
}

func unbindAdgInstanceEip(client *golangsdk.ServiceClient, instanceId, eipId string) error {
	httpUrl := "v1/{project_id}/instances/{instance_id}/unbind-eip"
	unbindPath := client.Endpoint + httpUrl
	unbindPath = strings.ReplaceAll(unbindPath, "{project_id}", client.ProjectID)
	unbindPath = strings.ReplaceAll(unbindPath, "{instance_id}", instanceId)

	unbindOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"Content-Type": "application/json"},
		JSONBody: map[string]interface{}{
			"eip_id": eipId,
		},
	}

	_, err := client.Request("POST", unbindPath, &unbindOpt)
	return err
}

func resourceDscAdgInstanceDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	client, err := cfg.NewServiceClient("dsc", region)
	if err != nil {
		return diag.Errorf("error creating DSC client: %s", err)
	}

	instanceId := d.Id()

	if err = common.UnsubscribePrePaidResource(d, cfg, []string{instanceId}); err != nil {
		return diag.Errorf("error unsubscribing DSC ADG instance (%s): %s", instanceId, err)
	}

	// wait for instance to be deleted
	stateConf := &retry.StateChangeConf{
		Pending:      []string{"PENDING"},
		Target:       []string{"COMPLETED"},
		Refresh:      refreshAdgInstanceDeleted(client, instanceId),
		Timeout:      d.Timeout(schema.TimeoutDelete),
		Delay:        10 * time.Second,
		PollInterval: 10 * time.Second,
	}

	_, err = stateConf.WaitForStateContext(ctx)
	if err != nil {
		return diag.Errorf("error waiting for DSC ADG instance (%s) to be deleted: %s", instanceId, err)
	}

	return nil
}

func refreshAdgInstanceDeleted(client *golangsdk.ServiceClient, instanceId string) retry.StateRefreshFunc {
	return func() (interface{}, string, error) {
		respBody, err := GetAdgInstanceById(client, instanceId)
		if err != nil {
			if _, ok := err.(golangsdk.ErrDefault404); ok {
				return "deleted", "COMPLETED", nil
			}
			return nil, "ERROR", err
		}

		return respBody, "PENDING", nil
	}
}
