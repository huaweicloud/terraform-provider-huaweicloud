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
	var (
		conf    = meta.(*config.Config)
		region  = conf.GetRegion(d)
		mErr    *multierror.Error
		httpUrl = "v1.0/cdn/statistics/domain-location-stats"
		product = "cdn"
	)
	client, err := conf.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating CDN client: %s", err)
	}

	requestPath := client.Endpoint + httpUrl
	requestPath += buildDomainStatisticsQueryParams(d)
	requestOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	resp, err := client.Request("GET", requestPath, &requestOpt)
	if err != nil {
		return diag.Errorf("error retrieving CDN statistics: %s", err)
	}

	respBody, err := utils.FlattenResponse(resp)
	if err != nil {
		return diag.FromErr(err)
	}

	generateUUID, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(generateUUID)

	results, err := utils.JsonMarshal(utils.PathSearch("result", respBody, nil))
	if err != nil {
		return diag.Errorf("error marshaling statistics result: %s", err)
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
