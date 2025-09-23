package iotda

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

// @API IoTDA GET /v5/iot/{project_id}/device-proxies
func DataSourceDeviceProxies() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceDeviceProxiesRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"space_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"proxies": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"space_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"effective_time_range": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"start_time": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"end_time": {
										Type:     schema.TypeString,
										Computed: true,
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

func buildDeviceProxiesQueryParams(d *schema.ResourceData) string {
	queryParams := "?limit=50"
	if v, ok := d.GetOk("space_id"); ok {
		queryParams = fmt.Sprintf("%s&app_id=%v", queryParams, v)
	}
	if v, ok := d.GetOk("name"); ok {
		queryParams = fmt.Sprintf("%s&proxy_name=%v", queryParams, v)
	}

	return queryParams
}

func dataSourceDeviceProxiesRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg       = meta.(*config.Config)
		region    = cfg.GetRegion(d)
		isDerived = WithDerivedAuth(cfg, region)
		product   = "iotda"
		httpUrl   = "v5/iot/{project_id}/device-proxies"
		offset    = 0
		result    = make([]interface{}, 0)
	)

	client, err := cfg.NewServiceClientWithDerivedAuth(product, region, isDerived)
	if err != nil {
		return diag.Errorf("error creating IoTDA client: %s", err)
	}

	getPath := client.Endpoint + httpUrl
	getPath = strings.ReplaceAll(getPath, "{project_id}", client.ProjectID)
	getPath += buildDeviceProxiesQueryParams(d)
	getOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	for {
		currentPath := fmt.Sprintf("%s&offset=%v", getPath, offset)
		getResp, err := client.Request("GET", currentPath, &getOpt)
		if err != nil {
			return diag.Errorf("error retrieving IoTDA device proxies: %s", err)
		}

		getRespBody, err := utils.FlattenResponse(getResp)
		if err != nil {
			return diag.FromErr(err)
		}

		proxiesResp := utils.PathSearch("device_proxies", getRespBody, make([]interface{}, 0)).([]interface{})
		if len(proxiesResp) == 0 {
			break
		}

		result = append(result, proxiesResp...)
		offset += len(proxiesResp)
	}

	generateUUID, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}

	d.SetId(generateUUID)

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("proxies", flattenDeviceProxies(result)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenDeviceProxies(proxiesResp []interface{}) []interface{} {
	rst := make([]interface{}, 0, len(proxiesResp))
	for _, v := range proxiesResp {
		rst = append(rst, map[string]interface{}{
			"id":                   utils.PathSearch("proxy_id", v, nil),
			"name":                 utils.PathSearch("proxy_name", v, nil),
			"space_id":             utils.PathSearch("app_id", v, nil),
			"effective_time_range": flattenEffectiveTimeRange(utils.PathSearch("effective_time_range", v, nil)),
		})
	}

	return rst
}
