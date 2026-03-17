package mrs

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

var (
	scalingPolicyV2NonUpdatableParams = []string{
		"cluster_id",
		"node_group_name",
		"resource_pool_name",
	}

	scalingPolicyV2NotFoundCodes = []string{
		"MRS.00005016", // The scaling policy does not exist.
	}
)

// @API MRS POST /v2/{project_id}/autoscaling-policy/{cluster_id}
// @API MRS GET /v2/{project_id}/autoscaling-policy/{cluster_id}
// @API MRS PUT /v2/{project_id}/autoscaling-policy/{cluster_id}
// @API MRS DELETE /v2/{project_id}/autoscaling-policy/{cluster_id}
func ResourceScalingPolicyV2() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceScalingPolicyV2Create,
		ReadContext:   resourceScalingPolicyV2Read,
		UpdateContext: resourceScalingPolicyV2Update,
		DeleteContext: resourceScalingPolicyV2Delete,

		CustomizeDiff: config.FlexibleForceNew(scalingPolicyV2NonUpdatableParams),

		Importer: &schema.ResourceImporter{
			StateContext: resourceScalingPolicyV2ImportState,
		},

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: `The region where the scaling policy is located.`,
			},
			"cluster_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The ID of the cluster.`,
			},
			"node_group_name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The name of the node group.`,
			},
			"resource_pool_name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The name of the resource pool.`,
			},
			"auto_scaling_policy": {
				Type:     schema.TypeList,
				Required: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"auto_scaling_enable": {
							Type:        schema.TypeBool,
							Optional:    true,
							Description: `Whether to enable the auto scaling policy.`,
						},
						"min_capacity": {
							Type:        schema.TypeInt,
							Optional:    true,
							Description: `The minimum number of nodes in the node group.`,
						},
						"max_capacity": {
							Type:        schema.TypeInt,
							Optional:    true,
							Description: `The maximum number of nodes in the node group.`,
						},
						"resources_plans": {
							Type:     schema.TypeList,
							Optional: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"period_type": {
										Type:        schema.TypeString,
										Required:    true,
										Description: `The period type of the resource plan.`,
									},
									"start_time": {
										Type:        schema.TypeString,
										Required:    true,
										Description: `The start time of the resource plan, in 'HH:mm' format.`,
									},
									"end_time": {
										Type:        schema.TypeString,
										Required:    true,
										Description: `The end time of the resource plan, in 'HH:mm' format.`,
									},
									"min_capacity": {
										Type:        schema.TypeInt,
										Optional:    true,
										Description: `The minimum number of retained nodes for the node group in the resource plan.`,
									},
									"max_capacity": {
										Type:        schema.TypeInt,
										Optional:    true,
										Description: `The maximum number of retained nodes for the node group in the resource plan.`,
									},
									"effective_days": {
										Type:        schema.TypeList,
										Optional:    true,
										Elem:        &schema.Schema{Type: schema.TypeString},
										Description: `The effective day list of the resource plan.`,
									},
								},
							},
							Description: `The list of resource plans.`,
						},
						"rules": {
							Type:     schema.TypeList,
							Optional: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"name": {
										Type:        schema.TypeString,
										Required:    true,
										Description: `The name of the rule.`,
									},
									"adjustment_type": {
										Type:        schema.TypeString,
										Required:    true,
										Description: `The adjustment type of the rule.`,
									},
									"cool_down_minutes": {
										Type:        schema.TypeInt,
										Required:    true,
										Description: `The cool down time of the cluster after the scaling rule is triggered, in minutes.`,
									},
									"scaling_adjustment": {
										Type:        schema.TypeInt,
										Required:    true,
										Description: `The number of adjusted nodes in one scaling action.`,
									},
									"trigger": {
										Type:     schema.TypeList,
										Required: true,
										MaxItems: 1,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"metric_name": {
													Type:        schema.TypeString,
													Required:    true,
													Description: `The name of the metric.`,
												},
												"metric_value": {
													Type:        schema.TypeString,
													Required:    true,
													Description: `The threshold value of the metric.`,
												},
												"evaluation_periods": {
													Type:        schema.TypeInt,
													Required:    true,
													Description: `The number of consecutive periods that meet the metric threshold.`,
												},
												"comparison_operator": {
													Type:        schema.TypeString,
													Optional:    true,
													Computed:    true,
													Description: `The comparison operator of the metric judgment logic.`,
												},
											},
										},
										Description: `The trigger condition list of the rule.`,
									},
									"description": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: `The description of the rule.`,
									},
								},
							},
							Description: `The list of auto scaling rules.`,
						},
						"tags": {
							Type:        schema.TypeMap,
							Optional:    true,
							Elem:        &schema.Schema{Type: schema.TypeString},
							Description: `The key/value pairs associated with auto scaling policy.`,
						},
					},
				},
				Description: `The configurations of the auto scaling policy.`,
			},
			// Internal parameters.
			"enable_force_new": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"true", "false"}, false),
				Description:  utils.SchemaDesc("", utils.SchemaDescInput{Internal: true}),
			},
		},
	}
}

func buildScalingPolicyV2Params(d *schema.ResourceData) map[string]interface{} {
	return map[string]interface{}{
		"node_group_name":     d.Get("node_group_name"),
		"resource_pool_name":  d.Get("resource_pool_name"),
		"auto_scaling_policy": buildScalingPolicyV2AutoScalingPolicy(d.Get("auto_scaling_policy").([]interface{})),
	}
}

func buildScalingPolicyV2AutoScalingPolicy(policies []interface{}) map[string]interface{} {
	if len(policies) < 1 {
		return nil
	}

	policy := policies[0]
	return map[string]interface{}{
		"auto_scaling_enable": utils.PathSearch("auto_scaling_enable", policy, nil),
		"min_capacity":        utils.ValueIgnoreEmpty(utils.PathSearch("min_capacity", policy, nil)),
		"max_capacity":        utils.ValueIgnoreEmpty(utils.PathSearch("max_capacity", policy, nil)),
		"resources_plans": buildScalingPolicyV2ResourcesPlans(utils.PathSearch("resources_plans",
			policy, make([]interface{}, 0)).([]interface{})),
		"rules": buildScalingPolicyV2Rules(utils.PathSearch("rules",
			policy, make([]interface{}, 0)).([]interface{})),
		"tags": utils.ValueIgnoreEmpty(utils.ExpandResourceTags(utils.PathSearch("tags",
			policy, make(map[string]interface{})).(map[string]interface{}))),
	}
}

func buildScalingPolicyV2ResourcesPlans(plans []interface{}) []map[string]interface{} {
	if len(plans) == 0 {
		return nil
	}

	result := make([]map[string]interface{}, 0, len(plans))
	for _, v := range plans {
		result = append(result, map[string]interface{}{
			"period_type":    utils.PathSearch("period_type", v, nil),
			"start_time":     utils.PathSearch("start_time", v, nil),
			"end_time":       utils.PathSearch("end_time", v, nil),
			"min_capacity":   utils.ValueIgnoreEmpty(utils.PathSearch("min_capacity", v, nil)),
			"max_capacity":   utils.ValueIgnoreEmpty(utils.PathSearch("max_capacity", v, nil)),
			"effective_days": utils.ValueIgnoreEmpty(utils.PathSearch("effective_days", v, nil)),
		})
	}

	return result
}

func buildScalingPolicyV2Rules(rules []interface{}) []map[string]interface{} {
	if len(rules) == 0 {
		return nil
	}

	result := make([]map[string]interface{}, 0, len(rules))
	for _, v := range rules {
		result = append(result, map[string]interface{}{
			"name":               utils.PathSearch("name", v, nil),
			"adjustment_type":    utils.PathSearch("adjustment_type", v, nil),
			"cool_down_minutes":  utils.PathSearch("cool_down_minutes", v, nil),
			"scaling_adjustment": utils.PathSearch("scaling_adjustment", v, nil),
			"trigger":            buildScalingPolicyV2RuleTrigger(utils.PathSearch("trigger", v, make([]interface{}, 0)).([]interface{})),
			"description":        utils.ValueIgnoreEmpty(utils.PathSearch("description", v, nil)),
		})
	}

	return result
}

func buildScalingPolicyV2RuleTrigger(triggers []interface{}) map[string]interface{} {
	if len(triggers) == 0 {
		return nil
	}

	trigger := triggers[0]
	return map[string]interface{}{
		"metric_name":         utils.PathSearch("metric_name", trigger, nil),
		"metric_value":        utils.PathSearch("metric_value", trigger, nil),
		"evaluation_periods":  utils.PathSearch("evaluation_periods", trigger, nil),
		"comparison_operator": utils.ValueIgnoreEmpty(utils.PathSearch("comparison_operator", trigger, nil)),
	}
}

func resourceScalingPolicyV2Create(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	client, err := cfg.NewServiceClient("mrs", cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating MRS client: %s", err)
	}

	createPath := client.Endpoint + "v2/{project_id}/autoscaling-policy/{cluster_id}"
	createPath = strings.ReplaceAll(createPath, "{project_id}", client.ProjectID)
	createPath = strings.ReplaceAll(createPath, "{cluster_id}", d.Get("cluster_id").(string))

	opts := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"Content-Type": "application/json"},
		JSONBody:         utils.RemoveNil(buildScalingPolicyV2Params(d)),
	}

	clusterId := d.Get("cluster_id").(string)
	_, err = client.Request("POST", createPath, &opts)
	if err != nil {
		return diag.Errorf("error creating scaling policy for the cluster (%s): %s", clusterId, err)
	}

	d.SetId(fmt.Sprintf("%s/%s/%s", clusterId, d.Get("node_group_name").(string), d.Get("resource_pool_name").(string)))

	return resourceScalingPolicyV2Read(ctx, d, meta)
}

// There are no paging parameters, query all data.
func listScalingPolicyV2(client *golangsdk.ServiceClient, clusterId string) (interface{}, error) {
	httpUrl := "v2/{project_id}/autoscaling-policy/{cluster_id}"
	listPath := client.Endpoint + httpUrl
	listPath = strings.ReplaceAll(listPath, "{project_id}", client.ProjectID)
	listPath = strings.ReplaceAll(listPath, "{cluster_id}", clusterId)

	opts := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"Content-Type": "application/json"},
	}

	resp, err := client.Request("GET", listPath, &opts)
	if err != nil {
		return nil, err
	}

	return utils.FlattenResponse(resp)
}

func GetScalingPolicyV2(client *golangsdk.ServiceClient, clusterId, nodeGroupName, resourcePoolName string) (interface{}, error) {
	policies, err := listScalingPolicyV2(client, clusterId)
	if err != nil {
		return nil, err
	}

	policy := utils.PathSearch(fmt.Sprintf("[?node_group_name=='%s'&&resource_pool_name=='%s']|[0]", nodeGroupName, resourcePoolName),
		policies, nil)
	if policy == nil {
		return nil, golangsdk.ErrDefault404{
			ErrUnexpectedResponseCode: golangsdk.ErrUnexpectedResponseCode{
				Method:    "GET",
				URL:       "/v2/{project_id}/autoscaling-policy/{cluster_id}",
				RequestId: "NONE",
				Body:      []byte("the scaling policy does not exist"),
			},
		}
	}

	return policy, nil
}

func resourceScalingPolicyV2Read(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg    = meta.(*config.Config)
		region = cfg.GetRegion(d)
	)

	client, err := cfg.NewServiceClient("mrs", region)
	if err != nil {
		return diag.Errorf("error creating MRS client: %s", err)
	}

	clusterId := d.Get("cluster_id").(string)
	policy, err := GetScalingPolicyV2(client, clusterId, d.Get("node_group_name").(string), d.Get("resource_pool_name").(string))
	if err != nil {
		return common.CheckDeletedDiag(
			d,
			common.ConvertExpected400ErrInto404Err(err, "error_code", clusterNotFoundCodes...),
			"error retrieving scaling policy",
		)
	}

	mErr := multierror.Append(
		d.Set("region", region),
		d.Set("cluster_id", clusterId),
		d.Set("node_group_name", utils.PathSearch("node_group_name", policy, nil)),
		d.Set("resource_pool_name", utils.PathSearch("resource_pool_name", policy, nil)),
		d.Set("auto_scaling_policy", flattenScalingPolicyV2AutoScalingPolicy(utils.PathSearch("auto_scaling_policy", policy, nil))),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenScalingPolicyV2AutoScalingPolicy(policy interface{}) []interface{} {
	return []interface{}{
		map[string]interface{}{
			"auto_scaling_enable": utils.PathSearch("auto_scaling_enable", policy, nil),
			"min_capacity":        utils.PathSearch("min_capacity", policy, nil),
			"max_capacity":        utils.PathSearch("max_capacity", policy, nil),
			"resources_plans": flattenScalingPolicyV2ResourcesPlans(utils.PathSearch("resources_plans",
				policy, make([]interface{}, 0)).([]interface{})),
			"rules": flattenScalingPolicyV2Rules(utils.PathSearch("rules",
				policy, make([]interface{}, 0)).([]interface{})),
			"tags": utils.FlattenTagsToMap(utils.PathSearch("tags", policy, make([]interface{}, 0))),
		},
	}
}

func flattenScalingPolicyV2ResourcesPlans(plans []interface{}) []interface{} {
	if len(plans) == 0 {
		return nil
	}

	result := make([]interface{}, 0, len(plans))
	for _, v := range plans {
		result = append(result, map[string]interface{}{
			"period_type":    utils.PathSearch("period_type", v, nil),
			"start_time":     utils.PathSearch("start_time", v, nil),
			"end_time":       utils.PathSearch("end_time", v, nil),
			"min_capacity":   utils.PathSearch("min_capacity", v, nil),
			"max_capacity":   utils.PathSearch("max_capacity", v, nil),
			"effective_days": utils.PathSearch("effective_days", v, nil),
		})
	}

	return result
}

func flattenScalingPolicyV2Rules(rules []interface{}) []interface{} {
	if len(rules) == 0 {
		return nil
	}

	result := make([]interface{}, 0, len(rules))
	for _, v := range rules {
		result = append(result, map[string]interface{}{
			"name":               utils.PathSearch("name", v, nil),
			"adjustment_type":    utils.PathSearch("adjustment_type", v, nil),
			"cool_down_minutes":  utils.PathSearch("cool_down_minutes", v, nil),
			"scaling_adjustment": utils.PathSearch("scaling_adjustment", v, nil),
			"trigger":            flattenScalingPolicyV2RuleTrigger(utils.PathSearch("trigger", v, nil)),
			"description":        utils.PathSearch("description", v, nil),
		})
	}

	return result
}

func flattenScalingPolicyV2RuleTrigger(trigger interface{}) []interface{} {
	return []interface{}{
		map[string]interface{}{
			"metric_name":         utils.PathSearch("metric_name", trigger, nil),
			"metric_value":        utils.PathSearch("metric_value", trigger, nil),
			"evaluation_periods":  utils.PathSearch("evaluation_periods", trigger, nil),
			"comparison_operator": utils.PathSearch("comparison_operator", trigger, nil),
		},
	}
}

func resourceScalingPolicyV2Update(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	client, err := cfg.NewServiceClient("mrs", cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating MRS client: %s", err)
	}

	updatePath := client.Endpoint + "v2/{project_id}/autoscaling-policy/{cluster_id}"
	updatePath = strings.ReplaceAll(updatePath, "{project_id}", client.ProjectID)
	updatePath = strings.ReplaceAll(updatePath, "{cluster_id}", d.Get("cluster_id").(string))

	updateOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"Content-Type": "application/json"},
		JSONBody:         utils.RemoveNil(buildScalingPolicyV2Params(d)),
	}

	_, err = client.Request("PUT", updatePath, &updateOpt)
	if err != nil {
		return diag.Errorf("error updating scaling policy: %s", err)
	}

	return resourceScalingPolicyV2Read(ctx, d, meta)
}

func resourceScalingPolicyV2Delete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	client, err := cfg.NewServiceClient("mrs", cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating MRS client: %s", err)
	}

	deletePath := client.Endpoint + "v2/{project_id}/autoscaling-policy/{cluster_id}"
	deletePath = strings.ReplaceAll(deletePath, "{project_id}", client.ProjectID)
	deletePath = strings.ReplaceAll(deletePath, "{cluster_id}", d.Get("cluster_id").(string))

	deleteOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"Content-Type": "application/json"},
		JSONBody: map[string]interface{}{
			"node_group_name":    d.Get("node_group_name"),
			"resource_pool_name": d.Get("resource_pool_name"),
		},
	}

	_, err = client.Request("DELETE", deletePath, &deleteOpt)
	if err != nil {
		return common.CheckDeletedDiag(
			d,
			common.ConvertExpected400ErrInto404Err(err, "error_code", append(clusterNotFoundCodes, scalingPolicyV2NotFoundCodes...)...),
			"error deleting scaling policy",
		)
	}

	return nil
}

func resourceScalingPolicyV2ImportState(_ context.Context, d *schema.ResourceData, _ interface{}) ([]*schema.ResourceData, error) {
	parts := strings.SplitN(d.Id(), "/", 3)
	if len(parts) != 3 {
		return nil, fmt.Errorf("invalid format specified for import ID, must be `<cluster_id>/<node_group_name>/<resource_pool_name>`, "+
			"but got '%s'", d.Id())
	}

	mErr := multierror.Append(
		d.Set("cluster_id", parts[0]),
		d.Set("node_group_name", parts[1]),
		d.Set("resource_pool_name", parts[2]),
	)

	return []*schema.ResourceData{d}, mErr.ErrorOrNil()
}
