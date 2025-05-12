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

// @API LTS GET /v2/{project_id}/lts/alarms/keywords-alarm-rule
// @API LTS POST /v2/{project_id}/lts/alarms/keywords-alarm-rule
// @API LTS PUT /v2/{project_id}/lts/alarms/keywords-alarm-rule
// @API LTS DELETE /v2/{project_id}/lts/alarms/keywords-alarm-rule/{keywords_alarm_rule_id}
// @API LTS PUT /v2/{project_id}/lts/alarms/status
func ResourceKeywordsAlarmRule() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceKeywordsAlarmRuleCreate,
		UpdateContext: resourceKeywordsAlarmRuleUpdate,
		ReadContext:   resourceKeywordsAlarmRuleRead,
		DeleteContext: resourceKeywordsAlarmRuleDelete,
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
				Description: `Specifies the name of the keywords alarm rule.`,
			},
			"keywords_requests": {
				Type:        schema.TypeList,
				Elem:        keywordsAlarmRuleKeywordsRequestsSchema(),
				Required:    true,
				Description: `Specifies the keywords requests.`,
			},
			"frequency": {
				Type:        schema.TypeList,
				MaxItems:    1,
				Elem:        keywordsAlarmRuleFrequencySchema(),
				Required:    true,
				Description: `Specifies the alarm frequency configurations.`,
			},
			"alarm_level": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `Specifies the alarm level.`,
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specifies the description of the keywords alarm rule.`,
			},
			"send_notifications": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: `Specifies whether to send notifications.`,
			},
			"alarm_action_rule_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The name of the alarm action rule associated with the keyword alarm rule.`,
			},
			"notification_save_rule": {
				Type:        schema.TypeList,
				MaxItems:    1,
				Elem:        keywordsAlarmRuleNotificationRuleSchema(),
				Optional:    true,
				Description: `The notification rule of the keyword alarm rule.`,
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
				Description: `Specifies alias name of the keyword alarm rule.`,
			},
			"notification_frequency": {
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
				Description: `The notification frequency of the keyword alarm rule, in minutes.`,
			},
			"status": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `Specifies the status.`,
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
			"condition_expression": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The condition expression of the keyword alarm rule.`,
			},
			// Deprecated parameters.
			"notification_rule": {
				Type:     schema.TypeList,
				MaxItems: 1,
				Elem:     keywordsAlarmRuleNotificationRuleSchema(),
				Optional: true,
				ForceNew: true,
				Computed: true,
				Description: utils.SchemaDesc(
					`The notification rule of the keyword alarm rule.`,
					utils.SchemaDescInput{
						Deprecated: true,
					},
				),
			},
		},
	}
}

func keywordsAlarmRuleKeywordsRequestsSchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"keywords": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `Specifies the keywords.`,
			},
			"condition": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `Specifies the keywords request condition.`,
			},
			"number": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: `Specifies the line number.`,
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

func keywordsAlarmRuleFrequencySchema() *schema.Resource {
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

func keywordsAlarmRuleNotificationRuleSchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"template_name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `Specifies the notification template name.`,
			},
			"user_name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `Specifies the user name.`,
			},
			"topics": {
				Type:        schema.TypeList,
				Elem:        keywordsAlarmRuleNotificationRuleTopicSchema(),
				Required:    true,
				Description: `Specifies the SMN topics.`,
			},
			"timezone": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specifies the timezone.`,
			},
			"language": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `Specifies the notification language.`,
			},
		},
	}
	return &sc
}

func keywordsAlarmRuleNotificationRuleTopicSchema() *schema.Resource {
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

func resourceKeywordsAlarmRuleCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	// createKeywordsAlarmRule: Create an LTS KeywordsAlarmRule.
	var (
		createKeywordsAlarmRuleHttpUrl = "v2/{project_id}/lts/alarms/keywords-alarm-rule"
		createKeywordsAlarmRuleProduct = "lts"
	)
	createKeywordsAlarmRuleClient, err := cfg.NewServiceClient(createKeywordsAlarmRuleProduct, region)
	if err != nil {
		return diag.Errorf("error creating LTS client: %s", err)
	}

	createKeywordsAlarmRulePath := createKeywordsAlarmRuleClient.Endpoint + createKeywordsAlarmRuleHttpUrl
	createKeywordsAlarmRulePath = strings.ReplaceAll(createKeywordsAlarmRulePath, "{project_id}", createKeywordsAlarmRuleClient.ProjectID)

	createKeywordsAlarmRuleOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"Content-Type": "application/json"},
	}

	createKeywordsAlarmRuleOpt.JSONBody = utils.RemoveNil(buildCreateKeywordsAlarmRuleBodyParams(d, cfg))
	createKeywordsAlarmRuleResp, err := createKeywordsAlarmRuleClient.Request("POST", createKeywordsAlarmRulePath, &createKeywordsAlarmRuleOpt)
	if err != nil {
		return diag.Errorf("error creating Keywords alarm rule: %s", err)
	}

	createKeywordsAlarmRuleRespBody, err := utils.FlattenResponse(createKeywordsAlarmRuleResp)
	if err != nil {
		return diag.FromErr(err)
	}

	ruleId := utils.PathSearch("keywords_alarm_rule_id", createKeywordsAlarmRuleRespBody, "").(string)
	if ruleId == "" {
		return diag.Errorf("unable to find the LTS keywords alarm rule ID from the API response")
	}
	d.SetId(ruleId)

	if _, ok := d.GetOk("alarm_rule_alias"); ok {
		err = updateKeywordAlarmRule(createKeywordsAlarmRuleClient, buildUpdateKeywordsAlarmRuleBodyParams(d, cfg))
		if err != nil {
			return diag.Errorf("error updating alias name of the keyword alarm rule (%s): %s", err, ruleId)
		}
	}

	if d.Get("status").(string) == "STOPPING" {
		err = updateAlarmRuleStatus(createKeywordsAlarmRuleClient, ruleId, "keywords", "STOPPING")
		if err != nil {
			return diag.Errorf("error stopping keyword alarm rule: %s", err)
		}
	}

	return resourceKeywordsAlarmRuleRead(ctx, d, meta)
}

func buildCreateKeywordsAlarmRuleBodyParams(d *schema.ResourceData, cfg *config.Config) map[string]interface{} {
	bodyParams := map[string]interface{}{
		// Required parameter.
		"keywords_alarm_rule_name": d.Get("name"),
		"keywords_requests":        buildKeywordsAlarmRuleRequestBodyKeywordsRequests(d.Get("keywords_requests")),
		"frequency":                buildKeywordsAlarmRuleRequestBodyFrequency(d.Get("frequency")),
		"keywords_alarm_level":     d.Get("alarm_level"),
		"domain_id":                cfg.DomainID,
		// Optional parameter.
		"keywords_alarm_send":             utils.ValueIgnoreEmpty(d.Get("send_notifications")),
		"keywords_alarm_rule_description": utils.ValueIgnoreEmpty(d.Get("description")),
		"alarm_action_rule_name":          utils.ValueIgnoreEmpty(d.Get("alarm_action_rule_name")),
		"notification_rule":               buildKeywordsAlarmRuleRequestBodyNotificationRule(d.Get("notification_rule")),
		"notification_save_rule":          buildKeywordsAlarmRuleRequestBodyNotificationRule(d.Get("notification_save_rule")),
		"trigger_condition_count":         utils.ValueIgnoreEmpty(d.Get("trigger_condition_count")),
		"trigger_condition_frequency":     utils.ValueIgnoreEmpty(d.Get("trigger_condition_frequency")),
		"whether_recovery_policy":         utils.ValueIgnoreEmpty(d.Get("send_recovery_notifications")),
		"recovery_policy":                 utils.ValueIgnoreEmpty(d.Get("recovery_frequency")),
		"notification_frequency":          d.Get("notification_frequency"),
	}
	return bodyParams
}

func buildKeywordsAlarmRuleRequestBodyKeywordsRequests(rawParams interface{}) []map[string]interface{} {
	if rawArray, ok := rawParams.([]interface{}); ok {
		if len(rawArray) == 0 {
			return nil
		}

		rst := make([]map[string]interface{}, len(rawArray))
		for i, v := range rawArray {
			raw := v.(map[string]interface{})
			rst[i] = map[string]interface{}{
				"keywords":               utils.ValueIgnoreEmpty(raw["keywords"]),
				"condition":              utils.ValueIgnoreEmpty(raw["condition"]),
				"number":                 utils.ValueIgnoreEmpty(raw["number"]),
				"log_stream_id":          utils.ValueIgnoreEmpty(raw["log_stream_id"]),
				"log_group_id":           utils.ValueIgnoreEmpty(raw["log_group_id"]),
				"search_time_range_unit": utils.ValueIgnoreEmpty(raw["search_time_range_unit"]),
				"search_time_range":      utils.ValueIgnoreEmpty(raw["search_time_range"]),
				"log_group_name":         utils.ValueIgnoreEmpty(raw["log_group_name"]),
				"log_stream_name":        utils.ValueIgnoreEmpty(raw["log_stream_name"]),
			}
		}
		return rst
	}
	return nil
}

func buildKeywordsAlarmRuleRequestBodyFrequency(rawParams interface{}) map[string]interface{} {
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

func buildKeywordsAlarmRuleRequestBodyNotificationRule(rawParams interface{}) map[string]interface{} {
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
			"user_name":     utils.ValueIgnoreEmpty(raw["user_name"]),
			"topics":        buildKeywordsNotificationRuleTopic(raw["topics"]),
			"timezone":      utils.ValueIgnoreEmpty(raw["timezone"]),
			"language":      utils.ValueIgnoreEmpty(raw["language"]),
		}
		return params
	}
	return nil
}

func buildKeywordsNotificationRuleTopic(rawParams interface{}) []map[string]interface{} {
	if rawArray, ok := rawParams.([]interface{}); ok {
		if len(rawArray) == 0 {
			return nil
		}

		rst := make([]map[string]interface{}, len(rawArray))
		for i, v := range rawArray {
			raw := v.(map[string]interface{})
			rst[i] = map[string]interface{}{
				"name":         raw["name"],
				"topic_urn":    raw["topic_urn"],
				"display_name": utils.ValueIgnoreEmpty(raw["display_name"]),
				// Valid values: `0` and `1`.
				"push_policy": raw["push_policy"],
			}
		}
		return rst
	}
	return nil
}

func GetKeywordsAlarmRuleById(client *golangsdk.ServiceClient, alarmRuleId string) (interface{}, error) {
	httpUrl := "v2/{project_id}/lts/alarms/keywords-alarm-rule"
	getPath := client.Endpoint + httpUrl
	getPath = strings.ReplaceAll(getPath, "{project_id}", client.ProjectID)
	getOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"Content-Type": "application/json"},
	}

	resp, err := client.Request("GET", getPath, &getOpt)
	if err != nil {
		return nil, err
	}

	respBody, err := utils.FlattenResponse(resp)
	if err != nil {
		return nil, err
	}

	jsonPath := fmt.Sprintf("keywords_alarm_rules[?keywords_alarm_rule_id =='%s']|[0]", alarmRuleId)
	keywordsAlarmRule := utils.PathSearch(jsonPath, respBody, nil)
	if keywordsAlarmRule == nil {
		return nil, golangsdk.ErrDefault404{}
	}

	return keywordsAlarmRule, nil
}

func resourceKeywordsAlarmRuleRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var mErr *multierror.Error

	client, err := cfg.NewServiceClient("lts", region)
	if err != nil {
		return diag.Errorf("error creating LTS client: %s", err)
	}

	getKeywordsAlarmRuleRespBody, err := GetKeywordsAlarmRuleById(client, d.Id())
	if err != nil {
		return common.CheckDeletedDiag(d, err, "error retrieving keyword alarm rule")
	}

	mErr = multierror.Append(
		mErr,
		d.Set("region", region),
		d.Set("name", utils.PathSearch("keywords_alarm_rule_name", getKeywordsAlarmRuleRespBody, nil)),
		d.Set("description", utils.PathSearch("keywords_alarm_rule_description", getKeywordsAlarmRuleRespBody, nil)),
		d.Set("keywords_requests", flattenGetKeywordsAlarmRuleResponseBodyKeywordsRequests(getKeywordsAlarmRuleRespBody)),
		d.Set("frequency", flattenGetKeywordsAlarmRuleResponseBodyFrequency(getKeywordsAlarmRuleRespBody)),
		d.Set("alarm_level", utils.PathSearch("keywords_alarm_level", getKeywordsAlarmRuleRespBody, nil)),
		d.Set("send_notifications", utils.PathSearch("keywords_alarm_send", getKeywordsAlarmRuleRespBody, nil)),
		d.Set("alarm_action_rule_name", utils.PathSearch("alarm_action_rule_name", getKeywordsAlarmRuleRespBody, nil)),
		d.Set("notification_save_rule", flattenNotificationSaveRule(d, getKeywordsAlarmRuleRespBody)),
		d.Set("domain_id", utils.PathSearch("domain_id", getKeywordsAlarmRuleRespBody, nil)),
		d.Set("trigger_condition_count", utils.PathSearch("trigger_condition_count", getKeywordsAlarmRuleRespBody, nil)),
		d.Set("trigger_condition_frequency", utils.PathSearch("trigger_condition_frequency", getKeywordsAlarmRuleRespBody, nil)),
		d.Set("send_recovery_notifications", utils.PathSearch("whether_recovery_policy", getKeywordsAlarmRuleRespBody, nil)),
		d.Set("recovery_frequency", utils.PathSearch("recovery_policy", getKeywordsAlarmRuleRespBody, nil)),
		d.Set("alarm_rule_alias", utils.PathSearch("alarm_rule_alias", getKeywordsAlarmRuleRespBody, nil)),
		d.Set("notification_frequency", utils.PathSearch("notification_frequency", getKeywordsAlarmRuleRespBody, nil)),
		d.Set("status", utils.PathSearch("status", getKeywordsAlarmRuleRespBody, nil)),
		d.Set("created_at", utils.FormatTimeStampUTC(
			int64(utils.PathSearch("create_time", getKeywordsAlarmRuleRespBody, float64(0)).(float64)/1000))),
		d.Set("updated_at", utils.FormatTimeStampUTC(
			int64(utils.PathSearch("update_time", getKeywordsAlarmRuleRespBody, float64(0)).(float64)/1000))),
		d.Set("condition_expression", utils.PathSearch("condition_expression", getKeywordsAlarmRuleRespBody, nil)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenGetKeywordsAlarmRuleResponseBodyKeywordsRequests(resp interface{}) []interface{} {
	if resp == nil {
		return nil
	}
	curJson := utils.PathSearch("keywords_requests", resp, make([]interface{}, 0))
	curArray := curJson.([]interface{})
	rst := make([]interface{}, 0, len(curArray))
	for _, v := range curArray {
		rst = append(rst, map[string]interface{}{
			"keywords":               utils.PathSearch("keywords", v, nil),
			"condition":              utils.PathSearch("condition", v, nil),
			"number":                 utils.PathSearch("number", v, nil),
			"log_stream_id":          utils.PathSearch("log_stream_id", v, nil),
			"log_group_id":           utils.PathSearch("log_group_id", v, nil),
			"search_time_range_unit": utils.PathSearch("search_time_range_unit", v, nil),
			"search_time_range":      utils.PathSearch("search_time_range", v, nil),
			"log_group_name":         utils.PathSearch("log_group_name", v, nil),
			"log_stream_name":        utils.PathSearch("log_stream_name", v, nil),
		})
	}
	return rst
}

func flattenGetKeywordsAlarmRuleResponseBodyFrequency(resp interface{}) []interface{} {
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

func flattenNotificationSaveRule(d *schema.ResourceData, resp interface{}) []interface{} {
	topics := utils.PathSearch("topics", resp, make([]interface{}, 0)).([]interface{})
	if len(topics) == 0 {
		return nil
	}

	return []interface{}{
		map[string]interface{}{
			"template_name": utils.PathSearch("template_name", resp, nil),
			"user_name":     d.Get("notification_save_rule.0.user_name"),
			"topics":        flattenNotificationSaveRuleTopics(topics),
			"timezone":      d.Get("notification_save_rule.0.timezone"),
			"language":      utils.PathSearch("language", resp, nil),
		},
	}
}

func flattenNotificationSaveRuleTopics(topics []interface{}) []interface{} {
	rest := make([]interface{}, len(topics))
	for i, v := range topics {
		rest[i] = map[string]interface{}{
			"name":         utils.PathSearch("name", v, nil),
			"topic_urn":    utils.PathSearch("topic_urn", v, nil),
			"display_name": utils.PathSearch("display_name", v, nil),
			"push_policy":  utils.PathSearch("push_policy", v, nil),
		}
	}

	return rest
}

func resourceKeywordsAlarmRuleUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg                = meta.(*config.Config)
		keywordAlarmRuleId = d.Id()
	)

	client, err := cfg.NewServiceClient("lts", cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating LTS client: %s", err)
	}

	updateKeywordsAlarmRuleChanges := []string{
		"description",
		"keywords_requests",
		"frequency",
		"alarm_level",
		"trigger_condition_count",
		"trigger_condition_frequency",
		"send_recovery_notifications",
		"recovery_frequency",
		"send_notifications",
		"alarm_action_rule_name",
		"notification_save_rule",
		"alarm_rule_alias",
		"notification_frequency",
	}

	if d.HasChanges(updateKeywordsAlarmRuleChanges...) {
		params := buildUpdateKeywordsAlarmRuleBodyParams(d, cfg)
		params["keywords_alarm_send_code"] = getAlarmSendCode(d)
		err = updateKeywordAlarmRule(client, params)
		if err != nil {
			return diag.Errorf("error updating keyword alarm rule (%s): %s", err, keywordAlarmRuleId)
		}
	}

	if d.HasChange("status") {
		err = updateAlarmRuleStatus(client, keywordAlarmRuleId, "keywords", d.Get("status").(string))
		if err != nil {
			return diag.Errorf("error updating status of the keyword alarm rule (%s): %s", err, keywordAlarmRuleId)
		}
	}
	return resourceKeywordsAlarmRuleRead(ctx, d, meta)
}

// `0` means do not modify the topics.
// `1` means add the topics.
// `2` means modify the topics.
// `3` means delete the topics.
func getAlarmSendCode(d *schema.ResourceData) int {
	if !d.HasChange("notification_save_rule.0.topics") {
		return 0
	}

	oRaw, nRaw := d.GetChange("notification_save_rule.0.topics")
	oTopics := oRaw.([]interface{})
	nTopics := nRaw.([]interface{})
	if len(oTopics) == 0 && len(nTopics) != 0 {
		return 1
	}

	if len(nTopics) == 0 && !d.Get("send_notifications").(bool) {
		return 3
	}
	return 2
}

func buildUpdateKeywordsAlarmRuleBodyParams(d *schema.ResourceData, cfg *config.Config) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"keywords_alarm_rule_id":          d.Id(),
		"keywords_alarm_rule_name":        d.Get("name"),
		"keywords_alarm_rule_description": d.Get("description"),
		"keywords_requests":               buildKeywordsAlarmRuleRequestBodyKeywordsRequests(d.Get("keywords_requests")),
		"frequency":                       buildKeywordsAlarmRuleRequestBodyFrequency(d.Get("frequency")),
		"keywords_alarm_level":            d.Get("alarm_level"),
		"keywords_alarm_send":             utils.ValueIgnoreEmpty(d.Get("send_notifications")),
		// `0` means do not modify the topics.
		"keywords_alarm_send_code":    0,
		"domain_id":                   cfg.DomainID,
		"trigger_condition_count":     utils.ValueIgnoreEmpty(d.Get("trigger_condition_count")),
		"trigger_condition_frequency": utils.ValueIgnoreEmpty(d.Get("trigger_condition_frequency")),
		"whether_recovery_policy":     utils.ValueIgnoreEmpty(d.Get("send_recovery_notifications")),
		"recovery_policy":             utils.ValueIgnoreEmpty(d.Get("recovery_frequency")),
		"alarm_action_rule_name":      utils.ValueIgnoreEmpty(d.Get("alarm_action_rule_name")),
		"alarm_rule_alias":            utils.ValueIgnoreEmpty(d.Get("alarm_rule_alias")),
		"notification_save_rule":      buildKeywordsAlarmRuleRequestBodyNotificationRule(d.Get("notification_save_rule")),
		"notification_frequency":      d.Get("notification_frequency"),
	}

	return bodyParams
}

func updateKeywordAlarmRule(client *golangsdk.ServiceClient, params map[string]interface{}) error {
	updateHttpUrl := "v2/{project_id}/lts/alarms/keywords-alarm-rule"
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

func updateAlarmRuleStatus(client *golangsdk.ServiceClient, alarmRuleId, alarmRuleType, status string) error {
	updateAlarmRuleStatusHttpUrl := "v2/{project_id}/lts/alarms/status"
	updateAlarmRuleStatusPath := client.Endpoint + updateAlarmRuleStatusHttpUrl
	updateAlarmRuleStatusPath = strings.ReplaceAll(updateAlarmRuleStatusPath, "{project_id}", client.ProjectID)
	updateAlarmRuleStatusOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"Content-Type": "application/json"},
		JSONBody:         utils.RemoveNil(buildUpdateAlarmRuleStatusBodyParams(alarmRuleId, alarmRuleType, status)),
	}
	_, err := client.Request("PUT", updateAlarmRuleStatusPath, &updateAlarmRuleStatusOpt)
	if err != nil {
		return err
	}

	return nil
}

func buildUpdateAlarmRuleStatusBodyParams(alarmRuleId, alarmRuletype, status string) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"alarm_rule_id": alarmRuleId,
		"type":          alarmRuletype,
		"status":        status,
	}
	return bodyParams
}

func resourceKeywordsAlarmRuleDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	// deleteKeywordsAlarmRule: Delete an existing LTS KeywordsAlarmRule
	var (
		deleteKeywordsAlarmRuleHttpUrl = "v2/{project_id}/lts/alarms/keywords-alarm-rule/{id}"
		deleteKeywordsAlarmRuleProduct = "lts"
	)
	deleteKeywordsAlarmRuleClient, err := cfg.NewServiceClient(deleteKeywordsAlarmRuleProduct, region)
	if err != nil {
		return diag.Errorf("error creating LTS client: %s", err)
	}

	deleteKeywordsAlarmRulePath := deleteKeywordsAlarmRuleClient.Endpoint + deleteKeywordsAlarmRuleHttpUrl
	deleteKeywordsAlarmRulePath = strings.ReplaceAll(deleteKeywordsAlarmRulePath, "{project_id}", deleteKeywordsAlarmRuleClient.ProjectID)
	deleteKeywordsAlarmRulePath = strings.ReplaceAll(deleteKeywordsAlarmRulePath, "{id}", d.Id())

	deleteKeywordsAlarmRuleOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		OkCodes: []int{
			200,
		},
		MoreHeaders: map[string]string{"Content-Type": "application/json"},
	}

	_, err = deleteKeywordsAlarmRuleClient.Request("DELETE", deleteKeywordsAlarmRulePath, &deleteKeywordsAlarmRuleOpt)
	if err != nil {
		return diag.Errorf("error deleting Keywords alarm rule: %s", err)
	}

	return nil
}
