package cdn

import (
	"context"
	"fmt"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API CDN GET /v1.0/cdn/statistics/top-url
func DataSourceTopUrlStatistics() *schema.Resource {
	return &schema.Resource{
		ReadContext: resourceTopUrlStatisticsRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `The region where the CDN service is located.`,
			},

			// Required parameters.
			"start_time": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The start time of the query, in UTC format.`,
			},
			"end_time": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The end time of the query, in UTC format.`,
			},
			"domain_name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The list of queried domain names`,
			},
			"stat_type": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The query type.`,
			},

			// Optional parameters.
			"service_area": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `The service area of the query.`,
			},
			"enterprise_project_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The ID of the enterprise project to which the resources belong.`,
			},

			// Attributes
			"top_url_summary": {
				Type:        schema.TypeList,
				Computed:    true,
				Elem:        topUrlSummarySchema(),
				Description: `The list of TOP100 URL statistics that matched filter parameters.`,
			},
		},
	}
}

func topUrlSummarySchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"url": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The URL name.`,
			},
			"value": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: `The value corresponding to the query type.`,
			},
			"start_time": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The start time of the query, in UTC format.`,
			},
			"end_time": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The end time of the query, in UTC format.`,
			},
			"stat_type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The query type.`,
			},
		},
	}
}

func buildTopUrlStatisticsQueryParams(d *schema.ResourceData) string {
	res := "?"

	// Convert UTC format str to timestamp
	if v, ok := d.GetOk("start_time"); ok {
		res = fmt.Sprintf("%s&start_time=%v", res, utils.ConvertTimeStrToNanoTimestamp(v.(string)))
	}
	// Convert UTC format str to timestamp
	if v, ok := d.GetOk("end_time"); ok {
		res = fmt.Sprintf("%s&end_time=%v", res, utils.ConvertTimeStrToNanoTimestamp(v.(string)))
	}
	if v, ok := d.GetOk("domain_name"); ok {
		res = fmt.Sprintf("%s&domain_name=%v", res, v)
	}
	if v, ok := d.GetOk("stat_type"); ok {
		res = fmt.Sprintf("%s&stat_type=%v", res, v)
	}
	if v, ok := d.GetOk("service_area"); ok {
		res = fmt.Sprintf("%s&service_area=%v", res, v)
	}
	if v, ok := d.GetOk("enterprise_project_id"); ok {
		res = fmt.Sprintf("%s&enterprise_project_id=%v", res, v)
	}

	return res
}

func flattenTopUrlSummary(items []interface{}) []map[string]interface{} {
	if len(items) < 1 {
		return nil
	}

	result := make([]map[string]interface{}, 0, len(items))
	for _, item := range items {
		// Convert timestamp to UTC format str
		result = append(result, map[string]interface{}{
			"url":        utils.PathSearch("url", item, nil),
			"value":      utils.PathSearch("value", item, nil),
			"start_time": utils.FormatTimeStampUTC(utils.PathSearch("start_time", item, int64(0)).(int64)),
			"end_time":   utils.FormatTimeStampUTC(utils.PathSearch("end_time", item, int64(0)).(int64)),
			"stat_type":  utils.PathSearch("stat_type", item, nil),
		})
	}

	return result
}

func resourceTopUrlStatisticsRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		httpUrl = "v1.0/cdn/statistics/top-url"
	)

	client, err := cfg.NewServiceClient("cdn", region)
	if err != nil {
		return diag.Errorf("error creating CDN client: %s", err)
	}

	listPath := client.Endpoint + httpUrl
	listPath += buildTopUrlStatisticsQueryParams(d)

	opt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json",
		},
	}

	requestResp, err := client.Request("GET", listPath, &opt)
	if err != nil {
		return diag.Errorf("error querying TOP100 URL statistics: %s", err)
	}

	resp, err := utils.FlattenResponse(requestResp)
	if err != nil {
		return diag.Errorf("error parsing TOP100 URL statistics: %s", err)
	}

	randomUUID, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(randomUUID)

	mErr := multierror.Append(
		d.Set("region", region),
		d.Set("service_area", utils.PathSearch("service_area", resp, d.Get("service_area").(string))),
		d.Set("top_url_summary", flattenTopUrlSummary(utils.PathSearch("top_url_summary", resp,
			make([]interface{}, 0)).([]interface{}))),
	)
	return diag.FromErr(mErr.ErrorOrNil())
}
