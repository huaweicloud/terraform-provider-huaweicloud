package cdn

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
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API CDN GET /v1.0/cdn/logs
func DataSourceLogs() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceLogsRead,

		Schema: map[string]*schema.Schema{
			// Required parameters
			"domain_name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The domain name to which the logs belong.`,
			},

			// Optional parameters
			"start_time": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: `The start time for querying logs.`,
			},
			"end_time": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: `The end time for querying logs.`,
			},
			"enterprise_project_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The enterprise project ID to which the resource belongs.`,
			},

			// Attributes
			"logs": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"domain_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The domain name to which the log belongs.`,
						},
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The name of the log file.`,
						},
						"size": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: `The size of the log file.`,
						},
						"link": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The log file download link.`,
						},
						"start_time": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: `The start time for querying log.`,
						},
						"end_time": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: `The end time for querying log.`,
						},
					},
				},
				Description: `The list of logs that matched filter parameters.`,
			},
		},
	}
}

type LogsDSWrapper struct {
	*schemas.ResourceDataWrapper
	Config *config.Config
}

func newLogsDSWrapper(d *schema.ResourceData, meta interface{}) *LogsDSWrapper {
	return &LogsDSWrapper{
		ResourceDataWrapper: schemas.NewSchemaWrapper(d),
		Config:              meta.(*config.Config),
	}
}

func (w *LogsDSWrapper) ShowLogs() (*gjson.Result, error) {
	client, err := w.NewClient(w.Config, "cdn")
	if err != nil {
		return nil, err
	}

	uri := "/v1.0/cdn/logs"
	params := map[string]any{
		"domain_name":           w.Get("domain_name"),
		"start_time":            w.Get("start_time"),
		"end_time":              w.Get("end_time"),
		"enterprise_project_id": w.Get("enterprise_project_id"),
	}
	params = utils.RemoveNil(params)
	return httphelper.New(client).
		Method("GET").
		URI(uri).
		Query(params).
		PageSizePager("logs", "page_number", "page_size", 1000).
		Request().
		Result()
}

func (w *LogsDSWrapper) showLogsToSchema(body *gjson.Result) error {
	d := w.ResourceData
	mErr := multierror.Append(nil,
		d.Set("logs", schemas.SliceToList(body.Get("logs"),
			func(logs gjson.Result) any {
				return map[string]any{
					"domain_name": logs.Get("domain_name").Value(),
					"name":        logs.Get("name").Value(),
					"size":        w.setLogsSize(logs),
					"link":        logs.Get("link").Value(),
					"start_time":  logs.Get("start_time").Value(),
					"end_time":    logs.Get("end_time").Value(),
				}
			},
		)),
	)
	return mErr.ErrorOrNil()
}

func (*LogsDSWrapper) setLogsSize(data gjson.Result) float64 {
	return data.Get("size").Float() / 1024
}

func dataSourceLogsRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	wrapper := newLogsDSWrapper(d, meta)
	showLogsRst, err := wrapper.ShowLogs()
	if err != nil {
		return diag.FromErr(err)
	}

	id, err := uuid.GenerateUUID()
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(id)

	err = wrapper.showLogsToSchema(showLogsRst)
	if err != nil {
		return diag.FromErr(err)
	}

	return nil
}
