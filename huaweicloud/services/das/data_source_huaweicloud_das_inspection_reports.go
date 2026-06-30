package das

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/google/uuid"
	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API DAS GET /v3/{project_id}/batch-inspection/health-report-list
func DataSourceInspectionReports() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceInspectionReportsRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `The region where the inspection reports are located.`,
			},

			// Required parameters.
			"start_time": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The start time of the inspection report, in RFC3339 format.`,
			},
			"end_time": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The end time of the inspection report, in RFC3339 format.`,
			},
			"datastore_type": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The database type.`,
			},

			// Optional parameters.
			"health_rank": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The health rank of the inspection report.`,
			},
			"sort_field": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The field used for sorting.`,
			},
			"asc": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: `Whether to sort in ascending order.`,
			},

			// Attributes.
			"reports": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: `The list of inspection reports that matched filter parameters.`,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"task_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The ID of the inspection report.`,
						},
						"instance_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The ID of the instance.`,
						},
						"instance_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The name of the instance.`,
						},
						"cpu": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: `The CPU size.`,
						},
						"mem": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: `The memory size in GB.`,
						},
						"disk_size": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: `The disk size in GB.`,
						},
						"created_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The generation time of the inspection report, in RFC3339 format.`,
						},
						"start_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The start time of the diagnosis, in RFC3339 format.`,
						},
						"end_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The end time of the diagnosis, in RFC3339 format.`,
						},
						"health_rank": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The health rank of the instance.`,
						},
						"score": {
							Type:        schema.TypeFloat,
							Computed:    true,
							Description: `The score of the inspection.`,
						},
						"lost_points_details": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: `The list of lost points details.`,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"risk_level": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: `The risk level.`,
									},
									"metric": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: `The metric name.`,
									},
									"metric_value": {
										Type:        schema.TypeFloat,
										Computed:    true,
										Description: `The value of the metric.`,
									},
									"deducted_points": {
										Type:        schema.TypeFloat,
										Computed:    true,
										Description: `The deducted points.`,
									},
									"deducted_condition": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: `The deducted condition.`,
									},
									"deducted_formula": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: `The deducted formula.`,
									},
									"suggestions": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: `The optimization suggestions.`,
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

func dataSourceInspectionReportsRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	client, err := cfg.NewServiceClient("das", region)
	if err != nil {
		return diag.Errorf("error creating DAS client: %s", err)
	}

	resp, err := listInspectionReports(d, client)
	if err != nil {
		return diag.Errorf("error querying DAS inspection reports: %s", err)
	}

	randomUUID, err := uuid.NewRandom()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(randomUUID.String())

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("reports", flattenInspectionReports(resp)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func buildInspectionReportsQueryParams(d *schema.ResourceData) string {
	startTime := utils.ConvertTimeStrToNanoTimestamp(d.Get("start_time").(string))
	endTime := utils.ConvertTimeStrToNanoTimestamp(d.Get("end_time").(string))

	res := fmt.Sprintf("&start_at=%v&end_at=%v&datastore_type=%v",
		startTime, endTime, d.Get("datastore_type"))

	if v, ok := d.GetOk("health_rank"); ok {
		res = fmt.Sprintf("%s&health_rank=%v", res, v)
	}

	if v, ok := d.GetOk("sort_field"); ok {
		res = fmt.Sprintf("%s&sort_field=%v", res, v)
	}

	if v, ok := d.GetOk("asc"); ok {
		res = fmt.Sprintf("%s&asc=%v", res, v)
	}

	return res
}

func listInspectionReports(d *schema.ResourceData, client *golangsdk.ServiceClient) ([]interface{}, error) {
	var (
		httpUrl = "v3/{project_id}/batch-inspection/health-report-list?limit={limit}"
		limit   = 200
		offset  = 0
		result  = make([]interface{}, 0)
	)

	listPath := client.Endpoint + httpUrl
	listPath = strings.ReplaceAll(listPath, "{project_id}", client.ProjectID)
	listPath = strings.ReplaceAll(listPath, "{limit}", strconv.Itoa(limit))
	listPath += buildInspectionReportsQueryParams(d)

	opt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json",
		},
	}

	for {
		listPathWithOffset := fmt.Sprintf("%s&offset=%d", listPath, offset)
		requestResp, err := client.Request("GET", listPathWithOffset, &opt)
		if err != nil {
			return nil, err
		}

		respBody, err := utils.FlattenResponse(requestResp)
		if err != nil {
			return nil, err
		}

		reports := utils.PathSearch("batch_inspection_report_list", respBody, make([]interface{}, 0)).([]interface{})
		result = append(result, reports...)
		if len(reports) < limit {
			break
		}

		offset += len(reports)
	}

	return result, nil
}

func flattenInspectionReports(items []interface{}) []map[string]interface{} {
	if len(items) < 1 {
		return nil
	}

	result := make([]map[string]interface{}, 0, len(items))
	for _, item := range items {
		result = append(result, map[string]interface{}{
			"task_id":       utils.PathSearch("task_id", item, nil),
			"instance_id":   utils.PathSearch("instance_id", item, nil),
			"instance_name": utils.PathSearch("instance_name", item, nil),
			"cpu":           utils.PathSearch("cpu", item, nil),
			"mem":           utils.PathSearch("mem", item, nil),
			"disk_size":     utils.PathSearch("disk_size", item, nil),
			"health_rank":   utils.PathSearch("health_rank", item, nil),
			"score":         utils.PathSearch("score", item, nil),
			"created_time": utils.FormatTimeStampRFC3339(utils.ConvertTimeStrToNanoTimestamp(
				utils.PathSearch("create_time", item, "").(string))/1000, false),
			"start_time": utils.FormatTimeStampRFC3339(utils.ConvertTimeStrToNanoTimestamp(
				utils.PathSearch("start_time", item, "").(string))/1000, false),
			"end_time": utils.FormatTimeStampRFC3339(utils.ConvertTimeStrToNanoTimestamp(
				utils.PathSearch("end_time", item, "").(string))/1000, false),
			"lost_points_details": flattenLostPointsDetails(
				utils.PathSearch("lost_points_detail_list", item, make([]interface{}, 0)).([]interface{})),
		})
	}

	return result
}

func flattenLostPointsDetails(items []interface{}) []map[string]interface{} {
	if len(items) < 1 {
		return nil
	}

	result := make([]map[string]interface{}, 0, len(items))
	for _, item := range items {
		result = append(result, map[string]interface{}{
			"risk_level":         utils.PathSearch("risk_level", item, nil),
			"metric":             utils.PathSearch("metric", item, nil),
			"metric_value":       utils.PathSearch("metric_value", item, nil),
			"deducted_points":    utils.PathSearch("deducted_points", item, nil),
			"deducted_condition": utils.PathSearch("deducted_condition", item, nil),
			"deducted_formula":   utils.PathSearch("deducted_formula", item, nil),
			"suggestions":        utils.PathSearch("suggestions", item, nil),
		})
	}

	return result
}
