package dli

import (
	"context"
	"encoding/json"
	"strconv"
	"time"

	"github.com/chnsz/golangsdk"
	"github.com/chnsz/golangsdk/openstack/common/tags"
	"github.com/chnsz/golangsdk/openstack/dli/v1/flinkjob"
	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils/fmtp"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils/logp"
)

func ResourceFlinkSqlJob() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceFlinkSqlJobCreate,
		ReadContext:   resourceFlinkSqlJobRead,
		UpdateContext: resourceFlinkSqlJobUpdate,
		DeleteContext: resourceFlinkSqlJobDelete,
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

			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringLenBetween(1, 57),
			},

			"type": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  flinkjob.JobTypeFlinkSql,
				ForceNew: true,
			},

			"run_mode": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  flinkjob.RunModeSharedCluster,
			},

			"description": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringLenBetween(1, 512),
			},

			"queue_name": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"sql": {
				Type:     schema.TypeString,
				Optional: true,
			},

			"cu_number": {
				Type:     schema.TypeInt,
				Optional: true,
				Default:  2,
			},

			"parallel_number": {
				Type:     schema.TypeInt,
				Optional: true,
				Default:  1,
			},

			"checkpoint_enabled": {
				Type:         schema.TypeBool,
				Optional:     true,
				Default:      false,
				RequiredWith: []string{"obs_bucket"},
			},

			"checkpoint_mode": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "exactly_once",
				ValidateFunc: validation.StringInSlice(
					[]string{flinkjob.CheckpointModeExactlyOnce, flinkjob.CheckpointModeAtLeastOnce}, true),
			},

			"checkpoint_interval": {
				Type:     schema.TypeInt,
				Optional: true,
				Default:  10,
			},

			"obs_bucket": {
				Type:     schema.TypeString,
				Optional: true,
			},

			"log_enabled": {
				Type:         schema.TypeBool,
				Optional:     true,
				Default:      false,
				RequiredWith: []string{"obs_bucket"},
			},

			"smn_topic": {
				Type:     schema.TypeString,
				Optional: true,
			},

			"restart_when_exception": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},

			"idle_state_retention": {
				Type:     schema.TypeInt,
				Optional: true,
				Default:  1,
			},

			"edge_group_ids": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},

			"dirty_data_strategy": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "0",
			},
			"udf_jar_url": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"manager_cu_number": {
				Type:     schema.TypeInt,
				Optional: true,
				Default:  1,
			},

			"tm_cus": {
				Type:     schema.TypeInt,
				Optional: true,
				Default:  1,
			},

			"tm_slot_num": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},

			"resume_checkpoint": {
				Type:     schema.TypeBool,
				Optional: true,
			},

			"resume_max_num": {
				Type:     schema.TypeInt,
				Optional: true,
				Default:  -1,
			},

			"runtime_config": common.TagsSchema(),

			"tags": common.TagsForceNewSchema(),

			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
			Update: schema.DefaultTimeout(20 * time.Minute),
			Delete: schema.DefaultTimeout(20 * time.Minute),
		},
	}
}

func resourceFlinkSqlJobCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*config.Config)
	region := config.GetRegion(d)
	client, err := config.DliV1Client(region)
	if err != nil {
		return diag.Errorf("error creating DLI v1 client, err=%s", err)
	}

	aErr := authorizeObsBucket(client, d)
	if aErr != nil {
		return aErr
	}

	opts := flinkjob.CreateSqlJobOpts{
		Name:                 d.Get("name").(string),
		JobType:              d.Get("type").(string),
		RunMode:              d.Get("run_mode").(string),
		Desc:                 d.Get("description").(string),
		QueueName:            d.Get("queue_name").(string),
		SqlBody:              d.Get("sql").(string),
		CuNumber:             golangsdk.IntToPointer(d.Get("cu_number").(int)),
		ParallelNumber:       golangsdk.IntToPointer(d.Get("parallel_number").(int)),
		CheckpointEnabled:    utils.Bool(d.Get("checkpoint_enabled").(bool)),
		CheckpointInterval:   golangsdk.IntToPointer(d.Get("checkpoint_interval").(int)),
		ObsBucket:            d.Get("obs_bucket").(string),
		LogEnabled:           utils.Bool(d.Get("log_enabled").(bool)),
		SmnTopic:             d.Get("smn_topic").(string),
		RestartWhenException: utils.Bool(d.Get("restart_when_exception").(bool)),
		IdleStateRetention:   golangsdk.IntToPointer(d.Get("idle_state_retention").(int)),
		DirtyDataStrategy:    d.Get("dirty_data_strategy").(string),
		UdfJarUrl:            d.Get("udf_jar_url").(string),
		ManagerCuNumber:      golangsdk.IntToPointer(d.Get("manager_cu_number").(int)),
		TmCus:                golangsdk.IntToPointer(d.Get("tm_cus").(int)),
		TmSlotNum:            golangsdk.IntToPointer(d.Get("tm_slot_num").(int)),
		ResumeCheckpoint:     utils.Bool(d.Get("resume_checkpoint").(bool)),
		ResumeMaxNum:         golangsdk.IntToPointer(d.Get("resume_max_num").(int)),
		Tags:                 utils.ExpandResourceTags(d.Get("tags").(map[string]interface{})),
	}

	if mode := d.Get("checkpoint_mode").(string); mode == flinkjob.CheckpointModeAtLeastOnce {
		opts.CheckpointMode = golangsdk.IntToPointer(2)
	} else {
		opts.CheckpointMode = golangsdk.IntToPointer(1)
	}

	if edgeGroupIds, ok := d.GetOk("edge_group_ids"); ok {
		var ids []string
		for _, v := range edgeGroupIds.([]interface{}) {
			ids = append(ids, v.(string))
		}
		opts.EdgeGroupIds = ids
	}

	if runtimConfig, ok := d.GetOk("runtime_config"); ok {
		config := utils.ExpandResourceTags(runtimConfig.(map[string]interface{}))
		configStr, _ := json.Marshal(config)
		opts.RuntimeConfig = string(configStr)
	}

	logp.Printf("[DEBUG] Creating new DLI flink job opts: %#v", opts)

	rst, createErr := flinkjob.CreateSqlJob(client, opts)
	if createErr != nil {
		return fmtp.DiagErrorf("Error creating DLI flink job: %s", createErr)
	}

	if rst != nil && !rst.IsSuccess {
		return fmtp.DiagErrorf("Error creating DLI flink job: %s", rst.Message)
	}

	d.SetId(strconv.Itoa(rst.Job.JobId))

	// run the flink job
	_, runErr := flinkjob.Run(client, flinkjob.RunJobOpts{
		JobIds:          []int{rst.Job.JobId},
		ResumeSavepoint: utils.Bool(false),
	})
	if runErr != nil {
		return fmtp.DiagErrorf("Error run DLI flink job: %s", runErr)
	}

	checkCreateErr := checkFlinkJobRunResult(ctx, client, rst.Job.JobId, d.Timeout(schema.TimeoutCreate))
	if checkCreateErr != nil {
		return diag.FromErr(checkCreateErr)
	}
	return resourceFlinkSqlJobRead(ctx, d, meta)
}

func resourceFlinkSqlJobRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*config.Config)
	region := config.GetRegion(d)
	client, err := config.DliV1Client(region)
	if err != nil {
		return diag.Errorf("error creating DLI v1 client, err=%s", err)
	}
	id, aErr := strconv.Atoi(d.Id())
	if aErr != nil {
		return fmtp.DiagErrorf("the DLI flink job_id must be number. actual id=%s", d.Id())
	}

	detailRsp, dErr := flinkjob.Get(client, id)
	if dErr != nil {
		return fmtp.DiagErrorf("Error query DLI flink job %q:%s", id, dErr)
	}

	if detailRsp != nil && !detailRsp.IsSuccess {
		return fmtp.DiagErrorf("Error query DLI flink job: %s", detailRsp.Message)
	}
	detail := detailRsp.JobDetail
	mErr := multierror.Append(
		d.Set("name", detail.Name),
		d.Set("type", detail.JobType),
		d.Set("run_mode", detail.RunMode),
		d.Set("description", detail.Desc),
		d.Set("queue_name", detail.QueueName),
		d.Set("sql", detail.SqlBody),
		d.Set("cu_number", detail.JobConfig.CuNumber),
		d.Set("parallel_number", detail.JobConfig.ParallelNumber),
		d.Set("checkpoint_enabled", detail.JobConfig.CheckpointEnabled),
		d.Set("checkpoint_mode", detail.JobConfig.CheckpointMode),
		d.Set("checkpoint_interval", detail.JobConfig.CheckpointInterval),
		d.Set("obs_bucket", detail.JobConfig.ObsBucket),
		d.Set("log_enabled", detail.JobConfig.LogEnabled),
		d.Set("smn_topic", detail.JobConfig.SmnTopic),
		d.Set("restart_when_exception", detail.JobConfig.RestartWhenException),
		d.Set("idle_state_retention", detail.JobConfig.IdleStateRetention),
		d.Set("edge_group_ids", detail.JobConfig.EdgeGroupIds),
		d.Set("dirty_data_strategy", detail.JobConfig.DirtyDataStrategy),
		d.Set("udf_jar_url", detail.JobConfig.UdfJarUrl),
		d.Set("manager_cu_number", detail.JobConfig.ManagerCuNumber),
		d.Set("tm_cus", detail.JobConfig.TmCus),
		d.Set("tm_slot_num", detail.JobConfig.TmSlotNum),
		d.Set("resume_checkpoint", detail.JobConfig.ResumeCheckpoint),
		d.Set("resume_max_num", detail.JobConfig.ResumeMaxNum),
		d.Set("runtime_config", utils.TagsToMap(parseConfig(detail.JobConfig.RuntimeConfig))),
		d.Set("status", detail.Status),
	)
	if setSdErr := mErr.ErrorOrNil(); setSdErr != nil {
		return fmtp.DiagErrorf("Error setting vault fields: %s", setSdErr)
	}

	return nil
}

// This API is used to cancel a submitted job. If execution of a job completes or fails, this job cannot be canceled.
func resourceFlinkSqlJobDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*config.Config)
	region := config.GetRegion(d)
	client, err := config.DliV1Client(region)
	if err != nil {
		return diag.Errorf("error creating DLI v1 client, err=%s", err)
	}

	jobId, aErr := strconv.Atoi(d.Id())
	if aErr != nil {
		return fmtp.DiagErrorf("the DLI flink job_id must be number. actual id=%s", d.Id())
	}

	deleteRst, dErr := flinkjob.Delete(client, jobId)
	if dErr != nil {
		return fmtp.DiagErrorf("delete DLI flink job failed. %q:%s", jobId, dErr)
	}
	if deleteRst != nil && !deleteRst.IsSuccess {
		return fmtp.DiagErrorf("delete DLI flink job failed. %q:%s", jobId, dErr)
	}

	d.SetId("")

	return nil
}

func resourceFlinkSqlJobUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*config.Config)
	region := config.GetRegion(d)
	client, err := config.DliV1Client(region)
	if err != nil {
		return diag.Errorf("error creating DLI v1 client, err=%s", err)
	}

	jobId, iErr := strconv.Atoi(d.Id())
	if iErr != nil {
		return fmtp.DiagErrorf("the DLI flink job_id must be number. actual id=%s", d.Id())
	}

	aErr := authorizeObsBucket(client, d)
	if aErr != nil {
		return aErr
	}

	uDiagErr := updateFlinkSqlJobInRunning(client, jobId, d)
	if uDiagErr != nil {
		return uDiagErr
	}

	uDiagErr = updateFlinkSqlJobWithStop(ctx, client, jobId, d)
	if uDiagErr != nil {
		return uDiagErr
	}

	checkCreateErr := checkFlinkJobRunResult(ctx, client, jobId, d.Timeout(schema.TimeoutCreate))
	if checkCreateErr != nil {
		return diag.FromErr(checkCreateErr)
	}
	return resourceFlinkSqlJobRead(ctx, d, meta)
}

//updated in "job_running": smn_topic,restart_when_exception,resume_checkpoint,resume_max_num,checkpoint_path,obs_bucket
func updateFlinkSqlJobInRunning(client *golangsdk.ServiceClient, jobId int, d *schema.ResourceData) diag.Diagnostics {
	if d.HasChanges("smn_topic", "restart_when_exception", "resume_checkpoint", "resume_max_num", "obs_bucket") {
		opts := flinkjob.UpdateSqlJobOpts{
			ObsBucket:            d.Get("obs_bucket").(string),
			SmnTopic:             d.Get("smn_topic").(string),
			RestartWhenException: utils.Bool(d.Get("restart_when_exception").(bool)),
			ResumeCheckpoint:     utils.Bool(d.Get("resume_checkpoint").(bool)),
			ResumeMaxNum:         golangsdk.IntToPointer(d.Get("resume_max_num").(int)),
		}

		logp.Printf("[DEBUG] update DLI flink job opts: %#v", opts)

		rst, uErr := flinkjob.UpdateSqlJob(client, jobId, opts)
		if uErr != nil {
			return fmtp.DiagErrorf("Error update DLI flink job=%d: %s", jobId, uErr)
		}

		if rst != nil && !rst.IsSuccess {
			return fmtp.DiagErrorf("Error update DLI flink job=%d: %s", rst.Message)
		}
	}

	return nil
}

func checkFlinkJobRunResult(ctx context.Context, client *golangsdk.ServiceClient, id int,
	timeout time.Duration) error {
	stateConf := &resource.StateChangeConf{
		Pending: []string{"job_init", "job_submitting"},
		Target:  []string{"job_running", "job_finish"},
		Refresh: func() (interface{}, string, error) {
			job, err := flinkjob.Get(client, id)
			logp.Printf("[DEBUG] the flink job info in create check func: %#v,%s", job, err)
			if err != nil {
				return nil, "", err
			}
			if job.JobDetail.Status == "job_submit_fail" {
				return job, "failed", fmtp.Errorf("%s:%s", job.JobDetail.Status, job.JobDetail.StatusDesc)
			}
			return job, job.JobDetail.Status, nil
		},
		Timeout:      timeout,
		PollInterval: 20 * timeout,
		Delay:        20 * time.Second,
	}
	_, err := stateConf.WaitForStateContext(ctx)
	if err != nil {
		return fmtp.Errorf("error waiting for DLI flink job (%s) to be created: %s", id, err)
	}
	return nil
}

func checkFlinkJobStopResult(ctx context.Context, client *golangsdk.ServiceClient, id int,
	timeout time.Duration) error {
	stateConf := &resource.StateChangeConf{
		Pending: []string{"job_submitting", "job_running", "job_canceling", "job_savepointing",
			"job_arrearage_recovering"},
		Target: []string{"job_init", "job_cancel_success"},
		Refresh: func() (interface{}, string, error) {
			job, err := flinkjob.Get(client, id)
			logp.Printf("[DEBUG] the flink job info in stop check func: %#v,%s", job, err)
			if err != nil {
				return nil, "", err
			}
			return job, job.JobDetail.Status, nil
		},
		Timeout:      timeout,
		PollInterval: 20 * timeout,
		Delay:        20 * time.Second,
	}
	_, err := stateConf.WaitForStateContext(ctx)
	if err != nil {
		return fmtp.Errorf("error waiting for DLI flink job (%s) to be stoped: %s", id, err)
	}
	return nil
}

func authorizeObsBucket(client *golangsdk.ServiceClient, d *schema.ResourceData) diag.Diagnostics {
	if value, ok := d.GetOk("obs_bucket"); ok {
		opts := flinkjob.ObsBucketsOpts{
			Buckets: []string{value.(string)},
		}

		logp.Printf("[DEBUG] update DLI flink job opts: %#v", opts)

		rst, uErr := flinkjob.AuthorizeBucket(client, opts)
		if uErr != nil {
			return fmtp.DiagErrorf("DLI Authorization on the following OBS buckets failed= %s: %s", value, uErr)
		}

		if rst != nil && !rst.IsSuccess {
			return fmtp.DiagErrorf("DLI Authorization on the following OBS buckets failed= %s: %s", value, uErr)
		}
	}

	return nil
}

// After the savepoint is triggered, the job status information will be saved in the OBS bucket under
// jobs/savepoint/{job_id}/{yyyy-mm-dd_hH-mm-ss}/.
func updateFlinkSqlJobWithStop(ctx context.Context, client *golangsdk.ServiceClient, jobId int,
	d *schema.ResourceData) diag.Diagnostics {

	if d.HasChangesExcept("smn_topic", "restart_when_exception", "resume_checkpoint", "resume_max_num", "obs_bucket") {
		// 1. stop the job
		_, err := flinkjob.Stop(client, flinkjob.StopFlinkJobInBatch{
			TriggerSavepoint: utils.Bool(false),
			JobIds:           []int{jobId},
		})

		if err != nil {
			return fmtp.DiagErrorf("stop job exception during update DLI flink job=%d: %s", jobId, err)
		}

		checkStopErr := checkFlinkJobStopResult(ctx, client, jobId, d.Timeout(schema.TimeoutUpdate))
		if checkStopErr != nil {
			return diag.FromErr(checkStopErr)
		}

		//2. update job
		opts := flinkjob.UpdateSqlJobOpts{
			Name:               d.Get("name").(string),
			RunMode:            d.Get("run_mode").(string),
			Desc:               d.Get("description").(string),
			QueueName:          d.Get("queue_name").(string),
			SqlBody:            d.Get("sql").(string),
			CuNumber:           golangsdk.IntToPointer(d.Get("cu_number").(int)),
			ParallelNumber:     golangsdk.IntToPointer(d.Get("parallel_number").(int)),
			CheckpointEnabled:  utils.Bool(d.Get("checkpoint_enabled").(bool)),
			CheckpointInterval: golangsdk.IntToPointer(d.Get("checkpoint_interval").(int)),
			LogEnabled:         utils.Bool(d.Get("log_enabled").(bool)),
			IdleStateRetention: golangsdk.IntToPointer(d.Get("idle_state_retention").(int)),
			DirtyDataStrategy:  d.Get("dirty_data_strategy").(string),
			UdfJarUrl:          d.Get("udf_jar_url").(string),
			ManagerCuNumber:    golangsdk.IntToPointer(d.Get("manager_cu_number").(int)),
			TmCus:              golangsdk.IntToPointer(d.Get("tm_cus").(int)),
			TmSlotNum:          golangsdk.IntToPointer(d.Get("tm_slot_num").(int)),
		}

		if runtimConfig, ok := d.GetOk("runtime_config"); ok {
			config := utils.ExpandResourceTags(runtimConfig.(map[string]interface{}))
			configStr, _ := json.Marshal(config)
			opts.RuntimeConfig = string(configStr)
		}

		if mode := d.Get("checkpoint_mode").(string); mode == flinkjob.CheckpointModeAtLeastOnce {
			opts.CheckpointMode = golangsdk.IntToPointer(2)
		} else {
			opts.CheckpointMode = golangsdk.IntToPointer(1)
		}

		if edgeGroupIds, ok := d.GetOk("edge_group_ids"); ok {
			var ids []string
			for _, v := range edgeGroupIds.([]interface{}) {
				ids = append(ids, v.(string))
			}
			opts.EdgeGroupIds = ids
		}

		logp.Printf("[DEBUG] update DLI flink job opts: %#v", opts)

		rst, uErr := flinkjob.UpdateSqlJob(client, jobId, opts)
		if uErr != nil {
			return fmtp.DiagErrorf("Error update DLI flink job=%d: %s", jobId, uErr)
		}

		if rst != nil && !rst.IsSuccess {
			return fmtp.DiagErrorf("Error update DLI flink job=%d: %s", rst.Message)
		}

		//3. run the flink job
		_, runErr := flinkjob.Run(client, flinkjob.RunJobOpts{
			JobIds:          []int{jobId},
			ResumeSavepoint: utils.Bool(d.Get("resume_checkpoint").(bool)),
		})
		if runErr != nil {
			return fmtp.DiagErrorf("Error run DLI flink job: %s", runErr)
		}
	}

	return nil
}

func parseConfig(configStr string) []tags.ResourceTag {
	var rst []tags.ResourceTag
	json.Unmarshal([]byte(configStr), &rst)
	return rst
}
