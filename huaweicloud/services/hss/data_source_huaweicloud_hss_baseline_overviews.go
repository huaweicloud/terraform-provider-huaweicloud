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

// @API HSS GET /v5/{project_id}/baseline/overview
func DataSourceBaselineOverviews() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceBaselineOverviewsRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"group_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"enterprise_project_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"scan_time": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"host_num": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"failed_host_num": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"check_type_num": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"check_rule_num": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"check_rule_pass_rate": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"cn_standard_check_rule_pass_rate": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"hw_standard_check_rule_pass_rate": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"check_rule_failed_num": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"check_rule_high_risk": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"check_rule_medium_risk": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"check_rule_low_risk": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"weak_pwd_total_host": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"weak_pwd_risk": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"weak_pwd_normal": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"weak_pwd_not_protected": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"host_risks": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"host_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"host_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"host_ip": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"scan_time": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"high_risk_num": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"medium_risk_num": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"low_risk_num": {
							Type:     schema.TypeInt,
							Computed: true,
						},
					},
				},
			},
			"weak_pwd_risk_hosts": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"host_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"host_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"host_ip": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"weak_pwd_num": {
							Type:     schema.TypeInt,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func buildBaselineOverviewsQueryParams(d *schema.ResourceData, epsId string) string {
	queryParams := ""

	if v, ok := d.GetOk("group_id"); ok {
		queryParams = fmt.Sprintf("%s&group_id=%v", queryParams, v)
	}

	if epsId != "" {
		queryParams = fmt.Sprintf("%s&enterprise_project_id=%v", queryParams, epsId)
	}

	if queryParams != "" {
		queryParams = "?" + queryParams[1:]
	}

	return queryParams
}

func dataSourceBaselineOverviewsRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		httpUrl = "v5/{project_id}/baseline/overview"
		epsId   = cfg.GetEnterpriseProjectID(d)
	)

	client, err := cfg.NewServiceClient("hss", region)
	if err != nil {
		return diag.Errorf("error creating HSS client: %s", err)
	}

	getPath := client.Endpoint + httpUrl
	getPath = strings.ReplaceAll(getPath, "{project_id}", client.ProjectID)
	getPath += buildBaselineOverviewsQueryParams(d, epsId)
	getOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	getResp, err := client.Request("GET", getPath, &getOpt)
	if err != nil {
		return diag.Errorf("error retrieving baseline overviews: %s", err)
	}

	getRespBody, err := utils.FlattenResponse(getResp)
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
		d.Set("scan_time", utils.PathSearch("scan_time", getRespBody, nil)),
		d.Set("host_num", utils.PathSearch("host_num", getRespBody, nil)),
		d.Set("failed_host_num", utils.PathSearch("failed_host_num", getRespBody, nil)),
		d.Set("check_type_num", utils.PathSearch("check_type_num", getRespBody, nil)),
		d.Set("check_rule_num", utils.PathSearch("check_rule_num", getRespBody, nil)),
		d.Set("check_rule_pass_rate", utils.PathSearch("check_rule_pass_rate", getRespBody, nil)),
		d.Set("cn_standard_check_rule_pass_rate", utils.PathSearch("cn_standard_check_rule_pass_rate", getRespBody, nil)),
		d.Set("hw_standard_check_rule_pass_rate", utils.PathSearch("hw_standard_check_rule_pass_rate", getRespBody, nil)),
		d.Set("check_rule_failed_num", utils.PathSearch("check_rule_failed_num", getRespBody, nil)),
		d.Set("check_rule_high_risk", utils.PathSearch("check_rule_high_risk", getRespBody, nil)),
		d.Set("check_rule_medium_risk", utils.PathSearch("check_rule_medium_risk", getRespBody, nil)),
		d.Set("check_rule_low_risk", utils.PathSearch("check_rule_low_risk", getRespBody, nil)),
		d.Set("weak_pwd_total_host", utils.PathSearch("weak_pwd_total_host", getRespBody, nil)),
		d.Set("weak_pwd_risk", utils.PathSearch("weak_pwd_risk", getRespBody, nil)),
		d.Set("weak_pwd_normal", utils.PathSearch("weak_pwd_normal", getRespBody, nil)),
		d.Set("weak_pwd_not_protected", utils.PathSearch("weak_pwd_not_protected", getRespBody, nil)),
		d.Set("host_risks", flattenBaselineHostRiskNumInfo(
			utils.PathSearch("host_risks", getRespBody, make([]interface{}, 0)).([]interface{}))),
		d.Set("weak_pwd_risk_hosts", flattenBaselineHostWeakPwdRiskNumInfo(
			utils.PathSearch("weak_pwd_risk_hosts", getRespBody, make([]interface{}, 0)).([]interface{}))),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenBaselineHostRiskNumInfo(dataResp []interface{}) []interface{} {
	if len(dataResp) == 0 {
		return nil
	}

	result := make([]interface{}, 0, len(dataResp))
	for _, v := range dataResp {
		result = append(result, map[string]interface{}{
			"host_id":         utils.PathSearch("host_id", v, nil),
			"host_name":       utils.PathSearch("host_name", v, nil),
			"host_ip":         utils.PathSearch("host_ip", v, nil),
			"scan_time":       utils.PathSearch("scan_time", v, nil),
			"high_risk_num":   utils.PathSearch("high_risk_num", v, nil),
			"medium_risk_num": utils.PathSearch("medium_risk_num", v, nil),
			"low_risk_num":    utils.PathSearch("low_risk_num", v, nil),
		})
	}

	return result
}

func flattenBaselineHostWeakPwdRiskNumInfo(dataResp []interface{}) []interface{} {
	if len(dataResp) == 0 || dataResp[0] == nil {
		return nil
	}

	result := make([]interface{}, 0, len(dataResp))
	for _, v := range dataResp {
		result = append(result, map[string]interface{}{
			"host_id":      utils.PathSearch("host_id", v, nil),
			"host_name":    utils.PathSearch("host_name", v, nil),
			"host_ip":      utils.PathSearch("host_ip", v, nil),
			"weak_pwd_num": utils.PathSearch("weak_pwd_num", v, nil),
		})
	}

	return result
}
