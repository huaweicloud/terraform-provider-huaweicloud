package lts

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

// @API LTS POST /v2/{project_id}/{domain_id}/lts/alarms/sql-alarm/query
func DataSourceAlarms() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceAlarmsRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"type": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The type of the alarm to be queried.`,
			},
			"whether_custom_field": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: `Whether to customize the query time range.`,
			},
			"time_range": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The time range of the alarm to be queried, in minutes.`,
			},
			"search": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The keyword search criteria.`,
			},
			"alarm_level_ids": {
				Type:        schema.TypeList,
				Optional:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: `The list of alarm levels.`,
			},
			"start_time": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: `The start time of a customized time segment, in milliseconds.`,
			},
			"end_time": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: `The end time of a customized time segment, in milliseconds.`,
			},
			"sort": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"order_by": {
							Type:        schema.TypeList,
							Required:    true,
							Elem:        &schema.Schema{Type: schema.TypeString},
							Description: `The fields to be sorted.`,
						},
						"order": {
							Type:        schema.TypeString,
							Required:    true,
							Description: `The sort mode of the alarm.`,
						},
					},
				},
				Description: `The sort criteria of the queried alarms.`,
			},
			"step": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: `The step of the query, in milliseconds.`,
			},
			"alarms": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: `The list of the queried alarms.`,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The ID of the alarm.`,
						},
						"type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The type of the alarm.`,
						},
						"timeout": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: `The time when the alarm is automatically cleared, in milliseconds.`,
						},
						"arrives_at": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: `The time when the alarm arrives, in milliseconds.`,
						},
						"ends_at": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: `The time when the alarm is cleared, in milliseconds.`,
						},
						"starts_at": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: `The time when the alarm is generated, in milliseconds.`,
						},
						"annotations": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: `The details of the alarm.`,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"type": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: `The type of the alarm rule.`,
									},
									"message": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: `The detail information of the alarm.`,
									},
									"log_info": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: `The log information of the alarm.`,
									},
									"current_value": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: `The current value of the alarm.`,
									},
									"old_annotations": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: `The raw data of the alarm detail.`,
									},
									"alarm_action_rule_name": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: `The name of the alarm action rule.`,
									},
									"alarm_rule_alias": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: `The alias of the alarm rule.`,
									},
									"alarm_rule_url": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: `The URL of the alarm rule.`,
									},
									"alarm_status": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: `The status of the alarm trigger.`,
									},
									"condition_expression": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: `The condition expression of the alarm trigger.`,
									},
									"condition_expression_with_value": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: `The condition of the alarm trigger.`,
									},
									"notification_frequency": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: `The notification frequency of the alarm.`,
									},
									"recovery_policy": {
										Type:        schema.TypeBool,
										Computed:    true,
										Description: `Whether the alarm is recovered.`,
									},
									"frequency": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: `The frequency of the alarm.`,
									},
								},
							},
						},
						"metadata": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: `The metadata of the alarm.`,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"event_id": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: `The ID of the alarm rule.`,
									},
									"event_name": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: `The name of the alarm rule.`,
									},
									"event_type": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: `The mode of the alarm.`,
									},
									"event_severity": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: `The level of the alarm.`,
									},
									"resource_provider": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: `The source of the alarm.`,
									},
									"lts_alarm_type": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: `The type of the alarm rule.`,
									},
									"resource_id": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: `The ID of the resource.`,
									},
									"resource_type": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: `The type of the resource.`,
									},
									"log_group_name": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: `The original name of the log group.`,
									},
									"log_stream_name": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: `The original name of the log stream.`,
									},
									"event_subtype": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: `The type of the alarm.`,
									},
								},
							},
						},
					},
				},
			},
		},
	}
}

func queryAlarms(client *golangsdk.ServiceClient, d *schema.ResourceData, domainId string) ([]interface{}, error) {
	var (
		httpUrl = "v2/{project_id}/{domain_id}/lts/alarms/sql-alarm/query?type={alarm_type}"
		result  = make([]interface{}, 0)
		marker  = ""
	)

	listPath := client.Endpoint + httpUrl
	listPath = strings.ReplaceAll(listPath, "{project_id}", client.ProjectID)
	listPath = strings.ReplaceAll(listPath, "{domain_id}", domainId)
	listPath = strings.ReplaceAll(listPath, "{alarm_type}", d.Get("type").(string))
	opt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"Content-Type": "application/json;charset=UTF-8"},
		JSONBody:         utils.RemoveNil(buildGetAlarmsBodyParams(d)),
	}

	for {
		listPathWithMarker := listPath
		if marker != "" {
			listPathWithMarker += fmt.Sprintf("&marker=%s", marker)
		}

		resp, err := client.Request("POST", listPathWithMarker, &opt)
		if err != nil {
			return nil, err
		}

		respBody, err := utils.FlattenResponse(resp)
		if err != nil {
			return nil, err
		}

		alarms := utils.PathSearch("events", respBody, make([]interface{}, 0)).([]interface{})
		if len(alarms) == 0 {
			break
		}

		result = append(result, alarms...)
		marker = utils.PathSearch("page_info.next_marker", respBody, "").(string)
		if marker == "" {
			break
		}
	}

	return result, nil
}

func buildGetAlarmsBodyParams(d *schema.ResourceData) map[string]interface{} {
	return map[string]interface{}{
		"whether_custom_field": utils.ValueIgnoreEmpty(d.Get("whether_custom_field")),
		"time_range":           utils.ValueIgnoreEmpty(d.Get("time_range")),
		"search":               utils.ValueIgnoreEmpty(d.Get("search")),
		"alarm_level_ids":      utils.ValueIgnoreEmpty(utils.ExpandToStringList(d.Get("alarm_level_ids").([]interface{}))),
		"start_time":           utils.ValueIgnoreEmpty(d.Get("start_time")),
		"end_time":             utils.ValueIgnoreEmpty(d.Get("end_time")),
		"sort":                 buildAlarmsSort(d.Get("sort").([]interface{})),
		"step":                 utils.ValueIgnoreEmpty(d.Get("step")),
	}
}

func buildAlarmsSort(sort []interface{}) map[string]interface{} {
	if len(sort) == 0 {
		return nil
	}
	return map[string]interface{}{
		"order":    utils.PathSearch("order", sort[0], nil),
		"order_by": utils.PathSearch("order_by", sort[0], nil),
	}
}

func dataSourceAlarmsRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg    = meta.(*config.Config)
		region = cfg.GetRegion(d)
	)

	client, err := cfg.NewServiceClient("lts", region)
	if err != nil {
		return diag.Errorf("error creating LTS client: %s", err)
	}

	alarms, err := queryAlarms(client, d, cfg.DomainID)
	if err != nil {
		return diag.Errorf("error retrieving LTS alarms: %s", err)
	}

	id, err := uuid.GenerateUUID()
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(id)

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("alarms", flattenAlarms(alarms)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenAlarms(alarms []interface{}) []map[string]interface{} {
	if len(alarms) == 0 {
		return nil
	}

	result := make([]map[string]interface{}, len(alarms))
	for i, v := range alarms {
		result[i] = map[string]interface{}{
			"id":          utils.PathSearch("id", v, nil),
			"type":        utils.PathSearch("type", v, nil),
			"timeout":     utils.PathSearch("timeout", v, nil),
			"arrives_at":  utils.PathSearch("arrives_at", v, nil),
			"ends_at":     utils.PathSearch("ends_at", v, nil),
			"starts_at":   utils.PathSearch("starts_at", v, nil),
			"annotations": flattenAlarmAnnotations(utils.PathSearch("annotations", v, nil)),
			"metadata":    flattenAlarmmetadata(utils.PathSearch("metadata", v, nil)),
		}
	}

	return result
}

func flattenAlarmAnnotations(annotations interface{}) []interface{} {
	if annotations == nil {
		return nil
	}

	return []interface{}{
		map[string]interface{}{
			"type":                            utils.PathSearch("type", annotations, nil),
			"message":                         utils.PathSearch("message", annotations, nil),
			"log_info":                        utils.PathSearch("log_info", annotations, nil),
			"current_value":                   utils.PathSearch("current_value", annotations, nil),
			"old_annotations":                 utils.PathSearch("old_annotations", annotations, nil),
			"alarm_action_rule_name":          utils.PathSearch("alarm_action_rule_name", annotations, nil),
			"alarm_rule_alias":                utils.PathSearch("alarm_rule_alias", annotations, nil),
			"alarm_rule_url":                  utils.PathSearch("alarm_rule_url", annotations, nil),
			"alarm_status":                    utils.PathSearch("alarm_status", annotations, nil),
			"condition_expression":            utils.PathSearch("condition_expression", annotations, nil),
			"condition_expression_with_value": utils.PathSearch("condition_expression_with_value", annotations, nil),
			"notification_frequency":          utils.PathSearch("notification_frequency", annotations, nil),
			"recovery_policy":                 utils.PathSearch("recovery_policy", annotations, nil),
			"frequency":                       utils.PathSearch("frequency", annotations, nil),
		},
	}
}

func flattenAlarmmetadata(metadata interface{}) []interface{} {
	if metadata == nil {
		return nil
	}

	return []interface{}{
		map[string]interface{}{
			"event_id":          utils.PathSearch("event_id", metadata, nil),
			"event_name":        utils.PathSearch("event_name", metadata, nil),
			"event_type":        utils.PathSearch("event_type", metadata, nil),
			"event_severity":    utils.PathSearch("event_severity", metadata, nil),
			"resource_provider": utils.PathSearch("resource_provider", metadata, nil),
			"lts_alarm_type":    utils.PathSearch("lts_alarm_type", metadata, nil),
			"resource_id":       utils.PathSearch("resource_id", metadata, nil),
			"resource_type":     utils.PathSearch("resource_type", metadata, nil),
			"log_group_name":    utils.PathSearch("log_group_name", metadata, nil),
			"log_stream_name":   utils.PathSearch("log_stream_name", metadata, nil),
			"event_subtype":     utils.PathSearch("event_subtype", metadata, nil),
		},
	}
}
