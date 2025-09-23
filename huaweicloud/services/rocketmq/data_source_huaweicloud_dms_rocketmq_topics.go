package rocketmq

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk/pagination"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API RocketMQ GET /v2/{project_id}/instances/{instance_id}/topics
func DataSourceDmsRocketMQTopics() *schema.Resource {
	return &schema.Resource{
		ReadContext: resourceDmsRocketMQTopicsRead,
		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"instance_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"total_read_queue_num": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"total_write_queue_num": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"permission": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"topics": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     topicsSchema(),
			},
		},
	}
}

func topicsSchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"total_read_queue_num": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"total_write_queue_num": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"permission": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"brokers": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     brokersSchema(),
			},
		},
	}
	return &sc
}

func brokersSchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"broker_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"read_queue_num": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"write_queue_num": {
				Type:     schema.TypeInt,
				Computed: true,
			},
		},
	}
	return &sc
}

func resourceDmsRocketMQTopicsRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var (
		getRocketmqTopicsHttpUrl = "v2/{project_id}/instances/{instance_id}/topics"
		getRocketmqTopicsProduct = "dmsv2"
	)
	getRocketmqTopicsClient, err := cfg.NewServiceClient(getRocketmqTopicsProduct, region)
	if err != nil {
		return diag.Errorf("error creating DMS Client: %s", err)
	}

	getRocketmqTopicsPath := getRocketmqTopicsClient.Endpoint + getRocketmqTopicsHttpUrl
	getRocketmqTopicsPath = strings.ReplaceAll(getRocketmqTopicsPath, "{project_id}", getRocketmqTopicsClient.ProjectID)
	getRocketmqTopicsPath = strings.ReplaceAll(getRocketmqTopicsPath, "{instance_id}", fmt.Sprintf("%v", d.Get("instance_id")))

	listTopicsResp, err := pagination.ListAllItems(
		getRocketmqTopicsClient,
		"offset",
		getRocketmqTopicsPath,
		&pagination.QueryOpts{MarkerField: ""})
	if err != nil {
		return common.CheckDeletedDiag(d, err, "error retrieving DMS rocketMQ topics")
	}

	listTopicsRespJson, err := json.Marshal(listTopicsResp)
	if err != nil {
		return diag.FromErr(err)
	}
	var listTopicsRespBody interface{}
	err = json.Unmarshal(listTopicsRespJson, &listTopicsRespBody)
	if err != nil {
		return diag.FromErr(err)
	}

	dataSourceId, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(dataSourceId)

	mErr := multierror.Append(
		d.Set("region", region),
		d.Set("topics", flattenListTopicsBody(filterTopics(d, listTopicsRespBody))),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenListTopicsBody(resp []interface{}) []interface{} {
	if resp == nil {
		return nil
	}
	rst := make([]interface{}, 0, len(resp))

	for _, v := range resp {
		rst = append(rst, map[string]interface{}{
			"name":                  utils.PathSearch("name", v, nil),
			"total_read_queue_num":  utils.PathSearch("total_read_queue_num", v, nil),
			"total_write_queue_num": utils.PathSearch("total_write_queue_num", v, nil),
			"permission":            utils.PathSearch("permission", v, nil),
			"brokers":               utils.PathSearch("brokers", v, nil),
		})
	}
	return rst
}

func filterTopics(d *schema.ResourceData, resp interface{}) []interface{} {
	topicJson := utils.PathSearch("topics", resp, make([]interface{}, 0))
	topicArray := topicJson.([]interface{})
	if len(topicArray) < 1 {
		return nil
	}
	result := make([]interface{}, 0, len(topicArray))

	rawName, rawNameOK := d.GetOk("name")
	rawTotalReadQueueNum, rawTotalReadQueueNumOK := d.GetOk("total_read_queue_num")
	rawTotalWriteQueueNum, rawTotalWriteQueueNumOK := d.GetOk("total_write_queue_num")
	rawPermission, rawPermissionOK := d.GetOk("permission")
	for _, topic := range topicArray {
		name := utils.PathSearch("name", topic, nil)
		totalReadQueueNum, _ := utils.PathSearch("total_read_queue_num", topic, float64(0)).(float64)
		totalWriteQueueNum, _ := utils.PathSearch("total_write_queue_num", topic, float64(0)).(float64)
		permission := utils.PathSearch("permission", topic, nil)
		if rawNameOK && rawName != name {
			continue
		}
		if rawTotalReadQueueNumOK && rawTotalReadQueueNum != int(totalReadQueueNum) {
			continue
		}
		if rawTotalWriteQueueNumOK && rawTotalWriteQueueNum != int(totalWriteQueueNum) {
			continue
		}
		if rawPermissionOK && rawPermission != permission {
			continue
		}
		result = append(result, topic)
	}

	return result
}
