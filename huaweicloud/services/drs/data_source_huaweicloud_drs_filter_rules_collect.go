package drs

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API DRS POST /v5/{project_id}/jobs/{job_id}/filter-rules
// @API DRS GET /v5/{project_id}/jobs/{job_id}/filter-rules/result
func DataSourceDrsFilterRulesCollect() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceDrsFilterRulesCollectRead,

		Timeouts: &schema.ResourceTimeout{
			Read: schema.DefaultTimeout(10 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Specifies the region in which to collect the filter rules. If omitted, the provider-level region will be used.",
			},
			"job_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Specifies the ID of the DRS job to collect filter rules from.",
			},
			"filter_type": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringInSlice([]string{"common", "config"}, false),
				Description:  "Specifies the type of the data filter rules to collect.",
			},
			"is_all_rules": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     true,
				Description: "Specifies whether to query all filter rules (including synchronized and unsent ones).",
			},
			"limit": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "Specifies the number of records to return.",
			},
			"offset": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "Specifies the offset of the records to query.",
			},
			"query_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The async query task ID returned by the collect request.",
			},
			"status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The query progress of the filter rules.",
			},
			"total_count": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The total number of the filter rules.",
			},
			"rules": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "The list of the data filter rules.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"batch_id": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "The batch ID of the data filter rule.",
						},
						"filter_object_list": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "The list of the filter objects.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"batch_id": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "The batch ID of the data filter object.",
									},
									"id": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The unique identifier of the data filter object.",
									},
									"parent_id": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The parent object identifier of the data filter object.",
									},
									"object_name": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The name of the data filter object.",
									},
									"object_alias_name": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The alias name of the data filter object.",
									},
									"object_type": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The type of the data filter object.",
									},
									"db_name": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The database name of the data filter object.",
									},
									"schema_name": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The schema name of the data filter object.",
									},
									"table_name": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The table name of the data filter object.",
									},
									"is_data_filter": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Whether the data filter is enabled for the object.",
									},
									"data_filter_type": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The data filter type of the object.",
									},
									"filter_condition": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The filter condition of the object.",
									},
									"is_synchronized": {
										Type:        schema.TypeBool,
										Computed:    true,
										Description: "Whether the data filter condition is synchronized.",
									},
								},
							},
						},
						"filter_rules_resp": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "The filter condition response of the rule.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"source_filter_condition": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The source database filter condition SQL expression.",
									},
									"target_filter_condition": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The target database filter condition SQL expression.",
									},
									"condition_separated": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Whether the source and target database filter conditions are configured separately.",
									},
								},
							},
						},
						"advanced_setting_resp": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "The advanced setting response of the rule.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"db_name": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The database name.",
									},
									"table_name": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The table name.",
									},
									"col_names": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The column names.",
									},
									"prim_key_or_indexes": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The primary key or unique index.",
									},
									"indexes": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The indexes used for optimized queries.",
									},
									"values": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The filter condition.",
									},
								},
							},
						},
					},
				},
			},
		},
	}
}

func dataSourceDrsFilterRulesCollectRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	var (
		httpUrl = "v5/{project_id}/jobs/{job_id}/filter-rules"
		product = "drs"
	)

	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating DRS client: %s", err)
	}

	// Step 1: submit the async collect request and get the query ID.
	collectPath := client.Endpoint + httpUrl
	collectPath = strings.ReplaceAll(collectPath, "{project_id}", client.ProjectID)
	collectPath = strings.ReplaceAll(collectPath, "{job_id}", d.Get("job_id").(string))

	collectOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"Content-Type": "application/json"},
	}
	collectOpt.JSONBody = utils.RemoveNil(buildCollectDrsFilterRulesBodyParams(d))

	collectResp, err := client.Request("POST", collectPath, &collectOpt)
	if err != nil {
		return diag.Errorf("error collecting DRS filter rules: %s", err)
	}

	collectRespBody, err := utils.FlattenResponse(collectResp)
	if err != nil {
		return diag.FromErr(err)
	}

	queryId := utils.PathSearch("id", collectRespBody, "").(string)
	if queryId == "" {
		return diag.Errorf("unable to find the query ID from the collect response")
	}

	// Step 2: poll the result until the async collect task is completed.
	respBody, err := waitForDrsFilterRulesCollectCompleted(ctx, client, d.Get("job_id").(string), queryId, d.Timeout(schema.TimeoutRead))
	if err != nil {
		return diag.FromErr(err)
	}

	randomUUID, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(randomUUID)

	mErr := multierror.Append(
		d.Set("region", region),
		d.Set("query_id", queryId),
		d.Set("status", utils.PathSearch("status", respBody, nil)),
		d.Set("total_count", utils.PathSearch("total_count", respBody, nil)),
		d.Set("rules", flattenFilterRulesResult(respBody)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func buildCollectDrsFilterRulesBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"filter_type":  d.Get("filter_type"),
		"is_all_rules": d.Get("is_all_rules"),
		"limit":        utils.ValueIgnoreEmpty(d.Get("limit")),
		"offset":       utils.ValueIgnoreEmpty(d.Get("offset")),
	}
	return bodyParams
}

func waitForDrsFilterRulesCollectCompleted(ctx context.Context, client *golangsdk.ServiceClient,
	jobId, queryId string, timeout time.Duration) (interface{}, error) {
	stateConf := &retry.StateChangeConf{
		Pending: []string{"PENDING"},
		Target:  []string{"SUCCESS"},
		Refresh: func() (interface{}, string, error) {
			respBody, err := getDrsFilterRulesCollectResult(client, jobId, queryId)
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

func getDrsFilterRulesCollectResult(client *golangsdk.ServiceClient, jobId, queryId string) (interface{}, error) {
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

func flattenFilterRulesResult(resp interface{}) []interface{} {
	curJson := utils.PathSearch("rules", resp, make([]interface{}, 0))
	curArray := curJson.([]interface{})
	res := make([]interface{}, 0, len(curArray))
	for _, v := range curArray {
		res = append(res, map[string]interface{}{
			"batch_id":              utils.PathSearch("batch_id", v, nil),
			"filter_object_list":    flattenFilterObjectList(v),
			"filter_rules_resp":     flattenFilterRulesResp(v),
			"advanced_setting_resp": flattenAdvancedSettingResp(v),
		})
	}
	return res
}

func flattenFilterObjectList(rule interface{}) []interface{} {
	curJson := utils.PathSearch("filter_object_list", rule, make([]interface{}, 0))
	curArray := curJson.([]interface{})
	res := make([]interface{}, 0, len(curArray))
	for _, v := range curArray {
		res = append(res, map[string]interface{}{
			"batch_id":          utils.PathSearch("batch_id", v, nil),
			"id":                utils.PathSearch("id", v, nil),
			"parent_id":         utils.PathSearch("parent_id", v, nil),
			"object_name":       utils.PathSearch("object_name", v, nil),
			"object_alias_name": utils.PathSearch("object_alias_name", v, nil),
			"object_type":       utils.PathSearch("object_type", v, nil),
			"db_name":           utils.PathSearch("db_name", v, nil),
			"schema_name":       utils.PathSearch("schema_name", v, nil),
			"table_name":        utils.PathSearch("table_name", v, nil),
			"is_data_filter":    utils.PathSearch("is_data_filter", v, nil),
			"data_filter_type":  utils.PathSearch("data_filter_type", v, nil),
			"filter_condition":  utils.PathSearch("filter_condition", v, nil),
			"is_synchronized":   utils.PathSearch("is_synchronized", v, nil),
		})
	}
	return res
}

func flattenFilterRulesResp(rule interface{}) []interface{} {
	curJson := utils.PathSearch("filter_rules_resp", rule, nil)
	if curJson == nil {
		return nil
	}
	return []interface{}{
		map[string]interface{}{
			"source_filter_condition": utils.PathSearch("source_filter_condition", curJson, nil),
			"target_filter_condition": utils.PathSearch("target_filter_condition", curJson, nil),
			"condition_separated":     utils.PathSearch("condition_separated", curJson, nil),
		},
	}
}

func flattenAdvancedSettingResp(rule interface{}) []interface{} {
	curJson := utils.PathSearch("advanced_setting_resp", rule, nil)
	if curJson == nil {
		return nil
	}
	return []interface{}{
		map[string]interface{}{
			"db_name":             utils.PathSearch("db_name", curJson, nil),
			"table_name":          utils.PathSearch("table_name", curJson, nil),
			"col_names":           utils.PathSearch("col_names", curJson, nil),
			"prim_key_or_indexes": utils.PathSearch("prim_key_or_indexes", curJson, nil),
			"indexes":             utils.PathSearch("indexes", curJson, nil),
			"values":              utils.PathSearch("values", curJson, nil),
		},
	}
}
