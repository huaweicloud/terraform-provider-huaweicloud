package as

import (
	"context"
	"regexp"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/tidwall/gjson"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/helper/filters"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/helper/httphelper"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/helper/schemas"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API AS GET /autoscaling-api/v2/{project_id}/scaling_policy
func DataSourceASPolicies() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceASPoliciesRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"scaling_policy_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"scaling_policy_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"scaling_policy_type": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"scaling_group_id": {
				Type:          schema.TypeString,
				Optional:      true,
				ConflictsWith: []string{"scaling_resource_id"},
			},
			"scaling_resource_id": {
				Type:          schema.TypeString,
				Optional:      true,
				ConflictsWith: []string{"scaling_group_id"},
			},
			"scaling_resource_type": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"status": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"alarm_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"policies": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"scaling_group_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"status": {
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
						"scaling_resource_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"scaling_resource_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"alarm_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"scheduled_policy": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"launch_time": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"recurrence_type": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"recurrence_value": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"start_time": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"end_time": {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						},
						"action": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"operation": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"instance_number": {
										Type:     schema.TypeInt,
										Computed: true,
									},
									"instance_percentage": {
										Type:     schema.TypeInt,
										Computed: true,
									},
									"limits": {
										Type:     schema.TypeInt,
										Computed: true,
									},
								},
							},
						},
						"meta_data": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"bandwidth_share_type": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"eip_id": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"eip_address": {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						},
						"cool_down_time": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"created_at": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

type PoliciesDSWrapper struct {
	*schemas.ResourceDataWrapper
	Config *config.Config
}

func newPoliciesDSWrapper(d *schema.ResourceData, meta interface{}) *PoliciesDSWrapper {
	return &PoliciesDSWrapper{
		ResourceDataWrapper: schemas.NewSchemaWrapper(d),
		Config:              meta.(*config.Config),
	}
}

func dataSourceASPoliciesRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	wrapper := newPoliciesDSWrapper(d, meta)
	lisAllScaV2PolRst, err := wrapper.ListAllScalingV2Policies()
	if err != nil {
		return diag.FromErr(err)
	}

	id, err := uuid.GenerateUUID()
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(id)

	err = wrapper.listAllScalingV2PoliciesToSchema(lisAllScaV2PolRst)
	if err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func (w *PoliciesDSWrapper) ListAllScalingV2Policies() (*gjson.Result, error) {
	client, err := w.NewClient(w.Config, "autoscaling")
	if err != nil {
		return nil, err
	}

	var scalingResourceId interface{}
	if w.Get("scaling_resource_id") == nil {
		scalingResourceId = w.Get("scaling_group_id")
	} else {
		scalingResourceId = w.Get("scaling_resource_id")
	}

	uri := "/autoscaling-api/v2/{project_id}/scaling_policy"
	params := map[string]any{
		"scaling_resource_id":   scalingResourceId,
		"scaling_resource_type": w.Get("scaling_resource_type"),
		"scaling_policy_name":   w.Get("scaling_policy_name"),
		"scaling_policy_type":   w.Get("scaling_policy_type"),
		"scaling_policy_id":     w.Get("scaling_policy_id"),
		"alarm_id":              w.Get("alarm_id"),
	}

	params = utils.RemoveNil(params)
	return httphelper.New(client).
		Method("GET").
		URI(uri).
		Query(params).
		OffsetPager("scaling_policies", "start_number", "limit", 100).
		Filter(
			filters.New().From("scaling_policies").
				Where("policy_status", "=", w.Get("status")),
		).
		OkCode(200).
		Request().
		Result()
}

func (w *PoliciesDSWrapper) listAllScalingV2PoliciesToSchema(body *gjson.Result) error {
	d := w.ResourceData
	mErr := multierror.Append(nil,
		d.Set("region", w.Config.GetRegion(w.ResourceData)),
		d.Set("policies", schemas.SliceToList(body.Get("scaling_policies"),
			func(policies gjson.Result) any {
				return map[string]any{
					"id":                    policies.Get("scaling_policy_id").Value(),
					"scaling_group_id":      w.setScalingGroupId(policies),
					"name":                  policies.Get("scaling_policy_name").Value(),
					"status":                policies.Get("policy_status").Value(),
					"type":                  policies.Get("scaling_policy_type").Value(),
					"description":           policies.Get("description").Value(),
					"scaling_resource_id":   policies.Get("scaling_resource_id").Value(),
					"scaling_resource_type": policies.Get("scaling_resource_type").Value(),
					"alarm_id":              policies.Get("alarm_id").Value(),
					"scheduled_policy": schemas.SliceToList(policies.Get("scheduled_policy"),
						func(scheduledPolicy gjson.Result) any {
							return map[string]any{
								"launch_time":      w.setScaPolSchPolLauTime(scheduledPolicy),
								"recurrence_type":  scheduledPolicy.Get("recurrence_type").Value(),
								"recurrence_value": scheduledPolicy.Get("recurrence_value").Value(),
								"start_time":       w.setScaPolSchPolStaTime(scheduledPolicy),
								"end_time":         w.setScaPolSchPolEndTime(scheduledPolicy),
							}
						},
					),
					"action": schemas.SliceToList(policies.Get("scaling_policy_action"),
						func(action gjson.Result) any {
							return map[string]any{
								"operation":           action.Get("operation").Value(),
								"instance_number":     action.Get("size").Value(),
								"instance_percentage": action.Get("percentage").Value(),
								"limits":              action.Get("limits").Value(),
							}
						},
					),
					"meta_data": schemas.SliceToList(policies.Get("meta_data"),
						func(metaData gjson.Result) any {
							return map[string]any{
								"bandwidth_share_type": metaData.Get("metadata_bandwidth_share_type").Value(),
								"eip_id":               metaData.Get("metadata_eip_id").Value(),
								"eip_address":          metaData.Get("metadata_eip_address").Value(),
							}
						},
					),
					"cool_down_time": policies.Get("cool_down_time").Value(),
					"created_at":     w.setScaPolCreTime(policies),
				}
			},
		)),
	)
	return mErr.ErrorOrNil()
}

func (*PoliciesDSWrapper) setScalingGroupId(data gjson.Result) string {
	resourceType := data.Get("scaling_resource_type").Value()
	if resourceType == "SCALING_GROUP" {
		return data.Get("scaling_resource_id").String()
	}
	return ""
}

func (*PoliciesDSWrapper) setScaPolSchPolLauTime(data gjson.Result) string {
	var time string
	rex := regexp.MustCompile(`^\d{2}:\d{2}$`)
	timeStr := data.Get("launch_time").String()
	if rex.MatchString(timeStr) {
		time = timeStr
	} else {
		timeStamp := utils.ConvertTimeStrToNanoTimestamp(timeStr, "2006-01-02T15:04Z")
		time = utils.FormatTimeStampRFC3339(timeStamp/1000, false)
	}
	return time
}

func (*PoliciesDSWrapper) setScaPolSchPolStaTime(data gjson.Result) string {
	return utils.FormatTimeStampRFC3339(utils.ConvertTimeStrToNanoTimestamp(data.Get("start_time").String(), "2006-01-02T15:04Z")/1000, false)
}

func (*PoliciesDSWrapper) setScaPolSchPolEndTime(data gjson.Result) string {
	return utils.FormatTimeStampRFC3339(utils.ConvertTimeStrToNanoTimestamp(data.Get("end_time").String(), "2006-01-02T15:04Z")/1000, false)
}

func (*PoliciesDSWrapper) setScaPolCreTime(data gjson.Result) string {
	return utils.FormatTimeStampRFC3339(utils.ConvertTimeStrToNanoTimestamp(data.Get("create_time").String())/1000, false)
}
