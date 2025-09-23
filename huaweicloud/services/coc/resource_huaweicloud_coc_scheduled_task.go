package coc

import (
	"context"
	"log"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API COC POST /v1/schedule/task
// @API COC PUT /v1/schedule/task/{task_id}
// @API COC GET /v1/schedule/task/{task_id}
// @API COC DELETE /v1/schedule/task/{task_id}
// @API COC POST /v1/schedule/task/{task_id}/enable
// @API COC POST /v1/schedule/task/{task_id}/disable
// @API COC GET /v1/schedule/task
func ResourceScheduledTask() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceScheduledTaskCreate,
		ReadContext:   resourceScheduledTaskRead,
		UpdateContext: resourceScheduledTaskUpdate,
		DeleteContext: resourceScheduledTaskDelete,

		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"version_no": {
				Type:     schema.TypeString,
				Required: true,
			},
			"trigger_time": {
				Type:     schema.TypeList,
				Required: true,
				MaxItems: 1,
				Elem:     scheduledTaskTriggerTimeSchema(),
			},
			"task_type": {
				Type:     schema.TypeString,
				Required: true,
			},
			"associated_task_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"associated_task_type": {
				Type:     schema.TypeString,
				Required: true,
			},
			"associated_task_name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"risk_level": {
				Type:     schema.TypeString,
				Required: true,
			},
			"input_param": {
				Type:     schema.TypeMap,
				Required: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"target_instances": {
				Type:     schema.TypeList,
				Required: true,
				Elem:     scheduledTaskScheduleInstanceSchema(),
			},
			"enable_approve": {
				Type:     schema.TypeBool,
				Required: true,
			},
			"enable_message_notification": {
				Type:     schema.TypeBool,
				Required: true,
			},
			"ticket_infos": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     scheduledTaskTicketInfoSchema(),
			},
			"enterprise_project_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"agency_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"associated_task_name_en": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"associated_task_enterprise_project_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"runbook_instance_mode": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"reviewer_notification": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem:     scheduledTaskMessageNotificationSchema(),
			},
			"reviewer_user_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"message_notification": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem:     scheduledTaskMessageNotificationSchema(),
			},
			"enabled": {
				Type:         schema.TypeString,
				ValidateFunc: validation.StringInSlice([]string{"true", "false"}, false),
				Optional:     true,
				Computed:     true,
			},
			"created_user_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"approve_status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"approve_comments": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func scheduledTaskMessageNotificationSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"notification_endpoint_type": {
				Type:     schema.TypeString,
				Required: true,
			},
			"policy": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"schedule_scene_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"schedule_role_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"recipients": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"protocol": {
				Type:     schema.TypeString,
				Optional: true,
			},
		},
	}
}

func scheduledTaskTicketInfoSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"ticket_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"ticket_type": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"target_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"scope_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
		},
	}
}

func scheduledTaskTriggerTimeSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"time_zone": {
				Type:     schema.TypeString,
				Required: true,
			},
			"policy": {
				Type:     schema.TypeString,
				Required: true,
			},
			"single_scheduled_time": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"periodic_scheduled_time": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"period": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"cron": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"scheduled_close_time": {
				Type:     schema.TypeInt,
				Optional: true,
			},
		},
	}
}

func scheduledTaskScheduleInstanceSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"target_selection": {
				Type:     schema.TypeString,
				Required: true,
			},
			"order_no": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"target_resource": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem:     scheduledTaskTargetResourceSchema(),
			},
			"target_instances": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"batch_strategy": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"sub_target_instances": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"target_selection": {
							Type:     schema.TypeString,
							Required: true,
						},
						"order_no": {
							Type:     schema.TypeInt,
							Required: true,
						},
						"target_resource": {
							Type:     schema.TypeList,
							Optional: true,
							MaxItems: 1,
							Elem:     scheduledTaskTargetResourceSchema(),
						},
						"target_instances": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"batch_strategy": {
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
			},
			"id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"schedule_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"create_time": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"update_time": {
				Type:     schema.TypeInt,
				Computed: true,
			},
		},
	}
}

func scheduledTaskTargetResourceSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"type": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"app_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"region_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"params": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     scheduledTaskResourceQuerySchema(),
			},
		},
	}
}

func scheduledTaskResourceQuerySchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"key": {
				Type:     schema.TypeString,
				Required: true,
			},
			"value": {
				Type:     schema.TypeString,
				Required: true,
			},
		},
	}
}

func resourceScheduledTaskCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)

	client, err := cfg.NewServiceClient("coc", cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating COC client: %s", err)
	}

	createHttpUrl := "v1/schedule/task"
	createPath := client.Endpoint + createHttpUrl
	createOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		JSONBody:         utils.RemoveNil(buildCreateAndUpdateScheduledTaskBodyParams(d)),
	}

	createResp, err := client.Request("POST", createPath, &createOpt)
	if err != nil {
		return diag.Errorf("error creating COC scheduled task: %s", err)
	}

	createRespBody, err := utils.FlattenResponse(createResp)
	if err != nil {
		return diag.Errorf("error flattening COC scheduled task response: %s", err)
	}

	id := utils.PathSearch("data", createRespBody, "").(string)
	if id == "" {
		return diag.Errorf("unable to find the COC scheduled task ID from the API response")
	}

	d.SetId(id)

	if v, ok := d.GetOk("enabled"); ok {
		enabled := v.(string) == "true"
		enableApprove := d.Get("enable_approve").(bool)
		err = postUpdateScheduledTaskEnabled(client, id, d, enabled, enableApprove)
		if err != nil {
			return diag.Errorf("error updating COC scheduled task enabled status: %s", err)
		}
	}

	return resourceScheduledTaskRead(ctx, d, meta)
}

func postUpdateScheduledTaskEnabled(client *golangsdk.ServiceClient, taskId string, d *schema.ResourceData,
	enabled bool, enableApprove bool) error {
	if enabled && enableApprove {
		return enableScheduledTask(client, taskId, d)
	}
	if !enabled && !enableApprove {
		return disableScheduledTask(client, taskId)
	}
	return nil
}

func disableScheduledTask(client *golangsdk.ServiceClient, taskId string) error {
	createHttpUrl := "v1/schedule/task/{task_id}/disable"
	createPath := client.Endpoint + createHttpUrl
	createPath = strings.ReplaceAll(createPath, "{task_id}", taskId)
	createOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json",
		},
	}

	_, err := client.Request("POST", createPath, &createOpt)
	return err
}

func enableScheduledTask(client *golangsdk.ServiceClient, taskId string, d *schema.ResourceData) error {
	createHttpUrl := "v1/schedule/task/{task_id}/enable"
	createPath := client.Endpoint + createHttpUrl
	createPath = strings.ReplaceAll(createPath, "{task_id}", taskId)
	createOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		JSONBody:         utils.RemoveNil(buildEnabledScheduledTaskBodyParams(d)),
	}

	_, err := client.Request("POST", createPath, &createOpt)
	return err
}

func buildEnabledScheduledTaskBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"ticket_infos": buildScheduledTaskTicketInfosBodyParams(d.Get("ticket_infos")),
	}

	return bodyParams
}

func buildCreateAndUpdateScheduledTaskBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"name":                                  d.Get("name"),
		"version_no":                            d.Get("version_no"),
		"trigger_time":                          buildScheduledTaskTriggerTimeBodyParams(d.Get("trigger_time")),
		"task_type":                             d.Get("task_type"),
		"associated_task_id":                    d.Get("associated_task_id"),
		"associated_task_type":                  d.Get("associated_task_type"),
		"associated_task_name":                  d.Get("associated_task_name"),
		"risk_level":                            d.Get("risk_level"),
		"input_param":                           d.Get("input_param"),
		"target_instances":                      buildScheduledTaskTargetInstancesBodyParams(d.Get("target_instances")),
		"enable_approve":                        d.Get("enable_approve"),
		"enable_message_notification":           d.Get("enable_message_notification"),
		"ticket_infos":                          buildScheduledTaskTicketInfosBodyParams(d.Get("ticket_infos")),
		"enterprise_project_id":                 utils.ValueIgnoreEmpty(d.Get("enterprise_project_id")),
		"agency_name":                           utils.ValueIgnoreEmpty(d.Get("agency_name")),
		"associated_task_name_en":               utils.ValueIgnoreEmpty(d.Get("associated_task_name_en")),
		"associated_task_enterprise_project_id": utils.ValueIgnoreEmpty(d.Get("associated_task_enterprise_project_id")),
		"runbook_instance_mode":                 utils.ValueIgnoreEmpty(d.Get("runbook_instance_mode")),
		"reviewer_notification":                 buildScheduledTaskMessageNotificationBodyParams(d.Get("reviewer_notification")),
		"reviewer_user_name":                    utils.ValueIgnoreEmpty(d.Get("reviewer_user_name")),
		"message_notification":                  buildScheduledTaskMessageNotificationBodyParams(d.Get("message_notification")),
	}

	return bodyParams
}

func buildScheduledTaskTicketInfosBodyParams(rawParams interface{}) []map[string]interface{} {
	if rawArray, ok := rawParams.([]interface{}); ok {
		if len(rawArray) == 0 {
			return nil
		}

		params := make([]map[string]interface{}, len(rawArray))
		for i, v := range rawArray {
			raw := v.(map[string]interface{})
			params[i] = map[string]interface{}{
				"ticket_id":   utils.ValueIgnoreEmpty(raw["ticket_id"]),
				"ticket_type": utils.ValueIgnoreEmpty(raw["ticket_type"]),
				"target_id":   utils.ValueIgnoreEmpty(raw["target_id"]),
				"scope_id":    utils.ValueIgnoreEmpty(raw["scope_id"]),
			}
		}
		return params
	}

	return nil
}

func buildScheduledTaskMessageNotificationBodyParams(rawParam interface{}) map[string]interface{} {
	if rawArray, ok := rawParam.([]interface{}); ok {
		if len(rawArray) != 1 {
			return nil
		}

		raw := rawArray[0].(map[string]interface{})
		param := map[string]interface{}{
			"notification_endpoint_type": raw["notification_endpoint_type"],
			"policy":                     utils.ValueIgnoreEmpty(raw["policy"]),
			"schedule_scene_id":          utils.ValueIgnoreEmpty(raw["schedule_scene_id"]),
			"schedule_role_id":           utils.ValueIgnoreEmpty(raw["schedule_role_id"]),
			"recipients":                 utils.ValueIgnoreEmpty(raw["recipients"]),
			"protocol":                   utils.ValueIgnoreEmpty(raw["protocol"]),
		}

		return param
	}

	return nil
}

func buildScheduledTaskTriggerTimeBodyParams(rawParam interface{}) map[string]interface{} {
	if rawArray, ok := rawParam.([]interface{}); ok {
		if len(rawArray) != 1 {
			return nil
		}

		raw := rawArray[0].(map[string]interface{})
		param := map[string]interface{}{
			"time_zone":               raw["time_zone"],
			"policy":                  raw["policy"],
			"single_scheduled_time":   utils.ValueIgnoreEmpty(raw["single_scheduled_time"]),
			"periodic_scheduled_time": utils.ValueIgnoreEmpty(raw["periodic_scheduled_time"]),
			"period":                  utils.ValueIgnoreEmpty(raw["period"]),
			"cron":                    utils.ValueIgnoreEmpty(raw["cron"]),
			"scheduled_close_time":    utils.ValueIgnoreEmpty(raw["scheduled_close_time"]),
		}

		return param
	}

	return nil
}

func buildScheduledTaskTargetInstancesBodyParams(rawParams interface{}) []map[string]interface{} {
	if rawArray, ok := rawParams.([]interface{}); ok {
		if len(rawArray) == 0 {
			return nil
		}

		params := make([]map[string]interface{}, len(rawArray))
		for i, v := range rawArray {
			raw := v.(map[string]interface{})
			params[i] = map[string]interface{}{
				"target_selection": raw["target_selection"],
				"order_no":         raw["order_no"],
				"target_resource":  buildScheduledTaskTargetInstancesTargetResourceBodyParams(raw["target_resource"]),
				"target_instances": utils.ValueIgnoreEmpty(raw["target_instances"]),
				"batch_strategy":   utils.ValueIgnoreEmpty(raw["batch_strategy"]),
				"sub_target_instances": buildScheduledTaskTargetInstancesSubTargetInstancesBodyParams(
					raw["sub_target_instances"]),
			}
		}
		return params
	}

	return nil
}

func buildScheduledTaskTargetInstancesSubTargetInstancesBodyParams(rawParams interface{}) []map[string]interface{} {
	if rawArray, ok := rawParams.([]interface{}); ok {
		if len(rawArray) == 0 {
			return nil
		}

		params := make([]map[string]interface{}, len(rawArray))
		for i, v := range rawArray {
			raw := v.(map[string]interface{})
			params[i] = map[string]interface{}{
				"target_selection": raw["target_selection"],
				"order_no":         raw["order_no"],
				"target_resource":  buildScheduledTaskTargetInstancesTargetResourceBodyParams(raw["target_resource"]),
				"target_instances": utils.ValueIgnoreEmpty(raw["target_instances"]),
				"batch_strategy":   utils.ValueIgnoreEmpty(raw["batch_strategy"]),
			}
		}
		return params
	}

	return nil
}

func buildScheduledTaskTargetInstancesTargetResourceBodyParams(rawParam interface{}) map[string]interface{} {
	if rawArray, ok := rawParam.([]interface{}); ok {
		if len(rawArray) != 1 {
			return nil
		}

		raw := rawArray[0].(map[string]interface{})
		param := map[string]interface{}{
			"type":      utils.ValueIgnoreEmpty(raw["type"]),
			"id":        utils.ValueIgnoreEmpty(raw["id"]),
			"app_name":  utils.ValueIgnoreEmpty(raw["app_name"]),
			"region_id": utils.ValueIgnoreEmpty(raw["region_id"]),
			"params":    buildScheduledTaskTargetInstancesTargetResourceParamsBodyParams(raw["params"]),
		}

		return param
	}

	return nil
}

func buildScheduledTaskTargetInstancesTargetResourceParamsBodyParams(rawParams interface{}) []map[string]interface{} {
	if rawArray, ok := rawParams.([]interface{}); ok {
		if len(rawArray) == 0 {
			return nil
		}

		params := make([]map[string]interface{}, len(rawArray))
		for i, v := range rawArray {
			raw := v.(map[string]interface{})
			params[i] = map[string]interface{}{
				"key":   raw["key"],
				"value": raw["value"],
			}
		}
		return params
	}

	return nil
}

func resourceScheduledTaskRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	var mErr *multierror.Error

	client, err := cfg.NewServiceClient("coc", cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating COC client: %s", err)
	}

	scheduledTask, err := GetScheduledTask(client, d.Id())
	if err != nil {
		return common.CheckDeletedDiag(d, common.ConvertExpected400ErrInto404Err(err, "error_code", "COC.00014102"),
			"error retrieving scheduled task")
	}
	currentEnabled, err := getScheduledTaskEnabled(client, d.Id())
	if err != nil {
		log.Printf("[ERROR] error retrieving COC scheduled task enabled: %s", err)
	}

	inputParamMap := parseJson(utils.PathSearch("input_param", scheduledTask, "").(string))

	mErr = multierror.Append(mErr,
		d.Set("name", utils.PathSearch("name", scheduledTask, nil)),
		d.Set("version_no", utils.PathSearch("version_no", scheduledTask, nil)),
		d.Set("trigger_time", flattenCocScheduledTaskTriggerTime(
			utils.PathSearch("trigger_time", scheduledTask, nil))),
		d.Set("task_type", utils.PathSearch("task_type", scheduledTask, nil)),
		d.Set("associated_task_id", utils.PathSearch("associated_task_id", scheduledTask, nil)),
		d.Set("associated_task_type", utils.PathSearch("associated_task_type", scheduledTask, nil)),
		d.Set("associated_task_name", utils.PathSearch("associated_task_name", scheduledTask, nil)),
		d.Set("risk_level", utils.PathSearch("risk_level", scheduledTask, nil)),
		d.Set("input_param", inputParamMap),
		d.Set("target_instances", parseJson(utils.PathSearch("target_instances", scheduledTask, "").(string))),
		d.Set("enable_approve", utils.PathSearch("enable_approve", scheduledTask, nil)),
		d.Set("enable_message_notification", utils.PathSearch("enable_message_notification", scheduledTask, nil)),
		d.Set("enterprise_project_id", utils.PathSearch("enterprise_project_id", scheduledTask, nil)),
		d.Set("agency_name", utils.PathSearch("agency_name", scheduledTask, nil)),
		d.Set("associated_task_name_en", utils.PathSearch("associated_task_name_en", scheduledTask, nil)),
		d.Set("runbook_instance_mode", utils.PathSearch("runbook_instance_mode", scheduledTask, nil)),
		d.Set("reviewer_notification", flattenCocScheduledTaskMessageNotification(
			utils.PathSearch("reviewer_notification", scheduledTask, nil))),
		d.Set("reviewer_user_name", utils.PathSearch("reviewer_user_name", scheduledTask, nil)),
		d.Set("message_notification", flattenCocScheduledTaskMessageNotification(
			utils.PathSearch("message_notification", scheduledTask, nil))),
		d.Set("created_user_name", utils.PathSearch("created_user_name", scheduledTask, nil)),
		d.Set("approve_status", utils.PathSearch("approve_status", scheduledTask, nil)),
		d.Set("approve_comments", utils.PathSearch("approve_comments", scheduledTask, nil)),
		d.Set("enabled", currentEnabled),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenCocScheduledTaskTriggerTime(param interface{}) interface{} {
	if param == nil {
		return nil
	}
	rst := []map[string]interface{}{
		{
			"time_zone":               utils.PathSearch("time_zone", param, nil),
			"policy":                  utils.PathSearch("policy", param, nil),
			"single_scheduled_time":   utils.PathSearch("single_scheduled_time", param, nil),
			"periodic_scheduled_time": utils.PathSearch("periodic_scheduled_time", param, nil),
			"period":                  utils.PathSearch("period", param, nil),
			"cron":                    utils.PathSearch("cron", param, nil),
			"scheduled_close_time":    utils.PathSearch("scheduled_close_time", param, nil),
		},
	}

	return rst
}

func flattenCocScheduledTaskMessageNotification(param interface{}) interface{} {
	if param == nil {
		return nil
	}
	rst := []map[string]interface{}{
		{
			"policy":                     utils.PathSearch("policy", param, nil),
			"notification_endpoint_type": utils.PathSearch("notification_endpoint_type", param, nil),
			"schedule_scene_id":          utils.PathSearch("schedule_scene_id", param, nil),
			"schedule_role_id":           utils.PathSearch("schedule_role_id", param, nil),
			"recipients":                 utils.PathSearch("recipients", param, nil),
			"protocol":                   utils.PathSearch("protocol", param, nil),
		},
	}

	return rst
}

func GetScheduledTask(client *golangsdk.ServiceClient, scheduledTaskID string) (interface{}, error) {
	requestPath := client.Endpoint + "v1/schedule/task/{task_id}"
	requestPath = strings.ReplaceAll(requestPath, "{task_id}", scheduledTaskID)
	requestOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json",
		},
	}

	resp, err := client.Request("GET", requestPath, &requestOpt)
	if err != nil {
		return nil, err
	}

	return utils.FlattenResponse(resp)
}

func getScheduledTaskEnabled(client *golangsdk.ServiceClient, scheduledTaskID string) (interface{}, error) {
	requestPath := client.Endpoint + "v1/schedule/task?task_id={task_id}"
	requestPath = strings.ReplaceAll(requestPath, "{task_id}", scheduledTaskID)
	requestOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json",
		},
	}

	resp, err := client.Request("GET", requestPath, &requestOpt)
	if err != nil {
		return nil, err
	}

	respBody, err := utils.FlattenResponse(resp)
	if err != nil {
		return nil, err
	}

	enabled := utils.PathSearch("scheduled_tasks[0].enabled", respBody, nil)
	if enabled == nil {
		return nil, golangsdk.ErrDefault404{}
	}

	if enabled.(bool) {
		return "true", nil
	}
	return "false", nil
}

func resourceScheduledTaskUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)

	client, err := cfg.NewServiceClient("coc", cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating COC client: %s", err)
	}

	if d.HasChange("enabled") && d.Get("enabled") == "false" {
		err = disableScheduledTask(client, d.Id())
		if err != nil {
			return diag.Errorf("error updating COC scheduled task enabled status: %s", err)
		}
	}

	changeList := []string{
		"name", "version_no", "trigger_time", "task_type", "associated_task_id", "associated_task_type",
		"associated_task_name", "input_param", "target_instances", "enable_approve", "enable_message_notification",
		"ticket_infos", "enterprise_project_id", "agency_name", "associated_task_name_en",
		"associated_task_enterprise_project_id", "runbook_instance_mode", "reviewer_notification", "reviewer_user_name",
		"message_notification",
	}
	if d.HasChanges(changeList...) {
		updateHttpUrl := "v1/schedule/task/{task_id}"
		updatePath := client.Endpoint + updateHttpUrl
		updatePath = strings.ReplaceAll(updatePath, "{task_id}", d.Id())
		updateOpt := golangsdk.RequestOpts{
			KeepResponseBody: true,
			JSONBody:         utils.RemoveNil(buildCreateAndUpdateScheduledTaskBodyParams(d)),
		}

		_, err = client.Request("PUT", updatePath, &updateOpt)
		if err != nil {
			return diag.Errorf("error updating scheduled task: %s", err)
		}
	}

	if v, ok := d.GetOk("enabled"); ok {
		enabled := v.(string) == "true"
		enableApprove := d.Get("enable_approve").(bool)
		if d.HasChange("enabled") && enabled {
			if !d.HasChanges(changeList...) || (d.HasChanges(changeList...) && enableApprove) {
				err = enableScheduledTask(client, d.Id(), d)
				if err != nil {
					return diag.Errorf("error updating COC scheduled task enabled status: %s", err)
				}
			}
		}

		if (!d.HasChange("enabled") || !enabled) && d.HasChanges(changeList...) {
			err = postUpdateScheduledTaskEnabled(client, d.Id(), d, enabled, enableApprove)
			if err != nil {
				return diag.Errorf("error updating COC scheduled task enabled status: %s", err)
			}
		}
	}

	return resourceScheduledTaskRead(ctx, d, meta)
}

func resourceScheduledTaskDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)

	client, err := cfg.NewServiceClient("coc", cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating COC client: %s", err)
	}

	if v, ok := d.GetOk("enabled"); ok && v.(string) == "true" {
		err = disableScheduledTask(client, d.Id())
		if err != nil {
			return diag.Errorf("error updating COC scheduled task enabled status: %s", err)
		}
	}

	deleteHttpUrl := "v1/schedule/task/{task_id}"
	deletePath := client.Endpoint + deleteHttpUrl
	deletePath = strings.ReplaceAll(deletePath, "{task_id}", d.Id())
	deleteOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json",
		},
	}

	_, err = client.Request("DELETE", deletePath, &deleteOpt)
	if err != nil {
		return common.CheckDeletedDiag(d, common.ConvertExpected400ErrInto404Err(err, "error_code", "COC.00014102"),
			"error deleting COC scheduled task")
	}

	return nil
}
