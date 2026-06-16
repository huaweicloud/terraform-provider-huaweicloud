package das

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/google/uuid"
	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API DAS GET /v3/{project_id}/transaction/{instance_id}/get-export-task-list
func DataSourceHistoryTransactionExportedTasks() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceHistoryTransactionExportedTasksRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "The region where the history transaction exported tasks are located.",
			},

			// Required parameters.
			"instance_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The instance ID.",
			},

			// Attributes.
			"tasks": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "The list of history transaction exported tasks.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "The exported task ID.",
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
						"created_time": {
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
					},
				},
			},
		},
	}
}

func dataSourceHistoryTransactionExportedTasksRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	client, err := cfg.NewServiceClient("das", region)
	if err != nil {
		return diag.Errorf("error creating DAS client: %s", err)
	}

	taskExports, err := listHistoryTransactionExportedTasks(client, d)
	if err != nil {
		return diag.Errorf("error querying DAS history transaction exported tasks: %s", err)
	}

	randomUUID, err := uuid.NewRandom()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(randomUUID.String())

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("tasks", flattenHistoryTransactionExportedTasks(taskExports)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func listHistoryTransactionExportedTasks(client *golangsdk.ServiceClient, d *schema.ResourceData) ([]interface{}, error) {
	var (
		httpUrl = "v3/{project_id}/transaction/{instance_id}/get-export-task-list?limit={limit}"
		limit   = 100
		offset  = 0
		result  = make([]interface{}, 0)
	)

	listPath := client.Endpoint + httpUrl
	listPath = strings.ReplaceAll(listPath, "{project_id}", client.ProjectID)
	listPath = strings.ReplaceAll(listPath, "{instance_id}", d.Get("instance_id").(string))
	listPath = strings.ReplaceAll(listPath, "{limit}", strconv.Itoa(limit))

	opt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json",
		},
	}

	for {
		listPathWithOffset := fmt.Sprintf("%s&offset=%d", listPath, offset)
		requestResp, err := client.Request("GET", listPathWithOffset, &opt)
		if err != nil {
			return nil, err
		}
		respBody, err := utils.FlattenResponse(requestResp)
		if err != nil {
			return nil, err
		}

		taskList := utils.PathSearch("task_list", respBody, make([]interface{}, 0)).([]interface{})
		result = append(result, taskList...)

		if len(taskList) < limit {
			break
		}
		offset += len(taskList)
	}

	return result, nil
}

func flattenHistoryTransactionExportedTasks(items []interface{}) []map[string]interface{} {
	if len(items) < 1 {
		return nil
	}

	result := make([]map[string]interface{}, 0, len(items))
	for _, item := range items {
		result = append(result, map[string]interface{}{
			"id":              utils.PathSearch("task_id", item, nil),
			"instance_id":     utils.PathSearch("instance_id", item, nil),
			"status":          utils.PathSearch("task_status", item, nil),
			"export_line_num": utils.PathSearch("export_line_num", item, nil),
			"download_url":    utils.PathSearch("download_url", item, nil),
			"start_time": utils.FormatTimeStampRFC3339(
				int64(utils.PathSearch("start_time", item, float64(0)).(float64))/1000, false),
			"end_time": utils.FormatTimeStampRFC3339(
				int64(utils.PathSearch("end_time", item, float64(0)).(float64))/1000, false),
			"created_time": utils.FormatTimeStampRFC3339(
				int64(utils.PathSearch("create_at", item, float64(0)).(float64))/1000, false),
		})
	}

	return result
}
