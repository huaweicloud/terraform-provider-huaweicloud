package gaussdb

import (
	"context"
	"fmt"
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

// @API GaussDB POST /v3/{project_id}/instances/{instance_id}/top-sql-list
func DataSourceTopSqlStatements() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceTopSqlStatementsRead,

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
			"node_ids": {
				Type:     schema.TypeList,
				Required: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"start_time": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"end_time": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"start_time_utc": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"end_time_utc": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"support_system": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"sql_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"db_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"user_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"sql_text": {
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
							MinItems: 1,
							MaxItems: 5,
						},
						"is_fuzzy": {
							Type:         schema.TypeString,
							Optional:     true,
							ValidateFunc: validation.StringInSlice([]string{"true", "false"}, false),
						},
					},
				},
			},
			"top_sql_infos": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"sql_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"user_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"sql_text": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"calls_percent": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"cpu_percent": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"io_percent": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"calls": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"returned_rows": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"tuple_read": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"avg_elapse_time": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"total_elapse_time": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"cpu_time": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"io_time": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"min_elapse_time": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"max_elapse_time": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"sql_hit_ratio": {
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
						"db_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceTopSqlStatementsRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var mErr *multierror.Error

	client, err := cfg.NewServiceClient("opengauss", region)
	if err != nil {
		return diag.Errorf("error creating GaussDB client: %s", err)
	}

	httpUrl := "v3/{project_id}/instances/{instance_id}/top-sql-list"
	getPath := client.Endpoint + httpUrl
	getPath = strings.ReplaceAll(getPath, "{project_id}", client.ProjectID)
	getPath = strings.ReplaceAll(getPath, "{instance_id}", d.Get("instance_id").(string))

	reqBody := buildTopSqlStatementsRequestBody(d)
	getOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json",
		},
		JSONBody: utils.RemoveNil(reqBody),
	}

	getResp, err := client.Request("POST", getPath, &getOpt)
	if err != nil {
		return diag.Errorf("error retrieving GaussDB Top SQL statements: %s", err)
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

	topSqlInfos := utils.PathSearch("top_sql_infos", getRespBody, make([]interface{}, 0))
	topSqlInfosRaw, err := flattenTopSqlStatements(topSqlInfos)
	if err != nil {
		return diag.FromErr(err)
	}

	mErr = multierror.Append(
		d.Set("region", region),
		d.Set("top_sql_infos", topSqlInfosRaw),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func buildTopSqlStatementsRequestBody(d *schema.ResourceData) map[string]interface{} {
	body := map[string]interface{}{
		"instance_id":    d.Get("instance_id").(string),
		"node_ids":       utils.ExpandToStringList(d.Get("node_ids").([]interface{})),
		"start_time":     d.Get("start_time").(int),
		"end_time":       d.Get("end_time").(int),
		"support_system": d.Get("support_system").(bool),
	}

	if v, ok := d.GetOk("start_time_utc"); ok {
		body["start_time_utc"] = v.(string)
	}
	if v, ok := d.GetOk("end_time_utc"); ok {
		body["end_time_utc"] = v.(string)
	}
	if v, ok := d.GetOk("sql_id"); ok {
		body["sql_id"] = v.(string)
	}
	if v, ok := d.GetOk("db_name"); ok {
		body["db_name"] = v.(string)
	}
	if v, ok := d.GetOk("user_name"); ok {
		body["user_name"] = v.(string)
	}
	if v, ok := d.GetOk("sql_text"); ok {
		body["sql_text"] = v.(string)
	}
	if v, ok := d.GetOk("multi_queries"); ok {
		body["multi_queries"] = buildMultiQueriesRequestBody(v.([]interface{}))
	}

	return body
}

func buildMultiQueriesRequestBody(raw []interface{}) []map[string]interface{} {
	result := make([]map[string]interface{}, 0, len(raw))
	for _, item := range raw {
		itemMap := item.(map[string]interface{})
		isFuzzy, _ := strconv.ParseBool(itemMap["is_fuzzy"].(string))
		query := map[string]interface{}{
			"name":      itemMap["name"].(string),
			"condition": itemMap["condition"].(string),
			"values":    utils.ExpandToStringList(itemMap["values"].([]interface{})),
			"is_fuzzy":  isFuzzy,
		}
		result = append(result, query)
	}
	return result
}

func flattenTopSqlStatements(raw interface{}) ([]map[string]interface{}, error) {
	var result []map[string]interface{}

	rawArray, ok := raw.([]interface{})
	if !ok {
		return nil, fmt.Errorf("failed to parse top_sql_infos, expected []interface{}")
	}

	for _, item := range rawArray {
		itemMap, ok := item.(map[string]interface{})
		if !ok {
			continue
		}

		result = append(result, map[string]interface{}{
			"sql_id":            utils.PathSearch("sql_id", itemMap, nil),
			"user_name":         utils.PathSearch("user_name", itemMap, nil),
			"sql_text":          utils.PathSearch("sql_text", itemMap, nil),
			"calls_percent":     utils.PathSearch("calls_percent", itemMap, nil),
			"cpu_percent":       utils.PathSearch("cpu_percent", itemMap, nil),
			"io_percent":        utils.PathSearch("io_percent", itemMap, nil),
			"calls":             utils.PathSearch("calls", itemMap, nil),
			"returned_rows":     utils.PathSearch("returned_rows", itemMap, nil),
			"tuple_read":        utils.PathSearch("tuple_read", itemMap, nil),
			"avg_elapse_time":   utils.PathSearch("avg_elapse_time", itemMap, nil),
			"total_elapse_time": utils.PathSearch("total_elapse_time", itemMap, nil),
			"cpu_time":          utils.PathSearch("cpu_time", itemMap, nil),
			"io_time":           utils.PathSearch("io_time", itemMap, nil),
			"min_elapse_time":   utils.PathSearch("min_elapse_time", itemMap, nil),
			"max_elapse_time":   utils.PathSearch("max_elapse_time", itemMap, nil),
			"sql_hit_ratio":     utils.PathSearch("sql_hit_ratio", itemMap, nil),
			"node_id":           utils.PathSearch("node_id", itemMap, nil),
			"node_name":         utils.PathSearch("node_name", itemMap, nil),
			"db_name":           utils.PathSearch("db_name", itemMap, nil),
		})
	}

	return result, nil
}
