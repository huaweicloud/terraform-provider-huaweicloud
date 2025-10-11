package kafka

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

// @API Kafka GET /v2/{engine}/{project_id}/instances/{instance_id}/groups/{group}/message-offset
func DataSourceConsumerGroupMessageOffsets() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceConsumerGroupMessageOffsetsRead,
		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `The region where the consumer group message offsets are located.`,
			},
			"instance_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The ID of the Kafka instance.`,
			},
			"group": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The name of the consumer group.`,
			},
			"topic": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The name of the topic.`,
			},
			"message_offsets": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"partition": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: `The name of the partition.`,
						},
						"message_current_offset": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: `The current offset of the message.`,
						},
						"message_log_start_offset": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: `The start offset of the message.`,
						},
						"message_log_end_offset": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: `The end offset of the message.`,
						},
						"consumer_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The consumer ID of the consumed message.`,
						},
						"host": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The consumer address of the consumed message.`,
						},
						"client_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The ID of the client.`,
						},
					},
				},
				Description: `The list of consumer group message offsets.`,
			},
		},
	}
}

func listConsumerGroupMessageOffsets(client *golangsdk.ServiceClient, d *schema.ResourceData) ([]interface{}, error) {
	var (
		httpUrl = "v2/kafka/{project_id}/instances/{instance_id}/groups/{group}/message-offset"
		result  = make([]interface{}, 0)
		offset  = 0
		// The limit maximum value is 50, default is 10.
		limit   = 50
		listOpt = golangsdk.RequestOpts{
			KeepResponseBody: true,
			MoreHeaders: map[string]string{
				"Content-Type": "application/json;charset=utf-8",
			},
		}
	)

	listPath := client.Endpoint + httpUrl
	listPath = strings.ReplaceAll(listPath, "{project_id}", client.ProjectID)
	listPath = strings.ReplaceAll(listPath, "{instance_id}", d.Get("instance_id").(string))
	listPath = strings.ReplaceAll(listPath, "{group}", d.Get("group").(string))
	listPath = fmt.Sprintf("%s?topic=%s&limit=%d", listPath, d.Get("topic").(string), limit)

	for {
		listPathWithOffset := listPath + fmt.Sprintf("&offset=%d", offset)
		resp, err := client.Request("GET", listPathWithOffset, &listOpt)
		if err != nil {
			return nil, err
		}

		respBody, err := utils.FlattenResponse(resp)
		if err != nil {
			return nil, err
		}

		messageOffsets := utils.PathSearch("group_message_offsets", respBody, make([]interface{}, 0)).([]interface{})
		result = append(result, messageOffsets...)
		if len(messageOffsets) < limit {
			break
		}

		offset += len(messageOffsets)
	}

	return result, nil
}

func dataSourceConsumerGroupMessageOffsetsRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg    = meta.(*config.Config)
		region = cfg.GetRegion(d)
	)
	client, err := cfg.NewServiceClient("dmsv2", region)
	if err != nil {
		return diag.Errorf("error creating DMS client: %s", err)
	}

	groupMessageOffsets, err := listConsumerGroupMessageOffsets(client, d)
	if err != nil {
		return diag.Errorf("error querying message offset list under consumer group (%s): %s", d.Get("group").(string), err)
	}

	randomUUID, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(randomUUID)

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("message_offsets", flattenConsumerGroupMessageOffsets(groupMessageOffsets)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenConsumerGroupMessageOffsets(groupMessageOffsets []interface{}) []interface{} {
	if len(groupMessageOffsets) == 0 {
		return nil
	}

	rst := make([]interface{}, 0, len(groupMessageOffsets))
	for _, v := range groupMessageOffsets {
		rst = append(rst, map[string]interface{}{
			"partition":                utils.PathSearch("partition", v, nil),
			"message_current_offset":   utils.PathSearch("message_current_offset", v, nil),
			"message_log_start_offset": utils.PathSearch("message_log_start_offset", v, nil),
			"message_log_end_offset":   utils.PathSearch("message_log_end_offset", v, nil),
			"consumer_id":              utils.PathSearch("consumer_id", v, nil),
			"host":                     utils.PathSearch("host", v, nil),
			"client_id":                utils.PathSearch("client_id", v, nil),
		})
	}
	return rst
}
