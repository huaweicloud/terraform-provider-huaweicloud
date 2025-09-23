package elb

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API ELB POST /v3/{project_id}/elb/master-slave-pools
// @API ELB GET /v3/{project_id}/elb/master-slave-pools/{pool_id}
// @API ELB DELETE /v3/{project_id}/elb/master-slave-pools/{pool_id}
func ResourceActiveStandbyPool() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceActiveStandbyPoolCreate,
		ReadContext:   resourceActiveStandbyPoolRead,
		DeleteContext: resourceActiveStandbyPoolDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
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
				ForceNew: true,
			},
			"members": {
				Type:     schema.TypeSet,
				Required: true,
				ForceNew: true,
				Elem:     activeStandbyMemberRefSchema(),
			},
			"healthmonitor": {
				Type:     schema.TypeList,
				Required: true,
				ForceNew: true,
				MaxItems: 1,
				Elem:     activeStandbyHealthMonitorRefSchema(),
			},
			"lb_algorithm": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Default:     "ROUND_ROBIN",
				Description: "schema: Required",
			},
			"name": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Computed: true,
			},
			"vpc_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Computed: true,
			},
			"loadbalancer_id": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				Computed:     true,
				AtLeastOneOf: []string{"loadbalancer_id", "listener_id", "type"},
			},
			"listener_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Computed: true,
			},
			"type": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Computed: true,
			},
			"any_port_enable": {
				Type:     schema.TypeBool,
				Optional: true,
				ForceNew: true,
				Computed: true,
			},
			"connection_drain_enabled": {
				Type:     schema.TypeBool,
				Optional: true,
				ForceNew: true,
				Computed: true,
			},
			"connection_drain_timeout": {
				Type:         schema.TypeInt,
				Optional:     true,
				ForceNew:     true,
				Computed:     true,
				RequiredWith: []string{"connection_drain_enabled"},
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Computed: true,
			},
			"ip_version": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"quic_cid_hash_strategy": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     activeStandbyQuicCidHashStrategyRefSchema(),
			},
			"created_at": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"updated_at": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func activeStandbyMemberRefSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"address": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"role": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"protocol_port": {
				Type:     schema.TypeInt,
				Optional: true,
				ForceNew: true,
				Computed: true,
			},
			"name": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Computed: true,
			},
			"subnet_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Computed: true,
			},
			"id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"member_type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"instance_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"operating_status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"ip_version": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"status": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     activeStandbyMemberStatusRefSchema(),
			},
			"reason": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     activeStandbyMemberReasonRefSchema(),
			},
		},
	}
}

func activeStandbyMemberStatusRefSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"listener_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"operating_status": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func activeStandbyMemberReasonRefSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"reason_code": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"expected_response": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"healthcheck_response": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func activeStandbyHealthMonitorRefSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"delay": {
				Type:     schema.TypeInt,
				Required: true,
				ForceNew: true,
			},
			"max_retries": {
				Type:     schema.TypeInt,
				Required: true,
				ForceNew: true,
			},
			"timeout": {
				Type:     schema.TypeInt,
				Required: true,
				ForceNew: true,
			},
			"type": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"domain_name": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Computed: true,
			},
			"expected_codes": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Computed: true,
			},
			"http_method": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Computed: true,
			},
			"max_retries_down": {
				Type:     schema.TypeInt,
				Optional: true,
				ForceNew: true,
				Computed: true,
			},
			"monitor_port": {
				Type:     schema.TypeInt,
				Optional: true,
				ForceNew: true,
				Computed: true,
			},
			"name": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Computed: true,
			},
			"url_path": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Computed: true,
			},
			"id": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func activeStandbyQuicCidHashStrategyRefSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"len": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"offset": {
				Type:     schema.TypeInt,
				Computed: true,
			},
		},
	}
}

func resourceActiveStandbyPoolCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conf := meta.(*config.Config)
	region := conf.GetRegion(d)

	var (
		createActiveStandbyPoolUrl     = "v3/{project_id}/elb/master-slave-pools"
		createActiveStandbyPoolProduct = "elb"
	)
	elbClient, err := conf.NewServiceClient(createActiveStandbyPoolProduct, region)
	if err != nil {
		return diag.Errorf("error creating ELB client: %s", err)
	}
	createActiveStandbyPoolPath := elbClient.Endpoint + createActiveStandbyPoolUrl
	createActiveStandbyPoolPath = strings.ReplaceAll(createActiveStandbyPoolPath, "{project_id}", elbClient.ProjectID)

	createActiveStandbyPoolOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}
	createActiveStandbyPoolOpt.JSONBody = utils.RemoveNil(buildCreateActiveStandbyPoolBodyParams(d))
	createActiveStandbyPoolResp, err := elbClient.Request("POST", createActiveStandbyPoolPath,
		&createActiveStandbyPoolOpt)
	if err != nil {
		return diag.Errorf("error creating ELB active standby pool: %s", err)
	}

	createActiveStandbyPoolRespBody, err := utils.FlattenResponse(createActiveStandbyPoolResp)
	if err != nil {
		return diag.Errorf("error retrieving ELB active standby pool: %s", err)
	}
	poolId := utils.PathSearch("pool.id", createActiveStandbyPoolRespBody, "").(string)
	if poolId == "" {
		return diag.Errorf("unable to find the ELB active standby pool ID from the API response")
	}

	d.SetId(poolId)

	timeout := d.Timeout(schema.TimeoutCreate)
	err = waitForElbActiveStandbyPool(ctx, elbClient, poolId, "ACTIVE", nil, timeout)
	if err != nil {
		return diag.FromErr(err)
	}

	return resourceActiveStandbyPoolRead(ctx, d, meta)
}

func resourceActiveStandbyPoolRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var mErr *multierror.Error

	var (
		getActiveStandbyPoolUrl     = "v3/{project_id}/elb/master-slave-pools/{pool_id}"
		getActiveStandbyPoolProduct = "elb"
	)
	elbClient, err := cfg.NewServiceClient(getActiveStandbyPoolProduct, region)
	if err != nil {
		return diag.Errorf("error creating ELB client: %s", err)
	}
	getActiveStandbyPoolPath := elbClient.Endpoint + getActiveStandbyPoolUrl
	getActiveStandbyPoolPath = strings.ReplaceAll(getActiveStandbyPoolPath, "{project_id}", elbClient.ProjectID)
	getActiveStandbyPoolPath = strings.ReplaceAll(getActiveStandbyPoolPath, "{pool_id}", d.Id())

	getActiveStandbyPoolOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	getActiveStandbyPoolResp, err := elbClient.Request("GET", getActiveStandbyPoolPath, &getActiveStandbyPoolOpt)

	if err != nil {
		return common.CheckDeletedDiag(d, err, "error retrieving ELB active standby pool")
	}

	getActiveStandbyPoolBody, err := utils.FlattenResponse(getActiveStandbyPoolResp)
	if err != nil {
		return diag.FromErr(err)
	}
	mErr = multierror.Append(
		mErr,
		d.Set("region", region),
		d.Set("protocol", utils.PathSearch("pool.protocol", getActiveStandbyPoolBody, nil)),
		d.Set("description", utils.PathSearch("pool.description", getActiveStandbyPoolBody, nil)),
		d.Set("name", utils.PathSearch("pool.name", getActiveStandbyPoolBody, nil)),
		d.Set("lb_algorithm", utils.PathSearch("pool.lb_algorithm", getActiveStandbyPoolBody, nil)),
		d.Set("loadbalancer_id", utils.PathSearch("pool.loadbalancers|[0].id", getActiveStandbyPoolBody, nil)),
		d.Set("listener_id", utils.PathSearch("pool.listeners|[0].id", getActiveStandbyPoolBody, nil)),
		d.Set("type", utils.PathSearch("pool.type", getActiveStandbyPoolBody, nil)),
		d.Set("any_port_enable", utils.PathSearch("pool.any_port_enable", getActiveStandbyPoolBody, nil)),
		d.Set("vpc_id", utils.PathSearch("pool.vpc_id", getActiveStandbyPoolBody, nil)),
		d.Set("members", flattenActiveStandbyPoolMembers(getActiveStandbyPoolBody)),
		d.Set("healthmonitor", flattenActiveStandbyPoolHealthMonitor(getActiveStandbyPoolBody)),
		d.Set("ip_version", utils.PathSearch("pool.ip_version", getActiveStandbyPoolBody, nil)),
		d.Set("connection_drain_enabled", utils.PathSearch("pool.connection_drain.enable", getActiveStandbyPoolBody, nil)),
		d.Set("connection_drain_timeout", utils.PathSearch("pool.connection_drain.timeout", getActiveStandbyPoolBody, nil)),
		d.Set("quic_cid_hash_strategy", flattenActiveStandbyQuicCidHashStrategy(getActiveStandbyPoolBody)),
		d.Set("created_at", utils.PathSearch("pool.created_at", getActiveStandbyPoolBody, nil)),
		d.Set("updated_at", utils.PathSearch("pool.updated_at", getActiveStandbyPoolBody, nil)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func resourceActiveStandbyPoolDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var (
		deleteActiveStandbyPoolUrl     = "v3/{project_id}/elb/master-slave-pools/{pool_id}"
		deleteActiveStandbyPoolProduct = "elb"
	)

	elbClient, err := cfg.NewServiceClient(deleteActiveStandbyPoolProduct, region)
	if err != nil {
		return diag.Errorf("error creating ELB client: %s", err)
	}
	deleteActiveStandbyPoolPath := elbClient.Endpoint + deleteActiveStandbyPoolUrl
	deleteActiveStandbyPoolPath = strings.ReplaceAll(deleteActiveStandbyPoolPath, "{project_id}", elbClient.ProjectID)
	deleteActiveStandbyPoolPath = strings.ReplaceAll(deleteActiveStandbyPoolPath, "{pool_id}", d.Id())

	deleteActiveStandbyPoolOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}
	_, err = elbClient.Request("DELETE", deleteActiveStandbyPoolPath, &deleteActiveStandbyPoolOpt)
	if err != nil {
		return diag.Errorf("error deleting ELB active standby pool: %s", err)
	}
	// Wait for Pool to delete
	timeout := d.Timeout(schema.TimeoutDelete)
	err = waitForElbActiveStandbyPool(ctx, elbClient, d.Id(), "DELETED", nil, timeout)
	if err != nil {
		return diag.FromErr(err)
	}
	return nil
}

func flattenActiveStandbyPoolMembers(resp interface{}) []interface{} {
	if resp == nil {
		return nil
	}
	curJson := utils.PathSearch("pool.members", resp, make([]interface{}, 0))
	if curJson == nil {
		return nil
	}
	curArray := curJson.([]interface{})
	rst := make([]interface{}, 0, len(curArray))
	for _, v := range curArray {
		rst = append(rst, map[string]interface{}{
			"address":          utils.PathSearch("address", v, nil),
			"protocol_port":    utils.PathSearch("protocol_port", v, nil),
			"role":             utils.PathSearch("role", v, nil),
			"id":               utils.PathSearch("id", v, nil),
			"name":             utils.PathSearch("name", v, nil),
			"subnet_id":        utils.PathSearch("subnet_cidr_id", v, nil),
			"member_type":      utils.PathSearch("member_type", v, nil),
			"instance_id":      utils.PathSearch("instance_id", v, nil),
			"operating_status": utils.PathSearch("operating_status", v, nil),
			"ip_version":       utils.PathSearch("ip_version", v, nil),
			"status":           flattenActiveStandbyPoolMemberStatus(v),
			"reason":           flattenActiveStandbyPoolMemberReason(v),
		})
	}
	return rst
}

func flattenActiveStandbyPoolMemberStatus(resp interface{}) []interface{} {
	if resp == nil {
		return nil
	}
	curJson := utils.PathSearch("status", resp, make([]interface{}, 0))
	if curJson == nil {
		return nil
	}
	curArray := curJson.([]interface{})
	rst := make([]interface{}, 0, len(curArray))
	for _, v := range curArray {
		rst = append(rst, map[string]interface{}{
			"listener_id":      utils.PathSearch("listener_id", v, nil),
			"operating_status": utils.PathSearch("operating_status", v, nil),
		})
	}
	return rst
}

func flattenActiveStandbyPoolMemberReason(resp interface{}) []interface{} {
	if resp == nil {
		return nil
	}

	curJson := utils.PathSearch("reason", resp, nil)
	if curJson == nil {
		return nil
	}

	rst := make([]interface{}, 0, 1)
	rst = append(rst, map[string]interface{}{
		"reason_code":          utils.PathSearch("reason_code", curJson, nil),
		"expected_response":    utils.PathSearch("expected_response", curJson, nil),
		"healthcheck_response": utils.PathSearch("healthcheck_response", curJson, nil),
	})
	return rst
}

func flattenActiveStandbyPoolHealthMonitor(resp interface{}) []interface{} {
	var rst []interface{}
	curJson := utils.PathSearch("pool.healthmonitor", resp, make(map[string]interface{})).(map[string]interface{})
	if len(curJson) < 1 {
		return nil
	}

	rst = []interface{}{
		map[string]interface{}{
			"name":             utils.PathSearch("name", curJson, nil),
			"monitor_port":     utils.PathSearch("monitor_port", curJson, nil),
			"domain_name":      utils.PathSearch("domain_name", curJson, nil),
			"type":             utils.PathSearch("type", curJson, nil),
			"url_path":         utils.PathSearch("url_path", curJson, nil),
			"max_retries_down": utils.PathSearch("max_retries_down", curJson, nil),
			"max_retries":      utils.PathSearch("max_retries", curJson, nil),
			"expected_codes":   utils.PathSearch("expected_codes", curJson, nil),
			"http_method":      utils.PathSearch("http_method", curJson, nil),
			"timeout":          utils.PathSearch("timeout", curJson, nil),
			"delay":            utils.PathSearch("delay", curJson, nil),
			"id":               utils.PathSearch("id", curJson, nil),
		},
	}
	return rst
}

func flattenActiveStandbyQuicCidHashStrategy(resp interface{}) []interface{} {
	curJson := utils.PathSearch("pool.quic_cid_hash_strategy", resp, nil)
	if curJson == nil {
		return nil
	}

	rst := []interface{}{
		map[string]interface{}{
			"len":    utils.PathSearch("len", curJson, nil),
			"offset": utils.PathSearch("offset", curJson, nil),
		},
	}
	return rst
}

func waitForElbActiveStandbyPool(ctx context.Context, elbClient *golangsdk.ServiceClient, id string, target string,
	pending []string, timeout time.Duration) error {
	log.Printf("[DEBUG] Waiting for pool %s to become %s.", id, target)

	stateConf := &resource.StateChangeConf{
		Target:     []string{target},
		Pending:    pending,
		Refresh:    resourceElbPoolRefreshFunc(elbClient, id),
		Timeout:    timeout,
		Delay:      5 * time.Second,
		MinTimeout: 3 * time.Second,
	}

	_, err := stateConf.WaitForStateContext(ctx)
	if err != nil {
		if _, ok := err.(golangsdk.ErrDefault404); ok {
			if target == "DELETED" {
				return nil
			}
			return fmt.Errorf("error: pool %s not found: %s", id, err)
		}
		return fmt.Errorf("error waiting for pool %s to become %s: %s", id, target, err)
	}
	return nil
}

func resourceElbPoolRefreshFunc(elbClient *golangsdk.ServiceClient, poolID string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		var (
			getActiveStandbyPoolUrl = "v3/{project_id}/elb/master-slave-pools/{pool_id}"
		)
		getActiveStandbyPoolPath := elbClient.Endpoint + getActiveStandbyPoolUrl
		getActiveStandbyPoolPath = strings.ReplaceAll(getActiveStandbyPoolPath, "{project_id}", elbClient.ProjectID)
		getActiveStandbyPoolPath = strings.ReplaceAll(getActiveStandbyPoolPath, "{pool_id}", poolID)

		getActiveStandbyPoolOpt := golangsdk.RequestOpts{
			KeepResponseBody: true,
		}

		getActiveStandbyPoolResp, err := elbClient.Request("GET", getActiveStandbyPoolPath, &getActiveStandbyPoolOpt)
		if err != nil {
			return nil, "", err
		}

		getActiveStandbyPoolBody, err := utils.FlattenResponse(getActiveStandbyPoolResp)
		if err != nil {
			return nil, "", err
		}

		// The pool resource has no Status attribute, so a successful Get is the best we can do
		return getActiveStandbyPoolBody, "ACTIVE", nil
	}
}

func buildCreateActiveStandbyPoolBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"lb_algorithm":     d.Get("lb_algorithm"),
		"protocol":         d.Get("protocol"),
		"name":             utils.ValueIgnoreEmpty(d.Get("name")),
		"loadbalancer_id":  utils.ValueIgnoreEmpty(d.Get("loadbalancer_id")),
		"listener_id":      utils.ValueIgnoreEmpty(d.Get("listener_id")),
		"type":             utils.ValueIgnoreEmpty(d.Get("type")),
		"any_port_enable":  utils.ValueIgnoreEmpty(d.Get("any_port_enable")),
		"vpc_id":           utils.ValueIgnoreEmpty(d.Get("vpc_id")),
		"description":      utils.ValueIgnoreEmpty(d.Get("description")),
		"ip_version":       utils.ValueIgnoreEmpty(d.Get("ip_version")),
		"members":          buildActiveStandbyPoolMembers(d.Get("members").(*schema.Set).List()),
		"healthmonitor":    buildActiveStandbyPoolHealthMonitor(d.Get("healthmonitor")),
		"connection_drain": buildActiveStandbyPoolConnectionDrain(d),
	}
	return map[string]interface{}{"pool": bodyParams}
}

func buildActiveStandbyPoolMembers(rawMembers []interface{}) []map[string]interface{} {
	if len(rawMembers) == 0 {
		return nil
	}
	members := make([]map[string]interface{}, 0, len(rawMembers))
	for _, member := range rawMembers {
		if v, ok := member.(map[string]interface{}); ok {
			members = append(members, map[string]interface{}{
				"address":        v["address"],
				"role":           v["role"],
				"protocol_port":  utils.ValueIgnoreEmpty(v["protocol_port"]),
				"name":           utils.ValueIgnoreEmpty(v["name"]),
				"subnet_cidr_id": utils.ValueIgnoreEmpty(v["subnet_id"]),
			})
		}
	}
	return members
}

func buildActiveStandbyPoolHealthMonitor(h interface{}) map[string]interface{} {
	if rawArray, ok := h.([]interface{}); ok {
		if len(rawArray) == 0 {
			return nil
		}
		raw, ok := rawArray[0].(map[string]interface{})
		if !ok {
			return nil
		}
		params := map[string]interface{}{
			"delay":            raw["delay"],
			"max_retries":      raw["max_retries"],
			"timeout":          raw["timeout"],
			"type":             raw["type"],
			"domain_name":      utils.ValueIgnoreEmpty(raw["domain_name"]),
			"expected_codes":   utils.ValueIgnoreEmpty(raw["expected_codes"]),
			"http_method":      utils.ValueIgnoreEmpty(raw["http_method"]),
			"max_retries_down": utils.ValueIgnoreEmpty(raw["max_retries_down"]),
			"monitor_port":     utils.ValueIgnoreEmpty(raw["monitor_port"]),
			"name":             utils.ValueIgnoreEmpty(raw["name"]),
			"url_path":         utils.ValueIgnoreEmpty(raw["url_path"]),
		}
		return params
	}
	return nil
}

func buildActiveStandbyPoolConnectionDrain(d *schema.ResourceData) map[string]interface{} {
	if v, ok := d.GetOk("connection_drain_enabled"); ok {
		params := map[string]interface{}{
			"enable":  v,
			"timeout": d.Get("connection_drain_timeout"),
		}
		return params
	}
	return nil
}
