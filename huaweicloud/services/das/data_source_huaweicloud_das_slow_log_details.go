package das

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

// @API DAS POST /v3/{project_id}/instances/{instance_id}/slow-log/get-slow-log-detail-list
func DataSourceSlowLogDetails() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceSlowLogDetailsRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `The region where the slow log details are located.`,
			},

			// Required parameters.
			"instance_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The instance ID.`,
			},
			"start_time": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The start time of the query range, in RFC3339 format.`,
			},
			"end_time": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The end time of the query range, in RFC3339 format.`,
			},

			// Optional parameters.
			"node_ids": {
				Type:        schema.TypeList,
				Optional:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: `The list of node IDs.`,
			},
			"db_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The database name.`,
			},
			"sort_field": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The field to sort by.`,
			},
			"sort_asc": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: `The sort order. true means ascending, false means descending.`,
			},
			"client_ip_address": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The client IP address.`,
			},
			"user_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The user name.`,
			},
			"killed": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The execution status.`,
			},
			"execute_time_min": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: `The minimum execution time, in milliseconds.`,
			},
			"execute_time_max": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: `The maximum execution time, in milliseconds.`,
			},
			"rows_max_examined": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: `The maximum number of rows examined.`,
			},
			"rows_min_examined": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: `The minimum number of rows examined.`,
			},
			"fuzzy_sql": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The fuzzy SQL pattern.`,
			},
			"operation": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The SQL operation types, separated by commas.`,
			},

			// Attributes.
			"details": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: `The list of slow log details that matched the filter parameters.`,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"occurrence_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The occurrence time, in RFC3339 format.`,
						},
						"sql_template_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The SQL template ID.`,
						},
						"original_sql": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The original SQL statement.`,
						},
						"db_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The database name.`,
						},
						"client_ip_address": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The client IP address.`,
						},
						"user_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The user name.`,
						},
						"execute_time": {
							Type:        schema.TypeFloat,
							Computed:    true,
							Description: `The execution time, in seconds.`,
						},
						"lock_wait_time": {
							Type:        schema.TypeFloat,
							Computed:    true,
							Description: `The lock wait time, in seconds.`,
						},
						"rows_examined": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: `The number of rows examined.`,
						},
						"rows_sent": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: `The number of rows sent.`,
						},
						"tunable": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: `Whether the SQL can be tuned.`,
						},
						"end_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The end time, in RFC3339 format.`,
						},
						"app_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The application name.`,
						},
						"rows_affected": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: `The number of rows affected.`,
						},
						"cpu_time": {
							Type:        schema.TypeFloat,
							Computed:    true,
							Description: `The CPU time, in milliseconds.`,
						},
						"logical_reads": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: `The number of logical reads.`,
						},
						"physical_reads": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: `The number of physical reads.`,
						},
						"writes": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: `The number of writes.`,
						},
						"sql_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The SQL type.`,
						},
						"collection": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The collection name (for MongoDB).`,
						},
						"key_examined": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: `The number of keys examined (for MongoDB).`,
						},
						"node_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The node ID (for MongoDB).`,
						},
						"node_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The node name (for MongoDB).`,
						},
						"killed": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The execution status.`,
						},
					},
				},
			},
		},
	}
}

func dataSourceSlowLogDetailsRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	client, err := cfg.NewServiceClient("das", region)
	if err != nil {
		return diag.Errorf("error creating DAS client: %s", err)
	}

	details, err := listSlowLogDetails(client, d)
	if err != nil {
		return diag.Errorf("error querying DAS slow log details: %s", err)
	}

	randomUUID, err := uuid.NewRandom()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(randomUUID.String())

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("details", flattenSlowLogDetails(details)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func listSlowLogDetails(client *golangsdk.ServiceClient, d *schema.ResourceData) ([]interface{}, error) {
	var (
		httpUrl = "v3/{project_id}/instances/{instance_id}/slow-log/get-slow-log-detail-list"
		perPage = 100
		curPage = 1
		result  = make([]interface{}, 0)
	)

	listPath := client.Endpoint + httpUrl
	listPath = strings.ReplaceAll(listPath, "{project_id}", client.ProjectID)
	listPath = strings.ReplaceAll(listPath, "{instance_id}", d.Get("instance_id").(string))

	opt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type":     "application/json",
			"X-Source-Service": "das",
		},
	}

	for {
		requestBody := utils.RemoveNil(buildSlowLogDetailsRequestBody(d, curPage, perPage))
		opt.JSONBody = requestBody

		requestResp, err := client.Request("POST", listPath, &opt)
		if err != nil {
			return nil, err
		}
		respBody, err := utils.FlattenResponse(requestResp)
		if err != nil {
			return nil, err
		}

		detailList := utils.PathSearch("detail_list", respBody, make([]interface{}, 0)).([]interface{})
		result = append(result, detailList...)

		if len(detailList) < perPage {
			break
		}
		curPage++
	}

	return result, nil
}

func buildSlowLogDetailsRequestBody(d *schema.ResourceData, curPage, perPage int) map[string]interface{} {
	// The `start_time` and `end_time` fields of the API are both in milliseconds.
	// So they don't need to be divided by 1000.
	body := map[string]interface{}{
		// required
		"start_time": utils.ConvertTimeStrToNanoTimestamp(d.Get("start_time").(string)),
		"end_time":   utils.ConvertTimeStrToNanoTimestamp(d.Get("end_time").(string)),
		"cur_page":   curPage,
		"per_page":   perPage,

		// optional
		"node_ids":          utils.ValueIgnoreEmpty(d.Get("node_ids")),
		"db_name":           utils.ValueIgnoreEmpty(d.Get("db_name")),
		"sort_field":        utils.ValueIgnoreEmpty(d.Get("sort_field")),
		"sort_asc":          utils.ValueIgnoreEmpty(d.Get("sort_asc")),
		"client":            utils.ValueIgnoreEmpty(d.Get("client_ip_address")),
		"user":              utils.ValueIgnoreEmpty(d.Get("user_name")),
		"killed":            utils.ValueIgnoreEmpty(d.Get("killed")),
		"execute_time_min":  utils.ValueIgnoreEmpty(d.Get("execute_time_min")),
		"execute_time_max":  utils.ValueIgnoreEmpty(d.Get("execute_time_max")),
		"rows_max_examined": utils.ValueIgnoreEmpty(d.Get("rows_max_examined")),
		"rows_min_examined": utils.ValueIgnoreEmpty(d.Get("rows_min_examined")),
		"fuzzy_sql":         utils.ValueIgnoreEmpty(d.Get("fuzzy_sql")),
		"operation":         utils.ValueIgnoreEmpty(d.Get("operation")),
	}

	return body
}

func flattenSlowLogDetails(items []interface{}) []map[string]interface{} {
	if len(items) < 1 {
		return nil
	}

	result := make([]map[string]interface{}, 0, len(items))
	for _, item := range items {
		result = append(result, map[string]interface{}{
			"sql_template_id":   utils.PathSearch("sql_template_id", item, nil),
			"original_sql":      utils.PathSearch("original_sql", item, nil),
			"db_name":           utils.PathSearch("db_name", item, nil),
			"client_ip_address": utils.PathSearch("client", item, nil),
			"user_name":         utils.PathSearch("user", item, nil),
			"execute_time":      utils.PathSearch("execute_time", item, nil),
			"lock_wait_time":    utils.PathSearch("lock_wait_time", item, nil),
			"rows_examined":     utils.PathSearch("rows_examined", item, nil),
			"rows_sent":         utils.PathSearch("rows_sent", item, nil),
			"tunable":           utils.PathSearch("tunable", item, nil),
			"app_name":          utils.PathSearch("app_name", item, nil),
			"rows_affected":     utils.PathSearch("rows_affected", item, nil),
			"cpu_time":          utils.PathSearch("cpu_time", item, nil),
			"logical_reads":     utils.PathSearch("logical_reads", item, nil),
			"physical_reads":    utils.PathSearch("physical_reads", item, nil),
			"writes":            utils.PathSearch("writes", item, nil),
			"sql_type":          utils.PathSearch("sql_type", item, nil),
			"collection":        utils.PathSearch("collection", item, nil),
			"key_examined":      utils.PathSearch("key_examined", item, nil),
			"node_id":           utils.PathSearch("node_id", item, nil),
			"node_name":         utils.PathSearch("node_name", item, nil),
			"killed":            utils.PathSearch("killed", item, nil),
			"end_time": utils.FormatTimeStampRFC3339(
				int64(utils.PathSearch("end_time", item, float64(0)).(float64))/1000, false),
			"occurrence_time": utils.FormatTimeStampRFC3339(
				int64(utils.PathSearch("occurrence_time", item, float64(0)).(float64))/1000, false),
		})
	}

	return result
}
