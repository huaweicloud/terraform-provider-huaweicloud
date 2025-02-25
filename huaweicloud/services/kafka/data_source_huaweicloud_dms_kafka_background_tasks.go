package kafka

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API Kafka GET /v2/{project_id}/instances/{instance_id}/tasks
func DataSourceDmsKafkaBackgroundTasks() *schema.Resource {
	return &schema.Resource{
		ReadContext: DataSourceDmsKafkaBackgroundTasksRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `Specifies the region in which to query the resource. If omitted, the provider-level region will be used.`,
			},
			"instance_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `Specifies the instance ID.`,
			},
			"begin_time": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specifies the time of task where the query starts. The format is YYYYMMDDHHmmss.`,
			},
			"end_time": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specifies the time of task where the query ends. The format is YYYYMMDDHHmmss.`,
			},
			"tasks": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: `Indicates the task list.`,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `Indicates the task ID.`,
						},
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `Indicates the task name.`,
						},
						"params": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `Indicates the task parameters.`,
						},
						"status": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `Indicates the task status.`,
						},
						"user_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `Indicates the user ID.`,
						},
						"user_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `Indicates the username.`,
						},
						"created_at": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `Indicates the start time.`,
						},
						"updated_at": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `Indicates the end time.`,
						},
					},
				},
			},
		},
	}
}

func DataSourceDmsKafkaBackgroundTasksRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	client, err := cfg.NewServiceClient("dmsv2", region)
	if err != nil {
		return diag.Errorf("error creating DMS client: %s", err)
	}

	listHttpUrl := "v2/{project_id}/instances/{instance_id}/tasks"
	listPath := client.Endpoint + listHttpUrl
	listPath = strings.ReplaceAll(listPath, "{project_id}", client.ProjectID)
	listPath = strings.ReplaceAll(listPath, "{instance_id}", d.Get("instance_id").(string))
	listOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"Content-Type": "application/json"},
	}

	// pagelimit is `10`
	listPath += fmt.Sprintf("?limit=%v", pageLimit)
	listPath = buildQueryBackgroundTasksListPath(d, listPath)

	// `start` counts from `1`
	start := 1
	results := make([]map[string]interface{}, 0)
	for {
		currentPath := listPath + fmt.Sprintf("&start=%d", start)
		listResp, err := client.Request("GET", currentPath, &listOpt)
		if err != nil {
			return diag.Errorf("error retrieving tasks: %s", err)
		}
		listRespBody, err := utils.FlattenResponse(listResp)
		if err != nil {
			return diag.Errorf("error flatten response: %s", err)
		}

		tasks := utils.PathSearch("tasks", listRespBody, make([]interface{}, 0)).([]interface{})
		for _, task := range tasks {
			results = append(results, map[string]interface{}{
				"id":         utils.PathSearch("id", task, nil),
				"name":       utils.PathSearch("name", task, nil),
				"params":     utils.PathSearch("params", task, nil),
				"status":     utils.PathSearch("status", task, nil),
				"user_id":    utils.PathSearch("user_id", task, nil),
				"user_name":  utils.PathSearch("user_name", task, nil),
				"created_at": utils.PathSearch("created_at", task, nil),
				"updated_at": utils.PathSearch("updated_at", task, nil),
			})
		}

		// `task_count` means the number of all `tasks`, and type is string.
		start += len(tasks)
		taskCount := utils.PathSearch("task_count", listRespBody, "0").(string)
		totalCount, _ := strconv.Atoi(taskCount)
		if totalCount <= start-1 {
			break
		}
	}

	id, err := uuid.GenerateUUID()
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(id)

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("tasks", results),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func buildQueryBackgroundTasksListPath(d *schema.ResourceData, listPath string) string {
	if beginTime, ok := d.GetOk("begin_time"); ok {
		listPath += fmt.Sprintf("&begin_time=%s", beginTime)
	}
	if endTime, ok := d.GetOk("end_time"); ok {
		listPath += fmt.Sprintf("&end_time=%s", endTime)
	}

	return listPath
}
