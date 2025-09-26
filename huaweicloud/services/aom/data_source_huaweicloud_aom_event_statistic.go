package aom

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

// @API AOM POST /v2/{project_id}/events/statistic
func DataSourceEventStatistic() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceEventStatisticRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `The region where the event statistics are located.`,
			},

			// Required parameters.
			"time_range": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The time range for querying event and alarm statistics.`,
			},

			// Optional parameters.
			"step": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: `The statistical step size in milliseconds.`,
			},

			// Attributes.
			"step_result": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: `The statistical step size in milliseconds.`,
			},
			"timestamps": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: `The time series corresponding to the statistical results.`,
				Elem:        &schema.Schema{Type: schema.TypeInt},
			},
			"series": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: `The statistical results for different severity levels at the same time series.`,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"event_severity": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The event or alarm severity level.`,
						},
						"values": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: `The statistical results for events or alarms at each time point.`,
							Elem:        &schema.Schema{Type: schema.TypeInt},
						},
					},
				},
			},
			"summary": {
				Type:        schema.TypeMap,
				Computed:    true,
				Description: `The summary of various alarm information quantities.`,
				Elem:        &schema.Schema{Type: schema.TypeInt},
			},
		},
	}
}

func buildEventStatisticBodyParams(d *schema.ResourceData) map[string]interface{} {
	return utils.RemoveNil(map[string]interface{}{
		"time_range": d.Get("time_range"),
		"step":       utils.ValueIgnoreEmpty(d.Get("step")),
	})
}

func getEventStatistic(client *golangsdk.ServiceClient, d *schema.ResourceData) (interface{}, error) {
	httpUrl := "v2/{project_id}/events/statistic"
	getPath := client.Endpoint + httpUrl
	getPath = strings.ReplaceAll(getPath, "{project_id}", client.ProjectID)

	getOpts := golangsdk.RequestOpts{
		KeepResponseBody: true,
		JSONBody:         buildEventStatisticBodyParams(d),
		MoreHeaders: map[string]string{
			"Content-Type": "application/json",
		},
	}

	requestResp, err := client.Request("POST", getPath, &getOpts)
	if err != nil {
		return nil, err
	}

	return utils.FlattenResponse(requestResp)
}

func flattenEventStatisticSeries(series []interface{}) []map[string]interface{} {
	if len(series) < 1 {
		return nil
	}

	result := make([]map[string]interface{}, 0, len(series))
	for _, item := range series {
		result = append(result, map[string]interface{}{
			"event_severity": utils.PathSearch("event_severity", item, nil),
			"values":         utils.PathSearch("values", item, make([]interface{}, 0)),
		})
	}

	return result
}

func dataSourceEventStatisticRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	client, err := cfg.NewServiceClient("aom", region)
	if err != nil {
		return diag.Errorf("error creating AOM client: %s", err)
	}

	respBody, err := getEventStatistic(client, d)
	if err != nil {
		return diag.Errorf("error querying event statistic information: %s", err)
	}

	randomUUID, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(randomUUID)

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("step_result", utils.PathSearch("step", respBody, nil)),
		d.Set("timestamps", utils.PathSearch("timestamps", respBody, make([]interface{}, 0))),
		d.Set("series", flattenEventStatisticSeries(utils.PathSearch("series", respBody, make([]interface{}, 0)).([]interface{}))),
		d.Set("summary", utils.PathSearch("summary", respBody, make(map[string]interface{}))),
	)
	return diag.FromErr(mErr.ErrorOrNil())
}
