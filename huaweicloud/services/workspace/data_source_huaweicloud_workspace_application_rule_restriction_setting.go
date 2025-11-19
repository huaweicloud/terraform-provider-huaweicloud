package workspace

import (
	"context"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API Workspace GET /v1/{project_id}/app-center/app-rules/actions/get-rule-restriction
func DataSourceApplicationRuleRestrictionSetting() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceApplicationRuleRestrictionSettingRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `The region where the application rule restriction setting is located.`,
			},
			"app_restrict_rule_switch": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: `Whether the application restriction rule switch is enabled.`,
			},
			"app_control_mode": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: `The application control mode.`,
			},
			"app_periodic_switch": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: `Whether the periodic monitoring switch is enabled.`,
			},
			"app_periodic_interval": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: `The periodic monitoring interval time, in minutes.`,
			},
			"app_force_kill_proc_switch": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: `Whether the force kill application switch is enabled.`,
			},
		},
	}
}

func dataSourceApplicationRuleRestrictionSettingRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		httpUrl = "v1/{project_id}/app-center/app-rules/actions/get-rule-restriction"
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

	requestResp, err := client.Request("GET", getPath, &opt)
	if err != nil {
		return diag.Errorf("error querying application rule restriction setting: %s", err)
	}

	respBody, err := utils.FlattenResponse(requestResp)
	if err != nil {
		return diag.Errorf("error flattening application rule restriction setting response: %s", err)
	}

	randomUUID, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(randomUUID)

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("app_restrict_rule_switch", utils.PathSearch("app_restrict_rule_switch", respBody, nil)),
		d.Set("app_control_mode", utils.PathSearch("app_control_mode", respBody, nil)),
		d.Set("app_periodic_switch", utils.PathSearch("app_periodic_switch", respBody, nil)),
		d.Set("app_periodic_interval", utils.PathSearch("app_periodic_interval", respBody, nil)),
		d.Set("app_force_kill_proc_switch", utils.PathSearch("app_force_kill_proc_switch", respBody, nil)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}
