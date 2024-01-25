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

// @API EIP GET /v3/{domain_id}/geip/access-sites
func DataSourceGlobalEIPAccessSites() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceGlobalEIPAccessSiteRead,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"proxy_region": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"iec_az_code": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"access_sites": {
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
						"proxy_region": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"iec_az_code": {
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
						"created_at": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"updated_at": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceGlobalEIPAccessSiteRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	client, err := cfg.NewServiceClient("geip", region)
	if err != nil {
		return diag.Errorf("error creating GEIP client: %s", err)
	}

	getGlobalEIPAccessSiteHttpUrl := "v3/{domain_id}/geip/access-sites"
	getGlobalEIPAccessSitePath := client.Endpoint + getGlobalEIPAccessSiteHttpUrl
	getGlobalEIPAccessSitePath = strings.ReplaceAll(getGlobalEIPAccessSitePath, "{domain_id}", cfg.DomainID)
	getGlobalEIPAccessSiteOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	getGlobalEIPAccessSitePath += "?limit=10"
	if name, ok := d.GetOk("name"); ok {
		getGlobalEIPAccessSitePath += fmt.Sprintf("&name=%s", name)
	}
	if proxyRegion, ok := d.GetOk("proxy_region"); ok {
		getGlobalEIPAccessSitePath += fmt.Sprintf("&proxy_region=%s", proxyRegion)
	}
	if iecAZCode, ok := d.GetOk("iec_az_code"); ok {
		getGlobalEIPAccessSitePath += fmt.Sprintf("&iec_az_code=%s", iecAZCode)
	}
	currentTotal := 0
	getGlobalEIPAccessSitePath += fmt.Sprintf("&offset=%v", currentTotal)

	results := make([]map[string]interface{}, 0)
	for {
		getGlobalEIPAccessSiteResp, err := client.Request("GET", getGlobalEIPAccessSitePath, &getGlobalEIPAccessSiteOpt)
		if err != nil {
			return diag.FromErr(err)
		}
		getGlobalEIPAccessSiteRespBody, err := utils.FlattenResponse(getGlobalEIPAccessSiteResp)
		if err != nil {
			return diag.FromErr(err)
		}

		accessSites := utils.PathSearch("access_sites", getGlobalEIPAccessSiteRespBody, make([]interface{}, 0)).([]interface{})
		for _, accessSite := range accessSites {
			results = append(results, map[string]interface{}{
				"id":           utils.PathSearch("id", accessSite, nil),
				"name":         utils.PathSearch("name", accessSite, nil),
				"en_name":      utils.PathSearch("en_name", accessSite, nil),
				"cn_name":      utils.PathSearch("cn_name", accessSite, nil),
				"proxy_region": utils.PathSearch("proxy_region", accessSite, nil),
				"iec_az_code":  utils.PathSearch("iec_az_code", accessSite, nil),
				"created_at":   utils.PathSearch("created_at", accessSite, nil),
				"updated_at":   utils.PathSearch("updated_at", accessSite, nil),
			})
		}

		// `current_count` means the number of `access_sites` in this page, and the limit of page is `10`.
		currentCount := utils.PathSearch("page_info.current_count", getGlobalEIPAccessSiteRespBody, float64(0))
		if currentCount.(float64) < 10 {
			break
		}

		currentTotal += len(accessSites)
		index := strings.Index(getGlobalEIPAccessSitePath, "offset")
		getGlobalEIPAccessSitePath = fmt.Sprintf("%soffset=%v", getGlobalEIPAccessSitePath[:index], currentTotal)
	}

	uuid, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(uuid)

	mErr := multierror.Append(nil,
		d.Set("access_sites", results),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}
