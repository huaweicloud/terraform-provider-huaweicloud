package dsc

import (
	"context"
	"fmt"
	"strings"

	"github.com/google/uuid"
	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API DSC GET /v2/{project_id}/data-map/risk-level-statistics
func DataSourceDataMapRiskLevelInfo() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceDataMapRiskLevelInfoRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"security_group_ids": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"security_level_ids": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"type": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"data_level": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     dataMapRiskLevelDataLevelSchema(),
			},
		},
	}
}

func dataMapRiskLevelDataLevelSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"level_color_number": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"level_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"level_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"risk_list": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     dataMapRiskLevelRiskListSchema(),
			},
			"total": {
				Type:     schema.TypeInt,
				Computed: true,
			},
		},
	}
}

func dataMapRiskLevelRiskListSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"detail_list": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     dataMapRiskLevelCommonInfoSchema(),
			},
			"total": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"type": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataMapRiskLevelCommonInfoSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"asset_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"asset_type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"create_time": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"ins_type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"is_authorized": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"is_scaned": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"scan_detail": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     dataMapRiskLevelScanDetailSchema(),
			},
			"security_strategy": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     dataMapRiskLevelSecurityStrategySchema(),
			},
			"threat_analysis": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     dataMapRiskLevelThreatAnalysisSchema(),
			},
			"vpc_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataMapRiskLevelScanDetailSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"job_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"last_scan_time": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"object_num": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"scan_risk": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"scan_template_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"scan_template_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"security_level_color": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"security_level_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"security_level_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"sensitive_object_num": {
				Type:     schema.TypeInt,
				Computed: true,
			},
		},
	}
}

func dataMapRiskLevelSecurityStrategySchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"ssl_enabled": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"access_strategy": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"access_strategy_level": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"authority_enable": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"authority_level": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"backup_and_restore": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"backup_enable": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"backup_level": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"data_volume_encrypt_enable": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"data_volume_encrypt_level": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"dbss_audit_security_level": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"dbss_audit_status": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"disk_encrypted": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"disk_encrypted_enable": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"encrypt_level": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"https_enable": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"https_level": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"ik_enable": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"is_encrypt": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"obs_audit_level": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"obs_audit_status": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"public_network_access": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"public_network_enable": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"security_group_level": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"ssl_status": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"total_risk": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"working_mode": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"working_type": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataMapRiskLevelThreatAnalysisSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"abnormal_login_level": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"risky_operation_level": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"sql_inject_level": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func buildDataMapRiskLevelInfoQueryParams(d *schema.ResourceData, limit, offset int) string {
	queryParams := fmt.Sprintf("?limit=%d&offset=%d", limit, offset)

	for _, v := range utils.ExpandToStringList(d.Get("security_group_ids").([]interface{})) {
		queryParams = fmt.Sprintf("%s&security_group_ids=%v", queryParams, v)
	}

	for _, v := range utils.ExpandToStringList(d.Get("security_level_ids").([]interface{})) {
		queryParams = fmt.Sprintf("%s&security_level_ids=%v", queryParams, v)
	}

	if v, ok := d.GetOk("type"); ok {
		queryParams = fmt.Sprintf("%s&type=%v", queryParams, v)
	}

	return queryParams
}

func dataSourceDataMapRiskLevelInfoRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		product = "dsc"
		httpUrl = "v2/{project_id}/data-map/risk-level-statistics"
		offset  = 0
		limit   = 1000
		result  = make([]interface{}, 0)
	)

	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating DSC client: %s", err)
	}

	requestPath := client.Endpoint + httpUrl
	requestPath = strings.ReplaceAll(requestPath, "{project_id}", client.ProjectID)
	requestOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"content-type": "application/json;charset=UTF-8",
		},
	}

	for {
		currentPath := requestPath + buildDataMapRiskLevelInfoQueryParams(d, limit, offset)

		resp, err := client.Request("GET", currentPath, &requestOpt)
		if err != nil {
			return diag.Errorf("error retrieving DSC data map risk level info: %s", err)
		}

		respBody, err := utils.FlattenResponse(resp)
		if err != nil {
			return diag.FromErr(err)
		}

		dataLevelList := utils.PathSearch("data_level", respBody, make([]interface{}, 0)).([]interface{})
		result = append(result, dataLevelList...)

		if len(dataLevelList) < limit {
			break
		}

		offset += len(dataLevelList)
	}

	randomUUID, err := uuid.NewRandom()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(randomUUID.String())

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("data_level", flattenDataMapRiskLevelDataLevelList(result)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenDataMapRiskLevelDataLevelList(dataLevelList []interface{}) []interface{} {
	if len(dataLevelList) == 0 {
		return nil
	}

	rst := make([]interface{}, 0, len(dataLevelList))
	for _, v := range dataLevelList {
		rst = append(rst, map[string]interface{}{
			"level_color_number": utils.PathSearch("level_color_number", v, nil),
			"level_id":           utils.PathSearch("level_id", v, nil),
			"level_name":         utils.PathSearch("level_name", v, nil),
			"risk_list": flattenDataMapRiskLevelRiskList(
				utils.PathSearch("risk_list", v, make([]interface{}, 0)).([]interface{})),
			"total": utils.PathSearch("total", v, nil),
		})
	}

	return rst
}

func flattenDataMapRiskLevelRiskList(riskList []interface{}) []interface{} {
	if len(riskList) == 0 {
		return nil
	}

	rst := make([]interface{}, 0, len(riskList))
	for _, v := range riskList {
		rst = append(rst, map[string]interface{}{
			"detail_list": flattenDataMapRiskLevelCommonInfoList(
				utils.PathSearch("detail_list", v, make([]interface{}, 0)).([]interface{})),
			"total": utils.PathSearch("total", v, nil),
			"type":  utils.PathSearch("type", v, nil),
		})
	}

	return rst
}

func flattenDataMapRiskLevelCommonInfoList(detailList []interface{}) []interface{} {
	if len(detailList) == 0 {
		return nil
	}

	rst := make([]interface{}, 0, len(detailList))
	for _, v := range detailList {
		rst = append(rst, map[string]interface{}{
			"asset_name":        utils.PathSearch("asset_name", v, nil),
			"asset_type":        utils.PathSearch("asset_type", v, nil),
			"create_time":       utils.PathSearch("create_time", v, nil),
			"id":                utils.PathSearch("id", v, nil),
			"ins_type":          utils.PathSearch("ins_type", v, nil),
			"is_authorized":     utils.PathSearch("is_authorized", v, nil),
			"is_scaned":         utils.PathSearch("is_scaned", v, nil),
			"scan_detail":       flattenDataMapRiskLevelScanDetail(utils.PathSearch("scan_detail", v, nil)),
			"security_strategy": flattenDataMapRiskLevelSecurityStrategy(utils.PathSearch("security_strategy", v, nil)),
			"threat_analysis":   flattenDataMapRiskLevelThreatAnalysis(utils.PathSearch("threat_analysis", v, nil)),
			"vpc_id":            utils.PathSearch("vpc_id", v, nil),
		})
	}

	return rst
}

func flattenDataMapRiskLevelScanDetail(scanDetail interface{}) []interface{} {
	if scanDetail == nil {
		return nil
	}

	return []interface{}{
		map[string]interface{}{
			"job_id":               utils.PathSearch("job_id", scanDetail, nil),
			"last_scan_time":       utils.PathSearch("last_scan_time", scanDetail, nil),
			"object_num":           utils.PathSearch("object_num", scanDetail, nil),
			"scan_risk":            utils.PathSearch("scan_risk", scanDetail, nil),
			"scan_template_id":     utils.PathSearch("scan_template_id", scanDetail, nil),
			"scan_template_name":   utils.PathSearch("scan_template_name", scanDetail, nil),
			"security_level_color": utils.PathSearch("security_level_color", scanDetail, nil),
			"security_level_id":    utils.PathSearch("security_level_id", scanDetail, nil),
			"security_level_name":  utils.PathSearch("security_level_name", scanDetail, nil),
			"sensitive_object_num": utils.PathSearch("sensitive_object_num", scanDetail, nil),
		},
	}
}

func flattenDataMapRiskLevelSecurityStrategy(securityStrategy interface{}) []interface{} {
	if securityStrategy == nil {
		return nil
	}

	return []interface{}{
		map[string]interface{}{
			"ssl_enabled":                utils.PathSearch("ssl_enabled", securityStrategy, nil),
			"access_strategy":            utils.PathSearch("access_strategy", securityStrategy, nil),
			"access_strategy_level":      utils.PathSearch("access_strategy_level", securityStrategy, nil),
			"authority_enable":           utils.PathSearch("authority_enable", securityStrategy, nil),
			"authority_level":            utils.PathSearch("authority_level", securityStrategy, nil),
			"backup_and_restore":         utils.PathSearch("backup_and_restore", securityStrategy, nil),
			"backup_enable":              utils.PathSearch("backup_enable", securityStrategy, nil),
			"backup_level":               utils.PathSearch("backup_level", securityStrategy, nil),
			"data_volume_encrypt_enable": utils.PathSearch("data_volume_encrypt_enable", securityStrategy, nil),
			"data_volume_encrypt_level":  utils.PathSearch("data_volume_encrypt_level", securityStrategy, nil),
			"dbss_audit_security_level":  utils.PathSearch("dbss_audit_security_level", securityStrategy, nil),
			"dbss_audit_status":          utils.PathSearch("dbss_audit_status", securityStrategy, nil),
			"disk_encrypted":             utils.PathSearch("disk_encrypted", securityStrategy, nil),
			"disk_encrypted_enable":      utils.PathSearch("disk_encrypted_enable", securityStrategy, nil),
			"encrypt_level":              utils.PathSearch("encrypt_level", securityStrategy, nil),
			"https_enable":               utils.PathSearch("https_enable", securityStrategy, nil),
			"https_level":                utils.PathSearch("https_level", securityStrategy, nil),
			"ik_enable":                  utils.PathSearch("ik_enable", securityStrategy, nil),
			"is_encrypt":                 utils.PathSearch("is_encrypt", securityStrategy, nil),
			"obs_audit_level":            utils.PathSearch("obs_audit_level", securityStrategy, nil),
			"obs_audit_status":           utils.PathSearch("obs_audit_status", securityStrategy, nil),
			"public_network_access":      utils.PathSearch("public_network_access", securityStrategy, nil),
			"public_network_enable":      utils.PathSearch("public_network_enable", securityStrategy, nil),
			"security_group_level":       utils.PathSearch("security_group_level", securityStrategy, nil),
			"ssl_status":                 utils.PathSearch("ssl_status", securityStrategy, nil),
			"total_risk":                 utils.PathSearch("total_risk", securityStrategy, nil),
			"working_mode":               utils.PathSearch("working_mode", securityStrategy, nil),
			"working_type":               utils.PathSearch("working_type", securityStrategy, nil),
		},
	}
}

func flattenDataMapRiskLevelThreatAnalysis(threatAnalysis interface{}) []interface{} {
	if threatAnalysis == nil {
		return nil
	}

	return []interface{}{
		map[string]interface{}{
			"abnormal_login_level":  utils.PathSearch("abnormal_login_level", threatAnalysis, nil),
			"risky_operation_level": utils.PathSearch("risky_operation_level", threatAnalysis, nil),
			"sql_inject_level":      utils.PathSearch("sql_inject_level", threatAnalysis, nil),
		},
	}
}
