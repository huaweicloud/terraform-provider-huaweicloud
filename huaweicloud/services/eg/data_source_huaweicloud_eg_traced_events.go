package eg

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

// @API EG GET /v1/{project_id}/traced-events
func DataSourceTracedEvents() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceTracedEventsRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `The region where the events are located.`,
			},
			"channel_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The ID of the event channel.`,
			},
			"start_time": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The start time of the search time range, in UTC format.`,
			},
			"end_time": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The end time of the search time range, in UTC format.`,
			},
			"event_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The ID of the event.`,
			},
			"source_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The name of the event source.`,
			},
			"event_type": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The type of the event.`,
			},
			"subscription_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The name of the event subscription.`,
			},
			"events": {
				Type:        schema.TypeList,
				Computed:    true,
				Elem:        tracedEventSchema(),
				Description: `The list of traced events that matched filter parameters.`,
			},
		},
	}
}

func tracedEventSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The ID of the event.`,
			},
			"type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The type of the event.`,
			},
			"source_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The name of the event source.`,
			},
			"subscription_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The name of the event subscription.`,
			},
			"received_time": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The time when the event to be received, in UTC format.`,
			},
		},
	}
}

func buildTracedEventsQueryParams(d *schema.ResourceData) string {
	res := ""

	// Required parameters
	res = fmt.Sprintf("%s&channel_id=%v", res, d.Get("channel_id"))

	timestamp := utils.ConvertTimeStrToNanoTimestamp(d.Get("start_time").(string))
	res = fmt.Sprintf("%s&start_time=%v", res, timestamp)

	timestamp = utils.ConvertTimeStrToNanoTimestamp(d.Get("end_time").(string))
	res = fmt.Sprintf("%s&end_time=%v", res, timestamp)

	// Optional parameters
	if v, ok := d.GetOk("event_id"); ok {
		res = fmt.Sprintf("%s&event_id=%v", res, v)
	}
	if v, ok := d.GetOk("source_name"); ok {
		res = fmt.Sprintf("%s&source_name=%v", res, v)
	}
	if v, ok := d.GetOk("event_type"); ok {
		res = fmt.Sprintf("%s&event_type=%v", res, v)
	}
	if v, ok := d.GetOk("subscription_name"); ok {
		res = fmt.Sprintf("%s&subscription_name=%v", res, v)
	}
	return res
}

func listTracedEvents(client *golangsdk.ServiceClient, d *schema.ResourceData) ([]interface{}, error) {
	var (
		httpUrl = "v1/{project_id}/traced-events?limit={limit}"
		offset  = 0
		limit   = 100
		result  = make([]interface{}, 0)
	)

	listPathWithLimit := client.Endpoint + httpUrl
	listPathWithLimit = strings.ReplaceAll(listPathWithLimit, "{project_id}", client.ProjectID)
	listPathWithLimit = strings.ReplaceAll(listPathWithLimit, "{limit}", strconv.Itoa(limit))
	listPathWithLimit += buildTracedEventsQueryParams(d)

	opt := &golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json",
		},
	}

	for {
		listPathWithOffset := fmt.Sprintf("%s&offset=%v", listPathWithLimit, strconv.Itoa(offset))
		requestResp, err := client.Request("GET", listPathWithOffset, opt)
		if err != nil {
			return nil, err
		}

		respBody, err := utils.FlattenResponse(requestResp)
		if err != nil {
			return nil, err
		}

		events := utils.PathSearch("items", respBody, make([]interface{}, 0)).([]interface{})
		result = append(result, events...)
		if len(events) < limit {
			break
		}
		offset += len(events)
	}

	return result, nil
}

func flattenTracedEvents(tracedEvents []interface{}) []interface{} {
	if len(tracedEvents) < 1 {
		return nil
	}

	result := make([]interface{}, 0, len(tracedEvents))
	for _, item := range tracedEvents {
		result = append(result, map[string]interface{}{
			"id":                utils.PathSearch("event_id", item, nil),
			"type":              utils.PathSearch("event_type", item, nil),
			"source_name":       utils.PathSearch("source_name", item, nil),
			"subscription_name": utils.PathSearch("subscription_name", item, nil),
			"received_time":     utils.FormatTimeStampUTC(int64(utils.PathSearch("event_received_time", item, float64(0)).(float64))),
		})
	}

	return result
}

func dataSourceTracedEventsRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg    = meta.(*config.Config)
		region = cfg.GetRegion(d)
	)

	client, err := cfg.NewServiceClient("eg", region)
	if err != nil {
		return diag.Errorf("error creating EG client: %s", err)
	}

	events, err := listTracedEvents(client, d)
	if err != nil {
		return diag.Errorf("error querying EG traced events: %s", err)
	}

	dataSourceId, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(dataSourceId)

	mErr := multierror.Append(
		d.Set("region", region),
		d.Set("events", flattenTracedEvents(events)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}
