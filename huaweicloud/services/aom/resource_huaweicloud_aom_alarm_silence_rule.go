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

// @API AOM POST /v2/{project_id}/alert/mute-rules
// @API AOM PUT /v2/{project_id}/alert/mute-rules
// @API AOM DELETE /v2/{project_id}/alert/mute-rules
// @API AOM GET /v2/{project_id}/alert/mute-rules
func ResourceAlarmSilenceRule() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceAlarmSilenceRuleCreate,
		UpdateContext: resourceAlarmSilenceRuleUpdate,
		ReadContext:   resourceAlarmSilenceRuleRead,
		DeleteContext: resourceAlarmSilenceRuleDelete,
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
				Description: `Specifies the rule name.`,
			},
			"time_zone": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: `Specifies the time zone.`,
			},
			"silence_time": {
				Type:     schema.TypeList,
				MaxItems: 1,
				Elem:     AlarmSilenceRuleSilenceTimeSchema(),
				Required: true,
			},
			"silence_conditions": {
				Type:        schema.TypeList,
				Elem:        AlarmSilenceRuleSilenceConditionsSchema(),
				Required:    true,
				Description: `Specifies the silence conditions of the rule.`,
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `Specifies the description.`,
			},
			"created_at": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: `The creation time.`,
			},
			"updated_at": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: `The last update time.`,
			},
		},
	}
}

func AlarmSilenceRuleSilenceTimeSchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"type": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `Specifies the effective time type of the silence rule.`,
			},
			"starts_at": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: `Specifies the start time of the silence rule.`,
			},
			"ends_at": {
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
				Description: `Specifies the end time of the silence rule.`,
			},
			"scope": {
				Type:        schema.TypeList,
				Elem:        &schema.Schema{Type: schema.TypeInt},
				Optional:    true,
				Computed:    true,
				Description: `Specifies the silence time of the rule.`,
			},
		},
	}
	return &sc
}

func AlarmSilenceRuleSilenceConditionsSchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"conditions": {
				Type:        schema.TypeList,
				Elem:        AlarmSilenceRuleconditionsSchema(),
				Required:    true,
				Description: `Specifies the serial conditions.`,
			},
		},
	}
	return &sc
}

func AlarmSilenceRuleconditionsSchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"key": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `Specifies the key of the match condition.`,
			},
			"operate": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `Specifies the operate of the match condition.`,
			},
			"value": {
				Type:        schema.TypeList,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Optional:    true,
				Computed:    true,
				Description: `Specifies the value list of the match condition.`,
			},
		},
	}
	return &sc
}

func resourceAlarmSilenceRuleCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	// createAlarmSilenceRule: create a Alarm Silence Rule.
	var (
		createAlarmSilenceRuleHttpUrl = "v2/{project_id}/alert/mute-rules"
		createAlarmSilenceRuleProduct = "aom"
	)
	createAlarmSilenceRuleClient, err := cfg.NewServiceClient(createAlarmSilenceRuleProduct, region)
	if err != nil {
		return diag.Errorf("error creating AOM Client: %s", err)
	}

	createAlarmSilenceRulePath := createAlarmSilenceRuleClient.Endpoint + createAlarmSilenceRuleHttpUrl
	createAlarmSilenceRulePath = strings.ReplaceAll(createAlarmSilenceRulePath, "{project_id}", createAlarmSilenceRuleClient.ProjectID)

	createAlarmSilenceRuleOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		OkCodes: []int{
			204,
		},
	}
	createAlarmSilenceRuleOpt.JSONBody = utils.RemoveNil(buildAlarmSilenceRuleBodyParams(d, cfg))
	_, err = createAlarmSilenceRuleClient.Request("POST", createAlarmSilenceRulePath, &createAlarmSilenceRuleOpt)
	if err != nil {
		return diag.Errorf("error creating AlarmSilenceRule: %s", err)
	}

	d.SetId(d.Get("name").(string))

	return resourceAlarmSilenceRuleRead(ctx, d, meta)
}

func buildAlarmSilenceRuleBodyParams(d *schema.ResourceData, cfg *config.Config) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"name":        utils.ValueIgnoreEmpty(d.Get("name")),
		"desc":        utils.ValueIgnoreEmpty(d.Get("description")),
		"user_id":     cfg.GetProjectID(cfg.GetRegion(d)),
		"timezone":    utils.ValueIgnoreEmpty(d.Get("time_zone")),
		"mute_config": buildAlarmSilenceRuleRequestBodyMuteConfig(d.Get("silence_time")),
		"match":       buildAlarmSilenceRuleSilenceConditions(d.Get("silence_conditions")),
	}
	return bodyParams
}

func buildAlarmSilenceRuleRequestBodyMuteConfig(rawParams interface{}) map[string]interface{} {
	if rawArray, ok := rawParams.([]interface{}); ok {
		if len(rawArray) == 0 {
			return nil
		}
		raw := rawArray[0].(map[string]interface{})
		params := map[string]interface{}{
			"type":      utils.ValueIgnoreEmpty(raw["type"]),
			"starts_at": raw["starts_at"],
			"ends_at":   utils.ValueIgnoreEmpty(raw["ends_at"]),
			"scope":     utils.ValueIgnoreEmpty(raw["scope"]),
		}
		return params
	}
	return nil
}

func buildAlarmSilenceRuleSilenceConditions(rawParams interface{}) [][]map[string]interface{} {
	if rawArray, ok := rawParams.([]interface{}); ok {
		if len(rawArray) == 0 {
			return nil
		}

		rst := make([][]map[string]interface{}, len(rawArray))
		for i, v := range rawArray {
			raw := v.(map[string]interface{})
			rst[i] = buildAlarmSilenceRuleconditions(raw["conditions"])
		}
		return rst
	}
	return nil
}

func buildAlarmSilenceRuleconditions(rawParams interface{}) []map[string]interface{} {
	if rawArray, ok := rawParams.([]interface{}); ok {
		if len(rawArray) == 0 {
			return nil
		}

		rst := make([]map[string]interface{}, len(rawArray))
		for i, v := range rawArray {
			raw := v.(map[string]interface{})
			rst[i] = map[string]interface{}{
				"key":     utils.ValueIgnoreEmpty(raw["key"]),
				"operate": utils.ValueIgnoreEmpty(raw["operate"]),
				"value":   utils.ValueIgnoreEmpty(raw["value"]),
			}
		}
		return rst
	}
	return nil
}

func resourceAlarmSilenceRuleRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var mErr *multierror.Error

	// getAlarmSilenceRule: Query the Alarm Silence Rule
	var (
		getAlarmSilenceRuleHttpUrl = "v2/{project_id}/alert/mute-rules"
		getAlarmSilenceRuleProduct = "aom"
	)
	getAlarmSilenceRuleClient, err := cfg.NewServiceClient(getAlarmSilenceRuleProduct, region)
	if err != nil {
		return diag.Errorf("error creating AOM Client: %s", err)
	}

	getAlarmSilenceRulePath := getAlarmSilenceRuleClient.Endpoint + getAlarmSilenceRuleHttpUrl
	getAlarmSilenceRulePath = strings.ReplaceAll(getAlarmSilenceRulePath, "{project_id}", getAlarmSilenceRuleClient.ProjectID)

	getAlarmSilenceRuleOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		OkCodes: []int{
			200,
		},
	}

	getAlarmSilenceRuleOpt.MoreHeaders = map[string]string{
		"Content-Type": "application/json",
	}

	getAlarmSilenceRuleResp, err := getAlarmSilenceRuleClient.Request("GET", getAlarmSilenceRulePath, &getAlarmSilenceRuleOpt)

	if err != nil {
		return common.CheckDeletedDiag(d, err, "error retrieving AlarmSilenceRule")
	}

	getAlarmSilenceRuleRespBody, err := utils.FlattenResponse(getAlarmSilenceRuleResp)
	if err != nil {
		return diag.FromErr(err)
	}

	rules := FilterListAlarmSilenceRules(getAlarmSilenceRuleRespBody.([]interface{}), d.Id())
	if len(rules) < 1 {
		return common.CheckDeletedDiag(d, golangsdk.ErrDefault404{}, "error retrieving AlarmSilenceRule")
	}

	rule := rules[0]

	mErr = multierror.Append(
		mErr,
		d.Set("region", region),
		d.Set("name", utils.PathSearch("name", rule, nil)),
		d.Set("description", utils.PathSearch("desc", rule, nil)),
		d.Set("time_zone", utils.PathSearch("timezone", rule, nil)),
		d.Set("silence_time", flattenSilenceRuleSilenceTime(rule)),
		d.Set("silence_conditions", flattenSilenceRuleSilenceConditions(rule)),
		d.Set("created_at", utils.PathSearch("create_time", rule, nil)),
		d.Set("updated_at", utils.PathSearch("update_time", rule, nil)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenSilenceRuleSilenceTime(resp interface{}) []interface{} {
	var rst []interface{}
	curJson := utils.PathSearch("mute_config", resp, nil)
	if curJson == nil {
		return rst
	}

	rst = []interface{}{
		map[string]interface{}{
			"type":      utils.PathSearch("type", curJson, nil),
			"starts_at": utils.PathSearch("starts_at", curJson, nil),
			"ends_at":   utils.PathSearch("ends_at", curJson, nil),
			"scope":     utils.PathSearch("scope", curJson, nil),
		},
	}
	return rst
}

func flattenSilenceRuleSilenceConditions(resp interface{}) []interface{} {
	if resp == nil {
		return nil
	}
	curJson := utils.PathSearch("match", resp, make([]interface{}, 0))
	curArray := curJson.([]interface{})
	rst := make([]interface{}, 0, len(curArray))
	for _, v := range curArray {
		rst = append(rst, map[string]interface{}{
			"conditions": flattenSilenceRuleConditions(v),
		})
	}
	return rst
}

func flattenSilenceRuleConditions(resp interface{}) []interface{} {
	if resp == nil {
		return nil
	}
	curArray := resp.([]interface{})
	rst := make([]interface{}, 0, len(curArray))
	for _, v := range curArray {
		rst = append(rst, map[string]interface{}{
			"key":     utils.PathSearch("key", v, nil),
			"operate": utils.PathSearch("operate", v, nil),
			"value":   utils.PathSearch("value", v, nil),
		})
	}
	return rst
}

func FilterListAlarmSilenceRules(all []interface{}, id string) []interface{} {
	rst := make([]interface{}, 0, len(all))
	for _, v := range all {
		if id != fmt.Sprint(utils.PathSearch("name", v, nil)) {
			continue
		}

		rst = append(rst, v)
	}
	return rst
}

func resourceAlarmSilenceRuleUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	updateAlarmSilenceRuleChanges := []string{
		"description",
		"silence_time",
		"silence_conditions",
	}

	if d.HasChanges(updateAlarmSilenceRuleChanges...) {
		// updateAlarmSilenceRule: update the Alarm Silence Rule
		var (
			updateAlarmSilenceRuleHttpUrl = "v2/{project_id}/alert/mute-rules"
			updateAlarmSilenceRuleProduct = "aom"
		)
		updateAlarmSilenceRuleClient, err := cfg.NewServiceClient(updateAlarmSilenceRuleProduct, region)
		if err != nil {
			return diag.Errorf("error creating AOM Client: %s", err)
		}

		updateAlarmSilenceRulePath := updateAlarmSilenceRuleClient.Endpoint + updateAlarmSilenceRuleHttpUrl
		updateAlarmSilenceRulePath = strings.ReplaceAll(updateAlarmSilenceRulePath, "{project_id}", updateAlarmSilenceRuleClient.ProjectID)

		updateAlarmSilenceRuleOpt := golangsdk.RequestOpts{
			KeepResponseBody: true,
			OkCodes: []int{
				204,
			},
		}
		updateAlarmSilenceRuleOpt.JSONBody = utils.RemoveNil(buildAlarmSilenceRuleBodyParams(d, cfg))
		_, err = updateAlarmSilenceRuleClient.Request("PUT", updateAlarmSilenceRulePath, &updateAlarmSilenceRuleOpt)
		if err != nil {
			return diag.Errorf("error updating AlarmSilenceRule: %s", err)
		}
	}
	return resourceAlarmSilenceRuleRead(ctx, d, meta)
}

func resourceAlarmSilenceRuleDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	// deleteAlarmSilenceRule: delete the Alarm Silence Rule
	var (
		deleteAlarmSilenceRuleHttpUrl = "v2/{project_id}/alert/mute-rules"
		deleteAlarmSilenceRuleProduct = "aom"
	)
	deleteAlarmSilenceRuleClient, err := cfg.NewServiceClient(deleteAlarmSilenceRuleProduct, region)
	if err != nil {
		return diag.Errorf("error creating AOM Client: %s", err)
	}

	deleteAlarmSilenceRulePath := deleteAlarmSilenceRuleClient.Endpoint + deleteAlarmSilenceRuleHttpUrl
	deleteAlarmSilenceRulePath = strings.ReplaceAll(deleteAlarmSilenceRulePath, "{project_id}", deleteAlarmSilenceRuleClient.ProjectID)

	deleteAlarmSilenceRuleOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		OkCodes: []int{
			204,
		},
	}
	deleteAlarmSilenceRuleOpt.JSONBody = []map[string]interface{}{
		{
			"name": d.Id(),
		},
	}
	_, err = deleteAlarmSilenceRuleClient.Request("DELETE", deleteAlarmSilenceRulePath, &deleteAlarmSilenceRuleOpt)
	if err != nil {
		return diag.Errorf("error deleting AlarmSilenceRule: %s", err)
	}

	return nil
}
