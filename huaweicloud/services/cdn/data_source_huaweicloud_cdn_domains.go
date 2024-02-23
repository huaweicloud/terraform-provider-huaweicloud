package cdn

import (
	"context"
	"fmt"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API CDN GET /v1.0/cdn/domains
func DataSourceCdnDomains() *schema.Resource {
	return &schema.Resource{
		ReadContext: datasourceDomainsRead,

		Schema: map[string]*schema.Schema{
			"domain_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"type": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"domain_status": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"service_area": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"enterprise_project_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"domains": {
				Type:     schema.TypeList,
				Elem:     cdnDomainSchema(),
				Computed: true,
			},
		},
	}
}

func cdnDomainSchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"domain_status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"cname": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"sources": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     sourceSchema(),
			},
			"domain_origin_host": {
				Type:     schema.TypeMap,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Computed: true,
			},
			"https_status": {
				Type:     schema.TypeInt,
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
			"disabled": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"locked": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"auto_refresh_preheat": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"service_area": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"range_based_retrieval_enabled": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"follow_status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"origin_status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"banned_reason": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"locked_reason": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"enterprise_project_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"tags": common.TagsComputedSchema(),
		},
	}
	return &sc
}

func sourceSchema() *schema.Resource {
	sc := &schema.Resource{
		Schema: map[string]*schema.Schema{
			"origin": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"origin_type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"active": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"obs_web_hosting_enabled": {
				Type:     schema.TypeBool,
				Computed: true,
			},
		},
	}

	return sc
}

func datasourceDomainsRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var mErr *multierror.Error

	var (
		getCDNDomainsHttpUrl = "v1.0/cdn/domains"
		getCDNDomainsProduct = "cdn"
	)
	getCDNDomainsClient, err := cfg.NewServiceClient(getCDNDomainsProduct, region)
	if err != nil {
		return diag.Errorf("error creating CDN client: %s", err)
	}

	getCDNDomainsPath := getCDNDomainsClient.Endpoint + getCDNDomainsHttpUrl + "?page_size=10"
	epsId := cfg.GetEnterpriseProjectID(d)
	if epsId != "" {
		getCDNDomainsPath += fmt.Sprintf("&enterprise_project_id=%v", epsId)
	}
	if v, ok := d.GetOk("name"); ok {
		getCDNDomainsPath += fmt.Sprintf("&domain_name=%v", v)
	}
	if v, ok := d.GetOk("type"); ok {
		getCDNDomainsPath += fmt.Sprintf("&business_type=%v", v)
	}
	if v, ok := d.GetOk("domain_status"); ok {
		getCDNDomainsPath += fmt.Sprintf("&domain_status=%v", v)
	}
	if v, ok := d.GetOk("service_area"); ok {
		getCDNDomainsPath += fmt.Sprintf("&service_area=%v", v)
	}
	getCDNDomainsOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	currentTotal := 1
	rst := make([]interface{}, 0)
	for {
		currentPath := fmt.Sprintf("%s&page_number=%v", getCDNDomainsPath, currentTotal)
		getCDNDomainsResp, err := getCDNDomainsClient.Request("GET", currentPath, &getCDNDomainsOpt)
		if err != nil {
			return diag.FromErr(err)
		}
		getCDNDomainsRespBody, err := utils.FlattenResponse(getCDNDomainsResp)
		if err != nil {
			return diag.FromErr(err)
		}
		domains := utils.PathSearch("domains", getCDNDomainsRespBody, make([]interface{}, 0)).([]interface{})
		if len(domains) == 0 {
			break
		}
		rst = append(rst, domains...)
		currentTotal++
	}
	dataSourceId, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(dataSourceId)

	mErr = multierror.Append(
		mErr,
		d.Set("domains", filterListDomainsBody(flattenListDomainsBody(rst), d)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenListDomainsBody(domains []interface{}) []interface{} {
	rst := make([]interface{}, 0)
	for _, v := range domains {
		createTime := utils.PathSearch("create_time", v, 0)
		updateTime := utils.PathSearch("modify_time", v, 0)
		rst = append(rst, map[string]interface{}{
			"id":                            utils.PathSearch("id", v, nil),
			"name":                          utils.PathSearch("domain_name", v, nil),
			"type":                          utils.PathSearch("business_type", v, nil),
			"domain_status":                 utils.PathSearch("domain_status", v, nil),
			"cname":                         utils.PathSearch("cname", v, nil),
			"sources":                       flattenSources(v),
			"domain_origin_host":            utils.PathSearch("domain_origin_host", v, nil),
			"https_status":                  utils.PathSearch("https_status", v, nil),
			"created_at":                    utils.FormatTimeStampRFC3339(int64(createTime.(float64))/1000, false),
			"updated_at":                    utils.FormatTimeStampRFC3339(int64(updateTime.(float64))/1000, false),
			"disabled":                      utils.PathSearch("disabled", v, nil),
			"locked":                        utils.PathSearch("locked", v, nil),
			"auto_refresh_preheat":          utils.PathSearch("auto_refresh_preheat", v, nil),
			"service_area":                  utils.PathSearch("service_area", v, nil),
			"range_based_retrieval_enabled": converseRangeStatusToBool(utils.PathSearch("range_status", v, "").(string)),
			"follow_status":                 utils.PathSearch("follow_status", v, nil),
			"origin_status":                 utils.PathSearch("origin_status", v, nil),
			"banned_reason":                 utils.PathSearch("banned_reason", v, nil),
			"locked_reason":                 utils.PathSearch("locked_reason", v, nil),
			"enterprise_project_id":         utils.PathSearch("enterprise_project_id", v, nil),
			"tags":                          utils.FlattenTagsToMap(utils.PathSearch("tags", v, make(map[string]interface{}, 0))),
		})
	}
	return rst
}

func converseRangeStatusToBool(status interface{}) bool {
	return status == "on"
}

func converseOBSWebHostStatusToBool(status interface{}) bool {
	return status == 1
}

func filterListDomainsBody(all []interface{}, d *schema.ResourceData) []interface{} {
	rst := make([]interface{}, 0, len(all))
	for _, v := range all {
		if param, ok := d.GetOk("domain_id"); ok &&
			fmt.Sprint(param) != fmt.Sprint(utils.PathSearch("id", v, nil)) {
			continue
		}
		rst = append(rst, v)
	}
	return rst
}

func flattenSources(resp interface{}) []interface{} {
	if resp == nil {
		return nil
	}
	curJson := utils.PathSearch("sources", resp, make([]interface{}, 0))
	curArray := curJson.([]interface{})
	rst := make([]interface{}, 0, len(curArray))
	for _, v := range curArray {
		rst = append(rst, map[string]interface{}{
			"origin":                  utils.PathSearch("ip_or_domain", v, nil),
			"origin_type":             utils.PathSearch("origin_type", v, nil),
			"active":                  utils.PathSearch("active_standby", v, nil),
			"obs_web_hosting_enabled": converseOBSWebHostStatusToBool(utils.PathSearch("enable_obs_web_hosting", v, nil)),
		})
	}
	return rst
}
