package eip

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

// @API EIP GET /v3/{domain_id}/geip/geip-pools
func DataSourceGlobalEIPPools() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceGlobalEIPPoolsRead,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"access_site": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"isp": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"ip_version": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"type": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"geip_pools": {
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
						"en_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"cn_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"isp": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"ip_version": {
							Type:     schema.TypeFloat,
							Computed: true,
						},
						"access_site": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"created_at": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"updated_at": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"allowed_bandwidth_types": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"type": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"cn_name": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"en_name": {
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

func dataSourceGlobalEIPPoolsRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	client, err := cfg.NewServiceClient("geip", region)
	if err != nil {
		return diag.Errorf("error creating GEIP client: %s", err)
	}

	getGlobalEIPPoolsHttpUrl := "v3/{domain_id}/geip/geip-pools"
	getGlobalEIPPoolsPath := client.Endpoint + getGlobalEIPPoolsHttpUrl
	getGlobalEIPPoolsPath = strings.ReplaceAll(getGlobalEIPPoolsPath, "{domain_id}", cfg.DomainID)
	getGlobalEIPPoolsOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	getGlobalEIPPoolsPath += "?limit=10"
	if name, ok := d.GetOk("name"); ok {
		getGlobalEIPPoolsPath += fmt.Sprintf("&name=%s", name)
	}
	if accessSite, ok := d.GetOk("access_site"); ok {
		getGlobalEIPPoolsPath += fmt.Sprintf("&access_site=%s", accessSite)
	}
	if isp, ok := d.GetOk("isp"); ok {
		getGlobalEIPPoolsPath += fmt.Sprintf("&isp=%s", isp)
	}
	if ipVersion, ok := d.GetOk("ip_version"); ok {
		getGlobalEIPPoolsPath += fmt.Sprintf("&ip_version=%v", ipVersion)
	}
	if poolType, ok := d.GetOk("type"); ok {
		getGlobalEIPPoolsPath += fmt.Sprintf("&type=%s", poolType)
	}

	currentTotal := 0
	getGlobalEIPPoolsPath += fmt.Sprintf("&offset=%v", currentTotal)

	results := make([]map[string]interface{}, 0)

	for {
		getGlobalEIPPoolsResp, err := client.Request("GET", getGlobalEIPPoolsPath, &getGlobalEIPPoolsOpt)
		if err != nil {
			return diag.FromErr(err)
		}
		getGlobalEIPPoolsRespBody, err := utils.FlattenResponse(getGlobalEIPPoolsResp)
		if err != nil {
			return diag.FromErr(err)
		}

		geipPools := utils.PathSearch("geip_pools", getGlobalEIPPoolsRespBody, make([]interface{}, 0)).([]interface{})
		for _, pool := range geipPools {
			results = append(results, map[string]interface{}{
				"name":        utils.PathSearch("name", pool, nil),
				"id":          utils.PathSearch("id", pool, nil),
				"en_name":     utils.PathSearch("en_name", pool, nil),
				"cn_name":     utils.PathSearch("cn_name", pool, nil),
				"isp":         utils.PathSearch("isp", pool, nil),
				"ip_version":  utils.PathSearch("ip_version", pool, float64(0)),
				"access_site": utils.PathSearch("access_site", pool, nil),
				"type":        utils.PathSearch("type", pool, nil),
				"created_at":  utils.PathSearch("created_at", pool, nil),
				"updated_at":  utils.PathSearch("updated_at", pool, nil),
				"allowed_bandwidth_types": flattenAllowedBandwidthTypes(
					utils.PathSearch("allowed_bandwidth_types", pool, make([]interface{}, 0))),
			})
		}

		// `current_count` means the number of pools in this page, and the limit of page is `10`.
		currentCount := utils.PathSearch("page_info.current_count", getGlobalEIPPoolsRespBody, 0)
		if currentCount.(float64) < 10 {
			break
		}

		currentTotal += len(geipPools)
		index := strings.Index(getGlobalEIPPoolsPath, "offset")
		getGlobalEIPPoolsPath = fmt.Sprintf("%soffset=%v", getGlobalEIPPoolsPath[:index], currentTotal)
	}

	uuid, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(uuid)

	mErr := multierror.Append(nil,
		d.Set("geip_pools", results),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenAllowedBandwidthTypes(rawParams interface{}) interface{} {
	rawArray, _ := rawParams.([]interface{})
	if len(rawArray) == 0 {
		return nil
	}
	rst := make([]map[string]interface{}, len(rawArray))
	for i, val := range rawArray {
		params := map[string]interface{}{
			"type":    utils.PathSearch("type", val, nil),
			"cn_name": utils.PathSearch("cn_name", val, nil),
			"en_name": utils.PathSearch("en_name", val, nil),
		}
		rst[i] = params
	}
	return rst
}
