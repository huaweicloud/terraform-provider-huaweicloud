package rds

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/google/uuid"
	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk/pagination"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API RDS GET /v3/{project_id}/schedule-events
func DataSourceEvents() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceEventsRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"event_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"instance_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"status": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"type": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"level": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"sort_field": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"order": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"total_count": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"inquiring_count": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"schedule_count": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"executing_count": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"failed_count": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"events": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     scheduleEventsSchema(),
			},
		},
	}
}

func scheduleEventsSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"instance_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"instance_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"db_type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"created_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"update_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"impact": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"reason": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"level": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"execute_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"latest_execution_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"execution_time_window": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     scheduleExecutionTimeWindowSchema(),
			},
		},
	}
}

func scheduleExecutionTimeWindowSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"planned_execution_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"start_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"end_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceEventsRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var mErr *multierror.Error

	client, err := cfg.NewServiceClient("rds", region)
	if err != nil {
		return diag.Errorf("error creating RDS client: %s", err)
	}

	httpUrl := "v3/{project_id}/schedule-events"
	listPath := client.Endpoint + httpUrl
	listPath = strings.ReplaceAll(listPath, "{project_id}", client.ProjectID)
	listPath += buildGetEventsQueryParams(d)

	listResp, err := pagination.ListAllItems(
		client,
		"offset",
		listPath,
		&pagination.QueryOpts{MarkerField: ""},
	)
	if err != nil {
		return diag.Errorf("error retrieving RDS events: %s", err)
	}

	listRespJson, err := json.Marshal(listResp)
	if err != nil {
		return diag.FromErr(err)
	}

	var listRespBody interface{}
	err = json.Unmarshal(listRespJson, &listRespBody)
	if err != nil {
		return diag.FromErr(err)
	}

	dataSourceId, err := uuid.NewRandom()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(dataSourceId.String())

	mErr = multierror.Append(
		d.Set("region", region),
		d.Set("total_count", utils.PathSearch("total_count", listRespBody, nil)),
		d.Set("inquiring_count", utils.PathSearch("inquiring_count", listRespBody, nil)),
		d.Set("schedule_count", utils.PathSearch("schedule_count", listRespBody, nil)),
		d.Set("executing_count", utils.PathSearch("executing_count", listRespBody, nil)),
		d.Set("failed_count", utils.PathSearch("failed_count", listRespBody, nil)),
		d.Set("events", flattenGetEventsBody(listRespBody)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func buildGetEventsQueryParams(d *schema.ResourceData) string {
	res := ""
	if v, ok := d.GetOk("event_id"); ok {
		res = fmt.Sprintf("%s&id=%v", res, v)
	}
	if v, ok := d.GetOk("instance_id"); ok {
		res = fmt.Sprintf("%s&instance_id=%v", res, v)
	}
	if v, ok := d.GetOk("status"); ok {
		res = fmt.Sprintf("%s&status=%v", res, v)
	}
	if v, ok := d.GetOk("type"); ok {
		res = fmt.Sprintf("%s&type=%v", res, v)
	}
	if v, ok := d.GetOk("level"); ok {
		res = fmt.Sprintf("%s&level=%v", res, v)
	}
	if v, ok := d.GetOk("sort_field"); ok {
		res = fmt.Sprintf("%s&sort_field=%v", res, v)
	}
	if v, ok := d.GetOk("order"); ok {
		res = fmt.Sprintf("%s&order=%v", res, v)
	}

	if res != "" {
		res = "?" + res[1:]
	}
	return res
}

func flattenGetEventsBody(resp interface{}) []interface{} {
	curJson := utils.PathSearch("events", resp, make([]interface{}, 0))
	curArray := curJson.([]interface{})
	res := make([]interface{}, 0, len(curArray))
	for _, v := range curArray {
		res = append(res, map[string]interface{}{
			"id":                    utils.PathSearch("id", v, nil),
			"instance_id":           utils.PathSearch("instance_id", v, nil),
			"instance_name":         utils.PathSearch("instance_name", v, nil),
			"db_type":               utils.PathSearch("db_type", v, nil),
			"created_time":          utils.PathSearch("created_time", v, nil),
			"update_time":           utils.PathSearch("update_time", v, nil),
			"type":                  utils.PathSearch("type", v, nil),
			"impact":                utils.PathSearch("impact", v, nil),
			"status":                utils.PathSearch("status", v, nil),
			"reason":                utils.PathSearch("reason", v, nil),
			"level":                 utils.PathSearch("level", v, nil),
			"execute_time":          utils.PathSearch("execute_time", v, nil),
			"latest_execution_time": utils.PathSearch("latest_execution_time", v, nil),
			"execution_time_window": flattenEventsExecutionTimeWindow(v),
		})
	}
	return res
}

func flattenEventsExecutionTimeWindow(resp interface{}) []interface{} {
	curJson := utils.PathSearch("execution_time_window", resp, nil)
	if curJson == nil {
		return nil
	}
	res := []interface{}{
		map[string]interface{}{
			"planned_execution_time": utils.PathSearch("planned_execution_time", curJson, nil),
			"start_time":             utils.PathSearch("start_time", curJson, nil),
			"end_time":               utils.PathSearch("end_time", curJson, nil),
		},
	}
	return res
}
