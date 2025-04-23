package iotda

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"regexp"
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

// @API IoTDA POST /v5/iot/{project_id}/rules
// @API IoTDA PUT /v5/iot/{project_id}/rules/{rule_id}/status
// @API IoTDA DELETE /v5/iot/{project_id}/rules/{rule_id}
// @API IoTDA GET /v5/iot/{project_id}/rules/{rule_id}
// @API IoTDA PUT /v5/iot/{project_id}/rules/{rule_id}
func ResourceDeviceLinkageRule() *schema.Resource {
	return &schema.Resource{
		CreateContext: ResourceDeviceLinkageRuleCreate,
		UpdateContext: ResourceDeviceLinkageRuleUpdate,
		DeleteContext: ResourceDeviceLinkageRuleDelete,
		ReadContext:   ResourceDeviceLinkageRuleRead,
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
				Type:     schema.TypeString,
				Required: true,
			},

			"space_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"triggers": {
				Type:     schema.TypeSet,
				Required: true,
				MaxItems: 10,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"type": {
							Type:     schema.TypeString,
							Required: true,
						},

						"device_data_condition": {
							Type:     schema.TypeList,
							Optional: true,
							Computed: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"device_id": {
										Type:     schema.TypeString,
										Optional: true,
									},

									"product_id": {
										Type:     schema.TypeString,
										Optional: true,
									},

									"path": {
										Type:     schema.TypeString,
										Required: true,
									},

									"operator": {
										Type:     schema.TypeString,
										Required: true,
									},

									"value": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"in_values": {
										Type:     schema.TypeList,
										Optional: true,
										Elem:     &schema.Schema{Type: schema.TypeString},
									},
									"trigger_strategy": {
										Type:         schema.TypeString,
										Optional:     true,
										ValidateFunc: validation.StringInSlice([]string{"pulse", "reverse"}, false),
										Computed:     true,
									},

									"data_validatiy_period": {
										Type:     schema.TypeInt,
										Optional: true,
										Default:  300,
									},
								},
							},
						},

						"simple_timer_condition": {
							Type:     schema.TypeList,
							Optional: true,
							Computed: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"start_time": {
										Type:     schema.TypeString,
										Required: true,
									},

									"repeat_interval": {
										Type:     schema.TypeInt,
										Required: true,
									},

									"repeat_count": {
										Type:     schema.TypeInt,
										Required: true,
									},
								},
							},
						},

						"daily_timer_condition": {
							Type:     schema.TypeList,
							Optional: true,
							Computed: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"start_time": {
										Type:     schema.TypeString,
										Required: true,
										ValidateFunc: validation.StringMatch(
											regexp.MustCompile(`^([0-1][0-9]|2[0-3]):([0-5][0-9])$`),
											"The format is: `HH:mm`"),
									},

									"days_of_week": {
										Type:     schema.TypeString,
										Optional: true,
										Computed: true,
									},
								},
							},
						},
						"device_linkage_status_condition": {
							Type:     schema.TypeList,
							Optional: true,
							Computed: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"device_id": {
										Type:     schema.TypeString,
										Optional: true,
										Computed: true,
									},
									// When both `device_id` and `product_id` are set simultaneously, the API does not
									// return the `product_id` field, so there is no need to add the Computed attribute.
									"product_id": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"status_list": {
										Type:     schema.TypeList,
										Optional: true,
										Computed: true,
										Elem:     &schema.Schema{Type: schema.TypeString},
									},
									"duration": {
										Type:     schema.TypeInt,
										Optional: true,
										Computed: true,
									},
								},
							},
						},
					},
				},
			},

			"actions": {
				Type:     schema.TypeSet,
				Required: true,
				MaxItems: 10,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"type": {
							Type:     schema.TypeString,
							Required: true,
							ValidateFunc: validation.StringInSlice([]string{
								"DEVICE_CMD",
								"SMN_FORWARDING",
								"DEVICE_ALARM",
							}, false),
						},

						"device_command": {
							Type:     schema.TypeList,
							Optional: true,
							Computed: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"device_id": {
										Type:     schema.TypeString,
										Required: true,
									},

									"service_id": {
										Type:     schema.TypeString,
										Required: true,
									},

									"command_name": {
										Type:     schema.TypeString,
										Required: true,
									},

									"command_body": {
										Type:     schema.TypeString,
										Required: true,
									},
									"buffer_timeout": {
										Type:     schema.TypeInt,
										Optional: true,
									},
									"response_timeout": {
										Type:     schema.TypeInt,
										Optional: true,
									},
									"mode": {
										Type:     schema.TypeString,
										Optional: true,
									},
								},
							},
						},

						"smn_forwarding": {
							Type:     schema.TypeList,
							Optional: true,
							Computed: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"region": {
										Type:     schema.TypeString,
										Required: true,
									},

									"topic_name": {
										Type:     schema.TypeString,
										Required: true,
									},

									"topic_urn": {
										Type:     schema.TypeString,
										Required: true,
									},

									"message_title": {
										Type:     schema.TypeString,
										Required: true,
									},

									"message_content": {
										Type:     schema.TypeString,
										Optional: true,
										Computed: true,
									},
									"message_template_name": {
										Type:     schema.TypeString,
										Optional: true,
										Computed: true,
									},
									"project_id": {
										Type:     schema.TypeString,
										Optional: true,
										Computed: true,
									},
								},
							},
						},

						"device_alarm": {
							Type:     schema.TypeList,
							Optional: true,
							Computed: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"name": {
										Type:     schema.TypeString,
										Required: true,
									},

									"type": {
										Type:     schema.TypeString,
										Required: true,
										ValidateFunc: validation.StringInSlice([]string{
											"fault",
											"recovery",
										}, false),
									},

									"severity": {
										Type:     schema.TypeString,
										Required: true,
										ValidateFunc: validation.StringInSlice([]string{
											"warning",
											"minor",
											"major",
											"critical",
										}, false),
									},
									"dimension": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"description": {
										Type:     schema.TypeString,
										Optional: true,
									},
								},
							},
						},
					},
				},
			},

			"trigger_logic": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"and", "or"}, false),
				Default:      "and",
			},

			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},

			"enabled": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
			},

			"effective_period": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"start_time": {
							Type:     schema.TypeString,
							Required: true,
							ValidateFunc: validation.StringMatch(regexp.MustCompile(`^([0-1][0-9]|2[0-3]):([0-5][0-9])$`),
								"The format is: `HH:mm`"),
						},

						"end_time": {
							Type:     schema.TypeString,
							Required: true,
							ValidateFunc: validation.StringMatch(regexp.MustCompile(`^([0-1][0-9]|2[0-3]):([0-5][0-9])$`),
								"The format is: `HH:mm`"),
						},

						"days_of_week": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func buildTriggerDeviceDataBodyParams(triggerMap map[string]interface{}, triggerType string) map[string]interface{} {
	return map[string]interface{}{
		"type": triggerType,
		"device_property_condition": map[string]interface{}{
			"device_id":  utils.ValueIgnoreEmpty(triggerMap["device_id"]),
			"product_id": utils.ValueIgnoreEmpty(triggerMap["product_id"]),
			"filters": []map[string]interface{}{
				{
					"path":      triggerMap["path"],
					"operator":  triggerMap["operator"],
					"value":     utils.ValueIgnoreEmpty(triggerMap["value"]),
					"in_values": triggerMap["in_values"],
					"strategy": map[string]interface{}{
						"trigger":          utils.ValueIgnoreEmpty(triggerMap["trigger_strategy"]),
						"event_valid_time": triggerMap["data_validatiy_period"],
					},
				},
			},
		},
	}
}

func buildTriggerSimpleTimerBodyParams(triggerMap map[string]interface{}, triggerType string) map[string]interface{} {
	return map[string]interface{}{
		"type": triggerType,
		"simple_timer_condition": map[string]interface{}{
			"start_time":      triggerMap["start_time"],
			"repeat_interval": triggerMap["repeat_interval"].(int) * 60,
			"repeat_count":    triggerMap["repeat_count"],
		},
	}
}

func buildTriggerDailyTimerBodyParams(triggerMap map[string]interface{}, triggerType string) map[string]interface{} {
	return map[string]interface{}{
		"type": triggerType,
		"daily_timer_condition": map[string]interface{}{
			"time":         triggerMap["start_time"],
			"days_of_week": utils.ValueIgnoreEmpty(triggerMap["days_of_week"]),
		},
	}
}

func buildTriggerDeviceLinkageStatusBodyParams(triggerMap map[string]interface{}, triggerType string) map[string]interface{} {
	return map[string]interface{}{
		"type": triggerType,
		"device_linkage_status_condition": map[string]interface{}{
			"device_id":   utils.ValueIgnoreEmpty(triggerMap["device_id"]),
			"product_id":  utils.ValueIgnoreEmpty(triggerMap["product_id"]),
			"status_list": triggerMap["status_list"],
			"duration":    utils.ValueIgnoreEmpty(triggerMap["duration"]),
		},
	}
}

func buildLinkageRuleTriggersParams(d *schema.ResourceData) ([]map[string]interface{}, error) {
	rawArray := d.Get("triggers").(*schema.Set).List()
	rst := make([]map[string]interface{}, 0, len(rawArray))
	for _, v := range rawArray {
		rawMap := v.(map[string]interface{})
		triggerType := rawMap["type"].(string)

		switch triggerType {
		case "DEVICE_DATA":
			triggerArray := rawMap["device_data_condition"].([]interface{})
			if len(triggerArray) == 0 {
				return nil, errors.New("device_data_condition is Required when the trigger type is DEVICE_DATA")
			}
			rst = append(rst, buildTriggerDeviceDataBodyParams(triggerArray[0].(map[string]interface{}), triggerType))
		case "SIMPLE_TIMER":
			triggerArray := rawMap["simple_timer_condition"].([]interface{})
			if len(triggerArray) == 0 {
				return nil, errors.New("simple_timer_condition is Required when the target type is SIMPLE_TIMER")
			}
			rst = append(rst, buildTriggerSimpleTimerBodyParams(triggerArray[0].(map[string]interface{}), triggerType))
		case "DAILY_TIMER":
			triggerArray := rawMap["daily_timer_condition"].([]interface{})
			if len(triggerArray) == 0 {
				return nil, errors.New("daily_timer_condition is Required when the target type is DAILY_TIMER")
			}
			rst = append(rst, buildTriggerDailyTimerBodyParams(triggerArray[0].(map[string]interface{}), triggerType))
		case "DEVICE_LINKAGE_STATUS":
			triggerArray := rawMap["device_linkage_status_condition"].([]interface{})
			if len(triggerArray) == 0 {
				return nil, errors.New("device_linkage_status_condition is Required when the target type is DEVICE_LINKAGE_STATUS")
			}
			rst = append(rst, buildTriggerDeviceLinkageStatusBodyParams(triggerArray[0].(map[string]interface{}), triggerType))
		default:
			return nil, fmt.Errorf("the trigger type (%s) is not support", triggerType)
		}
	}

	return rst, nil
}

func buildActionDeviceCmdBodyParams(actionMap map[string]interface{}, actionType string) map[string]interface{} {
	return map[string]interface{}{
		"type": actionType,
		"device_command": map[string]interface{}{
			"device_id": actionMap["device_id"],
			"cmd": map[string]interface{}{
				"command_name":     actionMap["command_name"],
				"service_id":       actionMap["service_id"],
				"command_body":     actionMap["command_body"],
				"buffer_timeout":   utils.ValueIgnoreEmpty(actionMap["buffer_timeout"]),
				"response_timeout": utils.ValueIgnoreEmpty(actionMap["response_timeout"]),
				"mode":             utils.ValueIgnoreEmpty(actionMap["mode"]),
			},
		},
	}
}

func buildActionSmnForwardingProjectIDParam(actionMap map[string]interface{}, projectID string) string {
	rawProjectID := actionMap["project_id"].(string)
	if rawProjectID == "" {
		rawProjectID = projectID
	}

	return rawProjectID
}

func buildActionSmnForwardingBodyParams(actionMap map[string]interface{}, actionType, projectID string) map[string]interface{} {
	return map[string]interface{}{
		"type": actionType,
		"smn_forwarding": map[string]interface{}{
			"region_name":           actionMap["region"],
			"project_id":            buildActionSmnForwardingProjectIDParam(actionMap, projectID),
			"theme_name":            actionMap["topic_name"],
			"topic_urn":             actionMap["topic_urn"],
			"message_title":         actionMap["message_title"],
			"message_content":       utils.ValueIgnoreEmpty(actionMap["message_content"]),
			"message_template_name": utils.ValueIgnoreEmpty(actionMap["message_template_name"]),
		},
	}
}

func buildActionDeviceAlarmBodyParams(actionMap map[string]interface{}, actionType string) map[string]interface{} {
	return map[string]interface{}{
		"type": actionType,
		"device_alarm": map[string]interface{}{
			"name":         actionMap["name"],
			"alarm_status": actionMap["type"],
			"severity":     actionMap["severity"],
			"dimension":    utils.ValueIgnoreEmpty(actionMap["dimension"]),
			"description":  utils.ValueIgnoreEmpty(actionMap["description"]),
		},
	}
}

func buildLinkageRuleActionsParams(d *schema.ResourceData, projectID string) ([]map[string]interface{}, error) {
	rawArray := d.Get("actions").(*schema.Set).List()
	rst := make([]map[string]interface{}, 0, len(rawArray))
	for _, v := range rawArray {
		rawMap := v.(map[string]interface{})
		actionType := rawMap["type"].(string)
		switch actionType {
		case "DEVICE_CMD":
			actionArray := rawMap["device_command"].([]interface{})
			if len(actionArray) == 0 {
				return nil, errors.New("device_command is Required when the trigger type is DEVICE_CMD")
			}
			rst = append(rst, buildActionDeviceCmdBodyParams(actionArray[0].(map[string]interface{}), actionType))
		case "SMN_FORWARDING":
			actionArray := rawMap["smn_forwarding"].([]interface{})
			if len(actionArray) == 0 {
				return nil, errors.New("smn_forwarding is Required when the target type is SMN_FORWARDING")
			}
			rst = append(rst, buildActionSmnForwardingBodyParams(actionArray[0].(map[string]interface{}), actionType, projectID))
		case "DEVICE_ALARM":
			actionArray := rawMap["device_alarm"].([]interface{})
			if len(actionArray) == 0 {
				return nil, errors.New("device_alarm is Required when the target type is DEVICE_ALARM")
			}
			rst = append(rst, buildActionDeviceAlarmBodyParams(actionArray[0].(map[string]interface{}), actionType))
		default:
			return nil, fmt.Errorf("the action type= %s is not support", actionType)
		}
	}
	return rst, nil
}

func buildLinkageRuleStatusParam(d *schema.ResourceData) string {
	if d.Get("enabled").(bool) {
		return "active"
	}

	return "inactive"
}

func buildLinkageRuleConditionTimeRangeParams(d *schema.ResourceData) map[string]interface{} {
	rawArray := d.Get("effective_period").([]interface{})
	if len(rawArray) != 1 {
		return nil
	}

	rawMap := rawArray[0].(map[string]interface{})
	return map[string]interface{}{
		"start_time":   rawMap["start_time"],
		"end_time":     rawMap["end_time"],
		"days_of_week": utils.ValueIgnoreEmpty(rawMap["days_of_week"]),
	}
}

func buildLinkageRuleBodyParams(d *schema.ResourceData, conditions, actions []map[string]interface{}) map[string]interface{} {
	return map[string]interface{}{
		"name":        d.Get("name"),
		"description": utils.ValueIgnoreEmpty(d.Get("description")),
		"rule_type":   "DEVICE_LINKAGE",
		"status":      buildLinkageRuleStatusParam(d),
		"app_id":      utils.ValueIgnoreEmpty(d.Get("space_id")),
		"condition_group": map[string]interface{}{
			"logic":      d.Get("trigger_logic"),
			"time_range": buildLinkageRuleConditionTimeRangeParams(d),
			"conditions": conditions,
		},
		"actions": actions,
	}
}

func ResourceDeviceLinkageRuleCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		httpUrl = "v5/iot/{project_id}/rules"
		product = "iotda"
	)

	isDerived := WithDerivedAuth(cfg, region)
	client, err := cfg.NewServiceClientWithDerivedAuth(product, region, isDerived)
	if err != nil {
		return diag.Errorf("error creating IoTDA client: %s", err)
	}

	conditions, err := buildLinkageRuleTriggersParams(d)
	if err != nil {
		return diag.FromErr(err)
	}

	actions, err := buildLinkageRuleActionsParams(d, client.ProjectID)
	if err != nil {
		return diag.FromErr(err)
	}

	requestPath := client.Endpoint + httpUrl
	requestPath = strings.ReplaceAll(requestPath, "{project_id}", client.ProjectID)
	requestOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		JSONBody:         utils.RemoveNil(buildLinkageRuleBodyParams(d, conditions, actions)),
	}

	resp, err := client.Request("POST", requestPath, &requestOpt)
	if err != nil {
		return diag.Errorf("error creating IoTDA device linkage rule: %s", err)
	}

	respBody, err := utils.FlattenResponse(resp)
	if err != nil {
		return diag.FromErr(err)
	}

	ruleID := utils.PathSearch("rule_id", respBody, "").(string)
	if ruleID == "" {
		return diag.Errorf("error creating IoTDA device linkage rule: ID is not found in API response")
	}

	d.SetId(ruleID)
	return ResourceDeviceLinkageRuleRead(ctx, d, meta)
}

func flattenDeviceDataAttribute(respBody interface{}, conditionType string, filters []interface{}) map[string]interface{} {
	return map[string]interface{}{
		"type": conditionType,
		"device_data_condition": []interface{}{
			map[string]interface{}{
				"device_id":             utils.PathSearch("device_property_condition.device_id", respBody, nil),
				"product_id":            utils.PathSearch("device_property_condition.product_id", respBody, nil),
				"path":                  utils.PathSearch("[0]|path", filters, nil),
				"operator":              utils.PathSearch("[0]|operator", filters, nil),
				"value":                 utils.PathSearch("[0]|value", filters, nil),
				"in_values":             utils.PathSearch("[0]|in_values", filters, nil),
				"trigger_strategy":      utils.PathSearch("[0]|strategy.trigger", filters, nil),
				"data_validatiy_period": utils.PathSearch("[0]|strategy.event_valid_time", filters, nil),
			},
		},
	}
}

func flattenSimpleTimerAttribute(condition interface{}, conditionType string) map[string]interface{} {
	return map[string]interface{}{
		"type": conditionType,
		"simple_timer_condition": []interface{}{
			map[string]interface{}{
				"start_time":      utils.PathSearch("start_time", condition, nil),
				"repeat_interval": utils.PathSearch("repeat_interval", condition, float64(0)).(float64) / 60,
				"repeat_count":    utils.PathSearch("repeat_count", condition, nil),
			},
		},
	}
}

func flattenDailyTimerAttribute(condition interface{}, conditionType string) map[string]interface{} {
	return map[string]interface{}{
		"type": conditionType,
		"daily_timer_condition": []interface{}{
			map[string]interface{}{
				"start_time":   utils.PathSearch("time", condition, nil),
				"days_of_week": utils.PathSearch("days_of_week", condition, nil),
			},
		},
	}
}

func flattenDeviceLinkageStatusAttribute(condition interface{}, conditionType string) map[string]interface{} {
	return map[string]interface{}{
		"type": conditionType,
		"device_linkage_status_condition": []interface{}{
			map[string]interface{}{
				"device_id":   utils.PathSearch("device_id", condition, nil),
				"product_id":  utils.PathSearch("product_id", condition, nil),
				"status_list": utils.PathSearch("status_list", condition, nil),
				"duration":    utils.PathSearch("duration", condition, nil),
			},
		},
	}
}

func flattenTriggerAttributes(respBody interface{}) []interface{} {
	conditions := utils.PathSearch("condition_group.conditions", respBody, make([]interface{}, 0)).([]interface{})
	rst := make([]interface{}, 0, len(conditions))
	for _, v := range conditions {
		conditionType := utils.PathSearch("type", v, "").(string)
		switch conditionType {
		case "DEVICE_DATA":
			filters := utils.PathSearch("device_property_condition.filters", v, make([]interface{}, 0)).([]interface{})
			if len(filters) == 0 {
				continue
			}
			rst = append(rst, flattenDeviceDataAttribute(v, conditionType, filters))
		case "SIMPLE_TIMER":
			condition := utils.PathSearch("simple_timer_condition", v, nil)
			if condition == nil {
				continue
			}
			rst = append(rst, flattenSimpleTimerAttribute(condition, conditionType))
		case "DAILY_TIMER":
			condition := utils.PathSearch("daily_timer_condition", v, nil)
			if condition == nil {
				continue
			}
			rst = append(rst, flattenDailyTimerAttribute(condition, conditionType))
		case "DEVICE_LINKAGE_STATUS":
			condition := utils.PathSearch("device_linkage_status_condition", v, nil)
			if condition == nil {
				continue
			}
			rst = append(rst, flattenDeviceLinkageStatusAttribute(condition, conditionType))
		default:
			log.Printf("[ERROR] API returned unknown trigger type= %s", conditionType)
		}
	}

	return rst
}

func flattenDeviceCommandAttribute(cmd interface{}, actionType string, actionResp interface{}) map[string]interface{} {
	jsonStr, err := json.Marshal(utils.PathSearch("command_body", cmd, nil))
	if err != nil {
		log.Printf("[ERROR] Convert the command_body to string failed: %s", err)
	}

	return map[string]interface{}{
		"type": actionType,
		"device_command": []interface{}{
			map[string]interface{}{
				"service_id":       utils.PathSearch("service_id", cmd, nil),
				"command_name":     utils.PathSearch("command_name", cmd, nil),
				"command_body":     string(jsonStr),
				"buffer_timeout":   utils.PathSearch("buffer_timeout", cmd, nil),
				"response_timeout": utils.PathSearch("response_timeout", cmd, nil),
				"mode":             utils.PathSearch("mode", cmd, nil),
				"device_id":        utils.PathSearch("device_command.device_id", actionResp, nil),
			},
		},
	}
}

func flattenSmnForwardingAttribute(smnForwarding interface{}, actionType string) map[string]interface{} {
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

func flattenDeviceAlarmAttribute(deviceAlarm interface{}, actionType string) map[string]interface{} {
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

func flattenActionsAttributes(respBody interface{}) []interface{} {
	actions := utils.PathSearch("actions", respBody, make([]interface{}, 0)).([]interface{})
	rst := make([]interface{}, 0, len(actions))
	for _, v := range actions {
		actionType := utils.PathSearch("type", v, "").(string)
		switch actionType {
		case "DEVICE_CMD":
			cmd := utils.PathSearch("device_command.cmd", v, nil)
			if cmd == nil {
				continue
			}

			rst = append(rst, flattenDeviceCommandAttribute(cmd, actionType, v))
		case "SMN_FORWARDING":
			smnForwarding := utils.PathSearch("smn_forwarding", v, nil)
			if smnForwarding == nil {
				continue
			}
			rst = append(rst, flattenSmnForwardingAttribute(smnForwarding, actionType))
		case "DEVICE_ALARM":
			deviceAlarm := utils.PathSearch("device_alarm", v, nil)
			if deviceAlarm == nil {
				continue
			}
			rst = append(rst, flattenDeviceAlarmAttribute(deviceAlarm, actionType))
		default:
			log.Printf("[ERROR] API returned unknown action type= %s", actionType)
		}
	}

	return rst
}

func flattenEffectivePeriodAttribute(respBody interface{}) []interface{} {
	timeRange := utils.PathSearch("condition_group.time_range", respBody, nil)
	if timeRange == nil {
		return nil
	}

	return []interface{}{
		map[string]interface{}{
			"start_time":   utils.PathSearch("start_time", timeRange, nil),
			"end_time":     utils.PathSearch("end_time", timeRange, nil),
			"days_of_week": utils.PathSearch("days_of_week", timeRange, nil),
		},
	}
}

func ResourceDeviceLinkageRuleRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		httpUrl = "v5/iot/{project_id}/rules/{rule_id}"
		product = "iotda"
	)

	isDerived := WithDerivedAuth(cfg, region)
	client, err := cfg.NewServiceClientWithDerivedAuth(product, region, isDerived)
	if err != nil {
		return diag.Errorf("error creating IoTDA client: %s", err)
	}

	requestPath := client.Endpoint + httpUrl
	requestPath = strings.ReplaceAll(requestPath, "{project_id}", client.ProjectID)
	requestPath = strings.ReplaceAll(requestPath, "{rule_id}", d.Id())
	requestOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	resp, err := client.Request("GET", requestPath, &requestOpt)
	if err != nil {
		return common.CheckDeletedDiag(d, err, "error retrieving IoTDA device linkage rule")
	}

	respBody, err := utils.FlattenResponse(resp)
	if err != nil {
		return diag.FromErr(err)
	}

	mErr := multierror.Append(
		d.Set("region", region),
		d.Set("name", utils.PathSearch("name", respBody, nil)),
		d.Set("triggers", flattenTriggerAttributes(respBody)),
		d.Set("actions", flattenActionsAttributes(respBody)),
		d.Set("trigger_logic", utils.PathSearch("condition_group.logic", respBody, nil)),
		d.Set("enabled", utils.PathSearch("status", respBody, "").(string) == "active"),
		d.Set("description", utils.PathSearch("description", respBody, nil)),
		d.Set("space_id", utils.PathSearch("app_id", respBody, nil)),
		d.Set("effective_period", flattenEffectivePeriodAttribute(respBody)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func ResourceDeviceLinkageRuleUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		product = "iotda"
	)

	isDerived := WithDerivedAuth(cfg, region)
	client, err := cfg.NewServiceClientWithDerivedAuth(product, region, isDerived)
	if err != nil {
		return diag.Errorf("error creating IoTDA client: %s", err)
	}

	if d.HasChangeExcept("enabled") {
		conditions, err := buildLinkageRuleTriggersParams(d)
		if err != nil {
			return diag.FromErr(err)
		}

		actions, err := buildLinkageRuleActionsParams(d, client.ProjectID)
		if err != nil {
			return diag.FromErr(err)
		}

		requestPath := client.Endpoint + "v5/iot/{project_id}/rules/{rule_id}"
		requestPath = strings.ReplaceAll(requestPath, "{project_id}", client.ProjectID)
		requestPath = strings.ReplaceAll(requestPath, "{rule_id}", d.Id())
		requestOpt := golangsdk.RequestOpts{
			KeepResponseBody: true,
			JSONBody:         utils.RemoveNil(buildLinkageRuleBodyParams(d, conditions, actions)),
		}

		_, err = client.Request("PUT", requestPath, &requestOpt)
		if err != nil {
			return diag.Errorf("error updating IoTDA device linkage rule: %s", err)
		}
	}

	if d.HasChange("enabled") {
		requestPath := client.Endpoint + "v5/iot/{project_id}/rules/{rule_id}/status"
		requestPath = strings.ReplaceAll(requestPath, "{project_id}", client.ProjectID)
		requestPath = strings.ReplaceAll(requestPath, "{rule_id}", d.Id())
		requestOpt := golangsdk.RequestOpts{
			KeepResponseBody: true,
			JSONBody: map[string]interface{}{
				"status": buildLinkageRuleStatusParam(d),
			},
		}

		_, err = client.Request("PUT", requestPath, &requestOpt)
		if err != nil {
			return diag.Errorf("error updating the IoTDA device linkage rule status : %s", err)
		}
	}

	return ResourceDeviceLinkageRuleRead(ctx, d, meta)
}

func ResourceDeviceLinkageRuleDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		httpUrl = "v5/iot/{project_id}/rules/{rule_id}"
		product = "iotda"
	)

	isDerived := WithDerivedAuth(cfg, region)
	client, err := cfg.NewServiceClientWithDerivedAuth(product, region, isDerived)
	if err != nil {
		return diag.Errorf("error creating IoTDA client: %s", err)
	}

	requestPath := client.Endpoint + httpUrl
	requestPath = strings.ReplaceAll(requestPath, "{project_id}", client.ProjectID)
	requestPath = strings.ReplaceAll(requestPath, "{rule_id}", d.Id())
	requestOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	_, err = client.Request("DELETE", requestPath, &requestOpt)
	if err != nil {
		return diag.Errorf("error deleting IoTDA device linkage rule: %s", err)
	}

	return nil
}
