package hss

import (
	"context"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API HSS GET /v5/{project_id}/overview/hot/information
func DataSourceOverviewHotInformation() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceOverviewHotInformationRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
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
						"title": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"update_time": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"severity_level": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceOverviewHotInformationRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		product = "hss"
		httpUrl = "v5/{project_id}/overview/hot/information"
	)

	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating HSS client: %s", err)
	}

	requestPath := client.Endpoint + httpUrl
	requestPath = strings.ReplaceAll(requestPath, "{project_id}", client.ProjectID)
	requestOpts := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	resp, err := client.Request("GET", requestPath, &requestOpts)
	if err != nil {
		return diag.Errorf("error retrieving HSS overview hot information: %s", err)
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
		d.Set("data_list", flattenOverviewHotInformationDataList(
			utils.PathSearch("data_list", respBody, make([]interface{}, 0)).([]interface{}))),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenOverviewHotInformationDataList(resp []interface{}) []interface{} {
	if len(resp) == 0 {
		return nil
	}

	rst := make([]interface{}, 0, len(resp))
	for _, v := range resp {
		rst = append(rst, map[string]interface{}{
			"title":          utils.PathSearch("title", v, nil),
			"update_time":    utils.PathSearch("update_time", v, nil),
			"severity_level": utils.PathSearch("severity_level", v, nil),
		})
	}

	return rst
}
