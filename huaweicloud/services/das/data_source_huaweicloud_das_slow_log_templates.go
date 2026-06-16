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

// @API DAS POST /v3/{project_id}/instances/{instance_id}/slow-log/get-slow-log-tpl-list
func DataSourceSlowLogTemplates() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceSlowLogTemplatesRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `The region where the slow log templates are located.`,
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
			"template_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The SQL template ID.`,
			},
			"node_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The node ID of the instance.`,
			},
			"db_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The database name.`,
			},
			"min_avg_execute_time": {
				Type:        schema.TypeFloat,
				Optional:    true,
				Description: `The minimum average execution time, in milliseconds.`,
			},
			"max_avg_execute_time": {
				Type:        schema.TypeFloat,
				Optional:    true,
				Description: `The maximum average execution time, in milliseconds.`,
			},
			"operation": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The SQL operation types, separated by commas.`,
			},

			// Attributes.
			"templates": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: `The list of slow log templates that matched filter parameters.`,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"template_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The name of the SQL template.`,
						},
						"template_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The ID of the SQL template.`,
						},
						"sql_sample": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The SQL sample.`,
						},
						"db_names": {
							Type:        schema.TypeList,
							Computed:    true,
							Elem:        &schema.Schema{Type: schema.TypeString},
							Description: `The database names.`,
						},
						"execute_count": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: `The execution count.`,
						},
						"avg_execute_time": {
							Type:        schema.TypeFloat,
							Computed:    true,
							Description: `The average execution time, in milliseconds.`,
						},
						"max_execute_time": {
							Type:        schema.TypeFloat,
							Computed:    true,
							Description: `The maximum execution time, in milliseconds.`,
						},
						"avg_lock_wait_time": {
							Type:        schema.TypeFloat,
							Computed:    true,
							Description: `The average lock wait time, in milliseconds.`,
						},
						"max_lock_wait_time": {
							Type:        schema.TypeFloat,
							Computed:    true,
							Description: `The maximum lock wait time, in milliseconds.`,
						},
						"avg_rows_examined": {
							Type:        schema.TypeFloat,
							Computed:    true,
							Description: `The average number of rows examined.`,
						},
						"max_rows_examined": {
							Type:        schema.TypeFloat,
							Computed:    true,
							Description: `The maximum number of rows examined.`,
						},
						"avg_rows_sent": {
							Type:        schema.TypeFloat,
							Computed:    true,
							Description: `The average number of rows sent.`,
						},
						"max_rows_sent": {
							Type:        schema.TypeFloat,
							Computed:    true,
							Description: `The maximum number of rows sent.`,
						},
						"tunable": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: `Whether the SQL can be tuned.`,
						},
						"node_ids": {
							Type:        schema.TypeList,
							Computed:    true,
							Elem:        &schema.Schema{Type: schema.TypeString},
							Description: `The node IDs.`,
						},
						"avg_cpu_time": {
							Type:        schema.TypeFloat,
							Computed:    true,
							Description: `The average CPU time, in milliseconds.`,
						},
						"max_cpu_time": {
							Type:        schema.TypeFloat,
							Computed:    true,
							Description: `The maximum CPU time, in milliseconds.`,
						},
						"avg_rows_affected": {
							Type:        schema.TypeFloat,
							Computed:    true,
							Description: `The average number of rows affected.`,
						},
						"max_rows_affected": {
							Type:        schema.TypeFloat,
							Computed:    true,
							Description: `The maximum number of rows affected.`,
						},
						"avg_logical_reads": {
							Type:        schema.TypeFloat,
							Computed:    true,
							Description: `The average number of logical reads.`,
						},
						"max_logical_reads": {
							Type:        schema.TypeFloat,
							Computed:    true,
							Description: `The maximum number of logical reads.`,
						},
						"avg_physical_reads": {
							Type:        schema.TypeFloat,
							Computed:    true,
							Description: `The average number of physical reads.`,
						},
						"max_physical_reads": {
							Type:        schema.TypeFloat,
							Computed:    true,
							Description: `The maximum number of physical reads.`,
						},
						"avg_writes": {
							Type:        schema.TypeFloat,
							Computed:    true,
							Description: `The average number of writes.`,
						},
						"max_writes": {
							Type:        schema.TypeFloat,
							Computed:    true,
							Description: `The maximum number of writes.`,
						},
						"instance_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The instance ID.`,
						},
						"total_execute_time_ratio": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The total execution time ratio.`,
						},
						"execute_count_ratio": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The execution count ratio.`,
						},
					},
				},
			},
		},
	}
}

func dataSourceSlowLogTemplatesRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	client, err := cfg.NewServiceClient("das", region)
	if err != nil {
		return diag.Errorf("error creating DAS client: %s", err)
	}

	templates, err := listSlowLogTemplates(client, d)
	if err != nil {
		return diag.Errorf("error querying DAS slow log templates: %s", err)
	}

	randomUUID, err := uuid.NewRandom()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(randomUUID.String())

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("templates", flattenSlowLogTemplates(templates)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func listSlowLogTemplates(client *golangsdk.ServiceClient, d *schema.ResourceData) ([]interface{}, error) {
	var (
		httpUrl = "v3/{project_id}/instances/{instance_id}/slow-log/get-slow-log-tpl-list"
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
		requestBody := utils.RemoveNil(buildSlowLogTemplatesRequestBody(d, curPage, perPage))
		opt.JSONBody = requestBody

		requestResp, err := client.Request("POST", listPath, &opt)
		if err != nil {
			return nil, err
		}
		respBody, err := utils.FlattenResponse(requestResp)
		if err != nil {
			return nil, err
		}

		tplList := utils.PathSearch("tpl_list", respBody, make([]interface{}, 0)).([]interface{})
		result = append(result, tplList...)

		if len(tplList) < perPage {
			break
		}
		curPage++
	}

	return result, nil
}

func buildSlowLogTemplatesRequestBody(d *schema.ResourceData, curPage, perPage int) map[string]interface{} {
	body := map[string]interface{}{
		// required
		"start_time": utils.ConvertTimeStrToNanoTimestamp(d.Get("start_time").(string)),
		"end_time":   utils.ConvertTimeStrToNanoTimestamp(d.Get("end_time").(string)),
		"cur_page":   curPage,
		"per_page":   perPage,

		// optional
		"sql_template_id":      utils.ValueIgnoreEmpty(d.Get("template_id")),
		"node_id":              utils.ValueIgnoreEmpty(d.Get("node_id")),
		"db_name":              utils.ValueIgnoreEmpty(d.Get("db_name")),
		"min_avg_execute_time": utils.ValueIgnoreEmpty(d.Get("min_avg_execute_time")),
		"max_avg_execute_time": utils.ValueIgnoreEmpty(d.Get("max_avg_execute_time")),
		"operation":            utils.ValueIgnoreEmpty(d.Get("operation")),
	}

	return body
}

func flattenSlowLogTemplates(items []interface{}) []map[string]interface{} {
	if len(items) < 1 {
		return nil
	}

	result := make([]map[string]interface{}, 0, len(items))
	for _, item := range items {
		result = append(result, map[string]interface{}{
			"template_name":            utils.PathSearch("sql_template", item, nil),
			"template_id":              utils.PathSearch("sql_template_id", item, nil),
			"sql_sample":               utils.PathSearch("sql_sample", item, nil),
			"db_names":                 utils.PathSearch("db_names", item, nil),
			"execute_count":            utils.PathSearch("execute_count", item, nil),
			"avg_execute_time":         utils.PathSearch("avg_execute_time", item, nil),
			"max_execute_time":         utils.PathSearch("max_execute_time", item, nil),
			"avg_lock_wait_time":       utils.PathSearch("avg_lock_wait_time", item, nil),
			"max_lock_wait_time":       utils.PathSearch("max_lock_wait_time", item, nil),
			"avg_rows_examined":        utils.PathSearch("avg_rows_examined", item, nil),
			"max_rows_examined":        utils.PathSearch("max_rows_examined", item, nil),
			"avg_rows_sent":            utils.PathSearch("avg_rows_sent", item, nil),
			"max_rows_sent":            utils.PathSearch("max_rows_sent", item, nil),
			"tunable":                  utils.PathSearch("tunable", item, nil),
			"node_ids":                 utils.PathSearch("node_ids", item, nil),
			"avg_cpu_time":             utils.PathSearch("avg_cpu_time", item, nil),
			"max_cpu_time":             utils.PathSearch("max_cpu_time", item, nil),
			"avg_rows_affected":        utils.PathSearch("avg_rows_affected", item, nil),
			"max_rows_affected":        utils.PathSearch("max_rows_affected", item, nil),
			"avg_logical_reads":        utils.PathSearch("avg_logical_reads", item, nil),
			"max_logical_reads":        utils.PathSearch("max_logical_reads", item, nil),
			"avg_physical_reads":       utils.PathSearch("avg_physical_reads", item, nil),
			"max_physical_reads":       utils.PathSearch("max_physical_reads", item, nil),
			"avg_writes":               utils.PathSearch("avg_writes", item, nil),
			"max_writes":               utils.PathSearch("max_writes", item, nil),
			"instance_id":              utils.PathSearch("instance_id", item, nil),
			"total_execute_time_ratio": utils.PathSearch("total_execute_time_ratio", item, nil),
			"execute_count_ratio":      utils.PathSearch("execute_count_ratio", item, nil),
		})
	}

	return result
}
