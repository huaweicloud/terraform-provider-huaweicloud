package ces

import (
	"bytes"
	"context"
	"fmt"
	"log"
	"regexp"
	"time"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/chnsz/golangsdk/openstack/cloudeyeservice/v1/alarmrule"
	alarmrulev2 "github.com/chnsz/golangsdk/openstack/cloudeyeservice/v2/alarmrule"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/helper/hashcode"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

const nameCESAR = "CES-AlarmRule"

var cesAlarmActions = schema.Schema{
	Type:     schema.TypeList,
	Optional: true,
	Elem: &schema.Resource{
		Schema: map[string]*schema.Schema{
			"type": {
				Type:     schema.TypeString,
				Required: true,
				ValidateFunc: validation.StringInSlice([]string{
					"notification", "autoscaling",
				}, false),
			},

			"notification_list": {
				Type:     schema.TypeList,
				MaxItems: 5,
				Required: true,
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
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
							Description: "schema: Deprecated",
						},

						"dimensions": {
							Type:     schema.TypeSet,
							Optional: true,
							Set:      resourceDimensionsHash,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"name": {
										Type:     schema.TypeString,
										Required: true,
									},

									"value": {
										Type:     schema.TypeString,
										Required: true,
									},
								},
							},
						},
					},
				},
			},

			"condition": {
				Type:     schema.TypeSet,
				Required: true,
				Set:      resourceConditionHash,
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
							Type:     schema.TypeFloat,
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
							ValidateFunc: validation.StringLenBetween(1, 32),
						},

						"suppress_duration": {
							Type:     schema.TypeInt,
							Optional: true,
							ValidateFunc: validation.IntInSlice([]int{
								0, 300, 600, 900, 1800, 3600, 10800, 21600, 43200, 86400,
							}),
						},

						"metric_name": {
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
							Description: "schema: Required",
						},

						"alarm_level": {
							Type:         schema.TypeInt,
							Optional:     true,
							Computed:     true,
							ValidateFunc: validation.IntBetween(1, 4),
						},
					},
				},
			},

			"alarm_actions": &cesAlarmActions,
			"ok_actions":    &cesAlarmActions,

			"notification_begin_time": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Computed: true,
			},

			"notification_end_time": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Computed: true,
			},

			"alarm_enabled": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
			},

			"alarm_level": {
				Type:         schema.TypeInt,
				Optional:     true,
				Computed:     true,
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

func resourceDimensionsHash(v interface{}) int {
	var buf bytes.Buffer
	m := v.(map[string]interface{})

	if m["name"] != nil {
		buf.WriteString(fmt.Sprintf("%s-", m["name"].(string)))
	}

	if m["value"] != nil {
		buf.WriteString(fmt.Sprintf("%s-", m["value"].(string)))
	}

	return hashcode.String(buf.String())
}

func resourceConditionHash(v interface{}) int {
	var buf bytes.Buffer
	m := v.(map[string]interface{})

	if m["metric_name"] != nil {
		buf.WriteString(fmt.Sprintf("%s-", m["metric_name"].(string)))
	}

	if m["period"] != nil {
		buf.WriteString(fmt.Sprintf("%d-", m["period"].(int)))
	}

	if m["filter"] != nil {
		buf.WriteString(fmt.Sprintf("%s-", m["filter"].(string)))
	}

	if m["comparison_operator"] != nil {
		buf.WriteString(fmt.Sprintf("%s-", m["comparison_operator"].(string)))
	}

	if m["value"] != nil {
		buf.WriteString(fmt.Sprintf("%f-", m["value"].(float64)))
	}

	if m["count"] != nil {
		buf.WriteString(fmt.Sprintf("%d-", m["count"].(int)))
	}

	if m["unit"] != nil {
		buf.WriteString(fmt.Sprintf("%s-", m["unit"].(string)))
	}

	if m["suppress_duration"] != nil {
		buf.WriteString(fmt.Sprintf("%d-", m["suppress_duration"].(int)))
	}

	if m["alarm_level"] != nil {
		buf.WriteString(fmt.Sprintf("%d-", m["alarm_level"].(int)))
	}

	return hashcode.String(buf.String())
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

func buildDimensionsOpts(dimensionsRaw []interface{}) [][]alarmrulev2.DimensionOpts {
	if len(dimensionsRaw) < 1 {
		return [][]alarmrulev2.DimensionOpts{}
	}
	resources := make([][]alarmrulev2.DimensionOpts, len(dimensionsRaw))
	for i, dimensionRaw := range dimensionsRaw {
		dimension := dimensionRaw.(map[string]interface{})
		resources[i] = []alarmrulev2.DimensionOpts{
			{
				Name:  dimension["name"].(string),
				Value: dimension["value"].(string),
			},
		}
	}

	return resources
}

func buildResourcesOpts(d *schema.ResourceData) ([][]alarmrulev2.DimensionOpts, string, string) {
	metricRaw := d.Get("metric").([]interface{})
	if len(metricRaw) != 1 {
		return nil, "", ""
	}

	metric := metricRaw[0].(map[string]interface{})
	dimensionsRaw := metric["dimensions"].(*schema.Set).List()

	return buildDimensionsOpts(dimensionsRaw), metric["namespace"].(string), metric["metric_name"].(string)
}

func buildNotificationsOpts(d *schema.ResourceData, name string) []alarmrulev2.NotificationOpts {
	if v, ok := d.GetOk(name); ok {
		notificationOptsRaw := v.([]interface{})
		notificationOpts := make([]alarmrulev2.NotificationOpts, len(notificationOptsRaw))
		for i, notificationOptRaw := range notificationOptsRaw {
			notificationOpt := notificationOptRaw.(map[string]interface{})

			notificationListRaw := notificationOpt["notification_list"].([]interface{})
			notificationList := make([]string, len(notificationListRaw))
			for j, notification := range notificationListRaw {
				notificationList[j] = notification.(string)
			}

			notificationOpts[i] = alarmrulev2.NotificationOpts{
				Type:             notificationOpt["type"].(string),
				NotificationList: notificationList,
			}
		}
		return notificationOpts
	}
	return nil
}

func buildPoliciesOpts(d *schema.ResourceData, globalMetricName string) []alarmrulev2.PolicyOpts {
	rawCondition := d.Get("condition").(*schema.Set).List()

	if len(rawCondition) < 1 {
		return nil
	}

	globalLevel := 2
	if v, ok := d.GetOk("alarm_level"); ok {
		globalLevel = v.(int)
	}
	policyOpts := make([]alarmrulev2.PolicyOpts, len(rawCondition))

	for i, v := range rawCondition {
		condition := v.(map[string]interface{})

		policyOpts[i] = alarmrulev2.PolicyOpts{
			Period:             condition["period"].(int),
			Filter:             condition["filter"].(string),
			ComparisonOperator: condition["comparison_operator"].(string),
			Value:              condition["value"].(float64),
			Unit:               condition["unit"].(string),
			Count:              condition["count"].(int),
			SuppressDuration:   condition["suppress_duration"].(int),
		}

		if condition["metric_name"].(string) != "" {
			policyOpts[i].MetricName = condition["metric_name"].(string)
		} else {
			policyOpts[i].MetricName = globalMetricName
		}

		if condition["alarm_level"].(int) != 0 {
			policyOpts[i].Level = condition["alarm_level"].(int)
		} else {
			policyOpts[i].Level = globalLevel
		}
	}

	return policyOpts
}

func resourceAlarmRuleCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*config.Config)
	clientV2, err := config.CesV2Client(config.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating Cloud Eye Service v2 client: %s", err)
	}

	resources, namespace, metricName := buildResourcesOpts(d)

	createOpts := alarmrulev2.CreateOpts{
		Name:                  d.Get("alarm_name").(string),
		Description:           d.Get("alarm_description").(string),
		Namespace:             namespace,
		Resources:             resources,
		Policies:              buildPoliciesOpts(d, metricName),
		Type:                  d.Get("alarm_type").(string),
		AlarmNotifications:    buildNotificationsOpts(d, "alarm_actions"),
		OkNotifications:       buildNotificationsOpts(d, "ok_actions"),
		NotificationBeginTime: d.Get("notification_begin_time").(string),
		NotificationEndTime:   d.Get("notification_end_time").(string),
		NotificationEnabled:   d.Get("alarm_action_enabled").(bool),
		Enabled:               d.Get("alarm_enabled").(bool),
		EnterpriseProjectID:   config.GetEnterpriseProjectID(d),
	}

	log.Printf("[DEBUG] Create %s Options: %#v", nameCESAR, createOpts)

	r, err := alarmrulev2.Create(clientV2, createOpts).Extract()
	if err != nil {
		return diag.Errorf("error creating %s: %s", nameCESAR, err)
	}
	log.Printf("[DEBUG] Create %s: %#v", nameCESAR, *r)

	d.SetId(r.AlarmID)

	return resourceAlarmRuleRead(ctx, d, meta)
}

func resourceAlarmRuleRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*config.Config)
	clientV1, err := config.CesV1Client(config.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating Cloud Eye Service v1 client: %s", err)
	}
	clientV2, err := config.CesV2Client(config.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating Cloud Eye Service v2 client: %s", err)
	}

	rV1, err := alarmrule.Get(clientV1, d.Id()).Extract()
	if err != nil {
		return common.CheckDeletedDiag(d, err, "error retrieving CES alarm rule")
	}
	log.Printf("[DEBUG] Retrieved with v1 API %s %s: %#v", nameCESAR, d.Id(), rV1)

	m, err := utils.ConvertStructToMap(rV1, map[string]string{"notificationList": "notification_list"})
	if err != nil {
		return diag.FromErr(err)
	}

	mErr := multierror.Append(nil,
		d.Set("alarm_name", m["alarm_name"]),
		d.Set("alarm_description", m["alarm_description"]),
		d.Set("alarm_type", m["alarm_type"]),
		d.Set("alarm_actions", m["alarm_actions"]),
		d.Set("ok_actions", m["ok_actions"]),
		d.Set("alarm_enabled", m["alarm_enabled"]),
		d.Set("alarm_action_enabled", m["alarm_action_enabled"]),
		d.Set("alarm_state", m["alarm_state"]),
		d.Set("update_time", m["update_time"]),
		d.Set("enterprise_project_id", m["enterprise_project_id"]),
	)

	rV2, err := alarmrulev2.Get(clientV2, d.Id()).Extract()
	if err != nil {
		return common.CheckDeletedDiag(d, err, "error retrieving CES alarm rule")
	}
	log.Printf("[DEBUG] Retrieved with v2 API %s %s: %#v", nameCESAR, d.Id(), rV2)

	conditions, metricName, alarmLevel := flattenCondition(rV2.Policies)

	// get resources
	resources, err := alarmrulev2.GetResources(clientV2, d.Id()).Extract()
	if err != nil {
		return common.CheckDeletedDiag(d, err, "error retrieving CES alarm resources")
	}

	dimensions := make([]map[string]interface{}, 0, len(*resources))
	if len(*resources) > 0 {
		for _, v := range *resources {
			if len(v) > 0 {
				dimensions = append(dimensions, map[string]interface{}{
					"name":  v[0].Name,
					"value": v[0].Value,
				})
			}
		}
	}

	mErr = multierror.Append(mErr,
		d.Set("notification_begin_time", rV2.NotificationBeginTime),
		d.Set("notification_end_time", rV2.NotificationEndTime),
		d.Set("condition", conditions),
		d.Set("metric", flattenMetric(dimensions, metricName, rV2.Namespace)),
		d.Set("alarm_level", alarmLevel),
	)

	if mErr.ErrorOrNil() != nil {
		return diag.FromErr(err)
	}

	return nil
}

func flattenMetric(dimensions []map[string]interface{}, metricName, namespace string) []map[string]interface{} {
	metric := map[string]interface{}{
		"metric_name": metricName,
		"namespace":   namespace,
	}

	if len(dimensions) > 0 {
		metric["dimensions"] = dimensions
	}

	return []map[string]interface{}{metric}
}

func flattenCondition(policies []alarmrulev2.PolicyOpts) ([]map[string]interface{}, string, int) {
	if len(policies) > 0 {
		conditions := make([]map[string]interface{}, len(policies))
		for i, v := range policies {
			conditions[i] = map[string]interface{}{
				"metric_name":         v.MetricName,
				"period":              v.Period,
				"filter":              v.Filter,
				"comparison_operator": v.ComparisonOperator,
				"value":               v.Value,
				"count":               v.Count,
				"unit":                v.Unit,
				"suppress_duration":   v.SuppressDuration,
				"alarm_level":         v.Level,
			}
		}

		return conditions, conditions[0]["metric_name"].(string), conditions[0]["alarm_level"].(int)
	}

	return nil, "", 0
}

func buildUpdatePoliciesOptsWithAlarmLevel(d *schema.ResourceData, level int, metricName string) []alarmrulev2.PolicyOpts {
	rawCondition := d.Get("condition").(*schema.Set).List()

	if len(rawCondition) < 1 {
		return nil
	}

	policyOpts := make([]alarmrulev2.PolicyOpts, len(rawCondition))

	for i, v := range rawCondition {
		condition := v.(map[string]interface{})

		policyOpts[i] = alarmrulev2.PolicyOpts{
			Period:             condition["period"].(int),
			Filter:             condition["filter"].(string),
			ComparisonOperator: condition["comparison_operator"].(string),
			Value:              condition["value"].(float64),
			Unit:               condition["unit"].(string),
			Count:              condition["count"].(int),
			SuppressDuration:   condition["suppress_duration"].(int),
			Level:              level,
		}

		if condition["metric_name"].(string) == "" {
			policyOpts[i].MetricName = metricName
		} else {
			policyOpts[i].MetricName = condition["metric_name"].(string)
		}
	}

	return policyOpts
}

func buildUpdatePoliciesOptsWithMetricName(d *schema.ResourceData, level int, metricName string) []alarmrulev2.PolicyOpts {
	rawCondition := d.Get("condition").(*schema.Set).List()

	if len(rawCondition) < 1 {
		return nil
	}

	policyOpts := make([]alarmrulev2.PolicyOpts, len(rawCondition))

	for i, v := range rawCondition {
		condition := v.(map[string]interface{})

		policyOpts[i] = alarmrulev2.PolicyOpts{
			Period:             condition["period"].(int),
			Filter:             condition["filter"].(string),
			ComparisonOperator: condition["comparison_operator"].(string),
			Value:              condition["value"].(float64),
			Unit:               condition["unit"].(string),
			Count:              condition["count"].(int),
			SuppressDuration:   condition["suppress_duration"].(int),
			MetricName:         metricName,
		}

		if condition["alarm_level"].(int) == 0 {
			policyOpts[i].Level = level
		} else {
			policyOpts[i].Level = condition["alarm_level"].(int)
		}
	}

	return policyOpts
}

func resourceAlarmRuleUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*config.Config)
	clientV1, err := config.CesV1Client(config.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating Cloud Eye Service v1 client: %s", err)
	}
	clientV2, err := config.CesV2Client(config.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating Cloud Eye Service v2 client: %s", err)
	}

	arId := d.Id()

	if d.HasChanges("alarm_name", "alarm_description", "alarm_action_enabled", "alarm_actions", "ok_actions") {
		updateOpts := alarmrule.UpdateOpts{
			Name:         d.Get("alarm_name").(string),
			AlarmActions: buildAlarmAction(d, "alarm_actions"),
			OkActions:    buildAlarmAction(d, "ok_actions"),
		}

		description := d.Get("alarm_description").(string)
		updateOpts.Description = &description

		// add alarm_action_enabled to the updateOpts only when it's changed
		// this can avoid API error
		if d.HasChange("alarm_action_enabled") {
			actionEnabled := d.Get("alarm_action_enabled").(bool)
			updateOpts.ActionEnabled = &actionEnabled
		}

		log.Printf("[DEBUG] Updating %s %s opts: %#v", nameCESAR, arId, updateOpts)
		err := alarmrule.Update(clientV1, arId, updateOpts).ExtractErr()
		if err != nil {
			return diag.Errorf("error updating %s %s: %s", nameCESAR, arId, err)
		}
	}

	if d.HasChange("metric.0.dimensions") {
		oldDimensions, newDimensions := d.GetChange("metric.0.dimensions")

		if len(oldDimensions.(*schema.Set).List()) > 0 {
			updateDimensionsOpts := alarmrulev2.UpdateResourcesOpts{
				Resources: buildDimensionsOpts(oldDimensions.(*schema.Set).List()),
			}

			err := alarmrulev2.BatchResources(clientV2, arId, "batch-delete", updateDimensionsOpts).ExtractErr()
			if err != nil {
				return diag.Errorf("error deleting old dimensions of %s %s: %s", nameCESAR, arId, err)
			}
		}

		if len(oldDimensions.(*schema.Set).List()) > 0 {
			updateDimensionsOpts := alarmrulev2.UpdateResourcesOpts{
				Resources: buildDimensionsOpts(newDimensions.(*schema.Set).List()),
			}
			err := alarmrulev2.BatchResources(clientV2, arId, "batch-create", updateDimensionsOpts).ExtractErr()
			if err != nil {
				return diag.Errorf("error creating new dimensions of %s %s: %s", nameCESAR, arId, err)
			}
		}
	}

	level := 2
	if v, ok := d.GetOk("alarm_level"); ok {
		level = v.(int)
	}

	_, _, metricName := buildResourcesOpts(d)

	// update condition if alarm_level changed
	if d.HasChange("alarm_level") {
		updatePoliciesOpts := alarmrulev2.UpdatePoliciesOpts{
			Policies: buildUpdatePoliciesOptsWithAlarmLevel(d, level, metricName),
		}

		err := alarmrulev2.PoliciesModify(clientV2, arId, updatePoliciesOpts).ExtractErr()
		if err != nil {
			return diag.Errorf("error updating condition of %s %s: %s", nameCESAR, arId, err)
		}
	}

	// update condition if metric.0.metric_name changed
	if d.HasChange("metric.0.metric_name") {
		updatePoliciesOpts := alarmrulev2.UpdatePoliciesOpts{
			Policies: buildUpdatePoliciesOptsWithMetricName(d, level, metricName),
		}

		err := alarmrulev2.PoliciesModify(clientV2, arId, updatePoliciesOpts).ExtractErr()
		if err != nil {
			return diag.Errorf("error updating condition of %s %s: %s", nameCESAR, arId, err)
		}
	}

	// update condition
	if d.HasChange("condition") {
		updatePoliciesOpts := alarmrulev2.UpdatePoliciesOpts{
			Policies: buildPoliciesOpts(d, metricName),
		}

		err := alarmrulev2.PoliciesModify(clientV2, arId, updatePoliciesOpts).ExtractErr()
		if err != nil {
			return diag.Errorf("error updating condition of %s %s: %s", nameCESAR, arId, err)
		}
	}

	if d.HasChange("alarm_enabled") {
		enabled := d.Get("alarm_enabled").(bool)
		actionOpts := alarmrulev2.ActionOpts{
			AlarmIDs:     []string{arId},
			AlarmEnabled: enabled,
		}
		log.Printf("[DEBUG] Updating %s %s to %#v", nameCESAR, arId, enabled)

		err := alarmrulev2.Action(clientV2, arId, actionOpts).ExtractErr()
		if err != nil {
			return diag.Errorf("error updating %s %s: %s", nameCESAR, arId, err)
		}
	}

	return resourceAlarmRuleRead(ctx, d, meta)
}

func resourceAlarmRuleDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*config.Config)
	clientV2, err := config.CesV2Client(config.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating Cloud Eye v2 Service client: %s", err)
	}

	arId := d.Id()
	log.Printf("[DEBUG] Deleting %s %s", nameCESAR, arId)

	deleteOpts := alarmrulev2.DeleteOpts{
		AlarmIDs: []string{arId},
	}
	err = alarmrulev2.Delete(clientV2, deleteOpts).ExtractErr()
	if err != nil {
		return common.CheckDeletedDiag(d, err, fmt.Sprintf("error deleting %s %s: %s", nameCESAR, arId, err))
	}

	return nil
}
