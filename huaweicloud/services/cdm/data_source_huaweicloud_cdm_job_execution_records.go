package cdm

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

// @API CDM GET /v1.1/{project_id}/clusters/{cluster_id}/cdm/submissions
func DataSourceCdmJobExecutionRecords() *schema.Resource {
	return &schema.Resource{
		ReadContext: resourceCdmJobExecutionRecordsRead,
		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"cluster_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `Specifies the cluster ID.`,
			},
			"job_name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `Specifies the job name.`,
			},
			"records": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: `Indicates the records.`,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"is_incrementing": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: `Indicates whether the job migrates incremental data.`,
						},
						"is_stoping_increment": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `Indicates whether incremental data migration stopped.`,
						},
						"is_execute_auto": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: `Indicates whether the job executed as scheduled.`,
						},
						"last_update_date": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `Indicates the time when the job was last updated.`,
						},
						"last_udpate_user": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `Indicates the user who last updated the job status.`,
						},
						"is_delete_job": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: `Indicates whether the job to be deleted after it is executed.`,
						},
						"creation_user": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `Indicates the user who created the job.`,
						},
						"creation_date": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `Indicates the creation time.`,
						},
						"external_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `Indicates the job ID.`,
						},
						"progress": {
							Type:        schema.TypeFloat,
							Computed:    true,
							Description: `Indicates the Job progress.`,
						},
						"submission_id": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: `Indicates the job submission ID.`,
						},
						"delete_rows": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: `Indicates the number of deleted rows.`,
						},
						"update_rows": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: `Indicates the number of updated rows.`,
						},
						"write_rows": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: `Indicates the number of write rows.`,
						},
						"execute_date": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `Indicates the execution time.`,
						},
						"status": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `Indicates the Job status.`,
						},
						"error_details": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `Indicates the error details.`,
						},
						"error_summary": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `Indicates the error summary.`,
						},
						"counters": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: `Indicates the job running result statistics. Only return when status is SUCCEEDED.`,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"bytes_written": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: `Indicates the number of bytes that are written.`,
									},
									"bytes_read": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: `Indicates the number of bytes that are read.`,
									},
									"total_files": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: `Indicates the total number of files.`,
									},
									"rows_read": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: `Indicates the number of rows that are read.`,
									},
									"rows_written": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: `Indicates the number of rows that are written.`,
									},
									"files_written": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: `Indicates the number of files that are written.`,
									},
									"files_read": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: `Indicates the number of files that are read.`,
									},
									"total_size": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: `Indicates the total number of bytes.`,
									},
									"file_skipped": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: `Indicates the number of files that are skipped.`,
									},
									"rows_written_skipped": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: `Indicates the number of rows that are skipped.`,
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

func resourceCdmJobExecutionRecordsRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	client, err := cfg.NewServiceClient("cdm", region)
	if err != nil {
		return diag.Errorf("error creating CDM client: %s", err)
	}

	listJobExecutionRecordsHttpUrl := "v1.1/{project_id}/clusters/{cluster_id}/cdm/submissions?jname={job_name}"
	listJobExecutionRecordsPath := client.Endpoint + listJobExecutionRecordsHttpUrl
	listJobExecutionRecordsPath = strings.ReplaceAll(listJobExecutionRecordsPath, "{project_id}", client.ProjectID)
	listJobExecutionRecordsPath = strings.ReplaceAll(listJobExecutionRecordsPath, "{cluster_id}", d.Get("cluster_id").(string))
	listJobExecutionRecordsPath = strings.ReplaceAll(listJobExecutionRecordsPath, "{job_name}", d.Get("job_name").(string))

	listJobExecutionRecordsOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"Content-Type": "application/json"},
	}

	listJobExecutionRecordsResp, err := client.Request("GET", listJobExecutionRecordsPath, &listJobExecutionRecordsOpt)
	if err != nil {
		return diag.FromErr(err)
	}
	listJobExecutionRecordsRespBody, err := utils.FlattenResponse(listJobExecutionRecordsResp)
	if err != nil {
		return diag.FromErr(err)
	}

	uuid, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(uuid)

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("records", flattenJobExecutionRecords(listJobExecutionRecordsRespBody)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenJobExecutionRecords(resp interface{}) []interface{} {
	if resp == nil {
		return nil
	}

	curJson := utils.PathSearch("submissions", resp, make([]interface{}, 0))
	curArray := curJson.([]interface{})
	rst := make([]interface{}, 0, len(curArray))
	for _, v := range curArray {
		rst = append(rst, map[string]interface{}{
			"is_incrementing":      utils.PathSearch("isIncrementing", v, nil),
			"is_stoping_increment": utils.PathSearch("isStopingIncrement", v, nil),
			"is_execute_auto":      utils.PathSearch(`"is-execute-auto"`, v, nil),
			"last_update_date": utils.FormatTimeStampRFC3339(
				int64(utils.PathSearch(`"last-update-date"`, v, float64(0)).(float64))/1000, true),
			"last_udpate_user": utils.PathSearch(`"last-udpate-user"`, v, nil),
			"is_delete_job":    utils.PathSearch("isDeleteJob", v, nil),
			"creation_user":    utils.PathSearch(`"creation-user"`, v, nil),
			"creation_date": utils.FormatTimeStampRFC3339(
				int64(utils.PathSearch(`"creation-date"`, v, float64(0)).(float64))/1000, true),
			"external_id":   utils.PathSearch(`"external-id"`, v, nil),
			"progress":      utils.PathSearch("progress", v, nil),
			"submission_id": utils.PathSearch(`"submission-id"`, v, nil),
			"delete_rows":   utils.PathSearch("delete_rows", v, nil),
			"update_rows":   utils.PathSearch("update_rows", v, nil),
			"write_rows":    utils.PathSearch("write_rows", v, nil),
			"execute_date": utils.FormatTimeStampRFC3339(
				int64(utils.PathSearch(`"execute-date"`, v, float64(0)).(float64))/1000, true),
			"status":        utils.PathSearch("status", v, nil),
			"error_details": utils.PathSearch(`"error-details"`, v, nil),
			"error_summary": utils.PathSearch(`"error-summary"`, v, nil),
			"counters": flattenCounters(
				utils.PathSearch(`counters."org.apache.sqoop.submission.counter.SqoopCounters"`, v, nil)),
		})
	}

	return rst
}

func flattenCounters(curJson interface{}) interface{} {
	if curJson == nil {
		return nil
	}

	rst := map[string]interface{}{
		"bytes_written":        utils.PathSearch("BYTES_WRITTEN", curJson, nil),
		"bytes_read":           utils.PathSearch("BYTES_READ", curJson, nil),
		"total_files":          utils.PathSearch("TOTAL_FILES", curJson, nil),
		"rows_read":            utils.PathSearch("ROWS_READ", curJson, nil),
		"rows_written":         utils.PathSearch("ROWS_WRITTEN", curJson, nil),
		"files_written":        utils.PathSearch("FILES_WRITTEN", curJson, nil),
		"files_read":           utils.PathSearch("FILES_READ", curJson, nil),
		"total_size":           utils.PathSearch("TOTAL_SIZE", curJson, nil),
		"file_skipped":         utils.PathSearch("FILES_SKIPPED", curJson, nil),
		"rows_written_skipped": utils.PathSearch("ROWS_WRITTEN_SKIPPED", curJson, nil),
	}

	return []interface{}{rst}
}
