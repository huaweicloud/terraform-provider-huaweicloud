package iotda

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"regexp"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/services/iotda/v5/model"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

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
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringLenBetween(1, 128),
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
							ValidateFunc: validation.StringInSlice([]string{
								"DEVICE_DATA",
								"SIMPLE_TIMER",
								"DAILY_TIMER",
							}, false),
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
										ValidateFunc: validation.StringInSlice(
											[]string{">", "<", ">=", "<=", "=", "between"}, false),
									},

									"value": {
										Type:     schema.TypeString,
										Required: true,
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
										Type:         schema.TypeString,
										Required:     true,
										ValidateFunc: validation.StringLenBetween(0, 256),
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

									"description": {
										Type:         schema.TypeString,
										Optional:     true,
										ValidateFunc: validation.StringLenBetween(0, 256),
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
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringLenBetween(0, 256),
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

func ResourceDeviceLinkageRuleCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*config.Config)
	region := c.GetRegion(d)
	client, err := c.HcIoTdaV5Client(region)
	if err != nil {
		return diag.Errorf("error creating IoTDA v5 client: %s", err)
	}

	projectId := c.RegionProjectIDMap[region]
	rule, err := buildDeviceLinkageRuleParams(d, projectId)
	if err != nil {
		return diag.FromErr(err)
	}
	log.Printf("[DEBUG] Create IoTDA device linkage rule params: %#v", rule)

	resp, err := client.CreateRule(&model.CreateRuleRequest{Body: rule})
	if err != nil {
		return diag.Errorf("error creating IoTDA device linkage rule: %s", err)
	}

	if resp.RuleId == nil {
		return diag.Errorf("error creating IoTDA device linkage rule: ID is not found in API response")
	}

	d.SetId(*resp.RuleId)
	return ResourceDeviceLinkageRuleRead(ctx, d, meta)
}

func ResourceDeviceLinkageRuleRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*config.Config)
	region := c.GetRegion(d)
	client, err := c.HcIoTdaV5Client(region)
	if err != nil {
		return diag.Errorf("error creating IoTDA v5 client: %s", err)
	}

	response, err := client.ShowRule(&model.ShowRuleRequest{RuleId: d.Id()})
	if err != nil {
		return common.CheckDeletedDiag(d, err, "error retrieving IoTDA device linkage rule")
	}

	mErr := multierror.Append(
		d.Set("region", region),
		d.Set("name", response.Name),
		d.Set("triggers", flattenLinkageTriggers(response.ConditionGroup)),
		d.Set("actions", flattenLinkageActions(response.Actions)),
		d.Set("trigger_logic", flattenLinkageLogic(response.ConditionGroup)),
		d.Set("enabled", utils.StringValue(response.Status) == "active"),
		d.Set("description", response.Description),
		d.Set("space_id", response.AppId),
		d.Set("effective_period", flattenLinkageEffectivePeriod(response.ConditionGroup)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func ResourceDeviceLinkageRuleUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*config.Config)
	region := c.GetRegion(d)
	client, err := c.HcIoTdaV5Client(region)
	if err != nil {
		return diag.Errorf("error creating IoTDA v5 client: %s", err)
	}

	projectId := c.RegionProjectIDMap[region]
	status := buildLinkageStatus(d.Get("enabled").(bool))
	if d.HasChangeExcept("enabled") {
		// not only change enabled, use update API
		rule, err := buildDeviceLinkageRuleParams(d, projectId)
		if err != nil {
			return diag.FromErr(err)
		}

		_, err = client.UpdateRule(&model.UpdateRuleRequest{
			RuleId: d.Id(),
			Body:   rule,
		})

		if err != nil {
			return diag.Errorf("error updating IoTDA device linkage rule: %s", err)
		}
	} else {
		// only change enabled, use changeStatus API
		_, err = client.ChangeRuleStatus(&model.ChangeRuleStatusRequest{
			RuleId: d.Id(),
			Body: &model.RuleStatus{
				Status: *status,
			},
		})
		if err != nil {
			return diag.Errorf("error updating the IoTDA device linkage status to %s: %s", *status, err)
		}
	}

	return ResourceDeviceLinkageRuleRead(ctx, d, meta)
}

func ResourceDeviceLinkageRuleDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*config.Config)
	region := c.GetRegion(d)
	client, err := c.HcIoTdaV5Client(region)
	if err != nil {
		return diag.Errorf("error creating IoTDA v5 client: %s", err)
	}

	deleteOpts := &model.DeleteRuleRequest{RuleId: d.Id()}
	_, err = client.DeleteRule(deleteOpts)
	if err != nil {
		return diag.Errorf("error deleting IoTDA device linkage rule: %s", err)
	}

	return nil
}

func buildDeviceLinkageRuleParams(d *schema.ResourceData, projectId string) (*model.Rule, error) {
	conditions, err := buildLinkageTriggers(d.Get("triggers").(*schema.Set).List())
	if err != nil {
		return nil, err
	}

	actions, err := buildLinkageActions(d.Get("actions").(*schema.Set).List(), projectId)
	if err != nil {
		return nil, err
	}

	req := model.Rule{
		Name:        d.Get("name").(string),
		Description: utils.StringIgnoreEmpty(d.Get("description").(string)),
		RuleType:    "DEVICE_LINKAGE",
		Status:      buildLinkageStatus(d.Get("enabled").(bool)),
		AppId:       utils.StringIgnoreEmpty(d.Get("space_id").(string)),
		ConditionGroup: &model.ConditionGroup{
			Logic:      utils.String(d.Get("trigger_logic").(string)),
			TimeRange:  buildLinkageTimaRange(d.Get("effective_period").([]interface{})),
			Conditions: &conditions,
		},
		Actions: actions,
	}
	return &req, nil
}

func buildLinkageTimaRange(raw []interface{}) *model.TimeRange {
	if len(raw) != 1 {
		return nil
	}
	target := raw[0].(map[string]interface{})
	rst := model.TimeRange{
		StartTime:  target["start_time"].(string),
		EndTime:    target["end_time"].(string),
		DaysOfWeek: utils.StringIgnoreEmpty(target["days_of_week"].(string)),
	}

	return &rst
}

func buildLinkageTriggers(raw []interface{}) ([]model.RuleCondition, error) {
	rst := make([]model.RuleCondition, len(raw))
	for i, v := range raw {
		target := v.(map[string]interface{})
		triggerType := target["type"].(string)
		ruleCondition, err := buildlinkageTrigger(target, triggerType)
		if err != nil {
			return nil, err
		}
		rst[i] = *ruleCondition
	}

	return rst, nil
}

func buildlinkageTrigger(raw map[string]interface{}, triggerType string) (*model.RuleCondition, error) {
	switch triggerType {
	case "DEVICE_DATA":
		trigger := raw["device_data_condition"].([]interface{})
		if len(trigger) == 0 {
			return nil, fmt.Errorf("device_data_condition is Required when the trigger type is DEVICE_DATA")
		}
		f := trigger[0].(map[string]interface{})
		d := model.RuleCondition{
			Type: triggerType,
			DevicePropertyCondition: &model.DeviceDataCondition{
				DeviceId:  utils.StringIgnoreEmpty(f["device_id"].(string)),
				ProductId: utils.StringIgnoreEmpty(f["product_id"].(string)),
				Filters: &[]model.PropertyFilter{
					{
						Path:     f["path"].(string),
						Operator: f["operator"].(string),
						Value:    f["value"].(string),
						Strategy: &model.Strategy{
							Trigger:        utils.StringIgnoreEmpty(f["trigger_strategy"].(string)),
							EventValidTime: utils.Int32(int32(f["data_validatiy_period"].(int))),
						},
					},
				},
			},
		}
		return &d, nil

	case "SIMPLE_TIMER":
		trigger := raw["simple_timer_condition"].([]interface{})
		if len(trigger) == 0 {
			return nil, fmt.Errorf("simple_timer_condition is Required when the target type is SIMPLE_TIMER")
		}
		f := trigger[0].(map[string]interface{})
		d := model.RuleCondition{
			Type: triggerType,
			SimpleTimerCondition: &model.SimpleTimerType{
				StartTime:      f["start_time"].(string),
				RepeatInterval: int32(f["repeat_interval"].(int) * 60),
				RepeatCount:    int32(f["repeat_count"].(int)),
			},
		}
		return &d, nil

	case "DAILY_TIMER":
		trigger := raw["daily_timer_condition"].([]interface{})
		if len(trigger) == 0 {
			return nil, fmt.Errorf("daily_timer_condition is Required when the target type is DAILY_TIMER")
		}
		f := trigger[0].(map[string]interface{})
		d := model.RuleCondition{
			Type: triggerType,
			DailyTimerCondition: &model.DailyTimerType{
				Time:       f["start_time"].(string),
				DaysOfWeek: utils.StringIgnoreEmpty(f["days_of_week"].(string)),
			},
		}
		return &d, nil

	default:
		return nil, fmt.Errorf("the trigger type= %q is not support", triggerType)
	}
}

func buildLinkageActions(raw []interface{}, projectId string) ([]model.RuleAction, error) {
	rst := make([]model.RuleAction, len(raw))
	for i, v := range raw {
		target := v.(map[string]interface{})
		action, err := buildlinkageAction(target, projectId)
		if err != nil {
			return nil, err
		}
		rst[i] = *action
	}

	return rst, nil
}

func buildlinkageAction(raw map[string]interface{}, projectId string) (*model.RuleAction, error) {
	actionType := raw["type"].(string)
	switch actionType {
	case "DEVICE_CMD":
		action := raw["device_command"].([]interface{})
		if len(action) == 0 {
			return nil, fmt.Errorf("device_command is Required when the trigger type is DEVICE_CMD")
		}
		f := action[0].(map[string]interface{})

		var commandBody interface{} = f["command_body"]
		d := model.RuleAction{
			Type: actionType,
			DeviceCommand: &model.ActionDeviceCommand{
				DeviceId: utils.String(f["device_id"].(string)),
				Cmd: &model.Cmd{
					CommandName: f["command_name"].(string),
					ServiceId:   f["service_id"].(string),
					CommandBody: &commandBody, // Json string and Map, all support.
				},
			},
		}
		return &d, nil

	case "SMN_FORWARDING":
		trigger := raw["smn_forwarding"].([]interface{})
		if len(trigger) == 0 {
			return nil, fmt.Errorf("smn_forwarding is Required when the target type is SMN_FORWARDING")
		}
		f := trigger[0].(map[string]interface{})

		projectIdStr := f["project_id"].(string)
		if projectIdStr == "" {
			projectIdStr = projectId
		}
		d := model.RuleAction{
			Type: actionType,
			SmnForwarding: &model.ActionSmnForwarding{
				RegionName:     f["region"].(string),
				ProjectId:      projectIdStr,
				ThemeName:      f["topic_name"].(string),
				TopicUrn:       f["topic_urn"].(string),
				MessageTitle:   f["message_title"].(string),
				MessageContent: f["message_content"].(string),
			},
		}
		return &d, nil

	case "DEVICE_ALARM":
		trigger := raw["device_alarm"].([]interface{})
		if len(trigger) == 0 {
			return nil, fmt.Errorf("device_alarm is Required when the target type is DEVICE_ALARM")
		}
		f := trigger[0].(map[string]interface{})
		d := model.RuleAction{
			Type: actionType,
			DeviceAlarm: &model.ActionDeviceAlarm{
				Name:        f["name"].(string),
				AlarmStatus: f["type"].(string),
				Severity:    f["severity"].(string),
				Description: utils.StringIgnoreEmpty(f["description"].(string)),
			},
		}
		return &d, nil

	default:
		return nil, fmt.Errorf("the action type= %q is not support", actionType)
	}
}

func buildLinkageStatus(enabled bool) *string {
	status := "active"
	if !enabled {
		status = "inactive"
	}
	return &status
}

func flattenLinkageTriggers(conditionGroup *model.ConditionGroup) []interface{} {
	var rst []interface{}
	if conditionGroup == nil || conditionGroup.Conditions == nil {
		return rst
	}

	rst = make([]interface{}, 0, len(*conditionGroup.Conditions))
	for _, v := range *conditionGroup.Conditions {
		switch v.Type {
		case "DEVICE_DATA":
			if v.DevicePropertyCondition != nil && v.DevicePropertyCondition.Filters != nil &&
				len(*v.DevicePropertyCondition.Filters) != 0 {
				filter := *v.DevicePropertyCondition.Filters
				rst = append(rst, map[string]interface{}{
					"type": v.Type,
					"device_data_condition": []interface{}{
						map[string]interface{}{
							"device_id":             v.DevicePropertyCondition.DeviceId,
							"product_id":            v.DevicePropertyCondition.ProductId,
							"path":                  filter[0].Path,
							"operator":              filter[0].Operator,
							"value":                 filter[0].Value,
							"trigger_strategy":      filter[0].Strategy.Trigger,
							"data_validatiy_period": filter[0].Strategy.EventValidTime,
						},
					},
				})
			}
		case "SIMPLE_TIMER":
			if v.SimpleTimerCondition != nil {
				rst = append(rst, map[string]interface{}{
					"type": v.Type,
					"simple_timer_condition": []interface{}{
						map[string]interface{}{
							"start_time":      v.SimpleTimerCondition.StartTime,
							"repeat_interval": v.SimpleTimerCondition.RepeatInterval / 60,
							"repeat_count":    v.SimpleTimerCondition.RepeatCount,
						},
					},
				})
			}
		case "DAILY_TIMER":
			if v.DailyTimerCondition != nil {
				rst = append(rst, map[string]interface{}{
					"type": v.Type,
					"daily_timer_condition": []interface{}{
						map[string]interface{}{
							"start_time":   v.DailyTimerCondition.Time,
							"days_of_week": v.DailyTimerCondition.DaysOfWeek,
						},
					},
				})
			}
		default:
			log.Printf("[ERROR] API returned unknown trigger type= %s", v.Type)
		}

	}

	return rst
}

func flattenLinkageActions(actions *[]model.RuleAction) []interface{} {
	var rst []interface{}
	if actions == nil || len(*actions) == 0 {
		return rst
	}

	rst = make([]interface{}, 0, len(*actions))
	for _, v := range *actions {
		switch v.Type {
		case "DEVICE_CMD":
			if v.DeviceCommand != nil && v.DeviceCommand.Cmd != nil {
				jsonStr, err := json.Marshal(v.DeviceCommand.Cmd.CommandBody)
				if err != nil {
					log.Printf("[ERROR] Convert the command_body to string failed: %s", err)
				}

				rst = append(rst, map[string]interface{}{
					"type": v.Type,
					"device_command": []interface{}{
						map[string]interface{}{
							"service_id":   v.DeviceCommand.Cmd.ServiceId,
							"command_name": v.DeviceCommand.Cmd.CommandName,
							"command_body": string(jsonStr),
							"device_id":    v.DeviceCommand.DeviceId,
						},
					},
				})
			}
		case "SMN_FORWARDING":
			if v.SmnForwarding != nil {
				rst = append(rst, map[string]interface{}{
					"type": v.Type,
					"smn_forwarding": []interface{}{
						map[string]interface{}{
							"region":          v.SmnForwarding.RegionName,
							"project_id":      v.SmnForwarding.ProjectId,
							"topic_name":      v.SmnForwarding.ThemeName,
							"topic_urn":       v.SmnForwarding.TopicUrn,
							"message_content": v.SmnForwarding.MessageContent,
							"message_title":   v.SmnForwarding.MessageTitle,
						},
					},
				})
			}
		case "DEVICE_ALARM":
			if v.DeviceAlarm != nil {
				rst = append(rst, map[string]interface{}{
					"type": v.Type,
					"device_alarm": []interface{}{
						map[string]interface{}{
							"name":        v.DeviceAlarm.Name,
							"type":        v.DeviceAlarm.AlarmStatus,
							"severity":    v.DeviceAlarm.Severity,
							"description": v.DeviceAlarm.Description,
						},
					},
				})
			}
		default:
			log.Printf("[ERROR] API returned unknown action type= %s", v.Type)
		}

	}

	return rst
}

func flattenLinkageEffectivePeriod(conditionGroup *model.ConditionGroup) []interface{} {
	var rst []interface{}
	if conditionGroup == nil || conditionGroup.TimeRange == nil {
		return rst
	}
	rst = []interface{}{
		map[string]interface{}{
			"start_time":   conditionGroup.TimeRange.StartTime,
			"end_time":     conditionGroup.TimeRange.EndTime,
			"days_of_week": conditionGroup.TimeRange.DaysOfWeek,
		},
	}
	return rst
}

func flattenLinkageLogic(v *model.ConditionGroup) *string {
	if v == nil {
		return nil
	}
	return v.Logic
}
