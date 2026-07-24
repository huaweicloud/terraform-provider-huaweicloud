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

// @API DSC GET /v2/{project_id}/data-map/data-infos/{ins_id}
func DataSourceDscDataDetail() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceDscDataDetailRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"ins_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"type": {
				Type:     schema.TypeString,
				Required: true,
			},
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
				Elem:     dscDataDetailScanDetailSchema(),
			},
			"security_strategy": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     dscDataDetailSecurityStrategySchema(),
			},
			"threat_analysis": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     dscDataDetailThreatAnalysisSchema(),
			},
			"vpc_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dscDataDetailScanDetailSchema() *schema.Resource {
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

func dscDataDetailSecurityStrategySchema() *schema.Resource {
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

func dscDataDetailThreatAnalysisSchema() *schema.Resource {
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

func buildDscDataDetailQueryParams(d *schema.ResourceData) string {
	return fmt.Sprintf("?type=%v", d.Get("type"))
}

func dataSourceDscDataDetailRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		product = "dsc"
		httpUrl = "v2/{project_id}/data-map/data-infos/{ins_id}"
		insID   = d.Get("ins_id").(string)
	)

	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating DSC client: %s", err)
	}

	requestPath := client.Endpoint + httpUrl
	requestPath = strings.ReplaceAll(requestPath, "{project_id}", client.ProjectID)
	requestPath = strings.ReplaceAll(requestPath, "{ins_id}", insID)
	requestPath += buildDscDataDetailQueryParams(d)
	requestOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"content-type": "application/json;charset=UTF-8",
		},
	}

	resp, err := client.Request("GET", requestPath, &requestOpt)
	if err != nil {
		return diag.Errorf("error retrieving DSC data detail: %s", err)
	}

	respBody, err := utils.FlattenResponse(resp)
	if err != nil {
		return diag.FromErr(err)
	}

	randomUUID, err := uuid.NewRandom()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(randomUUID.String())

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("asset_name", utils.PathSearch("asset_name", respBody, nil)),
		d.Set("asset_type", utils.PathSearch("asset_type", respBody, nil)),
		d.Set("create_time", utils.PathSearch("create_time", respBody, nil)),
		d.Set("ins_type", utils.PathSearch("ins_type", respBody, nil)),
		d.Set("is_authorized", utils.PathSearch("is_authorized", respBody, nil)),
		d.Set("is_scaned", utils.PathSearch("is_scaned", respBody, nil)),
		d.Set("scan_detail", flattenDscDataDetailScanDetail(utils.PathSearch("scan_detail", respBody, nil))),
		d.Set("security_strategy", flattenDscDataDetailSecurityStrategy(utils.PathSearch("security_strategy", respBody, nil))),
		d.Set("threat_analysis", flattenDscDataDetailThreatAnalysis(utils.PathSearch("threat_analysis", respBody, nil))),
		d.Set("vpc_id", utils.PathSearch("vpc_id", respBody, nil)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenDscDataDetailScanDetail(scanDetail interface{}) []interface{} {
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

func flattenDscDataDetailSecurityStrategy(securityStrategy interface{}) []interface{} {
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

func flattenDscDataDetailThreatAnalysis(threatAnalysis interface{}) []interface{} {
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
