package dds

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

// @API DDS GET /v3.1/{project_id}/jobs
func DataSourceDdsInstantTasks() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceDdsInstantTasksRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `Specifies the region in which to query the resource. If omitted, the provider-level region will be used.`,
			},
			"start_time": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `Specifies the start time. The format of the start time is **yyyy-mm-ddThh:mm:ssZ**.`,
			},
			"end_time": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `Specifies the end time. The format of the end time is **yyyy-mm-ddThh:mm:ssZ**`,
			},
			"status": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specifies the task status.`,
			},
			"name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specifies the task name. The value can be:`,
			},
			"jobs": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: `Indicates the tasks list.`,
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
						"fail_reason": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `Indicates the task failure information.`,
						},
						"instance_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `Indicates the instance ID.`,
						},
						"instance_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `Indicates the instance name.`,
						},
						"progress": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `Indicates the task execution progress.`,
						},
						"status": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `Indicates the task status.`,
						},
						"created_at": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `Indicates the task creation time.`,
						},
						"ended_at": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `Indicates the task end time.`,
						},
					},
				},
			},
		},
	}
}

func dataSourceDdsInstantTasksRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	client, err := cfg.NewServiceClient("dds", region)
	if err != nil {
		return diag.Errorf("error creating DDS client: %s", err)
	}

	getTasksHttpUrl := "v3.1/{project_id}/jobs"
	getTasksPath := client.Endpoint + getTasksHttpUrl
	getTasksPath = strings.ReplaceAll(getTasksPath, "{project_id}", client.ProjectID)
	getTasksOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	// pagelimit is `10`
	getTasksPath += fmt.Sprintf("?limit=%v", pageLimit)
	getTasksPath = buildQueryTasksListPath(d, getTasksPath)

	currentTotal := 0
	results := make([]map[string]interface{}, 0)
	for {
		currentPath := getTasksPath + fmt.Sprintf("&offset=%d", currentTotal)
		getTasksResp, err := client.Request("GET", currentPath, &getTasksOpt)
		if err != nil {
			return diag.Errorf("error retrieving tasks: %s", err)
		}
		getTasksRespBody, err := utils.FlattenResponse(getTasksResp)
		if err != nil {
			return diag.Errorf("error flatten response: %s", err)
		}

		tasks := utils.PathSearch("jobs", getTasksRespBody, make([]interface{}, 0)).([]interface{})
		for _, task := range tasks {
			results = append(results, map[string]interface{}{
				"id":            utils.PathSearch("id", task, nil),
				"name":          utils.PathSearch("name", task, nil),
				"instance_id":   utils.PathSearch("instance_id", task, nil),
				"instance_name": utils.PathSearch("instance_name", task, nil),
				"fail_reason":   utils.PathSearch("fail_reason", task, nil),
				"progress":      utils.PathSearch("progress", task, nil),
				"created_at":    utils.PathSearch("created_at", task, nil),
				"status":        utils.PathSearch("status", task, nil),
				"ended_at":      utils.PathSearch("ended_at", task, nil),
			})
		}

		// `total_count` is actual not in return
		if len(tasks) < pageLimit {
			break
		}

		currentTotal += len(tasks)
	}

	id, err := uuid.GenerateUUID()
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(id)

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("jobs", results),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func buildQueryTasksListPath(d *schema.ResourceData, getTasksPath string) string {
	getTasksPath += fmt.Sprintf("&start_time=%v", d.Get("start_time"))
	getTasksPath += fmt.Sprintf("&end_time=%v", d.Get("end_time"))

	if status, ok := d.GetOk("status"); ok {
		getTasksPath += fmt.Sprintf("&status=%s", status)
	}
	if name, ok := d.GetOk("name"); ok {
		getTasksPath += fmt.Sprintf("&name=%s", name)
	}

	return getTasksPath
}
