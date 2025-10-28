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

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API CDN GET /v1.0/cdn/domains
func DataSourceDomains() *schema.Resource {
	return &schema.Resource{
		ReadContext: datasourceDomainsRead,

		Schema: map[string]*schema.Schema{
			// Optional parameters
			"domain_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The ID of the domain.`,
			},
			"name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The name of the domain.`,
			},
			"type": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The business type of the domain.`,
			},
			"domain_status": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The status of the domain.`,
			},
			"service_area": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The accelerated coverage area for the domain.`,
			},
			"enterprise_project_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The ID of the enterprise project to which the resource belongs.`,
			},

			// Attributes
			"domains": {
				Type:        schema.TypeList,
				Computed:    true,
				Elem:        domainSchema(),
				Description: `The list of domains that matched filter parameters.`,
			},
		},
	}
}

func domainSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The ID of the domain.`,
			},
			"name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The name of the domain.`,
			},
			"type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The business type of the domain.`,
			},
			"domain_status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The status of the domain.`,
			},
			"cname": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The CNAME of the domain.`,
			},
			"sources": {
				Type:        schema.TypeList,
				Computed:    true,
				Elem:        sourceSchema(),
				Description: `An array of one or more objects specifies the domain of the origin server.`,
			},
			"domain_origin_host": {
				Type:        schema.TypeMap,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: `The back-to-origin HOST configuration of the domain.`,
			},
			"https_status": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: `Whether the HTTPS protocol is enabled.`,
			},
			"created_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The creation time, in RFC3339 format.`,
			},
			"updated_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The update time, in RFC3339 format.`,
			},
			"disabled": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: `The ban status of the domain.`,
			},
			"locked": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: `The lock status of the domain.`,
			},
			"auto_refresh_preheat": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: `Whether to automatically refresh preheating.`,
			},
			"service_area": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The accelerated coverage area for the domain.`,
			},
			"range_based_retrieval_enabled": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: `Whether to enable range-based retrieval.`,
			},
			"follow_status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The status of back-to-source following.`,
			},
			"origin_status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Whether to pause origin site return to origin.`,
			},
			"banned_reason": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The reason why the domain was banned.`,
			},
			"locked_reason": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The reason why the domain was locked.`,
			},
			"enterprise_project_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The ID of the enterprise project to which the resource belongs.`,
			},
			"tags": common.TagsComputedSchema(),
		},
	}
}

func sourceSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"origin": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The domain name or IP address of the origin server.`,
			},
			"origin_type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The origin server type.`,
			},
			"active": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: `Whether an origin server is active or standby.`,
			},
			"obs_web_hosting_enabled": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: `Whether to enable static website hosting for the OBS bucket.`,
			},
		},
	}
}

func buildDomainsQueryParams(cfg *config.Config, d *schema.ResourceData) string {
	res := ""

	if epsId := cfg.GetEnterpriseProjectID(d); epsId != "" {
		res = fmt.Sprintf("%s&enterprise_project_id=%v", res, epsId)
	}
	if v, ok := d.GetOk("name"); ok {
		res = fmt.Sprintf("%s&domain_name=%v", res, v)
	}
	if v, ok := d.GetOk("type"); ok {
		res = fmt.Sprintf("%s&business_type=%v", res, v)
	}
	if v, ok := d.GetOk("domain_status"); ok {
		res = fmt.Sprintf("%s&domain_status=%v", res, v)
	}
	if v, ok := d.GetOk("service_area"); ok {
		res = fmt.Sprintf("%s&service_area=%v", res, v)
	}

	return res
}

func converseRangeStatusToBool(status interface{}) bool {
	return status == "on"
}

func converseOBSWebHostStatusToBool(status interface{}) bool {
	return status == 1
}

func filterListDomainsBody(all []interface{}, d *schema.ResourceData) []interface{} {
	var (
		domainID = d.Get("domain_id").(string)
		result   = make([]interface{}, 0, len(all))
	)

	for _, v := range all {
		if domainID != "" && domainID != fmt.Sprint(utils.PathSearch("id", v, nil)) {
			continue
		}
		result = append(result, v)
	}
	return result
}

func flattenSources(sources []interface{}) []interface{} {
	if len(sources) < 1 {
		return nil
	}

	result := make([]interface{}, 0, len(sources))
	for _, v := range sources {
		result = append(result, map[string]interface{}{
			"origin":                  utils.PathSearch("ip_or_domain", v, nil),
			"origin_type":             utils.PathSearch("origin_type", v, nil),
			"active":                  utils.PathSearch("active_standby", v, nil),
			"obs_web_hosting_enabled": converseOBSWebHostStatusToBool(utils.PathSearch("enable_obs_web_hosting", v, nil)),
		})
	}
	return result
}

func flattenListDomainsBody(domains []interface{}) []interface{} {
	result := make([]interface{}, 0, len(domains))

	for _, v := range domains {
		createTime := utils.PathSearch("create_time", v, 0)
		updateTime := utils.PathSearch("modify_time", v, 0)
		result = append(result, map[string]interface{}{
			"id":                            utils.PathSearch("id", v, nil),
			"name":                          utils.PathSearch("domain_name", v, nil),
			"type":                          utils.PathSearch("business_type", v, nil),
			"domain_status":                 utils.PathSearch("domain_status", v, nil),
			"cname":                         utils.PathSearch("cname", v, nil),
			"sources":                       flattenSources(utils.PathSearch("sources", v, make([]interface{}, 0)).([]interface{})),
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

	return result
}

func datasourceDomainsRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg        = meta.(*config.Config)
		httpUrl    = "v1.0/cdn/domains?page_size={page_size}"
		pageNumber = 1
		pageSize   = 10
		result     = make([]interface{}, 0)
	)

	client, err := cfg.NewServiceClient("cdn", "")
	if err != nil {
		return diag.Errorf("error creating CDN client: %s", err)
	}

	listPathWithSize := client.Endpoint + httpUrl
	listPathWithSize = strings.ReplaceAll(listPathWithSize, "{page_size}", strconv.Itoa(pageSize))
	listPathWithSize += buildDomainsQueryParams(cfg, d)
	opt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	for {
		listPathWithPage := fmt.Sprintf("%s&page_number=%v", listPathWithSize, pageNumber)
		resp, err := client.Request("GET", listPathWithPage, &opt)
		if err != nil {
			return diag.FromErr(err)
		}

		respBody, err := utils.FlattenResponse(resp)
		if err != nil {
			return diag.FromErr(err)
		}

		domains := utils.PathSearch("domains", respBody, make([]interface{}, 0)).([]interface{})
		if len(domains) == 0 {
			break
		}
		result = append(result, domains...)
		pageNumber++
	}

	dataSourceId, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(dataSourceId)

	mErr := multierror.Append(
		d.Set("domains", filterListDomainsBody(flattenListDomainsBody(result), d)),
	)
	return diag.FromErr(mErr.ErrorOrNil())
}
