package dws

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

var logicalClusterPlanNonUpdatableParams = []string{
	"cluster_id",
	"plan_type",
	"logical_cluster_name",
	"user",
	"node_num",
	"main_logical_cluster",
	"start_time",
	"end_time",
}

// @API DWS POST /v1/{project_id}/clusters/{cluster_id}/logical-cluster-plans
// @API DWS GET /v1/{project_id}/clusters/{cluster_id}/logical-cluster-plans
// @API DWS PUT /v1/{project_id}/clusters/{cluster_id}/logical-cluster-plans/{plan_id}
// @API DWS POST /v1/{project_id}/clusters/{cluster_id}/logical-cluster-plans/{plan_id}/enable
// @API DWS POST /v1/{project_id}/clusters/{cluster_id}/logical-cluster-plans/{plan_id}/disable
// @API DWS DELETE /v1/{project_id}/clusters/{cluster_id}/logical-cluster-plans/{plan_id}
func ResourceLogicalClusterPlan() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceLogicalClusterPlanCreate,
		ReadContext:   resourceLogicalClusterPlanRead,
		UpdateContext: resourceLogicalClusterPlanUpdate,
		DeleteContext: resourceLogicalClusterPlanDelete,

		Importer: &schema.ResourceImporter{
			StateContext: resourceLogicalClusterPlanImportState,
		},

		CustomizeDiff: config.FlexibleForceNew(logicalClusterPlanNonUpdatableParams),

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(30 * time.Minute),
			Update: schema.DefaultTimeout(30 * time.Minute),
			Delete: schema.DefaultTimeout(30 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: `The region where the logical cluster plan is located.`,
			},

			// Required parameters.
			"cluster_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The ID of the cluster to which the logical cluster plan belongs.`,
			},
			"plan_type": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The plan type.`,
			},
			"actions": {
				Type:        schema.TypeSet,
				Required:    true,
				Elem:        logicalClusterPlanActionSchema(),
				Description: `The list of logical cluster plan actions.`,
			},

			// Optional parameters.
			"logical_cluster_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The logical cluster name.`,
			},
			"user": {
				Type:          schema.TypeString,
				Optional:      true,
				ConflictsWith: []string{"main_logical_cluster"},
				Description:   `The user bound to the logical cluster.`,
			},
			"node_num": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: `The number of logical cluster nodes.`,
			},
			"main_logical_cluster": {
				Type:          schema.TypeString,
				Optional:      true,
				Computed:      true,
				ConflictsWith: []string{"user"},
				Description:   `The main logical cluster bound to the logical cluster.`,
			},
			"start_time": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The start time of the logical cluster plan, in Unix timestamp format.`,
			},
			"end_time": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The end time of the logical cluster plan, in Unix timestamp format.`,
			},
			"enabled": {
				Type:        schema.TypeBool,
				Optional:    true,
				Computed:    true,
				Description: `Whether the logical cluster plan is enabled.`,
			},

			// Attributes.
			"status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The status of the logical cluster plan.`,
			},

			// Internal parameters.
			"enable_force_new": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"true", "false"}, false),
				Description: utils.SchemaDesc(
					`Whether to allow parameters that do not support changes to have their change-triggered behavior set to 'ForceNew'.`,
					utils.SchemaDescInput{
						Internal: true,
					}),
			},
		},
	}
}

func logicalClusterPlanActionSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			// Required parameters.
			"type": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The action type.`,
			},
			"strategy": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The strategy expression or timestamp for the action.`,
			},

			// Attributes.
			"id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The ID of the action.`,
			},
		},
	}
}

func buildLogicalClusterPlanActions(items []interface{}) []map[string]interface{} {
	if len(items) < 1 {
		return nil
	}

	result := make([]map[string]interface{}, 0, len(items))
	for _, item := range items {
		result = append(result, map[string]interface{}{
			// It is not required to create a new action, but it is required to update the existing action.
			"id":       utils.ValueIgnoreEmpty(utils.PathSearch("id", item, nil)),
			"type":     utils.PathSearch("type", item, nil),
			"strategy": utils.ValueIgnoreEmpty(utils.PathSearch("strategy", item, nil)),
		})
	}

	return result
}

func buildLogicalClusterPlanBody(d *schema.ResourceData) map[string]interface{} {
	body := map[string]interface{}{
		"plan_type":            d.Get("plan_type"),
		"actions":              buildLogicalClusterPlanActions(d.Get("actions").(*schema.Set).List()),
		"logical_cluster_name": utils.ValueIgnoreEmpty(d.Get("logical_cluster_name")),
		"user":                 utils.ValueIgnoreEmpty(d.Get("user")),
		"node_num":             utils.ValueIgnoreEmpty(d.Get("node_num")),
		"main_logical_cluster": utils.ValueIgnoreEmpty(d.Get("main_logical_cluster")),
		"start_time":           utils.ValueIgnoreEmpty(d.Get("start_time")),
		"end_time":             utils.ValueIgnoreEmpty(d.Get("end_time")),
	}
	return body
}

func operateLogicalClusterPlan(client *golangsdk.ServiceClient, clusterId, planId string, action string) error {
	httpUrl := "v1/{project_id}/clusters/{cluster_id}/logical-cluster-plans/{plan_id}/{action}"

	operationPath := client.Endpoint + httpUrl
	operationPath = strings.ReplaceAll(operationPath, "{project_id}", client.ProjectID)
	operationPath = strings.ReplaceAll(operationPath, "{cluster_id}", clusterId)
	operationPath = strings.ReplaceAll(operationPath, "{plan_id}", planId)
	operationPath = strings.ReplaceAll(operationPath, "{action}", action)

	opt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json",
		},
	}
	_, err := client.Request("POST", operationPath, &opt)
	return err
}

func logicalClusterPlanStatusRefreshFunc(client *golangsdk.ServiceClient, clusterId, planId string, isEnabled bool) retry.StateRefreshFunc {
	return func() (interface{}, string, error) {
		plan, err := GetLogicalClusterPlanById(client, clusterId, planId)
		if err != nil {
			return nil, "", err
		}
		status := utils.PathSearch("status", plan, "").(string)

		if status == "failed" {
			return plan, "ERROR", fmt.Errorf("unexpected status (%s)", status)
		}

		completedStatuses := []string{"waiting", "finished"}
		if !isEnabled {
			completedStatuses = []string{"disabled"}
		}

		if utils.StrSliceContains(completedStatuses, status) {
			return plan, "COMPLETED", nil
		}
		return plan, "PENDING", nil
	}
}

func waitForLogicalClusterPlanStatus(ctx context.Context, client *golangsdk.ServiceClient, clusterId, planId string,
	isEnabled bool, timeout time.Duration) error {
	stateConf := &retry.StateChangeConf{
		Pending:      []string{"PENDING"},
		Target:       []string{"COMPLETED"},
		Refresh:      logicalClusterPlanStatusRefreshFunc(client, clusterId, planId, isEnabled),
		Timeout:      timeout,
		Delay:        5 * time.Second,
		PollInterval: 10 * time.Second,
	}
	_, err := stateConf.WaitForStateContext(ctx)
	return err
}

func createLogicalClusterPlan(client *golangsdk.ServiceClient, clusterId string, d *schema.ResourceData) (string, error) {
	httpUrl := "v1/{project_id}/clusters/{cluster_id}/logical-cluster-plans"

	createPath := client.Endpoint + httpUrl
	createPath = strings.ReplaceAll(createPath, "{project_id}", client.ProjectID)
	createPath = strings.ReplaceAll(createPath, "{cluster_id}", clusterId)
	opt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json",
		},
		JSONBody: utils.RemoveNil(buildLogicalClusterPlanBody(d)),
	}

	resp, err := client.Request("POST", createPath, &opt)
	if err != nil {
		return "", err
	}

	respBody, err := utils.FlattenResponse(resp)
	if err != nil {
		return "", err
	}

	planId := utils.PathSearch("plan_id", respBody, "").(string)
	return planId, nil
}

func resourceLogicalClusterPlanCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg       = meta.(*config.Config)
		region    = cfg.GetRegion(d)
		clusterId = d.Get("cluster_id").(string)
	)

	client, err := cfg.NewServiceClient("dws", region)
	if err != nil {
		return diag.Errorf("error creating DWS client: %s", err)
	}

	planId, err := createLogicalClusterPlan(client, clusterId, d)
	if err != nil {
		return diag.Errorf("error creating logical cluster plan: %s", err)
	}
	d.SetId(planId)

	// If the logical cluster plan is not enabled, disable it.
	isEnabled := utils.GetNestedObjectFromRawConfig(d.GetRawConfig(), "enabled")
	if isEnabled != nil && !isEnabled.(bool) {
		if err = operateLogicalClusterPlan(client, clusterId, d.Id(), "disable"); err != nil {
			return diag.Errorf("error disabling logical cluster plan: %s", err)
		}
		if err = waitForLogicalClusterPlanStatus(ctx, client, clusterId, planId,
			false, d.Timeout(schema.TimeoutCreate)); err != nil {
			return diag.Errorf("error waiting for logical cluster plan (%s) to be disabled: %s", planId, err)
		}
	}

	return resourceLogicalClusterPlanRead(ctx, d, meta)
}

func GetLogicalClusterPlanById(client *golangsdk.ServiceClient, clusterId, planId string) (interface{}, error) {
	httpUrl := "v1/{project_id}/clusters/{cluster_id}/logical-cluster-plans"

	getPath := client.Endpoint + httpUrl
	getPath = strings.ReplaceAll(getPath, "{project_id}", client.ProjectID)
	getPath = strings.ReplaceAll(getPath, "{cluster_id}", clusterId)
	opt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json",
		},
	}

	resp, err := client.Request("GET", getPath, &opt)
	if err != nil {
		return nil, err
	}

	respBody, err := utils.FlattenResponse(resp)
	if err != nil {
		return nil, err
	}

	plan := utils.PathSearch(fmt.Sprintf("plans[?id=='%s']|[0]", planId), respBody, nil)
	if plan == nil {
		return nil, golangsdk.ErrDefault404{
			ErrUnexpectedResponseCode: golangsdk.ErrUnexpectedResponseCode{
				Method: "GET",
				URL:    "/v1/{project_id}/clusters/{cluster_id}/logical-cluster-plans",
				Body:   []byte(fmt.Sprintf("the logical cluster plan (%s) is not found", planId)),
			},
		}
	}
	return plan, nil
}

func flattenLogicalClusterPlanActions(actions []interface{}) []map[string]interface{} {
	if len(actions) < 1 {
		return nil
	}

	result := make([]map[string]interface{}, 0, len(actions))
	for _, a := range actions {
		result = append(result, map[string]interface{}{
			"id":       utils.PathSearch("id", a, nil),
			"type":     utils.PathSearch("type", a, nil),
			"strategy": utils.ValueIgnoreEmpty(utils.PathSearch("strategy", a, nil)),
		})
	}
	return result
}

func resourceLogicalClusterPlanRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg       = meta.(*config.Config)
		region    = cfg.GetRegion(d)
		clusterId = d.Get("cluster_id").(string)
		planId    = d.Id()
	)

	client, err := cfg.NewServiceClient("dws", region)
	if err != nil {
		return diag.Errorf("error creating DWS client: %s", err)
	}

	plan, err := GetLogicalClusterPlanById(client, clusterId, planId)
	if err != nil {
		return common.CheckDeletedDiag(d, err, "error retrieving logical cluster plan")
	}

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("cluster_id", clusterId),
		d.Set("logical_cluster_name", utils.PathSearch("logical_cluster_name", plan, nil)),
		d.Set("user", utils.PathSearch("user", plan, nil)),
		d.Set("node_num", utils.PathSearch("node_num", plan, nil)),
		d.Set("main_logical_cluster", utils.ValueIgnoreEmpty(utils.PathSearch("main_logical_cluster", plan, d.Get("main_logical_cluster")))),
		d.Set("plan_type", utils.PathSearch("plan_type", plan, nil)),
		d.Set("start_time", fmt.Sprintf("%d", utils.ConvertTimeStrToNanoTimestamp(utils.PathSearch("start_time", plan, "").(string)))),
		d.Set("end_time", fmt.Sprintf("%d", utils.ConvertTimeStrToNanoTimestamp(utils.PathSearch("end_time", plan, "").(string)))),
		d.Set("status", utils.PathSearch("status", plan, "").(string)),
		d.Set("enabled", utils.PathSearch("status", plan, "").(string) != "disabled"),
		d.Set("actions", flattenLogicalClusterPlanActions(
			utils.PathSearch("actions", plan, make([]interface{}, 0)).([]interface{}))),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func buildLogicalClusterPlanUpdateBody(d *schema.ResourceData) map[string]interface{} {
	return map[string]interface{}{
		"actions": buildLogicalClusterPlanActions(d.Get("actions").(*schema.Set).List()),
	}
}

func updateLogicalClusterPlan(client *golangsdk.ServiceClient, clusterId, planId string, d *schema.ResourceData) error {
	httpUrl := "v1/{project_id}/clusters/{cluster_id}/logical-cluster-plans/{plan_id}"

	updatePath := client.Endpoint + httpUrl
	updatePath = strings.ReplaceAll(updatePath, "{project_id}", client.ProjectID)
	updatePath = strings.ReplaceAll(updatePath, "{cluster_id}", clusterId)
	updatePath = strings.ReplaceAll(updatePath, "{plan_id}", planId)
	opt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json",
		},
		JSONBody: utils.RemoveNil(buildLogicalClusterPlanUpdateBody(d)),
	}

	_, err := client.Request("PUT", updatePath, &opt)
	return err
}

func resourceLogicalClusterPlanUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg       = meta.(*config.Config)
		region    = cfg.GetRegion(d)
		clusterId = d.Get("cluster_id").(string)
		planId    = d.Id()
	)

	client, err := cfg.NewServiceClient("dws", region)
	if err != nil {
		return diag.Errorf("error creating DWS client: %s", err)
	}

	if d.HasChange("enabled") {
		changeAction := "enable"
		isEnabled := utils.GetNestedObjectFromRawConfig(d.GetRawConfig(), "enabled").(bool)
		if !isEnabled {
			changeAction = "disable"
		}

		if err = operateLogicalClusterPlan(client, clusterId, planId, changeAction); err != nil {
			return diag.Errorf("error operating logical cluster plan (%s): %s", planId, err)
		}

		if err = waitForLogicalClusterPlanStatus(ctx, client, clusterId, planId,
			isEnabled, d.Timeout(schema.TimeoutUpdate)); err != nil {
			return diag.Errorf("error waiting for logical cluster plan (%s) operation (%s) to be completed: %s",
				planId, changeAction, err)
		}
	}

	// Plan actions can be updated only when the logical cluster plan is enabled.
	// so, we should enable the logical cluster plan before updating the actions.
	if d.HasChange("actions") {
		if err = updateLogicalClusterPlan(client, clusterId, planId, d); err != nil {
			return diag.Errorf("error updating logical cluster plan: %s", err)
		}
	}

	return resourceLogicalClusterPlanRead(ctx, d, meta)
}

func deleteLogicalClusterPlan(client *golangsdk.ServiceClient, clusterId, planId string) error {
	httpUrl := "v1/{project_id}/clusters/{cluster_id}/logical-cluster-plans/{plan_id}"

	deletePath := client.Endpoint + httpUrl
	deletePath = strings.ReplaceAll(deletePath, "{project_id}", client.ProjectID)
	deletePath = strings.ReplaceAll(deletePath, "{cluster_id}", clusterId)
	deletePath = strings.ReplaceAll(deletePath, "{plan_id}", planId)
	opt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json",
		},
	}

	_, err := client.Request("DELETE", deletePath, &opt)
	return err
}

func resourceLogicalClusterPlanDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg       = meta.(*config.Config)
		region    = cfg.GetRegion(d)
		clusterId = d.Get("cluster_id").(string)
		planId    = d.Id()
	)

	client, err := cfg.NewServiceClient("dws", region)
	if err != nil {
		return diag.Errorf("error creating DWS client: %s", err)
	}

	if err = deleteLogicalClusterPlan(client, clusterId, planId); err != nil {
		return diag.Errorf("error deleting logical cluster plan: %s", err)
	}
	return nil
}

func resourceLogicalClusterPlanImportState(_ context.Context, d *schema.ResourceData, _ interface{}) ([]*schema.ResourceData, error) {
	parts := strings.Split(d.Id(), "/")
	if len(parts) != 2 {
		return nil, fmt.Errorf("invalid format specified for import ID, want '<cluster_id>/<plan_id>', but got '%s'", d.Id())
	}

	d.SetId(parts[1])
	return []*schema.ResourceData{d}, d.Set("cluster_id", parts[0])
}
