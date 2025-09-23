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

// @API DWS POST /v2/{project_id}/clusters/{cluster_id}/workload/plans
// @API DWS GET /v2/{project_id}/clusters/{cluster_id}/workload/plans
// @API DWS GET /v2/{project_id}/clusters/{cluster_id}/workload/plans/{plan_id}
// @API DWS DELETE /v2/{project_id}/clusters/{cluster_id}/workload/plans/{plan_id}
func ResourceWorkLoadPlan() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceWorkLoadPlanCreate,
		ReadContext:   resourceWorkLoadPlanRead,
		DeleteContext: resourceWorkLoadPlanDelete,

		Importer: &schema.ResourceImporter{
			StateContext: resourceWorkLoadPlanImportState,
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
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"logical_cluster_name": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: utils.SchemaDesc("", utils.SchemaDescInput{Internal: true}),
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"current_stage_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"stages": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func resourceWorkLoadPlanCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		httpUrl = "v2/{project_id}/clusters/{cluster_id}/workload/plans"
		product = "dws"
	)

	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating DWS client: %s", err)
	}

	createPath := client.Endpoint + httpUrl
	createPath = strings.ReplaceAll(createPath, "{project_id}", client.ProjectID)
	createPath = strings.ReplaceAll(createPath, "{cluster_id}", d.Get("cluster_id").(string))
	createOpt := golangsdk.RequestOpts{
		MoreHeaders:      requestOpts.MoreHeaders,
		KeepResponseBody: true,
		JSONBody:         buildCreateWorkLoadPlanBodyParams(d),
	}

	_, err = client.Request("POST", createPath, &createOpt)
	if err != nil {
		return diag.Errorf("error creating DWS workload plan: %s", err)
	}

	// The create API does not return the plan ID, this method is needed to refresh plan ID.
	err = refreshWorkLoadPlanID(client, d)
	if err != nil {
		return diag.FromErr(err)
	}

	return resourceWorkLoadPlanRead(ctx, d, meta)
}

func buildCreateWorkLoadPlanBodyParams(d *schema.ResourceData) map[string]interface{} {
	return map[string]interface{}{
		"plan_name":            d.Get("name").(string),
		"logical_cluster_name": utils.ValueIgnoreEmpty(d.Get("logical_cluster_name").(string)),
	}
}

func refreshWorkLoadPlanID(client *golangsdk.ServiceClient, d *schema.ResourceData) error {
	httpUrl := "v2/{project_id}/clusters/{cluster_id}/workload/plans"
	listPath := client.Endpoint + httpUrl
	listPath = strings.ReplaceAll(listPath, "{project_id}", client.ProjectID)
	listPath = strings.ReplaceAll(listPath, "{cluster_id}", d.Get("cluster_id").(string))
	listOpt := golangsdk.RequestOpts{
		MoreHeaders:      requestOpts.MoreHeaders,
		KeepResponseBody: true,
	}

	listResp, err := client.Request("GET", listPath, &listOpt)
	if err != nil {
		return fmt.Errorf("error querying DWS workload plans: %s", err)
	}

	listRespBody, err := utils.FlattenResponse(listResp)
	if err != nil {
		return fmt.Errorf("%s", err)
	}

	name := d.Get("name")
	jsonPaths := fmt.Sprintf("plan_list[?plan_name=='%s']", name)
	plans := utils.PathSearch(jsonPaths, listRespBody, make([]interface{}, 0)).([]interface{})
	if len(plans) == 0 {
		return fmt.Errorf("the DWS workload plan (%s) does not exist", name)
	}

	id := utils.PathSearch("plan_id", plans[0], "")
	d.SetId(id.(string))

	return nil
}

// When the cluster ID does not exist, the API returns a 401 status code.
// When the cluster ID is illegal, the API returns a 400 status code.
// Both of the above situations need to be processed as 404 error codes.
func parseWorkLoadPlanError(err error) error {
	if errors.As(err, &golangsdk.ErrDefault400{}) {
		return common.ConvertExpected400ErrInto404Err(err, "error_code", "DWS.0001")
	}

	if errors.As(err, &golangsdk.ErrDefault401{}) {
		return common.ConvertExpected401ErrInto404Err(err, "error_code", "DWS.0047")
	}

	return err
}

func resourceWorkLoadPlanRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
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
		MoreHeaders:      requestOpts.MoreHeaders,
		KeepResponseBody: true,
	}

	getResp, err := getClient.Request("GET", getPath, &getOpt)
	if err != nil {
		return common.CheckDeletedDiag(d, parseWorkLoadPlanError(err), "DWS workload plan")
	}

	getRespBody, err := utils.FlattenResponse(getResp)
	if err != nil {
		return diag.FromErr(err)
	}

	// When calling the API, there is a situation where the plan ID does not exist but still returns a 200 status code.
	// So, it's necessary to check the errCode.
	errCode := utils.PathSearch("error_code", getRespBody, "").(string)
	plan := utils.PathSearch("workload_plan", getRespBody, nil)

	if errCode != "" || plan == nil {
		return common.CheckDeletedDiag(d, golangsdk.ErrDefault404{}, "DWS workload plan")
	}

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("cluster_id", utils.PathSearch("cluster_id", plan, nil)),
		d.Set("name", utils.PathSearch("plan_name", plan, nil)),
		d.Set("logical_cluster_name", utils.PathSearch("logical_cluster_name", plan, nil)),
		d.Set("status", convertWorkLoadPlanStatus(utils.PathSearch("status", plan, float64(0)).(float64))),
		d.Set("current_stage_name", utils.PathSearch("current_stage", plan, nil)),
		d.Set("stages", flattenWorkLoadPlanStages(plan)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenWorkLoadPlanStages(plan interface{}) []interface{} {
	curJson := utils.PathSearch("stage_list", plan, make([]interface{}, 0))
	stageList := curJson.([]interface{})
	result := make([]interface{}, 0, len(stageList))
	for _, stage := range stageList {
		result = append(result, map[string]interface{}{
			"id":   utils.PathSearch("stage_id", stage, nil),
			"name": utils.PathSearch("stage_name", stage, nil),
		})
	}

	return result
}

// The API actually returns 0 and 1, where 1 indicates the workload plan has been started.
func convertWorkLoadPlanStatus(statusInt float64) string {
	var status = "disabled"
	if statusInt == 1 {
		status = "enabled"
	}

	return status
}

func resourceWorkLoadPlanDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		httpUrl = "v2/{project_id}/clusters/{cluster_id}/workload/plans/{plan_id}"
		product = "dws"
	)

	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating DWS client: %s", err)
	}

	deletePath := client.Endpoint + httpUrl
	deletePath = strings.ReplaceAll(deletePath, "{project_id}", client.ProjectID)
	deletePath = strings.ReplaceAll(deletePath, "{cluster_id}", d.Get("cluster_id").(string))
	deletePath = strings.ReplaceAll(deletePath, "{plan_id}", d.Id())
	deleteOpt := golangsdk.RequestOpts{
		MoreHeaders:      requestOpts.MoreHeaders,
		KeepResponseBody: true,
	}

	resp, err := client.Request("DELETE", deletePath, &deleteOpt)
	if err != nil {
		return diag.Errorf("error deleting DWS workload plan: %s", err)
	}

	respBody, err := utils.FlattenResponse(resp)
	if err != nil {
		return diag.FromErr(err)
	}

	// When calling the API, there is a situation where deletion fails but still returns a 200 status code.
	errCode := utils.PathSearch("error_code", respBody, "").(string)
	if errCode != "" {
		errMsg := utils.PathSearch("workload_res_str", respBody, "").(string)
		return diag.Errorf("error deleting DWS workload plan: error code: %s, error message: %s", errCode, errMsg)
	}

	return nil
}

func resourceWorkLoadPlanImportState(_ context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	parts := strings.SplitN(d.Id(), "/", 2)
	if len(parts) != 2 {
		return nil, fmt.Errorf("invalid format specified for import id, must be <cluster_id>/<name>")
	}
	d.Set("cluster_id", parts[0])
	d.Set("name", parts[1])

	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		product = "dws"
	)

	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return []*schema.ResourceData{d}, fmt.Errorf("error creating DWS client: %s", err)
	}

	return []*schema.ResourceData{d}, refreshWorkLoadPlanID(client, d)
}
