package codeartspipeline

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

// @API CodeArtsPipeline GET /v5/{project_id}/api/pipelines/{pipeline_id}/pipeline-runs/detail
func DataSourceCodeartsPipelineRunDetail() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceCodeartsPipelineRunDetailRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `Specifies the region in which to query the resource. If omitted, the provider-level region will be used.`,
			},
			"project_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `Specifies the CodeArts project ID.`,
			},
			"pipeline_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `Specifies the pipeline ID.`,
			},
			"pipeline_run_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specifies the pipeline run ID.`,
			},
			"manifest_version": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Indicates the pipeline version.`,
			},
			"name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Indicates the pipeline name.`,
			},
			"description": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Indicates the pipeline running description.`,
			},
			"is_publish": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: `Indicates whether the pipeline is a change-triggered pipeline.`,
			},
			"executor_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Indicates the executor ID.`,
			},
			"executor_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Indicates the executor name.`,
			},
			"status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Indicates the pipeline run status.`,
			},
			"trigger_type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Indicates the trigger type.`,
			},
			"run_number": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: `Indicates the pipeline running sequence number.`,
			},
			"start_time": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: `Indicates the start time.`,
			},
			"end_time": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: `Indicates the end time.`,
			},
			"component_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Indicates the microservice ID.`,
			},
			"language": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Indicates the language.`,
			},
			"subject_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Indicates the pipeline run ID.`,
			},
			"group_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Indicates the pipeline group ID.`,
			},
			"group_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Indicates the pipeline group name.`,
			},
			"detail_url": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Indicates the pipeline detail URL.`,
			},
			"current_system_time": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: `Indicates the current system time.`,
			},
			"stages": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: `Indicates the stage running information.`,
				Elem:        schemaPipelineRunDetailStages(),
			},
			"sources": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: `Indicates the pipeline source information.`,
				Elem:        schemaPipelineRunDetailSources(),
			},
			"artifacts": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: `Indicates the artifacts after running a pipeline.`,
				Elem:        schemaPipelineRunDetailArtifact(),
			},
		},
	}
}

func schemaPipelineRunDetailStages() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Indicates the stage ID.`,
			},
			"category": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Indicates the stage type.`,
			},
			"name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Indicates the stage name.`,
			},
			"identifier": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Indicates the unique identifier of a stage.`,
			},
			"run_always": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: `Indicates whether to always run.`,
			},
			"parallel": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Indicates whether to execute jobs in parallel.`,
			},
			"is_select": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: `Indicates whether to select.`,
			},
			"sequence": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: `Indicates the serial number.`,
			},
			"depends_on": {
				Type:        schema.TypeList,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: `Indicates the dependency.`,
			},
			"condition": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Indicates the running conditions.`,
			},
			"status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Indicates the stage status.`,
			},
			"start_time": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: `Indicates the start time.`,
			},
			"end_time": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: `Indicates the end time.`,
			},
			"pre": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: `Indicates the stage entry.`,
				Elem:        schemaPipelineRunDetailStepRunObject(),
			},
			"post": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: `Indicates the stage exit.`,
				Elem:        schemaPipelineRunDetailStepRunObject(),
			},
			"jobs": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: `Indicates the job running information.`,
				Elem:        schemaPipelineRunDetailStagesJobs(),
			},
		},
	}
}

func schemaPipelineRunDetailStagesJobs() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Indicates the job ID.`,
			},
			"category": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Indicates the job type.`,
			},
			"sequence": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: `Indicates the serial number.`,
			},
			"async": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: `Indicates whether it is asynchronous.`,
			},
			"name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Indicates the job name.`,
			},
			"identifier": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Indicates the unique identifier of a job.`,
			},
			"depends_on": {
				Type:        schema.TypeList,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: `Indicates the dependency.`,
			},
			"condition": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Indicates the running conditions.`,
			},
			"resource": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Indicates the execution resources.`,
			},
			"is_select": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: `Indicates whether the parameter is selected.`,
			},
			"timeout": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Indicates the job timeout settings.`,
			},
			"last_dispatch_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Indicates the ID of the job delivered last time.`,
			},
			"status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Indicates the job status.`,
			},
			"message": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Indicates the error message.`,
			},
			"start_time": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: `Indicates the start time.`,
			},
			"end_time": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: `Indicates the end time.`,
			},
			"steps": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: `Indicates the step running information.`,
				Elem:        schemaPipelineRunDetailStepRunObject(),
			},
			"exec_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Indicates the job execution ID.`,
			},
		},
	}
}

func schemaPipelineRunDetailStepRunObject() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Indicates the step ID.`,
			},
			"name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Indicates the step name.`,
			},
			"task": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Indicates the step extension name.`,
			},
			"business_type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Indicates the extension type.`,
			},
			"inputs": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: `Indicates the step running information.`,
				Elem:        schemaPipelineRunDetailStepRunObjectInput(),
			},
			"sequence": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: `Indicates the serial number.`,
			},
			"official_task_version": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Indicates the official extension version.`,
			},
			"identifier": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Indicates the unique identifier.`,
			},
			"multi_step_editable": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: `Indicates whether the parameter is editable.`,
			},
			"endpoint_ids": {
				Type:        schema.TypeList,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: `Indicates the step name.`,
			},
			"last_dispatch_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Indicates the ID of the job delivered last time.`,
			},
			"message": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Indicates the error message.`,
			},
			"status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Indicates the step status.`,
			},
			"start_time": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: `Indicates the start time.`,
			},
			"end_time": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: `Indicates the end time.`,
			},
		},
	}
}

func schemaPipelineRunDetailStepRunObjectInput() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"key": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Indicates the parameter name.`,
			},
			"value": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Indicates the parameter value.`,
			},
		},
	}
}

func schemaPipelineRunDetailSources() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Indicates the source type.`,
			},
			"params": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: `Indicates the source parameters.`,
				Elem:        schemaPipelineRunDetailSourcesParam(),
			},
		},
	}
}

func schemaPipelineRunDetailSourcesParam() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"git_type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Indicates the code repository type.`,
			},
			"git_url": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Indicates the HTTPS address of the Git repository.`,
			},
			"ssh_git_url": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Indicates the SSH address of the Git repository.`,
			},
			"web_url": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Indicates the address of the code repository page.`,
			},
			"repo_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Indicates the code repository name.`,
			},
			"alias": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Indicates the code repository alias.`,
			},
			"codehub_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Indicates the CodeArts Repo code repository ID.`,
			},
			"default_branch": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Indicates the default branch.`,
			},
			"endpoint_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Indicates the ID of the code source endpoint.`,
			},
			"build_params": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: `Indicates the build parameters.`,
				Elem:        schemaPipelineRunDetailSourcesParamBuildParams(),
			},
		},
	}
}

func schemaPipelineRunDetailSourcesParamBuildParams() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"action": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Indicates the action.`,
			},
			"build_type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Indicates the code repository trigger type.`,
			},
			"commit_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Indicates the commit ID.`,
			},
			"event_type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Indicates the event type.`,
			},
			"merge_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Indicates the merge ID.`,
			},
			"message": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Indicates the commit message.`,
			},
			"target_branch": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Indicates the branch that triggers the pipeline execution.`,
			},
			"tag": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Indicates the tag that triggers the pipeline execution.`,
			},
			"source_branch": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Indicates the source branch.`,
			},
			"codehub_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Indicates the CodeArts Repo code repository ID.`,
			},
			"source_codehub_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Indicates the ID of the source Repo code repository.`,
			},
			"source_codehub_url": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Indicates the address of the source Repo code repository.`,
			},
			"source_codehub_http_url": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Indicates the HTTP address of the source Repo code repository.`,
			},
		},
	}
}

func schemaPipelineRunDetailArtifact() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Indicates the artifact name.`,
			},
			"download_url": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Indicates the artifact download address.`,
			},
			"version": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Indicates the artifact version number.`,
			},
			"package_type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Indicates the artifact type.`,
			},
		},
	}
}

func dataSourceCodeartsPipelineRunDetailRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	client, err := cfg.NewServiceClient("codearts_pipeline", region)
	if err != nil {
		return diag.Errorf("error creating CodeArts Pipeline client: %s", err)
	}

	projectId := d.Get("project_id").(string)
	pipelineId := d.Get("pipeline_id").(string)

	getHttpUrl := "v5/{project_id}/api/pipelines/{pipeline_id}/pipeline-runs/detail"
	getPath := client.Endpoint + getHttpUrl
	getPath = strings.ReplaceAll(getPath, "{project_id}", projectId)
	getPath = strings.ReplaceAll(getPath, "{pipeline_id}", pipelineId)

	if v, ok := d.GetOk("pipeline_run_id"); ok {
		getPath += "?pipeline_run_id=" + v.(string)
	}

	getOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	getResp, err := client.Request("GET", getPath, &getOpt)
	if err != nil {
		return diag.Errorf("error retrieving pipeline run detail: %s", err)
	}
	getRespBody, err := utils.FlattenResponse(getResp)
	if err != nil {
		return diag.Errorf("error flattening response: %s", err)
	}

	if err := checkResponseError(getRespBody, ""); err != nil {
		return diag.Errorf("error retrieving pipeline run detail: %s", err)
	}

	id, err := uuid.GenerateUUID()
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(id)

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("manifest_version", utils.PathSearch("manifest_version", getRespBody, nil)),
		d.Set("name", utils.PathSearch("name", getRespBody, nil)),
		d.Set("description", utils.PathSearch("description", getRespBody, nil)),
		d.Set("is_publish", utils.PathSearch("is_publish", getRespBody, nil)),
		d.Set("executor_id", utils.PathSearch("executor_id", getRespBody, nil)),
		d.Set("executor_name", utils.PathSearch("executor_name", getRespBody, nil)),
		d.Set("status", utils.PathSearch("status", getRespBody, nil)),
		d.Set("trigger_type", utils.PathSearch("trigger_type", getRespBody, nil)),
		d.Set("run_number", utils.PathSearch("run_number", getRespBody, nil)),
		d.Set("start_time", utils.PathSearch("start_time", getRespBody, nil)),
		d.Set("end_time", utils.PathSearch("end_time", getRespBody, nil)),
		d.Set("component_id", utils.PathSearch("component_id", getRespBody, nil)),
		d.Set("language", utils.PathSearch("language", getRespBody, nil)),
		d.Set("subject_id", utils.PathSearch("subject_id", getRespBody, nil)),
		d.Set("group_id", utils.PathSearch("group_id", getRespBody, nil)),
		d.Set("group_name", utils.PathSearch("group_name", getRespBody, nil)),
		d.Set("detail_url", utils.PathSearch("detail_url", getRespBody, nil)),
		d.Set("current_system_time", utils.PathSearch("current_system_time", getRespBody, nil)),
		d.Set("stages", flattenPipelineRunStages(getRespBody)),
		d.Set("sources", flattenPipelineRunSources(getRespBody)),
		d.Set("artifacts", flattenPipelineRunArtifacts(getRespBody)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenPipelineRunStages(resp interface{}) []map[string]interface{} {
	stages := utils.PathSearch("stages", resp, make([]interface{}, 0)).([]interface{})
	if len(stages) == 0 {
		return nil
	}

	result := make([]map[string]interface{}, 0, len(stages))
	for _, stage := range stages {
		result = append(result, map[string]interface{}{
			"id":         utils.PathSearch("id", stage, nil),
			"category":   utils.PathSearch("category", stage, nil),
			"name":       utils.PathSearch("name", stage, nil),
			"status":     utils.PathSearch("status", stage, nil),
			"start_time": utils.PathSearch("start_time", stage, nil),
			"end_time":   utils.PathSearch("end_time", stage, nil),
			"identifier": utils.PathSearch("identifier", stage, nil),
			"run_always": utils.PathSearch("run_always", stage, nil),
			"parallel":   utils.PathSearch("parallel", stage, nil),
			"is_select":  utils.PathSearch("is_select", stage, nil),
			"sequence":   utils.PathSearch("sequence", stage, nil),
			"depends_on": utils.PathSearch("depends_on", stage, nil),
			"condition":  utils.PathSearch("condition", stage, nil),
			"pre":        flattenPipelineRunSteps(utils.PathSearch("pre", stage, make([]interface{}, 0)).([]interface{})),
			"post":       flattenPipelineRunSteps(utils.PathSearch("post", stage, make([]interface{}, 0)).([]interface{})),
			"jobs":       flattenPipelineRunJobs(stage),
		})
	}

	return result
}

func flattenPipelineRunJobs(stage interface{}) []map[string]interface{} {
	jobs := utils.PathSearch("jobs", stage, make([]interface{}, 0)).([]interface{})
	if len(jobs) == 0 {
		return nil
	}

	result := make([]map[string]interface{}, 0, len(jobs))
	for _, job := range jobs {
		result = append(result, map[string]interface{}{
			"id":               utils.PathSearch("id", job, nil),
			"name":             utils.PathSearch("name", job, nil),
			"status":           utils.PathSearch("status", job, nil),
			"start_time":       utils.PathSearch("start_time", job, nil),
			"end_time":         utils.PathSearch("end_time", job, nil),
			"category":         utils.PathSearch("category", job, nil),
			"sequence":         utils.PathSearch("sequence", job, nil),
			"async":            utils.PathSearch("async", job, nil),
			"identifier":       utils.PathSearch("identifier", job, nil),
			"depends_on":       utils.PathSearch("depends_on", job, nil),
			"condition":        utils.PathSearch("condition", job, nil),
			"resource":         utils.PathSearch("resource", job, nil),
			"is_select":        utils.PathSearch("is_select", job, nil),
			"timeout":          utils.PathSearch("timeout", job, nil),
			"last_dispatch_id": utils.PathSearch("last_dispatch_id", job, nil),
			"message":          utils.PathSearch("message", job, nil),
			"steps":            flattenPipelineRunSteps(utils.PathSearch("steps", job, make([]interface{}, 0)).([]interface{})),
			"exec_id":          utils.PathSearch("exec_id", job, nil),
		})
	}

	return result
}

func flattenPipelineRunSteps(step []interface{}) []map[string]interface{} {
	if len(step) == 0 {
		return nil
	}

	result := make([]map[string]interface{}, 0, len(step))
	for _, step := range step {
		result = append(result, map[string]interface{}{
			"id":                    utils.PathSearch("id", step, nil),
			"name":                  utils.PathSearch("name", step, nil),
			"task":                  utils.PathSearch("task", step, nil),
			"business_type":         utils.PathSearch("business_type", step, nil),
			"inputs":                flattenPipelineRunStepsInputs(step),
			"sequence":              utils.PathSearch("sequence", step, nil),
			"official_task_version": utils.PathSearch("official_task_version", step, nil),
			"identifier":            utils.PathSearch("identifier", step, nil),
			"multi_step_editable":   utils.PathSearch("multi_step_editable", step, nil),
			"endpoint_ids":          utils.PathSearch("endpoint_ids", step, nil),
			"last_dispatch_id":      utils.PathSearch("last_dispatch_id", step, nil),
			"message":               utils.PathSearch("message", step, nil),
			"status":                utils.PathSearch("status", step, nil),
			"start_time":            utils.PathSearch("start_time", step, nil),
			"end_time":              utils.PathSearch("end_time", step, nil),
		})
	}

	return result
}

func flattenPipelineRunStepsInputs(resp interface{}) []map[string]interface{} {
	inputs := utils.PathSearch("inputs", resp, make([]interface{}, 0)).([]interface{})
	if len(inputs) == 0 {
		return nil
	}

	result := make([]map[string]interface{}, 0, len(inputs))
	for _, input := range inputs {
		result = append(result, map[string]interface{}{
			"key": utils.PathSearch("key", input, nil),
			// value is an object
			"value": encodeIntoJson(utils.PathSearch("value", input, nil)),
		})
	}

	return result
}

func flattenPipelineRunSources(resp interface{}) []map[string]interface{} {
	sources := utils.PathSearch("sources", resp, make([]interface{}, 0)).([]interface{})
	if len(sources) == 0 {
		return nil
	}

	result := make([]map[string]interface{}, 0, len(sources))
	for _, source := range sources {
		result = append(result, map[string]interface{}{
			"type":   utils.PathSearch("type", source, nil),
			"params": flattenPipelineRunSourceParams(source),
		})
	}

	return result
}

func flattenPipelineRunSourceParams(source interface{}) []map[string]interface{} {
	params := utils.PathSearch("params", source, nil)
	if params == nil {
		return nil
	}

	result := []map[string]interface{}{
		{
			"git_type":       utils.PathSearch("git_type", params, nil),
			"git_url":        utils.PathSearch("git_url", params, nil),
			"ssh_git_url":    utils.PathSearch("ssh_git_url", params, nil),
			"web_url":        utils.PathSearch("web_url", params, nil),
			"repo_name":      utils.PathSearch("repo_name", params, nil),
			"alias":          utils.PathSearch("alias", params, nil),
			"codehub_id":     utils.PathSearch("codehub_id", params, nil),
			"default_branch": utils.PathSearch("default_branch", params, nil),
			"endpoint_id":    utils.PathSearch("endpoint_id", params, nil),
			"build_params":   flattenPipelineRunBuildParams(params),
		},
	}

	return result
}

func flattenPipelineRunBuildParams(params interface{}) []map[string]interface{} {
	buildParams := utils.PathSearch("build_params", params, nil)
	if buildParams == nil {
		return nil
	}

	result := []map[string]interface{}{
		{
			"action":                  utils.PathSearch("action", buildParams, nil),
			"build_type":              utils.PathSearch("build_type", buildParams, nil),
			"commit_id":               utils.PathSearch("commit_id", buildParams, nil),
			"event_type":              utils.PathSearch("event_type", buildParams, nil),
			"merge_id":                utils.PathSearch("merge_id", buildParams, nil),
			"message":                 utils.PathSearch("message", buildParams, nil),
			"target_branch":           utils.PathSearch("target_branch", buildParams, nil),
			"tag":                     utils.PathSearch("tag", buildParams, nil),
			"source_branch":           utils.PathSearch("source_branch", buildParams, nil),
			"codehub_id":              utils.PathSearch("codehub_id", buildParams, nil),
			"source_codehub_id":       utils.PathSearch("source_codehub_id", buildParams, nil),
			"source_codehub_url":      utils.PathSearch("source_codehub_url", buildParams, nil),
			"source_codehub_http_url": utils.PathSearch("source_codehub_http_url", buildParams, nil),
		},
	}

	return result
}

func flattenPipelineRunArtifacts(resp interface{}) []map[string]interface{} {
	artifacts := utils.PathSearch("artifacts", resp, make([]interface{}, 0)).([]interface{})
	if len(artifacts) == 0 {
		return nil
	}

	result := make([]map[string]interface{}, 0, len(artifacts))
	for _, artifact := range artifacts {
		result = append(result, map[string]interface{}{
			"name":         utils.PathSearch("name", artifact, nil),
			"download_url": utils.PathSearch("downloadUrl", artifact, nil),
			"version":      utils.PathSearch("version", artifact, nil),
			"package_type": utils.PathSearch("packageType", artifact, nil),
		})
	}

	return result
}
