package drs

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API DRS PUT /v5/{project_id}/jobs/{job_id}/filter-rules
// @API DRS DELETE /v5/{project_id}/jobs/{job_id}/filter-rules
// @API DRS POST /v5/{project_id}/jobs/{job_id}/check-filter-rules
// @API DRS POST /v5/{project_id}/jobs/{job_id}/filter-rules
// @API DRS GET /v5/{project_id}/jobs/{job_id}/filter-rules/result
func ResourceFilterRules() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceFilterRulesCreate,
		ReadContext:   resourceFilterRulesRead,
		UpdateContext: resourceFilterRulesUpdate,
		DeleteContext: resourceFilterRulesDelete,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"job_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"data_filter_type": {
				Type:         schema.TypeString,
				Optional:     true,
				Default:      "common",
				ValidateFunc: validation.StringInSlice([]string{"common", "config"}, false),
			},
			"source": {
				Type:         schema.TypeString,
				Optional:     true,
				Default:      "job",
				ValidateFunc: validation.StringInSlice([]string{"job", "compare"}, false),
			},
			"general_filtering_list": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"filter_object_list": {
							Type:     schema.TypeList,
							Optional: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"id": {
										Type:     schema.TypeString,
										Required: true,
									},
									"parent_id": {
										Type:     schema.TypeString,
										Required: true,
									},
									"object_name": {
										Type:     schema.TypeString,
										Required: true,
									},
									"object_alias_name": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"object_type": {
										Type:     schema.TypeString,
										Optional: true,
									},
								},
							},
						},
						"source_filter_condition": {
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
			},
			"advanced_setting_list": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"db_name": {
							Type:     schema.TypeString,
							Required: true,
						},
						"table_name": {
							Type:     schema.TypeString,
							Required: true,
						},
						"col_names": {
							Type:     schema.TypeString,
							Required: true,
						},
						"prim_key_or_indexes": {
							Type:     schema.TypeString,
							Required: true,
						},
						"indexes": {
							Type:     schema.TypeString,
							Required: true,
						},
						"values": {
							Type:     schema.TypeString,
							Required: true,
						},
					},
				},
			},
			"query_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The ID of the check task submitted by the filter rules check request.",
			},
		},
	}
}

func buildFilterRulesBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"general_filtering_list": buildGeneralFilteringListBodyParams(d.Get("general_filtering_list")),
		"advanced_setting_list":  buildAdvancedSettingListBodyParams(d.Get("advanced_setting_list")),
		"source":                 utils.ValueIgnoreEmpty(d.Get("source")),
		"data_filter_type":       utils.ValueIgnoreEmpty(d.Get("data_filter_type")),
	}
	return bodyParams
}

func buildGeneralFilteringListBodyParams(rawParams interface{}) []interface{} {
	rawArray, ok := rawParams.([]interface{})
	if !ok {
		return nil
	}

	result := make([]interface{}, 0, len(rawArray))
	for _, rawItem := range rawArray {
		item := rawItem.(map[string]interface{})
		result = append(result, map[string]interface{}{
			"filter_object_list":      buildFilterObjectListBodyParams(item["filter_object_list"]),
			"source_filter_condition": utils.ValueIgnoreEmpty(item["source_filter_condition"]),
		})
	}
	return result
}

func buildFilterObjectListBodyParams(rawParams interface{}) []interface{} {
	rawArray, ok := rawParams.([]interface{})
	if !ok {
		return nil
	}

	result := make([]interface{}, 0, len(rawArray))
	for _, rawItem := range rawArray {
		item := rawItem.(map[string]interface{})
		result = append(result, map[string]interface{}{
			"id":                item["id"],
			"parent_id":         item["parent_id"],
			"object_name":       item["object_name"],
			"object_alias_name": utils.ValueIgnoreEmpty(item["object_alias_name"]),
			"object_type":       utils.ValueIgnoreEmpty(item["object_type"]),
		})
	}
	return result
}

func buildAdvancedSettingListBodyParams(rawParams interface{}) []interface{} {
	rawArray, ok := rawParams.([]interface{})
	if !ok {
		return nil
	}

	result := make([]interface{}, 0, len(rawArray))
	for _, rawItem := range rawArray {
		item := rawItem.(map[string]interface{})
		result = append(result, map[string]interface{}{
			"db_name":             item["db_name"],
			"table_name":          item["table_name"],
			"col_names":           item["col_names"],
			"prim_key_or_indexes": item["prim_key_or_indexes"],
			"indexes":             item["indexes"],
			"values":              item["values"],
		})
	}
	return result
}

func buildDeleteFilterRulesBodyParams(ruleIds []int, filterType string) map[string]interface{} {
	return map[string]interface{}{
		"rule_ids": ruleIds,
		"type":     filterType,
	}
}

func buildFilterRulesCheckBodyParams(d *schema.ResourceData) map[string]interface{} {
	return map[string]interface{}{
		"general_filtering_list": buildGeneralFilteringListBodyParams(d.Get("general_filtering_list")),
		"advanced_setting_list":  buildAdvancedSettingListBodyParams(d.Get("advanced_setting_list")),
		"source":                 utils.ValueIgnoreEmpty(d.Get("source")),
		"data_filter_type":       utils.ValueIgnoreEmpty(d.Get("data_filter_type")),
	}
}

func buildFilterRulesCollectBodyParams(d *schema.ResourceData) map[string]interface{} {
	return map[string]interface{}{
		"filter_type":  d.Get("data_filter_type"),
		"is_all_rules": true,
		"limit":        utils.ValueIgnoreEmpty(d.Get("limit")),
		"offset":       utils.ValueIgnoreEmpty(d.Get("offset")),
	}
}

func waitForFilterRulesCollectCompleted(ctx context.Context, client *golangsdk.ServiceClient,
	jobId, queryId string, timeout time.Duration) (interface{}, error) {
	stateConf := &retry.StateChangeConf{
		Pending: []string{"PENDING"},
		Target:  []string{"SUCCESS"},
		Refresh: func() (interface{}, string, error) {
			respBody, err := getFilterRulesCollectResult(client, jobId, queryId)
			if err != nil {
				return nil, "ERROR", err
			}

			status := utils.PathSearch("status", respBody, "").(string)
			if status == "FAILED" {
				return respBody, "FAILED", errors.New("the filter rules collect task failed")
			}

			return respBody, status, nil
		},
		Timeout:      timeout,
		Delay:        5 * time.Second,
		PollInterval: 5 * time.Second,
	}

	resp, err := stateConf.WaitForStateContext(ctx)
	if err != nil {
		return nil, fmt.Errorf("error waiting for DRS filter rules collect (%s) to complete: %s", queryId, err)
	}

	return resp, nil
}

func getFilterRulesCollectResult(client *golangsdk.ServiceClient, jobId, queryId string) (interface{}, error) {
	httpUrl := "v5/{project_id}/jobs/{job_id}/filter-rules/result"
	getPath := client.Endpoint + httpUrl
	getPath = strings.ReplaceAll(getPath, "{project_id}", client.ProjectID)
	getPath = strings.ReplaceAll(getPath, "{job_id}", jobId)
	getPath += fmt.Sprintf("?query_id=%s", queryId)

	getOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		OkCodes:          []int{200, 202},
		MoreHeaders:      map[string]string{"Content-Type": "application/json"},
	}

	getResp, err := client.Request("GET", getPath, &getOpt)
	if err != nil {
		return nil, fmt.Errorf("error retrieving DRS filter rules collect result: %s", err)
	}

	return utils.FlattenResponse(getResp)
}

func resourceFilterRulesCreate(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		httpUrl = "v5/{project_id}/jobs/{job_id}/filter-rules"
	)

	client, err := cfg.DrsV5Client(region)
	if err != nil {
		return diag.Errorf("error creating DRS client: %s", err)
	}

	requestPath := client.Endpoint + httpUrl
	requestPath = strings.ReplaceAll(requestPath, "{project_id}", client.ProjectID)
	requestPath = strings.ReplaceAll(requestPath, "{job_id}", d.Get("job_id").(string))

	requestOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"Content-Type": "application/json"},
		JSONBody:         utils.RemoveNil(buildFilterRulesBodyParams(d)),
	}

	_, err = client.Request("PUT", requestPath, &requestOpt)
	if err != nil {
		return diag.Errorf("error saving DRS data filter rules: %s", err)
	}

	d.SetId(d.Get("job_id").(string))

	if diagErr := checkDrsFilterRules(d, client); diagErr != nil {
		return diagErr
	}

	return nil
}

func resourceFilterRulesRead(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	// The query API for data filter rules is an asynchronous collection task and does not directly return
	// the current resource state, so no processing is performed in the 'Read()' method.
	return nil
}

func resourceFilterRulesUpdate(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		httpUrl = "v5/{project_id}/jobs/{job_id}/filter-rules"
	)

	client, err := cfg.DrsV5Client(region)
	if err != nil {
		return diag.Errorf("error creating DRS client: %s", err)
	}

	requestPath := client.Endpoint + httpUrl
	requestPath = strings.ReplaceAll(requestPath, "{project_id}", client.ProjectID)
	requestPath = strings.ReplaceAll(requestPath, "{job_id}", d.Get("job_id").(string))

	requestOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"Content-Type": "application/json"},
		JSONBody:         utils.RemoveNil(buildFilterRulesBodyParams(d)),
	}

	_, err = client.Request("PUT", requestPath, &requestOpt)
	if err != nil {
		return diag.Errorf("error updating DRS data filter rules: %s", err)
	}

	if diagErr := checkDrsFilterRules(d, client); diagErr != nil {
		return diagErr
	}

	return nil
}

func checkDrsFilterRules(d *schema.ResourceData, client *golangsdk.ServiceClient) diag.Diagnostics {
	httpUrl := "v5/{project_id}/jobs/{job_id}/check-filter-rules"
	checkPath := client.Endpoint + httpUrl
	checkPath = strings.ReplaceAll(checkPath, "{project_id}", client.ProjectID)
	checkPath = strings.ReplaceAll(checkPath, "{job_id}", d.Get("job_id").(string))

	checkOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"Content-Type": "application/json"},
		JSONBody:         utils.RemoveNil(buildFilterRulesCheckBodyParams(d)),
	}

	checkResp, err := client.Request("POST", checkPath, &checkOpt)
	if err != nil {
		return diag.Errorf("error checking DRS data filter rules: %s", err)
	}

	checkRespBody, err := utils.FlattenResponse(checkResp)
	if err != nil {
		return diag.FromErr(err)
	}

	queryId := utils.PathSearch("id", checkRespBody, "").(string)
	if queryId == "" {
		return diag.Errorf("error checking DRS data filter rules: ID is not found in API response")
	}

	if err := d.Set("query_id", queryId); err != nil {
		return diag.Errorf("error setting query_id: %s", err)
	}

	return nil
}

func resourceFilterRulesDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		httpUrl = "v5/{project_id}/jobs/{job_id}/filter-rules"
	)

	client, err := cfg.DrsV5Client(region)
	if err != nil {
		return diag.Errorf("error creating DRS client: %s", err)
	}

	// Query the filter rules first to get the batch IDs, which are required by the delete API.
	ruleIds, err := listFilterRuleBatchIds(ctx, d, client)
	if err != nil {
		return diag.Errorf("error querying DRS data filter rules before deleting: %s", err)
	}
	if len(ruleIds) == 0 {
		// No rules to delete, remove the resource from the state directly.
		return nil
	}

	requestPath := client.Endpoint + httpUrl
	requestPath = strings.ReplaceAll(requestPath, "{project_id}", client.ProjectID)
	requestPath = strings.ReplaceAll(requestPath, "{job_id}", d.Get("job_id").(string))

	requestOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"Content-Type": "application/json"},
		JSONBody:         utils.RemoveNil(buildDeleteFilterRulesBodyParams(ruleIds, d.Get("data_filter_type").(string))),
	}

	_, err = client.Request("DELETE", requestPath, &requestOpt)
	if err != nil {
		return diag.Errorf("error deleting DRS data filter rules: %s", err)
	}

	return nil
}

// listFilterRuleBatchIds queries the filter rules of the job and returns the batch IDs of the rules
// matching the configured data filter type. The batch IDs are generated by the backend after saving
// the rules and are required by the delete API.
func listFilterRuleBatchIds(ctx context.Context, d *schema.ResourceData, client *golangsdk.ServiceClient) ([]int, error) {
	var (
		httpUrl = "v5/{project_id}/jobs/{job_id}/filter-rules"
	)

	collectPath := client.Endpoint + httpUrl
	collectPath = strings.ReplaceAll(collectPath, "{project_id}", client.ProjectID)
	collectPath = strings.ReplaceAll(collectPath, "{job_id}", d.Get("job_id").(string))

	collectOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"Content-Type": "application/json"},
		JSONBody:         utils.RemoveNil(buildFilterRulesCollectBodyParams(d)),
	}

	collectResp, err := client.Request("POST", collectPath, &collectOpt)
	if err != nil {
		return nil, err
	}

	collectRespBody, err := utils.FlattenResponse(collectResp)
	if err != nil {
		return nil, err
	}

	queryId := utils.PathSearch("id", collectRespBody, "").(string)
	if queryId == "" {
		return nil, fmt.Errorf("unable to find the query ID from the collect response")
	}

	respBody, err := waitForFilterRulesCollectCompleted(ctx, client, d.Get("job_id").(string), queryId,
		d.Timeout(schema.TimeoutDelete))
	if err != nil {
		return nil, err
	}

	rules := utils.PathSearch("rules", respBody, make([]interface{}, 0)).([]interface{})
	ruleIds := make([]int, 0, len(rules))
	for _, rule := range rules {
		ruleIds = append(ruleIds, int(utils.PathSearch("batch_id", rule, float64(0)).(float64)))
	}

	return ruleIds, nil
}
