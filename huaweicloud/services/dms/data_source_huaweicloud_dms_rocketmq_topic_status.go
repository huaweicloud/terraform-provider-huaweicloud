package dms

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

func DataSourceDmsRocketMQTopicStatus() *schema.Resource {
	return &schema.Resource{
		ReadContext: resourceDmsRocketMQTopicStatusRead,
		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"instance_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `Specifies the ID of the RocketMQ instance.`,
			},
			"topic": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specifies the name of the RocketMQ topic.`,
			},
			"brokers": {
				Type:        schema.TypeList,
				Elem:        DmsRocketMQTopicStatusBrokerSchema(),
				Computed:    true,
				Description: `Indicates the broker list of RocketMQ topic associated with.`,
			},
		},
	}
}

func DmsRocketMQTopicStatusBrokerSchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"queues": {
				Type:        schema.TypeList,
				Elem:        DmsRocketMQTopicStatusBrokerQueueSchema(),
				Computed:    true,
				Description: `Indicates the queue list owned by the RocketMQ topic`,
			},
			"name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Indicates the name of broker.`,
			},
		},
	}
	return &sc
}

func DmsRocketMQTopicStatusBrokerQueueSchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: `Indicates the ID of queue.`,
			},
			"min_offset": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: `Indicates the minimum offset of queue.`,
			},
			"max_offset": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: `Indicates the maximum offset of queue.`,
			},
			"last_message_time": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: `Indicates the time of the last message of queue.`,
			},
		},
	}
	return &sc
}

func resourceDmsRocketMQTopicStatusRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*config.Config)
	region := config.GetRegion(d)

	var mErr *multierror.Error

	// getRocketmqTopicStatus: Query DMS RocketMQ topic status
	var (
		getRocketmqTopicStatusHttpUrl = "v2/{project_id}/instances/{instance_id}/topics/{topic}/status"
		getRocketmqTopicStatusProduct = "dms"
	)
	getRocketmqTopicStatusClient, err := config.NewServiceClient(getRocketmqTopicStatusProduct, region)
	if err != nil {
		return diag.Errorf("error creating DmsRocketMQTopicStatus Client: %s", err)
	}

	getRocketmqTopicStatusPath := getRocketmqTopicStatusClient.Endpoint + getRocketmqTopicStatusHttpUrl
	getRocketmqTopicStatusPath = strings.ReplaceAll(getRocketmqTopicStatusPath, "{project_id}", getRocketmqTopicStatusClient.ProjectID)
	getRocketmqTopicStatusPath = strings.ReplaceAll(getRocketmqTopicStatusPath, "{instance_id}", fmt.Sprintf("%v", d.Get("instance_id")))
	getRocketmqTopicStatusPath = strings.ReplaceAll(getRocketmqTopicStatusPath, "{topic}", fmt.Sprintf("%v", d.Get("topic")))

	getRocketmqTopicStatusOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		OkCodes: []int{
			200,
		},
	}
	getRocketmqTopicStatusResp, err := getRocketmqTopicStatusClient.Request("GET", getRocketmqTopicStatusPath, &getRocketmqTopicStatusOpt)

	if err != nil {
		return common.CheckDeletedDiag(d, err, "error retrieving DmsRocketMQTopicStatus")
	}

	getRocketmqTopicStatusRespBody, err := utils.FlattenResponse(getRocketmqTopicStatusResp)
	if err != nil {
		return diag.FromErr(err)
	}

	uuid, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(uuid)

	mErr = multierror.Append(
		mErr,
		d.Set("region", region),
		d.Set("brokers", flattenGetRocketmqTopicStatusResponseBodyBroker(getRocketmqTopicStatusRespBody)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenGetRocketmqTopicStatusResponseBodyBroker(resp interface{}) []interface{} {
	if resp == nil {
		return nil
	}
	curJson := utils.PathSearch("brokers", resp, make([]interface{}, 0))
	curArray := curJson.([]interface{})
	rst := make([]interface{}, 0, len(curArray))
	for _, v := range curArray {
		rst = append(rst, map[string]interface{}{
			"queues": flattenBrokerQueue(v),
			"name":   utils.PathSearch("broker_name", v, nil),
		})
	}
	return rst
}

func flattenBrokerQueue(resp interface{}) []interface{} {
	if resp == nil {
		return nil
	}
	curJson := utils.PathSearch("queues", resp, make([]interface{}, 0))
	curArray := curJson.([]interface{})
	rst := make([]interface{}, 0, len(curArray))
	for _, v := range curArray {
		rst = append(rst, map[string]interface{}{
			"id":                utils.PathSearch("id", v, nil),
			"min_offset":        utils.PathSearch("min_offset", v, nil),
			"max_offset":        utils.PathSearch("max_offset", v, nil),
			"last_message_time": utils.PathSearch("last_message_time", v, nil),
		})
	}
	return rst
}
