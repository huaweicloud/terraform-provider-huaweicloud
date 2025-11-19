package workspace

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API Workspace GET /v2/{project_id}/scheduled-tasks/{task_id}/records
func DataSourceScheduledTaskRecords() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceScheduledTaskRecordsRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `The region where the scheduled task records are located.`,
			},

			// Required parameters.
			"task_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The ID of the scheduled task to be queried.`,
			},

			// Attributes.
			"records": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The ID of the scheduled task execution record.`,
						},
						"start_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The execution time, in RFC3339 format.`,
						},
						"task_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The type of the scheduled task.`,
						},
						"scheduled_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The execution cycle type.`,
						},
						"life_cycle_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The trigger scenario type.`,
						},
						"status": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The execution status of this execution.`,
						},
						"success_num": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: `The number of successful desktops.`,
						},
						"failed_num": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: `The number of failed desktops.`,
						},
						"skip_num": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: `The number of skipped desktops.`,
						},
						"time_zone": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The time zone information.`,
						},
						"execute_task_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The task ID of executing the scheduled task.`,
						},
						"execute_object_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The object type of executing the scheduled task.`,
						},
					},
				},
				Description: `The list of scheduled task execution records.`,
			},
		},
	}
}

func listScheduledTaskRecords(client *golangsdk.ServiceClient, taskId string) ([]interface{}, error) {
	var (
		httpUrl = "v2/{project_id}/scheduled-tasks/{task_id}/records?limit={limit}"
		offset  = 0
		// For API, limit default value is 10.
		limit   = 100
		result  = make([]interface{}, 0)
		listOpt = golangsdk.RequestOpts{
			KeepResponseBody: true,
		}
	)

	listPath := client.Endpoint + httpUrl
	listPath = strings.ReplaceAll(listPath, "{project_id}", client.ProjectID)
	listPath = strings.ReplaceAll(listPath, "{task_id}", taskId)
	listPath = strings.ReplaceAll(listPath, "{limit}", strconv.Itoa(limit))

	for {
		listPathWithOffset := fmt.Sprintf("%s&offset=%d", listPath, offset)
		listResp, err := client.Request("GET", listPathWithOffset, &listOpt)
		if err != nil {
			return nil, err
		}

		respBody, err := utils.FlattenResponse(listResp)
		if err != nil {
			return nil, err
		}

		tasksRecords := utils.PathSearch("tasks_records", respBody, make([]interface{}, 0)).([]interface{})
		result = append(result, tasksRecords...)
		if len(tasksRecords) < limit {
			break
		}

		offset += len(tasksRecords)
	}

	return result, nil
}

func flattenScheduledTaskRecords(items []interface{}) []map[string]interface{} {
	if len(items) < 1 {
		return nil
	}

	result := make([]map[string]interface{}, 0, len(items))
	for _, item := range items {
		result = append(result, map[string]interface{}{
			"id": utils.PathSearch("id", item, nil),
			"start_time": utils.FormatTimeStampRFC3339(utils.ConvertTimeStrToNanoTimestamp(utils.PathSearch("start_time",
				item, "").(string))/1000, false),
			"task_type":           utils.PathSearch("task_type", item, nil),
			"scheduled_type":      utils.PathSearch("scheduled_type", item, nil),
			"life_cycle_type":     utils.PathSearch("life_cycle_type", item, nil),
			"status":              utils.PathSearch("status", item, nil),
			"success_num":         utils.PathSearch("success_num", item, nil),
			"failed_num":          utils.PathSearch("failed_num", item, nil),
			"skip_num":            utils.PathSearch("skip_num", item, nil),
			"time_zone":           utils.PathSearch("time_zone", item, nil),
			"execute_task_id":     utils.PathSearch("execute_task_id", item, nil),
			"execute_object_type": utils.PathSearch("execute_object_type", item, nil),
		})
	}

	return result
}

func dataSourceScheduledTaskRecordsRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg    = meta.(*config.Config)
		region = cfg.GetRegion(d)
	)

	client, err := cfg.NewServiceClient("workspace", region)
	if err != nil {
		return diag.Errorf("error creating Workspace client: %s", err)
	}

	taskId := d.Get("task_id").(string)
	items, err := listScheduledTaskRecords(client, taskId)
	if err != nil {
		return diag.Errorf("error querying Workspace scheduled task records: %s", err)
	}

	randUUID, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate data source ID: %s", err)
	}
	d.SetId(randUUID)

	mErr := multierror.Append(
		d.Set("region", region),
		d.Set("records", flattenScheduledTaskRecords(items)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}
