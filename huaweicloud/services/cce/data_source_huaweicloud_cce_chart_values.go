package cce

import (
	"context"
	"encoding/json"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/tidwall/gjson"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/helper/httphelper"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/helper/schemas"
)

func DataSourceCCEShowChartValues() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceCCEShowChartValuesRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Specifies the region in which to query the resource. If omitted, the provider-level region will be used.",
			},
			"chart_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Specifies the chart template id in which to query the resource.",
			},
			"values": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The values of the chart template.",
			},
		},
	}
}

type ShowChartValuesDSWrapper struct {
	*schemas.ResourceDataWrapper
	Config *config.Config
}

func newCCEShowChartValuesDSWrapper(d *schema.ResourceData, meta interface{}) *ShowChartValuesDSWrapper {
	return &ShowChartValuesDSWrapper{
		ResourceDataWrapper: schemas.NewSchemaWrapper(d),
		Config:              meta.(*config.Config),
	}
}

func dataSourceCCEShowChartValuesRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	wrapper := newCCEShowChartValuesDSWrapper(d, meta)
	showChartValuesRst, err := wrapper.ShowChartValues()
	if err != nil {
		return diag.FromErr(err)
	}

	id, err := uuid.GenerateUUID()
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(id)

	err = wrapper.showChartValuesToSchema(showChartValuesRst)
	if err != nil {
		return diag.FromErr(err)
	}

	return nil
}

// @API CCE GET /v2/charts/{chart_id}/values
func (w *ShowChartValuesDSWrapper) ShowChartValues() (*gjson.Result, error) {
	client, err := w.NewClient(w.Config, "cce")
	if err != nil {
		return nil, err
	}

	uri := "/v2/charts/{chart_id}/values"
	uri = strings.ReplaceAll(uri, "{chart_id}", w.Get("chart_id").(string))
	return httphelper.New(client).
		Method("GET").
		URI(uri).
		Request().
		Result()
}

func (w *ShowChartValuesDSWrapper) showChartValuesToSchema(body *gjson.Result) error {
	d := w.ResourceData
	mErr := multierror.Append(nil,
		d.Set("region", w.Config.GetRegion(w.ResourceData)),
		d.Set("values", w.setChartValues(body)),
	)
	return mErr.ErrorOrNil()
}

func (*ShowChartValuesDSWrapper) setChartValues(data *gjson.Result) string {
	chartValues := data.Get("values")
	jsonBytes, _ := json.Marshal(chartValues)

	return string(jsonBytes)
}
