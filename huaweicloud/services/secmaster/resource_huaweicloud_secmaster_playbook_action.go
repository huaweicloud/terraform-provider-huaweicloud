package secmaster

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

const (
	GetPlaybookActionNotFound    = "SecMaster.20030005"
	DeletePlaybookActionNotFound = "SecMaster.20048004"
)

// @API SecMaster GET /v1/{project_id}/workspaces/{workspace_id}/soc/playbooks/versions/{version_id}/actions
// @API SecMaster POST /v1/{project_id}/workspaces/{workspace_id}/soc/playbooks/versions/{version_id}/actions
// @API SecMaster DELETE /v1/{project_id}/workspaces/{workspace_id}/soc/playbooks/versions/{version_id}/actions/{id}
// @API SecMaster PUT /v1/{project_id}/workspaces/{workspace_id}/soc/playbooks/versions/{version_id}/actions/{id}
// @API SecMaster GET /v1/{project_id}/workspaces/{workspace_id}/soc/playbooks/versions/{version_id}
// @API SecMaster PUT /v1/{project_id}/workspaces/{workspace_id}/soc/playbooks/versions/{id}
func ResourcePlaybookAction() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourcePlaybookActionCreate,
		UpdateContext: resourcePlaybookActionUpdate,
		ReadContext:   resourcePlaybookActionRead,
		DeleteContext: resourcePlaybookActionDelete,
		Importer: &schema.ResourceImporter{
			StateContext: resourcePlaybookActionImportState,
		},

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"workspace_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: `Specifies the ID of the workspace to which the playbook action belongs.`,
			},
			"version_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: `Specifies the playbook version ID of the action.`,
			},
			"action_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `Specifies the workflow ID.`,
			},
			"name": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `Specifies the playbook action name.`,
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specifies the description of the playbook action.`,
			},
			"action_type": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `Specifies the playbook action type.`,
			},
			"sort_order": {
				Type:        schema.TypeFloat,
				Optional:    true,
				Computed:    true,
				Description: `Specifies the sort order of the playbook action.`,
			},
			"playbook_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Indicates the playbook ID of the action.`,
			},
			"created_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Indicates the created time of the playbook action.`,
			},
			"updated_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Indicates the updated time of the playbook action.`,
			},
		},
	}
}

func resourcePlaybookActionCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	// createPlaybookAction: Create a SecMaster playbook action.
	var (
		createPlaybookActionHttpUrl = "v1/{project_id}/workspaces/{workspace_id}/soc/playbooks/versions/{version_id}/actions"
		createPlaybookActionProduct = "secmaster"
	)
	createPlaybookActionClient, err := cfg.NewServiceClient(createPlaybookActionProduct, region)
	if err != nil {
		return diag.Errorf("error creating SecMaster client: %s", err)
	}

	createPlaybookActionPath := createPlaybookActionClient.Endpoint + createPlaybookActionHttpUrl
	createPlaybookActionPath = strings.ReplaceAll(createPlaybookActionPath, "{project_id}", createPlaybookActionClient.ProjectID)
	createPlaybookActionPath = strings.ReplaceAll(createPlaybookActionPath, "{workspace_id}", d.Get("workspace_id").(string))
	createPlaybookActionPath = strings.ReplaceAll(createPlaybookActionPath, "{version_id}", d.Get("version_id").(string))

	createPlaybookActionOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"Content-Type": "application/json"},
	}

	createPlaybookActionOpt.JSONBody = []interface{}{utils.RemoveNil(buildCreatePlaybookActionBodyParams(d))}
	createPlaybookActionResp, err := createPlaybookActionClient.Request("POST", createPlaybookActionPath, &createPlaybookActionOpt)
	if err != nil {
		return diag.Errorf("error creating SecMaster playbook action: %s", err)
	}

	createPlaybookActionRespBody, err := utils.FlattenResponse(createPlaybookActionResp)
	if err != nil {
		return diag.FromErr(err)
	}

	id := utils.PathSearch("data[0].id", createPlaybookActionRespBody, "").(string)
	if id == "" {
		return diag.Errorf("error creating SecMaster playbook action: ID is not found in API response")
	}
	d.SetId(id)

	return resourcePlaybookActionRead(ctx, d, meta)
}

func buildCreatePlaybookActionBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"action_id":   d.Get("action_id"),
		"name":        utils.ValueIgnoreEmpty(d.Get("name")),
		"description": utils.ValueIgnoreEmpty(d.Get("description")),
		"action_type": utils.ValueIgnoreEmpty(d.Get("action_type")),
		"sort_order":  utils.ValueIgnoreEmpty(d.Get("sort_order")),
	}
	return bodyParams
}

func resourcePlaybookActionRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	// getPlaybookAction: Query the SecMaster playbook action detail
	var (
		getPlaybookActionHttpUrl = "v1/{project_id}/workspaces/{workspace_id}/soc/playbooks/versions/{version_id}/actions"
		getPlaybookActionProduct = "secmaster"
	)
	getPlaybookActionClient, err := cfg.NewServiceClient(getPlaybookActionProduct, region)
	if err != nil {
		return diag.Errorf("error creating SecMaster client: %s", err)
	}

	getPlaybookActionPath := getPlaybookActionClient.Endpoint + getPlaybookActionHttpUrl
	getPlaybookActionPath = strings.ReplaceAll(getPlaybookActionPath, "{project_id}", getPlaybookActionClient.ProjectID)
	getPlaybookActionPath = strings.ReplaceAll(getPlaybookActionPath, "{workspace_id}", d.Get("workspace_id").(string))
	getPlaybookActionPath = strings.ReplaceAll(getPlaybookActionPath, "{version_id}", d.Get("version_id").(string))
	getPlaybookActionPath += "?limit=1000"

	getPlaybookActionOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"Content-Type": "application/json"},
	}

	getPlaybookActionResp, err := getPlaybookActionClient.Request("GET", getPlaybookActionPath, &getPlaybookActionOpt)
	if err != nil {
		// "SecMaster.20010001": workspace ID not found
		// "SecMaster.20030005": resource not found
		err = common.ConvertExpected403ErrInto404Err(err, "code", WorkspaceNotFound)
		err = common.ConvertExpected400ErrInto404Err(err, "code", GetPlaybookActionNotFound)
		return common.CheckDeletedDiag(d, err, "error retrieving SecMaster playbook action")
	}

	getPlaybookActionRespBody, err := utils.FlattenResponse(getPlaybookActionResp)
	if err != nil {
		return diag.FromErr(err)
	}

	jsonPath := fmt.Sprintf("data[?id=='%s']|[0]", d.Id())
	getPlaybookActionRespBody = utils.PathSearch(jsonPath, getPlaybookActionRespBody, nil)
	if getPlaybookActionRespBody == nil {
		return common.CheckDeletedDiag(d, golangsdk.ErrDefault404{}, "")
	}

	mErr := multierror.Append(
		d.Set("region", region),
		d.Set("name", utils.PathSearch("name", getPlaybookActionRespBody, nil)),
		d.Set("description", utils.PathSearch("description", getPlaybookActionRespBody, nil)),
		d.Set("action_type", utils.PathSearch("action_type", getPlaybookActionRespBody, nil)),
		d.Set("sort_order", utils.PathSearch("sort_order", getPlaybookActionRespBody, nil)),
		d.Set("action_id", utils.PathSearch("action_id", getPlaybookActionRespBody, nil)),
		d.Set("playbook_id", utils.PathSearch("playbook_id", getPlaybookActionRespBody, nil)),
		d.Set("created_at", utils.PathSearch("create_time", getPlaybookActionRespBody, nil)),
		d.Set("updated_at", utils.PathSearch("update_time", getPlaybookActionRespBody, nil)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func resourcePlaybookActionUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	// updatePlaybookAction: Update the configuration of SecMaster playbook action
	var (
		updatePlaybookActionHttpUrl = "v1/{project_id}/workspaces/{workspace_id}/soc/playbooks/versions/{version_id}/actions/{id}"
		updatePlaybookActionProduct = "secmaster"
	)
	updatePlaybookActionClient, err := cfg.NewServiceClient(updatePlaybookActionProduct, region)
	if err != nil {
		return diag.Errorf("error creating SecMaster client: %s", err)
	}

	updatePlaybookActionPath := updatePlaybookActionClient.Endpoint + updatePlaybookActionHttpUrl
	updatePlaybookActionPath = strings.ReplaceAll(updatePlaybookActionPath, "{project_id}", updatePlaybookActionClient.ProjectID)
	updatePlaybookActionPath = strings.ReplaceAll(updatePlaybookActionPath, "{workspace_id}", d.Get("workspace_id").(string))
	updatePlaybookActionPath = strings.ReplaceAll(updatePlaybookActionPath, "{version_id}", d.Get("version_id").(string))
	updatePlaybookActionPath = strings.ReplaceAll(updatePlaybookActionPath, "{id}", d.Id())

	updatePlaybookActionOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"Content-Type": "application/json"},
	}

	updatePlaybookActionOpt.JSONBody = utils.RemoveNil(buildUpdatePlaybookActionBodyParams(d))
	_, err = updatePlaybookActionClient.Request("PUT", updatePlaybookActionPath, &updatePlaybookActionOpt)
	if err != nil {
		return diag.Errorf("error updating SecMaster playbook action: %s", err)
	}
	return resourcePlaybookActionRead(ctx, d, meta)
}

func buildUpdatePlaybookActionBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"name":        utils.ValueIgnoreEmpty(d.Get("name")),
		"description": d.Get("description"),
		"action_type": utils.ValueIgnoreEmpty(d.Get("action_type")),
		"sort_order":  utils.ValueIgnoreEmpty(d.Get("sort_order")),
		"action_id":   utils.ValueIgnoreEmpty(d.Get("action_id")),
	}
	return bodyParams
}

func resourcePlaybookActionDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	workspaceID := d.Get("workspace_id").(string)
	versionID := d.Get("version_id").(string)

	client, err := cfg.NewServiceClient("secmaster", region)
	if err != nil {
		return diag.Errorf("error creating SecMaster client: %s", err)
	}

	// Check whether the version is enabled.
	// Before deleting this version, you need to ensure that the version is not enabled.
	playbookVersion, err := GetPlaybookVersion(client, workspaceID, versionID)
	if err != nil {
		return common.CheckDeletedDiag(d, err, "error querying SecMaster playbook version")
	}

	if utils.PathSearch("enabled", playbookVersion, false).(bool) {
		bodyParams := backfillUpdateBodyParams(playbookVersion)
		bodyParams["enabled"] = false
		err = updatePlaybookVersion(client, workspaceID, versionID, bodyParams)
		if err != nil {
			return diag.Errorf("error disabling SecMaster playbook version: %s", err)
		}
	}

	// deletePlaybookAction: Delete an existing SecMaster playbook action
	deletePlaybookActionHttpUrl := "v1/{project_id}/workspaces/{workspace_id}/soc/playbooks/versions/{version_id}/actions/{id}"
	deletePlaybookActionPath := client.Endpoint + deletePlaybookActionHttpUrl
	deletePlaybookActionPath = strings.ReplaceAll(deletePlaybookActionPath, "{project_id}", client.ProjectID)
	deletePlaybookActionPath = strings.ReplaceAll(deletePlaybookActionPath, "{workspace_id}", workspaceID)
	deletePlaybookActionPath = strings.ReplaceAll(deletePlaybookActionPath, "{version_id}", versionID)
	deletePlaybookActionPath = strings.ReplaceAll(deletePlaybookActionPath, "{id}", d.Id())

	deletePlaybookActionOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"Content-Type": "application/json"},
	}

	_, err = client.Request("DELETE", deletePlaybookActionPath, &deletePlaybookActionOpt)
	if err != nil {
		// "SecMaster.20010001": workspace ID not found
		// "SecMaster.20048004": playbook action not found
		err = common.ConvertExpected403ErrInto404Err(err, "code", WorkspaceNotFound)
		err = common.ConvertExpected400ErrInto404Err(err, "code", DeletePlaybookActionNotFound)
		return common.CheckDeletedDiag(d, err, "error deleting SecMaster playbook action")
	}

	return nil
}

func resourcePlaybookActionImportState(_ context.Context, d *schema.ResourceData, _ interface{}) ([]*schema.ResourceData, error) {
	parts := strings.Split(d.Id(), "/")
	if len(parts) != 3 {
		return nil, fmt.Errorf("invalid format specified for import id, must be <workspace_id>/<playbook_version_id>/<playbook_action_id>")
	}

	d.SetId(parts[2])

	var mErr *multierror.Error
	mErr = multierror.Append(
		mErr,
		d.Set("workspace_id", parts[0]),
		d.Set("version_id", parts[1]),
	)

	err := mErr.ErrorOrNil()
	if err != nil {
		return nil, err
	}

	return []*schema.ResourceData{d}, nil
}
