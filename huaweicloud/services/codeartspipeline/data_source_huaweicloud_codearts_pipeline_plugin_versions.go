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

// @API CodeArtsPipeline GET /v1/{domain_id}/agent-plugin/query
func DataSourceCodeartsPipelinePluginVersions() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceCodeartsPipelinePluginVersionsRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `Specifies the region in which to query the resource. If omitted, the provider-level region will be used.`,
			},
			"plugin_name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `Specifies the plugin name.`,
			},
			"versions": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: `Indicates the version list.`,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"plugin_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `Indicates the extension name.`,
						},
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
							Description: `Indicates the version attribute.`,
						},
						"description": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `Indicates the description.`,
						},
						"plugin_attribution": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `Indicates the attribute.`,
						},
						"plugin_composition_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `Indicates the combination type.`,
						},
						"icon_url": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `Indicates the icon URL.`,
						},
						"workspace_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `Indicates the tenant ID.`,
						},
						"refer_count": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: `Indicates the number of references.`,
						},
						"active": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: `Indicates the activated or not.`,
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
						"maintainers": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `Indicates the maintenance engineer.`,
						},
						"usage_count": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: `Indicates the number of usages.`,
						},
						"runtime_attribution": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `Indicates the runtime attributes.`,
						},
					},
				},
			},
		},
	}
}

func dataSourceCodeartsPipelinePluginVersionsRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	client, err := cfg.NewServiceClient("codearts_pipeline", region)
	if err != nil {
		return diag.Errorf("error creating CodeArts Pipeline client: %s", err)
	}

	httpUrl := "v1/{domain_id}/agent-plugin/query?plugin_name={plugin_name}&limit=10"
	getPath := client.Endpoint + httpUrl
	getPath = strings.ReplaceAll(getPath, "{domain_id}", cfg.DomainID)
	getPath = strings.ReplaceAll(getPath, "{plugin_name}", d.Get("plugin_name").(string))
	getOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	offset := 0
	rst := make([]map[string]interface{}, 0)
	for {
		currentPath := getPath + fmt.Sprintf("&offset=%d", offset)
		getResp, err := client.Request("GET", currentPath, &getOpt)
		if err != nil {
			return diag.Errorf("error getting plugin versions: %s", err)
		}
		getRespBody, err := utils.FlattenResponse(getResp)
		if err != nil {
			return diag.Errorf("error flattening response: %s", err)
		}
		if err := checkResponseError(getRespBody, ""); err != nil {
			return diag.Errorf("error getting plugin versions: %s", err)
		}

		versions := utils.PathSearch("data", getRespBody, make([]interface{}, 0)).([]interface{})
		if len(versions) == 0 {
			break
		}

		for _, version := range versions {
			rst = append(rst, map[string]interface{}{
				"plugin_name":                utils.PathSearch("plugin_name", version, nil),
				"unique_id":                  utils.PathSearch("unique_id", version, nil),
				"display_name":               utils.PathSearch("display_name", version, nil),
				"version":                    utils.PathSearch("version", version, nil),
				"version_description":        utils.PathSearch("version_description", version, nil),
				"version_attribution":        utils.PathSearch("version_attribution", version, nil),
				"description":                utils.PathSearch("description", version, nil),
				"plugin_attribution":         utils.PathSearch("plugin_attribution", version, nil),
				"plugin_composition_type":    utils.PathSearch("plugin_composition_type", version, nil),
				"icon_url":                   utils.PathSearch("icon_url", version, nil),
				"workspace_id":               utils.PathSearch("workspace_id", version, nil),
				"refer_count":                utils.PathSearch("refer_count", version, nil),
				"active":                     utils.PathSearch("active", version, nil),
				"business_type":              utils.PathSearch("business_type", version, nil),
				"business_type_display_name": utils.PathSearch("business_type_display_name", version, nil),
				"op_user":                    utils.PathSearch("op_user", version, nil),
				"op_time":                    utils.PathSearch("op_time", version, nil),
				"maintainers":                utils.PathSearch("maintainers", version, nil),
				"usage_count":                utils.PathSearch("usage_count", version, nil),
				"runtime_attribution":        utils.PathSearch("runtime_attribution", version, nil),
			})
		}

		offset += 10
	}

	id, err := uuid.GenerateUUID()
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(id)

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("versions", rst),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}
