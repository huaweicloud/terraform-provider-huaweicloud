package gaussdb

import (
	"context"
	"fmt"
	"strings"

	"github.com/google/uuid"
	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API GaussDB GET /v3/{project_id}/instances/{instance_id}/full-sqls/{sql_exec_id}
func DataSourceSingleSqlDetails() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceSingleSqlDetailsRead,

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
			"sql_exec_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"key_id": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"sql_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"trace_statistics": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     singleSqlDetailTraceStatisticsSchema(),
			},
			"components": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     singleSqlDetailComponentsSchema(),
			},
		},
	}
}

func singleSqlDetailTraceStatisticsSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"hit_rate": {
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
			"io_time": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"execution_time": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"scan_rows": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"update_rows": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"insert_rows": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"delete_rows": {
				Type:     schema.TypeInt,
				Computed: true,
			},
		},
	}
}

func singleSqlDetailComponentsSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
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
			"origin_node": {
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
				Type:     schema.TypeString,
				Computed: true,
			},
			"parent_sql_id": {
				Type:     schema.TypeString,
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
			"thread_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"session_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"start_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"finish_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"slow_query_threshold": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"n_soft_parse": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"n_hard_parse": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"query_plan": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"n_returned_rows": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"n_tuples_fetched": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"n_tuples_returned": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"n_tuples_inserted": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"n_tuples_updated": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"n_tuples_deleted": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"n_blocks_fetched": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"n_blocks_hit": {
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
			"execution_time": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"parse_time": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"plan_time": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"rewrite_time": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"pl_execution_time": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"pl_compilation_time": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"data_io_time": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"lock_count": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"lock_time": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"lock_wait_count": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"lock_wait_time": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"details": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"is_slow_sql": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"advise": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"finish_status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"net_send_info": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"net_recv_info": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"net_stream_send_info": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"net_stream_recv_info": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func buildSingleSqlDetailsQueryParams(d *schema.ResourceData) string {
	res := ""

	if v, ok := d.GetOk("key_id"); ok {
		res = fmt.Sprintf("%s&id=%v", res, v)
	}

	if v, ok := d.GetOk("sql_id"); ok {
		res = fmt.Sprintf("%s&sql_id=%v", res, v)
	}

	if res != "" {
		res = "?" + res[1:]
	}

	return res
}

func dataSourceSingleSqlDetailsRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		httpUrl = "v3/{project_id}/instances/{instance_id}/full-sqls/{sql_exec_id}"
	)

	client, err := cfg.NewServiceClient("opengauss", region)
	if err != nil {
		return diag.Errorf("error creating GaussDB client: %s", err)
	}

	getPath := client.Endpoint + httpUrl
	getPath = strings.ReplaceAll(getPath, "{project_id}", client.ProjectID)
	getPath = strings.ReplaceAll(getPath, "{instance_id}", d.Get("instance_id").(string))
	getPath = strings.ReplaceAll(getPath, "{sql_exec_id}", d.Get("sql_exec_id").(string))
	getPath += buildSingleSqlDetailsQueryParams(d)

	getOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json",
			"X-Language":   "en-us",
		},
	}

	getResp, err := client.Request("GET", getPath, &getOpt)
	if err != nil {
		return diag.Errorf("error retrieving GaussDB single SQL statement details: %s", err)
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
		d.Set("trace_statistics", flattenSingleSqlDetailsTraceStatistics(
			utils.PathSearch("trace_statistics", getRespBody, nil))),
		d.Set("components", flattenSingleSqlDetailsComponents(
			utils.PathSearch("components", getRespBody, make([]interface{}, 0)).([]interface{}))),
	)
	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenSingleSqlDetailsTraceStatistics(resp interface{}) []interface{} {
	if resp == nil {
		return nil
	}

	result := map[string]interface{}{
		"hit_rate":       utils.PathSearch("hit_rate", resp, nil),
		"db_time":        utils.PathSearch("db_time", resp, nil),
		"cpu_time":       utils.PathSearch("cpu_time", resp, nil),
		"io_time":        utils.PathSearch("io_time", resp, nil),
		"execution_time": utils.PathSearch("execution_time", resp, nil),
		"scan_rows":      utils.PathSearch("scan_rows", resp, nil),
		"update_rows":    utils.PathSearch("update_rows", resp, nil),
		"insert_rows":    utils.PathSearch("insert_rows", resp, nil),
		"delete_rows":    utils.PathSearch("delete_rows", resp, nil),
	}

	return []interface{}{result}
}

func flattenSingleSqlDetailsComponents(resp []interface{}) []interface{} {
	if len(resp) == 0 {
		return nil
	}

	result := make([]interface{}, 0, len(resp))
	for _, v := range resp {
		result = append(result, map[string]interface{}{
			"component_id":         utils.PathSearch("component_id", v, nil),
			"db_name":              utils.PathSearch("db_name", v, nil),
			"schema_name":          utils.PathSearch("schema_name", v, nil),
			"origin_node":          utils.PathSearch("origin_node", v, nil),
			"username":             utils.PathSearch("username", v, nil),
			"application_name":     utils.PathSearch("application_name", v, nil),
			"client_addr":          utils.PathSearch("client_addr", v, nil),
			"client_port":          utils.PathSearch("client_port", v, nil),
			"parent_sql_id":        utils.PathSearch("parent_sql_id", v, nil),
			"sql_id":               utils.PathSearch("sql_id", v, nil),
			"sql_exec_id":          utils.PathSearch("sql_exec_id", v, nil),
			"transaction_id":       utils.PathSearch("transaction_id", v, nil),
			"trace_id":             utils.PathSearch("trace_id", v, nil),
			"query":                utils.PathSearch("query", v, nil),
			"sql":                  utils.PathSearch("sql", v, nil),
			"thread_id":            utils.PathSearch("thread_id", v, nil),
			"session_id":           utils.PathSearch("session_id", v, nil),
			"start_time":           utils.PathSearch("start_time", v, nil),
			"finish_time":          utils.PathSearch("finish_time", v, nil),
			"slow_query_threshold": utils.PathSearch("slow_query_threshold", v, nil),
			"n_soft_parse":         utils.PathSearch("n_soft_parse", v, nil),
			"n_hard_parse":         utils.PathSearch("n_hard_parse", v, nil),
			"query_plan":           utils.PathSearch("query_plan", v, nil),
			"n_returned_rows":      utils.PathSearch("n_returned_rows", v, nil),
			"n_tuples_fetched":     utils.PathSearch("n_tuples_fetched", v, nil),
			"n_tuples_returned":    utils.PathSearch("n_tuples_returned", v, nil),
			"n_tuples_inserted":    utils.PathSearch("n_tuples_inserted", v, nil),
			"n_tuples_updated":     utils.PathSearch("n_tuples_updated", v, nil),
			"n_tuples_deleted":     utils.PathSearch("n_tuples_deleted", v, nil),
			"n_blocks_fetched":     utils.PathSearch("n_blocks_fetched", v, nil),
			"n_blocks_hit":         utils.PathSearch("n_blocks_hit", v, nil),
			"db_time":              utils.PathSearch("db_time", v, nil),
			"cpu_time":             utils.PathSearch("cpu_time", v, nil),
			"execution_time":       utils.PathSearch("execution_time", v, nil),
			"parse_time":           utils.PathSearch("parse_time", v, nil),
			"plan_time":            utils.PathSearch("plan_time", v, nil),
			"rewrite_time":         utils.PathSearch("rewrite_time", v, nil),
			"pl_execution_time":    utils.PathSearch("pl_execution_time", v, nil),
			"pl_compilation_time":  utils.PathSearch("pl_compilation_time", v, nil),
			"data_io_time":         utils.PathSearch("data_io_time", v, nil),
			"lock_count":           utils.PathSearch("lock_count", v, nil),
			"lock_time":            utils.PathSearch("lock_time", v, nil),
			"lock_wait_count":      utils.PathSearch("lock_wait_count", v, nil),
			"lock_wait_time":       utils.PathSearch("lock_wait_time", v, nil),
			"details":              utils.PathSearch("details", v, nil),
			"is_slow_sql":          utils.PathSearch("is_slow_sql", v, nil),
			"advise":               utils.PathSearch("advise", v, nil),
			"finish_status":        utils.PathSearch("finish_status", v, nil),
			"net_send_info":        utils.PathSearch("net_send_info", v, nil),
			"net_recv_info":        utils.PathSearch("net_recv_info", v, nil),
			"net_stream_send_info": utils.PathSearch("net_stream_send_info", v, nil),
			"net_stream_recv_info": utils.PathSearch("net_stream_recv_info", v, nil),
		})
	}

	return result
}
