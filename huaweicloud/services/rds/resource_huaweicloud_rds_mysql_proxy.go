package rds

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

var mysqlProxyNonUpdatableParams = []string{"instance_id", "flavor", "node_num", "proxy_name", "proxy_mode", "subnet_id"}

// @API RDS POST /v3/{project_id}/instances/{instance_id}/proxy/open
// @API RDS GET /v3/{project_id}/instances
// @API RDS GET /v3/{project_id}/jobs
// @API RDS GET /v3/{project_id}/instances/{instance_id}/proxy-list
// @API RDS POST /v3/{project_id}/instances/{instance_id}/proxy/{proxy_id}/route-mode
// @API RDS DELETE /v3/{project_id}/instances/{instance_id}/proxy/{proxy_id}
func ResourceMysqlProxy() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceMysqlProxyCreate,
		ReadContext:   resourceMysqlProxyRead,
		UpdateContext: resourceMysqlProxyUpdate,
		DeleteContext: resourceMysqlProxyDelete,

		Importer: &schema.ResourceImporter{
			StateContext: resourceMysqlMySQLProxyImportState,
		},

		CustomizeDiff: config.FlexibleForceNew(mysqlProxyNonUpdatableParams),

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(30 * time.Minute),
			Update: schema.DefaultTimeout(30 * time.Minute),
			Delete: schema.DefaultTimeout(10 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Computed: true,
			},
			"instance_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"flavor": {
				Type:     schema.TypeString,
				Required: true,
			},
			"node_num": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"proxy_name": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"proxy_mode": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"subnet_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"route_mode": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},
			"master_node_weight": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem:     mysqlProxyNodeWeightSchema(),
			},
			"readonly_nodes_weight": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem:     mysqlProxyNodeWeightSchema(),
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"address": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"port": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"delay_threshold_in_seconds": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"vcpus": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"memory": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"nodes": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     mysqlProxyNodeSchema(),
			},
			"mode": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"flavor_group_type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"transaction_split": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"connection_pool_type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"pay_mode": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"dns_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"seconds_level_monitor_fun_status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"alt_flag": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"force_read_only": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"ssl_option": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"support_balance_route_mode": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"support_proxy_ssl": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"support_switch_connection_pool_type": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"support_transaction_split": {
				Type:     schema.TypeBool,
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

func mysqlProxyNodeWeightSchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"weight": {
				Type:     schema.TypeInt,
				Required: true,
			},
		},
	}
	return &sc
}

func mysqlProxyNodeSchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"role": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"az_code": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"frozen_flag": {
				Type:     schema.TypeInt,
				Computed: true,
			},
		},
	}
	return &sc
}

func resourceMysqlProxyCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var (
		httpUrl = "v3/{project_id}/instances/{instance_id}/proxy/open"
		product = "rds"
	)
	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating RDS client: %s", err)
	}

	instanceId := d.Get("instance_id").(string)
	createPath := client.Endpoint + httpUrl
	createPath = strings.ReplaceAll(createPath, "{project_id}", client.ProjectID)
	createPath = strings.ReplaceAll(createPath, "{instance_id}", instanceId)

	createOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}
	createOpt.JSONBody = utils.RemoveNil(buildCreateMysqlProxyBodyParams(d))
	retryFunc := func() (interface{}, bool, error) {
		res, err := client.Request("POST", createPath, &createOpt)
		retry, err := handleMultiOperationsError(err)
		return res, retry, err
	}
	createResp, err := common.RetryContextWithWaitForState(&common.RetryContextWithWaitForStateParam{
		Ctx:          ctx,
		RetryFunc:    retryFunc,
		WaitFunc:     rdsInstanceStateRefreshFunc(client, instanceId),
		WaitTarget:   []string{"ACTIVE"},
		Timeout:      d.Timeout(schema.TimeoutCreate),
		DelayTimeout: 30 * time.Second,
		PollInterval: 10 * time.Second,
	})
	if err != nil {
		return diag.Errorf("error creating RDS Mysql database: %s", err)
	}

	createRespBody, err := utils.FlattenResponse(createResp.(*http.Response))
	if err != nil {
		return diag.FromErr(err)
	}

	jobId := utils.PathSearch("job_id", createRespBody, nil)
	if jobId == nil {
		return diag.Errorf("error creating RDS MySQL proxy: job_id is not found in API response")
	}

	proxyId := utils.PathSearch("proxy_id", createRespBody, nil)
	if proxyId == nil {
		return diag.Errorf("error creating RDS MySQL proxy: proxy_id is not found in API response")
	}
	d.SetId(proxyId.(string))

	err = checkRDSInstanceJobFinish(client, jobId.(string), d.Timeout(schema.TimeoutCreate))
	if err != nil {
		return diag.FromErr(err)
	}

	return resourceMysqlProxyRead(ctx, d, meta)
}

func buildCreateMysqlProxyBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"flavor_ref":        d.Get("flavor"),
		"node_num":          d.Get("node_num"),
		"proxy_name":        utils.ValueIgnoreEmpty(d.Get("proxy_name")),
		"proxy_mode":        utils.ValueIgnoreEmpty(d.Get("proxy_mode")),
		"route_mode":        utils.ValueIgnoreEmpty(d.Get("route_mode")),
		"nodes_read_weight": buildCreateMySQLProxyNodesReadWeightBody(d),
		"subnet_id":         utils.ValueIgnoreEmpty(d.Get("subnet_id")),
	}
	return bodyParams
}

func buildCreateMySQLProxyNodesReadWeightBody(d *schema.ResourceData) []map[string]interface{} {
	masterNodeRawParams := d.Get("master_node_weight").([]interface{})
	readonlyNodesRawParams := d.Get("readonly_nodes_weight").(*schema.Set).List()
	length := len(masterNodeRawParams) + len(readonlyNodesRawParams)
	if length == 0 {
		return nil
	}
	rst := make([]map[string]interface{}, 0, length)
	for _, v := range masterNodeRawParams {
		raw := v.(map[string]interface{})
		rst = append(rst, map[string]interface{}{
			"instance_id": raw["id"],
			"weight":      raw["weight"],
		})
	}
	for _, v := range readonlyNodesRawParams {
		raw := v.(map[string]interface{})
		rst = append(rst, map[string]interface{}{
			"instance_id": raw["id"],
			"weight":      raw["weight"],
		})
	}
	return rst
}

func resourceMysqlProxyRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var (
		httpUrl = "v3/{project_id}/instances/{instance_id}/proxy-list"
		product = "rds"
	)
	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating RDS client: %s", err)
	}
	getPath := client.Endpoint + httpUrl
	getPath = strings.ReplaceAll(getPath, "{project_id}", client.ProjectID)
	getPath = strings.ReplaceAll(getPath, "{instance_id}", d.Get("instance_id").(string))
	getOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"Content-Type": "application/json"},
	}
	getResp, err := client.Request("GET", getPath, &getOpt)
	if err != nil {
		return common.CheckDeletedDiag(d, err, "error retrieving RDS MySQL proxy")
	}

	getRespBody, err := utils.FlattenResponse(getResp)
	if err != nil {
		return diag.FromErr(err)
	}

	searchExpression := fmt.Sprintf("proxy_query_info_list[?proxy.pool_id=='%s']|[0]", d.Id())
	proxy := utils.PathSearch(searchExpression, getRespBody, nil)
	if proxy == nil {
		return common.CheckDeletedDiag(d, golangsdk.ErrDefault404{}, "error retrieving RDS MySQL proxy")
	}

	mErr := multierror.Append(
		d.Set("region", region),
		d.Set("instance_id", utils.PathSearch("master_instance.id", proxy, nil)),
		d.Set("flavor", utils.PathSearch("proxy.flavor_info.code", proxy, nil)),
		d.Set("node_num", utils.PathSearch("proxy.node_num", proxy, nil)),
		d.Set("proxy_name", utils.PathSearch("proxy.name", proxy, nil)),
		d.Set("proxy_mode", utils.PathSearch("proxy.proxy_mode", proxy, nil)),
		d.Set("subnet_id", utils.PathSearch("proxy.subnet_id", proxy, nil)),
		d.Set("route_mode", utils.PathSearch("proxy.route_mode", proxy, nil)),
		d.Set("master_node_weight", flattenMysqlProxyResponseBodyMasterNodeWeight(proxy, d)),
		d.Set("readonly_nodes_weight", flattenMysqlProxyResponseBodyReadonlyNodesWeight(proxy, d)),
		d.Set("status", utils.PathSearch("proxy.status", proxy, nil)),
		d.Set("address", utils.PathSearch("proxy.address", proxy, nil)),
		d.Set("port", utils.PathSearch("proxy.port", proxy, nil)),
		d.Set("delay_threshold_in_seconds",
			utils.PathSearch("proxy.delay_threshold_in_seconds", proxy, nil)),
		d.Set("vcpus", utils.PathSearch("proxy.cpu", proxy, nil)),
		d.Set("memory", utils.PathSearch("proxy.mem", proxy, nil)),
		d.Set("nodes", flattenMysqlProxyResponseBodyNodes(proxy)),
		d.Set("mode", utils.PathSearch("proxy.mode", proxy, nil)),
		d.Set("flavor_group_type", utils.PathSearch("proxy.flavor_info.group_type", proxy, nil)),
		d.Set("transaction_split", utils.PathSearch("proxy.transaction_split", proxy, nil)),
		d.Set("connection_pool_type", utils.PathSearch("proxy.connection_pool_type", proxy, nil)),
		d.Set("pay_mode", utils.PathSearch("proxy.pay_mode", proxy, nil)),
		d.Set("dns_name", utils.PathSearch("proxy.dns_name", proxy, nil)),
		d.Set("seconds_level_monitor_fun_status",
			utils.PathSearch("proxy.seconds_level_monitor_fun_status", proxy, nil)),
		d.Set("alt_flag", utils.PathSearch("proxy.alt_flag", proxy, nil)),
		d.Set("force_read_only", utils.PathSearch("proxy.force_read_only", proxy, nil)),
		d.Set("ssl_option", utils.PathSearch("proxy.ssl_option", proxy, nil)),
		d.Set("support_balance_route_mode",
			utils.PathSearch("proxy.support_balance_route_mode", proxy, nil)),
		d.Set("support_proxy_ssl",
			utils.PathSearch("proxy.support_balance_route_mode", proxy, nil)),
		d.Set("support_switch_connection_pool_type",
			utils.PathSearch("proxy.support_switch_connection_pool_type", proxy, nil)),
		d.Set("support_transaction_split",
			utils.PathSearch("proxy.support_transaction_split", proxy, nil)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenMysqlProxyResponseBodyMasterNodeWeight(proxy interface{}, d *schema.ResourceData) []interface{} {
	masterNodeRawParams := d.Get("master_node_weight").([]interface{})
	if len(masterNodeRawParams) < 1 {
		return nil
	}
	rst := make([]interface{}, 0, 1)
	rst = append(rst, map[string]interface{}{
		"id":     utils.PathSearch("master_instance.id", proxy, nil),
		"weight": utils.PathSearch("master_instance.weight", proxy, nil),
	})
	return rst
}

func flattenMysqlProxyResponseBodyReadonlyNodesWeight(proxy interface{}, d *schema.ResourceData) []interface{} {
	readonlyNodesJson := utils.PathSearch("readonly_instances", proxy, make([]interface{}, 0))
	readonlyNodeArray := readonlyNodesJson.([]interface{})
	if len(readonlyNodeArray) < 1 {
		return nil
	}

	readonlyNodesWeightRaw := d.Get("readonly_nodes_weight").(*schema.Set).List()
	readonlyNodesRawMap := make(map[string]bool)
	for _, v := range readonlyNodesWeightRaw {
		readonlyNodesRawMap[v.(map[string]interface{})["id"].(string)] = true
	}
	rst := make([]interface{}, 0, len(readonlyNodesWeightRaw))
	for _, v := range readonlyNodeArray {
		id := utils.PathSearch("id", v, "").(string)
		if readonlyNodesRawMap[id] {
			rst = append(rst, map[string]interface{}{
				"id":     id,
				"weight": utils.PathSearch("weight", v, nil),
			})
		}
	}
	return rst
}

func flattenMysqlProxyResponseBodyNodes(proxy interface{}) []interface{} {
	nodesJson := utils.PathSearch("proxy.nodes", proxy, make([]interface{}, 0))
	nodeArray := nodesJson.([]interface{})
	if len(nodeArray) < 1 {
		return nil
	}
	rst := make([]interface{}, 0, len(nodeArray))
	for _, v := range nodeArray {
		rst = append(rst, map[string]interface{}{
			"id":          utils.PathSearch("id", v, nil),
			"status":      utils.PathSearch("status", v, nil),
			"role":        utils.PathSearch("role", v, nil),
			"az_code":     utils.PathSearch("az_code", v, nil),
			"frozen_flag": utils.PathSearch("frozen_flag", v, nil),
		})
	}
	return rst
}

func resourceMysqlProxyUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var (
		product = "rds"
	)
	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating RDS client: %s", err)
	}

	if d.HasChanges("route_mode", "master_node_weight", "readonly_nodes_weight") {
		err = updateMysqlRouteMode(d, client)
		if err != nil {
			return diag.FromErr(err)
		}
	}

	return resourceMysqlProxyRead(ctx, d, meta)
}

func updateMysqlRouteMode(d *schema.ResourceData, client *golangsdk.ServiceClient) error {
	var (
		httpUrl = "v3/{project_id}/instances/{instance_id}/proxy/{proxy_id}/route-mode"
	)
	updatePath := client.Endpoint + httpUrl
	updatePath = strings.ReplaceAll(updatePath, "{project_id}", client.ProjectID)
	updatePath = strings.ReplaceAll(updatePath, "{instance_id}", d.Get("instance_id").(string))
	updatePath = strings.ReplaceAll(updatePath, "{proxy_id}", d.Id())

	updateOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}
	updateOpt.JSONBody = buildUpdateMysqlRouteModeBodyParams(d)

	updateResp, err := client.Request("POST", updatePath, &updateOpt)
	if err != nil {
		return fmt.Errorf("error updating RDS MySQL route mode: %s", err)
	}

	_, err = utils.FlattenResponse(updateResp)
	if err != nil {
		return err
	}

	return nil
}

func buildUpdateMysqlRouteModeBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := make(map[string]interface{})
	bodyParams["route_mode"] = d.Get("route_mode")
	if len(d.Get("master_node_weight").([]interface{})) > 0 {
		masterNodeWeight := d.Get("master_node_weight").([]interface{})[0].(map[string]interface{})["weight"]
		bodyParams["master_weight"] = masterNodeWeight
	}
	if len(d.Get("readonly_nodes_weight").(*schema.Set).List()) > 0 {
		readonlyNodesRawParams := d.Get("readonly_nodes_weight").(*schema.Set).List()
		readonlyNodesWeightParam := make([]map[string]interface{}, len(readonlyNodesRawParams))
		for i, v := range readonlyNodesRawParams {
			raw := v.(map[string]interface{})
			readonlyNodesWeightParam[i] = map[string]interface{}{
				"instance_id": raw["id"],
				"weight":      raw["weight"],
			}
		}
		bodyParams["readonly_instances"] = readonlyNodesWeightParam
	}
	return bodyParams
}

func resourceMysqlProxyDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var (
		httpUrl = "v3/{project_id}/instances/{instance_id}/proxy/{proxy_id}"
		product = "rds"
	)
	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating RDS client: %s", err)
	}

	deletePath := client.Endpoint + httpUrl
	deletePath = strings.ReplaceAll(deletePath, "{project_id}", client.ProjectID)
	deletePath = strings.ReplaceAll(deletePath, "{instance_id}", d.Get("instance_id").(string))
	deletePath = strings.ReplaceAll(deletePath, "{proxy_id}", d.Id())

	deleteOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"Content-Type": "application/json"},
	}

	deleteResp, err := client.Request("DELETE", deletePath, &deleteOpt)
	if err != nil {
		return common.CheckDeletedDiag(d, common.ConvertExpected400ErrInto404Err(err, "error_code", "DBS.200939"),
			"error deleting RDS MySQL proxy")
	}

	deleteRespBody, err := utils.FlattenResponse(deleteResp)
	if err != nil {
		return diag.FromErr(err)
	}

	jobId := utils.PathSearch("job_id", deleteRespBody, nil)
	if jobId == nil {
		return diag.Errorf("error deleting RDS MySQL proxy: job_id is not found in API response")
	}

	err = checkRDSInstanceJobFinish(client, jobId.(string), d.Timeout(schema.TimeoutDelete))
	if err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func resourceMysqlMySQLProxyImportState(_ context.Context, d *schema.ResourceData, _ interface{}) ([]*schema.ResourceData,
	error) {
	parts := strings.Split(d.Id(), "/")

	if len(parts) != 2 {
		return nil, fmt.Errorf("invalid format specified for import ID, must be <instance_id>/<id>")
	}

	d.SetId(parts[1])
	mErr := multierror.Append(nil,
		d.Set("instance_id", parts[0]),
	)

	return []*schema.ResourceData{d}, mErr.ErrorOrNil()
}
