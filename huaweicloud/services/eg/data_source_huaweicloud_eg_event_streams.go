package eg

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

// @API EG GET /v1/{project_id}/eventstreamings
func DataSourceEventStreams() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceEventStreamsRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The region in which to obtain the EG event streams resource.",
			},
			"event_streams": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "The list of event streams.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The ID of the event stream.",
						},
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The name of the event stream.",
						},
						"status": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The status of the event stream.",
						},
						"description": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The description of the event stream.",
						},
						"source": {
							Type:        schema.TypeList,
							Computed:    true,
							Elem:        dataSourceEventStreamSource(),
							Description: "The event source configuration.",
						},
						"sink": {
							Type:        schema.TypeList,
							Computed:    true,
							Elem:        dataSourceEventStreamSink(),
							Description: "The event sink configuration.",
						},
						"rule_config": {
							Type:        schema.TypeList,
							Computed:    true,
							Elem:        dataSourceEventStreamRuleConfig(),
							Description: "The configuration of event rules.",
						},
						"option": {
							Type:        schema.TypeList,
							Computed:    true,
							Elem:        dataSourceEventStreamOption(),
							Description: "The running configuration.",
						},
						"created_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The creation time of the event stream, in RFC3339 format.",
						},
						"updated_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The latest update time of the event stream, in RFC3339 format.",
						},
					},
				},
			},
		},
	}
}

func dataSourceEventStreamSource() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The name of the event source type.",
			},
			"source_kafka": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The configuration of kafka event source, in JSON format.",
			},
			"source_mobile_rocketmq": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The configuration of mobile rocketmq event source, in JSON format.",
			},
			"source_community_rocketmq": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The configuration of community RocketMQ event source, in JSON format.",
			},
			"source_dms_rocketmq": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The configuration of DMS RocketMQ event source, in JSON format.",
			},
		},
	}
}

func dataSourceEventStreamSink() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The name of the event sink type.",
			},
			"sink_fg": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The configuration of function graph event sink, in JSON format.",
			},
			"sink_kafka": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The configuration of Kafka event sink, in JSON format.",
			},
			"sink_obs": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The configuration of OBS event sink, in JSON format.",
			},
		},
	}
}

func dataSourceEventStreamRuleConfig() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"transform": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "The transformation rules.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The type of transformation rule.",
						},
						"value": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The value of transformation rule.",
						},
						"template": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The template of transformation rule.",
						},
					},
				},
			},
			"filter": {
				Type:        schema.TypeMap,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "The filter rules.",
			},
		},
	}
}

func dataSourceEventStreamOption() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"thread_num": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The number of concurrent threads.",
			},
			"batch_window": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"count": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "The number of batch push messages.",
						},
						"time": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "The number of retries.",
						},
						"interval": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "The batch push interval in seconds.",
						},
					},
				},
				Description: "The batch push configuration.",
			},
		},
	}
}

func flattenSource(source interface{}) []interface{} {
	if source == nil {
		return nil
	}

	return []interface{}{
		map[string]interface{}{
			"name":                      utils.PathSearch("name", source, ""),
			"source_kafka":              utils.JsonToString(utils.PathSearch("source_kafka", source, nil)),
			"source_mobile_rocketmq":    utils.JsonToString(utils.PathSearch("source_mobile_rocketmq", source, nil)),
			"source_community_rocketmq": utils.JsonToString(utils.PathSearch("source_community_rocketmq", source, nil)),
			"source_dms_rocketmq":       utils.JsonToString(utils.PathSearch("source_dms_rocketmq", source, nil)),
		},
	}
}

func flattenSink(sink interface{}) []interface{} {
	if sink == nil {
		return nil
	}

	return []interface{}{
		map[string]interface{}{
			"name":       utils.PathSearch("name", sink, ""),
			"sink_fg":    utils.JsonToString(utils.PathSearch("sink_fg", sink, nil)),
			"sink_kafka": utils.JsonToString(utils.PathSearch("sink_kafka", sink, nil)),
			"sink_obs":   utils.JsonToString(utils.PathSearch("sink_obs", sink, nil)),
		},
	}
}

func flattenRuleConfig(rule interface{}) []interface{} {
	if rule == nil {
		return nil
	}

	return []interface{}{
		map[string]interface{}{
			"filter":    utils.PathSearch("filter", rule, nil),
			"transform": flattenTransform(utils.PathSearch("transform", rule, nil)),
		},
	}
}

func flattenTransform(transform interface{}) []interface{} {
	if transform == nil {
		return nil
	}

	return []interface{}{
		map[string]interface{}{
			"type":     utils.PathSearch("type", transform, nil),
			"value":    utils.PathSearch("value", transform, nil),
			"template": utils.PathSearch("template", transform, nil),
		},
	}
}

func flattenOption(option interface{}) []interface{} {
	if option == nil {
		return nil
	}

	return []interface{}{
		map[string]interface{}{
			"thread_num":   utils.PathSearch("thread_num", option, nil),
			"batch_window": flattenBatchWindow(utils.PathSearch("batch_window", option, nil)),
		},
	}
}

func flattenBatchWindow(batchWindow interface{}) []interface{} {
	if batchWindow == nil {
		return nil
	}

	return []interface{}{
		map[string]interface{}{
			"count":    utils.PathSearch("count", batchWindow, nil),
			"time":     utils.PathSearch("time", batchWindow, nil),
			"interval": utils.PathSearch("interval", batchWindow, nil),
		},
	}
}

func flattenEventStreams(streams []interface{}) []interface{} {
	if len(streams) == 0 {
		return nil
	}

	result := make([]interface{}, 0, len(streams))
	for _, item := range streams {
		stream := item.(map[string]interface{})
		eventStream := map[string]interface{}{
			"id":           utils.PathSearch("id", stream, nil),
			"name":         utils.PathSearch("name", stream, nil),
			"status":       utils.PathSearch("status", stream, nil),
			"description":  utils.PathSearch("description", stream, nil),
			"source":       flattenSource(utils.PathSearch("source", stream, nil)),
			"sink":         flattenSink(utils.PathSearch("sink", stream, nil)),
			"rule_config":  flattenRuleConfig(utils.PathSearch("rule_config", stream, nil)),
			"option":       flattenOption(utils.PathSearch("option", stream, nil)),
			"created_time": utils.PathSearch("created_time", stream, nil),
			"updated_time": utils.PathSearch("updated_time", stream, nil),
		}

		result = append(result, eventStream)
	}
	return result
}

func listEventStreams(client *golangsdk.ServiceClient) ([]interface{}, error) {
	var (
		httpUrl = "v1/{project_id}/eventstreamings"
		offset  = 0
		limit   = 500
		result  = make([]interface{}, 0)
	)
	listPath := client.Endpoint + httpUrl
	listPath = strings.ReplaceAll(listPath, "{project_id}", client.ProjectID)

	opt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json",
		},
	}

	for {
		listPathWithOffset := listPath + fmt.Sprintf("?limit=%d&offset=%d", limit, offset)
		requestResp, err := client.Request("GET", listPathWithOffset, &opt)
		if err != nil {
			return nil, err
		}
		respBody, err := utils.FlattenResponse(requestResp)
		if err != nil {
			return nil, err
		}
		eventStreams := utils.PathSearch("items", respBody, make([]interface{}, 0)).([]interface{})
		result = append(result, eventStreams...)
		if len(eventStreams) < limit {
			break
		}
		offset += len(eventStreams)
	}

	return result, nil
}

func dataSourceEventStreamsRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg    = meta.(*config.Config)
		region = cfg.GetRegion(d)
	)

	client, err := cfg.NewServiceClient("eg", region)
	if err != nil {
		return diag.Errorf("error creating EG client: %s", err)
	}

	eventStreams, err := listEventStreams(client)
	if err != nil {
		return diag.Errorf("error querying event streams: %s", err)
	}

	dataSourceId, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(dataSourceId)

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("event_streams", flattenEventStreams(eventStreams)),
	)
	return diag.FromErr(mErr.ErrorOrNil())
}
