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
	PlaybookVersionNotFound = "SecMaster.20048004"
)

// @API SecMaster GET /v1/{project_id}/workspaces/{workspace_id}/soc/playbooks/versions/{version_id}
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
				Description: `Indicates the status of the playbook version.`,
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
		MoreHeaders:      map[string]string{"Content-Type": "application/json"},
	}

	createPlaybookVersionOpt.JSONBody = utils.RemoveNil(buildCreatePlaybookVersionBodyParams(d))
	createPlaybookVersionResp, err := createPlaybookVersionClient.Request("POST", createPlaybookVersionPath, &createPlaybookVersionOpt)
	if err != nil {
		return diag.Errorf("error creating SecMaster playbook version: %s", err)
	}

	createPlaybookVersionRespBody, err := utils.FlattenResponse(createPlaybookVersionResp)
	if err != nil {
		return diag.FromErr(err)
	}

	id := utils.PathSearch("data.id", createPlaybookVersionRespBody, "").(string)
	if id == "" {
		return diag.Errorf("error creating SecMaster playbook version: ID is not found in API response")
	}
	d.SetId(id)

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

	client, err := cfg.NewServiceClient("secmaster", region)
	if err != nil {
		return diag.Errorf("error creating SecMaster client: %s", err)
	}

	// getPlaybookVersion: Query the SecMaster playbook version detail
	playbookVersion, err := GetPlaybookVersion(client, d.Get("workspace_id").(string), d.Id())
	if err != nil {
		return common.CheckDeletedDiag(d,
			common.ConvertExpected403ErrInto404Err(err, "code", WorkspaceNotFound), "error retrieving SecMaster playbook version")
	}

	mErr := multierror.Append(
		d.Set("region", region),
		d.Set("description", utils.PathSearch("description", playbookVersion, nil)),
		d.Set("dataclass_id", utils.PathSearch("dataclass_id", playbookVersion, nil)),
		d.Set("rule_enable", utils.PathSearch("rule_enable", playbookVersion, nil)),
		d.Set("rule_id", utils.PathSearch("rule_id", playbookVersion, nil)),
		d.Set("trigger_type", utils.PathSearch("trigger_type", playbookVersion, nil)),
		d.Set("dataobject_create", utils.PathSearch("dataobject_create", playbookVersion, nil)),
		d.Set("dataobject_delete", utils.PathSearch("dataobject_delete", playbookVersion, nil)),
		d.Set("dataobject_update", utils.PathSearch("dataobject_update", playbookVersion, nil)),
		d.Set("action_strategy", utils.PathSearch("action_strategy", playbookVersion, nil)),
		d.Set("created_at", utils.PathSearch("create_time", playbookVersion, nil)),
		d.Set("updated_at", utils.PathSearch("update_time", playbookVersion, nil)),
		d.Set("approve_name", utils.PathSearch("approve_name", playbookVersion, nil)),
		d.Set("creator_id", utils.PathSearch("creator_id", playbookVersion, nil)),
		d.Set("dataclass_name", utils.PathSearch("dataclass_name", playbookVersion, nil)),
		d.Set("enabled", utils.PathSearch("enabled", playbookVersion, nil)),
		d.Set("modifier_id", utils.PathSearch("modifier_id", playbookVersion, nil)),
		d.Set("playbook_id", utils.PathSearch("playbook_id", playbookVersion, nil)),
		d.Set("status", utils.PathSearch("status", playbookVersion, nil)),
		d.Set("version", utils.PathSearch("version", playbookVersion, nil)),
		d.Set("version_type", utils.PathSearch("version_type", playbookVersion, nil)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func resourcePlaybookVersionUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	client, err := cfg.NewServiceClient("secmaster", region)
	if err != nil {
		return diag.Errorf("error creating SecMaster client: %s", err)
	}

	// updatePlaybookVersion: Update the configuration of SecMaster playbook version
	bodyParams := utils.RemoveNil(buildUpdatePlaybookVersionBodyParams(d))
	err = updatePlaybookVersion(client, d.Get("workspace_id").(string), d.Id(), bodyParams)
	if err != nil {
		return diag.Errorf("error updating SecMaster playbook version: %s", err)
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
		"status":            utils.ValueIgnoreEmpty(d.Get("status")),
		"enabled":           utils.ValueIgnoreEmpty(d.Get("enabled")),
	}
	return bodyParams
}

func resourcePlaybookVersionDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	client, err := cfg.NewServiceClient("secmaster", region)
	if err != nil {
		return diag.Errorf("error creating SecMaster client: %s", err)
	}

	// Check whether the version is enabled.
	// Before deleting this version, you need to ensure that it is not enabled.
	playbookVersion, err := GetPlaybookVersion(client, d.Get("workspace_id").(string), d.Id())
	if err != nil {
		return common.CheckDeletedDiag(d, err, "error querying SecMaster playbook version")
	}

	if utils.PathSearch("enabled", playbookVersion, false).(bool) {
		bodyParams := backfillUpdateBodyParams(playbookVersion)
		bodyParams["enabled"] = false
		err = updatePlaybookVersion(client, d.Get("workspace_id").(string), d.Id(), bodyParams)
		if err != nil {
			return diag.Errorf("error disabling SecMaster playbook version: %s", err)
		}
	}

	// deletePlaybookVersion: Delete an existing SecMaster playbook version
	deletePlaybookVersionHttpUrl := "v1/{project_id}/workspaces/{workspace_id}/soc/playbooks/versions/{id}"
	deletePlaybookVersionPath := client.Endpoint + deletePlaybookVersionHttpUrl
	deletePlaybookVersionPath = strings.ReplaceAll(deletePlaybookVersionPath, "{project_id}", client.ProjectID)
	deletePlaybookVersionPath = strings.ReplaceAll(deletePlaybookVersionPath, "{workspace_id}", d.Get("workspace_id").(string))
	deletePlaybookVersionPath = strings.ReplaceAll(deletePlaybookVersionPath, "{id}", d.Id())

	deletePlaybookVersionOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"Content-Type": "application/json"},
	}

	_, err = client.Request("DELETE", deletePlaybookVersionPath, &deletePlaybookVersionOpt)
	if err != nil {
		// SecMaster.20048004：the version ID not found
		// SecMaster.20010001: the workspace ID not found
		err = common.ConvertExpected400ErrInto404Err(err, "code", PlaybookVersionNotFound)
		err = common.ConvertExpected403ErrInto404Err(err, "code", WorkspaceNotFound)
		return common.CheckDeletedDiag(d, err, "error deleting SecMaster playbook version")
	}

	return nil
}

func backfillUpdateBodyParams(playbookVersion interface{}) map[string]interface{} {
	bodyParam := map[string]interface{}{
		"description":       utils.PathSearch("description", playbookVersion, "").(string),
		"dataclass_id":      utils.PathSearch("dataclass_id", playbookVersion, nil).(string),
		"rule_enable":       utils.PathSearch("rule_enable", playbookVersion, false).(bool),
		"rule_id":           utils.PathSearch("rule_id", playbookVersion, "").(string),
		"trigger_type":      utils.PathSearch("trigger_type", playbookVersion, "").(string),
		"dataobject_create": utils.PathSearch("dataobject_create", playbookVersion, false).(bool),
		"dataobject_delete": utils.PathSearch("dataobject_delete", playbookVersion, false).(bool),
		"dataobject_update": utils.PathSearch("dataobject_update", playbookVersion, false).(bool),
		"action_strategy":   utils.PathSearch("action_strategy", playbookVersion, "").(string),
		"playbook_id":       utils.PathSearch("playbook_id", playbookVersion, "").(string),
		"enabled":           utils.PathSearch("enabled", playbookVersion, false).(bool),
		"status":            utils.PathSearch("status", playbookVersion, "").(string),
	}

	return bodyParam
}

func GetPlaybookVersion(client *golangsdk.ServiceClient, workspaceId, id string) (interface{}, error) {
	getPlaybookVersionHttpUrl := "v1/{project_id}/workspaces/{workspace_id}/soc/playbooks/versions/{version_id}"
	getPlaybookVersionPath := client.Endpoint + getPlaybookVersionHttpUrl
	getPlaybookVersionPath = strings.ReplaceAll(getPlaybookVersionPath, "{project_id}", client.ProjectID)
	getPlaybookVersionPath = strings.ReplaceAll(getPlaybookVersionPath, "{workspace_id}", workspaceId)
	getPlaybookVersionPath = strings.ReplaceAll(getPlaybookVersionPath, "{version_id}", id)

	getPlaybookVersionOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"Content-Type": "application/json"},
	}

	getPlaybookVersionResp, err := client.Request("GET", getPlaybookVersionPath, &getPlaybookVersionOpt)
	if err != nil {
		return nil, err
	}

	getPlaybookVersionRespBody, err := utils.FlattenResponse(getPlaybookVersionResp)
	if err != nil {
		return nil, err
	}

	playbookVersion := utils.PathSearch("data", getPlaybookVersionRespBody, nil)
	if playbookVersion == nil {
		return nil, golangsdk.ErrDefault404{}
	}

	return playbookVersion, nil
}

func updatePlaybookVersion(client *golangsdk.ServiceClient, workspaceId, id string, bodyParam interface{}) error {
	updatePlaybookVersionHttpUrl := "v1/{project_id}/workspaces/{workspace_id}/soc/playbooks/versions/{id}"
	updatePlaybookVersionPath := client.Endpoint + updatePlaybookVersionHttpUrl
	updatePlaybookVersionPath = strings.ReplaceAll(updatePlaybookVersionPath, "{project_id}", client.ProjectID)
	updatePlaybookVersionPath = strings.ReplaceAll(updatePlaybookVersionPath, "{workspace_id}", workspaceId)
	updatePlaybookVersionPath = strings.ReplaceAll(updatePlaybookVersionPath, "{id}", id)

	updatePlaybookVersionOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"Content-Type": "application/json"},
		JSONBody:         bodyParam,
	}

	_, err := client.Request("PUT", updatePlaybookVersionPath, &updatePlaybookVersionOpt)
	return err
}

func resourcePlaybookVersionImportState(_ context.Context, d *schema.ResourceData, _ interface{}) ([]*schema.ResourceData, error) {
	parts := strings.Split(d.Id(), "/")
	if len(parts) != 3 {
		return nil, fmt.Errorf("invalid format specified for import id, must be <workspace_id>/<playbook_id>/<playbook_version_id>")
	}

	d.SetId(parts[2])

	mErr := multierror.Append(
		d.Set("workspace_id", parts[0]),
		d.Set("playbook_id", parts[1]),
	)

	return []*schema.ResourceData{d}, mErr.ErrorOrNil()
}
