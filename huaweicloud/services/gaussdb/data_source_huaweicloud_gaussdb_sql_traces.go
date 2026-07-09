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

// @API GaussDB GET /v3/{project_id}/instances/{instance_id}/full-sql/sql-trace
func DataSourceSqlTraces() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceSqlTracesRead,

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
			"traces": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     sqlTracesSchema(),
			},
		},
	}
}

func sqlTracesSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"component_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"node_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"transaction_id": {
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
			"db_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"schema_name": {
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
			"all_time": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"user_name": {
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
			"trace_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"application_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"session_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"is_slow_sql": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"execution_time_details": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     sqlTracesExecutionTimeDetailsSchema(),
			},
		},
	}
}

func sqlTracesExecutionTimeDetailsSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"resource_time": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     sqlTracesResourceTimeSchema(),
			},
			"kernel_time": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     sqlTracesKernelTimeSchema(),
			},
			"kernel_execution_time": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     sqlTracesKernelExecutionTimeSchema(),
			},
			"wait_event_time": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     sqlTracesWaitEventTimeSchema(),
			},
		},
	}
}

func sqlTracesResourceTimeSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"all_time": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"resource_time_details": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"cpu_time": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"data_io_time": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"other_time": {
							Type:     schema.TypeInt,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func sqlTracesKernelTimeSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"all_time": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"kernel_time_details": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"parse_time": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"rewrite_time": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"plan_time": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"execution_time": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"other_time": {
							Type:     schema.TypeInt,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func sqlTracesKernelExecutionTimeSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"all_time": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"kernel_execution_time_details": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"execution_time": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"other_time": {
							Type:     schema.TypeInt,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func sqlTracesWaitEventTimeSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"code_wait_event_time": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"all_time": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"code_wait_event_time_details": {
							Type:     schema.TypeList,
							Computed: true,
							Elem:     eventTimeDetailsSchema(),
						},
					},
				},
			},
			"resource_wait_event_time": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"all_time": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"resource_wait_event_time_details": {
							Type:     schema.TypeList,
							Computed: true,
							Elem:     resourceWaitEventTimeDetailsSchema(),
						},
						"other_time": {
							Type:     schema.TypeInt,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func resourceWaitEventTimeDetailsSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"data_io_time": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"all_time": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"data_io_time_details": {
							Type:     schema.TypeList,
							Computed: true,
							Elem:     eventTimeDetailsSchema(),
						},
					},
				},
			},
			"lock_time": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"all_time": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"lock_time_details": {
							Type:     schema.TypeList,
							Computed: true,
							Elem:     eventTimeDetailsSchema(),
						},
					},
				},
			},
			"lwlock_time": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"all_time": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"lwlock_time_details": {
							Type:     schema.TypeList,
							Computed: true,
							Elem:     eventTimeDetailsSchema(),
						},
					},
				},
			},
		},
	}
}

func eventTimeDetailsSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"events": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"event_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"event_time": {
							Type:     schema.TypeInt,
							Computed: true,
						},
					},
				},
			},
			"left_time": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"other_time": {
				Type:     schema.TypeInt,
				Computed: true,
			},
		},
	}
}

func buildSqlTracesQueryParams(d *schema.ResourceData) string {
	res := ""

	if v, ok := d.GetOk("sql_id"); ok {
		res = fmt.Sprintf("%s&sql_id=%v", res, v)
	}
	if v, ok := d.GetOk("sql_exec_id"); ok {
		res = fmt.Sprintf("%s&sql_exec_id=%v", res, v)
	}
	if v, ok := d.GetOk("transaction_id"); ok {
		res = fmt.Sprintf("%s&transaction_id=%v", res, v)
	}
	if v, ok := d.GetOk("trace_id"); ok {
		res = fmt.Sprintf("%s&trace_id=%v", res, v)
	}
	if res != "" {
		res = "?" + res[1:]
	}

	return res
}

func dataSourceSqlTracesRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		httpUrl = "v3/{project_id}/instances/{instance_id}/full-sql/sql-trace"
	)

	client, err := cfg.NewServiceClient("opengauss", region)
	if err != nil {
		return diag.Errorf("error creating GaussDB client: %s", err)
	}

	getPath := client.Endpoint + httpUrl
	getPath = strings.ReplaceAll(getPath, "{project_id}", client.ProjectID)
	getPath = strings.ReplaceAll(getPath, "{instance_id}", d.Get("instance_id").(string))
	getPath += buildSqlTracesQueryParams(d)

	getOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json",
			"X-Language":   "en-us",
		},
	}

	getResp, err := client.Request("GET", getPath, &getOpt)
	if err != nil {
		return diag.Errorf("error retrieving GaussDB instance SQL traces: %s", err)
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
		d.Set("traces", flattenSqlTraces(
			utils.PathSearch("[*]", getRespBody, make([]interface{}, 0)).([]interface{}))),
	)
	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenSqlTraces(resp []interface{}) []interface{} {
	if len(resp) == 0 {
		return nil
	}

	result := make([]interface{}, 0, len(resp))
	for _, v := range resp {
		result = append(result, map[string]interface{}{
			"component_id":           utils.PathSearch("component_id", v, nil),
			"node_id":                utils.PathSearch("node_id", v, nil),
			"transaction_id":         utils.PathSearch("transaction_id", v, nil),
			"sql_id":                 utils.PathSearch("sql_id", v, nil),
			"sql_exec_id":            utils.PathSearch("sql_exec_id", v, nil),
			"db_name":                utils.PathSearch("db_name", v, nil),
			"schema_name":            utils.PathSearch("schema_name", v, nil),
			"start_time":             utils.PathSearch("start_time", v, nil),
			"finish_time":            utils.PathSearch("finish_time", v, nil),
			"all_time":               utils.PathSearch("all_time", v, nil),
			"user_name":              utils.PathSearch("user_name", v, nil),
			"client_addr":            utils.PathSearch("client_addr", v, nil),
			"client_port":            utils.PathSearch("client_port", v, nil),
			"trace_id":               utils.PathSearch("trace_id", v, nil),
			"application_name":       utils.PathSearch("application_name", v, nil),
			"session_id":             utils.PathSearch("session_id", v, nil),
			"is_slow_sql":            utils.PathSearch("is_slow_sql", v, nil),
			"execution_time_details": flattenExecutionTimeDetails(utils.PathSearch("execution_time_details", v, nil)),
		})
	}

	return result
}

func flattenExecutionTimeDetails(resp interface{}) []interface{} {
	if resp == nil {
		return nil
	}

	result := map[string]interface{}{
		"resource_time":         flattenResourceTime(utils.PathSearch("resource_time", resp, nil)),
		"kernel_time":           flattenKernelTime(utils.PathSearch("kernel_time", resp, nil)),
		"kernel_execution_time": flattenKernelExecutionTime(utils.PathSearch("kernel_execution_time", resp, nil)),
		"wait_event_time":       flattenWaitEventTime(utils.PathSearch("wait_event_time", resp, nil)),
	}

	return []interface{}{result}
}

func flattenResourceTime(resp interface{}) []interface{} {
	if resp == nil {
		return nil
	}

	result := map[string]interface{}{
		"all_time":              utils.PathSearch("all_time", resp, nil),
		"resource_time_details": flattenResourceTimeDetails(utils.PathSearch("resource_time_details", resp, nil)),
	}

	return []interface{}{result}
}

func flattenResourceTimeDetails(resp interface{}) []interface{} {
	if resp == nil {
		return nil
	}

	result := map[string]interface{}{
		"cpu_time":     utils.PathSearch("cpu_time", resp, nil),
		"data_io_time": utils.PathSearch("data_io_time", resp, nil),
		"other_time":   utils.PathSearch("other_time", resp, nil),
	}

	return []interface{}{result}
}

func flattenKernelTime(resp interface{}) []interface{} {
	if resp == nil {
		return nil
	}

	result := map[string]interface{}{
		"all_time":            utils.PathSearch("all_time", resp, nil),
		"kernel_time_details": flattenKernelTimeDetails(utils.PathSearch("kernel_time_details", resp, nil)),
	}

	return []interface{}{result}
}

func flattenKernelTimeDetails(resp interface{}) []interface{} {
	if resp == nil {
		return nil
	}

	result := map[string]interface{}{
		"parse_time":     utils.PathSearch("parse_time", resp, nil),
		"rewrite_time":   utils.PathSearch("rewrite_time", resp, nil),
		"plan_time":      utils.PathSearch("plan_time", resp, nil),
		"execution_time": utils.PathSearch("execution_time", resp, nil),
		"other_time":     utils.PathSearch("other_time", resp, nil),
	}

	return []interface{}{result}
}

func flattenKernelExecutionTime(resp interface{}) []interface{} {
	if resp == nil {
		return nil
	}

	result := map[string]interface{}{
		"all_time": utils.PathSearch("all_time", resp, nil),
		"kernel_execution_time_details": flattenKernelExecutionTimeDetails(
			utils.PathSearch("kernel_execution_time_details", resp, nil)),
	}

	return []interface{}{result}
}

func flattenKernelExecutionTimeDetails(resp interface{}) []interface{} {
	if resp == nil {
		return nil
	}

	result := map[string]interface{}{
		"execution_time": utils.PathSearch("execution_time", resp, nil),
		"other_time":     utils.PathSearch("other_time", resp, nil),
	}

	return []interface{}{result}
}

func flattenWaitEventTime(resp interface{}) []interface{} {
	if resp == nil {
		return nil
	}

	result := map[string]interface{}{
		"code_wait_event_time":     flattenCodeWaitEventTime(utils.PathSearch("code_wait_event_time", resp, nil)),
		"resource_wait_event_time": flattenResourceWaitEventTime(utils.PathSearch("resource_wait_event_time", resp, nil)),
	}

	return []interface{}{result}
}

func flattenCodeWaitEventTime(resp interface{}) []interface{} {
	if resp == nil {
		return nil
	}

	result := map[string]interface{}{
		"all_time": utils.PathSearch("all_time", resp, nil),
		"code_wait_event_time_details": flattenEventTimeInfo(
			utils.PathSearch("code_wait_event_time_details", resp, nil)),
	}

	return []interface{}{result}
}

func flattenEventTimeInfo(resp interface{}) []interface{} {
	if resp == nil {
		return nil
	}

	result := map[string]interface{}{
		"events":     flattenTopEventsInfo(utils.PathSearch("events", resp, make([]interface{}, 0)).([]interface{})),
		"left_time":  utils.PathSearch("left_time", resp, nil),
		"other_time": utils.PathSearch("other_time", resp, nil),
	}

	return []interface{}{result}
}

func flattenTopEventsInfo(resp []interface{}) []interface{} {
	if len(resp) == 0 {
		return nil
	}

	result := make([]interface{}, 0, len(resp))
	for _, v := range resp {
		result = append(result, map[string]interface{}{
			"event_name": utils.PathSearch("event_name", v, nil),
			"event_time": utils.PathSearch("event_time", v, nil),
		})
	}

	return result
}

func flattenResourceWaitEventTime(resp interface{}) []interface{} {
	if resp == nil {
		return nil
	}

	result := map[string]interface{}{
		"all_time":   utils.PathSearch("all_time", resp, nil),
		"other_time": utils.PathSearch("other_time", resp, nil),
		"resource_wait_event_time_details": flattenResourceWaitEventTimeDetails(
			utils.PathSearch("resource_wait_event_time_details", resp, nil)),
	}

	return []interface{}{result}
}

func flattenResourceWaitEventTimeDetails(resp interface{}) []interface{} {
	if resp == nil {
		return nil
	}

	result := map[string]interface{}{
		"data_io_time": flattenDataIoTime(utils.PathSearch("data_io_time", resp, nil)),
		"lock_time":    flattenLockTime(utils.PathSearch("lock_time", resp, nil)),
		"lwlock_time":  flattenLwlockTime(utils.PathSearch("lwlock_time", resp, nil)),
	}

	return []interface{}{result}
}

func flattenDataIoTime(resp interface{}) []interface{} {
	if resp == nil {
		return nil
	}

	result := map[string]interface{}{
		"all_time":             utils.PathSearch("all_time", resp, nil),
		"data_io_time_details": flattenEventTimeInfo(utils.PathSearch("data_io_time_details", resp, nil)),
	}

	return []interface{}{result}
}

func flattenLockTime(resp interface{}) []interface{} {
	if resp == nil {
		return nil
	}

	result := map[string]interface{}{
		"all_time":          utils.PathSearch("all_time", resp, nil),
		"lock_time_details": flattenEventTimeInfo(utils.PathSearch("lock_time_details", resp, nil)),
	}

	return []interface{}{result}
}

func flattenLwlockTime(resp interface{}) []interface{} {
	if resp == nil {
		return nil
	}

	result := map[string]interface{}{
		"all_time":            utils.PathSearch("all_time", resp, nil),
		"lwlock_time_details": flattenEventTimeInfo(utils.PathSearch("lwlock_time_details", resp, nil)),
	}

	return []interface{}{result}
}
