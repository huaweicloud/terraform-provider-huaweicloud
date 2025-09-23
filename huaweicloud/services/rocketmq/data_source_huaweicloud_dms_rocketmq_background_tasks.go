package rocketmq

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

// @API RocketMQ GET /v2/{project_id}/instances/{instance_id}/tasks
func DataSourceBackgroundTasks() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceBackgroundTasksRead,
		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "The region where the background tasks are located.",
			},
			"instance_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The ID of the RocketMQ instance.",
			},
			"begin_time": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The start time of the background task, in UTC format.",
			},
			"end_time": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The end time of the background task, in UTC format.",
			},
			"tasks": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The ID of the background task.",
						},
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The name of the background task.",
						},
						"user_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The username of the user who executed the background task.",
						},
						"user_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The ID of the user who executed the background task.",
						},
						"params": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The parameters of the background task.",
						},
						"status": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The status of the background task.",
						},
						"created_at": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The creation time of the background task, in UTC format.",
						},
						"updated_at": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The latest update time of the background task, in UTC format.",
						},
					},
				},
				Description: "The list of the background tasks under the specified RocketMQ instance.",
			},
		},
	}
}

func buildQueryBackgroundTasksListPath(d *schema.ResourceData) string {
	res := ""

	if v, ok := d.GetOk("begin_time"); ok {
		res = fmt.Sprintf("%s&begin_time=%s", res, v)
	}

	if v, ok := d.GetOk("end_time"); ok {
		res = fmt.Sprintf("%s&end_time=%s", res, v)
	}
	return res
}

func queryBackgroundTasks(client *golangsdk.ServiceClient, d *schema.ResourceData, instanceId string) ([]interface{}, error) {
	var (
		listPath = "v2/{project_id}/instances/{instance_id}/tasks"
		// The `start` starts from `1`.
		start = 1
		// The `limit` default value is `10`.
		limit   = 100
		results = make([]interface{}, 0)
		listOpt = golangsdk.RequestOpts{
			KeepResponseBody: true,
		}
	)

	listPath = client.Endpoint + listPath
	listPath = strings.ReplaceAll(listPath, "{project_id}", client.ProjectID)
	listPath = strings.ReplaceAll(listPath, "{instance_id}", instanceId)
	listPath = fmt.Sprintf("%s?limit=%d", listPath, limit)
	listPath += buildQueryBackgroundTasksListPath(d)

	for {
		listPathWithStart := fmt.Sprintf("%s&start=%d", listPath, start)
		listResp, err := client.Request("GET", listPathWithStart, &listOpt)
		if err != nil {
			return nil, err
		}

		respBody, err := utils.FlattenResponse(listResp)
		if err != nil {
			return nil, err
		}

		backgroundTasks := utils.PathSearch("tasks", respBody, make([]interface{}, 0)).([]interface{})
		results = append(results, backgroundTasks...)
		if len(backgroundTasks) < limit {
			break
		}

		start += len(backgroundTasks)
	}

	return results, nil
}

func dataSourceBackgroundTasksRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg        = meta.(*config.Config)
		region     = cfg.GetRegion(d)
		instanceId = d.Get("instance_id").(string)
	)

	client, err := cfg.NewServiceClient("dmsv2", region)
	if err != nil {
		return diag.Errorf("error creating DMS client: %s", err)
	}

	tasks, err := queryBackgroundTasks(client, d, instanceId)
	if err != nil {
		return diag.Errorf("error getting background tasks under the specified RocketMQ instance (%s): %s",
			instanceId, err)
	}

	randomUUID, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate data source ID: %s", err)
	}
	d.SetId(randomUUID)

	mErr := multierror.Append(
		d.Set("region", region),
		d.Set("tasks", flattenBackgroundTasks(tasks)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenBackgroundTasks(tasks []interface{}) []map[string]interface{} {
	result := make([]map[string]interface{}, len(tasks))
	for i, v := range tasks {
		result[i] = map[string]interface{}{
			"id":         utils.PathSearch("id", v, nil),
			"name":       utils.PathSearch("name", v, nil),
			"user_name":  utils.PathSearch("user_name", v, nil),
			"user_id":    utils.PathSearch("user_id", v, nil),
			"params":     utils.PathSearch("params", v, nil),
			"status":     utils.PathSearch("status", v, nil),
			"created_at": utils.PathSearch("created_at", v, nil),
			"updated_at": utils.PathSearch("updated_at", v, nil),
		}
	}

	return result
}
