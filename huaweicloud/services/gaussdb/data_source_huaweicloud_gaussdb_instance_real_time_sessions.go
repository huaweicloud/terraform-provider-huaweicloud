package gaussdb

import (
	"context"
	"strings"

	"github.com/google/uuid"
	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API GaussDB POST /v3/{project_id}/instances/{instance_id}/real-time-session
func DataSourceGaussDBInstanceRealTimeSessions() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceGaussDBInstanceRealTimeSessionsRead,

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
			"node_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"component_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"query_info": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem:     realTimeSessionsQueryInfoSchema(),
			},
			"sessions": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     realTimeSessionsSessionSchema(),
			},
		},
	}
}

func realTimeSessionsQueryInfoSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"database_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"client_ip": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"user_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
		},
	}
}

func realTimeSessionsSessionSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"session_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"pid": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"unique_sql_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"database_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"client_ip": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"user_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"wait": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"block_session": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"wait_event": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"state": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"query_runtime": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"query": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"back_end_start": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"transaction_start": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"query_start": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"application_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"exec_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"trans_num": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"rollback_num": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"sql_num": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"client_port": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"query_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"transaction_time_cost": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"trace_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"global_session_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"top_transaction_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"current_transaction_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"xlog_quantity_pretty": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"wait_status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"lwt_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"thread_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"lock_mode": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"parent_session_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"smp_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"lock_tag": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"component_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceGaussDBInstanceRealTimeSessionsRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var mErr *multierror.Error

	client, err := cfg.NewServiceClient("opengauss", region)
	if err != nil {
		return diag.Errorf("error creating GaussDB client: %s", err)
	}

	httpUrl := "v3/{project_id}/instances/{instance_id}/real-time-session"
	getPath := client.Endpoint + httpUrl
	getPath = strings.ReplaceAll(getPath, "{project_id}", client.ProjectID)
	getPath = strings.ReplaceAll(getPath, "{instance_id}", d.Get("instance_id").(string))

	getOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json",
		},
	}
	getOpt.JSONBody = utils.RemoveNil(buildGetRealTimeSessionsBodyParams(d))

	getResp, err := client.Request("POST", getPath, &getOpt)
	if err != nil {
		return diag.Errorf("error retrieving GaussDB instance real time sessions: %s", err)
	}

	getRespBody, err := utils.FlattenResponse(getResp)
	if err != nil {
		return diag.FromErr(err)
	}

	dataSourceId, err := uuid.NewRandom()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(dataSourceId.String())

	mErr = multierror.Append(
		d.Set("region", region),
		d.Set("sessions", flattenGetRealTimeSessionsBody(getRespBody)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func buildGetRealTimeSessionsBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"node_id":      d.Get("node_id").(string),
		"component_id": d.Get("component_id").(string),
	}

	if v, ok := d.GetOk("query_info"); ok {
		bodyParams["query_info"] = buildGetRealTimeSessionsQueryInfoBodyParams(v.([]interface{})[0].(map[string]interface{}))
	}

	return bodyParams
}

func buildGetRealTimeSessionsQueryInfoBodyParams(params map[string]interface{}) map[string]interface{} {
	queryInfo := map[string]interface{}{
		"database_name": utils.ValueIgnoreEmpty(params["database_name"]),
		"client_ip":     utils.ValueIgnoreEmpty(params["client_ip"]),
		"user_name":     utils.ValueIgnoreEmpty(params["user_name"]),
	}
	return queryInfo
}

func flattenGetRealTimeSessionsBody(resp interface{}) []interface{} {
	curJson := utils.PathSearch("sessions", resp, make([]interface{}, 0))
	curArray := curJson.([]interface{})
	res := make([]interface{}, 0, len(curArray))
	for _, v := range curArray {
		res = append(res, map[string]interface{}{
			"session_id":             utils.PathSearch("session_id", v, nil),
			"pid":                    utils.PathSearch("pid", v, nil),
			"unique_sql_id":          utils.PathSearch("unique_sql_id", v, nil),
			"database_name":          utils.PathSearch("database_name", v, nil),
			"client_ip":              utils.PathSearch("client_ip", v, nil),
			"user_name":              utils.PathSearch("user_name", v, nil),
			"wait":                   utils.PathSearch("wait", v, nil),
			"block_session":          utils.PathSearch("block_session", v, nil),
			"wait_event":             utils.PathSearch("wait_event", v, nil),
			"state":                  utils.PathSearch("state", v, nil),
			"query_runtime":          utils.PathSearch("query_runtime", v, nil),
			"query":                  utils.PathSearch("query", v, nil),
			"back_end_start":         utils.PathSearch("back_end_start", v, nil),
			"transaction_start":      utils.PathSearch("transaction_start", v, nil),
			"query_start":            utils.PathSearch("query_start", v, nil),
			"application_name":       utils.PathSearch("application_name", v, nil),
			"exec_time":              utils.PathSearch("exec_time", v, nil),
			"trans_num":              utils.PathSearch("trans_num", v, nil),
			"rollback_num":           utils.PathSearch("rollback_num", v, nil),
			"sql_num":                utils.PathSearch("sql_num", v, nil),
			"client_port":            utils.PathSearch("client_port", v, nil),
			"query_id":               utils.PathSearch("query_id", v, nil),
			"transaction_time_cost":  utils.PathSearch("transaction_time_cost", v, nil),
			"trace_id":               utils.PathSearch("trace_id", v, nil),
			"global_session_id":      utils.PathSearch("global_session_id", v, nil),
			"top_transaction_id":     utils.PathSearch("top_transaction_id", v, nil),
			"current_transaction_id": utils.PathSearch("current_transaction_id", v, nil),
			"xlog_quantity_pretty":   utils.PathSearch("xlog_quantity_pretty", v, nil),
			"wait_status":            utils.PathSearch("wait_status", v, nil),
			"lwt_id":                 utils.PathSearch("lwt_id", v, nil),
			"thread_name":            utils.PathSearch("thread_name", v, nil),
			"lock_mode":              utils.PathSearch("lock_mode", v, nil),
			"parent_session_id":      utils.PathSearch("parent_session_id", v, nil),
			"smp_id":                 utils.PathSearch("smp_id", v, nil),
			"lock_tag":               utils.PathSearch("lock_tag", v, nil),
			"component_name":         utils.PathSearch("component_name", v, nil),
		})
	}
	return res
}
