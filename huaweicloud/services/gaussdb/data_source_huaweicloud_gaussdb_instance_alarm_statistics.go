package gaussdb

import (
	"context"
	"fmt"
	"strings"

	"github.com/google/uuid"
	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API GaussDB GET /v3/{project_id}/instances/alarm-statistics
func DataSourceInstanceAlarmStatistics() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceInstanceAlarmStatisticsRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"start_time": {
				Type:     schema.TypeString,
				Required: true,
			},
			"top_num": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"ring_percentage": {
				Type:     schema.TypeFloat,
				Computed: true,
			},
			"instance_alarm_level_statistics": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     gaussDbInstanceAlarmLevelStatisticsSchema(),
			},
			"total_alarm_level_statistics": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     gaussDbTotalAlarmLevelStatisticsSchema(),
			},
		},
	}
}

func gaussDbInstanceAlarmLevelStatisticsSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"instance_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"instance_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"total_count": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"alarm_level_statistics": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     gaussDbTotalAlarmLevelStatisticsSchema(),
			},
		},
	}
}

func gaussDbTotalAlarmLevelStatisticsSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"count": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"level_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func queryInstanceAlarmStatistics(client *golangsdk.ServiceClient, d *schema.ResourceData) (interface{}, error) {
	httpUrl := "v3/{project_id}/instances/alarm-statistics"

	listPath := client.Endpoint + httpUrl
	listPath = strings.ReplaceAll(listPath, "{project_id}", client.ProjectID)

	startTime := strings.ReplaceAll(d.Get("start_time").(string), "+", "%2B")
	listPath += fmt.Sprintf("?start_time=%s&top_num=%d", startTime, d.Get("top_num").(int))

	opt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"Content-Type": "application/json"},
	}

	requestResp, err := client.Request("GET", listPath, &opt)
	if err != nil {
		return nil, err
	}

	return utils.FlattenResponse(requestResp)
}

func dataSourceInstanceAlarmStatisticsRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg    = meta.(*config.Config)
		region = cfg.GetRegion(d)
	)

	client, err := cfg.NewServiceClient("opengauss", region)
	if err != nil {
		return diag.Errorf("error creating GaussDB client: %s", err)
	}

	resp, err := queryInstanceAlarmStatistics(client, d)
	if err != nil {
		return diag.Errorf("error querying GaussDB instance alarm statistics: %s", err)
	}

	dataSourceId, err := uuid.NewRandom()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(dataSourceId.String())

	mErr := multierror.Append(
		d.Set("region", region),
		d.Set("ring_percentage", utils.PathSearch("ring_percentage", resp, nil)),
		d.Set("instance_alarm_level_statistics", flattenInstanceAlarmLevelStatistics(resp)),
		d.Set("total_alarm_level_statistics", flattenTotalAlarmLevelStatistics(resp)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenInstanceAlarmLevelStatistics(resp interface{}) []interface{} {
	if resp == nil {
		return nil
	}

	curJson := utils.PathSearch("instance_alarm_level_statistics", resp, make([]interface{}, 0))
	curArray := curJson.([]interface{})
	res := make([]interface{}, 0, len(curArray))

	for _, v := range curArray {
		res = append(
			res, map[string]interface{}{
				"instance_id":            utils.PathSearch("instance_id", v, nil),
				"instance_name":          utils.PathSearch("instance_name", v, nil),
				"total_count":            utils.PathSearch("total_count", v, nil),
				"alarm_level_statistics": flattenTotalAlarmLevelStatistics(v),
			},
		)
	}
	return res
}

func flattenTotalAlarmLevelStatistics(resp interface{}) []interface{} {
	if resp == nil {
		return nil
	}

	curJson := utils.PathSearch("alarm_level_statistics", resp, make([]interface{}, 0))
	curArray := curJson.([]interface{})
	res := make([]interface{}, 0, len(curArray))

	for _, v := range curArray {
		res = append(
			res, map[string]interface{}{
				"count":      utils.PathSearch("count", v, nil),
				"level_name": utils.PathSearch("level_name", v, nil),
			},
		)
	}
	return res
}
