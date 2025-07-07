package dli

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/chnsz/golangsdk"
	"github.com/chnsz/golangsdk/openstack/common/tags"
	"github.com/chnsz/golangsdk/openstack/dli/v1/flinkjob"
	v3flinkjob "github.com/chnsz/golangsdk/openstack/dli/v3/flinkjob"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API DLI POST /v1.0/{project_id}/dli/obs-authorize
// @API DLI POST /v1.0/{project_id}/streaming/sql-jobs
// @API DLI POST /v1.0/{project_id}/streaming/jobs/run
// @API DLI GET /v1.0/{project_id}/streaming/jobs/{job_id}
// @API DLI PUT /v1.0/{project_id}/streaming/sql-jobs/{job_id}
// @API DLI POST /v1.0/{project_id}/streaming/jobs/stop
// @API DLI DELETE /v1.0/{project_id}/streaming/jobs/{job_id}
// @API DLI GET /v3/{project_id}/dli_flink_job/{resource_id}/tags
// @API DLI POST /v3/{project_id}/dli_flink_job/{resource_id}/tags/create
// @API DLI POST /v3/{project_id}/dli_flink_job/{resource_id}/tags/delete
// @API DLI POST /v3/{project_id}/streaming/jobs/{job_id}/gen-graph
func ResourceFlinkSqlJob() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceFlinkSqlJobCreate,
		ReadContext:   resourceFlinkSqlJobRead,
		UpdateContext: resourceFlinkSqlJobUpdate,
		DeleteContext: resourceFlinkSqlJobDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		CustomizeDiff: config.MergeDefaultTags(),

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},

			"name": {
				Type:     schema.TypeString,
				Required: true,
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
				Type:     schema.TypeString,
				Optional: true,
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

			"tags": common.TagsSchema(),
			"flink_version": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"graph_type": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"operator_config": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"static_estimator": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"static_estimator_config": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"stream_graph": {
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
		FlinkVersion:         d.Get("flink_version").(string),
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

	log.Printf("[DEBUG] Creating new DLI flink job opts: %#v", opts)

	rst, err := flinkjob.CreateSqlJob(client, opts)
	if err != nil {
		return diag.Errorf("error creating DLI flink job: %s", err)
	}

	if rst != nil && !rst.IsSuccess {
		return diag.Errorf("error creating DLI flink job: %s", rst.Message)
	}

	jobId := rst.Job.JobId
	d.SetId(strconv.Itoa(jobId))

	if err = addTagsToResource(config, region, d); err != nil {
		return diag.FromErr(err)
	}

	operatorConfig, hasOperatorConfig := d.GetOk("operator_config")
	staticEstimatorConfig, hasStaticEstimatorConfig := d.GetOk("static_estimator_config")
	if hasOperatorConfig || hasStaticEstimatorConfig {
		err = updateJobConfig(client, jobId, operatorConfig.(string), staticEstimatorConfig.(string))
		if err != nil {
			return diag.Errorf("error updating DLI flink job (%d): %s", jobId, err)
		}
	}
	// run the flink job
	_, err = flinkjob.Run(client, flinkjob.RunJobOpts{
		JobIds:          []int{jobId},
		ResumeSavepoint: utils.Bool(false),
	})
	if err != nil {
		return diag.Errorf("error run DLI flink job: %s", err)
	}

	err = checkFlinkJobRunResult(ctx, client, jobId, d.Timeout(schema.TimeoutCreate))
	if err != nil {
		return diag.FromErr(err)
	}
	return resourceFlinkSqlJobRead(ctx, d, meta)
}

func updateJobConfig(client *golangsdk.ServiceClient, jobId int, operatorConfig, staticEstimatorConfig string) error {
	opts := flinkjob.UpdateSqlJobOpts{
		OperatorConfig:        operatorConfig,
		StaticEstimatorConfig: staticEstimatorConfig,
	}
	resp, err := flinkjob.UpdateSqlJob(client, jobId, opts)
	if err != nil {
		return err
	}

	if !resp.IsSuccess {
		return fmt.Errorf(resp.Message)
	}
	return nil
}

func addTagsToResource(cfg *config.Config, region string, d *schema.ResourceData) error {
	if raw, ok := d.GetOk("tags"); ok {
		v3Client, err := cfg.DliV3Client(region)
		if err != nil {
			return fmt.Errorf("error creating DLI v3 client: %s", err)
		}

		id := d.Id()
		if err := addTags(v3Client, id, "dli_flink_job", raw.(map[string]interface{})); err != nil {
			return fmt.Errorf("error setting tags of the flink job (%s): %s", id, err)
		}
	}

	return nil
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
		return diag.Errorf("the DLI flink job_id must be number, but the actual ID is '%s'", d.Id())
	}

	detailRsp, err := flinkjob.Get(client, id)
	if err != nil {
		return common.CheckDeletedDiag(d, parseDliFlinkErrToError404(err), "DLI flink sql-job")
	}

	if detailRsp != nil && !detailRsp.IsSuccess {
		return diag.Errorf("error query DLI flink job: %s", detailRsp.Message)
	}
	detail := detailRsp.JobDetail
	mErr := multierror.Append(
		d.Set("name", detail.Name),
		d.Set("flink_version", detail.JobConfig.FlinkVersion),
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
		setRuntimeConfigToState(d, detail.JobConfig.RuntimeConfig),
		d.Set("operator_config", detail.JobConfig.OperatorConfig),
		d.Set("static_estimator_config", detail.JobConfig.StaticEstimatorConfig),
		d.Set("status", detail.Status),
	)

	if err = setTagsToResource(config, region, d); err != nil {
		return diag.FromErr(err)
	}

	_, ok := d.GetOk("graph_type")
	if d.Get("type").(string) == flinkjob.JobTypeFlinkOpenSourceSql && ok {
		v3Client, err := config.DliV3Client(region)
		if err != nil {
			return diag.Errorf("error creating DLI v3 client: %s", err)
		}

		jobId := d.Id()
		streamGraph, err := getSteramGraphById(v3Client, d, jobId)
		if err != nil {
			return diag.Errorf("error getting stream graph of flink sql (%s): %s", jobId, err)
		}
		mErr = multierror.Append(mErr, d.Set("stream_graph", streamGraph))
	}
	return diag.FromErr(mErr.ErrorOrNil())
}

func setTagsToResource(cfg *config.Config, region string, d *schema.ResourceData) error {
	v3Client, err := cfg.DliV3Client(region)
	if err != nil {
		return fmt.Errorf("error creating DLI v3 client: %s", err)
	}

	return utils.SetResourceTagsToState(d, v3Client, "dli_flink_job", d.Id())
}

func getSteramGraphById(client *golangsdk.ServiceClient, d *schema.ResourceData, jobId string) (string, error) {
	opts := v3flinkjob.StreamGraphOpts{
		JobId:                 jobId,
		SqlBody:               d.Get("sql").(string),
		FlinkVersion:          d.Get("flink_version").(string),
		CuNumber:              utils.Int(d.Get("cu_number").(int)),
		ManagerCuNumber:       utils.Int(d.Get("manager_cu_number").(int)),
		JobType:               flinkjob.JobTypeFlinkOpenSourceSql,
		GraphType:             d.Get("graph_type").(string),
		ParallelNumber:        utils.Int(d.Get("parallel_number").(int)),
		TmCus:                 utils.Int(d.Get("tm_cus").(int)),
		TmSlotNum:             utils.Int(d.Get("tm_slot_num").(int)),
		OperatorConfig:        d.Get("operator_config").(string),
		StaticEstimator:       utils.Bool(d.Get("static_estimator").(bool)),
		StaticEstimatorConfig: d.Get("static_estimator_config").(string),
	}
	resp, err := v3flinkjob.CreateFlinkSqlJobGraph(client, opts)
	if err != nil {
		return "", err
	}
	if !resp.IsSuccess {
		return "", fmt.Errorf(resp.Message)
	}
	return resp.StreamGraph, nil
}

// This API is used to cancel a submitted job. If execution of a job completes or fails, this job cannot be canceled.
func resourceFlinkSqlJobDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*config.Config)
	region := config.GetRegion(d)
	client, err := config.DliV1Client(region)
	if err != nil {
		return diag.Errorf("error creating DLI v1 client, err=%s", err)
	}

	jobId, err := strconv.Atoi(d.Id())
	if err != nil {
		return diag.Errorf("the DLI flink job_id must be number, but the actual ID is '%s'", d.Id())
	}

	deleteRst, err := flinkjob.Delete(client, jobId)
	if err != nil {
		return diag.Errorf("delete DLI flink job failed. %q:%s", jobId, err)
	}
	if deleteRst != nil && !deleteRst.IsSuccess {
		return diag.Errorf("delete DLI flink job (%d) failed: %s", jobId, deleteRst.Message)
	}

	return nil
}

func updateTagsToResource(cfg *config.Config, region string, d *schema.ResourceData) error {
	if d.HasChange("tags") {
		v3Client, err := cfg.DliV3Client(region)
		if err != nil {
			return fmt.Errorf("error creating DLI v3 client: %s", err)
		}

		id := d.Id()
		oldTags, newTags := d.GetChange("tags")
		err = updateResourceTags(v3Client, id, "dli_flink_job", oldTags, newTags)
		if err != nil {
			return fmt.Errorf("error updating tags of the flink job (%s): %s", id, err)
		}
	}
	return nil
}

func resourceFlinkSqlJobUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*config.Config)
	region := config.GetRegion(d)
	if err := updateTagsToResource(config, region, d); err != nil {
		return diag.FromErr(err)
	}

	if d.HasChangesExcept("tags", "graph_type") {
		client, err := config.DliV1Client(region)
		if err != nil {
			return diag.Errorf("error creating DLI v1 client, err=%s", err)
		}

		jobId, err := strconv.Atoi(d.Id())
		if err != nil {
			return diag.Errorf("the DLI flink job_id must be number, but the actual ID is '%s'", d.Id())
		}

		diagErr := authorizeObsBucket(client, d)
		if diagErr != nil {
			return diagErr
		}

		diagErr = updateFlinkSqlJobInRunning(client, jobId, d)
		if diagErr != nil {
			return diagErr
		}

		diagErr = updateFlinkSqlJobWithStop(ctx, client, jobId, d)
		if diagErr != nil {
			return diagErr
		}

		err = checkFlinkJobRunResult(ctx, client, jobId, d.Timeout(schema.TimeoutUpdate))
		if err != nil {
			return diag.FromErr(err)
		}
	}

	return resourceFlinkSqlJobRead(ctx, d, meta)
}

// updated in "job_running": smn_topic,restart_when_exception,resume_checkpoint,resume_max_num,checkpoint_path,obs_bucket
func updateFlinkSqlJobInRunning(client *golangsdk.ServiceClient, jobId int, d *schema.ResourceData) diag.Diagnostics {
	if d.HasChanges("smn_topic", "restart_when_exception", "resume_checkpoint", "resume_max_num", "obs_bucket") {
		opts := flinkjob.UpdateSqlJobOpts{
			ObsBucket:            d.Get("obs_bucket").(string),
			SmnTopic:             d.Get("smn_topic").(string),
			RestartWhenException: utils.Bool(d.Get("restart_when_exception").(bool)),
			ResumeCheckpoint:     utils.Bool(d.Get("resume_checkpoint").(bool)),
			ResumeMaxNum:         golangsdk.IntToPointer(d.Get("resume_max_num").(int)),
		}

		log.Printf("[DEBUG] update DLI flink job opts: %#v", opts)

		rst, uErr := flinkjob.UpdateSqlJob(client, jobId, opts)
		if uErr != nil {
			return diag.Errorf("error update DLI flink job=%d: %s", jobId, uErr)
		}

		if rst != nil && !rst.IsSuccess {
			return diag.Errorf("error update DLI flink job=%d: %s", jobId, rst.Message)
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
			log.Printf("[DEBUG] the flink job info in create check func: %#v,%s", job, err)
			if err != nil {
				return nil, "", err
			}
			if job.JobDetail.Status == "job_submit_fail" {
				return job, "failed", fmt.Errorf("%s:%s", job.JobDetail.Status, job.JobDetail.StatusDesc)
			}
			return job, job.JobDetail.Status, nil
		},
		Timeout:      timeout,
		PollInterval: 20 * timeout,
		Delay:        20 * time.Second,
	}
	_, err := stateConf.WaitForStateContext(ctx)
	if err != nil {
		return fmt.Errorf("error waiting for DLI flink job (%d) to be created: %s", id, err)
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
			log.Printf("[DEBUG] the flink job info in stop check func: %#v,%s", job, err)
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
		return fmt.Errorf("error waiting for DLI flink job (%d) to be stoped: %s", id, err)
	}
	return nil
}

func authorizeObsBucket(client *golangsdk.ServiceClient, d *schema.ResourceData) diag.Diagnostics {
	if value, ok := d.GetOk("obs_bucket"); ok {
		opts := flinkjob.ObsBucketsOpts{
			Buckets: []string{value.(string)},
		}

		log.Printf("[DEBUG] update DLI flink job opts: %#v", opts)

		rst, uErr := flinkjob.AuthorizeBucket(client, opts)
		if uErr != nil {
			return diag.Errorf("DLI Authorization on the following OBS buckets failed= %s: %s", value, uErr)
		}

		if rst != nil && !rst.IsSuccess {
			return diag.Errorf("DLI Authorization on the following OBS buckets failed: %s", rst.Message)
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
			return diag.Errorf("stop job exception during update DLI flink job=%d: %s", jobId, err)
		}

		checkStopErr := checkFlinkJobStopResult(ctx, client, jobId, d.Timeout(schema.TimeoutUpdate))
		if checkStopErr != nil {
			return diag.FromErr(checkStopErr)
		}

		// 2. update job
		opts := flinkjob.UpdateSqlJobOpts{
			Name:                  d.Get("name").(string),
			RunMode:               d.Get("run_mode").(string),
			Desc:                  d.Get("description").(string),
			QueueName:             d.Get("queue_name").(string),
			SqlBody:               d.Get("sql").(string),
			CuNumber:              golangsdk.IntToPointer(d.Get("cu_number").(int)),
			ParallelNumber:        golangsdk.IntToPointer(d.Get("parallel_number").(int)),
			CheckpointEnabled:     utils.Bool(d.Get("checkpoint_enabled").(bool)),
			CheckpointInterval:    golangsdk.IntToPointer(d.Get("checkpoint_interval").(int)),
			LogEnabled:            utils.Bool(d.Get("log_enabled").(bool)),
			IdleStateRetention:    golangsdk.IntToPointer(d.Get("idle_state_retention").(int)),
			DirtyDataStrategy:     d.Get("dirty_data_strategy").(string),
			UdfJarUrl:             d.Get("udf_jar_url").(string),
			ManagerCuNumber:       golangsdk.IntToPointer(d.Get("manager_cu_number").(int)),
			TmCus:                 golangsdk.IntToPointer(d.Get("tm_cus").(int)),
			TmSlotNum:             golangsdk.IntToPointer(d.Get("tm_slot_num").(int)),
			FlinkVersion:          d.Get("flink_version").(string),
			OperatorConfig:        d.Get("operator_config").(string),
			StaticEstimatorConfig: d.Get("static_estimator_config").(string),
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

		log.Printf("[DEBUG] update DLI flink job opts: %#v", opts)

		rst, uErr := flinkjob.UpdateSqlJob(client, jobId, opts)
		if uErr != nil {
			return diag.Errorf("error update DLI flink job=%d: %s", jobId, uErr)
		}

		if rst != nil && !rst.IsSuccess {
			return diag.Errorf("error update DLI flink job=%d: %s", jobId, rst.Message)
		}

		// 3. run the flink job
		_, runErr := flinkjob.Run(client, flinkjob.RunJobOpts{
			JobIds:          []int{jobId},
			ResumeSavepoint: utils.Bool(d.Get("resume_checkpoint").(bool)),
		})
		if runErr != nil {
			return diag.Errorf("error run DLI flink job: %s", runErr)
		}
	}

	return nil
}

func setRuntimeConfigToState(d *schema.ResourceData, configStr string) error {
	if len(configStr) == 0 {
		return nil
	}
	var rst []tags.ResourceTag
	err := json.Unmarshal([]byte(configStr), &rst)
	if err != nil {
		return fmt.Errorf("error parse runtime_config from API response: %s", err)
	}

	return d.Set("runtime_config", utils.TagsToMap(rst))
}

func parseDliFlinkErrToError404(respErr error) error {
	var apiError flinkjob.DliError

	if errCode, ok := respErr.(golangsdk.ErrDefault400); ok {
		pErr := json.Unmarshal(errCode.Body, &apiError)
		if pErr == nil && apiError.ErrorCode == "DLI.16001" {
			return golangsdk.ErrDefault404(errCode)
		}
	}
	return respErr
}
