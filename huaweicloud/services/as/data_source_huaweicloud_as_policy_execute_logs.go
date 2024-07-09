package as

import (
	"context"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk/openstack/autoscaling/v1/policyexecutelogs"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
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

func buildDataSourcePolicyExecuteLogOpts(d *schema.ResourceData) policyexecutelogs.ListOpts {
	return policyexecutelogs.ListOpts{
		PolicyID:     d.Get("scaling_policy_id").(string),
		LogID:        d.Get("log_id").(string),
		ResourceID:   d.Get("scaling_resource_id").(string),
		ResourceType: d.Get("scaling_resource_type").(string),
		ExecuteType:  d.Get("execute_type").(string),
		StartTime:    d.Get("start_time").(string),
		EndTime:      d.Get("end_time").(string),
	}
}

func dataSourcePolicyExecuteLogsRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg    = meta.(*config.Config)
		region = cfg.GetRegion(d)
		opts   = buildDataSourcePolicyExecuteLogOpts(d)
	)
	client, err := cfg.AutoscalingV1Client(region)
	if err != nil {
		return diag.Errorf("error creating AS v1 client: %s", err)
	}

	executeLogList, err := policyexecutelogs.List(client, opts)
	if err != nil {
		return diag.Errorf("error retrieving AS policy execute logs: %s", err)
	}

	randUUID, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}

	d.SetId(randUUID)
	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("execute_logs", flattenDataSourcePolicyExecuteLogs(executeLogList, d)),
	)
	if mErr.ErrorOrNil() != nil {
		return diag.Errorf("error saving AS policy execute logs data source fields: %s", mErr)
	}
	return nil
}

func flattenDataSourcePolicyExecuteLogs(executeLogList []policyexecutelogs.ExecuteLog, d *schema.ResourceData) []map[string]interface{} {
	executeLogs := make([]map[string]interface{}, 0, len(executeLogList))
	for _, executeLog := range executeLogList {
		if val, ok := d.GetOk("status"); ok && val.(string) != executeLog.Status {
			continue
		}
		executeLogMap := map[string]interface{}{
			"id":                    executeLog.ID,
			"status":                executeLog.Status,
			"failed_reason":         executeLog.FailedReason,
			"execute_type":          executeLog.ExecuteType,
			"execute_time":          executeLog.ExecuteTime,
			"scaling_policy_id":     executeLog.PolicyID,
			"scaling_resource_id":   executeLog.ResourceID,
			"scaling_resource_type": executeLog.ResourceType,
			"type":                  executeLog.Type,
			"old_value":             executeLog.OldValue,
			"desire_value":          executeLog.DesireValue,
			"limit_value":           executeLog.LimitValue,
			"job_records":           flattenJobRecords(executeLog.JobRecords),
			"metadata":              executeLog.MetaData,
		}
		executeLogs = append(executeLogs, executeLogMap)
	}
	return executeLogs
}

func flattenJobRecords(jobRecords []policyexecutelogs.JobRecord) []map[string]interface{} {
	if len(jobRecords) == 0 {
		return nil
	}

	jobRecordList := make([]map[string]interface{}, 0, len(jobRecords))
	for _, jobRecord := range jobRecords {
		job := map[string]interface{}{
			"job_name":    jobRecord.JobName,
			"job_status":  jobRecord.JobStatus,
			"record_type": jobRecord.RecordType,
			"record_time": jobRecord.RecordTime,
			"request":     jobRecord.Request,
			"response":    jobRecord.Response,
			"code":        jobRecord.Code,
			"message":     jobRecord.Message,
		}
		jobRecordList = append(jobRecordList, job)
	}

	return jobRecordList
}
