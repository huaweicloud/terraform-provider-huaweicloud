package cdn

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

// @API CDN GET /v1/cdn/statistics/stats-configs
func DataSourceStatisticConfiguration() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceStatisticConfigurationRead,

		Schema: map[string]*schema.Schema{
			// Required parameters.
			"config_type": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: `The configuration category.`,
			},

			// Attributes.
			"configurations": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"config_type": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: `The configuration category.`,
						},
						"resource_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The resource type.`,
						},
						"resource_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The resource name.`,
						},
						"config_info": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"url": {
										Type:        schema.TypeList,
										Computed:    true,
										Description: `The top URL statistics configuration.`,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"enable": {
													Type:        schema.TypeBool,
													Computed:    true,
													Description: `Whether the top URL statistics configuration is enabled.`,
												},
												"limit": {
													Type:        schema.TypeInt,
													Computed:    true,
													Description: `The number of top URL statistics to report.`,
												},
												"sort_by_code": {
													Type:        schema.TypeBool,
													Computed:    true,
													Description: `Whether to support reporting by status code.`,
												},
											},
										},
									},
									"ua": {
										Type:        schema.TypeList,
										Computed:    true,
										Description: `The top UA statistics configuration.`,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"enable": {
													Type:        schema.TypeBool,
													Computed:    true,
													Description: `Whether the top UA statistics configuration is enabled.`,
												},
											},
										},
									},
								},
							},
							Description: `The statistics configuration information.`,
						},
						"expired_time": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: `The expiration time of the statistics configuration, in seconds timestamp.`,
						},
					},
				},
				Description: `The list of statistic configurations.`,
			},
		},
	}
}

func listStatisticConfigurations(client *golangsdk.ServiceClient, configType int) ([]interface{}, error) {
	var (
		httpUrl = "v1/cdn/statistics/stats-configs?config_type={config_type}&limit={limit}"
		limit   = 100
		offset  = 0
		result  = make([]interface{}, 0)
	)

	listPath := client.Endpoint + httpUrl
	listPath = strings.ReplaceAll(listPath, "{config_type}", strconv.Itoa(configType))
	listPath = strings.ReplaceAll(listPath, "{limit}", strconv.Itoa(limit))

	opt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json",
		},
	}

	for {
		listPathWithMarker := fmt.Sprintf("%s&offset=%v", listPath, offset)
		requestResp, err := client.Request("GET", listPathWithMarker, &opt)
		if err != nil {
			return nil, err
		}
		respBody, err := utils.FlattenResponse(requestResp)
		if err != nil {
			return nil, err
		}
		configurations := utils.PathSearch("data", respBody, make([]interface{}, 0)).([]interface{})
		result = append(result, configurations...)
		if len(configurations) < limit {
			break
		}

		offset += len(configurations)
	}

	return result, nil
}

func flattenStatisticConfigurationConfigInfoUrl(urlConfig map[string]interface{}) []map[string]interface{} {
	if len(urlConfig) < 1 {
		return nil
	}

	return []map[string]interface{}{
		{
			"enable":       utils.PathSearch("enable", urlConfig, nil),
			"limit":        utils.PathSearch("limit", urlConfig, nil),
			"sort_by_code": utils.PathSearch("sort_by_code", urlConfig, nil),
		},
	}
}

func flattenStatisticConfigurationConfigInfoUa(uaConfig map[string]interface{}) []map[string]interface{} {
	if len(uaConfig) < 1 {
		return nil
	}

	return []map[string]interface{}{
		{
			"enable": utils.PathSearch("enable", uaConfig, nil),
		},
	}
}

func flattenStatisticConfigurationConfigInfo(configInfo map[string]interface{}) []map[string]interface{} {
	if len(configInfo) < 1 {
		return nil
	}

	return []map[string]interface{}{
		{
			"url": flattenStatisticConfigurationConfigInfoUrl(utils.PathSearch("url", configInfo,
				make(map[string]interface{})).(map[string]interface{})),
			"ua": flattenStatisticConfigurationConfigInfoUa(utils.PathSearch("ua", configInfo,
				make(map[string]interface{})).(map[string]interface{})),
		},
	}
}

func flattenStatisticConfigurations(configurations []interface{}) []map[string]interface{} {
	if len(configurations) < 1 {
		return nil
	}

	result := make([]map[string]interface{}, 0, len(configurations))
	for _, configuration := range configurations {
		result = append(result, map[string]interface{}{
			"config_type":   utils.PathSearch("config_type", configuration, nil),
			"resource_type": utils.PathSearch("resource_type", configuration, nil),
			"resource_name": utils.PathSearch("resource_name", configuration, nil),
			"config_info": flattenStatisticConfigurationConfigInfo(utils.PathSearch("config_info", configuration,
				make(map[string]interface{})).(map[string]interface{})),
			"expired_time": utils.PathSearch("expired_time", configuration, nil),
		})
	}

	return result
}

func dataSourceStatisticConfigurationRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	client, err := cfg.NewServiceClient("cdn", "")
	if err != nil {
		return diag.Errorf("error creating CDN client: %s", err)
	}

	configType := d.Get("config_type").(int)

	configurations, err := listStatisticConfigurations(client, configType)
	if err != nil {
		return diag.Errorf("error querying CDN statistic configurations: %s", err)
	}

	randomUUID, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(randomUUID)

	mErr := multierror.Append(nil,
		d.Set("configurations", flattenStatisticConfigurations(configurations)),
	)
	return diag.FromErr(mErr.ErrorOrNil())
}
