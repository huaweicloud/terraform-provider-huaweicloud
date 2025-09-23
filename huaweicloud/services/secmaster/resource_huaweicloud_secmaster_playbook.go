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
	WorkspaceNotFound = "SecMaster.20010001"
)

// @API SecMaster GET /v1/{project_id}/workspaces/{workspace_id}/soc/playbooks/{playbook_id}
// @API SecMaster PUT /v1/{project_id}/workspaces/{workspace_id}/soc/playbooks/{id}
// @API SecMaster DELETE /v1/{project_id}/workspaces/{workspace_id}/soc/playbooks/{id}
// @API SecMaster POST /v1/{project_id}/workspaces/{workspace_id}/soc/playbooks
func ResourcePlaybook() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourcePlaybookCreate,
		UpdateContext: resourcePlaybookUpdate,
		ReadContext:   resourcePlaybookRead,
		DeleteContext: resourcePlaybookDelete,
		Importer: &schema.ResourceImporter{
			StateContext: resourcePlaybookImportState,
		},

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
				ForceNew: true,
			},
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `Specifies the playbook name.`,
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specifies the description of the playbook.`,
			},
			"approve_role": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Indicates the approve role of the playbook.`,
			},
			"created_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Indicates the created time of the playbook.`,
			},
			"updated_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Indicates the updated time of the playbook.`,
			},
			"dataclass_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Indicates the data class ID.`,
			},
			"dataclass_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Indicates the data class name.`,
			},
			"edit_role": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Indicates the edit role.`,
			},
			"owner_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Indicates the owner ID.`,
			},
			"reject_version_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Indicates the rejected version ID.`,
			},
			"unaudited_version_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Indicates the unaudited version ID.`,
			},
			"user_role": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Indicates the user role.`,
			},
			"version": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Indicates the version.`,
			},
			"version_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Indicates the version ID.`,
			},
			"enabled": {
				Type:        schema.TypeBool,
				Optional:    true,
				Computed:    true,
				Description: `schema: Deprecated; Specifies whether to enable the playbook.`,
			},
			"active_version_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `schema: Deprecated; Specifies the active version ID.`,
			},
		},
	}
}

func resourcePlaybookCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	// createPlaybook: Create a SecMaster playbook.
	var (
		createPlaybookHttpUrl = "v1/{project_id}/workspaces/{workspace_id}/soc/playbooks"
		createPlaybookProduct = "secmaster"
	)
	createPlaybookClient, err := cfg.NewServiceClient(createPlaybookProduct, region)
	if err != nil {
		return diag.Errorf("error creating SecMaster client: %s", err)
	}

	createPlaybookPath := createPlaybookClient.Endpoint + createPlaybookHttpUrl
	createPlaybookPath = strings.ReplaceAll(createPlaybookPath, "{project_id}", createPlaybookClient.ProjectID)
	createPlaybookPath = strings.ReplaceAll(createPlaybookPath, "{workspace_id}", d.Get("workspace_id").(string))

	createPlaybookOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"Content-Type": "application/json"},
	}

	createPlaybookOpt.JSONBody = utils.RemoveNil(buildCreatePlaybookBodyParams(d))
	createPlaybookResp, err := createPlaybookClient.Request("POST", createPlaybookPath, &createPlaybookOpt)
	if err != nil {
		return diag.Errorf("error creating SecMaster playbook: %s", err)
	}

	createPlaybookRespBody, err := utils.FlattenResponse(createPlaybookResp)
	if err != nil {
		return diag.FromErr(err)
	}

	id := utils.PathSearch("data.id", createPlaybookRespBody, "").(string)
	if id == "" {
		return diag.Errorf("error creating SecMaster playbook: ID is not found in API response")
	}
	d.SetId(id)

	return resourcePlaybookUpdate(ctx, d, meta)
}

func buildCreatePlaybookBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"name":        d.Get("name"),
		"description": utils.ValueIgnoreEmpty(d.Get("description")),
	}
	return bodyParams
}

func resourcePlaybookRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	client, err := cfg.NewServiceClient("secmaster", region)
	if err != nil {
		return diag.Errorf("error creating SecMaster client: %s", err)
	}

	playbook, err := GetPlaybook(client, d.Get("workspace_id").(string), d.Id())
	if err != nil {
		return common.CheckDeletedDiag(d,
			common.ConvertExpected403ErrInto404Err(err, "code", WorkspaceNotFound), "error retrieving SecMaster playbook")
	}

	mErr := multierror.Append(
		d.Set("region", region),
		d.Set("name", utils.PathSearch("name", playbook, nil)),
		d.Set("description", utils.PathSearch("description", playbook, nil)),
		d.Set("enabled", utils.PathSearch("enabled", playbook, nil)),
		d.Set("approve_role", utils.PathSearch("approve_role", playbook, nil)),
		d.Set("created_at", utils.PathSearch("create_time", playbook, nil)),
		d.Set("updated_at", utils.PathSearch("update_time", playbook, nil)),
		d.Set("dataclass_id", utils.PathSearch("dataclass_id", playbook, nil)),
		d.Set("dataclass_name", utils.PathSearch("dataclass_name", playbook, nil)),
		d.Set("edit_role", utils.PathSearch("edit_role", playbook, nil)),
		d.Set("owner_id", utils.PathSearch("owner_id", playbook, nil)),
		d.Set("reject_version_id", utils.PathSearch("reject_version_id", playbook, nil)),
		d.Set("unaudited_version_id", utils.PathSearch("unaudited_version_id", playbook, nil)),
		d.Set("user_role", utils.PathSearch("user_role", playbook, nil)),
		d.Set("version", utils.PathSearch("version", playbook, nil)),
		d.Set("version_id", utils.PathSearch("version_id", playbook, nil)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func resourcePlaybookUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	// updatePlaybook: Update the configuration of SecMaster playbook
	var (
		updatePlaybookHttpUrl = "v1/{project_id}/workspaces/{workspace_id}/soc/playbooks/{id}"
		updatePlaybookProduct = "secmaster"
	)
	updatePlaybookClient, err := cfg.NewServiceClient(updatePlaybookProduct, region)
	if err != nil {
		return diag.Errorf("error creating SecMaster client: %s", err)
	}

	updatePlaybookPath := updatePlaybookClient.Endpoint + updatePlaybookHttpUrl
	updatePlaybookPath = strings.ReplaceAll(updatePlaybookPath, "{project_id}", updatePlaybookClient.ProjectID)
	updatePlaybookPath = strings.ReplaceAll(updatePlaybookPath, "{workspace_id}", d.Get("workspace_id").(string))
	updatePlaybookPath = strings.ReplaceAll(updatePlaybookPath, "{id}", d.Id())

	updatePlaybookOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"Content-Type": "application/json"},
	}

	updatePlaybookOpt.JSONBody = utils.RemoveNil(buildUpdatePlaybookBodyParams(d))
	_, err = updatePlaybookClient.Request("PUT", updatePlaybookPath, &updatePlaybookOpt)
	if err != nil {
		return diag.Errorf("error updating SecMaster playbook: %s", err)
	}

	return resourcePlaybookRead(ctx, d, meta)
}

func buildUpdatePlaybookBodyParams(d *schema.ResourceData) map[string]interface{} {
	// `enabled` and `active_version_id` parameters are not updated here, which will produce circular dependencies.
	// These two parameters will be split into a new resource.
	// However, `enabled` need to be filled in the body, otherwise the update interface will not be available.
	bodyParams := map[string]interface{}{
		"name":              utils.ValueIgnoreEmpty(d.Get("name")),
		"description":       d.Get("description"),
		"enabled":           utils.ValueIgnoreEmpty(d.Get("enabled")),
		"active_version_id": utils.ValueIgnoreEmpty(d.Get("active_version_id")),
	}
	return bodyParams
}

func resourcePlaybookDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	workspaceID := d.Get("workspace_id").(string)

	// deletePlaybook: Delete an existing SecMaster playbook
	var (
		deletePlaybookHttpUrl = "v1/{project_id}/workspaces/{workspace_id}/soc/playbooks/{id}"
		deletePlaybookProduct = "secmaster"
	)
	deletePlaybookClient, err := cfg.NewServiceClient(deletePlaybookProduct, region)
	if err != nil {
		return diag.Errorf("error creating SecMaster client: %s", err)
	}

	deletePlaybookPath := deletePlaybookClient.Endpoint + deletePlaybookHttpUrl
	deletePlaybookPath = strings.ReplaceAll(deletePlaybookPath, "{project_id}", deletePlaybookClient.ProjectID)
	deletePlaybookPath = strings.ReplaceAll(deletePlaybookPath, "{workspace_id}", workspaceID)
	deletePlaybookPath = strings.ReplaceAll(deletePlaybookPath, "{id}", d.Id())

	deletePlaybookOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"Content-Type": "application/json"},
	}

	_, err = deletePlaybookClient.Request("DELETE", deletePlaybookPath, &deletePlaybookOpt)
	if err != nil {
		return common.CheckDeletedDiag(d,
			common.ConvertExpected403ErrInto404Err(err, "code", WorkspaceNotFound), "error deleting SecMaster playbook")
	}
	// Successful deletion API call does not guarantee that the resource is successfully deleted.
	// Call the details API to confirm that the resource has been successfully deleted.
	_, err = GetPlaybook(deletePlaybookClient, workspaceID, d.Id())
	if err == nil {
		return diag.Errorf("error deleting SecMaster playbook")
	}

	return nil
}

func GetPlaybook(client *golangsdk.ServiceClient, workspaceId, id string) (interface{}, error) {
	getPlaybookHttpUrl := "v1/{project_id}/workspaces/{workspace_id}/soc/playbooks/{playbook_id}"

	getPlaybookPath := client.Endpoint + getPlaybookHttpUrl
	getPlaybookPath = strings.ReplaceAll(getPlaybookPath, "{project_id}", client.ProjectID)
	getPlaybookPath = strings.ReplaceAll(getPlaybookPath, "{workspace_id}", workspaceId)
	getPlaybookPath = strings.ReplaceAll(getPlaybookPath, "{playbook_id}", id)

	getPlaybookOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"Content-Type": "application/json"},
	}

	getPlaybookResp, err := client.Request("GET", getPlaybookPath, &getPlaybookOpt)
	if err != nil {
		return nil, err
	}

	getPlaybookRespBody, err := utils.FlattenResponse(getPlaybookResp)
	if err != nil {
		return nil, err
	}

	playbook := utils.PathSearch("data", getPlaybookRespBody, nil)
	if playbook == nil {
		return nil, golangsdk.ErrDefault404{}
	}

	return playbook, nil
}

func resourcePlaybookImportState(_ context.Context, d *schema.ResourceData, _ interface{}) ([]*schema.ResourceData, error) {
	parts := strings.Split(d.Id(), "/")
	if len(parts) != 2 {
		return nil, fmt.Errorf("invalid format specified for import id, must be <workspace_id>/<playbook_id>")
	}

	d.SetId(parts[1])
	mErr := multierror.Append(d.Set("workspace_id", parts[0]))

	return []*schema.ResourceData{d}, mErr.ErrorOrNil()
}
