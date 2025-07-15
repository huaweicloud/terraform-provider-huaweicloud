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

// @API CodeArtsPipeline POST /v5/{tenant_id}/api/pipeline-templates/list
func DataSourceCodeArtsPipelineTemplates() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceCodeArtsPipelineTemplatesRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specifies the template name.`,
			},
			"language": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specifies the template language.`,
			},
			"is_system": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: `Specifies whether the template is a system template.`,
			},
			"sort_key": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specifies the sorting attribute.`,
			},
			"sort_dir": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specifies the sorting sequence.`,
			},
			"templates": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: `Indicates the template list.`,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `Indicates the template ID.`,
						},
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `Indicates the template name.`,
						},
						"icon": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `Indicates the template icon.`,
						},
						"manifest_version": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `Indicates the manifest version.`,
						},
						"language": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `Indicates the template language.`,
						},
						"description": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `Indicates the template description.`,
						},
						"is_system": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: `Indicates whether the template is a system template.`,
						},
						"create_time": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: `Indicates the creation time.`,
						},
						"update_time": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: `Indicates the last update time.`,
						},
						"creator_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `Indicates the creator.`,
						},
						"creator_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `Indicates the creator name.`,
						},
						"updater_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `Indicates the last updater.`,
						},
						"is_favorite": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: `Indicates whether it is a favorite template.`,
						},
						"is_show_source": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: `Indicates whether to display the pipeline source.`,
						},
						"stages": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: `Indicates the stage running information.`,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"name": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: `Indicates the stage name.`,
									},
									"sequence": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: `Indicates the serial number.`,
									},
								},
							},
						},
					},
				},
			},
		},
	}
}

func dataSourceCodeArtsPipelineTemplatesRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	client, err := cfg.NewServiceClient("codearts_pipeline", region)
	if err != nil {
		return diag.Errorf("error creating CodeArts Pipeline client: %s", err)
	}

	getHttpUrl := "v5/{tenant_id}/api/pipeline-templates/list"
	getPath := client.Endpoint + getHttpUrl
	getPath = strings.ReplaceAll(getPath, "{tenant_id}", cfg.DomainID)
	getOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	offset := 0
	rst := make([]map[string]interface{}, 0)
	for {
		getOpt.JSONBody = utils.RemoveNil(buildPipelineCodeArtsPipelineTemplatesQueryParams(d, offset))
		getResp, err := client.Request("POST", getPath, &getOpt)
		if err != nil {
			return diag.Errorf("error retrieving pipeline templates: %s", err)
		}
		getRespBody, err := utils.FlattenResponse(getResp)
		if err != nil {
			return diag.Errorf("error flatten response: %s", err)
		}

		templates := utils.PathSearch("templates", getRespBody, make([]interface{}, 0)).([]interface{})
		if len(templates) == 0 {
			break
		}

		for _, template := range templates {
			rst = append(rst, map[string]interface{}{
				"id":               utils.PathSearch("id", template, nil),
				"name":             utils.PathSearch("name", template, nil),
				"description":      utils.PathSearch("description", template, nil),
				"language":         utils.PathSearch("language", template, nil),
				"is_system":        utils.PathSearch("is_system", template, nil),
				"is_show_source":   utils.PathSearch("is_show_source", template, nil),
				"creator_id":       utils.PathSearch("creator_id", template, nil),
				"creator_name":     utils.PathSearch("creator_name", template, nil),
				"updater_id":       utils.PathSearch("updater_id", template, nil),
				"create_time":      utils.PathSearch("create_time", template, nil),
				"update_time":      utils.PathSearch("update_time", template, nil),
				"icon":             utils.PathSearch("icon", template, nil),
				"manifest_version": utils.PathSearch("manifest_version", template, nil),
				"is_favorite":      utils.PathSearch("is_collect", template, nil),
				"stages":           flattenDataSourcePipelineTemplatesStages(template),
			})
		}

		offset += len(templates)
	}

	id, err := uuid.GenerateUUID()
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(id)

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("templates", rst),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenDataSourcePipelineTemplatesStages(resp interface{}) []interface{} {
	stages := utils.PathSearch("stages", resp, make([]interface{}, 0)).([]interface{})
	if len(stages) > 0 {
		result := make([]interface{}, 0, len(stages))
		for _, v := range stages {
			stage := v.(map[string]interface{})
			m := map[string]interface{}{
				"name":     utils.PathSearch("name", stage, nil),
				"sequence": utils.PathSearch("sequence", stage, nil),
			}
			result = append(result, m)
		}

		return result
	}

	return nil
}

func buildPipelineCodeArtsPipelineTemplatesQueryParams(d *schema.ResourceData, offset int) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"sort_dir":  utils.ValueIgnoreEmpty(d.Get("sort_dir")),
		"sort_key":  utils.ValueIgnoreEmpty(d.Get("sort_key")),
		"is_system": utils.ValueIgnoreEmpty(d.Get("is_system")),
		"language":  utils.ValueIgnoreEmpty(d.Get("language")),
		"name":      utils.ValueIgnoreEmpty(d.Get("name")),
		"limit":     100,
		"offset":    offset,
	}

	return bodyParams
}
