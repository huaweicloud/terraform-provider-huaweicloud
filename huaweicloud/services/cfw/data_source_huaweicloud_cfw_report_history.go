package cfw

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

// @API CFW GET /v1/{project_id}/report/history/{report_profile_id}
func DataSourceReportHistory() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceReportHistoryRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"report_profile_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"fw_instance_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"records": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"report_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"send_time": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"statistic_period": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"end_time": {
										Type:     schema.TypeInt,
										Computed: true,
									},
									"start_time": {
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

func buildReportHistoryQueryParams(d *schema.ResourceData) string {
	return fmt.Sprintf("?fw_instance_id=%s&limit=1024", d.Get("fw_instance_id"))
}

func dataSourceReportHistoryRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		product = "cfw"
		httpUrl = "v1/{project_id}/report/history/{report_profile_id}"
		offset  = 0
		result  = make([]interface{}, 0)
	)

	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating CFW client: %s", err)
	}

	requestPath := client.Endpoint + httpUrl
	requestPath = strings.ReplaceAll(requestPath, "{project_id}", client.ProjectID)
	requestPath = strings.ReplaceAll(requestPath, "{report_profile_id}", d.Get("report_profile_id").(string))
	requestPath += buildReportHistoryQueryParams(d)
	requestOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	for {
		requestPathWithOffset := fmt.Sprintf("%s&offset=%d", requestPath, offset)
		resp, err := client.Request("GET", requestPathWithOffset, &requestOpt)
		if err != nil {
			return diag.Errorf("error retrieving CFW report history: %s", err)
		}

		respBody, err := utils.FlattenResponse(resp)
		if err != nil {
			return diag.FromErr(err)
		}

		records := utils.PathSearch("data.records", respBody, make([]interface{}, 0)).([]interface{})
		if len(records) == 0 {
			break
		}

		result = append(result, records...)
		offset += len(records)
	}

	generateUUID, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}

	d.SetId(generateUUID)

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("records", flattenReportHistoryRecords(result)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenReportHistoryRecords(recordsResp []interface{}) []interface{} {
	if len(recordsResp) == 0 {
		return nil
	}

	result := make([]interface{}, 0, len(recordsResp))
	for _, raw := range recordsResp {
		result = append(result, map[string]interface{}{
			"report_id":        utils.PathSearch("report_id", raw, nil),
			"send_time":        utils.PathSearch("send_time", raw, nil),
			"statistic_period": flattenReportHistoryStatisticPeriod(utils.PathSearch("statistic_period", raw, nil)),
		})
	}

	return result
}

func flattenReportHistoryStatisticPeriod(raw interface{}) []interface{} {
	if raw == nil {
		return nil
	}

	return []interface{}{map[string]interface{}{
		"end_time":   utils.PathSearch("end_time", raw, nil),
		"start_time": utils.PathSearch("start_time", raw, nil),
	}}
}
