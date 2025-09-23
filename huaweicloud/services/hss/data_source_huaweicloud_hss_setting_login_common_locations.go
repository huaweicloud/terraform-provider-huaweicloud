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

// @API HSS GET /v5/{project_id}/setting/login-common-location
func DataSourceSettingLoginCommonLocations() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceSettingLoginCommonLocationsRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"area_code": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"enterprise_project_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"total_num": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"data_list": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"area_code": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"total_num": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"host_id_list": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
					},
				},
			},
		},
	}
}

func buildSettingLoginCommonLocationsQueryParams(d *schema.ResourceData, epsId string) string {
	queryParams := ""

	if v, ok := d.GetOk("area_code"); ok {
		queryParams = fmt.Sprintf("%s&area_code=%v", queryParams, v)
	}

	if epsId != "" {
		queryParams = fmt.Sprintf("%s&enterprise_project_id=%v", queryParams, epsId)
	}

	if queryParams != "" {
		queryParams = "?" + queryParams[1:]
	}

	return queryParams
}

func dataSourceSettingLoginCommonLocationsRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		epsId   = cfg.GetEnterpriseProjectID(d)
		httpUrl = "v5/{project_id}/setting/login-common-location"
	)

	client, err := cfg.NewServiceClient("hss", region)
	if err != nil {
		return diag.Errorf("error creating HSS client: %s", err)
	}

	getPath := client.Endpoint + httpUrl
	getPath = strings.ReplaceAll(getPath, "{project_id}", client.ProjectID)
	getPath += buildSettingLoginCommonLocationsQueryParams(d, epsId)

	getOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	getResp, err := client.Request("GET", getPath, &getOpt)
	if err != nil {
		return diag.Errorf("error retrieving common login locations: %s", err)
	}

	getRespBody, err := utils.FlattenResponse(getResp)
	if err != nil {
		return diag.FromErr(err)
	}

	dataList := utils.PathSearch("data_list", getRespBody, make([]interface{}, 0)).([]interface{})

	generateUUID, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}

	d.SetId(generateUUID)

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("total_num", utils.PathSearch("total_num", getRespBody, nil)),
		d.Set("data_list", flattenSettingLoginCommonLocations(dataList)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenSettingLoginCommonLocations(dataListResp []interface{}) []interface{} {
	rst := make([]interface{}, 0, len(dataListResp))
	for _, v := range dataListResp {
		rst = append(rst, map[string]interface{}{
			"area_code":    utils.PathSearch("area_code", v, nil),
			"total_num":    utils.PathSearch("total_num", v, nil),
			"host_id_list": utils.PathSearch("host_id_list", v, nil),
		})
	}

	return rst
}
