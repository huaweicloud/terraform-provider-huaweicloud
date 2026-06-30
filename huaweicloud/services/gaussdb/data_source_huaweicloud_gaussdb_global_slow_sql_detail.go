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

// @API GaussDB POST /v3/{project_id}/instances/{instance_id}/global-slow-sql-detail
func DataSourceGaussDbGlobalSlowSqlDetail() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceGaussDbGlobalSlowSqlDetailRead,

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
			"start_time": {
				Type:     schema.TypeString,
				Required: true,
			},
			"end_time": {
				Type:     schema.TypeString,
				Required: true,
			},
			"sql_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"node_ids": {
				Type:     schema.TypeList,
				Required: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"component_type": {
				Type:     schema.TypeString,
				Required: true,
			},
			"slow_sql_details": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     globalSlowSqlDetailInfosSchema(),
			},
		},
	}
}

func globalSlowSqlDetailInfosSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"db_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"schema_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"sql_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"user_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"client_ip": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"client_port": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"node_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"node_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"sql_text": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"sql": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"query_plan": {
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
			"returned_rows": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"fetched_rows": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"fetched_pages": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"hit_pages": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"total_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"cpu_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"plan_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"io_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"lock_count": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"lock_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceGaussDbGlobalSlowSqlDetailRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var mErr *multierror.Error

	client, err := cfg.NewServiceClient("opengauss", region)
	if err != nil {
		return diag.Errorf("error creating GaussDB client: %s", err)
	}

	httpUrl := "v3/{project_id}/instances/{instance_id}/global-slow-sql-detail"
	getPath := client.Endpoint + httpUrl
	getPath = strings.ReplaceAll(getPath, "{project_id}", client.ProjectID)
	getPath = strings.ReplaceAll(getPath, "{instance_id}", d.Get("instance_id").(string))

	getOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json",
		},
	}

	getOpt.JSONBody = utils.RemoveNil(buildGaussDbGlobalSlowSqlDetailBodyParams(d))
	getResp, err := client.Request("POST", getPath, &getOpt)
	if err != nil {
		return diag.Errorf("error retrieving GaussDB global slow SQL detail: %s", err)
	}

	getRespBody, err := utils.FlattenResponse(getResp)
	if err != nil {
		return diag.Errorf("error flattening GaussDB global slow SQL detail response: %s", err)
	}

	dataSourceId, err := uuid.NewRandom()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(dataSourceId.String())

	mErr = multierror.Append(
		d.Set("region", region),
		d.Set("slow_sql_details", flattenGaussDbGlobalSlowSqlDetailBody(getRespBody)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func buildGaussDbGlobalSlowSqlDetailBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"start_time":     d.Get("start_time"),
		"end_time":       d.Get("end_time"),
		"instance_id":    d.Get("instance_id"),
		"sql_id":         d.Get("sql_id"),
		"node_ids":       d.Get("node_ids"),
		"component_type": d.Get("component_type"),
	}

	return bodyParams
}

func flattenGaussDbGlobalSlowSqlDetailBody(resp interface{}) []interface{} {
	curJson := utils.PathSearch("slow_sql_details", resp, make([]interface{}, 0))
	curArray := curJson.([]interface{})
	res := make([]interface{}, 0, len(curArray))
	for _, v := range curArray {
		res = append(res, map[string]interface{}{
			"db_name":       utils.PathSearch("db_name", v, nil),
			"schema_name":   utils.PathSearch("schema_name", v, nil),
			"sql_id":        utils.PathSearch("sql_id", v, nil),
			"user_name":     utils.PathSearch("user_name", v, nil),
			"client_ip":     utils.PathSearch("client_ip", v, nil),
			"client_port":   utils.PathSearch("client_port", v, nil),
			"node_id":       utils.PathSearch("node_id", v, nil),
			"node_name":     utils.PathSearch("node_name", v, nil),
			"sql_text":      utils.PathSearch("sql_text", v, nil),
			"sql":           utils.PathSearch("sql", v, nil),
			"query_plan":    utils.PathSearch("query_plan", v, nil),
			"start_time":    utils.PathSearch("start_time", v, nil),
			"finish_time":   utils.PathSearch("finish_time", v, nil),
			"returned_rows": utils.PathSearch("returned_rows", v, nil),
			"fetched_rows":  utils.PathSearch("fetched_rows", v, nil),
			"fetched_pages": utils.PathSearch("fetched_pages", v, nil),
			"hit_pages":     utils.PathSearch("hit_pages", v, nil),
			"total_time":    utils.PathSearch("total_time", v, nil),
			"cpu_time":      utils.PathSearch("cpu_time", v, nil),
			"plan_time":     utils.PathSearch("plan_time", v, nil),
			"io_time":       utils.PathSearch("io_time", v, nil),
			"lock_count":    utils.PathSearch("lock_count", v, nil),
			"lock_time":     utils.PathSearch("lock_time", v, nil),
		})
	}
	return res
}
