package rocketmq

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

// @API RocketMQ GET /v2/{engine}/{project_id}/instances/{instance_id}/messages
func DataSourceDmsRocketMQMessages() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceDmsRocketMQMessagesRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `Specifies the region in which to query the resource. If omitted, the provider-level region will be used.`,
			},
			"instance_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `Specifies the instance ID.`,
			},
			"topic": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `Specifies the topic name.`,
			},
			"start_time": {
				Type:         schema.TypeString,
				Optional:     true,
				RequiredWith: []string{"end_time"},
				Description:  `Specifies the start time, a Unix timestamp in millisecond.`,
			},
			"end_time": {
				Type:         schema.TypeString,
				Optional:     true,
				RequiredWith: []string{"start_time"},
				Description:  `Specifies the end time, a Unix timestamp in millisecond.`,
			},
			"key": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specifies the message key.`,
			},
			"message_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specifies the message ID.`,
			},
			"messages": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: `Indicates the message list.`,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"message_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `Indicates the message ID.`,
						},
						"store_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `Indicates the message stored time.`,
						},
						"born_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `Indicates the message generated time.`,
						},
						"reconsume_times": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: `Indicates the number of retry times.`,
						},
						"body": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `Indicates the message body.`,
						},
						"body_crc": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: `Indicates the message body checksum.`,
						},
						"store_size": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: `Indicates the storage size.`,
						},
						"born_host": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `Indicates the IP address of the host that generates the message.`,
						},
						"store_host": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `Indicates the IP address of the host that stores the message.`,
						},
						"queue_id": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: `Indicates the queue ID.`,
						},
						"queue_offset": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: `Indicates the offset in the queue.`,
						},
						"property_list": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: `Indicates the property list.`,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"name": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: `Indicates the property name.`,
									},
									"value": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: `Indicates the property value.`,
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

func dataSourceDmsRocketMQMessagesRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	client, err := cfg.NewServiceClient("dmsv2", region)
	if err != nil {
		return diag.Errorf("error creating DMS client: %s", err)
	}

	listMessagesHttpUrl := "v2/{engine}/{project_id}/instances/{instance_id}/messages"
	listMessagesPath := client.Endpoint + listMessagesHttpUrl
	listMessagesPath = strings.ReplaceAll(listMessagesPath, "{engine}", "reliability")
	listMessagesPath = strings.ReplaceAll(listMessagesPath, "{project_id}", client.ProjectID)
	listMessagesPath = strings.ReplaceAll(listMessagesPath, "{instance_id}", d.Get("instance_id").(string))
	listMessagesOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"Content-Type": "application/json"},
	}

	// pagelimit is `10`
	listMessagesPath += fmt.Sprintf("?limit=%v", pageLimit)
	listMessagesPath = buildQueryRocketMQMessagesListPath(d, listMessagesPath)

	currentTotal := 0
	results := make([]map[string]interface{}, 0)
	for {
		currentPath := listMessagesPath + fmt.Sprintf("&offset=%d", currentTotal)
		listMessagesResp, err := client.Request("GET", currentPath, &listMessagesOpt)
		if err != nil {
			return diag.Errorf("error retrieving messages: %s", err)
		}
		listMessagesRespBody, err := utils.FlattenResponse(listMessagesResp)
		if err != nil {
			return diag.Errorf("error flatten response: %s", err)
		}

		messages := utils.PathSearch("messages", listMessagesRespBody, make([]interface{}, 0)).([]interface{})
		for _, message := range messages {
			results = append(results, map[string]interface{}{
				"message_id":      utils.PathSearch("msg_id", message, nil),
				"reconsume_times": utils.PathSearch("reconsume_times", message, nil),
				"body":            utils.PathSearch("body", message, nil),
				"body_crc":        utils.PathSearch("body_crc", message, nil),
				"store_size":      utils.PathSearch("store_size", message, nil),
				"born_host":       utils.PathSearch("born_host", message, nil),
				"store_host":      utils.PathSearch("store_host", message, nil),
				"queue_id":        utils.PathSearch("queue_id", message, nil),
				"queue_offset":    utils.PathSearch("queue_offset", message, nil),
				"store_time": utils.FormatTimeStampRFC3339(
					int64(utils.PathSearch("store_timestamp", message, float64(0)).(float64))/1000, true),
				"born_time": utils.FormatTimeStampRFC3339(
					int64(utils.PathSearch("born_timestamp", message, float64(0)).(float64))/1000, true),
				"property_list": flattenMessagesPropertyList(
					utils.PathSearch("property_list", message, make([]interface{}, 0)).([]interface{})),
			})
		}

		// `totalCount` means the number of all `messages`, and type is float64.
		currentTotal += len(messages)
		totalCount := utils.PathSearch("total", listMessagesRespBody, float64(0))
		if int(totalCount.(float64)) <= currentTotal {
			break
		}
	}

	id, err := uuid.GenerateUUID()
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(id)

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("messages", results),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenMessagesPropertyList(paramsList []interface{}) interface{} {
	if len(paramsList) == 0 {
		return nil
	}
	rst := make([]interface{}, 0, len(paramsList))
	for _, params := range paramsList {
		rst = append(rst, map[string]interface{}{
			"name":  utils.PathSearch("name", params, nil),
			"value": utils.PathSearch("value", params, nil),
		})
	}
	return rst
}

func buildQueryRocketMQMessagesListPath(d *schema.ResourceData, listMessagesPath string) string {
	listMessagesPath += fmt.Sprintf("&topic=%v", d.Get("topic"))
	if startTime, ok := d.GetOk("start_time"); ok {
		listMessagesPath += fmt.Sprintf("&start_time=%s", startTime)
		listMessagesPath += fmt.Sprintf("&end_time=%s", d.Get("end_time"))
	}
	if key, ok := d.GetOk("key"); ok {
		listMessagesPath += fmt.Sprintf("&key=%v", key)
	}
	if messageID, ok := d.GetOk("message_id"); ok {
		listMessagesPath += fmt.Sprintf("&msg_id=%v", messageID)
	}

	return listMessagesPath
}
