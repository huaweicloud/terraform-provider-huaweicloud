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

var workflowNonUpdatableParams = []string{"workspace_id", "engine_type", "aop_type", "dataclass_id"}

// @API SecMaster POST /v1/{project_id}/workspaces/{workspace_id}/soc/workflows
// @API SecMaster GET /v1/{project_id}/workspaces/{workspace_id}/soc/workflows/{workflow_id}
// @API SecMaster PUT /v1/{project_id}/workspaces/{workspace_id}/soc/workflows/{workflow_id}
// @API SecMaster DELETE /v1/{project_id}/workspaces/{workspace_id}/soc/workflows/{workflow_id}
func ResourceWorkflow() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceWorkflowCreate,
		ReadContext:   resourceWorkflowRead,
		UpdateContext: resourceWorkflowUpdate,
		DeleteContext: resourceWorkflowDelete,

		Importer: &schema.ResourceImporter{
			StateContext: resourceWorkflowImportState,
		},

		CustomizeDiff: config.FlexibleForceNew(workflowNonUpdatableParams),

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
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"engine_type": {
				Type:     schema.TypeString,
				Required: true,
			},
			"aop_type": {
				Type:     schema.TypeString,
				Required: true,
			},
			"dataclass_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"labels": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			// This parameter only can use in `update`
			"enabled": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},
			// This parameter only can use in `update`
			"version_id": {
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
			"dataclass_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"version": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"project_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"owner_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"creator_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"creator_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"modifier_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"modifier_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"create_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"update_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"use_role": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"edit_role": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"approve_role": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"current_approval_version_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"current_rejected_version_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func buildWorkflowBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"name":         d.Get("name"),
		"engine_type":  d.Get("engine_type"),
		"aop_type":     d.Get("aop_type"),
		"dataclass_id": d.Get("dataclass_id"),
		"labels":       utils.ValueIgnoreEmpty(d.Get("labels")),
		"description":  utils.ValueIgnoreEmpty(d.Get("description")),
	}

	return bodyParams
}

func resourceWorkflowCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg           = meta.(*config.Config)
		region        = cfg.GetRegion(d)
		createHttpUrl = "v1/{project_id}/workspaces/{workspace_id}/soc/workflows"
		workspaceId   = d.Get("workspace_id").(string)
	)

	client, err := cfg.NewServiceClient("secmaster", region)
	if err != nil {
		return diag.Errorf("error creating SecMaster client: %s", err)
	}

	createPath := client.Endpoint + createHttpUrl
	createPath = strings.ReplaceAll(createPath, "{project_id}", client.ProjectID)
	createPath = strings.ReplaceAll(createPath, "{workspace_id}", workspaceId)

	createOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"content-type": "application/json;charset=UTF-8"},
		JSONBody:         utils.RemoveNil(buildWorkflowBodyParams(d)),
	}

	createResp, err := client.Request("POST", createPath, &createOpt)
	if err != nil {
		return diag.Errorf("error creating workflow: %s", err)
	}

	respBody, err := utils.FlattenResponse(createResp)
	if err != nil {
		return diag.FromErr(err)
	}

	workflowId := utils.PathSearch("data.id", respBody, "").(string)
	if workflowId == "" {
		return diag.Errorf("unable to find the workflow ID from the API response")
	}

	d.SetId(workflowId)

	return resourceWorkflowRead(ctx, d, meta)
}

func resourceWorkflowRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg         = meta.(*config.Config)
		region      = cfg.GetRegion(d)
		workspaceId = d.Get("workspace_id").(string)
	)

	client, err := cfg.NewServiceClient("secmaster", region)
	if err != nil {
		return diag.Errorf("error creating SecMaster client: %s", err)
	}

	workflow, err := GetWorkflowInfo(client, workspaceId, d.Id())
	if err != nil {
		// When the workflow does not exist, the response HTTP status code of the query API is `400`
		return common.CheckDeletedDiag(d, common.ConvertExpected400ErrInto404Err(err, "code", "SecMaster.20040611"),
			"error retrieving workflow")
	}

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("name", utils.PathSearch("data.name", workflow, nil)),
		d.Set("engine_type", utils.PathSearch("data.engine_type", workflow, nil)),
		d.Set("aop_type", utils.PathSearch("data.aop_type", workflow, nil)),
		d.Set("dataclass_id", utils.PathSearch("data.dataclass_id", workflow, nil)),
		d.Set("labels", utils.PathSearch("data.labels", workflow, nil)),
		d.Set("description", utils.PathSearch("data.description", workflow, nil)),
		d.Set("enabled", utils.PathSearch("data.enabled", workflow, nil)),
		d.Set("version_id", utils.PathSearch("data.version_id", workflow, nil)),
		d.Set("dataclass_name", utils.PathSearch("data.dataclass_name", workflow, nil)),
		d.Set("version", utils.PathSearch("data.version", workflow, nil)),
		d.Set("project_id", utils.PathSearch("data.project_id", workflow, nil)),
		d.Set("owner_id", utils.PathSearch("data.owner_id", workflow, nil)),
		d.Set("creator_id", utils.PathSearch("data.creator_id", workflow, nil)),
		d.Set("creator_name", utils.PathSearch("data.creator_name", workflow, nil)),
		d.Set("modifier_id", utils.PathSearch("data.modifier_id", workflow, nil)),
		d.Set("modifier_name", utils.PathSearch("data.modifier_name", workflow, nil)),
		d.Set("create_time", utils.PathSearch("data.create_time", workflow, nil)),
		d.Set("update_time", utils.PathSearch("data.update_time", workflow, nil)),
		d.Set("use_role", utils.PathSearch("data.use_role", workflow, nil)),
		d.Set("edit_role", utils.PathSearch("data.edit_role", workflow, nil)),
		d.Set("approve_role", utils.PathSearch("data.approve_role", workflow, nil)),
		d.Set("current_approval_version_id", utils.PathSearch("data.current_approval_version_id", workflow, nil)),
		d.Set("current_rejected_version_id", utils.PathSearch("data.current_rejected_version_id", workflow, nil)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func GetWorkflowInfo(client *golangsdk.ServiceClient, workspaceId, workflowId string) (interface{}, error) {
	httpUrl := "v1/{project_id}/workspaces/{workspace_id}/soc/workflows/{workflow_id}"
	getPath := client.Endpoint + httpUrl
	getPath = strings.ReplaceAll(getPath, "{project_id}", client.ProjectID)
	getPath = strings.ReplaceAll(getPath, "{workspace_id}", workspaceId)
	getPath = strings.ReplaceAll(getPath, "{workflow_id}", workflowId)
	getOpts := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"content-type": "application/json;charset=UTF-8"},
	}

	resp, err := client.Request("GET", getPath, &getOpts)
	if err != nil {
		return nil, err
	}

	return utils.FlattenResponse(resp)
}

func buildUpdateWorkflowBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"name":        d.Get("name"),
		"description": d.Get("description"),
		"enabled":     d.Get("enabled"),
		"version_id":  utils.ValueIgnoreEmpty(d.Get("version_id")),
		"labels":      utils.ValueIgnoreEmpty(d.Get("labels")),
	}

	return bodyParams
}

func resourceWorkflowUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg         = meta.(*config.Config)
		region      = cfg.GetRegion(d)
		workspaceId = d.Get("workspace_id").(string)
		httpUrl     = "v1/{project_id}/workspaces/{workspace_id}/soc/workflows/{workflow_id}"
	)

	client, err := cfg.NewServiceClient("secmaster", region)
	if err != nil {
		return diag.Errorf("error creating SecMaster client: %s", err)
	}

	updatePath := client.Endpoint + httpUrl
	updatePath = strings.ReplaceAll(updatePath, "{project_id}", client.ProjectID)
	updatePath = strings.ReplaceAll(updatePath, "{workspace_id}", workspaceId)
	updatePath = strings.ReplaceAll(updatePath, "{workflow_id}", d.Id())

	updateOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"content-type": "application/json;charset=UTF-8"},
		JSONBody:         utils.RemoveNil(buildUpdateWorkflowBodyParams(d)),
	}

	_, err = client.Request("PUT", updatePath, &updateOpt)
	if err != nil {
		return diag.Errorf("error updating workflow: %s", err)
	}

	return resourceWorkflowRead(ctx, d, meta)
}

func resourceWorkflowDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg         = meta.(*config.Config)
		region      = cfg.GetRegion(d)
		workspaceId = d.Get("workspace_id").(string)
		httpUrl     = "v1/{project_id}/workspaces/{workspace_id}/soc/workflows/{workflow_id}"
	)

	client, err := cfg.NewServiceClient("secmaster", region)
	if err != nil {
		return diag.Errorf("error creating SecMaster client: %s", err)
	}

	deletePath := client.Endpoint + httpUrl
	deletePath = strings.ReplaceAll(deletePath, "{project_id}", client.ProjectID)
	deletePath = strings.ReplaceAll(deletePath, "{workspace_id}", workspaceId)
	deletePath = strings.ReplaceAll(deletePath, "{workflow_id}", d.Id())
	deleteOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"content-type": "application/json;charset=UTF-8"},
	}

	_, err = client.Request("DELETE", deletePath, &deleteOpt)
	if err != nil {
		// If the workflow does not exist, the response HTTP status code of the deletion API is `400`.
		return common.CheckDeletedDiag(d, common.ConvertExpected400ErrInto404Err(err, "code", "SecMaster.20040611"),
			"error deleting workflow")
	}

	return nil
}

func resourceWorkflowImportState(_ context.Context, d *schema.ResourceData, _ interface{}) ([]*schema.ResourceData,
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
