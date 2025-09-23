package apig

import (
	"context"
	"fmt"
	"net/url"
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

// @API APIG GET /v2/{project_id}/apigw/instances/{instance_id}/plugins/{plugin_id}/attachable-apis
func DataSourcePluginAssociableApis() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourcePluginAssociableApisRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `The region where the plugin and the associable APIs are located.`,
			},

			// Required parameters.
			"instance_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The ID of the dedicated instance to which the plugin and the associable APIs belong.`,
			},
			"plugin_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The ID of the plugin to be queried.`,
			},
			"env_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The ID of the environment where the associable APIs are published.`,
			},

			// Optional parameters.
			"api_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The name of the associable APIs to be queried.`,
			},
			"api_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The ID of the specified associable API to be queried.`,
			},
			"group_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The ID of the API group to be queried to which the associable APIs belong.`,
			},
			"req_method": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The request method of the associable APIs to be queried.`,
			},
			"tags": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The tags of the associable APIs to be queried.`,
			},

			// Attributes.
			"apis": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: `The list of APIs that can be bound to the specified plugin.`,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The ID of the associable API.`,
						},
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The name of the associable API.`,
						},
						"type": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: `The type of the associable API.`,
						},
						"req_protocol": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The request protocol of the associable API.`,
						},
						"req_method": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The request method of the associable API.`,
						},
						"req_uri": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The request path of the associable API.`,
						},
						"auth_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The authentication type of the associable API.`,
						},
						"match_mode": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The match mode of the associable API.`,
						},
						"description": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The description of the associable API.`,
						},
						"group_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The ID of the API group to which the associable API belongs.`,
						},
						"group_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The name of the API group to which the associable API belongs.`,
						},
						"tags": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: `The tag list bound to the associable API.`,
							Elem:        &schema.Schema{Type: schema.TypeString},
						},
					},
				},
			},
		},
	}
}

func buildPluginAssociableApisQueryParams(d *schema.ResourceData) string {
	res := "&env_id={env_id}"
	res = strings.ReplaceAll(res, "{env_id}", d.Get("env_id").(string))

	if v, ok := d.GetOk("api_name"); ok {
		res = fmt.Sprintf("%s&api_name=%v", res, v)
	}
	if v, ok := d.GetOk("api_id"); ok {
		res = fmt.Sprintf("%s&api_id=%v", res, v)
	}
	if v, ok := d.GetOk("group_id"); ok {
		res = fmt.Sprintf("%s&group_id=%v", res, v)
	}
	if v, ok := d.GetOk("req_method"); ok {
		res = fmt.Sprintf("%s&req_method=%v", res, v)
	}
	if v, ok := d.GetOk("tags"); ok {
		encodedTags := url.QueryEscape(v.(string))
		res = fmt.Sprintf("%s&tags=%v", res, encodedTags)
	}

	return res
}

func listPluginAssociableApis(client *golangsdk.ServiceClient, d *schema.ResourceData) ([]interface{}, error) {
	var (
		httpUrl = "v2/{project_id}/apigw/instances/{instance_id}/plugins/{plugin_id}/attachable-apis?limit={limit}"
		limit   = 100
		offset  = 0
		result  = make([]interface{}, 0)
	)

	listPath := client.Endpoint + httpUrl
	listPath = strings.ReplaceAll(listPath, "{project_id}", client.ProjectID)
	listPath = strings.ReplaceAll(listPath, "{instance_id}", d.Get("instance_id").(string))
	listPath = strings.ReplaceAll(listPath, "{plugin_id}", d.Get("plugin_id").(string))
	listPath = strings.ReplaceAll(listPath, "{limit}", strconv.Itoa(limit))
	listPath += buildPluginAssociableApisQueryParams(d)

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
		apis := utils.PathSearch("apis", respBody, make([]interface{}, 0)).([]interface{})
		result = append(result, apis...)
		if len(apis) < limit {
			break
		}
		offset += len(apis)
	}

	return result, nil
}

func flattenPluginAssociableApis(items []interface{}) []map[string]interface{} {
	if len(items) < 1 {
		return nil
	}

	result := make([]map[string]interface{}, 0, len(items))
	for _, item := range items {
		result = append(result, map[string]interface{}{
			"id":           utils.PathSearch("api_id", item, nil),
			"name":         utils.PathSearch("api_name", item, nil),
			"type":         utils.PathSearch("type", item, nil),
			"req_protocol": utils.PathSearch("req_protocol", item, nil),
			"req_method":   utils.PathSearch("req_method", item, nil),
			"req_uri":      utils.PathSearch("req_uri", item, nil),
			"auth_type":    utils.PathSearch("auth_type", item, nil),
			"match_mode":   utils.PathSearch("match_mode", item, nil),
			"description":  utils.PathSearch("remark", item, nil),
			"group_id":     utils.PathSearch("group_id", item, nil),
			"group_name":   utils.PathSearch("group_name", item, nil),
			"tags":         utils.PathSearch("tags", item, make([]interface{}, 0)),
		})
	}

	return result
}

func dataSourcePluginAssociableApisRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	client, err := cfg.NewServiceClient("apig", region)
	if err != nil {
		return diag.Errorf("error creating APIG client: %s", err)
	}

	pluginId := d.Get("plugin_id").(string)
	resp, err := listPluginAssociableApis(client, d)
	if err != nil {
		return diag.Errorf("error querying associable APIs for specified plugin (%s): %s", pluginId, err)
	}

	randomUUID, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(randomUUID)

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("apis", flattenPluginAssociableApis(resp)),
	)
	return diag.FromErr(mErr.ErrorOrNil())
}
