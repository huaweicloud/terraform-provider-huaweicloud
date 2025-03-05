package iotda

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API IoTDA GET /v5/iot/{project_id}/rules
func DataSourceDeviceLinkageRules() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceDeviceLinkageRulesRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"rule_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"type": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"space_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"status": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"rules": {
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
						"description": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"status": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"space_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"triggers": {
							Type:     schema.TypeList,
							Elem:     ruleConditionSchema(),
							Computed: true,
						},
						"actions": {
							Type:     schema.TypeList,
							Elem:     ruleActionSchema(),
							Computed: true,
						},
						"effective_period": {
							Type:     schema.TypeList,
							Elem:     ruleTimeRangeSchema(),
							Computed: true,
						},
						"trigger_logic": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"updated_at": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func ruleConditionSchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"device_data_condition": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"device_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"product_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"path": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"operator": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"value": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"in_values": {
							Type:     schema.TypeList,
							Computed: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
						"trigger_strategy": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"data_validatiy_period": {
							Type:     schema.TypeInt,
							Computed: true,
						},
					},
				},
			},
			"simple_timer_condition": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"start_time": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"repeat_interval": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"repeat_count": {
							Type:     schema.TypeInt,
							Computed: true,
						},
					},
				},
			},
			"daily_timer_condition": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"start_time": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"days_of_week": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			"device_linkage_status_condition": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"device_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"product_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"status_list": {
							Type:     schema.TypeList,
							Elem:     &schema.Schema{Type: schema.TypeString},
							Computed: true,
						},
						"duration": {
							Type:     schema.TypeInt,
							Computed: true,
						},
					},
				},
			},
		},
	}
	return &sc
}

func ruleActionSchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"device_command": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"device_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"command_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"command_body": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"buffer_timeout": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"response_timeout": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"mode": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"service_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			"smn_forwarding": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"region": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"project_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"topic_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"topic_urn": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"message_title": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"message_content": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"message_template_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			"device_alarm": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"severity": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"dimension": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"description": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
	return &sc
}

func ruleTimeRangeSchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"start_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"end_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"days_of_week": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
	return &sc
}

func buildLinkageRulesQueryParams(d *schema.ResourceData) string {
	req := ""

	if v, ok := d.GetOk("space_id"); ok {
		req = fmt.Sprintf("%s&app_id=%v", req, v)
	}

	if v, ok := d.GetOk("type"); ok {
		req = fmt.Sprintf("%s&rule_type=%v", req, v)
	}

	return req
}

func dataSourceDeviceLinkageRulesRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg       = meta.(*config.Config)
		region    = cfg.GetRegion(d)
		isDerived = WithDerivedAuth(cfg, region)
		httpUrl   = "v5/iot/{project_id}/rules?limit=50"
		allRules  = make([]interface{}, 0)
		offset    = 0
	)

	client, err := cfg.NewServiceClientWithDerivedAuth("iotda", region, isDerived)
	if err != nil {
		return diag.Errorf("error creating IoTDA client: %s", err)
	}

	listPath := client.Endpoint + httpUrl
	listPath = strings.ReplaceAll(listPath, "{project_id}", client.ProjectID)
	listPath += buildLinkageRulesQueryParams(d)
	listOpts := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	for {
		listPathWithOffset := fmt.Sprintf("%s&offset=%d", listPath, offset)
		resp, err := client.Request("GET", listPathWithOffset, &listOpts)
		if err != nil {
			return diag.Errorf("error querying IoTDA device linkage rules: %s", err)
		}

		respBody, err := utils.FlattenResponse(resp)
		if err != nil {
			return diag.FromErr(err)
		}

		rules := utils.PathSearch("rules", respBody, make([]interface{}, 0)).([]interface{})
		if len(rules) == 0 {
			break
		}

		allRules = append(allRules, rules...)
		offset += len(rules)
	}

	uuId, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}

	d.SetId(uuId)

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("rules", flattenDeviceLinkageRules(filterListDeviceLinkageRules(allRules, d))),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func filterListDeviceLinkageRules(rules []interface{}, d *schema.ResourceData) []interface{} {
	if len(rules) == 0 {
		return nil
	}

	rst := make([]interface{}, 0, len(rules))
	for _, v := range rules {
		if ruleId, ok := d.GetOk("rule_id"); ok &&
			fmt.Sprint(ruleId) != utils.PathSearch("rule_id", v, "").(string) {
			continue
		}

		if ruleName, ok := d.GetOk("name"); ok &&
			fmt.Sprint(ruleName) != utils.PathSearch("name", v, "").(string) {
			continue
		}

		if status, ok := d.GetOk("status"); ok &&
			fmt.Sprint(status) != utils.PathSearch("status", v, "").(string) {
			continue
		}

		rst = append(rst, v)
	}

	return rst
}

func flattenDeviceLinkageRules(rules []interface{}) []interface{} {
	if len(rules) == 0 {
		return nil
	}

	rst := make([]interface{}, 0, len(rules))
	for _, v := range rules {
		rst = append(rst, map[string]interface{}{
			"id":               utils.PathSearch("rule_id", v, nil),
			"name":             utils.PathSearch("name", v, nil),
			"description":      utils.PathSearch("description", v, nil),
			"type":             utils.PathSearch("rule_type", v, nil),
			"status":           utils.PathSearch("status", v, nil),
			"space_id":         utils.PathSearch("app_id", v, nil),
			"triggers":         flattenLinkageRuleTriggers(utils.PathSearch("condition_group", v, nil)),
			"actions":          flattenLinkageRuleActions(utils.PathSearch("actions", v, make([]interface{}, 0)).([]interface{})),
			"effective_period": flattenLinkageRuleEffectivePeriod(utils.PathSearch("condition_group", v, nil)),
			"trigger_logic":    utils.PathSearch("condition_group.logic", v, nil),
			"updated_at":       utils.PathSearch("last_update_time", v, nil),
		})
	}

	return rst
}

func flattenLinkageRuleTriggers(conditionGroup interface{}) []interface{} {
	if conditionGroup == nil {
		return nil
	}

	conditions := utils.PathSearch("conditions", conditionGroup, make([]interface{}, 0)).([]interface{})
	if len(conditions) == 0 {
		return nil
	}

	rst := make([]interface{}, 0, len(conditions))
	for _, v := range conditions {
		conditionType := utils.PathSearch("type", v, "").(string)
		switch conditionType {
		case "DEVICE_DATA":
			deviceProperty := utils.PathSearch("device_property_condition", v, nil)
			filters := utils.PathSearch("device_property_condition.filters", v, make([]interface{}, 0)).([]interface{})
			if deviceProperty != nil && len(filters) != 0 {
				rst = append(rst, flattenDeviceDataCondition(conditionType, deviceProperty))
			}
		case "SIMPLE_TIMER":
			simpleTimer := utils.PathSearch("simple_timer_condition", v, nil)
			if simpleTimer != nil {
				rst = append(rst, flattenSimpleTimerCondition(conditionType, simpleTimer))
			}
		case "DAILY_TIMER":
			dailyTimer := utils.PathSearch("daily_timer_condition", v, nil)
			if dailyTimer != nil {
				rst = append(rst, flattenDailyTimerCondition(conditionType, dailyTimer))
			}
		case "DEVICE_LINKAGE_STATUS":
			deviceLinkageStatus := utils.PathSearch("device_linkage_status_condition", v, nil)
			if deviceLinkageStatus != nil {
				rst = append(rst, flattenDeviceLinkageCondition(conditionType, deviceLinkageStatus))
			}
		default:
			log.Printf("[ERROR] API returned unknown trigger type= %s", conditionType)
		}
	}

	return rst
}

func flattenDeviceDataCondition(conditionType string, deviceDataCondition interface{}) map[string]interface{} {
	filter := utils.PathSearch("filters|[0]", deviceDataCondition, nil)
	return map[string]interface{}{
		"type": conditionType,
		"device_data_condition": []interface{}{
			map[string]interface{}{
				"device_id":             utils.PathSearch("device_id", deviceDataCondition, nil),
				"product_id":            utils.PathSearch("product_id", deviceDataCondition, nil),
				"path":                  utils.PathSearch("path", filter, nil),
				"operator":              utils.PathSearch("operator", filter, nil),
				"value":                 utils.PathSearch("value", filter, nil),
				"in_values":             utils.PathSearch("in_values", filter, nil),
				"trigger_strategy":      utils.PathSearch("strategy.trigger", filter, nil),
				"data_validatiy_period": utils.PathSearch("strategy.event_valid_time", filter, nil),
			},
		},
	}
}

func flattenSimpleTimerCondition(conditionType string, simpleTimerCondition interface{}) map[string]interface{} {
	return map[string]interface{}{
		"type": conditionType,
		"simple_timer_condition": []interface{}{
			map[string]interface{}{
				"start_time":      utils.PathSearch("start_time", simpleTimerCondition, nil),
				"repeat_interval": int(utils.PathSearch("repeat_interval", simpleTimerCondition, float64(0)).(float64)) / 60,
				"repeat_count":    utils.PathSearch("repeat_count", simpleTimerCondition, nil),
			},
		},
	}
}

func flattenDailyTimerCondition(conditionType string, dailyTimerCondition interface{}) map[string]interface{} {
	return map[string]interface{}{
		"type": conditionType,
		"daily_timer_condition": []interface{}{
			map[string]interface{}{
				"start_time":   utils.PathSearch("time", dailyTimerCondition, nil),
				"days_of_week": utils.PathSearch("days_of_week", dailyTimerCondition, nil),
			},
		},
	}
}

func flattenDeviceLinkageCondition(conditionType string, deviceLinkageCondition interface{}) map[string]interface{} {
	return map[string]interface{}{
		"type": conditionType,
		"device_linkage_status_condition": []interface{}{
			map[string]interface{}{
				"device_id":   utils.PathSearch("device_id", deviceLinkageCondition, nil),
				"product_id":  utils.PathSearch("product_id", deviceLinkageCondition, nil),
				"status_list": utils.PathSearch("status_list", deviceLinkageCondition, nil),
				"duration":    utils.PathSearch("duration", deviceLinkageCondition, nil),
			},
		},
	}
}

func flattenLinkageRuleActions(ruleActions []interface{}) []interface{} {
	if len(ruleActions) == 0 {
		return nil
	}

	rst := make([]interface{}, 0, len(ruleActions))
	for _, v := range ruleActions {
		actionType := utils.PathSearch("type", v, "").(string)
		switch actionType {
		case "DEVICE_CMD":
			deviceCommand := utils.PathSearch("device_command", v, nil)
			commandCmd := utils.PathSearch("device_command.cmd", v, nil)
			if deviceCommand != nil && commandCmd != nil {
				commandBody := utils.PathSearch("command_body", commandCmd, nil)
				jsonStr, err := json.Marshal(commandBody)
				if err != nil {
					log.Printf("[ERROR] Convert the command_body to string failed: %s", err)
				}

				rst = append(rst, flattenActionDeviceCommand(actionType, string(jsonStr), deviceCommand))
			}
		case "SMN_FORWARDING":
			smnForwarding := utils.PathSearch("smn_forwarding", v, nil)
			if smnForwarding != nil {
				rst = append(rst, flattenActionSmnForwording(actionType, smnForwarding))
			}
		case "DEVICE_ALARM":
			deviceAlarm := utils.PathSearch("device_alarm", v, nil)
			if deviceAlarm != nil {
				rst = append(rst, flattenActionDeviceAlarm(actionType, deviceAlarm))
			}
		default:
			log.Printf("[ERROR] API returned unknown action type= %s", actionType)
		}
	}

	return rst
}

func flattenActionDeviceCommand(actionType, commandBody string, deviceCommand interface{}) map[string]interface{} {
	return map[string]interface{}{
		"type": actionType,
		"device_command": []interface{}{
			map[string]interface{}{
				"service_id":       utils.PathSearch("cmd.service_id", deviceCommand, nil),
				"command_name":     utils.PathSearch("cmd.command_name", deviceCommand, nil),
				"command_body":     commandBody,
				"buffer_timeout":   utils.PathSearch("cmd.buffer_timeout", deviceCommand, nil),
				"response_timeout": utils.PathSearch("cmd.response_timeout", deviceCommand, nil),
				"mode":             utils.PathSearch("cmd.mode", deviceCommand, nil),
				"device_id":        utils.PathSearch("device_id", deviceCommand, nil),
			},
		},
	}
}

func flattenActionSmnForwording(actionType string, smnForwarding interface{}) map[string]interface{} {
	return map[string]interface{}{
		"type": actionType,
		"smn_forwarding": []interface{}{
			map[string]interface{}{
				"region":                utils.PathSearch("region_name", smnForwarding, nil),
				"project_id":            utils.PathSearch("project_id", smnForwarding, nil),
				"topic_name":            utils.PathSearch("theme_name", smnForwarding, nil),
				"topic_urn":             utils.PathSearch("topic_urn", smnForwarding, nil),
				"message_content":       utils.PathSearch("message_content", smnForwarding, nil),
				"message_template_name": utils.PathSearch("message_template_name", smnForwarding, nil),
				"message_title":         utils.PathSearch("message_title", smnForwarding, nil),
			},
		},
	}
}

func flattenActionDeviceAlarm(actionType string, deviceAlarm interface{}) map[string]interface{} {
	return map[string]interface{}{
		"type": actionType,
		"device_alarm": []interface{}{
			map[string]interface{}{
				"name":        utils.PathSearch("name", deviceAlarm, nil),
				"type":        utils.PathSearch("alarm_status", deviceAlarm, nil),
				"severity":    utils.PathSearch("severity", deviceAlarm, nil),
				"dimension":   utils.PathSearch("dimension", deviceAlarm, nil),
				"description": utils.PathSearch("description", deviceAlarm, nil),
			},
		},
	}
}

func flattenLinkageRuleEffectivePeriod(conditionGroup interface{}) []interface{} {
	if conditionGroup == nil {
		return nil
	}

	timeRange := utils.PathSearch("time_range", conditionGroup, nil)
	if timeRange == nil {
		return nil
	}

	rst := []interface{}{
		map[string]interface{}{
			"start_time":   utils.PathSearch("start_time", timeRange, nil),
			"end_time":     utils.PathSearch("end_time", timeRange, nil),
			"days_of_week": utils.PathSearch("days_of_week", timeRange, nil),
		},
	}

	return rst
}
