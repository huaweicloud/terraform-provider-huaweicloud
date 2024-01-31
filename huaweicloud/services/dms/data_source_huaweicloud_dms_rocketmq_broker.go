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

// @API RocketMQ GET /v2/{project_id}/instances/{instance_id}/brokers
func DataSourceDmsRocketMQBroker() *schema.Resource {
	return &schema.Resource{
		ReadContext: resourceDmsRocketMQBrokerRead,
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
			"brokers": {
				Type:        schema.TypeList,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Computed:    true,
				Description: `Indicates the list of the brokers.`,
			},
		},
	}
}

func resourceDmsRocketMQBrokerRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var mErr *multierror.Error

	// getRocketmqBroker: Query the List of rocketMQ broker
	var (
		getRocketmqBrokerHttpUrl = "v2/{project_id}/instances/{instance_id}/brokers"
		getRocketmqBrokerProduct = "dmsv2"
	)
	getRocketmqBrokerClient, err := cfg.NewServiceClient(getRocketmqBrokerProduct, region)
	if err != nil {
		return diag.Errorf("error creating DmsRocketMQBroker Client: %s", err)
	}

	getRocketmqBrokerPath := getRocketmqBrokerClient.Endpoint + getRocketmqBrokerHttpUrl
	getRocketmqBrokerPath = strings.ReplaceAll(getRocketmqBrokerPath, "{project_id}", getRocketmqBrokerClient.ProjectID)
	getRocketmqBrokerPath = strings.ReplaceAll(getRocketmqBrokerPath, "{instance_id}", fmt.Sprintf("%v", d.Get("instance_id")))

	getRocketmqBrokerOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		OkCodes: []int{
			200,
		},
	}
	getRocketmqBrokerResp, err := getRocketmqBrokerClient.Request("GET", getRocketmqBrokerPath, &getRocketmqBrokerOpt)

	if err != nil {
		return common.CheckDeletedDiag(d, err, "error retrieving DmsRocketMQBroker")
	}

	getRocketmqBrokerRespBody, err := utils.FlattenResponse(getRocketmqBrokerResp)
	if err != nil {
		return diag.FromErr(err)
	}

	dataSourceId, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(dataSourceId)

	mErr = multierror.Append(
		mErr,
		d.Set("region", region),
		d.Set("brokers", flattenGetRocketmqBrokerResponseBodyBroker(getRocketmqBrokerRespBody)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenGetRocketmqBrokerResponseBodyBroker(resp interface{}) []interface{} {
	if resp == nil {
		return nil
	}
	curJson := utils.PathSearch("brokers", resp, make([]interface{}, 0))
	curArray := curJson.([]interface{})
	rst := make([]interface{}, 0, len(curArray))
	for _, v := range curArray {
		rst = append(rst, utils.PathSearch("broker_name", v, nil))
	}
	return rst
}
