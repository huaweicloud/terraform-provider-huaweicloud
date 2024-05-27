package apig

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

// @API APIG GET /v2/{project_id}/apigw/instances/{instance_id}/app-auths/binded-apps
func DataSourceApiAssociatedPlugins() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceApiAssociatedPluginsRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"instance_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The ID of the dedicated instance to which the plugins belong.`,
			},
			"api_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The ID of the API bound to the plugin.`,
			},
			"plugin_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The ID of the plugin.`,
			},
			"name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The name of the plugin.`,
			},
			"type": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The type of the plugin.`,
			},
			"env_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The ID of the environment where the API is published.`,
			},
			"env_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The name of the environment where the API is published.`,
			},
			"plugins": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: `All plugins that match the filter parameters.`,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The ID of the plugin.`,
						},
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The name of the plugin.`,
						},
						"type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The type of the plugin.`,
						},
						"description": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The description of the plugin.`,
						},
						"content": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The configuration details for the plugin.`,
						},
						"env_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The ID of the environment where the API is published.`,
						},
						"env_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The name of the environment where the API is published.`,
						},
						"bind_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The bind ID.`,
						},
						"bind_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The time that the plugin is bound to the API.`,
						},
					},
				},
			},
		},
	}
}

func buildListApiAssociatedPluginsParams(d *schema.ResourceData) string {
	res := ""
	if v, ok := d.GetOk("env_id"); ok {
		res = fmt.Sprintf("%s&env_id=%v", res, v)
	}
	if v, ok := d.GetOk("env_name"); ok {
		res = fmt.Sprintf("%s&env_name=%v", res, v)
	}
	if v, ok := d.GetOk("plugin_id"); ok {
		res = fmt.Sprintf("%s&plugin_id=%v", res, v)
	}
	if v, ok := d.GetOk("name"); ok {
		res = fmt.Sprintf("%s&plugin_name=%v", res, v)
	}
	if v, ok := d.GetOk("type"); ok {
		res = fmt.Sprintf("%s&plugin_type=%v", res, v)
	}
	return res
}

func queryApiAssociatedPlugins(client *golangsdk.ServiceClient, d *schema.ResourceData) ([]interface{}, error) {
	var (
		httpUrl    = "v2/{project_id}/apigw/instances/{instance_id}/apis/{api_id}/attached-plugins?limit=100"
		instanceId = d.Get("instance_id").(string)
		apiId      = d.Get("api_id").(string)
		offset     = 0
		result     = make([]interface{}, 0)
	)

	listPath := client.Endpoint + httpUrl
	listPath = strings.ReplaceAll(listPath, "{project_id}", client.ProjectID)
	listPath = strings.ReplaceAll(listPath, "{instance_id}", instanceId)
	listPath = strings.ReplaceAll(listPath, "{api_id}", apiId)

	queryParams := buildListApiAssociatedPluginsParams(d)
	listPath += queryParams

	opt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	for {
		listPathWithOffset := fmt.Sprintf("%s&offset=%d", listPath, offset)
		requestResp, err := client.Request("GET", listPathWithOffset, &opt)
		if err != nil {
			return nil, fmt.Errorf("error retrieving associated plugins (bound to the API: %s) under specified "+
				"dedicated instance (%s): %s", apiId, instanceId, err)
		}
		respBody, err := utils.FlattenResponse(requestResp)
		if err != nil {
			return nil, err
		}
		plugins := utils.PathSearch("plugins", respBody, make([]interface{}, 0)).([]interface{})
		if len(plugins) < 1 {
			break
		}
		result = append(result, plugins...)
		offset += len(plugins)
	}
	return result, nil
}

func dataSourceApiAssociatedPluginsRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg    = meta.(*config.Config)
		region = cfg.GetRegion(d)
	)

	client, err := cfg.NewServiceClient("apig", region)
	if err != nil {
		return diag.Errorf("error creating APIG client: %s", err)
	}
	plugins, err := queryApiAssociatedPlugins(client, d)
	if err != nil {
		return diag.FromErr(err)
	}

	dataSourceId, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(dataSourceId)

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("plugins", flattenAssociatedPlugins(plugins)),
	)
	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenAssociatedPlugins(plugins []interface{}) []interface{} {
	if len(plugins) < 1 {
		return nil
	}

	result := make([]interface{}, 0, len(plugins))
	for _, plugin := range plugins {
		result = append(result, map[string]interface{}{
			"id":          utils.PathSearch("plugin_id", plugin, nil),
			"name":        utils.PathSearch("plugin_name", plugin, nil),
			"type":        utils.PathSearch("plugin_type", plugin, nil),
			"description": utils.PathSearch("remark", plugin, nil),
			"content":     utils.PathSearch("plugin_content", plugin, nil),
			"env_id":      utils.PathSearch("env_id", plugin, nil),
			"env_name":    utils.PathSearch("env_name", plugin, nil),
			"bind_id":     utils.PathSearch("plugin_attach_id", plugin, nil),
			"bind_time":   utils.PathSearch("attached_time", plugin, nil),
		})
	}
	return result
}
