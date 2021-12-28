package dli

import (
	"context"
	"fmt"
	"time"

	"github.com/chnsz/golangsdk"
	"github.com/chnsz/golangsdk/openstack/dli/v1/sqljob"
	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils/fmtp"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils/logp"
)

func ResourceSqlJob() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceSqlJobCreate,
		ReadContext:   resourceSqlJobRead,
		DeleteContext: resourceSqlJobDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"sql": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"database_name": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"queue_name": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Computed: true,
			},
			"conf": {
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"spark_sql_max_records_per_file": {
							Type:     schema.TypeInt,
							Optional: true,
							ForceNew: true,
						},
						"spark_sql_auto_broadcast_join_threshold": {
							Type:     schema.TypeInt,
							Optional: true,
							ForceNew: true,
						},
						"spark_sql_shuffle_partitions": {
							Type:     schema.TypeInt,
							Optional: true,
							ForceNew: true,
						},
						"spark_sql_dynamic_partition_overwrite_enabled": {
							Type:     schema.TypeBool,
							Optional: true,
							ForceNew: true,
						},
						"spark_sql_files_max_partition_bytes": {
							Type:     schema.TypeInt,
							Optional: true,
							ForceNew: true,
						},
						"spark_sql_bad_records_path": {
							Type:     schema.TypeString,
							Optional: true,
							ForceNew: true,
						},
						"dli_sql_sqlasync_enabled": {
							Type:     schema.TypeBool,
							Optional: true,
							ForceNew: true,
						},
						"dli_sql_job_timeout": {
							Type:     schema.TypeInt,
							Optional: true,
							ForceNew: true,
						},
					},
				},
			},
			"tags": common.TagsForceNewSchema(),
			"job_type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"schema": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeMap,
					Elem: &schema.Schema{
						Type: schema.TypeString,
					},
				},
			},
			"rows": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeList,
					Elem: &schema.Schema{
						Type: schema.TypeString,
					},
				},
			},
			"owner": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"start_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"duration": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(20 * time.Minute),
			Delete: schema.DefaultTimeout(45 * time.Minute),
		},
	}
}

func resourceSqlJobCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*config.Config)
	region := config.GetRegion(d)
	client, err := config.DliV1Client(region)
	if err != nil {
		return diag.Errorf("error creating DLI v1 client, err=%s", err)
	}

	opts := sqljob.SqlJobOpts{
		Sql:       d.Get("sql").(string),
		Currentdb: d.Get("database_name").(string),
		QueueName: d.Get("queue_name").(string),
		Tags:      utils.ExpandResourceTags(d.Get("tags").(map[string]interface{})),
	}

	if _, ok := d.GetOk("conf"); ok {
		opts.Conf = buildConfParam(d)
	}

	logp.Printf("[DEBUG] Creating new DLI sql job opts: %#v", opts)
	rst, createErr := sqljob.Submit(client, opts)
	if createErr != nil {
		return fmtp.DiagErrorf("Error creating DLI sql job: %s", createErr)
	}

	if rst != nil && !rst.IsSuccess {
		return fmtp.DiagErrorf("Error creating DLI sql job: %s", rst.Message)
	}

	d.SetId(rst.JobId)
	d.Set("schema", rst.Schema)
	d.Set("rows", rst.Rows)

	errCheckRt := waitingforJobRunning(ctx, client, rst.JobId, d.Timeout(schema.TimeoutCreate))
	if errCheckRt != nil {
		return diag.FromErr(errCheckRt)
	}

	return resourceSqlJobRead(ctx, d, meta)
}

func resourceSqlJobRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*config.Config)
	region := config.GetRegion(d)
	client, err := config.DliV1Client(region)
	if err != nil {
		return diag.Errorf("error creating DLI v1 client, err=%s", err)
	}

	listResp, lErr := sqljob.List(client, sqljob.ListJobsOpts{
		JobId: d.Id(),
	})

	if lErr != nil {
		return fmtp.DiagErrorf("Error query DLI sql job %q:%s", d.Id(), lErr)
	}

	if listResp != nil && !listResp.IsSuccess && listResp.JobCount != 1 {
		return fmtp.DiagErrorf("Error query DLI sql job: %s", listResp.Message)
	}
	dt := listResp.Jobs[0]
	mErr := multierror.Append(
		d.Set("sql", dt.Statement),
		d.Set("database_name", dt.DatabaseName),
		d.Set("queue_name", dt.QueueName),
		d.Set("job_type", dt.JobType),
		d.Set("owner", dt.Owner),
		d.Set("start_time", utils.FormatTimeStampRFC3339(int64(dt.StartTime))),
		d.Set("duration", dt.Duration),
		d.Set("status", dt.Status),
		d.Set("tags", utils.TagsToMap(dt.Tags)),
	)
	if setSdErr := mErr.ErrorOrNil(); setSdErr != nil {
		return fmtp.DiagErrorf("Error setting vault fields: %s", setSdErr)
	}

	return nil
}

// This API is used to cancel a submitted job. If execution of a job completes or fails, this job cannot be canceled.
func resourceSqlJobDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*config.Config)
	region := config.GetRegion(d)
	client, err := config.DliV1Client(region)
	if err != nil {
		return diag.Errorf("error creating DLI v1 client, err=%s", err)
	}
	jobId := d.Id()
	detail, dErr := sqljob.Status(client, jobId)
	if dErr != nil {
		return fmtp.DiagErrorf("Error query DLI sql job %q:%s", jobId, dErr)
	}

	if detail != nil && !detail.IsSuccess {
		return fmtp.DiagErrorf("Error query DLI sql job: %s", detail.Message)
	}

	if detail.Status != sqljob.JobStatusFinished &&
		detail.Status != sqljob.JobStatusFailed &&
		detail.Status != sqljob.JobStatusCancelled {

		cancelRst, cancelErr := sqljob.Cancel(client, jobId)
		if cancelErr != nil {
			return fmtp.DiagErrorf("cancel DLI sql job failed. %q:%s", jobId, dErr)
		}
		if cancelRst != nil && !cancelRst.IsSuccess {
			return fmtp.DiagErrorf("cancel DLI sql job failed. %q:%s", jobId, dErr)
		}

	}

	errCheckRt := checkSqlJobCancelledResult(ctx, client, jobId, d.Timeout(schema.TimeoutDelete))
	if errCheckRt != nil {
		return fmtp.DiagErrorf("Failed to check the result of deletion %s", errCheckRt)
	}

	d.SetId("")
	return nil
}

func buildConfParam(d *schema.ResourceData) []string {
	var rt []string

	if v, ok := d.GetOk("conf.0.spark_sql_max_records_per_file"); ok {
		rt = append(rt, fmt.Sprint(sqljob.ConfigSparkSqlFilesMaxRecordsPerFile, "=", v))
	}
	if v, ok := d.GetOk("conf.0.spark_sql_auto_broadcast_join_threshold"); ok {
		rt = append(rt, fmt.Sprint(sqljob.ConfigSparkSqlAutoBroadcastJoinThreshold, "=", v))
	}
	if v, ok := d.GetOk("conf.0.spark_sql_shuffle_partitions"); ok {
		rt = append(rt, fmt.Sprint(sqljob.ConfigSparkSqlShufflePartitions, "=", v))
	}
	if v, ok := d.GetOk("conf.0.spark_sql_dynamic_partition_overwrite_enabled"); ok {
		rt = append(rt, fmt.Sprint(sqljob.ConfigSparkSqlDynamicPartitionOverwriteEnabled, "=", v))
	}
	if v, ok := d.GetOk("conf.0.spark_sql_files_max_partition_bytes"); ok {
		rt = append(rt, fmt.Sprint(sqljob.ConfigSparkSqlMaxPartitionBytes, "=", v))
	}
	if v, ok := d.GetOk("conf.0.spark_sql_bad_records_path"); ok {
		rt = append(rt, fmt.Sprint(sqljob.ConfigSparkSqlBadRecordsPath, "=", v))
	}
	if v, ok := d.GetOk("conf.0.dli_sql_sqlasync_enabled"); ok {
		rt = append(rt, fmt.Sprint(sqljob.ConfigDliSqlasyncEnabled, "=", v))
	}
	if v, ok := d.GetOk("conf.0.dli_sql_job_timeout"); ok {
		rt = append(rt, fmt.Sprint(sqljob.ConfigDliSqljobTimeout, "=", v))
	}

	return rt
}

func checkSqlJobCancelledResult(ctx context.Context, client *golangsdk.ServiceClient, id string,
	timeout time.Duration) error {
	stateConf := &resource.StateChangeConf{
		Pending: []string{"Pending"},
		Target:  []string{"Done"},
		Refresh: func() (interface{}, string, error) {
			jobStatus, err := sqljob.Status(client, id)
			if err == nil {
				if jobStatus.Status == sqljob.JobStatusCancelled ||
					jobStatus.Status == sqljob.JobStatusFinished ||
					jobStatus.Status == sqljob.JobStatusFailed {
					return true, "Done", nil
				}
			}
			if err != nil {
				if _, ok := err.(golangsdk.ErrDefault400); ok {
					return true, "Done", nil
				}
				return nil, "", nil
			}
			return true, "Pending", nil
		},
		Timeout:      timeout,
		PollInterval: 20 * timeout,
		Delay:        20 * time.Second,
	}
	_, err := stateConf.WaitForStateContext(ctx)
	if err != nil {
		return fmtp.Errorf("error waiting for Dli sql job (%s) to be cancelled: %s", id, err)
	}
	return nil
}

func waitingforJobRunning(ctx context.Context, client *golangsdk.ServiceClient, id string, timeout time.Duration) error {
	stateConf := &resource.StateChangeConf{
		Pending: []string{"Pending"},
		Target:  []string{"Done"},
		Refresh: func() (interface{}, string, error) {
			jobStatus, err := sqljob.Status(client, id)
			if err != nil {
				return nil, "failed", err
			}

			if jobStatus.Status == sqljob.JobStatusLaunching {
				return true, "Pending", nil
			}

			if jobStatus.Status == sqljob.JobStatusCancelled || jobStatus.Status == sqljob.JobStatusFailed {
				return true, "failed", fmtp.Errorf("current status is %s", jobStatus.Status)
			}

			return true, "Done", nil
		},
		Timeout:      timeout,
		PollInterval: 20 * timeout,
		Delay:        20 * time.Second,
	}
	_, err := stateConf.WaitForStateContext(ctx)
	if err != nil {
		return fmtp.Errorf("error waiting for Dli sql job (%s) to be running: %s", id, err)
	}
	return nil
}
