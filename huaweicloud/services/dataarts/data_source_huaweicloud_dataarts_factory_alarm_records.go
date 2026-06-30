package dataarts

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

// @API DataArtsStudio GET /v2/{project_id}/factory/alarm-info
func DataSourceFactoryAlarmRecords() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceFactoryAlarmRecordsRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `The region where the alarm records are located.`,
			},

			// Optional parameters.
			"workspace": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The workspace ID.`,
			},
			"start_time": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The start time of the alarm records, in RFC3339 format.`,
			},
			"end_time": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The end time of the alarm records, in RFC3339 format.`,
			},

			// Attributes.
			"records": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: `The list of alarm records.`,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"alarm_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The alarm notification time, in RFC3339 format.`,
						},
						"job_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The name of the job.`,
						},
						"schedule_type": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: `The job instance scheduling mode.`,
						},
						"send_msg": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The send message.`,
						},
						"plan_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The plan time, in RFC3339 format.`,
						},
						"remind_type": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: `The alarm notification type.`,
						},
						"send_status": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: `The send status.`,
						},
						"job_id": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: `The ID of the job.`,
						},
					},
				},
			},
		},
	}
}

func buildAlarmRecordsQueryParams(d *schema.ResourceData) string {
	res := ""

	if startTime, ok := d.GetOk("start_time"); ok {
		// Convert RFC3339 time string to 13-bit timestamp (milliseconds)
		timestamp := utils.ConvertTimeStrToNanoTimestamp(startTime.(string))
		res = fmt.Sprintf("%s&start_time=%d", res, timestamp)
	}
	if endTime, ok := d.GetOk("end_time"); ok {
		// Convert RFC3339 time string to 13-bit timestamp (milliseconds)
		timestamp := utils.ConvertTimeStrToNanoTimestamp(endTime.(string))
		res = fmt.Sprintf("%s&end_time=%d", res, timestamp)
	}

	if res != "" {
		res = "&" + res[1:]
	}
	return res
}

func listAlarmRecords(client *golangsdk.ServiceClient, d *schema.ResourceData) ([]interface{}, error) {
	var (
		httpUrl = "v2/{project_id}/factory/alarm-info?limit={limit}"
		limit   = 100
		offset  = 0
		result  = make([]interface{}, 0)
	)

	listPathWithLimit := client.Endpoint + httpUrl
	listPathWithLimit = strings.ReplaceAll(listPathWithLimit, "{project_id}", client.ProjectID)
	listPathWithLimit = strings.ReplaceAll(listPathWithLimit, "{limit}", strconv.Itoa(limit))
	listPathWithLimit += buildAlarmRecordsQueryParams(d)

	listOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      buildArchitectureMoreHeaders(d.Get("workspace").(string)),
	}

	for {
		listPathWithOffset := fmt.Sprintf("%s&offset=%d", listPathWithLimit, offset)
		requestResp, err := client.Request("GET", listPathWithOffset, &listOpt)
		if err != nil {
			return nil, err
		}

		respBody, err := utils.FlattenResponse(requestResp)
		if err != nil {
			return nil, err
		}

		records := utils.PathSearch("alarm_info", respBody, make([]interface{}, 0)).([]interface{})
		if len(records) < 1 {
			break
		}
		result = append(result, records...)

		if len(records) < limit {
			break
		}
		offset += len(records)
	}

	return result, nil
}

func flattenAlarmRecords(records []interface{}) []interface{} {
	if len(records) < 1 {
		return nil
	}

	result := make([]interface{}, 0, len(records))
	for _, record := range records {
		result = append(result, map[string]interface{}{
			"alarm_time":    utils.FormatTimeStampRFC3339(int64(utils.PathSearch("alarm_time", record, float64(0)).(float64))/1000, false),
			"job_name":      utils.PathSearch("job_name", record, nil),
			"schedule_type": utils.PathSearch("schedule_type", record, nil),
			"send_msg":      utils.PathSearch("send_msg", record, nil),
			"plan_time":     utils.FormatTimeStampRFC3339(int64(utils.PathSearch("plan_time", record, float64(0)).(float64))/1000, false),
			"remind_type":   utils.PathSearch("remind_type", record, nil),
			"send_status":   utils.PathSearch("send_status", record, nil),
			"job_id":        utils.PathSearch("job_id", record, nil),
		})
	}

	return result
}

func dataSourceFactoryAlarmRecordsRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg    = meta.(*config.Config)
		region = cfg.GetRegion(d)
	)

	client, err := cfg.NewServiceClient("dataarts", region)
	if err != nil {
		return diag.Errorf("error creating DataArts Studio client: %s", err)
	}

	records, err := listAlarmRecords(client, d)
	if err != nil {
		return diag.FromErr(err)
	}

	randomUUID, err := uuid.NewRandom()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(randomUUID.String())

	mErr := multierror.Append(
		d.Set("region", region),
		d.Set("records", flattenAlarmRecords(records)),
	)
	return diag.FromErr(mErr.ErrorOrNil())
}
