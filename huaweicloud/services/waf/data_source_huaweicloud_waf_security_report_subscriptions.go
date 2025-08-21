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

// @API WAF GET /v1/{project_id}/waf/security-report/subscriptions
func DataSourceSecurityReportSubscriptions() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceSecurityReportSubscriptionsRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"report_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"report_category": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"report_status": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"enterprise_project_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"items": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"subscription_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"report_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"report_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"report_category": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"report_status": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"sending_period": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"is_all_enterprise_project": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"enterprise_project_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"template_eps_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"is_report_created": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"latest_create_time": {
							Type:     schema.TypeInt,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func buildSecurityReportSubscriptionsQueryParams(d *schema.ResourceData, epsId string) string {
	queryParams := "?limit=1000"

	if v, ok := d.GetOk("report_name"); ok {
		queryParams = fmt.Sprintf("%s&report_name=%v", queryParams, v)
	}
	if v, ok := d.GetOk("report_category"); ok {
		queryParams = fmt.Sprintf("%s&report_category=%v", queryParams, v)
	}
	if v, ok := d.GetOk("report_status"); ok {
		queryParams = fmt.Sprintf("%s&report_status=%v", queryParams, v)
	}
	if epsId != "" {
		queryParams = fmt.Sprintf("%s&enterprise_project_id=%v", queryParams, epsId)
	}

	return queryParams
}

func dataSourceSecurityReportSubscriptionsRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		httpUrl = "v1/{project_id}/waf/security-report/subscriptions"
		epsId   = cfg.GetEnterpriseProjectID(d)
		offset  = 0
		result  = make([]interface{}, 0)
	)

	client, err := cfg.NewServiceClient("waf", region)
	if err != nil {
		return diag.Errorf("error creating WAF client: %s", err)
	}

	getPath := client.Endpoint + httpUrl
	getPath = strings.ReplaceAll(getPath, "{project_id}", client.ProjectID)
	getPath += buildSecurityReportSubscriptionsQueryParams(d, epsId)
	getOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json;charset=utf8",
		},
	}

	for {
		currentPath := fmt.Sprintf("%s&offset=%v", getPath, offset)
		getResp, err := client.Request("GET", currentPath, &getOpt)
		if err != nil {
			return diag.Errorf("error retrieving security report subscriptions: %s", err)
		}

		getRespBody, err := utils.FlattenResponse(getResp)
		if err != nil {
			return diag.FromErr(err)
		}

		subscription := utils.PathSearch("items", getRespBody, make([]interface{}, 0)).([]interface{})
		if len(subscription) == 0 {
			break
		}

		result = append(result, subscription...)
		offset += len(subscription)
	}

	generateUUID, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}

	d.SetId(generateUUID)

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("items", flattenSecurityReportSubscriptions(result)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenSecurityReportSubscriptions(subscriptions []interface{}) []interface{} {
	result := make([]interface{}, 0, len(subscriptions))
	for _, v := range subscriptions {
		result = append(result, map[string]interface{}{
			"subscription_id":           utils.PathSearch("subscription_id", v, nil),
			"report_id":                 utils.PathSearch("report_id", v, nil),
			"report_name":               utils.PathSearch("report_name", v, nil),
			"report_category":           utils.PathSearch("report_category", v, nil),
			"report_status":             utils.PathSearch("report_status", v, nil),
			"sending_period":            utils.PathSearch("sending_period", v, nil),
			"is_all_enterprise_project": utils.PathSearch("is_all_enterprise_project", v, nil),
			"enterprise_project_id":     utils.PathSearch("enterprise_project_id", v, nil),
			"template_eps_id":           utils.PathSearch("template_eps_id", v, nil),
			"is_report_created":         utils.PathSearch("is_report_created", v, nil),
			"latest_create_time":        utils.PathSearch("latest_create_time", v, nil),
		})
	}

	return result
}
