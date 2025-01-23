package secmaster

import (
	"context"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

var nonUpdatableParamsVersionAction = []string{"workspace_id", "version_id", "status", "enabled"}

// @API SecMaster GET /v1/{project_id}/workspaces/{workspace_id}/soc/playbooks/{playbook_id}/versions
// @API SecMaster PUT /v1/{project_id}/workspaces/{workspace_id}/soc/playbooks/versions/{id}
func ResourcePlaybookVersionAction() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourcePlaybookVersionActionCreate,
		UpdateContext: resourcePlaybookVersionActionUpdate,
		ReadContext:   resourcePlaybookVersionActionRead,
		DeleteContext: resourcePlaybookVersionActionDelete,

		CustomizeDiff: config.FlexibleForceNew(nonUpdatableParamsVersionAction),

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
				Description: `Specifies the ID of the workspace to which the playbook version belongs.`,
			},
			"version_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `Specifies the playbook version ID.`,
			},
			"status": {
				Type:          schema.TypeString,
				Optional:      true,
				Computed:      true,
				ConflictsWith: []string{"enabled"},
				Description:   `Specifies the playbook version status.`,
			},
			"enabled": {
				Type:          schema.TypeBool,
				Optional:      true,
				Computed:      true,
				ConflictsWith: []string{"status"},
				Description:   `Specifies whether the playbook version is enabled.`,
			},
			"version": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The playbook version.`,
			},
			"type": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: `The playbook version type.`,
			},
			"description": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The description.`,
			},
			"created_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The creation time.`,
			},
			"data_object_create": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: `Whether to trigger a playbook when a data object is created.`,
			},
			"data_class_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The data class ID.`,
			},
			"playbook_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The playbook ID.`,
			},
			"trigger_type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The triggering type.`,
			},
			"modifier_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The ID of the user who updated the information.`,
			},
			"project_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The project ID.`,
			},
			"rule_enabled": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: `Whether the filtering rule is enabled.`,
			},
			"data_object_delete": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: `Whether to trigger a playbook when a data object is deleted.`,
			},
			"data_object_update": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: `Whether to trigger a playbook when a data object is updated.`,
			},
			"rule_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The rule ID.`,
			},
			"data_class_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The data class name.`,
			},
			"updated_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The update time.`,
			},
			"creator_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The creator ID.`,
			},
			"action_strategy": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The execution policy.`,
			},
			"approve_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The reviewer.`,
			},
			"enable_force_new": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"true", "false"}, false),
				Description:  utils.SchemaDesc("", utils.SchemaDescInput{Internal: true}),
			},
		},
	}
}

func resourcePlaybookVersionActionCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	playbookVersionID := d.Get("version_id").(string)

	client, err := cfg.NewServiceClient("secmaster", region)
	if err != nil {
		return diag.Errorf("error creating SecMaster client: %s", err)
	}

	// getPlaybookVersion: Query the SecMaster playbook version detail
	playbookVersion, err := GetPlaybookVersion(client, d.Get("workspace_id").(string), playbookVersionID)
	if err != nil {
		return diag.FromErr(err)
	}

	enabled := utils.PathSearch("enabled", playbookVersion, false).(bool)

	bodyParams := backfillUpdateBodyParams(playbookVersion)
	if v, ok := d.GetOk("status"); ok {
		if enabled {
			return diag.Errorf("this version has been activated, and cannot be edited")
		}
		bodyParams["status"] = v.(string)
	} else {
		inputEnabled := d.Get("enabled").(bool)
		if enabled == inputEnabled {
			return diag.Errorf("this version is already active or inactive and there is no need to perform this action again")
		}
		bodyParams["enabled"] = inputEnabled
	}

	// updatePlaybookVersion: Update the configuration of SecMaster playbook version
	err = updatePlaybookVersion(client, d.Get("workspace_id").(string), playbookVersionID, bodyParams)
	if err != nil {
		return diag.Errorf("error updating SecMaster playbook version: %s", err)
	}

	d.SetId(playbookVersionID)

	return resourcePlaybookVersionActionRead(ctx, d, meta)
}

func resourcePlaybookVersionActionRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
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
		d.Set("data_class_id", utils.PathSearch("dataclass_id", playbookVersion, nil)),
		d.Set("rule_enabled", utils.PathSearch("rule_enable", playbookVersion, nil)),
		d.Set("rule_id", utils.PathSearch("rule_id", playbookVersion, nil)),
		d.Set("trigger_type", utils.PathSearch("trigger_type", playbookVersion, nil)),
		d.Set("data_object_create", utils.PathSearch("dataobject_create", playbookVersion, nil)),
		d.Set("data_object_delete", utils.PathSearch("dataobject_delete", playbookVersion, nil)),
		d.Set("data_object_update", utils.PathSearch("dataobject_update", playbookVersion, nil)),
		d.Set("action_strategy", utils.PathSearch("action_strategy", playbookVersion, nil)),
		d.Set("created_at", utils.PathSearch("create_time", playbookVersion, nil)),
		d.Set("updated_at", utils.PathSearch("update_time", playbookVersion, nil)),
		d.Set("approve_name", utils.PathSearch("approve_name", playbookVersion, nil)),
		d.Set("creator_id", utils.PathSearch("creator_id", playbookVersion, nil)),
		d.Set("data_class_name", utils.PathSearch("dataclass_name", playbookVersion, nil)),
		d.Set("enabled", utils.PathSearch("enabled", playbookVersion, nil)),
		d.Set("modifier_id", utils.PathSearch("modifier_id", playbookVersion, nil)),
		d.Set("playbook_id", utils.PathSearch("playbook_id", playbookVersion, nil)),
		d.Set("status", utils.PathSearch("status", playbookVersion, nil)),
		d.Set("version", utils.PathSearch("version", playbookVersion, nil)),
		d.Set("type", utils.PathSearch("version_type", playbookVersion, nil)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func resourcePlaybookVersionActionUpdate(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return nil
}

func resourcePlaybookVersionActionDelete(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	errorMsg := `This resource is only a one-time action resource for playbook version action resource. Deleting this resource will not change
		the status of the currently playbook version action resource, but will only remove the resource information from the tfstate file.`
	return diag.Diagnostics{
		diag.Diagnostic{
			Severity: diag.Warning,
			Summary:  errorMsg,
		},
	}
}
