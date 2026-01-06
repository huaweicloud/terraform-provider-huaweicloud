package aom

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
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

// @API AOM DELETE /v2/{project_id}/alarm-rules/{alarm_rule_id}
// @API AOM GET /v2/{project_id}/alarm-rules/{alarm_rule_id}
// @API AOM POST /v2/{project_id}/alarm-rules
// @API AOM PUT /v2/{project_id}/alarm-rules
func ResourceAlarmRule() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceAlarmRuleCreate,
		ReadContext:   resourceAlarmRuleRead,
		UpdateContext: resourceAlarmRuleUpdate,
		DeleteContext: resourceAlarmRuleDelete,
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
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"metric_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"namespace": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"dimensions": {
				Type:     schema.TypeList,
				Required: true,
				ForceNew: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:     schema.TypeString,
							Required: true,
							ForceNew: true,
						},
						"value": {
							Type:     schema.TypeString,
							Required: true,
							ForceNew: true,
						},
					},
				},
			},
			"period": {
				Type:     schema.TypeInt,
				Required: true,
				ValidateFunc: validation.IntInSlice([]int{
					60000, 300000, 900000, 3600000,
				}),
			},
			"unit": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"comparison_operator": {
				Type:     schema.TypeString,
				Required: true,
				ValidateFunc: validation.StringInSlice([]string{
					">=", ">", "<=", "<", "=",
				}, false),
			},
			"statistic": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
				ValidateFunc: validation.StringInSlice([]string{
					"maximum", "minimum", "average", "sum", "sampleCount",
				}, false),
			},
			"threshold": {
				Type:     schema.TypeString,
				Required: true,
			},
			"evaluation_periods": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"alarm_level": {
				Type:     schema.TypeInt,
				Optional: true,
				Default:  2,
			},
			"alarm_actions": {
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"alarm_action_enabled": {
				Type:     schema.TypeBool,
				Optional: true,
				ForceNew: true,
				Default:  true,
			},
			"ok_actions": {
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"insufficient_data_actions": {
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},

			"alarm_enabled": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"state_value": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"state_reason": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

// flattenAlarmRuleCreateResponse returns the api response body for creating alarm rule,
// json parse `alarm_rule_id` into float64 in default, which make a loss of precision,
// parse it into json.Number to avoid this situation.
func flattenAlarmRuleCreateResponse(resp *http.Response) (interface{}, error) {
	var respBody interface{}
	defer resp.Body.Close()
	// Don't decode JSON when there is no content
	if resp.StatusCode == http.StatusNoContent {
		_, err := io.Copy(io.Discard, resp.Body)
		return resp, err
	}

	decoder := json.NewDecoder(resp.Body)
	decoder.UseNumber()
	if err := decoder.Decode(&respBody); err != nil {
		return nil, err
	}

	return respBody, nil
}

func resourceAlarmRuleCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	client, err := cfg.NewServiceClient("aom", cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating AOM client: %s", err)
	}

	createHttpUrl := "v2/{project_id}/alarm-rules"
	createPath := client.Endpoint + createHttpUrl
	createPath = strings.ReplaceAll(createPath, "{project_id}", client.ProjectID)
	createOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json",
		},
		JSONBody: utils.RemoveNil(buildV2AlarmRuleBodyParams(d)),
	}

	createResp, err := client.Request("POST", createPath, &createOpt)
	if err != nil {
		return diag.Errorf("error creating the alarm rule: %s", err)
	}

	createRespBody, err := flattenAlarmRuleCreateResponse(createResp)
	if err != nil {
		return diag.Errorf("error flattening the response: %s", err)
	}

	alarmRuleId, ok := utils.PathSearch("alarm_rule_id", createRespBody, json.Number("")).(json.Number)
	if !ok {
		return diag.Errorf("error asserting alarm rule ID value to json.Number")
	}
	if alarmRuleId.String() == "" {
		return diag.Errorf("unable to find alarm rule ID from API response")
	}

	d.SetId(alarmRuleId.String())

	return resourceAlarmRuleRead(ctx, d, meta)
}

func buildV2AlarmRuleBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"alarm_rule_name":           d.Get("name"),
		"alarm_description":         d.Get("description"),
		"alarm_level":               utils.ValueIgnoreEmpty(d.Get("alarm_level")),
		"is_turn_on":                true,
		"action_enabled":            utils.ValueIgnoreEmpty(d.Get("alarm_action_enabled")),
		"alarm_actions":             utils.ValueIgnoreEmpty(d.Get("alarm_actions")),
		"ok_actions":                utils.ValueIgnoreEmpty(d.Get("ok_actions")),
		"insufficient_data_actions": utils.ValueIgnoreEmpty(d.Get("insufficient_data_actions")),
		"namespace":                 d.Get("namespace"),
		"metric_name":               d.Get("metric_name"),
		"dimensions":                buildDimensionsOpts(d.Get("dimensions").([]interface{})),
		"unit":                      d.Get("unit"),
		"threshold":                 d.Get("threshold"),
		"statistic":                 d.Get("statistic"),
		"period":                    d.Get("period"),
		"evaluation_periods":        d.Get("evaluation_periods"),
		"comparison_operator":       d.Get("comparison_operator"),
	}

	return bodyParams
}

func buildDimensionsOpts(rawDimensions []interface{}) interface{} {
	if len(rawDimensions) == 0 {
		return nil
	}
	dimensions := make([]interface{}, len(rawDimensions))
	for i, rawdimension := range rawDimensions {
		dimension := rawdimension.(map[string]interface{})
		dimensions[i] = map[string]interface{}{
			"name":  dimension["name"].(string),
			"value": dimension["value"].(string),
		}
	}
	return dimensions
}

func resourceAlarmRuleRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	client, err := cfg.NewServiceClient("aom", cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating AOM client: %s", err)
	}

	rule, err := GetV2AlarmRule(client, d.Id())
	if err != nil {
		return common.CheckDeletedDiag(d, err, "error getting the alarm rule")
	}

	alarmLevel, _ := strconv.Atoi(utils.PathSearch("alarm_level", rule, "0").(string))

	mErr := multierror.Append(nil,
		d.Set("region", cfg.GetRegion(d)),
		d.Set("name", utils.PathSearch("alarm_rule_name", rule, nil)),
		d.Set("description", utils.PathSearch("alarm_description", rule, nil)),
		d.Set("alarm_level", alarmLevel),
		d.Set("metric_name", utils.PathSearch("metric_name", rule, nil)),
		d.Set("alarm_actions", utils.PathSearch("alarm_actions", rule, nil)),
		d.Set("ok_actions", utils.PathSearch("ok_actions", rule, nil)),
		d.Set("alarm_enabled", utils.PathSearch("id_turn_on", rule, nil)),
		d.Set("alarm_action_enabled", utils.PathSearch("alarm_enabled", rule, nil)),
		d.Set("comparison_operator", utils.PathSearch("comparison_operator", rule, nil)),
		d.Set("evaluation_periods", utils.PathSearch("evaluation_periods", rule, nil)),
		d.Set("insufficient_data_actions", utils.PathSearch("insufficient_data_actions", rule, nil)),
		d.Set("namespace", utils.PathSearch("namespace", rule, nil)),
		d.Set("period", utils.PathSearch("period", rule, nil)),
		d.Set("state_value", utils.PathSearch("state_value", rule, nil)),
		d.Set("state_reason", utils.PathSearch("state_reason", rule, nil)),
		d.Set("statistic", utils.PathSearch("statistic", rule, nil)),
		d.Set("threshold", utils.PathSearch("threshold", rule, nil)),
		d.Set("unit", utils.PathSearch("unit", rule, nil)),
		d.Set("dimensions", flattenV2AlarmRuleDimensions(
			utils.PathSearch("dimensions", rule, make([]interface{}, 0)).([]interface{}))),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func GetV2AlarmRule(client *golangsdk.ServiceClient, ruleId string) (interface{}, error) {
	getHttpUrl := "v2/{project_id}/alarm-rules/{alarm_rule_id}"
	getPath := client.Endpoint + getHttpUrl
	getPath = strings.ReplaceAll(getPath, "{project_id}", client.ProjectID)
	getPath = strings.ReplaceAll(getPath, "{alarm_rule_id}", ruleId)
	getOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json",
		},
	}

	getResp, err := client.Request("GET", getPath, &getOpt)
	if err != nil {
		return nil, err
	}
	getRespBody, err := utils.FlattenResponse(getResp)
	if err != nil {
		return nil, err
	}

	rule := utils.PathSearch("thresholds|[0]", getRespBody, nil)
	if rule == nil {
		return nil, golangsdk.ErrDefault404{
			ErrUnexpectedResponseCode: golangsdk.ErrUnexpectedResponseCode{
				Method:    "GET",
				URL:       "/v2/{project_id}/alarm-rules/{alarm_rule_id}",
				RequestId: "NONE",
				Body:      []byte(fmt.Sprintf("the alarm rule (%s) does not exist", ruleId)),
			},
		}
	}

	return rule, nil
}

func flattenV2AlarmRuleDimensions(paramsList []interface{}) interface{} {
	if len(paramsList) == 0 {
		return nil
	}
	rst := make([]map[string]interface{}, 0, len(paramsList))
	for _, params := range paramsList {
		m := map[string]interface{}{
			"name":  utils.PathSearch("name", params, nil),
			"value": utils.PathSearch("value", params, nil),
		}
		rst = append(rst, m)
	}

	return rst
}

func resourceAlarmRuleUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	client, err := cfg.NewServiceClient("aom", cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating AOM client: %s", err)
	}

	updateHttpUrl := "v2/{project_id}/alarm-rules"
	updatePath := client.Endpoint + updateHttpUrl
	updatePath = strings.ReplaceAll(updatePath, "{project_id}", client.ProjectID)
	updateOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json",
		},
		JSONBody: utils.RemoveNil(buildV2AlarmRuleBodyParams(d)),
	}

	_, err = client.Request("PUT", updatePath, &updateOpt)
	if err != nil {
		return diag.Errorf("error updating the alarm rule: %s", err)
	}

	return resourceAlarmRuleRead(ctx, d, meta)
}

func resourceAlarmRuleDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	client, err := cfg.NewServiceClient("aom", cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating AOM client: %s", err)
	}

	deleteHttpUrl := "v2/{project_id}/alarm-rules/{alarm_rule_id}"
	deletePath := client.Endpoint + deleteHttpUrl
	deletePath = strings.ReplaceAll(deletePath, "{project_id}", client.ProjectID)
	deletePath = strings.ReplaceAll(deletePath, "{alarm_rule_id}", d.Id())

	deleteOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json",
		},
	}

	_, err = client.Request("DELETE", deletePath, &deleteOpt)
	if err != nil {
		return common.CheckDeletedDiag(d, common.ConvertExpected400ErrInto404Err(err, "error_code", "AOM.02005112"),
			"error deleting the alarm rule")
	}

	return nil
}
