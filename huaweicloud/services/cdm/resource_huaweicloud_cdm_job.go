package cdm

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/chnsz/golangsdk"
	"github.com/chnsz/golangsdk/openstack/cdm/v1/job"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
)

const (
	fromJobConfig = "fromJobConfig"
	toJobConfig   = "toJobConfig"
)

// @API CDM PUT /v1.1/{project_id}/clusters/{clusterId}/cdm/job/{jobName}/start
// @API CDM PUT /v1.1/{project_id}/clusters/{clusterId}/cdm/job/{jobName}/stop
// @API CDM GET /v1.1/{project_id}/clusters/{clusterId}/cdm/job/{jobName}
// @API CDM PUT /v1.1/{project_id}/clusters/{clusterId}/cdm/job/{jobName}
// @API CDM DELETE /v1.1/{project_id}/clusters/{clusterId}/cdm/job/{jobName}
// @API CDM POST /v1.1/{project_id}/clusters/{clusterId}/cdm/job
func ResourceCdmJob() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceCdmJobCreate,
		ReadContext:   resourceCdmJobRead,
		UpdateContext: resourceCdmJobUpdate,
		DeleteContext: resourceCdmJobDelete,
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
				Type:     schema.TypeString,
				Required: true,
			},

			"cluster_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"job_type": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"NORMAL_JOB", "BATCH_JOB", "SCENARIO_JOB"}, false),
			},

			"source_connector": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"source_link_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"source_job_config": {
				Type:     schema.TypeMap,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Required: true,
			},

			"destination_connector": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"destination_link_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"destination_job_config": {
				Type:     schema.TypeMap,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Required: true,
			},

			"config": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"throttling_extractors_number": {
							Type:     schema.TypeInt,
							Optional: true,
							Default:  1,
						},

						"group_name": {
							Type:     schema.TypeString,
							Optional: true,
							Default:  "DEFAULT",
						},

						"throttling_loader_number": {
							Type:     schema.TypeInt,
							Optional: true,
						},

						"throttling_record_dirty_data": {
							Type:     schema.TypeBool,
							Optional: true,
						},

						"throttling_dirty_write_to_link": {
							Type:     schema.TypeString,
							Optional: true,
						},

						"throttling_dirty_write_to_bucket": {
							Type:     schema.TypeString,
							Optional: true,
						},

						"throttling_dirty_write_to_directory": {
							Type:     schema.TypeString,
							Optional: true,
						},

						"throttling_max_error_records": {
							Type:     schema.TypeInt,
							Optional: true,
						},

						"scheduler_enabled": {
							Type:     schema.TypeBool,
							Optional: true,
							Default:  false,
						},

						"scheduler_cycle_type": {
							Type:     schema.TypeString,
							Optional: true,
							ValidateFunc: validation.StringInSlice([]string{"minute", "hour", "day", "week", "month"},
								false),
						},

						"scheduler_cycle": {
							Type:     schema.TypeInt,
							Optional: true,
						},

						"scheduler_run_at": {
							Type:     schema.TypeString,
							Optional: true,
						},

						"scheduler_start_date": {
							Type:     schema.TypeString,
							Optional: true,
						},

						"scheduler_stop_date": {
							Type:     schema.TypeString,
							Optional: true,
						},

						"scheduler_disposable_type": {
							Type:     schema.TypeString,
							Optional: true,
							ValidateFunc: validation.StringInSlice([]string{"NONE", "DELETE_AFTER_SUCCEED", "DELETE"},
								false),
						},

						"retry_type": {
							Type:         schema.TypeString,
							Optional:     true,
							ValidateFunc: validation.StringInSlice([]string{"NONE", "RETRY_TRIPLE"}, false),
							Default:      "NONE",
						},
					},
				},
			},

			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(20 * time.Minute),
			Delete: schema.DefaultTimeout(20 * time.Minute),
			Update: schema.DefaultTimeout(20 * time.Minute),
		},
	}
}

func resourceCdmJobCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	client, err := cfg.CdmV11Client(region)
	if err != nil {
		return diag.Errorf("error creating CDM v1.1 client, err=%s", err)
	}

	fromConfig, err := buildConfigParamter(d, "source_job_config", fromJobConfig)
	if err != nil {
		return diag.FromErr(err)
	}

	toConfig, err := buildConfigParamter(d, "destination_job_config", toJobConfig)
	if err != nil {
		return diag.FromErr(err)
	}

	opts := job.JobCreateOpts{
		Jobs: []job.Job{
			{
				Name:               d.Get("name").(string),
				JobType:            d.Get("job_type").(string),
				FromLinkName:       d.Get("source_link_name").(string),
				FromConnectorName:  d.Get("source_connector").(string),
				FromConfigValues:   *fromConfig,
				ToLinkName:         d.Get("destination_link_name").(string),
				ToConnectorName:    d.Get("destination_connector").(string),
				ToConfigValues:     *toConfig,
				DriverConfigValues: buildDriverConfigParamter(d),
			},
		},
	}

	log.Printf("[DEBUG] Creating CDM job opts: %#v", opts)
	clusterId := d.Get("cluster_id").(string)
	rst, createErr := job.Create(client, clusterId, opts)
	if createErr != nil {
		return diag.Errorf("error creating CDM job: %s", createErr)
	}

	d.SetId(fmt.Sprintf("%s/%s", clusterId, rst.Name))

	checkErr := waitingforJobRunning(ctx, client, clusterId, rst.Name, d.Timeout(schema.TimeoutCreate))
	if checkErr != nil {
		return diag.FromErr(checkErr)
	}
	return resourceCdmJobRead(ctx, d, meta)
}

func resourceCdmJobRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	client, err := cfg.CdmV11Client(region)
	if err != nil {
		return diag.Errorf("error creating CDM v1.1 client, err=%s", err)
	}

	clusterId, jobName, err := ParseJobInfoFromId(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	rst, gErr := job.Get(client, clusterId, jobName, job.GetJobsOpts{})
	log.Printf("[DEBUG] read CDM job opts: %v", gErr)

	if gErr != nil {
		return common.CheckDeletedDiag(d, parseCdmJobErrorToError404(gErr), "Error retrieving CDM job")
	}

	if len(rst.Jobs) < 1 {
		d.SetId("")
		return nil
	}

	detail := rst.Jobs[0]
	mErr := multierror.Append(
		d.Set("region", region),
		d.Set("name", detail.Name),
		d.Set("cluster_id", clusterId),
		d.Set("job_type", detail.JobType),
		d.Set("source_connector", detail.FromConnectorName),
		d.Set("source_link_name", detail.FromLinkName),
		d.Set("source_job_config", flattenFromOrToConfig(fromJobConfig, detail.FromConfigValues.Configs)),
		d.Set("destination_connector", detail.ToConnectorName),
		d.Set("destination_link_name", detail.ToLinkName),
		d.Set("destination_job_config", flattenFromOrToConfig(toJobConfig, detail.ToConfigValues.Configs)),
		setJobConfigtoState(d, detail.DriverConfigValues.Configs),
		d.Set("status", detail.Status),
	)

	if mErr.ErrorOrNil() != nil {
		return diag.Errorf("error setting CDM job fields: %s", mErr)
	}

	return nil
}

func resourceCdmJobUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	client, err := cfg.CdmV11Client(region)
	if err != nil {
		return diag.Errorf("error creating CDM v1.1 client, err=%s", err)
	}

	clusterId, jobName, err := ParseJobInfoFromId(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	rst, gErr := job.Get(client, clusterId, jobName, job.GetJobsOpts{})
	if gErr != nil {
		return common.CheckDeletedDiag(d, parseCdmJobErrorToError404(gErr), "Error retrieving CDM job")
	}

	if d.HasChanges("name", "source_job_config", "destination_job_config", "config") {
		status := rst.Jobs[0].Status
		// shutdown job
		if status == "BOOTING" || status == "RUNNING" {
			sErr := job.Stop(client, clusterId, jobName)
			if sErr != nil {
				return diag.Errorf("stop job failed when update CDM job. %q:%s", d.Id(), sErr)
			}
		}

		fromConfig, err := buildConfigParamter(d, "source_job_config", fromJobConfig)
		if err != nil {
			return diag.FromErr(err)
		}

		toConfig, err := buildConfigParamter(d, "destination_job_config", toJobConfig)
		if err != nil {
			return diag.FromErr(err)
		}

		opts := job.JobCreateOpts{
			Jobs: []job.Job{
				{
					Name:               d.Get("name").(string),
					JobType:            d.Get("job_type").(string),
					FromLinkName:       d.Get("source_link_name").(string),
					FromConnectorName:  d.Get("source_connector").(string),
					FromConfigValues:   *fromConfig,
					ToLinkName:         d.Get("destination_link_name").(string),
					ToConnectorName:    d.Get("destination_connector").(string),
					ToConfigValues:     *toConfig,
					DriverConfigValues: buildDriverConfigParamter(d),
				},
			},
		}

		log.Printf("[DEBUG] update CDM job opts: %#v", opts)

		_, uErr := job.Update(client, clusterId, jobName, opts)
		if uErr != nil {
			return diag.Errorf("error update CDM job: %s", uErr)
		}

		checkErr := waitingforJobRunning(ctx, client, clusterId, jobName, d.Timeout(schema.TimeoutUpdate))
		if checkErr != nil {
			return diag.FromErr(checkErr)
		}
	}

	return nil
}

func resourceCdmJobDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	client, err := cfg.CdmV11Client(region)
	if err != nil {
		return diag.Errorf("error creating CDM v1.1 client, err=%s", err)
	}

	clusterId, jobName, err := ParseJobInfoFromId(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	rst, gErr := job.Get(client, clusterId, jobName, job.GetJobsOpts{})
	if gErr != nil {
		return common.CheckDeletedDiag(d, parseCdmJobErrorToError404(gErr), "Error retrieving CDM job")
	}

	status := rst.Jobs[0].Status
	// shutdown job
	if status == "BOOTING" || status == "RUNNING" {
		stopResp := job.Stop(client, clusterId, jobName)
		if stopResp.Err != nil {
			return diag.Errorf("stop job failed when delete CDM job %q. response= %s", d.Id(), stopResp)
		}
	}

	// delete job
	resp, dErr := job.Delete(client, clusterId, jobName)
	if dErr != nil {
		return diag.Errorf("delete CDM job %q failed. err= %s, response= %s", d.Id(), dErr, resp)
	}

	d.SetId("")

	return nil
}

func buildConfigParamter(d *schema.ResourceData, inputConfigName, outConfigName string) (*job.JobConfigs, error) {
	var confs []job.Input
	configRaw := d.Get(inputConfigName).(map[string]interface{})

	if len(configRaw) < 1 {
		return nil, fmt.Errorf("the %s is Required", inputConfigName)
	}

	for k, v := range configRaw {
		conf := job.Input{
			Name:  fmt.Sprintf("%s.%s", outConfigName, k),
			Value: v.(string),
		}
		confs = append(confs, conf)
	}

	rst := job.JobConfigs{
		Configs: []job.Configs{
			{
				Name:   outConfigName,
				Inputs: confs,
			},
		},
	}

	return &rst, nil
}

func buildDriverConfigParamter(d *schema.ResourceData) job.JobConfigs {
	throttlingConfigs := buildThrottlingConfigParamter(d)
	schedulerConfigs := buildSchedulerConfigsParamter(d)

	var retryConfigs []job.Input
	if v, ok := d.GetOk("config.0.retry_type"); ok {
		retryConfigs = append(retryConfigs, job.Input{
			Name:  "retryJobConfig.retryJobType",
			Value: fmt.Sprintf("%s", v),
		})
	}

	var groupConfigs []job.Input
	if v, ok := d.GetOk("config.0.group_name"); ok {
		groupConfigs = append(groupConfigs, job.Input{
			Name:  "groupJobConfig.groupName",
			Value: fmt.Sprintf("%s", v),
		})
	}

	var configs []job.Configs
	if len(throttlingConfigs) > 0 {
		configs = append(configs, job.Configs{
			Name:   "throttlingConfig",
			Inputs: throttlingConfigs,
		})
	}
	if len(schedulerConfigs) > 0 {
		configs = append(configs, job.Configs{
			Name:   "schedulerConfig",
			Inputs: schedulerConfigs,
		})
	}
	if len(retryConfigs) > 0 {
		configs = append(configs, job.Configs{
			Name:   "retryJobConfig",
			Inputs: retryConfigs,
		})
	}
	if len(groupConfigs) > 0 {
		configs = append(configs, job.Configs{
			Name:   "groupJobConfig",
			Inputs: groupConfigs,
		})
	}

	log.Printf("[DEBUG] create CDM job opts: %#v", configs)
	return job.JobConfigs{Configs: configs}
}

func buildThrottlingConfigParamter(d *schema.ResourceData) []job.Input {
	var throttlingConfigs []job.Input

	if v, ok := d.GetOk("config.0.throttling_extractors_number"); ok {
		throttlingConfigs = append(throttlingConfigs, job.Input{
			Name:  "throttlingConfig.numExtractors",
			Value: fmt.Sprintf("%d", v),
		})
	}

	if v, ok := d.GetOk("config.0.throttling_loader_number"); ok {
		throttlingConfigs = append(throttlingConfigs, job.Input{
			Name:  "throttlingConfig.numLoaders",
			Value: fmt.Sprintf("%d", v),
		})
	}

	if v, ok := d.GetOk("config.0.throttling_record_dirty_data"); ok {
		throttlingConfigs = append(throttlingConfigs, job.Input{
			Name:  "throttlingConfig.recordDirtyData",
			Value: fmt.Sprintf("%t", v),
		})
	}

	if v, ok := d.GetOk("config.0.throttling_dirty_write_to_link"); ok {
		throttlingConfigs = append(throttlingConfigs, job.Input{
			Name:  "throttlingConfig.writeToLink",
			Value: fmt.Sprintf("%s", v),
		})
	}

	if v, ok := d.GetOk("config.0.throttling_dirty_write_to_bucket"); ok {
		throttlingConfigs = append(throttlingConfigs, job.Input{
			Name:  "throttlingConfig.obsBucket",
			Value: fmt.Sprintf("%s", v),
		})
	}

	if v, ok := d.GetOk("config.0.throttling_dirty_write_to_directory"); ok {
		throttlingConfigs = append(throttlingConfigs, job.Input{
			Name:  "throttlingConfig.dirtyDataDirectory",
			Value: fmt.Sprintf("%s", v),
		})
	}

	if v, ok := d.GetOk("config.0.throttling_max_error_records"); ok {
		throttlingConfigs = append(throttlingConfigs, job.Input{
			Name:  "throttlingConfig.maxErrorRecords",
			Value: fmt.Sprintf("%d", v),
		})
	}
	return throttlingConfigs
}

func buildSchedulerConfigsParamter(d *schema.ResourceData) []job.Input {
	var schedulerConfigs []job.Input

	if v, ok := d.GetOk("config.0.scheduler_enabled"); ok {
		schedulerConfigs = append(schedulerConfigs, job.Input{
			Name:  "schedulerConfig.isSchedulerJob",
			Value: fmt.Sprintf("%t", v),
		})

		if v, ok := d.GetOk("config.0.scheduler_cycle_type"); ok {
			schedulerConfigs = append(schedulerConfigs, job.Input{
				Name:  "schedulerConfig.cycleType",
				Value: fmt.Sprintf("%s", v),
			})
		}

		if v, ok := d.GetOk("config.0.scheduler_cycle"); ok {
			schedulerConfigs = append(schedulerConfigs, job.Input{
				Name:  "schedulerConfig.cycle",
				Value: fmt.Sprintf("%d", v),
			})
		}

		if v, ok := d.GetOk("config.0.scheduler_run_at"); ok {
			schedulerConfigs = append(schedulerConfigs, job.Input{
				Name:  "schedulerConfig.runAt",
				Value: fmt.Sprintf("%s", v),
			})
		}

		if v, ok := d.GetOk("config.0.scheduler_start_date"); ok {
			schedulerConfigs = append(schedulerConfigs, job.Input{
				Name:  "schedulerConfig.startDate",
				Value: fmt.Sprintf("%s", v),
			})
		}

		if v, ok := d.GetOk("config.0.scheduler_stop_date"); ok {
			schedulerConfigs = append(schedulerConfigs, job.Input{
				Name:  "schedulerConfig.stopDate",
				Value: fmt.Sprintf("%s", v),
			})
		}

		if v, ok := d.GetOk("config.0.scheduler_disposable_type"); ok {
			schedulerConfigs = append(schedulerConfigs, job.Input{
				Name:  "schedulerConfig.disposableType",
				Value: fmt.Sprintf("%s", v),
			})
		}
	}
	return schedulerConfigs
}

func flattenFromOrToConfig(configName string, configs []job.Configs) map[string]interface{} {
	configPref := fmt.Sprintf("%s.", configName)
	result := make(map[string]interface{})
	for _, item := range configs {
		if item.Name == configName {
			for _, v := range item.Inputs {
				if v.Value != "" {
					key := strings.Replace(v.Name, configPref, "", 1)
					// Value in return is encoded, use `url.PathUnescape` to decode it.
					result[key], _ = url.PathUnescape(v.Value)
				}
			}
		}
	}
	return result
}

// nolint:gocyclo
func setJobConfigtoState(d *schema.ResourceData, configs []job.Configs) error {
	var err *multierror.Error
	var pErr error
	result := make(map[string]interface{})
	for _, item := range configs {
		for _, v := range item.Inputs {
			if v.Value != "" {
				switch v.Name {
				case "throttlingConfig.numExtractors":
					result["throttling_extractors_number"], pErr = strconv.Atoi(v.Value)
					err = multierror.Append(err, pErr)
				case "throttlingConfig.numLoaders":
					result["throttling_loader_number"], pErr = strconv.Atoi(v.Value)
					err = multierror.Append(err, pErr)
				case "throttlingConfig.recordDirtyData":
					result["throttling_record_dirty_data"], pErr = strconv.ParseBool(v.Value)
					err = multierror.Append(err, pErr)
				case "throttlingConfig.writeToLink":
					result["throttling_dirty_write_to_link"] = v.Value
				case "throttlingConfig.obsBucket":
					result["throttling_dirty_write_to_bucket"] = v.Value
				case "throttlingConfig.dirtyDataDirectory":
					// Value in return is encoded, use `url.PathUnescape` to decode it.
					result["throttling_dirty_write_to_directory"], pErr = url.PathUnescape(v.Value)
					err = multierror.Append(err, pErr)
				case "throttlingConfig.maxErrorRecords":
					result["throttling_max_error_records"], pErr = strconv.Atoi(v.Value)
					err = multierror.Append(err, pErr)
				case "schedulerConfig.isSchedulerJob":
					result["scheduler_enabled"], pErr = strconv.ParseBool(v.Value)
					err = multierror.Append(err, pErr)
				case "schedulerConfig.cycleType":
					result["scheduler_cycle_type"] = v.Value
				case "schedulerConfig.cycle":
					result["scheduler_cycle"], pErr = strconv.Atoi(v.Value)
					err = multierror.Append(err, pErr)
				case "schedulerConfig.runAt":
					result["scheduler_run_at"] = v.Value
				case "schedulerConfig.startDate":
					// Value in return is encoded, use `url.PathUnescape` to decode it.
					result["scheduler_start_date"], pErr = url.PathUnescape(v.Value)
					err = multierror.Append(err, pErr)
				case "schedulerConfig.stopDate":
					// Value in return is encoded, use `url.PathUnescape` to decode it.
					result["scheduler_stop_date"], pErr = url.PathUnescape(v.Value)
					err = multierror.Append(err, pErr)
				case "schedulerConfig.disposableType":
					result["scheduler_disposable_type"] = v.Value
				case "retryJobConfig.retryJobType":
					result["retry_type"] = v.Value
				case "groupJobConfig.groupName":
					result["group_name"] = v.Value
				}
			}
		}
	}

	if err.ErrorOrNil() != nil {
		return err.ErrorOrNil()
	}

	return d.Set("config", []map[string]interface{}{result})
}

func waitingforJobRunning(ctx context.Context, client *golangsdk.ServiceClient, clusterId, jobName string,
	timeout time.Duration) error {
	// start job
	startResp, startErr := job.Start(client, clusterId, jobName)
	if startErr != nil {
		return fmt.Errorf("error start CDM job: %s", startErr)
	}
	if startResp.Submissions[0].Status == "FAILURE_ON_SUBMIT" || startResp.Submissions[0].Status == "FAILED" ||
		startResp.Submissions[0].Status == "NEVER_EXECUTED" {
		return fmt.Errorf("error start CDM job:%f,%s", startResp.Submissions[0].Progress,
			startResp.Submissions[0].Status)
	}

	// check job status
	stateConf := &resource.StateChangeConf{
		Pending: []string{job.StatusBooting},
		Target:  []string{job.StatusRunning, job.StatusSucceeded},
		Refresh: func() (interface{}, string, error) {
			rst, gErr := job.Get(client, clusterId, jobName, job.GetJobsOpts{})
			log.Printf("[DEBUG] query CDM job in running check func:%s", gErr)
			if gErr != nil {
				return nil, "", gErr
			}
			detail := rst.Jobs[0]
			if detail.Status == job.StatusFailed || detail.Status == job.StatusFailureOnSubmit {
				return detail, "failed", fmt.Errorf("%s", detail.Status)
			}
			return detail, detail.Status, nil
		},
		Timeout:      timeout,
		PollInterval: 20 * timeout,
		Delay:        20 * time.Second,
	}
	_, err := stateConf.WaitForStateContext(ctx)
	if err != nil {
		return fmt.Errorf("error waiting for CDM job (%s/%s) to be created: %s", clusterId, jobName, err)
	}
	return nil
}

func ParseJobInfoFromId(id string) (clusterId, jobName string, err error) {
	idArrays := strings.SplitN(id, "/", 2)
	if len(idArrays) != 2 {
		err = fmt.Errorf("invalid format specified for ID. Format must be <cluster_id>/<job_name>")
		return
	}

	clusterId = idArrays[0]
	jobName = idArrays[1]
	return
}

func parseCdmJobErrorToError404(respErr error) error {
	var apiError job.ErrorResponse

	if errCode, ok := respErr.(golangsdk.ErrDefault400); ok {
		pErr := json.Unmarshal(errCode.Body, &apiError)
		if pErr == nil && (apiError.ErrCode == "Cdm.0100" || apiError.ErrCode == "Cdm.0054") {
			return golangsdk.ErrDefault404(errCode)
		}
	}
	return respErr
}
