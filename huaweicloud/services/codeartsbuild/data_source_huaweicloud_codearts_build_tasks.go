package codeartsbuild

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

// @API CodeArtsBuild GET /v1/job/{project_id}/list
func DataSourceCodeArtsBuildTasks() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceCodeArtsBuildTasksRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"project_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `Specifies the CodeArts project ID.`,
			},
			"search": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specifies the search condition.`,
			},
			"sort_field": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specifies the sorting field.`,
			},
			"sort_order": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specifies the sorting order.`,
			},
			"creator_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specifies the creator ID.`,
			},
			"build_status": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specifies the build status filter condition.`,
			},
			"by_group": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: `Specifies whether to group.`,
			},
			"group_path_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specifies the group ID.`,
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
						"creator": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `Indicates the task creator.`,
						},
						"user_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `Indicates the user name.`,
						},
						"last_build_time": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: `Indicates the latest execution time.`,
						},
						"health_score": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: `Indicates the health score.`,
						},
						"source_code": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `Indicates the code source.`,
						},
						"last_build_status": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `Indicates the latest build status.`,
						},
						"is_finished": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: `Indicates whether it has ended.`,
						},
						"disabled": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: `Indicates whether it is disabled.`,
						},
						"favorite": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: `Indicates whether it is favorited.`,
						},
						"is_modify": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: `Indicates whether there is permission to modify the task.`,
						},
						"is_delete": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: `Indicates whether there is permission to delete the task.`,
						},
						"is_execute": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: `Indicates whether there is permission to execute the task.`,
						},
						"is_copy": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: `Indicates whether there is permission to copy the task.`,
						},
						"is_forbidden": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: `Indicates whether there is permission to disable the task.`,
						},
						"is_view": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: `Indicates whether there is permission to view the task.`,
						},
						"last_build_user": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `Indicates the last build user.`,
						},
						"trigger_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `Indicates the trigger type.`,
						},
						"build_time": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: `Indicates the build time.`,
						},
						"scm_web_url": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `Indicates the code repository web address.`,
						},
						"scm_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `Indicates the code repository type.`,
						},
						"repo_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `Indicates the code repository ID.`,
						},
						"build_project_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `Indicates the build project ID.`,
						},
						"last_job_running_status": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `Indicates the last build time.`,
						},
						"last_build_user_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `Indicates the last build user ID.`,
						},
					},
				},
			},
		},
	}
}

func dataSourceCodeArtsBuildTasksRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	client, err := cfg.NewServiceClient("codearts_build", region)
	if err != nil {
		return diag.Errorf("error creating CodeArts Build client: %s", err)
	}

	projectId := d.Get("project_id").(string)
	getHttpUrl := "v1/job/{project_id}/list"
	getPath := client.Endpoint + getHttpUrl
	getPath = strings.ReplaceAll(getPath, "{project_id}", projectId)
	getOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	getPath += fmt.Sprintf("?page_size=%v", 100)
	getPath += buildCodeArtsBuildTasksQueryParams(d)
	pageIndex := 0
	rst := make([]map[string]interface{}, 0)
	for {
		currentPath := getPath + fmt.Sprintf("&page_index=%d", pageIndex)
		getResp, err := client.Request("GET", currentPath, &getOpt)
		if err != nil {
			return diag.Errorf("error retrieving build tasks: %s", err)
		}
		getRespBody, err := utils.FlattenResponse(getResp)
		if err != nil {
			return diag.Errorf("error flatten response: %s", err)
		}

		tasks := utils.PathSearch("result.job_list", getRespBody, make([]interface{}, 0)).([]interface{})
		if len(tasks) == 0 {
			break
		}

		for _, task := range tasks {
			rst = append(rst, map[string]interface{}{
				"id":                      utils.PathSearch("id", task, nil),
				"name":                    utils.PathSearch("job_name", task, nil),
				"creator":                 utils.PathSearch("job_creator", task, nil),
				"user_name":               utils.PathSearch("user_name", task, nil),
				"last_build_time":         utils.PathSearch("last_build_time", task, nil),
				"health_score":            utils.PathSearch("health_score", task, nil),
				"source_code":             utils.PathSearch("source_code", task, nil),
				"last_build_status":       utils.PathSearch("last_build_status", task, nil),
				"is_finished":             utils.PathSearch("is_finished", task, nil),
				"disabled":                utils.PathSearch("disabled", task, nil),
				"favorite":                utils.PathSearch("favorite", task, nil),
				"is_modify":               utils.PathSearch("is_modify", task, nil),
				"is_delete":               utils.PathSearch("is_delete", task, nil),
				"is_execute":              utils.PathSearch("is_execute", task, nil),
				"is_copy":                 utils.PathSearch("is_copy", task, nil),
				"is_forbidden":            utils.PathSearch("is_forbidden", task, nil),
				"is_view":                 utils.PathSearch("is_view", task, nil),
				"last_build_user":         utils.PathSearch("last_build_user", task, nil),
				"trigger_type":            utils.PathSearch("trigger_type", task, nil),
				"build_time":              utils.PathSearch("build_time", task, nil),
				"scm_web_url":             utils.PathSearch("scm_web_url", task, nil),
				"scm_type":                utils.PathSearch("scm_type", task, nil),
				"repo_id":                 utils.PathSearch("repo_id", task, nil),
				"build_project_id":        utils.PathSearch("build_project_id", task, nil),
				"last_job_running_status": utils.PathSearch("last_job_running_status", task, nil),
				"last_build_user_id":      utils.PathSearch("last_build_user_id", task, nil),
			})
		}

		pageIndex++
	}

	id, err := uuid.GenerateUUID()
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(id)

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("tasks", rst),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func buildCodeArtsBuildTasksQueryParams(d *schema.ResourceData) string {
	res := ""
	params := []string{"search", "sort_field", "sort_order", "creator_id", "build_status", "by_group", "group_path_id"}

	for _, param := range params {
		if v, ok := d.GetOk(param); ok {
			res = fmt.Sprintf("%s&%s=%v", res, param, v)
		}
	}

	return res
}
