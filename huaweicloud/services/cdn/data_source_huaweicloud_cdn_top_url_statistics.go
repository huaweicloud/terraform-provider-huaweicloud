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
		ReadContext: dataSourceTopUrlStatisticsRead,

		Schema: map[string]*schema.Schema{
			// Required parameters.
			"start_time": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The start time of the query range, in RFC3339 format.`,
			},
			"end_time": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The end time of the query range, in RFC3339 format.`,
			},
			"domain_name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The list of queried domain names.`,
			},
			"stat_type": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The statistical type of the query.`,
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
				Description: `The ID of the enterprise project to which the domains belong.`,
			},

			// Attributes
			"statistics": {
				Type:        schema.TypeList,
				Computed:    true,
				Elem:        topUrlStatisticsSchema(),
				Description: `The list of TOP100 URL statistics that matched filter parameters.`,
			},
		},
	}
}

func topUrlStatisticsSchema() *schema.Resource {
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
				Description: `The start time of the query range, in RFC3339 format.`,
			},
			"end_time": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The end time of the query range, in RFC3339 format.`,
			},
			"stat_type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The statistical type of the query.`,
			},
		},
	}
}

func buildTopUrlStatisticsQueryParams(d *schema.ResourceData) string {
	res := "?"

	// Convert RFC3339 format str to timestamp
	res = fmt.Sprintf("%s&start_time=%v", res, utils.ConvertTimeStrToNanoTimestamp(d.Get("start_time").(string)))
	// Convert RFC3339 format str to timestamp
	res = fmt.Sprintf("%s&end_time=%v", res, utils.ConvertTimeStrToNanoTimestamp(d.Get("end_time").(string)))
	res = fmt.Sprintf("%s&domain_name=%v", res, d.Get("domain_name").(string))
	res = fmt.Sprintf("%s&stat_type=%v", res, d.Get("stat_type").(string))

	if v, ok := d.GetOk("service_area"); ok {
		res = fmt.Sprintf("%s&service_area=%v", res, v)
	}
	if v, ok := d.GetOk("enterprise_project_id"); ok {
		res = fmt.Sprintf("%s&enterprise_project_id=%v", res, v)
	}

	return res
}

func flattenTopUrlStatistics(items []interface{}) []map[string]interface{} {
	if len(items) < 1 {
		return nil
	}

	result := make([]map[string]interface{}, 0, len(items))
	for _, item := range items {
		result = append(result, map[string]interface{}{
			"url":        utils.PathSearch("url", item, nil),
			"value":      utils.PathSearch("value", item, nil),
			"start_time": utils.FormatTimeStampRFC3339(int64(utils.PathSearch("start_time", item, float64(0)).(float64))/1000, true),
			"end_time":   utils.FormatTimeStampRFC3339(int64(utils.PathSearch("end_time", item, float64(0)).(float64))/1000, true),
			"stat_type":  utils.PathSearch("stat_type", item, nil),
		})
	}

	return result
}

func dataSourceTopUrlStatisticsRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		httpUrl = "v1.0/cdn/statistics/top-url"
	)

	client, err := cfg.NewServiceClient("cdn", "")
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
		d.Set("statistics", flattenTopUrlStatistics(utils.PathSearch("top_url_summary", resp,
			make([]interface{}, 0)).([]interface{}))),
	)
	return diag.FromErr(mErr.ErrorOrNil())
}
