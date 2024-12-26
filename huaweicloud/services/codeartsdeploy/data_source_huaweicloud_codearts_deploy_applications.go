package codeartsdeploy

import (
	"context"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API CodeArtsDeploy POST /v1/applications/list
func DataSourceCodeartsDeployApplications() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceCodeartsDeployApplicationsRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `Specifies the region in which to query the resource. If omitted, the provider-level region will be used.`,
			},
			"project_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `Specifies the project ID.`,
			},
			"states": {
				Type:        schema.TypeList,
				Optional:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: `Specifies the application status list.`,
			},
			"group_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specifies the application group ID.`,
			},
			"applications": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: `Indicates the application list`,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"project_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `Indicates the project name.`,
						},
						"id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `Indicates the application ID.`,
						},
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `Indicates the application name.`,
						},
						"deploy_system": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `Indicates the deployment type.`,
						},
						"release_id": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: `Indicates the release ID.`,
						},
						"can_modify": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: `Indicates whether the user has permission to modify application.`,
						},
						"can_manage": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: `Indicates whether the user has permission to modify application permission.`,
						},
						"can_create_env": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: `Indicates whether the user has permission to create environment in application.`,
						},
						"can_execute": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: `Indicates whether the user has permission to deploy.`,
						},
						"can_copy": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: `Indicates whether the user has permission to clone application.`,
						},
						"can_delete": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: `Indicates whether the user has permission to delete application.`,
						},
						"can_view": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: `Indicates whether the user has permission to view application.`,
						},
						"can_disable": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: `Indicates whether the user has permission to disable application.`,
						},
						"is_care": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: `Indicates whether application is saved to favorites.`,
						},
						"is_disable": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: `Indicates whether the application is disabled.`,
						},
						"create_user_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `Indicates the creator user ID.`,
						},
						"create_tenant_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `Indicates the created tenant ID.`,
						},
						"created_at": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `Indicates the created time.`,
						},
						"updated_at": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `Indicates the updated time.`,
						},
						"arrange_infos": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: `Indicates the deployment task information.`,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"id": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: `Indicates the task ID.`,
									},
									"state": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: `Indicates the deployment task status.`,
									},
									"deploy_system": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: `Indicates the deployment task type.`,
									},
								},
							},
						},
						"duration": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `Indicates the deployment duration.`,
						},
						"execution_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `Indicates the latest execution time.`,
						},
						"end_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `Indicates the deployment end time.`,
						},
						"executor_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `Indicates the executor user ID.`,
						},
						"executor_nick_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `Indicates the executor user name.`,
						},
						"execution_state": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `Indicates the execution status.`,
						},
					},
				},
			},
		},
	}
}

func dataSourceCodeartsDeployApplicationsRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	client, err := cfg.NewServiceClient("codearts_deploy", region)
	if err != nil {
		return diag.Errorf("error creating CodeArts deploy client: %s", err)
	}

	getHttpUrl := "v1/applications/list"
	getPath := client.Endpoint + getHttpUrl
	getOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json;charset=utf-8",
		},
	}

	pageIndex := 1

	rst := make([]map[string]interface{}, 0)
	for {
		getOpt.JSONBody = utils.RemoveNil(buildGetDeployApplicationsBodyParams(d, pageIndex))
		getResp, err := client.Request("POST", getPath, &getOpt)
		if err != nil {
			return diag.Errorf("error retrieving applications: %s", err)
		}
		getRespBody, err := utils.FlattenResponse(getResp)
		if err != nil {
			return diag.Errorf("error flatten response: %s", err)
		}

		applications := utils.PathSearch("result", getRespBody, make([]interface{}, 0)).([]interface{})
		for _, application := range applications {
			rst = append(rst, map[string]interface{}{
				"project_name":       utils.PathSearch("project_name", application, nil),
				"id":                 utils.PathSearch("id", application, nil),
				"name":               utils.PathSearch("name", application, nil),
				"deploy_system":      utils.PathSearch("deploy_system", application, nil),
				"release_id":         utils.PathSearch("release_id", application, nil),
				"can_modify":         utils.PathSearch("can_modify", application, nil),
				"can_delete":         utils.PathSearch("can_delete", application, nil),
				"can_view":           utils.PathSearch("can_view", application, nil),
				"can_execute":        utils.PathSearch("can_execute", application, nil),
				"can_copy":           utils.PathSearch("can_copy", application, nil),
				"can_manage":         utils.PathSearch("can_manage", application, nil),
				"can_create_env":     utils.PathSearch("can_create_env", application, nil),
				"can_disable":        utils.PathSearch("can_disable", application, nil),
				"is_care":            utils.PathSearch("is_care", application, nil),
				"is_disable":         utils.PathSearch("is_disable", application, nil),
				"create_user_id":     utils.PathSearch("create_user_id", application, nil),
				"create_tenant_id":   utils.PathSearch("create_tenant_id", application, nil),
				"created_at":         utils.PathSearch("create_time", application, nil),
				"updated_at":         utils.PathSearch("update_time", application, nil),
				"duration":           utils.PathSearch("duration", application, nil),
				"execution_time":     utils.PathSearch("execution_time", application, nil),
				"end_time":           utils.PathSearch("end_time", application, nil),
				"executor_id":        utils.PathSearch("executor_id", application, nil),
				"executor_nick_name": utils.PathSearch("executor_nick_name", application, nil),
				"execution_state":    utils.PathSearch("execution_state", application, nil),
				"arrange_infos": flattenDataSourceDeployApplicationArranges(
					utils.PathSearch("arrange_infos", application, make([]interface{}, 0)).([]interface{})),
			})
		}

		total := utils.PathSearch("total_num", getRespBody, float64(0)).(float64)
		if pageSize*(pageIndex-1)+len(applications) >= int(total) {
			break
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
		d.Set("applications", rst),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func buildGetDeployApplicationsBodyParams(d *schema.ResourceData, page int) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"project_id": d.Get("project_id"),
		"page":       page,
		"size":       pageSize,
		"states":     utils.ValueIgnoreEmpty(d.Get("states")),
		"group_id":   utils.ValueIgnoreEmpty(d.Get("group_id")),
	}
	return bodyParams
}

func flattenDataSourceDeployApplicationArranges(resp []interface{}) []interface{} {
	rst := make([]interface{}, 0, len(resp))
	for _, v := range resp {
		rst = append(rst, map[string]interface{}{
			"id":            utils.PathSearch("id", v, nil),
			"state":         utils.PathSearch("state", v, nil),
			"deploy_system": utils.PathSearch("deploy_system", v, nil),
		})
	}

	return rst
}
