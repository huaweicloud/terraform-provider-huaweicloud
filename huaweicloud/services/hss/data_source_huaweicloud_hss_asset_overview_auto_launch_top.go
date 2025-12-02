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

// @API HSS GET /v5/{project_id}/asset/overview/statistics/top/auto-launch
func DataSourceAssetOverviewAutoLaunchTop() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceAssetOverviewAutoLaunchTopRead,

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

func buildAssetOverviewAutoLaunchTopQueryParams(epsId string) string {
	if epsId != "" {
		return fmt.Sprintf("?enterprise_project_id=%v", epsId)
	}

	return ""
}

func dataSourceAssetOverviewAutoLaunchTopRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		product = "hss"
		epsId   = cfg.GetEnterpriseProjectID(d)
		httpUrl = "v5/{project_id}/asset/overview/statistics/top/auto-launch"
	)

	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating HSS client: %s", err)
	}

	requestPath := client.Endpoint + httpUrl
	requestPath = strings.ReplaceAll(requestPath, "{project_id}", client.ProjectID)
	requestPath += buildAssetOverviewAutoLaunchTopQueryParams(epsId)
	requestOpts := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	resp, err := client.Request("GET", requestPath, &requestOpts)
	if err != nil {
		return diag.Errorf("error retrieving HSS asset overview auto launch top: %s", err)
	}

	respBody, err := utils.FlattenResponse(resp)
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
		d.Set("total_num", utils.PathSearch("total_num", respBody, nil)),
		d.Set("data_list", flattenAssetOverviewAutoLaunchTopDataList(
			utils.PathSearch("data_list", respBody, make([]interface{}, 0)).([]interface{}))),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenAssetOverviewAutoLaunchTopDataList(resp []interface{}) []interface{} {
	if len(resp) == 0 {
		return nil
	}

	rst := make([]interface{}, 0, len(resp))
	for _, v := range resp {
		rst = append(rst, map[string]interface{}{
			"name":       utils.PathSearch("name", v, nil),
			"host_num":   utils.PathSearch("host_num", v, nil),
			"percentage": utils.PathSearch("percentage", v, nil),
		})
	}

	return rst
}
