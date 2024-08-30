package aom

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API AOM GET /v4/{project_id}/alarm-rules-template
func DataSourceAlarmRulesTemplates() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceAlarmRulesTemplatesRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"template_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"enterprise_project_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"templates": {
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
						"type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"alarm_template_spec_list": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"related_cloud_service": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"related_cce_clusters": {
										Type:     schema.TypeList,
										Computed: true,
										Elem:     &schema.Schema{Type: schema.TypeString},
									},
									"related_prometheus_instances": {
										Type:     schema.TypeList,
										Computed: true,
										Elem:     &schema.Schema{Type: schema.TypeString},
									},
									"alarm_notification": dataSourceSchemeAlarmNotifications(),
									"alarm_template_spec_items": {
										Type:     schema.TypeList,
										Computed: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"alarm_rule_name": {
													Type:     schema.TypeString,
													Computed: true,
												},
												"alarm_rule_type": {
													Type:     schema.TypeString,
													Computed: true,
												},
												"alarm_rule_description": {
													Type:     schema.TypeString,
													Computed: true,
												},
												"event_alarm_spec":  dataSourceSchemeEventAlarmTemplateSpec(),
												"metric_alarm_spec": dataSourceSchemeMetricAlarmTemplateSpec(),
											},
										},
									},
								},
							},
						},
						"templating": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"list": {
										Type:     schema.TypeList,
										Computed: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"name": {
													Type:     schema.TypeString,
													Computed: true,
												},
												"query": {
													Type:     schema.TypeString,
													Computed: true,
												},
												"type": {
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
							},
						},
						"description": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"enterprise_project_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"created_at": {
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

func dataSourceSchemeEventAlarmTemplateSpec() *schema.Schema {
	return &schema.Schema{
		Type:     schema.TypeList,
		Computed: true,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"alarm_subtype": {
					Type:     schema.TypeString,
					Computed: true,
				},
				"alarm_source": {
					Type:     schema.TypeString,
					Computed: true,
				},
				"event_source": {
					Type:     schema.TypeString,
					Computed: true,
				},
				"trigger_conditions": {
					Type:     schema.TypeList,
					Computed: true,
					Elem: &schema.Resource{
						Schema: map[string]*schema.Schema{
							"trigger_type": {
								Type:     schema.TypeString,
								Computed: true,
							},
							"event_name": {
								Type:     schema.TypeString,
								Computed: true,
							},
							"thresholds": {
								Type:     schema.TypeMap,
								Computed: true,
								Elem:     &schema.Schema{Type: schema.TypeInt},
							},
							"aggregation_window": {
								Type:     schema.TypeInt,
								Computed: true,
							},
							"operator": {
								Type:     schema.TypeString,
								Computed: true,
							},
							"frequency": {
								Type:     schema.TypeString,
								Computed: true,
							},
						},
					},
				},
				"monitor_objects": {
					Type:     schema.TypeList,
					Computed: true,
					Elem: &schema.Schema{
						Type: schema.TypeMap,
						Elem: &schema.Schema{Type: schema.TypeString},
					},
				},
				"monitor_object_templates": {
					Type:     schema.TypeList,
					Computed: true,
					Elem:     &schema.Schema{Type: schema.TypeString},
				},
			},
		},
	}
}

func dataSourceSchemeMetricAlarmTemplateSpec() *schema.Schema {
	return &schema.Schema{
		Type:     schema.TypeList,
		Computed: true,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"alarm_subtype": {
					Type:     schema.TypeString,
					Computed: true,
				},
				"alarm_source": {
					Type:     schema.TypeString,
					Computed: true,
				},
				"monitor_type": {
					Type:     schema.TypeString,
					Computed: true,
				},
				"recovery_conditions": {
					Type:     schema.TypeList,
					Computed: true,
					Elem: &schema.Resource{
						Schema: map[string]*schema.Schema{
							"recovery_timeframe": {
								Type:     schema.TypeInt,
								Computed: true,
							},
						},
					},
				},
				"trigger_conditions": dataSourceSchemeTemplateMetricTriggerConditions(),
				"no_data_conditions": {
					Type:     schema.TypeList,
					Computed: true,
					Elem: &schema.Resource{
						Schema: map[string]*schema.Schema{
							"no_data_timeframe": {
								Type:     schema.TypeInt,
								Computed: true,
							},
							"no_data_alert_state": {
								Type:     schema.TypeString,
								Computed: true,
							},
							"notify_no_data": {
								Type:     schema.TypeBool,
								Computed: true,
							},
						},
					},
				},
				"alarm_tags": {
					Type:     schema.TypeList,
					Computed: true,
					Elem: &schema.Resource{
						Schema: map[string]*schema.Schema{
							"auto_tags": {
								Type:     schema.TypeList,
								Computed: true,
								Elem:     &schema.Schema{Type: schema.TypeString},
							},
							"custom_tags": {
								Type:     schema.TypeList,
								Computed: true,
								Elem:     &schema.Schema{Type: schema.TypeString},
							},
							"custom_annotations": {
								Type:     schema.TypeList,
								Computed: true,
								Elem:     &schema.Schema{Type: schema.TypeString},
							},
						},
					},
				},
			},
		},
	}
}

func dataSourceSchemeTemplateMetricTriggerConditions() *schema.Schema {
	return &schema.Schema{
		Type:     schema.TypeList,
		Computed: true,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"metric_query_mode": {
					Type:     schema.TypeString,
					Computed: true,
				},
				"metric_name": {
					Type:     schema.TypeString,
					Computed: true,
				},
				"promql": {
					Type:     schema.TypeString,
					Computed: true,
				},
				"aggregation_window": {
					Type:     schema.TypeString,
					Computed: true,
				},
				"query_match": {
					Type:     schema.TypeString,
					Computed: true,
				},
				"aggregate_type": {
					Type:     schema.TypeString,
					Computed: true,
				},
				"metric_labels": {
					Type:     schema.TypeList,
					Computed: true,
					Elem:     &schema.Schema{Type: schema.TypeString},
				},
				"aggregation_type": {
					Type:     schema.TypeString,
					Computed: true,
				},
				"operator": {
					Type:     schema.TypeString,
					Computed: true,
				},
				"thresholds": {
					Type:     schema.TypeMap,
					Computed: true,
					Elem:     &schema.Schema{Type: schema.TypeString},
				},
				"trigger_times": {
					Type:     schema.TypeString,
					Computed: true,
				},
				"trigger_type": {
					Type:     schema.TypeString,
					Computed: true,
				},
				"trigger_interval": {
					Type:     schema.TypeString,
					Computed: true,
				},
				"expression": {
					Type:     schema.TypeString,
					Computed: true,
				},
				"mix_promql": {
					Type:     schema.TypeString,
					Computed: true,
				},
				"metric_statistic_method": {
					Type:     schema.TypeString,
					Computed: true,
				},
				"metric_namespace": {
					Type:     schema.TypeString,
					Computed: true,
				},
				"metric_unit": {
					Type:     schema.TypeString,
					Computed: true,
				},
				"promql_expr": {
					Type:     schema.TypeList,
					Computed: true,
					Elem:     &schema.Schema{Type: schema.TypeString},
				},
				"promql_for": {
					Type:     schema.TypeString,
					Computed: true,
				},
				"aom_monitor_level": {
					Type:     schema.TypeString,
					Computed: true,
				},
			},
		},
	}
}

func dataSourceAlarmRulesTemplatesRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	client, err := cfg.NewServiceClient("aom", region)
	if err != nil {
		return diag.Errorf("error creating AOM client: %s", err)
	}

	templates, err := getAlarmRulesTemplates(cfg, client, d)
	if err != nil {
		return diag.FromErr(err)
	}

	id, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("error generating UUID")
	}
	d.SetId(id)

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("templates", templates),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func getAlarmRulesTemplates(cfg *config.Config, client *golangsdk.ServiceClient, d *schema.ResourceData) (interface{}, error) {
	listHttpUrl := "v4/{project_id}/alarm-rules-template"
	listPath := client.Endpoint + listHttpUrl
	listPath = strings.ReplaceAll(listPath, "{project_id}", client.ProjectID)
	listPath += buildListAlarmRulesTemplatesQueryParams(d)
	listOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      buildHeadersForDataSource(cfg, d),
	}

	listResp, err := client.Request("GET", listPath, &listOpt)
	if err != nil {
		return nil, fmt.Errorf("error retrieving alarm rules templates: %s", err)
	}
	listRespBody, err := utils.FlattenResponse(listResp)
	if err != nil {
		return nil, fmt.Errorf("error flattening alarm rules templates: %s", err)
	}

	templates := listRespBody.([]interface{})
	results := make([]map[string]interface{}, 0, len(templates))
	for _, template := range templates {
		results = append(results, flattenAlarmRulesTemplates(template))
	}

	return results, nil
}

func flattenAlarmRulesTemplates(template interface{}) map[string]interface{} {
	rst := map[string]interface{}{
		"name":                  utils.PathSearch("alarm_rule_template_name", template, nil),
		"id":                    utils.PathSearch("alarm_rule_template_id", template, nil),
		"type":                  utils.PathSearch("alarm_rule_template_type", template, nil),
		"description":           utils.PathSearch("alarm_rule_template_description", template, nil),
		"enterprise_project_id": utils.PathSearch("enterprise_project_id", template, nil),
		"alarm_template_spec_list": flattenTemplateSpecList(
			utils.PathSearch("alarm_template_spec_list", template, make([]interface{}, 0)).([]interface{})),
		"templating": flattenTemplating(utils.PathSearch("templating", template, nil)),
		"created_at": utils.FormatTimeStampRFC3339(
			int64(utils.PathSearch("create_time", template, float64(0)).(float64))/1000, true),
		"updated_at": utils.FormatTimeStampRFC3339(
			int64(utils.PathSearch("modify_time", template, float64(0)).(float64))/1000, true),
	}
	return rst
}

func buildListAlarmRulesTemplatesQueryParams(d *schema.ResourceData) string {
	res := ""
	if v, ok := d.GetOk("template_id"); ok {
		res = fmt.Sprintf("%s?id=%v", res, v)
	}

	return res
}
