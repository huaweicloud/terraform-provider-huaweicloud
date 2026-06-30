package modelarts

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/google/uuid"
	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API ModelArts GET /v2/{project_id}/workflows/{workflow_id}/executions
func DataSourceV2WorkflowExecutions() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceV2WorkflowExecutionsRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `The region where the workflow executions are located.`,
			},

			// Required parameters.
			"workflow_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The workflow ID of the execution records to be queried.`,
			},

			// Optional parameters.
			"workspace_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The workspace ID of the execution records to be queried.`,
			},
			"labels": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The labels of the execution records to be queried.`,
			},
			"status": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The status of the execution records to be queried.`,
			},
			"scene_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The scene ID of the execution records to be queried.`,
			},

			// Attributes.
			"executions": {
				Type:        schema.TypeList,
				Computed:    true,
				Elem:        dataV2WorkflowExecutionsElemSchema(),
				Description: `The list of the workflow executions that matched filter parameters.`,
			},
		},
	}
}

func dataV2WorkflowExecutionsElemSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The ID of the workflow execution.`,
			},
			"name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The name of the workflow execution.`,
			},
			"description": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The description of the workflow execution.`,
			},
			"status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The status of the workflow execution.`,
			},
			"workspace_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The workspace ID to which the workflow execution belongs.`,
			},
			"workflow_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The workflow ID.`,
			},
			"workflow_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The workflow name.`,
			},
			"scene_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The custom scene ID.`,
			},
			"scene_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The custom scene name.`,
			},
			"labels": {
				Type:        schema.TypeList,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: `The labels of the workflow execution.`,
			},
			"data_requirements": {
				Type:        schema.TypeList,
				Computed:    true,
				Elem:        dataV2WorkflowExecutionDataRequirementSchema(),
				Description: `The data requirements used by the workflow execution steps.`,
			},
			"parameters": {
				Type:        schema.TypeList,
				Computed:    true,
				Elem:        dataV2WorkflowExecutionParameterSchema(),
				Description: `The parameters used by the workflow execution steps.`,
			},
			"policies": {
				Type:        schema.TypeList,
				Computed:    true,
				Elem:        dataV2WorkflowExecutionPolicySchema(),
				Description: `The execution policies used by the workflow execution.`,
			},
			"steps_execution": {
				Type:        schema.TypeList,
				Computed:    true,
				Elem:        dataV2WorkflowExecutionStepExecutionSchema(),
				Description: `The step execution information of the workflow execution.`,
			},
			"sub_graphs": {
				Type:        schema.TypeList,
				Computed:    true,
				Elem:        dataV2WorkflowExecutionSubGraphSchema(),
				Description: `The subgraph information of the workflow execution.`,
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
		},
	}
}

func dataV2WorkflowExecutionDataRequirementConditionSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"attribute": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The condition attribute.`,
			},
			"operator": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The operator of the condition.`,
			},
			"value": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The value of the condition, in JSON format.`,
			},
		},
	}
}

func dataV2WorkflowExecutionDataRequirementSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The name of the training data.`,
			},
			"type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The type of the data source.`,
			},
			"conditions": {
				Type:        schema.TypeList,
				Computed:    true,
				Elem:        dataV2WorkflowExecutionDataRequirementConditionSchema(),
				Description: `The data constraint conditions.`,
			},
			"value": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The value of the data, in JSON format.`,
			},
			"used_steps": {
				Type:        schema.TypeList,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: `The workflow steps that use this data.`,
			},
			"delay": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: `Whether the data is delayed.`,
			},
		},
	}
}

func dataV2WorkflowExecutionParameterSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The name of the parameter.`,
			},
			"type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The type of the parameter.`,
			},
			"description": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The description of the parameter.`,
			},
			"example": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The example of the parameter, in JSON format.`,
			},
			"delay": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: `Whether the parameter is delayed.`,
			},
			"default": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The default value of the parameter, in JSON format.`,
			},
			"value": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The value of the parameter, in JSON format.`,
			},
			"enum": {
				Type:        schema.TypeList,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: `The enumeration values of the parameter, in JSON format.`,
			},
			"used_steps": {
				Type:        schema.TypeList,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: `The workflow steps that use this parameter.`,
			},
			"format": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The data format of the parameter.`,
			},
			"constraint": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The constraint of the parameter, in JSON format.`,
			},
		},
	}
}

func dataV2WorkflowExecutionPolicySchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"use_cache": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: `Whether to use cache.`,
			},
		},
	}
}

func dataV2WorkflowExecutionStepInputSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The name of the input data.`,
			},
			"type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The type of the input.`,
			},
			"data": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The input data, in JSON format.`,
			},
			"value": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The value of the input, in JSON format.`,
			},
		},
	}
}

func dataV2WorkflowExecutionStepOutputSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The name of the output data.`,
			},
			"type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The type of the output.`,
			},
			"config": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The output configuration, in JSON format.`,
			},
		},
	}
}

func dataV2WorkflowExecutionStepErrorInfoSchema() *schema.Resource {
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

func dataV2WorkflowExecutionStepPolicySchema() *schema.Resource {
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

func dataV2WorkflowExecutionMetricPairSchema() *schema.Resource {
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

func dataV2WorkflowExecutionConditionExecutionSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"result": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The execution result, in JSON format.`,
			},
			"metric_list": {
				Type:        schema.TypeList,
				Computed:    true,
				Elem:        dataV2WorkflowExecutionMetricPairSchema(),
				Description: `The list of workflow metric information.`,
			},
		},
	}
}

func dataV2WorkflowExecutionStepConditionSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The condition type.`,
			},
			"left": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The left branch when the condition is true, in JSON format.`,
			},
			"right": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The right branch when the condition is false, in JSON format.`,
			},
		},
	}
}

func dataV2WorkflowExecutionStepExecutionSchema() *schema.Resource {
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
				Description: `The name of the execution record.`,
			},
			"name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The name of the step in this execution.`,
			},
			"uuid": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The UUID of the step in the execution instance.`,
			},
			"execution_uuid": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The UUID of the execution instance.`,
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
				Description: `The execution duration of the step, in seconds.`,
			},
			"type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The type of the step.`,
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
				Elem:        dataV2WorkflowExecutionStepInputSchema(),
				Description: `The inputs of the step.`,
			},
			"outputs": {
				Type:        schema.TypeList,
				Computed:    true,
				Elem:        dataV2WorkflowExecutionStepOutputSchema(),
				Description: `The outputs of the step.`,
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
				Description: `The events of the step.`,
			},
			"error_info": {
				Type:        schema.TypeList,
				Computed:    true,
				Elem:        dataV2WorkflowExecutionStepErrorInfoSchema(),
				Description: `The error information of the step execution.`,
			},
			"policy": {
				Type:        schema.TypeList,
				Computed:    true,
				Elem:        dataV2WorkflowExecutionStepPolicySchema(),
				Description: `The execution policy of the step.`,
			},
			"conditions_execution": {
				Type:        schema.TypeList,
				Computed:    true,
				Elem:        dataV2WorkflowExecutionConditionExecutionSchema(),
				Description: `The condition execution of the step.`,
			},
			"step_title": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The title of the step.`,
			},
			"conditions": {
				Type:        schema.TypeList,
				Computed:    true,
				Elem:        dataV2WorkflowExecutionStepConditionSchema(),
				Description: `The conditions of the step.`,
			},
		},
	}
}

func dataV2WorkflowExecutionSubGraphSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The name of the subgraph.`,
			},
			"steps": {
				Type:        schema.TypeList,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: `The step members of the subgraph.`,
			},
		},
	}
}

func buildListV2WorkflowExecutionsQueryParams(d *schema.ResourceData) string {
	res := ""

	if v, ok := d.GetOk("workspace_id"); ok {
		res = fmt.Sprintf("%s&workspace_id=%v", res, v)
	}
	if v, ok := d.GetOk("labels"); ok {
		res = fmt.Sprintf("%s&labels=%v", res, v)
	}
	if v, ok := d.GetOk("status"); ok {
		res = fmt.Sprintf("%s&status=%v", res, v)
	}
	if v, ok := d.GetOk("scene_id"); ok {
		res = fmt.Sprintf("%s&scene_id=%v", res, v)
	}

	return res
}

func listV2WorkflowExecutions(client *golangsdk.ServiceClient, d *schema.ResourceData) ([]interface{}, error) {
	var (
		httpUrl    = "v2/{project_id}/workflows/{workflow_id}/executions?limit={limit}"
		workflowID = d.Get("workflow_id").(string)
		limit      = 100
		offset     = 0
		result     = make([]interface{}, 0)
	)

	listPathWithLimit := client.Endpoint + httpUrl
	listPathWithLimit = strings.ReplaceAll(listPathWithLimit, "{project_id}", client.ProjectID)
	listPathWithLimit = strings.ReplaceAll(listPathWithLimit, "{workflow_id}", workflowID)
	listPathWithLimit = strings.ReplaceAll(listPathWithLimit, "{limit}", strconv.Itoa(limit))
	listPathWithLimit += buildListV2WorkflowExecutionsQueryParams(d)

	listOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	for {
		listPathWithOffset := fmt.Sprintf("%s&offset=%d", listPathWithLimit, offset)
		requestResp, err := client.Request("GET", listPathWithOffset, &listOpt)
		if err != nil {
			return nil, err
		}

		respBody, err := utils.FlattenResponse(requestResp)
		if err != nil {
			return nil, err
		}

		items := utils.PathSearch("items", respBody, make([]interface{}, 0)).([]interface{})
		result = append(result, items...)

		if len(items) < limit {
			break
		}
		offset += len(items)
	}
	return result, nil
}

func flattenDataV2WorkflowExecutionDataRequirementConditions(conditions []interface{}) []interface{} {
	if len(conditions) < 1 {
		return nil
	}

	result := make([]interface{}, 0, len(conditions))
	for _, condition := range conditions {
		result = append(result, map[string]interface{}{
			"attribute": utils.PathSearch("attribute", condition, nil),
			"operator":  utils.PathSearch("operator", condition, nil),
			"value":     utils.JsonToString(utils.PathSearch("value", condition, nil)),
		})
	}
	return result
}

func flattenDataV2WorkflowExecutionDataRequirements(dataRequirements []interface{}) []interface{} {
	if len(dataRequirements) < 1 {
		return nil
	}

	result := make([]interface{}, 0, len(dataRequirements))
	for _, req := range dataRequirements {
		result = append(result, map[string]interface{}{
			"name": utils.PathSearch("name", req, nil),
			"type": utils.PathSearch("type", req, nil),
			"conditions": flattenDataV2WorkflowExecutionDataRequirementConditions(
				utils.PathSearch("conditions", req, make([]interface{}, 0)).([]interface{})),
			"value":      utils.JsonToString(utils.PathSearch("value", req, nil)),
			"used_steps": utils.ExpandToStringList(utils.PathSearch("used_steps", req, make([]interface{}, 0)).([]interface{})),
			"delay":      utils.PathSearch("delay", req, nil),
		})
	}
	return result
}

func flattenDataV2WorkflowExecutionParametersEnum(enumList []interface{}) []string {
	if len(enumList) < 1 {
		return nil
	}

	result := make([]string, 0, len(enumList))
	for _, enum := range enumList {
		result = append(result, utils.JsonToString(enum))
	}
	return result
}

func flattenDataV2WorkflowExecutionParameters(parameters []interface{}) []interface{} {
	if len(parameters) < 1 {
		return nil
	}

	result := make([]interface{}, 0, len(parameters))
	for _, param := range parameters {
		result = append(result, map[string]interface{}{
			"name":        utils.PathSearch("name", param, nil),
			"type":        utils.PathSearch("type", param, nil),
			"description": utils.PathSearch("description", param, nil),
			"example":     utils.JsonToString(utils.PathSearch("example", param, nil)),
			"delay":       utils.PathSearch("delay", param, nil),
			"default":     utils.JsonToString(utils.PathSearch("default", param, nil)),
			"value":       utils.JsonToString(utils.PathSearch("value", param, nil)),
			"enum": flattenDataV2WorkflowExecutionParametersEnum(
				utils.PathSearch("enum", param, make([]interface{}, 0)).([]interface{})),
			"used_steps": utils.ExpandToStringList(utils.PathSearch("used_steps", param, make([]interface{}, 0)).([]interface{})),
			"format":     utils.PathSearch("format", param, nil),
			"constraint": utils.JsonToString(utils.PathSearch("constraint", param, nil)),
		})
	}
	return result
}

func flattenDataV2WorkflowExecutionPolicies(policies interface{}) []interface{} {
	if policies == nil {
		return nil
	}

	return []interface{}{
		map[string]interface{}{
			"use_cache": utils.PathSearch("use_cache", policies, nil),
		},
	}
}

func flattenDataV2WorkflowExecutionStepInputs(inputs []interface{}) []interface{} {
	if len(inputs) < 1 {
		return nil
	}

	result := make([]interface{}, 0, len(inputs))
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

func flattenDataV2WorkflowExecutionStepOutputs(outputs []interface{}) []interface{} {
	if len(outputs) < 1 {
		return nil
	}

	result := make([]interface{}, 0, len(outputs))
	for _, output := range outputs {
		result = append(result, map[string]interface{}{
			"name":   utils.PathSearch("name", output, nil),
			"type":   utils.PathSearch("type", output, nil),
			"config": utils.JsonToString(utils.PathSearch("config", output, nil)),
		})
	}
	return result
}

func flattenDataV2WorkflowExecutionStepErrorInfo(errorInfo interface{}) []interface{} {
	if errorInfo == nil {
		return nil
	}

	return []interface{}{
		map[string]interface{}{
			"error_code":    utils.PathSearch("error_code", errorInfo, nil),
			"error_message": utils.PathSearch("error_message", errorInfo, nil),
		},
	}
}

func flattenDataV2WorkflowExecutionStepPolicy(policy interface{}) []interface{} {
	if policy == nil {
		return nil
	}

	return []interface{}{
		map[string]interface{}{
			"execution_policy": utils.PathSearch("execution_policy", policy, nil),
			"use_cache":        utils.PathSearch("use_cache", policy, nil),
		},
	}
}

func flattenDataV2WorkflowExecutionMetricPairs(metricList []interface{}) []interface{} {
	if len(metricList) < 1 {
		return nil
	}

	result := make([]interface{}, 0, len(metricList))
	for _, metric := range metricList {
		result = append(result, map[string]interface{}{
			"key":   utils.PathSearch("key", metric, nil),
			"value": utils.JsonToString(utils.PathSearch("value", metric, nil)),
		})
	}
	return result
}

func flattenDataV2WorkflowExecutionConditionExecution(conditionExecution interface{}) []interface{} {
	if conditionExecution == nil {
		return nil
	}

	return []interface{}{
		map[string]interface{}{
			"result": utils.JsonToString(utils.PathSearch("result", conditionExecution, nil)),
			"metric_list": flattenDataV2WorkflowExecutionMetricPairs(
				utils.PathSearch("metric_list", conditionExecution, make([]interface{}, 0)).([]interface{})),
		},
	}
}

func flattenDataV2WorkflowExecutionStepConditions(conditions []interface{}) []interface{} {
	if len(conditions) < 1 {
		return nil
	}

	result := make([]interface{}, 0, len(conditions))
	for _, condition := range conditions {
		result = append(result, map[string]interface{}{
			"type":  utils.PathSearch("type", condition, nil),
			"left":  utils.JsonToString(utils.PathSearch("left", condition, nil)),
			"right": utils.JsonToString(utils.PathSearch("right", condition, nil)),
		})
	}
	return result
}

func flattenDataV2WorkflowExecutionStepExecutions(stepExecutions []interface{}) []interface{} {
	if len(stepExecutions) < 1 {
		return nil
	}

	result := make([]interface{}, 0, len(stepExecutions))
	for _, stepExec := range stepExecutions {
		result = append(result, map[string]interface{}{
			"step_name":      utils.PathSearch("step_name", stepExec, nil),
			"execution_name": utils.PathSearch("execution_name", stepExec, nil),
			"name":           utils.PathSearch("name", stepExec, nil),
			"uuid":           utils.PathSearch("uuid", stepExec, nil),
			"execution_uuid": utils.PathSearch("execution_uuid", stepExec, nil),
			"duration":       utils.PathSearch("duration", stepExec, nil),
			"type":           utils.PathSearch("type", stepExec, nil),
			"instance_id":    utils.PathSearch("instance_id", stepExec, nil),
			"status":         utils.PathSearch("status", stepExec, nil),
			"inputs": flattenDataV2WorkflowExecutionStepInputs(
				utils.PathSearch("inputs", stepExec, make([]interface{}, 0)).([]interface{})),
			"outputs": flattenDataV2WorkflowExecutionStepOutputs(
				utils.PathSearch("outputs", stepExec, make([]interface{}, 0)).([]interface{})),
			"step_uuid":            utils.PathSearch("step_uuid", stepExec, nil),
			"properties":           utils.JsonToString(utils.PathSearch("properties", stepExec, nil)),
			"events":               utils.ExpandToStringList(utils.PathSearch("events", stepExec, make([]interface{}, 0)).([]interface{})),
			"error_info":           flattenDataV2WorkflowExecutionStepErrorInfo(utils.PathSearch("error_info", stepExec, nil)),
			"policy":               flattenDataV2WorkflowExecutionStepPolicy(utils.PathSearch("policy", stepExec, nil)),
			"conditions_execution": flattenDataV2WorkflowExecutionConditionExecution(utils.PathSearch("conditions_execution", stepExec, nil)),
			"step_title":           utils.PathSearch("step_title", stepExec, nil),
			"conditions": flattenDataV2WorkflowExecutionStepConditions(
				utils.PathSearch("conditions", stepExec, make([]interface{}, 0)).([]interface{})),
			"created_at": utils.FormatTimeStampRFC3339(
				utils.ConvertTimeStrToNanoTimestamp(utils.PathSearch("created_at", stepExec, "").(string))/1000, false),
			"updated_at": utils.FormatTimeStampRFC3339(
				utils.ConvertTimeStrToNanoTimestamp(utils.PathSearch("updated_at", stepExec, "").(string))/1000, false),
		})
	}
	return result
}

func flattenDataV2WorkflowExecutionSubGraphs(subGraphs []interface{}) []interface{} {
	if len(subGraphs) < 1 {
		return nil
	}

	result := make([]interface{}, 0, len(subGraphs))
	for _, subGraph := range subGraphs {
		result = append(result, map[string]interface{}{
			"name":  utils.PathSearch("name", subGraph, nil),
			"steps": utils.ExpandToStringList(utils.PathSearch("steps", subGraph, make([]interface{}, 0)).([]interface{})),
		})
	}
	return result
}

func flattenDataV2WorkflowExecutions(executions []interface{}) []interface{} {
	if len(executions) < 1 {
		return nil
	}

	result := make([]interface{}, 0, len(executions))
	for _, item := range executions {
		result = append(result, map[string]interface{}{
			"id":            utils.PathSearch("execution_id", item, nil),
			"name":          utils.PathSearch("name", item, nil),
			"description":   utils.PathSearch("description", item, nil),
			"status":        utils.PathSearch("status", item, nil),
			"workspace_id":  utils.PathSearch("workspace_id", item, nil),
			"workflow_id":   utils.PathSearch("workflow_id", item, nil),
			"workflow_name": utils.PathSearch("workflow_name", item, nil),
			"scene_id":      utils.PathSearch("scene_id", item, nil),
			"scene_name":    utils.PathSearch("scene_name", item, nil),
			"labels":        utils.ExpandToStringList(utils.PathSearch("labels", item, make([]interface{}, 0)).([]interface{})),
			"data_requirements": flattenDataV2WorkflowExecutionDataRequirements(
				utils.PathSearch("data_requirements", item, make([]interface{}, 0)).([]interface{})),
			"parameters": flattenDataV2WorkflowExecutionParameters(
				utils.PathSearch("parameters", item, make([]interface{}, 0)).([]interface{})),
			"policies": flattenDataV2WorkflowExecutionPolicies(utils.PathSearch("policies", item, nil)),
			"steps_execution": flattenDataV2WorkflowExecutionStepExecutions(
				utils.PathSearch("steps_execution", item, make([]interface{}, 0)).([]interface{})),
			"sub_graphs": flattenDataV2WorkflowExecutionSubGraphs(
				utils.PathSearch("sub_graphs", item, make([]interface{}, 0)).([]interface{})),
			"events":   utils.ExpandToStringList(utils.PathSearch("events", item, make([]interface{}, 0)).([]interface{})),
			"duration": fmt.Sprintf("%v", utils.PathSearch("duration", item, nil)),
			"created_at": utils.FormatTimeStampRFC3339(
				utils.ConvertTimeStrToNanoTimestamp(utils.PathSearch("created_at", item, "").(string))/1000, false),
		})
	}

	return result
}

func dataSourceV2WorkflowExecutionsRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg    = meta.(*config.Config)
		region = cfg.GetRegion(d)
	)

	client, err := cfg.NewServiceClient("modelarts", region)
	if err != nil {
		return diag.Errorf("error creating ModelArts client: %s", err)
	}

	executions, err := listV2WorkflowExecutions(client, d)
	if err != nil {
		return diag.Errorf("error querying workflow executions: %s", err)
	}

	randomUUID, err := uuid.NewRandom()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(randomUUID.String())

	mErr := multierror.Append(
		d.Set("region", region),
		d.Set("executions", flattenDataV2WorkflowExecutions(executions)),
	)
	return diag.FromErr(mErr.ErrorOrNil())
}
