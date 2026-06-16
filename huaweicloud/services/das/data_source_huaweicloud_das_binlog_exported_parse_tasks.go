package das

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

// @API DAS GET /v3/{project_id}/connections/{connection_id}/binlog-parse/export-list
func DataSourceBinlogExportedParseTasks() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceBinlogExportedParseTasksRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "The region where the binlog exported parse tasks are located.",
			},

			// Required parameters.
			"user_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The database user ID.",
			},

			// Attributes.
			"tasks": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "The list of binlog exported parse tasks.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"exported_task_id": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "The exported task ID.",
						},
						"parsed_task_id": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "The parsed task ID.",
						},
						"instance_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The instance ID.",
						},
						"status": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "The task status.",
						},
						"start_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The start time, in RFC3339 format.",
						},
						"end_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The end time, in RFC3339 format.",
						},
						"last_record_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The last record time, in RFC3339 format.",
						},
						"created_at": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The task creation time, in RFC3339 format.",
						},
						"export_line_num": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "The number of exported lines.",
						},
						"download_url": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The download URL of the exported file.",
						},
						"source_file_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The binlog source file name.",
						},
					},
				},
			},
		},
	}
}

func dataSourceBinlogExportedParseTasksRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	client, err := cfg.NewServiceClient("das", region)
	if err != nil {
		return diag.Errorf("error creating DAS client: %s", err)
	}

	taskExports, err := listBinlogExportedParseTasks(client, d)
	if err != nil {
		return diag.Errorf("error querying DAS binlog exported parse tasks: %s", err)
	}

	randomUUID, err := uuid.NewRandom()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(randomUUID.String())

	// The GET API will return a 'task_id', but the meaning of this 'task_id' is exported binlog task
	// ID, not the parsed binlog task ID.
	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("tasks", flattenBinlogExportedParseTasks(taskExports)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func listBinlogExportedParseTasks(client *golangsdk.ServiceClient, d *schema.ResourceData) ([]interface{}, error) {
	var (
		httpUrl = "v3/{project_id}/connections/{connection_id}/binlog-parse/export-list"
		perPage = 100
		curPage = 1
		result  = make([]interface{}, 0)
	)

	listPath := client.Endpoint + httpUrl
	listPath = strings.ReplaceAll(listPath, "{project_id}", client.ProjectID)
	listPath = strings.ReplaceAll(listPath, "{connection_id}", d.Get("user_id").(string))

	opt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json",
		},
	}

	for {
		listPathWithPage := listPath + fmt.Sprintf("?cur_page=%d&per_page=%d", curPage, perPage)

		requestResp, err := client.Request("GET", listPathWithPage, &opt)
		if err != nil {
			return nil, err
		}
		respBody, err := utils.FlattenResponse(requestResp)
		if err != nil {
			return nil, err
		}

		taskList := utils.PathSearch("task_list", respBody, make([]interface{}, 0)).([]interface{})
		result = append(result, taskList...)

		if len(taskList) < perPage {
			break
		}
		curPage++
	}

	return result, nil
}

func flattenBinlogExportedParseTasks(items []interface{}) []map[string]interface{} {
	if len(items) < 1 {
		return nil
	}

	result := make([]map[string]interface{}, 0, len(items))
	for _, item := range items {
		result = append(result, map[string]interface{}{
			"exported_task_id": utils.PathSearch("task_id", item, nil),
			"instance_id":      utils.PathSearch("instance_id", item, nil),
			"status":           utils.PathSearch("task_status", item, nil),
			"export_line_num":  utils.PathSearch("export_line_num", item, nil),
			"download_url":     utils.PathSearch("download_url", item, nil),
			"source_file_name": utils.PathSearch("source_file_name", item, nil),
			"parsed_task_id":   utils.PathSearch("parse_task_id", item, nil),
			"start_time": utils.FormatTimeStampRFC3339(
				int64(utils.PathSearch("start_time", item, float64(0)).(float64))/1000, false),
			"end_time": utils.FormatTimeStampRFC3339(
				int64(utils.PathSearch("end_time", item, float64(0)).(float64))/1000, false),
			"last_record_time": utils.FormatTimeStampRFC3339(
				int64(utils.PathSearch("last_record_time", item, float64(0)).(float64))/1000, false),
			"created_at": utils.FormatTimeStampRFC3339(
				int64(utils.PathSearch("create_at", item, float64(0)).(float64))/1000, false),
		})
	}

	return result
}
