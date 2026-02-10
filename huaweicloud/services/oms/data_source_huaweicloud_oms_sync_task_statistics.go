package oms

import (
	"context"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API OMS GET /v2/{project_id}/sync-tasks/{sync_task_id}/statistics
func DataSourceSyncTaskStatistics() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceSyncTaskStatisticsRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"sync_task_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"data_type": {
				Type:     schema.TypeString,
				Required: true,
			},
			"start_time": {
				Type:     schema.TypeString,
				Required: true,
			},
			"end_time": {
				Type:     schema.TypeString,
				Required: true,
			},
			"task_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"statistic_time_type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"statistic_datas": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"data_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"data": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"time_stamp": {
										Type:     schema.TypeInt,
										Computed: true,
									},
									"statistic_num": {
										Type:     schema.TypeInt,
										Computed: true,
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

func listSyncTaskStatistics(client *golangsdk.ServiceClient, d *schema.ResourceData) (interface{}, error) {
	var (
		httpUrl   = "v2/{project_id}/sync-tasks/{sync_task_id}/statistics?data_type={data_type}&start_time={start_time}&end_time={end_time}"
		taskId    = d.Get("sync_task_id").(string)
		dataType  = d.Get("data_type").(string)
		startTime = d.Get("start_time").(string)
		endTime   = d.Get("end_time").(string)
	)

	getPath := client.Endpoint + httpUrl
	getPath = strings.ReplaceAll(getPath, "{project_id}", client.ProjectID)
	getPath = strings.ReplaceAll(getPath, "{sync_task_id}", taskId)
	getPath = strings.ReplaceAll(getPath, "{data_type}", dataType)
	getPath = strings.ReplaceAll(getPath, "{start_time}", startTime)
	getPath = strings.ReplaceAll(getPath, "{end_time}", endTime)
	getOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json",
		},
	}

	getResp, err := client.Request("GET", getPath, &getOpt)
	if err != nil {
		return nil, err
	}

	getRespBody, err := utils.FlattenResponse(getResp)
	if err != nil {
		return nil, err
	}

	return getRespBody, nil
}

func dataSourceSyncTaskStatisticsRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg    = meta.(*config.Config)
		region = cfg.GetRegion(d)
	)

	client, err := cfg.NewServiceClient("oms", region)
	if err != nil {
		return diag.Errorf("error creating OMS client: %s", err)
	}

	resp, err := listSyncTaskStatistics(client, d)
	if err != nil {
		return diag.Errorf("error retrieving synchronization task statistics: %s", err)
	}

	generateUUID, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}

	d.SetId(generateUUID)

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("task_id", utils.PathSearch("sync_task_id", resp, nil)),
		d.Set("statistic_time_type", utils.PathSearch("statistic_time_type", resp, nil)),
		d.Set("statistic_datas", flattenSyncTaskStatistics(
			utils.PathSearch("statistic_datas", resp, make([]interface{}, 0)).([]interface{}))),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenSyncTaskStatistics(resp []interface{}) []interface{} {
	if len(resp) == 0 {
		return nil
	}

	result := make([]interface{}, 0, len(resp))
	for _, v := range resp {
		result = append(result, map[string]interface{}{
			"data_type": utils.PathSearch("data_type", v, nil),
			"data":      flattenSyncTaskStatisticsData(utils.PathSearch("data", v, make([]interface{}, 0)).([]interface{})),
		})
	}

	return result
}

func flattenSyncTaskStatisticsData(resp []interface{}) []interface{} {
	if len(resp) == 0 {
		return nil
	}

	result := make([]interface{}, 0, len(resp))
	for _, v := range resp {
		result = append(result, map[string]interface{}{
			"time_stamp":    utils.PathSearch("time_stamp", v, nil),
			"statistic_num": utils.PathSearch("statistic_num", v, nil),
		})
	}

	return result
}
