package workspace

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API Workspace POST /v1/{project_id}/schedule-task
// @API Workspace GET /v1/{project_id}/schedule-task/{task_id}
// @API Workspace PATCH /v1/{project_id}/schedule-task/{task_id}
// @API Workspace DELETE /v1/{project_id}/schedule-task/{task_id}
func ResourceAppScheduleTask() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceAppScheduleTaskCreate,
		ReadContext:   resourceAppScheduleTaskRead,
		UpdateContext: resourceAppScheduleTaskUpdate,
		DeleteContext: resourceAppScheduleTaskDelete,

		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: `The region where the schedule task is located.`,
			},
			"task_name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The name of the schedule task.`,
			},
			"task_type": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The type of the schedule task.`,
			},
			"scheduled_type": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The execution cycle of the schedule task.`,
			},
			"scheduled_time": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The execution time of the schedule task.`,
			},
			"target_infos": {
				Type:     schema.TypeSet,
				Required: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"target_id": {
							Type:        schema.TypeString,
							Required:    true,
							Description: `The ID of the target object.`,
						},
						"target_type": {
							Type:        schema.TypeString,
							Required:    true,
							Description: `The type of the target object.`,
						},
					},
				},
				Description: `The target object list of the schedule task.`,
			},
			"day_interval": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: `The execution interval of the scheduled task, in days.`,
				DiffSuppressFunc: func(_, _, newRaw string, d *schema.ResourceData) bool {
					// When scheduled_type is not "DAY", the update interface cannot update `day_interval` to null.
					// so the change needs to be suppressed.
					return newRaw == "0" && d.Get("scheduled_type") != "DAY"
				},
			},
			"week_list": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The days of week of the schedule task.`,
			},
			"month_list": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The months of the schedule task.`,
			},
			"date_list": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The days of month of the schedule task.`,
			},
			"scheduled_date": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The fixed date of the schedule task.`,
			},
			"time_zone": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `The time zone of the schedule task.`,
			},
			"expire_time": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The expiration time of the schedule task, in UTC format.`,
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The description of the schedule task.`,
			},
			"schedule_task_policy": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"enforcement_enable": {
							Type:        schema.TypeBool,
							Required:    true,
							Description: `Whether to forcefully execute the task when there are active sessions.`,
						},
					},
				},
				Description: `The policy of the schedule task.`,
			},
			"is_enable": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     true,
				Description: `Whether to enable the schedule task.`,
			},
		},
	}
}

func buildAppScheduleTaskPolicy(scheduleTaskPolicy []interface{}) map[string]interface{} {
	if len(scheduleTaskPolicy) < 1 {
		return nil
	}

	return map[string]interface{}{
		"enforcement_enable": utils.PathSearch("enforcement_enable", scheduleTaskPolicy[0], nil),
	}
}

func buildAppScheduleTaskTargetInfos(rawTargetInfos *schema.Set) []map[string]interface{} {
	if rawTargetInfos.Len() < 1 {
		return nil
	}

	result := make([]map[string]interface{}, rawTargetInfos.Len())
	for i, v := range rawTargetInfos.List() {
		result[i] = map[string]interface{}{
			"target_type": utils.PathSearch("target_type", v, nil),
			"target_id":   utils.PathSearch("target_id", v, nil),
		}
	}

	return result
}

func buildCreateAppScheduleTaskBodyParams(d *schema.ResourceData) map[string]interface{} {
	return map[string]interface{}{
		"task_name":            d.Get("task_name"),
		"task_type":            d.Get("task_type"),
		"scheduled_time":       d.Get("scheduled_time"),
		"scheduled_type":       d.Get("scheduled_type"),
		"target_infos":         buildAppScheduleTaskTargetInfos(d.Get("target_infos").(*schema.Set)),
		"day_interval":         utils.ValueIgnoreEmpty(d.Get("day_interval")),
		"week_list":            utils.ValueIgnoreEmpty(d.Get("week_list")),
		"month_list":           utils.ValueIgnoreEmpty(d.Get("month_list")),
		"date_list":            utils.ValueIgnoreEmpty(d.Get("date_list")),
		"time_zone":            utils.ValueIgnoreEmpty(d.Get("time_zone")),
		"scheduled_date":       utils.ValueIgnoreEmpty(d.Get("scheduled_date")),
		"expire_time":          utils.ValueIgnoreEmpty(d.Get("expire_time")),
		"description":          utils.ValueIgnoreEmpty(d.Get("description")),
		"schedule_task_policy": buildAppScheduleTaskPolicy(d.Get("schedule_task_policy").([]interface{})),
	}
}

func resourceAppScheduleTaskCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	client, err := cfg.NewServiceClient("appstream", cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating Workspace APP client: %s", err)
	}

	createPath := "v1/{project_id}/schedule-task"
	createPath = client.Endpoint + createPath
	createPath = strings.ReplaceAll(createPath, "{project_id}", client.ProjectID)
	createOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		JSONBody:         utils.RemoveNil(buildCreateAppScheduleTaskBodyParams(d)),
	}

	createResp, err := client.Request("POST", createPath, &createOpt)
	if err != nil {
		return diag.Errorf("error creating Workspace APP schedule task: %s", err)
	}
	respBody, err := utils.FlattenResponse(createResp)
	if err != nil {
		return diag.Errorf("error creating Workspace APP schedule task: %s", err)
	}

	scheduleTaskId := utils.PathSearch("id", respBody, "").(string)
	if scheduleTaskId == "" {
		return diag.Errorf("unable to find schedule task ID in the response")
	}
	d.SetId(scheduleTaskId)

	isEnabled := d.Get("is_enable").(bool)
	if !isEnabled {
		params := map[string]interface{}{
			"is_enable": isEnabled,
		}
		err = updateAppScheduleTask(client, scheduleTaskId, params)
		if err != nil {
			return diag.Errorf("unable to update the status of the schedule task (%s): %s", scheduleTaskId, err)
		}
	}

	return resourceAppScheduleTaskRead(ctx, d, meta)
}

// GetAppScheduleTaskById is a method is used to get the schedule task by its ID.
func GetAppScheduleTaskById(client *golangsdk.ServiceClient, scheduleTaskId string) (interface{}, error) {
	getPath := "v1/{project_id}/schedule-task/{task_id}"
	getPath = client.Endpoint + getPath
	getPath = strings.ReplaceAll(getPath, "{project_id}", client.ProjectID)
	getPath = strings.ReplaceAll(getPath, "{task_id}", scheduleTaskId)
	getOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	getResp, err := client.Request("GET", getPath, &getOpt)
	if err != nil {
		return nil, err
	}

	return utils.FlattenResponse(getResp)
}

func resourceAppScheduleTaskRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg            = meta.(*config.Config)
		region         = cfg.GetRegion(d)
		scheduleTaskId = d.Id()
	)

	client, err := cfg.NewServiceClient("appstream", region)
	if err != nil {
		return diag.Errorf("error creating Workspace APP client: %s", err)
	}

	respBody, err := GetAppScheduleTaskById(client, scheduleTaskId)
	if err != nil {
		return common.CheckDeletedDiag(d, err, "error retrieving Workspace APP schedule task")
	}

	mErr := multierror.Append(nil,
		d.Set("region", region),
		// Required parameters.
		d.Set("task_name", utils.PathSearch("task_name", respBody, nil)),
		d.Set("task_type", utils.PathSearch("task_type", respBody, nil)),
		d.Set("scheduled_type", utils.PathSearch("scheduled_type", respBody, nil)),
		d.Set("scheduled_time", utils.PathSearch("scheduled_time", respBody, nil)),
		d.Set("target_infos", flattenAppScheduleTaskTargetInfos(utils.PathSearch("target_infos", respBody, make([]interface{}, 0)).([]interface{}))),
		// Optional parameters.
		d.Set("day_interval", utils.PathSearch("day_interval", respBody, nil)),
		d.Set("week_list", utils.PathSearch("week_list", respBody, nil)),
		d.Set("month_list", utils.PathSearch("month_list", respBody, nil)),
		d.Set("date_list", utils.PathSearch("date_list", respBody, nil)),
		d.Set("time_zone", utils.PathSearch("time_zone", respBody, nil)),
		d.Set("scheduled_date", utils.PathSearch("scheduled_date", respBody, nil)),
		d.Set("expire_time", utils.PathSearch("expire_time", respBody, nil)),
		d.Set("description", utils.PathSearch("description", respBody, nil)),
		d.Set("schedule_task_policy", flattenAppScheduleTaskPolicy(utils.PathSearch("schedule_task_policy", respBody, nil))),
		d.Set("is_enable", utils.PathSearch("is_enable", respBody, nil)),
	)

	scheduleTasks, err := getAppScheduleTasks(client, d.Get("task_type").(string))
	if err != nil {
		return diag.Errorf("error getting Workspace APP schedule task by its ID (%s): %s", scheduleTaskId, err)
	}

	isEnabled := utils.PathSearch(fmt.Sprintf("[?id=='%s']|[0].is_enable", scheduleTaskId), scheduleTasks, false)
	mErr = multierror.Append(
		mErr,
		d.Set("is_enable", isEnabled),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenAppScheduleTaskPolicy(policy interface{}) []map[string]interface{} {
	if policy == nil {
		return nil
	}

	return []map[string]interface{}{
		{
			"enforcement_enable": utils.PathSearch("enforcement_enable", policy, nil),
		},
	}
}

func getAppScheduleTasks(client *golangsdk.ServiceClient, taskType string) ([]interface{}, error) {
	var (
		listPath = "v1/{project_id}/schedule-task"
		offset   = 0
		limit    = 100
		results  = make([]interface{}, 0)
	)

	listPath = client.Endpoint + listPath
	listPath = strings.ReplaceAll(listPath, "{project_id}", client.ProjectID)
	listPath = fmt.Sprintf("%s?limit=%d&task_type=%s", listPath, limit, taskType)

	listOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

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

		scheduleTasks := utils.PathSearch("items", respBody, make([]interface{}, 0)).([]interface{})
		results = append(results, scheduleTasks...)
		if len(scheduleTasks) < limit {
			break
		}

		offset += len(scheduleTasks)
	}

	return results, nil
}

func flattenAppScheduleTaskTargetInfos(targetInfos []interface{}) []map[string]interface{} {
	if len(targetInfos) < 1 {
		return nil
	}

	result := make([]map[string]interface{}, len(targetInfos))
	for i, targetObj := range targetInfos {
		result[i] = map[string]interface{}{
			"target_type": utils.PathSearch("target_type", targetObj, nil),
			"target_id":   utils.PathSearch("target_id", targetObj, nil),
		}
	}

	return result
}

func buildUpdateAppScheduleTaskBodyParams(d *schema.ResourceData) map[string]interface{} {
	return map[string]interface{}{
		// Required parameters.
		"task_name":      d.Get("task_name"),
		"task_type":      d.Get("task_type"),
		"scheduled_type": d.Get("scheduled_type"),
		"scheduled_time": d.Get("scheduled_time"),
		"target_infos":   buildAppScheduleTaskTargetInfos(d.Get("target_infos").(*schema.Set)),
		// Optional parameters.
		"day_interval":         utils.ValueIgnoreEmpty(d.Get("day_interval")),
		"week_list":            d.Get("week_list"),
		"month_list":           d.Get("month_list"),
		"date_list":            d.Get("date_list"),
		"time_zone":            d.Get("time_zone"),
		"scheduled_date":       d.Get("scheduled_date"),
		"expire_time":          d.Get("expire_time"),
		"description":          d.Get("description"),
		"schedule_task_policy": buildAppScheduleTaskPolicy(d.Get("schedule_task_policy").([]interface{})),
		"is_enable":            d.Get("is_enable"),
	}
}

func updateAppScheduleTask(client *golangsdk.ServiceClient, scheduleTaskId string, params map[string]interface{}) error {
	updatePath := "v1/{project_id}/schedule-task/{task_id}"
	updatePath = client.Endpoint + updatePath
	updatePath = strings.ReplaceAll(updatePath, "{project_id}", client.ProjectID)
	updatePath = strings.ReplaceAll(updatePath, "{task_id}", scheduleTaskId)
	updateOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		JSONBody:         utils.RemoveNil(params),
	}
	_, err := client.Request("PATCH", updatePath, &updateOpt)
	return err
}

func resourceAppScheduleTaskUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	scheduleTaskId := d.Id()
	client, err := cfg.NewServiceClient("appstream", cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating Workspace APP client: %s", err)
	}

	err = updateAppScheduleTask(client, scheduleTaskId, buildUpdateAppScheduleTaskBodyParams(d))
	if err != nil {
		return diag.Errorf("error updating Workspace APP schedule task (%s): %s", scheduleTaskId, err)
	}

	return resourceAppScheduleTaskRead(ctx, d, meta)
}

func resourceAppScheduleTaskDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg            = meta.(*config.Config)
		deletePath     = "v1/{project_id}/schedule-task/{task_id}"
		scheduleTaskId = d.Id()
	)

	client, err := cfg.NewServiceClient("appstream", cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating Workspace APP client: %s", err)
	}

	deletePath = client.Endpoint + deletePath
	deletePath = strings.ReplaceAll(deletePath, "{project_id}", client.ProjectID)
	deletePath = strings.ReplaceAll(deletePath, "{task_id}", scheduleTaskId)
	deleteOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}
	_, err = client.Request("DELETE", deletePath, &deleteOpt)
	if err != nil {
		return common.CheckDeletedDiag(d, err, fmt.Sprintf("error deleting Workspace APP schedule task (%s)", scheduleTaskId))
	}

	return nil
}
