package drs

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API DRS GET /v5/{project_id}/jobs/{job_id}/filter-rules/result
func DataSourceDrsFilterRulesResult() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceDrsFilterRulesResultRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Specifies the region in which to query the data source. If omitted, the provider-level region will be used.",
			},
			"job_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Specifies the ID of the DRS job to query filter rules result.",
			},
			"query_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Specifies the request ID used to associate the previous filter rules query request.",
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
										Type:        schema.TypeString,
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

func dataSourceDrsFilterRulesResultRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	var (
		httpUrl = "v5/{project_id}/jobs/{job_id}/filter-rules/result"
		product = "drs"
	)

	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating DRS client: %s", err)
	}

	getPath := client.Endpoint + httpUrl
	getPath = strings.ReplaceAll(getPath, "{project_id}", client.ProjectID)
	getPath = strings.ReplaceAll(getPath, "{job_id}", d.Get("job_id").(string))
	getPath += buildListFilterRulesResultQueryParams(d)

	getOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"Content-Type": "application/json"},
	}

	getResp, err := client.Request("GET", getPath, &getOpt)
	if err != nil {
		return diag.Errorf("error retrieving DRS filter rules result: %s", err)
	}

	getRespBody, err := utils.FlattenResponse(getResp)
	if err != nil {
		return diag.FromErr(err)
	}

	randomUUID, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(randomUUID)

	mErr := multierror.Append(
		d.Set("status", utils.PathSearch("status", getRespBody, nil)),
		d.Set("total_count", utils.PathSearch("total_count", getRespBody, nil)),
		d.Set("rules", flattenFilterRulesResult(getRespBody)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func buildListFilterRulesResultQueryParams(d *schema.ResourceData) string {
	res := ""
	if v, ok := d.GetOk("query_id"); ok {
		res = fmt.Sprintf("%s&query_id=%v", res, v)
	}

	if res != "" {
		res = "?" + res[1:]
	}
	return res
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
