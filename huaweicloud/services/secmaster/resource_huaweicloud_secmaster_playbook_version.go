// ---------------------------------------------------------------
// *** AUTO GENERATED CODE ***
// @Product SecMaster
// ---------------------------------------------------------------

package secmaster

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/jmespath/go-jmespath"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API SecMaster GET /v1/{project_id}/workspaces/{workspace_id}/soc/playbooks/{playbook_id}/versions
// @API SecMaster POST /v1/{project_id}/workspaces/{workspace_id}/soc/playbooks/{playbook_id}/versions
// @API SecMaster PUT /v1/{project_id}/workspaces/{workspace_id}/soc/playbooks/versions/{id}
// @API SecMaster DELETE /v1/{project_id}/workspaces/{workspace_id}/soc/playbooks/versions/{id}
func ResourcePlaybookVersion() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourcePlaybookVersionCreate,
		UpdateContext: resourcePlaybookVersionUpdate,
		ReadContext:   resourcePlaybookVersionRead,
		DeleteContext: resourcePlaybookVersionDelete,
		Importer: &schema.ResourceImporter{
			StateContext: resourcePlaybookVersionImportState,
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
				Description: `Specifies the ID of the workspace to which the playbook version belongs.`,
			},
			"playbook_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: `Specifies playbook ID of the playbook version.`,
			},
			"dataclass_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `Specifies the data class ID of the playbook version.`,
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specifies the description of the playbook version.`,
			},
			"rule_enable": {
				Type:        schema.TypeBool,
				Optional:    true,
				Computed:    true,
				Description: `Specifies whether to enable playbook rule.`,
			},
			"rule_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `Specifies the playbook rule ID.`,
			},
			"trigger_type": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `Specifies the trigger type.`,
			},
			"dataobject_create": {
				Type:        schema.TypeBool,
				Optional:    true,
				Computed:    true,
				Description: `Specifies whether to trigger the playbook when data object is created.`,
			},
			"dataobject_delete": {
				Type:        schema.TypeBool,
				Optional:    true,
				Computed:    true,
				Description: `Specifies whether to trigger the playbook when data object is deleted.`,
			},
			"dataobject_update": {
				Type:        schema.TypeBool,
				Optional:    true,
				Computed:    true,
				Description: `Specifies whether to trigger the playbook when data object is updated.`,
			},
			"action_strategy": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `Specifies the action strategy of the playbook version.`,
			},
			"created_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Indicates the created time of the playbook version.`,
			},
			"updated_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Indicates the updated time of the playbook version.`,
			},
			"approve_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Indicates the approver name.`,
			},
			"creator_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Indicates the creator ID.`,
			},
			"dataclass_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Indicates the data class name.`,
			},
			"enabled": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: `Indicates whether the playbook version is enabled.`,
			},
			"modifier_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Indicates the modifier ID.`,
			},
			"status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Indicates the status of the plaubook.`,
			},
			"version": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Indicates the version number.`,
			},
			"version_type": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: `Indicates the version type.`,
			},
		},
	}
}

func resourcePlaybookVersionCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	// createPlaybookVersion: Create a SecMaster playbook version.
	var (
		createPlaybookVersionHttpUrl = "v1/{project_id}/workspaces/{workspace_id}/soc/playbooks/{playbook_id}/versions"
		createPlaybookVersionProduct = "secmaster"
	)
	createPlaybookVersionClient, err := cfg.NewServiceClient(createPlaybookVersionProduct, region)
	if err != nil {
		return diag.Errorf("error creating SecMaster client: %s", err)
	}

	createPlaybookVersionPath := createPlaybookVersionClient.Endpoint + createPlaybookVersionHttpUrl
	createPlaybookVersionPath = strings.ReplaceAll(createPlaybookVersionPath, "{project_id}", createPlaybookVersionClient.ProjectID)
	createPlaybookVersionPath = strings.ReplaceAll(createPlaybookVersionPath, "{workspace_id}", d.Get("workspace_id").(string))
	createPlaybookVersionPath = strings.ReplaceAll(createPlaybookVersionPath, "{playbook_id}", d.Get("playbook_id").(string))

	createPlaybookVersionOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		OkCodes: []int{
			200,
		},
		MoreHeaders: map[string]string{"Content-Type": "application/json"},
	}

	createPlaybookVersionOpt.JSONBody = utils.RemoveNil(buildCreatePlaybookVersionBodyParams(d))
	createPlaybookVersionResp, err := createPlaybookVersionClient.Request("POST", createPlaybookVersionPath, &createPlaybookVersionOpt)
	if err != nil {
		return diag.Errorf("error creating PlaybookVersion: %s", err)
	}

	createPlaybookVersionRespBody, err := utils.FlattenResponse(createPlaybookVersionResp)
	if err != nil {
		return diag.FromErr(err)
	}

	id, err := jmespath.Search("data.id", createPlaybookVersionRespBody)
	if err != nil {
		return diag.Errorf("error creating PlaybookVersion: ID is not found in API response")
	}
	d.SetId(id.(string))

	return resourcePlaybookVersionRead(ctx, d, meta)
}

func buildCreatePlaybookVersionBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"dataclass_id":      d.Get("dataclass_id"),
		"description":       utils.ValueIgnoreEmpty(d.Get("description")),
		"rule_enable":       utils.ValueIgnoreEmpty(d.Get("rule_enable")),
		"rule_id":           utils.ValueIgnoreEmpty(d.Get("rule_id")),
		"trigger_type":      utils.ValueIgnoreEmpty(d.Get("trigger_type")),
		"dataobject_create": utils.ValueIgnoreEmpty(d.Get("dataobject_create")),
		"dataobject_delete": utils.ValueIgnoreEmpty(d.Get("dataobject_delete")),
		"dataobject_update": utils.ValueIgnoreEmpty(d.Get("dataobject_update")),
		"action_strategy":   utils.ValueIgnoreEmpty(d.Get("action_strategy")),
	}
	return bodyParams
}

func resourcePlaybookVersionRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var mErr *multierror.Error

	// getPlaybookVersion: Query the SecMaster playbook version detail
	var (
		getPlaybookVersionHttpUrl = "v1/{project_id}/workspaces/{workspace_id}/soc/playbooks/{playbook_id}/versions"
		getPlaybookVersionProduct = "secmaster"
	)
	getPlaybookVersionClient, err := cfg.NewServiceClient(getPlaybookVersionProduct, region)
	if err != nil {
		return diag.Errorf("error creating SecMaster client: %s", err)
	}

	getPlaybookVersionPath := getPlaybookVersionClient.Endpoint + getPlaybookVersionHttpUrl
	getPlaybookVersionPath = strings.ReplaceAll(getPlaybookVersionPath, "{project_id}", getPlaybookVersionClient.ProjectID)
	getPlaybookVersionPath = strings.ReplaceAll(getPlaybookVersionPath, "{workspace_id}", d.Get("workspace_id").(string))
	getPlaybookVersionPath = strings.ReplaceAll(getPlaybookVersionPath, "{playbook_id}", d.Get("playbook_id").(string))
	getPlaybookVersionPath += "?limit=1000"

	getPlaybookVersionOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		OkCodes: []int{
			200,
		},
		MoreHeaders: map[string]string{"Content-Type": "application/json"},
	}

	getPlaybookVersionResp, err := getPlaybookVersionClient.Request("GET", getPlaybookVersionPath, &getPlaybookVersionOpt)

	if err != nil {
		return common.CheckDeletedDiag(d, err, "error retrieving PlaybookVersion")
	}

	getPlaybookVersionRespBody, err := utils.FlattenResponse(getPlaybookVersionResp)
	if err != nil {
		return diag.FromErr(err)
	}

	jsonPath := fmt.Sprintf("data[?id=='%s']|[0]", d.Id())
	getPlaybookVersionRespBody = utils.PathSearch(jsonPath, getPlaybookVersionRespBody, nil)
	if getPlaybookVersionRespBody == nil {
		return common.CheckDeletedDiag(d, golangsdk.ErrDefault404{}, "no data found")
	}

	mErr = multierror.Append(
		mErr,
		d.Set("region", region),
		d.Set("description", utils.PathSearch("description", getPlaybookVersionRespBody, nil)),
		d.Set("dataclass_id", utils.PathSearch("dataclass_id", getPlaybookVersionRespBody, nil)),
		d.Set("rule_enable", utils.PathSearch("rule_enable", getPlaybookVersionRespBody, nil)),
		d.Set("rule_id", utils.PathSearch("rule_id", getPlaybookVersionRespBody, nil)),
		d.Set("trigger_type", utils.PathSearch("trigger_type", getPlaybookVersionRespBody, nil)),
		d.Set("dataobject_create", utils.PathSearch("dataobject_create", getPlaybookVersionRespBody, nil)),
		d.Set("dataobject_delete", utils.PathSearch("dataobject_delete", getPlaybookVersionRespBody, nil)),
		d.Set("dataobject_update", utils.PathSearch("dataobject_update", getPlaybookVersionRespBody, nil)),
		d.Set("action_strategy", utils.PathSearch("action_strategy", getPlaybookVersionRespBody, nil)),
		d.Set("created_at", utils.PathSearch("create_time", getPlaybookVersionRespBody, nil)),
		d.Set("updated_at", utils.PathSearch("update_time", getPlaybookVersionRespBody, nil)),
		d.Set("approve_name", utils.PathSearch("approve_name", getPlaybookVersionRespBody, nil)),
		d.Set("creator_id", utils.PathSearch("creator_id", getPlaybookVersionRespBody, nil)),
		d.Set("dataclass_name", utils.PathSearch("dataclass_name", getPlaybookVersionRespBody, nil)),
		d.Set("enabled", utils.PathSearch("enabled", getPlaybookVersionRespBody, nil)),
		d.Set("modifier_id", utils.PathSearch("modifier_id", getPlaybookVersionRespBody, nil)),
		d.Set("playbook_id", utils.PathSearch("playbook_id", getPlaybookVersionRespBody, nil)),
		d.Set("status", utils.PathSearch("status", getPlaybookVersionRespBody, nil)),
		d.Set("version", utils.PathSearch("version", getPlaybookVersionRespBody, nil)),
		d.Set("version_type", utils.PathSearch("version_type", getPlaybookVersionRespBody, nil)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func resourcePlaybookVersionUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	updatePlaybookVersionChanges := []string{
		"description",
		"dataclass_id",
		"rule_enable",
		"rule_id",
		"trigger_type",
		"dataobject_create",
		"dataobject_delete",
		"dataobject_update",
		"action_strategy",
	}

	if d.HasChanges(updatePlaybookVersionChanges...) {
		// updatePlaybookVersion: Update the configuration of SecMaster playbook version
		var (
			updatePlaybookVersionHttpUrl = "v1/{project_id}/workspaces/{workspace_id}/soc/playbooks/versions/{id}"
			updatePlaybookVersionProduct = "secmaster"
		)
		updatePlaybookVersionClient, err := cfg.NewServiceClient(updatePlaybookVersionProduct, region)
		if err != nil {
			return diag.Errorf("error creating SecMaster client: %s", err)
		}

		updatePlaybookVersionPath := updatePlaybookVersionClient.Endpoint + updatePlaybookVersionHttpUrl
		updatePlaybookVersionPath = strings.ReplaceAll(updatePlaybookVersionPath, "{project_id}", updatePlaybookVersionClient.ProjectID)
		updatePlaybookVersionPath = strings.ReplaceAll(updatePlaybookVersionPath, "{workspace_id}", d.Get("workspace_id").(string))
		updatePlaybookVersionPath = strings.ReplaceAll(updatePlaybookVersionPath, "{id}", d.Id())

		updatePlaybookVersionOpt := golangsdk.RequestOpts{
			KeepResponseBody: true,
			OkCodes: []int{
				200,
			},
			MoreHeaders: map[string]string{"Content-Type": "application/json"},
		}

		updatePlaybookVersionOpt.JSONBody = utils.RemoveNil(buildUpdatePlaybookVersionBodyParams(d))
		_, err = updatePlaybookVersionClient.Request("PUT", updatePlaybookVersionPath, &updatePlaybookVersionOpt)
		if err != nil {
			return diag.Errorf("error updating PlaybookVersion: %s", err)
		}
	}
	return resourcePlaybookVersionRead(ctx, d, meta)
}

func buildUpdatePlaybookVersionBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"description":       d.Get("description"),
		"dataclass_id":      utils.ValueIgnoreEmpty(d.Get("dataclass_id")),
		"rule_enable":       utils.ValueIgnoreEmpty(d.Get("rule_enable")),
		"rule_id":           utils.ValueIgnoreEmpty(d.Get("rule_id")),
		"trigger_type":      utils.ValueIgnoreEmpty(d.Get("trigger_type")),
		"dataobject_create": utils.ValueIgnoreEmpty(d.Get("dataobject_create")),
		"dataobject_delete": utils.ValueIgnoreEmpty(d.Get("dataobject_delete")),
		"dataobject_update": utils.ValueIgnoreEmpty(d.Get("dataobject_update")),
		"action_strategy":   utils.ValueIgnoreEmpty(d.Get("action_strategy")),
	}
	return bodyParams
}

func resourcePlaybookVersionDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	// deletePlaybookVersion: Delete an existing SecMaster playbook version
	var (
		deletePlaybookVersionHttpUrl = "v1/{project_id}/workspaces/{workspace_id}/soc/playbooks/versions/{id}"
		deletePlaybookVersionProduct = "secmaster"
	)
	deletePlaybookVersionClient, err := cfg.NewServiceClient(deletePlaybookVersionProduct, region)
	if err != nil {
		return diag.Errorf("error creating SecMaster client: %s", err)
	}

	deletePlaybookVersionPath := deletePlaybookVersionClient.Endpoint + deletePlaybookVersionHttpUrl
	deletePlaybookVersionPath = strings.ReplaceAll(deletePlaybookVersionPath, "{project_id}", deletePlaybookVersionClient.ProjectID)
	deletePlaybookVersionPath = strings.ReplaceAll(deletePlaybookVersionPath, "{workspace_id}", d.Get("workspace_id").(string))
	deletePlaybookVersionPath = strings.ReplaceAll(deletePlaybookVersionPath, "{id}", d.Id())

	deletePlaybookVersionOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		OkCodes: []int{
			200,
		},
		MoreHeaders: map[string]string{"Content-Type": "application/json"},
	}

	_, err = deletePlaybookVersionClient.Request("DELETE", deletePlaybookVersionPath, &deletePlaybookVersionOpt)
	if err != nil {
		return diag.Errorf("error deleting PlaybookVersion: %s", err)
	}

	return nil
}

func resourcePlaybookVersionImportState(_ context.Context, d *schema.ResourceData, _ interface{}) ([]*schema.ResourceData, error) {
	parts := strings.Split(d.Id(), "/")
	if len(parts) != 3 {
		return nil, fmt.Errorf("invalid format specified for import id, must be <workspace_id>/<playbook_id>/<playbook_version_id>")
	}

	d.SetId(parts[2])

	var mErr *multierror.Error
	mErr = multierror.Append(
		mErr,
		d.Set("workspace_id", parts[0]),
		d.Set("playbook_id", parts[1]),
	)

	err := mErr.ErrorOrNil()
	if err != nil {
		return nil, err
	}

	return []*schema.ResourceData{d}, nil
}
