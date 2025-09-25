package aom

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

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API AOM POST /v2/{project_id}/events
func DataSourceEvents() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceEventsRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `The region where the events are located.`,
			},

			// Required parameters.
			"time_range": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The time range for querying events and alarms.`,
			},

			// Optional parameters.
			"step": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: `The statistical step size in milliseconds.`,
			},

			// Attributes.
			"events": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: `The list of events and alarms that matched the filter parameters.`,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The ID of the event or alarm.`,
						},
						"event_sn": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The alarm serial number.`,
						},
						"starts_at": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: `The time when the event or alarm occurred, CST millisecond timestamp.`,
						},
						"ends_at": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: `The time when the event or alarm was cleared, CST millisecond timestamp, 0 means not cleared.`,
						},
						"arrives_at": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: `The time when the event arrived at the system, CST millisecond timestamp.`,
						},
						"timeout": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: `The automatic clearing time for alarms in milliseconds.`,
						},
						"enterprise_project_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The enterprise project ID to which the event or alarm belongs.`,
						},
						"metadata": {
							Type:        schema.TypeMap,
							Computed:    true,
							Elem:        &schema.Schema{Type: schema.TypeString},
							Description: `The detailed information (key/value pair) of the event or alarm.`,
						},
						"annotations": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The additional fields of the event or alarm, in JSON format.`,
						},
						"policy": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The open alarm policy, in JSON format.`,
						},
					},
				},
			},
		},
	}
}

func buildEventsRequestBody(d *schema.ResourceData) map[string]interface{} {
	return utils.RemoveNil(map[string]interface{}{
		"time_range": d.Get("time_range"),
		"step":       utils.ValueIgnoreEmpty(d.Get("step")),
	})
}

func listEvents(client *golangsdk.ServiceClient, d *schema.ResourceData) ([]interface{}, error) {
	var (
		httpUrl = "v2/{project_id}/events?limit={limit}"
		limit   = 1000
		marker  = "0"
		result  = make([]interface{}, 0)
	)

	listPath := client.Endpoint + httpUrl
	listPath = strings.ReplaceAll(listPath, "{project_id}", client.ProjectID)
	listPath = strings.ReplaceAll(listPath, "{limit}", strconv.Itoa(limit))

	opt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json",
		},
		JSONBody: buildEventsRequestBody(d),
	}

	for {
		listPathWithMarker := listPath + fmt.Sprintf("&marker=%s", marker)

		requestResp, err := client.Request("POST", listPathWithMarker, &opt)
		if err != nil {
			return nil, err
		}

		respBody, err := utils.FlattenResponse(requestResp)
		if err != nil {
			return nil, err
		}

		events := utils.PathSearch("events", respBody, make([]interface{}, 0)).([]interface{})
		result = append(result, events...)
		if len(events) < limit {
			break
		}
		marker = utils.PathSearch("page_info.next_marker", respBody, "").(string)
		if marker == "" {
			break
		}
	}

	return result, nil
}

func dataSourceEventsRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	client, err := cfg.NewServiceClient("aom", region)
	if err != nil {
		return diag.Errorf("error creating AOM client: %s", err)
	}

	events, err := listEvents(client, d)
	if err != nil {
		return diag.Errorf("error querying AOM events: %s", err)
	}

	randomUUID, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(randomUUID)

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("events", flattenEvents(events)),
	)
	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenEvents(events []interface{}) []map[string]interface{} {
	if len(events) < 1 {
		return nil
	}

	result := make([]map[string]interface{}, 0, len(events))
	for _, event := range events {
		result = append(result, map[string]interface{}{
			"id":                    utils.PathSearch("id", event, nil),
			"event_sn":              utils.PathSearch("event_sn", event, nil),
			"starts_at":             utils.PathSearch("starts_at", event, nil),
			"ends_at":               utils.PathSearch("ends_at", event, nil),
			"arrives_at":            utils.PathSearch("arrives_at", event, nil),
			"timeout":               utils.PathSearch("timeout", event, nil),
			"enterprise_project_id": utils.PathSearch("enterprise_project_id", event, nil),
			"metadata":              utils.PathSearch("metadata", event, make(map[string]interface{})),
			"annotations":           utils.JsonToString(utils.PathSearch("annotations", event, nil)),
			"policy":                utils.JsonToString(utils.PathSearch("policy", event, nil)),
		})
	}

	return result
}
