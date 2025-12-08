package workspace

import (
	"context"
	"strings"

	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

var desktopMaintenanceBatchManageNonUpdatableParams = []string{
	"desktop_ids",
	"in_maintenance_mode",
}

// @API Workspace PUT /v2/{project_id}/desktops/maintenance-mode
func ResourceDesktopMaintenanceBatchManage() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceDesktopMaintenanceBatchManageCreate,
		ReadContext:   resourceDesktopMaintenanceBatchManageRead,
		UpdateContext: resourceDesktopMaintenanceBatchManageUpdate,
		DeleteContext: resourceDesktopMaintenanceBatchManageDelete,

		CustomizeDiff: config.FlexibleForceNew(desktopMaintenanceBatchManageNonUpdatableParams),

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: `The region where the desktops are located.`,
			},

			// Required parameters.
			"desktop_ids": {
				Type:        schema.TypeList,
				Required:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: `The list of desktop IDs to set maintenance mode.`,
			},
			"in_maintenance_mode": {
				Type:        schema.TypeBool,
				Required:    true,
				Description: `Whether to enter or exit maintenance mode.`,
			},

			// Internal parameters.
			"enable_force_new": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"true", "false"}, false),
				Description:  utils.SchemaDesc("", utils.SchemaDescInput{Internal: true}),
			},
		},
	}
}

func buildDesktopMaintenanceBatchManageBodyParams(d *schema.ResourceData) map[string]interface{} {
	body := map[string]interface{}{
		"desktop_ids":         d.Get("desktop_ids").([]interface{}),
		"in_maintenance_mode": d.Get("in_maintenance_mode").(bool),
	}

	return body
}

func resourceDesktopMaintenanceBatchManageCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		httpUrl = "v2/{project_id}/desktops/maintenance-mode"
	)

	client, err := cfg.NewServiceClient("workspace", region)
	if err != nil {
		return diag.Errorf("error creating Workspace client: %s", err)
	}

	createPath := client.Endpoint + httpUrl
	createPath = strings.ReplaceAll(createPath, "{project_id}", client.ProjectID)

	opt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json",
		},
		JSONBody: buildDesktopMaintenanceBatchManageBodyParams(d),
	}

	_, err = client.Request("PUT", createPath, &opt)
	if err != nil {
		return diag.Errorf("error creating maintenance batch manage: %s", err)
	}

	randomUUID, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(randomUUID)

	return resourceDesktopMaintenanceBatchManageRead(ctx, d, meta)
}

func resourceDesktopMaintenanceBatchManageRead(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return nil
}

func resourceDesktopMaintenanceBatchManageUpdate(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return nil
}

func resourceDesktopMaintenanceBatchManageDelete(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	errorMsg := `This resource is a one-time action resource for batch managing desktops' maintenance mode. Deleting
    this resource will not clear the corresponding request record, but will only remove the resource information
    from the tfstate file.`
	return diag.Diagnostics{
		diag.Diagnostic{
			Severity: diag.Warning,
			Summary:  errorMsg,
		},
	}
}
