package ces

import (
	"context"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

var oneClickAlarmRulePolicyActionNonUpdatableFields = []string{
	"one_click_alarm_id", "alarm_id", "alarm_policy_ids", "enabled",
}

// @API CES PUT /v2/{project_id}/one-click-alarms/{one_click_alarm_id}/alarms/{alarm_id}/policies/action
func ResourceOneClickAlarmRulePolicyAction() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceOneClickAlarmRulePolicyActionCreate,
		UpdateContext: resourceOneClickAlarmRulePolicyActionUpdate,
		ReadContext:   resourceOneClickAlarmRulePolicyActionRead,
		DeleteContext: resourceOneClickAlarmRulePolicyActionDelete,

		CustomizeDiff: config.FlexibleForceNew(oneClickAlarmRulePolicyActionNonUpdatableFields),

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"one_click_alarm_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"alarm_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"alarm_policy_ids": {
				Type:     schema.TypeSet,
				Required: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"enabled": {
				Type:     schema.TypeBool,
				Required: true,
			},
			"enable_force_new": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"true", "false"}, false),
				Description:  utils.SchemaDesc("", utils.SchemaDescInput{Internal: true}),
			},
			"success_alarm_policy_ids": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
		},
	}
}

func resourceOneClickAlarmRulePolicyActionCreate(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var (
		httpUrl = "v2/{project_id}/one-click-alarms/{one_click_alarm_id}/alarms/{alarm_id}/policies/action"
		product = "ces"
	)
	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating CES client: %s", err)
	}

	createPath := client.Endpoint + httpUrl
	createPath = strings.ReplaceAll(createPath, "{project_id}", client.ProjectID)
	createPath = strings.ReplaceAll(createPath, "{one_click_alarm_id}", d.Get("one_click_alarm_id").(string))
	createPath = strings.ReplaceAll(createPath, "{alarm_id}", d.Get("alarm_id").(string))

	createOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		JSONBody:         utils.RemoveNil(buildOneClickAlarmRulePolicyActionBodyParams(d)),
	}

	createResp, err := client.Request("PUT", createPath, &createOpt)
	if err != nil {
		return diag.Errorf("error creating CES batch enable or disable alarm rule policies: %s", err)
	}
	createRespBody, err := utils.FlattenResponse(createResp)
	if err != nil {
		return diag.Errorf("error flattening batch enable or disable alarm rule policies: %s", err)
	}

	id, err := uuid.GenerateUUID()
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(id)

	mErr := multierror.Append(nil,
		d.Set("success_alarm_policy_ids", utils.PathSearch("alarm_policy_ids", createRespBody, nil)),
		d.Set("region", region),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func buildOneClickAlarmRulePolicyActionBodyParams(d *schema.ResourceData) map[string]interface{} {
	param := map[string]interface{}{
		"alarm_policy_ids": d.Get("alarm_policy_ids").(*schema.Set).List(),
		"enabled":          d.Get("enabled"),
	}

	return param
}

func resourceOneClickAlarmRulePolicyActionRead(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return nil
}

func resourceOneClickAlarmRulePolicyActionUpdate(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return nil
}

func resourceOneClickAlarmRulePolicyActionDelete(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	errorMsg := "Deleting the batch enable or disable alarm policies in alarm rules for one service in one-click" +
		" monitoring resource is not supported. The batch enable or disable alarm policies in alarm rules for one" +
		" service in one-click monitoring resource is only removed from the state."
	return diag.Diagnostics{
		diag.Diagnostic{
			Severity: diag.Warning,
			Summary:  errorMsg,
		},
	}
}
