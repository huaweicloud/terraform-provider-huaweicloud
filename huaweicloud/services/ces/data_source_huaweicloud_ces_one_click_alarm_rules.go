package ces

import (
	"context"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/tidwall/gjson"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/helper/httphelper"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/helper/schemas"
)

func DataSourceCesOneClickAlarmRules() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceCesOneClickAlarmRulesRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `Specifies the region in which to query the resource. If omitted, the provider-level region will be used.`,
			},
			"one_click_alarm_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `Specifies the one-click monitoring ID for a service.`,
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
							Description: `The ID of an alarm rule.`,
						},
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The alarm rule name.`,
						},
						"description": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The supplementary information about an alarm rule.`,
						},
						"namespace": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The metric namespace.`,
						},
						"policies": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: `The alarm policy list.`,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"alarm_policy_id": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: `The alarm policy ID.`,
									},
									"metric_name": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: `The metric name.`,
									},
									"period": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: `How often to generate an alarm.`,
									},
									"comparison_operator": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: `The operator of an alarm threshold.`,
									},
									"filter": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: `The roll up method.`,
									},
									"value": {
										Type:        schema.TypeFloat,
										Computed:    true,
										Description: `The threshold.`,
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
										Description: `The suppression period.`,
									},
									"level": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: `The alarm severity.`,
									},
									"enabled": {
										Type:        schema.TypeBool,
										Computed:    true,
										Description: `Whether the one-click monitoring is enabled.`,
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
										Elem:        alarmResourcesDimensionsElem(),
									},
								},
							},
						},
						"type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The alarm rule type.`,
						},
						"enabled": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: `Whether to generate alarms when the alarm triggering conditions are met.`,
						},
						"notification_enabled": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: `Whether the alarm notification is enabled.`,
						},
						"alarm_notifications": {
							Type:        schema.TypeList,
							Computed:    true,
							Elem:        actionSchema(),
							Description: `The action to be triggered by an alarm.`,
						},
						"ok_notifications": {
							Type:        schema.TypeList,
							Computed:    true,
							Elem:        actionSchema(),
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
					},
				},
			},
		},
	}
}

// alarmResourcesDimensionsElem
// The Elem of "alarm.resources.dimensions"
func alarmResourcesDimensionsElem() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The name of the metric dimension.`,
			},
			"value": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The value of the metric dimension.`,
			},
		},
	}
}

func actionSchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The notification type.`,
			},
			"notification_list": {
				Type:        schema.TypeList,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: `The list of objects to be notified if the alarm status changes.`,
			},
		},
	}
	return &sc
}

type OneClickAlarmRulesDSWrapper struct {
	*schemas.ResourceDataWrapper
	Config *config.Config
}

func newOneClickAlarmRulesDSWrapper(d *schema.ResourceData, meta interface{}) *OneClickAlarmRulesDSWrapper {
	return &OneClickAlarmRulesDSWrapper{
		ResourceDataWrapper: schemas.NewSchemaWrapper(d),
		Config:              meta.(*config.Config),
	}
}

func dataSourceCesOneClickAlarmRulesRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	wrapper := newOneClickAlarmRulesDSWrapper(d, meta)
	lisOneCliAlaRulRst, err := wrapper.ListOneClickAlarmRules()
	if err != nil {
		return diag.FromErr(err)
	}

	id, err := uuid.GenerateUUID()
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(id)

	err = wrapper.listOneClickAlarmRulesToSchema(lisOneCliAlaRulRst)
	if err != nil {
		return diag.FromErr(err)
	}

	return nil
}

// @API CES GET /v2/{project_id}/one-click-alarms/{one_click_alarm_id}/alarms
func (w *OneClickAlarmRulesDSWrapper) ListOneClickAlarmRules() (*gjson.Result, error) {
	client, err := w.NewClient(w.Config, "ces")
	if err != nil {
		return nil, err
	}

	uri := "/v2/{project_id}/one-click-alarms/{one_click_alarm_id}/alarms"
	uri = strings.ReplaceAll(uri, "{one_click_alarm_id}", w.Get("one_click_alarm_id").(string))
	return httphelper.New(client).
		Method("GET").
		URI(uri).
		Request().
		Result()
}

func (w *OneClickAlarmRulesDSWrapper) listOneClickAlarmRulesToSchema(body *gjson.Result) error {
	d := w.ResourceData
	mErr := multierror.Append(nil,
		d.Set("region", w.Config.GetRegion(w.ResourceData)),
		d.Set("alarms", schemas.SliceToList(body.Get("alarms"),
			func(alarms gjson.Result) any {
				return map[string]any{
					"alarm_id":                alarms.Get("alarm_id").Value(),
					"name":                    alarms.Get("name").Value(),
					"namespace":               alarms.Get("namespace").Value(),
					"type":                    alarms.Get("type").Value(),
					"notification_enabled":    alarms.Get("notification_enabled").Value(),
					"notification_begin_time": alarms.Get("notification_begin_time").Value(),
					"notification_end_time":   alarms.Get("notification_end_time").Value(),
					"description":             alarms.Get("description").Value(),
					"policies": schemas.SliceToList(alarms.Get("policies"),
						func(policies gjson.Result) any {
							return map[string]any{
								"alarm_policy_id":     policies.Get("alarm_policy_id").Value(),
								"metric_name":         policies.Get("metric_name").Value(),
								"period":              policies.Get("period").Value(),
								"comparison_operator": policies.Get("comparison_operator").Value(),
								"filter":              policies.Get("filter").Value(),
								"value":               policies.Get("value").Value(),
								"unit":                policies.Get("unit").Value(),
								"count":               policies.Get("count").Value(),
								"suppress_duration":   policies.Get("suppress_duration").Value(),
								"level":               policies.Get("level").Value(),
								"enabled":             policies.Get("enabled").Value(),
							}
						},
					),
					"resources": schemas.SliceToList(alarms.Get("resources"),
						func(resources gjson.Result) any {
							return map[string]any{
								"resource_group_id":   resources.Get("resource_group_id").Value(),
								"resource_group_name": resources.Get("resource_group_name").Value(),
								"dimensions":          w.setAlaResDim(resources),
							}
						},
					),
					"enabled": alarms.Get("enabled").Value(),
					"alarm_notifications": schemas.SliceToList(alarms.Get("alarm_notifications"),
						func(alarmNotifications gjson.Result) any {
							return map[string]any{
								"type":              alarmNotifications.Get("type").Value(),
								"notification_list": schemas.SliceToStrList(alarmNotifications.Get("notification_list")),
							}
						},
					),
					"ok_notifications": schemas.SliceToList(alarms.Get("ok_notifications"),
						func(okNotifications gjson.Result) any {
							return map[string]any{
								"type":              okNotifications.Get("type").Value(),
								"notification_list": schemas.SliceToStrList(okNotifications.Get("notification_list")),
							}
						},
					),
				}
			},
		)),
	)
	return mErr.ErrorOrNil()
}

func (*OneClickAlarmRulesDSWrapper) setAlaResDim(resources gjson.Result) any {
	return schemas.SliceToList(resources.Get("dimensions"), func(dimensions gjson.Result) any {
		return map[string]any{
			"name":  dimensions.Get("name").Value(),
			"value": dimensions.Get("value").Value(),
		}
	})
}
