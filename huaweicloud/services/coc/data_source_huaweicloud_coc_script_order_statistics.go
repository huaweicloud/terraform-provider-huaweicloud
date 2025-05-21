package coc

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

func DataSourceCocScriptOrderStatistics() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceCocScriptOrderStatisticsRead,

		Schema: map[string]*schema.Schema{
			"execute_uuid": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `Specifies the execution ID of a script order.`,
			},
			"execute_statistics": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: `Indicates the statistical details.`,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"instance_status": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `Indicates the status of the execution instance.`,
						},
						"instance_count": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: `Indicates the number of instances executed in this state.`,
						},
						"batch_indexes": {
							Type:        schema.TypeList,
							Computed:    true,
							Elem:        &schema.Schema{Type: schema.TypeInt},
							Description: `Indicates a list of batch indexes in this state.`,
						},
					},
				},
			},
		},
	}
}

type ScriptOrderStatisticsDSWrapper struct {
	*schemas.ResourceDataWrapper
	Config *config.Config
}

func newScriptOrderStatisticsDSWrapper(d *schema.ResourceData, meta interface{}) *ScriptOrderStatisticsDSWrapper {
	return &ScriptOrderStatisticsDSWrapper{
		ResourceDataWrapper: schemas.NewSchemaWrapper(d),
		Config:              meta.(*config.Config),
	}
}

func dataSourceCocScriptOrderStatisticsRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	wrapper := newScriptOrderStatisticsDSWrapper(d, meta)
	getScrJobStaRst, err := wrapper.GetScriptJobStatistics()
	if err != nil {
		return diag.FromErr(err)
	}

	id, err := uuid.GenerateUUID()
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(id)

	err = wrapper.getScriptJobStatisticsToSchema(getScrJobStaRst)
	if err != nil {
		return diag.FromErr(err)
	}

	return nil
}

// @API COC GET /v1/job/script/orders/{execute_uuid}/statistics
func (w *ScriptOrderStatisticsDSWrapper) GetScriptJobStatistics() (*gjson.Result, error) {
	client, err := w.NewClient(w.Config, "coc")
	if err != nil {
		return nil, err
	}

	uri := "/v1/job/script/orders/{execute_uuid}/statistics"
	uri = strings.ReplaceAll(uri, "{execute_uuid}", w.Get("execute_uuid").(string))
	return httphelper.New(client).
		Method("GET").
		URI(uri).
		Request().
		Result()
}

func (w *ScriptOrderStatisticsDSWrapper) getScriptJobStatisticsToSchema(body *gjson.Result) error {
	d := w.ResourceData
	mErr := multierror.Append(nil,
		d.Set("execute_statistics", schemas.SliceToList(body.Get("data.execute_statistics"),
			func(executeStatistics gjson.Result) any {
				return map[string]any{
					"instance_status": executeStatistics.Get("instance_status").Value(),
					"instance_count":  executeStatistics.Get("instance_count").Value(),
					"batch_indexes":   schemas.SliceToIntList(executeStatistics.Get("batch_indexes")),
				}
			},
		)),
	)
	return mErr.ErrorOrNil()
}
