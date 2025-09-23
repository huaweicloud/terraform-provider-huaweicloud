package cce

import (
	"context"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API CCE GET /v2/charts
func DataSourceCCECharts() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceCCEChartsRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"charts": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"values": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"translate": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"instruction": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"version": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"description": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"source": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"icon_url": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"public": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"chart_url": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"created_at": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"updated_at": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceCCEChartsRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var mErr *multierror.Error

	var (
		getChartsHttpUrl = "v2/charts"
		getChartsProduct = "cce"
	)
	getChartsClient, err := cfg.NewServiceClient(getChartsProduct, region)
	if err != nil {
		return diag.Errorf("error creating CCE client: %s", err)
	}

	getChartsHttpPath := getChartsClient.Endpoint + getChartsHttpUrl

	getChartsOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	getChartsResp, err := getChartsClient.Request("GET", getChartsHttpPath, &getChartsOpt)
	if err != nil {
		return diag.Errorf("error retrieving CCE charts: %s", err)
	}

	getChartsRespBody, err := utils.FlattenResponse(getChartsResp)
	if err != nil {
		return diag.FromErr(err)
	}

	dataSourceId, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(dataSourceId)

	mErr = multierror.Append(
		mErr,
		d.Set("charts", flattenCharts(getChartsRespBody)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenCharts(resp interface{}) []map[string]interface{} {
	charts := resp.([]interface{})
	if len(charts) == 0 {
		return nil
	}
	res := make([]map[string]interface{}, len(charts))
	for i, chart := range charts {
		res[i] = map[string]interface{}{
			"id":          utils.PathSearch("id", chart, nil),
			"name":        utils.PathSearch("name", chart, nil),
			"values":      utils.PathSearch("values", chart, nil),
			"translate":   utils.PathSearch("translate", chart, nil),
			"instruction": utils.PathSearch("instruction", chart, nil),
			"version":     utils.PathSearch("version", chart, nil),
			"description": utils.PathSearch("description", chart, nil),
			"source":      utils.PathSearch("source", chart, nil),
			"icon_url":    utils.PathSearch("icon_url", chart, nil),
			"public":      utils.PathSearch("public", chart, nil),
			"chart_url":   utils.PathSearch("chart_url", chart, nil),
			"created_at":  utils.PathSearch("create_at", chart, nil),
			"updated_at":  utils.PathSearch("update_at", chart, nil),
		}
	}
	return res
}
