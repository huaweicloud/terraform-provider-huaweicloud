package ces

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API CES GET /V1.0/{project_id}/event/{event_name}
func DataSourceCesEventDetails() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceCesEventDetailsRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `Specifies the region in which to query the resource. If omitted, the provider-level region will be used.`,
			},
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `Specifies the event name.`,
			},
			"type": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `Specifies the event type.`,
			},
			"source": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specifies the event source.`,
			},
			"level": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specifies the event severity.`,
			},
			"user": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specifies the name of the user for reporting event monitoring data.`,
			},
			"state": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specifies the event status.`,
			},
			"from": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specifies the start time of the query.`,
			},
			"to": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specifies the end time of the query.`,
			},
			"event_info": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: `The event information.`,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"event_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The event ID.`,
						},
						"event_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The event name.`,
						},
						"event_source": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The event source.`,
						},
						"time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The time when the event occurred.`,
						},
						"detail": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: `The event detail.`,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"event_state": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: `The event status.`,
									},
									"event_level": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: `The event level.`,
									},
									"event_user": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: `The event user.`,
									},
									"content": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: `The event content.`,
									},
									"group_id": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: `The group that the event belongs to.`,
									},
									"resource_id": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: `The resource ID.`,
									},
									"resource_name": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: `The resource name.`,
									},
									"event_type": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: `The event type.`,
									},
									"dimensions": {
										Type:        schema.TypeList,
										Computed:    true,
										Description: `The resource dimensions.`,
										Elem:        eventInfoDetailDimensionsElem(),
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

// eventInfoDetailDimensionsElem
// The Elem of "event_info.detail.dimensions"
func eventInfoDetailDimensionsElem() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The resource dimension name.`,
			},
			"value": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The resource dimension value.`,
			},
		},
	}
}

func dataSourceCesEventDetailsRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	client, err := cfg.NewServiceClient("ces", region)
	if err != nil {
		return diag.Errorf("error creating CES client: %s", err)
	}

	result, err := getDetails(client, d)
	if err != nil {
		return diag.Errorf("error retrieving CES event details: %s", err)
	}

	uuid, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(uuid)

	var mErr *multierror.Error
	mErr = multierror.Append(
		mErr,
		d.Set("region", region),
		d.Set("event_info", flattenEventDetails(result)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func getDetails(client *golangsdk.ServiceClient, d *schema.ResourceData) ([]interface{}, error) {
	httpUrl := "V1.0/{project_id}/event/{event_name}"
	path := client.Endpoint + httpUrl
	path = strings.ReplaceAll(path, "{project_id}", client.ProjectID)
	path = strings.ReplaceAll(path, "{event_name}", d.Get("name").(string))
	path += fmt.Sprintf("?event_type=%v", d.Get("type"))

	params, err := buildCesEventDetailsQueryParams(d)
	if err != nil {
		return nil, fmt.Errorf("error building CES event details query params: %s", err)
	}
	path += params

	opt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"Content-Type": "application/json; charset=UTF-8"},
	}

	rst := make([]interface{}, 0)
	start := 0
	for {
		path := fmt.Sprintf("%s&limit=100&start=%d", path, start)
		resp, err := client.Request("GET", path, &opt)

		if err != nil {
			return nil, err
		}

		respBody, err := utils.FlattenResponse(resp)
		if err != nil {
			return nil, err
		}

		rawEventInfo := utils.PathSearch("event_info", respBody, make([]interface{}, 0))
		eventInfo := rawEventInfo.([]interface{})
		if len(eventInfo) > 0 {
			rst = append(rst, eventInfo...)
		}

		start += 100
		total := utils.PathSearch("meta_data.total", respBody, float64(0))
		if int(total.(float64)) <= start {
			return rst, nil
		}
	}
}

func buildCesEventDetailsQueryParams(d *schema.ResourceData) (string, error) {
	res := ""

	if v, ok := d.GetOk("source"); ok {
		res = fmt.Sprintf("%s&event_source=%v", res, v)
	}
	if v, ok := d.GetOk("level"); ok {
		res = fmt.Sprintf("%s&event_level=%v", res, v)
	}
	if v, ok := d.GetOk("user"); ok {
		res = fmt.Sprintf("%s&event_user=%v", res, v)
	}
	if v, ok := d.GetOk("state"); ok {
		res = fmt.Sprintf("%s&event_state=%v", res, v)
	}
	if v, ok := d.GetOk("from"); ok {
		startTime, err := parseTimeToTimestamp(v.(string))
		if err != nil {
			return "", err
		}
		res = fmt.Sprintf("%s&from=%v", res, startTime*1000)
	}
	if v, ok := d.GetOk("to"); ok {
		endTime, err := parseTimeToTimestamp(v.(string))
		if err != nil {
			return "", err
		}
		res = fmt.Sprintf("%s&to=%v", res, endTime*1000)
	}
	return res, nil
}

func flattenEventDetails(resp []interface{}) []interface{} {
	if resp == nil {
		return nil
	}

	rst := make([]interface{}, 0, len(resp))
	for _, v := range resp {
		eventTime := int64(utils.PathSearch("time", v, float64(0)).(float64) / 1000)
		rst = append(rst, map[string]interface{}{
			"event_name":   utils.PathSearch("event_name", v, nil),
			"event_source": utils.PathSearch("event_source", v, nil),
			"time":         utils.FormatTimeStampRFC3339(eventTime, true),
			"detail":       flattenEventDetail(v),
			"event_id":     utils.PathSearch("event_id", v, ""),
		})
	}
	return rst
}

func flattenEventDetail(resp interface{}) []interface{} {
	detail := utils.PathSearch("detail", resp, nil)
	if detail == nil {
		return nil
	}

	rst := map[string]interface{}{
		"event_state":   utils.PathSearch("event_state", detail, nil),
		"event_level":   utils.PathSearch("event_level", detail, nil),
		"event_user":    utils.PathSearch("event_user", detail, nil),
		"content":       utils.PathSearch("content", detail, nil),
		"group_id":      utils.PathSearch("group_id", detail, nil),
		"resource_id":   utils.PathSearch("resource_id", detail, nil),
		"resource_name": utils.PathSearch("resource_name", detail, nil),
		"event_type":    utils.PathSearch("event_type", detail, nil),
		"dimensions":    flattenEventDetailDimensions(detail),
	}

	return []interface{}{rst}
}

func flattenEventDetailDimensions(resp interface{}) []interface{} {
	raw := utils.PathSearch("dimensions", resp, nil)
	if raw == nil {
		return nil
	}

	dimensions := raw.([]interface{})
	rst := make([]interface{}, 0, len(dimensions))
	for _, v := range dimensions {
		rst = append(rst, map[string]interface{}{
			"name":  utils.PathSearch("name", v, nil),
			"value": utils.PathSearch("value", v, nil),
		})
	}
	return rst
}
