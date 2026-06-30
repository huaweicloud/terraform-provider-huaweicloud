package modelarts

import (
	"context"
	"fmt"
	"strings"

	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

var workflowExecutionActionNonUpdatableParams = []string{
	"workflow_id",
	"execution_id",
	"action_name",
	"data_requirements",
	"data_requirements.*.name",
	"data_requirements.*.type",
	"data_requirements.*.conditions",
	"data_requirements.*.conditions.*.attribute",
	"data_requirements.*.conditions.*.operator",
	"data_requirements.*.conditions.*.value",
	"data_requirements.*.value",
	"data_requirements.*.used_steps",
	"data_requirements.*.delay",
	"parameters",
	"parameters.*.name",
	"parameters.*.type",
	"parameters.*.description",
	"parameters.*.example",
	"parameters.*.delay",
	"parameters.*.default",
	"parameters.*.value",
	"parameters.*.enum",
	"parameters.*.used_steps",
	"parameters.*.format",
	"parameters.*.constraint",
	"policies",
	"policies.*.rerun_steps",
}

// @API ModelArts POST /v2/{project_id}/workflows/{workflow_id}/executions/{execution_id}/actions
func ResourceV2WorkflowExecutionAction() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceWorkflowExecutionActionCreate,
		ReadContext:   resourceWorkflowExecutionActionRead,
		UpdateContext: resourceWorkflowExecutionActionUpdate,
		DeleteContext: resourceWorkflowExecutionActionDelete,

		CustomizeDiff: config.FlexibleForceNew(workflowExecutionActionNonUpdatableParams),

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: `The region where the workflow execution is located.`,
			},

			// Required parameters.
			"workflow_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The workflow ID.`,
			},
			"execution_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The workflow execution ID.`,
			},
			"action_name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The action name.`,
			},

			// Optional parameters.
			"data_requirements": {
				Type:        schema.TypeList,
				Optional:    true,
				Elem:        workflowExecutionActionDataRequirementsSchema(),
				Description: `The data requirements used by the workflow steps.`,
			},
			"parameters": {
				Type:        schema.TypeList,
				Optional:    true,
				Elem:        workflowExecutionActionParameterSchema(),
				Description: `The parameters used by the workflow steps.`,
			},
			"policies": {
				Type:        schema.TypeList,
				Optional:    true,
				MaxItems:    1,
				Elem:        workflowExecutionActionPoliciesSchema(),
				Description: `The execution policies used by the execution record.`,
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

func workflowExecutionActionDataRequirementConditionSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"attribute": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The attribute of the constraint.`,
			},
			"operator": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The operator of the constraint.`,
			},
			"value": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The value of the constraint, in JSON format.`,
			},
		},
	}
}

func workflowExecutionActionDataRequirementsSchema() *schema.Resource {
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
				Elem:        workflowExecutionActionDataRequirementConditionSchema(),
				Description: `The constraint conditions of the data.`,
			},
			"value": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The value of the data, in JSON format.`,
			},
			"used_steps": {
				Type:        schema.TypeList,
				Optional:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: `The steps that use this data.`,
			},
			"delay": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: `Whether this is a delayed parameter.`,
			},
		},
	}
}

func workflowExecutionActionParameterSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The name of the parameter.`,
			},
			"type": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The type of the parameter.`,
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The description of the parameter.`,
			},
			"example": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The example of the parameter, in JSON format.`,
			},
			"delay": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: `Whether this is a delayed input parameter.`,
			},
			"default": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The default value of the parameter, in JSON format.`,
			},
			"value": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The value of the parameter, in JSON format.`,
			},
			"enum": {
				Type:        schema.TypeList,
				Optional:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: `The enum values of the parameter, in JSON format.`,
			},
			"used_steps": {
				Type:        schema.TypeList,
				Optional:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: `The steps that use this parameter.`,
			},
			"format": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The format of the parameter data.`,
			},
			"constraint": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The constraint of the parameter, in JSON format.`,
			},
		},
	}
}

func workflowExecutionActionPoliciesSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"rerun_steps": {
				Type:        schema.TypeList,
				Optional:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: `The steps to rerun.`,
			},
		},
	}
}

func buildWorkflowExecutionActionDataRequirements(dataRequirements []interface{}) []map[string]interface{} {
	if len(dataRequirements) < 1 {
		return nil
	}

	result := make([]map[string]interface{}, 0, len(dataRequirements))
	for _, req := range dataRequirements {
		result = append(result, map[string]interface{}{
			"name": utils.PathSearch("name", req, "").(string),
			"type": utils.PathSearch("type", req, "").(string),
			"conditions": utils.ValueIgnoreEmpty(
				buildWorkflowExecutionActionDataRequirementConditions(utils.PathSearch("conditions", req, make([]interface{}, 0)).([]interface{}))),
			"value": utils.ValueIgnoreEmpty(utils.StringToJson(utils.PathSearch("value", req, "").(string))),
			"used_steps": utils.ValueIgnoreEmpty(
				utils.ExpandToStringList(utils.PathSearch("used_steps", req, make([]interface{}, 0)).([]interface{}))),
			"delay": utils.ValueIgnoreEmpty(utils.PathSearch("delay", req, nil)),
		})
	}
	return result
}

func buildWorkflowExecutionActionDataRequirementConditions(conditions []interface{}) []map[string]interface{} {
	if len(conditions) < 1 || conditions[0] == nil {
		return nil
	}

	result := make([]map[string]interface{}, 0, len(conditions))
	for _, condition := range conditions {
		result = append(result, map[string]interface{}{
			"attribute": utils.ValueIgnoreEmpty(utils.PathSearch("attribute", condition, nil)),
			"operator":  utils.ValueIgnoreEmpty(utils.PathSearch("operator", condition, nil)),
			"value":     utils.ValueIgnoreEmpty(utils.StringToJson(utils.PathSearch("value", condition, "").(string))),
		})
	}
	return result
}

func buildWorkflowExecutionActionParameterEnum(enumList []interface{}) []interface{} {
	if len(enumList) < 1 {
		return nil
	}

	result := make([]interface{}, 0)
	for _, enumItem := range enumList {
		result = append(result, utils.StringToJson(enumItem.(string)))
	}
	return result
}

func buildWorkflowExecutionActionParameters(parameters []interface{}) []map[string]interface{} {
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
				buildWorkflowExecutionActionParameterEnum(utils.PathSearch("enum", param, make([]interface{}, 0)).([]interface{}))),
			"used_steps": utils.ValueIgnoreEmpty(
				utils.ExpandToStringList(utils.PathSearch("used_steps", param, make([]interface{}, 0)).([]interface{}))),
			"format":     utils.ValueIgnoreEmpty(utils.PathSearch("format", param, nil)),
			"constraint": utils.ValueIgnoreEmpty(utils.StringToJson(utils.PathSearch("constraint", param, "").(string))),
		})
	}
	return result
}

func buildWorkflowExecutionActionPolicies(policies []interface{}) map[string]interface{} {
	if len(policies) < 1 || policies[0] == nil {
		return nil
	}

	return map[string]interface{}{
		"rerun_steps": utils.ValueIgnoreEmpty(
			utils.ExpandToStringList(utils.PathSearch("rerun_steps", policies[0], make([]interface{}, 0)).([]interface{}))),
	}
}

func buildWorkflowExecutionActionBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"action_name":       d.Get("action_name"),
		"data_requirements": utils.ValueIgnoreEmpty(buildWorkflowExecutionActionDataRequirements(d.Get("data_requirements").([]interface{}))),
		"parameters":        utils.ValueIgnoreEmpty(buildWorkflowExecutionActionParameters(d.Get("parameters").([]interface{}))),
		"policies":          utils.ValueIgnoreEmpty(buildWorkflowExecutionActionPolicies(d.Get("policies").([]interface{}))),
	}
	return bodyParams
}

func createWorkflowExecutionAction(client *golangsdk.ServiceClient, d *schema.ResourceData) error {
	var (
		httpUrl     = "v2/{project_id}/workflows/{workflow_id}/executions/{execution_id}/actions"
		workflowId  = d.Get("workflow_id").(string)
		executionId = d.Get("execution_id").(string)
	)

	createPath := client.Endpoint + httpUrl
	createPath = strings.ReplaceAll(createPath, "{project_id}", client.ProjectID)
	createPath = strings.ReplaceAll(createPath, "{workflow_id}", workflowId)
	createPath = strings.ReplaceAll(createPath, "{execution_id}", executionId)

	createOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json",
		},
		JSONBody: utils.RemoveNil(buildWorkflowExecutionActionBodyParams(d)),
	}

	resp, err := client.Request("POST", createPath, &createOpt)
	if err != nil {
		return err
	}

	respBody, err := utils.FlattenResponse(resp)
	if err != nil {
		return err
	}

	errorCode := utils.PathSearch("error_code", respBody, "").(string)
	if errorCode != "" {
		return fmt.Errorf("error_code=%s, error_msg=%s", errorCode, utils.PathSearch("error_msg", respBody, ""))
	}

	return nil
}

func resourceWorkflowExecutionActionCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg    = meta.(*config.Config)
		region = cfg.GetRegion(d)
	)

	client, err := cfg.NewServiceClient("modelarts", region)
	if err != nil {
		return diag.Errorf("error creating ModelArts client: %s", err)
	}

	err = createWorkflowExecutionAction(client, d)
	if err != nil {
		return diag.Errorf("error creating workflow execution action: %s", err)
	}

	randomUUID, err := uuid.NewRandom()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(randomUUID.String())

	return resourceWorkflowExecutionActionRead(ctx, d, meta)
}

func resourceWorkflowExecutionActionRead(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return nil
}

func resourceWorkflowExecutionActionUpdate(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return nil
}

func resourceWorkflowExecutionActionDelete(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	errorMsg := `This resource is only a one-time action resource for operating the workflow execution. Deleting this resource will
not clear the corresponding request record, but will only remove the resource information from the tfstate file.`
	return diag.Diagnostics{
		diag.Diagnostic{
			Severity: diag.Warning,
			Summary:  errorMsg,
		},
	}
}
