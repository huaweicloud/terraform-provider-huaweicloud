package codeartspipeline

import (
	"context"
	"fmt"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API CodeArtsPipeline GET /v2/extensions/modules
func DataSourceCodeArtsPipelineModules() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceCodeArtsPipelineModulesRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specifies the module name.`,
			},
			"product_line": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specifies the product line.`,
			},
			"tags": {
				Type:        schema.TypeList,
				Optional:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: `Specifies the tags.`,
			},
			"modules": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: `Indicates the module list.`,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: `Indicates the module ID.`,
						},
						"base_url": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `Indicates the module base URL.`,
						},
						"description": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `Indicates the module description.`,
						},
						"location": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `Indicates the endpoint.`,
						},
						"module_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `Indicates the module ID.`,
						},
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `Indicates the module name.`,
						},
						"properties": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `Indicates the properties.`,
						},
						"publisher": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `Indicates the publisher.`,
						},
						"type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `Indicates the module type.`,
						},
						"version": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `Indicates the module version.`,
						},
						"tags": {
							Type:        schema.TypeList,
							Computed:    true,
							Elem:        &schema.Schema{Type: schema.TypeString},
							Description: `Indicates the tags.`,
						},
						"url_relative": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `Indicates the extension URL.`,
						},
						"properties_list": {
							Type:        schema.TypeList,
							Computed:    true,
							Elem:        &schema.Schema{Type: schema.TypeString},
							Description: `Indicates the properties list.`,
						},
						"manifest_version": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `Indicates the summary version.`,
						},
					},
				},
			},
		},
	}
}

func dataSourceCodeArtsPipelineModulesRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	client, err := cfg.NewServiceClient("codearts_pipeline", region)
	if err != nil {
		return diag.Errorf("error creating CodeArts Pipeline client: %s", err)
	}

	getHttpUrl := "v2/extensions/modules"
	getPath := client.Endpoint + getHttpUrl
	getOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	getPath += buildCodeArtsPipelineModulesQueryParams(d)
	getResp, err := client.Request("GET", getPath, &getOpt)
	if err != nil {
		return diag.Errorf("error retrieving pipeline modules: %s", err)
	}
	getRespBody, err := utils.FlattenResponse(getResp)
	if err != nil {
		return diag.Errorf("error flatten response: %s", err)
	}

	result := utils.PathSearch("result", getRespBody, nil)
	modules := utils.PathSearch(`"devcloud.open.endpoint"`, result, nil)
	moduleDatas := utils.PathSearch("data", modules, make([]interface{}, 0)).([]interface{})
	rst := make([]map[string]interface{}, 0, len(moduleDatas))
	for _, module := range moduleDatas {
		rst = append(rst, map[string]interface{}{
			"id":               utils.PathSearch("id", module, nil),
			"base_url":         utils.PathSearch("baseUrl", module, nil),
			"description":      utils.PathSearch("description", module, nil),
			"location":         utils.PathSearch("location", module, nil),
			"module_id":        utils.PathSearch("module_id", module, nil),
			"name":             utils.PathSearch("name", module, nil),
			"properties":       encodeIntoJson(utils.PathSearch("properties", module, nil)),
			"publisher":        utils.PathSearch("publisher", module, nil),
			"type":             utils.PathSearch("type", module, nil),
			"version":          utils.PathSearch("version", module, nil),
			"tags":             utils.PathSearch("tags", module, nil),
			"url_relative":     utils.PathSearch("url_relative", module, nil),
			"properties_list":  flattenPipelineModeulesPropertiesList(module),
			"manifest_version": utils.PathSearch("manifest_version", module, nil),
		})
	}

	id, err := uuid.GenerateUUID()
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(id)

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("modules", rst),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenPipelineModeulesPropertiesList(resp interface{}) []interface{} {
	propertiesList, ok := utils.PathSearch("properties_list", resp, make([]interface{}, 0)).([]interface{})
	if ok && len(propertiesList) > 0 {
		result := make([]interface{}, 0, len(propertiesList))
		for _, properties := range propertiesList {
			result = append(result, encodeIntoJson(properties))
		}

		return result
	}

	return nil
}

func buildCodeArtsPipelineModulesQueryParams(d *schema.ResourceData) string {
	res := "?locations=devcloud.open.endpoint"

	if v, ok := d.GetOk("name"); ok {
		res += fmt.Sprintf("&name=%v", v)
	}
	if v, ok := d.GetOk("product_line"); ok {
		res += fmt.Sprintf("&productLine=%v", v)
	}
	if v, ok := d.GetOk("tags"); ok {
		tags := v.([]interface{})
		for _, tag := range tags {
			res += fmt.Sprintf("&tags=%v", tag)
		}
	}

	return res
}
