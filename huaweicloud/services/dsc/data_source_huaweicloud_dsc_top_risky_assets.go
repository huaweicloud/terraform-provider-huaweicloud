package dsc

import (
	"context"
	"strings"

	"github.com/google/uuid"
	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API DSC GET /v2/{project_id}/sec-ops/situation-dashboard/top-risky-assets
func DataSourceDscTopRiskyAssets() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceDscTopRiskyAssetsRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"risk_asset_list": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"asset_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"asset_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"asset_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"data_source": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"deducted_point": {
							Type:     schema.TypeInt,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceDscTopRiskyAssetsRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		product = "dsc"
		httpUrl = "v2/{project_id}/sec-ops/situation-dashboard/top-risky-assets"
	)

	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating DSC client: %s", err)
	}

	requestPath := client.Endpoint + httpUrl
	requestPath = strings.ReplaceAll(requestPath, "{project_id}", client.ProjectID)
	requestOpts := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	resp, err := client.Request("GET", requestPath, &requestOpts)
	if err != nil {
		return diag.Errorf("error retrieving DSC top risky assets: %s", err)
	}

	respBody, err := utils.FlattenResponse(resp)
	if err != nil {
		return diag.FromErr(err)
	}

	id, err := uuid.NewRandom()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(id.String())

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("risk_asset_list", flattenDscTopRiskyAssets(
			utils.PathSearch("risk_asset_list", respBody, make([]interface{}, 0)).([]interface{}))),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenDscTopRiskyAssets(riskAssets []interface{}) []interface{} {
	if len(riskAssets) == 0 {
		return nil
	}

	rst := make([]interface{}, 0, len(riskAssets))
	for _, v := range riskAssets {
		rst = append(rst, map[string]interface{}{
			"asset_id":       utils.PathSearch("asset_id", v, nil),
			"asset_name":     utils.PathSearch("asset_name", v, nil),
			"asset_type":     utils.PathSearch("asset_type", v, nil),
			"data_source":    utils.PathSearch("data_source", v, nil),
			"deducted_point": utils.PathSearch("deducted_point", v, nil),
		})
	}

	return rst
}
