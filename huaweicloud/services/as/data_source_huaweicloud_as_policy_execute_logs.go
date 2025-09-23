package as

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

// @API AS GET /autoscaling-api/v1/{project_id}/scaling_policy_execute_log/{scaling_policy_id}
func DataSourcePolicyExecuteLogs() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourcePolicyExecuteLogsRead,
		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"scaling_policy_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"log_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"scaling_resource_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"scaling_resource_type": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"execute_type": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"start_time": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"end_time": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"status": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"execute_logs": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"status": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"failed_reason": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"execute_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"execute_time": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"scaling_policy_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"scaling_resource_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"scaling_resource_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"old_value": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"desire_value": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"limit_value": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"job_records": {
							Type:     schema.TypeList,
							Computed: true,
							Elem:     policyExecuteLogDataSourceJobRecordsSchema(),
						},
						"metadata": {
							Type:     schema.TypeMap,
							Computed: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
					},
				},
			},
		},
	}
}

func policyExecuteLogDataSourceJobRecordsSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"job_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"job_status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"record_type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"record_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"request": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"response": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"code": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"message": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func buildPolicyExecuteLogsQueryParams(d *schema.ResourceData) string {
	queryParam := ""
	if logID := d.Get("log_id").(string); logID != "" {
		queryParam = fmt.Sprintf("%s&log_id=%s", queryParam, logID)
	}

	if scalingResourceID := d.Get("scaling_resource_id").(string); scalingResourceID != "" {
		queryParam = fmt.Sprintf("%s&scaling_resource_id=%s", queryParam, scalingResourceID)
	}

	if scalingResourceType := d.Get("scaling_resource_type").(string); scalingResourceType != "" {
		queryParam = fmt.Sprintf("%s&scaling_resource_type=%s", queryParam, scalingResourceType)
	}

	if executeType := d.Get("execute_type").(string); executeType != "" {
		queryParam = fmt.Sprintf("%s&execute_type=%s", queryParam, executeType)
	}

	if startTime := d.Get("start_time").(string); startTime != "" {
		queryParam = fmt.Sprintf("%s&start_time=%s", queryParam, startTime)
	}

	if endTime := d.Get("end_time").(string); endTime != "" {
		queryParam = fmt.Sprintf("%s&end_time=%s", queryParam, endTime)
	}

	if queryParam == "" {
		return ""
	}

	return fmt.Sprintf("?%s", queryParam[1:])
}

// There is currently a problem with openAPI paging, so paging query cannot be implemented.
func dataSourcePolicyExecuteLogsRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg             = meta.(*config.Config)
		region          = cfg.GetRegion(d)
		product         = "autoscaling"
		httpUrl         = "autoscaling-api/v1/{project_id}/scaling_policy_execute_log/{scaling_policy_id}"
		scalingPolicyID = d.Get("scaling_policy_id").(string)
	)

	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating AS client: %s", err)
	}

	requestPath := client.Endpoint + httpUrl
	requestPath = strings.ReplaceAll(requestPath, "{project_id}", client.ProjectID)
	requestPath = strings.ReplaceAll(requestPath, "{scaling_policy_id}", scalingPolicyID)
	requestPath += buildPolicyExecuteLogsQueryParams(d)
	requestOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"Content-Type": "application/json"},
	}

	resp, err := client.Request("GET", requestPath, &requestOpt)
	if err != nil {
		return diag.Errorf("error retrieving AS policy execute logs: %s", err)
	}

	respBody, err := utils.FlattenResponse(resp)
	if err != nil {
		return diag.FromErr(err)
	}

	executeLogsResp := utils.PathSearch("scaling_policy_execute_log", respBody, make([]interface{}, 0)).([]interface{})
	randUUID, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}

	d.SetId(randUUID)
	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("execute_logs", flattenDataSourcePolicyExecuteLogs(executeLogsResp, d)),
	)
	if mErr.ErrorOrNil() != nil {
		return diag.Errorf("error saving AS policy execute logs data source fields: %s", mErr)
	}
	return nil
}

func flattenDataSourcePolicyExecuteLogs(executeLogsResp []interface{}, d *schema.ResourceData) []map[string]interface{} {
	executeLogs := make([]map[string]interface{}, 0, len(executeLogsResp))
	for _, executeLog := range executeLogsResp {
		status := utils.PathSearch("status", executeLog, "").(string)
		if val, ok := d.GetOk("status"); ok && val.(string) != status {
			continue
		}

		executeLogMap := map[string]interface{}{
			"id":                    utils.PathSearch("id", executeLog, nil),
			"status":                status,
			"failed_reason":         utils.PathSearch("failed_reason", executeLog, nil),
			"execute_type":          utils.PathSearch("execute_type", executeLog, nil),
			"execute_time":          utils.PathSearch("execute_time", executeLog, nil),
			"scaling_policy_id":     utils.PathSearch("scaling_policy_id", executeLog, nil),
			"scaling_resource_id":   utils.PathSearch("scaling_resource_id", executeLog, nil),
			"scaling_resource_type": utils.PathSearch("scaling_resource_type", executeLog, nil),
			"type":                  utils.PathSearch("type", executeLog, nil),
			"old_value":             utils.PathSearch("old_value", executeLog, nil),
			"desire_value":          utils.PathSearch("desire_value", executeLog, nil),
			"limit_value":           utils.PathSearch("limit_value", executeLog, nil),
			"job_records":           flattenJobRecords(utils.PathSearch("job_records", executeLog, make([]interface{}, 0)).([]interface{})),
			"metadata":              utils.PathSearch("meta_data", executeLog, nil),
		}
		executeLogs = append(executeLogs, executeLogMap)
	}
	return executeLogs
}

func flattenJobRecords(jobRecords []interface{}) []map[string]interface{} {
	if len(jobRecords) == 0 {
		return nil
	}

	jobRecordList := make([]map[string]interface{}, 0, len(jobRecords))
	for _, jobRecord := range jobRecords {
		job := map[string]interface{}{
			"job_name":    utils.PathSearch("job_name", jobRecord, nil),
			"job_status":  utils.PathSearch("job_status", jobRecord, nil),
			"record_type": utils.PathSearch("record_type", jobRecord, nil),
			"record_time": utils.PathSearch("record_time", jobRecord, nil),
			"request":     utils.PathSearch("request", jobRecord, nil),
			"response":    utils.PathSearch("response", jobRecord, nil),
			"code":        utils.PathSearch("code", jobRecord, nil),
			"message":     utils.PathSearch("message", jobRecord, nil),
		}
		jobRecordList = append(jobRecordList, job)
	}

	return jobRecordList
}
