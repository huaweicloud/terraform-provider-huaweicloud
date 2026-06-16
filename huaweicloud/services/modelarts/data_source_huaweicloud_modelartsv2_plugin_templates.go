package modelarts

import (
	"context"
	"fmt"
	"strings"

	"github.com/google/uuid"
	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API ModelArts GET /v2/{project_id}/plugintemplates
func DataSourceV2PluginTemplates() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceV2PluginTemplatesRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "The region where the plugin templates are located.",
			},

			// Optional parameters.
			"template_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The template name of the plugin templates.",
			},
			"pool_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The pool name of the plugin templates.",
			},

			// Attributes.
			"plugin_templates": {
				Type:        schema.TypeList,
				Computed:    true,
				Elem:        dataSourceV2PluginTemplatesSchema(),
				Description: "The list of plugin templates that matched filter parameters.",
			},
		},
	}
}

func dataSourceV2PluginTemplatesSchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"metadata": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The metadata name of the plugin template.",
						},
						"annotations": {
							Type:        schema.TypeMap,
							Computed:    true,
							Elem:        &schema.Schema{Type: schema.TypeString},
							Description: "The metadata annotations of the plugin template.",
						},
					},
				},
				Description: "The metadata of the plugin template.",
			},
			"spec": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"optional": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "The spec optional of the plugin template.",
						},
						"type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The spec type of the plugin template.",
						},
						"logo_url": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The spec logo url of the plugin template.",
						},
						"description": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The spec description of the plugin template.",
						},
						"versions": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"version": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The version of the plugin template.",
									},
									"creation_timestamp": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The creation timestamp of the plugin template, in RFC3339 format.",
									},
									"inputs": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The inputs of the plugin template.",
									},
									"translate": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The translate of the plugin template.",
									},
								},
							},
							Description: "The versions of the plugin template.",
						},
					},
				},
				Description: "The spec of the plugin template.",
			},
		},
	}
	return &sc
}

func buildV2PluginTemplatesQueryParams(d *schema.ResourceData) string {
	res := ""

	if v, ok := d.GetOk("template_name"); ok {
		res = fmt.Sprintf("%s&templateName=%v", res, v)
	}
	if v, ok := d.GetOk("pool_name"); ok {
		res = fmt.Sprintf("%s&poolName=%v", res, v)
	}

	if res != "" {
		res = "?" + res
	}
	return res
}

func listV2PluginTemplates(client *golangsdk.ServiceClient, d *schema.ResourceData) ([]interface{}, error) {
	httpUrl := "v2/{project_id}/plugintemplates"

	listPath := client.Endpoint + httpUrl
	listPath = strings.ReplaceAll(listPath, "{project_id}", client.ProjectID)
	listPath += buildV2PluginTemplatesQueryParams(d)

	listOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json",
		},
	}

	requestResp, err := client.Request("GET", listPath, &listOpt)
	if err != nil {
		return nil, err
	}

	respBody, err := utils.FlattenResponse(requestResp)
	if err != nil {
		return nil, err
	}

	pluginTemplates := utils.PathSearch("items", respBody, make([]interface{}, 0)).([]interface{})
	result := make([]interface{}, 0, len(pluginTemplates))
	result = append(result, pluginTemplates...)

	return result, nil
}

func flattenV2PluginTemplatesMetadata(metadata interface{}) []map[string]interface{} {
	if metadata == nil {
		return nil
	}

	item := map[string]interface{}{
		"name":        utils.PathSearch("name", metadata, nil),
		"annotations": utils.PathSearch("annotations", metadata, nil),
	}

	return []map[string]interface{}{item}
}

func flattenV2PluginTemplatesVersions(versions []interface{}) []map[string]interface{} {
	if len(versions) < 1 {
		return nil
	}

	result := make([]map[string]interface{}, 0, len(versions))
	for _, version := range versions {
		result = append(result, map[string]interface{}{
			"version":            utils.PathSearch("version", version, nil),
			"creation_timestamp": utils.PathSearch("creation_timestamp", version, nil),
			"inputs":             utils.JsonToString(utils.PathSearch("inputs", version, nil)),
			"translate":          utils.JsonToString(utils.PathSearch("translate", version, nil)),
		})
	}

	return result
}

func flattenV2PluginTemplatesSpec(spec interface{}) []map[string]interface{} {
	if spec == nil {
		return nil
	}

	item := map[string]interface{}{
		"optional":    utils.PathSearch("optional", spec, nil),
		"type":        utils.PathSearch("type", spec, nil),
		"logo_url":    utils.PathSearch("logo_url", spec, nil),
		"description": utils.PathSearch("description", spec, nil),
		"versions": flattenV2PluginTemplatesVersions(
			utils.PathSearch("versions", spec, make([]interface{}, 0)).([]interface{})),
	}

	return []map[string]interface{}{item}
}

func flattenV2PluginTemplates(pluginTemplates []interface{}) []map[string]interface{} {
	if len(pluginTemplates) < 1 {
		return nil
	}

	result := make([]map[string]interface{}, 0, len(pluginTemplates))
	for _, pluginTemplate := range pluginTemplates {
		result = append(result, map[string]interface{}{
			"metadata": utils.ValueIgnoreEmpty(flattenV2PluginTemplatesMetadata(
				utils.PathSearch("metadata", pluginTemplate, nil))),
			"spec": utils.ValueIgnoreEmpty(flattenV2PluginTemplatesSpec(
				utils.PathSearch("spec", pluginTemplate, nil))),
		})
	}
	return result
}

func dataSourceV2PluginTemplatesRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg    = meta.(*config.Config)
		region = cfg.GetRegion(d)
	)

	client, err := cfg.NewServiceClient("modelarts", region)
	if err != nil {
		return diag.Errorf("error creating ModelArts client: %s", err)
	}

	pluginTemplates, err := listV2PluginTemplates(client, d)
	if err != nil {
		return diag.Errorf("error querying ModelArts Plugin templates: %s", err)
	}

	randUUID, err := uuid.NewRandom()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(randUUID.String())

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("plugin_templates", flattenV2PluginTemplates(pluginTemplates)),
	)
	return diag.FromErr(mErr.ErrorOrNil())
}
