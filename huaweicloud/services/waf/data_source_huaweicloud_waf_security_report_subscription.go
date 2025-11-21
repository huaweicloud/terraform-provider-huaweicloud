package waf

import (
	"context"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API WAF GET /v1/{project_id}/waf/security-report/subscriptions/{subscription_id}
func DataSourceWafSecurityReportSubscription() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceWafSecurityReportSubscriptionRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `Specifies the region in which to query the resource. If omitted, the provider-level region will be used.`,
			},
			"subscription_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `Specifies the ID of the security report subscription.`,
			},
			"sending_period": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The sending period of the security report.`,
			},
			"report_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The name of the security report.`,
			},
			"report_category": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The category of the security report.`,
			},
			"topic_urn": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The URN of the SMN topic for receiving reports.`,
			},
			"subscription_type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The subscription type of the security report.`,
			},
			"report_content_subscription": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: `The content subscription configuration of the security report.`,
				Elem:        buildReportContentSubscriptionSchema(),
			},
			"stat_period": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: `The statistical period of the security report.`,
				Elem:        buildStatPeriodSchema(),
			},
			"is_all_enterprise_project": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: `Whether the subscription applies to all enterprise projects.`,
			},
			"enterprise_project_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The enterprise project ID associated with the subscription.`,
			},
		},
	}
}

func buildReportContentSubscriptionSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"overview_statistics_enable": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: `Whether to enable overview statistics.`,
			},
			"group_by_day_enable": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: `Whether to enable group by day statistics.`,
			},
			"request_statistics_enable": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: `Whether to enable request statistics.`,
			},
			"qps_statistics_enable": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: `Whether to enable QPS statistics.`,
			},
			"bandwidth_statistics_enable": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: `Whether to enable bandwidth statistics.`,
			},
			"response_code_statistics_enable": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: `Whether to enable response code statistics.`,
			},
			"attack_type_distribution_enable": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: `Whether to enable attack type distribution statistics.`,
			},
			"top_attacked_domains_enable": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: `Whether to enable top attacked domains statistics.`,
			},
			"top_attack_source_ips_enable": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: `Whether to enable top attack source IPs statistics.`,
			},
			"top_attacked_urls_enable": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: `Whether to enable top attacked URLs statistics.`,
			},
			"top_attack_source_locations_enable": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: `Whether to enable top attack source locations statistics.`,
			},
			"top_abnormal_urls_enable": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: `Whether to enable top abnormal URLs statistics.`,
			},
		},
	}
}

func buildStatPeriodSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"begin_time": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: `The start time of the statistical period.`,
			},
			"end_time": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: `The end time of the statistical period.`,
			},
		},
	}
}

func dataSourceWafSecurityReportSubscriptionRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		httpUrl = "v1/{project_id}/waf/security-report/subscriptions/{subscription_id}"
	)

	client, err := cfg.NewServiceClient("waf", region)
	if err != nil {
		return diag.Errorf("error creating WAF client: %s", err)
	}

	requestPath := client.Endpoint + httpUrl
	requestPath = strings.ReplaceAll(requestPath, "{project_id}", client.ProjectID)
	requestPath = strings.ReplaceAll(requestPath, "{subscription_id}", d.Get("subscription_id").(string))
	requestOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json;charset=utf8",
		},
	}

	resp, err := client.Request("GET", requestPath, &requestOpt)
	if err != nil {
		return diag.Errorf("error retrieving security report subscription: %s", err)
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

	reportContentSubscription := utils.PathSearch("report_content_subscription", respBody, nil)
	statPeriod := utils.PathSearch("stat_period", respBody, nil)

	mErr := multierror.Append(
		d.Set("region", region),
		d.Set("sending_period", utils.PathSearch("sending_period", respBody, nil)),
		d.Set("report_name", utils.PathSearch("report_name", respBody, nil)),
		d.Set("report_category", utils.PathSearch("report_category", respBody, nil)),
		d.Set("topic_urn", utils.PathSearch("topic_urn", respBody, nil)),
		d.Set("subscription_type", utils.PathSearch("subscription_type", respBody, nil)),
		d.Set("report_content_subscription", flattenReportContentSubscriptionAttribute(reportContentSubscription)),
		d.Set("stat_period", flattenStatPeriodAttribute(statPeriod)),
		d.Set("is_all_enterprise_project", utils.PathSearch("is_all_enterprise_project", respBody, nil)),
		d.Set("enterprise_project_id", utils.PathSearch("enterprise_project_id", respBody, nil)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenReportContentSubscriptionAttribute(rawMap interface{}) []interface{} {
	if rawMap == nil {
		return nil
	}

	rstMap := map[string]interface{}{
		"overview_statistics_enable":         utils.PathSearch("overview_statistics_enable", rawMap, nil),
		"group_by_day_enable":                utils.PathSearch("group_by_day_enable", rawMap, nil),
		"request_statistics_enable":          utils.PathSearch("request_statistics_enable", rawMap, nil),
		"qps_statistics_enable":              utils.PathSearch("qps_statistics_enable", rawMap, nil),
		"bandwidth_statistics_enable":        utils.PathSearch("bandwidth_statistics_enable", rawMap, nil),
		"response_code_statistics_enable":    utils.PathSearch("response_code_statistics_enable", rawMap, nil),
		"attack_type_distribution_enable":    utils.PathSearch("attack_type_distribution_enable", rawMap, nil),
		"top_attacked_domains_enable":        utils.PathSearch("top_attacked_domains_enable", rawMap, nil),
		"top_attack_source_ips_enable":       utils.PathSearch("top_attack_source_ips_enable", rawMap, nil),
		"top_attacked_urls_enable":           utils.PathSearch("top_attacked_urls_enable", rawMap, nil),
		"top_attack_source_locations_enable": utils.PathSearch("top_attack_source_locations_enable", rawMap, nil),
		"top_abnormal_urls_enable":           utils.PathSearch("top_abnormal_urls_enable", rawMap, nil),
	}

	return []interface{}{rstMap}
}

func flattenStatPeriodAttribute(rawMap interface{}) []interface{} {
	if rawMap == nil {
		return nil
	}

	rstMap := map[string]interface{}{
		"begin_time": utils.PathSearch("begin_time", rawMap, nil),
		"end_time":   utils.PathSearch("end_time", rawMap, nil),
	}

	return []interface{}{rstMap}
}
