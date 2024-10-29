package iotda

import (
	"context"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/huaweicloud/huaweicloud-sdk-go-v3/services/iotda/v5/model"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API IoTDA GET /v5/iot/{project_id}/batchtasks
func DataSourceBatchTasks() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceBatchTasksRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"type": {
				Type:     schema.TypeString,
				Required: true,
			},
			"space_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"status": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"batchtasks": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"targets": {
							Type:     schema.TypeList,
							Computed: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
						"targets_filter": {
							Type:     schema.TypeMap,
							Computed: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
						"task_policy": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"schedule_time": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"retry_count": {
										Type:     schema.TypeInt,
										Computed: true,
									},
									"retry_interval": {
										Type:     schema.TypeInt,
										Computed: true,
									},
								},
							},
						},
						"status": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"status_desc": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"task_progress": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"total": {
										Type:     schema.TypeInt,
										Computed: true,
									},
									"processing": {
										Type:     schema.TypeInt,
										Computed: true,
									},
									"success": {
										Type:     schema.TypeInt,
										Computed: true,
									},
									"fail": {
										Type:     schema.TypeInt,
										Computed: true,
									},
									"waitting": {
										Type:     schema.TypeInt,
										Computed: true,
									},
									"fail_wait_retry": {
										Type:     schema.TypeInt,
										Computed: true,
									},
									"stopped": {
										Type:     schema.TypeInt,
										Computed: true,
									},
									"removed": {
										Type:     schema.TypeInt,
										Computed: true,
									},
								},
							},
						},
						"created_at": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceBatchTasksRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg       = meta.(*config.Config)
		region    = cfg.GetRegion(d)
		isDerived = WithDerivedAuth(cfg, region)
	)

	client, err := cfg.HcIoTdaV5Client(region, isDerived)
	if err != nil {
		return diag.Errorf("error creating IoTDA v5 client: %s", err)
	}

	var (
		allTasks []model.Task
		limit    = int32(50)
		offset   int32
	)

	for {
		listOpts := model.ListBatchTasksRequest{
			AppId:    utils.StringIgnoreEmpty(d.Get("space_id").(string)),
			TaskType: d.Get("type").(string),
			Status:   utils.StringIgnoreEmpty(d.Get("status").(string)),
			Limit:    utils.Int32(limit),
			Offset:   &offset,
		}

		listResp, listErr := client.ListBatchTasks(&listOpts)
		if listErr != nil {
			return diag.Errorf("error querying IoTDA batch tasks: %s", listErr)
		}

		if listResp == nil || listResp.Batchtasks == nil {
			break
		}

		if len(*listResp.Batchtasks) == 0 {
			break
		}

		allTasks = append(allTasks, *listResp.Batchtasks...)
		//nolint:gosec
		offset += int32(len(*listResp.Batchtasks))
	}

	uuId, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}

	d.SetId(uuId)

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("batchtasks", flattenBatchTasks(allTasks)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenBatchTasks(tasks []model.Task) []interface{} {
	if len(tasks) == 0 {
		return nil
	}

	rst := make([]interface{}, 0, len(tasks))
	for _, v := range tasks {
		rst = append(rst, map[string]interface{}{
			"id":             v.TaskId,
			"name":           v.TaskName,
			"type":           v.TaskType,
			"targets":        v.Targets,
			"targets_filter": v.TargetsFilter,
			"task_policy":    flattenDataSourceTaskPolicy(v.TaskPolicy),
			"status":         v.Status,
			"status_desc":    v.StatusDesc,
			"task_progress":  flattenDataSourceTaskProgress(v.TaskProgress),
			"created_at":     v.CreateTime,
		})
	}

	return rst
}

func flattenDataSourceTaskPolicy(taskPolicy *model.TaskPolicy) []interface{} {
	if taskPolicy == nil {
		return nil
	}

	rst := map[string]interface{}{
		"schedule_time":  taskPolicy.ScheduleTime,
		"retry_count":    taskPolicy.RetryCount,
		"retry_interval": taskPolicy.RetryInterval,
	}

	return []interface{}{rst}
}

func flattenDataSourceTaskProgress(taskProgress *model.TaskProgress) []interface{} {
	if taskProgress == nil {
		return nil
	}

	rst := map[string]interface{}{
		"total":           taskProgress.Total,
		"processing":      taskProgress.Processing,
		"success":         taskProgress.Success,
		"fail":            taskProgress.Fail,
		"waitting":        taskProgress.Waitting,
		"fail_wait_retry": taskProgress.FailWaitRetry,
		"stopped":         taskProgress.Stopped,
		"removed":         taskProgress.Removed,
	}

	return []interface{}{rst}
}
