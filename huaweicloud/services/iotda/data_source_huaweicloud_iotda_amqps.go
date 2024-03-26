package iotda

import (
	"context"
	"fmt"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/huaweicloud/huaweicloud-sdk-go-v3/services/iotda/v5/model"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API IoTDA GET /v5/iot/{project_id}/amqp-queues
func DataSourceAMQPQueues() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceAMQPQueuesRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"queue_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"queues": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"created_at": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"updated_at": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceAMQPQueuesRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	isDerived := WithDerivedAuth(cfg, region)
	client, err := cfg.HcIoTdaV5Client(region, isDerived)
	if err != nil {
		return diag.Errorf("error creating IoTDA v5 client: %s", err)
	}

	var (
		allQueues []model.QueryQueueBase
		limit     = int32(50)
		offset    int32
	)

	for {
		listOpts := model.BatchShowQueueRequest{
			QueueName: utils.StringIgnoreEmpty(d.Get("name").(string)),
			Limit:     utils.Int32(limit),
			Offset:    &offset,
		}

		listResp, listErr := client.BatchShowQueue(&listOpts)
		if listErr != nil {
			return diag.Errorf("error querying IoTDA AMQP queues: %s", listErr)
		}

		if len(*listResp.Queues) == 0 {
			break
		}

		allQueues = append(allQueues, *listResp.Queues...)
		offset += limit
	}

	uuId, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(uuId)

	targetQueues := filterListQueues(allQueues, d)
	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("queues", flattenAMQPQueues(targetQueues)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func filterListQueues(queues []model.QueryQueueBase, d *schema.ResourceData) []model.QueryQueueBase {
	if len(queues) == 0 {
		return nil
	}

	rst := make([]model.QueryQueueBase, 0, len(queues))
	for _, v := range queues {
		if queueId, ok := d.GetOk("queue_id"); ok &&
			fmt.Sprint(queueId) != utils.StringValue(v.QueueId) {
			continue
		}

		rst = append(rst, v)
	}

	return rst
}

func flattenAMQPQueues(queues []model.QueryQueueBase) []interface{} {
	if len(queues) == 0 {
		return nil
	}

	rst := make([]interface{}, 0, len(queues))
	for _, v := range queues {
		rst = append(rst, map[string]interface{}{
			"id":         v.QueueId,
			"name":       v.QueueName,
			"created_at": v.CreateTime,
			"updated_at": v.LastModifyTime,
		})
	}

	return rst
}
