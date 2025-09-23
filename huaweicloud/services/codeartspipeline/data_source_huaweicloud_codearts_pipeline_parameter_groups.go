package codeartspipeline

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

// @API CodeArtsPipeline POST /v5/{project_id}/api/pipeline/variable/group/list
func DataSourceCodeArtsPipelineParameterGroups() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceCodeArtsPipelineParameterGroupsRead,

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
				Description: `Specifies the CodeArts project ID.`,
			},
			"name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specifies the parameter group name.`,
			},
			"groups": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: `Indicates the parameter group list.`,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `Indicates the parameter group ID.`,
						},
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `Indicates the parameter group name.`,
						},
						"description": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `Indicates the parameter group description.`,
						},
						"variables": {
							Type:        schema.TypeSet,
							Computed:    true,
							Description: `Indicates the parameter list.`,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"name": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: `Indicates the custom variable name.`,
									},
									"sequence": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: `Indicates the parameter sequence, starting from 1.`,
									},
									"type": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: `Indicates the custom parameter type.`,
									},
									"value": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: `Indicates the custom parameter default value.`,
									},
									"is_secret": {
										Type:        schema.TypeBool,
										Computed:    true,
										Description: `Indicates whether it is a private parameter.`,
									},
									"description": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: `Indicates the parameter description.`,
									},
								},
							},
						},
						"related_pipelines": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: `Indicates the associated pipeline.`,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"name": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: `Indicates the pipeline name.`,
									},
									"id": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: `Indicates the pipeline ID.`,
									},
								},
							},
						},
						"creator_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `Indicates the creator ID.`,
						},
						"updater_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `Indicates the updater ID.`,
						},
						"creator_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `Indicates the creator name.`,
						},
						"updater_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `Indicates the updater name.`,
						},
						"create_time": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: `Indicates the create time.`,
						},
						"update_time": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: `Indicates the update time.`,
						},
					},
				},
			},
		},
	}
}

func dataSourceCodeArtsPipelineParameterGroupsRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	client, err := cfg.NewServiceClient("codearts_pipeline", region)
	if err != nil {
		return diag.Errorf("error creating CodeArts Pipeline client: %s", err)
	}

	getHttpUrl := "v5/{project_id}/api/pipeline/variable/group/list"
	getPath := client.Endpoint + getHttpUrl
	getPath = strings.ReplaceAll(getPath, "{project_id}", d.Get("project_id").(string))
	getOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	offset := 0
	rst := make([]map[string]interface{}, 0)
	for {
		getOpt.JSONBody = utils.RemoveNil(buildPipelineCodeArtsPipelineParameterGroupsQueryParams(d, offset))
		getResp, err := client.Request("POST", getPath, &getOpt)
		if err != nil {
			return diag.Errorf("error retrieving pipeline parameter groups: %s", err)
		}
		getRespBody, err := utils.FlattenResponse(getResp)
		if err != nil {
			return diag.Errorf("error flatten response: %s", err)
		}

		groups := utils.PathSearch("pipeline_variable_groups", getRespBody, make([]interface{}, 0)).([]interface{})
		if len(groups) == 0 {
			break
		}

		for _, group := range groups {
			rst = append(rst, map[string]interface{}{
				"id":                utils.PathSearch("id", group, nil),
				"name":              utils.PathSearch("name", group, nil),
				"description":       utils.PathSearch("description", group, nil),
				"variables":         flattenPipelineParameterGroupVariables(group),
				"related_pipelines": flattenPipelineParameterGroupRelatedPipelines(group),
				"creator_id":        utils.PathSearch("creator_id", group, nil),
				"updater_id":        utils.PathSearch("updater_id", group, nil),
				"creator_name":      utils.PathSearch("creator_name", group, nil),
				"updater_name":      utils.PathSearch("updater_name", group, nil),
				"create_time":       utils.PathSearch("create_time", group, nil),
				"update_time":       utils.PathSearch("update_time", group, nil),
			})
		}

		offset += len(groups)
	}

	id, err := uuid.GenerateUUID()
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(id)

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("groups", rst),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func buildPipelineCodeArtsPipelineParameterGroupsQueryParams(d *schema.ResourceData, offset int) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"name":   utils.ValueIgnoreEmpty(d.Get("name")),
		"limit":  50,
		"offset": offset,
	}

	return bodyParams
}
