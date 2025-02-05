package ces

import (
	"context"
	"fmt"
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

var oneClickAlarmNonUpdatableParams = []string{"one_click_alarm_id", "dimension_names", "dimension_names.*.metric", "dimension_names.*.event"}

// @API CES POST /v2/{project_id}/one-click-alarms
// @API CES GET /v2/{project_id}/one-click-alarms
// @API CES PUT /v2/{project_id}/one-click-alarms/{one_click_alarm_id}/notifications
// @API CES POST /v2/{project_id}/one-click-alarms/batch-delete
func ResourceOneClickAlarm() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceOneClickAlarmCreate,
		UpdateContext: resourceOneClickAlarmUpdate,
		ReadContext:   resourceOneClickAlarmRead,
		DeleteContext: resourceOneClickAlarmDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		CustomizeDiff: config.FlexibleForceNew(oneClickAlarmNonUpdatableParams),

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"one_click_alarm_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `Specifies the default one-click monitoring ID.`,
			},
			"dimension_names": {
				Type:        schema.TypeList,
				Required:    true,
				MaxItems:    1,
				Elem:        dimensionNamesSchema(),
				Description: `Specifies dimensions in metric and event alarm rules that have one-click monitoring enabled.`,
			},
			"notification_enabled": {
				Type:        schema.TypeBool,
				Required:    true,
				Description: `Specifies whether to enable the alarm notification.`,
			},
			"alarm_notifications": {
				Type:        schema.TypeList,
				Optional:    true,
				Elem:        SMNActionSchema(),
				Description: `Specifies the action to be triggered by an alarm.`,
			},
			"ok_notifications": {
				Type:        schema.TypeList,
				Optional:    true,
				Elem:        SMNActionSchema(),
				Description: `Specifies the action to be triggered after an alarm is cleared.`,
			},
			"notification_begin_time": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specifies the time when the alarm notification was enabled.`,
			},
			"notification_end_time": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specifies the time when the alarm notification was disabled.`,
			},
			"enable_force_new": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"true", "false"}, false),
				Description:  utils.SchemaDesc("", utils.SchemaDescInput{Internal: true}),
			},
			"namespace": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The metric namespace.`,
			},
			"description": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The supplementary information about one-click monitoring.`,
			},
			"enabled": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: `Whether to enable one-click monitoring.`,
			},
		},
	}
}

func dimensionNamesSchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"metric": {
				Type:         schema.TypeList,
				Optional:     true,
				Elem:         &schema.Schema{Type: schema.TypeString},
				AtLeastOneOf: []string{"dimension_names.0.metric", "dimension_names.0.event"},
				Description:  `Specifies dimensions in metric alarm rules that have one-click monitoring enabled.`,
			},
			"event": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: `Specifies whether to enable the event alarm rules.`,
			},
		},
	}
	return &sc
}

func SMNActionSchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"type": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `Specifies the notification type.`,
			},
			"notification_list": {
				Type:        schema.TypeList,
				Required:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: `Specifies the list of objects to be notified if the alarm status changes.`,
			},
		},
	}
	return &sc
}

func resourceOneClickAlarmCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var (
		createOneClickAlarmHttpUrl = "v2/{project_id}/one-click-alarms"
		createOneClickAlarmProduct = "ces"
	)
	client, err := cfg.NewServiceClient(createOneClickAlarmProduct, region)
	if err != nil {
		return diag.Errorf("error creating CES client: %s", err)
	}

	createOneClickAlarmPath := client.Endpoint + createOneClickAlarmHttpUrl
	createOneClickAlarmPath = strings.ReplaceAll(createOneClickAlarmPath, "{project_id}", client.ProjectID)

	createOneClickAlarmOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}
	createOneClickAlarmOpt.JSONBody = utils.RemoveNil(buildCreateOneClickAlarmBodyParams(d))
	createOneClickAlarmResp, err := client.Request("POST", createOneClickAlarmPath, &createOneClickAlarmOpt)
	if err != nil {
		return diag.Errorf("error creating CES one-click alarm: %s", err)
	}

	createOneClickAlarmRespBody, err := utils.FlattenResponse(createOneClickAlarmResp)
	if err != nil {
		return diag.FromErr(err)
	}

	id := utils.PathSearch("one_click_alarm_id", createOneClickAlarmRespBody, "").(string)
	if id == "" {
		return diag.Errorf("error creating CES one-click alarm: ID is not found in API response")
	}
	d.SetId(id)

	return resourceOneClickAlarmRead(ctx, d, meta)
}

func buildCreateOneClickAlarmBodyParams(d *schema.ResourceData) map[string]interface{} {
	return map[string]interface{}{
		"one_click_alarm_id":      d.Get("one_click_alarm_id"),
		"dimension_names":         buildDimensionNamesBodyParams(d.Get("dimension_names")),
		"notification_enabled":    d.Get("notification_enabled"),
		"alarm_notifications":     buildSMNActionBodyParams(d.Get("alarm_notifications")),
		"ok_notifications":        buildSMNActionBodyParams(d.Get("ok_notifications")),
		"notification_begin_time": utils.ValueIgnoreEmpty(d.Get("notification_begin_time")),
		"notification_end_time":   utils.ValueIgnoreEmpty(d.Get("notification_end_time")),
	}
}

func buildDimensionNamesBodyParams(rawParam interface{}) map[string]interface{} {
	if rawArray, ok := rawParam.([]interface{}); ok {
		if len(rawArray) == 0 {
			return nil
		}
		raw := rawArray[0].(map[string]interface{})
		param := map[string]interface{}{
			"metric": utils.ValueIgnoreEmpty(raw["metric"]),
		}
		hasEvent := raw["event"].(bool)
		if hasEvent {
			param["event"] = []string{""}
		}
		return param
	}
	return nil
}

func buildSMNActionBodyParams(rawParam interface{}) []map[string]interface{} {
	if rawArray, ok := rawParam.([]interface{}); ok && len(rawArray) > 0 {
		rst := make([]map[string]interface{}, len(rawArray))
		for i, v := range rawArray {
			raw := v.(map[string]interface{})
			rst[i] = map[string]interface{}{
				"type":              utils.ValueIgnoreEmpty(raw["type"]),
				"notification_list": utils.ValueIgnoreEmpty(raw["notification_list"]),
			}
		}
		return rst
	}
	return nil
}

func resourceOneClickAlarmRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var mErr *multierror.Error

	getOneClickAlarmProduct := "ces"
	client, err := cfg.NewServiceClient(getOneClickAlarmProduct, region)
	if err != nil {
		return diag.Errorf("error creating CES client: %s", err)
	}

	alarm, err := GetOneClickAlarm(client, d.Id())
	if err != nil {
		return common.CheckDeletedDiag(d, err, "no CES one-click alarm found")
	}

	mErr = multierror.Append(
		mErr,
		d.Set("region", region),
		d.Set("namespace", utils.PathSearch("namespace", alarm, nil)),
		d.Set("description", utils.PathSearch("description", alarm, nil)),
		d.Set("enabled", utils.PathSearch("enabled", alarm, false)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func GetOneClickAlarm(client *golangsdk.ServiceClient, id string) (interface{}, error) {
	getOneClickAlarmHttpUrl := "v2/{project_id}/one-click-alarms"
	getOneClickAlarmPath := client.Endpoint + getOneClickAlarmHttpUrl
	getOneClickAlarmPath = strings.ReplaceAll(getOneClickAlarmPath, "{project_id}", client.ProjectID)

	getOneClickAlarmOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}
	getOneClickAlarmResp, err := client.Request("GET", getOneClickAlarmPath, &getOneClickAlarmOpt)
	if err != nil {
		return nil, fmt.Errorf("error retrieve CES one-click alarm: %s", err)
	}

	getOneClickAlarmRespBody, err := utils.FlattenResponse(getOneClickAlarmResp)
	if err != nil {
		return nil, err
	}

	findAlarmStr := fmt.Sprintf("one_click_alarms[?one_click_alarm_id=='%s']|[0]", id)
	alarm := utils.PathSearch(findAlarmStr, getOneClickAlarmRespBody, nil)
	if alarm == nil {
		return nil, golangsdk.ErrDefault404{}
	}

	return alarm, nil
}

func resourceOneClickAlarmUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	alarmId := d.Id()
	var (
		updateOneClickAlarmHttpUrl = "v2/{project_id}/one-click-alarms/{one_click_alarm_id}/notifications"
		updateOneClickAlarmProduct = "ces"
	)
	updateOneClickAlarmClient, err := cfg.NewServiceClient(updateOneClickAlarmProduct, region)
	if err != nil {
		return diag.Errorf("error creating CES client: %s", err)
	}

	updateOneClickAlarmPath := updateOneClickAlarmClient.Endpoint + updateOneClickAlarmHttpUrl
	updateOneClickAlarmPath = strings.ReplaceAll(updateOneClickAlarmPath, "{project_id}", updateOneClickAlarmClient.ProjectID)
	updateOneClickAlarmPath = strings.ReplaceAll(updateOneClickAlarmPath, "{one_click_alarm_id}", alarmId)

	updateOneClickAlarmOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		OkCodes:          []int{204},
	}
	updateOneClickAlarmOpt.JSONBody = utils.RemoveNil(buildUpdateOneClickAlarmBodyParams(d))
	_, err = updateOneClickAlarmClient.Request("PUT", updateOneClickAlarmPath, &updateOneClickAlarmOpt)
	if err != nil {
		return diag.Errorf("error updating CES one-click alarm: %s", err)
	}

	return resourceOneClickAlarmRead(ctx, d, meta)
}

func buildUpdateOneClickAlarmBodyParams(d *schema.ResourceData) map[string]interface{} {
	return map[string]interface{}{
		"notification_enabled":    d.Get("notification_enabled"),
		"alarm_notifications":     buildSMNActionBodyParams(d.Get("alarm_notifications")),
		"ok_notifications":        buildSMNActionBodyParams(d.Get("ok_notifications")),
		"notification_begin_time": utils.ValueIgnoreEmpty(d.Get("notification_begin_time")),
		"notification_end_time":   utils.ValueIgnoreEmpty(d.Get("notification_end_time")),
	}
}

func resourceOneClickAlarmDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var (
		deleteOneClickAlarmHttpUrl = "v2/{project_id}/one-click-alarms/batch-delete"
		deleteOneClickAlarmProduct = "ces"
	)
	deleteOneClickAlarmClient, err := cfg.NewServiceClient(deleteOneClickAlarmProduct, region)
	if err != nil {
		return diag.Errorf("error creating CES client: %s", err)
	}

	deleteOneClickAlarmPath := deleteOneClickAlarmClient.Endpoint + deleteOneClickAlarmHttpUrl
	deleteOneClickAlarmPath = strings.ReplaceAll(deleteOneClickAlarmPath, "{project_id}", deleteOneClickAlarmClient.ProjectID)

	deleteOneClickAlarmOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}
	deleteOneClickAlarmOpt.JSONBody = buildDeleteOneClickAlarmBodyParams(d)
	_, err = deleteOneClickAlarmClient.Request("POST", deleteOneClickAlarmPath, &deleteOneClickAlarmOpt)
	if err != nil {
		return diag.Errorf("error deleting CES one-click alarm: %s", err)
	}

	// Successful deletion API call does not guarantee that the resource is successfully deleted.
	// Call the details API to confirm that the resource has been successfully deleted.
	_, err = GetOneClickAlarm(deleteOneClickAlarmClient, d.Id())
	if err != nil {
		return common.CheckDeletedDiag(d, err, "error retrieving CES one-click alarm")
	}

	return diag.Errorf("error deleting CES one-click alarm")
}

func buildDeleteOneClickAlarmBodyParams(d *schema.ResourceData) map[string]interface{} {
	return map[string]interface{}{
		"one_click_alarm_ids": []string{d.Id()},
	}
}
