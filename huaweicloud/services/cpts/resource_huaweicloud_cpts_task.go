package cpts

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

const (
	clusterTypeShare   = "shared-cluster-internet"
	clusterTypePrivate = "private-cluster"

	runStatusFinished  = 2
	runStatusToRunning = 9

	operationEnable = "enable"
	operationStop   = "stop"
)

// @API CPTS PUT /v1/{project_id}/tasks/{task_id}
// @API CPTS DELETE /v1/{project_id}/tasks/{task_id}
// @API CPTS GET /v1/{project_id}/tasks/{task_id}
// @API CPTS POST /v1/{project_id}/tasks
// @API CPTS POST /v1/{project_id}/test-suites/{test_suite_id}/tasks/{task_id}
func ResourceTask() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceTaskCreate,
		UpdateContext: resourceTaskUpdate,
		DeleteContext: resourceTaskDelete,
		ReadContext:   resourceTaskRead,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Timeouts: &schema.ResourceTimeout{
			Update: schema.DefaultTimeout(30 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},

			"name": {
				Type:     schema.TypeString,
				Required: true,
			},

			"project_id": {
				Type:     schema.TypeInt,
				Required: true,
				ForceNew: true,
			},

			"benchmark_concurrency": {
				Type:     schema.TypeInt,
				Optional: true,
				Default:  100,
			},

			"cluster_id": {
				Type:     schema.TypeInt,
				Optional: true,
			},

			"operation": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{operationEnable, operationStop}, false),
			},

			"status": {
				Type:     schema.TypeInt,
				Computed: true,
			},
		},
	}
}

func buildCreateTaskBodyParams(d *schema.ResourceData) map[string]interface{} {
	return map[string]interface{}{
		"name":             d.Get("name"),
		"project_id":       d.Get("project_id"),
		"bench_concurrent": d.Get("benchmark_concurrency"),
	}
}

func resourceTaskCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		httpUrl = "v1/{project_id}/tasks"
		product = "cpts"
	)

	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating CPTS client: %s", err)
	}

	requestPath := client.Endpoint + httpUrl
	requestPath = strings.ReplaceAll(requestPath, "{project_id}", client.ProjectID)
	requestOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		JSONBody:         buildCreateTaskBodyParams(d),
	}

	resp, err := client.Request("POST", requestPath, &requestOpt)
	if err != nil {
		return diag.Errorf("error creating CPTS task: %s", err)
	}

	respBody, err := utils.FlattenResponse(resp)
	if err != nil {
		return diag.FromErr(err)
	}

	taskID := utils.PathSearch("task_id", respBody, nil)
	if taskID == nil {
		return diag.Errorf("error creating CPTS task: ID is not found in API response")
	}

	// The `task_id` field is a numeric type.
	d.SetId(strconv.Itoa(int(taskID.(float64))))
	return resourceTaskRead(ctx, d, meta)
}

func resourceTaskRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		product = "cpts"
	)

	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating CPTS client: %s", err)
	}

	respBody, err := GetTaskDetail(client, d.Id())
	if err != nil {
		return common.CheckDeletedDiag(d,
			common.ConvertExpected403ErrInto404Err(err, "code", "SVCSTG.CPTS.4032002"),
			"error retrieving CPTS task")
	}

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("name", utils.PathSearch("taskInfo.name", respBody, nil)),
		d.Set("project_id", utils.PathSearch("taskInfo.project_id", respBody, nil)),
		d.Set("benchmark_concurrency", utils.PathSearch("taskInfo.bench_concurrent", respBody, nil)),
		d.Set("status", utils.PathSearch("taskInfo.run_status", respBody, nil)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func buildUpdateTaskBodyParams(d *schema.ResourceData, idInt64 int64) map[string]interface{} {
	return map[string]interface{}{
		"id":               idInt64,
		"name":             d.Get("name"),
		"project_id":       d.Get("project_id"),
		"bench_concurrent": d.Get("benchmark_concurrency"),
	}
}

func updateTask(client *golangsdk.ServiceClient, d *schema.ResourceData) error {
	idInt64, err := strconv.ParseInt(d.Id(), 10, 32)
	if err != nil {
		return fmt.Errorf("the task ID must be integer: %s", err)
	}

	requestPath := client.Endpoint + "v1/{project_id}/tasks/{task_id}"
	requestPath = strings.ReplaceAll(requestPath, "{project_id}", client.ProjectID)
	requestPath = strings.ReplaceAll(requestPath, "{task_id}", d.Id())
	requestOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		JSONBody:         buildUpdateTaskBodyParams(d, idInt64),
	}

	_, err = client.Request("PUT", requestPath, &requestOpt)
	if err != nil {
		return fmt.Errorf("error updating CPTS task: %s", err)
	}

	return nil
}

func buildEnableTaskClusterType(d *schema.ResourceData) string {
	clusterId := d.Get("cluster_id").(int)
	if clusterId > 0 {
		return clusterTypePrivate
	}

	return clusterTypeShare
}

func buildEnableTaskClusterID(d *schema.ResourceData) int {
	clusterId := d.Get("cluster_id").(int)
	if clusterId > 0 {
		return clusterId
	}

	return 0
}

func buildEnableTaskNetworkInfo(d *schema.ResourceData) map[string]interface{} {
	clusterId := d.Get("cluster_id").(int)
	if clusterId > 0 {
		return make(map[string]interface{})
	}

	return map[string]interface{}{
		"network_type": "internet",
	}
}

func buildEnableTaskBodyParams(d *schema.ResourceData) map[string]interface{} {
	return map[string]interface{}{
		"cluster_type": buildEnableTaskClusterType(d),
		"cluster_id":   buildEnableTaskClusterID(d),
		"status":       runStatusToRunning,
		"network_info": buildEnableTaskNetworkInfo(d),
	}
}

func enableTask(client *golangsdk.ServiceClient, d *schema.ResourceData) error {
	requestPath := client.Endpoint + "v1/{project_id}/test-suites/{test_suite_id}/tasks/{task_id}"
	requestPath = strings.ReplaceAll(requestPath, "{project_id}", client.ProjectID)
	requestPath = strings.ReplaceAll(requestPath, "{test_suite_id}", fmt.Sprintf("%d", d.Get("project_id")))
	requestPath = strings.ReplaceAll(requestPath, "{task_id}", d.Id())
	requestOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		JSONBody:         buildEnableTaskBodyParams(d),
	}

	_, err := client.Request("POST", requestPath, &requestOpt)
	if err != nil {
		return fmt.Errorf("error enabling CPTS task: %s", err)
	}

	return nil
}

func buildDisableTaskBodyParams() map[string]interface{} {
	return map[string]interface{}{
		"cluster_type": clusterTypeShare,
		"cluster_id":   -1,
		"status":       runStatusFinished,
		"network_info": make(map[string]interface{}),
	}
}

func disableTask(client *golangsdk.ServiceClient, d *schema.ResourceData) error {
	requestPath := client.Endpoint + "v1/{project_id}/test-suites/{test_suite_id}/tasks/{task_id}"
	requestPath = strings.ReplaceAll(requestPath, "{project_id}", client.ProjectID)
	requestPath = strings.ReplaceAll(requestPath, "{test_suite_id}", fmt.Sprintf("%d", d.Get("project_id")))
	requestPath = strings.ReplaceAll(requestPath, "{task_id}", d.Id())
	requestOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		JSONBody:         buildDisableTaskBodyParams(),
	}

	_, err := client.Request("POST", requestPath, &requestOpt)
	if err != nil {
		return fmt.Errorf("error disabling CPTS task: %s", err)
	}

	return nil
}

func GetTaskDetail(client *golangsdk.ServiceClient, taskID string) (interface{}, error) {
	requestPath := client.Endpoint + "v1/{project_id}/tasks/{task_id}"
	requestPath = strings.ReplaceAll(requestPath, "{project_id}", client.ProjectID)
	requestPath = strings.ReplaceAll(requestPath, "{task_id}", taskID)
	requestOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	resp, err := client.Request("GET", requestPath, &requestOpt)
	if err != nil {
		return nil, err
	}

	return utils.FlattenResponse(resp)
}

func waitingForTaskFinished(ctx context.Context, client *golangsdk.ServiceClient, d *schema.ResourceData,
	timeout time.Duration) error {
	stateConf := &resource.StateChangeConf{
		Pending: []string{"PENDING"},
		Target:  []string{"COMPLETED"},
		Refresh: func() (interface{}, string, error) {
			respBody, err := GetTaskDetail(client, d.Id())
			if err != nil {
				return nil, "ERROR", err
			}

			runStatus := utils.PathSearch("taskInfo.run_status", respBody, "").(int)
			if runStatus == runStatusFinished {
				return respBody, "COMPLETED", nil
			}

			return respBody, "PENDING", nil
		},
		Timeout:      timeout,
		Delay:        10 * time.Second,
		PollInterval: 10 * time.Second,
	}

	_, err := stateConf.WaitForStateContext(ctx)
	if err != nil {
		return fmt.Errorf("error waiting for CPTS task (%s) to be finished: %s", d.Id(), err)
	}
	return nil
}

func resourceTaskUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		product = "cpts"
	)

	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating CPTS client: %s", err)
	}

	if d.HasChanges("benchmark_concurrency", "name") {
		if err := updateTask(client, d); err != nil {
			return diag.FromErr(err)
		}
	}

	// Enable or stop task
	if d.HasChange("operation") {
		var err error
		switch d.Get("operation").(string) {
		case operationEnable:
			err = enableTask(client, d)
		case operationStop:
			err = disableTask(client, d)
		}

		if err != nil {
			return diag.FromErr(err)
		}

		if err := waitingForTaskFinished(ctx, client, d, d.Timeout(schema.TimeoutUpdate)); err != nil {
			return diag.FromErr(err)
		}
	}

	return resourceTaskRead(ctx, d, meta)
}

func resourceTaskDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		httpUrl = "v1/{project_id}/tasks/{task_id}"
		product = "cpts"
	)

	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating CPTS client: %s", err)
	}

	requestPath := client.Endpoint + httpUrl
	requestPath = strings.ReplaceAll(requestPath, "{project_id}", client.ProjectID)
	requestPath = strings.ReplaceAll(requestPath, "{task_id}", d.Id())
	requestOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	_, err = client.Request("DELETE", requestPath, &requestOpt)
	if err != nil {
		return common.CheckDeletedDiag(d,
			common.ConvertExpected403ErrInto404Err(err, "code", "SVCSTG.CPTS.4032002"),
			"error deleting CPTS task")
	}

	return nil
}
