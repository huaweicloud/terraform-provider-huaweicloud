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

// @API HSS GET /v5/{project_id}/asset/overview/status/os
func DataSourceAssetOverviewStatusOs() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceAssetOverviewStatusOsRead,

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
			"win_num": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"linux_num": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"os_list": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"os_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"os_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"number": {
							Type:     schema.TypeInt,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func buildAssetOverviewStatusOsQueryParams(epsId string) string {
	if epsId != "" {
		return fmt.Sprintf("?enterprise_project_id=%v", epsId)
	}

	return ""
}

func dataSourceAssetOverviewStatusOsRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		epsId   = cfg.GetEnterpriseProjectID(d)
		product = "hss"
		httpUrl = "v5/{project_id}/asset/overview/status/os"
	)

	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating HSS client: %s", err)
	}

	requestPath := client.Endpoint + httpUrl
	requestPath = strings.ReplaceAll(requestPath, "{project_id}", client.ProjectID)
	requestPath += buildAssetOverviewStatusOsQueryParams(epsId)
	requestOpts := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	resp, err := client.Request("GET", requestPath, &requestOpts)
	if err != nil {
		return diag.Errorf("error retrieving HSS asset overview status OS: %s", err)
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
		d.Set("win_num", utils.PathSearch("win_num", respBody, nil)),
		d.Set("linux_num", utils.PathSearch("linux_num", respBody, nil)),
		d.Set("os_list", flattenOsList(respBody)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenOsList(resp interface{}) []interface{} {
	osList := utils.PathSearch("os_list", resp, make([]interface{}, 0)).([]interface{})
	if len(osList) == 0 {
		return nil
	}

	rst := make([]interface{}, 0, len(osList))
	for _, v := range osList {
		rst = append(rst, map[string]interface{}{
			"os_name": utils.PathSearch("os_name", v, nil),
			"os_type": utils.PathSearch("os_type", v, nil),
			"number":  utils.PathSearch("number", v, nil),
		})
	}

	return rst
}
