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

// @API GaussDB POST /v3/{project_id}/instances/{instance_id}/list-full-sqls
func DataSourceSingleFullSqls() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceSingleFullSqlsRead,

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
			"full_sqls": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"instance_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"node_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"component_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"db_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"schema_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"username": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"application_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"client_addr": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"client_port": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"sql_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"sql_exec_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"transaction_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"trace_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"query": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"sql": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"begin_time": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"end_time": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"all_time": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"db_time": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"cpu_time": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"data_io_time": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"execution_time": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"scan_lines": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"insert_rows": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"update_rows": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"delete_rows": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"is_slow_sql": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"start_timestamp": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"finish_timestamp": {
							Type:     schema.TypeInt,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func buildQuerySingleFullSqlsBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"begin_time":         d.Get("begin_time"),
		"end_time":           d.Get("end_time"),
		"node_id":            utils.ValueIgnoreEmpty(d.Get("node_id")),
		"query":              utils.ValueIgnoreEmpty(d.Get("query")), // This parameter not take effect.
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
		"multi_queries":      buildMultiMergeConditionInfo(d.Get("multi_queries").([]interface{})),
		"compare_conditions": buildCompareConditionInfo(d.Get("compare_conditions").([]interface{})),
	}

	return bodyParams
}

func buildMultiMergeConditionInfo(multiQueries []interface{}) []map[string]interface{} {
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

func buildCompareConditionInfo(compareConditions []interface{}) []map[string]interface{} {
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
				log.Printf("[ERROR] error parsing 'is_fuzzy' field to Boolean: %s", err)
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

func convertStringtoInt(s string) int {
	i, err := strconv.Atoi(s)
	if err != nil {
		log.Printf("[ERROR] convert the string %s to int failed.", s)
	}

	return i
}

func dataSourceSingleFullSqlsRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		httpUrl = "v3/{project_id}/instances/{instance_id}/list-full-sqls"
	)

	client, err := cfg.NewServiceClient("opengauss", region)
	if err != nil {
		return diag.Errorf("error creating GaussDB client: %s", err)
	}

	getPath := client.Endpoint + httpUrl
	getPath = strings.ReplaceAll(getPath, "{project_id}", client.ProjectID)
	getPath = strings.ReplaceAll(getPath, "{instance_id}", d.Get("instance_id").(string))
	getOpt := golangsdk.RequestOpts{
		MoreHeaders: map[string]string{
			"Content-Type": "application/json",
			"X-Language":   "en-us",
		},
		KeepResponseBody: true,
		JSONBody:         utils.RemoveNil(buildQuerySingleFullSqlsBodyParams(d)),
	}

	getResp, err := client.Request("POST", getPath, &getOpt)
	if err != nil {
		return diag.Errorf("error retrieving the GaussDB full data of a single SQL statement: %s", err)
	}

	getRespBody, err := utils.FlattenResponse(getResp)
	if err != nil {
		return diag.FromErr(err)
	}

	randomUUID, err := uuid.NewRandom()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}

	d.SetId(randomUUID.String())

	mErr := multierror.Append(
		d.Set("region", region),
		d.Set("full_sqls", flattenSingleFullSqls(
			utils.PathSearch("full_sqls", getRespBody, make([]interface{}, 0)).([]interface{}))),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenSingleFullSqls(sqlInfos []interface{}) []interface{} {
	if len(sqlInfos) == 0 {
		return nil
	}

	rst := make([]interface{}, 0, len(sqlInfos))
	for _, v := range sqlInfos {
		rst = append(rst, map[string]interface{}{
			"id":               utils.PathSearch("id", v, nil),
			"instance_id":      utils.PathSearch("instance_id", v, nil),
			"node_id":          utils.PathSearch("node_id", v, nil),
			"component_id":     utils.PathSearch("component_id", v, nil),
			"db_name":          utils.PathSearch("db_name", v, nil),
			"schema_name":      utils.PathSearch("schema_name", v, nil),
			"username":         utils.PathSearch("username", v, nil),
			"application_name": utils.PathSearch("application_name", v, nil),
			"client_addr":      utils.PathSearch("client_addr", v, nil),
			"client_port":      utils.PathSearch("client_port", v, nil),
			"sql_id":           utils.PathSearch("sql_id", v, nil),
			"sql_exec_id":      utils.PathSearch("sql_exec_id", v, nil),
			"transaction_id":   utils.PathSearch("transaction_id", v, nil),
			"trace_id":         utils.PathSearch("trace_id", v, nil),
			"query":            utils.PathSearch("query", v, nil),
			"sql":              utils.PathSearch("sql", v, nil),
			"begin_time":       utils.PathSearch("begin_time", v, nil),
			"end_time":         utils.PathSearch("end_time", v, nil),
			"all_time":         utils.PathSearch("all_time", v, nil),
			"db_time":          utils.PathSearch("db_time", v, nil),
			"cpu_time":         utils.PathSearch("cpu_time", v, nil),
			"data_io_time":     utils.PathSearch("data_io_time", v, nil),
			"execution_time":   utils.PathSearch("execution_time", v, nil),
			"scan_lines":       utils.PathSearch("scan_lines", v, nil),
			"insert_rows":      utils.PathSearch("insert_rows", v, nil),
			"update_rows":      utils.PathSearch("update_rows", v, nil),
			"delete_rows":      utils.PathSearch("delete_rows", v, nil),
			"is_slow_sql":      utils.PathSearch("is_slow_sql", v, nil),
			"start_timestamp":  utils.PathSearch("start_timestamp", v, nil),
			"finish_timestamp": utils.PathSearch("finish_timestamp", v, nil),
		})
	}

	return rst
}
