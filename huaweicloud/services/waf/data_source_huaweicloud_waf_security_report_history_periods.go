package waf

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

// @API WAF GET /v1/{project_id}/waf/security-report/history-periods
func DataSourceSecurityReportHistoryPeriods() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceSecurityReportHistoryPeriodsRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"subscription_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"items": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"report_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"subscription_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"stat_period": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"begin_time": {
										Type:     schema.TypeInt,
										Computed: true,
									},
									"end_time": {
										Type:     schema.TypeInt,
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

func buildSecurityReportHistoryPeriodsQueryParams(d *schema.ResourceData, offset int) string {
	queryParams := fmt.Sprintf("?subscription_id=%s", d.Get("subscription_id").(string))
	if offset > 0 {
		queryParams += fmt.Sprintf("&offset=%d", offset)
	}

	return queryParams
}

func dataSourceSecurityReportHistoryPeriodsRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		httpUrl = "v1/{project_id}/waf/security-report/history-periods"
		offset  = 0
		result  = make([]interface{}, 0)
	)

	client, err := cfg.NewServiceClient("waf", region)
	if err != nil {
		return diag.Errorf("error creating WAF client: %s", err)
	}

	requestPath := client.Endpoint + httpUrl
	requestPath = strings.ReplaceAll(requestPath, "{project_id}", client.ProjectID)
	requestOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json;charset=utf8",
		},
	}

	for {
		currentPath := requestPath + buildSecurityReportHistoryPeriodsQueryParams(d, offset)
		resp, err := client.Request("GET", currentPath, &requestOpt)
		if err != nil {
			return diag.Errorf("error retrieving security report history periods: %s", err)
		}

		respBody, err := utils.FlattenResponse(resp)
		if err != nil {
			return diag.FromErr(err)
		}

		periods := utils.PathSearch("items", respBody, make([]interface{}, 0)).([]interface{})
		if len(periods) == 0 {
			break
		}

		result = append(result, periods...)
		offset += len(periods)
	}

	generateUUID, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}

	d.SetId(generateUUID)

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("items", flattenSecurityReportHistoryPeriods(result)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenSecurityReportHistoryPeriods(respArray []interface{}) []interface{} {
	if len(respArray) == 0 {
		return nil
	}

	result := make([]interface{}, 0, len(respArray))
	for _, v := range respArray {
		result = append(result, map[string]interface{}{
			"report_id":       utils.PathSearch("report_id", v, nil),
			"subscription_id": utils.PathSearch("subscription_id", v, nil),
			"stat_period":     flattenStatPeriod(utils.PathSearch("stat_period", v, nil)),
		})
	}
	return result
}

func flattenStatPeriod(statPeriod interface{}) []map[string]interface{} {
	if statPeriod == nil {
		return nil
	}
	return []map[string]interface{}{
		{
			"begin_time": utils.PathSearch("begin_time", statPeriod, nil),
			"end_time":   utils.PathSearch("end_time", statPeriod, nil),
		},
	}
}
