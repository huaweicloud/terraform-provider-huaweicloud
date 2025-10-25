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
		ReadContext: dataSourceStatisticsRead,

		Schema: map[string]*schema.Schema{
			// Required parameters
			"domain_name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The domain name list.`,
			},
			"stat_type": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The status type.`,
			},
			"action": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The action name.`,
				ValidateFunc: validation.StringInSlice([]string{
					"location_summary", "location_detail",
				}, false),
			},
			"start_time": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: `The start timestamp of the query.`,
			},
			"end_time": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: `The end timestamp of the query.`,
			},

			// Optional parameters
			"interval": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: `The query time interval.`,
			},
			"group_by": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The data grouping mode.`,
			},
			"country": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The country or region code.`,
			},
			"province": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The province code.`,
			},
			"isp": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The carrier code.`,
			},
			"enterprise_project_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The ID of the enterprise project that the resource belongs.`,
			},

			// Attribute
			"result": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The data organized according to the specified grouping mode.`,
			},
		},
	}
}

func buildDomainStatisticsQueryParams(d *schema.ResourceData) string {
	res := "?"

	res = fmt.Sprintf("%s&domain_name=%v", res, d.Get("domain_name"))
	res = fmt.Sprintf("%s&stat_type=%v", res, d.Get("stat_type"))
	res = fmt.Sprintf("%s&action=%v", res, d.Get("action"))
	res = fmt.Sprintf("%s&start_time=%v", res, d.Get("start_time"))
	res = fmt.Sprintf("%s&end_time=%v", res, d.Get("end_time"))

	if v, ok := d.GetOk("interval"); ok {
		res = fmt.Sprintf("%s&interval=%v", res, v)
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

	return res
}

func dataSourceStatisticsRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		httpUrl = "v1.0/cdn/statistics/domain-location-stats"
	)

	client, err := cfg.NewServiceClient("cdn", "")
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
		return diag.Errorf("error querying domain statistics: %s", err)
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

	mErr := multierror.Append(
		d.Set("result", strings.TrimSpace(string(results))),
	)
	return diag.FromErr(mErr.ErrorOrNil())
}
