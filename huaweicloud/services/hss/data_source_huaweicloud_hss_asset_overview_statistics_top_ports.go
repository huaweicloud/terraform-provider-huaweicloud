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

// @API HSS GET /v5/{project_id}/asset/overview/statistics/top/port
func DataSourceAssetOverviewStatisticsTopPorts() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceAssetOverviewStatisticsTopPortsRead,

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
			"data_list": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"host_num": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"percentage": {
							Type:     schema.TypeInt,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceAssetOverviewStatisticsTopPortsRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		epsId   = cfg.GetEnterpriseProjectID(d)
		httpUrl = "v5/{project_id}/asset/overview/statistics/top/port"
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
		return diag.Errorf("error retrieving Top ports information: %s", err)
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
		d.Set("data_list", flattenAssetOverviewStatisticsTopPorts(dataList)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenAssetOverviewStatisticsTopPorts(dataList []interface{}) []interface{} {
	if len(dataList) == 0 {
		return nil
	}

	rst := make([]interface{}, 0, len(dataList))
	for _, v := range dataList {
		rst = append(rst, map[string]interface{}{
			"name":       utils.PathSearch("name", v, nil),
			"host_num":   utils.PathSearch("host_num", v, nil),
			"percentage": utils.PathSearch("percentage", v, nil),
		})
	}

	return rst
}
