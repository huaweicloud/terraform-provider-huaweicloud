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

var applicationRuleRestrictionSettingNonUpdatableParams = []string{
	"app_restrict_rule_switch",
	"app_control_mode",
	"app_periodic_switch",
	"app_periodic_interval",
	"app_force_kill_proc_switch",
}

// @API Workspace POST /v1/{project_id}/app-center/app-rules/actions/set-rule-restriction
func ResourceApplicationRuleRestrictionSetting() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceApplicationRuleRestrictionSettingCreate,
		ReadContext:   resourceApplicationRuleRestrictionSettingRead,
		UpdateContext: resourceApplicationRuleRestrictionSettingUpdate,
		DeleteContext: resourceApplicationRuleRestrictionSettingDelete,

		CustomizeDiff: config.FlexibleForceNew(applicationRuleRestrictionSettingNonUpdatableParams),

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: `The region where the application rule restriction is located.`,
			},

			// Required parameters.
			"app_restrict_rule_switch": {
				Type:        schema.TypeBool,
				Required:    true,
				Description: `Whether to enable the application restriction rule switch.`,
			},

			// Optional parameters.
			"app_control_mode": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: `The application control mode.`,
			},
			"app_periodic_switch": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: `Whether to enable the periodic monitoring switch.`,
			},
			"app_periodic_interval": {
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
				Description: `The periodic monitoring interval time, in minutes.`,
			},
			"app_force_kill_proc_switch": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: `Whether to enable the force kill application switch.`,
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

func buildApplicationRuleRestrictionSettingBodyParams(d *schema.ResourceData) map[string]interface{} {
	return map[string]interface{}{
		"app_restrict_rule_switch":   d.Get("app_restrict_rule_switch").(bool),
		"app_control_mode":           d.Get("app_control_mode").(int),
		"app_periodic_switch":        d.Get("app_periodic_switch").(bool),
		"app_periodic_interval":      utils.ValueIgnoreEmpty(d.Get("app_periodic_interval").(int)),
		"app_force_kill_proc_switch": d.Get("app_force_kill_proc_switch").(bool),
	}
}

func resourceApplicationRuleRestrictionSettingCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		httpUrl = "v1/{project_id}/app-center/app-rules/actions/set-rule-restriction"
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
		JSONBody: buildApplicationRuleRestrictionSettingBodyParams(d),
	}

	_, err = client.Request("POST", createPath, &opt)
	if err != nil {
		return diag.Errorf("error setting application rule restriction: %s", err)
	}

	randUUID, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(randUUID)

	return resourceApplicationRuleRestrictionSettingRead(ctx, d, meta)
}

func resourceApplicationRuleRestrictionSettingRead(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return nil
}

func resourceApplicationRuleRestrictionSettingUpdate(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return nil
}

func resourceApplicationRuleRestrictionSettingDelete(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	errorMsg := `This resource is only a one-time action resource for setting the application rule restriction. Deleting this resource will
not clear the corresponding request record, but will only remove the resource information from the tfstate file.`
	return diag.Diagnostics{
		diag.Diagnostic{
			Severity: diag.Warning,
			Summary:  errorMsg,
		},
	}
}
