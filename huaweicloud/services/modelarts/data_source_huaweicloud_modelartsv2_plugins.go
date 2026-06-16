package modelarts

import (
	"context"
	"strings"

	"github.com/google/uuid"
	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API ModelArts GET /v2/{project_id}/pools/{pool_name}/plugins
func DataSourceV2Plugins() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceV2PluginsRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `The region where the plugins are located.`,
			},

			// Required parameters
			"pool_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The ID of the resource pool.`,
			},

			// Attributes
			"plugins": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: `The list of the plugins that matched filter parameters.`,
				Elem:        v2PluginElemSchema(),
			},
		},
	}
}

func v2PluginElemSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"api_version": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The API version of the plugin.`,
			},
			"kind": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The type of the plugin instance.`,
			},
			"metadata": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: `The metadata information of the plugin.`,
				Elem:        v2PluginMetadataElemSchema(),
			},
			"spec": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: `The spec information of the plugin.`,
				Elem:        v2PluginSpecElemSchema(),
			},
			"status": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: `The status information of the plugin.`,
				Elem:        v2PluginStatusElemSchema(),
			},
		},
	}
}

func v2PluginMetadataElemSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The name of the plugin instance.`,
			},
			"created_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The creation time of the plugin instance.`,
			},
		},
	}
}

func v2PluginSpecElemSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"template": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: `The template information of the plugin.`,
				Elem:        v2PluginTemplateElemSchema(),
			},
		},
	}
}

func v2PluginTemplateElemSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The name of the plugin template.`,
			},
			"version": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The version of the plugin template.`,
			},
			"inputs": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The installation parameters of the plugin template.`,
			},
		},
	}
}

func v2PluginStatusElemSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"phase": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The status of the plugin instance.`,
			},
			"version": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The version of the plugin instance.`,
			},
			"reason": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The reason for the plugin instance installation failure.`,
			},
			"values": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The installation parameters of the plugin instance.`,
			},
			"resources": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: `The resources occupied by the plugin instance.`,
				Elem:        v2PluginResourcesElemSchema(),
			},
		},
	}
}

func v2PluginResourcesElemSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"involved_object": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: `The referenced resource object of the plugin.`,
				Elem:        v2PluginInvolvedObjectElemSchema(),
			},
			"replicas": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: `The number of replicas of the resource object.`,
			},
			"limits": {
				Type:        schema.TypeMap,
				Computed:    true,
				Description: `The resource limits of the plugin.`,
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
			"requests": {
				Type:        schema.TypeMap,
				Computed:    true,
				Description: `The resource requests of the plugin.`,
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
		},
	}
}

func v2PluginInvolvedObjectElemSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"kind": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The API type of the resource object.`,
			},
			"api_version": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The API version of the resource object.`,
			},
			"namespace": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The namespace of the resource object.`,
			},
			"name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The name of the resource object.`,
			},
			"uid": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The unique identifier of the resource object.`,
			},
			"resource_version": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The current version of the resource object.`,
			},
		},
	}
}

func listV2Plugins(client *golangsdk.ServiceClient, poolId string) ([]interface{}, error) {
	httpUrl := "v2/{project_id}/pools/{pool_name}/plugins"
	listPath := client.Endpoint + httpUrl
	listPath = strings.ReplaceAll(listPath, "{project_id}", client.ProjectID)
	listPath = strings.ReplaceAll(listPath, "{pool_name}", poolId)

	opt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json",
		},
	}

	requestResp, err := client.Request("GET", listPath, &opt)
	if err != nil {
		return nil, err
	}
	respBody, err := utils.FlattenResponse(requestResp)
	if err != nil {
		return nil, err
	}

	return utils.PathSearch("items", respBody, make([]interface{}, 0)).([]interface{}), nil
}

func flattenV2PluginObjectReference(objRef map[string]interface{}) []map[string]interface{} {
	if len(objRef) < 1 {
		return nil
	}

	return []map[string]interface{}{
		{
			"kind":             utils.PathSearch("kind", objRef, nil),
			"api_version":      utils.PathSearch("apiVersion", objRef, nil),
			"namespace":        utils.PathSearch("namespace", objRef, nil),
			"name":             utils.PathSearch("name", objRef, nil),
			"uid":              utils.PathSearch("uid", objRef, nil),
			"resource_version": utils.PathSearch("resourceVersion", objRef, nil),
		},
	}
}

func flattenV2PluginTemplate(template map[string]interface{}) []map[string]interface{} {
	if len(template) < 1 {
		return nil
	}

	return []map[string]interface{}{
		{
			"name":    utils.PathSearch("name", template, nil),
			"version": utils.PathSearch("version", template, nil),
			"inputs":  utils.JsonToString(utils.PathSearch("inputs", template, nil)),
		},
	}
}

func flattenV2PluginResources(resources []interface{}) []map[string]interface{} {
	if len(resources) < 1 {
		return nil
	}

	result := make([]map[string]interface{}, 0, len(resources))
	for _, resource := range resources {
		result = append(result, map[string]interface{}{
			"involved_object": flattenV2PluginObjectReference(
				utils.PathSearch("involvedObject", resource, make(map[string]interface{})).(map[string]interface{})),
			"replicas": utils.PathSearch("replicas", resource, nil),
			"limits":   utils.PathSearch("limits", resource, make(map[string]interface{})),
			"requests": utils.PathSearch("requests", resource, make(map[string]interface{})),
		})
	}

	return result
}

func flattenV2PluginMetadata(metadata map[string]interface{}) []map[string]interface{} {
	if len(metadata) < 1 {
		return nil
	}

	return []map[string]interface{}{
		{
			"name":       utils.PathSearch("name", metadata, nil),
			"created_at": utils.PathSearch("creationTimestamp", metadata, nil),
		},
	}
}

func flattenV2PluginSpec(spec map[string]interface{}) []map[string]interface{} {
	if len(spec) < 1 {
		return nil
	}

	return []map[string]interface{}{
		{
			"template": flattenV2PluginTemplate(
				utils.PathSearch("template", spec, make(map[string]interface{})).(map[string]interface{})),
		},
	}
}

func flattenV2PluginStatus(status map[string]interface{}) []map[string]interface{} {
	if len(status) < 1 {
		return nil
	}

	return []map[string]interface{}{
		{
			"phase":   utils.PathSearch("phase", status, nil),
			"version": utils.PathSearch("version", status, nil),
			"reason":  utils.PathSearch("reason", status, nil),
			"values":  utils.PathSearch("values", status, nil),
			"resources": flattenV2PluginResources(
				utils.PathSearch("resources", status, make([]interface{}, 0)).([]interface{})),
		},
	}
}

func flattenV2Plugins(plugins []interface{}) []map[string]interface{} {
	if len(plugins) < 1 {
		return nil
	}

	result := make([]map[string]interface{}, 0, len(plugins))
	for _, plugin := range plugins {
		result = append(result, map[string]interface{}{
			"api_version": utils.PathSearch("apiVersion", plugin, nil),
			"kind":        utils.PathSearch("kind", plugin, nil),
			"metadata": flattenV2PluginMetadata(
				utils.PathSearch("metadata", plugin, make(map[string]interface{})).(map[string]interface{})),
			"spec": flattenV2PluginSpec(
				utils.PathSearch("spec", plugin, make(map[string]interface{})).(map[string]interface{})),
			"status": flattenV2PluginStatus(
				utils.PathSearch("status", plugin, make(map[string]interface{})).(map[string]interface{})),
		})
	}

	return result
}

func dataSourceV2PluginsRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	client, err := cfg.NewServiceClient("modelarts", region)
	if err != nil {
		return diag.Errorf("error creating ModelArts client: %s", err)
	}

	poolId := d.Get("pool_id").(string)
	plugins, err := listV2Plugins(client, poolId)
	if err != nil {
		return diag.Errorf("error querying pool plugins: %s", err)
	}

	randomUUID, err := uuid.NewRandom()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(randomUUID.String())

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("plugins", flattenV2Plugins(plugins)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}
