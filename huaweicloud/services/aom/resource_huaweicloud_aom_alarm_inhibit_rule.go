package aom

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

var alarmInhibitRuleNonUpdatableParams = []string{"name", "enterprise_project_id"}

// @API AOM POST /v2/{project_id}/alert/inhibit-rules
// @API AOM GET /v2/{project_id}/alert/inhibit-rules
// @API AOM PUT /v2/{project_id}/alert/inhibit-rules
// @API AOM DELETE /v2/{project_id}/alert/inhibit-rules
func ResourceAlarmInhibitRule() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceAlarmInhibitRuleCreate,
		ReadContext:   resourceAlarmInhibitRuleRead,
		UpdateContext: resourceAlarmInhibitRuleUpdate,
		DeleteContext: resourceAlarmInhibitRuleDelete,

		Importer: &schema.ResourceImporter{
			StateContext: resourceAlarmInhibitRuleImportState,
		},

		CustomizeDiff: config.FlexibleForceNew(alarmInhibitRuleNonUpdatableParams),

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: "The region where the alarm inhibit rule is located.",
			},
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The name of the alarm inhibit rule.",
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The description of the alarm inhibit rule.",
			},
			"binding_group_rule": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The rule name associated with the alarm inhibit rule.",
			},
			"match_v3": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringIsJSON,
				Description:  "The orchestrated alarm inhibit rule definition, in JSON format.",
			},
			"source_matches": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"conditions": {
							Type:        schema.TypeList,
							Required:    true,
							Description: "The serial conditions within a parallel condition group.",
							Elem:        alarmInhibitRuleMatchSchema(),
						},
					},
				},
				Description: "The parallel match conditions for root alerts that suppress other alerts.",
			},
			"target_matches": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"conditions": {
							Type:        schema.TypeList,
							Required:    true,
							Description: "The serial conditions within a parallel condition group.",
							Elem:        alarmInhibitRuleMatchSchema(),
						},
					},
				},
				Description: "The parallel match conditions for target alerts that will be suppressed.",
			},
			"enterprise_project_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "The ID of the enterprise project to which the alarm inhibit rule belongs.",
			},
			// Internal parameter(s).
			"enable_force_new": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"true", "false"}, false),
				Description:  utils.SchemaDesc("", utils.SchemaDescInput{Internal: true}),
			},
		},
	}
}

func alarmInhibitRuleMatchSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"key": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The key of the alarm.",
			},
			"operate": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The match operator for the alarm key.",
			},
			"values": {
				Type:        schema.TypeSet,
				Optional:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "The value list corresponding to the key of the alarm.",
			},
		},
	}
}

func buildAlarmInhibitRuleBodyParams(d *schema.ResourceData) map[string]interface{} {
	return map[string]interface{}{
		"name":               d.Get("name"),
		"desc":               utils.ValueIgnoreEmpty(d.Get("description")),
		"binding_group_rule": utils.ValueIgnoreEmpty(d.Get("binding_group_rule")),
		"match_v3":           utils.StringToJson(d.Get("match_v3").(string)),
		"source_match":       buildAlarmInhibitRuleMatches(d.Get("source_matches").([]interface{})),
		"target_match":       buildAlarmInhibitRuleMatches(d.Get("target_matches").([]interface{})),
	}
}

func buildAlarmInhibitRuleMatches(matches []interface{}) [][]map[string]interface{} {
	if len(matches) == 0 {
		return nil
	}

	result := make([][]map[string]interface{}, 0, len(matches))
	for _, parallelGroup := range matches {
		conditions := utils.PathSearch("conditions", parallelGroup, make([]interface{}, 0)).([]interface{})
		serialConditions := make([]map[string]interface{}, 0, len(conditions))
		for _, condition := range conditions {
			serialConditions = append(serialConditions, map[string]interface{}{
				"key":     utils.PathSearch("key", condition, nil),
				"operate": utils.PathSearch("operate", condition, nil),
				"value": utils.ExpandToStringListBySet(utils.PathSearch("values", condition,
					schema.NewSet(schema.HashString, nil)).(*schema.Set)),
			})
		}

		if len(serialConditions) > 0 {
			result = append(result, serialConditions)
		}
	}

	return result
}

func resourceAlarmInhibitRuleCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		httpUrl = "v2/{project_id}/alert/inhibit-rules"
	)
	client, err := cfg.NewServiceClient("aom", cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating AOM client: %s", err)
	}

	createPath := client.Endpoint + httpUrl
	createPath = strings.ReplaceAll(createPath, "{project_id}", client.ProjectID)

	opt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      buildRequestMoreHeaders(cfg.GetEnterpriseProjectID(d)),
		JSONBody:         utils.RemoveNil(buildAlarmInhibitRuleBodyParams(d)),
		OkCodes:          []int{204},
	}

	_, err = client.Request("POST", createPath, &opt)
	if err != nil {
		return diag.Errorf("error creating alarm inhibit rule: %s", err)
	}

	d.SetId(d.Get("name").(string))

	return resourceAlarmInhibitRuleRead(ctx, d, meta)
}

func GetAlarmInhibitRuleByName(client *golangsdk.ServiceClient, epsId, name string) (interface{}, error) {
	httpUrl := "v2/{project_id}/alert/inhibit-rules"
	path := client.Endpoint + httpUrl
	path = strings.ReplaceAll(path, "{project_id}", client.ProjectID)

	getOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      buildRequestMoreHeaders(epsId),
	}

	resp, err := client.Request("GET", path, &getOpt)
	if err != nil {
		return nil, err
	}

	inhibitRules, err := utils.FlattenResponse(resp)
	if err != nil {
		return nil, err
	}

	inhibitRule := utils.PathSearch(fmt.Sprintf("[?name=='%s']|[0]", name), inhibitRules, nil)
	if inhibitRule == nil {
		return nil, golangsdk.ErrDefault404{
			ErrUnexpectedResponseCode: golangsdk.ErrUnexpectedResponseCode{
				Method:    "GET",
				URL:       "/v2/{project_id}/alert/inhibit-rules",
				RequestId: "NONE",
				Body:      []byte(fmt.Sprintf("The alarm inhibit rule (%s) not found", name)),
			},
		}
	}

	return inhibitRule, nil
}

func resourceAlarmInhibitRuleRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg      = meta.(*config.Config)
		region   = cfg.GetRegion(d)
		ruleName = d.Id()
	)
	client, err := cfg.NewServiceClient("aom", region)
	if err != nil {
		return diag.Errorf("error creating AOM client: %s", err)
	}

	inhibitRule, err := GetAlarmInhibitRuleByName(client, cfg.GetEnterpriseProjectID(d), ruleName)
	if err != nil {
		return common.CheckDeletedDiag(d, err, "error retrieving alarm inhibit rule")
	}

	mErr := multierror.Append(
		d.Set("region", region),
		d.Set("name", utils.PathSearch("name", inhibitRule, nil)),
		d.Set("description", utils.PathSearch("desc", inhibitRule, nil)),
		d.Set("binding_group_rule", utils.PathSearch("binding_group_rule", inhibitRule, nil)),
		d.Set("match_v3", utils.JsonToString(utils.PathSearch("match_v3", inhibitRule, nil))),
		d.Set("source_matches", flattenAlarmInhibitRuleMatches(utils.PathSearch("source_match", inhibitRule,
			make([]interface{}, 0)).([]interface{}))),
		d.Set("target_matches", flattenAlarmInhibitRuleMatches(utils.PathSearch("target_match", inhibitRule,
			make([]interface{}, 0)).([]interface{}))),
		d.Set("enterprise_project_id", utils.PathSearch("enterprise_project_id", inhibitRule, nil)),
	)
	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenAlarmInhibitRuleMatches(matches []interface{}) []map[string]interface{} {
	if len(matches) == 0 {
		return nil
	}

	result := make([]map[string]interface{}, 0, len(matches))
	for _, parallelGroup := range matches {
		serialConditions, ok := parallelGroup.([]interface{})
		if !ok {
			continue
		}

		conditions := make([]map[string]interface{}, 0, len(serialConditions))
		for _, condition := range serialConditions {
			conditions = append(conditions, map[string]interface{}{
				"key":     utils.PathSearch("key", condition, nil),
				"operate": utils.PathSearch("operate", condition, nil),
				"values":  utils.PathSearch("value", condition, nil),
			})
		}

		if len(conditions) > 0 {
			result = append(result, map[string]interface{}{
				"conditions": conditions,
			})
		}
	}

	return result
}

func buildUpdateAlarmInhibitRuleBodyParams(d *schema.ResourceData) map[string]interface{} {
	params := map[string]interface{}{
		"name":               d.Get("name"),
		"desc":               d.Get("description"),
		"binding_group_rule": d.Get("binding_group_rule"),
		"match_v3":           utils.ValueIgnoreEmpty(utils.StringToJson(d.Get("match_v3").(string))),
	}

	if v, ok := d.GetOk("source_matches"); ok {
		params["source_match"] = buildAlarmInhibitRuleMatches(v.([]interface{}))
	}

	if v, ok := d.GetOk("target_matches"); ok {
		params["target_match"] = buildAlarmInhibitRuleMatches(v.([]interface{}))
	}

	return params
}

func resourceAlarmInhibitRuleUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		httpUrl = "v2/{project_id}/alert/inhibit-rules"
	)
	client, err := cfg.NewServiceClient("aom", cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating AOM client: %s", err)
	}

	if d.HasChanges("description", "binding_group_rule", "match_v3", "source_matches", "target_matches") {
		updatePath := client.Endpoint + httpUrl
		updatePath = strings.ReplaceAll(updatePath, "{project_id}", client.ProjectID)

		opt := golangsdk.RequestOpts{
			KeepResponseBody: true,
			MoreHeaders:      buildRequestMoreHeaders(cfg.GetEnterpriseProjectID(d)),
			OkCodes:          []int{204},
			JSONBody:         utils.RemoveNil(buildUpdateAlarmInhibitRuleBodyParams(d)),
		}

		_, err = client.Request("PUT", updatePath, &opt)
		if err != nil {
			return diag.Errorf("error updating alarm inhibit rule: %s", err)
		}
	}

	return resourceAlarmInhibitRuleRead(ctx, d, meta)
}

func resourceAlarmInhibitRuleDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg      = meta.(*config.Config)
		httpUrl  = "v2/{project_id}/alert/inhibit-rules"
		ruleName = d.Get("name").(string)
	)
	client, err := cfg.NewServiceClient("aom", cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating AOM client: %s", err)
	}

	deletePath := client.Endpoint + httpUrl
	deletePath = strings.ReplaceAll(deletePath, "{project_id}", client.ProjectID)
	deleteOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      buildRequestMoreHeaders(cfg.GetEnterpriseProjectID(d)),
		JSONBody:         []string{ruleName},
	}

	_, err = client.Request("DELETE", deletePath, &deleteOpt)
	if err != nil {
		// AOM.08010002: The alarm inhibit rule does not exist
		return common.CheckDeletedDiag(
			d,
			common.ConvertExpected400ErrInto404Err(err, "error_code", "AOM.08010002"),
			fmt.Sprintf("error deleting alarm inhibit rule (%s)", ruleName),
		)
	}
	return nil
}

func resourceAlarmInhibitRuleImportState(_ context.Context, d *schema.ResourceData, _ interface{}) ([]*schema.ResourceData, error) {
	importedId := d.Id()
	parts := strings.Split(importedId, "/")
	switch len(parts) {
	case 1:
		d.SetId(parts[0])
		return []*schema.ResourceData{d}, nil
	case 2:
		d.SetId(parts[0])
		return []*schema.ResourceData{d}, d.Set("enterprise_project_id", parts[1])
	}
	return nil, fmt.Errorf("invalid format specified for import ID, want '<id>' or "+"'<id>/<enterprise_project_id>', but got '%s'",
		importedId)
}
