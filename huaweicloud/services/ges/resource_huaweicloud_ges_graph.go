// ---------------------------------------------------------------
// *** AUTO GENERATED CODE ***
// @Product GES
// ---------------------------------------------------------------

package ges

import (
	"context"
	"fmt"
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

// @API GES POST /v2/{project_id}/graphs/{id}/resize
// @API GES POST /v2/{project_id}/graphs
// @API GES GET /v2/{project_id}/graphs/{id}
// @API GES DELETE /v2/{project_id}/graphs/{id}
// @API GES POST /v2/{project_id}/graphs/{id}/expand
func ResourceGesGraph() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceGesGraphCreate,
		UpdateContext: resourceGesGraphUpdate,
		ReadContext:   resourceGesGraphRead,
		DeleteContext: resourceGesGraphDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(60 * time.Minute),
			Update: schema.DefaultTimeout(60 * time.Minute),
			Delete: schema.DefaultTimeout(60 * time.Minute),
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
				ForceNew:    true,
				Description: `The graph name.`,
			},
			"graph_size_type_index": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `Graph size type index.`,
			},
			"vpc_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: `The VPC ID.`,
			},
			"subnet_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: `The subnet ID.`,
			},
			"security_group_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: `The security group ID.`,
			},
			"crypt_algorithm": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: `Graph instance cryptography algorithm.`,
			},
			"enable_https": {
				Type:        schema.TypeBool,
				Required:    true,
				ForceNew:    true,
				Description: `Whether to enable the security mode. This mode may damage GES performance greatly.`,
			},
			"cpu_arch": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: `Graph instance's CPU architecture type.`,
			},
			"public_ip": {
				Type:        schema.TypeList,
				MaxItems:    1,
				Elem:        graphPublicIpSchema(),
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: `The information about public IP.`,
			},
			"enable_multi_az": {
				Type:        schema.TypeBool,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: `Whether the created graph supports the cross-AZ mode. The default value is false.`,
			},
			"encryption": {
				Type:        schema.TypeList,
				MaxItems:    1,
				Elem:        graphEncryptionSchema(),
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: `Whether to enable data encryption The value can be true or false.`,
			},
			"lts_operation_trace": {
				Type:     schema.TypeList,
				MaxItems: 1,
				Elem:     graphLtsOperationTraceSchema(),
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"enterprise_project_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: `The enterprise project ID.`,
			},
			"tags": common.TagsForceNewSchema(),
			"enable_rbac": {
				Type:        schema.TypeBool,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: `Whether to enable granular permission control for the created graph.`,
			},
			"enable_full_text_index": {
				Type:        schema.TypeBool,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: `Whether to enable full-text index control for the created graph.`,
			},
			"enable_hyg": {
				Type:        schema.TypeBool,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: `Whether to enable HyG for the graph. This parameter is available for database edition graphs only.`,
			},
			"product_type": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: `Graph product type`,
			},
			"vertex_id_type": {
				Type:        schema.TypeList,
				MaxItems:    1,
				Elem:        graphVertexIdTypeSchema(),
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: `ID type of vertices. This parameter is mandatory only for database edition graphs.`,
			},
			"replication": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},
			"keep_backup": {
				Type:        schema.TypeBool,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: `Whether to retain the backups of a graph after it is deleted.`,
			},
			"az_code": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `AZ code`,
			},
			"status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Status of a graph.`,
			},
			"private_ip": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Floating IP address of a graph instance.`,
			},
			"traffic_ip_list": {
				Type:        schema.TypeList,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Computed:    true,
				Description: `Physical addresses of a graph instance for access from private networks.`,
			},
		},
	}
}

func graphPublicIpSchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"public_bind_type": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `The bind type of public IP.`,
			},
			"eip_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `The EIP ID.`,
			},
		},
	}
	return &sc
}

func graphEncryptionSchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"enable": {
				Type:        schema.TypeBool,
				Optional:    true,
				Computed:    true,
				Description: `Whether to enable data encryption The value can be true or false. The default value is false.`,
			},
			"master_key_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `ID of the customer master key created by DEW in the project corresponding to the graph creation.`,
			},
		},
	}
	return &sc
}

func graphLtsOperationTraceSchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"enable_audit": {
				Type:        schema.TypeBool,
				Optional:    true,
				Computed:    true,
				Description: `Whether to enable graph audit. The default value is false.`,
			},
			"audit_log_group_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `LTS log group name.`,
			},
		},
	}
	return &sc
}

func graphVertexIdTypeSchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"id_type": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `Vertex ID type.`,
			},
			"id_length": {
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
				Description: `The length of ID.`,
			},
		},
	}
	return &sc
}

func resourceGesGraphCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	// createGraph: create a GES graph.
	var (
		createGraphHttpUrl = "v2/{project_id}/graphs"
		createGraphProduct = "ges"
	)
	createGraphClient, err := cfg.NewServiceClient(createGraphProduct, region)
	if err != nil {
		return diag.Errorf("error creating GES Client: %s", err)
	}

	createGraphPath := createGraphClient.Endpoint + createGraphHttpUrl
	createGraphPath = strings.ReplaceAll(createGraphPath, "{project_id}", createGraphClient.ProjectID)

	createGraphOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		OkCodes: []int{
			200,
		},
		MoreHeaders: map[string]string{"Content-Type": "application/json;charset=UTF-8"},
	}

	createGraphOpt.JSONBody = utils.RemoveNil(buildCreateGraphBodyParams(d, cfg))
	createGraphResp, err := createGraphClient.Request("POST", createGraphPath, &createGraphOpt)
	if err != nil {
		return diag.Errorf("error creating GesGraph: %s", err)
	}

	createGraphRespBody, err := utils.FlattenResponse(createGraphResp)
	if err != nil {
		return diag.FromErr(err)
	}

	graphId := utils.PathSearch("graph_id", createGraphRespBody, "").(string)
	if graphId == "" {
		return diag.Errorf("unable to find the GES graph ID from the API response")
	}
	d.SetId(graphId)

	err = graphWaitingForStateCompleted(ctx, d, meta, d.Timeout(schema.TimeoutCreate))
	if err != nil {
		return diag.Errorf("error waiting for the Create of GesGraph (%s) to complete: %s", graphId, err)
	}
	return resourceGesGraphRead(ctx, d, meta)
}

func buildCreateGraphBodyParams(d *schema.ResourceData, cfg *config.Config) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"graph": map[string]interface{}{
			"name":                   utils.ValueIgnoreEmpty(d.Get("name")),
			"graph_size_type_index":  utils.ValueIgnoreEmpty(d.Get("graph_size_type_index")),
			"arch":                   utils.ValueIgnoreEmpty(d.Get("cpu_arch")),
			"vpc_id":                 utils.ValueIgnoreEmpty(d.Get("vpc_id")),
			"subnet_id":              utils.ValueIgnoreEmpty(d.Get("subnet_id")),
			"security_group_id":      utils.ValueIgnoreEmpty(d.Get("security_group_id")),
			"public_ip":              buildCreateGraphReqBodyPublicIp(d.Get("public_ip")),
			"enable_multi_az":        utils.ValueIgnoreEmpty(d.Get("enable_multi_az")),
			"encryption":             buildCreateGraphReqBodyEncryption(d.Get("encryption")),
			"lts_operation_trace":    buildCreateGraphReqBodyLtsOperationTrace(d.Get("lts_operation_trace")),
			"sys_tags":               utils.BuildSysTags(cfg.GetEnterpriseProjectID(d)),
			"tags":                   utils.ExpandResourceTags(d.Get("tags").(map[string]interface{})),
			"enable_rbac":            utils.ValueIgnoreEmpty(d.Get("enable_rbac")),
			"enable_full_text_index": utils.ValueIgnoreEmpty(d.Get("enable_full_text_index")),
			"enable_hyg":             utils.ValueIgnoreEmpty(d.Get("enable_hyg")),
			"crypt_algorithm":        utils.ValueIgnoreEmpty(d.Get("crypt_algorithm")),
			"enable_https":           utils.ValueIgnoreEmpty(d.Get("enable_https")),
			"product_type":           utils.ValueIgnoreEmpty(d.Get("product_type")),
			"vertex_id_type":         buildCreateGraphReqBodyvertexIdType(d.Get("vertex_id_type")),
		},
	}
	return bodyParams
}

func buildCreateGraphReqBodyPublicIp(rawParams interface{}) map[string]interface{} {
	if rawArray, ok := rawParams.([]interface{}); ok {
		if len(rawArray) == 0 {
			return nil
		}
		raw := rawArray[0].(map[string]interface{})
		params := map[string]interface{}{
			"public_bind_type": utils.ValueIgnoreEmpty(raw["public_bind_type"]),
			"eip_id":           utils.ValueIgnoreEmpty(raw["eip_id"]),
		}
		return params
	}
	return nil
}

func buildCreateGraphReqBodyEncryption(rawParams interface{}) map[string]interface{} {
	if rawArray, ok := rawParams.([]interface{}); ok {
		if len(rawArray) == 0 {
			return nil
		}
		raw := rawArray[0].(map[string]interface{})
		params := map[string]interface{}{
			"enable":        utils.ValueIgnoreEmpty(raw["enable"]),
			"master_key_id": utils.ValueIgnoreEmpty(raw["master_key_id"]),
		}
		return params
	}
	return nil
}

func buildCreateGraphReqBodyLtsOperationTrace(rawParams interface{}) map[string]interface{} {
	if rawArray, ok := rawParams.([]interface{}); ok {
		if len(rawArray) == 0 {
			return nil
		}
		raw := rawArray[0].(map[string]interface{})
		params := map[string]interface{}{
			"enable_audit":         utils.ValueIgnoreEmpty(raw["enable_audit"]),
			"audit_log_group_name": utils.ValueIgnoreEmpty(raw["audit_log_group_name"]),
		}
		return params
	}
	return nil
}

func buildCreateGraphReqBodyvertexIdType(rawParams interface{}) map[string]interface{} {
	if rawArray, ok := rawParams.([]interface{}); ok {
		if len(rawArray) == 0 {
			return nil
		}
		raw := rawArray[0].(map[string]interface{})
		params := map[string]interface{}{
			"id_type":   utils.ValueIgnoreEmpty(raw["id_type"]),
			"id_length": utils.ValueIgnoreEmpty(raw["id_length"]),
		}
		return params
	}
	return nil
}

func graphWaitingForStateCompleted(ctx context.Context, d *schema.ResourceData, meta interface{}, t time.Duration) error {
	stateConf := &resource.StateChangeConf{
		Pending:                   []string{"PENDING"},
		Target:                    []string{"COMPLETED"},
		ContinuousTargetOccurence: 3,
		Refresh: func() (interface{}, string, error) {
			cfg := meta.(*config.Config)
			region := cfg.GetRegion(d)
			// createGraphWaiting: waiting graph is available
			var (
				createGraphWaitingHttpUrl = "v2/{project_id}/graphs/{id}"
				createGraphWaitingProduct = "ges"
			)
			createGraphWaitingClient, err := cfg.NewServiceClient(createGraphWaitingProduct, region)
			if err != nil {
				return nil, "ERROR", fmt.Errorf("error creating GES Client: %s", err)
			}

			createGraphWaitingPath := createGraphWaitingClient.Endpoint + createGraphWaitingHttpUrl
			createGraphWaitingPath = strings.ReplaceAll(createGraphWaitingPath, "{project_id}", createGraphWaitingClient.ProjectID)
			createGraphWaitingPath = strings.ReplaceAll(createGraphWaitingPath, "{id}", d.Id())

			createGraphWaitingOpt := golangsdk.RequestOpts{
				KeepResponseBody: true,
				OkCodes: []int{
					200,
				},
				MoreHeaders: map[string]string{"Content-Type": "application/json;charset=UTF-8"},
			}

			createGraphWaitingResp, err := createGraphWaitingClient.Request("GET", createGraphWaitingPath, &createGraphWaitingOpt)
			if err != nil {
				return nil, "ERROR", err
			}

			createGraphWaitingRespBody, err := utils.FlattenResponse(createGraphWaitingResp)
			if err != nil {
				return nil, "ERROR", err
			}
			statusRaw := utils.PathSearch(`graph.status`, createGraphWaitingRespBody, nil)

			status := fmt.Sprintf("%v", statusRaw)

			targetStatus := []string{
				"200",
			}
			if utils.StrSliceContains(targetStatus, status) {
				return createGraphWaitingRespBody, "COMPLETED", nil
			}

			unexpectedStatus := []string{
				"300",
				"303",
				"800",
			}
			if utils.StrSliceContains(unexpectedStatus, status) {
				return createGraphWaitingRespBody, status, nil
			}

			return createGraphWaitingRespBody, "PENDING", nil
		},
		Timeout:      t,
		Delay:        30 * time.Second,
		PollInterval: 30 * time.Second,
	}
	_, err := stateConf.WaitForStateContext(ctx)
	return err
}

func resourceGesGraphRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var mErr *multierror.Error

	// getGraph: Query the GES graph.
	var (
		getGraphHttpUrl = "v2/{project_id}/graphs/{id}"
		getGraphProduct = "ges"
	)
	getGraphClient, err := cfg.NewServiceClient(getGraphProduct, region)
	if err != nil {
		return diag.Errorf("error creating GES Client: %s", err)
	}

	getGraphPath := getGraphClient.Endpoint + getGraphHttpUrl
	getGraphPath = strings.ReplaceAll(getGraphPath, "{project_id}", getGraphClient.ProjectID)
	getGraphPath = strings.ReplaceAll(getGraphPath, "{id}", d.Id())

	getGraphOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		OkCodes: []int{
			200,
		},
		MoreHeaders: map[string]string{"Content-Type": "application/json;charset=UTF-8"},
	}

	getGraphResp, err := getGraphClient.Request("GET", getGraphPath, &getGraphOpt)

	if err != nil {
		return common.CheckDeletedDiag(d, err, "error retrieving GesGraph")
	}

	getGraphRespBody, err := utils.FlattenResponse(getGraphResp)
	if err != nil {
		return diag.FromErr(err)
	}

	mErr = multierror.Append(
		mErr,
		d.Set("region", region),
		d.Set("name", utils.PathSearch("graph.name", getGraphRespBody, nil)),
		d.Set("graph_size_type_index", utils.PathSearch("graph.graph_size_type_index", getGraphRespBody, nil)),
		d.Set("cpu_arch", utils.PathSearch("graph.arch", getGraphRespBody, nil)),
		d.Set("vpc_id", utils.PathSearch("graph.vpc_id", getGraphRespBody, nil)),
		d.Set("subnet_id", utils.PathSearch("graph.subnet_id", getGraphRespBody, nil)),
		d.Set("security_group_id", utils.PathSearch("graph.security_group_id", getGraphRespBody, nil)),
		d.Set("public_ip", flattenGetGraphRespBodyPublicIp(getGraphRespBody)),
		d.Set("enable_multi_az", utils.StringToBool(utils.PathSearch("graph.is_multi_az", getGraphRespBody, nil))),
		d.Set("az_code", utils.PathSearch("graph.az_code", getGraphRespBody, nil)),
		d.Set("enterprise_project_id", utils.PathSearch("graph.sys_tags[0]", getGraphRespBody, nil)),
		d.Set("tags", utils.FlattenTagsToMap(utils.PathSearch("graph.tags", getGraphRespBody, nil))),
		d.Set("enable_rbac", utils.PathSearch("graph.enable_rbac", getGraphRespBody, nil)),
		d.Set("enable_full_text_index", utils.PathSearch("graph.enable_full_text_index", getGraphRespBody, nil)),
		d.Set("enable_hyg", utils.PathSearch("graph.enable_hyg", getGraphRespBody, nil)),
		d.Set("crypt_algorithm", utils.PathSearch("graph.crypt_algorithm", getGraphRespBody, nil)),
		d.Set("enable_https", utils.PathSearch("graph.enable_https", getGraphRespBody, nil)),
		d.Set("product_type", utils.PathSearch("graph.product_type", getGraphRespBody, nil)),
		d.Set("vertex_id_type", flattenGetGraphRespBodyvertexIdType(getGraphRespBody)),
		d.Set("status", utils.PathSearch("graph.status", getGraphRespBody, nil)),
		d.Set("replication", utils.PathSearch("graph.replication", getGraphRespBody, nil)),
		d.Set("private_ip", utils.PathSearch("graph.private_ip", getGraphRespBody, nil)),
		d.Set("traffic_ip_list", utils.PathSearch("graph.traffic_ip_list", getGraphRespBody, nil)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenGetGraphRespBodyPublicIp(resp interface{}) []interface{} {
	var rst []interface{}
	curJson := utils.PathSearch("graph.public_ip", resp, make(map[string]interface{})).(map[string]interface{})
	if len(curJson) < 1 {
		return rst
	}

	rst = []interface{}{
		map[string]interface{}{
			"public_bind_type": utils.PathSearch("public_bind_type", curJson, nil),
			"eip_id":           utils.PathSearch("eip_id", curJson, nil),
		},
	}
	return rst
}

func flattenGetGraphRespBodyvertexIdType(resp interface{}) []interface{} {
	var rst []interface{}
	curJson := utils.PathSearch("graph.vertex_id_type", resp, make(map[string]interface{})).(map[string]interface{})
	if len(curJson) < 1 {
		return rst
	}

	rst = []interface{}{
		map[string]interface{}{
			"id_type":   utils.PathSearch("id_type", curJson, nil),
			"id_length": utils.PathSearch("id_length", curJson, nil),
		},
	}
	return rst
}

func resourceGesGraphUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	resizeGraphChanges := []string{
		"graph_size_type_index",
	}

	if d.HasChanges(resizeGraphChanges...) {
		// resizeGraph: resize GES graph
		var (
			resizeGraphHttpUrl = "v2/{project_id}/graphs/{id}/resize"
			resizeGraphProduct = "ges"
		)
		resizeGraphClient, err := cfg.NewServiceClient(resizeGraphProduct, region)
		if err != nil {
			return diag.Errorf("error creating GES Client: %s", err)
		}

		resizeGraphPath := resizeGraphClient.Endpoint + resizeGraphHttpUrl
		resizeGraphPath = strings.ReplaceAll(resizeGraphPath, "{project_id}", resizeGraphClient.ProjectID)
		resizeGraphPath = strings.ReplaceAll(resizeGraphPath, "{id}", d.Id())

		resizeGraphOpt := golangsdk.RequestOpts{
			KeepResponseBody: true,
			OkCodes: []int{
				200,
			},
			MoreHeaders: map[string]string{"Content-Type": "application/json;charset=UTF-8"},
		}

		resizeGraphOpt.JSONBody = utils.RemoveNil(buildResizeGraphBodyParams(d))
		_, err = resizeGraphClient.Request("POST", resizeGraphPath, &resizeGraphOpt)
		if err != nil {
			return diag.Errorf("error updating GesGraph: %s", err)
		}

		err = graphWaitingForStateCompleted(ctx, d, meta, d.Timeout(schema.TimeoutUpdate))
		if err != nil {
			return diag.Errorf("error waiting for the Update of GesGraph (%s) to complete: %s", d.Id(), err)
		}
	}

	expandGraphReplicationChanges := []string{
		"replication",
	}

	if d.HasChanges(expandGraphReplicationChanges...) {
		// expandGraphReplication: expand GES graph
		var (
			expandGraphReplicationHttpUrl = "v2/{project_id}/graphs/{id}/expand"
			expandGraphReplicationProduct = "ges"
		)
		expandGraphReplicationClient, err := cfg.NewServiceClient(expandGraphReplicationProduct, region)
		if err != nil {
			return diag.Errorf("error creating GES Client: %s", err)
		}

		expandGraphReplicationPath := expandGraphReplicationClient.Endpoint + expandGraphReplicationHttpUrl
		expandGraphReplicationPath = strings.ReplaceAll(expandGraphReplicationPath, "{project_id}", expandGraphReplicationClient.ProjectID)
		expandGraphReplicationPath = strings.ReplaceAll(expandGraphReplicationPath, "{id}", d.Id())

		expandGraphReplicationOpt := golangsdk.RequestOpts{
			KeepResponseBody: true,
			OkCodes: []int{
				200,
			},
			MoreHeaders: map[string]string{"Content-Type": "application/json;charset=UTF-8"},
		}

		expandGraphReplicationOpt.JSONBody = utils.RemoveNil(buildExpandGraphReplicationBodyParams(d))
		_, err = expandGraphReplicationClient.Request("POST", expandGraphReplicationPath, &expandGraphReplicationOpt)
		if err != nil {
			return diag.Errorf("error updating GesGraph: %s", err)
		}
		err = graphWaitingForStateCompleted(ctx, d, meta, d.Timeout(schema.TimeoutUpdate))
		if err != nil {
			return diag.Errorf("error waiting for the Update of GesGraph (%s) to complete: %s", d.Id(), err)
		}
	}

	return resourceGesGraphRead(ctx, d, meta)
}

func buildExpandGraphReplicationBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"expand": map[string]interface{}{
			"replication": utils.ValueIgnoreEmpty(d.Get("replication")),
		},
	}
	return bodyParams
}

func buildResizeGraphBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"resize": map[string]interface{}{
			"graph_size_type_index": utils.ValueIgnoreEmpty(d.Get("graph_size_type_index")),
		},
	}
	return bodyParams
}

func resourceGesGraphDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	// deleteGraph: delete GES graph
	var (
		deleteGraphHttpUrl = "v2/{project_id}/graphs/{id}"
		deleteGraphProduct = "ges"
	)
	deleteGraphClient, err := cfg.NewServiceClient(deleteGraphProduct, region)
	if err != nil {
		return diag.Errorf("error creating GES Client: %s", err)
	}

	deleteGraphPath := deleteGraphClient.Endpoint + deleteGraphHttpUrl
	deleteGraphPath = strings.ReplaceAll(deleteGraphPath, "{project_id}", deleteGraphClient.ProjectID)
	deleteGraphPath = strings.ReplaceAll(deleteGraphPath, "{id}", d.Id())

	deleteGraphOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		OkCodes: []int{
			200,
		},
		MoreHeaders: map[string]string{"Content-Type": "application/json;charset=UTF-8"},
	}

	deleteGraphOpt.JSONBody = utils.RemoveNil(buildDeleteGraphBodyParams(d))
	_, err = deleteGraphClient.Request("DELETE", deleteGraphPath, &deleteGraphOpt)
	if err != nil {
		return diag.Errorf("error deleting GesGraph: %s", err)
	}

	err = deleteGraphWaitingForStateCompleted(ctx, d, meta, d.Timeout(schema.TimeoutDelete))
	if err != nil {
		return diag.Errorf("error waiting for the Delete of GesGraph (%s) to complete: %s", d.Id(), err)
	}
	return nil
}

func buildDeleteGraphBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"keep_backup": utils.ValueIgnoreEmpty(d.Get("keep_backup")),
	}
	return bodyParams
}

func deleteGraphWaitingForStateCompleted(ctx context.Context, d *schema.ResourceData, meta interface{}, t time.Duration) error {
	stateConf := &resource.StateChangeConf{
		Pending:                   []string{"PENDING"},
		Target:                    []string{"COMPLETED"},
		ContinuousTargetOccurence: 3,
		Refresh: func() (interface{}, string, error) {
			cfg := meta.(*config.Config)
			region := cfg.GetRegion(d)
			// deleteGraphWaiting: missing operation notes
			var (
				deleteGraphWaitingHttpUrl = "v2/{project_id}/graphs/{id}"
				deleteGraphWaitingProduct = "ges"
			)
			deleteGraphWaitingClient, err := cfg.NewServiceClient(deleteGraphWaitingProduct, region)
			if err != nil {
				return nil, "ERROR", fmt.Errorf("error creating GES Client: %s", err)
			}

			deleteGraphWaitingPath := deleteGraphWaitingClient.Endpoint + deleteGraphWaitingHttpUrl
			deleteGraphWaitingPath = strings.ReplaceAll(deleteGraphWaitingPath, "{project_id}", deleteGraphWaitingClient.ProjectID)
			deleteGraphWaitingPath = strings.ReplaceAll(deleteGraphWaitingPath, "{id}", d.Id())

			deleteGraphWaitingOpt := golangsdk.RequestOpts{
				KeepResponseBody: true,
				OkCodes: []int{
					200,
				},
				MoreHeaders: map[string]string{"Content-Type": "application/json;charset=UTF-8"},
			}

			deleteGraphWaitingResp, err := deleteGraphWaitingClient.Request("GET", deleteGraphWaitingPath, &deleteGraphWaitingOpt)
			if err != nil {
				if _, ok := err.(golangsdk.ErrDefault404); ok {
					// When the error code is 404, the value of respBody is nil, and a non-null value is returned to avoid continuing the loop check.
					return "Resource Not Found", "COMPLETED", nil
				}

				return nil, "ERROR", err
			}

			deleteGraphWaitingRespBody, err := utils.FlattenResponse(deleteGraphWaitingResp)
			if err != nil {
				return nil, "ERROR", err
			}
			statusRaw := utils.PathSearch(`graph.status`, deleteGraphWaitingRespBody, nil)

			status := fmt.Sprintf("%v", statusRaw)

			targetStatus := []string{
				"400",
			}
			if utils.StrSliceContains(targetStatus, status) {
				return deleteGraphWaitingRespBody, "COMPLETED", nil
			}

			unexpectedStatus := []string{
				"300",
				"303",
				"800",
			}
			if utils.StrSliceContains(unexpectedStatus, status) {
				return deleteGraphWaitingRespBody, status, nil
			}

			return deleteGraphWaitingRespBody, "PENDING", nil
		},
		Timeout:      t,
		Delay:        30 * time.Second,
		PollInterval: 30 * time.Second,
	}
	_, err := stateConf.WaitForStateContext(ctx)
	return err
}
