package hss

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API HSS GET /v5/{project_id}/overview/security/risk/list
func DataSourceOverviewSecurityRisks() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceOverviewSecurityRisksRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"enterprise_project_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"alarm_risk": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"risk_list": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"severity": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"risk_num": {
										Type:     schema.TypeInt,
										Computed: true,
									},
									"effected_host_num": {
										Type:     schema.TypeInt,
										Computed: true,
									},
								},
							},
						},
						"deduct_score": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"policy_list": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"policy_id": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"policy_name": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"os_type": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"host_num": {
										Type:     schema.TypeInt,
										Computed: true,
									},
									"rule_name": {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						},
						"total_risk_num": {
							Type:     schema.TypeInt,
							Computed: true,
						},
					},
				},
			},
			"baseline_risk": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"risk_list": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"severity": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"risk_num": {
										Type:     schema.TypeInt,
										Computed: true,
									},
									"effected_host_num": {
										Type:     schema.TypeInt,
										Computed: true,
									},
								},
							},
						},
						"deduct_score": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"policy_list": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"policy_id": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"policy_name": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"os_type": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"host_num": {
										Type:     schema.TypeInt,
										Computed: true,
									},
									"rule_name": {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						},
						"existed_pwd_host_num": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"un_scanned_baseline_host_num": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"total_risk_num": {
							Type:     schema.TypeInt,
							Computed: true,
						},
					},
				},
			},
			"asset_risk": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"existed_danger_port_host_num": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"policy_list": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"policy_id": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"policy_name": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"os_type": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"host_num": {
										Type:     schema.TypeInt,
										Computed: true,
									},
									"rule_name": {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						},
						"deduct_score": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"total_risk_num": {
							Type:     schema.TypeInt,
							Computed: true,
						},
					},
				},
			},
			"security_protect_risk": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"un_open_protection_host_num": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"deduct_score": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"total_risk_num": {
							Type:     schema.TypeInt,
							Computed: true,
						},
					},
				},
			},
			"vul_risk": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"risk_list": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"severity": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"risk_num": {
										Type:     schema.TypeInt,
										Computed: true,
									},
									"effected_host_num": {
										Type:     schema.TypeInt,
										Computed: true,
									},
								},
							},
						},
						"deduct_score": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"un_scanned_host_num": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"total_risk_num": {
							Type:     schema.TypeInt,
							Computed: true,
						},
					},
				},
			},
			"image_risk": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"deduct_score": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"un_scanned_image_num": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"risk_list": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"severity": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"image_num": {
										Type:     schema.TypeInt,
										Computed: true,
									},
								},
							},
						},
						"total_risk_num": {
							Type:     schema.TypeInt,
							Computed: true,
						},
					},
				},
			},
			"total_risk_num": {
				Type:     schema.TypeInt,
				Computed: true,
			},
		},
	}
}

func dataSourceOverviewSecurityRisksRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		httpUrl = "v5/{project_id}/overview/security/risk/list"
		epsId   = cfg.GetEnterpriseProjectID(d)
	)

	client, err := cfg.NewServiceClient("hss", region)
	if err != nil {
		return diag.Errorf("error creating HSS client: %s", err)
	}

	getPath := client.Endpoint + httpUrl
	getPath = strings.ReplaceAll(getPath, "{project_id}", client.ProjectID)
	if epsId != "" {
		getPath = fmt.Sprintf("%s?enterprise_project_id=%s", getPath, epsId)
	}

	getOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	getResp, err := client.Request("GET", getPath, &getOpt)
	if err != nil {
		return diag.Errorf("error retrieving the security risks: %s", err)
	}

	respBody, err := utils.FlattenResponse(getResp)
	if err != nil {
		return diag.FromErr(err)
	}

	generateUUID, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}

	d.SetId(generateUUID)

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("alarm_risk", flattenAlarmRiskInfos(utils.PathSearch("alarm_risk", respBody, nil))),
		d.Set("baseline_risk", flattenBaselineRiskInfos(utils.PathSearch("baseline_risk", respBody, nil))),
		d.Set("asset_risk", flattenAssetRiskInfos(utils.PathSearch("asset_risk", respBody, nil))),
		d.Set("security_protect_risk", flattenSecurityProtectRiskInfos(utils.PathSearch("security_protect_risk", respBody, nil))),
		d.Set("vul_risk", flattenVulRiskInfos(utils.PathSearch("vul_risk", respBody, nil))),
		d.Set("image_risk", flattenImageRiskInfos(utils.PathSearch("image_risk", respBody, nil))),
		d.Set("total_risk_num", utils.PathSearch("total_risk_num", respBody, nil)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenAlarmRiskInfos(alarmRisk interface{}) []map[string]interface{} {
	if alarmRisk == nil {
		return nil
	}

	result := map[string]interface{}{
		"risk_list": flattenOverviewRiskListInfos(
			utils.PathSearch("risk_list", alarmRisk, make([]interface{}, 0)).([]interface{})),
		"deduct_score": utils.PathSearch("deduct_score", alarmRisk, nil),
		"policy_list": flattenOverviewPolicyListInfos(
			utils.PathSearch("policy_list", alarmRisk, make([]interface{}, 0)).([]interface{})),
		"total_risk_num": utils.PathSearch("total_risk_num", alarmRisk, nil),
	}

	return []map[string]interface{}{result}
}

func flattenBaselineRiskInfos(baselineRisk interface{}) []map[string]interface{} {
	if baselineRisk == nil {
		return nil
	}

	result := map[string]interface{}{
		"risk_list": flattenOverviewRiskListInfos(
			utils.PathSearch("risk_list", baselineRisk, make([]interface{}, 0)).([]interface{})),
		"deduct_score": utils.PathSearch("deduct_score", baselineRisk, nil),
		"policy_list": flattenOverviewPolicyListInfos(
			utils.PathSearch("policy_list", baselineRisk, make([]interface{}, 0)).([]interface{})),
		"existed_pwd_host_num":         utils.PathSearch("existed_pwd_host_num", baselineRisk, nil),
		"un_scanned_baseline_host_num": utils.PathSearch("un_scanned_baseline_host_num", baselineRisk, nil),
		"total_risk_num":               utils.PathSearch("total_risk_num", baselineRisk, nil),
	}

	return []map[string]interface{}{result}
}

func flattenAssetRiskInfos(assetRisk interface{}) []map[string]interface{} {
	if assetRisk == nil {
		return nil
	}

	result := map[string]interface{}{
		"existed_danger_port_host_num": utils.PathSearch("existed_danger_port_host_num", assetRisk, nil),
		"policy_list": flattenOverviewPolicyListInfos(
			utils.PathSearch("policy_list", assetRisk, make([]interface{}, 0)).([]interface{})),
		"deduct_score":   utils.PathSearch("deduct_score", assetRisk, nil),
		"total_risk_num": utils.PathSearch("total_risk_num", assetRisk, nil),
	}

	return []map[string]interface{}{result}
}

func flattenSecurityProtectRiskInfos(securityProtectRisk interface{}) []map[string]interface{} {
	if securityProtectRisk == nil {
		return nil
	}

	result := map[string]interface{}{
		"un_open_protection_host_num": utils.PathSearch("un_open_protection_host_num", securityProtectRisk, nil),
		"deduct_score":                utils.PathSearch("deduct_score", securityProtectRisk, nil),
		"total_risk_num":              utils.PathSearch("total_risk_num", securityProtectRisk, nil),
	}

	return []map[string]interface{}{result}
}

func flattenVulRiskInfos(vulRisk interface{}) []map[string]interface{} {
	if vulRisk == nil {
		return nil
	}

	result := map[string]interface{}{
		"risk_list": flattenOverviewRiskListInfos(
			utils.PathSearch("risk_list", vulRisk, make([]interface{}, 0)).([]interface{})),
		"deduct_score":        utils.PathSearch("deduct_score", vulRisk, nil),
		"un_scanned_host_num": utils.PathSearch("un_scanned_host_num", vulRisk, nil),
		"total_risk_num":      utils.PathSearch("total_risk_num", vulRisk, nil),
	}

	return []map[string]interface{}{result}
}

func flattenImageRiskInfos(imageRisk interface{}) []map[string]interface{} {
	if imageRisk == nil {
		return nil
	}

	result := map[string]interface{}{
		"deduct_score":         utils.PathSearch("deduct_score", imageRisk, nil),
		"un_scanned_image_num": utils.PathSearch("un_scanned_image_num", imageRisk, nil),
		"risk_list": flattenOverviewImageRiskListInfos(
			utils.PathSearch("risk_list", imageRisk, make([]interface{}, 0)).([]interface{})),
		"total_risk_num": utils.PathSearch("total_risk_num", imageRisk, nil),
	}

	return []map[string]interface{}{result}
}

func flattenOverviewRiskListInfos(riskList []interface{}) []interface{} {
	if len(riskList) == 0 {
		return nil
	}

	rst := make([]interface{}, 0, len(riskList))
	for _, v := range riskList {
		rst = append(rst, map[string]interface{}{
			"severity":          utils.PathSearch("severity", v, nil),
			"risk_num":          utils.PathSearch("risk_num", v, nil),
			"effected_host_num": utils.PathSearch("effected_host_num", v, nil),
		})
	}

	return rst
}

func flattenOverviewPolicyListInfos(policyList []interface{}) []interface{} {
	if len(policyList) == 0 {
		return nil
	}

	result := make([]interface{}, 0, len(policyList))
	for _, v := range policyList {
		result = append(result, map[string]interface{}{
			"policy_id":   utils.PathSearch("policy_id", v, nil),
			"policy_name": utils.PathSearch("policy_name", v, nil),
			"os_type":     utils.PathSearch("os_type", v, nil),
			"host_num":    utils.PathSearch("host_num", v, nil),
			"rule_name":   utils.PathSearch("rule_name", v, nil),
		})
	}

	return result
}

func flattenOverviewImageRiskListInfos(riskList []interface{}) []interface{} {
	if len(riskList) == 0 {
		return nil
	}

	rst := make([]interface{}, 0, len(riskList))
	for _, v := range riskList {
		rst = append(rst, map[string]interface{}{
			"severity":  utils.PathSearch("severity", v, nil),
			"image_num": utils.PathSearch("image_num", v, nil),
		})
	}

	return rst
}
