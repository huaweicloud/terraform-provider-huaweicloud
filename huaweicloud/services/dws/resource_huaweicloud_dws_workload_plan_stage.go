package dws

import (
	"context"
	"errors"
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

var defaultErrMsg = "the resource is illegal, please contact technical support"

// @API DWS POST /v2/{project_id}/clusters/{cluster_id}/workload/plans/{plan_id}/stages
// @API DWS GET /v2/{project_id}/clusters/{cluster_id}/workload/plans/{plan_id}
// @API DWS DELETE /v2/{project_id}/clusters/{cluster_id}/workload/plans/{plan_id}/stages/{stage_id}
// @API DWS GET /v2/{project_id}/clusters/{cluster_id}/workload/plans/{plan_id}/stages/{stage_id}
func ResourceWorkLoadPlanStage() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceWorkLoadPlanStageCreate,
		ReadContext:   resourceWorkLoadPlanStageRead,
		DeleteContext: resourceWorkLoadPlanStageDelete,

		Importer: &schema.ResourceImporter{
			StateContext: resourceWorkloadPlanStageImportState,
		},

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
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"start_time": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"end_time": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"queues": {
				Type:     schema.TypeList,
				Required: true,
				ForceNew: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:     schema.TypeString,
							Required: true,
							ForceNew: true,
						},
						"configuration": {
							Type:     schema.TypeSet,
							Required: true,
							ForceNew: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"resource_name": {
										Type:     schema.TypeString,
										Required: true,
										ForceNew: true,
									},
									"resource_value": {
										Type:     schema.TypeInt,
										Required: true,
										ForceNew: true,
									},
									"value_unit": {
										Type:     schema.TypeString,
										Optional: true,
										ForceNew: true,
									},
									"resource_description": {
										Type:     schema.TypeString,
										Optional: true,
										ForceNew: true,
									},
								},
							},
						},
					},
				},
			},
			"month": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"day": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
		},
	}
}

func resourceWorkLoadPlanStageCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		httpUrl = "v2/{project_id}/clusters/{cluster_id}/workload/plans/{plan_id}/stages"
		product = "dws"
	)

	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating DWS client: %s", err)
	}

	createPath := client.Endpoint + httpUrl
	createPath = strings.ReplaceAll(createPath, "{project_id}", client.ProjectID)
	createPath = strings.ReplaceAll(createPath, "{cluster_id}", d.Get("cluster_id").(string))
	createPath = strings.ReplaceAll(createPath, "{plan_id}", d.Get("plan_id").(string))
	createOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		JSONBody:         utils.RemoveNil(buildCreateWorkloadStageBodyParams(d)),
	}

	resp, err := client.Request("POST", createPath, &createOpt)
	if err != nil {
		return diag.Errorf("error creating DWS workload plan stage: %s", err)
	}
	respBody, err := utils.FlattenResponse(resp)
	if err != nil {
		return diag.FromErr(err)
	}
	respErrCode := utils.PathSearch("error_code", respBody, "")
	if respErrCode != "" {
		return diag.Errorf("error creating DWS workload plan stage: %s", utils.PathSearch("workload_res_str", respBody, defaultErrMsg))
	}

	// The create API does not return the stage ID, this method is needed to refresh stage ID.
	err = refreshWorkloadPlanStageID(client, d)
	if err != nil {
		return diag.FromErr(err)
	}

	return resourceWorkLoadPlanStageRead(ctx, d, meta)
}

func buildCreateWorkloadStageBodyParams(d *schema.ResourceData) map[string]interface{} {
	return map[string]interface{}{
		"workload_plan_stage": map[string]interface{}{
			"stage_name": d.Get("name"),
			"month":      d.Get("month"),
			"day":        d.Get("day"),
			"start_time": d.Get("start_time"),
			"end_time":   d.Get("end_time"),
			"queue_list": buildCreateQueueListBodyParam(d),
		},
	}
}

func buildCreateQueueListBodyParam(d *schema.ResourceData) []map[string]interface{} {
	queueList := d.Get("queues").([]interface{})
	if len(queueList) == 0 {
		return nil
	}
	params := make([]map[string]interface{}, 0, len(queueList))
	for _, queue := range queueList {
		queueObject := queue.(map[string]interface{})
		param := map[string]interface{}{
			"queue_name":      queueObject["name"],
			"queue_resources": buildQueueConfigParam(queueObject["configuration"]),
		}
		params = append(params, param)
	}
	return params
}

func buildQueueConfigParam(d interface{}) []map[string]interface{} {
	rawParams := d.(*schema.Set)
	if rawParams.Len() == 0 {
		return nil
	}
	params := make([]map[string]interface{}, 0, rawParams.Len())
	for _, rawParam := range rawParams.List() {
		raw := rawParam.(map[string]interface{})
		param := map[string]interface{}{
			"resource_name":        raw["resource_name"],
			"resource_value":       raw["resource_value"],
			"value_unit":           raw["value_unit"],
			"resource_description": raw["resource_description"],
		}
		params = append(params, param)
	}

	return params
}

func refreshWorkloadPlanStageID(client *golangsdk.ServiceClient, d *schema.ResourceData) error {
	httpUrl := "v2/{project_id}/clusters/{cluster_id}/workload/plans/{plan_id}"
	getPath := client.Endpoint + httpUrl
	getPath = strings.ReplaceAll(getPath, "{project_id}", client.ProjectID)
	getPath = strings.ReplaceAll(getPath, "{cluster_id}", d.Get("cluster_id").(string))
	getPath = strings.ReplaceAll(getPath, "{plan_id}", d.Get("plan_id").(string))
	getOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	resp, err := client.Request("GET", getPath, &getOpt)
	if err != nil {
		return fmt.Errorf("error querying DWS workload plan stage: %s", err)
	}

	respBody, err := utils.FlattenResponse(resp)
	if err != nil {
		return err
	}

	respErrCode := utils.PathSearch("error_code", respBody, "")
	if respErrCode != "" {
		return fmt.Errorf("error retrieving DWS workload plan stage: %s", utils.PathSearch("workload_res_str", respBody, defaultErrMsg))
	}

	// Get the stage by name form plan.
	name := d.Get("name")
	expression := fmt.Sprintf("workload_plan.stage_list[?stage_name=='%s']|[0]", name)
	stage := utils.PathSearch(expression, respBody, nil)
	if stage == nil {
		return fmt.Errorf("the DWS workload plan stage (%s) does not exist", name)
	}

	d.SetId(utils.PathSearch("stage_id", stage, "").(string))

	return nil
}

func parseWorkLoadPlanStageError(err error) error {
	if errors.As(err, &golangsdk.ErrDefault403{}) {
		return common.ConvertExpected400ErrInto404Err(err, "error_code", "DWS.0001")
	}

	if errors.As(err, &golangsdk.ErrDefault400{}) {
		return common.ConvertExpected403ErrInto404Err(err, "error_code", "DWS.0015")
	}
	return err
}

func resourceWorkLoadPlanStageRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		httpUrl = "v2/{project_id}/clusters/{cluster_id}/workload/plans/{plan_id}/stages/{stage_id}"
		product = "dws"
	)

	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating DWS client: %s", err)
	}
	getPath := client.Endpoint + httpUrl
	getPath = strings.ReplaceAll(getPath, "{project_id}", client.ProjectID)
	getPath = strings.ReplaceAll(getPath, "{cluster_id}", d.Get("cluster_id").(string))
	getPath = strings.ReplaceAll(getPath, "{plan_id}", d.Get("plan_id").(string))
	getPath = strings.ReplaceAll(getPath, "{stage_id}", d.Id())
	getOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	resp, err := client.Request("GET", getPath, &getOpt)
	if err != nil {
		return common.CheckDeletedDiag(d, parseWorkLoadPlanStageError(err), "DWS workload plan stage")
	}

	respBody, err := utils.FlattenResponse(resp)
	if err != nil {
		return diag.FromErr(err)
	}

	respErrCode := utils.PathSearch("error_code", respBody, "")
	if respErrCode != "" {
		// When stage ID is invalid, api return workload_res_code 110.
		workloadResCode := utils.PathSearch("workload_res_code", respBody, "").(float64)
		if workloadResCode == 110 {
			return common.CheckDeletedDiag(d, golangsdk.ErrDefault404{}, "DWS workload plan stage")
		}
		return diag.Errorf("error retrieving DWS workload plan stage: %s", utils.PathSearch("workload_res_str", respBody, defaultErrMsg))
	}

	stage := utils.PathSearch("workload_plan_stage", respBody, nil)
	if stage == nil {
		return common.CheckDeletedDiag(d, golangsdk.ErrDefault404{}, "DWS workload plan stage")
	}

	mErr := multierror.Append(
		d.Set("region", region),
		d.Set("name", utils.PathSearch("stage_name", stage, "")),
		d.Set("month", utils.PathSearch("month", stage, "")),
		d.Set("day", utils.PathSearch("day", stage, "")),
		d.Set("start_time", utils.PathSearch("start_time", stage, "")),
		d.Set("end_time", utils.PathSearch("end_time", stage, "")),
		d.Set("queues", flattenPoolList(stage)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenPoolList(stage interface{}) []interface{} {
	queueList := utils.PathSearch("queue_list", stage, make([]interface{}, 0)).([]interface{})
	if len(queueList) == 0 {
		return nil
	}
	rst := make([]interface{}, len(queueList))
	for i, queue := range queueList {
		rst[i] = map[string]interface{}{
			"name":          utils.PathSearch("queue_name", queue, ""),
			"configuration": flattenPoolConfiguration(queue),
		}
	}
	return rst
}

func flattenPoolConfiguration(queue interface{}) []interface{} {
	poolCon := utils.PathSearch("queue_resources", queue, make([]interface{}, 0)).([]interface{})
	if len(poolCon) == 0 {
		return nil
	}
	rst := make([]interface{}, len(poolCon))
	for i, queue := range poolCon {
		rst[i] = map[string]interface{}{
			"resource_name":        utils.PathSearch("resource_name", queue, ""),
			"resource_value":       utils.PathSearch("resource_value", queue, nil),
			"value_unit":           utils.PathSearch("value_unit", queue, ""),
			"resource_description": utils.PathSearch("resource_description", queue, ""),
		}
	}
	return rst
}
func resourceWorkLoadPlanStageDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		httpUrl = "v2/{project_id}/clusters/{cluster_id}/workload/plans/{plan_id}/stages/{stage_id}"
		product = "dws"
	)

	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating DWS client: %s", err)
	}
	delPath := client.Endpoint + httpUrl
	delPath = strings.ReplaceAll(delPath, "{project_id}", client.ProjectID)
	delPath = strings.ReplaceAll(delPath, "{cluster_id}", d.Get("cluster_id").(string))
	delPath = strings.ReplaceAll(delPath, "{plan_id}", d.Get("plan_id").(string))
	delPath = strings.ReplaceAll(delPath, "{stage_id}", d.Id())
	delStageOpts := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	delRes, err := client.Request("DELETE", delPath, &delStageOpts)
	if err != nil {
		return common.CheckDeletedDiag(d, parseWorkLoadPlanStageError(err), "DWS workload plan stage")
	}
	respBody, err := utils.FlattenResponse(delRes)
	if err != nil {
		return diag.FromErr(err)
	}
	// The workload_res_code 103 means plan not exists, 110 means stage not exists.
	respErrCode := utils.PathSearch("error_code", respBody, "")
	if respErrCode != "" {
		workloadResCode := utils.PathSearch("workload_res_code", respBody, "").(float64)
		if workloadResCode == 110 || workloadResCode == 103 {
			return common.CheckDeletedDiag(d, golangsdk.ErrDefault404{}, "DWS workload plan stage")
		}
		return diag.Errorf("error deleting DWS workload plan stage: %s", utils.PathSearch("workload_res_str", respBody, defaultErrMsg))
	}
	return nil
}

func resourceWorkloadPlanStageImportState(_ context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	parts := strings.SplitN(d.Id(), "/", 3)
	if len(parts) != 3 {
		return nil, fmt.Errorf("invalid format specified for import ID, must be <cluster_id>/<plan_id>/<name>")
	}

	d.Set("cluster_id", parts[0])
	d.Set("plan_id", parts[1])
	d.Set("name", parts[2])
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		product = "dws"
	)

	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return []*schema.ResourceData{d}, fmt.Errorf("error creating DWS client: %s", err)
	}
	return []*schema.ResourceData{d}, refreshWorkloadPlanStageID(client, d)
}
