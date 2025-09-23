package elb

import (
	"context"
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

var poolNonUpdatableParams = []string{"protocol", "loadbalancer_id", "listener_id", "ip_version", "any_port_enable",
	"public_border_group"}

// @API ELB POST /v3/{project_id}/elb/pools
// @API ELB GET /v3/{project_id}/elb/pools/{pool_id}
// @API ELB PUT /v3/{project_id}/elb/pools/{pool_id}
// @API ELB DELETE /v3/{project_id}/elb/pools/{pool_id}
func ResourcePoolV3() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourcePoolV3Create,
		ReadContext:   resourcePoolV3Read,
		UpdateContext: resourcePoolV3Update,
		DeleteContext: resourcePoolV3Delete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		CustomizeDiff: config.FlexibleForceNew(poolNonUpdatableParams),

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
			Update: schema.DefaultTimeout(10 * time.Minute),
			Delete: schema.DefaultTimeout(10 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"protocol": {
				Type:     schema.TypeString,
				Required: true,
			},
			"loadbalancer_id": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				AtLeastOneOf: []string{"loadbalancer_id", "listener_id", "type"},
			},
			"listener_id": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				AtLeastOneOf: []string{"loadbalancer_id", "listener_id", "type"},
			},
			"type": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				AtLeastOneOf: []string{"loadbalancer_id", "listener_id", "type"},
			},
			"lb_method": {
				Type:     schema.TypeString,
				Required: true,
			},
			"persistence": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"type": {
							Type:     schema.TypeString,
							Required: true,
						},
						"cookie_name": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"timeout": {
							Type:     schema.TypeInt,
							Optional: true,
							Computed: true,
						},
					},
				},
			},
			"name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"vpc_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"protection_status": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"protection_reason": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"slow_start_enabled": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},
			"slow_start_duration": {
				Type:         schema.TypeInt,
				Optional:     true,
				Computed:     true,
				RequiredWith: []string{"slow_start_enabled"},
			},
			"ip_version": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"any_port_enable": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},
			"deletion_protection_enable": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"connection_drain_enabled": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},
			"connection_drain_timeout": {
				Type:         schema.TypeInt,
				Optional:     true,
				Computed:     true,
				RequiredWith: []string{"connection_drain_enabled"},
			},
			"minimum_healthy_member_count": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},
			"public_border_group": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"enable_force_new": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"true", "false"}, false),
				Description:  utils.SchemaDesc("", utils.SchemaDescInput{Internal: true}),
			},
			"monitor_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"enterprise_project_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"updated_at": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"created_at": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourcePoolV3Create(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conf := meta.(*config.Config)
	region := conf.GetRegion(d)

	var (
		httpUrl = "v3/{project_id}/elb/pools"
		product = "elb"
	)
	client, err := conf.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating ELB client: %s", err)
	}

	createPath := client.Endpoint + httpUrl
	createPath = strings.ReplaceAll(createPath, "{project_id}", client.ProjectID)

	createOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}
	createOpt.JSONBody = utils.RemoveNil(buildCreatePoolBodyParams(d))
	createResp, err := client.Request("POST", createPath, &createOpt)
	if err != nil {
		return diag.Errorf("error creating ELB pool: %s", err)
	}

	createRespBody, err := utils.FlattenResponse(createResp)
	if err != nil {
		return diag.Errorf("error retrieving ELB pool: %s", err)
	}
	poolId := utils.PathSearch("pool.id", createRespBody, "").(string)
	if poolId == "" {
		return diag.Errorf("error creating ELB pool: ID is not found in API response")
	}

	d.SetId(poolId)

	err = waitForPool(ctx, client, d.Id(), "ACTIVE", nil, d.Timeout(schema.TimeoutCreate))
	if err != nil {
		return diag.FromErr(err)
	}

	return resourcePoolV3Read(ctx, d, meta)
}

func buildCreatePoolBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"name":                              utils.ValueIgnoreEmpty(d.Get("name")),
		"description":                       utils.ValueIgnoreEmpty(d.Get("description")),
		"protocol":                          d.Get("protocol"),
		"loadbalancer_id":                   utils.ValueIgnoreEmpty(d.Get("loadbalancer_id")),
		"listener_id":                       utils.ValueIgnoreEmpty(d.Get("listener_id")),
		"lb_algorithm":                      d.Get("lb_method"),
		"protection_status":                 utils.ValueIgnoreEmpty(d.Get("protection_status")),
		"protection_reason":                 utils.ValueIgnoreEmpty(d.Get("protection_reason")),
		"type":                              utils.ValueIgnoreEmpty(d.Get("type")),
		"vpc_id":                            utils.ValueIgnoreEmpty(d.Get("vpc_id")),
		"ip_version":                        utils.ValueIgnoreEmpty(d.Get("ip_version")),
		"any_port_enable":                   utils.ValueIgnoreEmpty(d.Get("any_port_enable")),
		"member_deletion_protection_enable": utils.ValueIgnoreEmpty(d.Get("deletion_protection_enable")),
		"public_border_group":               utils.ValueIgnoreEmpty(d.Get("public_border_group")),
		"session_persistence":               buildPersistence(d),
	}
	if _, ok := d.GetOk("slow_start_enabled"); ok {
		bodyParams["slow_start"] = buildSlowStart(d)
	}
	if _, ok := d.GetOk("connection_drain_enabled"); ok {
		bodyParams["connection_drain"] = buildConnectionDrain(d)
	}
	if _, ok := d.GetOk("minimum_healthy_member_count"); ok {
		bodyParams["pool_health"] = buildPoolHealth(d)
	}
	return map[string]interface{}{"pool": bodyParams}
}

func buildSlowStart(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"enable":   d.Get("slow_start_enabled").(bool),
		"duration": d.Get("slow_start_duration").(int),
	}
	return bodyParams
}

func buildConnectionDrain(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"enable":  d.Get("connection_drain_enabled").(bool),
		"timeout": d.Get("connection_drain_timeout").(int),
	}
	return bodyParams
}

func buildPoolHealth(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"minimum_healthy_member_count": d.Get("minimum_healthy_member_count").(int),
	}
	return bodyParams
}

func buildPersistence(d *schema.ResourceData) map[string]interface{} {
	rawPersistence, ok := d.GetOk("persistence")
	if !ok {
		return nil
	}

	persistence, ok := rawPersistence.([]interface{})[0].(map[string]interface{})
	if !ok {
		return nil
	}

	bodyParams := map[string]interface{}{
		"type": persistence["type"],
	}
	if persistence["timeout"].(int) != 0 {
		bodyParams["persistence_timeout"] = persistence["timeout"]
	}
	if persistence["cookie_name"].(string) != "" {
		bodyParams["cookie_name"] = persistence["cookie_name"]
	}
	return bodyParams
}

func resourcePoolV3Read(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var mErr *multierror.Error

	var (
		httpUrl = "v3/{project_id}/elb/pools/{pool_id}"
		product = "elb"
	)
	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating ELB client: %s", err)
	}

	getPath := client.Endpoint + httpUrl
	getPath = strings.ReplaceAll(getPath, "{project_id}", client.ProjectID)
	getPath = strings.ReplaceAll(getPath, "{pool_id}", d.Id())

	getOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	getResp, err := client.Request("GET", getPath, &getOpt)
	if err != nil {
		return common.CheckDeletedDiag(d, err, "error retrieving ELB pool")
	}

	getRespBody, err := utils.FlattenResponse(getResp)
	if err != nil {
		return diag.FromErr(err)
	}

	mErr = multierror.Append(
		mErr,
		d.Set("region", region),
		d.Set("lb_method", utils.PathSearch("pool.lb_algorithm", getRespBody, nil)),
		d.Set("protocol", utils.PathSearch("pool.protocol", getRespBody, nil)),
		d.Set("description", utils.PathSearch("pool.description", getRespBody, nil)),
		d.Set("name", utils.PathSearch("pool.name", getRespBody, nil)),
		d.Set("type", utils.PathSearch("pool.type", getRespBody, nil)),
		d.Set("vpc_id", utils.PathSearch("pool.vpc_id", getRespBody, nil)),
		d.Set("protection_status", utils.PathSearch("pool.protection_status", getRespBody, nil)),
		d.Set("protection_reason", utils.PathSearch("pool.protection_reason", getRespBody, nil)),
		d.Set("slow_start_enabled", utils.PathSearch("pool.slow_start.enable", getRespBody, nil)),
		d.Set("slow_start_duration", utils.PathSearch("pool.slow_start.duration", getRespBody, nil)),
		d.Set("connection_drain_enabled", utils.PathSearch("pool.connection_drain.enable", getRespBody, nil)),
		d.Set("connection_drain_timeout", utils.PathSearch("pool.connection_drain.timeout", getRespBody, nil)),
		d.Set("minimum_healthy_member_count", utils.PathSearch("pool.pool_health.minimum_healthy_member_count",
			getRespBody, nil)),
		d.Set("ip_version", utils.PathSearch("pool.ip_version", getRespBody, nil)),
		d.Set("any_port_enable", utils.PathSearch("pool.any_port_enable", getRespBody, nil)),
		d.Set("deletion_protection_enable", utils.PathSearch("pool.member_deletion_protection_enable",
			getRespBody, nil)),
		d.Set("public_border_group", utils.PathSearch("pool.public_border_group", getRespBody, nil)),
		d.Set("loadbalancer_id", utils.PathSearch("pool.loadbalancers[0].id", getRespBody, nil)),
		d.Set("listener_id", utils.PathSearch("pool.listeners[0].id", getRespBody, nil)),
		d.Set("persistence", flattenPoolPersistence(getRespBody)),
		d.Set("monitor_id", utils.PathSearch("pool.healthmonitor_id", getRespBody, nil)),
		d.Set("enterprise_project_id", utils.PathSearch("pool.enterprise_project_id", getRespBody, nil)),
		d.Set("created_at", utils.PathSearch("pool.created_at", getRespBody, nil)),
		d.Set("updated_at", utils.PathSearch("pool.updated_at", getRespBody, nil)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenPoolPersistence(resp interface{}) []map[string]interface{} {
	persistence := utils.PathSearch("pool.session_persistence", resp, nil)
	if persistence == nil {
		return nil
	}

	rst := []map[string]interface{}{
		{
			"cookie_name": utils.PathSearch("cookie_name", persistence, nil),
			"type":        utils.PathSearch("type", persistence, nil),
			"timeout":     utils.PathSearch("persistence_timeout", persistence, nil),
		},
	}
	return rst
}

func resourcePoolV3Update(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var (
		httpUrl = "v3/{project_id}/elb/pools/{pool_id}"
		product = "elb"
	)
	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating ELB client: %s", err)
	}

	updatePath := client.Endpoint + httpUrl
	updatePath = strings.ReplaceAll(updatePath, "{project_id}", client.ProjectID)
	updatePath = strings.ReplaceAll(updatePath, "{pool_id}", d.Id())

	updateOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}
	updateOpt.JSONBody = buildUpdatePoolBodyParams(d)
	_, err = client.Request("PUT", updatePath, &updateOpt)
	if err != nil {
		return diag.Errorf("error updating ELB pool: %s", err)
	}

	err = waitForPool(ctx, client, d.Id(), "ACTIVE", nil, d.Timeout(schema.TimeoutUpdate))
	if err != nil {
		return diag.FromErr(err)
	}

	return resourcePoolV3Read(ctx, d, meta)
}

func buildUpdatePoolBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"lb_algorithm":                      utils.ValueIgnoreEmpty(d.Get("lb_method")),
		"name":                              d.Get("name"),
		"description":                       d.Get("description"),
		"session_persistence":               buildPersistence(d),
		"protection_status":                 d.Get("protection_status"),
		"protection_reason":                 d.Get("protection_reason"),
		"member_deletion_protection_enable": d.Get("deletion_protection_enable"),
	}
	if d.HasChange("type") {
		bodyParams["type"] = d.Get("type")
	}
	if d.HasChange("vpc_id") {
		bodyParams["vpc_id"] = d.Get("vpc_id")
	}
	if d.HasChanges("slow_start_enabled", "slow_start_duration") {
		bodyParams["slow_start"] = buildSlowStart(d)
	}
	if d.HasChanges("connection_drain_enabled", "connection_drain_timeout") {
		bodyParams["connection_drain"] = buildConnectionDrain(d)
	}
	if d.HasChange("minimum_healthy_member_count") {
		bodyParams["pool_health"] = buildPoolHealth(d)
	}
	return map[string]interface{}{"pool": bodyParams}
}

func resourcePoolV3Delete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var (
		httpUrl = "v3/{project_id}/elb/pools/{pool_id}"
		product = "elb"
	)

	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating ELB client: %s", err)
	}

	deletePath := client.Endpoint + httpUrl
	deletePath = strings.ReplaceAll(deletePath, "{project_id}", client.ProjectID)
	deletePath = strings.ReplaceAll(deletePath, "{pool_id}", d.Id())

	deleteOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}
	_, err = client.Request("DELETE", deletePath, &deleteOpt)
	if err != nil {
		return common.CheckDeletedDiag(d, err, "error deleting ELB pool")
	}

	err = waitForPool(ctx, client, d.Id(), "DELETED", nil, d.Timeout(schema.TimeoutDelete))
	if err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func waitForPool(ctx context.Context, client *golangsdk.ServiceClient, id string, target string, pending []string,
	timeout time.Duration) error {
	stateConf := &resource.StateChangeConf{
		Target:     []string{target},
		Pending:    pending,
		Refresh:    resourcePoolRefreshFunc(client, id),
		Timeout:    timeout,
		MinTimeout: 3 * time.Second,
	}

	_, err := stateConf.WaitForStateContext(ctx)
	if err != nil {
		if _, ok := err.(golangsdk.ErrDefault404); ok {
			switch target {
			case "DELETED":
				return nil
			default:
				return fmt.Errorf("error: pool %s not found: %s", id, err)
			}
		}
		return fmt.Errorf("error waiting for pool %s to become %s: %s", id, target, err)
	}

	return nil
}

func resourcePoolRefreshFunc(client *golangsdk.ServiceClient, id string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		pool, err := getPool(client, id)
		if err != nil {
			return nil, "", err
		}

		// The pool resource has no Status attribute, so a successful Get is the best we can do
		return pool, "ACTIVE", nil
	}
}

func getPool(client *golangsdk.ServiceClient, id string) (interface{}, error) {
	var (
		httpUrl = "v3/{project_id}/elb/pools/{pool_id}"
	)

	getPath := client.Endpoint + httpUrl
	getPath = strings.ReplaceAll(getPath, "{project_id}", client.ProjectID)
	getPath = strings.ReplaceAll(getPath, "{pool_id}", id)

	getOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	getResp, err := client.Request("GET", getPath, &getOpt)
	if err != nil {
		return nil, err
	}

	return utils.FlattenResponse(getResp)
}
