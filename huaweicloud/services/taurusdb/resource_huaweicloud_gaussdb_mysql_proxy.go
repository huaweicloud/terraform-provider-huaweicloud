package taurusdb

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"
	"github.com/chnsz/golangsdk/pagination"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API GaussDBforMySQL POST /v3/{project_id}/instances/{instance_id}/proxy
// @API GaussDBforMySQL GET /v3/{project_id}/jobs
// @API GaussDBforMySQL PUT /v3/{project_id}/instances/{instance_id}/proxy/{proxy_id}/port
// @API GaussDBforMySQL PUT /v3/{project_id}/instances/{instance_id}/proxy/{proxy_id}/configurations
// @API GaussDBforMySQL PUT /v3/{project_id}/instances/{instance_id}/proxy/{proxy_id}/flavor
// @API GaussDBforMySQL POST /v3/{project_id}/instances/{instance_id}/proxy/enlarge
// @API GaussDBforMySQL PUT /v3/{project_id}/instances/{instance_id}/proxy/{proxy_id}/reduce
// @API GaussDBforMySQL PUT /v3/{project_id}/instances/{instance_id}/proxy/{proxy_id}/rename
// @API GaussDBforMySQL PUT /v3/{project_id}/instances/{instance_id}/proxy/{proxy_id}/weight
// @API GaussDBforMySQL POST /v3/{project_id}/instances/{instance_id}/proxy/{proxy_id}/new-node-auto-add
// @API GaussDBforMySQL POST /v3/{project_id}/instances/{instance_id}/proxy/transaction-split
// @API GaussDBforMySQL PUT /v3/{project_id}/instances/{instance_id}/proxy/{proxy_id}/session-consistence
// @API GaussDBforMySQL PUT /v3/{project_id}/instances/{instance_id}/proxy/{proxy_id}/connection-pool-type
// @API GaussDBforMySQL POST /v3/{project_id}/instances/{instance_id}/proxy/{proxy_id}/access-control-switch
// @API GaussDBforMySQL POST /v3/{project_id}/instances/{instance_id}/proxy/{proxy_id}/access-control
// @API GaussDBforMySQL GET /v3/{project_id}/instances/{instance_id}/proxy/{proxy_id}/ipgroup
// @API GaussDBforMySQL GET /v3/{project_id}/instances/{instance_id}/proxies
// @API GaussDBforMySQL GET /v3/{project_id}/instances/{instance_id}/proxy/{proxy_id}/{engine_name}/proxy-version
// @API GaussDBforMySQL GET /v3/{project_id}/instances/{instance_id}/proxy/{proxy_id}/configurations
// @API GaussDBforMySQL DELETE /v3/{project_id}/instances/{instance_id}/proxy
func ResourceGaussDBProxy() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceGaussDBProxyCreate,
		ReadContext:   resourceGaussDBProxyRead,
		UpdateContext: resourceGaussDBProxyUpdate,
		DeleteContext: resourceGaussDBProxyDelete,
		Importer: &schema.ResourceImporter{
			StateContext: resourceGaussDBMySQLProxyImportState,
		},

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
				ForceNew: true,
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
				ForceNew: true,
			},
			"route_mode": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"subnet_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"master_node_weight": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				MaxItems: 1,
				Elem:     gaussDBMysqlProxyNodeWeightSchema(),
			},
			"readonly_nodes_weight": {
				Type:     schema.TypeSet,
				Optional: true,
				Computed: true,
				Elem:     gaussDBMysqlProxyNodeWeightSchema(),
			},
			"new_node_auto_add_status": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"new_node_weight": {
				Type:         schema.TypeInt,
				Optional:     true,
				RequiredWith: []string{"new_node_auto_add_status"},
			},
			"port": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},
			"parameters": {
				Type: schema.TypeSet,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:     schema.TypeString,
							Required: true,
						},
						"value": {
							Type:     schema.TypeString,
							Required: true,
						},
						"elem_type": {
							Type:     schema.TypeString,
							Required: true,
						},
					},
				},
				Optional: true,
				Computed: true,
			},
			"transaction_split": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"consistence_mode": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"connection_pool_type": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"open_access_control": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},
			"access_control_type": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"access_control_ip_list": {
				Type:         schema.TypeSet,
				Optional:     true,
				Computed:     true,
				Elem:         gaussDBMysqlProxyAccessControlIpListSchema(),
				RequiredWith: []string{"access_control_type"},
			},
			"address": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"switch_connection_pool_type_enabled": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"nodes": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     gaussDBMysqlProxyNodeSchema(),
			},
			"current_version": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"can_upgrade": {
				Type:     schema.TypeBool,
				Computed: true,
			},
		},
	}
}

func gaussDBMysqlProxyNodeWeightSchema() *schema.Resource {
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

func gaussDBMysqlProxyAccessControlIpListSchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"ip": {
				Type:     schema.TypeString,
				Required: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
		},
	}
	return &sc
}

func gaussDBMysqlProxyNodeSchema() *schema.Resource {
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
			"name": {
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

func resourceGaussDBProxyCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var (
		httpUrl = "v3/{project_id}/instances/{instance_id}/proxy"
		product = "gaussdb"
	)
	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating GaussDB client: %s", err)
	}

	createPath := client.Endpoint + httpUrl
	createPath = strings.ReplaceAll(createPath, "{project_id}", client.ProjectID)
	createPath = strings.ReplaceAll(createPath, "{instance_id}", d.Get("instance_id").(string))

	createOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}
	createOpt.JSONBody = utils.RemoveNil(buildCreateGaussDBMysqlProxyBodyParams(d))
	createResp, err := client.Request("POST", createPath, &createOpt)
	if err != nil {
		return diag.Errorf("error creating GaussDB MySQL proxy: %s", err)
	}

	createRespBody, err := utils.FlattenResponse(createResp)
	if err != nil {
		return diag.FromErr(err)
	}

	jobId := utils.PathSearch("job_id", createRespBody, nil)
	if jobId == nil {
		return diag.Errorf("error creating GaussDB MySQL proxy: job_id is not found in API response")
	}

	searchExpression := "proxy_list[?proxy.status=='ENABLING PROXY']|[0].proxy.pool_id"
	proxyId, err := getGaussDBProxy(client, d.Get("instance_id").(string), searchExpression)
	if err != nil {
		return diag.Errorf("error retrieving GaussDB MySQL proxy: %s", err)
	}
	d.SetId(proxyId.(string))

	err = checkGaussDBMySQLProxyJobFinish(ctx, client, jobId.(string), d.Timeout(schema.TimeoutCreate))
	if err != nil {
		return diag.FromErr(err)
	}

	if _, ok := d.GetOk("port"); ok {
		err = updateGaussDBMySQLPort(ctx, d, client)
		if err != nil {
			return diag.FromErr(err)
		}
	}

	if _, ok := d.GetOk("parameters"); ok {
		err = updateGaussDBMySQLParameters(d, client)
		if err != nil {
			return diag.FromErr(err)
		}
	}

	if transactionSplit, ok := d.GetOk("transaction_split"); ok && transactionSplit.(string) == "ON" {
		err = updateGaussDBMySQLTransactionSplit(ctx, d, client)
		if err != nil {
			return diag.FromErr(err)
		}
	}

	if consistenceMode, ok := d.GetOk("consistence_mode"); ok && consistenceMode != "eventual" {
		err = updateGaussDBMySQLConsistenceMode(ctx, d, client)
		if err != nil {
			return diag.FromErr(err)
		}
	}

	if connectionPoolType, ok := d.GetOk("connection_pool_type"); ok && connectionPoolType == "SESSION" {
		err = updateGaussDBMySQLConnectionPoolType(ctx, d, client)
		if err != nil {
			return diag.FromErr(err)
		}
	}

	if _, ok := d.GetOk("open_access_control"); ok {
		err = updateGaussDBMySQLOpenAccessControl(d, client)
		if err != nil {
			return diag.FromErr(err)
		}
	}

	if _, ok := d.GetOk("access_control_type"); ok {
		err = updateGaussDBMySQLAccessControlRule(d, client)
		if err != nil {
			return diag.FromErr(err)
		}
	}

	return resourceGaussDBProxyRead(ctx, d, meta)
}

func buildCreateGaussDBMysqlProxyBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"flavor_ref":               d.Get("flavor"),
		"node_num":                 d.Get("node_num"),
		"proxy_name":               utils.ValueIgnoreEmpty(d.Get("proxy_name")),
		"proxy_mode":               utils.ValueIgnoreEmpty(d.Get("proxy_mode")),
		"route_mode":               utils.ValueIgnoreEmpty(d.Get("route_mode")),
		"nodes_read_weight":        buildCreateGaussDBMySQLProxyNodesReadWeightBody(d),
		"subnet_id":                utils.ValueIgnoreEmpty(d.Get("subnet_id")),
		"new_node_auto_add_status": utils.ValueIgnoreEmpty(d.Get("new_node_auto_add_status")),
		"new_node_weight":          utils.ValueIgnoreEmpty(d.Get("new_node_weight")),
	}
	return bodyParams
}

func buildCreateGaussDBMySQLProxyNodesReadWeightBody(d *schema.ResourceData) []map[string]interface{} {
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
			"id":     raw["id"],
			"weight": raw["weight"],
		})
	}
	for _, v := range readonlyNodesRawParams {
		raw := v.(map[string]interface{})
		rst = append(rst, map[string]interface{}{
			"id":     raw["id"],
			"weight": raw["weight"],
		})
	}
	return rst
}

func resourceGaussDBProxyRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var (
		product = "gaussdb"
	)
	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating GaussDB client: %s", err)
	}
	searchExpression := fmt.Sprintf("proxy_list[?proxy.pool_id=='%s']|[0]", d.Id())
	proxy, err := getGaussDBProxy(client, d.Get("instance_id").(string), searchExpression)
	if err != nil {
		return common.CheckDeletedDiag(d, parseMysqlProxyError(err), "error retrieving GaussDB MySQL proxy")
	}

	mErr := multierror.Append(
		d.Set("region", region),
		d.Set("flavor", utils.PathSearch("proxy.flavor_ref", proxy, nil)),
		d.Set("node_num", utils.PathSearch("proxy.node_num", proxy, nil)),
		d.Set("proxy_name", utils.PathSearch("proxy.name", proxy, nil)),
		d.Set("route_mode", utils.PathSearch("proxy.route_mode", proxy, nil)),
		d.Set("subnet_id", utils.PathSearch("proxy.subnet_id", proxy, nil)),
		d.Set("master_node_weight", flattenGaussDBProxyResponseBodyMasterNodeWeight(proxy)),
		d.Set("readonly_nodes_weight", flattenGaussDBProxyResponseBodyReadonlyNodesWeight(proxy, d)),
		d.Set("new_node_auto_add_status", utils.PathSearch("proxy.new_node_auto_add_status", proxy, nil)),
		d.Set("port", utils.PathSearch("proxy.port", proxy, nil)),
		d.Set("consistence_mode", utils.PathSearch("proxy.consistence_mode", proxy, nil)),
		d.Set("connection_pool_type", utils.PathSearch("proxy.connection_pool_type", proxy, nil)),
		d.Set("switch_connection_pool_type_enabled", utils.PathSearch("proxy.switch_connection_pool_type_enabled",
			proxy, nil)),
		d.Set("address", utils.PathSearch("proxy.address", proxy, nil)),
		d.Set("nodes", flattenGaussDBProxyResponseBodyNodes(proxy)),
	)
	transactionSplit := utils.PathSearch("proxy.transaction_split", proxy, "").(string)
	if transactionSplit == "true" {
		mErr = multierror.Append(mErr, d.Set("transaction_split", "ON"))
	} else {
		mErr = multierror.Append(mErr, d.Set("transaction_split", "OFF"))
	}

	parameters, err := getGaussDBProxyParameters(d, client)
	if err != nil {
		log.Printf("[WARN] fetching GaussDB MySQL proxy paremeters failed: %s", err)
	} else {
		mErr = multierror.Append(d.Set("parameters", parameters))
	}

	version, err := getGaussDBProxyVersion(d, client)
	if err != nil {
		log.Printf("[WARN] fetching GaussDB MySQL proxy version failed: %s", err)
	} else {
		mErr = multierror.Append(
			d.Set("current_version", utils.PathSearch("current_version", version, nil)),
			d.Set("can_upgrade", utils.PathSearch("can_upgrade", version, nil)),
		)
	}

	accessControl, err := getGaussDBProxyAccessControl(d, client)
	if err != nil {
		log.Printf("[WARN] fetching GaussDB MySQL proxy access control failed: %s", err)
	} else {
		openAccessControl := utils.PathSearch("enable_ip_group", accessControl, nil)
		accessControlType := utils.PathSearch("type", accessControl, nil)
		accessControlIpList := flattenGaussDBProxyAccessControlIpList(accessControl)
		mErr = multierror.Append(
			d.Set("open_access_control", openAccessControl),
			d.Set("access_control_type", accessControlType),
			d.Set("access_control_ip_list", accessControlIpList),
		)
	}

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenGaussDBProxyResponseBodyMasterNodeWeight(proxy interface{}) []interface{} {
	rst := make([]interface{}, 0, 1)
	rst = append(rst, map[string]interface{}{
		"id":     utils.PathSearch("master_node.id", proxy, nil),
		"weight": utils.PathSearch("master_node.weight", proxy, nil),
	})
	return rst
}

func flattenGaussDBProxyResponseBodyReadonlyNodesWeight(proxy interface{}, d *schema.ResourceData) []interface{} {
	readonlyNodesJson := utils.PathSearch("readonly_nodes", proxy, make([]interface{}, 0))
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

func flattenGaussDBProxyResponseBodyNodes(proxy interface{}) []interface{} {
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
			"name":        utils.PathSearch("name", v, nil),
			"role":        utils.PathSearch("role", v, nil),
			"az_code":     utils.PathSearch("az_code", v, nil),
			"frozen_flag": utils.PathSearch("frozen_flag", v, nil),
		})
	}
	return rst
}

func getGaussDBProxy(client *golangsdk.ServiceClient, instanceId, expression string) (interface{}, error) {
	var (
		httpUrl = "v3/{project_id}/instances/{instance_id}/proxies"
	)

	listPath := client.Endpoint + httpUrl
	listPath = strings.ReplaceAll(listPath, "{project_id}", client.ProjectID)
	listPath = strings.ReplaceAll(listPath, "{instance_id}", instanceId)

	listMysqlDatabasesResp, err := pagination.ListAllItems(
		client,
		"offset",
		listPath,
		&pagination.QueryOpts{MarkerField: ""})

	if err != nil {
		return nil, err
	}

	listRespJson, err := json.Marshal(listMysqlDatabasesResp)
	if err != nil {
		return nil, err
	}
	var listRespBody interface{}
	err = json.Unmarshal(listRespJson, &listRespBody)
	if err != nil {
		return nil, err
	}
	proxy := utils.PathSearch(expression, listRespBody, nil)
	if proxy == nil {
		return nil, golangsdk.ErrDefault404{}
	}
	return proxy, nil
}

func getGaussDBProxyParameters(d *schema.ResourceData, client *golangsdk.ServiceClient) (interface{}, error) {
	var (
		httpUrl = "v3/{project_id}/instances/{instance_id}/proxy/{proxy_id}/configurations"
	)

	listPath := client.Endpoint + httpUrl
	listPath = strings.ReplaceAll(listPath, "{project_id}", client.ProjectID)
	listPath = strings.ReplaceAll(listPath, "{instance_id}", d.Get("instance_id").(string))
	listPath = strings.ReplaceAll(listPath, "{proxy_id}", d.Id())

	listMysqlDatabasesResp, err := pagination.ListAllItems(
		client,
		"offset",
		listPath,
		&pagination.QueryOpts{MarkerField: ""})

	if err != nil {
		return nil, err
	}

	listRespJson, err := json.Marshal(listMysqlDatabasesResp)
	if err != nil {
		return nil, err
	}
	var listRespBody interface{}
	err = json.Unmarshal(listRespJson, &listRespBody)
	if err != nil {
		return nil, err
	}
	return flattenGaussDBProxyParameters(listRespBody, d), nil
}

func getGaussDBProxyVersion(d *schema.ResourceData, client *golangsdk.ServiceClient) (interface{}, error) {
	var (
		httpUrl = "v3/{project_id}/instances/{instance_id}/proxy/{proxy_id}/{engine_name}/proxy-version"
	)

	getPath := client.Endpoint + httpUrl
	getPath = strings.ReplaceAll(getPath, "{project_id}", client.ProjectID)
	getPath = strings.ReplaceAll(getPath, "{instance_id}", d.Get("instance_id").(string))
	getPath = strings.ReplaceAll(getPath, "{proxy_id}", d.Id())
	getPath = strings.ReplaceAll(getPath, "{engine_name}", "taurusproxy")

	getOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"Content-Type": "application/json"},
	}
	getResp, err := client.Request("GET", getPath, &getOpt)
	if err != nil {
		return nil, err
	}

	return utils.FlattenResponse(getResp)
}

func flattenGaussDBProxyParameters(resp interface{}, d *schema.ResourceData) []interface{} {
	if resp == nil {
		return nil
	}
	parameters := d.Get("parameters").(*schema.Set).List()
	if len(parameters) == 0 {
		return nil
	}
	parametersMap := make(map[string]bool)
	for _, v := range parameters {
		name := v.(map[string]interface{})["name"].(string)
		parametersMap[name] = true
	}

	curJson := utils.PathSearch("configurations", resp, make([]interface{}, 0))
	curArray := curJson.([]interface{})
	rst := make([]interface{}, 0, len(curArray))
	for _, v := range curArray {
		name := utils.PathSearch("name", v, "").(string)
		if parametersMap[name] {
			rst = append(rst, map[string]interface{}{
				"name":      utils.PathSearch("name", v, nil),
				"value":     utils.PathSearch("value", v, nil),
				"elem_type": utils.PathSearch("elem_type", v, nil),
			})
		}
	}
	return rst
}

func getGaussDBProxyAccessControl(d *schema.ResourceData, client *golangsdk.ServiceClient) (interface{}, error) {
	var (
		httpUrl = "v3/{project_id}/instances/{instance_id}/proxy/{proxy_id}/ipgroup"
	)

	getPath := client.Endpoint + httpUrl
	getPath = strings.ReplaceAll(getPath, "{project_id}", client.ProjectID)
	getPath = strings.ReplaceAll(getPath, "{instance_id}", d.Get("instance_id").(string))
	getPath = strings.ReplaceAll(getPath, "{proxy_id}", d.Id())

	getOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"Content-Type": "application/json"},
	}

	getResp, err := client.Request("GET", getPath, &getOpt)
	if err != nil {
		return nil, fmt.Errorf("error retrieving GaussDB MySQL access control: %s", err)
	}

	return utils.FlattenResponse(getResp)
}

func flattenGaussDBProxyAccessControlIpList(resp interface{}) []interface{} {
	if resp == nil {
		return nil
	}

	curJson := utils.PathSearch("ip_group.ip_list", resp, make([]interface{}, 0))
	curArray := curJson.([]interface{})
	rst := make([]interface{}, 0, len(curArray))
	for _, v := range curArray {
		rst = append(rst, map[string]interface{}{
			"ip":          utils.PathSearch("ip", v, nil),
			"description": utils.PathSearch("description", v, nil),
		})
	}
	return rst
}

func resourceGaussDBProxyUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	client, err := cfg.GaussdbV3Client(cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating GaussDB client: %s ", err)
	}

	if d.HasChange("flavor") {
		err = updateGaussDBMySQLFlavor(ctx, d, client)
		if err != nil {
			return diag.FromErr(err)
		}
	}

	if d.HasChange("node_num") {
		oldNum, newNum := d.GetChange("node_num")
		if oldNum.(int) < newNum.(int) {
			err = enlargeGaussDBMySQLProxyNumber(ctx, d, client, newNum.(int)-oldNum.(int))
			if err != nil {
				return diag.FromErr(err)
			}
		} else {
			err = reduceGaussDBMySQLProxyNumber(ctx, d, client, oldNum.(int)-newNum.(int))
			if err != nil {
				return diag.FromErr(err)
			}
		}
	}

	if d.HasChange("proxy_name") {
		err = updateGaussDBMySQLProxyName(d, client)
		if err != nil {
			return diag.FromErr(err)
		}
	}

	if d.HasChange("port") {
		err = updateGaussDBMySQLPort(ctx, d, client)
		if err != nil {
			return diag.FromErr(err)
		}
	}

	if d.HasChanges("master_node_weight", "readonly_nodes_weight") {
		err = updateGaussDBMySQLNodesWeight(ctx, d, client)
		if err != nil {
			return diag.FromErr(err)
		}
	}

	if d.HasChanges("new_node_auto_add_status", "new_node_weight") {
		err = updateGaussDBMySQLNewNode(d, client)
		if err != nil {
			return diag.FromErr(err)
		}
	}

	if d.HasChange("parameters") {
		err = updateGaussDBMySQLParameters(d, client)
		if err != nil {
			return diag.FromErr(err)
		}
	}

	if d.HasChange("transaction_split") {
		err = updateGaussDBMySQLTransactionSplit(ctx, d, client)
		if err != nil {
			return diag.FromErr(err)
		}
	}

	if d.HasChange("consistence_mode") {
		err = updateGaussDBMySQLConsistenceMode(ctx, d, client)
		if err != nil {
			return diag.FromErr(err)
		}
	}

	if d.HasChange("connection_pool_type") {
		err = updateGaussDBMySQLConnectionPoolType(ctx, d, client)
		if err != nil {
			return diag.FromErr(err)
		}
	}

	if d.HasChange("open_access_control") {
		err = updateGaussDBMySQLOpenAccessControl(d, client)
		if err != nil {
			return diag.FromErr(err)
		}
	}

	if d.HasChanges("access_control_type", "access_control_ip_list") {
		err = updateGaussDBMySQLAccessControlRule(d, client)
		if err != nil {
			return diag.FromErr(err)
		}
	}

	return resourceGaussDBProxyRead(ctx, d, meta)
}

func updateGaussDBMySQLFlavor(ctx context.Context, d *schema.ResourceData, client *golangsdk.ServiceClient) error {
	var (
		httpUrl = "v3/{project_id}/instances/{instance_id}/proxy/{proxy_id}/flavor"
	)
	updatePath := client.Endpoint + httpUrl
	updatePath = strings.ReplaceAll(updatePath, "{project_id}", client.ProjectID)
	updatePath = strings.ReplaceAll(updatePath, "{instance_id}", d.Get("instance_id").(string))
	updatePath = strings.ReplaceAll(updatePath, "{proxy_id}", d.Id())

	updateOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}
	updateOpt.JSONBody = utils.RemoveNil(buildUpdateGaussDBMySQLFlavorBodyParams(d))

	updateResp, err := client.Request("PUT", updatePath, &updateOpt)
	if err != nil {
		return fmt.Errorf("error updating GaussDB MySQL proxy flavor: %s", err)
	}

	updateRespBody, err := utils.FlattenResponse(updateResp)
	if err != nil {
		return err
	}

	jobId := utils.PathSearch("job_id", updateRespBody, nil)
	if jobId == nil {
		return fmt.Errorf("error updating GaussDB MySQL proxy flavor: job_id is not found in API response")
	}

	err = checkGaussDBMySQLProxyJobFinish(ctx, client, jobId.(string), d.Timeout(schema.TimeoutUpdate))
	if err != nil {
		return err
	}

	return nil
}

func buildUpdateGaussDBMySQLFlavorBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"flavor_ref": d.Get("flavor"),
	}
	return bodyParams
}

func enlargeGaussDBMySQLProxyNumber(ctx context.Context, d *schema.ResourceData, client *golangsdk.ServiceClient,
	number int) error {
	var (
		httpUrl = "v3/{project_id}/instances/{instance_id}/proxy/enlarge"
	)
	updatePath := client.Endpoint + httpUrl
	updatePath = strings.ReplaceAll(updatePath, "{project_id}", client.ProjectID)
	updatePath = strings.ReplaceAll(updatePath, "{instance_id}", d.Get("instance_id").(string))

	updateOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}
	updateOpt.JSONBody = utils.RemoveNil(buildEnlargeGaussDBMySQLProxyNumberBodyParams(d, number))

	updateResp, err := client.Request("POST", updatePath, &updateOpt)
	if err != nil {
		return fmt.Errorf("error enlarging GaussDB MySQL proxy number: %s", err)
	}

	updateRespBody, err := utils.FlattenResponse(updateResp)
	if err != nil {
		return err
	}

	jobId := utils.PathSearch("job_id", updateRespBody, nil)
	if jobId == nil {
		return fmt.Errorf("error enlarging GaussDB MySQL proxy number: job_id is not found in API response")
	}

	err = checkGaussDBMySQLProxyJobFinish(ctx, client, jobId.(string), d.Timeout(schema.TimeoutUpdate))
	if err != nil {
		return err
	}

	return nil
}

func buildEnlargeGaussDBMySQLProxyNumberBodyParams(d *schema.ResourceData, number int) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"node_num": number,
		"proxy_id": d.Id(),
	}
	return bodyParams
}

func reduceGaussDBMySQLProxyNumber(ctx context.Context, d *schema.ResourceData, client *golangsdk.ServiceClient,
	number int) error {
	var (
		httpUrl = "v3/{project_id}/instances/{instance_id}/proxy/{proxy_id}/reduce"
	)
	updatePath := client.Endpoint + httpUrl
	updatePath = strings.ReplaceAll(updatePath, "{project_id}", client.ProjectID)
	updatePath = strings.ReplaceAll(updatePath, "{instance_id}", d.Get("instance_id").(string))
	updatePath = strings.ReplaceAll(updatePath, "{proxy_id}", d.Id())

	updateOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}
	updateOpt.JSONBody = utils.RemoveNil(buildReduceGaussDBMySQLProxyNumberBodyParams(number))

	updateResp, err := client.Request("PUT", updatePath, &updateOpt)
	if err != nil {
		return fmt.Errorf("error reducing GaussDB MySQL proxy number: %s", err)
	}

	updateRespBody, err := utils.FlattenResponse(updateResp)
	if err != nil {
		return err
	}

	jobId := utils.PathSearch("job_id", updateRespBody, nil)
	if jobId == nil {
		return fmt.Errorf("error reducing GaussDB MySQL proxy number: job_id is not found in API response")
	}

	err = checkGaussDBMySQLProxyJobFinish(ctx, client, jobId.(string), d.Timeout(schema.TimeoutUpdate))
	if err != nil {
		return err
	}

	return nil
}

func buildReduceGaussDBMySQLProxyNumberBodyParams(number int) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"node_num": number,
	}
	return bodyParams
}

func updateGaussDBMySQLProxyName(d *schema.ResourceData, client *golangsdk.ServiceClient) error {
	var (
		httpUrl = "v3/{project_id}/instances/{instance_id}/proxy/{proxy_id}/rename"
	)
	updatePath := client.Endpoint + httpUrl
	updatePath = strings.ReplaceAll(updatePath, "{project_id}", client.ProjectID)
	updatePath = strings.ReplaceAll(updatePath, "{instance_id}", d.Get("instance_id").(string))
	updatePath = strings.ReplaceAll(updatePath, "{proxy_id}", d.Id())

	updateOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}
	updateOpt.JSONBody = utils.RemoveNil(buildUpdateGaussDBMySQLProxyNameBodyParams(d))

	updateResp, err := client.Request("PUT", updatePath, &updateOpt)
	if err != nil {
		return fmt.Errorf("error updating GaussDB MySQL proxy name: %s", err)
	}

	updateRespBody, err := utils.FlattenResponse(updateResp)
	if err != nil {
		return err
	}

	result := utils.PathSearch("result", updateRespBody, "").(string)
	if result != "success" {
		return fmt.Errorf("error updating GaussDB MySQL proxy name: result is: %s", result)
	}

	return nil
}

func buildUpdateGaussDBMySQLProxyNameBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"alias": d.Get("proxy_name"),
	}
	return bodyParams
}

func updateGaussDBMySQLPort(ctx context.Context, d *schema.ResourceData, client *golangsdk.ServiceClient) error {
	var (
		httpUrl = "v3/{project_id}/instances/{instance_id}/proxy/{proxy_id}/port"
	)
	updatePath := client.Endpoint + httpUrl
	updatePath = strings.ReplaceAll(updatePath, "{project_id}", client.ProjectID)
	updatePath = strings.ReplaceAll(updatePath, "{instance_id}", d.Get("instance_id").(string))
	updatePath = strings.ReplaceAll(updatePath, "{proxy_id}", d.Id())

	updateOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}
	updateOpt.JSONBody = utils.RemoveNil(buildUpdateGaussDBMySQLPortBodyParams(d))

	updateResp, err := client.Request("PUT", updatePath, &updateOpt)
	if err != nil {
		return fmt.Errorf("error updating GaussDB MySQL proxy port: %s", err)
	}

	updateRespBody, err := utils.FlattenResponse(updateResp)
	if err != nil {
		return err
	}

	jobId := utils.PathSearch("job_id", updateRespBody, nil)
	if jobId == nil {
		// if the port is not change, then the job_id is null
		return nil
	}

	err = checkGaussDBMySQLProxyJobFinish(ctx, client, jobId.(string), d.Timeout(schema.TimeoutUpdate))
	if err != nil {
		return err
	}

	return nil
}

func buildUpdateGaussDBMySQLPortBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"port": d.Get("port"),
	}
	return bodyParams
}

func updateGaussDBMySQLNodesWeight(ctx context.Context, d *schema.ResourceData, client *golangsdk.ServiceClient) error {
	var (
		httpUrl = "v3/{project_id}/instances/{instance_id}/proxy/{proxy_id}/weight"
	)
	updatePath := client.Endpoint + httpUrl
	updatePath = strings.ReplaceAll(updatePath, "{project_id}", client.ProjectID)
	updatePath = strings.ReplaceAll(updatePath, "{instance_id}", d.Get("instance_id").(string))
	updatePath = strings.ReplaceAll(updatePath, "{proxy_id}", d.Id())

	updateOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}
	updateOpt.JSONBody = buildUpdateGaussDBMySQLNodesWeightBodyParams(d)

	updateResp, err := client.Request("PUT", updatePath, &updateOpt)
	if err != nil {
		return fmt.Errorf("error updating GaussDB MySQL nodes weight associated with proxy: %s", err)
	}

	updateRespBody, err := utils.FlattenResponse(updateResp)
	if err != nil {
		return err
	}

	jobId := utils.PathSearch("job_id", updateRespBody, nil)
	if jobId == nil {
		return fmt.Errorf("error updating GaussDB MySQL nodes weight associated with proxy: job_id is not " +
			"found in API response")
	}

	err = checkGaussDBMySQLProxyJobFinish(ctx, client, jobId.(string), d.Timeout(schema.TimeoutUpdate))
	if err != nil {
		return err
	}

	return nil
}

func buildUpdateGaussDBMySQLNodesWeightBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := make(map[string]interface{})
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
				"id":     raw["id"],
				"weight": raw["weight"],
			}
		}
		bodyParams["readonly_nodes"] = readonlyNodesWeightParam
	}
	return bodyParams
}

func updateGaussDBMySQLNewNode(d *schema.ResourceData, client *golangsdk.ServiceClient) error {
	var (
		httpUrl = "v3/{project_id}/instances/{instance_id}/proxy/{proxy_id}/new-node-auto-add"
	)
	updatePath := client.Endpoint + httpUrl
	updatePath = strings.ReplaceAll(updatePath, "{project_id}", client.ProjectID)
	updatePath = strings.ReplaceAll(updatePath, "{instance_id}", d.Get("instance_id").(string))
	updatePath = strings.ReplaceAll(updatePath, "{proxy_id}", d.Id())

	updateOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}
	updateOpt.JSONBody = utils.RemoveNil(buildUpdateGaussDBMySQLNewNodeBodyParams(d))

	updateResp, err := client.Request("POST", updatePath, &updateOpt)
	if err != nil {
		return fmt.Errorf("error updating GaussDB MySQL new node automatically associate with proxy: %s", err)
	}

	updateRespBody, err := utils.FlattenResponse(updateResp)
	if err != nil {
		return err
	}

	result := utils.PathSearch("result", updateRespBody, "").(string)
	if result != "success" {
		return fmt.Errorf("error updating GaussDB MySQL new node automatically associate with proxy: result is: %s",
			result)
	}

	return nil
}

func buildUpdateGaussDBMySQLNewNodeBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"switch_status": d.Get("new_node_auto_add_status"),
		"weight":        utils.ValueIgnoreEmpty(d.Get("new_node_weight")),
	}
	return bodyParams
}

func updateGaussDBMySQLParameters(d *schema.ResourceData, client *golangsdk.ServiceClient) error {
	var (
		httpUrl = "v3/{project_id}/instances/{instance_id}/proxy/{proxy_id}/configurations"
	)
	updatePath := client.Endpoint + httpUrl
	updatePath = strings.ReplaceAll(updatePath, "{project_id}", client.ProjectID)
	updatePath = strings.ReplaceAll(updatePath, "{instance_id}", d.Get("instance_id").(string))
	updatePath = strings.ReplaceAll(updatePath, "{proxy_id}", d.Id())

	updateOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}
	updateOpt.JSONBody = utils.RemoveNil(buildUpdateGaussDBMySQLParametersBodyParams(d))

	_, err := client.Request("PUT", updatePath, &updateOpt)
	if err != nil {
		return fmt.Errorf("error updating GaussDB MySQL proxy parameters: %s", err)
	}

	return nil
}

func buildUpdateGaussDBMySQLParametersBodyParams(d *schema.ResourceData) map[string]interface{} {
	parametersRaw := d.Get("parameters").(*schema.Set).List()
	parameters := make([]map[string]interface{}, len(parametersRaw))
	for i, v := range parametersRaw {
		raw := v.(map[string]interface{})
		parameters[i] = map[string]interface{}{
			"name":      raw["name"],
			"value":     raw["value"],
			"elem_type": raw["elem_type"],
		}
	}
	bodyParams := map[string]interface{}{
		"configurations": parameters,
	}
	return bodyParams
}

func updateGaussDBMySQLTransactionSplit(ctx context.Context, d *schema.ResourceData, client *golangsdk.ServiceClient) error {
	var (
		httpUrl = "v3/{project_id}/instances/{instance_id}/proxy/transaction-split"
	)
	updatePath := client.Endpoint + httpUrl
	updatePath = strings.ReplaceAll(updatePath, "{project_id}", client.ProjectID)
	updatePath = strings.ReplaceAll(updatePath, "{instance_id}", d.Get("instance_id").(string))

	updateOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}
	updateOpt.JSONBody = utils.RemoveNil(buildUpdateGaussDBMySQLTransactionSplitBodyParams(d))

	updateResp, err := client.Request("POST", updatePath, &updateOpt)
	if err != nil {
		return fmt.Errorf("error updating GaussDB MySQL proxy transaction split: %s", err)
	}

	updateRespBody, err := utils.FlattenResponse(updateResp)
	if err != nil {
		return err
	}

	jobId := utils.PathSearch("job_id", updateRespBody, nil)
	if jobId == nil {
		return fmt.Errorf("error updating GaussDB MySQL proxy transaction split: job_id is not found in API response")
	}

	err = checkGaussDBMySQLProxyJobFinish(ctx, client, jobId.(string), d.Timeout(schema.TimeoutUpdate))
	if err != nil {
		return err
	}

	return nil
}

func buildUpdateGaussDBMySQLTransactionSplitBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"transaction_split": d.Get("transaction_split"),
		"proxy_id_list":     []string{d.Id()},
	}
	return bodyParams
}

func updateGaussDBMySQLConsistenceMode(ctx context.Context, d *schema.ResourceData, client *golangsdk.ServiceClient) error {
	var (
		httpUrl = "v3/{project_id}/instances/{instance_id}/proxy/{proxy_id}/session-consistence"
	)
	updatePath := client.Endpoint + httpUrl
	updatePath = strings.ReplaceAll(updatePath, "{project_id}", client.ProjectID)
	updatePath = strings.ReplaceAll(updatePath, "{instance_id}", d.Get("instance_id").(string))
	updatePath = strings.ReplaceAll(updatePath, "{proxy_id}", d.Id())

	updateOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}
	updateOpt.JSONBody = buildUpdateGaussDBMySQLConsistenceModeBodyParams(d)

	updateResp, err := client.Request("PUT", updatePath, &updateOpt)
	if err != nil {
		return fmt.Errorf("error updating GaussDB MySQL proxy consistence mode: %s", err)
	}

	updateRespBody, err := utils.FlattenResponse(updateResp)
	if err != nil {
		return err
	}

	jobId := utils.PathSearch("job_id", updateRespBody, nil)
	if jobId == nil {
		return fmt.Errorf("error updating GaussDB MySQL proxy consistence mode: job_id is not found in API response")
	}

	err = checkGaussDBMySQLProxyJobFinish(ctx, client, jobId.(string), d.Timeout(schema.TimeoutUpdate))
	if err != nil {
		return err
	}

	return nil
}

func buildUpdateGaussDBMySQLConsistenceModeBodyParams(d *schema.ResourceData) map[string]interface{} {
	consistenceMode := d.Get("consistence_mode").(string)
	bodyParams := map[string]interface{}{
		"consistence_mode": consistenceMode,
	}
	if consistenceMode == "session" {
		bodyParams["session_consistence"] = true
	} else {
		bodyParams["session_consistence"] = false
	}
	return bodyParams
}

func updateGaussDBMySQLConnectionPoolType(ctx context.Context, d *schema.ResourceData, client *golangsdk.ServiceClient) error {
	var (
		httpUrl = "v3/{project_id}/instances/{instance_id}/proxy/{proxy_id}/connection-pool-type"
	)
	updatePath := client.Endpoint + httpUrl
	updatePath = strings.ReplaceAll(updatePath, "{project_id}", client.ProjectID)
	updatePath = strings.ReplaceAll(updatePath, "{instance_id}", d.Get("instance_id").(string))
	updatePath = strings.ReplaceAll(updatePath, "{proxy_id}", d.Id())

	updateOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}
	updateOpt.JSONBody = utils.RemoveNil(buildUpdateGaussDBMySQLConnectionPoolTypeBodyParams(d))

	updateResp, err := client.Request("PUT", updatePath, &updateOpt)
	if err != nil {
		return fmt.Errorf("error updating GaussDB MySQL proxy connection pool type: %s", err)
	}

	updateRespBody, err := utils.FlattenResponse(updateResp)
	if err != nil {
		return err
	}

	jobId := utils.PathSearch("job_id", updateRespBody, nil)
	if jobId == nil {
		return fmt.Errorf("error updating GaussDB MySQL proxy connection pool type: job_id is not found in API response")
	}

	err = checkGaussDBMySQLProxyJobFinish(ctx, client, jobId.(string), d.Timeout(schema.TimeoutUpdate))
	if err != nil {
		return err
	}

	return nil
}

func buildUpdateGaussDBMySQLConnectionPoolTypeBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"connection_pool_type": d.Get("connection_pool_type"),
	}
	return bodyParams
}

func updateGaussDBMySQLOpenAccessControl(d *schema.ResourceData, client *golangsdk.ServiceClient) error {
	var (
		httpUrl = "v3/{project_id}/instances/{instance_id}/proxy/{proxy_id}/access-control-switch"
	)
	updatePath := client.Endpoint + httpUrl
	updatePath = strings.ReplaceAll(updatePath, "{project_id}", client.ProjectID)
	updatePath = strings.ReplaceAll(updatePath, "{instance_id}", d.Get("instance_id").(string))
	updatePath = strings.ReplaceAll(updatePath, "{proxy_id}", d.Id())

	updateOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}
	updateOpt.JSONBody = utils.RemoveNil(buildUpdateGaussDBMySQLOpenAccessControlBodyParams(d))

	_, err := client.Request("POST", updatePath, &updateOpt)
	if err != nil {
		return fmt.Errorf("error updating GaussDB MySQL proxy access control switch: %s", err)
	}

	return nil
}

func buildUpdateGaussDBMySQLOpenAccessControlBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"open_access_control": d.Get("open_access_control"),
	}
	return bodyParams
}

func updateGaussDBMySQLAccessControlRule(d *schema.ResourceData, client *golangsdk.ServiceClient) error {
	var (
		httpUrl = "v3/{project_id}/instances/{instance_id}/proxy/{proxy_id}/access-control"
	)
	updatePath := client.Endpoint + httpUrl
	updatePath = strings.ReplaceAll(updatePath, "{project_id}", client.ProjectID)
	updatePath = strings.ReplaceAll(updatePath, "{instance_id}", d.Get("instance_id").(string))
	updatePath = strings.ReplaceAll(updatePath, "{proxy_id}", d.Id())

	updateOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}
	updateOpt.JSONBody = utils.RemoveNil(buildUpdateGaussDBMySQLAccessControlRuleBodyParams(d))

	_, err := client.Request("POST", updatePath, &updateOpt)
	if err != nil {
		return fmt.Errorf("error updating GaussDB MySQL proxy access control: %s", err)
	}

	return nil
}

func buildUpdateGaussDBMySQLAccessControlRuleBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"type": d.Get("access_control_type"),
	}
	accessControlIpList := make([]interface{}, 0)
	for _, v := range d.Get("access_control_ip_list").(*schema.Set).List() {
		raw := v.(map[string]interface{})
		accessControlIpList = append(accessControlIpList, map[string]interface{}{
			"ip":          raw["ip"],
			"description": raw["description"],
		})
	}
	bodyParams["ip_list"] = accessControlIpList
	return bodyParams
}

func resourceGaussDBProxyDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var (
		httpUrl = "v3/{project_id}/instances/{instance_id}/proxy"
		product = "gaussdb"
	)
	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating GaussDB client: %s", err)
	}

	deletePath := client.Endpoint + httpUrl
	deletePath = strings.ReplaceAll(deletePath, "{project_id}", client.ProjectID)
	deletePath = strings.ReplaceAll(deletePath, "{instance_id}", d.Get("instance_id").(string))

	deleteOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}
	deleteOpt.JSONBody = utils.RemoveNil(buildDeleteGaussDBMySQLProxyBodyParams(d))

	deleteResp, err := client.Request("DELETE", deletePath, &deleteOpt)
	if err != nil {
		return common.CheckDeletedDiag(d, parseMysqlProxyError(err), "error deleting GaussDB MySQL proxy")
	}

	deleteRespBody, err := utils.FlattenResponse(deleteResp)
	if err != nil {
		return diag.FromErr(err)
	}

	jobId := utils.PathSearch("job_id", deleteRespBody, nil)
	if jobId == nil {
		return diag.Errorf("error deleting GaussDB MySQL proxy: job_id is not found in API response")
	}

	err = checkGaussDBMySQLProxyJobFinish(ctx, client, jobId.(string), d.Timeout(schema.TimeoutDelete))
	if err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func buildDeleteGaussDBMySQLProxyBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"proxy_ids": []string{d.Id()},
	}
	return bodyParams
}

func checkGaussDBMySQLProxyJobFinish(ctx context.Context, client *golangsdk.ServiceClient, jobID string,
	timeout time.Duration) error {
	stateConf := &resource.StateChangeConf{
		Pending:      []string{"Pending", "Running"},
		Target:       []string{"Completed"},
		Refresh:      gaussDBMysqlDatabaseStatusRefreshFunc(client, jobID),
		Timeout:      timeout,
		PollInterval: 10 * time.Second,
	}
	if _, err := stateConf.WaitForStateContext(ctx); err != nil {
		return fmt.Errorf("error waiting for GaussDB MySQL proxy job (%s) to be completed: %s ", jobID, err)
	}
	return nil
}

func parseMysqlProxyError(err error) error {
	if parsedErr, ok := err.(golangsdk.ErrDefault400); ok {
		return common.ConvertExpected400ErrInto404Err(parsedErr, "error_code", "DBS.201028")
	}
	return common.ConvertUndefinedErrInto404Err(err, 409, "error_code", "DBS.200932")
}

func resourceGaussDBMySQLProxyImportState(_ context.Context, d *schema.ResourceData, _ interface{}) ([]*schema.ResourceData,
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
