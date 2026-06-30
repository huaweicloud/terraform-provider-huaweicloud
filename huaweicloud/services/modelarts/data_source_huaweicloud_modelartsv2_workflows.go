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

// @API ModelArts GET /v2/{project_id}/workflows
func DataSourceV2Workflows() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceV2WorkflowsRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `The region where the workflows are located.`,
			},

			// Optional parameters.
			"name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The name of the workflow to be queried for fuzzy matching.`,
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The description of the workflow to be queried.`,
			},
			"labels": {
				Type:        schema.TypeList,
				Optional:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: `The labels of the workflows to be queried.`,
			},
			"template_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The template ID of the workflows to be queried.`,
			},
			"search_type": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The search type of the workflows to be queried.`,
			},

			// Attributes.
			"workflows": {
				Type:        schema.TypeList,
				Computed:    true,
				Elem:        dataV2WorkflowsElemSchema(),
				Description: `The list of the workflows that matched filter parameters.`,
			},
		},
	}
}

func dataV2WorkflowsElemSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The ID of the workflow.`,
			},
			"name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The name of the workflow.`,
			},
			"description": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The description of the workflow.`,
			},
			"steps": {
				Type:        schema.TypeList,
				Computed:    true,
				Elem:        dataV2WorkflowStepSchema(),
				Description: `The steps of the workflow.`,
			},
			"workspace_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The workspace ID to which the workflow belongs.`,
			},
			"data_requirements": {
				Type:        schema.TypeList,
				Computed:    true,
				Elem:        dataV2WorkflowDataRequirementSchema(),
				Description: `The data requirements of the workflow.`,
			},
			"data": {
				Type:        schema.TypeList,
				Computed:    true,
				Elem:        dataV2WorkflowDataSchema(),
				Description: `The data of the workflow.`,
			},
			"parameters": {
				Type:        schema.TypeList,
				Computed:    true,
				Elem:        dataV2WorkflowParameterSchema(),
				Description: `The parameters of the workflow.`,
			},
			"source_workflow_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The source workflow ID for copying.`,
			},
			"gallery_subscription": {
				Type:        schema.TypeList,
				Computed:    true,
				Elem:        dataV2WorkflowGallerySubscriptionSchema(),
				Description: `The gallery subscription information of the workflow.`,
			},
			"storages": {
				Type:        schema.TypeList,
				Computed:    true,
				Elem:        dataV2WorkflowStorageSchema(),
				Description: `The unified storage definitions of the workflow.`,
			},
			"labels": {
				Type:        schema.TypeList,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: `The labels of the workflow.`,
			},
			"assets": {
				Type:        schema.TypeList,
				Computed:    true,
				Elem:        dataV2WorkflowAssetSchema(),
				Description: `The assets bound to the workflow.`,
			},
			"sub_graphs": {
				Type:        schema.TypeList,
				Computed:    true,
				Elem:        dataV2WorkflowSubGraphSchema(),
				Description: `The subgraphs of the workflow.`,
			},
			"extend": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The extended fields of the billing workflow, in JSON format.`,
			},
			"policy": {
				Type:        schema.TypeList,
				Computed:    true,
				Elem:        dataV2WorkflowPolicySchema(),
				Description: `The partial running policy of the workflow.`,
			},
			"with_subscription": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: `Whether to enable the SMN message subscription of the workflow.`,
			},
			"smn_switch": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: `Whether to enable the SMN switch of the workflow.`,
			},
			"subscription_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The SMN message subscription ID of the workflow.`,
			},
			"exeml_template_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The auto learning template ID of the workflow.`,
			},
			"package": {
				Type:        schema.TypeList,
				Computed:    true,
				Elem:        dataV2WorkflowPackageSchema(),
				Description: `The billing workflow subscription package information.`,
			},
			"created_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The creation time of the workflow, in RFC3339 format.`,
			},
			"user_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The user name that created the workflow.`,
			},
			"latest_execution": {
				Type:        schema.TypeList,
				Computed:    true,
				Elem:        dataV2WorkflowLatestExecutionSchema(),
				Description: `The latest execution information of the workflow.`,
			},
			"run_count": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: `The number of times the workflow has been run.`,
			},
			"param_ready": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: `Whether all required parameters of the workflow are filled in.`,
			},
			"source": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The source of the workflow.`,
			},
			"last_modified_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The last modified time of the workflow, in RFC3339 format.`,
			},
		},
	}
}

func dataV2WorkflowStepInputSchema() *schema.Resource {
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

func dataV2WorkflowStepOutputSchema() *schema.Resource {
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

func dataV2WorkflowStepConditionSchema() *schema.Resource {
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

func dataV2WorkflowStepPolicySchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"poll_interval_seconds": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: `The execution interval of the workflow step, in seconds.`,
			},
			"max_execution_minutes": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: `The maximum execution time of the workflow step, in minutes.`,
			},
		},
	}
}

func dataV2WorkflowStepSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The name of the workflow step.`,
			},
			"type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The type of the workflow step.`,
			},
			"inputs": {
				Type:        schema.TypeList,
				Computed:    true,
				Elem:        dataV2WorkflowStepInputSchema(),
				Description: `The inputs of the workflow step.`,
			},
			"outputs": {
				Type:        schema.TypeList,
				Computed:    true,
				Elem:        dataV2WorkflowStepOutputSchema(),
				Description: `The outputs of the workflow step.`,
			},
			"title": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The title of the workflow step.`,
			},
			"description": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The description of the workflow step.`,
			},
			"properties": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The properties of the workflow step, in JSON format.`,
			},
			"depend_steps": {
				Type:        schema.TypeList,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: `The dependent steps of the workflow step.`,
			},
			"conditions": {
				Type:        schema.TypeList,
				Computed:    true,
				Elem:        dataV2WorkflowStepConditionSchema(),
				Description: `The execution conditions of the workflow step.`,
			},
			"if_then_steps": {
				Type:        schema.TypeList,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: `The conditional branch steps of the workflow step.`,
			},
			"else_then_steps": {
				Type:        schema.TypeList,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: `The other conditional branch steps of the workflow step.`,
			},
			"policy": {
				Type:        schema.TypeList,
				Computed:    true,
				Elem:        dataV2WorkflowStepPolicySchema(),
				Description: `The execution policy of the workflow step.`,
			},
			"created_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The creation time of the workflow step, in RFC3339 format.`,
			},
		},
	}
}

func dataV2WorkflowDataRequirementConditionSchema() *schema.Resource {
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

func dataV2WorkflowDataRequirementSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The name of the data requirement.`,
			},
			"type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The type of the data source.`,
			},
			"conditions": {
				Type:        schema.TypeList,
				Computed:    true,
				Elem:        dataV2WorkflowDataRequirementConditionSchema(),
				Description: `The data constraint conditions.`,
			},
			"value": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The value of the data requirement, in JSON format.`,
			},
			"used_steps": {
				Type:        schema.TypeList,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: `The steps that use this data requirement.`,
			},
			"delay": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: `Whether the data requirement is delayed.`,
			},
		},
	}
}

func dataV2WorkflowDataSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The name of the data.`,
			},
			"type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The type of the data source.`,
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
				Description: `The steps that use this data.`,
			},
		},
	}
}

func dataV2WorkflowParameterSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The name of the workflow parameter.`,
			},
			"type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The type of the workflow parameter.`,
			},
			"description": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The description of the workflow parameter.`,
			},
			"example": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The example of the workflow parameter, in JSON format.`,
			},
			"delay": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: `Whether the workflow parameter is delayed.`,
			},
			"default": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The default value of the workflow parameter, in JSON format.`,
			},
			"value": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The value of the workflow parameter, in JSON format.`,
			},
			"enum": {
				Type:        schema.TypeList,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: `The enumeration items of the workflow parameter, in JSON format.`,
			},
			"used_steps": {
				Type:        schema.TypeList,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: `The steps that use this parameter.`,
			},
			"format": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The data format of the workflow parameter.`,
			},
			"constraint": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The constraint of the workflow parameter, in JSON format.`,
			},
		},
	}
}

func dataV2WorkflowGallerySubscriptionSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"content_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The asset ID of the gallery subscription.`,
			},
			"version_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The version ID of the gallery subscription.`,
			},
			"expired_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The expiration time of the gallery subscription, in RFC3339 format.`,
			},
		},
	}
}

func dataV2WorkflowStorageSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The name of the workflow storage.`,
			},
			"type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The type of the workflow storage.`,
			},
			"path": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The root path of the unified storage.`,
			},
		},
	}
}

func dataV2WorkflowAssetSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The name of the asset.`,
			},
			"type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The type of the asset.`,
			},
			"content_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The asset ID.`,
			},
			"subscription_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The subscription ID of the asset.`,
			},
			"expired_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The expiration time of the asset, in RFC3339 format.`,
			},
		},
	}
}

func dataV2WorkflowSubGraphSchema() *schema.Resource {
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

func dataV2WorkflowPolicySceneSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The scene ID.`,
			},
			"name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The scene name.`,
			},
			"steps": {
				Type:        schema.TypeList,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: `The step list of the scene.`,
			},
		},
	}
}

func dataV2WorkflowPolicySchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"use_scene": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The usage scenario of the workflow policy.`,
			},
			"scene_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The scene ID of the workflow policy.`,
			},
			"scenes": {
				Type:        schema.TypeList,
				Computed:    true,
				Elem:        dataV2WorkflowPolicySceneSchema(),
				Description: `The scenes of the workflow policy.`,
			},
		},
	}
}

func dataV2WorkflowPackageOrderSkuSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"code": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The billing code.`,
			},
			"period": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: `The billing period.`,
			},
			"queries_limit": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: `The query limit.`,
			},
			"price": {
				Type:        schema.TypeFloat,
				Computed:    true,
				Description: `The price.`,
			},
		},
	}
}

func dataV2WorkflowPackageOrderSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The subscription ID.`,
			},
			"sku": {
				Type:        schema.TypeList,
				Computed:    true,
				Elem:        dataV2WorkflowPackageOrderSkuSchema(),
				Description: `The subscription billing information.`,
			},
			"sku_count": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: `The subscription count.`,
			},
		},
	}
}

func dataV2WorkflowPackageSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"package_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The resource package UUID.`,
			},
			"pool_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The resource pool ID.`,
			},
			"service_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The service ID.`,
			},
			"workflow_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The workflow ID.`,
			},
			"order": {
				Type:        schema.TypeList,
				Computed:    true,
				Elem:        dataV2WorkflowPackageOrderSchema(),
				Description: `The subscription information.`,
			},
			"consume_limit": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: `The subscription limit.`,
			},
			"current_consume": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: `The current subscription consumption.`,
			},
			"current_date": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: `The current date.`,
			},
			"limit_enable": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: `Whether the limit is enabled.`,
			},
			"status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The status of the resource package.`,
			},
			"created_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The creation time of the resource package, in RFC3339 format.`,
			},
		},
	}
}

func dataV2WorkflowLatestExecutionSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"execution_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The execution ID of the workflow.`,
			},
			"created_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The creation time of the workflow execution, in RFC3339 format.`,
			},
			"status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The status of the workflow execution.`,
			},
			"running_steps": {
				Type:        schema.TypeList,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: `The running steps of the workflow execution.`,
			},
			"current_steps": {
				Type:        schema.TypeList,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: `The current steps of the workflow execution.`,
			},
			"duration": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: `The duration of the workflow execution.`,
			},
		},
	}
}

func buildListV2WorkflowsQueryParams(d *schema.ResourceData) string {
	res := ""

	if v, ok := d.GetOk("name"); ok {
		res = fmt.Sprintf("%s&name=%v", res, v)
	}
	if v, ok := d.GetOk("description"); ok {
		res = fmt.Sprintf("%s&description=%v", res, v)
	}
	if v, ok := d.GetOk("labels"); ok {
		labels := v.([]interface{})
		for _, label := range labels {
			res = fmt.Sprintf("%s&labels=%v", res, label)
		}
	}
	if v, ok := d.GetOk("template_id"); ok {
		res = fmt.Sprintf("%s&template_id=%v", res, v)
	}
	if v, ok := d.GetOk("search_type"); ok {
		res = fmt.Sprintf("%s&search_type=%v", res, v)
	}

	return res
}

func listV2Workflows(client *golangsdk.ServiceClient, d *schema.ResourceData) ([]interface{}, error) {
	var (
		httpUrl = "v2/{project_id}/workflows?limit={limit}"
		limit   = 100
		offset  = 0
		result  = make([]interface{}, 0)
	)

	listPathWithLimit := client.Endpoint + httpUrl
	listPathWithLimit = strings.ReplaceAll(listPathWithLimit, "{project_id}", client.ProjectID)
	listPathWithLimit = strings.ReplaceAll(listPathWithLimit, "{limit}", strconv.Itoa(limit))
	listPathWithLimit += buildListV2WorkflowsQueryParams(d)

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
		offset += limit
	}
	return result, nil
}

func flattenDataV2WorkflowStepInputs(inputs []interface{}) []interface{} {
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

func flattenDataV2WorkflowStepOutputs(outputs []interface{}) []interface{} {
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

func flattenDataV2WorkflowStepConditions(conditions []interface{}) []interface{} {
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

func flattenDataV2WorkflowStepPolicy(policy interface{}) []interface{} {
	if policy == nil {
		return nil
	}

	return []interface{}{
		map[string]interface{}{
			"poll_interval_seconds": utils.PathSearch("poll_interval_seconds", policy, nil),
			"max_execution_minutes": utils.PathSearch("max_execution_minutes", policy, nil),
		},
	}
}

func flattenDataV2WorkflowSteps(steps []interface{}) []interface{} {
	if len(steps) < 1 {
		return nil
	}

	result := make([]interface{}, 0, len(steps))
	for _, step := range steps {
		result = append(result, map[string]interface{}{
			"name":            utils.PathSearch("name", step, nil),
			"type":            utils.PathSearch("type", step, nil),
			"inputs":          flattenDataV2WorkflowStepInputs(utils.PathSearch("inputs", step, make([]interface{}, 0)).([]interface{})),
			"outputs":         flattenDataV2WorkflowStepOutputs(utils.PathSearch("outputs", step, make([]interface{}, 0)).([]interface{})),
			"title":           utils.PathSearch("title", step, nil),
			"description":     utils.PathSearch("description", step, nil),
			"properties":      utils.JsonToString(utils.PathSearch("properties", step, nil)),
			"depend_steps":    utils.PathSearch("depend_steps", step, make([]interface{}, 0)).([]interface{}),
			"conditions":      flattenDataV2WorkflowStepConditions(utils.PathSearch("conditions", step, make([]interface{}, 0)).([]interface{})),
			"if_then_steps":   utils.PathSearch("if_then_steps", step, make([]interface{}, 0)).([]interface{}),
			"else_then_steps": utils.PathSearch("else_then_steps", step, make([]interface{}, 0)).([]interface{}),
			"policy":          flattenDataV2WorkflowStepPolicy(utils.PathSearch("policy", step, nil)),
			"created_at": utils.FormatTimeStampRFC3339(
				utils.ConvertTimeStrToNanoTimestamp(utils.PathSearch("created_at", step, "").(string))/1000, true),
		})
	}
	return result
}

func flattenDataV2WorkflowDataRequirementConditions(conditions []interface{}) []interface{} {
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

func flattenDataV2WorkflowDataRequirements(dataRequirements []interface{}) []interface{} {
	if len(dataRequirements) < 1 {
		return nil
	}

	result := make([]interface{}, 0, len(dataRequirements))
	for _, req := range dataRequirements {
		result = append(result, map[string]interface{}{
			"name": utils.PathSearch("name", req, nil),
			"type": utils.PathSearch("type", req, nil),
			"conditions": flattenDataV2WorkflowDataRequirementConditions(
				utils.PathSearch("conditions", req, make([]interface{}, 0)).([]interface{})),
			"value":      utils.JsonToString(utils.PathSearch("value", req, nil)),
			"used_steps": utils.PathSearch("used_steps", req, make([]interface{}, 0)).([]interface{}),
			"delay":      utils.PathSearch("delay", req, nil),
		})
	}
	return result
}

func flattenDataV2WorkflowData(dataList []interface{}) []interface{} {
	if len(dataList) < 1 {
		return nil
	}

	result := make([]interface{}, 0, len(dataList))
	for _, data := range dataList {
		result = append(result, map[string]interface{}{
			"name":       utils.PathSearch("name", data, nil),
			"type":       utils.PathSearch("type", data, nil),
			"value":      utils.JsonToString(utils.PathSearch("value", data, nil)),
			"used_steps": utils.PathSearch("used_steps", data, make([]interface{}, 0)).([]interface{}),
		})
	}
	return result
}

func flattenDataV2WorkflowParametersEnum(enumList []interface{}) []string {
	if len(enumList) < 1 {
		return nil
	}

	result := make([]string, 0, len(enumList))
	for _, enum := range enumList {
		result = append(result, utils.JsonToString(enum))
	}
	return result
}

func flattenDataV2WorkflowParameters(parameters []interface{}) []interface{} {
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
			"enum": flattenDataV2WorkflowParametersEnum(
				utils.PathSearch("enum", param, make([]interface{}, 0)).([]interface{})),
			"used_steps": utils.PathSearch("used_steps", param, make([]interface{}, 0)).([]interface{}),
			"format":     utils.PathSearch("format", param, nil),
			"constraint": utils.JsonToString(utils.PathSearch("constraint", param, nil)),
		})
	}
	return result
}

func flattenDataV2WorkflowGallerySubscription(subscription interface{}) []interface{} {
	if subscription == nil {
		return nil
	}

	return []interface{}{
		map[string]interface{}{
			"content_id": utils.PathSearch("content_id", subscription, nil),
			"version_id": utils.PathSearch("version_id", subscription, nil),
			"expired_at": utils.FormatTimeStampRFC3339(
				utils.ConvertTimeStrToNanoTimestamp(utils.PathSearch("expired_at", subscription, "").(string))/1000, true),
		},
	}
}

func flattenDataV2WorkflowStorages(storages []interface{}) []interface{} {
	if len(storages) < 1 {
		return nil
	}

	result := make([]interface{}, 0, len(storages))
	for _, storage := range storages {
		result = append(result, map[string]interface{}{
			"name": utils.PathSearch("name", storage, nil),
			"type": utils.PathSearch("type", storage, nil),
			"path": utils.PathSearch("path", storage, nil),
		})
	}
	return result
}

func flattenDataV2WorkflowAssets(assets []interface{}) []interface{} {
	if len(assets) < 1 {
		return nil
	}

	result := make([]interface{}, 0, len(assets))
	for _, asset := range assets {
		result = append(result, map[string]interface{}{
			"name":            utils.PathSearch("name", asset, nil),
			"type":            utils.PathSearch("type", asset, nil),
			"content_id":      utils.PathSearch("content_id", asset, nil),
			"subscription_id": utils.PathSearch("subscription_id", asset, nil),
			"expired_at": utils.FormatTimeStampRFC3339(
				utils.ConvertTimeStrToNanoTimestamp(utils.PathSearch("expired_at", asset, "").(string))/1000, true),
		})
	}
	return result
}

func flattenDataV2WorkflowSubGraphs(subGraphs []interface{}) []interface{} {
	if len(subGraphs) < 1 {
		return nil
	}

	result := make([]interface{}, 0, len(subGraphs))
	for _, subGraph := range subGraphs {
		result = append(result, map[string]interface{}{
			"name":  utils.PathSearch("name", subGraph, nil),
			"steps": utils.PathSearch("steps", subGraph, make([]interface{}, 0)).([]interface{}),
		})
	}
	return result
}

func flattenDataV2WorkflowPolicyScenes(scenes []interface{}) []interface{} {
	if len(scenes) < 1 {
		return nil
	}

	result := make([]interface{}, 0, len(scenes))
	for _, scene := range scenes {
		result = append(result, map[string]interface{}{
			"id":    utils.PathSearch("id", scene, nil),
			"name":  utils.PathSearch("name", scene, nil),
			"steps": utils.PathSearch("steps", scene, make([]interface{}, 0)).([]interface{}),
		})
	}
	return result
}

func flattenDataV2WorkflowPolicy(policy interface{}) []interface{} {
	if policy == nil {
		return nil
	}

	return []interface{}{
		map[string]interface{}{
			"use_scene": utils.PathSearch("use_scene", policy, nil),
			"scene_id":  utils.PathSearch("scene_id", policy, nil),
			"scenes": flattenDataV2WorkflowPolicyScenes(
				utils.PathSearch("scenes", policy, make([]interface{}, 0)).([]interface{})),
		},
	}
}

func flattenDataV2WorkflowPackageOrderSku(sku interface{}) []interface{} {
	if sku == nil {
		return nil
	}

	return []interface{}{
		map[string]interface{}{
			"code":          utils.PathSearch("code", sku, nil),
			"period":        utils.PathSearch("period", sku, nil),
			"queries_limit": utils.PathSearch("queries_limit", sku, nil),
			"price":         utils.PathSearch("price", sku, nil),
		},
	}
}

func flattenDataV2WorkflowPackageOrder(order interface{}) []interface{} {
	if order == nil {
		return nil
	}

	return []interface{}{
		map[string]interface{}{
			"id":        utils.PathSearch("id", order, nil),
			"sku":       flattenDataV2WorkflowPackageOrderSku(utils.PathSearch("sku", order, nil)),
			"sku_count": utils.PathSearch("sku_count", order, nil),
		},
	}
}

func flattenDataV2WorkflowPackage(packageItem interface{}) []interface{} {
	if packageItem == nil {
		return nil
	}

	return []interface{}{
		map[string]interface{}{
			"package_id":      utils.PathSearch("package_id", packageItem, nil),
			"pool_id":         utils.PathSearch("pool_id", packageItem, nil),
			"service_id":      utils.PathSearch("service_id", packageItem, nil),
			"workflow_id":     utils.PathSearch("workflow_id", packageItem, nil),
			"order":           flattenDataV2WorkflowPackageOrder(utils.PathSearch("order", packageItem, nil)),
			"consume_limit":   utils.PathSearch("consume_limit", packageItem, nil),
			"current_consume": utils.PathSearch("current_consume", packageItem, nil),
			"current_date":    utils.PathSearch("current_date", packageItem, nil),
			"limit_enable":    utils.PathSearch("limit_enable", packageItem, nil),
			"status":          utils.PathSearch("status", packageItem, nil),
			"created_at": utils.FormatTimeStampRFC3339(
				utils.ConvertTimeStrToNanoTimestamp(utils.PathSearch("created_at", packageItem, "").(string))/1000, true),
		},
	}
}

func flattenDataV2WorkflowLatestExecution(execution interface{}) []interface{} {
	if execution == nil {
		return nil
	}

	return []interface{}{
		map[string]interface{}{
			"execution_id": utils.PathSearch("execution_id", execution, nil),
			"created_at": utils.FormatTimeStampRFC3339(
				utils.ConvertTimeStrToNanoTimestamp(utils.PathSearch("created_at", execution, "").(string))/1000, true),
			"status":        utils.PathSearch("status", execution, nil),
			"running_steps": utils.PathSearch("running_steps", execution, make([]interface{}, 0)).([]interface{}),
			"current_steps": utils.PathSearch("current_steps", execution, make([]interface{}, 0)).([]interface{}),
			"duration":      utils.PathSearch("duration", execution, nil),
		},
	}
}

func flattenDataV2Workflows(workflows []interface{}) []interface{} {
	if len(workflows) < 1 {
		return nil
	}

	result := make([]interface{}, 0, len(workflows))
	for _, item := range workflows {
		result = append(result, map[string]interface{}{
			"id":           utils.PathSearch("workflow_id", item, nil),
			"name":         utils.PathSearch("name", item, nil),
			"description":  utils.PathSearch("description", item, nil),
			"steps":        flattenDataV2WorkflowSteps(utils.PathSearch("steps", item, make([]interface{}, 0)).([]interface{})),
			"workspace_id": utils.PathSearch("workspace_id", item, nil),
			"data_requirements": flattenDataV2WorkflowDataRequirements(
				utils.PathSearch("data_requirements", item, make([]interface{}, 0)).([]interface{})),
			"data":                 flattenDataV2WorkflowData(utils.PathSearch("data", item, make([]interface{}, 0)).([]interface{})),
			"parameters":           flattenDataV2WorkflowParameters(utils.PathSearch("parameters", item, make([]interface{}, 0)).([]interface{})),
			"source_workflow_id":   utils.PathSearch("source_workflow_id", item, nil),
			"gallery_subscription": flattenDataV2WorkflowGallerySubscription(utils.PathSearch("gallery_subscription", item, nil)),
			"storages":             flattenDataV2WorkflowStorages(utils.PathSearch("storages", item, make([]interface{}, 0)).([]interface{})),
			"labels":               utils.PathSearch("labels", item, make([]interface{}, 0)).([]interface{}),
			"assets":               flattenDataV2WorkflowAssets(utils.PathSearch("assets", item, make([]interface{}, 0)).([]interface{})),
			"sub_graphs":           flattenDataV2WorkflowSubGraphs(utils.PathSearch("sub_graphs", item, make([]interface{}, 0)).([]interface{})),
			"extend":               utils.JsonToString(utils.PathSearch("extend", item, nil)),
			"policy":               flattenDataV2WorkflowPolicy(utils.PathSearch("policy", item, nil)),
			"with_subscription":    utils.PathSearch("with_subscription", item, nil),
			"smn_switch":           utils.PathSearch("smn_switch", item, nil),
			"subscription_id":      utils.PathSearch("subscription_id", item, nil),
			"exeml_template_id":    utils.PathSearch("exeml_template_id", item, nil),
			"package":              flattenDataV2WorkflowPackage(utils.PathSearch("package", item, nil)),
			"created_at": utils.FormatTimeStampRFC3339(
				utils.ConvertTimeStrToNanoTimestamp(utils.PathSearch("created_at", item, "").(string))/1000, true),
			"user_name":        utils.PathSearch("user_name", item, nil),
			"latest_execution": flattenDataV2WorkflowLatestExecution(utils.PathSearch("latest_execution", item, nil)),
			"run_count":        utils.PathSearch("run_count", item, nil),
			"param_ready":      utils.PathSearch("param_ready", item, nil),
			"source":           utils.PathSearch("source", item, nil),
			"last_modified_at": utils.FormatTimeStampRFC3339(
				utils.ConvertTimeStrToNanoTimestamp(utils.PathSearch("last_modified_at", item, "").(string))/1000, true),
		})
	}

	return result
}

func dataSourceV2WorkflowsRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg    = meta.(*config.Config)
		region = cfg.GetRegion(d)
	)

	client, err := cfg.NewServiceClient("modelarts", region)
	if err != nil {
		return diag.Errorf("error creating ModelArts client: %s", err)
	}

	workflows, err := listV2Workflows(client, d)
	if err != nil {
		return diag.Errorf("error querying workflows: %s", err)
	}

	randomUUID, err := uuid.NewRandom()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(randomUUID.String())

	mErr := multierror.Append(
		d.Set("region", region),
		d.Set("workflows", flattenDataV2Workflows(workflows)),
	)
	return diag.FromErr(mErr.ErrorOrNil())
}
