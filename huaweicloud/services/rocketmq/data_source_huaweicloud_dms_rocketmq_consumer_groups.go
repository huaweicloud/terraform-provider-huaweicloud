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

// @API RocketMQ GET /v2/{project_id}/instances/{instance_id}/groups
func DataSourceDmsRocketMQConsumerGroups() *schema.Resource {
	return &schema.Resource{
		ReadContext: resourceDmsRocketMQConsumerGroupsRead,
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
			"enabled": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
			},
			"broadcast": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"retry_max_times": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"groups": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     groupsSchema(),
			},
		},
	}
}

func groupsSchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"enabled": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"broadcast": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"brokers": {
				Type:     schema.TypeList,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Computed: true,
			},
			"name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"description": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"retry_max_times": {
				Type:     schema.TypeInt,
				Computed: true,
			},
		},
	}
	return &sc
}

func resourceDmsRocketMQConsumerGroupsRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var (
		getRocketmqConsumerGroupsHttpUrl = "v2/{project_id}/instances/{instance_id}/groups"
		getRocketmqConsumerGroupsProduct = "dmsv2"
	)
	getRocketmqConsumerGroupsClient, err := cfg.NewServiceClient(getRocketmqConsumerGroupsProduct, region)
	if err != nil {
		return diag.Errorf("error creating DMS Client: %s", err)
	}

	getRocketmqConsumerGroupsPath := getRocketmqConsumerGroupsClient.Endpoint + getRocketmqConsumerGroupsHttpUrl
	getRocketmqConsumerGroupsPath = strings.ReplaceAll(getRocketmqConsumerGroupsPath, "{project_id}", getRocketmqConsumerGroupsClient.ProjectID)
	getRocketmqConsumerGroupsPath = strings.ReplaceAll(getRocketmqConsumerGroupsPath, "{instance_id}", fmt.Sprintf("%v", d.Get("instance_id")))

	getConsumerGroupsResp, err := pagination.ListAllItems(
		getRocketmqConsumerGroupsClient,
		"offset",
		getRocketmqConsumerGroupsPath,
		&pagination.QueryOpts{MarkerField: ""})
	if err != nil {
		return common.CheckDeletedDiag(d, err, "error retrieving DMS rocketMQ consumer groups")
	}

	getConsumerGroupsRespJson, err := json.Marshal(getConsumerGroupsResp)
	if err != nil {
		return diag.FromErr(err)
	}
	var getConsumerGroupsRespBody interface{}
	err = json.Unmarshal(getConsumerGroupsRespJson, &getConsumerGroupsRespBody)
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
		d.Set("groups", flattenGetConsumerGroupsBody(filterConsumerGroups(d, getConsumerGroupsRespBody))),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenGetConsumerGroupsBody(resp []interface{}) []interface{} {
	if resp == nil {
		return nil
	}
	rst := make([]interface{}, 0, len(resp))

	for _, v := range resp {
		rst = append(rst, map[string]interface{}{
			"enabled":         utils.PathSearch("enabled", v, nil),
			"broadcast":       utils.PathSearch("broadcast", v, nil),
			"brokers":         utils.PathSearch("brokers", v, nil),
			"name":            utils.PathSearch("name", v, nil),
			"description":     utils.PathSearch("group_desc", v, nil),
			"retry_max_times": utils.PathSearch("retry_max_time", v, nil),
		})
	}
	return rst
}

func filterConsumerGroups(d *schema.ResourceData, resp interface{}) []interface{} {
	groupJson := utils.PathSearch("groups", resp, make([]interface{}, 0))
	groupArray := groupJson.([]interface{})
	if len(groupArray) < 1 {
		return nil
	}
	result := make([]interface{}, 0, len(groupArray))

	rawEnabled := d.Get("enabled").(bool)
	rawBroadcast := d.Get("broadcast").(bool)
	rawName, rawNameOK := d.GetOk("name")
	rawRetryMaxTimes, rawRetryMaxTimesOK := d.GetOk("retry_max_times")

	for _, group := range groupArray {
		enabled := utils.PathSearch("enabled", group, false).(bool)
		broadcast := utils.PathSearch("broadcast", group, false).(bool)
		name := utils.PathSearch("name", group, nil)
		retryMaxTimes := utils.PathSearch("retry_max_time", group, float64(0)).(float64)
		if (rawBroadcast && !broadcast) || (!rawBroadcast && broadcast) {
			continue
		}
		if (rawEnabled && !enabled) || (!rawEnabled && enabled) {
			continue
		}
		if rawNameOK && rawName != name {
			continue
		}
		if rawRetryMaxTimesOK && rawRetryMaxTimes != int(retryMaxTimes) {
			continue
		}
		result = append(result, group)
	}

	return result
}
