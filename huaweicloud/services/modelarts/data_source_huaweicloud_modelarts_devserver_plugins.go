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

// @API ModelArts GET /v1/{project_id}/dev-servers/plugins
func DataSourceDevServerPlugins() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceDevServerPluginsRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `The region where the DevServer plugins are located.`,
			},

			// Optional parameters.
			"name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The name of the DevServer plugin to be queried.`,
			},

			// Attributes.
			"plugins": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: `The list of DevServer plugins that matched filter parameters.`,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The name of the plugin.`,
						},
						"infos": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: `The list of plugin details.`,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"id": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: `The ID of the plugin.`,
									},
									"version": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: `The version of the plugin.`,
									},
									"status": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: `The status of the plugin.`,
									},
									"url": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: `The download URL of the plugin.`,
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

func buildDevServerPluginsQueryParams(d *schema.ResourceData) string {
	res := ""
	if v, ok := d.GetOk("name"); ok {
		res = fmt.Sprintf("%s&name=%v", res, v)
	}
	if res != "" {
		res = "?" + res[1:]
	}
	return res
}

func listDevServerPlugins(client *golangsdk.ServiceClient, d *schema.ResourceData) ([]interface{}, error) {
	httpUrl := "v1/{project_id}/dev-servers/plugins"
	listPath := client.Endpoint + httpUrl
	listPath = strings.ReplaceAll(listPath, "{project_id}", client.ProjectID)
	listPath += buildDevServerPluginsQueryParams(d)

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

	return utils.PathSearch("data", respBody, make([]interface{}, 0)).([]interface{}), nil
}

func flattenDevServerPluginInfos(infos []interface{}) []map[string]interface{} {
	if len(infos) < 1 {
		return nil
	}

	result := make([]map[string]interface{}, 0, len(infos))
	for _, info := range infos {
		result = append(result, map[string]interface{}{
			"id":      utils.PathSearch("id", info, nil),
			"version": utils.PathSearch("version", info, nil),
			"status":  utils.PathSearch("status", info, nil),
			"url":     utils.PathSearch("url", info, nil),
		})
	}

	return result
}

func flattenDevServerPlugins(plugins []interface{}) []map[string]interface{} {
	if len(plugins) < 1 {
		return nil
	}

	result := make([]map[string]interface{}, 0, len(plugins))
	for _, plugin := range plugins {
		result = append(result, map[string]interface{}{
			"name": utils.PathSearch("name", plugin, nil),
			"infos": flattenDevServerPluginInfos(utils.PathSearch("infos", plugin,
				make([]interface{}, 0)).([]interface{})),
		})
	}

	return result
}

func dataSourceDevServerPluginsRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	client, err := cfg.NewServiceClient("modelarts", region)
	if err != nil {
		return diag.Errorf("error creating ModelArts client: %s", err)
	}

	plugins, err := listDevServerPlugins(client, d)
	if err != nil {
		return diag.Errorf("error querying DevServer plugins: %s", err)
	}

	randomUUID, err := uuid.NewRandom()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(randomUUID.String())

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("plugins", flattenDevServerPlugins(plugins)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}
