package das

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/google/uuid"
	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API DAS GET /v3/{project_id}/batch-inspection/email-record
func DataSourceEmailSendRecords() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceEmailSendRecordsRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "The region where the email send records are located.",
			},

			// Required parameters.
			"datastore_type": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The database type.",
			},

			// Attributes.
			"records": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "The list of email send records.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"send_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The send time, in RFC3339 format.",
						},
						"status": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "The send status.",
						},
						"email": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The email address.",
						},
						"topic_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The topic ID.",
						},
						"topic_urn": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The topic URN.",
						},
						"instance_health_reports": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "The list of instance health reports.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"task_id": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The report ID.",
									},
									"instance_id": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The instance ID.",
									},
									"instance_name": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The instance name.",
									},
									"start_time": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The diagnosis start time, in RFC3339 format.",
									},
									"end_time": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The diagnosis end time, in RFC3339 format.",
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

func dataSourceEmailSendRecordsRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	client, err := cfg.NewServiceClient("das", region)
	if err != nil {
		return diag.Errorf("error creating DAS client: %s", err)
	}

	emailRecords, err := listEmailSendRecords(client, d)
	if err != nil {
		return diag.Errorf("error querying DAS email send records: %s", err)
	}

	randomUUID, err := uuid.NewRandom()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(randomUUID.String())

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("records", flattenEmailSendRecords(emailRecords)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func listEmailSendRecords(client *golangsdk.ServiceClient, d *schema.ResourceData) ([]interface{}, error) {
	var (
		httpUrl = "v3/{project_id}/batch-inspection/email-record?limit={limit}"
		limit   = 100
		offset  = 0
		result  = make([]interface{}, 0)
	)

	listPath := client.Endpoint + httpUrl
	listPath = strings.ReplaceAll(listPath, "{project_id}", client.ProjectID)
	listPath = strings.ReplaceAll(listPath, "{limit}", strconv.Itoa(limit))
	listPath += buildEmailSendRecordsQueryParams(d)

	opt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type":     "application/json",
			"X-Source-Service": "das",
		},
	}

	for {
		listPathWithOffset := fmt.Sprintf("%s&offset=%d", listPath, offset)
		requestResp, err := client.Request("GET", listPathWithOffset, &opt)
		if err != nil {
			return nil, err
		}
		respBody, err := utils.FlattenResponse(requestResp)
		if err != nil {
			return nil, err
		}

		emailRecordList := utils.PathSearch("email_record_list", respBody, make([]interface{}, 0)).([]interface{})
		result = append(result, emailRecordList...)

		if len(emailRecordList) < limit {
			break
		}
		offset += len(emailRecordList)
	}

	return result, nil
}

func buildEmailSendRecordsQueryParams(d *schema.ResourceData) string {
	res := ""

	if v, ok := d.GetOk("datastore_type"); ok {
		res = fmt.Sprintf("%s&datastore_type=%v", res, v)
	}

	return res
}

func flattenEmailSendRecords(items []interface{}) []map[string]interface{} {
	if len(items) < 1 {
		return nil
	}

	result := make([]map[string]interface{}, 0, len(items))
	for _, item := range items {
		result = append(result, map[string]interface{}{
			"send_time": utils.FormatTimeStampRFC3339(
				int64(utils.PathSearch("send_at", item, float64(0)).(float64))/1000, false),
			"status":    utils.PathSearch("send_status", item, nil),
			"email":     utils.PathSearch("email", item, nil),
			"topic_id":  utils.PathSearch("topic", item, nil),
			"topic_urn": utils.PathSearch("topic_urn", item, nil),
			"instance_health_reports": flattenInstanceHealthReports(
				utils.PathSearch("instance_health_report_list", item, make([]interface{}, 0)).([]interface{})),
		})
	}

	return result
}

func flattenInstanceHealthReports(items []interface{}) []map[string]interface{} {
	if len(items) < 1 {
		return nil
	}

	result := make([]map[string]interface{}, 0, len(items))
	for _, item := range items {
		result = append(result, map[string]interface{}{
			"task_id":       utils.PathSearch("task_id", item, nil),
			"instance_id":   utils.PathSearch("instance_id", item, nil),
			"instance_name": utils.PathSearch("instance_name", item, nil),
			"start_time": utils.FormatTimeStampRFC3339(
				int64(utils.PathSearch("start_at", item, float64(0)).(float64))/1000, false),
			"end_time": utils.FormatTimeStampRFC3339(
				int64(utils.PathSearch("end_at", item, float64(0)).(float64))/1000, false),
		})
	}

	return result
}
