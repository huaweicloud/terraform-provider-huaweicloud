package ces

import (
	"context"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/tidwall/gjson"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/helper/httphelper"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/helper/schemas"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

func DataSourceCesAlarmRules() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceCesAlarmRulesRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `Specifies the region in which to query the resource. If omitted, the provider-level region will be used.`,
			},
			"alarm_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specifies the alarm rule ID.`,
			},
			"name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specifies the name of an alarm rule.`,
			},
			"namespace": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specifies the namespace of a service.`,
			},
			"resource_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specifies the alarm resource ID.`,
			},
			"alarms": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: `The alarm rule list.`,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"alarm_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The alarm rule ID.`,
						},
						"alarm_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The alarm rule name.`,
						},
						"alarm_description": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The alarm rule description.`,
						},
						"namespace": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The namespace of a service.`,
						},
						"condition": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: `The alarm triggering condition list.`,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"metric_name": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: `The metric name of a resource.`,
									},
									"period": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: `The monitoring period of a metric.`,
									},
									"filter": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: `The filter method.`,
									},
									"comparison_operator": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: `The comparison condition of alarm thresholds.`,
									},
									"value": {
										Type:        schema.TypeFloat,
										Computed:    true,
										Description: `The alarm threshold.`,
									},
									"unit": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: `The metric unit.`,
									},
									"count": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: `The number of times that the alarm triggering conditions are met.`,
									},
									"suppress_duration": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: `The interval for triggering an alarm if the alarm persists.`,
									},
									"alarm_level": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: `The alarm severity.`,
									},
								},
							},
						},
						"resources": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: `The resource list.`,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"resource_group_id": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: `The resource group ID.`,
									},
									"resource_group_name": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: `The resource group name.`,
									},
									"dimensions": {
										Type:        schema.TypeList,
										Computed:    true,
										Description: `The dimension information.`,
										Elem:        alarmsResourcesDimensionsElem(),
									},
								},
							},
						},
						"alarm_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The alarm rule type.`,
						},
						"alarm_enabled": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: `Whether to generate alarms when the alarm triggering conditions are met.`,
						},
						"alarm_action_enabled": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: `Whether to enable the action to be triggered by an alarm.`,
						},
						"alarm_actions": {
							Type:        schema.TypeList,
							Computed:    true,
							Elem:        alarmsActionsElem(),
							Description: `The action to be triggered by an alarm.`,
						},
						"ok_actions": {
							Type:        schema.TypeList,
							Computed:    true,
							Elem:        alarmsActionsElem(),
							Description: `The action to be triggered after an alarm is cleared.`,
						},
						"notification_begin_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The time when the alarm notification was enabled.`,
						},
						"notification_end_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The time when the alarm notification was disabled.`,
						},
						"enterprise_project_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The enterprise project ID.`,
						},
						"alarm_template_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The ID of an alarm template associated with an alarm rule.`,
						},
					},
				},
			},
		},
	}
}

// alarmsResourcesDimensionsElem
// The Elem of "alarms.resources.dimensions"
func alarmsResourcesDimensionsElem() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The name of the metric dimension.`,
			},
		},
	}
}

func alarmsActionsElem() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The notification type.`,
			},
			"notification_list": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: `The list of objects to be notified if the alarm status changes.`,
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
		},
	}
}

type AlarmRulesDSWrapper struct {
	*schemas.ResourceDataWrapper
	Config *config.Config
}

func newCesAlarmRulesDSWrapper(d *schema.ResourceData, meta interface{}) *AlarmRulesDSWrapper {
	return &AlarmRulesDSWrapper{
		ResourceDataWrapper: schemas.NewSchemaWrapper(d),
		Config:              meta.(*config.Config),
	}
}

func dataSourceCesAlarmRulesRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	wrapper := newCesAlarmRulesDSWrapper(d, meta)
	listAlarmRulesRst, err := wrapper.ListAlarmRules()
	if err != nil {
		return diag.FromErr(err)
	}

	id, err := uuid.GenerateUUID()
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(id)

	err = wrapper.listAlarmRulesToSchema(listAlarmRulesRst)
	if err != nil {
		return diag.FromErr(err)
	}

	return nil
}

// @API CES GET /v2/{project_id}/alarms
func (w *AlarmRulesDSWrapper) ListAlarmRules() (*gjson.Result, error) {
	client, err := w.NewClient(w.Config, "ces")
	if err != nil {
		return nil, err
	}

	uri := "/v2/{project_id}/alarms"
	params := map[string]any{
		"alarm_id":    w.Get("alarm_id"),
		"name":        w.Get("name"),
		"namespace":   w.Get("namespace"),
		"resource_id": w.Get("resource_id"),
	}
	params = utils.RemoveNil(params)
	return httphelper.New(client).
		Method("GET").
		URI(uri).
		Query(params).
		OffsetPager("alarms", "offset", "limit", 100).
		Request().
		Result()
}

func (w *AlarmRulesDSWrapper) listAlarmRulesToSchema(body *gjson.Result) error {
	d := w.ResourceData
	mErr := multierror.Append(nil,
		d.Set("region", w.Config.GetRegion(w.ResourceData)),
		d.Set("alarms", schemas.SliceToList(body.Get("alarms"),
			func(alarms gjson.Result) any {
				return map[string]any{
					"alarm_id":   alarms.Get("alarm_id").Value(),
					"alarm_name": alarms.Get("name").Value(),
					"resources": schemas.SliceToList(alarms.Get("resources"),
						func(resources gjson.Result) any {
							return map[string]any{
								"resource_group_id":   resources.Get("resource_group_id").Value(),
								"resource_group_name": resources.Get("resource_group_name").Value(),
								"dimensions":          w.setAlaResDim(resources),
							}
						},
					),
					"alarm_type":              alarms.Get("type").Value(),
					"namespace":               alarms.Get("namespace").Value(),
					"alarm_description":       alarms.Get("description").Value(),
					"alarm_enabled":           alarms.Get("enabled").Value(),
					"notification_begin_time": alarms.Get("notification_begin_time").Value(),
					"notification_end_time":   alarms.Get("notification_end_time").Value(),
					"alarm_actions": schemas.SliceToList(alarms.Get("alarm_notifications"),
						func(action gjson.Result) any {
							return map[string]any{
								"type":              action.Get("type").Value(),
								"notification_list": schemas.SliceToStrList(action.Get("notification_list")),
							}
						}),
					"ok_actions": schemas.SliceToList(alarms.Get("ok_notifications"),
						func(action gjson.Result) any {
							return map[string]any{
								"type":              action.Get("type").Value(),
								"notification_list": schemas.SliceToStrList(action.Get("notification_list")),
							}
						}),
					"condition": schemas.SliceToList(alarms.Get("policies"),
						func(condition gjson.Result) any {
							return map[string]any{
								"metric_name":         condition.Get("metric_name").Value(),
								"period":              condition.Get("period").Value(),
								"filter":              condition.Get("filter").Value(),
								"comparison_operator": condition.Get("comparison_operator").Value(),
								"value":               condition.Get("value").Value(),
								"unit":                condition.Get("unit").Value(),
								"count":               condition.Get("count").Value(),
								"suppress_duration":   condition.Get("suppress_duration").Value(),
								"alarm_level":         condition.Get("level").Value(),
							}
						},
					),
					"alarm_action_enabled":  alarms.Get("notification_enabled").Value(),
					"enterprise_project_id": alarms.Get("enterprise_project_id").Value(),
					"alarm_template_id":     alarms.Get("alarm_template_id").Value(),
				}
			},
		)),
	)
	return mErr.ErrorOrNil()
}

func (*AlarmRulesDSWrapper) setAlaResDim(resources gjson.Result) any {
	return schemas.SliceToList(resources.Get("dimensions"), func(dimensions gjson.Result) any {
		return map[string]any{
			"name": dimensions.Get("name").Value(),
		}
	})
}
