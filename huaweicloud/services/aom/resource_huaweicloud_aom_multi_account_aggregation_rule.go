package aom

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

// @API AOM POST /v1/{project_id}/aom/aggr-config
// @API AOM PUT /v1/{project_id}/aom/aggr-config/{metric_aggr_id}
// @API AOM DELETE /v1/{project_id}/aom/aggr-config/{metric_aggr_id}
// @API AOM GET /v1/{project_id}/aom/aggr-config
func ResourceMultiAccountAggregationRule() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceMultiAccountAggregationRuleCreate,
		ReadContext:   resourceMultiAccountAggregationRuleRead,
		UpdateContext: resourceMultiAccountAggregationRuleUpdate,
		DeleteContext: resourceMultiAccountAggregationRuleDelete,

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
			"instance_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"accounts": {
				Type:     schema.TypeSet,
				Required: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Required: true,
						},
						"name": {
							Type:     schema.TypeString,
							Required: true,
						},
						"urn": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"join_method": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"joined_at": {
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
			},
			"services": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"service": {
							Type:     schema.TypeString,
							Required: true,
						},
						"metrics": {
							Type:     schema.TypeSet,
							Required: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
					},
				},
			},
			"send_to_source_account": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},
		},
	}
}

func resourceMultiAccountAggregationRuleCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	client, err := cfg.NewServiceClient("aom", region)
	if err != nil {
		return diag.Errorf("error creating AOM client: %s", err)
	}

	createHttpUrl := "v1/{project_id}/aom/aggr-config"
	createPath := client.Endpoint + createHttpUrl
	createPath = strings.ReplaceAll(createPath, "{project_id}", client.ProjectID)
	createOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json",
		},
		JSONBody: utils.RemoveNil(buildMultiAccountAggregationRuleBodyParams(d)),
	}

	_, err = client.Request("POST", createPath, &createOpt)
	if err != nil {
		return diag.Errorf("error creating AOM multi account aggregation rule: %s", err)
	}

	d.SetId(d.Get("instance_id").(string))

	return resourceMultiAccountAggregationRuleRead(ctx, d, meta)
}

func buildMultiAccountAggregationRuleBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"dest_prometheus_id":     d.Get("instance_id"),
		"accounts":               buildMultiAccountAggregationRuleAccounts(d),
		"metrics":                buildMultiAccountAggregationRuleMetrics(d),
		"send_to_source_account": utils.ValueIgnoreEmpty(d.Get("send_to_source_account")),
	}

	return bodyParams
}

func buildMultiAccountAggregationRuleAccounts(d *schema.ResourceData) []map[string]interface{} {
	rawParams := d.Get("accounts").(*schema.Set).List()
	rst := make([]map[string]interface{}, 0, len(rawParams))
	for _, val := range rawParams {
		raw := val.(map[string]interface{})
		params := map[string]interface{}{
			"id":          raw["id"],
			"name":        raw["name"],
			"urn":         utils.ValueIgnoreEmpty(raw["urn"]),
			"join_method": utils.ValueIgnoreEmpty(raw["join_method"]),
			"joined_at":   utils.ValueIgnoreEmpty(raw["joined_at"]),
		}
		rst = append(rst, params)
	}

	return rst
}

func buildMultiAccountAggregationRuleMetrics(d *schema.ResourceData) map[string]interface{} {
	services := d.Get("services")
	if services == nil {
		return nil
	}
	rawParams := services.(*schema.Set).List()
	if len(rawParams) == 0 {
		return nil
	}

	rst := make(map[string]interface{})
	for _, val := range rawParams {
		raw := val.(map[string]interface{})
		rst[raw["service"].(string)] = raw["metrics"].(*schema.Set).List()
	}

	return rst
}

func resourceMultiAccountAggregationRuleRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	client, err := cfg.NewServiceClient("aom", region)
	if err != nil {
		return diag.Errorf("error creating AOM client: %s", err)
	}

	rule, err := getMultiAccountAggregationRule(client, d)
	if err != nil {
		return common.CheckDeletedDiag(d, err, "error retrieving multi account aggregation rule")
	}

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("instance_id", utils.PathSearch("dest_prometheus_id", rule, nil)),
		d.Set("accounts", flattenMultiAccountAggregationRuleResponseBodyAccounts(rule)),
		d.Set("services", flattenMultiAccountAggregationRuleResponseBodyMetrics(rule)),
		d.Set("send_to_source_account", utils.PathSearch("send_to_source_account", rule, nil)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func getMultiAccountAggregationRule(client *golangsdk.ServiceClient, d *schema.ResourceData) (interface{}, error) {
	listHttpUrl := "v1/{project_id}/aom/aggr-config"
	listPath := client.Endpoint + listHttpUrl
	listPath = strings.ReplaceAll(listPath, "{project_id}", client.ProjectID)
	listOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json",
		},
	}

	listResp, err := client.Request("GET", listPath, &listOpt)
	if err != nil {
		return nil, fmt.Errorf("error retrieving multi account aggregation rule: %s", err)
	}
	listRespBody, err := utils.FlattenResponse(listResp)
	if err != nil {
		return nil, fmt.Errorf("error flattening multi account aggregation rule: %s", err)
	}

	jsonPath := fmt.Sprintf("[?dest_prometheus_id=='%s']|[0]", d.Id())
	rule := utils.PathSearch(jsonPath, listRespBody, nil)
	if rule == nil {
		return nil, golangsdk.ErrDefault404{}
	}

	return rule, nil
}

func flattenMultiAccountAggregationRuleResponseBodyAccounts(resp interface{}) []interface{} {
	curArray := utils.PathSearch("accounts", resp, make([]interface{}, 0)).([]interface{})
	rst := make([]interface{}, 0, len(curArray))
	for _, v := range curArray {
		rst = append(rst, map[string]interface{}{
			"id":          utils.PathSearch("id", v, nil),
			"name":        utils.PathSearch("name", v, nil),
			"urn":         utils.PathSearch("urn", v, nil),
			"join_method": utils.PathSearch("join_method", v, nil),
			"joined_at":   utils.PathSearch("joined_at", v, nil),
		})
	}

	return rst
}

func flattenMultiAccountAggregationRuleResponseBodyMetrics(resp interface{}) []interface{} {
	curArray := utils.PathSearch("metrics", resp, make(map[string]interface{})).(map[string]interface{})
	rst := make([]interface{}, 0, len(curArray))
	for k, v := range curArray {
		rst = append(rst, map[string]interface{}{
			"service": k,
			"metrics": v,
		})
	}

	return rst
}

func resourceMultiAccountAggregationRuleUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	client, err := cfg.NewServiceClient("aom", region)
	if err != nil {
		return diag.Errorf("error creating AOM client: %s", err)
	}

	updateChanges := []string{
		"accounts",
		"services",
		"send_to_source_account",
	}

	if d.HasChanges(updateChanges...) {
		updateHttpUrl := "v1/{project_id}/aom/aggr-config/{metric_aggr_id}"
		updatePath := client.Endpoint + updateHttpUrl
		updatePath = strings.ReplaceAll(updatePath, "{project_id}", client.ProjectID)
		updatePath = strings.ReplaceAll(updatePath, "{metric_aggr_id}", d.Id())
		updateOpt := golangsdk.RequestOpts{
			KeepResponseBody: true,
			MoreHeaders: map[string]string{
				"Content-Type": "application/json",
			},
			JSONBody: utils.RemoveNil(buildMultiAccountAggregationRuleBodyParams(d)),
		}

		_, err = client.Request("PUT", updatePath, &updateOpt)
		if err != nil {
			return diag.Errorf("error updating multi account aggregation rule: %s", err)
		}
	}

	return resourceMultiAccountAggregationRuleRead(ctx, d, meta)
}

func resourceMultiAccountAggregationRuleDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	client, err := cfg.NewServiceClient("aom", region)
	if err != nil {
		return diag.Errorf("error creating AOM client: %s", err)
	}

	// DELETE will return 200 even deleting a non exist rule
	_, err = getMultiAccountAggregationRule(client, d)
	if err != nil {
		return common.CheckDeletedDiag(d, err, "error deleting multi account aggregation rule")
	}

	deleteHttpUrl := "v1/{project_id}/aom/aggr-config/{metric_aggr_id}"
	deletePath := client.Endpoint + deleteHttpUrl
	deletePath = strings.ReplaceAll(deletePath, "{project_id}", client.ProjectID)
	deletePath = strings.ReplaceAll(deletePath, "{metric_aggr_id}", d.Id())
	deleteOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json",
		},
	}

	_, err = client.Request("DELETE", deletePath, &deleteOpt)
	if err != nil {
		return common.CheckDeletedDiag(d, err, "error deleting multi account aggregation rule")
	}

	return nil
}
