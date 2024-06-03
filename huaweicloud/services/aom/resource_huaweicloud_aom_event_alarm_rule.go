// ---------------------------------------------------------------
// *** AUTO GENERATED CODE ***
// @Product AOM
// ---------------------------------------------------------------

package aom

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

// @API AOM POST /v2/{project_id}/event2alarm-rule
// @API AOM PUT /v2/{project_id}/event2alarm-rule
// @API AOM DELETE /v2/{project_id}/event2alarm-rule
// @API AOM GET /v2/{project_id}/event2alarm-rule
func ResourceEventAlarmRule() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceEventAlarmRuleCreate,
		UpdateContext: resourceEventAlarmRuleUpdate,
		ReadContext:   resourceEventAlarmRuleRead,
		DeleteContext: resourceEventAlarmRuleDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: `Specifies the name of the rule.`,
			},
			"alarm_source": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `Specifies the alarm source of the rule.`,
			},
			"select_object": {
				Type:        schema.TypeMap,
				Required:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: `Specifies the select object of the rule.`,
			},
			"alarm_type": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `Specifies the alarm type of the rule.`,
			},
			"trigger_type": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `Specifies the trigger type.`,
			},
			"enabled": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     true,
				Description: `Specifies whether the rule is enabled.`,
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `Specifies the description of the rule.`,
			},
			"action_rule": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `Specifies the action rule name.`,
			},
			"grouping_rule": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `Specifies the route grouping rule name.`,
			},
			"trigger_count": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: `Specifies the accumulated times to trigger the alarm.`,
			},
			"comparison_operator": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `Specifies the comparison condition of alarm.`,
			},
			"period": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: `Specifies the monitoring period in seconds.`,
			},
			"created_at": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: `The creation time.`,
			},
			"updated_at": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: `The last updated time.`,
			},
		},
	}
}

func resourceEventAlarmRuleCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	// createEventAlarmRule: create a Event Alarm Rule.
	var (
		createEventAlarmRuleHttpUrl = "v2/{project_id}/event2alarm-rule"
		createEventAlarmRuleProduct = "aom"
	)
	createEventAlarmRuleClient, err := cfg.NewServiceClient(createEventAlarmRuleProduct, region)
	if err != nil {
		return diag.Errorf("error creating AOM Client: %s", err)
	}

	createEventAlarmRulePath := createEventAlarmRuleClient.Endpoint + createEventAlarmRuleHttpUrl
	createEventAlarmRulePath = strings.ReplaceAll(createEventAlarmRulePath, "{project_id}", createEventAlarmRuleClient.ProjectID)

	createEventAlarmRuleOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		OkCodes: []int{
			204,
		},
	}
	createEventAlarmRuleOpt.JSONBody = utils.RemoveNil(buildEventAlarmRuleBodyParams(d, cfg))
	_, err = createEventAlarmRuleClient.Request("POST", createEventAlarmRulePath, &createEventAlarmRuleOpt)
	if err != nil {
		return diag.Errorf("error creating EventAlarmRule: %s", err)
	}

	d.SetId(d.Get("name").(string))

	return resourceEventAlarmRuleRead(ctx, d, meta)
}

func buildEventAlarmRuleBodyParams(d *schema.ResourceData, cfg *config.Config) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"name":              utils.ValueIgnoreEmpty(d.Get("name")),
		"description":       utils.ValueIgnoreEmpty(d.Get("description")),
		"alarm_type":        utils.ValueIgnoreEmpty(d.Get("alarm_type")),
		"action_rule":       utils.ValueIgnoreEmpty(d.Get("action_rule")),
		"route_group_rule":  utils.ValueIgnoreEmpty(d.Get("grouping_rule")),
		"enable":            utils.ValueIgnoreEmpty(d.Get("enabled")),
		"resource_provider": utils.ValueIgnoreEmpty(d.Get("alarm_source")),
		"metadata":          utils.ValueIgnoreEmpty(d.Get("select_object")),
		"trigger_policies":  buildEventAlarmRuleTriggerPolicies(d),
		"user_id":           cfg.GetProjectID(cfg.GetRegion(d)),
	}
	return bodyParams
}

func buildEventAlarmRuleTriggerPolicies(d *schema.ResourceData) []map[string]interface{} {
	bodyParams := map[string]interface{}{
		"trigger_type": utils.ValueIgnoreEmpty(d.Get("trigger_type")),
		"count":        utils.ValueIgnoreEmpty(d.Get("trigger_count")),
		"operator":     utils.ValueIgnoreEmpty(d.Get("comparison_operator")),
		"period":       utils.ValueIgnoreEmpty(d.Get("period")),
	}

	return []map[string]interface{}{bodyParams}
}

func resourceEventAlarmRuleRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var mErr *multierror.Error

	// getEventAlarmRule: Query the Event Alarm Rule
	var (
		getEventAlarmRuleHttpUrl = "v2/{project_id}/event2alarm-rule"
		getEventAlarmRuleProduct = "aom"
	)
	getEventAlarmRuleClient, err := cfg.NewServiceClient(getEventAlarmRuleProduct, region)
	if err != nil {
		return diag.Errorf("error creating AOM Client: %s", err)
	}

	getEventAlarmRulePath := getEventAlarmRuleClient.Endpoint + getEventAlarmRuleHttpUrl
	getEventAlarmRulePath = strings.ReplaceAll(getEventAlarmRulePath, "{project_id}", getEventAlarmRuleClient.ProjectID)

	getEventAlarmRuleOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		OkCodes: []int{
			200,
		},
	}

	getEventAlarmRuleOpt.MoreHeaders = map[string]string{
		"Content-Type": "application/json",
	}

	getEventAlarmRuleResp, err := getEventAlarmRuleClient.Request("GET", getEventAlarmRulePath, &getEventAlarmRuleOpt)
	if err != nil {
		return common.CheckDeletedDiag(d, err, "error retrieving EventAlarmRule")
	}

	getEventAlarmRuleRespBody, err := utils.FlattenResponse(getEventAlarmRuleResp)
	if err != nil {
		return diag.FromErr(err)
	}

	jsonPath := fmt.Sprintf("[?name=='%s']|[0]", d.Id())
	rule := utils.PathSearch(jsonPath, getEventAlarmRuleRespBody, nil)
	if rule == nil {
		return common.CheckDeletedDiag(d, golangsdk.ErrDefault404{}, "error retrieving EventAlarmRule")
	}

	mErr = multierror.Append(
		mErr,
		d.Set("region", region),
		d.Set("name", utils.PathSearch("name", rule, nil)),
		d.Set("description", utils.PathSearch("description", rule, nil)),
		d.Set("alarm_type", utils.PathSearch("alarm_type", rule, nil)),
		d.Set("action_rule", utils.PathSearch("action_rule", rule, nil)),
		d.Set("grouping_rule", utils.PathSearch("route_group_rule", rule, nil)),
		d.Set("enabled", utils.PathSearch("enable", rule, nil)),
		d.Set("alarm_source", utils.PathSearch("resource_provider", rule, nil)),
		d.Set("select_object", utils.PathSearch("metadata", rule, nil)),
		d.Set("trigger_type", utils.PathSearch("trigger_policies[0].trigger_type", rule, nil)),
		d.Set("trigger_count", utils.PathSearch("trigger_policies[0].count", rule, nil)),
		d.Set("comparison_operator", utils.PathSearch("trigger_policies[0].operator", rule, nil)),
		d.Set("period", utils.PathSearch("trigger_policies[0].period", rule, nil)),
		d.Set("created_at", utils.PathSearch("create_time", rule, nil)),
		d.Set("updated_at", utils.PathSearch("update_time", rule, nil)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func resourceEventAlarmRuleUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	updateEventAlarmRuleChanges := []string{
		"description",
		"alarm_type",
		"action_rule",
		"grouping_rule",
		"enabled",
		"alarm_source",
		"select_object",
		"trigger_type",
		"trigger_count",
		"comparison_operator",
		"period",
	}

	if d.HasChanges(updateEventAlarmRuleChanges...) {
		// updateEventAlarmRule: update the Event Alarm Rule
		var (
			updateEventAlarmRuleHttpUrl = "v2/{project_id}/event2alarm-rule"
			updateEventAlarmRuleProduct = "aom"
		)
		updateEventAlarmRuleClient, err := cfg.NewServiceClient(updateEventAlarmRuleProduct, region)
		if err != nil {
			return diag.Errorf("error creating AOM Client: %s", err)
		}

		updateEventAlarmRulePath := updateEventAlarmRuleClient.Endpoint + updateEventAlarmRuleHttpUrl
		updateEventAlarmRulePath = strings.ReplaceAll(updateEventAlarmRulePath, "{project_id}", updateEventAlarmRuleClient.ProjectID)
		updateEventAlarmRuleOpt := golangsdk.RequestOpts{
			KeepResponseBody: true,
			OkCodes: []int{
				204,
			},
		}
		updateEventAlarmRuleOpt.JSONBody = utils.RemoveNil(buildEventAlarmRuleBodyParams(d, cfg))
		_, err = updateEventAlarmRuleClient.Request("PUT", updateEventAlarmRulePath, &updateEventAlarmRuleOpt)
		if err != nil {
			return diag.Errorf("error updating EventAlarmRule: %s", err)
		}
	}
	return resourceEventAlarmRuleRead(ctx, d, meta)
}

func resourceEventAlarmRuleDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	// deleteEventAlarmRule: missing operation notes
	var (
		deleteEventAlarmRuleHttpUrl = "v2/{project_id}/event2alarm-rule"
		deleteEventAlarmRuleProduct = "aom"
	)
	deleteEventAlarmRuleClient, err := cfg.NewServiceClient(deleteEventAlarmRuleProduct, region)
	if err != nil {
		return diag.Errorf("error creating AOM Client: %s", err)
	}

	deleteEventAlarmRulePath := deleteEventAlarmRuleClient.Endpoint + deleteEventAlarmRuleHttpUrl
	deleteEventAlarmRulePath = strings.ReplaceAll(deleteEventAlarmRulePath, "{project_id}", deleteEventAlarmRuleClient.ProjectID)

	deleteEventAlarmRuleOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		OkCodes: []int{
			204,
		},
	}

	deleteEventAlarmRuleOpt.JSONBody = []string{d.Id()}
	_, err = deleteEventAlarmRuleClient.Request("DELETE", deleteEventAlarmRulePath, &deleteEventAlarmRuleOpt)
	if err != nil {
		return diag.Errorf("error deleting EventAlarmRule: %s", err)
	}

	return nil
}
