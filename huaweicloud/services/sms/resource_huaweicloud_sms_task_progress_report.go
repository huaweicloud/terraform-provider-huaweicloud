package sms

import (
	"context"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

var taskProgressReportNonUpdatableParams = []string{"task_id", "subtask_name", "progress", "replicatesize", "totalsize",
	"process_trace", "migrate_speed", "compress_rate", "remain_time", "total_cpu_usage", "agent_cpu_usage",
	"total_mem_usage", "agent_mem_usage", "total_disk_io", "agent_disk_io", "need_migration_test", "agent_time"}

// @API SMS PUT /v3/tasks/{task_id}/progress
func ResourceTaskProgressReport() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceTaskProgressReportCreate,
		ReadContext:   resourceTaskProgressReportRead,
		UpdateContext: resourceTaskProgressReportUpdate,
		DeleteContext: resourceTaskProgressReportDelete,

		CustomizeDiff: config.FlexibleForceNew(taskProgressReportNonUpdatableParams),

		Schema: map[string]*schema.Schema{
			"task_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"subtask_name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"progress": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"replicatesize": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"totalsize": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"process_trace": {
				Type:     schema.TypeString,
				Required: true,
			},
			"migrate_speed": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"compress_rate": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"remain_time": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"total_cpu_usage": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"agent_cpu_usage": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"total_mem_usage": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"agent_mem_usage": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"total_disk_io": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"agent_disk_io": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"need_migration_test": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"agent_time": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"enable_force_new": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"true", "false"}, false),
				Description:  utils.SchemaDesc("", utils.SchemaDescInput{Internal: true}),
			},
		},
	}
}

func resourceTaskProgressReportCreate(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	client, err := cfg.SmsV3Client(cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating SMS client: %s", err)
	}

	taskID := d.Get("task_id").(string)
	err = taskProgressReport(client, d, taskID)
	if err != nil {
		return diag.Errorf("error reporting SMS task progress: %s", err)
	}

	d.SetId(taskID)

	return nil
}

func taskProgressReport(client *golangsdk.ServiceClient, d *schema.ResourceData, taskID string) error {
	httpUrl := "v3/tasks/{task_id}/progress"

	createPath := client.Endpoint + httpUrl
	createPath = strings.ReplaceAll(createPath, "{task_id}", taskID)

	createOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json",
		},
		JSONBody: utils.RemoveNil(buildCreateTaskProgressReportBodyParams(d)),
	}

	_, err := client.Request("PUT", createPath, &createOpt)
	if err != nil {
		return err
	}
	return nil
}

func buildCreateTaskProgressReportBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"subtask_name":        d.Get("subtask_name"),
		"progress":            d.Get("progress"),
		"replicatesize":       d.Get("replicatesize"),
		"totalsize":           d.Get("totalsize"),
		"process_trace":       d.Get("process_trace"),
		"migrate_speed":       utils.ValueIgnoreEmpty(d.Get("migrate_speed")),
		"compress_rate":       utils.ValueIgnoreEmpty(d.Get("compress_rate")),
		"remain_time":         utils.ValueIgnoreEmpty(d.Get("remain_time")),
		"total_cpu_usage":     utils.ValueIgnoreEmpty(d.Get("total_cpu_usage")),
		"agent_cpu_usage":     utils.ValueIgnoreEmpty(d.Get("agent_cpu_usage")),
		"total_mem_usage":     utils.ValueIgnoreEmpty(d.Get("total_mem_usage")),
		"agent_mem_usage":     utils.ValueIgnoreEmpty(d.Get("agent_mem_usage")),
		"total_disk_io":       utils.ValueIgnoreEmpty(d.Get("total_disk_io")),
		"agent_disk_io":       utils.ValueIgnoreEmpty(d.Get("agent_disk_io")),
		"need_migration_test": utils.ValueIgnoreEmpty(d.Get("need_migration_test")),
		"agent_time":          utils.ValueIgnoreEmpty(d.Get("agent_time")),
	}

	return bodyParams
}

func resourceTaskProgressReportRead(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return nil
}

func resourceTaskProgressReportUpdate(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return nil
}

func resourceTaskProgressReportDelete(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	errorMsg := "Deleting the task progress report resource is not supported. The task progress report resource is only" +
		" removed from the state."
	return diag.Diagnostics{
		diag.Diagnostic{
			Severity: diag.Warning,
			Summary:  errorMsg,
		},
	}
}
