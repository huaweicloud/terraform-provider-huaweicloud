package cdn

import (
	"context"
	"fmt"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/tidwall/gjson"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/helper/httphelper"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/helper/schemas"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API CDN GET /v1.0/cdn/statistics/domain-stats
func DataSourceAnalytics() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceAnalyticsRead,

		Schema: map[string]*schema.Schema{
			// Required parameters
			"action": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The action name.`,
			},
			"start_time": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: `The start timestamp of the query.`,
			},
			"end_time": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: `The end timestamp of the query.`,
			},
			"domain_name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The domain name list.`,
			},
			"stat_type": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The network resource consumption statistics.`,
			},

			// Optional parameters
			"interval": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: `The query interval.`,
			},
			"group_by": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The data grouping mode.`,
			},
			"service_area": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The service area.`,
			},
			"enterprise_project_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The ID of the enterprise project that the resource belongs.`,
			},

			// Attributes
			"result": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The data organized according to the specified grouping mode.`,
			},
		},
	}
}

type AnalyticsDSWrapper struct {
	*schemas.ResourceDataWrapper
	Config *config.Config
}

func newAnalyticsDSWrapper(d *schema.ResourceData, meta interface{}) *AnalyticsDSWrapper {
	return &AnalyticsDSWrapper{
		ResourceDataWrapper: schemas.NewSchemaWrapper(d),
		Config:              meta.(*config.Config),
	}
}

func dataSourceAnalyticsRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	wrapper := newAnalyticsDSWrapper(d, meta)
	showDomainStatsRst, err := wrapper.ShowDomainStats()
	if err != nil {
		return diag.FromErr(err)
	}

	id, err := uuid.GenerateUUID()
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(id)

	err = wrapper.showDomainStatsToSchema(showDomainStatsRst)
	if err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func (w *AnalyticsDSWrapper) ShowDomainStats() (*gjson.Result, error) {
	client, err := w.NewClient(w.Config, "cdn")
	if err != nil {
		return nil, err
	}

	uri := "/v1.0/cdn/statistics/domain-stats"
	params := map[string]any{
		"action":                w.Get("action"),
		"start_time":            w.Get("start_time"),
		"end_time":              w.Get("end_time"),
		"domain_name":           w.Get("domain_name"),
		"stat_type":             w.Get("stat_type"),
		"interval":              w.Get("interval"),
		"group_by":              w.Get("group_by"),
		"service_area":          w.Get("service_area"),
		"enterprise_project_id": w.Get("enterprise_project_id"),
	}
	params = utils.RemoveNil(params)
	return httphelper.New(client).
		Method("GET").
		URI(uri).
		Query(params).
		Request().
		Result()
}

func (w *AnalyticsDSWrapper) showDomainStatsToSchema(body *gjson.Result) error {
	result := body.Get("result").Value()
	resultRaw, err := utils.JsonMarshal(result)
	if err != nil {
		return fmt.Errorf("error marshaling result value: %s", err)
	}

	d := w.ResourceData
	mErr := multierror.Append(nil,
		d.Set("result", string(resultRaw)),
	)
	return mErr.ErrorOrNil()
}
