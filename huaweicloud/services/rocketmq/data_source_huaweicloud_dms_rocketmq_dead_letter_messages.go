package rocketmq

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

// @API RocketMQ POST /v2/{project_id}/instances/{instance_id}/messages/export
func DataSourceDeadLetterMessages() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceDeadLetterMessagesRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `The region where the dead letter messages are located.`,
			},
			"instance_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The ID of the RocketMQ instance.`,
			},
			"topic": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The name of the topic to which the dead letter messages belong.`,
			},
			"msg_id_list": {
				Type:        schema.TypeList,
				Required:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: `The list of dead letter message IDs.`,
			},
			"messages": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"instance_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The ID of the RocketMQ instance.`,
						},
						"msg_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The ID of the dead letter message.`,
						},
						"topic": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The name of the topic to which the dead letter message belongs.`,
						},
						"store_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The time when the dead letter message was stored, in RFC3339 format.`,
						},
						"born_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The time when the dead letter message was generated, in RFC3339 format.`,
						},
						"reconsume_times": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: `The number of times the message has been retried.`,
						},
						"body": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The body of the message.`,
						},
						"body_crc": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: `The checksum of the message body.`,
						},
						"store_size": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: `The storage size of the message.`,
						},
						"born_host": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The IP address of the host that generated the message.`,
						},
						"store_host": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The IP address of the host that stored the message.`,
						},
						"queue_id": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: `The ID of the queue.`,
						},
						"queue_offset": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: `The offset in the queue.`,
						},
						"property_list": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"name": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: `The name of the property.`,
									},
									"value": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: `The value of the property.`,
									},
								},
							},
							Description: `The list of message properties.`,
						},
					},
				},
				Description: `All dead letter messages that match the filter parameters.`,
			},
		},
	}
}

func buildDeadLetterMessagesQueryParams(d *schema.ResourceData) map[string]interface{} {
	return map[string]interface{}{
		"topic":       d.Get("topic"),
		"msg_id_list": d.Get("msg_id_list"),
	}
}

func dataSourceDeadLetterMessagesRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	client, err := cfg.NewServiceClient("dmsv2", region)
	if err != nil {
		return diag.Errorf("error creating DMS client: %s", err)
	}

	httpUrl := "v2/{project_id}/instances/{instance_id}/messages/export"
	listPath := client.Endpoint + httpUrl
	listPath = strings.ReplaceAll(listPath, "{project_id}", client.ProjectID)
	listPath = strings.ReplaceAll(listPath, "{instance_id}", d.Get("instance_id").(string))

	listOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"Content-Type": "application/json"},
		JSONBody:         utils.RemoveNil(buildDeadLetterMessagesQueryParams(d)),
	}

	listResp, err := client.Request("POST", listPath, &listOpt)
	if err != nil {
		return diag.Errorf("error retrieving dead letter messages: %s", err)
	}

	listRespBody, err := utils.FlattenResponse(listResp)
	if err != nil {
		return diag.FromErr(err)
	}

	randomId, err := uuid.GenerateUUID()
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(randomId)

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("messages", flattenDeadLetterMessages(utils.PathSearch("[]", listRespBody, make([]interface{}, 0)).([]interface{}))),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenDeadLetterMessages(messages []interface{}) []interface{} {
	if len(messages) == 0 {
		return nil
	}

	results := make([]interface{}, 0, len(messages))
	for _, message := range messages {
		results = append(results, map[string]interface{}{
			"instance_id":     utils.PathSearch("instance_id", message, nil),
			"msg_id":          utils.PathSearch("msg_id", message, nil),
			"topic":           utils.PathSearch("topic", message, nil),
			"reconsume_times": utils.PathSearch("reconsume_times", message, nil),
			"body":            utils.PathSearch("body", message, nil),
			"body_crc":        utils.PathSearch("body_crc", message, nil),
			"store_size":      utils.PathSearch("store_size", message, nil),
			"born_host":       utils.PathSearch("born_host", message, nil),
			"store_host":      utils.PathSearch("store_host", message, nil),
			"queue_id":        utils.PathSearch("queue_id", message, nil),
			"queue_offset":    utils.PathSearch("queue_offset", message, nil),
			"store_time": utils.FormatTimeStampRFC3339(int64(utils.PathSearch("store_timestamp", message,
				float64(0)).(float64))/1000, false),
			"born_time": utils.FormatTimeStampRFC3339(int64(utils.PathSearch("born_timestamp", message,
				float64(0)).(float64))/1000, false),
			"property_list": flattenDeadLetterMessagePropertyList(
				utils.PathSearch("property_list", message, make([]interface{}, 0)).([]interface{})),
		})
	}
	return results
}

func flattenDeadLetterMessagePropertyList(properties []interface{}) interface{} {
	if len(properties) == 0 {
		return nil
	}

	rst := make([]interface{}, 0, len(properties))
	for _, property := range properties {
		rst = append(rst, map[string]interface{}{
			"name":  utils.PathSearch("name", property, nil),
			"value": utils.PathSearch("value", property, nil),
		})
	}
	return rst
}
