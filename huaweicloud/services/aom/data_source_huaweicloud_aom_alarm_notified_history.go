package aom

import (
	"context"
	"strconv"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API AOM POST /v2/{project_id}/events
func DataSourceAlarmNotifiedHistory() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceAlarmNotifiedHistoryRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `The region where the alarm histories are located.`,
			},

			// Required parameters
			"time_range": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The time range for querying alarm histories.`,
				// Description: `The time range for querying alarm histories in the format: startTimeInMillis.endTimeInMillis.durationInMinutes.`,
			},

			// Optional parameters
			// type
			"alarm_type": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The type of alarm to query.`,
				// Description: `The type of alarm to query. Valid values: active_alert, history_alert.`,
			},
			// event_type
			"type": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The event type to filter alarm histories.`,
			},
			// event_severity
			"severity": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The severity level to filter alarm histories.`,
				// Description: `The severity level to filter alarm histories. Valid values: Critical, Major, Minor, Info.`,
			},

			// Attributes
			"events": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						// event_id
						"id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The ID of the alarm event.`,
						},
						// event_name
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The name of the alarm event.`,
						},
						// event_type
						"type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The type of the alarm event.`,
						},
						// event_severity
						"severity": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The severity level of the alarm event.`,
						},
						"resource_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The ID of the resource associated with the alarm.`,
						},
						"resource_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The name of the resource associated with the alarm.`,
						},
						"start_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The start time of the alarm event.`,
						},
						"end_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The end time of the alarm event.`,
						},
						"detail": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The detailed information of the alarm event.`,
						},
					},
				},
				Description: `The list of alarm notification histories.`,
			},
		},
	}
}

func dataSourceAlarmNotifiedHistoryRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		httpUrl = "v2/{project_id}/events?limit={limit}"
		limit   = 1000
		marker  = 0
	)

	client, err := cfg.NewServiceClient("aom", region)
	if err != nil {
		return diag.Errorf("error creating AOM client: %s", err)
	}

	listPathWithLimit := client.Endpoint + httpUrl
	listPathWithLimit = strings.ReplaceAll(listPathWithLimit, "{project_id}", client.ProjectID)
	listPathWithLimit = strings.ReplaceAll(listPathWithLimit, "{limit}", strconv.Itoa(limit))
	listOpts := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json",
		},
		JSONBody: utils.RemoveNil(buildAlarmHistoryQueryParams(d)),
	}

	// Send API request
	resp, err := client.Request("POST", listPathWithLimit, &listOpts)
	if err != nil {
		return diag.Errorf("error querying alarm histories: %s", err)
	}

	respBody, err := utils.FlattenResponse(resp)
	if err != nil {
		return diag.FromErr(err)
	}

	randomUUID, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(randomUUID)

	mErr := multierror.Append(
		d.Set("region", region),
		d.Set("events", flattenAlarmEvents(respBody)),
	)
	return diag.FromErr(mErr.ErrorOrNil())
}

func buildAlarmHistoryQueryParams(d *schema.ResourceData) map[string]interface{} {
	queryParams := map[string]interface{}{
		"time_range": d.Get("time_range").(string),
	}

	if v, ok := d.GetOk("type"); ok {
		queryParams["type"] = v.(string)
	}

	// Build metadata filters
	var metadataFilters []map[string]interface{}
	if v, ok := d.GetOk("event_type"); ok {
		metadataFilters = append(metadataFilters, map[string]interface{}{
			"key":      "event_type",
			"value":    []string{v.(string)},
			"relation": "AND",
		})
	}

	if v, ok := d.GetOk("event_severity"); ok {
		metadataFilters = append(metadataFilters, map[string]interface{}{
			"key":      "event_severity",
			"value":    []string{v.(string)},
			"relation": "AND",
		})
	}

	if len(metadataFilters) > 0 {
		queryParams["metadata_relation"] = metadataFilters
	}

	if v, ok := d.GetOk("limit"); ok {
		queryParams["limit"] = v.(int)
	}

	if v, ok := d.GetOk("marker"); ok {
		queryParams["marker"] = v.(string)
	}

	return queryParams
}

func flattenAlarmEvents(resp interface{}) []map[string]interface{} {
	if resp == nil {
		return nil
	}

	events := utils.PathSearch("events", resp, make([]interface{}, 0)).([]interface{})
	result := make([]map[string]interface{}, len(events))

	for i, event := range events {
		result[i] = map[string]interface{}{
			"event_id":       utils.PathSearch("event_id", event, ""),
			"event_name":     utils.PathSearch("event_name", event, ""),
			"event_type":     utils.PathSearch("event_type", event, ""),
			"event_severity": utils.PathSearch("event_severity", event, ""),
			"resource_id":    utils.PathSearch("resource_id", event, ""),
			"resource_name":  utils.PathSearch("resource_name", event, ""),
			"start_time":     utils.FormatTimeStampRFC3339(utils.PathSearch("start_time", event, 0).(int64)/1000, false),
			"end_time":       utils.FormatTimeStampRFC3339(utils.PathSearch("end_time", event, 0).(int64)/1000, false),
			"detail":         utils.PathSearch("detail", event, ""),
		}
	}

	return result
}
