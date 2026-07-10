package gaussdb

import (
	"context"
	"log"
	"strconv"
	"strings"

	"github.com/google/uuid"
	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API GaussDB POST /v3/{project_id}/instances/{instance_id}/list-full-sql-statistics
func DataSourceFullSqlStatistics() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceFullSqlStatisticsRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"instance_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"begin_time": {
				Type:     schema.TypeString,
				Required: true,
			},
			"end_time": {
				Type:     schema.TypeString,
				Required: true,
			},
			"node_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"component_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"query": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"sql_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"sql_exec_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"transaction_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"trace_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"db_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"schema_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"username": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"client_addr": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"client_port": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"application_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			// This parameter not take effect
			"is_slow_sql": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"order_by": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"order": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"multi_queries": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:     schema.TypeString,
							Required: true,
						},
						"condition": {
							Type:     schema.TypeString,
							Required: true,
						},
						"values": {
							Type:     schema.TypeList,
							Required: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
						"is_fuzzy": {
							Type:         schema.TypeString,
							Optional:     true,
							ValidateFunc: validation.StringInSlice([]string{"true", "false"}, false),
						},
					},
				},
			},
			"compare_conditions": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:     schema.TypeString,
							Required: true,
						},
						"enable_equal": {
							Type:         schema.TypeString,
							Optional:     true,
							ValidateFunc: validation.StringInSlice([]string{"true", "false"}, false),
						},
						"min": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"max": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"value": {
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
			},
			"statistics": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     fullSqlStatisticsSchema(),
			},
		},
	}
}

func fullSqlStatisticsSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"template": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"sql_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"sql_count": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"total_sql_time": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"avg_sql_time": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"total_db_time": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"avg_db_time": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"total_cpu_time": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"avg_cpu_time": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"avg_execution_time": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"avg_parse_time": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"avg_plan_time": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"total_data_io_time": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"avg_data_io_time": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"avg_n_blocks_hit": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"avg_n_returned_rows": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"avg_n_tuples_fetched": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"start_time_stamp": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"end_time_stamp": {
				Type:     schema.TypeInt,
				Computed: true,
			},
		},
	}
}

func buildQueryFullSqlStatisticsBodyParams(d *schema.ResourceData, limit int) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"limit":              limit,
		"begin_time":         d.Get("begin_time"),
		"end_time":           d.Get("end_time"),
		"node_id":            utils.ValueIgnoreEmpty(d.Get("node_id")),
		"component_id":       utils.ValueIgnoreEmpty(d.Get("component_id")),
		"query":              utils.ValueIgnoreEmpty(d.Get("query")),
		"sql_id":             utils.ValueIgnoreEmpty(d.Get("sql_id")),
		"sql_exec_id":        utils.ValueIgnoreEmpty(d.Get("sql_exec_id")),
		"transaction_id":     utils.ValueIgnoreEmpty(d.Get("transaction_id")),
		"trace_id":           utils.ValueIgnoreEmpty(d.Get("trace_id")),
		"db_name":            utils.ValueIgnoreEmpty(d.Get("db_name")),
		"schema_name":        utils.ValueIgnoreEmpty(d.Get("schema_name")),
		"username":           utils.ValueIgnoreEmpty(d.Get("username")),
		"client_addr":        utils.ValueIgnoreEmpty(d.Get("client_addr")),
		"client_port":        utils.ValueIgnoreEmpty(d.Get("client_port")),
		"is_slow_sql":        utils.GetNestedObjectFromRawConfig(d.GetRawConfig(), "is_slow_sql"),
		"order_by":           utils.ValueIgnoreEmpty(d.Get("order_by")),
		"order":              utils.ValueIgnoreEmpty(d.Get("order")),
		"multi_queries":      buildFullSqlStatisticsMultiQueriesInfo(d.Get("multi_queries").([]interface{})),
		"compare_conditions": buildFullSqlStatisticsCompareConditionsInfo(d.Get("compare_conditions").([]interface{})),
	}

	return bodyParams
}

func buildFullSqlStatisticsMultiQueriesInfo(multiQueries []interface{}) []map[string]interface{} {
	if len(multiQueries) == 0 {
		return nil
	}

	multiQueryInfos := make([]map[string]interface{}, 0, len(multiQueries))
	for _, v := range multiQueries {
		raw, ok := v.(map[string]interface{})
		if !ok {
			return nil
		}

		params := map[string]interface{}{
			"name":      raw["name"],
			"condition": raw["condition"],
			"values":    utils.ExpandToStringList(raw["values"].([]interface{})),
		}

		if v, ok := raw["is_fuzzy"]; ok && v.(string) != "" {
			isFuzzy, err := strconv.ParseBool(v.(string))
			if err != nil {
				log.Printf("[ERROR] error parsing 'is_fuzzy' field to Boolean: %s", err)
			}
			params["is_fuzzy"] = isFuzzy
		}

		multiQueryInfos = append(multiQueryInfos, params)
	}

	return multiQueryInfos
}

func buildFullSqlStatisticsCompareConditionsInfo(compareConditions []interface{}) []map[string]interface{} {
	if len(compareConditions) == 0 {
		return nil
	}

	compareInfos := make([]map[string]interface{}, 0, len(compareConditions))
	for _, v := range compareConditions {
		raw, ok := v.(map[string]interface{})
		if !ok {
			return nil
		}

		params := map[string]interface{}{
			"name": raw["name"],
		}

		if v, ok := raw["enable_equal"]; ok && v.(string) != "" {
			enableEqual, err := strconv.ParseBool(v.(string))
			if err != nil {
				log.Printf("[ERROR] error parsing 'enable_equal' field to Boolean: %s", err)
			}
			params["enable_equal"] = enableEqual
		}

		comMix := utils.PathSearch("min", raw, "").(string)
		comMax := utils.PathSearch("max", raw, "").(string)
		comValue := utils.PathSearch("value", raw, "").(string)
		if comMix != "" {
			params["min"] = convertStringtoInt(comMix)
		}
		if comMax != "" {
			params["max"] = convertStringtoInt(comMax)
		}
		if comValue != "" {
			params["value"] = convertStringtoInt(comValue)
		}

		compareInfos = append(compareInfos, params)
	}

	return compareInfos
}

func dataSourceFullSqlStatisticsRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		limit   = 1000
		offset  = 0
		result  = make([]interface{}, 0)
		httpUrl = "v3/{project_id}/instances/{instance_id}/list-full-sql-statistics"
	)

	client, err := cfg.NewServiceClient("opengauss", region)
	if err != nil {
		return diag.Errorf("error creating GaussDB client: %s", err)
	}

	listPath := client.Endpoint + httpUrl
	listPath = strings.ReplaceAll(listPath, "{project_id}", client.ProjectID)
	listPath = strings.ReplaceAll(listPath, "{instance_id}", d.Get("instance_id").(string))
	bodyParams := utils.RemoveNil(buildQueryFullSqlStatisticsBodyParams(d, limit))
	listOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json",
			"X-Language":   "en-us",
		},
	}

	for {
		bodyParams["offset"] = offset
		listOpt.JSONBody = bodyParams

		resp, err := client.Request("POST", listPath, &listOpt)
		if err != nil {
			return diag.Errorf("error retrieving GaussDB instance full SQL statistics: %s", err)
		}

		respBody, err := utils.FlattenResponse(resp)
		if err != nil {
			return diag.FromErr(err)
		}

		statistics := utils.PathSearch("statistics", respBody, make([]interface{}, 0)).([]interface{})
		result = append(result, statistics...)

		if len(statistics) < limit {
			break
		}

		offset += len(statistics)
	}

	randomUUID, err := uuid.NewRandom()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}

	d.SetId(randomUUID.String())

	mErr := multierror.Append(
		d.Set("region", region),
		d.Set("statistics", flattenFullSqlStatistics(result)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenFullSqlStatistics(resp []interface{}) []interface{} {
	if len(resp) == 0 {
		return nil
	}

	result := make([]interface{}, 0, len(resp))
	for _, v := range resp {
		result = append(result, map[string]interface{}{
			"template":             utils.PathSearch("template", v, nil),
			"sql_id":               utils.PathSearch("sql_id", v, nil),
			"sql_count":            utils.PathSearch("sql_count", v, nil),
			"total_sql_time":       utils.PathSearch("total_sql_time", v, nil),
			"avg_sql_time":         utils.PathSearch("avg_sql_time", v, nil),
			"total_db_time":        utils.PathSearch("total_db_time", v, nil),
			"avg_db_time":          utils.PathSearch("avg_db_time", v, nil),
			"total_cpu_time":       utils.PathSearch("total_cpu_time", v, nil),
			"avg_cpu_time":         utils.PathSearch("avg_cpu_time", v, nil),
			"avg_execution_time":   utils.PathSearch("avg_execution_time", v, nil),
			"avg_parse_time":       utils.PathSearch("avg_parse_time", v, nil),
			"avg_plan_time":        utils.PathSearch("avg_plan_time", v, nil),
			"total_data_io_time":   utils.PathSearch("total_data_io_time", v, nil),
			"avg_data_io_time":     utils.PathSearch("avg_data_io_time", v, nil),
			"avg_n_blocks_hit":     utils.PathSearch("avg_n_blocks_hit", v, nil),
			"avg_n_returned_rows":  utils.PathSearch("avg_n_returned_rows", v, nil),
			"avg_n_tuples_fetched": utils.PathSearch("avg_n_tuples_fetched", v, nil),
			"start_time_stamp":     utils.PathSearch("start_time_stamp", v, nil),
			"end_time_stamp":       utils.PathSearch("end_time_stamp", v, nil),
		})
	}

	return result
}
