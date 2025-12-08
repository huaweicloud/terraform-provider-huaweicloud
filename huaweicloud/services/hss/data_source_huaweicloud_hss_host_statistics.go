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

// @API HSS GET /v5/{project_id}/host-management/host-statistics
func DataSourceHostStatistics() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceHostStatisticsRead,

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
			"total_num": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"risk_num": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"unprotected_num": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"not_installed_num": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"installed_failed_num": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"not_online_num": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"version_basic_num": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"version_advanced_num": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"version_enterprise_num": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"version_premium_num": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"version_wtp_num": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"version_container_num": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"host_group_num": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"server_group_num": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"asset_value_list": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"value_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"host_num": {
							Type:     schema.TypeInt,
							Computed: true,
						},
					},
				},
			},
			"server_group_list": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"server_group_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"server_group_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"host_num": {
							Type:     schema.TypeInt,
							Computed: true,
						},
					},
				},
			},
			"ignore_host_num": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"protected_num": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"protect_interrupt_num": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"idle_num": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"premium_non_sp_num": {
				Type:     schema.TypeInt,
				Computed: true,
			},
		},
	}
}

func dataSourceHostStatisticsRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		epsId   = cfg.GetEnterpriseProjectID(d)
		httpUrl = "v5/{project_id}/host-management/host-statistics"
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
		MoreHeaders:      map[string]string{"region": region},
	}

	getResp, err := client.Request("GET", getPath, &getOpt)
	if err != nil {
		return diag.Errorf("error retrieving server statistics: %s", err)
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
		d.Set("total_num", utils.PathSearch("total_num", getRespBody, nil)),
		d.Set("risk_num", utils.PathSearch("risk_num", getRespBody, nil)),
		d.Set("unprotected_num", utils.PathSearch("unprotected_num", getRespBody, nil)),
		d.Set("not_installed_num", utils.PathSearch("not_installed_num", getRespBody, nil)),
		d.Set("installed_failed_num", utils.PathSearch("installed_failed_num", getRespBody, nil)),
		d.Set("not_online_num", utils.PathSearch("not_online_num", getRespBody, nil)),
		d.Set("version_basic_num", utils.PathSearch("version_basic_num", getRespBody, nil)),
		d.Set("version_advanced_num", utils.PathSearch("version_advanced_num", getRespBody, nil)),
		d.Set("version_enterprise_num", utils.PathSearch("version_enterprise_num", getRespBody, nil)),
		d.Set("version_premium_num", utils.PathSearch("version_premium_num", getRespBody, nil)),
		d.Set("version_wtp_num", utils.PathSearch("version_wtp_num", getRespBody, nil)),
		d.Set("version_container_num", utils.PathSearch("version_container_num", getRespBody, nil)),
		d.Set("host_group_num", utils.PathSearch("host_group_num", getRespBody, nil)),
		d.Set("server_group_num", utils.PathSearch("server_group_num", getRespBody, nil)),
		d.Set("asset_value_list", flattenAssetValueHostInfo(
			utils.PathSearch("asset_value_list", getRespBody, make([]interface{}, 0)).([]interface{}))),
		d.Set("server_group_list", flattenServerGroupInfo(
			utils.PathSearch("server_group_list", getRespBody, make([]interface{}, 0)).([]interface{}))),
		d.Set("ignore_host_num", utils.PathSearch("ignore_host_num", getRespBody, nil)),
		d.Set("protected_num", utils.PathSearch("protected_num", getRespBody, nil)),
		d.Set("protect_interrupt_num", utils.PathSearch("protect_interrupt_num", getRespBody, nil)),
		d.Set("idle_num", utils.PathSearch("idle_num", getRespBody, nil)),
		d.Set("premium_non_sp_num", utils.PathSearch("premium_non_sp_num", getRespBody, nil)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenAssetValueHostInfo(assetList []interface{}) []interface{} {
	if len(assetList) == 0 {
		return nil
	}
	rst := make([]interface{}, 0, len(assetList))
	for _, v := range assetList {
		rst = append(rst, map[string]interface{}{
			"value_type": utils.PathSearch("value_type", v, nil),
			"host_num":   utils.PathSearch("host_num", v, nil),
		})
	}

	return rst
}

func flattenServerGroupInfo(groupInfo []interface{}) []interface{} {
	if len(groupInfo) == 0 {
		return nil
	}
	rst := make([]interface{}, 0, len(groupInfo))
	for _, v := range groupInfo {
		rst = append(rst, map[string]interface{}{
			"server_group_id":   utils.PathSearch("server_group_id", v, nil),
			"server_group_name": utils.PathSearch("server_group_name", v, nil),
			"host_num":          utils.PathSearch("host_num", v, nil),
		})
	}

	return rst
}
