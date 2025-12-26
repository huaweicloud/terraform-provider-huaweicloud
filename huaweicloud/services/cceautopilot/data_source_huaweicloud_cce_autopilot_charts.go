package cceautopilot

import (
	"context"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/tidwall/gjson"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/helper/httphelper"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/helper/schemas"
)

func DataSourceCceAutopilotCharts() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceCceAutopilotChartsRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `Specifies the region in which to query the resource. If omitted, the provider-level region will be used.`,
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
						"create_at": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"update_at": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

type ChartsDSWrapper struct {
	*schemas.ResourceDataWrapper
	Config *config.Config
}

func newChartsDSWrapper(d *schema.ResourceData, meta interface{}) *ChartsDSWrapper {
	return &ChartsDSWrapper{
		ResourceDataWrapper: schemas.NewSchemaWrapper(d),
		Config:              meta.(*config.Config),
	}
}

func dataSourceCceAutopilotChartsRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	wrapper := newChartsDSWrapper(d, meta)
	lisAutChaRst, err := wrapper.ListAutopilotCharts()
	if err != nil {
		return diag.FromErr(err)
	}

	id, err := uuid.GenerateUUID()
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(id)

	err = wrapper.listAutopilotChartsToSchema(lisAutChaRst)
	if err != nil {
		return diag.FromErr(err)
	}

	return nil
}

// @API CCE GET /autopilot/v2/charts
func (w *ChartsDSWrapper) ListAutopilotCharts() (*gjson.Result, error) {
	client, err := w.NewClient(w.Config, "cce")
	if err != nil {
		return nil, err
	}

	uri := "/autopilot/v2/charts"
	return httphelper.New(client).
		Method("GET").
		URI(uri).
		Request().
		Result()
}

func (w *ChartsDSWrapper) listAutopilotChartsToSchema(body *gjson.Result) error {
	d := w.ResourceData
	mErr := multierror.Append(nil,
		d.Set("region", w.Config.GetRegion(w.ResourceData)),
		d.Set("charts", schemas.SliceToList(*body,
			func(chart gjson.Result) any {
				return map[string]any{
					"id":          chart.Get("id").Value(),
					"name":        chart.Get("name").Value(),
					"values":      chart.Get("values").Value(),
					"translate":   chart.Get("translate").Value(),
					"instruction": chart.Get("instruction").Value(),
					"version":     chart.Get("version").Value(),
					"description": chart.Get("description").Value(),
					"source":      chart.Get("source").Value(),
					"icon_url":    chart.Get("icon_url").Value(),
					"public":      chart.Get("public").Value(),
					"chart_url":   chart.Get("chart_url").Value(),
					"create_at":   chart.Get("create_at").Value(),
					"update_at":   chart.Get("update_at").Value(),
				}
			},
		)),
	)
	return mErr.ErrorOrNil()
}
