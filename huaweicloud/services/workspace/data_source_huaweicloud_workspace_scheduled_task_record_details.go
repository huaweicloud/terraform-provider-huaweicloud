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

// @API Workspace GET /v2/{project_id}/scheduled-tasks/{task_id}/records/{record_id}
func DataSourceScheduledTaskRecordDetails() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceScheduledTaskRecordDetailsRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `The region where the scheduled task record details are located.`,
			},

			// Required parameters.
			"task_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The ID of the scheduled task to be queried.`,
			},
			"record_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The ID of the scheduled task execution record to be queried.`,
			},

			// Attributes.
			"details": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The ID of the scheduled task execution record detail.`,
						},
						"record_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The ID of the scheduled task execution record.`,
						},
						"desktop_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The ID of the desktop.`,
						},
						"desktop_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The name of the desktop.`,
						},
						"exec_status": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The execution status.`,
						},
						"exec_script_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The ID of the execution script.`,
						},
						"result_code": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The error code of the failure or skip reason.`,
						},
						"fail_reason": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The failure or skip reason.`,
						},
						"start_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The execution start time, in RFC3339 format.`,
						},
						"end_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The execution end time, in RFC3339 format.`,
						},
						"time_zone": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The time zone information.`,
						},
					},
				},
				Description: `The list of scheduled task execution record details.`,
			},
		},
	}
}

func listScheduledTaskRecordDetails(client *golangsdk.ServiceClient, taskId, recordId string) ([]interface{}, error) {
	var (
		httpUrl = "v2/{project_id}/scheduled-tasks/{task_id}/records/{record_id}?limit={limit}"
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
	listPath = strings.ReplaceAll(listPath, "{record_id}", recordId)
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

		tasksRecordsDetails := utils.PathSearch("tasks_records_details", respBody, make([]interface{}, 0)).([]interface{})
		result = append(result, tasksRecordsDetails...)
		if len(tasksRecordsDetails) < limit {
			break
		}

		offset += len(tasksRecordsDetails)
	}

	return result, nil
}

func flattenScheduledTaskRecordDetails(items []interface{}) []map[string]interface{} {
	if len(items) < 1 {
		return nil
	}

	result := make([]map[string]interface{}, 0, len(items))
	for _, item := range items {
		result = append(result, map[string]interface{}{
			"id":             utils.PathSearch("id", item, nil),
			"record_id":      utils.PathSearch("record_id", item, nil),
			"desktop_id":     utils.PathSearch("desktop_id", item, nil),
			"desktop_name":   utils.PathSearch("desktop_name", item, nil),
			"exec_status":    utils.PathSearch("exec_status", item, nil),
			"exec_script_id": utils.PathSearch("exec_script_id", item, nil),
			"result_code":    utils.PathSearch("result_code", item, nil),
			"fail_reason":    utils.PathSearch("fail_reason", item, nil),
			"start_time": utils.FormatTimeStampRFC3339(utils.ConvertTimeStrToNanoTimestamp(utils.PathSearch("start_time",
				item, "").(string), "2006-01-02T15:04:05.000Z")/1000, false),
			"end_time": utils.FormatTimeStampRFC3339(utils.ConvertTimeStrToNanoTimestamp(utils.PathSearch("end_time",
				item, "").(string), "2006-01-02T15:04:05.000Z")/1000, false),
			"time_zone": utils.PathSearch("time_zone", item, nil),
		})
	}

	return result
}

func dataSourceScheduledTaskRecordDetailsRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg      = meta.(*config.Config)
		region   = cfg.GetRegion(d)
		taskId   = d.Get("task_id").(string)
		recordId = d.Get("record_id").(string)
	)

	client, err := cfg.NewServiceClient("workspace", region)
	if err != nil {
		return diag.Errorf("error creating Workspace client: %s", err)
	}

	items, err := listScheduledTaskRecordDetails(client, taskId, recordId)
	if err != nil {
		return diag.Errorf("error querying Workspace scheduled task record details: %s", err)
	}

	randUUID, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate data source ID: %s", err)
	}
	d.SetId(randUUID)

	mErr := multierror.Append(
		d.Set("region", region),
		d.Set("details", flattenScheduledTaskRecordDetails(items)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}
