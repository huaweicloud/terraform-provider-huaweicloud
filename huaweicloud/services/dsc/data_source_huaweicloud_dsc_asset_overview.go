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

// @API DSC GET /v2/{project_id}/sec-ops/situation-dashboard/asset-overview
func DataSourceDscAssetOverview() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceDscAssetOverviewRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"asset_sensitive_num": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"asset_total_num": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"high_level_num": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"middle_level_num": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"low_level_num": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"un_classed_num": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"asset_classification_list": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     assetClassificationInfoSchema(),
			},
		},
	}
}

func assetClassificationInfoSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"color_num": {
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
			"sensitive_num": {
				Type:     schema.TypeInt,
				Computed: true,
			},
		},
	}
}

func dataSourceDscAssetOverviewRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		product = "dsc"
		httpUrl = "v2/{project_id}/sec-ops/situation-dashboard/asset-overview"
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
		return diag.Errorf("error retrieving DSC asset overview: %s", err)
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

	assetClassificationList := utils.PathSearch("asset_classification_list", respBody, make([]interface{}, 0)).([]interface{})

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("asset_sensitive_num", utils.PathSearch("asset_sensitive_num", respBody, nil)),
		d.Set("asset_total_num", utils.PathSearch("asset_total_num", respBody, nil)),
		d.Set("high_level_num", utils.PathSearch("high_level_num", respBody, nil)),
		d.Set("middle_level_num", utils.PathSearch("middle_level_num", respBody, nil)),
		d.Set("low_level_num", utils.PathSearch("low_level_num", respBody, nil)),
		d.Set("un_classed_num", utils.PathSearch("un_classed_num", respBody, nil)),
		d.Set("asset_classification_list", flattenAssetClassificationList(assetClassificationList)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenAssetClassificationList(classifications []interface{}) []interface{} {
	if len(classifications) == 0 {
		return nil
	}

	result := make([]interface{}, 0, len(classifications))
	for _, item := range classifications {
		result = append(result, map[string]interface{}{
			"color_num":     utils.PathSearch("color_num", item, nil),
			"level_id":      utils.PathSearch("level_id", item, nil),
			"level_name":    utils.PathSearch("level_name", item, nil),
			"sensitive_num": utils.PathSearch("sensitive_num", item, nil),
		})
	}

	return result
}
