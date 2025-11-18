package dli

import (
	"context"
	"encoding/json"
	"log"
	"strconv"
	"time"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/chnsz/golangsdk"
	"github.com/chnsz/golangsdk/openstack/dli/v1/flinkjob"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API DLI POST /v1.0/{project_id}/dli/obs-authorize
// @API DLI POST /v1.0/{project_id}/streaming/flink-jobs
// @API DLI POST /v1.0/{project_id}/streaming/jobs/run
// @API DLI GET /v1.0/{project_id}/streaming/jobs/{job_id}
// @API DLI PUT /v1.0/{project_id}/streaming/flink-jobs/{job_id}
// @API DLI POST /v1.0/{project_id}/streaming/jobs/stop
// @API DLI DELETE /v1.0/{project_id}/streaming/jobs/{job_id}
// @API DLI GET /v3/{project_id}/dli_flink_job/{resource_id}/tags
// @API DLI POST /v3/{project_id}/dli_flink_job/{resource_id}/tags/create
// @API DLI POST /v3/{project_id}/dli_flink_job/{resource_id}/tags/delete
func ResourceFlinkJarJob() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceFlinkJarJobCreate,
		ReadContext:   resourceFlinkJarJobRead,
		UpdateContext: resourceFlinkJarJobUpdate,
		DeleteContext: resourceFlinkJarJobDelete,
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

			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},

			"queue_name": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"main_class": {
				Type:     schema.TypeString,
				Optional: true,
			},

			"entrypoint": {
				Type:     schema.TypeString,
				Optional: true,
			},

			"entrypoint_args": {
				Type:     schema.TypeString,
				Optional: true,
			},

			"dependency_jars": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},

			"dependency_files": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},

			"feature": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"basic", "custom"}, true),
				Computed:     true,
			},

			"flink_version": {
				Type:         schema.TypeString,
				Optional:     true,
				RequiredWith: []string{"feature"},
				Computed:     true,
			},

			"image": {
				Type:     schema.TypeString,
				Optional: true,
			},

			"cu_num": {
				Type:     schema.TypeInt,
				Optional: true,
				Default:  2,
			},

			"parallel_num": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},

			"obs_bucket": {
				Type:         schema.TypeString,
				Optional:     true,
				RequiredWith: []string{"log_enabled"},
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

			"manager_cu_num": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},

			"tm_cu_num": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
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

			"checkpoint_path": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"tags": common.TagsSchema(),

			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(20 * time.Minute),
			Update: schema.DefaultTimeout(30 * time.Minute),
			Delete: schema.DefaultTimeout(40 * time.Minute),
		},
	}
}

func resourceFlinkJarJobCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
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

	opts := flinkjob.CreateJarJobOpts{
		Name:                 d.Get("name").(string),
		Desc:                 d.Get("description").(string),
		QueueName:            d.Get("queue_name").(string),
		MainClass:            d.Get("main_class").(string),
		Entrypoint:           d.Get("entrypoint").(string),
		EntrypointArgs:       d.Get("entrypoint_args").(string),
		Feature:              d.Get("feature").(string),
		FlinkVersion:         d.Get("flink_version").(string),
		Image:                d.Get("image").(string),
		CuNumber:             golangsdk.IntToPointer(d.Get("cu_num").(int)),
		ParallelNumber:       d.Get("parallel_num").(int),
		ObsBucket:            d.Get("obs_bucket").(string),
		LogEnabled:           utils.Bool(d.Get("log_enabled").(bool)),
		SmnTopic:             d.Get("smn_topic").(string),
		RestartWhenException: utils.Bool(d.Get("restart_when_exception").(bool)),
		ManagerCuNumber:      d.Get("manager_cu_num").(int),
		TmCus:                d.Get("tm_cu_num").(int),
		TmSlotNum:            golangsdk.IntToPointer(d.Get("tm_slot_num").(int)),
		ResumeCheckpoint:     utils.Bool(d.Get("resume_checkpoint").(bool)),
		ResumeMaxNum:         golangsdk.IntToPointer(d.Get("resume_max_num").(int)),
		CheckpointPath:       d.Get("checkpoint_path").(string),
	}

	if runtimConfig, ok := d.GetOk("runtime_config"); ok {
		config := utils.ExpandResourceTags(runtimConfig.(map[string]interface{}))
		configStr, _ := json.Marshal(config)
		opts.RuntimeConfig = string(configStr)
	}

	if dependencyJars, ok := d.GetOk("dependency_jars"); ok {
		var dependencyArray []string
		for _, v := range dependencyJars.([]interface{}) {
			dependencyArray = append(dependencyArray, v.(string))
		}
		opts.DependencyJars = dependencyArray
	}

	if dependencyFiles, ok := d.GetOk("dependency_files"); ok {
		var dependencyArray []string
		for _, v := range dependencyFiles.([]interface{}) {
			dependencyArray = append(dependencyArray, v.(string))
		}
		opts.DependencyFiles = dependencyArray
	}

	log.Printf("[DEBUG] Creating new DLI flink jar job opts: %#v", opts)

	rst, createErr := flinkjob.CreateJarJob(client, opts)
	if createErr != nil {
		return diag.Errorf("error creating DLI flink jar job: %s", createErr)
	}

	if rst != nil && !rst.IsSuccess {
		return diag.Errorf("error creating DLI flink jar job: %s", rst.Message)
	}

	d.SetId(strconv.Itoa(rst.Job.JobId))

	if err = addTagsToResource(config, region, d); err != nil {
		return diag.FromErr(err)
	}
	// run the flink jar job
	_, runErr := flinkjob.Run(client, flinkjob.RunJobOpts{
		JobIds:          []int{rst.Job.JobId},
		ResumeSavepoint: utils.Bool(false),
	})
	if runErr != nil {
		return diag.Errorf("error run DLI flink jar job: %s", runErr)
	}

	checkCreateErr := checkFlinkJobRunResult(ctx, client, rst.Job.JobId, d.Timeout(schema.TimeoutCreate))
	if checkCreateErr != nil {
		return diag.FromErr(checkCreateErr)
	}
	return resourceFlinkJarJobRead(ctx, d, meta)
}

func resourceFlinkJarJobRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
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
		return diag.Errorf("error query DLI flink jar job: %s", detailRsp.Message)
	}
	detail := detailRsp.JobDetail
	mErr := multierror.Append(
		d.Set("name", detail.Name),
		d.Set("description", detail.Desc),
		d.Set("queue_name", detail.QueueName),
		d.Set("main_class", detail.MainClass),
		d.Set("entrypoint", detail.JobConfig.Entrypoint),
		d.Set("entrypoint_args", detail.EntrypointArgs),
		d.Set("dependency_jars", detail.JobConfig.DependencyJars),
		d.Set("dependency_files", detail.JobConfig.DependencyFiles),
		d.Set("feature", detail.JobConfig.Feature),
		d.Set("flink_version", detail.JobConfig.FlinkVersion),
		d.Set("image", detail.JobConfig.Image),
		d.Set("cu_num", detail.JobConfig.CuNumber),
		d.Set("parallel_num", detail.JobConfig.ParallelNumber),
		d.Set("obs_bucket", detail.JobConfig.ObsBucket),
		d.Set("log_enabled", detail.JobConfig.LogEnabled),
		d.Set("smn_topic", detail.JobConfig.SmnTopic),
		d.Set("restart_when_exception", detail.JobConfig.RestartWhenException),
		d.Set("manager_cu_num", detail.JobConfig.ManagerCuNumber),
		d.Set("tm_cu_num", detail.JobConfig.TmCus),
		d.Set("tm_slot_num", detail.JobConfig.TmSlotNum),
		d.Set("resume_checkpoint", detail.JobConfig.ResumeCheckpoint),
		d.Set("resume_max_num", detail.JobConfig.ResumeMaxNum),
		d.Set("checkpoint_path", detail.JobConfig.CheckpointPath),
		setRuntimeConfigToState(d, detail.JobConfig.RuntimeConfig),
		d.Set("status", detail.Status),
		d.Set("tags", d.Get("tags")),
	)

	if err = setTagsToResource(config, region, d); err != nil {
		return diag.FromErr(err)
	}

	return diag.FromErr(mErr.ErrorOrNil())
}

func resourceFlinkJarJobDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*config.Config)
	region := config.GetRegion(d)
	client, err := config.DliV1Client(region)
	if err != nil {
		return diag.Errorf("error creating DLI v1 client, err=%s", err)
	}

	jobId, aErr := strconv.Atoi(d.Id())
	if aErr != nil {
		return diag.Errorf("the DLI flink job_id must be number, but the actual ID is '%s'", d.Id())
	}

	deleteRst, dErr := flinkjob.Delete(client, jobId)
	if dErr != nil {
		return diag.Errorf("delete DLI flink jar job failed. %q:%s", jobId, dErr)
	}
	if deleteRst != nil && !deleteRst.IsSuccess {
		return diag.Errorf("delete DLI flink jar job (%d) failed: %s", jobId, deleteRst.Message)
	}

	return nil
}

func resourceFlinkJarJobUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*config.Config)
	region := config.GetRegion(d)
	if err := updateTagsToResource(config, region, d); err != nil {
		return diag.FromErr(err)
	}

	if d.HasChangesExcept("tags") {
		client, err := config.DliV1Client(region)
		if err != nil {
			return diag.Errorf("error creating DLI v1 client, err=%s", err)
		}

		jobId, iErr := strconv.Atoi(d.Id())
		if iErr != nil {
			return diag.Errorf("the DLI flink job_id must be number, but the actual ID is '%s'", d.Id())
		}

		diagErr := authorizeObsBucket(client, d)
		if diagErr != nil {
			return diagErr
		}

		diagErr = updateFlinkJarJobInRunning(client, jobId, d)
		if diagErr != nil {
			return diagErr
		}

		diagErr = updateFlinkJarJobWithStop(ctx, client, jobId, d)
		if diagErr != nil {
			return diagErr
		}

		err = checkFlinkJobRunResult(ctx, client, jobId, d.Timeout(schema.TimeoutCreate))
		if err != nil {
			return diag.FromErr(err)
		}
	}
	return resourceFlinkJarJobRead(ctx, d, meta)
}

func updateFlinkJarJobInRunning(client *golangsdk.ServiceClient, jobId int, d *schema.ResourceData) diag.Diagnostics {
	if d.HasChanges("smn_topic", "restart_when_exception", "resume_checkpoint", "resume_max_num", "obs_bucket",
		"checkpoint_path") {

		opts := flinkjob.UpdateJarJobOpts{
			ObsBucket:            d.Get("obs_bucket").(string),
			SmnTopic:             d.Get("smn_topic").(string),
			RestartWhenException: utils.Bool(d.Get("restart_when_exception").(bool)),
			ResumeCheckpoint:     utils.Bool(d.Get("resume_checkpoint").(bool)),
			ResumeMaxNum:         golangsdk.IntToPointer(d.Get("resume_max_num").(int)),
			CheckpointPath:       d.Get("checkpoint_path").(string),
		}

		log.Printf("[DEBUG] update DLI flink jar job opts: %#v", opts)

		rst, uErr := flinkjob.UpdateJarJob(client, jobId, opts)
		if uErr != nil {
			return diag.Errorf("error update DLI flink jar job=%d: %s", jobId, uErr)
		}

		if rst != nil && !rst.IsSuccess {
			return diag.Errorf("error update DLI flink jar job=%d: %s", jobId, rst.Message)
		}
	}

	return nil
}

func updateFlinkJarJobWithStop(ctx context.Context, client *golangsdk.ServiceClient, jobId int,
	d *schema.ResourceData) diag.Diagnostics {
	if d.HasChangesExcept("smn_topic", "restart_when_exception", "resume_checkpoint", "resume_max_num", "obs_bucket",
		"checkpoint_path") {

		// 1. stop the job
		_, err := flinkjob.Stop(client, flinkjob.StopFlinkJobInBatch{
			TriggerSavepoint: utils.Bool(d.Get("resume_checkpoint").(bool)),
			JobIds:           []int{jobId},
		})

		if err != nil {
			return diag.Errorf("stop job exception during update DLI flink jar job=%d: %s", jobId, err)
		}

		checkStopErr := checkFlinkJobStopResult(ctx, client, jobId, d.Timeout(schema.TimeoutUpdate))
		if checkStopErr != nil {
			return diag.FromErr(checkStopErr)
		}

		// 2. update job
		opts := flinkjob.UpdateJarJobOpts{
			Name:            d.Get("name").(string),
			Desc:            d.Get("description").(string),
			QueueName:       d.Get("queue_name").(string),
			MainClass:       d.Get("main_class").(string),
			Entrypoint:      d.Get("entrypoint").(string),
			EntrypointArgs:  d.Get("entrypoint_args").(string),
			Feature:         d.Get("feature").(string),
			FlinkVersion:    d.Get("flink_version").(string),
			Image:           d.Get("image").(string),
			CuNumber:        golangsdk.IntToPointer(d.Get("cu_num").(int)),
			ParallelNumber:  d.Get("parallel_num").(int),
			LogEnabled:      utils.Bool(d.Get("log_enabled").(bool)),
			ManagerCuNumber: d.Get("manager_cu_num").(int),
			TmCus:           d.Get("tm_cu_num").(int),
			TmSlotNum:       golangsdk.IntToPointer(d.Get("tm_slot_num").(int)),
		}

		if dependencyJars, ok := d.GetOk("dependency_jars"); ok {
			var dependencyArray []string
			for _, v := range dependencyJars.([]interface{}) {
				dependencyArray = append(dependencyArray, v.(string))
			}
			opts.DependencyJars = dependencyArray
		}

		if dependencyFiles, ok := d.GetOk("dependency_files"); ok {
			var dependencyArray []string
			for _, v := range dependencyFiles.([]interface{}) {
				dependencyArray = append(dependencyArray, v.(string))
			}
			opts.DependencyFiles = dependencyArray
		}

		if runtimConfig, ok := d.GetOk("runtime_config"); ok {
			config := utils.ExpandResourceTags(runtimConfig.(map[string]interface{}))
			configStr, _ := json.Marshal(config)
			opts.RuntimeConfig = string(configStr)
		}

		log.Printf("[DEBUG] update DLI flink jar job opts: %#v", opts)

		rst, uErr := flinkjob.UpdateJarJob(client, jobId, opts)
		if uErr != nil {
			return diag.Errorf("error update DLI flink jar job=%d: %s", jobId, uErr)
		}

		if rst != nil && !rst.IsSuccess {
			return diag.Errorf("error update DLI flink jar job=%d: %s", jobId, rst.Message)
		}

		// 3. run the flink jar job
		_, runErr := flinkjob.Run(client, flinkjob.RunJobOpts{
			JobIds:          []int{jobId},
			ResumeSavepoint: utils.Bool(d.Get("resume_checkpoint").(bool)),
		})
		if runErr != nil {
			return diag.Errorf("error run DLI flink jar job: %s", runErr)
		}
	}

	return nil
}
