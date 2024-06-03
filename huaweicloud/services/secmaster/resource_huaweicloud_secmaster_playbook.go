// ---------------------------------------------------------------
// *** AUTO GENERATED CODE ***
// @Product SecMaster
// ---------------------------------------------------------------

package secmaster

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/jmespath/go-jmespath"

	"github.com/chnsz/golangsdk"
	"github.com/chnsz/golangsdk/pagination"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API SecMaster GET /v1/{project_id}/workspaces/{workspace_id}/soc/playbooks
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
			"enabled": {
				Type:        schema.TypeBool,
				Optional:    true,
				Computed:    true,
				Description: `Specifies whether to enable the playbook.`,
			},
			"active_version_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `Specifies the active version ID.`,
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
		OkCodes: []int{
			200,
		},
		MoreHeaders: map[string]string{"Content-Type": "application/json"},
	}

	createPlaybookOpt.JSONBody = utils.RemoveNil(buildCreatePlaybookBodyParams(d))
	createPlaybookResp, err := createPlaybookClient.Request("POST", createPlaybookPath, &createPlaybookOpt)
	if err != nil {
		return diag.Errorf("error creating Playbook: %s", err)
	}

	createPlaybookRespBody, err := utils.FlattenResponse(createPlaybookResp)
	if err != nil {
		return diag.FromErr(err)
	}

	id, err := jmespath.Search("data.id", createPlaybookRespBody)
	if err != nil {
		return diag.Errorf("error creating Playbook: ID is not found in API response")
	}
	d.SetId(id.(string))

	return resourcePlaybookUpdate(ctx, d, meta)
}

func buildCreatePlaybookBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"name":        d.Get("name"),
		"description": utils.ValueIgnoreEmpty(d.Get("description")),
		"enabled":     utils.ValueIgnoreEmpty(d.Get("enabled")),
	}
	return bodyParams
}

func resourcePlaybookRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var mErr *multierror.Error

	// getPlaybook: Query the SecMaster playbook detail
	var (
		getPlaybookHttpUrl = "v1/{project_id}/workspaces/{workspace_id}/soc/playbooks"
		getPlaybookProduct = "secmaster"
	)
	getPlaybookClient, err := cfg.NewServiceClient(getPlaybookProduct, region)
	if err != nil {
		return diag.Errorf("error creating SecMaster client: %s", err)
	}

	getPlaybookPath := getPlaybookClient.Endpoint + getPlaybookHttpUrl
	getPlaybookPath = strings.ReplaceAll(getPlaybookPath, "{project_id}", getPlaybookClient.ProjectID)
	getPlaybookPath = strings.ReplaceAll(getPlaybookPath, "{workspace_id}", d.Get("workspace_id").(string))

	getPlaybookqueryParams := buildGetPlaybookQueryParams(d)
	getPlaybookPath += getPlaybookqueryParams

	getPlaybookResp, err := pagination.ListAllItems(
		getPlaybookClient,
		"offset",
		getPlaybookPath,
		&pagination.QueryOpts{MarkerField: ""})

	if err != nil {
		return common.CheckDeletedDiag(d, err, "error retrieving Playbook")
	}

	getPlaybookRespJson, err := json.Marshal(getPlaybookResp)
	if err != nil {
		return diag.FromErr(err)
	}
	var getPlaybookRespBody interface{}
	err = json.Unmarshal(getPlaybookRespJson, &getPlaybookRespBody)
	if err != nil {
		return diag.FromErr(err)
	}

	jsonPath := fmt.Sprintf("data[?id=='%s']|[0]", d.Id())
	getPlaybookRespBody = utils.PathSearch(jsonPath, getPlaybookRespBody, nil)
	if getPlaybookRespBody == nil {
		return common.CheckDeletedDiag(d, golangsdk.ErrDefault404{}, "no data found")
	}

	mErr = multierror.Append(
		mErr,
		d.Set("region", region),
		d.Set("name", utils.PathSearch("name", getPlaybookRespBody, nil)),
		d.Set("description", utils.PathSearch("description", getPlaybookRespBody, nil)),
		d.Set("enabled", utils.PathSearch("enabled", getPlaybookRespBody, nil)),
		d.Set("approve_role", utils.PathSearch("approve_role", getPlaybookRespBody, nil)),
		d.Set("created_at", utils.PathSearch("create_time", getPlaybookRespBody, nil)),
		d.Set("updated_at", utils.PathSearch("update_time", getPlaybookRespBody, nil)),
		d.Set("dataclass_id", utils.PathSearch("dataclass_id", getPlaybookRespBody, nil)),
		d.Set("dataclass_name", utils.PathSearch("dataclass_name", getPlaybookRespBody, nil)),
		d.Set("edit_role", utils.PathSearch("edit_role", getPlaybookRespBody, nil)),
		d.Set("owner_id", utils.PathSearch("owner_id", getPlaybookRespBody, nil)),
		d.Set("reject_version_id", utils.PathSearch("reject_version_id", getPlaybookRespBody, nil)),
		d.Set("unaudited_version_id", utils.PathSearch("unaudited_version_id", getPlaybookRespBody, nil)),
		d.Set("user_role", utils.PathSearch("user_role", getPlaybookRespBody, nil)),
		d.Set("version", utils.PathSearch("version", getPlaybookRespBody, nil)),
		d.Set("version_id", utils.PathSearch("version_id", getPlaybookRespBody, nil)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func buildGetPlaybookQueryParams(d *schema.ResourceData) string {
	res := ""
	if v, ok := d.GetOk("name"); ok {
		res = fmt.Sprintf("%s&name=%v", res, v)
	}

	if res != "" {
		res = "?" + res[1:]
	}
	return res
}

func resourcePlaybookUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	updatePlaybookChanges := []string{
		"name",
		"description",
		"enabled",
		"active_version_id",
	}

	if d.HasChanges(updatePlaybookChanges...) {
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
			OkCodes: []int{
				200,
			},
			MoreHeaders: map[string]string{"Content-Type": "application/json"},
		}

		updatePlaybookOpt.JSONBody = utils.RemoveNil(buildUpdatePlaybookBodyParams(d))
		_, err = updatePlaybookClient.Request("PUT", updatePlaybookPath, &updatePlaybookOpt)
		if err != nil {
			return diag.Errorf("error updating Playbook: %s", err)
		}
	}
	return resourcePlaybookRead(ctx, d, meta)
}

func buildUpdatePlaybookBodyParams(d *schema.ResourceData) map[string]interface{} {
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
	deletePlaybookPath = strings.ReplaceAll(deletePlaybookPath, "{workspace_id}", d.Get("workspace_id").(string))
	deletePlaybookPath = strings.ReplaceAll(deletePlaybookPath, "{id}", d.Id())

	deletePlaybookOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		OkCodes: []int{
			200,
		},
		MoreHeaders: map[string]string{"Content-Type": "application/json"},
	}

	_, err = deletePlaybookClient.Request("DELETE", deletePlaybookPath, &deletePlaybookOpt)
	if err != nil {
		return diag.Errorf("error deleting Playbook: %s", err)
	}

	return nil
}

func resourcePlaybookImportState(_ context.Context, d *schema.ResourceData, _ interface{}) ([]*schema.ResourceData, error) {
	parts := strings.Split(d.Id(), "/")
	if len(parts) != 2 {
		return nil, fmt.Errorf("invalid format specified for import id, must be <workspace_id>/<playbook_id>")
	}

	d.SetId(parts[1])

	err := d.Set("workspace_id", parts[0])
	if err != nil {
		return nil, err
	}

	return []*schema.ResourceData{d}, nil
}
