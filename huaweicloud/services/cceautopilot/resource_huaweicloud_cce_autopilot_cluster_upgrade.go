package cceautopilot

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API CCE POST /autopilot/v3/projects/{project_id}/clusters/{cluster_id}/operation/upgradeworkflows
// @API CCE POST /autopilot/v3/projects/{project_id}/clusters/{cluster_id}/operation/upgrade
// @API CCE GET /autopilot/v3/projects/{project_id}/clusters/{cluster_id}/operation/upgrade/tasks/{task_id}
// @API CCE POST /autopilot/v3/projects/{project_id}/clusters/{cluster_id}/operation/precheck
// @API CCE GET /autopilot/v3/projects/{project_id}/clusters/{cluster_id}/operation/precheck/tasks/{task_id}
// @API CCE POST /autopilot/v3.1/projects/{project_id}/clusters/{cluster_id}/operation/snapshot
// @API CCE GET /autopilot/v3.1/projects/{project_id}/clusters/{cluster_id}/operation/snapshot/tasks
// @API CCE POST /autopilot/v3/projects/{project_id}/clusters/{cluster_id}/operation/postcheck
var autopilotClusterUpgradeNonUpdatableParams = []string{"cluster_id", "target_version", "current_version", "addons", "is_snapshot",
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

func ResourceAutopilotClusterUpgrade() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceAutopilotClusterUpgradeCreate,
		ReadContext:   resourceAutopilotClusterUpgradeRead,
		UpdateContext: resourceAutopilotClusterUpgradeUpdate,
		DeleteContext: resourceAutopilotClusterUpgradeDelete,

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(60 * time.Minute),
		},

		CustomizeDiff: config.FlexibleForceNew(autopilotClusterUpgradeNonUpdatableParams),

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
			"current_version": {
				Type:     schema.TypeString,
				Optional: true,
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
			"is_snapshot": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"is_postcheck": {
				Type:     schema.TypeBool,
				Optional: true,
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

func buildAutopilotClusterUpgradeCreateOpts(d *schema.ResourceData, targetVersion string) (map[string]interface{}, error) {
	nodeOrder, err := buildAutopilotClusterUpgradeNodeOrderOpts(d)
	if err != nil {
		return nil, fmt.Errorf("error building node_order Opts: %s", err)
	}

	addons, err := buildAutopilotClusterUpgradeAddonsOpts(d)
	if err != nil {
		return nil, fmt.Errorf("error building addons Opts: %s", err)
	}

	result := map[string]interface{}{
		"metadata": map[string]interface{}{
			"kind":       "UpgradeTask",
			"apiVersion": "v3",
		},
		"spec": map[string]interface{}{
			"autopilotClusterUpgradeAction": map[string]interface{}{
				"addons":        addons,
				"nodeOrder":     nodeOrder,
				"nodePoolOrder": d.Get("nodepool_order"),
				"strategy":      buildAutopilotClusterUpgradeStrategyOpts(d),
				"targetVersion": targetVersion,
			},
		},
	}
	return result, nil
}

func buildAutopilotClusterUpgradeAddonsOpts(d *schema.ResourceData) ([]map[string]interface{}, error) {
	addonsRaw := d.Get("addons").([]interface{})
	if len(addonsRaw) == 0 {
		return nil, nil
	}

	result := make([]map[string]interface{}, len(addonsRaw))

	for i, v := range addonsRaw {
		if addon, ok := v.(map[string]interface{}); ok {
			values, err := buildAutopilotClusterUpgradeAddonsValuesOpts(addon["values"].([]interface{}))
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

func buildAutopilotClusterUpgradeAddonsValuesOpts(valuesRaw []interface{}) (map[string]interface{}, error) {
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

func buildAutopilotClusterUpgradeNodeOrderOpts(d *schema.ResourceData) (map[string]interface{}, error) {
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

func buildAutopilotClusterUpgradeStrategyOpts(d *schema.ResourceData) map[string]interface{} {
	strategyRaw := d.Get("strategy").([]interface{})
	if len(strategyRaw) == 0 {
		return nil
	}

	if strategy, ok := strategyRaw[0].(map[string]interface{}); ok {
		return map[string]interface{}{
			"type":                 strategy["type"],
			"inPlaceRollingUpdate": buildAutopilotClusterUpgradeInPlaceRollingUpdateOpts(strategy["in_place_rolling_update"].([]interface{})),
		}
	}

	return nil
}

func buildAutopilotClusterUpgradeInPlaceRollingUpdateOpts(inPlaceRollingUpdateRaw []interface{}) map[string]interface{} {
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

func resourceAutopilotClusterUpgradeCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	client, err := cfg.CceV3Client(cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating CCE v3 client: %s", err)
	}

	clusterID := d.Get("cluster_id").(string)
	currentVersion := d.Get("current_version").(string)
	targetVersion := d.Get("target_version").(string)

	// workflow
	workflowResp, err := createWorkflow(client, clusterID, currentVersion, targetVersion)
	if err != nil {
		return diag.Errorf("error creating CCE Autopilot cluster precheck: %s", err)
	}

	exactCurrentVersion := utils.PathSearch("spec.clusterVersion", workflowResp, "").(string)
	exactTargetVersion := utils.PathSearch("spec.targetVersion", workflowResp, "").(string)
	if exactCurrentVersion == "" || exactTargetVersion == "" {
		return diag.Errorf("unable to get clusterVersion or targetVersion in workflow response: %s", workflowResp)
	}

	// precheck
	createPreCheckResp, err := createPreCheck(client, clusterID, exactCurrentVersion, exactTargetVersion)
	if err != nil {
		return diag.Errorf("error creating CCE Autopilot cluster precheck: %s", err)
	}
	preCheckTaskId := utils.PathSearch("metadata.uid", createPreCheckResp, "").(string)
	err = clusterPreCheckForStateCompleted(ctx, client, clusterID, preCheckTaskId, d.Timeout(schema.TimeoutCreate))
	if err != nil {
		return diag.Errorf("error waiting for cluster precheck to complete: %s", err)
	}

	// cluster snapshot
	if d.Get("is_snapshot").(bool) {
		createSnapshotResp, err := createSnapshot(client, clusterID)
		if err != nil {
			return diag.Errorf("error creating CCE Autopilot cluster snapshot: %s", err)
		}
		snapTaskId := utils.PathSearch("uid", createSnapshotResp, "").(string)
		err = clusterSnapshotWaitingForStateCompleted(ctx, client, clusterID, snapTaskId, d.Timeout(schema.TimeoutCreate))
		if err != nil {
			return diag.Errorf("error waiting for cluster snapshot task to complete: %s", err)
		}
	}

	var (
		createAutopilotClusterUpgradeHttpUrl = "autopilot/v3/projects/{project_id}/clusters/{cluster_id}/operation/upgrade"
	)

	createAutopilotClusterUpgradePath := client.Endpoint + createAutopilotClusterUpgradeHttpUrl
	createAutopilotClusterUpgradePath = strings.ReplaceAll(createAutopilotClusterUpgradePath, "{project_id}", client.ProjectID)
	createAutopilotClusterUpgradePath = strings.ReplaceAll(createAutopilotClusterUpgradePath, "{cluster_id}", clusterID)

	createAutopilotClusterUpgradeOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	createOpts, err := buildAutopilotClusterUpgradeCreateOpts(d, exactTargetVersion)
	if err != nil {
		return nil
	}

	createAutopilotClusterUpgradeOpt.JSONBody = utils.RemoveNil(createOpts)
	createAutopilotClusterUpgradeResp, err := client.Request("POST",
		createAutopilotClusterUpgradePath, &createAutopilotClusterUpgradeOpt)
	if err != nil {
		return diag.Errorf("error upgrading CCE Autopilot cluster: %s", err)
	}

	createAutopilotClusterUpgradeRespBody, err := utils.FlattenResponse(createAutopilotClusterUpgradeResp)
	if err != nil {
		return diag.FromErr(err)
	}

	taskID := utils.PathSearch("metadata.uid", createAutopilotClusterUpgradeRespBody, "").(string)
	if taskID == "" {
		return diag.Errorf("error upgrading CCE Autopilot cluster: task_id is not found in API response")
	}

	d.SetId(taskID)

	err = autopilotClusterUpgradeWaitingForStateCompleted(ctx, client, clusterID, taskID, d.Timeout(schema.TimeoutCreate))
	if err != nil {
		return diag.Errorf("error waiting for upgrading cluster task to complete: %s", err)
	}

	// postcheck
	if d.Get("is_postcheck").(bool) {
		err = checkoutAfterUpgrade(client, clusterID, exactCurrentVersion, exactTargetVersion)
		if err != nil {
			return diag.FromErr(err)
		}
	}

	return resourceAutopilotClusterUpgradeRead(ctx, d, meta)
}

func createWorkflow(client *golangsdk.ServiceClient, clusterID, currentVersion, targetVersion string) (interface{}, error) {
	workflowHttpUrl := "autopilot/v3/projects/{project_id}/clusters/{cluster_id}/operation/upgradeworkflows"
	workflowPath := client.Endpoint + workflowHttpUrl
	workflowPath = strings.ReplaceAll(workflowPath, "{project_id}", client.ProjectID)
	workflowPath = strings.ReplaceAll(workflowPath, "{cluster_id}", clusterID)

	workflowOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		JSONBody: map[string]interface{}{
			"kind":       "WorkFlowTask",
			"apiVersion": "v3",
			"spec": map[string]interface{}{
				"clusterID":      clusterID,
				"clusterVersion": currentVersion,
				"targetVersion":  targetVersion,
			},
		},
	}
	workflowResp, err := client.Request("POST", workflowPath, &workflowOpt)
	if err != nil {
		return nil, fmt.Errorf("error preccheck CCE Autopilot cluster: %s", err)
	}

	return utils.FlattenResponse(workflowResp)
}

func createPreCheck(client *golangsdk.ServiceClient, clusterID, currentVersion, targetVersion string) (interface{}, error) {
	preCheckHttpUrl := "autopilot/v3/projects/{project_id}/clusters/{cluster_id}/operation/precheck"
	preCheckPath := client.Endpoint + preCheckHttpUrl
	preCheckPath = strings.ReplaceAll(preCheckPath, "{project_id}", client.ProjectID)
	preCheckPath = strings.ReplaceAll(preCheckPath, "{cluster_id}", clusterID)

	preCheckOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		JSONBody: map[string]interface{}{
			"kind":       "PreCheckTask",
			"apiVersion": "v3",
			"spec": map[string]interface{}{
				"clusterID":      clusterID,
				"clusterVersion": currentVersion,
				"targetVersion":  targetVersion,
			},
		},
	}
	preCheckResp, err := client.Request("POST", preCheckPath, &preCheckOpt)
	if err != nil {
		return nil, fmt.Errorf("error precheck CCE Autopilot cluster: %s", err)
	}

	return utils.FlattenResponse(preCheckResp)
}

func clusterPreCheckForStateCompleted(ctx context.Context, client *golangsdk.ServiceClient,
	clusterID, preCheckTaskId string, t time.Duration) error {
	stateConf := &resource.StateChangeConf{
		Pending:      []string{"PENDING"},
		Target:       []string{"COMPLETED"},
		Refresh:      refreshClusterPreCheckState(client, clusterID, preCheckTaskId),
		Timeout:      t,
		Delay:        10 * time.Second,
		PollInterval: 10 * time.Second,
	}
	_, err := stateConf.WaitForStateContext(ctx)
	return err
}

func refreshClusterPreCheckState(client *golangsdk.ServiceClient, clusterID, preCheckTaskId string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		preCheckTaskHttpUrl := "autopilot/v3/projects/{project_id}/clusters/{cluster_id}/operation/precheck/tasks/{task_id}"
		preCheckTaskPath := client.Endpoint + preCheckTaskHttpUrl
		preCheckTaskPath = strings.ReplaceAll(preCheckTaskPath, "{project_id}", client.ProjectID)
		preCheckTaskPath = strings.ReplaceAll(preCheckTaskPath, "{cluster_id}", clusterID)
		preCheckTaskPath = strings.ReplaceAll(preCheckTaskPath, "{task_id}", preCheckTaskId)

		preCheckTaskOpt := golangsdk.RequestOpts{
			KeepResponseBody: true,
		}
		preCheckTaskResp, err := client.Request("GET", preCheckTaskPath, &preCheckTaskOpt)
		if err != nil {
			return nil, "ERROR", err
		}

		preCheckTaskRespBody, err := utils.FlattenResponse(preCheckTaskResp)
		if err != nil {
			return nil, "ERROR", err
		}
		status := utils.PathSearch("status.phase", preCheckTaskRespBody, "").(string)

		targetStatus := []string{
			"Success",
		}
		if utils.StrSliceContains(targetStatus, status) {
			return preCheckTaskRespBody, "COMPLETED", nil
		}

		unexpectedStatus := []string{
			"Failed", "Error",
		}
		if utils.StrSliceContains(unexpectedStatus, status) {
			message := utils.PathSearch("status.message", preCheckTaskRespBody, "").(string)
			return preCheckTaskRespBody, status, errors.New(message)
		}
		return preCheckTaskRespBody, "PENDING", nil
	}
}

func clusterSnapshotWaitingForStateCompleted(ctx context.Context, client *golangsdk.ServiceClient,
	clusterID, snapTaskId string, t time.Duration) error {
	stateConf := &resource.StateChangeConf{
		Pending:      []string{"PENDING"},
		Target:       []string{"COMPLETED"},
		Refresh:      refreshClusterSnapshotState(client, clusterID, snapTaskId),
		Timeout:      t,
		Delay:        10 * time.Second,
		PollInterval: 10 * time.Second,
	}
	_, err := stateConf.WaitForStateContext(ctx)
	return err
}

func refreshClusterSnapshotState(client *golangsdk.ServiceClient, clusterID, snapTaskId string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		snapshotTaskHttpUrl := "autopilot/v3.1/projects/{project_id}/clusters/{cluster_id}/operation/snapshot/tasks"
		snapshotTaskPath := client.Endpoint + snapshotTaskHttpUrl
		snapshotTaskPath = strings.ReplaceAll(snapshotTaskPath, "{project_id}", client.ProjectID)
		snapshotTaskPath = strings.ReplaceAll(snapshotTaskPath, "{cluster_id}", clusterID)

		snapshotTaskOpt := golangsdk.RequestOpts{
			KeepResponseBody: true,
		}
		snapshotTaskResp, err := client.Request("GET", snapshotTaskPath, &snapshotTaskOpt)
		if err != nil {
			return nil, "ERROR", err
		}

		snapshotTaskRespBody, err := utils.FlattenResponse(snapshotTaskResp)
		if err != nil {
			return nil, "ERROR", err
		}
		expression := fmt.Sprintf("items[?metadata.uid=='%s']|[0].status.phase", snapTaskId)
		status := utils.PathSearch(expression, snapshotTaskRespBody, "").(string)

		targetStatus := []string{
			"Success",
		}
		if utils.StrSliceContains(targetStatus, status) {
			return snapshotTaskRespBody, "COMPLETED", nil
		}

		unexpectedStatus := []string{
			"Failed",
		}
		if utils.StrSliceContains(unexpectedStatus, status) {
			return snapshotTaskRespBody, status, nil
		}
		return snapshotTaskRespBody, "PENDING", nil
	}
}

func createSnapshot(client *golangsdk.ServiceClient, clusterID string) (interface{}, error) {
	createSnapshotHttpUrl := "autopilot/v3.1/projects/{project_id}/clusters/{cluster_id}/operation/snapshot"
	createSnapshotPath := client.Endpoint + createSnapshotHttpUrl
	createSnapshotPath = strings.ReplaceAll(createSnapshotPath, "{project_id}", client.ProjectID)
	createSnapshotPath = strings.ReplaceAll(createSnapshotPath, "{cluster_id}", clusterID)

	createSnapshotOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}
	createSnapshotResp, err := client.Request("POST", createSnapshotPath, &createSnapshotOpt)
	if err != nil {
		return nil, err
	}

	return utils.FlattenResponse(createSnapshotResp)
}

func checkoutAfterUpgrade(client *golangsdk.ServiceClient, clusterID, currentVersion, targetVersion string) error {
	postCheckHttpUrl := "autopilot/v3/projects/{project_id}/clusters/{cluster_id}/operation/postcheck"
	postCheckPath := client.Endpoint + postCheckHttpUrl
	postCheckPath = strings.ReplaceAll(postCheckPath, "{project_id}", client.ProjectID)
	postCheckPath = strings.ReplaceAll(postCheckPath, "{cluster_id}", clusterID)

	postCheckOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		JSONBody: map[string]interface{}{
			"kind":       "PostCheckTask",
			"apiVersion": "v3",
			"spec": map[string]interface{}{
				"clusterID":      clusterID,
				"clusterVersion": currentVersion,
				"targetVersion":  targetVersion,
			},
		},
	}
	postCheckResp, err := client.Request("POST", postCheckPath, &postCheckOpt)
	if err != nil {
		return fmt.Errorf("error confirmation after CCE Autopilot cluster upgrade: %s", err)
	}

	postCheckRespBody, err := utils.FlattenResponse(postCheckResp)
	if err != nil {
		return err
	}

	status := utils.PathSearch("status.phase", postCheckRespBody, "").(string)
	if status != "Success" {
		return fmt.Errorf("error confirmation after CCE Autopilot cluster upgrade: %s", status)
	}

	return nil
}

func autopilotClusterUpgradeWaitingForStateCompleted(ctx context.Context, client *golangsdk.ServiceClient,
	clusterID, taskID string, t time.Duration) error {
	stateConf := &resource.StateChangeConf{
		Pending: []string{"PENDING"},
		Target:  []string{"COMPLETED"},
		Refresh: func() (interface{}, string, error) {
			autopilotClusterUpgradeWaitingRespBody, err := getAutopilotClusterUpgradeDetail(client, clusterID, taskID)
			if err != nil {
				return nil, "ERROR", err
			}
			status := utils.PathSearch(`status.phase`, autopilotClusterUpgradeWaitingRespBody, "").(string)

			targetStatus := []string{
				"Success",
			}
			if utils.StrSliceContains(targetStatus, status) {
				return autopilotClusterUpgradeWaitingRespBody, "COMPLETED", nil
			}

			unexpectedStatus := []string{
				"Failed",
			}
			if utils.StrSliceContains(unexpectedStatus, status) {
				return autopilotClusterUpgradeWaitingRespBody, status, nil
			}

			return autopilotClusterUpgradeWaitingRespBody, "PENDING", nil
		},
		Timeout:      t,
		Delay:        10 * time.Second,
		PollInterval: 10 * time.Second,
	}
	_, err := stateConf.WaitForStateContext(ctx)
	return err
}

func getAutopilotClusterUpgradeDetail(client *golangsdk.ServiceClient, clusterID, taskID string) (interface{}, error) {
	getUpgradeDetailHttpUrl := "autopilot/v3/projects/{project_id}/clusters/{cluster_id}/operation/upgrade/tasks/{task_id}"

	getUpgradeDetailPath := client.Endpoint + getUpgradeDetailHttpUrl
	getUpgradeDetailPath = strings.ReplaceAll(getUpgradeDetailPath, "{project_id}", client.ProjectID)
	getUpgradeDetailPath = strings.ReplaceAll(getUpgradeDetailPath, "{cluster_id}", clusterID)
	getUpgradeDetailPath = strings.ReplaceAll(getUpgradeDetailPath, "{task_id}", taskID)

	getUpgradeDetailOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}
	getUpgradeDetailResp, err := client.Request("GET", getUpgradeDetailPath, &getUpgradeDetailOpt)
	if err != nil {
		return nil, err
	}

	return utils.FlattenResponse(getUpgradeDetailResp)
}

func resourceAutopilotClusterUpgradeRead(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return nil
}

func resourceAutopilotClusterUpgradeUpdate(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return nil
}

func resourceAutopilotClusterUpgradeDelete(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	errorMsg := "Deleting cluster upgrade resource is not supported. The cluster upgrade resource is only removed from the state."
	return diag.Diagnostics{
		diag.Diagnostic{
			Severity: diag.Warning,
			Summary:  errorMsg,
		},
	}
}
