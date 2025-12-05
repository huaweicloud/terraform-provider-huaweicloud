package workspace

import (
	"context"
	"strings"

	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

var appActionNonUpdatableParams = []string{
	"app_restrict_rule_switch",
}

// @API Workspace PATCH /v1/{project_id}/app-center/profiles
// @API Workspace GET /v1/{project_id}/app-center/profiles
func ResourceAppAction() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceAppActionCreate,
		ReadContext:   resourceAppActionRead,
		UpdateContext: resourceAppActionUpdate,
		DeleteContext: resourceAppActionDelete,

		CustomizeDiff: config.FlexibleForceNew(appActionNonUpdatableParams),

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: `The region where the tenant profiles are located.`,
			},
			"app_restrict_rule_switch": {
				Type:        schema.TypeBool,
				Required:    true,
				Description: `Whether to enable the application restriction rule switch.`,
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

func buildAppActionBodyParams(d *schema.ResourceData) map[string]interface{} {
	body := map[string]interface{}{
		"app_restrict_rule_switch": d.Get("app_restrict_rule_switch"),
	}
	return body
}

func resourceAppActionCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		httpUrl = "v1/{project_id}/app-center/profiles"
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
		JSONBody: buildAppActionBodyParams(d),
	}

	_, err = client.Request("PATCH", createPath, &opt)
	if err != nil {
		return diag.Errorf("error creating app action: %s", err)
	}

	randomUUID, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(randomUUID)

	return resourceAppActionRead(ctx, d, meta)
}

func resourceAppActionRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		httpUrl = "v1/{project_id}/app-center/profiles"
	)

	client, err := cfg.NewServiceClient("workspace", region)
	if err != nil {
		return diag.Errorf("error creating Workspace client: %s", err)
	}

	getPath := client.Endpoint + httpUrl
	getPath = strings.ReplaceAll(getPath, "{project_id}", client.ProjectID)

	opt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json",
		},
	}

	response, err := client.Request("GET", getPath, &opt)
	if err != nil {
		return common.CheckDeletedDiag(d, err, "error retrieving app action")
	}

	respBody, err := utils.FlattenResponse(response)
	if err != nil {
		return diag.FromErr(err)
	}

	appRestrictRuleSwitch := utils.PathSearch("app_restrict_rule_switch", respBody, false)
	mErr := d.Set("app_restrict_rule_switch", appRestrictRuleSwitch.(bool))
	if mErr != nil {
		return diag.Errorf("error setting app_restrict_rule_switch: %s", mErr)
	}

	return nil
}

func resourceAppActionUpdate(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return nil
}

func resourceAppActionDelete(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	errorMsg := `This resource is used to manage the tenant profile settings. Deleting this resource will not 
	clear the corresponding configuration, but will only remove the resource information from the tfstate file.`
	return diag.Diagnostics{
		diag.Diagnostic{
			Severity: diag.Warning,
			Summary:  errorMsg,
		},
	}
}
