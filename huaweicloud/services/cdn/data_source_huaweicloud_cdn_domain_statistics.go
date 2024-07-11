// ---------------------------------------------------------------
// *** AUTO GENERATED CODE ***
// @Product CDN
// ---------------------------------------------------------------

package cdn

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API CDN GET /v1.0/cdn/statistics/domain-location-stats
func DataSourceStatistics() *schema.Resource {
	return &schema.Resource{
		ReadContext: resourceStatisticsRead,
		Schema: map[string]*schema.Schema{
			"domain_name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `Specifies the domain name list.`,
			},
			"stat_type": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `Specifies the status type.`,
			},
			"action": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `Specifies the action name.`,
				ValidateFunc: validation.StringInSlice([]string{
					"location_summary", "location_detail",
				}, false),
			},
			"start_time": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: `Specifies the start timestamp of the query.`,
			},
			"end_time": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: `Specifies the end timestamp of the query.`,
			},
			"interval": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: `Specifies the query time interval.`,
			},
			"group_by": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specifies the data grouping mode.`,
			},
			"country": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specifies the country or region code.`,
			},
			"province": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specifies the province code.`,
			},
			"isp": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specifies the carrier code.`,
			},
			"enterprise_project_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specifies the enterprise project that the resource belongs to.`,
			},
			"result": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The data organized according to the specified grouping mode.`,
			},
		},
	}
}

func resourceStatisticsRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conf := meta.(*config.Config)
	region := conf.GetRegion(d)

	var mErr *multierror.Error

	// domainStatistics: Query the statistics of CDN domain.
	var (
		domainStatisticsHttpUrl = "v1.0/cdn/statistics/domain-location-stats"
		domainStatisticsProduct = "cdn"
	)
	domainStatisticsClient, err := conf.NewServiceClient(domainStatisticsProduct, region)
	if err != nil {
		return diag.Errorf("error creating Statistics Client: %s", err)
	}

	domainStatisticsPath := domainStatisticsClient.Endpoint + domainStatisticsHttpUrl

	domainStatisticsqueryParams := buildDomainStatisticsQueryParams(d)
	domainStatisticsPath += domainStatisticsqueryParams

	domainStatisticsOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		OkCodes: []int{
			200,
		},
	}
	domainStatisticsResp, err := domainStatisticsClient.Request("GET", domainStatisticsPath, &domainStatisticsOpt)

	if err != nil {
		return diag.Errorf("error retrieving CDN statistics: %s", err)
	}

	domainStatisticsRespBody, err := utils.FlattenResponse(domainStatisticsResp)
	if err != nil {
		return diag.FromErr(err)
	}

	uuid, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(uuid)

	results, err := utils.JsonMarshal(utils.PathSearch("result", domainStatisticsRespBody, nil))
	if err != nil {
		return diag.Errorf("error marshaling Statistics result: %s", err)
	}

	mErr = multierror.Append(
		mErr,
		d.Set("result", strings.TrimSpace(string(results))),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func buildDomainStatisticsQueryParams(d *schema.ResourceData) string {
	res := ""
	if v, ok := d.GetOk("action"); ok {
		res = fmt.Sprintf("%s&action=%v", res, v)
	}
	if v, ok := d.GetOk("start_time"); ok {
		res = fmt.Sprintf("%s&start_time=%v", res, v)
	}
	if v, ok := d.GetOk("end_time"); ok {
		res = fmt.Sprintf("%s&end_time=%v", res, v)
	}
	if v, ok := d.GetOk("interval"); ok {
		res = fmt.Sprintf("%s&interval=%v", res, v)
	}
	if v, ok := d.GetOk("domain_name"); ok {
		res = fmt.Sprintf("%s&domain_name=%v", res, v)
	}
	if v, ok := d.GetOk("stat_type"); ok {
		res = fmt.Sprintf("%s&stat_type=%v", res, v)
	}
	if v, ok := d.GetOk("group_by"); ok {
		res = fmt.Sprintf("%s&group_by=%v", res, v)
	}
	if v, ok := d.GetOk("country"); ok {
		res = fmt.Sprintf("%s&country=%v", res, v)
	}
	if v, ok := d.GetOk("province"); ok {
		res = fmt.Sprintf("%s&province=%v", res, v)
	}
	if v, ok := d.GetOk("isp"); ok {
		res = fmt.Sprintf("%s&isp=%v", res, v)
	}
	if v, ok := d.GetOk("enterprise_project_id"); ok {
		res = fmt.Sprintf("%s&enterprise_project_id=%v", res, v)
	}
	if res != "" {
		res = "?" + res[1:]
	}
	return res
}
