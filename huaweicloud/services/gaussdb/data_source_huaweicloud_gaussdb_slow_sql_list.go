package gaussdb

import (
	"context"
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

// @API GaussDB POST /v3/{project_id}/instances/{instance_id}/slow-sql-list
func DataSourceGaussDbSlowSqlList() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceGaussDbSlowSqlListRead,

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
			"threshold": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"node_ids": {
				Type:     schema.TypeList,
				Required: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"sql_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"multi_queries": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     slowSqlListMultiQueriesSchema(),
			},
			"slow_sql_infos": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     slowSqlListInfosSchema(),
			},
		},
	}
}

func slowSqlListMultiQueriesSchema() *schema.Resource {
	return &schema.Resource{
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
	}
}

func slowSqlListInfosSchema() *schema.Resource {
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
			"calls": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"avg_exec_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"avg_cpu_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"avg_io_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"avg_returned_rows": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"avg_fetched_rows": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"buffer_hit_ratio": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"sql_hit_ratio": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceGaussDbSlowSqlListRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var mErr *multierror.Error

	client, err := cfg.NewServiceClient("opengauss", region)
	if err != nil {
		return diag.Errorf("error creating GaussDB client: %s", err)
	}

	httpUrl := "v3/{project_id}/instances/{instance_id}/slow-sql-list"
	getPath := client.Endpoint + httpUrl
	getPath = strings.ReplaceAll(getPath, "{project_id}", client.ProjectID)
	getPath = strings.ReplaceAll(getPath, "{instance_id}", d.Get("instance_id").(string))

	getOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json",
		},
	}

	getOpt.JSONBody = utils.RemoveNil(buildGaussDbSlowSqlListBodyParams(d))
	getResp, err := client.Request("POST", getPath, &getOpt)
	if err != nil {
		return diag.Errorf("error retrieving GaussDB slow SQL list: %s", err)
	}

	getRespBody, err := utils.FlattenResponse(getResp)
	if err != nil {
		return diag.Errorf("error flattening GaussDB slow SQL list response: %s", err)
	}

	dataSourceId, err := uuid.NewRandom()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(dataSourceId.String())

	mErr = multierror.Append(
		mErr,
		d.Set("region", region),
		d.Set("slow_sql_infos", flattenGaussDbSlowSqlListBody(getRespBody)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func buildGaussDbSlowSqlListBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"start_time":  d.Get("start_time"),
		"end_time":    d.Get("end_time"),
		"instance_id": d.Get("instance_id"),
		"threshold":   d.Get("threshold"),
		"node_ids":    d.Get("node_ids"),
		"sql_id":      d.Get("sql_id"),
	}

	if v, ok := d.GetOk("multi_queries"); ok {
		bodyParams["multi_queries"] = buildSlowSqlListMultiQueries(v.([]interface{}))
	}

	return bodyParams
}

func buildSlowSqlListMultiQueries(rawParams []interface{}) []map[string]interface{} {
	if len(rawParams) == 0 {
		return nil
	}
	queries := make([]map[string]interface{}, 0, len(rawParams))
	for _, rawParam := range rawParams {
		param := rawParam.(map[string]interface{})
		query := map[string]interface{}{
			"name":      param["name"],
			"condition": param["condition"],
			"values":    param["values"],
		}
		if v, ok := param["is_fuzzy"]; ok && v.(string) != "" {
			isFuzzy, _ := strconv.ParseBool(v.(string))
			query["is_fuzzy"] = isFuzzy
		}
		queries = append(queries, query)
	}
	return queries
}

func flattenGaussDbSlowSqlListBody(resp interface{}) []interface{} {
	curJson := utils.PathSearch("slow_sql_infos", resp, make([]interface{}, 0))
	curArray := curJson.([]interface{})
	res := make([]interface{}, 0, len(curArray))
	for _, v := range curArray {
		res = append(res, map[string]interface{}{
			"db_name":           utils.PathSearch("db_name", v, nil),
			"schema_name":       utils.PathSearch("schema_name", v, nil),
			"sql_id":            utils.PathSearch("sql_id", v, nil),
			"user_name":         utils.PathSearch("user_name", v, nil),
			"node_id":           utils.PathSearch("node_id", v, nil),
			"node_name":         utils.PathSearch("node_name", v, nil),
			"sql_text":          utils.PathSearch("sql_text", v, nil),
			"sql":               utils.PathSearch("sql", v, nil),
			"query_plan":        utils.PathSearch("query_plan", v, nil),
			"calls":             utils.PathSearch("calls", v, nil),
			"avg_exec_time":     utils.PathSearch("avg_exec_time", v, nil),
			"avg_cpu_time":      utils.PathSearch("avg_cpu_time", v, nil),
			"avg_io_time":       utils.PathSearch("avg_io_time", v, nil),
			"avg_returned_rows": utils.PathSearch("avg_returned_rows", v, nil),
			"avg_fetched_rows":  utils.PathSearch("avg_fetched_rows", v, nil),
			"buffer_hit_ratio":  utils.PathSearch("buffer_hit_ratio", v, nil),
			"sql_hit_ratio":     utils.PathSearch("sql_hit_ratio", v, nil),
		})
	}
	return res
}
