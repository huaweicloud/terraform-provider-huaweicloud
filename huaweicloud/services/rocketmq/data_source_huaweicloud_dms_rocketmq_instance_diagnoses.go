package rocketmq

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API RocketMQ GET /v2/{engine}/{project_id}/instances/{instance_id}/diagnosis
func DataSourceInstanceDiagnoses() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceInstanceDiagnosesRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `The region where the diagnosis reports are located.`,
			},
			"instance_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The ID of the RocketMQ instance.`,
			},
			"reports": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: `The list of the diagnosis reports.`,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"report_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The ID of the diagnosis report.`,
						},
						"group_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The name of the consumer group.`,
						},
						"status": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The status of the report.`,
						},
						"created_at": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The creation time of the report.`,
						},
						"abnormal_item_sum": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: `The number of abnormal items.`,
						},
						"faulted_node_sum": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: `The number of faulted nodes.`,
						},
					},
				},
			},
		},
	}
}

func listInstanceDiagnosisReports(client *golangsdk.ServiceClient, d *schema.ResourceData) ([]interface{}, error) {
	var (
		httpUrl = "v2/{engine}/{project_id}/instances/{instance_id}/diagnosis?limit={limit}"
		limit   = 50 // max limit is 50
		offset  = 0
		result  = make([]interface{}, 0)
	)

	listPath := client.Endpoint + httpUrl
	listPath = strings.ReplaceAll(listPath, "{engine}", "rocketmq")
	listPath = strings.ReplaceAll(listPath, "{project_id}", client.ProjectID)
	listPath = strings.ReplaceAll(listPath, "{instance_id}", d.Get("instance_id").(string))
	listPath = strings.ReplaceAll(listPath, "{limit}", strconv.Itoa(limit))

	opt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json",
		},
	}

	for {
		listPathWithOffset := listPath + fmt.Sprintf("&offset=%d", offset)
		requestResp, err := client.Request("GET", listPathWithOffset, &opt)
		if err != nil {
			return nil, err
		}

		respBody, err := utils.FlattenResponse(requestResp)
		if err != nil {
			return nil, err
		}

		reports := utils.PathSearch("diagnosis_report_list", respBody, make([]interface{}, 0)).([]interface{})
		result = append(result, reports...)

		if len(reports) < limit {
			break
		}

		offset += len(reports)
	}

	return result, nil
}

func dataSourceInstanceDiagnosesRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	client, err := cfg.NewServiceClient("dmsv2", region)
	if err != nil {
		return diag.Errorf("error creating DMS client: %s", err)
	}

	reports, err := listInstanceDiagnosisReports(client, d)
	if err != nil {
		return diag.Errorf("error retrieving instance diagnosis reports: %s", err)
	}

	if len(reports) == 0 {
		return diag.Errorf("unable to find any diagnosis reports")
	}

	randUUID, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(randUUID)

	mErr := multierror.Append(
		d.Set("region", region),
		d.Set("reports", flattenDiagnosisReports(reports)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenDiagnosisReports(reports []interface{}) []map[string]interface{} {
	if len(reports) == 0 {
		return nil
	}

	result := make([]map[string]interface{}, 0, len(reports))
	for _, report := range reports {
		result = append(result, map[string]interface{}{
			"report_id":         utils.PathSearch("report_id", report, nil),
			"group_name":        utils.PathSearch("group_name", report, nil),
			"status":            utils.PathSearch("status", report, nil),
			"created_at":        utils.PathSearch("created_at", report, nil),
			"abnormal_item_sum": utils.PathSearch("abnormal_item_sum", report, nil),
			"faulted_node_sum":  utils.PathSearch("faulted_node_sum", report, nil),
		})
	}

	return result
}
