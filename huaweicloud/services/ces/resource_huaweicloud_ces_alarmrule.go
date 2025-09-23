package ces

import (
	"bytes"
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"

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

// @API CES POST /v2/{project_id}/alarms
// @API CES GET /v2/{project_id}/alarms
// @API CES GET /v2/{project_id}/alarms/{id}/resources
// @API CES POST /v2/{project_id}/alarms/{id}/resources/batch-delete
// @API CES POST /v2/{project_id}/alarms/{id}/resources/batch-create
// @API CES PUT /v2/{project_id}/alarms/{id}/policies
// @API CES POST /v2/{project_id}/alarms/action
// @API CES POST /v2/{project_id}/alarms/{id}/resources/{operation}
// @API CES POST /v2/{project_id}/alarms/batch-delete
// @API CES GET /V1.0/{project_id}/alarms/{id}
// @API CES PUT /V1.0/{project_id}/alarms/{id}
// @API EPS POST /v1.0/enterprise-projects/{enterprise_project_id}/resources-migrate
// @API CES PUT /v2/{project_id}/alarms/{alarm_id}/notifications

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
			},

			"alarm_description": {
				Type:     schema.TypeString,
				Optional: true,
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
							Type:          schema.TypeSet,
							Optional:      true,
							Computed:      true,
							Set:           resourceDimensionsHash,
							ConflictsWith: []string{"resources"},
							Description:   "schema: Deprecated",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"name": {
										Type:     schema.TypeString,
										Required: true,
									},

									"value": {
										Type:     schema.TypeString,
										Optional: true,
									},
								},
							},
						},
					},
				},
			},

			"resource_group_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: "schema: Internal",
			},

			"resources": {
				Type:     schema.TypeSet,
				Optional: true,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"dimensions": {
							Type:     schema.TypeList,
							Optional: true,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"name": {
										Type:     schema.TypeString,
										Required: true,
									},

									"value": {
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

			"alarm_template_id": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ExactlyOneOf: []string{"alarm_template_id", "condition"},
			},

			"condition": {
				Type:         schema.TypeSet,
				Optional:     true,
				Computed:     true,
				Set:          resourceConditionHash,
				ExactlyOneOf: []string{"alarm_template_id", "condition"},
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"period": {
							Type:     schema.TypeInt,
							Required: true,
						},

						"filter": {
							Type:     schema.TypeString,
							Required: true,
						},

						"comparison_operator": {
							Type:     schema.TypeString,
							Required: true,
						},

						"value": {
							Type:     schema.TypeFloat,
							Required: true,
						},

						"count": {
							Type:     schema.TypeInt,
							Required: true,
						},

						"unit": {
							Type:     schema.TypeString,
							Optional: true,
						},

						"suppress_duration": {
							Type:     schema.TypeInt,
							Optional: true,
						},

						"metric_name": {
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
							Description: "schema: Required",
						},

						"alarm_level": {
							Type:     schema.TypeInt,
							Optional: true,
							Computed: true,
						},
					},
				},
			},

			"alarm_actions": &cesAlarmActions,
			"ok_actions":    &cesAlarmActions,

			"notification_begin_time": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"notification_end_time": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"effective_timezone": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"alarm_enabled": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
			},

			"alarm_level": {
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
				Description: "schema: Deprecated",
			},

			"alarm_type": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "MULTI_INSTANCE",
			},

			"alarm_action_enabled": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
			},

			"enterprise_project_id": {
				Type:     schema.TypeString,
				Optional: true,
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

func buildCreateAndUpdateAlarmActionBodyParams(rawParams interface{}, notificationListKey string) []map[string]interface{} {
	if rawArray, ok := rawParams.([]interface{}); ok {
		if len(rawArray) == 0 {
			return nil
		}

		params := make([]map[string]interface{}, len(rawArray))
		for i, v := range rawArray {
			raw := v.(map[string]interface{})
			params[i] = map[string]interface{}{
				"type":              raw["type"],
				notificationListKey: raw["notification_list"],
			}
		}
		return params
	}

	return nil
}

func buildDimensionsOpts(dimensionsRaw []interface{}) [][]map[string]interface{} {
	if len(dimensionsRaw) < 1 {
		return make([][]map[string]interface{}, 0)
	}
	resources := make([][]map[string]interface{}, len(dimensionsRaw))
	for i, dimensionRaw := range dimensionsRaw {
		dimension := dimensionRaw.(map[string]interface{})
		content := utils.RemoveNil(map[string]interface{}{
			"name":  dimension["name"],
			"value": utils.ValueIgnoreEmpty(dimension["value"]),
		})
		resources[i] = []map[string]interface{}{content}
	}

	return resources
}

func buildDimensionsOptsV2(resourcesRaw []interface{}) [][]map[string]interface{} {
	if len(resourcesRaw) < 1 {
		return make([][]map[string]interface{}, 0)
	}
	resources := make([][]map[string]interface{}, len(resourcesRaw))
	for i, resourceRaw := range resourcesRaw {
		resource := resourceRaw.(map[string]interface{})
		dimensionRaw := resource["dimensions"].([]interface{})
		res := make([]map[string]interface{}, len(dimensionRaw))
		for j, dimension := range dimensionRaw {
			dim := dimension.(map[string]interface{})
			res[j] = utils.RemoveNil(map[string]interface{}{
				"name":  dim["name"],
				"value": utils.ValueIgnoreEmpty(dim["value"]),
			})
		}
		resources[i] = res
	}

	return resources
}

func buildResourcesOpts(d *schema.ResourceData) [][]map[string]interface{} {
	metricRaw := d.Get("metric").([]interface{})
	if len(metricRaw) != 1 {
		return nil
	}

	metric := metricRaw[0].(map[string]interface{})
	dimensionsRaw := metric["dimensions"].(*schema.Set).List()

	var resources [][]map[string]interface{}
	if v, ok := d.GetOk("resources"); ok {
		resources = buildDimensionsOptsV2(v.(*schema.Set).List())
	} else {
		resources = buildDimensionsOpts(dimensionsRaw)
	}

	return resources
}

func buildPoliciesOpts(d *schema.ResourceData, globalMetricName string) []map[string]interface{} {
	rawCondition := d.Get("condition").(*schema.Set).List()

	if len(rawCondition) < 1 {
		return nil
	}

	globalLevel := 2
	if v, ok := d.GetOk("alarm_level"); ok {
		globalLevel = v.(int)
	}
	policyOpts := make([]map[string]interface{}, len(rawCondition))

	for i, v := range rawCondition {
		condition := v.(map[string]interface{})

		policyOpts[i] = map[string]interface{}{
			"period":              condition["period"],
			"filter":              condition["filter"],
			"comparison_operator": condition["comparison_operator"],
			"value":               utils.ValueIgnoreEmpty(condition["value"]),
			"unit":                utils.ValueIgnoreEmpty(condition["unit"]),
			"count":               condition["count"],
			"suppress_duration":   utils.ValueIgnoreEmpty(condition["suppress_duration"]),
		}

		if condition["metric_name"].(string) != "" {
			policyOpts[i]["metric_name"] = condition["metric_name"].(string)
		} else {
			policyOpts[i]["metric_name"] = globalMetricName
		}

		if condition["alarm_level"].(int) != 0 {
			policyOpts[i]["level"] = condition["alarm_level"].(int)
		} else {
			policyOpts[i]["level"] = globalLevel
		}
	}

	return policyOpts
}

func buildCreateAlarmRuleV2BodyParams(d *schema.ResourceData, enterpriseProjectID string) map[string]interface{} {
	resources := buildResourcesOpts(d)
	namespace := d.Get("metric.0.namespace")
	var metricName string
	if v, ok := d.GetOk("metric.0.metric_name"); ok {
		metricName = v.(string)
	}

	bodyParams := map[string]interface{}{
		"name":                    d.Get("alarm_name"),
		"description":             utils.ValueIgnoreEmpty(d.Get("alarm_description")),
		"namespace":               namespace,
		"resource_group_id":       utils.ValueIgnoreEmpty(d.Get("resource_group_id")),
		"resources":               resources,
		"alarm_template_id":       utils.ValueIgnoreEmpty(d.Get("alarm_template_id")),
		"policies":                buildPoliciesOpts(d, metricName),
		"type":                    d.Get("alarm_type"),
		"alarm_notifications":     buildCreateAndUpdateAlarmActionBodyParams(d.Get("alarm_actions"), "notification_list"),
		"ok_notifications":        buildCreateAndUpdateAlarmActionBodyParams(d.Get("ok_actions"), "notification_list"),
		"notification_begin_time": utils.ValueIgnoreEmpty(d.Get("notification_begin_time")),
		"notification_end_time":   utils.ValueIgnoreEmpty(d.Get("notification_end_time")),
		"effective_timezone":      utils.ValueIgnoreEmpty(d.Get("effective_timezone")),
		"notification_enabled":    d.Get("alarm_action_enabled"),
		"enabled":                 d.Get("alarm_enabled"),
		"enterprise_project_id":   utils.ValueIgnoreEmpty(enterpriseProjectID),
	}

	return bodyParams
}

func resourceAlarmRuleCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conf := meta.(*config.Config)
	clientV2, err := conf.CesV2Client(conf.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating Cloud Eye Service v2 client: %s", err)
	}

	createHttpUrl := "v2/{project_id}/alarms"
	createPath := clientV2.Endpoint + createHttpUrl
	createPath = strings.ReplaceAll(createPath, "{project_id}", clientV2.ProjectID)
	createOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		JSONBody:         utils.RemoveNil(buildCreateAlarmRuleV2BodyParams(d, conf.GetEnterpriseProjectID(d))),
	}

	createResp, err := clientV2.Request("POST", createPath, &createOpt)
	if err != nil {
		return diag.Errorf("error creating %s: %s", nameCESAR, err)
	}

	createRespBody, err := utils.FlattenResponse(createResp)
	if err != nil {
		return diag.Errorf("error flattening creating %s response: %s", nameCESAR, err)
	}

	id := utils.PathSearch("alarm_id", createRespBody, "").(string)
	if id == "" {
		return diag.Errorf("error creating CES %s: can not found %s id in return", nameCESAR, nameCESAR)
	}

	d.SetId(id)

	return resourceAlarmRuleRead(ctx, d, meta)
}

func flattenAlarmRuleAlarmAction(rawParams interface{}, notificationListKey string) []interface{} {
	if paramsList, ok := rawParams.([]interface{}); ok {
		if len(paramsList) == 0 {
			return nil
		}
		rst := make([]interface{}, 0, len(paramsList))
		for _, params := range paramsList {
			raw := params.(map[string]interface{})
			m := map[string]interface{}{
				"type":              utils.PathSearch("type", raw, nil),
				"notification_list": utils.PathSearch(notificationListKey, raw, nil),
			}
			rst = append(rst, m)
		}

		return rst
	}

	return nil
}

func resourceAlarmRuleRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conf := meta.(*config.Config)
	clientV1, err := conf.CesV1Client(conf.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating Cloud Eye Service v1 client: %s", err)
	}
	clientV2, err := conf.CesV2Client(conf.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating Cloud Eye Service v2 client: %s", err)
	}

	rV1, err := getAlarmV1(clientV1, d.Id())
	if err != nil {
		return common.CheckDeletedDiag(d, err, "error retrieving CES alarm rule")
	}

	mErr := multierror.Append(nil,
		d.Set("alarm_name", utils.PathSearch("alarm_name", rV1, nil)),
		d.Set("alarm_description", utils.PathSearch("alarm_description", rV1, nil)),
		d.Set("alarm_type", utils.PathSearch("alarm_type", rV1, nil)),
		d.Set("alarm_actions", flattenAlarmRuleAlarmAction(
			utils.PathSearch("alarm_actions", rV1, nil), "notificationList")),
		d.Set("ok_actions", flattenAlarmRuleAlarmAction(
			utils.PathSearch("ok_actions", rV1, nil), "notificationList")),
		d.Set("alarm_enabled", utils.PathSearch("alarm_enabled", rV1, nil)),
		d.Set("alarm_action_enabled", utils.PathSearch("alarm_action_enabled", rV1, nil)),
		d.Set("alarm_state", utils.PathSearch("alarm_state", rV1, nil)),
		d.Set("update_time", utils.PathSearch("update_time", rV1, nil)),
		d.Set("enterprise_project_id", utils.PathSearch("enterprise_project_id", rV1, nil)),
	)

	rV2, err := getAlarmV2(clientV2, d.Id())
	if err != nil {
		return common.CheckDeletedDiag(d, err, "error retrieving CES alarm rule")
	}

	conditions, metricName, alarmLevel := flattenCondition(utils.PathSearch("policies", rV2, nil))
	resourceGroupID := utils.PathSearch("resources.0.resource_group_id", rV2, "")

	// get resources
	resourcesResp, err := getAlarmResourcesV2(clientV2, d.Id())
	if err != nil {
		return common.CheckDeletedDiag(d, err, "error retrieving CES alarm resources")
	}
	dimensions, resourcesToSet := flattenResources(utils.PathSearch("resources", resourcesResp, nil))

	mErr = multierror.Append(mErr,
		d.Set("notification_begin_time", utils.PathSearch("notification_begin_time", rV2, nil)),
		d.Set("notification_end_time", utils.PathSearch("notification_end_time", rV2, nil)),
		d.Set("effective_timezone", utils.PathSearch("effective_timezone", rV2, nil)),
		d.Set("alarm_template_id", utils.PathSearch("alarm_template_id", rV2, nil)),
		d.Set("condition", conditions),
		d.Set("metric",
			flattenMetric(dimensions, metricName, utils.PathSearch("namespace", rV2, "").(string))),
		d.Set("alarm_level", alarmLevel),
		d.Set("resource_group_id", resourceGroupID),
		d.Set("resources", resourcesToSet),
	)

	if mErr.ErrorOrNil() != nil {
		return diag.FromErr(err)
	}

	return nil
}

func getAlarmV1(client *golangsdk.ServiceClient, alarmId string) (interface{}, error) {
	getHttpUrl := "V1.0/{project_id}/alarms/{id}"
	getPath := client.Endpoint + getHttpUrl
	getPath = strings.ReplaceAll(getPath, "{project_id}", client.ProjectID)
	getPath = strings.ReplaceAll(getPath, "{id}", alarmId)
	getOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json",
		},
	}

	getResp, err := client.Request("GET", getPath, &getOpt)
	if err != nil {
		return nil, err
	}
	getRespBody, err := utils.FlattenResponse(getResp)
	if err != nil {
		return nil, fmt.Errorf("error flattening %s: %s", nameCESAR, err)
	}

	alarm := utils.PathSearch("metric_alarms|[0]", getRespBody, nil)

	return alarm, nil
}

func getAlarmV2(client *golangsdk.ServiceClient, alarmId string) (interface{}, error) {
	getHttpUrl := "v2/{project_id}/alarms?alarm_id={id}"
	getPath := client.Endpoint + getHttpUrl
	getPath = strings.ReplaceAll(getPath, "{project_id}", client.ProjectID)
	getPath = strings.ReplaceAll(getPath, "{id}", alarmId)
	getOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json",
		},
	}

	getResp, err := client.Request("GET", getPath, &getOpt)
	if err != nil {
		return nil, err
	}
	getRespBody, err := utils.FlattenResponse(getResp)
	if err != nil {
		return nil, fmt.Errorf("error flattening %s: %s", nameCESAR, err)
	}

	alarm := utils.PathSearch("alarms|[0]", getRespBody, nil)
	if alarm == nil {
		return nil, golangsdk.ErrDefault404{}
	}

	return alarm, nil
}

func getAlarmResourcesV2(client *golangsdk.ServiceClient, alarmId string) (interface{}, error) {
	getHttpUrl := "v2/{project_id}/alarms/{id}/resources"
	getPath := client.Endpoint + getHttpUrl
	getPath = strings.ReplaceAll(getPath, "{project_id}", client.ProjectID)
	getPath = strings.ReplaceAll(getPath, "{id}", alarmId)
	getOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json",
		},
	}

	getResp, err := client.Request("GET", getPath, &getOpt)
	if err != nil {
		return nil, err
	}
	getRespBody, err := utils.FlattenResponse(getResp)
	if err != nil {
		return nil, fmt.Errorf("error flattening %s: %s", nameCESAR, err)
	}

	return getRespBody, nil
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

func flattenCondition(rawParams interface{}) ([]map[string]interface{}, string, int) {
	if paramsList, ok := rawParams.([]interface{}); ok {
		if len(paramsList) < 1 {
			return nil, "", 0
		}
		conditions := make([]map[string]interface{}, len(paramsList))
		for i, params := range paramsList {
			raw := params.(map[string]interface{})
			conditions[i] = map[string]interface{}{
				"metric_name":         utils.PathSearch("metric_name", raw, ""),
				"period":              utils.PathSearch("period", raw, nil),
				"filter":              utils.PathSearch("filter", raw, nil),
				"comparison_operator": utils.PathSearch("comparison_operator", raw, nil),
				"value":               utils.PathSearch("value", raw, nil),
				"count":               utils.PathSearch("count", raw, nil),
				"unit":                utils.PathSearch("unit", raw, nil),
				"suppress_duration":   utils.PathSearch("suppress_duration", raw, nil),
				"alarm_level":         int(utils.PathSearch("level", raw, float64(0)).(float64)),
			}
		}
		return conditions, conditions[0]["metric_name"].(string), conditions[0]["alarm_level"].(int)
	}

	return nil, "", 0
}

func flattenResources(rawParams interface{}) (dimensions []map[string]interface{}, resourcesToSet []map[string]interface{}) {
	if paramsList, ok := rawParams.([]interface{}); ok && len(paramsList) > 0 {
		dimensions = make([]map[string]interface{}, 0, len(paramsList))
		resourcesToSet = make([]map[string]interface{}, len(paramsList))
		for i, params := range paramsList {
			if raws, ok := params.([]interface{}); ok {
				resource := make([]map[string]interface{}, len(raws))
				for j, v := range raws {
					raw := v.(map[string]interface{})
					if j == 0 {
						dimensions = append(dimensions, map[string]interface{}{
							"name":  raw["name"],
							"value": raw["value"],
						})
					}
					resource[j] = map[string]interface{}{
						"name":  raw["name"],
						"value": raw["value"],
					}
				}
				resourcesToSet[i] = map[string]interface{}{
					"dimensions": resource,
				}
			}
		}

		return dimensions, resourcesToSet
	}

	dimensions = make([]map[string]interface{}, 0)
	resourcesToSet = make([]map[string]interface{}, 0)
	return dimensions, resourcesToSet
}

func buildUpdatePoliciesOptsWithAlarmLevel(d *schema.ResourceData, level int, metricName string) []map[string]interface{} {
	rawCondition := d.Get("condition").(*schema.Set).List()

	if len(rawCondition) < 1 {
		return nil
	}

	policyOpts := make([]map[string]interface{}, len(rawCondition))

	for i, v := range rawCondition {
		condition := v.(map[string]interface{})

		policyOpts[i] = map[string]interface{}{
			"period":              condition["period"],
			"filter":              condition["filter"],
			"comparison_operator": condition["comparison_operator"],
			"value":               utils.ValueIgnoreEmpty(condition["value"]),
			"unit":                utils.ValueIgnoreEmpty(condition["unit"]),
			"count":               condition["count"],
			"suppress_duration":   utils.ValueIgnoreEmpty(condition["suppress_duration"]),
			"level":               level,
		}

		if condition["metric_name"].(string) == "" {
			policyOpts[i]["metric_name"] = metricName
		} else {
			policyOpts[i]["metric_name"] = condition["metric_name"]
		}
	}

	return policyOpts
}

func buildUpdatePoliciesOptsWithMetricName(d *schema.ResourceData, level int, metricName string) []map[string]interface{} {
	rawCondition := d.Get("condition").(*schema.Set).List()

	if len(rawCondition) < 1 {
		return nil
	}

	policyOpts := make([]map[string]interface{}, len(rawCondition))

	for i, v := range rawCondition {
		condition := v.(map[string]interface{})

		policyOpts[i] = map[string]interface{}{
			"period":              condition["period"],
			"filter":              condition["filter"],
			"comparison_operator": condition["comparison_operator"],
			"value":               utils.ValueIgnoreEmpty(condition["value"]),
			"unit":                utils.ValueIgnoreEmpty(condition["unit"]),
			"count":               condition["count"],
			"suppress_duration":   utils.ValueIgnoreEmpty(condition["suppress_duration"]),
			"metric_name":         metricName,
		}

		if condition["alarm_level"].(int) == 0 {
			policyOpts[i]["level"] = level
		} else {
			policyOpts[i]["level"] = condition["alarm_level"]
		}
	}

	return policyOpts
}

func resourceAlarmRuleUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	clientV1, err := cfg.CesV1Client(region)
	if err != nil {
		return diag.Errorf("error creating Cloud Eye Service v1 client: %s", err)
	}
	clientV2, err := cfg.CesV2Client(region)
	if err != nil {
		return diag.Errorf("error creating Cloud Eye Service v2 client: %s", err)
	}

	arId := d.Id()

	if d.HasChanges("alarm_name", "alarm_description", "alarm_action_enabled", "alarm_actions", "ok_actions") {
		updateHttpUrl := "V1.0/{project_id}/alarms/{id}"
		updatePath := clientV1.Endpoint + updateHttpUrl
		updatePath = strings.ReplaceAll(updatePath, "{project_id}", clientV1.ProjectID)
		updatePath = strings.ReplaceAll(updatePath, "{id}", d.Id())
		updateOpt := golangsdk.RequestOpts{
			KeepResponseBody: true,
			JSONBody:         utils.RemoveNil(buildUpdateAlarmV1BodyParams(d)),
			OkCodes:          []int{204},
		}

		_, err = clientV1.Request("PUT", updatePath, &updateOpt)
		if err != nil {
			return diag.Errorf("error updating %s %s: %s", nameCESAR, arId, err)
		}
	}

	if d.HasChange("metric.0.dimensions") {
		oldDimensions, newDimensions := d.GetChange("metric.0.dimensions")

		if len(oldDimensions.(*schema.Set).List()) > 0 {
			err := batchCreateAndUpdateAlarmRuleResources(clientV2, arId, "batch-delete",
				buildDimensionsOpts(oldDimensions.(*schema.Set).List()))
			if err != nil {
				return diag.Errorf("error deleting old dimensions of %s %s: %s", nameCESAR, arId, err)
			}
		}

		if len(newDimensions.(*schema.Set).List()) > 0 {
			err := batchCreateAndUpdateAlarmRuleResources(clientV2, arId, "batch-create",
				buildDimensionsOpts(newDimensions.(*schema.Set).List()))
			if err != nil {
				return diag.Errorf("error creating new dimensions of %s %s: %s", nameCESAR, arId, err)
			}
		}
	}

	if d.HasChange("resources") {
		oldDimensions, newDimensions := d.GetChange("resources")

		if len(oldDimensions.(*schema.Set).List()) > 0 {
			err := batchCreateAndUpdateAlarmRuleResources(clientV2, arId, "batch-delete",
				buildDimensionsOptsV2(oldDimensions.(*schema.Set).List()))
			if err != nil {
				return diag.Errorf("error deleting old resources of %s %s: %s", nameCESAR, arId, err)
			}
		}

		if len(newDimensions.(*schema.Set).List()) > 0 {
			err := batchCreateAndUpdateAlarmRuleResources(clientV2, arId, "batch-create",
				buildDimensionsOptsV2(newDimensions.(*schema.Set).List()))
			if err != nil {
				return diag.Errorf("error creating new resources of %s %s: %s", nameCESAR, arId, err)
			}
		}
	}

	level := 2
	if v, ok := d.GetOk("alarm_level"); ok {
		level = v.(int)
	}

	var metricName string
	if v, ok := d.GetOk("metric.0.metric_name"); ok {
		metricName = v.(string)
	}

	// update condition if alarm_level changed
	if d.HasChange("alarm_level") {
		err := updateAlarmRulePolicies(clientV2, arId, buildUpdatePoliciesOptsWithAlarmLevel(d, level, metricName))
		if err != nil {
			return diag.Errorf("error updating condition of %s %s: %s", nameCESAR, arId, err)
		}
	}

	// update condition if metric.0.metric_name changed
	if d.HasChange("metric.0.metric_name") {
		err := updateAlarmRulePolicies(clientV2, arId, buildUpdatePoliciesOptsWithMetricName(d, level, metricName))
		if err != nil {
			return diag.Errorf("error updating condition of %s %s: %s", nameCESAR, arId, err)
		}
	}

	// update condition
	if d.HasChange("condition") {
		err := updateAlarmRulePolicies(clientV2, arId, buildPoliciesOpts(d, metricName))
		if err != nil {
			return diag.Errorf("error updating condition of %s %s: %s", nameCESAR, arId, err)
		}
	}

	if d.HasChange("alarm_enabled") {
		enabled := d.Get("alarm_enabled").(bool)
		updateHttpUrl := "v2/{project_id}/alarms/action"
		updatePath := clientV2.Endpoint + updateHttpUrl
		updatePath = strings.ReplaceAll(updatePath, "{project_id}", clientV2.ProjectID)
		updateOpt := golangsdk.RequestOpts{
			KeepResponseBody: true,
			JSONBody: map[string]interface{}{
				"alarm_ids":     []string{arId},
				"alarm_enabled": enabled,
			},
		}

		_, err := clientV2.Request("POST", updatePath, &updateOpt)
		if err != nil {
			return diag.Errorf("error updating %s %s: %s", nameCESAR, arId, err)
		}
	}

	if d.HasChange("enterprise_project_id") {
		migrateOpts := config.MigrateResourceOpts{
			ResourceId:   arId,
			ResourceType: "CES-alarm",
			RegionId:     region,
			ProjectId:    clientV1.ProjectID,
		}
		if err := cfg.MigrateEnterpriseProject(ctx, d, migrateOpts); err != nil {
			return diag.FromErr(err)
		}
	}

	if d.HasChanges("notification_begin_time", "notification_end_time", "effective_timezone") {
		updateHttpUrl := "v2/{project_id}/alarms/{id}/notifications"
		updatePath := clientV2.Endpoint + updateHttpUrl
		updatePath = strings.ReplaceAll(updatePath, "{project_id}", clientV2.ProjectID)
		updatePath = strings.ReplaceAll(updatePath, "{id}", arId)
		updateOpt := golangsdk.RequestOpts{
			KeepResponseBody: true,
			JSONBody:         utils.RemoveNil(buildUpdateAlarmNotificationsV2BodyParams(d)),
		}

		_, err := clientV2.Request("PUT", updatePath, &updateOpt)
		if err != nil {
			return diag.Errorf("error updating %s %s: %s", nameCESAR, arId, err)
		}
	}

	return resourceAlarmRuleRead(ctx, d, meta)
}

func updateAlarmRulePolicies(clientV2 *golangsdk.ServiceClient, alarmID string, policiesBody []map[string]interface{}) error {
	policyUpdateHttpUrl := "v2/{project_id}/alarms/{id}/policies"
	policyUpdatePath := clientV2.Endpoint + policyUpdateHttpUrl
	policyUpdatePath = strings.ReplaceAll(policyUpdatePath, "{project_id}", clientV2.ProjectID)
	policyUpdatePath = strings.ReplaceAll(policyUpdatePath, "{id}", alarmID)
	updateOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		JSONBody: utils.RemoveNil(map[string]interface{}{
			"policies": policiesBody,
		}),
	}

	_, err := clientV2.Request("PUT", policyUpdatePath, &updateOpt)
	return err
}

func batchCreateAndUpdateAlarmRuleResources(clientV2 *golangsdk.ServiceClient, alarmID string, method string,
	resourcesBody [][]map[string]interface{}) error {
	batchDeleteHttpUrl := "v2/{project_id}/alarms/{id}/resources/{batch-create-delete}"
	batchDeletePath := clientV2.Endpoint + batchDeleteHttpUrl
	batchDeletePath = strings.ReplaceAll(batchDeletePath, "{project_id}", clientV2.ProjectID)
	batchDeletePath = strings.ReplaceAll(batchDeletePath, "{id}", alarmID)
	batchDeletePath = strings.ReplaceAll(batchDeletePath, "{batch-create-delete}", method)
	createOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		JSONBody: map[string]interface{}{
			"resources": resourcesBody,
		},
	}

	_, err := clientV2.Request("POST", batchDeletePath, &createOpt)
	return err
}

func buildUpdateAlarmV1BodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"alarm_name":        utils.ValueIgnoreEmpty(d.Get("alarm_name")),
		"alarm_actions":     buildCreateAndUpdateAlarmActionBodyParams(d.Get("alarm_actions"), "notificationList"),
		"ok_actions":        buildCreateAndUpdateAlarmActionBodyParams(d.Get("ok_actions"), "notificationList"),
		"alarm_description": utils.ValueIgnoreEmpty(d.Get("alarm_description")),
	}

	// add alarm_action_enabled to the updateOpts only when it's changed
	// this can avoid API error
	if d.HasChange("alarm_action_enabled") {
		bodyParams["alarm_action_enabled"] = d.Get("alarm_action_enabled").(bool)
	}

	return bodyParams
}

func buildUpdateAlarmNotificationsV2BodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"notification_enabled":    d.Get("alarm_action_enabled"),
		"alarm_notifications":     buildCreateAndUpdateAlarmActionBodyParams(d.Get("alarm_actions"), "notification_list"),
		"ok_notifications":        buildCreateAndUpdateAlarmActionBodyParams(d.Get("ok_actions"), "notification_list"),
		"notification_begin_time": utils.ValueIgnoreEmpty(d.Get("notification_begin_time")),
		"notification_end_time":   utils.ValueIgnoreEmpty(d.Get("notification_end_time")),
		"effective_timezone":      utils.ValueIgnoreEmpty(d.Get("effective_timezone")),
	}

	return bodyParams
}

func resourceAlarmRuleDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conf := meta.(*config.Config)
	clientV2, err := conf.CesV2Client(conf.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating Cloud Eye v2 Service client: %s", err)
	}

	arId := d.Id()

	deleteHttpUrl := "v2/{project_id}/alarms/batch-delete"
	deletePath := clientV2.Endpoint + deleteHttpUrl
	deletePath = strings.ReplaceAll(deletePath, "{project_id}", clientV2.ProjectID)
	deleteOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		JSONBody: map[string]interface{}{
			"alarm_ids": []string{arId},
		},
	}

	_, err = clientV2.Request("POST", deletePath, &deleteOpt)
	if err != nil {
		return common.CheckDeletedDiag(d, err, fmt.Sprintf("error deleting %s %s: %s", nameCESAR, arId, err))
	}

	return nil
}
