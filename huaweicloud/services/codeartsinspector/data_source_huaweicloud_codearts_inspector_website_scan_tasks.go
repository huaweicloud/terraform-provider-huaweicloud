package codeartsinspector

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

// @API VSS GET /v3/{project_id}/webscan/tasks/histories
func DataSourceCodeartsInspectorWebsiteScanTasks() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceCodeartsInspectorWebsiteScanTasksRead,

		Schema: map[string]*schema.Schema{
			"domain_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `Specifies the website domain ID.`,
			},
			"tasks": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: `Indicates the tasks list.`,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"task_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `Indicates the task ID.`,
						},
						"task_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `Indicates the task name.`,
						},
						"task_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `Indicates the task type.`,
						},
						"url": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `Indicates the destination URL to scan.`,
						},
						"task_status": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `Indicates the task status.`,
						},
						"reason": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `Indicates the description of task status.`,
						},
						"created_at": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `Indicates the create time of the task.`,
						},
						"start_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `Indicates the start time of the task.`,
						},
						"end_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `Indicates the end time of the task.`,
						},
						"schedule_status": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `Indicates the monitor task status.`,
						},
						"safe_level": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `Indicates the security level.`,
						},
						"high": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: `Indicates the number of high-risk vulnerabilities.`,
						},
						"middle": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: `Indicates the number of medium-risk vulnerabilities.`,
						},
						"low": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: `Indicates the number of low-severity vulnerabilities.`,
						},
						"hint": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: `Indicates the number of hint-risk vulnerabilities.`,
						},
						"progress": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: `Indicates the task progress.`,
						},
						"pack_num": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: `Indicates the total number of packages.`,
						},
						"score": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: `Indicates the safety score.`,
						},
						"domain_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `Indicates the domain name.`,
						},
						"task_period": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `Indicates the scheduled trigger period of the monitor task.`,
						},
						"timer": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `Indicates the scheduled trigger time of the normal task.`,
						},
						"trigger_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `Indicates the scheduled trigger time of the monitor task.`,
						},
						"malicious_link": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: `Indicates whether to perform link health detection.`,
						},
						"scan_mode": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `Indicates the task scan mode.`,
						},
						"port_scan": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: `Indicates whether to perform port scanning.`,
						},
						"weak_pwd_scan": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: `Indicates whether to scan for weak passwords.`,
						},
						"cve_check": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: `Indicates whether to perform CVE vulnerability scanning.`,
						},
						"text_check": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: `Indicates whether to conduct website content compliance text detection.`,
						},
						"picture_check": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: `Indicates whether to conduct website content compliance image detection.`,
						},
						"malicious_code": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: `Indicates whether to perform malicious code scanning.`,
						},
					},
				},
			},
		},
	}
}

func dataSourceCodeartsInspectorWebsiteScanTasksRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	client, err := cfg.NewServiceClient("vss", region)
	if err != nil {
		return diag.Errorf("error creating CodeArts inspector client: %s", err)
	}

	getHttpUrl := "v3/{project_id}/webscan/tasks/histories"
	getPath := client.Endpoint + getHttpUrl
	getPath = strings.ReplaceAll(getPath, "{project_id}", client.ProjectID)
	getOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	// pageLimit is `10`
	getPath += fmt.Sprintf("?limit=%v", pageLimit)
	getPath += buildCodeartsInspectorWebsitesQueryParams(d)

	currentTotal := 0

	rst := make([]map[string]interface{}, 0)
	for {
		currentPath := getPath + fmt.Sprintf("&offset=%d", currentTotal)
		getResp, err := client.Request("GET", currentPath, &getOpt)
		if err != nil {
			return diag.Errorf("error retrieving tasks: %s", err)
		}
		getRespBody, err := utils.FlattenResponse(getResp)
		if err != nil {
			return diag.Errorf("error flatten response: %s", err)
		}

		tasks := utils.PathSearch("data", getRespBody, make([]interface{}, 0)).([]interface{})
		for _, task := range tasks {
			rst = append(rst, flattenCodeartsInspectorWebsiteScanTasks(task))
		}

		currentTotal += len(tasks)
		totalCount := utils.PathSearch("total", getRespBody, float64(0))
		if int(totalCount.(float64)) <= currentTotal {
			break
		}
	}

	id, err := uuid.GenerateUUID()
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(id)

	mErr := multierror.Append(nil,
		d.Set("tasks", rst),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenCodeartsInspectorWebsiteScanTasks(task interface{}) map[string]interface{} {
	return map[string]interface{}{
		"task_id":         utils.PathSearch("task_id", task, nil),
		"task_name":       utils.PathSearch("task_name", task, nil),
		"task_type":       utils.PathSearch("task_type", task, nil),
		"url":             utils.PathSearch("url", task, nil),
		"task_period":     utils.PathSearch("task_settings.task_period", task, nil),
		"timer":           utils.PathSearch("task_settings.timer", task, nil),
		"trigger_time":    utils.PathSearch("task_settings.trigger_time", task, nil),
		"malicious_link":  utils.PathSearch("task_settings.task_config.malicious_link", task, nil),
		"scan_mode":       utils.PathSearch("task_settings.task_config.scan_mode", task, nil),
		"port_scan":       utils.PathSearch("task_settings.task_config.port_scan", task, nil),
		"weak_pwd_scan":   utils.PathSearch("task_settings.task_config.weak_pwd_scan", task, nil),
		"cve_check":       utils.PathSearch("task_settings.task_config.cve_check", task, nil),
		"text_check":      utils.PathSearch("task_settings.task_config.text_check", task, nil),
		"picture_check":   utils.PathSearch("task_settings.task_config.picture_check", task, nil),
		"malicious_code":  utils.PathSearch("task_settings.task_config.malicious_code", task, nil),
		"task_status":     utils.PathSearch("task_status", task, nil),
		"reason":          utils.PathSearch("reason", task, nil),
		"created_at":      utils.PathSearch("create_time", task, nil),
		"start_time":      utils.PathSearch("start_time", task, nil),
		"end_time":        utils.PathSearch("end_time", task, nil),
		"schedule_status": utils.PathSearch("schedule_status", task, nil),
		"safe_level":      utils.PathSearch("safe_level", task, nil),
		"progress":        utils.PathSearch("progress", task, nil),
		"pack_num":        utils.PathSearch("pack_num", task, nil),
		"score":           utils.PathSearch("score", task, nil),
		"domain_name":     utils.PathSearch("domain_name", task, nil),
		"high":            utils.PathSearch("statistics.high", task, nil),
		"low":             utils.PathSearch("statistics.low", task, nil),
		"middle":          utils.PathSearch("statistics.middle", task, nil),
		"hint":            utils.PathSearch("statistics.hint", task, nil),
	}
}
