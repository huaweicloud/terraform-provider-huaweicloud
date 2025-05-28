package sms

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

var migrateProjectDefaultNonUpdatableParams = []string{"mig_project_id"}

// @API SMS PUT /v3/migprojects/{mig_project_id}/default
func ResourceMigrateProjectDefault() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceMigrateProjectDefaultCreate,
		ReadContext:   resourceMigrateProjectDefaultRead,
		UpdateContext: resourceMigrateProjectDefaultUpdate,
		DeleteContext: resourceMigrateProjectDefaultDelete,

		CustomizeDiff: config.FlexibleForceNew(migrateProjectDefaultNonUpdatableParams),

		Schema: map[string]*schema.Schema{
			"mig_project_id": {
				Type:     schema.TypeString,
				Required: true,
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

func resourceMigrateProjectDefaultCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	client, err := cfg.SmsV3Client(cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating SMS client: %s", err)
	}

	migrateProjectID := d.Get("mig_project_id").(string)
	err = defaultMigrateProject(client, migrateProjectID)
	if err != nil {
		return nil
	}

	d.SetId(migrateProjectID)

	return resourceMigrateProjectDefaultRead(ctx, d, meta)
}

func resourceMigrateProjectDefaultRead(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return nil
}

func resourceMigrateProjectDefaultUpdate(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return nil
}

func resourceMigrateProjectDefaultDelete(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	errorMsg := "Deleting the set default migration project resource is not supported." +
		" The set default migration project resource is only removed from the state."
	return diag.Diagnostics{
		diag.Diagnostic{
			Severity: diag.Warning,
			Summary:  errorMsg,
		},
	}
}
