// ---------------------------------------------------------------
// *** AUTO GENERATED CODE ***
// @Product ModelArts
// ---------------------------------------------------------------

package modelarts

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API ModelArts POST /v2/{project_id}/pools
// @API ModelArts DELETE /v2/{project_id}/pools/{id}
// @API ModelArts GET /v2/{project_id}/pools/{id}
// @API ModelArts PATCH /v2/{project_id}/pools/{id}
func ResourceModelartsResourcePool() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceModelartsResourcePoolCreate,
		UpdateContext: resourceModelartsResourcePoolUpdate,
		ReadContext:   resourceModelartsResourcePoolRead,
		DeleteContext: resourceModelartsResourcePoolDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(20 * time.Minute),
			Update: schema.DefaultTimeout(20 * time.Minute),
			Delete: schema.DefaultTimeout(20 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: `The name of the resource pool.`,
			},
			"scope": {
				Type:        schema.TypeList,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Required:    true,
				Description: `List of job types supported by the resource pool.`,
			},
			"resources": {
				Type:        schema.TypeList,
				Elem:        modelartsResourcePoolResourceFlavorSchema(),
				Required:    true,
				Description: `List of resource specifications in the resource pool.`,
			},
			"network_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: `The ModelArts network ID of the resource pool.`,
			},
			"workspace_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: `Workspace ID, which defaults to 0.`,
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `The description of the resource pool.`,
			},
		},
	}
}

func modelartsResourcePoolResourceFlavorSchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"flavor_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The resource flavor ID.`,
			},
			"count": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: `Number of resources of the corresponding flavors.`,
			},
			"azs": {
				Type:        schema.TypeList,
				Elem:        modelartsResourcePoolResourceFlavorAzsSchema(),
				Optional:    true,
				Computed:    true,
				Description: `AZs for resource pool nodes.`,
			},
		},
	}
	return &sc
}

func modelartsResourcePoolResourceFlavorAzsSchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"az": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `The AZ name.`,
			},
			"count": {
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
				Description: `Number of nodes.`,
			},
		},
	}
	return &sc
}

func resourceModelartsResourcePoolCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var (
		createResourcePoolHttpUrl = "v2/{project_id}/pools"
		createResourcePoolProduct = "modelarts"
	)
	createResourcePoolClient, err := cfg.NewServiceClient(createResourcePoolProduct, region)
	if err != nil {
		return diag.Errorf("error creating ModelArts client: %s", err)
	}

	createResourcePoolPath := createResourcePoolClient.Endpoint + createResourcePoolHttpUrl
	createResourcePoolPath = strings.ReplaceAll(createResourcePoolPath, "{project_id}", createResourcePoolClient.ProjectID)

	createResourcePoolOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		OkCodes: []int{
			200,
		},
		MoreHeaders: map[string]string{"Content-Type": "application/json"},
	}

	createResourcePoolOpt.JSONBody = utils.RemoveNil(buildCreateResourcePoolBodyParams(d))
	createResourcePoolResp, err := createResourcePoolClient.Request("POST", createResourcePoolPath, &createResourcePoolOpt)
	if err != nil {
		return diag.Errorf("error creating Modelarts resource pool: %s", err)
	}

	createResourcePoolRespBody, err := utils.FlattenResponse(createResourcePoolResp)
	if err != nil {
		return diag.FromErr(err)
	}

	id := utils.PathSearch("metadata.name", createResourcePoolRespBody, nil)
	if id == nil {
		return diag.Errorf("error creating Modelarts resource pool: ID is not found in API response")
	}
	d.SetId(id.(string))

	err = createResourcePoolWaitingForStateCompleted(ctx, d, meta, d.Timeout(schema.TimeoutCreate))
	if err != nil {
		return diag.Errorf("error waiting for the Modelarts resource pool (%s) creation to complete: %s", d.Id(), err)
	}
	return resourceModelartsResourcePoolRead(ctx, d, meta)
}

func buildCreateResourcePoolBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"apiVersion": "v2",
		"kind":       "Pool",
		"metadata": map[string]interface{}{
			"labels": map[string]interface{}{
				"os.modelarts/name":         d.Get("name"),
				"os.modelarts/workspace.id": utils.ValueIngoreEmpty(d.Get("workspace_id")),
			},
			"annotations": map[string]interface{}{
				"os.modelarts/description": utils.ValueIngoreEmpty(d.Get("description")),
			},
		},
		"spec": map[string]interface{}{
			"type":      "Dedicate",
			"scope":     d.Get("scope"),
			"resources": buildResourcePoolResourceFlavor(d.Get("resources")),
			"network": map[string]interface{}{
				"name": d.Get("network_id"),
			},
		},
	}
	return bodyParams
}

func createResourcePoolWaitingForStateCompleted(ctx context.Context, d *schema.ResourceData, meta interface{}, t time.Duration) error {
	stateConf := &resource.StateChangeConf{
		Pending: []string{"PENDING"},
		Target:  []string{"COMPLETED"},
		Refresh: func() (interface{}, string, error) {
			cfg := meta.(*config.Config)
			region := cfg.GetRegion(d)

			createResourcePoolWaitingRespBody, err := queryResourcePool(cfg, region, d)
			if err != nil {
				return nil, "ERROR", err
			}
			statusRaw := utils.PathSearch(`status.phase`, createResourcePoolWaitingRespBody, nil)
			if statusRaw == nil {
				return nil, "ERROR", fmt.Errorf("error parse %s from response body", `status.phase`)
			}

			status := fmt.Sprintf("%v", statusRaw)

			targetStatus := []string{
				"Running",
			}
			if utils.StrSliceContains(targetStatus, status) {
				return createResourcePoolWaitingRespBody, "COMPLETED", nil
			}

			pendingStatus := []string{
				"Creating",
				"Waiting",
			}
			if utils.StrSliceContains(pendingStatus, status) {
				return createResourcePoolWaitingRespBody, "PENDING", nil
			}

			return createResourcePoolWaitingRespBody, status, nil
		},
		Timeout:      t,
		Delay:        30 * time.Second,
		PollInterval: 20 * time.Second,
	}
	_, err := stateConf.WaitForStateContext(ctx)
	return err
}

func resourceModelartsResourcePoolRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var mErr *multierror.Error

	getModelartsResourcePoolRespBody, err := queryResourcePool(cfg, region, d)
	if err != nil {
		return common.CheckDeletedDiag(d, err, "error retrieving Modelarts resource pool")
	}

	mErr = multierror.Append(
		mErr,
		d.Set("region", region),
		d.Set("name", utils.PathSearch(`metadata.labels."os.modelarts/name"`, getModelartsResourcePoolRespBody, nil)),
		d.Set("workspace_id", utils.PathSearch(`metadata.labels."os.modelarts/workspace.id"`, getModelartsResourcePoolRespBody, nil)),
		d.Set("description", utils.PathSearch(`metadata.annotations."os.modelarts/description"`, getModelartsResourcePoolRespBody, nil)),
		d.Set("scope", utils.PathSearch("spec.scope", getModelartsResourcePoolRespBody, nil)),
		d.Set("resources", flattenGetResourcePoolResponseBodyResourceFlavor(getModelartsResourcePoolRespBody)),
		d.Set("network_id", utils.PathSearch("spec.network.name", getModelartsResourcePoolRespBody, nil)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func queryResourcePool(cfg *config.Config, region string, d *schema.ResourceData) (interface{}, error) {
	var (
		getModelartsResourcePoolHttpUrl = "v2/{project_id}/pools/{id}"
		getModelartsResourcePoolProduct = "modelarts"
	)
	getModelartsResourcePoolClient, err := cfg.NewServiceClient(getModelartsResourcePoolProduct, region)
	if err != nil {
		return nil, fmt.Errorf("error creating ModelArts client: %s", err)
	}

	getModelartsResourcePoolPath := getModelartsResourcePoolClient.Endpoint + getModelartsResourcePoolHttpUrl
	getModelartsResourcePoolPath = strings.ReplaceAll(getModelartsResourcePoolPath, "{project_id}", getModelartsResourcePoolClient.ProjectID)
	getModelartsResourcePoolPath = strings.ReplaceAll(getModelartsResourcePoolPath, "{id}", d.Id())

	getModelartsResourcePoolOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		OkCodes: []int{
			200,
		},
		MoreHeaders: map[string]string{"Content-Type": "application/json"},
	}

	getModelartsResourcePoolResp, err := getModelartsResourcePoolClient.Request("GET", getModelartsResourcePoolPath, &getModelartsResourcePoolOpt)

	if err != nil {
		return nil, err
	}

	getModelartsResourcePoolRespBody, err := utils.FlattenResponse(getModelartsResourcePoolResp)
	if err != nil {
		return nil, err
	}
	return getModelartsResourcePoolRespBody, nil
}

func flattenGetResourcePoolResponseBodyResourceFlavor(resp interface{}) []interface{} {
	if resp == nil {
		return nil
	}
	curJson := utils.PathSearch("spec.resources", resp, make([]interface{}, 0))
	curArray := curJson.([]interface{})
	rst := make([]interface{}, 0, len(curArray))
	for _, v := range curArray {
		rst = append(rst, map[string]interface{}{
			"flavor_id": utils.PathSearch("flavor", v, nil),
			"count":     utils.PathSearch("count", v, nil),
			"azs":       flattenResourcePoolsFlavorAzs(v),
		})
	}
	return rst
}

func flattenResourcePoolsFlavorAzs(resp interface{}) []interface{} {
	if resp == nil {
		return nil
	}
	curJson := utils.PathSearch("azs", resp, make([]interface{}, 0))
	curArray := curJson.([]interface{})
	rst := make([]interface{}, 0, len(curArray))
	for _, v := range curArray {
		rst = append(rst, map[string]interface{}{
			"az":    utils.PathSearch("az", v, nil),
			"count": utils.PathSearch("count", v, nil),
		})
	}
	return rst
}

func resourceModelartsResourcePoolUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	updateResourcePoolChanges := []string{
		"description",
		"scope",
		"resources",
	}

	if d.HasChanges(updateResourcePoolChanges...) {
		var (
			updateResourcePoolHttpUrl = "v2/{project_id}/pools/{id}"
			updateResourcePoolProduct = "modelarts"
		)
		updateResourcePoolClient, err := cfg.NewServiceClient(updateResourcePoolProduct, region)
		if err != nil {
			return diag.Errorf("error creating ModelArts client: %s", err)
		}

		updateResourcePoolPath := updateResourcePoolClient.Endpoint + updateResourcePoolHttpUrl
		updateResourcePoolPath = strings.ReplaceAll(updateResourcePoolPath, "{project_id}", updateResourcePoolClient.ProjectID)
		updateResourcePoolPath = strings.ReplaceAll(updateResourcePoolPath, "{id}", d.Id())

		updateResourcePoolOpt := golangsdk.RequestOpts{
			KeepResponseBody: true,
			OkCodes: []int{
				200,
			},
			MoreHeaders: map[string]string{"Content-Type": "application/merge-patch+json"},
		}

		updateResourcePoolOpt.JSONBody = utils.RemoveNil(buildUpdateResourcePoolBodyParams(d))
		_, err = updateResourcePoolClient.Request("PATCH", updateResourcePoolPath, &updateResourcePoolOpt)
		if err != nil {
			return diag.Errorf("error updating Modelarts resource pool: %s", err)
		}
		err = updateResourcePoolWaitingForStateCompleted(ctx, d, meta, d.Timeout(schema.TimeoutUpdate))
		if err != nil {
			return diag.Errorf("error waiting for the Modelarts resource pool (%s) update to complete: %s", d.Id(), err)
		}
	}
	return resourceModelartsResourcePoolRead(ctx, d, meta)
}

func buildUpdateResourcePoolBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"metadata": map[string]interface{}{
			"annotations": map[string]interface{}{
				"os.modelarts/description": utils.ValueIngoreEmpty(d.Get("description")),
			},
		},
		"spec": map[string]interface{}{
			"scope":     d.Get("scope"),
			"resources": buildResourcePoolResourceFlavor(d.Get("resources")),
		},
	}
	return bodyParams
}

func buildResourcePoolResourceFlavor(rawParams interface{}) []map[string]interface{} {
	if rawArray, ok := rawParams.([]interface{}); ok {
		if len(rawArray) == 0 {
			return nil
		}

		rst := make([]map[string]interface{}, len(rawArray))
		for i, v := range rawArray {
			if raw, ok := v.(map[string]interface{}); ok {
				rst[i] = map[string]interface{}{
					"flavor": utils.ValueIngoreEmpty(raw["flavor_id"]),
					"count":  utils.ValueIngoreEmpty(raw["count"]),
					"azs":    buildResourceFlavorResourceFlavorAzs(raw["azs"]),
				}
			}
		}
		return rst
	}
	return nil
}

func buildResourceFlavorResourceFlavorAzs(rawParams interface{}) []map[string]interface{} {
	if rawArray, ok := rawParams.([]interface{}); ok {
		if len(rawArray) == 0 {
			return nil
		}

		rst := make([]map[string]interface{}, len(rawArray))
		for i, v := range rawArray {
			if raw, ok := v.(map[string]interface{}); ok {
				rst[i] = map[string]interface{}{
					"az":    utils.ValueIngoreEmpty(raw["az"]),
					"count": utils.ValueIngoreEmpty(raw["count"]),
				}
			}
		}
		return rst
	}
	return nil
}

func updateResourcePoolWaitingForStateCompleted(ctx context.Context, d *schema.ResourceData, meta interface{}, t time.Duration) error {
	stateConf := &resource.StateChangeConf{
		Pending: []string{"PENDING"},
		Target:  []string{"COMPLETED"},
		Refresh: func() (interface{}, string, error) {
			cfg := meta.(*config.Config)
			region := cfg.GetRegion(d)

			getResourcePoolRespBody, err := queryResourcePool(cfg, region, d)
			if err != nil {
				return nil, "ERROR", err
			}
			statusRaw := utils.PathSearch(`status.phase`, getResourcePoolRespBody, nil)
			if statusRaw == nil {
				return nil, "ERROR", fmt.Errorf("error parse 'status.phase' from response body")
			}

			status := fmt.Sprintf("%v", statusRaw)

			unexpectedStatus := []string{
				"Abnormal", "Error", "ScalingFailed", "CreationFailed",
			}
			if utils.StrSliceContains(unexpectedStatus, status) {
				return getResourcePoolRespBody, status, nil
			}

			// check if the resource pool is in the process of expanding capacity
			if rawArray, ok := d.Get("resources").([]interface{}); ok {
				if utils.PathSearch("status.resources.abnormal", getResourcePoolRespBody, nil) != nil {
					return nil, "ERROR", fmt.Errorf("error updating resource pool")
				}

				for _, v := range rawArray {
					raw := v.(map[string]interface{})
					flavor := utils.ValueIngoreEmpty(raw["flavor_id"])
					count := utils.ValueIngoreEmpty(raw["count"])

					searchActiveJsonPath := fmt.Sprintf("length(status.resources.available[?flavor=='%s' && count==`%d`])",
						flavor, count)

					log.Println("searchActiveJsonPath: ", searchActiveJsonPath)
					if utils.PathSearch(searchActiveJsonPath, getResourcePoolRespBody, float64(0)).(float64) == 0 {
						return getResourcePoolRespBody, "PENDING", nil
					}
				}
			}

			// check if the resource pool is in the process of changing scope
			if rawArray, ok := d.Get("scope").([]string); ok {
				for _, v := range rawArray {
					scopeStatus := fmt.Sprintf("status.scope[?scopeType=='%s']|[0].state", v)
					log.Println("scopeStatus: ", scopeStatus)
					if utils.PathSearch(scopeStatus, getResourcePoolRespBody, "").(string) != "Enabled" {
						return getResourcePoolRespBody, "PENDING", nil
					}
				}
			}

			return getResourcePoolRespBody, "COMPLETED", nil
		},
		Timeout:      t,
		Delay:        30 * time.Second,
		PollInterval: 20 * time.Second,
	}
	_, err := stateConf.WaitForStateContext(ctx)
	return err
}

func resourceModelartsResourcePoolDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var (
		deleteResourcePoolHttpUrl = "v2/{project_id}/pools/{id}"
		deleteResourcePoolProduct = "modelarts"
	)
	deleteResourcePoolClient, err := cfg.NewServiceClient(deleteResourcePoolProduct, region)
	if err != nil {
		return diag.Errorf("error creating ModelArts client: %s", err)
	}

	deleteResourcePoolPath := deleteResourcePoolClient.Endpoint + deleteResourcePoolHttpUrl
	deleteResourcePoolPath = strings.ReplaceAll(deleteResourcePoolPath, "{project_id}", deleteResourcePoolClient.ProjectID)
	deleteResourcePoolPath = strings.ReplaceAll(deleteResourcePoolPath, "{id}", d.Id())

	deleteResourcePoolOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		OkCodes: []int{
			200,
		},
		MoreHeaders: map[string]string{"Content-Type": "application/json"},
	}

	_, err = deleteResourcePoolClient.Request("DELETE", deleteResourcePoolPath, &deleteResourcePoolOpt)
	if err != nil {
		return diag.Errorf("error deleting Modelarts resource pool: %s", err)
	}

	err = deleteResourcePoolWaitingForStateCompleted(ctx, d, meta, d.Timeout(schema.TimeoutDelete))
	if err != nil {
		return diag.Errorf("error waiting for the Modelarts resource pool (%s) deletion to complete: %s", d.Id(), err)
	}
	return nil
}

func deleteResourcePoolWaitingForStateCompleted(ctx context.Context, d *schema.ResourceData, meta interface{}, t time.Duration) error {
	stateConf := &resource.StateChangeConf{
		Pending: []string{"PENDING"},
		Target:  []string{"COMPLETED"},
		Refresh: func() (interface{}, string, error) {
			cfg := meta.(*config.Config)
			region := cfg.GetRegion(d)
			_, err := queryResourcePool(cfg, region, d)
			if err != nil {
				if _, ok := err.(golangsdk.ErrDefault404); ok {
					var obj = map[string]string{"code": "COMPLETED"}
					return obj, "COMPLETED", nil
				}

				return nil, "ERROR", err
			}

			return nil, "PENDING", nil
		},
		Timeout:      t,
		Delay:        30 * time.Second,
		PollInterval: 20 * time.Second,
	}
	_, err := stateConf.WaitForStateContext(ctx)
	return err
}
