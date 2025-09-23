package codeartspipeline

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

// @API CodeArtsPipeline POST /v1/{domain_id}/relation/plugins
func DataSourceCodeartsPipelineAvailablePlugins() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceCodeartsPipelineAvailablePluginsRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `Specifies the region in which to query the resource. If omitted, the provider-level region will be used.`,
			},
			"business_type": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specifies the service type.`,
			},
			"regex_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specifies the regex name.`,
			},
			"use_condition": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specifies the use condition.`,
			},
			"input_repo_type": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specifies the source code repository type.`,
			},
			"input_source_type": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specifies the input source type, whether a pipeline has one source or multiple sources.`,
			},
			"data": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: `Indicates the result set.`,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"unique_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `Indicates the unique ID.`,
						},
						"display_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `Indicates the display name.`,
						},
						"business_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `Indicates the service type.`,
						},
						"conditions": {
							Type:        schema.TypeList,
							Computed:    true,
							Elem:        &schema.Schema{Type: schema.TypeString},
							Description: `Indicates the conditions.`,
						},
						"removable": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: `Indicates whether it is removable.`,
						},
						"cloneable": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: `Indicates whether it is replicable.`,
						},
						"disabled": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: `Indicates whether it is disabled.`,
						},
						"editable": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: `Indicates whether it is editable.`,
						},
						"plugins_list": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: `Indicates the extension list.`,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"unique_id": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: `Indicates the unique ID.`,
									},
									"display_name": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: `Indicates the display name.`,
									},
									"plugin_name": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: `Indicates the extension name.`,
									},
									"disabled": {
										Type:        schema.TypeBool,
										Computed:    true,
										Description: `Indicates whether it is disabled.`,
									},
									"group_name": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: `Indicates the group name.`,
									},
									"group_type": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: `Indicates the group type.`,
									},
									"plugin_attribution": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: `Indicates the attribute.`,
									},
									"plugin_composition_type": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: `Indicates the combination extension.`,
									},
									"runtime_attribution": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: `Indicates the runtime attributes.`,
									},
									"description": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: `Indicates the description.`,
									},
									"version_attribution": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: `Indicates the version attribute.`,
									},
									"icon_url": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: `Indicates the icon URL.`,
									},
									"multi_step_editable": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: `Indicates whether it is editable.`,
									},
									"location": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: `Indicates the address.`,
									},
									"publisher_unique_id": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: `Indicates the publisher ID.`,
									},
									"manifest_version": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: `Indicates the version.`,
									},
									"all_steps": {
										Type:        schema.TypeList,
										Computed:    true,
										Description: `Indicates the basic extension list.`,
										Elem:        dataPluginsListAllStepsElem(),
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

// dataPluginsListAllStepsElem
// The Elem of "data.plugins_list.all_steps"
func dataPluginsListAllStepsElem() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"plugin_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Indicates the extension name.`,
			},
			"display_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Indicates the display name.`,
			},
			"version": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Indicates the version.`,
			},
		},
	}
}

func dataSourceCodeartsPipelineAvailablePluginsRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	client, err := cfg.NewServiceClient("codearts_pipeline", region)
	if err != nil {
		return diag.Errorf("error creating CodeArts Pipeline client: %s", err)
	}

	getHttpUrl := "v1/{domain_id}/relation/plugins?limit=100"
	getPath := client.Endpoint + getHttpUrl
	getPath = strings.ReplaceAll(getPath, "{domain_id}", cfg.DomainID)
	getOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		JSONBody:         buildPipelineCodeArtsPipelineAvailablePluginsQueryParams(d),
	}

	offset := 0
	rst := make([]map[string]interface{}, 0)
	for {
		currentPath := getPath + fmt.Sprintf("&offset=%d", offset)
		getResp, err := client.Request("POST", currentPath, &getOpt)
		if err != nil {
			return diag.Errorf("error retrieving pipeline available plugins: %s", err)
		}
		getRespBody, err := utils.FlattenResponse(getResp)
		if err != nil {
			return diag.Errorf("error flatten response: %s", err)
		}

		plugins := utils.PathSearch("data", getRespBody, make([]interface{}, 0)).([]interface{})
		if len(plugins) == 0 {
			break
		}

		for _, plugin := range plugins {
			id := utils.PathSearch("unique_id", plugin, "").(string)
			found := false
			for _, r := range rst {
				if r["unique_id"].(string) == id {
					temp := r["plugins_list"].([]interface{})
					temp = append(temp, flattenPipelineAvailablePluginsPluginsList(plugin)...)
					r["plugins_list"] = temp
					found = true
					break
				}
			}
			if !found {
				rst = append(rst, map[string]interface{}{
					"unique_id":     utils.PathSearch("unique_id", plugin, nil),
					"display_name":  utils.PathSearch("display_name", plugin, nil),
					"business_type": utils.PathSearch("business_type", plugin, nil),
					"conditions":    utils.PathSearch("conditions", plugin, nil),
					"removable":     utils.PathSearch("removable", plugin, nil),
					"cloneable":     utils.PathSearch("cloneable", plugin, nil),
					"disabled":      utils.PathSearch("disabled", plugin, nil),
					"editable":      utils.PathSearch("editable", plugin, nil),
					"plugins_list":  flattenPipelineAvailablePluginsPluginsList(plugin),
				})
			}
		}

		// use total, which is the total of `plugins_list` in actual
		total := utils.PathSearch("total", getRespBody, float64(0)).(float64)
		offset += 100
		if offset >= int(total) {
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
		d.Set("data", rst),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func buildPipelineCodeArtsPipelineAvailablePluginsQueryParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		// have to send empty value, otherwise, will get empty list
		"business_type":     d.Get("business_type"),
		"regex_name":        d.Get("regex_name"),
		"use_condition":     d.Get("use_condition"),
		"input_repo_type":   d.Get("input_repo_type"),
		"input_source_type": d.Get("input_source_type"),
	}

	return bodyParams
}

func flattenPipelineAvailablePluginsPluginsList(resp interface{}) []interface{} {
	params := utils.PathSearch("plugins_list", resp, nil)
	if params == nil {
		return nil
	}

	if paramsList, ok := params.([]interface{}); ok && len(paramsList) > 0 {
		result := make([]interface{}, 0, len(paramsList))
		for _, p := range paramsList {
			param := p.(map[string]interface{})
			res := map[string]interface{}{
				"unique_id":               utils.PathSearch("unique_id", param, nil),
				"display_name":            utils.PathSearch("display_name", param, nil),
				"plugin_name":             utils.PathSearch("plugin_name", param, nil),
				"disabled":                utils.PathSearch("disabled", param, nil),
				"group_name":              utils.PathSearch("group_name", param, nil),
				"group_type":              utils.PathSearch("group_type", param, nil),
				"plugin_attribution":      utils.PathSearch("plugin_attribution", param, nil),
				"plugin_composition_type": utils.PathSearch("plugin_composition_type", param, nil),
				"runtime_attribution":     utils.PathSearch("runtime_attribution", param, nil),
				"description":             utils.PathSearch("description", param, nil),
				"version_attribution":     utils.PathSearch("version_attribution", param, nil),
				"icon_url":                utils.PathSearch("icon_url", param, nil),
				"multi_step_editable":     utils.PathSearch("multi_step_editable", param, nil),
				"location":                utils.PathSearch("location", param, nil),
				"publisher_unique_id":     utils.PathSearch("publisher_unique_id", param, nil),
				"manifest_version":        utils.PathSearch("manifest_version", param, nil),
				"all_steps":               flattenPipelineAvailablePluginsPluginsListAllSteps(param),
			}
			result = append(result, res)
		}

		return result
	}
	return nil
}

func flattenPipelineAvailablePluginsPluginsListAllSteps(resp interface{}) []interface{} {
	params := utils.PathSearch("all_steps", resp, nil)
	if params == nil {
		return nil
	}

	if paramsList, ok := params.([]interface{}); ok && len(paramsList) > 0 {
		result := make([]interface{}, 0, len(paramsList))
		for _, p := range paramsList {
			param := p.(map[string]interface{})
			res := map[string]interface{}{
				"plugin_name":  utils.PathSearch("plugin_name", param, nil),
				"display_name": utils.PathSearch("display_name", param, nil),
				"version":      utils.PathSearch("version", param, nil),
			}
			result = append(result, res)
		}

		return result
	}
	return nil
}
