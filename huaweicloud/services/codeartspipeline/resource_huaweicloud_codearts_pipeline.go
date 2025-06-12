package codeartspipeline

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

const (
	banHttpUrl   = "v5/{project_id}/api/pipelines/{pipeline_id}/ban"
	unbanHttpUrl = "v5/{project_id}/api/pipelines/{pipeline_id}/unban"
)

var pipelineNonUpdatableParams = []string{
	"project_id", "component_id",
}

// @API CodeArtsPipeline POST /v5/{project_id}/api/pipelines
// @API CodeArtsPipeline GET /v5/{project_id}/api/pipelines/{pipeline_id}
// @API CodeArtsPipeline PUT /v5/{project_id}/api/pipelines/{pipeline_id}
// @API CodeArtsPipeline DELETE /v5/{project_id}/api/pipelines/{pipeline_id}
// @API CodeArtsPipeline PUT /v5/{project_id}/api/pipelines/{pipeline_id}/unban
// @API CodeArtsPipeline PUT /v5/{project_id}/api/pipelines/{pipeline_id}/ban
func ResourceCodeArtsPipeline() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourcePipelineCreate,
		ReadContext:   resourcePipelineRead,
		UpdateContext: resourcePipelineUpdate,
		DeleteContext: resourcePipelineDelete,

		Importer: &schema.ResourceImporter{
			StateContext: resourceImportStateFuncWithProjectIdAndId,
		},

		CustomizeDiff: config.FlexibleForceNew(pipelineNonUpdatableParams),

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
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `Specifies the pipeline name.`,
			},
			"definition": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `Specifies the pipeline definition JSON.`,
			},
			"is_publish": {
				Type:        schema.TypeBool,
				Required:    true,
				Description: `Specifies whether it is a change-triggered pipeline.`,
			},
			"manifest_version": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `Specifies the pipeline structure definition version.`,
			},
			"component_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specifies the microservice ID.`,
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specifies the pipeline description.`,
			},
			"group_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specifies the pipeline group ID.`,
			},
			"project_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `Specifies the project name.`,
			},
			"banned": {
				Type:        schema.TypeBool,
				Optional:    true,
				Computed:    true,
				Description: `Specifies whether the pipeline is banned.`,
			},
			"sources": {
				Type:        schema.TypeSet,
				Optional:    true,
				Description: `Specifies the pipeline source information.`,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"type": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: `Specifies the pipeline source type.`,
						},
						"params": {
							Type:        schema.TypeList,
							Optional:    true,
							MaxItems:    1,
							Description: `Specifies the pipeline source parameters.`,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"git_type": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: `Specifies the code repository type.`,
									},
									"codehub_id": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: `Specifies the CodeArts Repo code repository ID.`,
									},
									"endpoint_id": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: `Specifies the code source endpoint ID.`,
									},
									"default_branch": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: `Specifies the default branch.`,
									},
									"git_url": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: `Specifies the HTTPS address of the Git repository.`,
									},
									"ssh_git_url": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: `Specifies the SSH Git address,`,
									},
									"web_url": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: `Specifies the web page URL.`,
									},
									"repo_name": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: `Specifies the pipeline source name.`,
									},
									"alias": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: `Specifies the code repository alias.`,
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
				Description: `Specifies the custom variables.`,
				Elem:        resourceSchemePipelineVariables(),
			},
			"schedules": {
				Type:        schema.TypeSet,
				Optional:    true,
				Description: `Specifies the pipeline schedule settings.`,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"type": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: `Specifies the schedule job type.`,
						},
						"name": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: `Specifies the schedule job name.`,
						},
						"enable": {
							Type:        schema.TypeBool,
							Optional:    true,
							Description: `Specifies whether to enable the schedule job.`,
						},
						"days_of_week": {
							Type:        schema.TypeSet,
							Optional:    true,
							Elem:        &schema.Schema{Type: schema.TypeInt},
							Description: `Specifies the execution day in a week.`,
						},
						"time_zone": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: `Specifies the time zone.`,
						},
						"start_time": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: `Specifies the start time.`,
						},
						"end_time": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: `Specifies the end time.`,
						},
						"interval_time": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: `Specifies the interval time.`,
						},
						"interval_unit": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: `Specifies the interval unit.`,
						},
						"uuid": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `Indicates the ID of a scheduled task.`,
						},
					},
				},
			},
			"triggers": {
				Type:        schema.TypeSet,
				Optional:    true,
				Description: `Specifies the pipeline trigger settings.`,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"git_url": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: `Specifies the Git URL.`,
						},
						"git_type": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: `Specifies the Git repository type.`,
						},
						"is_auto_commit": {
							Type:        schema.TypeBool,
							Optional:    true,
							Description: `Specifies whether to automatically commit code.`,
						},
						"events": {
							Type:        schema.TypeSet,
							Optional:    true,
							Description: `Specifies the trigger event list.`,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"type": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: `Specifies the event type.`,
									},
									"enable": {
										Type:        schema.TypeBool,
										Optional:    true,
										Description: `Specifies whether it is available.`,
									},
								},
							},
						},
						"repo_id": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: `Specifies the repository ID.`,
						},
						"endpoint_id": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: `Specifies the code source endpoint ID.`,
						},
						"callback_url": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: `Specifies the callback URL.`,
						},
						"security_token": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: `Specifies the User token.`,
						},
						"hook_id": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: `Indicates the callback ID.`,
						},
					},
				},
			},
			"concurrency_control": {
				Type:        schema.TypeList,
				Optional:    true,
				MaxItems:    1,
				Description: `Specifies the pipeline concurrency control information.`,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"concurrency_number": {
							Type:        schema.TypeInt,
							Optional:    true,
							Description: `Specifies the number of concurrent instances.`,
						},
						"exceed_action": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: `Specifies the policy when the threshold is exceeded.`,
						},
						"enable": {
							Type:        schema.TypeBool,
							Optional:    true,
							Description: `Specifies whether to enable the strategy.`,
						},
					},
				},
			},
			"enable_force_new": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"true", "false"}, false),
				Description:  utils.SchemaDesc("", utils.SchemaDescInput{Internal: true}),
			},
			"creator_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Indicates the creator ID.`,
			},
			"creator_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Indicates the creator name.`,
			},
			"updater_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Indicates the last updater ID.`,
			},
			"create_time": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: `Indicates the creation time.`,
			},
			"update_time": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: `Indicates the last update time.`,
			},
			"is_collect": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: `Indicates whether the current user has collected it.`,
			},
		},
	}
}

func resourceSchemePipelineVariables() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specifies the custom variable name.`,
			},
			"sequence": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: `Specifies the parameter sequence, starting from 1.`,
			},
			"type": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specifies the custom parameter type.`,
			},
			"value": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specifies the custom parameter default value.`,
			},
			"is_secret": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: `Specifies whether it is a private parameter.`,
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specifies the parameter description.`,
			},
			"is_runtime": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: `Specifies whether to set parameters at runtime.`,
			},
			"is_reset": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: `Specifies whether to reset.`,
			},
			"latest_value": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specifies the last parameter value.`,
			},
			"runtime_value": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specifies the value passed in at runtime.`,
			},
		},
	}
}

func resourcePipelineCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	client, err := cfg.NewServiceClient("codearts_pipeline", region)
	if err != nil {
		return diag.Errorf("error creating CodeArts Pipeline client: %s", err)
	}

	projectId := d.Get("project_id").(string)

	httpUrl := "v5/{project_id}/api/pipelines"
	createPath := client.Endpoint + httpUrl
	createPath = strings.ReplaceAll(createPath, "{project_id}", projectId)
	if v, ok := d.GetOk("component_id"); ok {
		createPath += fmt.Sprintf("component_id=%v", v)
	}

	createOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		JSONBody:         utils.RemoveNil(buildCreateOrUpdatePipelineBodyParams(d)),
	}

	createResp, err := client.Request("POST", createPath, &createOpt)
	if err != nil {
		return diag.Errorf("error creating CodeArts Pipeline: %s", err)
	}
	createRespBody, err := utils.FlattenResponse(createResp)
	if err != nil {
		return diag.FromErr(err)
	}

	if err := checkResponseError(createRespBody, ""); err != nil {
		return diag.Errorf("error creating CodeArts Pipeline: %s", err)
	}

	id := utils.PathSearch("pipeline_id", createRespBody, "").(string)
	if id == "" {
		return diag.Errorf("unable to find the CodeArts Pipeline ID from the API response")
	}

	d.SetId(id)

	if d.Get("banned").(bool) {
		if err := updatePipelineBanned(client, banHttpUrl, projectId, d.Id()); err != nil {
			return diag.Errorf("error banning pipeline: %s", err)
		}
	}

	return resourcePipelineRead(ctx, d, meta)
}

func buildCreateOrUpdatePipelineBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"name":                d.Get("name"),
		"description":         d.Get("description"),
		"definition":          d.Get("definition"),
		"is_publish":          d.Get("is_publish"),
		"project_name":        utils.ValueIgnoreEmpty(d.Get("project_name")),
		"component_id":        utils.ValueIgnoreEmpty(d.Get("component_id")),
		"group_id":            utils.ValueIgnoreEmpty(d.Get("group_id")),
		"manifest_version":    utils.ValueIgnoreEmpty(d.Get("manifest_version")),
		"sources":             buildPipelineSources(d),
		"variables":           buildPipelineTemplateVariables(d),
		"schedules":           buildPipelineSchedules(d),
		"triggers":            buildPipelineTriggers(d),
		"concurrency_control": buildPipelineConcurrencyControlParams(d),
	}

	return bodyParams
}

func buildPipelineSources(d *schema.ResourceData) interface{} {
	rawSources := d.Get("sources").(*schema.Set).List()
	if len(rawSources) == 0 {
		return nil
	}

	sources := make([]map[string]interface{}, 0, len(rawSources))
	for _, s := range rawSources {
		if source, ok := s.(map[string]interface{}); ok {
			sourceMap := map[string]interface{}{
				"type":   utils.ValueIgnoreEmpty(source["type"]),
				"params": buildPipelineSourcesParams(source["params"].([]interface{})),
			}
			sources = append(sources, sourceMap)
		}
	}

	return sources
}

func buildPipelineSourcesParams(paramsList []interface{}) interface{} {
	if len(paramsList) == 0 {
		return nil
	}

	if params, ok := paramsList[0].(map[string]interface{}); ok {
		return map[string]interface{}{
			"git_type":       utils.ValueIgnoreEmpty(params["git_type"]),
			"codehub_id":     utils.ValueIgnoreEmpty(params["codehub_id"]),
			"endpoint_id":    utils.ValueIgnoreEmpty(params["endpoint_id"]),
			"default_branch": utils.ValueIgnoreEmpty(params["default_branch"]),
			"git_url":        utils.ValueIgnoreEmpty(params["git_url"]),
			"ssh_git_url":    utils.ValueIgnoreEmpty(params["ssh_git_url"]),
			"web_url":        utils.ValueIgnoreEmpty(params["web_url"]),
			"repo_name":      utils.ValueIgnoreEmpty(params["repo_name"]),
			"alias":          utils.ValueIgnoreEmpty(params["alias"]),
		}
	}

	return nil
}

func buildPipelineSchedules(d *schema.ResourceData) interface{} {
	rawSchedules := d.Get("schedules").(*schema.Set).List()
	if len(rawSchedules) == 0 {
		return nil
	}

	rst := make([]map[string]interface{}, 0, len(rawSchedules))
	for _, s := range rawSchedules {
		if schedule, ok := s.(map[string]interface{}); ok {
			scheduleMap := map[string]interface{}{
				"type":          utils.ValueIgnoreEmpty(schedule["type"]),
				"name":          utils.ValueIgnoreEmpty(schedule["name"]),
				"enable":        utils.ValueIgnoreEmpty(schedule["enable"]),
				"days_of_week":  utils.ValueIgnoreEmpty(schedule["days_of_week"].(*schema.Set).List()),
				"time_zone":     utils.ValueIgnoreEmpty(schedule["time_zone"]),
				"start_time":    utils.ValueIgnoreEmpty(schedule["start_time"]),
				"end_time":      utils.ValueIgnoreEmpty(schedule["end_time"]),
				"interval_time": utils.ValueIgnoreEmpty(schedule["interval_time"]),
				"interval_unit": utils.ValueIgnoreEmpty(schedule["interval_unit"]),
			}
			rst = append(rst, scheduleMap)
		}
	}

	return rst
}

func buildPipelineTriggers(d *schema.ResourceData) interface{} {
	rawTriggers := d.Get("triggers").(*schema.Set).List()
	if len(rawTriggers) == 0 {
		return nil
	}

	triggers := make([]map[string]interface{}, 0, len(rawTriggers))
	for _, t := range rawTriggers {
		if trigger, ok := t.(map[string]interface{}); ok {
			triggerMap := map[string]interface{}{
				"git_url":        utils.ValueIgnoreEmpty(trigger["git_url"]),
				"git_type":       utils.ValueIgnoreEmpty(trigger["git_type"]),
				"is_auto_commit": utils.ValueIgnoreEmpty(trigger["is_auto_commit"]),
				"repo_id":        utils.ValueIgnoreEmpty(trigger["repo_id"]),
				"endpoint_id":    utils.ValueIgnoreEmpty(trigger["endpoint_id"]),
				"callback_url":   utils.ValueIgnoreEmpty(trigger["callback_url"]),
				"security_token": utils.ValueIgnoreEmpty(trigger["security_token"]),
				"events":         buildPipelineTriggersEvents(trigger["events"].(*schema.Set).List()),
			}
			triggers = append(triggers, triggerMap)
		}
	}

	return triggers
}

func buildPipelineTriggersEvents(paramsList []interface{}) interface{} {
	if len(paramsList) == 0 {
		return nil
	}
	rst := make([]map[string]interface{}, 0, len(paramsList))
	for _, rawParams := range paramsList {
		if params, ok := rawParams.(map[string]interface{}); ok {
			m := map[string]interface{}{
				"type":   utils.ValueIgnoreEmpty(params["type"]),
				"enable": utils.ValueIgnoreEmpty(params["enable"]),
			}
			rst = append(rst, m)
		}
	}

	return rst
}

func buildPipelineConcurrencyControlParams(d *schema.ResourceData) interface{} {
	paramsList := d.Get("concurrency_control").([]interface{})
	if len(paramsList) == 0 {
		return nil
	}

	if params, ok := paramsList[0].(map[string]interface{}); ok {
		return map[string]interface{}{
			"concurrency_number": utils.ValueIgnoreEmpty(params["concurrency_number"]),
			"exceed_action":      utils.ValueIgnoreEmpty(params["exceed_action"]),
			"enable":             utils.ValueIgnoreEmpty(params["enable"]),
		}
	}

	return nil
}

func updatePipelineBanned(client *golangsdk.ServiceClient, httpUrl, projectId, id string) error {
	unbanPath := client.Endpoint + httpUrl
	unbanPath = strings.ReplaceAll(unbanPath, "{project_id}", projectId)
	unbanPath = strings.ReplaceAll(unbanPath, "{pipeline_id}", id)
	unbanOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	unbanResp, err := client.Request("PUT", unbanPath, &unbanOpt)
	if err != nil {
		return err
	}
	unbanRespBody, err := utils.FlattenResponse(unbanResp)
	if err != nil {
		return err
	}

	return checkResponseError(unbanRespBody, "")
}

func GetPipeline(client *golangsdk.ServiceClient, projectId, id string) (interface{}, error) {
	httpUrl := "v5/{project_id}/api/pipelines/{pipeline_id}"
	getPath := client.Endpoint + httpUrl
	getPath = strings.ReplaceAll(getPath, "{project_id}", projectId)
	getPath = strings.ReplaceAll(getPath, "{pipeline_id}", id)
	getOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	getResp, err := client.Request("GET", getPath, &getOpt)
	if err != nil {
		return nil, err
	}
	getRespBody, err := utils.FlattenResponse(getResp)
	if err != nil {
		return nil, err
	}

	if err := checkResponseError(getRespBody, pipelineNotFoundError); err != nil {
		return nil, err
	}

	return getRespBody, nil
}

func resourcePipelineRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	client, err := cfg.NewServiceClient("codearts_pipeline", region)
	if err != nil {
		return diag.Errorf("error creating CodeArts Pipeline client: %s", err)
	}

	projectId := d.Get("project_id").(string)
	getRespBody, err := GetPipeline(client, projectId, d.Id())
	if err != nil {
		return common.CheckDeletedDiag(d, err, "error retrieving CodeArts Pipeline")
	}

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("name", utils.PathSearch("name", getRespBody, nil)),
		d.Set("description", utils.PathSearch("description", getRespBody, nil)),
		d.Set("manifest_version", utils.PathSearch("manifest_version", getRespBody, nil)),
		d.Set("is_publish", utils.PathSearch("is_publish", getRespBody, nil)),
		d.Set("component_id", utils.PathSearch("component_id", getRespBody, nil)),
		d.Set("group_id", utils.PathSearch("group_id", getRespBody, nil)),
		d.Set("creator_id", utils.PathSearch("creator_id", getRespBody, nil)),
		d.Set("creator_name", utils.PathSearch("creator_name", getRespBody, nil)),
		d.Set("updater_id", utils.PathSearch("updater_id", getRespBody, nil)),
		d.Set("create_time", utils.PathSearch("create_time", getRespBody, nil)),
		d.Set("update_time", utils.PathSearch("update_time", getRespBody, nil)),
		d.Set("is_collect", utils.PathSearch("is_collect", getRespBody, nil)),
		d.Set("region", utils.PathSearch("region", getRespBody, nil)),
		d.Set("sources", flattenPipelineSources(getRespBody)),
		d.Set("variables", flattenPipelineTemplateVariables(getRespBody)),
		d.Set("schedules", flattenPipelineSchedules(getRespBody)),
		d.Set("triggers", flattenPipelineTriggers(getRespBody)),
		d.Set("concurrency_control", flattenPipelineConcurrencyControl(getRespBody)),
		d.Set("banned", utils.PathSearch("banned", getRespBody, nil)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenPipelineSources(resp interface{}) []interface{} {
	sources := utils.PathSearch("sources", resp, nil)
	if sources == nil {
		return nil
	}

	sourcesList, ok := sources.([]interface{})
	if !ok || len(sourcesList) == 0 {
		return nil
	}

	result := make([]interface{}, 0, len(sourcesList))
	for _, s := range sourcesList {
		source := s.(map[string]interface{})
		sourceMap := map[string]interface{}{
			"type":   utils.PathSearch("type", source, nil),
			"params": flattenPipelineSourcesParams(utils.PathSearch("params", source, nil)),
		}
		result = append(result, sourceMap)
	}

	return result
}

func flattenPipelineSourcesParams(params interface{}) interface{} {
	if params == nil {
		return nil
	}

	rst := map[string]interface{}{
		"git_type":       utils.PathSearch("git_type", params, nil),
		"codehub_id":     utils.PathSearch("codehub_id", params, nil),
		"endpoint_id":    utils.PathSearch("endpoint_id", params, nil),
		"default_branch": utils.PathSearch("default_branch", params, nil),
		"git_url":        utils.PathSearch("git_url", params, nil),
		"ssh_git_url":    utils.PathSearch("ssh_git_url", params, nil),
		"web_url":        utils.PathSearch("web_url", params, nil),
		"repo_name":      utils.PathSearch("repo_name", params, nil),
		"alias":          utils.PathSearch("alias", params, nil),
	}

	return []interface{}{rst}
}

func flattenPipelineSchedules(resp interface{}) []interface{} {
	schedules := utils.PathSearch("schedules", resp, nil)
	if schedules == nil {
		return nil
	}

	schedulesList, ok := schedules.([]interface{})
	if !ok || len(schedulesList) == 0 {
		return nil
	}

	result := make([]interface{}, 0, len(schedulesList))
	for _, s := range schedulesList {
		schedule := s.(map[string]interface{})
		scheduleMap := map[string]interface{}{
			"type":          utils.PathSearch("type", schedule, nil),
			"name":          utils.PathSearch("name", schedule, nil),
			"enable":        utils.PathSearch("enable", schedule, nil),
			"days_of_week":  utils.PathSearch("days_of_week", schedule, nil),
			"time_zone":     utils.PathSearch("time_zone", schedule, nil),
			"uuid":          utils.PathSearch("uuid", schedule, nil),
			"start_time":    utils.PathSearch("start_time", schedule, nil),
			"end_time":      utils.PathSearch("end_time", schedule, nil),
			"interval_time": utils.PathSearch("interval_time", schedule, nil),
			"interval_unit": utils.PathSearch("interval_unit", schedule, nil),
		}
		result = append(result, scheduleMap)
	}

	return result
}

func flattenPipelineTriggers(resp interface{}) []interface{} {
	triggers := utils.PathSearch("triggers", resp, nil)
	if triggers == nil {
		return nil
	}

	triggersList, ok := triggers.([]interface{})
	if !ok || len(triggersList) == 0 {
		return nil
	}

	result := make([]interface{}, 0, len(triggersList))
	for _, t := range triggersList {
		trigger := t.(map[string]interface{})
		triggerMap := map[string]interface{}{
			"git_url":        utils.PathSearch("git_url", trigger, nil),
			"git_type":       utils.PathSearch("git_type", trigger, nil),
			"is_auto_commit": utils.PathSearch("is_auto_commit", trigger, nil),
			"repo_id":        utils.PathSearch("repo_id", trigger, nil),
			"endpoint_id":    utils.PathSearch("endpoint_id", trigger, nil),
			"hook_id":        utils.PathSearch("hook_id", trigger, nil),
			"callback_url":   utils.PathSearch("callback_url", trigger, nil),
			"security_token": utils.PathSearch("security_token", trigger, nil),
			"events":         flattenPipelineTriggersEvents(trigger),
		}
		result = append(result, triggerMap)
	}

	return result
}

func flattenPipelineTriggersEvents(resp interface{}) []interface{} {
	events := utils.PathSearch("events", resp, nil)
	if events == nil {
		return nil
	}

	eventsList, ok := events.([]interface{})
	if !ok || len(eventsList) == 0 {
		return nil
	}

	result := make([]interface{}, 0, len(eventsList))
	for _, e := range eventsList {
		event := e.(map[string]interface{})
		eventMap := map[string]interface{}{
			"type":   utils.PathSearch("type", event, nil),
			"enable": utils.PathSearch("enable", event, nil),
		}
		result = append(result, eventMap)
	}

	return result
}

func flattenPipelineConcurrencyControl(resp interface{}) []interface{} {
	concurrencyControl := utils.PathSearch("concurrency_control", resp, nil)
	if concurrencyControl == nil {
		return nil
	}

	rst := map[string]interface{}{
		"concurrency_number": utils.PathSearch("concurrency_number", concurrencyControl, nil),
		"exceed_action":      utils.PathSearch("exceed_action", concurrencyControl, nil),
		"enable":             utils.PathSearch("enable", concurrencyControl, nil),
	}

	return []interface{}{rst}
}

func resourcePipelineUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	client, err := cfg.NewServiceClient("codearts_pipeline", region)
	if err != nil {
		return diag.Errorf("error creating CodeArts Pipeline client: %s", err)
	}

	projectId := d.Get("project_id").(string)

	// unban first, when pipeline is banned, cannot be updated.
	if d.HasChange("banned") {
		if !d.Get("banned").(bool) {
			if err := updatePipelineBanned(client, unbanHttpUrl, projectId, d.Id()); err != nil {
				return diag.Errorf("error unbanning pipeline: %s", err)
			}
		}
	}

	changes := []string{
		"name", "definition", "is_publish", "manifest_version", "description", "group_id", "project_name",
		"sources", "variables", "schedules", "triggers", "concurrency_control",
	}
	if d.HasChanges(changes...) {
		httpUrl := "v5/{project_id}/api/pipelines/{pipeline_id}"
		updatePath := client.Endpoint + httpUrl
		updatePath = strings.ReplaceAll(updatePath, "{project_id}", projectId)
		updatePath = strings.ReplaceAll(updatePath, "{pipeline_id}", d.Id())
		updateOpt := golangsdk.RequestOpts{
			KeepResponseBody: true,
			JSONBody:         utils.RemoveNil(buildCreateOrUpdatePipelineBodyParams(d)),
		}

		updateResp, err := client.Request("PUT", updatePath, &updateOpt)
		if err != nil {
			return diag.Errorf("error updating CodeArts Pipeline: %s", err)
		}
		updateRespBody, err := utils.FlattenResponse(updateResp)
		if err != nil {
			return diag.FromErr(err)
		}

		if err := checkResponseError(updateRespBody, ""); err != nil {
			return diag.Errorf("error updating CodeArts Pipeline: %s", err)
		}
	}

	// ban at the end, when pipeline is banned, cannot be updated.
	if d.HasChange("banned") {
		if d.Get("banned").(bool) {
			if err := updatePipelineBanned(client, banHttpUrl, projectId, d.Id()); err != nil {
				return diag.Errorf("error banning pipeline: %s", err)
			}
		}
	}

	return resourcePipelineRead(ctx, d, meta)
}

func resourcePipelineDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	client, err := cfg.NewServiceClient("codearts_pipeline", region)
	if err != nil {
		return diag.Errorf("error creating CodeArts Pipeline client: %s", err)
	}

	projectId := d.Get("project_id").(string)
	httpUrl := "v5/{project_id}/api/pipelines/{pipeline_id}"
	deletePath := client.Endpoint + httpUrl
	deletePath = strings.ReplaceAll(deletePath, "{project_id}", projectId)
	deletePath = strings.ReplaceAll(deletePath, "{pipeline_id}", d.Id())
	deleteOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	deleteResp, err := client.Request("DELETE", deletePath, &deleteOpt)
	if err != nil {
		return common.CheckDeletedDiag(d, err, "error deleting CodeArts Pipeline")
	}
	deleteRespBody, err := utils.FlattenResponse(deleteResp)
	if err != nil {
		return diag.FromErr(err)
	}

	if err := checkResponseError(deleteRespBody, projectNotFoundError2, pipelineNotFoundError); err != nil {
		return common.CheckDeletedDiag(d, err, "error deleting CodeArts Pipeline")
	}

	return nil
}
