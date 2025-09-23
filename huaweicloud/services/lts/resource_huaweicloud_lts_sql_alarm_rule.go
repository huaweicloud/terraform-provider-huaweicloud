package lts

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

// @API LTS POST /v2/{project_id}/lts/alarms/sql-alarm-rule
// @API LTS PUT /v2/{project_id}/lts/alarms/status
// @API LTS GET /v2/{project_id}/lts/alarms/sql-alarm-rule
// @API LTS PUT /v2/{project_id}/lts/alarms/sql-alarm-rule
// @API LTS DELETE /v2/{project_id}/lts/alarms/sql-alarm-rule/{sql_alarm_rule_id}
func ResourceSQLAlarmRule() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceSQLAlarmRuleCreate,
		UpdateContext: resourceSQLAlarmRuleUpdate,
		ReadContext:   resourceSQLAlarmRuleRead,
		DeleteContext: resourceSQLAlarmRuleDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: `Specifies the name of the SQL alarm rule.`,
			},
			"sql_requests": {
				Type:        schema.TypeList,
				Elem:        sqlAlarmRuleSQLRequestsSchema(),
				Required:    true,
				Description: `Specifies the SQL requests.`,
			},
			"frequency": {
				Type:        schema.TypeList,
				MaxItems:    1,
				Elem:        sqlAlarmRuleFrequencySchema(),
				Required:    true,
				Description: `Specifies the alarm frequency configurations.`,
			},
			"condition_expression": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `Specifies the condition expression.`,
			},
			"alarm_level": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `Specifies the alarm level.`,
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specifies the description of the SQL alarm rule.`,
			},
			"send_notifications": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: `Specifies whether to send notifications.`,
			},
			"alarm_action_rule_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The name of the alarm action rule associated with the SQL alarm rule.`,
			},
			"notification_save_rule": {
				Type:        schema.TypeList,
				MaxItems:    1,
				Optional:    true,
				Elem:        sqlAlarmRuleNotificationRuleSchema(),
				Description: `The notification rule of the SQL alarm rule.`,
			},
			"trigger_condition_count": {
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
				Description: `Specifies the count to trigger the alarm.`,
			},
			"trigger_condition_frequency": {
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
				Description: `Specifies the frequency to trigger the alarm.`,
			},
			"send_recovery_notifications": {
				Type:        schema.TypeBool,
				Optional:    true,
				Computed:    true,
				Description: `Specifies whether to send recovery notifications.`,
			},
			"recovery_frequency": {
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
				Description: `Specifies the frequency to recover the alarm.`,
			},
			"alarm_rule_alias": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `The alias name of the SQL alarm rule.`,
			},
			"notification_frequency": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: `The notification frequency of the SQL alarm rule, in minutes.`,
			},
			"status": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `Specifies the status of the alarm rule.`,
			},
			"domain_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The domain ID.`,
			},
			"created_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The creation time of the alarm rule.`,
			},
			"updated_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The last update time of the alarm rule.`,
			},
			// Deprecated parameters.
			"notification_rule": {
				Type:     schema.TypeList,
				MaxItems: 1,
				Elem:     sqlAlarmRuleNotificationRuleSchema(),
				Optional: true,
				ForceNew: true,
				Computed: true,
				Description: utils.SchemaDesc(
					`The notification rule of the SQL alarm rule.`,
					utils.SchemaDescInput{
						Deprecated: true,
					},
				)},
		},
	}
}

func sqlAlarmRuleSQLRequestsSchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"title": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `Specifies the SQL request title.`,
			},
			"sql": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `Specifies the SQL.`,
			},
			"log_stream_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `Specifies the log stream id.`,
			},
			"log_group_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `Specifies the log group id.`,
			},
			"search_time_range_unit": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `Specifies the unit of search time range.`,
			},
			"search_time_range": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: `Specifies the search time range.`,
			},
			"is_time_range_relative": {
				Type:        schema.TypeBool,
				Optional:    true,
				Computed:    true,
				Description: `Specifies the SQL request is relative to time range.`,
			},
			"log_stream_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `The name of the log stream.`,
			},
			"log_group_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `The name of the log group.`,
			},
		},
	}
	return &sc
}

func sqlAlarmRuleFrequencySchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"type": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `Specifies the frequency type.`,
			},
			"cron_expression": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `Specifies the cron expression.`,
			},
			"hour_of_day": {
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
				Description: `Specifies the hour of day.`,
			},
			"day_of_week": {
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
				Description: `Specifies the day of week.`,
			},
			"fixed_rate_unit": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `Specifies the unit of fixed rate.`,
			},
			"fixed_rate": {
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
				Description: `Specifies the unit fixed rate.`,
			},
		},
	}
	return &sc
}

func sqlAlarmRuleNotificationRuleSchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"template_name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `Specifies the notification template name.`,
			},
			"language": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `Specifies the notification language.`,
			},
			"user_name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `Specifies the user name.`,
			},
			"topics": {
				Type:        schema.TypeList,
				Elem:        sqlAlarmRuleNotificationRuleTopicSchema(),
				Required:    true,
				Description: `Specifies the SMN topics.`,
			},
			"timezone": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specifies the timezone.`,
			},
		},
	}
	return &sc
}

func sqlAlarmRuleNotificationRuleTopicSchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `Specifies the topic name.`,
			},
			"topic_urn": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `Specifies the topic URN.`,
			},
			"display_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `Specifies the display name.`,
			},
			"push_policy": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: `Specifies the push policy.`,
			},
		},
	}
	return &sc
}

func resourceSQLAlarmRuleCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	// createSQLAlarmRule: Create an LTS SQLAlarmRule.
	var (
		createSQLAlarmRuleHttpUrl = "v2/{project_id}/lts/alarms/sql-alarm-rule"
		createSQLAlarmRuleProduct = "lts"
	)
	createSQLAlarmRuleClient, err := cfg.NewServiceClient(createSQLAlarmRuleProduct, region)
	if err != nil {
		return diag.Errorf("error creating LTS client: %s", err)
	}

	createSQLAlarmRulePath := createSQLAlarmRuleClient.Endpoint + createSQLAlarmRuleHttpUrl
	createSQLAlarmRulePath = strings.ReplaceAll(createSQLAlarmRulePath, "{project_id}", createSQLAlarmRuleClient.ProjectID)

	createSQLAlarmRuleOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"Content-Type": "application/json"},
	}

	createSQLAlarmRuleOpt.JSONBody = utils.RemoveNil(buildCreateSQLAlarmRuleBodyParams(d, cfg))
	createSQLAlarmRuleResp, err := createSQLAlarmRuleClient.Request("POST", createSQLAlarmRulePath, &createSQLAlarmRuleOpt)
	if err != nil {
		return diag.Errorf("error creating SQL alarm rule: %s", err)
	}

	createSQLAlarmRuleRespBody, err := utils.FlattenResponse(createSQLAlarmRuleResp)
	if err != nil {
		return diag.FromErr(err)
	}

	ruleId := utils.PathSearch("sql_alarm_rule_id", createSQLAlarmRuleRespBody, "").(string)
	if ruleId == "" {
		return diag.Errorf("unable to find the SQL alarm rule ID from the API response")
	}
	d.SetId(ruleId)

	if _, ok := d.GetOk("alarm_rule_alias"); ok {
		err = updateSQLAlarmRule(createSQLAlarmRuleClient, buildUpdateSQLAlarmRuleBodyParams(d, cfg))
		if err != nil {
			return diag.Errorf("error updating alias name of the SQL alarm rule (%s): %s", err, ruleId)
		}
	}

	if d.Get("status").(string) == "STOPPING" {
		err = updateAlarmRuleStatus(createSQLAlarmRuleClient, ruleId, "sql", "STOPPING")
		if err != nil {
			return diag.Errorf("error stopping SQL alarm rule: %s", err)
		}
	}

	return resourceSQLAlarmRuleRead(ctx, d, meta)
}

func buildCreateSQLAlarmRuleBodyParams(d *schema.ResourceData, cfg *config.Config) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"sql_alarm_rule_name":  d.Get("name"),
		"sql_requests":         buildSQLAlarmRuleRequestBodySQLRequests(d.Get("sql_requests")),
		"frequency":            buildSQLAlarmRuleRequestBodyFrequency(d.Get("frequency")),
		"condition_expression": d.Get("condition_expression"),
		"sql_alarm_level":      d.Get("alarm_level"),
		// Optional parameter.
		"sql_alarm_rule_description":  utils.ValueIgnoreEmpty(d.Get("description")),
		"sql_alarm_send":              utils.ValueIgnoreEmpty(d.Get("send_notifications")),
		"domain_id":                   cfg.DomainID,
		"alarm_action_rule_name":      utils.ValueIgnoreEmpty(d.Get("alarm_action_rule_name")),
		"notification_rule":           buildSQLAlarmRuleRequestBodyNotificationRule(d.Get("notification_rule")),
		"notification_save_rule":      buildSQLAlarmRuleRequestBodyNotificationRule(d.Get("notification_save_rule")),
		"trigger_condition_count":     utils.ValueIgnoreEmpty(d.Get("trigger_condition_count")),
		"trigger_condition_frequency": utils.ValueIgnoreEmpty(d.Get("trigger_condition_frequency")),
		"whether_recovery_policy":     utils.ValueIgnoreEmpty(d.Get("send_recovery_notifications")),
		"recovery_policy":             utils.ValueIgnoreEmpty(d.Get("recovery_frequency")),
		"notification_frequency":      d.Get("notification_frequency"),
	}
	return bodyParams
}

func buildSQLAlarmRuleRequestBodySQLRequests(rawParams interface{}) []map[string]interface{} {
	if rawArray, ok := rawParams.([]interface{}); ok {
		if len(rawArray) == 0 {
			return nil
		}

		rst := make([]map[string]interface{}, len(rawArray))
		for i, v := range rawArray {
			raw := v.(map[string]interface{})
			rst[i] = map[string]interface{}{
				"sql_request_title":      utils.ValueIgnoreEmpty(raw["title"]),
				"sql":                    utils.ValueIgnoreEmpty(raw["sql"]),
				"log_stream_id":          utils.ValueIgnoreEmpty(raw["log_stream_id"]),
				"log_group_id":           utils.ValueIgnoreEmpty(raw["log_group_id"]),
				"search_time_range_unit": utils.ValueIgnoreEmpty(raw["search_time_range_unit"]),
				"search_time_range":      utils.ValueIgnoreEmpty(raw["search_time_range"]),
				"is_time_range_relative": utils.ValueIgnoreEmpty(raw["is_time_range_relative"]),
				"log_group_name":         utils.ValueIgnoreEmpty(raw["log_group_name"]),
				"log_stream_name":        utils.ValueIgnoreEmpty(raw["log_stream_name"]),
			}
		}
		return rst
	}
	return nil
}

func buildSQLAlarmRuleRequestBodyFrequency(rawParams interface{}) map[string]interface{} {
	if rawArray, ok := rawParams.([]interface{}); ok {
		if len(rawArray) == 0 {
			return nil
		}
		raw, ok := rawArray[0].(map[string]interface{})
		if !ok {
			return nil
		}

		params := map[string]interface{}{
			"type":            utils.ValueIgnoreEmpty(raw["type"]),
			"cron_expr":       utils.ValueIgnoreEmpty(raw["cron_expression"]),
			"hour_of_day":     utils.ValueIgnoreEmpty(raw["hour_of_day"]),
			"day_of_week":     utils.ValueIgnoreEmpty(raw["day_of_week"]),
			"fixed_rate_unit": utils.ValueIgnoreEmpty(raw["fixed_rate_unit"]),
			"fixed_rate":      utils.ValueIgnoreEmpty(raw["fixed_rate"]),
		}
		return params
	}
	return nil
}

func buildSQLAlarmRuleRequestBodyNotificationRule(rawParams interface{}) map[string]interface{} {
	if rawArray, ok := rawParams.([]interface{}); ok {
		if len(rawArray) == 0 {
			return nil
		}
		raw, ok := rawArray[0].(map[string]interface{})
		if !ok {
			return nil
		}

		params := map[string]interface{}{
			"template_name": utils.ValueIgnoreEmpty(raw["template_name"]),
			"language":      utils.ValueIgnoreEmpty(raw["language"]),
			"user_name":     utils.ValueIgnoreEmpty(raw["user_name"]),
			"timezone":      utils.ValueIgnoreEmpty(raw["timezone"]),
			"topics":        buildNotificationRuleTopic(raw["topics"]),
		}
		return params
	}
	return nil
}

func buildNotificationRuleTopic(rawParams interface{}) []map[string]interface{} {
	if rawArray, ok := rawParams.([]interface{}); ok {
		if len(rawArray) == 0 {
			return nil
		}

		rst := make([]map[string]interface{}, len(rawArray))
		for i, v := range rawArray {
			raw := v.(map[string]interface{})
			rst[i] = map[string]interface{}{
				"name":         utils.ValueIgnoreEmpty(raw["name"]),
				"topic_urn":    utils.ValueIgnoreEmpty(raw["topic_urn"]),
				"display_name": utils.ValueIgnoreEmpty(raw["display_name"]),
				"push_policy":  raw["push_policy"],
			}
		}
		return rst
	}
	return nil
}

// GetSQLAlarmRuleById is a method used to get SQL alarm rule detail by its ID.
func GetSQLAlarmRuleById(client *golangsdk.ServiceClient, alarmRuleId string) (interface{}, error) {
	httpUrl := "v2/{project_id}/lts/alarms/sql-alarm-rule"
	getPath := client.Endpoint + httpUrl
	getPath = strings.ReplaceAll(getPath, "{project_id}", client.ProjectID)
	getOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"Content-Type": "application/json"},
	}

	getResp, err := client.Request("GET", getPath, &getOpt)
	if err != nil {
		return nil, err
	}

	respBody, err := utils.FlattenResponse(getResp)
	if err != nil {
		return nil, err
	}

	alarmRule := utils.PathSearch(fmt.Sprintf("sql_alarm_rules[?sql_alarm_rule_id =='%s']|[0]", alarmRuleId), respBody, nil)
	if alarmRule == nil {
		return nil, golangsdk.ErrDefault404{}
	}
	return alarmRule, nil
}

func resourceSQLAlarmRuleRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg    = meta.(*config.Config)
		region = cfg.GetRegion(d)
		mErr   *multierror.Error
	)
	getSQLAlarmRuleClient, err := cfg.NewServiceClient("lts", region)
	if err != nil {
		return diag.Errorf("error creating LTS client: %s", err)
	}

	getSQLAlarmRuleRespBody, err := GetSQLAlarmRuleById(getSQLAlarmRuleClient, d.Id())
	if err != nil {
		return common.CheckDeletedDiag(d, err, "error retrieving SQL alarm rule")
	}

	mErr = multierror.Append(
		mErr,
		d.Set("region", region),
		d.Set("name", utils.PathSearch("sql_alarm_rule_name", getSQLAlarmRuleRespBody, nil)),
		d.Set("sql_requests", flattenGetSQLAlarmRuleResponseBodySQLRequests(getSQLAlarmRuleRespBody)),
		d.Set("frequency", flattenGetSQLAlarmRuleResponseBodyFrequency(getSQLAlarmRuleRespBody)),
		d.Set("condition_expression", utils.PathSearch("condition_expression", getSQLAlarmRuleRespBody, nil)),
		d.Set("alarm_level", utils.PathSearch("sql_alarm_level", getSQLAlarmRuleRespBody, nil)),
		// OPtional attributes.
		d.Set("description", utils.PathSearch("sql_alarm_rule_description", getSQLAlarmRuleRespBody, nil)),
		d.Set("send_notifications", utils.PathSearch("sql_alarm_send", getSQLAlarmRuleRespBody, nil)),
		d.Set("alarm_action_rule_name", utils.PathSearch("alarm_action_rule_name", getSQLAlarmRuleRespBody, nil)),
		d.Set("notification_save_rule", flattenNotificationSaveRule(d, getSQLAlarmRuleRespBody)),
		d.Set("trigger_condition_count", utils.PathSearch("trigger_condition_count", getSQLAlarmRuleRespBody, nil)),
		d.Set("trigger_condition_frequency", utils.PathSearch("trigger_condition_frequency", getSQLAlarmRuleRespBody, nil)),
		d.Set("send_recovery_notifications", utils.PathSearch("whether_recovery_policy", getSQLAlarmRuleRespBody, nil)),
		d.Set("recovery_frequency", utils.PathSearch("recovery_policy", getSQLAlarmRuleRespBody, nil)),
		d.Set("alarm_rule_alias", utils.PathSearch("alarm_rule_alias", getSQLAlarmRuleRespBody, nil)),
		d.Set("notification_frequency", utils.PathSearch("notification_frequency", getSQLAlarmRuleRespBody, nil)),
		d.Set("status", utils.PathSearch("status", getSQLAlarmRuleRespBody, nil)),
		// Attributes.
		d.Set("domain_id", utils.PathSearch("domain_id", getSQLAlarmRuleRespBody, nil)),
		d.Set("created_at", utils.FormatTimeStampUTC(
			int64(utils.PathSearch("create_time", getSQLAlarmRuleRespBody, float64(0)).(float64)/1000))),
		d.Set("updated_at", utils.FormatTimeStampUTC(
			int64(utils.PathSearch("update_time", getSQLAlarmRuleRespBody, float64(0)).(float64)/1000))),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenGetSQLAlarmRuleResponseBodySQLRequests(resp interface{}) []interface{} {
	if resp == nil {
		return nil
	}
	curJson := utils.PathSearch("sql_requests", resp, make([]interface{}, 0))
	curArray := curJson.([]interface{})
	rst := make([]interface{}, 0, len(curArray))
	for _, v := range curArray {
		rst = append(rst, map[string]interface{}{
			"title":                  utils.PathSearch("sql_request_title", v, nil),
			"sql":                    utils.PathSearch("sql", v, nil),
			"log_stream_id":          utils.PathSearch("log_stream_id", v, nil),
			"log_group_id":           utils.PathSearch("log_group_id", v, nil),
			"search_time_range_unit": utils.PathSearch("search_time_range_unit", v, nil),
			"search_time_range":      utils.PathSearch("search_time_range", v, nil),
			"is_time_range_relative": utils.PathSearch("is_time_range_relative", v, nil),
			"log_group_name":         utils.PathSearch("log_group_name", v, nil),
			"log_stream_name":        utils.PathSearch("log_stream_name", v, nil),
		})
	}
	return rst
}

func flattenGetSQLAlarmRuleResponseBodyFrequency(resp interface{}) []interface{} {
	var rst []interface{}
	curJson := utils.PathSearch("frequency", resp, make(map[string]interface{})).(map[string]interface{})
	if len(curJson) < 1 {
		return rst
	}

	rst = []interface{}{
		map[string]interface{}{
			"type":            utils.PathSearch("type", curJson, nil),
			"cron_expression": utils.PathSearch("cron_expr", curJson, nil),
			"hour_of_day":     utils.PathSearch("hour_of_day", curJson, nil),
			"day_of_week":     utils.PathSearch("day_of_week", curJson, nil),
			"fixed_rate_unit": utils.PathSearch("fixed_rate_unit", curJson, nil),
			"fixed_rate":      utils.PathSearch("fixed_rate", curJson, nil),
		},
	}
	return rst
}

func resourceSQLAlarmRuleUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	client, err := cfg.NewServiceClient("lts", cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating LTS client: %s", err)
	}

	updateSQLAlarmRuleChanges := []string{
		"description",
		"sql_requests",
		"frequency",
		"condition_expression",
		"alarm_level",
		"send_notifications",
		"alarm_action_rule_name",
		"notification_save_rule",
		"trigger_condition_count",
		"trigger_condition_frequency",
		"send_recovery_notifications",
		"recovery_frequency",
		"alarm_rule_alias",
		"notification_frequency",
	}

	if d.HasChanges(updateSQLAlarmRuleChanges...) {
		params := buildUpdateSQLAlarmRuleBodyParams(d, cfg)
		params["sql_alarm_send_code"] = getAlarmSendCode(d)
		err = updateSQLAlarmRule(client, params)
		if err != nil {
			return diag.Errorf("error updating SQL alarm rule: %s", err)
		}
	}

	if d.HasChanges("status") {
		err = updateAlarmRuleStatus(client, d.Id(), "sql", d.Get("status").(string))
		if err != nil {
			return diag.Errorf("error updating status of the SQL alarm rule: %s", err)
		}
	}
	return resourceSQLAlarmRuleRead(ctx, d, meta)
}

func updateSQLAlarmRule(client *golangsdk.ServiceClient, params map[string]interface{}) error {
	updateHttpUrl := "v2/{project_id}/lts/alarms/sql-alarm-rule"
	updatePath := client.Endpoint + updateHttpUrl
	updatePath = strings.ReplaceAll(updatePath, "{project_id}", client.ProjectID)

	updateOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"Content-Type": "application/json"},
		JSONBody:         utils.RemoveNil(params),
	}
	_, err := client.Request("PUT", updatePath, &updateOpt)
	return err
}

func buildUpdateSQLAlarmRuleBodyParams(d *schema.ResourceData, cfg *config.Config) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"sql_alarm_rule_id":    utils.ValueIgnoreEmpty(d.Id()),
		"sql_alarm_rule_name":  d.Get("name"),
		"sql_requests":         buildSQLAlarmRuleRequestBodySQLRequests(d.Get("sql_requests")),
		"frequency":            buildSQLAlarmRuleRequestBodyFrequency(d.Get("frequency")),
		"condition_expression": d.Get("condition_expression"),
		"sql_alarm_level":      d.Get("alarm_level"),
		// Optional parameters.
		"sql_alarm_rule_description":  d.Get("description"),
		"alarm_action_rule_name":      utils.ValueIgnoreEmpty(d.Get("alarm_action_rule_name")),
		"notification_save_rule":      buildSQLAlarmRuleRequestBodyNotificationRule(d.Get("notification_save_rule")),
		"sql_alarm_send":              utils.ValueIgnoreEmpty(d.Get("send_notifications")),
		"sql_alarm_send_code":         0,
		"domain_id":                   cfg.DomainID,
		"trigger_condition_count":     utils.ValueIgnoreEmpty(d.Get("trigger_condition_count")),
		"trigger_condition_frequency": utils.ValueIgnoreEmpty(d.Get("trigger_condition_frequency")),
		"whether_recovery_policy":     utils.ValueIgnoreEmpty(d.Get("send_recovery_notifications")),
		"recovery_policy":             utils.ValueIgnoreEmpty(d.Get("recovery_frequency")),
		"alarm_rule_alias":            utils.ValueIgnoreEmpty(d.Get("alarm_rule_alias")),
		"notification_frequency":      d.Get("notification_frequency"),
	}
	return bodyParams
}

func resourceSQLAlarmRuleDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	// deleteSQLAlarmRule: Delete an existing LTS SQLAlarmRule
	var (
		deleteSQLAlarmRuleHttpUrl = "v2/{project_id}/lts/alarms/sql-alarm-rule/{id}"
		deleteSQLAlarmRuleProduct = "lts"
	)
	deleteSQLAlarmRuleClient, err := cfg.NewServiceClient(deleteSQLAlarmRuleProduct, region)
	if err != nil {
		return diag.Errorf("error creating LTS client: %s", err)
	}

	deleteSQLAlarmRulePath := deleteSQLAlarmRuleClient.Endpoint + deleteSQLAlarmRuleHttpUrl
	deleteSQLAlarmRulePath = strings.ReplaceAll(deleteSQLAlarmRulePath, "{project_id}", deleteSQLAlarmRuleClient.ProjectID)
	deleteSQLAlarmRulePath = strings.ReplaceAll(deleteSQLAlarmRulePath, "{id}", d.Id())

	deleteSQLAlarmRuleOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"Content-Type": "application/json"},
	}

	_, err = deleteSQLAlarmRuleClient.Request("DELETE", deleteSQLAlarmRulePath, &deleteSQLAlarmRuleOpt)
	if err != nil {
		return diag.Errorf("error deleting SQL alarm rule: %s", err)
	}

	return nil
}
