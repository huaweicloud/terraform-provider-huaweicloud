package dataarts

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

var (
	actionResourceNotFoundCodes = []string{
		"DLF.0819", // The workspace ID does not exist.
	}
	jobActionNonUpdatableParams = []string{
		"job_name",
		"process_type",
		"workspace_id",
	}
)

// @API DataArtsStudio POST /v1/{project_id}/jobs/{job_name}/start
// @API DataArtsStudio POST /v1/{project_id}/jobs/{job_name}/stop
// @API DataArtsStudio POST /v1/{project_id}/jobs/{job_name}/run-immediate
// @API DataArtsStudio GET /v1/{project_id}/jobs
// @API DataArtsStudio GET /v1/{project_id}/jobs/{job_name}/instances/{instance_id}
func ResourceFactoryJobAction() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceFactoryJobActionCreate,
		ReadContext:   resourceFactoryJobActionRead,
		UpdateContext: resourceFactoryJobActionUpdate,
		DeleteContext: resourceFactoryJobActionDelete,

		CustomizeDiff: config.FlexibleForceNew(jobActionNonUpdatableParams),

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(20 * time.Minute),
			Update: schema.DefaultTimeout(20 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: `The region where the job is located.`,
			},

			// Required parameters.
			"job_name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The name of the job to be performed.`,
			},
			"process_type": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The type of the job to be performed.`,
			},
			"action": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The action type of the job to be performed.`,
			},

			// Optional parameters.
			"workspace_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The ID of the workspace to which the job belongs.`,
			},
			"job_params": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:        schema.TypeString,
							Required:    true,
							Description: `The name of the job parameter.`,
						},
						"value": {
							Type:        schema.TypeString,
							Required:    true,
							Description: `The value of the job parameter.`,
						},
						"type": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: `The type of the job parameter.`,
						},
					},
				},
				Description: `The parameters of the job action.`,
			},
			"start_date": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: `The start date of the job action when start job.`,
			},
			"ignore_first_self_dep": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: `Whether to ignore the first self dependence when start job.`,
			},
			"use_execution_user": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Whether to use the execution user to execute the job when start job immediately.`,
			},

			// Attribute.
			"status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The current status of the job.`,
			},
			"instance_status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The instance status after starting the job immediately.`,
			},

			// Internal parameters.
			"enable_force_new": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"true", "false"}, false),
				Description: utils.SchemaDesc(
					`Whether to allow parameters that do not support changes to have their change-triggered behavior set to 'ForceNew'.`,
					utils.SchemaDescInput{
						Internal: true,
					},
				),
			},
		},
	}
}

func buildFactoryRequestMoreHeaders(workspaceId string) map[string]string {
	results := map[string]string{
		"Content-Type": "application/json",
	}

	if workspaceId != "" {
		results["workspace"] = workspaceId
	}

	return results
}

func buildStartJobJobParams(jobParams []interface{}) []interface{} {
	if len(jobParams) < 1 {
		return nil
	}

	result := make([]interface{}, 0, len(jobParams))
	for _, jobParam := range jobParams {
		result = append(result, map[string]interface{}{
			// Required parameters.
			"name":  utils.PathSearch("name", jobParam, nil),
			"value": utils.PathSearch("value", jobParam, nil),
			// Optional parameters.
			"paramType": utils.ValueIgnoreEmpty(utils.PathSearch("type", jobParam, nil)),
		})
	}
	return result
}

func buildStartJobBodyParams(d *schema.ResourceData) map[string]interface{} {
	return map[string]interface{}{
		// Optional parameters.
		"jobParams":                    utils.ValueIgnoreEmpty(buildStartJobJobParams(d.Get("job_params").([]interface{}))),
		"start_date":                   utils.ValueIgnoreEmpty(d.Get("start_date").(int)),
		"ignore_first_self_dependence": utils.ValueIgnoreEmpty(d.Get("ignore_first_self_dep").(bool)),
	}
}

func startJob(client *golangsdk.ServiceClient, d *schema.ResourceData) error {
	var (
		httpUrl     = "v1/{project_id}/jobs/{job_name}/start"
		workspaceId = d.Get("workspace_id").(string)
		jobName     = d.Get("job_name").(string)
	)

	actionPath := client.Endpoint + httpUrl
	actionPath = strings.ReplaceAll(actionPath, "{project_id}", client.ProjectID)
	actionPath = strings.ReplaceAll(actionPath, "{job_name}", jobName)

	actionOpts := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      buildFactoryRequestMoreHeaders(workspaceId),
		JSONBody:         utils.RemoveNil(buildStartJobBodyParams(d)),
		OkCodes:          []int{204},
	}

	_, err := client.Request("POST", actionPath, &actionOpts)
	return err
}

func stopJob(client *golangsdk.ServiceClient, d *schema.ResourceData) error {
	var (
		httpUrl     = "v1/{project_id}/jobs/{job_name}/stop"
		workspaceId = d.Get("workspace_id").(string)
		jobName     = d.Get("job_name").(string)
	)

	actionPath := client.Endpoint + httpUrl
	actionPath = strings.ReplaceAll(actionPath, "{project_id}", client.ProjectID)
	actionPath = strings.ReplaceAll(actionPath, "{job_name}", jobName)

	actionOpts := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      buildFactoryRequestMoreHeaders(workspaceId),
		OkCodes:          []int{204},
	}

	_, err := client.Request("POST", actionPath, &actionOpts)
	return err
}

func buildStartJobImmediatelyJobParams(jobParams []interface{}) []interface{} {
	if len(jobParams) < 1 {
		return nil
	}

	result := make([]interface{}, 0, len(jobParams))
	for _, jobParam := range jobParams {
		result = append(result, map[string]interface{}{
			// Required parameters.
			"name":  utils.PathSearch("name", jobParam, nil),
			"value": utils.PathSearch("value", jobParam, nil),
			// Optional parameters.
			"type": utils.ValueIgnoreEmpty(utils.PathSearch("type", jobParam, nil)),
		})
	}
	return result
}

func buildStartJobImmediatelyBodyParams(d *schema.ResourceData) map[string]interface{} {
	return map[string]interface{}{
		// Optional parameters.
		"jobParams":        buildStartJobImmediatelyJobParams(d.Get("job_params").([]interface{})),
		"useExecutionUser": utils.ValueIgnoreEmpty(d.Get("use_execution_user").(string)),
	}
}

func startJobImmediately(client *golangsdk.ServiceClient, d *schema.ResourceData) (interface{}, error) {
	var (
		httpUrl     = "v1/{project_id}/jobs/{job_name}/run-immediate"
		workspaceId = d.Get("workspace_id").(string)
		jobName     = d.Get("job_name").(string)
	)

	actionPath := client.Endpoint + httpUrl
	actionPath = strings.ReplaceAll(actionPath, "{project_id}", client.ProjectID)
	actionPath = strings.ReplaceAll(actionPath, "{job_name}", jobName)

	actionOpts := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      buildFactoryRequestMoreHeaders(workspaceId),
		JSONBody:         utils.RemoveNil(buildStartJobImmediatelyBodyParams(d)),
	}

	requestBody, err := client.Request("POST", actionPath, &actionOpts)
	if err != nil {
		return nil, err
	}
	return utils.FlattenResponse(requestBody)
}

func getJobByName(client *golangsdk.ServiceClient, workspaceId, jobName, jobType string) (interface{}, error) {
	// The maximum value of limit is 100.
	httpUrl := "v1/{project_id}/jobs?limit=100"
	getPath := client.Endpoint + httpUrl
	getPath = strings.ReplaceAll(getPath, "{project_id}", client.ProjectID)
	// The job_name field supports fuzzy matching.
	getPath = fmt.Sprintf("%s&jobName=%v&jobType=%v", getPath, jobName, jobType)

	getOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"workspace": workspaceId,
		},
	}

	offset := 0
	for {
		listPathWithOffset := fmt.Sprintf("%s&offset=%d", getPath, offset)
		requestResp, err := client.Request("GET", listPathWithOffset, &getOpt)
		if err != nil {
			return nil, err
		}

		respBody, err := utils.FlattenResponse(requestResp)
		if err != nil {
			return nil, err
		}

		jobs := utils.PathSearch("jobs", respBody, make([]interface{}, 0)).([]interface{})
		if len(jobs) < 1 {
			break
		}

		job := utils.PathSearch(fmt.Sprintf("[?name=='%s']|[0]", jobName), jobs, nil)
		if job != nil {
			return job, nil
		}
		offset += len(jobs)
	}

	return nil, golangsdk.ErrDefault404{
		ErrUnexpectedResponseCode: golangsdk.ErrUnexpectedResponseCode{
			Method:    "GET",
			URL:       "/v1/{project_id}/jobs",
			RequestId: "NONE",
			Body:      []byte(fmt.Sprintf("the job (%s) is not found", jobName)),
		},
	}
}

func jobStateRefreshFunc(client *golangsdk.ServiceClient, workspaceId, jobName, jobType string,
	targets []string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		respBody, err := getJobByName(client, workspaceId, jobName, jobType)
		if err != nil {
			return respBody, "ERROR", err
		}

		var (
			jobStatus           = utils.PathSearch("status", respBody, "").(string)
			unexpectedJobStatus = []string{"EXCEPTION"}
		)
		if utils.StrSliceContains(unexpectedJobStatus, jobStatus) {
			return respBody, "ERROR", fmt.Errorf("unexpected job status (%s)", jobStatus)
		}

		if utils.StrSliceContains(targets, jobStatus) {
			return respBody, "COMPLETED", nil
		}

		return respBody, "PENDING", nil
	}
}

func getJobInstanceById(client *golangsdk.ServiceClient, workspaceId, jobName, instanceId string) (interface{}, error) {
	httpUrl := "v1/{project_id}/jobs/{job_name}/instances/{instance_id}"
	getPath := client.Endpoint + httpUrl
	getPath = strings.ReplaceAll(getPath, "{project_id}", client.ProjectID)
	getPath = strings.ReplaceAll(getPath, "{job_name}", jobName)
	getPath = strings.ReplaceAll(getPath, "{instance_id}", instanceId)

	getOpts := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      buildFactoryRequestMoreHeaders(workspaceId),
	}

	requestBody, err := client.Request("GET", getPath, &getOpts)
	if err != nil {
		return nil, err
	}
	return utils.FlattenResponse(requestBody)
}

func jobInstanceStateRefreshFunc(client *golangsdk.ServiceClient, workspaceId, jobName, jobType string,
	targets []string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		respBody, err := getJobInstanceById(client, workspaceId, jobName, jobType)
		if err != nil {
			return respBody, "ERROR", err
		}

		var (
			jobStatus           = utils.PathSearch("status", respBody, "").(string)
			unexpectedJobStatus = []string{"running-exception"}
		)
		if utils.StrSliceContains(unexpectedJobStatus, jobStatus) {
			return respBody, "ERROR", fmt.Errorf("unexpected job status (%s)", jobStatus)
		}

		if utils.StrSliceContains(targets, jobStatus) {
			return respBody, "COMPLETED", nil
		}

		return respBody, "PENDING", nil
	}
}

func doActionJob(ctx context.Context, client *golangsdk.ServiceClient, d *schema.ResourceData, timeout time.Duration) error {
	var (
		workspaceId = d.Get("workspace_id").(string)
		jobName     = d.Get("job_name").(string)
		processType = d.Get("process_type").(string)
		actionType  = d.Get("action").(string)
		err         error
		targets     []string
	)

	switch actionType {
	case "start":
		err = startJob(client, d)
		targets = []string{"NORMAL", "SCHEDULING"}
	case "stop":
		err = stopJob(client, d)
		targets = []string{"STOPPED", "PAUSED"}
	case "run-immediate":
		respBody, err := startJobImmediately(client, d)
		if err != nil {
			return err
		}
		// Immediate execution will not affect the current state, but it will generate an execution instance, whose
		// state needs to be monitored separately.
		instanceId := strconv.FormatFloat(utils.PathSearch("instanceId", respBody, float64(0)).(float64), 'f', -1, 64)
		instanceStateConf := &resource.StateChangeConf{
			Pending: []string{"PENDING"},
			Target:  []string{"COMPLETED"},
			// For the test run, failure is also a final status.
			Refresh:      jobInstanceStateRefreshFunc(client, workspaceId, jobName, instanceId, []string{"success", "fail"}),
			Timeout:      timeout,
			Delay:        10 * time.Second,
			PollInterval: 20 * time.Second,
		}
		respBody, err = instanceStateConf.WaitForStateContext(ctx)
		if err != nil {
			return fmt.Errorf("error waiting for the job (%s) action to become completed: %s", jobName, err)
		}
		// Set the instance status to the state file and skip the state refresh of the job (run immediately test just
		// checks the instance status).
		return d.Set("instance_status", utils.PathSearch("status", respBody, ""))
	default:
		return fmt.Errorf("invalid action type (%s)", actionType)
	}

	if err != nil {
		return err
	}

	stateConf := &resource.StateChangeConf{
		Pending:      []string{"PENDING"},
		Target:       []string{"COMPLETED"},
		Refresh:      jobStateRefreshFunc(client, workspaceId, jobName, processType, targets),
		Timeout:      timeout,
		Delay:        10 * time.Second,
		PollInterval: 20 * time.Second,
	}
	_, err = stateConf.WaitForStateContext(ctx)
	if err != nil {
		return fmt.Errorf("error waiting for the job (%s) action to become completed: %s", jobName, err)
	}
	return nil
}

func resourceFactoryJobActionCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		jobName = d.Get("job_name").(string)
	)
	client, err := cfg.NewServiceClient("dataarts-dlf", region)
	if err != nil {
		return diag.Errorf("error creating DataArts client: %s", err)
	}

	d.SetId(jobName)

	err = doActionJob(ctx, client, d, d.Timeout(schema.TimeoutCreate))
	if err != nil {
		return diag.Errorf("unable to operate status of the job (%s): %s", jobName, err)
	}

	return resourceFactoryJobActionRead(ctx, d, meta)
}

func resourceFactoryJobActionRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg         = meta.(*config.Config)
		region      = cfg.GetRegion(d)
		workspaceId = d.Get("workspace_id").(string)
		jobName     = d.Get("job_name").(string)
		jobType     = d.Get("process_type").(string)
		actionType  = d.Get("action").(string)
	)

	client, err := cfg.NewServiceClient("dataarts-dlf", region)
	if err != nil {
		return diag.Errorf("error creating DataArts client: %s", err)
	}

	job, err := getJobByName(client, workspaceId, jobName, jobType)
	if err != nil {
		return common.CheckDeletedDiag(d, common.ConvertExpected400ErrInto404Err(err, "error_code", actionResourceNotFoundCodes...),
			fmt.Sprintf("job (%s) not found: %s", jobName, err))
	}

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("job_name", utils.PathSearch("name", job, nil)),
		d.Set("status", utils.PathSearch("status", job, nil)),
		d.Set("action", actionType),
	)
	if err = mErr.ErrorOrNil(); err != nil {
		return diag.Errorf("error saving the fields of the job action: %s", err)
	}
	return nil
}

func resourceFactoryJobActionUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	client, err := cfg.NewServiceClient("dataarts-dlf", cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating DataArts client: %s", err)
	}

	if d.HasChange("action") {
		jobName := d.Get("job_name").(string)
		err = doActionJob(ctx, client, d, d.Timeout(schema.TimeoutUpdate))
		if err != nil {
			return diag.Errorf("error updating DataArts job (%s) status: %s", jobName, err)
		}
	}

	return resourceFactoryJobActionRead(ctx, d, meta)
}

func resourceFactoryJobActionDelete(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	errorMsg := `This resource is only a one-time action resource for changing job status. Deleting this resource will
not change the current status, but will only remove the resource information from the tfstate file.`
	return diag.Diagnostics{
		diag.Diagnostic{
			Severity: diag.Warning,
			Summary:  errorMsg,
		},
	}
}
