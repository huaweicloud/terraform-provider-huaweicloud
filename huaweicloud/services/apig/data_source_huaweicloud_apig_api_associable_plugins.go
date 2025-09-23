package apig

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API APIG GET /v2/{project_id}/apigw/instances/{instance_id}/apis/{api_id}/attachable-plugins
func DataSourceApiAssociablePlugins() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceApiAssociablePluginsRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `The region where the API and the associable plugins are located.`,
			},

			// Required parameters.
			"instance_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The ID of the dedicated instance to which the API and the associable plugins belong.`,
			},
			"api_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The ID of the API to be queried.`,
			},

			// Optional parameters.
			"env_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The ID of the environment where the API is published.`,
			},
			"plugin_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The name of the associable plugins to be queried.`,
			},
			"plugin_type": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The type of the associable plugins to be queried.`,
			},
			"plugin_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The ID of the specified associable plugin to be queried.`,
			},

			// Attributes.
			"plugins": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: `The list of plugins that can be bound to the specified API.`,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The ID of the associable plugin.`,
						},
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The name of the associable plugin.`,
						},
						"type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The type of the associable plugin.`,
						},
						"scope": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The scope of the associable plugin.`,
						},
						"content": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The content of the associable plugin configuration.`,
						},
						"description": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The description of the associable plugin.`,
						},
						"create_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The creation time of the associable plugin, in RFC3339 format.`,
						},
						"update_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The update time of the associable plugin, in RFC3339 format.`,
						},
					},
				},
			},
		},
	}
}

func buildApiAssociablePluginsQueryParams(d *schema.ResourceData) string {
	res := ""

	if v, ok := d.GetOk("env_id"); ok {
		res = fmt.Sprintf("%s&env_id=%v", res, v)
	}
	if v, ok := d.GetOk("plugin_name"); ok {
		res = fmt.Sprintf("%s&plugin_name=%v", res, v)
	}
	if v, ok := d.GetOk("plugin_type"); ok {
		res = fmt.Sprintf("%s&plugin_type=%v", res, v)
	}
	if v, ok := d.GetOk("plugin_id"); ok {
		res = fmt.Sprintf("%s&plugin_id=%v", res, v)
	}

	return res
}

func listApiAssociablePlugins(client *golangsdk.ServiceClient, d *schema.ResourceData) ([]interface{}, error) {
	var (
		httpUrl = "v2/{project_id}/apigw/instances/{instance_id}/apis/{api_id}/attachable-plugins?limit={limit}"
		limit   = 100
		offset  = 0
		result  = make([]interface{}, 0)
	)

	listPath := client.Endpoint + httpUrl
	listPath = strings.ReplaceAll(listPath, "{project_id}", client.ProjectID)
	listPath = strings.ReplaceAll(listPath, "{instance_id}", d.Get("instance_id").(string))
	listPath = strings.ReplaceAll(listPath, "{api_id}", d.Get("api_id").(string))
	listPath = strings.ReplaceAll(listPath, "{limit}", strconv.Itoa(limit))
	listPath += buildApiAssociablePluginsQueryParams(d)

	opt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json",
		},
	}

	for {
		listPathWithOffset := listPath + fmt.Sprintf("&offset=%d", offset)
		requestResp, err := client.Request("GET", listPathWithOffset, &opt)
		if err != nil {
			return nil, err
		}
		respBody, err := utils.FlattenResponse(requestResp)
		if err != nil {
			return nil, err
		}
		plugins := utils.PathSearch("plugins", respBody, make([]interface{}, 0)).([]interface{})
		result = append(result, plugins...)
		if len(plugins) < limit {
			break
		}
		offset += len(plugins)
	}

	return result, nil
}

func dataSourceApiAssociablePluginsRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	client, err := cfg.NewServiceClient("apig", region)
	if err != nil {
		return diag.Errorf("error creating APIG client: %s", err)
	}

	apiId := d.Get("api_id").(string)
	resp, err := listApiAssociablePlugins(client, d)
	if err != nil {
		return diag.Errorf("error error querying associable plugins for specified API (%s): %s", apiId, err)
	}

	randomUUID, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(randomUUID)

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("plugins", flattenApiAssociablePlugins(resp)),
	)
	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenApiAssociablePlugins(plugins []interface{}) []map[string]interface{} {
	if len(plugins) < 1 {
		return nil
	}

	result := make([]map[string]interface{}, 0, len(plugins))
	for _, plugin := range plugins {
		result = append(result, map[string]interface{}{
			"id":          utils.PathSearch("plugin_id", plugin, nil),
			"name":        utils.PathSearch("plugin_name", plugin, nil),
			"type":        utils.PathSearch("plugin_type", plugin, nil),
			"scope":       utils.PathSearch("plugin_scope", plugin, nil),
			"content":     utils.PathSearch("plugin_content", plugin, nil),
			"description": utils.PathSearch("remark", plugin, nil),
			"create_time": utils.PathSearch("create_time", plugin, nil),
			"update_time": utils.PathSearch("update_time", plugin, nil),
		})
	}

	return result
}
