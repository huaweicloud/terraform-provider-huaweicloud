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

var applicationRuleRestrictionSwitchNonUpdatableParams = []string{
	"action",
}

// @API Workspace POST /v1/{project_id}/app-center/app-rules/actions/enable-rule-restriction
// @API Workspace POST /v1/{project_id}/app-center/app-rules/actions/disable-rule-restriction
func ResourceApplicationRuleRestrictionSwitch() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceApplicationRuleRestrictionSwitchCreate,
		ReadContext:   resourceApplicationRuleRestrictionSwitchRead,
		UpdateContext: resourceApplicationRuleRestrictionSwitchUpdate,
		DeleteContext: resourceApplicationRuleRestrictionSwitchDelete,

		CustomizeDiff: config.FlexibleForceNew(applicationRuleRestrictionSwitchNonUpdatableParams),

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: `The region where the application rule restriction to be operated is located.`,
			},

			// Required parameters.
			"action": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The action type for the application rule restriction.`,
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

func resourceApplicationRuleRestrictionSwitchCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		action  = d.Get("action").(string)
		httpUrl = "v1/{project_id}/app-center/app-rules/actions/{action}-rule-restriction"
	)

	client, err := cfg.NewServiceClient("workspace", region)
	if err != nil {
		return diag.Errorf("error creating Workspace client: %s", err)
	}

	createPath := client.Endpoint + httpUrl
	createPath = strings.ReplaceAll(createPath, "{project_id}", client.ProjectID)
	createPath = strings.ReplaceAll(createPath, "{action}", action)

	opt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json",
		},
	}

	_, err = client.Request("POST", createPath, &opt)
	if err != nil {
		return diag.Errorf("error executing application rule restriction switch (%s): %s", action, err)
	}

	uuid, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(uuid)

	return resourceApplicationRuleRestrictionSwitchRead(ctx, d, meta)
}

func resourceApplicationRuleRestrictionSwitchRead(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return nil
}

func resourceApplicationRuleRestrictionSwitchUpdate(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return nil
}

func resourceApplicationRuleRestrictionSwitchDelete(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	errorMsg := `This resource is only a one-time action resource for enabling or disabling the application rule
restriction. Deleting this resource will not clear the corresponding request record, but will only remove the resource
information from the tfstate file.`
	return diag.Diagnostics{
		diag.Diagnostic{
			Severity: diag.Warning,
			Summary:  errorMsg,
		},
	}
}
