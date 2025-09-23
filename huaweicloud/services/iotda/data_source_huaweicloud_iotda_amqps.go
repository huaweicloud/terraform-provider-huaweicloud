package iotda

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

func buildAMQPQueuesQueryParams(d *schema.ResourceData) string {
	queryParams := "?limit=50"
	if v, ok := d.GetOk("name"); ok {
		return fmt.Sprintf("%s&queue_name=%v", queryParams, v)
	}

	return queryParams
}

func dataSourceAMQPQueuesRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg       = meta.(*config.Config)
		region    = cfg.GetRegion(d)
		isDerived = WithDerivedAuth(cfg, region)
		product   = "iotda"
		httpUrl   = "v5/iot/{project_id}/amqp-queues"
		offset    = 0
		result    = make([]interface{}, 0)
	)

	client, err := cfg.NewServiceClientWithDerivedAuth(product, region, isDerived)
	if err != nil {
		return diag.Errorf("error creating IoTDA client: %s", err)
	}

	getPath := client.Endpoint + httpUrl
	getPath = strings.ReplaceAll(getPath, "{project_id}", client.ProjectID)
	getPath += buildAMQPQueuesQueryParams(d)
	getOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	for {
		currentPath := fmt.Sprintf("%s&offset=%v", getPath, offset)
		getResp, err := client.Request("GET", currentPath, &getOpt)
		if err != nil {
			return diag.Errorf("error retrieving IoTDA AMQP queues: %s", err)
		}

		getRespBody, err := utils.FlattenResponse(getResp)
		if err != nil {
			return diag.FromErr(err)
		}

		queuesResp := utils.PathSearch("queues", getRespBody, make([]interface{}, 0)).([]interface{})
		if len(queuesResp) == 0 {
			break
		}

		result = append(result, queuesResp...)
		offset += len(queuesResp)
	}

	generateUUID, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}

	d.SetId(generateUUID)

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("queues", flattenAMQPQueues(filterAMQPQueues(result, d))),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func filterAMQPQueues(queuesResp []interface{}, d *schema.ResourceData) []interface{} {
	rst := make([]interface{}, 0, len(queuesResp))
	for _, v := range queuesResp {
		if queueId, ok := d.GetOk("queue_id"); ok &&
			fmt.Sprint(queueId) != utils.PathSearch("queue_id", v, "").(string) {
			continue
		}

		rst = append(rst, v)
	}

	return rst
}

func flattenAMQPQueues(queuesResp []interface{}) []interface{} {
	rst := make([]interface{}, 0, len(queuesResp))
	for _, v := range queuesResp {
		rst = append(rst, map[string]interface{}{
			"id":         utils.PathSearch("queue_id", v, nil),
			"name":       utils.PathSearch("queue_name", v, nil),
			"created_at": utils.PathSearch("create_time", v, nil),
			"updated_at": utils.PathSearch("last_modify_time", v, nil),
		})
	}

	return rst
}
