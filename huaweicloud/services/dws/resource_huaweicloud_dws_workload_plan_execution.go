package dws

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API DWS POST /v2/{project_id}/clusters/{cluster_id}/workload/plans/{plan_id}/start
// @API DWS POST /v2/{project_id}/clusters/{cluster_id}/workload/plans/{plan_id}/stage-switch
// @API DWS GET /v2/{project_id}/clusters/{cluster_id}/workload/plans/{plan_id}
// @API DWS POST /v2/{project_id}/clusters/{cluster_id}/workload/plans/{plan_id}/stop
func ResourceWorkLoadPlanExecution() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceWorkLoadPlanExecutionCreate,
		ReadContext:   resourceWorkLoadPlanExecutionRead,
		UpdateContext: resourceWorkLoadPlanExecutionUpdate,
		DeleteContext: resourceWorkLoadPlanExecutionDelete,

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
				ForceNew: true,
			},
			"plan_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"stage_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
		},
	}
}

func resourceWorkLoadPlanExecutionCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		httpUrl = "v2/{project_id}/clusters/{cluster_id}/workload/plans/{plan_id}/start"
		product = "dws"
	)

	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating DWS client: %s", err)
	}

	planID := d.Get("plan_id").(string)
	createPath := client.Endpoint + httpUrl
	createPath = strings.ReplaceAll(createPath, "{project_id}", client.ProjectID)
	createPath = strings.ReplaceAll(createPath, "{cluster_id}", d.Get("cluster_id").(string))
	createPath = strings.ReplaceAll(createPath, "{plan_id}", planID)
	createOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	resp, err := client.Request("POST", createPath, &createOpt)
	if err != nil {
		return diag.Errorf("error starting DWS workload plan (%s): %s", planID, err)
	}

	respBody, err := utils.FlattenResponse(resp)
	if err != nil {
		return diag.FromErr(err)
	}

	// When calling the API, there are situations where the call fails but still returns a 200 status code.
	// The workload_res_code = 0 indicates that the plan has been successfully started,
	// and workload_res_code = 113 indicates that the plan to be started has been started.
	// Both of the above situations consider resource creation successful.
	resCode := utils.PathSearch("workload_res_code", respBody, float64(0)).(float64)
	if resCode != 0 && resCode != 113 {
		resMsg := utils.PathSearch("workload_res_str", respBody, "").(string)
		return diag.Errorf("error starting DWS workload plan (%s): error code: %v, error message: %s", planID, resCode, resMsg)
	}

	d.SetId(planID)

	if v, ok := d.GetOk("stage_id"); ok && v.(string) != "" {
		err = switchWorkLoadPlanStage(client, d)
		if err != nil {
			return diag.FromErr(err)
		}
	}

	return resourceWorkLoadPlanExecutionRead(ctx, d, meta)
}

func switchWorkLoadPlanStage(client *golangsdk.ServiceClient, d *schema.ResourceData) error {
	planID := d.Get("plan_id").(string)
	stageID := d.Get("stage_id").(string)

	httpUrl := "v2/{project_id}/clusters/{cluster_id}/workload/plans/{plan_id}/stage-switch"
	switchPath := client.Endpoint + httpUrl
	switchPath = strings.ReplaceAll(switchPath, "{project_id}", client.ProjectID)
	switchPath = strings.ReplaceAll(switchPath, "{cluster_id}", d.Get("cluster_id").(string))
	switchPath = strings.ReplaceAll(switchPath, "{plan_id}", planID)
	switchOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		JSONBody: map[string]interface{}{
			"stage_id": stageID,
		},
	}

	resp, err := client.Request("POST", switchPath, &switchOpt)
	if err != nil {
		return fmt.Errorf("error switching plan stage (%s) for DWS workload plan (%s): %s", stageID, planID, err)
	}

	respBody, err := utils.FlattenResponse(resp)
	if err != nil {
		return err
	}

	// When calling the API,
	// there is a situation where the workload plan switch plan stage failed but still returns a 200 status code.
	resCode := utils.PathSearch("workload_res_code", respBody, float64(0)).(float64)
	if resCode != 0 {
		resMsg := utils.PathSearch("workload_res_str", respBody, "").(string)
		return fmt.Errorf("error switching plan stage (%s) for DWS workload plan (%s): error code: %v, error message: %s",
			stageID, planID, resCode, resMsg)
	}

	return nil
}

func resourceWorkLoadPlanExecutionRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		httpUrl = "v2/{project_id}/clusters/{cluster_id}/workload/plans/{plan_id}"
		product = "dws"
	)

	getClient, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating DWS client: %s", err)
	}

	getPath := getClient.Endpoint + httpUrl
	getPath = strings.ReplaceAll(getPath, "{project_id}", getClient.ProjectID)
	getPath = strings.ReplaceAll(getPath, "{cluster_id}", d.Get("cluster_id").(string))
	getPath = strings.ReplaceAll(getPath, "{plan_id}", d.Id())
	getOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	getResp, err := getClient.Request("GET", getPath, &getOpt)
	if err != nil {
		return common.CheckDeletedDiag(d, parseWorkLoadPlanError(err), "DWS workload plan execution")
	}

	getRespBody, err := utils.FlattenResponse(getResp)
	if err != nil {
		return diag.FromErr(err)
	}

	// When calling the API, there is a situation where the plan ID does not exist but still returns a 200 status code.
	// When the workload plan is successfully started, the value of the status attribute of the workload plan is 1.
	resCode := utils.PathSearch("workload_res_code", getRespBody, float64(0)).(float64)
	status := utils.PathSearch("workload_plan.status", getRespBody, float64(0)).(float64)
	if resCode != 0 || status != 1 {
		return common.CheckDeletedDiag(d, golangsdk.ErrDefault404{}, "DWS workload plan execution")
	}

	jsonPaths := fmt.Sprintf("workload_plan.stage_list[?stage_id=='%s'] | [0].stage_id", d.Get("stage_id").(string))
	stageID := utils.PathSearch(jsonPaths, getRespBody, "").(string)
	mErr := multierror.Append(
		d.Set("region", region),
		d.Set("cluster_id", utils.PathSearch("workload_plan.cluster_id", getRespBody, "").(string)),
		d.Set("plan_id", utils.PathSearch("workload_plan.plan_id", getRespBody, "").(string)),
		d.Set("stage_id", stageID),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func resourceWorkLoadPlanExecutionUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		product = "dws"
	)

	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating DWS client: %s", err)
	}

	if d.HasChange("stage_id") {
		if v, ok := d.GetOk("stage_id"); ok && v.(string) != "" {
			err = switchWorkLoadPlanStage(client, d)
			if err != nil {
				return diag.FromErr(err)
			}
		}
	}

	return resourceWorkLoadPlanExecutionRead(ctx, d, meta)
}

func resourceWorkLoadPlanExecutionDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		httpUrl = "v2/{project_id}/clusters/{cluster_id}/workload/plans/{plan_id}/stop"
		product = "dws"
	)

	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating DWS client: %s", err)
	}

	id := d.Id()
	deletePath := client.Endpoint + httpUrl
	deletePath = strings.ReplaceAll(deletePath, "{project_id}", client.ProjectID)
	deletePath = strings.ReplaceAll(deletePath, "{cluster_id}", d.Get("cluster_id").(string))
	deletePath = strings.ReplaceAll(deletePath, "{plan_id}", id)
	deleteOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	resp, err := client.Request("POST", deletePath, &deleteOpt)
	if err != nil {
		return diag.Errorf("error stopping DWS workload plan (%s): %s", id, err)
	}

	respBody, err := utils.FlattenResponse(resp)
	if err != nil {
		return diag.FromErr(err)
	}

	// When calling the API, there are situations where the plan ID does not exist or the plan is not started
	// but still returns a 200 status code.
	// The workload_res_code = 0 indicates that the plan was successfully stopped,
	// and workload_res_code = 116 indicates that the plan to be stopped has already been stopped.
	// Both of the above situations consider the deletion of resource successful.
	resCode := utils.PathSearch("workload_res_code", respBody, float64(0)).(float64)
	if resCode != 0 && resCode != 116 {
		resMsg := utils.PathSearch("workload_res_str", respBody, "").(string)
		return diag.Errorf("error stopping DWS workload plan (%s): error code: %v, error message: %s", id, resCode, resMsg)
	}

	return nil
}
