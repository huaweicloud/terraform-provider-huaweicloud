package lts

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API LTS GET /v2/{project_id}/groups/{log_group_id}/streams/{log_stream_id}/charts
func DataSourceStreamCharts() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceStreamChartsRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `The region in which to query the chart list.`,
			},
			"log_group_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The ID of the log group.`,
			},
			"log_stream_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The ID of the log stream.`,
			},
			"charts": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: `The list of charts.`,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The ID of the chart.`,
						},
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The name of the chart.`,
						},
						"type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The type of the chart.`,
						},
						"log_group_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The ID of the log group.`,
						},
						"log_group_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The name of the log group.`,
						},
						"log_stream_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The ID of the log stream.`,
						},
						"log_stream_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The name of the log stream.`,
						},
						"sql": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The SQL statement of the chart.`,
						},
						"config": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: `The configuration of the chart.`,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"page_size": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: `The page size of the chart.`,
									},
									"can_sort": {
										Type:        schema.TypeBool,
										Computed:    true,
										Description: `Whether to enable sorting.`,
									},
									"can_search": {
										Type:        schema.TypeBool,
										Computed:    true,
										Description: `Whether to enable search.`,
									},
								},
							},
						},
					},
				},
			},
		},
	}
}

func listStreamCharts(client *golangsdk.ServiceClient, d *schema.ResourceData) ([]interface{}, error) {
	var (
		httpUrl  = "v2/{project_id}/groups/{log_group_id}/streams/{log_stream_id}/charts?limit={limit}"
		limit    = 100
		offset   = 0
		result   = make([]interface{}, 0)
		respBody interface{}
	)

	listPath := client.Endpoint + httpUrl
	listPath = strings.ReplaceAll(listPath, "{project_id}", client.ProjectID)
	listPath = strings.ReplaceAll(listPath, "{log_group_id}", d.Get("log_group_id").(string))
	listPath = strings.ReplaceAll(listPath, "{log_stream_id}", d.Get("log_stream_id").(string))
	listPath = strings.ReplaceAll(listPath, "{limit}", strconv.Itoa(limit))

	opt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json",
		},
	}

	for {
		listPathWithOffset := listPath + fmt.Sprintf("&offset=%d", offset)
		requestResp, err := client.Request("GET", listPathWithOffset, &opt)
		if err != nil {
			return nil, err
		}

		respBody, err = utils.FlattenResponse(requestResp)
		if err != nil {
			return nil, err
		}

		charts := utils.PathSearch("results", respBody, make([]interface{}, 0)).([]interface{})
		result = append(result, charts...)

		if len(charts) < limit {
			break
		}
		offset += len(charts)
	}

	return result, nil
}

func flattenStreamCharts(charts []interface{}) []map[string]interface{} {
	if len(charts) < 1 {
		return nil
	}

	result := make([]map[string]interface{}, 0, len(charts))
	for _, chart := range charts {
		result = append(result, map[string]interface{}{
			"id":              utils.PathSearch("id", chart, nil),
			"name":            utils.PathSearch("title", chart, nil),
			"type":            utils.PathSearch("type", chart, nil),
			"log_group_id":    utils.PathSearch("log_group_id", chart, nil),
			"log_group_name":  utils.PathSearch("log_group_name", chart, nil),
			"log_stream_id":   utils.PathSearch("log_stream_id", chart, nil),
			"log_stream_name": utils.PathSearch("log_stream_name", chart, nil),
			"sql":             utils.PathSearch("sql", chart, nil),
			"config": flattenStreamChartConfig(
				utils.PathSearch("config", chart, make(map[string]interface{}, 0)).(map[string]interface{})),
		})
	}

	return result
}

func flattenStreamChartConfig(configs map[string]interface{}) []map[string]interface{} {
	if len(configs) < 1 {
		return nil
	}

	return []map[string]interface{}{
		{
			"page_size":  utils.PathSearch("pageSize", configs, nil),
			"can_sort":   utils.PathSearch("canSort", configs, nil),
			"can_search": utils.PathSearch("canSearch", configs, nil),
		},
	}
}

func dataSourceStreamChartsRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	client, err := cfg.NewServiceClient("lts", region)
	if err != nil {
		return diag.Errorf("error creating LTS client: %s", err)
	}

	charts, err := listStreamCharts(client, d)
	if err != nil {
		return common.CheckDeletedDiag(d, err, "error retrieving stream charts")
	}

	randomUUID, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(randomUUID)

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("charts", flattenStreamCharts(charts)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}
