package secmaster

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

var nonUpdatableParamsVersionAction = []string{"workspace_id", "playbook_version_id"}

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
				ConflictsWith: []string{"enabled"},
				Description:   `Specifies the playbook version status.`,
			},
			"enabled": {
				Type:          schema.TypeBool,
				Optional:      true,
				ConflictsWith: []string{"status"},
				Description:   `Specifies whether the playbook version is enabled.`,
			},
		},
	}
}

func resourcePlaybookVersionActionCreate(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
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

	return nil
}

func resourcePlaybookVersionActionRead(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return nil
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
