package codeartspipeline

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

var (
	//nolint:revive
	actionHttpUrl = map[string]string{
		// POST
		"run":            "v5/{project_id}/api/pipelines/{pipeline_id}/run",
		"stop":           "v5/{project_id}/api/pipelines/{pipeline_id}/pipeline-runs/{pipeline_run_id}/stop",
		"pass":           "v5/{project_id}/api/pipelines/{pipeline_id}/pipeline-runs/{pipeline_run_id}/jobs/{job_run_id}/steps/{step_run_id}/pass",
		"refuse":         "v5/{project_id}/api/pipelines/{pipeline_id}/pipeline-runs/{pipeline_run_id}/jobs/{job_run_id}/steps/{step_run_id}/refuse",
		"delay-pass":     "v5/{project_id}/api/pipelines/{pipeline_id}/pipeline-runs/{pipeline_run_id}/jobs/{job_run_id}/steps/{step_run_id}/delay-pass",
		"delay-refuse":   "v5/{project_id}/api/pipelines/{pipeline_id}/pipeline-runs/{pipeline_run_id}/jobs/{job_run_id}/steps/{step_run_id}/delay-refuse",
		"delay":          "v5/{project_id}/api/pipelines/{pipeline_id}/pipeline-runs/{pipeline_run_id}/jobs/{job_run_id}/steps/{step_run_id}/delay",
		"manual-pass":    "v5/{project_id}/api/pipelines/{pipeline_id}/pipeline-runs/{pipeline_run_id}/steps/{step_run_id}/manual/pass",
		"manual-refuse":  "v5/{project_id}/api/pipelines/{pipeline_id}/pipeline-runs/{pipeline_run_id}/steps/{step_run_id}/manual/refuse",
		"resume":         "v5/{project_id}/api/pipelines/{pipeline_id}/pipeline-runs/{pipeline_run_id}/jobs/{job_run_id}/steps/{step_run_id}/resume",
		"cancel-queuing": "v5/{project_id}/api/pipelines/{pipeline_id}/{pipeline_run_id}/cancel-queuing/{queue_id}",

		// PUT
		"retry": "v5/{project_id}/api/pipelines/{pipeline_id}/pipeline-runs/{pipeline_run_id}/retry",
	}

	queuingError1 = "DEVPIPE.00011169" // Maximun concurrency of a single pipeline is limited, enter the queue state
	queuingError2 = "DEVPIPE.00011170" // Maximun concurrency of single pipelines under a tenant is limited, enter the queue state

	pipelineActionNonUpdatableParams = []string{
		"project_id", "pipeline_id", "action", "pipeline_run_id", "job_run_id", "step_run_id", "queue_id", "description",
		"choose_jobs", "choose_stages", "variables", "sources",
	}
)

//nolint:revive
// @API CodeArtsPipeline POST /v5/{project_id}/api/pipelines/{pipeline_id}/run
// @API CodeArtsPipeline POST /v5/{project_id}/api/pipelines/{pipeline_id}/pipeline-runs/{pipeline_run_id}/stop
// @API CodeArtsPipeline POST /v5/{project_id}/api/pipelines/{pipeline_id}/pipeline-runs/{pipeline_run_id}/jobs/{job_run_id}/steps/{step_run_id}/pass
// @API CodeArtsPipeline POST /v5/{project_id}/api/pipelines/{pipeline_id}/pipeline-runs/{pipeline_run_id}/jobs/{job_run_id}/steps/{step_run_id}/refuse
// @API CodeArtsPipeline POST /v5/{project_id}/api/pipelines/{pipeline_id}/pipeline-runs/{pipelineRunId}/jobs/{jobRunId}/steps/{stepRunId}/delay-pass
// @API CodeArtsPipeline POST /v5/{project_id}/api/pipelines/{pipeline_id}/pipeline-runs/{pipelineRunId}/jobs/{jobRunId}/steps/{stepRunId}/delay-refuse
// @API CodeArtsPipeline POST /v5/{project_id}/api/pipelines/{pipeline_id}/pipeline-runs/{pipelineRunId}/jobs/{jobRunId}/steps/{stepRunId}/delay
// @API CodeArtsPipeline POST /v5/{project_id}/api/pipelines/{pipeline_id}/pipeline-runs/{pipeline_run_id}/steps/{step_run_id}/manual/pass
// @API CodeArtsPipeline POST /v5/{project_id}/api/pipelines/{pipeline_id}/pipeline-runs/{pipeline_run_id}/steps/{stepRunId}/manual/refuse
// @API CodeArtsPipeline POST /v5/{project_id}/api/pipelines/{pipeline_id}/pipeline-runs/{pipelineRunId}/jobs/{jobRunId}/steps/{stepRunId}/resume
// @API CodeArtsPipeline POST /v5/{project_id}/api/pipelines/{pipeline_id}/{pipelineRunId}/cancel-queuing/{id}
// @API CodeArtsPipeline PUT /v5/{project_id}/api/pipelines/{pipeline_id}/pipeline-runs/{pipeline_run_id}/retry

func ResourceCodeArtsPipelineAction() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourcePipelineActionCreate,
		ReadContext:   resourcePipelineActionRead,
		UpdateContext: resourcePipelineActionUpdate,
		DeleteContext: resourcePipelineActionDelete,

		CustomizeDiff: config.FlexibleForceNew(pipelineActionNonUpdatableParams),

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
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
			"action": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `Specifies the action.`,
			},
			"pipeline_run_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `Specifies the pipeline run ID.`,
			},
			"job_run_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specifies the pipeline job run ID.`,
			},
			"step_run_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specifies the pipeline step run ID.`,
			},
			"queue_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `Specifies the queued pipeline step run ID.`,
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specifies the running description.`,
			},
			"sources": {
				Type:        schema.TypeSet,
				Optional:    true,
				Description: `Specifies the code source information list.`,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"type": {
							Type:        schema.TypeString,
							Required:    true,
							Description: `Specifies the pipeline source type.`,
						},
						"params": {
							Type:        schema.TypeList,
							Required:    true,
							MaxItems:    1,
							Description: `Specifies the source parameters.`,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"git_type": {
										Type:        schema.TypeString,
										Required:    true,
										Description: `Specifies the code repository type.`,
									},
									"git_url": {
										Type:        schema.TypeString,
										Required:    true,
										Description: `Specifies the HTTPS address of the Git repository.`,
									},
									"alias": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: `Specifies the code repository alias.`,
									},
									"codehub_id": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: `Specifies the CodeArts Repo code repository ID.`,
									},
									"default_branch": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: `Specifies the default branch of the code repository for pipeline execution.`,
									},
									"endpoint_id": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: `Specifies the ID of the code source endpoint.`,
									},
									"build_params": {
										Type:        schema.TypeList,
										Optional:    true,
										MaxItems:    1,
										Description: `Specifies the detailed build parameters.`,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"build_type": {
													Type:        schema.TypeString,
													Required:    true,
													Description: `Specifies the code repository trigger type.`,
												},
												"target_branch": {
													Type:        schema.TypeString,
													Optional:    true,
													Description: `Specifies the branch that triggers the pipeline execution.`,
												},
												"tag": {
													Type:        schema.TypeString,
													Optional:    true,
													Description: `Specifies the tag that triggers the pipeline execution.`,
												},
												"event_type": {
													Type:        schema.TypeString,
													Optional:    true,
													Description: `Specifies the event type that triggers the pipeline execution.`,
												},
											},
										},
									},
									"change_request_ids": {
										Type:        schema.TypeSet,
										Optional:    true,
										Elem:        &schema.Schema{Type: schema.TypeString},
										Description: `Specifies the change IDs of the change-triggered pipeline.`,
									},
								},
							},
						},
					},
				},
			},
			"variables": {
				Type:        schema.TypeSet,
				Optional:    true,
				Description: `Specifies the custom parameters used.`,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:        schema.TypeString,
							Required:    true,
							Description: `Specifies the parameter name.`,
						},
						"value": {
							Type:        schema.TypeString,
							Required:    true,
							Description: `Specifies the parameter value.`,
						},
					},
				},
			},
			"choose_jobs": {
				Type:        schema.TypeSet,
				Optional:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: `Specifies the selected pipeline jobs.`,
			},
			"choose_stages": {
				Type:        schema.TypeSet,
				Optional:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: `Specifies the selected pipeline stages.`,
			},
			"enable_force_new": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"true", "false"}, false),
				Description:  utils.SchemaDesc("", utils.SchemaDescInput{Internal: true}),
			},
		},
	}
}

func resourcePipelineActionCreate(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	client, err := cfg.NewServiceClient("codearts_pipeline", region)
	if err != nil {
		return diag.Errorf("error creating CodeArts Pipeline client: %s", err)
	}

	projectId := d.Get("project_id").(string)
	pipelineId := d.Get("pipeline_id").(string)
	pipelineRunId := d.Get("pipeline_run_id").(string)
	jobRunId := d.Get("job_run_id").(string)
	stepRunId := d.Get("step_run_id").(string)
	queueId := d.Get("queue_id").(string)
	action := d.Get("action").(string)
	switch action {
	case "run":
		return runPipeline(client, d)
	case "stop":
		if err = handlePipeline(client, d); err != nil {
			return diag.FromErr(err)
		}

		d.SetId(fmt.Sprintf("%s/%s/%s", projectId, pipelineId, pipelineRunId))
	case "pass", "refuse", "delay-pass", "delay-refuse", "delay", "resume":
		if err = handlePipeline(client, d); err != nil {
			return diag.FromErr(err)
		}

		d.SetId(fmt.Sprintf("%s/%s/%s/%s/%s", projectId, pipelineId, pipelineRunId, jobRunId, stepRunId))
	case "manual-pass", "manual-refuse":
		if err = handlePipeline(client, d); err != nil {
			return diag.FromErr(err)
		}

		d.SetId(fmt.Sprintf("%s/%s/%s/%s", projectId, pipelineId, pipelineRunId, stepRunId))
	case "retry":
		return retryPipeline(client, d)
	case "cancel-queuing":
		if err = cancelQueuingPipeline(client, d); err != nil {
			return diag.FromErr(err)
		}

		d.SetId(fmt.Sprintf("%s/%s/%s/%s", projectId, pipelineId, pipelineRunId, queueId))
	default:
		return diag.Errorf("unsupport action")
	}

	return nil
}

func runPipeline(client *golangsdk.ServiceClient, d *schema.ResourceData) diag.Diagnostics {
	projectId := d.Get("project_id").(string)
	pipelineId := d.Get("pipeline_id").(string)

	httpUrl := "v5/{project_id}/api/pipelines/{pipeline_id}/run"
	createPath := client.Endpoint + httpUrl
	createPath = strings.ReplaceAll(createPath, "{project_id}", projectId)
	createPath = strings.ReplaceAll(createPath, "{pipeline_id}", pipelineId)
	createOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		JSONBody:         utils.RemoveNil(buildPipelineActionRunBodyParams(d)),
	}

	createResp, err := client.Request("POST", createPath, &createOpt)
	if err != nil {
		return diag.Errorf("error running pipeline: %s", err)
	}
	createRespBody, err := utils.FlattenResponse(createResp)
	if err != nil {
		return diag.Errorf("error flattening run pipeline response: %s", err)
	}

	errorCode := utils.PathSearch("error_code", createRespBody, "")
	if errorCode != "" {
		errorMsg := utils.PathSearch("error_msg", createRespBody, "").(string)
		if errorCode == queuingError1 || errorCode == queuingError2 {
			pipelineRunId := utils.PathSearch("error.pipeline_run_id", createRespBody, "").(string)
			if pipelineRunId == "" {
				return diag.Errorf("unable to find the CodeArts Pipeline run ID from the API response")
			}

			d.SetId(fmt.Sprintf("%s/%s/%s", projectId, pipelineId, pipelineRunId))
			d.Set("pipeline_run_id", pipelineRunId)
			d.Set("queue_id", utils.PathSearch("error.queue_id", createRespBody, nil))

			return diag.Diagnostics{
				diag.Diagnostic{
					Severity: diag.Warning,
					Summary:  errorMsg,
				},
			}
		}

		return diag.Errorf("error running pipeline, error code: %s, error message: %s", errorCode, errorMsg)
	}

	pipelineRunId := utils.PathSearch("pipeline_run_id", createRespBody, "").(string)
	if pipelineRunId == "" {
		return diag.Errorf("unable to find the CodeArts Pipeline run ID from the API response")
	}

	d.SetId(fmt.Sprintf("%s/%s/%s", projectId, pipelineId, pipelineRunId))
	d.Set("pipeline_run_id", pipelineRunId)

	return nil
}

func handlePipeline(client *golangsdk.ServiceClient, d *schema.ResourceData) error {
	action := d.Get("action").(string)
	httpUrl := actionHttpUrl[action]
	createPath := client.Endpoint + httpUrl
	createPath = strings.ReplaceAll(createPath, "{project_id}", d.Get("project_id").(string))
	createPath = strings.ReplaceAll(createPath, "{pipeline_id}", d.Get("pipeline_id").(string))
	createPath = strings.ReplaceAll(createPath, "{pipeline_run_id}", d.Get("pipeline_run_id").(string))
	createPath = strings.ReplaceAll(createPath, "{job_run_id}", d.Get("job_run_id").(string))
	createPath = strings.ReplaceAll(createPath, "{step_run_id}", d.Get("step_run_id").(string))
	createOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	createResp, err := client.Request("POST", createPath, &createOpt)
	if err != nil {
		return fmt.Errorf("error handling CodeArts Pipeline with action(%s): %s", action, err)
	}
	createRespBody, err := utils.FlattenResponse(createResp)
	if err != nil {
		return err
	}

	if err := checkResponseError(createRespBody, ""); err != nil {
		return fmt.Errorf("error handling CodeArts Pipeline with action(%s): %s", action, err)
	}

	return nil
}

// PUT
func retryPipeline(client *golangsdk.ServiceClient, d *schema.ResourceData) diag.Diagnostics {
	projectId := d.Get("project_id").(string)
	pipelineId := d.Get("pipeline_id").(string)
	pipelineRunId := d.Get("pipeline_run_id").(string)
	action := d.Get("action").(string)

	httpUrl := actionHttpUrl[action]
	createPath := client.Endpoint + httpUrl
	createPath = strings.ReplaceAll(createPath, "{project_id}", projectId)
	createPath = strings.ReplaceAll(createPath, "{pipeline_id}", pipelineId)
	createPath = strings.ReplaceAll(createPath, "{pipeline_run_id}", pipelineRunId)
	createOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	createResp, err := client.Request("PUT", createPath, &createOpt)
	if err != nil {
		return diag.Errorf("error retrying pipeline:%s", err)
	}
	createRespBody, err := utils.FlattenResponse(createResp)
	if err != nil {
		return diag.Errorf("error flattening retry pipeline response:%s", err)
	}

	errorCode := utils.PathSearch("error_code", createRespBody, "")
	if errorCode != "" {
		errorMsg := utils.PathSearch("error_msg", createRespBody, "").(string)
		if errorCode == queuingError1 || errorCode == queuingError2 {
			d.SetId(fmt.Sprintf("%s/%s/%s", projectId, pipelineId, pipelineRunId))
			d.Set("queue_id", utils.PathSearch("error.queue_id", createRespBody, nil))

			return diag.Diagnostics{
				diag.Diagnostic{
					Severity: diag.Warning,
					Summary:  errorMsg,
				},
			}
		}

		return diag.Errorf("error retrying pipeline, error code: %s, error message: %s", errorCode, errorMsg)
	}

	d.SetId(fmt.Sprintf("%s/%s/%s", projectId, pipelineId, pipelineRunId))

	return nil
}

func cancelQueuingPipeline(client *golangsdk.ServiceClient, d *schema.ResourceData) error {
	action := d.Get("action").(string)
	httpUrl := actionHttpUrl[action]
	createPath := client.Endpoint + httpUrl
	createPath = strings.ReplaceAll(createPath, "{project_id}", d.Get("project_id").(string))
	createPath = strings.ReplaceAll(createPath, "{pipeline_id}", d.Get("pipeline_id").(string))
	createPath = strings.ReplaceAll(createPath, "{pipeline_run_id}", d.Get("pipeline_run_id").(string))
	createPath = strings.ReplaceAll(createPath, "{queue_id}", d.Get("queue_id").(string))
	createOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		JSONBody:         map[string]interface{}{},
	}

	createResp, err := client.Request("POST", createPath, &createOpt)
	if err != nil {
		return fmt.Errorf("error handling CodeArts Pipeline with action(%s): %s", action, err)
	}
	createRespBody, err := utils.FlattenResponse(createResp)
	if err != nil {
		return err
	}

	if err := checkResponseError(createRespBody, ""); err != nil {
		return fmt.Errorf("error handling CodeArts Pipeline with action(%s): %s", action, err)
	}

	return nil
}

func buildPipelineActionRunBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"description":   d.Get("description"),
		"sources":       buildPipelineActionSources(d),
		"variables":     buildPipelineActionVariables(d),
		"choose_jobs":   utils.ValueIgnoreEmpty(d.Get("choose_jobs").(*schema.Set).List()),
		"choose_stages": utils.ValueIgnoreEmpty(d.Get("choose_stages").(*schema.Set).List()),
	}

	return bodyParams
}

func buildPipelineActionSources(d *schema.ResourceData) interface{} {
	rawSources := d.Get("sources").(*schema.Set).List()
	if len(rawSources) == 0 {
		return nil
	}

	sources := make([]map[string]interface{}, 0, len(rawSources))
	for _, s := range rawSources {
		if source, ok := s.(map[string]interface{}); ok {
			sourceMap := map[string]interface{}{
				"type":   utils.ValueIgnoreEmpty(source["type"]),
				"params": buildPipelineActionSourceParams(source["params"].([]interface{})),
			}
			sources = append(sources, sourceMap)
		}
	}

	return sources
}

func buildPipelineActionSourceParams(paramsList []interface{}) interface{} {
	if len(paramsList) == 0 {
		return nil
	}

	params := paramsList[0].(map[string]interface{})
	result := map[string]interface{}{
		"git_type":           params["git_type"],
		"git_url":            params["git_url"],
		"alias":              utils.ValueIgnoreEmpty(params["alias"]),
		"codehub_id":         utils.ValueIgnoreEmpty(params["codehub_id"]),
		"default_branch":     utils.ValueIgnoreEmpty(params["default_branch"]),
		"endpoint_id":        utils.ValueIgnoreEmpty(params["endpoint_id"]),
		"build_params":       buildPipelineActionBuildParams(params["build_params"].([]interface{})),
		"change_request_ids": utils.ValueIgnoreEmpty(params["change_request_ids"].(*schema.Set).List()),
	}

	return result
}

func buildPipelineActionBuildParams(paramsList []interface{}) interface{} {
	if len(paramsList) == 0 {
		return nil
	}

	params := paramsList[0].(map[string]interface{})
	result := map[string]interface{}{
		"build_type":    params["build_type"],
		"target_branch": utils.ValueIgnoreEmpty(params["target_branch"]),
		"tag":           utils.ValueIgnoreEmpty(params["tag"]),
		"event_type":    utils.ValueIgnoreEmpty(params["event_type"]),
	}

	return result
}

func buildPipelineActionVariables(d *schema.ResourceData) interface{} {
	rawVariables := d.Get("variables").(*schema.Set).List()
	if len(rawVariables) == 0 {
		return nil
	}

	variables := make([]map[string]interface{}, 0, len(rawVariables))
	for _, v := range rawVariables {
		variable := v.(map[string]interface{})
		variableMap := map[string]interface{}{
			"name":  variable["name"],
			"value": variable["value"],
		}
		variables = append(variables, variableMap)
	}

	return variables
}

func resourcePipelineActionRead(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return nil
}

func resourcePipelineActionUpdate(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return nil
}

func resourcePipelineActionDelete(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	errorMsg := "Deleting pipeline action resource is not supported. The resource is only removed from the state."
	return diag.Diagnostics{
		diag.Diagnostic{
			Severity: diag.Warning,
			Summary:  errorMsg,
		},
	}
}
