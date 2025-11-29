package esw

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

var instanceNonUpdatableParams = []string{"flavor_ref", "ha_mode", "availability_zones", "availability_zones.*.primary",
	"availability_zones.*.standby", "tunnel_info", "tunnel_info.*.vpc_id", "tunnel_info.*.virsubnet_id",
	"tunnel_info.*.tunnel_ip", "charge_infos", "charge_infos.*.charge_mode"}

// @API RDS POST /v3/{project_id}/l2cg/instances
// @API RDS GET /v3/{project_id}/l2cg/instances/{instance_id}
// @API RDS PUT /v3/{project_id}/l2cg/instances/{instance_id}
// @API RDS DELETE /v3/{project_id}/l2cg/instances/{instance_id}
func ResourceInstance() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceInstanceCreate,
		ReadContext:   resourceInstanceRead,
		UpdateContext: resourceInstanceUpdate,
		DeleteContext: resourceInstanceDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		CustomizeDiff: config.FlexibleForceNew(instanceNonUpdatableParams),

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(30 * time.Minute),
			Delete: schema.DefaultTimeout(30 * time.Minute),
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
			},
			"flavor_ref": {
				Type:     schema.TypeString,
				Required: true,
			},
			"ha_mode": {
				Type:     schema.TypeString,
				Required: true,
			},
			"availability_zones": {
				Type:     schema.TypeList,
				Required: true,
				MinItems: 1,
				Elem:     instanceAvailabilityZonesSchema(),
			},
			"tunnel_info": {
				Type:     schema.TypeList,
				Required: true,
				MinItems: 1,
				Elem:     instanceTunnelInfoSchema(),
			},
			"charge_infos": {
				Type:     schema.TypeList,
				Required: true,
				MinItems: 1,
				Elem:     instanceChargeInfosSchema(),
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"created_at": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"updated_at": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"enable_force_new": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"true", "false"}, false),
				Description:  utils.SchemaDesc("", utils.SchemaDescInput{Internal: true}),
			},
		},
	}
}

func instanceAvailabilityZonesSchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"primary": {
				Type:     schema.TypeString,
				Required: true,
			},
			"standby": {
				Type:     schema.TypeString,
				Required: true,
			},
		},
	}
	return &sc
}

func instanceTunnelInfoSchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"vpc_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"virsubnet_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"tunnel_ip": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"tunnel_port": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"tunnel_type": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
	return &sc
}

func instanceChargeInfosSchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"charge_mode": {
				Type:     schema.TypeString,
				Required: true,
			},
		},
	}
	return &sc
}

func resourceInstanceCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var (
		httpUrl = "v3/{project_id}/l2cg/instances"
		product = "esw"
	)
	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating ESW client: %s", err)
	}

	createPath := client.Endpoint + httpUrl
	createPath = strings.ReplaceAll(createPath, "{project_id}", client.ProjectID)

	createOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}
	createOpt.JSONBody = utils.RemoveNil(buildCreateInstanceBodyParams(d))
	createResp, err := client.Request("POST", createPath, &createOpt)
	if err != nil {
		return diag.Errorf("error creating ESW instance: %s", err)
	}

	createRespBody, err := utils.FlattenResponse(createResp)
	if err != nil {
		return diag.FromErr(err)
	}

	instanceId := utils.PathSearch("instance.id", createRespBody, "").(string)
	if instanceId == "" {
		return diag.Errorf("unable to find the ESW instance ID from the API response")
	}
	d.SetId(instanceId)

	err = waitForInstanceActive(ctx, client, d)
	if err != nil {
		return diag.FromErr(err)
	}
	return resourceInstanceRead(ctx, d, meta)
}

func buildCreateInstanceBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"name":               d.Get("name"),
		"flavor_ref":         d.Get("flavor_ref"),
		"ha_mode":            d.Get("ha_mode"),
		"availability_zones": buildCreateInstanceAvailabilityZonesBody(d),
		"tunnel_info":        buildCreateInstanceTunnelInfoBody(d),
		"charge_infos":       buildCreateInstanceChargeInfosBody(d),
		"description":        utils.ValueIgnoreEmpty(d.Get("description")),
	}
	return map[string]interface{}{
		"instance": bodyParams,
	}
}

func buildCreateInstanceAvailabilityZonesBody(d *schema.ResourceData) map[string]interface{} {
	rawParams := d.Get("availability_zones").([]interface{})
	if len(rawParams) == 0 {
		return nil
	}

	raw, ok := rawParams[0].(map[string]interface{})
	if !ok {
		return nil
	}

	rst := map[string]interface{}{
		"primary": raw["primary"],
		"standby": raw["standby"],
	}

	return rst
}

func buildCreateInstanceTunnelInfoBody(d *schema.ResourceData) map[string]interface{} {
	rawParams := d.Get("tunnel_info").([]interface{})
	if len(rawParams) == 0 {
		return nil
	}

	raw, ok := rawParams[0].(map[string]interface{})
	if !ok {
		return nil
	}

	rst := map[string]interface{}{
		"vpc_id":       raw["vpc_id"],
		"virsubnet_id": raw["virsubnet_id"],
		"tunnel_ip":    utils.ValueIgnoreEmpty(raw["tunnel_ip"]),
	}

	return rst
}

func buildCreateInstanceChargeInfosBody(d *schema.ResourceData) map[string]interface{} {
	rawParams := d.Get("charge_infos").([]interface{})
	if len(rawParams) == 0 {
		return nil
	}

	raw, ok := rawParams[0].(map[string]interface{})
	if !ok {
		return nil
	}

	rst := map[string]interface{}{
		"charge_mode": raw["charge_mode"],
	}

	return rst
}

func resourceInstanceRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var mErr *multierror.Error

	var (
		product = "esw"
	)
	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating ESW client: %s", err)
	}

	instance, err := getInstanceById(client, d.Id())
	if err != nil {
		return common.CheckDeletedDiag(d, err, "error retrieving ESW instance")
	}

	mErr = multierror.Append(
		mErr,
		d.Set("region", region),
		d.Set("name", utils.PathSearch("instance.name", instance, nil)),
		d.Set("flavor_ref", utils.PathSearch("instance.flavor_ref", instance, nil)),
		d.Set("ha_mode", utils.PathSearch("instance.ha_mode", instance, nil)),
		d.Set("availability_zones", flattenInstanceAvailabilityZones(instance)),
		d.Set("tunnel_info", flattenInstanceTunnelInfo(instance)),
		d.Set("charge_infos", flattenInstanceChargeInfos(instance)),
		d.Set("description", utils.PathSearch("instance.description", instance, nil)),
		d.Set("status", utils.PathSearch("instance.status", instance, nil)),
		d.Set("created_at", utils.PathSearch("instance.created_at", instance, nil)),
		d.Set("updated_at", utils.PathSearch("instance.updated_at", instance, nil)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenInstanceAvailabilityZones(resp interface{}) []interface{} {
	curJson := utils.PathSearch("instance.availability_zones", resp, nil)
	if curJson == nil {
		return nil
	}

	rst := []interface{}{
		map[string]interface{}{
			"primary": utils.PathSearch("primary", curJson, nil),
			"standby": utils.PathSearch("standby", curJson, nil),
		},
	}
	return rst
}

func flattenInstanceTunnelInfo(resp interface{}) []interface{} {
	curJson := utils.PathSearch("instance.tunnel_info", resp, nil)
	if curJson == nil {
		return nil
	}

	rst := []interface{}{
		map[string]interface{}{
			"vpc_id":       utils.PathSearch("vpc_id", curJson, nil),
			"virsubnet_id": utils.PathSearch("virsubnet_id", curJson, nil),
			"tunnel_ip":    utils.PathSearch("tunnel_ip", curJson, nil),
			"tunnel_port":  int(utils.PathSearch("tunnel_port", curJson, float64(0)).(float64)),
			"tunnel_type":  utils.PathSearch("tunnel_type", curJson, nil),
		},
	}
	return rst
}

func flattenInstanceChargeInfos(resp interface{}) []interface{} {
	curJson := utils.PathSearch("instance.charge_infos", resp, nil)
	if curJson == nil {
		return nil
	}

	rst := []interface{}{
		map[string]interface{}{
			"charge_mode": utils.PathSearch("charge_mode", curJson, nil),
		},
	}
	return rst
}

func getInstanceById(client *golangsdk.ServiceClient, instanceId string) (interface{}, error) {
	var (
		httpUrl = "v3/{project_id}/l2cg/instances/{instance_id}"
	)

	getPath := client.Endpoint + httpUrl
	getPath = strings.ReplaceAll(getPath, "{project_id}", client.ProjectID)
	getPath = strings.ReplaceAll(getPath, "{instance_id}", instanceId)

	getOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}
	getResp, err := client.Request("GET", getPath, &getOpt)
	if err != nil {
		return nil, err
	}

	return utils.FlattenResponse(getResp)
}

func resourceInstanceUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var (
		httpUrl = "v3/{project_id}/l2cg/instances/{instance_id}"
		product = "esw"
	)
	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating ESW Client: %s", err)
	}

	updatePath := client.Endpoint + httpUrl
	updatePath = strings.ReplaceAll(updatePath, "{project_id}", client.ProjectID)
	updatePath = strings.ReplaceAll(updatePath, "{instance_id}", d.Id())

	updateOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}
	updateOpt.JSONBody = buildUpdateInstanceBodyParams(d)

	_, err = client.Request("PUT", updatePath, &updateOpt)
	if err != nil {
		return diag.Errorf("error updating ESW instance: %s", err)
	}

	return resourceInstanceRead(ctx, d, meta)
}

func buildUpdateInstanceBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"name":        d.Get("name"),
		"description": d.Get("description"),
	}
	return map[string]interface{}{
		"instance": bodyParams,
	}
}

func resourceInstanceDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var (
		httpUrl = "v3/{project_id}/l2cg/instances/{instance_id}"
		product = "esw"
	)
	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating ESW client: %s", err)
	}

	deletePath := client.Endpoint + httpUrl
	deletePath = strings.ReplaceAll(deletePath, "{project_id}", client.ProjectID)
	deletePath = strings.ReplaceAll(deletePath, "{instance_id}", d.Id())

	deleteOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}
	_, err = client.Request("DELETE", deletePath, &deleteOpt)
	if err != nil {
		return common.CheckDeletedDiag(d, err, "error deleting ESW instance")
	}

	err = waitForInstanceDeleted(ctx, client, d)
	if err != nil {
		return diag.FromErr(err)
	}
	return nil
}

func waitForInstanceActive(ctx context.Context, client *golangsdk.ServiceClient, d *schema.ResourceData) error {
	stateConf := &resource.StateChangeConf{
		Pending:      []string{"Pending"},
		Target:       []string{"Active"},
		Refresh:      instanceStatusRefreshFunc(client, d.Id()),
		Timeout:      d.Timeout(schema.TimeoutCreate),
		Delay:        10 * time.Second,
		PollInterval: 10 * time.Second,
	}
	if _, err := stateConf.WaitForStateContext(ctx); err != nil {
		return fmt.Errorf("error waiting for ESW instance(%s) to active: %s ", d.Id(), err)
	}
	return nil
}

func waitForInstanceDeleted(ctx context.Context, client *golangsdk.ServiceClient, d *schema.ResourceData) error {
	stateConf := &resource.StateChangeConf{
		Pending:      []string{"Pending"},
		Target:       []string{"Deleted"},
		Refresh:      instanceStatusRefreshFunc(client, d.Id()),
		Timeout:      d.Timeout(schema.TimeoutDelete),
		Delay:        10 * time.Second,
		PollInterval: 10 * time.Second,
	}
	if _, err := stateConf.WaitForStateContext(ctx); err != nil {
		return fmt.Errorf("error waiting for ESW instance(%s) to be deleted: %s ", d.Id(), err)
	}
	return nil
}

func instanceStatusRefreshFunc(client *golangsdk.ServiceClient, instanceId string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		instance, err := getInstanceById(client, instanceId)
		if err != nil {
			var errDefault404 golangsdk.ErrDefault404
			if errors.As(err, &errDefault404) {
				return "", "Deleted", nil
			}
			return nil, "Failed", err
		}

		status := utils.PathSearch("instance.status", instance, "").(string)
		if status == "" {
			return nil, "Failed", errors.New("status is not found")
		}

		if utils.StrSliceContains([]string{"failed", "abnormal"}, status) {
			return instance, "Failed", fmt.Errorf("the isntance status is: %s", status)
		}
		if status == "active" {
			return instance, "Active", nil
		}

		return instance, "Pending", nil
	}
}
