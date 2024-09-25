package cce

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/jmespath/go-jmespath"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API CCE GET /api/v3/projects/{project_id}/clusters/{cluster_id}
// @API CCE POST /api/v3/projects/{project_id}/clusters/{cluster_id}/operation/upgrade
// @API CCE GET /api/v3/projects/{project_id}/clusters/{cluster_id}/operation/upgrade/tasks/{task_id}
var clusterUpgradeNonUpdatableParams = []string{"cluster_id", "target_version", "addons",
	"addons.*.addon_template_name",
	"addons.*.operation",
	"addons.*.version",
	"addons.*.values",
	"node_order", "nodepool_order",
	"strategy",
	"strategy.*.type",
	"strategy.*.in_place_rolling_update",
	"strategy.*.in_place_rolling_update.*.user_defined_step",
}

func ResourceClusterUpgrade() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceClusterUpgradeCreate,
		ReadContext:   resourceClusterUpgradeRead,
		UpdateContext: resourceClusterUpgradeUpdate,
		DeleteContext: resourceClusterUpgradeDelete,

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(20 * time.Minute),
		},

		CustomizeDiff: config.FlexibleForceNew(clusterUpgradeNonUpdatableParams),

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"cluster_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"target_version": {
				Type:     schema.TypeString,
				Required: true,
			},
			"addons": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"addon_template_name": {
							Type:     schema.TypeString,
							Required: true,
						},
						"operation": {
							Type:     schema.TypeString,
							Required: true,
						},
						"version": {
							Type:     schema.TypeString,
							Required: true,
						},
						"values": {
							Type:     schema.TypeList,
							Optional: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"basic_json": {
										Type:         schema.TypeString,
										Optional:     true,
										ValidateFunc: validation.StringIsJSON,
									},
									"custom_json": {
										Type:         schema.TypeString,
										Optional:     true,
										ValidateFunc: validation.StringIsJSON,
									},
									"flavor_json": {
										Type:         schema.TypeString,
										Optional:     true,
										ValidateFunc: validation.StringIsJSON,
									},
								},
							},
						},
					},
				},
			},
			"node_order": {
				Type:     schema.TypeMap,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"nodepool_order": {
				Type:     schema.TypeMap,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeInt},
			},
			"strategy": {
				Type:     schema.TypeList,
				Required: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"type": {
							Type:     schema.TypeString,
							Required: true,
						},
						"in_place_rolling_update": {
							Type:     schema.TypeList,
							Optional: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"user_defined_step": {
										Type:     schema.TypeInt,
										Optional: true,
									},
								}},
						},
					}},
			},
			"enable_force_new": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"true", "false"}, false),
				Description:  utils.SchemaDesc("", utils.SchemaDescInput{Internal: true}),
			},
		},
	}
}

func buildClusterUpgradeCreateOpts(d *schema.ResourceData) (map[string]interface{}, error) {
	nodeOrder, err := buildClusterUpgradeNodeOrderOpts(d)
	if err != nil {
		return nil, fmt.Errorf("error building node_order Opts: %s", err)
	}

	addons, err := buildClusterUpgradeAddonsOpts(d)
	if err != nil {
		return nil, fmt.Errorf("error building addons Opts: %s", err)
	}

	result := map[string]interface{}{
		"metadata": map[string]interface{}{
			"kind":       "UpgradeTask",
			"apiVersion": "v3",
		},
		"spec": map[string]interface{}{
			"clusterUpgradeAction": map[string]interface{}{
				"addons":        addons,
				"nodeOrder":     nodeOrder,
				"nodePoolOrder": d.Get("nodepool_order"),
				"strategy":      buildClusterUpgradeStrategyOpts(d),
				"targetVersion": d.Get("target_version"),
			},
		},
	}
	return result, nil
}

func buildClusterUpgradeAddonsOpts(d *schema.ResourceData) ([]map[string]interface{}, error) {
	addonsRaw := d.Get("addons").([]interface{})
	if len(addonsRaw) == 0 {
		return nil, nil
	}

	result := make([]map[string]interface{}, len(addonsRaw))

	for i, v := range addonsRaw {
		if addon, ok := v.(map[string]interface{}); ok {
			values, err := buildClusterUpgradeAddonsValuesOpts(addon["values"].([]interface{}))
			if err != nil {
				return nil, err
			}
			result[i] = map[string]interface{}{
				"addonTemplateName": addon["addon_template_name"],
				"operation":         addon["operation"],
				"version":           addon["version"],
				"values":            values,
			}
		}
	}
	return result, nil
}

func buildClusterUpgradeAddonsValuesOpts(valuesRaw []interface{}) (map[string]interface{}, error) {
	if len(valuesRaw) == 0 {
		return nil, nil
	}

	if valuesMap, ok := valuesRaw[0].(map[string]interface{}); ok {
		var basic, custom, flavor map[string]interface{}
		if basicJsonRaw := valuesMap["basic_json"].(string); basicJsonRaw != "" {
			err := json.Unmarshal([]byte(basicJsonRaw), &basic)
			if err != nil {
				return nil, fmt.Errorf("error unmarshalling basic json: %s", err)
			}
		}
		if customJsonRaw := valuesMap["custom_json"].(string); customJsonRaw != "" {
			err := json.Unmarshal([]byte(customJsonRaw), &custom)
			if err != nil {
				return nil, fmt.Errorf("error unmarshalling custom json: %s", err)
			}
		}
		if flavorJsonRaw := valuesMap["flavor_json"].(string); flavorJsonRaw != "" {
			err := json.Unmarshal([]byte(flavorJsonRaw), &flavor)
			if err != nil {
				return nil, fmt.Errorf("error unmarshalling flavor json %s", err)
			}
		}

		result := map[string]interface{}{
			"basic":  basic,
			"custom": custom,
			"flavor": flavor,
		}

		return result, nil
	}
	return nil, nil
}

func buildClusterUpgradeNodeOrderOpts(d *schema.ResourceData) (map[string]interface{}, error) {
	nodeOrderRaw := d.Get("node_order").(map[string]interface{})
	if len(nodeOrderRaw) == 0 {
		return nil, nil
	}
	result := make(map[string]interface{}, len(nodeOrderRaw))

	for k, v := range nodeOrderRaw {
		var value []map[string]interface{}
		err := json.Unmarshal([]byte(v.(string)), &value)
		if err != nil {
			return nil, err
		}
		result[k] = value
	}
	return result, nil
}

func buildClusterUpgradeStrategyOpts(d *schema.ResourceData) map[string]interface{} {
	strategyRaw := d.Get("strategy").([]interface{})
	if len(strategyRaw) == 0 {
		return nil
	}

	if strategy, ok := strategyRaw[0].(map[string]interface{}); ok {
		return map[string]interface{}{
			"type":                 strategy["type"],
			"inPlaceRollingUpdate": buildClusterUpgradeInPlaceRollingUpdateOpts(strategy["in_place_rolling_update"].([]interface{})),
		}
	}

	return nil
}

func buildClusterUpgradeInPlaceRollingUpdateOpts(inPlaceRollingUpdateRaw []interface{}) map[string]interface{} {
	if len(inPlaceRollingUpdateRaw) == 0 {
		return nil
	}

	if inPlaceRollingUpdate, ok := inPlaceRollingUpdateRaw[0].(map[string]interface{}); ok {
		return map[string]interface{}{
			"userDefinedStep": utils.ValueIgnoreEmpty(inPlaceRollingUpdate["user_defined_step"]),
		}
	}

	return nil
}

func resourceClusterUpgradeCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	client, err := cfg.CceV3Client(cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating CCE v3 client: %s", err)
	}

	// Wait for the cce cluster to become available
	clusterID := d.Get("cluster_id").(string)
	stateCluster := &resource.StateChangeConf{
		Pending:      []string{"PENDING"},
		Target:       []string{"COMPLETED"},
		Refresh:      clusterStateRefreshFunc(client, clusterID, []string{"Available"}),
		Timeout:      d.Timeout(schema.TimeoutCreate),
		Delay:        5 * time.Second,
		PollInterval: 5 * time.Second,
	}
	_, err = stateCluster.WaitForStateContext(ctx)
	if err != nil {
		return diag.Errorf("error waiting for CCE cluster to become available: %s", err)
	}

	var (
		createClusterUpgradeHttpUrl = "api/v3/projects/{project_id}/clusters/{cluster_id}/operation/upgrade"
	)

	createClusterUpgradePath := client.Endpoint + createClusterUpgradeHttpUrl
	createClusterUpgradePath = strings.ReplaceAll(createClusterUpgradePath, "{project_id}", client.ProjectID)
	createClusterUpgradePath = strings.ReplaceAll(createClusterUpgradePath, "{cluster_id}", clusterID)

	createClusterUpgradeOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	createOpts, err := buildClusterUpgradeCreateOpts(d)
	if err != nil {
		return nil
	}

	createClusterUpgradeOpt.JSONBody = utils.RemoveNil(createOpts)
	createClusterUpgradeResp, err := client.Request("POST",
		createClusterUpgradePath, &createClusterUpgradeOpt)
	if err != nil {
		return diag.Errorf("error upgrading CCE cluster: %s", err)
	}

	createClusterUpgradeRespBody, err := utils.FlattenResponse(createClusterUpgradeResp)
	if err != nil {
		return diag.FromErr(err)
	}

	taskID := utils.PathSearch("metadata.uid", createClusterUpgradeRespBody, "")
	if taskID == "" {
		return diag.Errorf("error upgrading CCE cluster: task_id is not found in API response")
	}

	d.SetId(taskID.(string))

	err = clusterUpgradeWaitingForStateCompleted(ctx, d, meta, d.Timeout(schema.TimeoutCreate))
	if err != nil {
		return diag.Errorf("error waiting for upgrading cluster task to complete: %s", err)
	}

	return resourceClusterUpgradeRead(ctx, d, meta)
}

func clusterUpgradeWaitingForStateCompleted(ctx context.Context, d *schema.ResourceData, meta interface{}, t time.Duration) error {
	stateConf := &resource.StateChangeConf{
		Pending: []string{"PENDING"},
		Target:  []string{"COMPLETED"},
		Refresh: func() (interface{}, string, error) {
			config := meta.(*config.Config)
			region := config.GetRegion(d)
			var (
				clusterUpgradeWaitingHttpUrl = "api/v3/projects/{project_id}/clusters/{cluster_id}/operation/upgrade/tasks/{task_id}"
				clusterUpgradeWaitingProduct = "cce"
			)
			client, err := config.NewServiceClient(clusterUpgradeWaitingProduct, region)
			if err != nil {
				return nil, "ERROR", fmt.Errorf("error creating CCE Client: %s", err)
			}

			clusterUpgradeWaitingPath := client.Endpoint + clusterUpgradeWaitingHttpUrl
			clusterUpgradeWaitingPath = strings.ReplaceAll(clusterUpgradeWaitingPath, "{project_id}", client.ProjectID)
			clusterUpgradeWaitingPath = strings.ReplaceAll(clusterUpgradeWaitingPath, "{cluster_id}", d.Get("cluster_id").(string))
			clusterUpgradeWaitingPath = strings.ReplaceAll(clusterUpgradeWaitingPath, "{task_id}", d.Id())

			clusterUpgradeWaitingOpt := golangsdk.RequestOpts{
				KeepResponseBody: true,
			}
			clusterUpgradeWaitingResp, err := client.Request("GET", clusterUpgradeWaitingPath, &clusterUpgradeWaitingOpt)
			if err != nil {
				return nil, "ERROR", err
			}

			clusterUpgradeWaitingRespBody, err := utils.FlattenResponse(clusterUpgradeWaitingResp)
			if err != nil {
				return nil, "ERROR", err
			}
			statusRaw, err := jmespath.Search(`status.phase`, clusterUpgradeWaitingRespBody)
			if err != nil {
				return nil, "ERROR", fmt.Errorf("error parse %s from response body", `status.phase`)
			}

			status := fmt.Sprintf("%v", statusRaw)

			targetStatus := []string{
				"Success",
			}
			if utils.StrSliceContains(targetStatus, status) {
				return clusterUpgradeWaitingRespBody, "COMPLETED", nil
			}

			unexpectedStatus := []string{
				"Failed",
			}
			if utils.StrSliceContains(unexpectedStatus, status) {
				return clusterUpgradeWaitingRespBody, status, nil
			}

			return clusterUpgradeWaitingRespBody, "PENDING", nil
		},
		Timeout:      t,
		Delay:        10 * time.Second,
		PollInterval: 10 * time.Second,
	}
	_, err := stateConf.WaitForStateContext(ctx)
	return err
}

func resourceClusterUpgradeRead(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return nil
}

func resourceClusterUpgradeUpdate(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return nil
}

func resourceClusterUpgradeDelete(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	errorMsg := "Deleting cluster upgrade resource is not supported. The cluster upgrade resource is only removed from the state."
	return diag.Diagnostics{
		diag.Diagnostic{
			Severity: diag.Warning,
			Summary:  errorMsg,
		},
	}
}
