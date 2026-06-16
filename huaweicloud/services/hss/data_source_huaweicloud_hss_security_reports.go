package hss

import (
	"context"
	"fmt"
	"strings"

	"github.com/google/uuid"
	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API HSS GET /v5/{project_id}/report/report-list
func DataSourceSecurityReports() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceSecurityReportsRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"enterprise_project_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"report_category": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"data_list": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"report_id": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"report_sub_id": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"default_report": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"latest_create_time": {
							Type:     schema.TypeInt,
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
						"report_create_time": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"sending_period": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func buildSecurityReportsQueryParams(d *schema.ResourceData, epsId string) string {
	queryParams := "?limit=200"

	if epsId != "" {
		queryParams = fmt.Sprintf("%s&enterprise_project_id=%v", queryParams, epsId)
	}
	if v, ok := d.GetOk("report_category"); ok {
		queryParams = fmt.Sprintf("%s&report_category=%v", queryParams, v)
	}

	return queryParams
}

func dataSourceSecurityReportsRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		product = "hss"
		epsId   = cfg.GetEnterpriseProjectID(d)
		result  = make([]interface{}, 0)
		offset  = 0
		httpUrl = "v5/{project_id}/report/report-list"
	)

	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating HSS client: %s", err)
	}

	requestPath := client.Endpoint + httpUrl
	requestPath = strings.ReplaceAll(requestPath, "{project_id}", client.ProjectID)
	requestPath += buildSecurityReportsQueryParams(d, epsId)
	requestOpts := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"region": region},
	}

	for {
		requestPathWithOffset := fmt.Sprintf("%s&offset=%v", requestPath, offset)
		resp, err := client.Request("GET", requestPathWithOffset, &requestOpts)
		if err != nil {
			return diag.Errorf("error retrieving HSS security reports: %s", err)
		}

		respBody, err := utils.FlattenResponse(resp)
		if err != nil {
			return diag.FromErr(err)
		}

		dataList := utils.PathSearch("data_list", respBody, make([]interface{}, 0)).([]interface{})
		if len(dataList) == 0 {
			break
		}

		result = append(result, dataList...)
		offset += len(dataList)
	}

	generateUUID, err := uuid.NewRandom()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}

	d.SetId(generateUUID.String())

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("data_list", flattenSecurityReportsDataList(result)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenSecurityReportsDataList(resp []interface{}) []interface{} {
	if len(resp) == 0 {
		return nil
	}

	rst := make([]interface{}, 0, len(resp))
	for _, v := range resp {
		rst = append(rst, map[string]interface{}{
			"report_id":          utils.PathSearch("report_id", v, nil),
			"report_sub_id":      utils.PathSearch("report_sub_id", v, nil),
			"default_report":     utils.PathSearch("default_report", v, nil),
			"latest_create_time": utils.PathSearch("latest_create_time", v, nil),
			"report_name":        utils.PathSearch("report_name", v, nil),
			"report_category":    utils.PathSearch("report_category", v, nil),
			"report_status":      utils.PathSearch("report_status", v, nil),
			"report_create_time": utils.PathSearch("report_create_time", v, nil),
			"sending_period":     utils.PathSearch("sending_period", v, nil),
		})
	}

	return rst
}
