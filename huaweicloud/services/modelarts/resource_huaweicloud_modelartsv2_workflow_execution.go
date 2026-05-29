package modelarts

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

var (
	v2WorkflowExecutionNonUpdatableParams = []string{
		"name",
		"description",
		"workspace_id",
		"workflow_id",
		"workflow_name",
		"scene_id",
		"scene_name",
		"policies",
		"policies.*.use_cache",
	}
	v2WorkflowExecutionNotFoundErrCodes = []string{
		"ModelArts.7511", // The workflow execution does not exist.
		"ModelArts.7512", // The workflow does not exist.
	}
)

// @API ModelArts POST /v2/{project_id}/workflows/{workflow_id}/executions
// @API ModelArts GET /v2/{project_id}/workflows/{workflow_id}/executions/{execution_id}
// @API ModelArts PUT /v2/{project_id}/workflows/{workflow_id}/executions/{execution_id}
// @API ModelArts DELETE /v2/{project_id}/workflows/{workflow_id}/executions/{execution_id}
func ResourceV2WorkflowExecution() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceV2WorkflowExecutionCreate,
		ReadContext:   resourceV2WorkflowExecutionRead,
		UpdateContext: resourceV2WorkflowExecutionUpdate,
		DeleteContext: resourceV2WorkflowExecutionDelete,

		Importer: &schema.ResourceImporter{
			StateContext: resourceV2WorkflowExecutionImportState,
		},

		CustomizeDiff: config.FlexibleForceNew(v2WorkflowExecutionNonUpdatableParams),

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: `The region where the workflow execution is located.`,
			},

			// Optional parameters.
			"name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The name of the workflow execution.`,
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The description of the workflow execution.`,
			},
			"workspace_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The workspace ID to which the workflow execution belongs.`,
			},
			"workflow_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The workflow ID.`,
			},
			"workflow_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The workflow name.`,
			},
			"scene_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The custom scene ID.`,
			},
			"scene_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The custom scene name.`,
			},
			"labels": {
				Type:        schema.TypeList,
				Optional:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: `The labels of the workflow execution.`,
			},
			"data_requirements": {
				Type:        schema.TypeList,
				Optional:    true,
				Computed:    true,
				Elem:        v2WorkflowExecutionDataRequirementsSchema(),
				Description: `The data requirements of the workflow execution.`,
			},
			"parameters": {
				Type:        schema.TypeList,
				Optional:    true,
				Computed:    true,
				Elem:        v2WorkflowExecutionParameterSchema(),
				Description: `The parameters of the workflow execution.`,
			},
			"policies": {
				Type:        schema.TypeList,
				Optional:    true,
				MaxItems:    1,
				Elem:        v2WorkflowExecutionPoliciesSchema(),
				Description: `The policies of the workflow execution.`,
			},

			// Attributes.
			"status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The status of the workflow execution.`,
			},
			"steps_execution": {
				Type:        schema.TypeList,
				Computed:    true,
				Elem:        v2WorkflowExecutionStepExecutionSchema(),
				Description: `The step execution information of the workflow execution.`,
			},
			"sub_graphs": {
				Type:        schema.TypeList,
				Computed:    true,
				Elem:        v2WorkflowExecutionSubGraphSchema(),
				Description: `The sub graph information of the workflow execution.`,
			},
			"events": {
				Type:        schema.TypeList,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: `The events of the workflow execution.`,
			},
			"duration": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The duration of the workflow execution.`,
			},
			"created_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The creation time of the workflow execution, in RFC3339 format.`,
			},

			// Internal parameters.
			"enable_force_new": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"true", "false"}, false),
				Description: utils.SchemaDesc(
					`Whether to allow parameters that do not support changes to have their change-triggered behavior set to 'ForceNew'.`,
					utils.SchemaDescInput{Internal: true},
				),
			},
		},
	}
}

func v2WorkflowExecutionDataRequirementsSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The name of the training data.`,
			},
			"type": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The type of the training data source.`,
			},
			"conditions": {
				Type:        schema.TypeList,
				Optional:    true,
				Computed:    true,
				Elem:        v2WorkflowExecutionDataRequirementConditionSchema(),
				Description: `The constraint conditions of the data.`,
			},
			"value": {
				Type:             schema.TypeString,
				Optional:         true,
				Computed:         true,
				DiffSuppressFunc: compareJsonStringValue,
				Description:      `The value of the data, in JSON format.`,
			},
			"used_steps": {
				Type:        schema.TypeList,
				Optional:    true,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: `The steps that use this data.`,
			},
			"delay": {
				Type:        schema.TypeBool,
				Optional:    true,
				Computed:    true,
				Description: `Whether this is a delayed parameter.`,
			},
		},
	}
}

func v2WorkflowExecutionDataRequirementConditionSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"attribute": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `The attribute of the constraint.`,
			},
			"operator": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `The operator of the constraint.`,
			},
			"value": {
				Type:             schema.TypeString,
				Optional:         true,
				Computed:         true,
				DiffSuppressFunc: compareJsonStringValue,
				Description:      `The value of the constraint, in JSON format.`,
			},
		},
	}
}

func v2WorkflowExecutionParameterSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `The name of the parameter.`,
			},
			"type": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `The type of the parameter.`,
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `The description of the parameter.`,
			},
			"example": {
				Type:             schema.TypeString,
				Optional:         true,
				Computed:         true,
				DiffSuppressFunc: compareJsonStringValue,
				Description:      `The example of the parameter, in JSON format.`,
			},
			"delay": {
				Type:        schema.TypeBool,
				Optional:    true,
				Computed:    true,
				Description: `Whether this is a delayed input parameter.`,
			},
			"default": {
				Type:             schema.TypeString,
				Optional:         true,
				Computed:         true,
				DiffSuppressFunc: compareJsonStringValue,
				Description:      `The default value of the parameter, in JSON format.`,
			},
			"value": {
				Type:             schema.TypeString,
				Optional:         true,
				Computed:         true,
				DiffSuppressFunc: compareJsonStringValue,
				Description:      `The value of the parameter, in JSON format.`,
			},
			"enum": {
				Type:             schema.TypeList,
				Optional:         true,
				Computed:         true,
				Elem:             &schema.Schema{Type: schema.TypeString},
				DiffSuppressFunc: compareJsonStringList,
				Description:      `The enum values of the parameter, in JSON format.`,
			},
			"used_steps": {
				Type:        schema.TypeList,
				Optional:    true,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: `The steps that use this parameter.`,
			},
			"format": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `The format of the parameter data.`,
			},
			"constraint": {
				Type:             schema.TypeString,
				Optional:         true,
				Computed:         true,
				DiffSuppressFunc: compareJsonStringValue,
				Description:      `The constraint of the parameter, in JSON format.`,
			},
		},
	}
}

func v2WorkflowExecutionPoliciesSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"use_cache": {
				Type:        schema.TypeBool,
				Optional:    true,
				Computed:    true,
				Description: `Whether to use cache.`,
			},
		},
	}
}

func v2WorkflowExecutionStepExecutionSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"step_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The name of the step.`,
			},
			"execution_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The name of the execution.`,
			},
			"name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The name of the step in this execution.`,
			},
			"uuid": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The UUID of the step in this execution.`,
			},
			"execution_uuid": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The UUID of the execution.`,
			},
			"created_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The creation time of the step execution, in RFC3339 format.`,
			},
			"updated_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The update time of the step execution, in RFC3339 format.`,
			},
			"duration": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: `The duration of the step execution.`,
			},
			"type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The type of the step execution.`,
			},
			"instance_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The instance ID of the step execution.`,
			},
			"status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The status of the step execution.`,
			},
			"inputs": {
				Type:        schema.TypeList,
				Computed:    true,
				Elem:        v2WorkflowExecutionJobInputSchema(),
				Description: `The inputs of the step execution.`,
			},
			"outputs": {
				Type:        schema.TypeList,
				Computed:    true,
				Elem:        v2WorkflowExecutionJobOutputSchema(),
				Description: `The outputs of the step execution.`,
			},
			"step_uuid": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The UUID of the step.`,
			},
			"properties": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The properties of the step, in JSON format.`,
			},
			"events": {
				Type:        schema.TypeList,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: `The events of the step execution.`,
			},
			"error_info": {
				Type:        schema.TypeList,
				Computed:    true,
				Elem:        v2WorkflowExecutionErrorInfoSchema(),
				Description: `The error information of the step execution.`,
			},
			"policy": {
				Type:        schema.TypeList,
				Computed:    true,
				Elem:        v2WorkflowExecutionStepPolicySchema(),
				Description: `The policy of the step execution.`,
			},
			"conditions_execution": {
				Type:        schema.TypeList,
				Computed:    true,
				Elem:        v2WorkflowExecutionConditionExecutionSchema(),
				Description: `The conditions execution of the step.`,
			},
			"step_title": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The title of the step.`,
			},
			"conditions": {
				Type:        schema.TypeList,
				Computed:    true,
				Elem:        v2WorkflowExecutionStepConditionSchema(),
				Description: `The conditions of the step execution.`,
			},
		},
	}
}

func v2WorkflowExecutionJobInputSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The name of the input.`,
			},
			"type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The type of the input.`,
			},
			"data": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The data of the input, in JSON format.`,
			},
			"value": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The value of the input, in JSON format.`,
			},
		},
	}
}

func v2WorkflowExecutionJobOutputSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The name of the output.`,
			},
			"type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The type of the output.`,
			},
			"config": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The config of the output, in JSON format.`,
			},
		},
	}
}

func v2WorkflowExecutionErrorInfoSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"error_code": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The error code.`,
			},
			"error_message": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The error message.`,
			},
		},
	}
}

func v2WorkflowExecutionStepPolicySchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"execution_policy": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The execution policy.`,
			},
			"use_cache": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: `Whether to use cache.`,
			},
		},
	}
}

func v2WorkflowExecutionConditionExecutionSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"result": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The result of the condition execution.`,
			},
			"metric_list": {
				Type:        schema.TypeList,
				Computed:    true,
				Elem:        v2WorkflowExecutionMetricPairSchema(),
				Description: `The metric list of the condition execution.`,
			},
		},
	}
}

func v2WorkflowExecutionMetricPairSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"key": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The key of the metric.`,
			},
			"value": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The value of the metric, in JSON format.`,
			},
		},
	}
}

func v2WorkflowExecutionStepConditionSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The type of the condition.`,
			},
			"left": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The left value of the condition, in JSON format.`,
			},
			"right": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The right value of the condition, in JSON format.`,
			},
		},
	}
}

func v2WorkflowExecutionSubGraphSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The name of the sub graph.`,
			},
			"steps": {
				Type:        schema.TypeList,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: `The steps of the sub graph.`,
			},
		},
	}
}

func buildV2WorkflowExecutionDataRequirementConditions(conditions []interface{}) []map[string]interface{} {
	if len(conditions) < 1 {
		return nil
	}

	result := make([]map[string]interface{}, 0, len(conditions))
	for _, condition := range conditions {
		result = append(result, map[string]interface{}{
			"attribute": utils.ValueIgnoreEmpty(utils.PathSearch("attribute", condition, "").(string)),
			"operator":  utils.ValueIgnoreEmpty(utils.PathSearch("operator", condition, "").(string)),
			"value":     utils.ValueIgnoreEmpty(utils.StringToJson(utils.PathSearch("value", condition, "").(string))),
		})
	}
	return result
}

func buildV2WorkflowExecutionDataRequirements(dataRequirements []interface{}) []map[string]interface{} {
	if len(dataRequirements) < 1 {
		return nil
	}

	result := make([]map[string]interface{}, 0, len(dataRequirements))
	for _, req := range dataRequirements {
		result = append(result, map[string]interface{}{
			"name": utils.PathSearch("name", req, "").(string),
			"type": utils.PathSearch("type", req, "").(string),
			"conditions": utils.ValueIgnoreEmpty(
				buildV2WorkflowExecutionDataRequirementConditions(utils.PathSearch("conditions", req, make([]interface{}, 0)).([]interface{}))),
			"value": utils.ValueIgnoreEmpty(utils.StringToJson(utils.PathSearch("value", req, "").(string))),
			"used_steps": utils.ValueIgnoreEmpty(
				utils.ExpandToStringList(utils.PathSearch("used_steps", req, make([]interface{}, 0)).([]interface{}))),
			"delay": utils.ValueIgnoreEmpty(utils.PathSearch("delay", req, nil)),
		})
	}
	return result
}

func buildV2WorkflowExecutionParameterEnum(enumList []interface{}) []map[string]interface{} {
	if len(enumList) < 1 {
		return nil
	}

	result := make([]map[string]interface{}, 0, len(enumList))
	for _, enumItem := range enumList {
		result = append(result, utils.StringToJson(enumItem.(string), "").(map[string]interface{}))
	}
	return result
}

func buildV2WorkflowExecutionParameters(parameters []interface{}) []map[string]interface{} {
	if len(parameters) < 1 {
		return nil
	}

	result := make([]map[string]interface{}, 0, len(parameters))
	for _, param := range parameters {
		result = append(result, map[string]interface{}{
			"name":        utils.ValueIgnoreEmpty(utils.PathSearch("name", param, nil)),
			"type":        utils.ValueIgnoreEmpty(utils.PathSearch("type", param, nil)),
			"description": utils.ValueIgnoreEmpty(utils.PathSearch("description", param, nil)),
			"example":     utils.ValueIgnoreEmpty(utils.StringToJson(utils.PathSearch("example", param, "").(string))),
			"delay":       utils.ValueIgnoreEmpty(utils.PathSearch("delay", param, nil)),
			"default":     utils.ValueIgnoreEmpty(utils.StringToJson(utils.PathSearch("default", param, "").(string))),
			"value":       utils.ValueIgnoreEmpty(utils.StringToJson(utils.PathSearch("value", param, "").(string))),
			"enum": utils.ValueIgnoreEmpty(
				buildV2WorkflowExecutionParameterEnum(utils.PathSearch("enum", param, make([]interface{}, 0)).([]interface{}))),
			"used_steps": utils.ValueIgnoreEmpty(
				utils.ExpandToStringList(utils.PathSearch("used_steps", param, make([]interface{}, 0)).([]interface{}))),
			"format":     utils.ValueIgnoreEmpty(utils.PathSearch("format", param, nil)),
			"constraint": utils.ValueIgnoreEmpty(utils.StringToJson(utils.PathSearch("constraint", param, "").(string))),
		})
	}
	return result
}

func buildV2WorkflowExecutionPolicies(policies []interface{}) map[string]interface{} {
	if len(policies) < 1 {
		return nil
	}

	return map[string]interface{}{
		"use_cache": utils.ValueIgnoreEmpty(utils.PathSearch("use_cache", policies[0], nil)),
	}
}

func buildV2WorkflowExecutionCreateBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"name":              utils.ValueIgnoreEmpty(d.Get("name")),
		"description":       utils.ValueIgnoreEmpty(d.Get("description")),
		"workspace_id":      utils.ValueIgnoreEmpty(d.Get("workspace_id")),
		"workflow_id":       utils.ValueIgnoreEmpty(d.Get("workflow_id")),
		"workflow_name":     utils.ValueIgnoreEmpty(d.Get("workflow_name")),
		"scene_id":          utils.ValueIgnoreEmpty(d.Get("scene_id")),
		"scene_name":        utils.ValueIgnoreEmpty(d.Get("scene_name")),
		"labels":            utils.ValueIgnoreEmpty(d.Get("labels")),
		"data_requirements": utils.ValueIgnoreEmpty(buildV2WorkflowExecutionDataRequirements(d.Get("data_requirements").([]interface{}))),
		"parameters":        utils.ValueIgnoreEmpty(buildV2WorkflowExecutionParameters(d.Get("parameters").([]interface{}))),
		"policies":          utils.ValueIgnoreEmpty(buildV2WorkflowExecutionPolicies(d.Get("policies").([]interface{}))),
	}
	return bodyParams
}

func createV2WorkflowExecution(client *golangsdk.ServiceClient, d *schema.ResourceData) (interface{}, error) {
	httpUrl := "v2/{project_id}/workflows/{workflow_id}/executions"

	createPath := client.Endpoint + httpUrl
	createPath = strings.ReplaceAll(createPath, "{project_id}", client.ProjectID)
	createPath = strings.ReplaceAll(createPath, "{workflow_id}", d.Get("workflow_id").(string))

	createOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json",
		},
		JSONBody: utils.RemoveNil(buildV2WorkflowExecutionCreateBodyParams(d)),
	}

	requestResp, err := client.Request("POST", createPath, &createOpt)
	if err != nil {
		return nil, err
	}

	return utils.FlattenResponse(requestResp)
}

func resourceV2WorkflowExecutionCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg    = meta.(*config.Config)
		region = cfg.GetRegion(d)
	)

	client, err := cfg.NewServiceClient("modelarts", region)
	if err != nil {
		return diag.Errorf("error creating ModelArts client: %s", err)
	}

	resp, err := createV2WorkflowExecution(client, d)
	if err != nil {
		return diag.Errorf("error creating workflow execution: %s", err)
	}

	executionId := utils.PathSearch("execution_id", resp, "").(string)
	if executionId == "" {
		return diag.Errorf("unable to find the workflow execution ID from the API response")
	}
	d.SetId(executionId)

	return resourceV2WorkflowExecutionRead(ctx, d, meta)
}

func GetV2WorkflowExecutionById(client *golangsdk.ServiceClient, workflowId, executionId string) (interface{}, error) {
	httpUrl := "v2/{project_id}/workflows/{workflow_id}/executions/{execution_id}"

	getPath := client.Endpoint + httpUrl
	getPath = strings.ReplaceAll(getPath, "{project_id}", client.ProjectID)
	getPath = strings.ReplaceAll(getPath, "{workflow_id}", workflowId)
	getPath = strings.ReplaceAll(getPath, "{execution_id}", executionId)

	opt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json",
		},
	}

	requestResp, err := client.Request("GET", getPath, &opt)
	if err != nil {
		return nil, err
	}
	return utils.FlattenResponse(requestResp)
}

func flattenV2WorkflowExecutionDataRequirementConditions(conditions []interface{}) []map[string]interface{} {
	if len(conditions) < 1 {
		return nil
	}

	result := make([]map[string]interface{}, 0, len(conditions))
	for _, condition := range conditions {
		result = append(result, map[string]interface{}{
			"attribute": utils.PathSearch("attribute", condition, nil),
			"operator":  utils.PathSearch("operator", condition, nil),
			"value":     utils.JsonToString(utils.PathSearch("value", condition, nil)),
		})
	}
	return result
}

func flattenV2WorkflowExecutionDataRequirements(dataRequirements []interface{}) []map[string]interface{} {
	if len(dataRequirements) < 1 {
		return nil
	}

	result := make([]map[string]interface{}, 0, len(dataRequirements))
	for _, req := range dataRequirements {
		result = append(result, map[string]interface{}{
			"name": utils.PathSearch("name", req, nil),
			"type": utils.PathSearch("type", req, nil),
			"conditions": flattenV2WorkflowExecutionDataRequirementConditions(
				utils.PathSearch("conditions", req, make([]interface{}, 0)).([]interface{})),
			"value":      utils.JsonToString(utils.PathSearch("value", req, nil)),
			"used_steps": utils.PathSearch("used_steps", req, make([]interface{}, 0)).([]interface{}),
			"delay":      utils.PathSearch("delay", req, nil),
		})
	}
	return result
}

func flattenV2WorkflowExecutionParametersEnum(enumList []interface{}) []string {
	if len(enumList) < 1 {
		return nil
	}

	result := make([]string, 0, len(enumList))
	for _, enum := range enumList {
		result = append(result, utils.JsonToString(enum))
	}
	return result
}

func flattenV2WorkflowExecutionParameters(parameters []interface{}) []map[string]interface{} {
	if len(parameters) < 1 {
		return nil
	}

	result := make([]map[string]interface{}, 0, len(parameters))
	for _, param := range parameters {
		result = append(result, map[string]interface{}{
			"name":        utils.PathSearch("name", param, nil),
			"type":        utils.PathSearch("type", param, nil),
			"description": utils.PathSearch("description", param, nil),
			"example":     utils.JsonToString(utils.PathSearch("example", param, nil)),
			"delay":       utils.PathSearch("delay", param, nil),
			"default":     utils.JsonToString(utils.PathSearch("default", param, nil)),
			"value":       utils.JsonToString(utils.PathSearch("value", param, nil)),
			"enum":        flattenV2WorkflowExecutionParametersEnum(utils.PathSearch("enum", param, make([]interface{}, 0)).([]interface{})),
			"used_steps":  utils.PathSearch("used_steps", param, make([]interface{}, 0)).([]interface{}),
			"format":      utils.PathSearch("format", param, nil),
			"constraint":  utils.JsonToString(utils.PathSearch("constraint", param, nil)),
		})
	}
	return result
}

func flattenV2WorkflowExecutionJobInputs(inputs []interface{}) []map[string]interface{} {
	if len(inputs) < 1 {
		return nil
	}

	result := make([]map[string]interface{}, 0, len(inputs))
	for _, input := range inputs {
		result = append(result, map[string]interface{}{
			"name":  utils.PathSearch("name", input, nil),
			"type":  utils.PathSearch("type", input, nil),
			"data":  utils.JsonToString(utils.PathSearch("data", input, nil)),
			"value": utils.JsonToString(utils.PathSearch("value", input, nil)),
		})
	}
	return result
}

func flattenV2WorkflowExecutionJobOutputs(outputs []interface{}) []map[string]interface{} {
	if len(outputs) < 1 {
		return nil
	}

	result := make([]map[string]interface{}, 0, len(outputs))
	for _, output := range outputs {
		result = append(result, map[string]interface{}{
			"name":   utils.PathSearch("name", output, nil),
			"type":   utils.PathSearch("type", output, nil),
			"config": utils.JsonToString(utils.PathSearch("config", output, nil)),
		})
	}
	return result
}

func flattenV2WorkflowExecutionErrorInfo(errorInfo interface{}) []map[string]interface{} {
	if errorInfo == nil {
		return nil
	}

	return []map[string]interface{}{
		{
			"error_code":    utils.PathSearch("error_code", errorInfo, nil),
			"error_message": utils.PathSearch("error_message", errorInfo, nil),
		},
	}
}

func flattenV2WorkflowExecutionStepPolicy(policy interface{}) []map[string]interface{} {
	if policy == nil {
		return nil
	}

	return []map[string]interface{}{
		{
			"execution_policy": utils.PathSearch("execution_policy", policy, nil),
			"use_cache":        utils.PathSearch("use_cache", policy, nil),
		},
	}
}

func flattenV2WorkflowExecutionMetricPairs(metricPairs []interface{}) []map[string]interface{} {
	if len(metricPairs) < 1 {
		return nil
	}

	result := make([]map[string]interface{}, 0, len(metricPairs))
	for _, v := range metricPairs {
		metricPair := v.(map[string]interface{})
		result = append(result, map[string]interface{}{
			"key":   utils.PathSearch("key", metricPair, nil),
			"value": utils.JsonToString(utils.PathSearch("value", metricPair, nil)),
		})
	}
	return result
}

func flattenV2WorkflowExecutionConditionExecution(conditionExecution interface{}) []map[string]interface{} {
	if conditionExecution == nil {
		return nil
	}

	return []map[string]interface{}{
		{
			"result": utils.PathSearch("result", conditionExecution, nil),
			"metric_list": flattenV2WorkflowExecutionMetricPairs(
				utils.PathSearch("metric_list", conditionExecution, make([]interface{}, 0)).([]interface{})),
		},
	}
}

func flattenV2WorkflowExecutionStepConditions(conditions []interface{}) []map[string]interface{} {
	if len(conditions) < 1 {
		return nil
	}

	result := make([]map[string]interface{}, 0, len(conditions))
	for _, condition := range conditions {
		result = append(result, map[string]interface{}{
			"type":  utils.PathSearch("type", condition, nil),
			"left":  utils.JsonToString(utils.PathSearch("left", condition, nil)),
			"right": utils.JsonToString(utils.PathSearch("right", condition, nil)),
		})
	}
	return result
}

func flattenV2WorkflowExecutionStepsExecution(stepsExecutions []interface{}) []map[string]interface{} {
	if len(stepsExecutions) < 1 {
		return nil
	}

	result := make([]map[string]interface{}, 0, len(stepsExecutions))
	for _, stepExecution := range stepsExecutions {
		result = append(result, map[string]interface{}{
			"step_name":      utils.PathSearch("step_name", stepExecution, nil),
			"execution_name": utils.PathSearch("execution_name", stepExecution, nil),
			"name":           utils.PathSearch("name", stepExecution, nil),
			"uuid":           utils.PathSearch("uuid", stepExecution, nil),
			"execution_uuid": utils.PathSearch("execution_uuid", stepExecution, nil),
			"duration":       utils.PathSearch("duration", stepExecution, nil),
			"type":           utils.PathSearch("type", stepExecution, nil),
			"instance_id":    utils.PathSearch("instance_id", stepExecution, nil),
			"status":         utils.PathSearch("status", stepExecution, nil),
			"inputs": flattenV2WorkflowExecutionJobInputs(
				utils.PathSearch("inputs", stepExecution, make([]interface{}, 0)).([]interface{})),
			"outputs": flattenV2WorkflowExecutionJobOutputs(
				utils.PathSearch("outputs", stepExecution, make([]interface{}, 0)).([]interface{})),
			"step_uuid":            utils.PathSearch("step_uuid", stepExecution, nil),
			"properties":           utils.JsonToString(utils.PathSearch("properties", stepExecution, nil)),
			"events":               utils.ExpandToStringList(utils.PathSearch("events", stepExecution, make([]interface{}, 0)).([]interface{})),
			"error_info":           flattenV2WorkflowExecutionErrorInfo(utils.PathSearch("error_info", stepExecution, nil)),
			"policy":               flattenV2WorkflowExecutionStepPolicy(utils.PathSearch("policy", stepExecution, nil)),
			"conditions_execution": flattenV2WorkflowExecutionConditionExecution(utils.PathSearch("conditions_execution", stepExecution, nil)),
			"step_title":           utils.PathSearch("step_title", stepExecution, nil),
			"conditions": flattenV2WorkflowExecutionStepConditions(
				utils.PathSearch("conditions", stepExecution, make([]interface{}, 0)).([]interface{})),
			"created_at": utils.FormatTimeStampRFC3339(
				utils.ConvertTimeStrToNanoTimestamp(utils.PathSearch("created_at", stepExecution, "").(string))/1000, false),
			"updated_at": utils.FormatTimeStampRFC3339(
				utils.ConvertTimeStrToNanoTimestamp(utils.PathSearch("updated_at", stepExecution, "").(string))/1000, false),
		})
	}
	return result
}

func flattenV2WorkflowExecutionSubGraphs(subGraphs []interface{}) []map[string]interface{} {
	if len(subGraphs) < 1 {
		return nil
	}

	result := make([]map[string]interface{}, 0, len(subGraphs))
	for _, subGraph := range subGraphs {
		result = append(result, map[string]interface{}{
			"name":  utils.PathSearch("name", subGraph, nil),
			"steps": utils.PathSearch("steps", subGraph, make([]interface{}, 0)),
		})
	}
	return result
}

func resourceV2WorkflowExecutionRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg         = meta.(*config.Config)
		region      = cfg.GetRegion(d)
		workflowId  = d.Get("workflow_id").(string)
		executionId = d.Id()
	)

	client, err := cfg.NewServiceClient("modelarts", region)
	if err != nil {
		return diag.Errorf("error creating ModelArts client: %s", err)
	}

	resp, err := GetV2WorkflowExecutionById(client, workflowId, executionId)
	if err != nil {
		return common.CheckDeletedDiag(d, common.ConvertExpected400ErrInto404Err(err, "error_code", v2WorkflowExecutionNotFoundErrCodes...),
			fmt.Sprintf("error retrieving workflow execution (%s)", executionId))
	}

	mErr := multierror.Append(
		d.Set("region", region),
		d.Set("data_requirements", flattenV2WorkflowExecutionDataRequirements(
			utils.PathSearch("data_requirements", resp, make([]interface{}, 0)).([]interface{}))),
		d.Set("parameters", flattenV2WorkflowExecutionParameters(
			utils.PathSearch("parameters", resp, make([]interface{}, 0)).([]interface{}))),
		d.Set("status", utils.PathSearch("status", resp, nil)),
		d.Set("steps_execution", flattenV2WorkflowExecutionStepsExecution(
			utils.PathSearch("steps_execution", resp, make([]interface{}, 0)).([]interface{}))),
		d.Set("sub_graphs", flattenV2WorkflowExecutionSubGraphs(
			utils.PathSearch("sub_graphs", resp, make([]interface{}, 0)).([]interface{}))),
		d.Set("events", utils.ExpandToStringList(utils.PathSearch("events", resp, make([]interface{}, 0)).([]interface{}))),
		d.Set("duration", fmt.Sprintf("%v", utils.PathSearch("duration", resp, nil))),
		d.Set("created_at", utils.FormatTimeStampRFC3339(
			utils.ConvertTimeStrToNanoTimestamp(utils.PathSearch("created_at", resp, "").(string))/1000, false)),
	)
	return diag.FromErr(mErr.ErrorOrNil())
}

func buildV2WorkflowExecutionUpdateBodyParams(d *schema.ResourceData) map[string]interface{} {
	return map[string]interface{}{
		"labels":            utils.ValueIgnoreEmpty(d.Get("labels")),
		"data_requirements": utils.ValueIgnoreEmpty(buildV2WorkflowExecutionDataRequirements(d.Get("data_requirements").([]interface{}))),
		"parameters":        utils.ValueIgnoreEmpty(buildV2WorkflowExecutionParameters(d.Get("parameters").([]interface{}))),
	}
}

func updateV2WorkflowExecution(client *golangsdk.ServiceClient, d *schema.ResourceData) error {
	var (
		httpUrl     = "v2/{project_id}/workflows/{workflow_id}/executions/{execution_id}"
		workflowId  = d.Get("workflow_id").(string)
		executionId = d.Id()
	)

	updatePath := client.Endpoint + httpUrl
	updatePath = strings.ReplaceAll(updatePath, "{project_id}", client.ProjectID)
	updatePath = strings.ReplaceAll(updatePath, "{workflow_id}", workflowId)
	updatePath = strings.ReplaceAll(updatePath, "{execution_id}", executionId)

	opt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json",
		},
		JSONBody: utils.RemoveNil(buildV2WorkflowExecutionUpdateBodyParams(d)),
	}

	_, err := client.Request("PUT", updatePath, &opt)
	return err
}

func resourceV2WorkflowExecutionUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg    = meta.(*config.Config)
		region = cfg.GetRegion(d)
	)

	client, err := cfg.NewServiceClient("modelarts", region)
	if err != nil {
		return diag.Errorf("error creating ModelArts client: %s", err)
	}

	err = updateV2WorkflowExecution(client, d)
	if err != nil {
		return diag.Errorf("error updating workflow execution: %s", err)
	}

	return resourceV2WorkflowExecutionRead(ctx, d, meta)
}

func deleteV2WorkflowExecution(client *golangsdk.ServiceClient, workflowId, executionId string) error {
	httpUrl := "v2/{project_id}/workflows/{workflow_id}/executions/{execution_id}"

	deletePath := client.Endpoint + httpUrl
	deletePath = strings.ReplaceAll(deletePath, "{project_id}", client.ProjectID)
	deletePath = strings.ReplaceAll(deletePath, "{workflow_id}", workflowId)
	deletePath = strings.ReplaceAll(deletePath, "{execution_id}", executionId)

	opt := golangsdk.RequestOpts{
		MoreHeaders: map[string]string{
			"Content-Type": "application/json",
		},
	}

	_, err := client.Request("DELETE", deletePath, &opt)
	return err
}

func resourceV2WorkflowExecutionDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg         = meta.(*config.Config)
		region      = cfg.GetRegion(d)
		workflowId  = d.Get("workflow_id").(string)
		executionId = d.Id()
	)

	client, err := cfg.NewServiceClient("modelarts", region)
	if err != nil {
		return diag.Errorf("error creating ModelArts client: %s", err)
	}

	err = deleteV2WorkflowExecution(client, workflowId, executionId)
	if err != nil {
		return common.CheckDeletedDiag(d, common.ConvertExpected400ErrInto404Err(err, "error_code", v2WorkflowExecutionNotFoundErrCodes...),
			fmt.Sprintf("error deleting workflow execution (%s)", executionId))
	}

	return nil
}

func resourceV2WorkflowExecutionImportState(_ context.Context, d *schema.ResourceData, _ interface{}) ([]*schema.ResourceData, error) {
	importedId := d.Id()
	parts := strings.SplitN(importedId, "/", 2)
	if len(parts) != 2 {
		return nil, fmt.Errorf("invalid format specified for import ID, want '<workflow_id>/<id>', but got '%s'", importedId)
	}

	d.SetId(parts[1])
	return []*schema.ResourceData{d}, d.Set("workflow_id", parts[0])
}
