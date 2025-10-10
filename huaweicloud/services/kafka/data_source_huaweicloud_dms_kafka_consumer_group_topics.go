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

// @API Kafka GET /v2/{engine}/{project_id}/instances/{instance_id}/groups/{group}/topics
func DataSourceConsumerGroupTopics() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceConsumerGroupTopicsRead,
		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `The region where the Kafka instance and consumer group are located.`,
			},
			"instance_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The ID of the Kafka instance.`,
			},
			"group": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The ID of the consumer group.`,
			},
			"topic": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The name of the topic to be queried.`,
			},
			"sort_key": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The sorting field for the query result.`,
			},
			"sort_dir": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The sorting order for the query result.`,
			},
			"topics": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: `The list of topics that match the filter parameters.`,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"topic": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The name of the topic.`,
						},
						"partitions": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: `The number of partitions.`,
						},
						"lag": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: `The number of message accumulations.`,
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

	if sortKey, ok := d.GetOk("sort_key"); ok {
		res += fmt.Sprintf("&sort_key=%s", sortKey)
	}
	if sortDir, ok := d.GetOk("sort_dir"); ok {
		res += fmt.Sprintf("&sort_dir=%s", sortDir)
	}

	return res
}

func listConsumerGroupTopics(client *golangsdk.ServiceClient, d *schema.ResourceData) ([]interface{}, error) {
	var (
		// The limit maximum value is 50, default is 10.
		httpUrl = "v2/kafka/{project_id}/instances/{instance_id}/groups/{group}/topics?limit=50"
		offset  = 0
		result  = make([]interface{}, 0)
	)

	listPath := client.Endpoint + httpUrl
	listPath = strings.ReplaceAll(listPath, "{project_id}", client.ProjectID)
	listPath = strings.ReplaceAll(listPath, "{instance_id}", d.Get("instance_id").(string))
	listPath = strings.ReplaceAll(listPath, "{group}", d.Get("group").(string))
	listPath += buildConsumerGroupTopicsQueryParams(d)

	getOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json;charset=utf-8",
		},
	}

	for {
		listPathWithOffset := fmt.Sprintf("%s&offset=%d", listPath, offset)
		resp, err := client.Request("GET", listPathWithOffset, &getOpt)
		if err != nil {
			return nil, err
		}

		respBody, err := utils.FlattenResponse(resp)
		if err != nil {
			return nil, err
		}

		topics := utils.PathSearch("topics", respBody, make([]interface{}, 0)).([]interface{})
		result = append(result, topics...)
		// The `offset` cannot be greater than or equal to the total number of topics. Otherwise, the response is as follows:
		// {"error_code": "DMS.00400062","error_msg": "Invalid {0} parameter in the request."}
		offset += len(topics)
		if offset >= int(utils.PathSearch("total", respBody, float64(0)).(float64)) {
			break
		}
	}
	return result, nil
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

	topics, err := listConsumerGroupTopics(client, d)
	if err != nil {
		return diag.Errorf("error querying topic list under consumer group (%s): %s", d.Get("group").(string), err)
	}

	randomUUID, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(randomUUID)

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("topics", flattenConsumerGroupTopics(topics)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenConsumerGroupTopics(topics []interface{}) []interface{} {
	if len(topics) == 0 {
		return nil
	}

	rst := make([]interface{}, 0, len(topics))
	for _, v := range topics {
		rst = append(rst, map[string]interface{}{
			"topic":      utils.PathSearch("topic", v, nil),
			"partitions": utils.PathSearch("partitions", v, nil),
			"lag":        utils.PathSearch("lag", v, nil),
		})
	}
	return rst
}
