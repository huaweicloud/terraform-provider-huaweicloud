package ces

import (
	"context"
	"regexp"
	"time"

	"github.com/chnsz/golangsdk/openstack/cloudeyeservice/alarmrule"
	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils/fmtp"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils/logp"
)

const nameCESAR = "CES-AlarmRule"

var cesAlarmActions = schema.Schema{
	Type:     schema.TypeList,
	Optional: true,
	ForceNew: true,
	Elem: &schema.Resource{
		Schema: map[string]*schema.Schema{
			"type": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
				ValidateFunc: validation.StringInSlice([]string{
					"notification", "autoscaling",
				}, false),
			},

			"notification_list": {
				Type:     schema.TypeList,
				MaxItems: 5,
				Required: true,
				ForceNew: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
		},
	},
}

func ResourceAlarmRule() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceAlarmRuleCreate,
		ReadContext:   resourceAlarmRuleRead,
		UpdateContext: resourceAlarmRuleUpdate,
		DeleteContext: resourceAlarmRuleDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
			Update: schema.DefaultTimeout(10 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},

			"alarm_name": {
				Type:     schema.TypeString,
				Required: true,
				ValidateFunc: validation.All(
					validation.StringLenBetween(1, 128),
					validation.StringMatch(regexp.MustCompile("^[\u4e00-\u9fa5-_A-Za-z0-9]+$"),
						"The name can only consist of letters, digits, underscores (_),"+
							" hyphens (-) and chinese characters."),
				),
			},
			"alarm_description": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringLenBetween(1, 256),
			},

			"metric": {
				Type:     schema.TypeList,
				Required: true,
				ForceNew: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"namespace": {
							Type:     schema.TypeString,
							Required: true,
							ForceNew: true,
						},

						"metric_name": {
							Type:     schema.TypeString,
							Required: true,
							ForceNew: true,
						},

						"dimensions": {
							Type:     schema.TypeList,
							Optional: true,
							ForceNew: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"name": {
										Type:     schema.TypeString,
										Required: true,
										ForceNew: true,
									},

									"value": {
										Type:     schema.TypeString,
										Required: true,
										ForceNew: true,
									},
								},
							},
						},
					},
				},
			},

			"condition": {
				Type:     schema.TypeList,
				Required: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"period": {
							Type:         schema.TypeInt,
							Required:     true,
							ValidateFunc: validation.IntInSlice([]int{0, 1, 300, 1200, 3600, 14400, 86400}),
						},

						"filter": {
							Type:     schema.TypeString,
							Required: true,
							ValidateFunc: validation.StringInSlice([]string{
								"max", "min", "average", "sum", "variance",
							}, false),
						},

						"comparison_operator": {
							Type:     schema.TypeString,
							Required: true,
							ValidateFunc: validation.StringInSlice([]string{
								">=", ">", "<=", "<", "=",
							}, false),
						},

						"value": {
							Type:     schema.TypeInt,
							Required: true,
						},

						"count": {
							Type:         schema.TypeInt,
							Required:     true,
							ValidateFunc: validation.IntBetween(1, 5),
						},

						"unit": {
							Type:         schema.TypeString,
							Optional:     true,
							ForceNew:     true,
							ValidateFunc: validation.StringLenBetween(1, 32),
						},
						"suppress_duration": {
							Type:     schema.TypeInt,
							Optional: true,
							ForceNew: true,
							ValidateFunc: validation.IntInSlice([]int{
								0, 300, 600, 900, 1800, 3600, 10800, 21600, 43200, 86400,
							}),
						},
					},
				},
			},

			"alarm_actions": &cesAlarmActions,
			"ok_actions":    &cesAlarmActions,

			"alarm_enabled": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
			},

			"alarm_level": {
				Type:         schema.TypeInt,
				Optional:     true,
				Default:      2,
				ValidateFunc: validation.IntBetween(1, 4),
			},

			"alarm_type": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.StringInSlice([]string{"EVENT.SYS", "EVENT.CUSTOM", "MULTI_INSTANCE"}, false),
			},

			"alarm_action_enabled": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
			},

			"enterprise_project_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Computed: true,
			},

			"alarm_state": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"update_time": {
				Type:     schema.TypeInt,
				Computed: true,
			},

			// deprecated
			"insufficientdata_actions": {
				Type:       schema.TypeList,
				Optional:   true,
				Deprecated: "insufficientdata_actions is deprecated",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"type": {
							Type:     schema.TypeString,
							Required: true,
						},
						"notification_list": {
							Type:     schema.TypeList,
							Required: true,
							MaxItems: 5,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
					},
				},
			},
		},
	}
}

func buildMetricOpts(d *schema.ResourceData) alarmrule.MetricOpts {
	var metricOpts alarmrule.MetricOpts
	metricRaw := d.Get("metric").([]interface{})

	if len(metricRaw) == 1 {
		metric := metricRaw[0].(map[string]interface{})
		dimensionsRaw := metric["dimensions"].([]interface{})
		dimensionsOpts := make([]alarmrule.DimensionOpts, len(dimensionsRaw))
		for i, dimensionRaw := range dimensionsRaw {
			dimension := dimensionRaw.(map[string]interface{})
			dimensionsOpts[i] = alarmrule.DimensionOpts{
				Name:  dimension["name"].(string),
				Value: dimension["value"].(string),
			}
		}

		metricOpts.Namespace = metric["namespace"].(string)
		metricOpts.MetricName = metric["metric_name"].(string)
		metricOpts.Dimensions = dimensionsOpts
	}

	return metricOpts
}

func buildAlarmAction(d *schema.ResourceData, name string) []alarmrule.ActionOpts {
	if v, ok := d.GetOk(name); ok {
		actionOptsRaw := v.([]interface{})
		actionOpts := make([]alarmrule.ActionOpts, len(actionOptsRaw))
		for i, actionOptRaw := range actionOptsRaw {
			actionOpt := actionOptRaw.(map[string]interface{})

			notificationListRaw := actionOpt["notification_list"].([]interface{})
			notificationList := make([]string, len(notificationListRaw))
			for j, notification := range notificationListRaw {
				notificationList[j] = notification.(string)
			}

			actionOpts[i] = alarmrule.ActionOpts{
				Type:             actionOpt["type"].(string),
				NotificationList: notificationList,
			}
		}
		return actionOpts
	}
	return nil
}

func buildAlarmCondition(d *schema.ResourceData) alarmrule.ConditionOpts {
	var opts alarmrule.ConditionOpts

	rawCondition := d.Get("condition").([]interface{})
	if len(rawCondition) == 1 {
		condition := rawCondition[0].(map[string]interface{})

		opts.Period = condition["period"].(int)
		opts.Filter = condition["filter"].(string)
		opts.ComparisonOperator = condition["comparison_operator"].(string)
		opts.Value = condition["value"].(int)
		opts.Unit = condition["unit"].(string)
		opts.Count = condition["count"].(int)
		opts.SuppressDuration = condition["suppress_duration"].(int)
	}

	return opts
}

func resourceAlarmRuleCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*config.Config)
	client, err := config.CesV1Client(config.GetRegion(d))
	if err != nil {
		return fmtp.DiagErrorf("Error creating Cloud Eye Service client: %s", err)
	}

	createOpts := alarmrule.CreateOpts{
		AlarmName:               d.Get("alarm_name").(string),
		AlarmDescription:        d.Get("alarm_description").(string),
		AlarmLevel:              d.Get("alarm_level").(int),
		AlarmType:               d.Get("alarm_type").(string),
		Metric:                  buildMetricOpts(d),
		Condition:               buildAlarmCondition(d),
		AlarmActions:            buildAlarmAction(d, "alarm_actions"),
		OkActions:               buildAlarmAction(d, "ok_actions"),
		InsufficientdataActions: buildAlarmAction(d, "insufficientdata_actions"),
		AlarmEnabled:            d.Get("alarm_enabled").(bool),
		AlarmActionEnabled:      d.Get("alarm_action_enabled").(bool),
		EnterpriseProjectID:     config.GetEnterpriseProjectID(d),
	}
	logp.Printf("[DEBUG] Create %s Options: %#v", nameCESAR, createOpts)

	r, err := alarmrule.Create(client, createOpts).Extract()
	if err != nil {
		return fmtp.DiagErrorf("Error creating %s: %s", nameCESAR, err)
	}
	logp.Printf("[DEBUG] Create %s: %#v", nameCESAR, *r)

	d.SetId(r.AlarmID)

	return resourceAlarmRuleRead(ctx, d, meta)
}

func resourceAlarmRuleRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*config.Config)
	client, err := config.CesV1Client(config.GetRegion(d))
	if err != nil {
		return fmtp.DiagErrorf("Error creating Cloud Eye Service client: %s", err)
	}

	r, err := alarmrule.Get(client, d.Id()).Extract()
	if err != nil {
		return common.CheckDeletedDiag(d, err, "Error retrieving HuaweiCloud CES alarm rule")
	}
	logp.Printf("[DEBUG] Retrieved %s %s: %#v", nameCESAR, d.Id(), r)

	m, err := utils.ConvertStructToMap(r, map[string]string{"notificationList": "notification_list"})
	if err != nil {
		return diag.FromErr(err)
	}

	alarmMetric := make([]interface{}, 1)
	alarmMetric[0] = m["metric"]
	alarmCondition := make([]interface{}, 1)
	alarmCondition[0] = m["condition"]

	mErr := multierror.Append(nil,
		d.Set("alarm_name", m["alarm_name"]),
		d.Set("alarm_description", m["alarm_description"]),
		d.Set("alarm_level", m["alarm_level"]),
		d.Set("alarm_type", m["alarm_type"]),
		d.Set("metric", alarmMetric),
		d.Set("condition", alarmCondition),
		d.Set("alarm_actions", m["alarm_actions"]),
		d.Set("ok_actions", m["ok_actions"]),
		d.Set("alarm_enabled", m["alarm_enabled"]),
		d.Set("alarm_action_enabled", m["alarm_action_enabled"]),
		d.Set("alarm_state", m["alarm_state"]),
		d.Set("update_time", m["update_time"]),
		d.Set("enterprise_project_id", m["enterprise_project_id"]),
	)
	if mErr.ErrorOrNil() != nil {
		return diag.FromErr(err)
	}

	return nil
}

func resourceAlarmRuleUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*config.Config)
	client, err := config.CesV1Client(config.GetRegion(d))
	if err != nil {
		return fmtp.DiagErrorf("Error creating Cloud Eye Service client: %s", err)
	}

	arId := d.Id()

	if d.HasChange("alarm_enabled") {
		enabled := d.Get("alarm_enabled").(bool)
		enableOpts := alarmrule.EnableOpts{
			AlarmEnabled: enabled,
		}
		logp.Printf("[DEBUG] Updating %s %s to %#v", nameCESAR, arId, enabled)

		timeout := d.Timeout(schema.TimeoutUpdate)
		//lintignore:R006
		err = resource.RetryContext(ctx, timeout, func() *resource.RetryError {
			err := alarmrule.Enable(client, arId, enableOpts).ExtractErr()
			if err != nil {
				return common.CheckForRetryableError(err)
			}
			return nil
		})
		if err != nil {
			return fmtp.DiagErrorf("Error updating %s %s: %s", nameCESAR, arId, err)
		}
	}

	updateOpts := alarmrule.UpdateOpts{}
	changed := false
	if d.HasChanges("alarm_name", "alarm_description", "alarm_level", "alarm_action_enabled") {
		description := d.Get("alarm_description").(string)
		actionEnabled := d.Get("alarm_action_enabled").(bool)

		updateOpts.Name = d.Get("alarm_name").(string)
		updateOpts.AlarmLevel = d.Get("alarm_level").(int)
		updateOpts.Description = &description
		updateOpts.ActionEnabled = &actionEnabled
		changed = true
	}

	if d.HasChange("condition") {
		condition := buildAlarmCondition(d)
		// unit field is not supported in Update
		condition.Unit = ""
		updateOpts.Condition = &condition
		changed = true
	}

	if changed {
		logp.Printf("[DEBUG] Updating %s %s opts: %#v", nameCESAR, arId, updateOpts)
		err := alarmrule.Update(client, arId, updateOpts).ExtractErr()
		if err != nil {
			return fmtp.DiagErrorf("Error updating %s %s: %s", nameCESAR, arId, err)
		}
	}

	return resourceAlarmRuleRead(ctx, d, meta)
}

func resourceAlarmRuleDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*config.Config)
	client, err := config.CesV1Client(config.GetRegion(d))
	if err != nil {
		return fmtp.DiagErrorf("Error creating Cloud Eye Service client: %s", err)
	}

	arId := d.Id()
	logp.Printf("[DEBUG] Deleting %s %s", nameCESAR, arId)

	timeout := d.Timeout(schema.TimeoutDelete)
	//lintignore:R006
	err = resource.RetryContext(ctx, timeout, func() *resource.RetryError {
		err := alarmrule.Delete(client, arId).ExtractErr()
		if err != nil {
			return common.CheckForRetryableError(err)
		}
		return nil
	})
	if err != nil {
		if utils.IsResourceNotFound(err) {
			logp.Printf("[INFO] deleting an unavailable %s: %s", nameCESAR, arId)
			return nil
		}
		return fmtp.DiagErrorf("Error deleting %s %s: %s", nameCESAR, arId, err)
	}

	return nil
}
