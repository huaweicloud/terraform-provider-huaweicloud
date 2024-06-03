// ---------------------------------------------------------------
// *** AUTO GENERATED CODE ***
// @Product LTS
// ---------------------------------------------------------------

package lts

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/jmespath/go-jmespath"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API LTS PUT /v2/{project_id}/lts/alarms/status
// @API LTS GET /v2/{project_id}/lts/alarms/sql-alarm-rule
// @API LTS POST /v2/{project_id}/lts/alarms/sql-alarm-rule
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
				Computed:    true,
				ForceNew:    true,
				Description: `Specifies whether to send notifications.`,
			},
			"notification_rule": {
				Type:        schema.TypeList,
				MaxItems:    1,
				Elem:        sqlAlarmRuleNotificationRuleSchema(),
				Optional:    true,
				ForceNew:    true,
				Computed:    true,
				Description: `Specifies the notification rule.`,
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
				ForceNew:    true,
				Description: `Specifies the notification template name.`,
			},
			"language": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: `Specifies the notification language.`,
			},
			"user_name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: `Specifies the user name.`,
			},
			"topics": {
				Type:        schema.TypeList,
				Elem:        sqlAlarmRuleNotificationRuleTopicSchema(),
				Required:    true,
				ForceNew:    true,
				Description: `Specifies the SMN topics.`,
			},
			"timezone": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Computed:    true,
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
				ForceNew:    true,
				Description: `Specifies the topic name.`,
			},
			"topic_urn": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: `Specifies the topic URN.`,
			},
			"display_name": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Computed:    true,
				Description: `Specifies the display name.`,
			},
			"push_policy": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Computed:    true,
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
		OkCodes: []int{
			200,
		},
		MoreHeaders: map[string]string{"Content-Type": "application/json"},
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

	id, err := jmespath.Search("sql_alarm_rule_id", createSQLAlarmRuleRespBody)
	if err != nil {
		return diag.Errorf("error creating SQL alarm rule: ID is not found in API response")
	}
	d.SetId(id.(string))

	if d.Get("status").(string) == "STOPPING" {
		// updateSQLAlarmRuleStatus: Update the LTS SQLAlarmRule status.
		var (
			updateSQLAlarmRuleStatusHttpUrl = "v2/{project_id}/lts/alarms/status"
		)

		updateSQLAlarmRuleStatusPath := createSQLAlarmRuleClient.Endpoint + updateSQLAlarmRuleStatusHttpUrl
		updateSQLAlarmRuleStatusPath = strings.ReplaceAll(updateSQLAlarmRuleStatusPath, "{project_id}", createSQLAlarmRuleClient.ProjectID)

		updateSQLAlarmRuleStatusOpt := golangsdk.RequestOpts{
			KeepResponseBody: true,
			OkCodes: []int{
				200,
			},
			MoreHeaders: map[string]string{"Content-Type": "application/json"},
		}

		updateSQLAlarmRuleStatusOpt.JSONBody = utils.RemoveNil(buildUpdateSQLAlarmRuleStatusBodyParams(d))
		_, err = createSQLAlarmRuleClient.Request("PUT", updateSQLAlarmRuleStatusPath, &updateSQLAlarmRuleStatusOpt)
		if err != nil {
			return diag.Errorf("error updating SQL alarm rule: %s", err)
		}
	}

	return resourceSQLAlarmRuleRead(ctx, d, meta)
}

func buildCreateSQLAlarmRuleBodyParams(d *schema.ResourceData, cfg *config.Config) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"sql_alarm_rule_name":         d.Get("name"),
		"sql_alarm_rule_description":  utils.ValueIgnoreEmpty(d.Get("description")),
		"sql_requests":                buildSQLAlarmRuleRequestBodySQLRequests(d.Get("sql_requests")),
		"frequency":                   buildSQLAlarmRuleRequestBodyFrequency(d.Get("frequency")),
		"condition_expression":        d.Get("condition_expression"),
		"sql_alarm_level":             d.Get("alarm_level"),
		"sql_alarm_send":              utils.ValueIgnoreEmpty(d.Get("send_notifications")),
		"domain_id":                   cfg.DomainID,
		"notification_rule":           buildSQLAlarmRuleRequestBodyNotificationRule(d.Get("notification_rule")),
		"trigger_condition_count":     utils.ValueIgnoreEmpty(d.Get("trigger_condition_count")),
		"trigger_condition_frequency": utils.ValueIgnoreEmpty(d.Get("trigger_condition_frequency")),
		"whether_recovery_policy":     utils.ValueIgnoreEmpty(d.Get("send_recovery_notifications")),
		"recovery_policy":             utils.ValueIgnoreEmpty(d.Get("recovery_frequency")),
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
				"push_policy":  utils.ValueIgnoreEmpty(raw["push_policy"]),
			}
		}
		return rst
	}
	return nil
}

func resourceSQLAlarmRuleRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var mErr *multierror.Error

	// getSQLAlarmRule: Query the LTS SQLAlarmRule detail
	var (
		getSQLAlarmRuleHttpUrl = "v2/{project_id}/lts/alarms/sql-alarm-rule"
		getSQLAlarmRuleProduct = "lts"
	)
	getSQLAlarmRuleClient, err := cfg.NewServiceClient(getSQLAlarmRuleProduct, region)
	if err != nil {
		return diag.Errorf("error creating LTS client: %s", err)
	}

	getSQLAlarmRulePath := getSQLAlarmRuleClient.Endpoint + getSQLAlarmRuleHttpUrl
	getSQLAlarmRulePath = strings.ReplaceAll(getSQLAlarmRulePath, "{project_id}", getSQLAlarmRuleClient.ProjectID)

	getSQLAlarmRuleOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		OkCodes: []int{
			200,
		},
		MoreHeaders: map[string]string{"Content-Type": "application/json"},
	}

	getSQLAlarmRuleResp, err := getSQLAlarmRuleClient.Request("GET", getSQLAlarmRulePath, &getSQLAlarmRuleOpt)

	if err != nil {
		return common.CheckDeletedDiag(d, err, "error retrieving SQLAlarmRule")
	}

	getSQLAlarmRuleRespBody, err := utils.FlattenResponse(getSQLAlarmRuleResp)
	if err != nil {
		return diag.FromErr(err)
	}

	jsonPath := fmt.Sprintf("sql_alarm_rules[?sql_alarm_rule_id =='%s']|[0]", d.Id())
	getSQLAlarmRuleRespBody = utils.PathSearch(jsonPath, getSQLAlarmRuleRespBody, nil)
	if getSQLAlarmRuleRespBody == nil {
		return common.CheckDeletedDiag(d, golangsdk.ErrDefault404{}, "no data found")
	}

	mErr = multierror.Append(
		mErr,
		d.Set("region", region),
		d.Set("name", utils.PathSearch("sql_alarm_rule_name", getSQLAlarmRuleRespBody, nil)),
		d.Set("description", utils.PathSearch("sql_alarm_rule_description", getSQLAlarmRuleRespBody, nil)),
		d.Set("sql_requests", flattenGetSQLAlarmRuleResponseBodySQLRequests(getSQLAlarmRuleRespBody)),
		d.Set("frequency", flattenGetSQLAlarmRuleResponseBodyFrequency(getSQLAlarmRuleRespBody)),
		d.Set("condition_expression", utils.PathSearch("condition_expression", getSQLAlarmRuleRespBody, nil)),
		d.Set("alarm_level", utils.PathSearch("sql_alarm_level", getSQLAlarmRuleRespBody, nil)),
		d.Set("send_notifications", utils.PathSearch("sql_alarm_send", getSQLAlarmRuleRespBody, nil)),
		d.Set("domain_id", utils.PathSearch("domain_id", getSQLAlarmRuleRespBody, nil)),
		d.Set("trigger_condition_count", utils.PathSearch("trigger_condition_count", getSQLAlarmRuleRespBody, nil)),
		d.Set("trigger_condition_frequency", utils.PathSearch("trigger_condition_frequency", getSQLAlarmRuleRespBody, nil)),
		d.Set("send_recovery_notifications", utils.PathSearch("whether_recovery_policy", getSQLAlarmRuleRespBody, nil)),
		d.Set("recovery_frequency", utils.PathSearch("recovery_policy", getSQLAlarmRuleRespBody, nil)),
		d.Set("status", utils.PathSearch("status", getSQLAlarmRuleRespBody, nil)),
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
		})
	}
	return rst
}

func flattenGetSQLAlarmRuleResponseBodyFrequency(resp interface{}) []interface{} {
	var rst []interface{}
	curJson, err := jmespath.Search("frequency", resp)
	if err != nil {
		log.Printf("[ERROR] error parsing frequency from response= %#v", resp)
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
	region := cfg.GetRegion(d)

	updateSQLAlarmRuleChanges := []string{
		"description",
		"sql_requests",
		"frequency",
		"condition_expression",
		"alarm_level",
		"trigger_condition_count",
		"trigger_condition_frequency",
		"send_recovery_notifications",
		"recovery_frequency",
	}

	if d.HasChanges(updateSQLAlarmRuleChanges...) {
		// updateSQLAlarmRule: Update the LTS SQLAlarmRule.
		var (
			updateSQLAlarmRuleHttpUrl = "v2/{project_id}/lts/alarms/sql-alarm-rule"
			updateSQLAlarmRuleProduct = "lts"
		)
		updateSQLAlarmRuleClient, err := cfg.NewServiceClient(updateSQLAlarmRuleProduct, region)
		if err != nil {
			return diag.Errorf("error creating LTS client: %s", err)
		}

		updateSQLAlarmRulePath := updateSQLAlarmRuleClient.Endpoint + updateSQLAlarmRuleHttpUrl
		updateSQLAlarmRulePath = strings.ReplaceAll(updateSQLAlarmRulePath, "{project_id}", updateSQLAlarmRuleClient.ProjectID)

		updateSQLAlarmRuleOpt := golangsdk.RequestOpts{
			KeepResponseBody: true,
			OkCodes: []int{
				200,
			},
			MoreHeaders: map[string]string{"Content-Type": "application/json"},
		}

		updateSQLAlarmRuleOpt.JSONBody = utils.RemoveNil(buildUpdateSQLAlarmRuleBodyParams(d, cfg))
		_, err = updateSQLAlarmRuleClient.Request("PUT", updateSQLAlarmRulePath, &updateSQLAlarmRuleOpt)
		if err != nil {
			return diag.Errorf("error updating SQL alarm rule: %s", err)
		}
	}
	updateSQLAlarmRuleStatusChanges := []string{
		"status",
	}

	if d.HasChanges(updateSQLAlarmRuleStatusChanges...) {
		// updateSQLAlarmRuleStatus: Update the LTS SQLAlarmRule status.
		var (
			updateSQLAlarmRuleStatusHttpUrl = "v2/{project_id}/lts/alarms/status"
			updateSQLAlarmRuleStatusProduct = "lts"
		)
		updateSQLAlarmRuleStatusClient, err := cfg.NewServiceClient(updateSQLAlarmRuleStatusProduct, region)
		if err != nil {
			return diag.Errorf("error creating LTS client: %s", err)
		}

		updateSQLAlarmRuleStatusPath := updateSQLAlarmRuleStatusClient.Endpoint + updateSQLAlarmRuleStatusHttpUrl
		updateSQLAlarmRuleStatusPath = strings.ReplaceAll(updateSQLAlarmRuleStatusPath, "{project_id}", updateSQLAlarmRuleStatusClient.ProjectID)

		updateSQLAlarmRuleStatusOpt := golangsdk.RequestOpts{
			KeepResponseBody: true,
			OkCodes: []int{
				200,
			},
			MoreHeaders: map[string]string{"Content-Type": "application/json"},
		}

		updateSQLAlarmRuleStatusOpt.JSONBody = utils.RemoveNil(buildUpdateSQLAlarmRuleStatusBodyParams(d))
		_, err = updateSQLAlarmRuleStatusClient.Request("PUT", updateSQLAlarmRuleStatusPath, &updateSQLAlarmRuleStatusOpt)
		if err != nil {
			return diag.Errorf("error updating SQL alarm rule: %s", err)
		}
	}
	return resourceSQLAlarmRuleRead(ctx, d, meta)
}

func buildUpdateSQLAlarmRuleBodyParams(d *schema.ResourceData, cfg *config.Config) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"sql_alarm_rule_id":           utils.ValueIgnoreEmpty(d.Id()),
		"sql_alarm_rule_name":         d.Get("name"),
		"sql_alarm_rule_description":  d.Get("description"),
		"sql_requests":                buildSQLAlarmRuleRequestBodySQLRequests(d.Get("sql_requests")),
		"frequency":                   buildSQLAlarmRuleRequestBodyFrequency(d.Get("frequency")),
		"condition_expression":        d.Get("condition_expression"),
		"sql_alarm_level":             d.Get("alarm_level"),
		"sql_alarm_send":              utils.ValueIgnoreEmpty(d.Get("send_notifications")),
		"sql_alarm_send_code":         0,
		"domain_id":                   cfg.DomainID,
		"trigger_condition_count":     utils.ValueIgnoreEmpty(d.Get("trigger_condition_count")),
		"trigger_condition_frequency": utils.ValueIgnoreEmpty(d.Get("trigger_condition_frequency")),
		"whether_recovery_policy":     utils.ValueIgnoreEmpty(d.Get("send_recovery_notifications")),
		"recovery_policy":             utils.ValueIgnoreEmpty(d.Get("recovery_frequency")),
	}
	return bodyParams
}

func buildUpdateSQLAlarmRuleStatusBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"alarm_rule_id": d.Id(),
		"type":          "sql",
		"status":        utils.ValueIgnoreEmpty(d.Get("status")),
	}
	log.Printf("xxxx bodyParams: %v", bodyParams)
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
		OkCodes: []int{
			200,
		},
		MoreHeaders: map[string]string{"Content-Type": "application/json"},
	}

	_, err = deleteSQLAlarmRuleClient.Request("DELETE", deleteSQLAlarmRulePath, &deleteSQLAlarmRuleOpt)
	if err != nil {
		return diag.Errorf("error deleting SQL alarm rule: %s", err)
	}

	return nil
}
