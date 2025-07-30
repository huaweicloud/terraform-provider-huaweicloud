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

// @API CodeArtsPipeline POST /v1/{domain_id}/agent-plugin/query-all
func DataSourceCodeArtsPipelinePlugins() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceCodeArtsPipelinePluginsRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `Specifies the region in which to query the resource. If omitted, the provider-level region will be used.`,
			},
			"plugin_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specifies the plugin name.`,
			},
			"regex_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specifies the match name.`,
			},
			"maintainer": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specifies the maintenance engineer.`,
			},
			"business_type": {
				Type:        schema.TypeList,
				Optional:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: `Specifies the service type.`,
			},
			"plugin_attribution": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specifies the extension attribute, official or custom.`,
			},
			"plugins": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: `Indicates the plugin list.`,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"plugin_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `Indicates the plugin name.`,
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
						"version_description": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `Indicates the version description.`,
						},
						"version_attribution": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `Indicates the version attribution.`,
						},
						"description": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `Indicates the description.`,
						},
						"unique_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `Indicates the unique ID.`,
						},
						"op_user": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `Indicates the operator.`,
						},
						"op_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `Indicates the operation time.`,
						},
						"plugin_composition_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `Indicates the combination type.`,
						},
						"plugin_attribution": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `Indicates the attribute.`,
						},
						"workspace_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `Indicates the tenant ID.`,
						},
						"business_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `Indicates the service type.`,
						},
						"business_type_display_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `Indicates the display name of service type.`,
						},
						"maintainers": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `Indicates the maintenance engineer.`,
						},
						"icon_url": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `Indicates the icon URL.`,
						},
						"refer_count": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: `Indicates the number of references.`,
						},
						"usage_count": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: `Indicates the number of usages.`,
						},
						"runtime_attribution": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `Indicates the runtime attribution.`,
						},
						"active": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: `Indicates whether the plugin is activate or not.`,
						},
					},
				},
			},
		},
	}
}

func dataSourceCodeArtsPipelinePluginsRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	client, err := cfg.NewServiceClient("codearts_pipeline", region)
	if err != nil {
		return diag.Errorf("error creating CodeArts Pipeline client: %s", err)
	}

	getHttpUrl := "v1/{domain_id}/agent-plugin/query-all?limit=10"
	getPath := client.Endpoint + getHttpUrl
	getPath = strings.ReplaceAll(getPath, "{domain_id}", cfg.DomainID)
	getOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		JSONBody:         buildPipelineCodeArtsPipelinePluginsQueryParams(d),
	}

	offset := 0
	rst := make([]map[string]interface{}, 0)
	for {
		currentPath := getPath + fmt.Sprintf("&offset=%d", offset)
		getResp, err := client.Request("POST", currentPath, &getOpt)
		if err != nil {
			return diag.Errorf("error retrieving pipeline plugins: %s", err)
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
			rst = append(rst, map[string]interface{}{
				"plugin_name":                utils.PathSearch("plugin_name", plugin, nil),
				"version":                    utils.PathSearch("version", plugin, nil),
				"version_description":        utils.PathSearch("version_description", plugin, nil),
				"version_attribution":        utils.PathSearch("version_attribution", plugin, nil),
				"description":                utils.PathSearch("description", plugin, nil),
				"unique_id":                  utils.PathSearch("unique_id", plugin, nil),
				"op_user":                    utils.PathSearch("op_user", plugin, nil),
				"op_time":                    utils.PathSearch("op_time", plugin, nil),
				"plugin_composition_type":    utils.PathSearch("plugin_composition_type", plugin, nil),
				"plugin_attribution":         utils.PathSearch("plugin_attribution", plugin, nil),
				"workspace_id":               utils.PathSearch("workspace_id", plugin, nil),
				"business_type":              utils.PathSearch("business_type", plugin, nil),
				"business_type_display_name": utils.PathSearch("business_type_display_name", plugin, nil),
				"maintainers":                utils.PathSearch("maintainers", plugin, nil),
				"icon_url":                   utils.PathSearch("icon_url", plugin, nil),
				"refer_count":                utils.PathSearch("refer_count", plugin, nil),
				"usage_count":                utils.PathSearch("usage_count", plugin, nil),
				"runtime_attribution":        utils.PathSearch("runtime_attribution", plugin, nil),
				"active":                     utils.PathSearch("active", plugin, nil),
			})
		}

		offset += len(plugins)
	}

	id, err := uuid.GenerateUUID()
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(id)

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("plugins", rst),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func buildPipelineCodeArtsPipelinePluginsQueryParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		// have to send empty value, otherwise, will get empty list
		"plugin_name":        d.Get("plugin_name"),
		"regex_name":         d.Get("regex_name"),
		"maintainer":         d.Get("maintainer"),
		"business_type":      d.Get("business_type"),
		"plugin_attribution": d.Get("plugin_attribution"),
	}

	return bodyParams
}
