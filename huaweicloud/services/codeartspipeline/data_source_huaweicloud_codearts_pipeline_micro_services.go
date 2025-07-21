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

// @API CodeArtsPipeline POST /v2/{cloudProjectId}/component/list/query
func DataSourceCodeArtsPipelineMicroServices() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceCodeArtsPipelineMicroServicesRead,

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
				Description: `Specifies the micro service name.`,
			},
			"sort_dir": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specifies the sorting sequence.`,
			},
			"micro_services": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: `Indicates the micro service list.`,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `Indicates the micro service ID.`,
						},
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `Indicates the micro service name.`,
						},
						"type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `Indicates the micro service type.`,
						},
						"description": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `Indicates the micro service description.`,
						},
						"status": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `Indicates the micro service status.`,
						},
						"parent_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `Indicates the micro service parent ID.`,
						},
						"repos": {
							Type:        schema.TypeSet,
							Computed:    true,
							Description: `Indicates the repository information.`,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"type": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: `Indicates the repository type.`,
									},
									"repo_id": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: `Indicates the repository ID.`,
									},
									"http_url": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: `Indicates the HTTP address of the Git repository.`,
									},
									"git_url": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: `Indicates the Git address of the Git repository.`,
									},
									"branch": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: `Indicates the branch.`,
									},
									"language": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: `Indicates the language.`,
									},
									"endpoint_id": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: `Indicates the endpoint ID.`,
									},
								},
							},
						},
						"is_followed": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: `Indicates whether the micro service is followed.`,
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
							Type:        schema.TypeString,
							Computed:    true,
							Description: `Indicates the create time.`,
						},
						"update_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `Indicates the update time.`,
						},
					},
				},
			},
		},
	}
}

func dataSourceCodeArtsPipelineMicroServicesRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	client, err := cfg.NewServiceClient("codearts_pipeline", region)
	if err != nil {
		return diag.Errorf("error creating CodeArts Pipeline client: %s", err)
	}

	getHttpUrl := "v2/{cloudProjectId}/component/list/query"
	getPath := client.Endpoint + getHttpUrl
	getPath = strings.ReplaceAll(getPath, "{cloudProjectId}", d.Get("project_id").(string))
	getOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	offset := 0
	rst := make([]map[string]interface{}, 0)
	for {
		getOpt.JSONBody = utils.RemoveNil(buildPipelineCodeArtsPipelineMicroServicesQueryParams(d, offset))
		getResp, err := client.Request("POST", getPath, &getOpt)
		if err != nil {
			return diag.Errorf("error retrieving pipeline micro services: %s", err)
		}
		getRespBody, err := utils.FlattenResponse(getResp)
		if err != nil {
			return diag.Errorf("error flatten response: %s", err)
		}

		microServices := utils.PathSearch("data", getRespBody, make([]interface{}, 0)).([]interface{})
		if len(microServices) == 0 {
			break
		}

		for _, microService := range microServices {
			rst = append(rst, map[string]interface{}{
				"id":           utils.PathSearch("id", microService, nil),
				"name":         utils.PathSearch("name", microService, nil),
				"description":  utils.PathSearch("description", microService, nil),
				"type":         utils.PathSearch("type", microService, nil),
				"parent_id":    utils.PathSearch("parent_id", microService, nil),
				"repos":        flattenPipelineMicroServiceRepos(microService),
				"creator_id":   utils.PathSearch("creator_id", microService, nil),
				"updater_id":   utils.PathSearch("updater_id", microService, nil),
				"creator_name": utils.PathSearch("creator_name", microService, nil),
				"updater_name": utils.PathSearch("updater_name", microService, nil),
				"create_time":  utils.PathSearch("create_time", microService, nil),
				"update_time":  utils.PathSearch("update_time", microService, nil),
				"status":       utils.PathSearch("status", microService, nil),
				"is_followed":  utils.PathSearch("is_followed", microService, nil),
			})
		}

		offset += len(microServices)
	}

	id, err := uuid.GenerateUUID()
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(id)

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("micro_services", rst),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func buildPipelineCodeArtsPipelineMicroServicesQueryParams(d *schema.ResourceData, offset int) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"sort_dir": utils.ValueIgnoreEmpty(d.Get("sort_dir")),
		"name":     utils.ValueIgnoreEmpty(d.Get("name")),
		"limit":    50,
		"offset":   offset,
	}

	return bodyParams
}
