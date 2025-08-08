package secmaster

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

// @API SECMASTER GET /v1/{project_id}/workspaces/{workspace_id}/sa/reports
func DataSourceSecurityReports() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceSecurityReportsRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"workspace_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"report_period": {
				Type:     schema.TypeString,
				Required: true,
			},
			"status": {
				Type:     schema.TypeString,
				Required: true,
			},
			"reports": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"report_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"report_period": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"report_range": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"start": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"end": {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						},
						"language": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"notification_task": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"layout_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"status": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"is_generated": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"report_rule_infos": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"id": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"project_id": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"workspace_id": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"cycle": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"rule": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"start_time": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"end_time": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"email_title": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"email_to": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"email_content": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"report_file_type": {
										Type:     schema.TypeString,
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

func dataSourceSecurityReportsRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg          = meta.(*config.Config)
		region       = cfg.GetRegion(d)
		workspaceId  = d.Get("workspace_id").(string)
		reportPeriod = d.Get("report_period").(string)
		status       = d.Get("status").(string)
		httpUrl      = "v1/{project_id}/workspaces/{workspace_id}/sa/reports"
	)

	client, err := cfg.NewServiceClient("secmaster", region)
	if err != nil {
		return diag.Errorf("error creating SecMaster client: %s", err)
	}

	listPath := client.Endpoint + httpUrl
	listPath = strings.ReplaceAll(listPath, "{project_id}", client.ProjectID)
	listPath = strings.ReplaceAll(listPath, "{workspace_id}", workspaceId)
	listPath = fmt.Sprintf("%s?report_period=%s&status=%s", listPath, reportPeriod, status)

	listOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	listResp, err := client.Request("GET", listPath, &listOpt)
	if err != nil {
		return diag.Errorf("error retrieving SecMaster security reports: %s", err)
	}

	listRespBody, err := utils.FlattenResponse(listResp)
	if err != nil {
		return diag.FromErr(err)
	}

	listArray, ok := listRespBody.([]interface{})
	if !ok {
		return diag.Errorf("convert interface array failed")
	}

	dataSourceId, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}

	d.SetId(dataSourceId)

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("reports", flattenSecurityReports(listArray)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenSecurityReports(rawReports []interface{}) []interface{} {
	if len(rawReports) == 0 {
		return nil
	}

	rst := make([]interface{}, 0, len(rawReports))
	for _, v := range rawReports {
		rst = append(rst, map[string]interface{}{
			"id":                utils.PathSearch("id", v, nil),
			"report_name":       utils.PathSearch("report_name", v, nil),
			"report_period":     utils.PathSearch("report_period", v, nil),
			"report_range":      flattenSecurityReportRange(utils.PathSearch("report_range", v, nil)),
			"language":          utils.PathSearch("language", v, nil),
			"notification_task": utils.PathSearch("notification_task", v, nil),
			"layout_id":         utils.PathSearch("layout_id", v, nil),
			"status":            utils.PathSearch("status", v, nil),
			"is_generated":      utils.PathSearch("is_generated", v, nil),
			"report_rule_infos": flattenSecurityReportRuleInfos(utils.PathSearch("report_rule_infos", v, make([]interface{}, 0)).([]interface{})),
		})
	}

	return rst
}

func flattenSecurityReportRange(rawRange interface{}) []interface{} {
	if rawRange == nil {
		return nil
	}

	return []interface{}{
		map[string]interface{}{
			"start": utils.PathSearch("start", rawRange, nil),
			"end":   utils.PathSearch("end", rawRange, nil),
		},
	}
}

func flattenSecurityReportRuleInfos(rawRuleInfos []interface{}) []interface{} {
	if len(rawRuleInfos) == 0 {
		return nil
	}

	rst := make([]interface{}, 0, len(rawRuleInfos))
	for _, v := range rawRuleInfos {
		rst = append(rst, map[string]interface{}{
			"id":               utils.PathSearch("id", v, nil),
			"project_id":       utils.PathSearch("projectId", v, nil),
			"workspace_id":     utils.PathSearch("workspaceId", v, nil),
			"cycle":            utils.PathSearch("cycle", v, nil),
			"rule":             utils.PathSearch("rule", v, nil),
			"start_time":       utils.PathSearch("startTime", v, nil),
			"end_time":         utils.PathSearch("endTime", v, nil),
			"email_title":      utils.PathSearch("emailTitle", v, nil),
			"email_to":         utils.PathSearch("emailTo", v, nil),
			"email_content":    utils.PathSearch("emailContent", v, nil),
			"report_file_type": utils.PathSearch("report_file_type", v, nil),
		})
	}

	return rst
}
