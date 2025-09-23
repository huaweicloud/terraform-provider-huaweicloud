// ---------------------------------------------------------------
// *** AUTO GENERATED CODE ***
// @Product ModelArts
// ---------------------------------------------------------------

package modelarts

import (
	"context"
	"fmt"
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

// @API ModelArts POST /v1/{project_id}/workspaces
// @API ModelArts DELETE /v1/{project_id}/workspaces/{id}
// @API ModelArts GET /v1/{project_id}/workspaces/{id}
// @API ModelArts PUT /v1/{project_id}/workspaces/{id}
func ResourceModelartsWorkspace() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceModelartsWorkspaceCreate,
		UpdateContext: resourceModelartsWorkspaceUpdate,
		ReadContext:   resourceModelartsWorkspaceRead,
		DeleteContext: resourceModelartsWorkspaceDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Timeouts: &schema.ResourceTimeout{
			Delete: schema.DefaultTimeout(10 * time.Minute),
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
				Description: `Workspace name, which consists of 4 to 64 characters.`,
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The description of the worksapce.`,
			},
			"enterprise_project_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: `The enterprise project id of the worksapce.`,
			},
			"auth_type": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `Inference mode.`,
			},
			"grants": {
				Type:        schema.TypeList,
				Elem:        modelartsWorkspaceGrantsSchema(),
				Optional:    true,
				Description: `List of authorized users.`,
			},
			"status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Workspace status.`,
			},
			"status_info": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Status details.`,
			},
		},
	}
}

func modelartsWorkspaceGrantsSchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"user_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `IAM user ID.`,
			},
			"user_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `IAM username.`,
			},
		},
	}
	return &sc
}

func resourceModelartsWorkspaceCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	// createWorkspace: create a Modelarts workspace.
	var (
		createWorkspaceHttpUrl = "v1/{project_id}/workspaces"
		createWorkspaceProduct = "modelarts"
	)
	createWorkspaceClient, err := cfg.NewServiceClient(createWorkspaceProduct, region)
	if err != nil {
		return diag.Errorf("error creating ModelArts Client: %s", err)
	}

	createWorkspacePath := createWorkspaceClient.Endpoint + createWorkspaceHttpUrl
	createWorkspacePath = strings.ReplaceAll(createWorkspacePath, "{project_id}", createWorkspaceClient.ProjectID)

	createWorkspaceOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		OkCodes: []int{
			200,
		},
	}

	createWorkspaceOpt.JSONBody = utils.RemoveNil(buildCreateWorkspaceBodyParams(d, cfg))
	createWorkspaceResp, err := createWorkspaceClient.Request("POST", createWorkspacePath, &createWorkspaceOpt)
	if err != nil {
		return diag.Errorf("error creating Modelarts workspace: %s", err)
	}

	createWorkspaceRespBody, err := utils.FlattenResponse(createWorkspaceResp)
	if err != nil {
		return diag.FromErr(err)
	}

	workspaceId := utils.PathSearch("id", createWorkspaceRespBody, "").(string)
	if workspaceId == "" {
		return diag.Errorf("unable to find the ModelArts workspace ID from the API response")
	}
	d.SetId(workspaceId)
	return resourceModelartsWorkspaceRead(ctx, d, meta)
}

func buildCreateWorkspaceBodyParams(d *schema.ResourceData, cfg *config.Config) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"name":                  d.Get("name"),
		"description":           d.Get("description"),
		"enterprise_project_id": utils.ValueIgnoreEmpty(cfg.GetEnterpriseProjectID(d)),
		"auth_type":             utils.ValueIgnoreEmpty(d.Get("auth_type")),
		"grants":                buildCreateWorkspaceRequestBodyGrants(d.Get("grants")),
	}
	return bodyParams
}

func buildCreateWorkspaceRequestBodyGrants(rawParams interface{}) []map[string]interface{} {
	if rawArray, ok := rawParams.([]interface{}); ok {
		if len(rawArray) == 0 {
			return nil
		}

		rst := make([]map[string]interface{}, len(rawArray))
		for i, v := range rawArray {
			if raw, ok := v.(map[string]interface{}); ok {
				rst[i] = map[string]interface{}{
					"user_id":   utils.ValueIgnoreEmpty(raw["user_id"]),
					"user_name": utils.ValueIgnoreEmpty(raw["user_name"]),
				}
			}
		}
		return rst
	}
	return nil
}

func resourceModelartsWorkspaceRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var mErr *multierror.Error

	// getWorkspace: Query the Modelarts workspace.
	var (
		getWorkspaceHttpUrl = "v1/{project_id}/workspaces/{id}"
		getWorkspaceProduct = "modelarts"
	)
	getWorkspaceClient, err := cfg.NewServiceClient(getWorkspaceProduct, region)
	if err != nil {
		return diag.Errorf("error creating ModelArts Client: %s", err)
	}

	getWorkspacePath := getWorkspaceClient.Endpoint + getWorkspaceHttpUrl
	getWorkspacePath = strings.ReplaceAll(getWorkspacePath, "{project_id}", getWorkspaceClient.ProjectID)
	getWorkspacePath = strings.ReplaceAll(getWorkspacePath, "{id}", d.Id())

	getWorkspaceOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		OkCodes: []int{
			200,
		},
	}

	getWorkspaceResp, err := getWorkspaceClient.Request("GET", getWorkspacePath, &getWorkspaceOpt)

	if err != nil {
		return common.CheckDeletedDiag(d, common.ConvertExpected400ErrInto404Err(err, "error_code", "ModelArts.7005"),
			"error retrieving Modelarts workspace")
	}

	getWorkspaceRespBody, err := utils.FlattenResponse(getWorkspaceResp)
	if err != nil {
		return diag.FromErr(err)
	}

	mErr = multierror.Append(
		mErr,
		d.Set("region", region),
		d.Set("name", utils.PathSearch("name", getWorkspaceRespBody, nil)),
		d.Set("description", utils.PathSearch("description", getWorkspaceRespBody, nil)),
		d.Set("enterprise_project_id", utils.PathSearch("enterprise_project_id", getWorkspaceRespBody, nil)),
		d.Set("auth_type", utils.PathSearch("auth_type", getWorkspaceRespBody, nil)),
		d.Set("grants", flattenGetWorkspaceResponseBodyGrants(getWorkspaceRespBody)),
		d.Set("status", utils.PathSearch("status", getWorkspaceRespBody, nil)),
		d.Set("status_info", utils.PathSearch("status_info", getWorkspaceRespBody, nil)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenGetWorkspaceResponseBodyGrants(resp interface{}) []interface{} {
	if resp == nil {
		return nil
	}
	curJson := utils.PathSearch("grants", resp, make([]interface{}, 0))
	curArray := curJson.([]interface{})
	rst := make([]interface{}, 0, len(curArray))
	for _, v := range curArray {
		rst = append(rst, map[string]interface{}{
			"user_id":   utils.PathSearch("user_id", v, nil),
			"user_name": utils.PathSearch("user_name", v, nil),
		})
	}
	return rst
}

func resourceModelartsWorkspaceUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	updateWorkspaceChanges := []string{
		"name",
		"description",
		"auth_type",
		"grants",
	}

	if d.HasChanges(updateWorkspaceChanges...) {
		// updateWorkspace: update the Modelarts workspace.
		var (
			updateWorkspaceHttpUrl = "v1/{project_id}/workspaces/{id}"
			updateWorkspaceProduct = "modelarts"
		)
		updateWorkspaceClient, err := cfg.NewServiceClient(updateWorkspaceProduct, region)
		if err != nil {
			return diag.Errorf("error creating ModelArts Client: %s", err)
		}

		updateWorkspacePath := updateWorkspaceClient.Endpoint + updateWorkspaceHttpUrl
		updateWorkspacePath = strings.ReplaceAll(updateWorkspacePath, "{project_id}", updateWorkspaceClient.ProjectID)
		updateWorkspacePath = strings.ReplaceAll(updateWorkspacePath, "{id}", d.Id())

		updateWorkspaceOpt := golangsdk.RequestOpts{
			KeepResponseBody: true,
			OkCodes: []int{
				200,
			},
		}

		updateWorkspaceOpt.JSONBody = utils.RemoveNil(buildUpdateWorkspaceBodyParams(d))
		_, err = updateWorkspaceClient.Request("PUT", updateWorkspacePath, &updateWorkspaceOpt)
		if err != nil {
			return diag.Errorf("error updating Modelarts workspace: %s", err)
		}
	}
	return resourceModelartsWorkspaceRead(ctx, d, meta)
}

func buildUpdateWorkspaceBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"name":        d.Get("name"),
		"description": d.Get("description"),
		"auth_type":   utils.ValueIgnoreEmpty(d.Get("auth_type")),
		"grants":      buildUpdateWorkspaceRequestBodyGrants(d.Get("grants")),
	}
	return bodyParams
}

func buildUpdateWorkspaceRequestBodyGrants(rawParams interface{}) []map[string]interface{} {
	if rawArray, ok := rawParams.([]interface{}); ok {
		if len(rawArray) == 0 {
			return nil
		}

		rst := make([]map[string]interface{}, len(rawArray))
		for i, v := range rawArray {
			if raw, ok := v.(map[string]interface{}); ok {
				rst[i] = map[string]interface{}{
					"user_id":   utils.ValueIgnoreEmpty(raw["user_id"]),
					"user_name": utils.ValueIgnoreEmpty(raw["user_name"]),
				}
			}
		}
		return rst
	}
	return nil
}

func resourceModelartsWorkspaceDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	// deleteWorkspace: delete Modelarts workspace
	var (
		deleteWorkspaceHttpUrl = "v1/{project_id}/workspaces/{id}"
		deleteWorkspaceProduct = "modelarts"
	)
	deleteWorkspaceClient, err := cfg.NewServiceClient(deleteWorkspaceProduct, region)
	if err != nil {
		return diag.Errorf("error creating ModelArts Client: %s", err)
	}

	deleteWorkspacePath := deleteWorkspaceClient.Endpoint + deleteWorkspaceHttpUrl
	deleteWorkspacePath = strings.ReplaceAll(deleteWorkspacePath, "{project_id}", deleteWorkspaceClient.ProjectID)
	deleteWorkspacePath = strings.ReplaceAll(deleteWorkspacePath, "{id}", d.Id())

	deleteWorkspaceOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		OkCodes: []int{
			200,
		},
	}

	_, err = deleteWorkspaceClient.Request("DELETE", deleteWorkspacePath, &deleteWorkspaceOpt)
	if err != nil {
		return diag.Errorf("error deleting Modelarts workspace: %s", err)
	}

	err = deleteWorkspaceWaitingForStateCompleted(ctx, d, meta, d.Timeout(schema.TimeoutDelete))
	if err != nil {
		return diag.Errorf("error waiting for the Modelarts workspace (%s) deletion to complete: %s", d.Id(), err)
	}

	return nil
}

func deleteWorkspaceWaitingForStateCompleted(ctx context.Context, d *schema.ResourceData, meta interface{}, t time.Duration) error {
	stateConf := &resource.StateChangeConf{
		Pending: []string{"PENDING"},
		Target:  []string{"COMPLETED"},
		Refresh: func() (interface{}, string, error) {
			cfg := meta.(*config.Config)
			region := cfg.GetRegion(d)
			var (
				deleteWorkspaceWaitingHttpUrl = "v1/{project_id}/workspaces/{id}"
				deleteWorkspaceWaitingProduct = "modelarts"
			)
			deleteWorkspaceWaitingClient, err := cfg.NewServiceClient(deleteWorkspaceWaitingProduct, region)
			if err != nil {
				return nil, "ERROR", fmt.Errorf("error creating ModelArts Client: %s", err)
			}

			deleteWorkspaceWaitingPath := deleteWorkspaceWaitingClient.Endpoint + deleteWorkspaceWaitingHttpUrl
			deleteWorkspaceWaitingPath = strings.ReplaceAll(deleteWorkspaceWaitingPath, "{project_id}", deleteWorkspaceWaitingClient.ProjectID)
			deleteWorkspaceWaitingPath = strings.ReplaceAll(deleteWorkspaceWaitingPath, "{id}", d.Id())

			deleteWorkspaceWaitingOpt := golangsdk.RequestOpts{
				KeepResponseBody: true,
				OkCodes: []int{
					200,
				},
			}

			deleteWorkspaceWaitingResp, err := deleteWorkspaceWaitingClient.Request("GET", deleteWorkspaceWaitingPath, &deleteWorkspaceWaitingOpt)
			if err != nil {
				parsedErr := common.ConvertExpected400ErrInto404Err(err, "error_code", "ModelArts.7005")
				if _, ok := parsedErr.(golangsdk.ErrDefault404); ok {
					// When the error code is 404, the value of respBody is nil, and a non-null value is returned to avoid continuing the loop check.
					return "Resource Not Found", "COMPLETED", nil
				}
				return nil, "ERROR", err
			}

			deleteWorkspaceWaitingRespBody, err := utils.FlattenResponse(deleteWorkspaceWaitingResp)
			if err != nil {
				return nil, "ERROR", err
			}
			status := utils.PathSearch(`status`, deleteWorkspaceWaitingRespBody, "").(string)

			pendingStatus := []string{
				"DELETING",
			}
			if utils.StrSliceContains(pendingStatus, status) {
				return deleteWorkspaceWaitingRespBody, "PENDING", nil
			}

			return deleteWorkspaceWaitingRespBody, status, nil
		},
		Timeout:      t,
		Delay:        10 * time.Second,
		PollInterval: 10 * time.Second,
	}
	_, err := stateConf.WaitForStateContext(ctx)
	return err
}
