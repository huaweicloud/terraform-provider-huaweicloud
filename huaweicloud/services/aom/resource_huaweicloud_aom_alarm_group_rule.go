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

// @API AOM POST /v2/{project_id}/alert/group-rules
// @API AOM PUT /v2/{project_id}/alert/group-rules
// @API AOM DELETE /v2/{project_id}/alert/group-rules
// @API AOM GET /v2/{project_id}/alert/group-rules
func ResourceAlarmGroupRule() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceAlarmGroupRuleCreate,
		ReadContext:   resourceAlarmGroupRuleRead,
		UpdateContext: resourceAlarmGroupRuleUpdate,
		DeleteContext: resourceAlarmGroupRuleDelete,

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
			"detail": {
				Type:     schema.TypeSet,
				Required: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"bind_notification_rule_ids": {
							Type:     schema.TypeSet,
							Required: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
						"match": {
							Type:     schema.TypeSet,
							Required: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"key": {
										Type:     schema.TypeString,
										Required: true,
									},
									"operate": {
										Type:     schema.TypeString,
										Required: true,
									},
									"value": {
										Type:     schema.TypeSet,
										Optional: true,
										Elem:     &schema.Schema{Type: schema.TypeString},
									},
								},
							},
						},
					},
				},
			},
			"group_by": {
				Type:     schema.TypeSet,
				Required: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"group_wait": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"group_interval": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"group_repeat_waiting": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"enterprise_project_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"created_at": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"updated_at": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceAlarmGroupRuleCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	client, err := cfg.NewServiceClient("aom", region)
	if err != nil {
		return diag.Errorf("error creating AOM client: %s", err)
	}

	createHttpUrl := "v2/{project_id}/alert/group-rules"
	createPath := client.Endpoint + createHttpUrl
	createPath = strings.ReplaceAll(createPath, "{project_id}", client.ProjectID)
	createOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		OkCodes:          []int{204},
		MoreHeaders:      buildRequestMoreHeaders(cfg.GetEnterpriseProjectID(d)),
		JSONBody:         utils.RemoveNil(buildCreateAlarmGroupRuleBodyParams(d)),
	}

	_, err = client.Request("POST", createPath, &createOpt)
	if err != nil {
		return diag.Errorf("error creating AOM alarm group rule: %s", err)
	}

	d.SetId(d.Get("name").(string))

	return resourceAlarmGroupRuleRead(ctx, d, meta)
}

func buildCreateAlarmGroupRuleBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"name":                 d.Get("name"),
		"detail":               buildAlarmGroupRuleDetail(d),
		"group_by":             d.Get("group_by").(*schema.Set).List(),
		"group_wait":           d.Get("group_wait"),
		"group_interval":       d.Get("group_interval"),
		"group_repeat_waiting": d.Get("group_repeat_waiting"),
		"desc":                 d.Get("description"),
	}

	return bodyParams
}

func buildAlarmGroupRuleDetail(d *schema.ResourceData) []interface{} {
	detail := d.Get("detail").(*schema.Set).List()
	rst := make([]interface{}, 0, len(detail))
	for _, v := range detail {
		params := v.(map[string]interface{})
		rst = append(rst, map[string]interface{}{
			"match":                      buildAlarmGroupRuleDetailMatch(params["match"].(*schema.Set).List()),
			"bind_notification_rule_ids": params["bind_notification_rule_ids"].(*schema.Set).List(),
		})
	}
	return rst
}

func buildAlarmGroupRuleDetailMatch(paramsList []interface{}) []interface{} {
	rst := make([]interface{}, 0, len(paramsList))
	for _, params := range paramsList {
		raw := params.(map[string]interface{})
		m := map[string]interface{}{
			"key":     raw["key"],
			"operate": raw["operate"],
		}
		value := []interface{}{}
		if v, ok := raw["value"].(*schema.Set); ok && v.Len() > 0 {
			value = v.List()
		}
		m["value"] = value

		rst = append(rst, m)
	}

	return rst
}

func resourceAlarmGroupRuleRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	client, err := cfg.NewServiceClient("aom", region)
	if err != nil {
		return diag.Errorf("error creating AOM client: %s", err)
	}

	rule, err := GetAlarmGroupRule(client, d.Id())
	if err != nil {
		return common.CheckDeletedDiag(d, err, "error retrieving alarm group rule")
	}

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("name", utils.PathSearch("name", rule, nil)),
		d.Set("detail", flattenAlarmGroupRuleDetail(utils.PathSearch("detail", rule, make([]interface{}, 0)).([]interface{}))),
		d.Set("group_by", utils.PathSearch("group_by", rule, nil)),
		d.Set("group_wait", utils.PathSearch("group_wait", rule, nil)),
		d.Set("group_interval", utils.PathSearch("group_interval", rule, nil)),
		d.Set("group_repeat_waiting", utils.PathSearch("group_repeat_waiting", rule, nil)),
		d.Set("description", utils.PathSearch("desc", rule, nil)),
		d.Set("enterprise_project_id", utils.PathSearch("enterprise_project_id", rule, nil)),
		d.Set("created_at", utils.FormatTimeStampRFC3339(
			int64(utils.PathSearch("create_time", rule, float64(0)).(float64))/1000, true)),
		d.Set("updated_at", utils.FormatTimeStampRFC3339(
			int64(utils.PathSearch("update_time", rule, float64(0)).(float64))/1000, true)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func GetAlarmGroupRule(client *golangsdk.ServiceClient, name string) (interface{}, error) {
	listHttpUrl := "v2/{project_id}/alert/group-rules"
	listPath := client.Endpoint + listHttpUrl
	listPath = strings.ReplaceAll(listPath, "{project_id}", client.ProjectID)
	listOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type":          "application/json",
			"Enterprise-Project-Id": "all_granted_eps",
		},
	}

	listResp, err := client.Request("GET", listPath, &listOpt)
	if err != nil {
		return nil, err
	}
	listRespBody, err := utils.FlattenResponse(listResp)
	if err != nil {
		return nil, err
	}

	jsonPath := fmt.Sprintf("[?name=='%s']|[0]", name)
	rule := utils.PathSearch(jsonPath, listRespBody, nil)
	if rule == nil {
		return nil, golangsdk.ErrDefault404{
			ErrUnexpectedResponseCode: golangsdk.ErrUnexpectedResponseCode{
				Method:    "GET",
				URL:       "/v2/{project_id}/alert/group-rules",
				RequestId: "NONE",
				Body:      []byte(fmt.Sprintf("the alarm group rule (%s) does not exist", name)),
			},
		}
	}

	return rule, nil
}

func flattenAlarmGroupRuleDetail(paramsList []interface{}) []interface{} {
	if len(paramsList) == 0 {
		return nil
	}
	rst := make([]interface{}, 0, len(paramsList))
	for _, params := range paramsList {
		m := map[string]interface{}{
			"bind_notification_rule_ids": utils.PathSearch("bind_notification_rule_ids", params, nil),
			"match": flattenAlarmGroupRuleDetailMatch(
				utils.PathSearch("match", params, make([]interface{}, 0)).([]interface{})),
		}
		rst = append(rst, m)
	}

	return rst
}

func flattenAlarmGroupRuleDetailMatch(paramsList []interface{}) []interface{} {
	if len(paramsList) == 0 {
		return nil
	}
	rst := make([]interface{}, 0, len(paramsList))
	for _, params := range paramsList {
		m := map[string]interface{}{
			"key":     utils.PathSearch("key", params, nil),
			"operate": utils.PathSearch("operate", params, nil),
			"value":   utils.PathSearch("value", params, nil),
		}
		rst = append(rst, m)
	}

	return rst
}

func resourceAlarmGroupRuleUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	client, err := cfg.NewServiceClient("aom", region)
	if err != nil {
		return diag.Errorf("error creating AOM client: %s", err)
	}

	updateHttpUrl := "v2/{project_id}/alert/group-rules"
	updatePath := client.Endpoint + updateHttpUrl
	updatePath = strings.ReplaceAll(updatePath, "{project_id}", client.ProjectID)
	updateOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		OkCodes:          []int{204},
		MoreHeaders:      buildRequestMoreHeaders(cfg.GetEnterpriseProjectID(d)),
		JSONBody:         buildCreateAlarmGroupRuleBodyParams(d),
	}

	_, err = client.Request("PUT", updatePath, &updateOpt)
	if err != nil {
		return diag.Errorf("error updating alarm group rule: %s", err)
	}

	return resourceAlarmGroupRuleRead(ctx, d, meta)
}

func resourceAlarmGroupRuleDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	client, err := cfg.NewServiceClient("aom", region)
	if err != nil {
		return diag.Errorf("error creating AOM client: %s", err)
	}

	deleteHttpUrl := "v2/{project_id}/alert/group-rules"
	deletePath := client.Endpoint + deleteHttpUrl
	deletePath = strings.ReplaceAll(deletePath, "{project_id}", client.ProjectID)
	deleteOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      buildRequestMoreHeaders(cfg.GetEnterpriseProjectID(d)),
		JSONBody:         []interface{}{d.Id()},
	}

	_, err = client.Request("DELETE", deletePath, &deleteOpt)
	if err != nil {
		return common.CheckDeletedDiag(d, common.ConvertExpected400ErrInto404Err(err, "error_code", "AOM.08002002"),
			"error deleting alarm group rule")
	}

	return nil
}
