// Generated by PMS #346
package secmaster

import (
	"context"
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

func DataSourceSecmasterPlaybookStatistics() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceSecmasterPlaybookStatisticsRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `Specifies the region in which to query the resource. If omitted, the provider-level region will be used.`,
			},
			"workspace_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `Specifies the workspace ID.`,
			},
			"data": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: `The playbook statistics.`,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"enabled_num": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: `The number of enabled playbooks.`,
						},
						"unapproved_num": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: `The number of unapproved playbooks.`,
						},
						"disabled_num": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: `The number of playbooks that are not enabled.`,
						},
					},
				},
			},
		},
	}
}

type PlaybookStatisticsDSWrapper struct {
	*schemas.ResourceDataWrapper
	Config *config.Config
}

func newPlaybookStatisticsDSWrapper(d *schema.ResourceData, meta interface{}) *PlaybookStatisticsDSWrapper {
	return &PlaybookStatisticsDSWrapper{
		ResourceDataWrapper: schemas.NewSchemaWrapper(d),
		Config:              meta.(*config.Config),
	}
}

func dataSourceSecmasterPlaybookStatisticsRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	wrapper := newPlaybookStatisticsDSWrapper(d, meta)
	shoPlaStaRst, err := wrapper.ShowPlaybookStatistics()
	if err != nil {
		return diag.FromErr(err)
	}

	id, err := uuid.GenerateUUID()
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(id)

	err = wrapper.showPlaybookStatisticsToSchema(shoPlaStaRst)
	if err != nil {
		return diag.FromErr(err)
	}

	return nil
}

// @API SecMaster GET /v1/{project_id}/workspaces/{workspace_id}/soc/playbooks/statistics
func (w *PlaybookStatisticsDSWrapper) ShowPlaybookStatistics() (*gjson.Result, error) {
	client, err := w.NewClient(w.Config, "secmaster")
	if err != nil {
		return nil, err
	}

	uri := "/v1/{project_id}/workspaces/{workspace_id}/soc/playbooks/statistics"
	uri = strings.ReplaceAll(uri, "{workspace_id}", w.Get("workspace_id").(string))
	return httphelper.New(client).
		Method("GET").
		URI(uri).
		Request().
		Result()
}

func (w *PlaybookStatisticsDSWrapper) showPlaybookStatisticsToSchema(body *gjson.Result) error {
	d := w.ResourceData
	mErr := multierror.Append(nil,
		d.Set("region", w.Config.GetRegion(w.ResourceData)),
		d.Set("data", schemas.ObjectToList(body.Get("data"),
			func(data gjson.Result) any {
				return map[string]any{
					"enabled_num":    data.Get("enabled_num").Value(),
					"unapproved_num": data.Get("unapproved_num").Value(),
					"disabled_num":   data.Get("disabled_num").Value(),
				}
			},
		)),
	)
	return mErr.ErrorOrNil()
}
