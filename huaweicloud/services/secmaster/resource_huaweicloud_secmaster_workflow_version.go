package secmaster

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

var nonUpdatableParamsWorkflowVersion = []string{"workspace_id", "workflow_id"}

// @API SecMaster POST /v1/{project_id}/workspaces/{workspace_id}/soc/workflows/{workflow_id}/versions
// @API SecMaster GET /v1/{project_id}/workspaces/{workspace_id}/soc/workflows/versions/{version_id}
// @API SecMaster PUT /v1/{project_id}/workspaces/{workspace_id}/soc/workflows/versions/{version_id}
// @API SecMaster DELETE /v1/{project_id}/workspaces/{workspace_id}/soc/workflows/versions/{version_id}
func ResourceWorkflowVersion() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceWorkflowVersionCreate,
		ReadContext:   resourceWorkflowVersionRead,
		UpdateContext: resourceWorkflowVersionUpdate,
		DeleteContext: resourceWorkflowVersionDelete,

		Importer: &schema.ResourceImporter{
			StateContext: resourceWorkflowVersionImportState,
		},

		CustomizeDiff: config.FlexibleForceNew(nonUpdatableParamsWorkflowVersion),

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"workspace_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"workflow_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"taskflow": {
				Type:     schema.TypeString,
				Required: true,
			},
			"taskconfig": {
				Type:     schema.TypeString,
				Required: true,
			},
			"taskflow_type": {
				Type:     schema.TypeString,
				Required: true,
			},
			"aop_type": {
				Type:     schema.TypeString,
				Required: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"status": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"enable_force_new": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"true", "false"}, false),
				Description:  utils.SchemaDesc("", utils.SchemaDescInput{Internal: true}),
			},
			"owner_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"creator_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"version": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"enabled": {
				Type:     schema.TypeBool,
				Computed: true,
			},
		},
	}
}

func buildCreateWorkflowVersionBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"name":          d.Get("name"),
		"taskflow":      d.Get("taskflow"),
		"taskconfig":    d.Get("taskconfig"),
		"taskflow_type": d.Get("taskflow_type"),
		"aop_type":      d.Get("aop_type"),
		"description":   utils.ValueIgnoreEmpty(d.Get("description")),
	}

	return bodyParams
}

func resourceWorkflowVersionCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg           = meta.(*config.Config)
		region        = cfg.GetRegion(d)
		workspaceId   = d.Get("workspace_id").(string)
		workflowId    = d.Get("workflow_id").(string)
		createHttpUrl = "v1/{project_id}/workspaces/{workspace_id}/soc/workflows/{workflow_id}/versions"
	)

	client, err := cfg.NewServiceClient("secmaster", region)
	if err != nil {
		return diag.Errorf("error creating SecMaster client: %s", err)
	}

	createPath := client.Endpoint + createHttpUrl
	createPath = strings.ReplaceAll(createPath, "{project_id}", client.ProjectID)
	createPath = strings.ReplaceAll(createPath, "{workspace_id}", workspaceId)
	createPath = strings.ReplaceAll(createPath, "{workflow_id}", workflowId)

	createOpt := golangsdk.RequestOpts{
		MoreHeaders:      map[string]string{"content-type": "application/json;charset=UTF-8"},
		KeepResponseBody: true,
		JSONBody:         utils.RemoveNil(buildCreateWorkflowVersionBodyParams(d)),
	}

	resp, err := client.Request("POST", createPath, &createOpt)
	if err != nil {
		return diag.Errorf("error creating workflow version: %s", err)
	}

	respBody, err := utils.FlattenResponse(resp)
	if err != nil {
		return diag.FromErr(err)
	}

	versionId := utils.PathSearch("data.id", respBody, "").(string)
	if versionId == "" {
		return diag.Errorf("error creating workflow version: unable to find workflow version ID")
	}

	d.SetId(versionId)

	return resourceWorkflowVersionRead(ctx, d, meta)
}

func resourceWorkflowVersionRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg         = meta.(*config.Config)
		region      = cfg.GetRegion(d)
		workspaceId = d.Get("workspace_id").(string)
	)

	client, err := cfg.NewServiceClient("secmaster", region)
	if err != nil {
		return diag.Errorf("error creating SecMaster client: %s", err)
	}

	version, err := GetWorkflowVersionInfo(client, workspaceId, d.Id())
	if err != nil {
		// When the workflow version does not exist, the response HTTP status code of the query API is 400
		return common.CheckDeletedDiag(d, common.ConvertExpected400ErrInto404Err(err, "code", "SecMaster.20040605"),
			"error retrieving workflow version")
	}

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("workflow_id", utils.PathSearch("data.aopworkflow_id", version, nil)),
		d.Set("name", utils.PathSearch("data.name", version, nil)),
		d.Set("taskflow", utils.PathSearch("data.taskflow", version, nil)),
		d.Set("taskconfig", utils.PathSearch("data.taskconfig", version, nil)),
		d.Set("taskflow_type", utils.PathSearch("data.taskflow_type", version, nil)),
		d.Set("aop_type", utils.PathSearch("data.aop_type", version, nil)),
		d.Set("description", utils.PathSearch("data.description", version, nil)),
		d.Set("status", utils.PathSearch("data.status", version, nil)),
		d.Set("owner_id", utils.PathSearch("data.owner_id", version, nil)),
		d.Set("creator_id", utils.PathSearch("data.creator_id", version, nil)),
		d.Set("version", utils.PathSearch("data.version", version, nil)),
		d.Set("enabled", utils.PathSearch("data.enabled", version, nil)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func GetWorkflowVersionInfo(client *golangsdk.ServiceClient, workspaceId, versionId string) (interface{}, error) {
	httpUrl := "v1/{project_id}/workspaces/{workspace_id}/soc/workflows/versions/{version_id}"
	getPath := client.Endpoint + httpUrl
	getPath = strings.ReplaceAll(getPath, "{project_id}", client.ProjectID)
	getPath = strings.ReplaceAll(getPath, "{workspace_id}", workspaceId)
	getPath = strings.ReplaceAll(getPath, "{version_id}", versionId)
	getOpts := golangsdk.RequestOpts{
		MoreHeaders:      map[string]string{"content-type": "application/json;charset=UTF-8"},
		KeepResponseBody: true,
	}

	resp, err := client.Request("GET", getPath, &getOpts)
	if err != nil {
		return nil, err
	}

	return utils.FlattenResponse(resp)
}

func buildUpdateWorkflowVersionBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"name":          d.Get("name"),
		"taskflow":      d.Get("taskflow"),
		"taskconfig":    d.Get("taskconfig"),
		"taskflow_type": d.Get("taskflow_type"),
		"aop_type":      d.Get("aop_type"),
		"description":   d.Get("description"),
	}

	return bodyParams
}

func resourceWorkflowVersionUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg         = meta.(*config.Config)
		region      = cfg.GetRegion(d)
		workspaceId = d.Get("workspace_id").(string)
		httpUrl     = "v1/{project_id}/workspaces/{workspace_id}/soc/workflows/versions/{version_id}"
	)

	client, err := cfg.NewServiceClient("secmaster", region)
	if err != nil {
		return diag.Errorf("error creating SecMaster client: %s", err)
	}

	updatePath := client.Endpoint + httpUrl
	updatePath = strings.ReplaceAll(updatePath, "{project_id}", client.ProjectID)
	updatePath = strings.ReplaceAll(updatePath, "{workspace_id}", workspaceId)
	updatePath = strings.ReplaceAll(updatePath, "{version_id}", d.Id())
	updateBody := buildUpdateWorkflowVersionBodyParams(d)

	if d.HasChangeExcept("status") {
		updateOpt := golangsdk.RequestOpts{
			MoreHeaders:      map[string]string{"content-type": "application/json;charset=UTF-8"},
			KeepResponseBody: true,
			JSONBody:         updateBody,
		}

		_, err = client.Request("PUT", updatePath, &updateOpt)
		if err != nil {
			return diag.Errorf("error updating workflow version: %s", err)
		}
	}

	if d.HasChange("status") {
		updateBody["status"] = utils.ValueIgnoreEmpty(d.Get("status"))
		updateOpt := golangsdk.RequestOpts{
			MoreHeaders:      map[string]string{"content-type": "application/json;charset=UTF-8"},
			KeepResponseBody: true,
			JSONBody:         utils.RemoveNil(updateBody),
		}

		_, err = client.Request("PUT", updatePath, &updateOpt)
		if err != nil {
			return diag.Errorf("error updating workflow version: %s", err)
		}
	}

	return resourceWorkflowVersionRead(ctx, d, meta)
}

func resourceWorkflowVersionDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg         = meta.(*config.Config)
		region      = cfg.GetRegion(d)
		workspaceId = d.Get("workspace_id").(string)
		httpUrl     = "v1/{project_id}/workspaces/{workspace_id}/soc/workflows/versions/{version_id}"
	)

	client, err := cfg.NewServiceClient("secmaster", region)
	if err != nil {
		return diag.Errorf("error creating SecMaster client: %s", err)
	}

	deletePath := client.Endpoint + httpUrl
	deletePath = strings.ReplaceAll(deletePath, "{project_id}", client.ProjectID)
	deletePath = strings.ReplaceAll(deletePath, "{workspace_id}", workspaceId)
	deletePath = strings.ReplaceAll(deletePath, "{version_id}", d.Id())
	deleteOpt := golangsdk.RequestOpts{
		MoreHeaders:      map[string]string{"content-type": "application/json;charset=UTF-8"},
		KeepResponseBody: true,
	}

	_, err = client.Request("DELETE", deletePath, &deleteOpt)
	if err != nil {
		// If the workflow version does not exist, the response HTTP status code of the deletion API is 400.
		return common.CheckDeletedDiag(d, common.ConvertExpected400ErrInto404Err(err, "code", "SecMaster.20040605"),
			fmt.Sprintf("error deleting workflow version, the error message: %s", err))
	}

	return nil
}

func resourceWorkflowVersionImportState(_ context.Context, d *schema.ResourceData, _ interface{}) ([]*schema.ResourceData,
	error) {
	importedId := d.Id()
	parts := strings.Split(importedId, "/")
	if len(parts) != 2 {
		return nil, fmt.Errorf("invalid format specified for import ID, want '<workspace_id>/<id>', but got '%s'",
			importedId)
	}

	d.SetId(parts[1])

	mErr := multierror.Append(nil,
		d.Set("workspace_id", parts[0]),
	)

	return []*schema.ResourceData{d}, mErr.ErrorOrNil()
}
