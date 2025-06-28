package codeartsbuild

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

var taskNonUpdatableParams = []string{
	"project_id",
}

// @API CodeArtsBuild POST /v1/job/create
// @API CodeArtsBuild GET /v1/job/{job_id}/config
// @API CodeArtsBuild POST /v1/job/update
// @API CodeArtsBuild DELETE /v1/job/{job_id}/delete
func ResourceCodeArtsBuildTask() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceBuildTaskCreate,
		ReadContext:   resourceBuildTaskRead,
		UpdateContext: resourceBuildTaskUpdate,
		DeleteContext: resourceBuildTaskDelete,

		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		CustomizeDiff: config.FlexibleForceNew(taskNonUpdatableParams),

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
				Description: `Specifies the name of the build task.`,
			},
			"arch": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `Specifies the architecture of the build machine.`,
			},
			"steps": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: `Specifies the build execution steps.`,
				Elem:        resourceSchemeTaskSteps(),
			},
			"auto_update_sub_module": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specifies whether to automatically update submodules.`,
			},
			"flavor": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specifies the specification of the execution machine.`,
			},
			"host_type": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `Specifies the host type.`,
			},
			"build_config_type": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specifies the build task configuration type.`,
			},
			"build_if_code_updated": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: `Specifies whether to enable the code commit trigger build switch.`,
			},
			"group_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specifies the task group ID.`,
			},
			"parameters": {
				Type:        schema.TypeSet,
				Optional:    true,
				Description: `Specifies the build execution parameter list.`,
				Elem:        resourceSchemeTaskParameters(),
			},
			"scms": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: `Specifies the build execution SCM.`,
				Elem:        resourceSchemeTaskScm(),
			},
			"triggers": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: `Specifies the collection of timed task triggers.`,
				Elem:        resourceSchemeTaskTriggers(),
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

func resourceSchemeTaskSteps() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"module_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `Specifies the build step module ID.`,
			},
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `Specifies the build step name.`,
			},
			// `properties` is Map<String, Object> in API actually
			"properties": {
				Type:        schema.TypeMap,
				Optional:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: `Specifies the build step properties.`,
			},
			"version": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specifies the build step version.`,
			},
			"enable": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: `Specifies whether to enable the step.`,
			},
			// add `properties_all` as computed attribute for the API will return extra properties which are not in input
			"properties_all": {
				Type:        schema.TypeMap,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: `Specifies the build step properties.`,
			},
		},
	}
}

func resourceSchemeTaskParameters() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specifies the parameter definition name.`,
			},
			"params": {
				Type:        schema.TypeSet,
				Optional:    true,
				Description: `Specifies the build execution sub-parameters.`,
				Elem:        resourceSchemeTaskParametersParams(),
			},
		},
	}
}

func resourceSchemeTaskParametersParams() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specifies the parameter field name.`,
			},
			"value": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specifies the parameter field value.`,
			},
			"limits": {
				Type:        schema.TypeSet,
				Optional:    true,
				Description: `Specifies the enumeration parameter restrictions.`,
				Elem:        resourceSchemeTaskParametersParamsLimits(),
			},
		},
	}
}

func resourceSchemeTaskParametersParamsLimits() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"disable": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specifies whether it is effective.`,
			},
			"display_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specifies the displayed name of the parameter.`,
			},
			"name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specifies the parameter name.`,
			},
		},
	}
}

func resourceSchemeTaskScm() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"url": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `Specifies the repository URL.`,
			},
			"repo_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `Specifies the repository ID.`,
			},
			"web_url": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `Specifies the web URL of the repository.`,
			},
			"scm_type": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `Specifies the source code management type.`,
			},
			"branch": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specifies the branch name.`,
			},
			"is_auto_build": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: `Specifies whether to automatically build.`,
			},
			"enable_git_lfs": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: `Specifies whether to enable Git LFS.`,
			},
			"build_type": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `Specifies the build type.`,
			},
			"depth": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `Specifies the depth.`,
			},
			"endpoint_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specifies the endpoint ID.`,
			},
			"source": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specifies the source type.`,
			},
			"group_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `Specifies the group name.`,
			},
			"repo_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `Specifies the repository name.`,
			},
		},
	}
}

func resourceSchemeTaskTriggers() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `Specifies the trigger type.`,
			},
			"parameters": {
				Type:        schema.TypeSet,
				Required:    true,
				Description: `Specifies the custom parameters.`,
				Elem:        resourceSchemeTaskTriggersParameters(),
			},
		},
	}
}

func resourceSchemeTaskTriggersParameters() *schema.Resource {
	return &schema.Resource{
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
	}
}

func resourceBuildTaskCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	client, err := cfg.NewServiceClient("codearts_build", cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating CodeArts Build client: %s", err)
	}

	httpUrl := "v1/job/create"
	createPath := client.Endpoint + httpUrl
	createOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		JSONBody:         utils.RemoveNil(buildCreateOrUpdateBuildTaskBodyParams(d)),
	}

	createResp, err := client.Request("POST", createPath, &createOpt)
	if err != nil {
		return diag.Errorf("error creating CodeArts Build task: %s", err)
	}
	createRespBody, err := utils.FlattenResponse(createResp)
	if err != nil {
		return diag.FromErr(err)
	}

	if err := checkResponseError(createRespBody); err != nil {
		return diag.Errorf("error creating CodeArts Build task: %s", err)
	}

	id := utils.PathSearch("result.job_id", createRespBody, "").(string)
	if id == "" {
		return diag.Errorf("unable to find the CodeArts Build task ID from the API response")
	}

	d.SetId(id)

	return resourceBuildTaskRead(ctx, d, meta)
}

func buildCreateOrUpdateBuildTaskBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"job_id":                 utils.ValueIgnoreEmpty(d.Id()),
		"arch":                   d.Get("arch"),
		"project_id":             d.Get("project_id"),
		"job_name":               d.Get("name"),
		"auto_update_sub_module": utils.ValueIgnoreEmpty(d.Get("auto_update_sub_module")),
		"flavor":                 utils.ValueIgnoreEmpty(d.Get("flavor")),
		"host_type":              utils.ValueIgnoreEmpty(d.Get("host_type")),
		"build_config_type":      utils.ValueIgnoreEmpty(d.Get("build_config_type")),
		"build_if_code_updated":  utils.ValueIgnoreEmpty(d.Get("build_if_code_updated")),
		"group_id":               utils.ValueIgnoreEmpty(d.Get("group_id")),
		"parameters":             buildBuildTaskParameters(d),
		"scms":                   buildBuildTaskScms(d),
		"steps":                  buildBuildTaskSteps(d),
		"triggers":               buildBuildTaskTriggers(d),
	}

	return bodyParams
}

func buildBuildTaskParameters(d *schema.ResourceData) interface{} {
	rawParameters := d.Get("parameters").(*schema.Set).List()
	if len(rawParameters) == 0 {
		return nil
	}

	parameters := make([]map[string]interface{}, 0, len(rawParameters))
	for _, p := range rawParameters {
		if parameter, ok := p.(map[string]interface{}); ok {
			parameterMap := map[string]interface{}{
				"name":   utils.ValueIgnoreEmpty(parameter["name"]),
				"params": buildBuildTaskParametersParams(parameter["params"].(*schema.Set).List()),
			}
			parameters = append(parameters, parameterMap)
		}
	}

	return parameters
}

func buildBuildTaskParametersParams(paramsList []interface{}) interface{} {
	if len(paramsList) == 0 {
		return nil
	}

	params := make([]map[string]interface{}, 0, len(paramsList))
	for _, p := range paramsList {
		if param, ok := p.(map[string]interface{}); ok {
			paramMap := map[string]interface{}{
				"name":   utils.ValueIgnoreEmpty(param["name"]),
				"value":  utils.ValueIgnoreEmpty(param["value"]),
				"limits": buildBuildTaskParametersLimits(param["limits"].(*schema.Set).List()),
			}
			params = append(params, paramMap)
		}
	}

	return params
}

func buildBuildTaskParametersLimits(limitsList []interface{}) interface{} {
	if len(limitsList) == 0 {
		return nil
	}

	limits := make([]map[string]interface{}, 0, len(limitsList))
	for _, l := range limitsList {
		if limit, ok := l.(map[string]interface{}); ok {
			limitMap := map[string]interface{}{
				"disable":      utils.ValueIgnoreEmpty(limit["disable"]),
				"display_name": utils.ValueIgnoreEmpty(limit["display_name"]),
				"name":         utils.ValueIgnoreEmpty(limit["name"]),
			}
			limits = append(limits, limitMap)
		}
	}

	return limits
}

func buildBuildTaskScms(d *schema.ResourceData) interface{} {
	rawScms := d.Get("scms").([]interface{})
	if len(rawScms) == 0 {
		return nil
	}

	scms := make([]map[string]interface{}, 0, len(rawScms))
	for _, s := range rawScms {
		if scm, ok := s.(map[string]interface{}); ok {
			scmMap := map[string]interface{}{
				"url":            scm["url"],
				"repo_id":        scm["repo_id"],
				"web_url":        scm["web_url"],
				"scm_type":       scm["scm_type"],
				"branch":         utils.ValueIgnoreEmpty(scm["branch"]),
				"is_auto_build":  utils.ValueIgnoreEmpty(scm["is_auto_build"]),
				"enable_git_lfs": utils.ValueIgnoreEmpty(scm["enable_git_lfs"]),
				"build_type":     utils.ValueIgnoreEmpty(scm["build_type"]),
				"depth":          utils.ValueIgnoreEmpty(scm["depth"]),
				"endpoint_id":    utils.ValueIgnoreEmpty(scm["end_point_id"]),
				"source":         utils.ValueIgnoreEmpty(scm["source"]),
				"group_name":     utils.ValueIgnoreEmpty(scm["group_name"]),
				"repo_name":      utils.ValueIgnoreEmpty(scm["repo_name"]),
			}
			scms = append(scms, scmMap)
		}
	}

	return scms
}

func buildBuildTaskSteps(d *schema.ResourceData) interface{} {
	rawSteps := d.Get("steps").([]interface{})
	if len(rawSteps) == 0 {
		// `steps` is required in API, but can be empty
		return make([]interface{}, 0)
	}

	steps := make([]map[string]interface{}, 0, len(rawSteps))
	for _, s := range rawSteps {
		if step, ok := s.(map[string]interface{}); ok {
			stepMap := map[string]interface{}{
				"module_id":  step["module_id"],
				"name":       step["name"],
				"properties": utils.ValueIgnoreEmpty(buildBuildTaskStepsProperties(step["properties"].(map[string]interface{}))),
				"version":    utils.ValueIgnoreEmpty(step["version"]),
				"enable":     utils.ValueIgnoreEmpty(step["enable"]),
			}
			steps = append(steps, stepMap)
		}
	}

	return steps
}

func buildBuildTaskStepsProperties(m map[string]interface{}) map[string]interface{} {
	rst := make(map[string]interface{}, len(m))
	for k, v := range m {
		rst[k] = parseJson(v.(string))
	}
	return rst
}

func buildBuildTaskTriggers(d *schema.ResourceData) interface{} {
	rawTriggers := d.Get("triggers").([]interface{})
	if len(rawTriggers) == 0 {
		return nil
	}

	triggers := make([]map[string]interface{}, 0, len(rawTriggers))
	for _, t := range rawTriggers {
		if trigger, ok := t.(map[string]interface{}); ok {
			triggerMap := map[string]interface{}{
				"parameters": buildBuildTaskTriggerParameters(trigger["parameters"].(*schema.Set).List()),
				"name":       trigger["name"],
			}
			triggers = append(triggers, triggerMap)
		}
	}

	return triggers
}

func buildBuildTaskTriggerParameters(paramsList []interface{}) interface{} {
	if len(paramsList) == 0 {
		return nil
	}

	params := make([]map[string]interface{}, 0, len(paramsList))
	for _, p := range paramsList {
		if param, ok := p.(map[string]interface{}); ok {
			paramMap := map[string]interface{}{
				"name":  param["name"],
				"value": param["value"],
			}
			params = append(params, paramMap)
		}
	}

	return params
}

func GetBuildTask(client *golangsdk.ServiceClient, id string) (interface{}, error) {
	httpUrl := "v1/job/{job_id}/config?get_all_params=true"
	getPath := client.Endpoint + httpUrl
	getPath = strings.ReplaceAll(getPath, "{job_id}", id)
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

	if err := checkResponseError(getRespBody); err != nil {
		return nil, err
	}

	return getRespBody, nil
}

func resourceBuildTaskRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	client, err := cfg.NewServiceClient("codearts_build", region)
	if err != nil {
		return diag.Errorf("error creating CodeArts Build client: %s", err)
	}

	getRespBody, err := GetBuildTask(client, d.Id())
	if err != nil {
		return common.CheckDeletedDiag(d, common.ConvertUndefinedErrInto404Err(err, 422, "error_code", buildTaskNotFoundErr),
			"error retrieving CodeArts Build task")
	}

	task := utils.PathSearch("result", getRespBody, nil)
	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("project_id", utils.PathSearch("project_id", task, nil)),
		d.Set("name", utils.PathSearch("job_name", task, nil)),
		d.Set("arch", utils.PathSearch("arch", task, nil)),
		d.Set("auto_update_sub_module", utils.PathSearch("auto_update_sub_module", task, nil)),
		d.Set("flavor", utils.PathSearch("flavor", task, nil)),
		d.Set("host_type", utils.PathSearch("host_type", task, nil)),
		d.Set("build_config_type", utils.PathSearch("build_config_type", task, nil)),
		d.Set("build_if_code_updated", utils.PathSearch("build_if_code_updated", task, nil)),
		d.Set("group_id", utils.PathSearch("group_id", task, nil)),
		d.Set("parameters", flattenBuildTaskParameters(task)),
		d.Set("scms", flattenBuildTaskScms(task)),
		d.Set("steps", flattenBuildTaskSteps(d, task)),
		d.Set("triggers", flattenBuildTaskTriggers(task)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenBuildTaskParameters(resp interface{}) []interface{} {
	parametersList := utils.PathSearch("parameters", resp, make([]interface{}, 0)).([]interface{})
	if len(parametersList) == 0 {
		return nil
	}

	result := make([]interface{}, 0, len(parametersList))
	for _, p := range parametersList {
		parameter := p.(map[string]interface{})
		parameterMap := map[string]interface{}{
			"name":   utils.PathSearch("name", parameter, nil),
			"params": flattenBuildTaskParametersParams(parameter),
		}
		result = append(result, parameterMap)
	}

	return result
}

func flattenBuildTaskParametersParams(parameter map[string]interface{}) []interface{} {
	paramsList := utils.PathSearch("params", parameter, make([]interface{}, 0)).([]interface{})
	if len(paramsList) == 0 {
		return nil
	}

	result := make([]interface{}, 0, len(paramsList))
	for _, p := range paramsList {
		param := p.(map[string]interface{})
		paramMap := map[string]interface{}{
			"name":   utils.PathSearch("name", param, nil),
			"value":  utils.PathSearch("value", param, nil),
			"limits": flattenBuildTaskParametersLimits(param),
		}
		result = append(result, paramMap)
	}

	return result
}

func flattenBuildTaskParametersLimits(param map[string]interface{}) []interface{} {
	limitsList := utils.PathSearch("limits", param, make([]interface{}, 0)).([]interface{})
	if len(limitsList) == 0 {
		return nil
	}

	result := make([]interface{}, 0, len(limitsList))
	for _, l := range limitsList {
		limit := l.(map[string]interface{})
		limitMap := map[string]interface{}{
			"disable":      utils.PathSearch("disable", limit, nil),
			"display_name": utils.PathSearch("display_name", limit, nil),
			"name":         utils.PathSearch("name", limit, nil),
		}
		result = append(result, limitMap)
	}

	return result
}

func flattenBuildTaskScms(resp interface{}) []interface{} {
	scmsList := utils.PathSearch("scms", resp, make([]interface{}, 0)).([]interface{})
	if len(scmsList) == 0 {
		return nil
	}

	result := make([]interface{}, 0, len(scmsList))
	for _, s := range scmsList {
		scm := s.(map[string]interface{})
		scmMap := map[string]interface{}{
			"branch":         utils.PathSearch("branch", scm, nil),
			"url":            utils.PathSearch("url", scm, nil),
			"repo_id":        utils.PathSearch("repo_id", scm, nil),
			"web_url":        utils.PathSearch("web_url", scm, nil),
			"scm_type":       utils.PathSearch("scm_type", scm, nil),
			"is_auto_build":  utils.PathSearch("is_auto_build", scm, nil),
			"build_type":     utils.PathSearch("build_type", scm, nil),
			"depth":          utils.PathSearch("depth", scm, nil),
			"enable_git_lfs": utils.PathSearch("enable_git_lfs", scm, nil),
			"endpoint_id":    utils.PathSearch("end_point_id", scm, nil),
			"source":         utils.PathSearch("source", scm, nil),
			"group_name":     utils.PathSearch("group_name", scm, nil),
			"repo_name":      utils.PathSearch("repo_name", scm, nil),
		}
		result = append(result, scmMap)
	}

	return result
}

func flattenBuildTaskSteps(d *schema.ResourceData, resp interface{}) []interface{} {
	stepsList := utils.PathSearch("steps", resp, make([]interface{}, 0)).([]interface{})
	if len(stepsList) == 0 {
		return nil
	}

	result := make([]interface{}, 0, len(stepsList))
	for i, s := range stepsList {
		step := s.(map[string]interface{})
		propertiesPath := fmt.Sprintf("steps.%d.properties", i)
		properties := d.Get(propertiesPath).(map[string]interface{})
		stepMap := map[string]interface{}{
			"properties": flattenBuildTaskStepsProperties(properties,
				utils.PathSearch("properties", step, make(map[string]interface{}, 0)).(map[string]interface{})),
			"module_id": utils.PathSearch("module_id", step, nil),
			"name":      utils.PathSearch("name", step, nil),
			"version":   utils.PathSearch("version", step, nil),
			"enable":    utils.PathSearch("enable", step, nil),
			"properties_all": flattenBuildTaskStepsPropertiesAll(
				utils.PathSearch("properties", step, make(map[string]interface{}, 0)).(map[string]interface{})),
		}
		result = append(result, stepMap)
	}

	return result
}

func flattenBuildTaskStepsProperties(o, m map[string]interface{}) map[string]interface{} {
	rst := make(map[string]interface{}, len(o))
	for k, v := range m {
		if _, ok := o[k]; ok {
			rst[k] = encodeJson(v)
		}
	}
	return rst
}

func flattenBuildTaskStepsPropertiesAll(m map[string]interface{}) map[string]interface{} {
	rst := make(map[string]interface{}, len(m))
	for k, v := range m {
		rst[k] = encodeJson(v)
	}
	return rst
}

func flattenBuildTaskTriggers(resp interface{}) []interface{} {
	triggersList := utils.PathSearch("triggers", resp, make([]interface{}, 0)).([]interface{})
	if len(triggersList) == 0 {
		return nil
	}

	result := make([]interface{}, 0, len(triggersList))
	for _, t := range triggersList {
		trigger := t.(map[string]interface{})
		triggerMap := map[string]interface{}{
			"name":       utils.PathSearch("name", trigger, nil),
			"parameters": flattenBuildTaskTriggerParameters(trigger),
		}
		result = append(result, triggerMap)
	}

	return result
}

func flattenBuildTaskTriggerParameters(trigger map[string]interface{}) []interface{} {
	parametersList := utils.PathSearch("parameters", trigger, make([]interface{}, 0)).([]interface{})
	if len(parametersList) == 0 {
		return nil
	}

	result := make([]interface{}, 0, len(parametersList))
	for _, p := range parametersList {
		param := p.(map[string]interface{})
		paramMap := map[string]interface{}{
			"name":  utils.PathSearch("name", param, nil),
			"value": utils.PathSearch("value", param, nil),
		}
		result = append(result, paramMap)
	}

	return result
}

func resourceBuildTaskUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	client, err := cfg.NewServiceClient("codearts_build", region)
	if err != nil {
		return diag.Errorf("error creating CodeArts Build client: %s", err)
	}

	httpUrl := "v1/job/update"
	updatePath := client.Endpoint + httpUrl
	updateOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		JSONBody:         utils.RemoveNil(buildCreateOrUpdateBuildTaskBodyParams(d)),
	}

	updateResp, err := client.Request("POST", updatePath, &updateOpt)
	if err != nil {
		return diag.Errorf("error updating CodeArts Build task: %s", err)
	}
	updateRespBody, err := utils.FlattenResponse(updateResp)
	if err != nil {
		return diag.FromErr(err)
	}

	if err := checkResponseError(updateRespBody); err != nil {
		return diag.Errorf("error updating CodeArts Build task: %s", err)
	}

	return resourceBuildTaskRead(ctx, d, meta)
}

func resourceBuildTaskDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	client, err := cfg.NewServiceClient("codearts_build", region)
	if err != nil {
		return diag.Errorf("error creating CodeArts Build client: %s", err)
	}

	httpUrl := "v1/job/{job_id}/delete"
	deletePath := client.Endpoint + httpUrl
	deletePath = strings.ReplaceAll(deletePath, "{job_id}", d.Id())
	deleteOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	_, err = client.Request("DELETE", deletePath, &deleteOpt)
	if err != nil {
		return common.CheckDeletedDiag(d, common.ConvertUndefinedErrInto404Err(err, 422, "error_code", buildTaskNotFoundErr),
			"error deleting CodeArts Build task")
	}

	return nil
}
