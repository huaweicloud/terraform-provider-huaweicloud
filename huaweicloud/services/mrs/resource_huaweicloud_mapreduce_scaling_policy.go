// ---------------------------------------------------------------
// *** AUTO GENERATED CODE ***
// @Product MRS
// ---------------------------------------------------------------

package mrs

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API MRS POST /v1.1/{project_id}/autoscaling-policy/{cluster_id}
// @API MRS GET /v2/{project_id}/autoscaling-policy/{cluster_id}
func ResourceScalingPolicy() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceScalingPolicyCreateOrUpdate,
		UpdateContext: resourceScalingPolicyCreateOrUpdate,
		ReadContext:   resourceScalingPolicyRead,
		DeleteContext: resourceScalingPolicyDelete,
		Importer: &schema.ResourceImporter{
			StateContext: resourceScalingPolicyImportState,
		},

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"cluster_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: `The MRS cluster ID to which the auto scaling policy applies.`,
			},
			"node_group": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `Name of the node to which an auto scaling rule applies.`,
			},
			"auto_scaling_enable": {
				Type:        schema.TypeBool,
				Required:    true,
				Description: `Whether to enable the auto scaling rule.`,
			},
			"min_capacity": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: `Minimum number of nodes in the node group. Value range: 0 to 500.`,
			},
			"max_capacity": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: `Maximum number of nodes in the node group. Value range: 0 to 500.`,
			},
			"resources_plans": {
				Type:        schema.TypeList,
				Elem:        scalingPolicyResourcesPlanSchema(),
				Optional:    true,
				Description: `The list of resources plans.`,
			},
			"rules": {
				Type:        schema.TypeList,
				Elem:        scalingPolicyRuleSchema(),
				Optional:    true,
				Description: `The list of auto scaling rules.`,
			},
			"exec_scripts": {
				Type:        schema.TypeList,
				Elem:        scalingPolicyExecScriptSchema(),
				Optional:    true,
				Description: `The list of custom scaling automation scripts.`,
			},
		},
	}
}

func scalingPolicyResourcesPlanSchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"period_type": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `Cycle type of a resource plan.`,
			},
			"start_time": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The start time of a resource plan.`,
			},
			"end_time": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `End time of a resource plan.`,
			},
			"min_capacity": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: `Minimum number of the preserved nodes in a node group in a resource plan. Value range: 0 to 500.`,
			},
			"max_capacity": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: `Maximum number of the preserved nodes in a node group in a resource plan. Value range: 0 to 500.`,
			},
		},
	}
	return &sc
}

func scalingPolicyRuleSchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `Name of an auto scaling rule.`,
			},
			"adjustment_type": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `Auto scaling rule adjustment type.`,
			},
			"cool_down_minutes": {
				Type:     schema.TypeInt,
				Required: true,
				Description: `Cluster cooling time after an auto scaling rule is triggered, when no auto scaling 
					operation is performed.`,
			},
			"scaling_adjustment": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: `Number of nodes that can be adjusted once. Value range: 1 to 100.`,
			},
			"trigger": {
				Type:     schema.TypeList,
				MaxItems: 1,
				Elem:     scalingPolicyRuleTriggerSchema(),
				Required: true,
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `Description about an auto scaling rule.`,
			},
		},
	}
	return &sc
}

func scalingPolicyRuleTriggerSchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"metric_name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `Metric name.`,
			},
			"metric_value": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `Metric threshold to trigger a rule.`,
			},
			"comparison_operator": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `Metric judgment logic operator.`,
			},
			"evaluation_periods": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: `Number of consecutive five-minute periods, during which a metric threshold is reached.`,
			},
		},
	}
	return &sc
}

func scalingPolicyExecScriptSchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `Name of a custom automation script.`,
			},
			"uri": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `Path of a custom automation script.`,
			},
			"parameters": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `Parameters of a custom automation script.`,
			},
			"nodes": {
				Type:        schema.TypeList,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Required:    true,
				Description: `Type of a node where the custom automation script is executed.`,
			},
			"active_master": {
				Type:        schema.TypeBool,
				Optional:    true,
				Computed:    true,
				Description: `Whether the custom automation script runs only on the active Master node.`,
			},
			"action_stage": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `Time when a script is executed.`,
			},
			"fail_action": {
				Type:     schema.TypeString,
				Required: true,
				Description: `Whether to continue to execute subsequent scripts and create a cluster after
					 the custom automation script fails to be executed.`,
			},
		},
	}
	return &sc
}

func resourceScalingPolicyCreateOrUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var (
		createScalingPolicyHttpUrl = "v1.1/{project_id}/autoscaling-policy/{cluster_id}"
		createScalingPolicyProduct = "mrs"
	)
	createScalingPolicyClient, err := cfg.NewServiceClient(createScalingPolicyProduct, region)
	if err != nil {
		return diag.Errorf("error creating MRS client: %s", err)
	}

	createScalingPolicyPath := createScalingPolicyClient.Endpoint + createScalingPolicyHttpUrl
	createScalingPolicyPath = strings.ReplaceAll(createScalingPolicyPath, "{project_id}", createScalingPolicyClient.ProjectID)
	createScalingPolicyPath = strings.ReplaceAll(createScalingPolicyPath, "{cluster_id}", d.Get("cluster_id").(string))

	createScalingPolicyOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		OkCodes: []int{
			200,
		},
		MoreHeaders: map[string]string{"Content-Type": "application/json"},
	}

	createScalingPolicyOpt.JSONBody = buildCreateScalingPolicyBodyParams(d)
	_, err = createScalingPolicyClient.Request("POST", createScalingPolicyPath, &createScalingPolicyOpt)
	if err != nil {
		return diag.Errorf("error creating or updating scaling policy: %s", err)
	}

	d.SetId(fmt.Sprintf("%s/%s", d.Get("cluster_id"), d.Get("node_group").(string)))

	return resourceScalingPolicyRead(ctx, d, meta)
}

func buildCreateScalingPolicyBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"node_group": d.Get("node_group"),
		"auto_scaling_policy": map[string]interface{}{
			"auto_scaling_enable": d.Get("auto_scaling_enable"),
			"min_capacity":        d.Get("min_capacity"),
			"max_capacity":        d.Get("max_capacity"),
			"resources_plans":     buildCreateScalingPolicyRequestBodyResourcesPlan(d.Get("resources_plans")),
			"rules":               buildCreateScalingPolicyRequestBodyRule(d.Get("rules")),
			"exec_scripts":        buildCreateScalingPolicyRequestBodyExecScript(d.Get("exec_scripts")),
		},
	}
	return bodyParams
}

func buildCreateScalingPolicyRequestBodyResourcesPlan(rawParams interface{}) []map[string]interface{} {
	if rawArray, ok := rawParams.([]interface{}); ok {
		if len(rawArray) == 0 {
			return nil
		}

		rst := make([]map[string]interface{}, len(rawArray))
		for i, v := range rawArray {
			if raw, ok := v.(map[string]interface{}); ok {
				rst[i] = map[string]interface{}{
					"period_type":  utils.ValueIgnoreEmpty(raw["period_type"]),
					"start_time":   utils.ValueIgnoreEmpty(raw["start_time"]),
					"end_time":     utils.ValueIgnoreEmpty(raw["end_time"]),
					"min_capacity": utils.ValueIgnoreEmpty(raw["min_capacity"]),
					"max_capacity": utils.ValueIgnoreEmpty(raw["max_capacity"]),
				}
			}
		}
		return rst
	}
	return nil
}

func buildCreateScalingPolicyRequestBodyRule(rawParams interface{}) []map[string]interface{} {
	if rawArray, ok := rawParams.([]interface{}); ok {
		if len(rawArray) == 0 {
			return nil
		}

		rst := make([]map[string]interface{}, len(rawArray))
		for i, v := range rawArray {
			if raw, ok := v.(map[string]interface{}); ok {
				rst[i] = map[string]interface{}{
					"name":               utils.ValueIgnoreEmpty(raw["name"]),
					"adjustment_type":    utils.ValueIgnoreEmpty(raw["adjustment_type"]),
					"cool_down_minutes":  utils.ValueIgnoreEmpty(raw["cool_down_minutes"]),
					"scaling_adjustment": utils.ValueIgnoreEmpty(raw["scaling_adjustment"]),
					"trigger":            buildRuleTrigger(raw["trigger"]),
					"description":        utils.ValueIgnoreEmpty(raw["description"]),
				}
			}
		}
		return rst
	}
	return nil
}

func buildRuleTrigger(rawParams interface{}) map[string]interface{} {
	if rawArray, ok := rawParams.([]interface{}); ok {
		if len(rawArray) == 0 {
			return nil
		}
		raw, ok := rawArray[0].(map[string]interface{})
		if !ok {
			return nil
		}

		params := map[string]interface{}{
			"metric_name":         utils.ValueIgnoreEmpty(raw["metric_name"]),
			"metric_value":        utils.ValueIgnoreEmpty(raw["metric_value"]),
			"comparison_operator": utils.ValueIgnoreEmpty(raw["comparison_operator"]),
			"evaluation_periods":  utils.ValueIgnoreEmpty(raw["evaluation_periods"]),
		}
		return params
	}
	return nil
}

func buildCreateScalingPolicyRequestBodyExecScript(rawParams interface{}) []map[string]interface{} {
	rst := make([]map[string]interface{}, 0, 10) // if want to remove all exec_scripts, set a array with 0 length.
	if rawArray, ok := rawParams.([]interface{}); ok {
		if len(rawArray) == 0 {
			return rst
		}

		rst := make([]map[string]interface{}, len(rawArray))
		for i, v := range rawArray {
			if raw, ok := v.(map[string]interface{}); ok {
				rst[i] = map[string]interface{}{
					"name":          utils.ValueIgnoreEmpty(raw["name"]),
					"uri":           utils.ValueIgnoreEmpty(raw["uri"]),
					"parameters":    utils.ValueIgnoreEmpty(raw["parameters"]),
					"nodes":         utils.ValueIgnoreEmpty(raw["nodes"]),
					"active_master": utils.ValueIgnoreEmpty(raw["active_master"]),
					"action_stage":  utils.ValueIgnoreEmpty(raw["action_stage"]),
					"fail_action":   utils.ValueIgnoreEmpty(raw["fail_action"]),
				}
			}
		}
		return rst
	}
	return rst
}

func resourceScalingPolicyRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var mErr *multierror.Error

	var (
		getScalingPolicyHttpUrl = "v2/{project_id}/autoscaling-policy/{cluster_id}"
		getScalingPolicyProduct = "mrs"
	)
	getScalingPolicyClient, err := cfg.NewServiceClient(getScalingPolicyProduct, region)
	if err != nil {
		return diag.Errorf("error creating MRS client: %s", err)
	}

	getScalingPolicyPath := getScalingPolicyClient.Endpoint + getScalingPolicyHttpUrl
	getScalingPolicyPath = strings.ReplaceAll(getScalingPolicyPath, "{project_id}", getScalingPolicyClient.ProjectID)
	getScalingPolicyPath = strings.ReplaceAll(getScalingPolicyPath, "{cluster_id}", d.Get("cluster_id").(string))

	getScalingPolicyOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		OkCodes: []int{
			200,
		},
		MoreHeaders: map[string]string{"Content-Type": "application/json"},
	}

	getScalingPolicyResp, err := getScalingPolicyClient.Request("GET", getScalingPolicyPath, &getScalingPolicyOpt)

	if err != nil {
		return common.CheckDeletedDiag(d, err, "error retrieving scaling policy")
	}

	getScalingPolicyRespBody, err := utils.FlattenResponse(getScalingPolicyResp)
	if err != nil {
		return diag.FromErr(err)
	}

	jsonPath := fmt.Sprintf("[?node_group_name =='%s']|[0]", d.Get("node_group").(string))
	scalingPolicy := utils.PathSearch(jsonPath, getScalingPolicyRespBody, nil)
	if scalingPolicy == nil {
		return common.CheckDeletedDiag(d, golangsdk.ErrDefault404{}, "no data found")
	}

	mErr = multierror.Append(
		mErr,
		d.Set("region", region),
		d.Set("node_group", utils.PathSearch("node_group_name", scalingPolicy, nil)),
		d.Set("auto_scaling_enable", utils.PathSearch("auto_scaling_policy.auto_scaling_enable", scalingPolicy, nil)),
		d.Set("min_capacity", utils.PathSearch("auto_scaling_policy.min_capacity", scalingPolicy, nil)),
		d.Set("max_capacity", utils.PathSearch("auto_scaling_policy.max_capacity", scalingPolicy, nil)),
		d.Set("resources_plans", flattenGetScalingPolicyResponseBodyResourcesPlan(scalingPolicy)),
		d.Set("rules", flattenGetScalingPolicyResponseBodyRule(scalingPolicy)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenGetScalingPolicyResponseBodyResourcesPlan(resp interface{}) []interface{} {
	if resp == nil {
		return nil
	}
	curJson := utils.PathSearch("auto_scaling_policy.resources_plans", resp, make([]interface{}, 0))
	curArray := curJson.([]interface{})
	rst := make([]interface{}, 0, len(curArray))
	for _, v := range curArray {
		rst = append(rst, map[string]interface{}{
			"period_type":  utils.PathSearch("period_type", v, nil),
			"start_time":   utils.PathSearch("start_time", v, nil),
			"end_time":     utils.PathSearch("end_time", v, nil),
			"min_capacity": utils.PathSearch("min_capacity", v, nil),
			"max_capacity": utils.PathSearch("max_capacity", v, nil),
		})
	}
	return rst
}

func flattenGetScalingPolicyResponseBodyRule(resp interface{}) []interface{} {
	if resp == nil {
		return nil
	}
	curJson := utils.PathSearch("auto_scaling_policy.rules", resp, make([]interface{}, 0))
	curArray := curJson.([]interface{})
	rst := make([]interface{}, 0, len(curArray))
	for _, v := range curArray {
		rst = append(rst, map[string]interface{}{
			"name":               utils.PathSearch("name", v, nil),
			"adjustment_type":    utils.PathSearch("adjustment_type", v, nil),
			"cool_down_minutes":  utils.PathSearch("cool_down_minutes", v, nil),
			"scaling_adjustment": utils.PathSearch("scaling_adjustment", v, nil),
			"trigger":            flattenRuleTrigger(v),
			"description":        utils.PathSearch("description", v, nil),
		})
	}
	return rst
}

func flattenRuleTrigger(resp interface{}) []interface{} {
	var rst []interface{}
	curJson := utils.PathSearch("trigger", resp, nil)
	if curJson == nil {
		log.Printf("[ERROR] error parsing trigger from response")
		return rst
	}

	rst = []interface{}{
		map[string]interface{}{
			"metric_name":         utils.PathSearch("metric_name", curJson, nil),
			"metric_value":        utils.PathSearch("metric_value", curJson, nil),
			"comparison_operator": utils.PathSearch("comparison_operator", curJson, nil),
			"evaluation_periods":  utils.PathSearch("evaluation_periods", curJson, nil),
		},
	}
	return rst
}

func resourceScalingPolicyDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var (
		deleteScalingPolicyHttpUrl = "v1.1/{project_id}/autoscaling-policy/{cluster_id}"
		deleteScalingPolicyProduct = "mrs"
	)
	deleteScalingPolicyClient, err := cfg.NewServiceClient(deleteScalingPolicyProduct, region)
	if err != nil {
		return diag.Errorf("error creating MRS client: %s", err)
	}

	deleteScalingPolicyPath := deleteScalingPolicyClient.Endpoint + deleteScalingPolicyHttpUrl
	deleteScalingPolicyPath = strings.ReplaceAll(deleteScalingPolicyPath, "{project_id}", deleteScalingPolicyClient.ProjectID)
	deleteScalingPolicyPath = strings.ReplaceAll(deleteScalingPolicyPath, "{cluster_id}", d.Get("cluster_id").(string))

	deleteScalingPolicyOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		OkCodes: []int{
			200,
		},
		MoreHeaders: map[string]string{"Content-Type": "application/json"},
	}

	deleteScalingPolicyOpt.JSONBody = buildDeleteScalingPolicyBodyParams(d)
	_, err = deleteScalingPolicyClient.Request("POST", deleteScalingPolicyPath, &deleteScalingPolicyOpt)
	if err != nil {
		return diag.Errorf("error deleting scaling policy: %s", err)
	}

	return nil
}

func buildDeleteScalingPolicyBodyParams(d *schema.ResourceData) map[string]interface{} {
	nodeGroup, _ := d.GetChange("node_group")
	minCapacity, _ := d.GetChange("min_capacity")
	maxCapacity, _ := d.GetChange("max_capacity")
	resourcePlans, _ := d.GetChange("resources_plans")
	rules, _ := d.GetChange("rules")

	bodyParams := map[string]interface{}{
		"node_group": nodeGroup,
		"auto_scaling_policy": map[string]interface{}{
			"auto_scaling_enable": false,
			"min_capacity":        minCapacity,
			"max_capacity":        maxCapacity,
			"resources_plans":     buildCreateScalingPolicyRequestBodyResourcesPlan(resourcePlans),
			"rules":               buildCreateScalingPolicyRequestBodyRule(rules),
			"exec_scripts":        []map[string]interface{}{},
		},
	}
	return bodyParams
}

func resourceScalingPolicyImportState(_ context.Context, d *schema.ResourceData, _ interface{}) ([]*schema.ResourceData, error) {
	parts := strings.SplitN(d.Id(), "/", 2)
	if len(parts) != 2 {
		return nil, fmt.Errorf("invalid format specified for import id, must be <cluster_id>/<node_group>")
	}

	d.Set("cluster_id", parts[0])
	d.Set("node_group", parts[1])

	return []*schema.ResourceData{d}, nil
}
