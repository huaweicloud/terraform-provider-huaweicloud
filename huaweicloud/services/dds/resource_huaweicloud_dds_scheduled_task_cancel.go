package dds

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API DDS DELETE /v3/{project_id}/scheduled-jobs/{job_id}
// @API DDS GET /v3/{project_id}/scheduled-jobs
func ResourceDDSScheduledTaskCancel() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceDDSMysqlScheduledTaskCancelCreate,
		ReadContext:   resourceDDSMysqlScheduledTaskCancelRead,
		DeleteContext: resourceDDSMysqlScheduledTaskCancelDelete,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"job_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: `Specifies the task ID.`,
			},
			"job_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Indicates the task name.`,
			},
			"create_time": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Indicates the create time.`,
			},
			"start_time": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Indicates the start time.`,
			},
			"end_time": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Indicates the end time.`,
			},
			"job_status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Indicates the task execution status.`,
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
			"instance_status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Indicates the instance status.`,
			},
		},
	}
}

func resourceDDSMysqlScheduledTaskCancelCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	client, err := cfg.DdsV3Client(region)
	if err != nil {
		return diag.Errorf("error creating DDS client: %s", err)
	}

	jobId := d.Get("job_id").(string)
	httpUrl := "v3/{project_id}/scheduled-jobs/{job_id}"
	createPath := client.Endpoint + httpUrl
	createPath = strings.ReplaceAll(createPath, "{project_id}", client.ProjectID)
	createPath = strings.ReplaceAll(createPath, "{job_id}", jobId)
	createOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"Content-Type": "application/json"},
	}

	_, err = client.Request("DELETE", createPath, &createOpt)
	if err != nil {
		return diag.Errorf("error canceling DDS scheduled task(%s): %s", jobId, err)
	}

	d.SetId(jobId)

	return resourceDDSMysqlScheduledTaskCancelRead(ctx, d, meta)
}

func resourceDDSMysqlScheduledTaskCancelRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	client, err := cfg.DdsV3Client(region)
	if err != nil {
		return diag.Errorf("error creating DDS client: %s", err)
	}

	getTasksHttpUrl := "v3/{project_id}/scheduled-jobs"
	getTasksPath := client.Endpoint + getTasksHttpUrl
	getTasksPath = strings.ReplaceAll(getTasksPath, "{project_id}", client.ProjectID)
	getTasksOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"Content-Type": "application/json"},
	}

	// pagelimit is `10`
	getTasksPath += fmt.Sprintf("?limit=%v", pageLimit)
	currentTotal := 0

	var scheduledTask interface{}
	for {
		currentPath := getTasksPath + fmt.Sprintf("&offset=%d", currentTotal)
		getTasksResp, err := client.Request("GET", currentPath, &getTasksOpt)
		if err != nil {
			return diag.Errorf("error retrieving scheduled tasks: %s", err)
		}
		getTasksRespBody, err := utils.FlattenResponse(getTasksResp)
		if err != nil {
			return diag.Errorf("error flatten response: %s", err)
		}

		searchPath := fmt.Sprintf("schedules[?job_id=='%s']|[0]", d.Id())
		scheduledTask = utils.PathSearch(searchPath, getTasksRespBody, nil)
		if scheduledTask != nil {
			break
		}

		// next page
		tasks := utils.PathSearch("schedules", getTasksRespBody, make([]interface{}, 0)).([]interface{})
		currentTotal += len(tasks)
		total := utils.PathSearch("total_count", getTasksRespBody, float64(0)).(float64)
		if currentTotal >= int(total) {
			break
		}
	}

	if scheduledTask == nil {
		return common.CheckDeletedDiag(d, golangsdk.ErrDefault404{}, "error retrieving scheduled task")
	}

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("instance_id", utils.PathSearch("instance_id", scheduledTask, nil)),
		d.Set("instance_name", utils.PathSearch("instance_name", scheduledTask, nil)),
		d.Set("instance_status", utils.PathSearch("instance_status", scheduledTask, nil)),
		d.Set("job_id", utils.PathSearch("job_id", scheduledTask, nil)),
		d.Set("job_name", utils.PathSearch("job_name", scheduledTask, nil)),
		d.Set("create_time", utils.PathSearch("create_time", scheduledTask, nil)),
		d.Set("start_time", utils.PathSearch("start_time", scheduledTask, nil)),
		d.Set("end_time", utils.PathSearch("end_time", scheduledTask, nil)),
		d.Set("job_status", utils.PathSearch("job_status", scheduledTask, nil)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func resourceDDSMysqlScheduledTaskCancelDelete(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	errorMsg := "Deleting DDS scheduled task cancel resource is not supported. The resource is only removed from the state," +
		" the DDS scheduled task remains in the cloud."
	return diag.Diagnostics{
		diag.Diagnostic{
			Severity: diag.Warning,
			Summary:  errorMsg,
		},
	}
}
