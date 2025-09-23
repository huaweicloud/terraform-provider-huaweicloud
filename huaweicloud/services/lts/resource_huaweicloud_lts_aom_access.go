// ---------------------------------------------------------------
// *** AUTO GENERATED CODE ***
// @Product LTS
// ---------------------------------------------------------------

package lts

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API LTS DELETE /v2/{project_id}/lts/aom-mapping
// @API LTS POST /v2/{project_id}/lts/aom-mapping
// @API LTS PUT /v2/{project_id}/lts/aom-mapping
// @API LTS GET /v2/{project_id}/lts/aom-mapping/{rule_id}
func ResourceAOMAccess() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceAOMAccessCreate,
		UpdateContext: resourceAOMAccessUpdate,
		ReadContext:   resourceAOMAccessRead,
		DeleteContext: resourceAOMAccessDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
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
				Description: `Specifies the rule name.`,
			},
			"cluster_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `Specifies the cluster ID.`,
			},
			"cluster_name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `Specifies the cluster name.`,
			},
			"namespace": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `Specifies the namespace.`,
			},
			"workloads": {
				Type:        schema.TypeList,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Required:    true,
				Description: `Specifies the deployments.`,
			},
			"access_rules": {
				Type:        schema.TypeList,
				Elem:        accessRuleSchema(),
				Required:    true,
				Description: `Specifies the log details.`,
			},
			"container_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `Specifies the container name.`,
			},
		},
	}
}

func accessRuleSchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"file_name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `Specifies the path name.`,
			},
			"log_group_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `Specifies the log group ID.`,
			},
			"log_group_name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `Specifies the log group name.`,
			},
			"log_stream_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `Specifies the log stream ID.`,
			},
			"log_stream_name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `Specifies the log stream name.`,
			},
		},
	}
	return &sc
}

func resourceAOMAccessCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		httpUrl = "v2/{project_id}/lts/aom-mapping?isBatch=false"
		product = "lts"
	)
	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating LTS client: %s", err)
	}

	createPath := client.Endpoint + httpUrl
	createPath = strings.ReplaceAll(createPath, "{project_id}", client.ProjectID)
	createOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		JSONBody:         buildCreateAOMAccessBodyParams(client.ProjectID, d),
	}

	createResp, err := client.Request("POST", createPath, &createOpt)
	if err != nil {
		return diag.Errorf("error creating AOM to LTS log mapping rule: %s", err)
	}

	createRespBody, err := utils.FlattenResponse(createResp)
	if err != nil {
		return diag.FromErr(err)
	}

	ruleId := utils.PathSearch("[0]|rule_id", createRespBody, "").(string)
	if ruleId == "" {
		return diag.Errorf("unable to find the rule ID of the AOM to LTS log mapping from the API response")
	}
	d.SetId(ruleId)

	return resourceAOMAccessRead(ctx, d, meta)
}

func buildCreateAOMAccessBodyParams(projectID string, d *schema.ResourceData) map[string]interface{} {
	return map[string]interface{}{
		"project_id": projectID,
		"rule_name":  d.Get("name"),
		"rule_info":  buildAOMAccessRuleInfo(d),
	}
}

func buildAOMAccessRuleInfo(d *schema.ResourceData) map[string]interface{} {
	rst := map[string]interface{}{
		"cluster_id":   d.Get("cluster_id"),
		"cluster_name": d.Get("cluster_name"),
		"namespace":    d.Get("namespace"),
		"deployments":  d.Get("workloads"),
		"files":        buildAccessRules(d.Get("access_rules")),
	}
	if v, ok := d.GetOk("container_name"); ok {
		rst["container_name"] = v
	}
	return rst
}

func buildAccessRules(rawParams interface{}) []map[string]interface{} {
	rawArray, ok := rawParams.([]interface{})
	if !ok || len(rawArray) == 0 {
		return nil
	}

	rst := make([]map[string]interface{}, 0, len(rawArray))
	for _, v := range rawArray {
		if raw, isMap := v.(map[string]interface{}); isMap {
			rst = append(rst, map[string]interface{}{
				"file_name":       raw["file_name"],
				"log_stream_info": buildLogStreamInfo(raw),
			})
		}
	}
	return rst
}

func buildLogStreamInfo(rawMap map[string]interface{}) map[string]interface{} {
	return map[string]interface{}{
		"target_log_group_id":    rawMap["log_group_id"],
		"target_log_group_name":  rawMap["log_group_name"],
		"target_log_stream_id":   rawMap["log_stream_id"],
		"target_log_stream_name": rawMap["log_stream_name"],
	}
}

func resourceAOMAccessRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		mErr    *multierror.Error
		httpUrl = "v2/{project_id}/lts/aom-mapping/{rule_id}"
		product = "lts"
	)
	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating LTS client: %s", err)
	}

	getPath := client.Endpoint + httpUrl
	getPath = strings.ReplaceAll(getPath, "{project_id}", client.ProjectID)
	getPath = strings.ReplaceAll(getPath, "{rule_id}", d.Id())
	getOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"Content-Type": "application/json;charset=UTF-8"},
	}

	getResp, err := client.Request("GET", getPath, &getOpt)
	if err != nil {
		// there is no special error code here
		return diag.Errorf("error retrieving AOM to LTS log mapping rule: %s", err)
	}

	getRespBody, err := utils.FlattenResponse(getResp)
	if err != nil {
		return diag.FromErr(err)
	}

	resultBody, err := parseGetResponseBody(getRespBody)
	if err != nil {
		return common.CheckDeletedDiag(d, err, "error retrieving AOM to LTS log mapping rule")
	}

	mErr = multierror.Append(
		mErr,
		d.Set("region", region),
		d.Set("name", utils.PathSearch("rule_name", resultBody, nil)),
		d.Set("cluster_name", utils.PathSearch("rule_info.cluster_name", resultBody, nil)),
		d.Set("cluster_id", utils.PathSearch("rule_info.cluster_id", resultBody, nil)),
		d.Set("namespace", utils.PathSearch("rule_info.namespace", resultBody, nil)),
		d.Set("container_name", utils.PathSearch("rule_info.container_name", resultBody, nil)),
		d.Set("workloads", utils.PathSearch("rule_info.deployments", resultBody, nil)),
		d.Set("access_rules", flattenAccessRules(resultBody)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

// parseGetResponseBody use to parse the return value of the details interface.
// The response body of the detail interface is an array.
// If the resource is not exist, the response will be an empty array.
func parseGetResponseBody(getRespBody interface{}) (interface{}, error) {
	arrayBody, ok := getRespBody.([]interface{})
	if !ok {
		return nil, fmt.Errorf("the API response is not array")
	}
	if len(arrayBody) == 0 {
		return nil, golangsdk.ErrDefault404{}
	}
	return arrayBody[0], nil
}

func flattenAccessRules(resp interface{}) []interface{} {
	if resp == nil {
		return nil
	}
	curJson := utils.PathSearch("rule_info.files", resp, make([]interface{}, 0))
	curArray := curJson.([]interface{})
	rst := make([]interface{}, len(curArray))
	for i, v := range curArray {
		rst[i] = map[string]interface{}{
			"file_name":       utils.PathSearch("file_name", v, nil),
			"log_group_id":    utils.PathSearch("log_stream_info.target_log_group_id", v, nil),
			"log_group_name":  utils.PathSearch("log_stream_info.target_log_group_name", v, nil),
			"log_stream_id":   utils.PathSearch("log_stream_info.target_log_stream_id", v, nil),
			"log_stream_name": utils.PathSearch("log_stream_info.target_log_stream_name", v, nil),
		}
	}
	return rst
}

func resourceAOMAccessUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		httpUrl = "v2/{project_id}/lts/aom-mapping"
		product = "lts"
	)
	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating LTS client: %s", err)
	}

	updatePath := client.Endpoint + httpUrl
	updatePath = strings.ReplaceAll(updatePath, "{project_id}", client.ProjectID)
	updateOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		JSONBody:         buildUpdateAOMAccessBodyParams(client.ProjectID, d),
	}

	_, err = client.Request("PUT", updatePath, &updateOpt)
	if err != nil {
		return diag.Errorf("error updating AOM to LTS log mapping rule: %s", err)
	}
	return resourceAOMAccessRead(ctx, d, meta)
}

func buildUpdateAOMAccessBodyParams(projectID string, d *schema.ResourceData) map[string]interface{} {
	return map[string]interface{}{
		"project_id": projectID,
		"rule_id":    d.Id(),
		"rule_name":  d.Get("name"),
		"rule_info":  buildAOMAccessRuleInfo(d),
	}
}

func resourceAOMAccessDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		httpUrl = "v2/{project_id}/lts/aom-mapping"
		product = "lts"
	)
	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating LTS client: %s", err)
	}

	deletePath := client.Endpoint + httpUrl
	deletePath = strings.ReplaceAll(deletePath, "{project_id}", client.ProjectID)
	deletePath += buildDeleteAOMAccessQueryParams(d)
	deleteOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"Content-Type": "application/json;charset=UTF-8"},
	}

	_, err = client.Request("DELETE", deletePath, &deleteOpt)
	if err != nil {
		return diag.Errorf("error deleting AOM to LTS log mapping rule: %s", err)
	}

	return nil
}

func buildDeleteAOMAccessQueryParams(d *schema.ResourceData) string {
	return fmt.Sprintf("?id=%s", d.Id())
}
