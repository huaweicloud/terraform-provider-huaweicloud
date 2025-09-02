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

// @API RocketMQ GET /v2/{project_id}/instances/{instance_id}/groups/{group}/topics
func DataSourceConsumerGroupTopics() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceConsumerGroupTopicsRead,
		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `The region where the RocketMQ instance and consumer group are located.`,
			},
			"instance_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The ID of the RocketMQ instance.`,
			},
			"group": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The name of the consumer group.`,
			},
			"topic": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The name of the topic to be queried.`,
			},
			"topics": {
				Type:        schema.TypeList,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: `The list of topics consumed by the consumer group.`,
			},
			"lag": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: `The number of consumption accumulations.`,
			},
			"max_offset": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: `The total number of messages.`,
			},
			"consumer_offset": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: `The number of consumed messages.`,
			},
			"brokers": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"broker_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The name of the broker.`,
						},
						"queues": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: `The queue details of the associated broker.`,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"id": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: `The ID of the queue.`,
									},
									"lag": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: `The number of consumption accumulations.`,
									},
									"broker_offset": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: `The total number of messages.`,
									},
									"consumer_offset": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: `The number of consumed messages.`,
									},
									"last_message_time": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: `The storage time of the latest consumed message, in RFC3339 format.`,
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

func buildConsumerGroupTopicsQueryParams(d *schema.ResourceData) string {
	res := ""

	if topic, ok := d.GetOk("topic"); ok {
		res += fmt.Sprintf("&topic=%s", topic)
	}
	return res
}

func getConsumerGroupTopics(client *golangsdk.ServiceClient, d *schema.ResourceData) ([]interface{}, interface{}, error) {
	var (
		// The limit maximum value is 50, default is 10.
		httpUrl  = "v2/{project_id}/instances/{instance_id}/groups/{group}/topics?limit=50"
		offset   = 0
		result   = make([]interface{}, 0)
		respBody interface{}
	)

	listPath := client.Endpoint + httpUrl
	listPath = strings.ReplaceAll(listPath, "{project_id}", client.ProjectID)
	listPath = strings.ReplaceAll(listPath, "{instance_id}", d.Get("instance_id").(string))
	listPath = strings.ReplaceAll(listPath, "{group}", d.Get("group").(string))
	listPath += buildConsumerGroupTopicsQueryParams(d)

	getOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json",
		},
	}

	for {
		listPathWithOffset := fmt.Sprintf("%s&offset=%d", listPath, offset)
		requestResp, err := client.Request("GET", listPathWithOffset, &getOpt)
		if err != nil {
			return nil, nil, err
		}

		respBody, err = utils.FlattenResponse(requestResp)
		if err != nil {
			return nil, nil, err
		}

		topics := utils.PathSearch("topics", respBody, make([]interface{}, 0)).([]interface{})
		result = append(result, topics...)
		// The `offset` cannot be greater than or equal to the total number of topics. Otherwise, the response is as follows:
		// {"error_code": "DMS.40050010","error_msg": "Offset parameter is invalid."}
		offset += len(topics)
		totalCount := utils.PathSearch("total", respBody, float64(0))
		if offset >= int(totalCount.(float64)) {
			break
		}
	}
	return result, respBody, nil
}

func dataSourceConsumerGroupTopicsRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg    = meta.(*config.Config)
		region = cfg.GetRegion(d)
	)
	client, err := cfg.NewServiceClient("dmsv2", region)
	if err != nil {
		return diag.Errorf("error creating DMS client: %s", err)
	}

	topics, respBody, err := getConsumerGroupTopics(client, d)
	if err != nil {
		return diag.Errorf("error getting topics under consumer group: %s", err)
	}

	randomUUID, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(randomUUID)

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("topics", topics),
		d.Set("lag", utils.PathSearch("lag", respBody, nil)),
		d.Set("max_offset", utils.PathSearch("max_offset", respBody, nil)),
		d.Set("consumer_offset", utils.PathSearch("consumer_offset", respBody, nil)),
		d.Set("brokers", flattenConsumerGroupTopicsBrokers(utils.PathSearch("brokers", respBody, make([]interface{}, 0)).([]interface{}))),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenConsumerGroupTopicsBrokers(brokers []interface{}) []interface{} {
	if len(brokers) == 0 {
		return nil
	}

	rst := make([]interface{}, 0, len(brokers))
	for _, v := range brokers {
		rst = append(rst, map[string]interface{}{
			"broker_name": utils.PathSearch("broker_name", v, nil),
			"queues":      flattenConsumerGroupTopicsBrokerQueues(utils.PathSearch("queues", v, make([]interface{}, 0)).([]interface{})),
		})
	}
	return rst
}

func flattenConsumerGroupTopicsBrokerQueues(queues []interface{}) []interface{} {
	if len(queues) == 0 {
		return nil
	}
	rst := make([]interface{}, 0, len(queues))
	for _, v := range queues {
		rst = append(rst, map[string]interface{}{
			"id":                utils.PathSearch("id", v, nil),
			"lag":               utils.PathSearch("lag", v, nil),
			"broker_offset":     utils.PathSearch("broker_offset", v, nil),
			"consumer_offset":   utils.PathSearch("consumer_offset", v, nil),
			"last_message_time": utils.FormatTimeStampRFC3339(int64(utils.PathSearch("last_message_time", v, float64(0)).(float64))/1000, false),
		})
	}
	return rst
}
