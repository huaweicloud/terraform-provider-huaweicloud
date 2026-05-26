package modelarts

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

var v2WorkflowNonUpdatableParams = []string{
	"workspace_id",
	"data",
	"data.*.name",
	"data.*.type",
	"data.*.value",
	"data.*.used_steps",
	"source_workflow_id",
	"gallery_subscription",
	"gallery_subscription.*.content_id",
	"gallery_subscription.*.version_id",
	"gallery_subscription.*.expired_at",
	"assets",
	"assets.*.name",
	"assets.*.type",
	"assets.*.content_id",
	"assets.*.subscription_id",
	"assets.*.expired_at",
	"sub_graphs",
	"sub_graphs.*.name",
	"sub_graphs.*.steps",
	"extend",
	"policy.use_scene",
	"policy.scene_id",
	"policy.*.scenes.*.id",
	"policy.*.scenes.*.name",
	"policy.*.scenes.*.steps",
	"with_subscription",
	"subscription_id",
	"exeml_template_id",
	"package",
	"package.package_id",
	"package.pool_id",
	"package.service_id",
	"package.workflow_id",
	"package.*.order",
	"package.*.order.id",
	"package.*.order.sku",
	"package.*.order.sku.code",
	"package.*.order.sku.period",
	"package.*.order.sku.queries_limit",
	"package.*.order.sku.price",
	"package.*.order.sku_count",
	"package.consume_limit",
	"package.current_consume",
	"package.current_date",
	"package.limit_enable",
}

var v2WorkflowNotFoundErrCodes = []string{
	"ModelArts.7512",
}

// @API ModelArts POST /v2/{project_id}/workflows
// @API ModelArts GET /v2/{project_id}/workflows/{workflow_id}
// @API ModelArts PUT /v2/{project_id}/workflows/{workflow_id}
// @API ModelArts DELETE /v2/{project_id}/workflows/{workflow_id}
func ResourceV2Workflow() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceV2WorkflowCreate,
		ReadContext:   resourceV2WorkflowRead,
		UpdateContext: resourceV2WorkflowUpdate,
		DeleteContext: resourceV2WorkflowDelete,

		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		CustomizeDiff: config.FlexibleForceNew(v2WorkflowNonUpdatableParams),

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: `The region where the workflow is located.`,
			},

			// Required parameters.
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The name of the workflow.`,
			},

			// Optional parameters.
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `The description of the workflow.`,
			},
			"steps": {
				Type:        schema.TypeList,
				Optional:    true,
				Computed:    true,
				Elem:        workflowStepsSchema(),
				Description: `The steps of the workflow.`,
			},
			"workspace_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The workspace ID to which the workflow belongs.`,
			},
			"data_requirements": {
				Type:        schema.TypeList,
				Optional:    true,
				Computed:    true,
				Elem:        workflowDataRequirementsSchema(),
				Description: `The data requirements of the workflow.`,
			},
			"data": {
				Type:        schema.TypeList,
				Optional:    true,
				Elem:        workflowDataSchema(),
				Description: `The data of the workflow.`,
			},
			"parameters": {
				Type:        schema.TypeList,
				Optional:    true,
				Computed:    true,
				Elem:        workflowParametersSchema(),
				Description: `The parameters of the workflow.`,
			},
			"source_workflow_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The source workflow ID for copying.`,
			},
			"gallery_subscription": {
				Type:        schema.TypeList,
				Optional:    true,
				MaxItems:    1,
				Elem:        workflowGallerySubscriptionSchema(),
				Description: `The gallery subscription information of the workflow.`,
			},
			"storages": {
				Type:        schema.TypeList,
				Optional:    true,
				Computed:    true,
				Elem:        workflowStoragesSchema(),
				Description: `The unified storage definitions of the workflow.`,
			},
			"labels": {
				Type:        schema.TypeList,
				Optional:    true,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: `The labels of the workflow.`,
			},
			"assets": {
				Type:        schema.TypeList,
				Optional:    true,
				Elem:        workflowAssetsSchema(),
				Description: `The assets bound to the workflow.`,
			},
			"sub_graphs": {
				Type:        schema.TypeList,
				Optional:    true,
				Elem:        workflowSubGraphsSchema(),
				Description: `The subgraphs of the workflow.`,
			},
			"extend": {
				Type:             schema.TypeString,
				Optional:         true,
				DiffSuppressFunc: compareJsonStringValue,
				Description:      `The extended fields of the billing workflow, in JSON format.`,
			},
			"policy": {
				Type:        schema.TypeList,
				Optional:    true,
				MaxItems:    1,
				Elem:        workflowPolicySchema(),
				Description: `The partial running policy of the workflow.`,
			},
			"with_subscription": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: `Whether to enable the SMN message subscription of the workflow.`,
			},
			"smn_switch": {
				Type:        schema.TypeBool,
				Optional:    true,
				Computed:    true,
				Description: `Whether to enable the SMN switch of the workflow.`,
			},
			"subscription_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The SMN message subscription ID of the workflow.`,
			},
			"exeml_template_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The auto learning template ID of the workflow.`,
			},
			"package": {
				Type:        schema.TypeList,
				Optional:    true,
				MaxItems:    1,
				Elem:        workflowPackageSchema(),
				Description: `The billing workflow subscription package information.`,
			},

			// Attributes.
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
				Elem:        workflowLatestExecutionSchema(),
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

func workflowStepInputsSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `The name of the input data.`,
			},
			"type": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `The type of the input.`,
			},
			"data": {
				Type:             schema.TypeString,
				Optional:         true,
				Computed:         true,
				DiffSuppressFunc: compareJsonStringValue,
				Description:      `The input data, in JSON format.`,
			},
			"value": {
				Type:             schema.TypeString,
				Optional:         true,
				Computed:         true,
				DiffSuppressFunc: compareJsonStringValue,
				Description:      `The value of the input, in JSON format.`,
			},
		},
	}
}

func workflowStepOutputsSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `The name of the output data.`,
			},
			"type": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `The type of the output.`,
			},
			"config": {
				Type:             schema.TypeString,
				Optional:         true,
				Computed:         true,
				DiffSuppressFunc: compareJsonStringValue,
				Description:      `The output configuration, in JSON format.`,
			},
		},
	}
}

func workflowStepConditionsSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"type": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `The condition type.`,
			},
			"left": {
				Type:             schema.TypeString,
				Optional:         true,
				Computed:         true,
				DiffSuppressFunc: compareJsonStringValue,
				Description:      `The left branch when the condition is true, in JSON format.`,
			},
			"right": {
				Type:             schema.TypeString,
				Optional:         true,
				Computed:         true,
				DiffSuppressFunc: compareJsonStringValue,
				Description:      `The right branch when the condition is false, in JSON format.`,
			},
		},
	}
}

func workflowStepPolicySchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"poll_interval_seconds": {
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
				Description: `The execution interval of the workflow step, in seconds.`,
			},
			"max_execution_minutes": {
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
				Description: `The maximum execution time of the workflow step, in minutes.`,
			},
		},
	}
}

func workflowStepsSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The name of the workflow step.`,
			},
			"type": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `The type of the workflow step.`,
			},
			"inputs": {
				Type:        schema.TypeList,
				Optional:    true,
				Computed:    true,
				Elem:        workflowStepInputsSchema(),
				Description: `The inputs of the workflow step.`,
			},
			"outputs": {
				Type:        schema.TypeList,
				Optional:    true,
				Computed:    true,
				Elem:        workflowStepOutputsSchema(),
				Description: `The outputs of the workflow step.`,
			},
			"title": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `The title of the workflow step.`,
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `The description of the workflow step.`,
			},
			"properties": {
				Type:             schema.TypeString,
				Optional:         true,
				Computed:         true,
				DiffSuppressFunc: compareJsonStringValue,
				Description:      `The properties of the workflow step, in JSON format.`,
			},
			"depend_steps": {
				Type:        schema.TypeList,
				Optional:    true,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: `The dependent steps of the workflow step.`,
			},
			"conditions": {
				Type:        schema.TypeList,
				Optional:    true,
				Computed:    true,
				Elem:        workflowStepConditionsSchema(),
				Description: `The execution conditions of the workflow step.`,
			},
			"if_then_steps": {
				Type:        schema.TypeList,
				Optional:    true,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: `The conditional branch steps of the workflow step.`,
			},
			"else_then_steps": {
				Type:        schema.TypeList,
				Optional:    true,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: `The other conditional branch steps of the workflow step.`,
			},
			"policy": {
				Type:        schema.TypeList,
				Optional:    true,
				Computed:    true,
				MaxItems:    1,
				Elem:        workflowStepPolicySchema(),
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

func workflowDataRequirementConditionsSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"attribute": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `The condition attribute.`,
			},
			"operator": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `The operator of the condition.`,
			},
			"value": {
				Type:             schema.TypeString,
				Optional:         true,
				Computed:         true,
				DiffSuppressFunc: compareJsonStringValue,
				Description:      `The value of the condition, in JSON format.`,
			},
		},
	}
}

func workflowDataRequirementsSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The name of the data requirement.`,
			},
			"type": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The type of the data source.`,
			},
			"conditions": {
				Type:        schema.TypeList,
				Optional:    true,
				Computed:    true,
				Elem:        workflowDataRequirementConditionsSchema(),
				Description: `The data constraint conditions.`,
			},
			"value": {
				Type:             schema.TypeString,
				Optional:         true,
				Computed:         true,
				DiffSuppressFunc: compareJsonStringValue,
				Description:      `The value of the data requirement, in JSON format.`,
			},
			"used_steps": {
				Type:        schema.TypeList,
				Optional:    true,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: `The steps that use this data requirement.`,
			},
			"delay": {
				Type:        schema.TypeBool,
				Optional:    true,
				Computed:    true,
				Description: `Whether the data requirement is delayed.`,
			},
		},
	}
}

func workflowDataSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The name of the data.`,
			},
			"type": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The type of the data source.`,
			},
			"value": {
				Type:             schema.TypeString,
				Optional:         true,
				DiffSuppressFunc: compareJsonStringValue,
				Description:      `The value of the data, in JSON format.`,
			},
			"used_steps": {
				Type:        schema.TypeList,
				Optional:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: `The steps that use this data.`,
			},
		},
	}
}

func workflowParametersSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `The name of the workflow parameter.`,
			},
			"type": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `The type of the workflow parameter.`,
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `The description of the workflow parameter.`,
			},
			"example": {
				Type:             schema.TypeString,
				Optional:         true,
				Computed:         true,
				DiffSuppressFunc: compareJsonStringValue,
				Description:      `The example of the workflow parameter, in JSON format.`,
			},
			"delay": {
				Type:        schema.TypeBool,
				Optional:    true,
				Computed:    true,
				Description: `Whether the workflow parameter is delayed.`,
			},
			"default": {
				Type:             schema.TypeString,
				Optional:         true,
				Computed:         true,
				DiffSuppressFunc: compareJsonStringValue,
				Description:      `The default value of the workflow parameter, in JSON format.`,
			},
			"value": {
				Type:             schema.TypeString,
				Optional:         true,
				Computed:         true,
				DiffSuppressFunc: compareJsonStringValue,
				Description:      `The value of the workflow parameter, in JSON format.`,
			},
			"enum": {
				Type:             schema.TypeList,
				Optional:         true,
				Computed:         true,
				Elem:             &schema.Schema{Type: schema.TypeString},
				DiffSuppressFunc: compareJsonStringList,
				Description:      `The enumeration items of the workflow parameter, in JSON format.`,
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
				Description: `The data format of the workflow parameter.`,
			},
			"constraint": {
				Type:             schema.TypeString,
				Optional:         true,
				Computed:         true,
				DiffSuppressFunc: compareJsonStringValue,
				Description:      `The constraint of the workflow parameter, in JSON format.`,
			},
		},
	}
}

func workflowGallerySubscriptionSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"content_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The asset ID of the gallery subscription.`,
			},
			"version_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The version ID of the gallery subscription.`,
			},
			"expired_at": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The expiration time of the gallery subscription, in RFC3339 format.`,
			},
		},
	}
}

func workflowStoragesSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `The name of the workflow storage.`,
			},
			"type": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `The type of the workflow storage.`,
			},
			"path": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `The root path of the unified storage.`,
			},
		},
	}
}

func workflowAssetsSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The name of the asset.`,
			},
			"type": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The type of the asset.`,
			},
			"content_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The asset ID.`,
			},
			"subscription_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The subscription ID of the asset.`,
			},
			"expired_at": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The expiration time of the asset, in RFC3339 format.`,
			},
		},
	}
}

func workflowSubGraphsSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The name of the subgraph.`,
			},
			"steps": {
				Type:        schema.TypeList,
				Optional:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: `The step members of the subgraph.`,
			},
		},
	}
}

func workflowPolicyScenesSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The scene ID.`,
			},
			"name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The scene name.`,
			},
			"steps": {
				Type:        schema.TypeList,
				Optional:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: `The step list of the scene.`,
			},
		},
	}
}

func workflowPolicySchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"use_scene": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The usage scenario of the workflow policy.`,
			},
			"scene_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The scene ID of the workflow policy.`,
			},
			"scenes": {
				Type:        schema.TypeList,
				Optional:    true,
				Elem:        workflowPolicyScenesSchema(),
				Description: `The scenes of the workflow policy.`,
			},
		},
	}
}

func workflowPackageOrderSkuSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"code": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The billing code.`,
			},
			"period": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: `The billing period.`,
			},
			"queries_limit": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: `The query limit.`,
			},
			"price": {
				Type:        schema.TypeFloat,
				Optional:    true,
				Description: `The price.`,
			},
		},
	}
}

func workflowPackageOrderSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The subscription ID.`,
			},
			"sku": {
				Type:        schema.TypeList,
				Required:    true,
				Elem:        workflowPackageOrderSkuSchema(),
				Description: `The subscription billing information.`,
			},
			"sku_count": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: `The subscription count.`,
			},
		},
	}
}

func workflowPackageSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"package_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The resource package UUID.`,
			},
			"pool_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The resource pool ID.`,
			},
			"service_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The service ID.`,
			},
			"workflow_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The workflow ID.`,
			},
			"order": {
				Type:        schema.TypeList,
				Optional:    true,
				Elem:        workflowPackageOrderSchema(),
				Description: `The subscription information.`,
			},
			"consume_limit": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: `The subscription limit.`,
			},
			"current_consume": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: `The current subscription consumption.`,
			},
			"current_date": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The current date.`,
			},
			"limit_enable": {
				Type:        schema.TypeBool,
				Optional:    true,
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

func workflowLatestExecutionSchema() *schema.Resource {
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

func compareJsonStringValue(_, o, n string, _ *schema.ResourceData) bool {
	return utils.JSONStringsEqual(o, n)
}

func compareJsonStringList(k, o, n string, d *schema.ResourceData) bool {
	if o == n {
		return true
	}

	oldCount, newCount := d.GetChange(k + ".#")
	if oldCount != newCount {
		return false
	}
	if oldCount == nil {
		return true
	}

	count := oldCount.(int)
	for i := 0; i < count; i++ {
		oldVal, newVal := d.GetChange(fmt.Sprintf("%s.%d", k, i))
		if oldVal == nil && newVal == nil {
			continue
		}
		if oldVal == nil || newVal == nil {
			return false
		}

		oldStr, ok1 := oldVal.(string)
		newStr, ok2 := newVal.(string)
		if !ok1 || !ok2 {
			return false
		}
		if !utils.JSONStringsEqual(oldStr, newStr) {
			return false
		}
	}
	return true
}

func buildV2WorkflowStepInputs(inputs []interface{}) []map[string]interface{} {
	if len(inputs) < 1 {
		return nil
	}

	result := make([]map[string]interface{}, 0, len(inputs))
	for _, input := range inputs {
		result = append(result, map[string]interface{}{
			"name":  utils.ValueIgnoreEmpty(utils.PathSearch("name", input, "").(string)),
			"type":  utils.ValueIgnoreEmpty(utils.PathSearch("type", input, "").(string)),
			"data":  utils.ValueIgnoreEmpty(utils.StringToJson(utils.PathSearch("data", input, "").(string))),
			"value": utils.ValueIgnoreEmpty(utils.StringToJson(utils.PathSearch("value", input, "").(string))),
		})
	}
	return result
}

func buildV2WorkflowStepOutputs(outputs []interface{}) []map[string]interface{} {
	if len(outputs) < 1 {
		return nil
	}

	result := make([]map[string]interface{}, 0, len(outputs))
	for _, output := range outputs {
		result = append(result, map[string]interface{}{
			"name":   utils.ValueIgnoreEmpty(utils.PathSearch("name", output, "").(string)),
			"type":   utils.ValueIgnoreEmpty(utils.PathSearch("type", output, "").(string)),
			"config": utils.ValueIgnoreEmpty(utils.StringToJson(utils.PathSearch("config", output, "").(string))),
		})
	}
	return result
}

func buildV2WorkflowStepConditions(conditions []interface{}) []map[string]interface{} {
	if len(conditions) < 1 {
		return nil
	}

	result := make([]map[string]interface{}, 0, len(conditions))
	for _, condition := range conditions {
		result = append(result, map[string]interface{}{
			"type":  utils.ValueIgnoreEmpty(utils.PathSearch("type", condition, "").(string)),
			"left":  utils.ValueIgnoreEmpty(utils.StringToJson(utils.PathSearch("left", condition, "").(string))),
			"right": utils.ValueIgnoreEmpty(utils.StringToJson(utils.PathSearch("right", condition, "").(string))),
		})
	}
	return result
}

func buildV2WorkflowStepPolicy(policy interface{}) map[string]interface{} {
	if policy == nil {
		return nil
	}

	return map[string]interface{}{
		"poll_interval_seconds": utils.ValueIgnoreEmpty(utils.PathSearch("poll_interval_seconds", policy, nil)),
		"max_execution_minutes": utils.ValueIgnoreEmpty(utils.PathSearch("max_execution_minutes", policy, nil)),
	}
}

func buildV2WorkflowSteps(steps []interface{}) []map[string]interface{} {
	if len(steps) < 1 {
		return nil
	}

	result := make([]map[string]interface{}, 0, len(steps))
	for _, step := range steps {
		result = append(result, map[string]interface{}{
			"name": utils.PathSearch("name", step, nil),
			"type": utils.ValueIgnoreEmpty(utils.PathSearch("type", step, nil)),
			"inputs": utils.ValueIgnoreEmpty(
				buildV2WorkflowStepInputs(utils.PathSearch("inputs", step, make([]interface{}, 0)).([]interface{}))),
			"outputs": utils.ValueIgnoreEmpty(
				buildV2WorkflowStepOutputs(utils.PathSearch("outputs", step, make([]interface{}, 0)).([]interface{}))),
			"title":       utils.ValueIgnoreEmpty(utils.PathSearch("title", step, nil)),
			"description": utils.ValueIgnoreEmpty(utils.PathSearch("description", step, nil)),
			"properties":  utils.ValueIgnoreEmpty(utils.StringToJson(utils.PathSearch("properties", step, "").(string))),
			"depend_steps": utils.ValueIgnoreEmpty(
				utils.ExpandToStringList(utils.PathSearch("depend_steps", step, make([]interface{}, 0)).([]interface{}))),
			"conditions": utils.ValueIgnoreEmpty(
				buildV2WorkflowStepConditions(utils.PathSearch("conditions", step, make([]interface{}, 0)).([]interface{}))),
			"if_then_steps": utils.ValueIgnoreEmpty(
				utils.ExpandToStringList(utils.PathSearch("if_then_steps", step, make([]interface{}, 0)).([]interface{}))),
			"else_then_steps": utils.ValueIgnoreEmpty(
				utils.ExpandToStringList(utils.PathSearch("else_then_steps", step, make([]interface{}, 0)).([]interface{}))),
			"policy": utils.ValueIgnoreEmpty(
				buildV2WorkflowStepPolicy(utils.PathSearch("policy", step, make([]interface{}, 0)).([]interface{}))),
		})
	}
	return result
}

func buildV2WorkflowDataRequirementConditions(conditions []interface{}) []map[string]interface{} {
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

func buildV2WorkflowDataRequirements(dataRequirements []interface{}) []map[string]interface{} {
	if len(dataRequirements) < 1 {
		return nil
	}

	result := make([]map[string]interface{}, 0, len(dataRequirements))
	for _, req := range dataRequirements {
		result = append(result, map[string]interface{}{
			"name": utils.PathSearch("name", req, "").(string),
			"type": utils.PathSearch("type", req, "").(string),
			"conditions": utils.ValueIgnoreEmpty(
				buildV2WorkflowDataRequirementConditions(utils.PathSearch("conditions", req, make([]interface{}, 0)).([]interface{}))),
			"value": utils.ValueIgnoreEmpty(utils.StringToJson(utils.PathSearch("value", req, "").(string))),
			"used_steps": utils.ValueIgnoreEmpty(
				utils.ExpandToStringList(utils.PathSearch("used_steps", req, make([]interface{}, 0)).([]interface{}))),
			"delay": utils.ValueIgnoreEmpty(utils.PathSearch("delay", req, nil)),
		})
	}
	return result
}

func buildV2WorkflowData(dataList []interface{}) []map[string]interface{} {
	if len(dataList) < 1 {
		return nil
	}

	result := make([]map[string]interface{}, 0, len(dataList))
	for _, data := range dataList {
		result = append(result, map[string]interface{}{
			"name":  utils.ValueIgnoreEmpty(utils.PathSearch("name", data, "").(string)),
			"type":  utils.ValueIgnoreEmpty(utils.PathSearch("type", data, "").(string)),
			"value": utils.ValueIgnoreEmpty(utils.StringToJson(utils.PathSearch("value", data, "").(string))),
			"used_steps": utils.ValueIgnoreEmpty(
				utils.ExpandToStringList(utils.PathSearch("used_steps", data, make([]interface{}, 0)).([]interface{}))),
		})
	}
	return result
}

func buildV2WorkflowParameterEnum(enumList []interface{}) []map[string]interface{} {
	if len(enumList) < 1 {
		return nil
	}

	result := make([]map[string]interface{}, 0, len(enumList))
	for _, enumItem := range enumList {
		result = append(result, utils.StringToJson(enumItem.(string), "").(map[string]interface{}))
	}
	return result
}

func buildV2WorkflowParameters(parameters []interface{}) []map[string]interface{} {
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
				buildV2WorkflowParameterEnum(utils.PathSearch("enum", param, make([]interface{}, 0)).([]interface{}))),
			"used_steps": utils.ValueIgnoreEmpty(
				utils.ExpandToStringList(utils.PathSearch("used_steps", param, make([]interface{}, 0)).([]interface{}))),
			"format":     utils.ValueIgnoreEmpty(utils.PathSearch("format", param, nil)),
			"constraint": utils.ValueIgnoreEmpty(utils.StringToJson(utils.PathSearch("constraint", param, "").(string))),
		})
	}
	return result
}

func buildV2WorkflowGallerySubscription(subscriptions []interface{}) map[string]interface{} {
	if len(subscriptions) < 1 {
		return nil
	}

	subscription := subscriptions[0]
	return map[string]interface{}{
		"content_id": utils.ValueIgnoreEmpty(utils.PathSearch("content_id", subscription, nil)),
		"version_id": utils.ValueIgnoreEmpty(utils.PathSearch("version_id", subscription, nil)),
		"expired_at": utils.ValueIgnoreEmpty(utils.PathSearch("expired_at", subscription, nil)),
	}
}

func buildV2WorkflowStorages(storages []interface{}) []map[string]interface{} {
	if len(storages) < 1 {
		return nil
	}

	result := make([]map[string]interface{}, 0, len(storages))
	for _, storage := range storages {
		result = append(result, map[string]interface{}{
			"name": utils.ValueIgnoreEmpty(utils.PathSearch("name", storage, nil)),
			"type": utils.ValueIgnoreEmpty(utils.PathSearch("type", storage, nil)),
			"path": utils.ValueIgnoreEmpty(utils.PathSearch("path", storage, nil)),
		})
	}
	return result
}

func buildV2WorkflowAssets(assets []interface{}) []map[string]interface{} {
	if len(assets) < 1 {
		return nil
	}

	result := make([]map[string]interface{}, 0, len(assets))
	for _, asset := range assets {
		result = append(result, map[string]interface{}{
			"name":            utils.ValueIgnoreEmpty(utils.PathSearch("name", asset, nil)),
			"type":            utils.ValueIgnoreEmpty(utils.PathSearch("type", asset, nil)),
			"content_id":      utils.ValueIgnoreEmpty(utils.PathSearch("content_id", asset, nil)),
			"subscription_id": utils.ValueIgnoreEmpty(utils.PathSearch("subscription_id", asset, nil)),
			"expired_at":      utils.ValueIgnoreEmpty(utils.PathSearch("expired_at", asset, nil)),
		})
	}
	return result
}

func buildV2WorkflowSubGraphs(subGraphs []interface{}) []map[string]interface{} {
	if len(subGraphs) < 1 {
		return nil
	}

	result := make([]map[string]interface{}, 0, len(subGraphs))
	for _, subGraph := range subGraphs {
		result = append(result, map[string]interface{}{
			"name": utils.ValueIgnoreEmpty(utils.PathSearch("name", subGraph, nil)),
			"steps": utils.ValueIgnoreEmpty(
				utils.ExpandToStringList(utils.PathSearch("steps", subGraph, make([]interface{}, 0)).([]interface{}))),
		})
	}
	return result
}

func buildV2WorkflowPolicyScenes(scenes []interface{}) []map[string]interface{} {
	if len(scenes) < 1 {
		return nil
	}

	result := make([]map[string]interface{}, 0, len(scenes))
	for _, scene := range scenes {
		result = append(result, map[string]interface{}{
			"id":   utils.ValueIgnoreEmpty(utils.PathSearch("id", scene, nil)),
			"name": utils.ValueIgnoreEmpty(utils.PathSearch("name", scene, nil)),
			"steps": utils.ValueIgnoreEmpty(
				utils.ExpandToStringList(utils.PathSearch("steps", scene, make([]interface{}, 0)).([]interface{}))),
		})
	}
	return result
}

func buildV2WorkflowPolicy(policies []interface{}) map[string]interface{} {
	if len(policies) < 1 {
		return nil
	}

	policy := policies[0]
	return map[string]interface{}{
		"use_scene": utils.ValueIgnoreEmpty(utils.PathSearch("use_scene", policy, nil)),
		"scene_id":  utils.ValueIgnoreEmpty(utils.PathSearch("scene_id", policy, nil)),
		"scenes": utils.ValueIgnoreEmpty(
			buildV2WorkflowPolicyScenes(utils.PathSearch("scenes", policy, make([]interface{}, 0)).([]interface{}))),
	}
}

func buildV2WorkflowPackageOrderSku(sku interface{}) map[string]interface{} {
	if sku == nil {
		return nil
	}

	return map[string]interface{}{
		"code":          utils.ValueIgnoreEmpty(utils.PathSearch("code", sku, nil)),
		"period":        utils.ValueIgnoreEmpty(utils.PathSearch("period", sku, nil)),
		"queries_limit": utils.ValueIgnoreEmpty(utils.PathSearch("queries_limit", sku, nil)),
		"price":         utils.ValueIgnoreEmpty(utils.PathSearch("price", sku, nil)),
	}
}

func buildV2WorkflowPackageOrder(order interface{}) map[string]interface{} {
	if order == nil {
		return nil
	}

	return map[string]interface{}{
		"id":        utils.ValueIgnoreEmpty(utils.PathSearch("id", order, nil)),
		"sku":       buildV2WorkflowPackageOrderSku(utils.PathSearch("sku", order, nil)),
		"sku_count": utils.PathSearch("sku_count", order, nil),
	}
}

func buildV2WorkflowPackage(packages []interface{}) []map[string]interface{} {
	if len(packages) < 1 {
		return nil
	}

	result := make([]map[string]interface{}, 0, len(packages))
	for _, packageItem := range packages {
		result = append(result, map[string]interface{}{
			"package_id":      utils.ValueIgnoreEmpty(utils.PathSearch("package_id", packageItem, nil)),
			"pool_id":         utils.ValueIgnoreEmpty(utils.PathSearch("pool_id", packageItem, nil)),
			"service_id":      utils.ValueIgnoreEmpty(utils.PathSearch("service_id", packageItem, nil)),
			"workflow_id":     utils.ValueIgnoreEmpty(utils.PathSearch("workflow_id", packageItem, nil)),
			"order":           utils.ValueIgnoreEmpty(buildV2WorkflowPackageOrder(utils.PathSearch("order", packageItem, nil))),
			"consume_limit":   utils.ValueIgnoreEmpty(utils.PathSearch("consume_limit", packageItem, nil)),
			"current_consume": utils.ValueIgnoreEmpty(utils.PathSearch("current_consume", packageItem, nil)),
			"current_date":    utils.ValueIgnoreEmpty(utils.PathSearch("current_date", packageItem, nil)),
			"limit_enable":    utils.ValueIgnoreEmpty(utils.PathSearch("limit_enable", packageItem, nil)),
		})
	}
	return result
}

func buildV2WorkflowCreateBodyParams(d *schema.ResourceData) map[string]interface{} {
	return map[string]interface{}{
		"name":                 d.Get("name"),
		"description":          utils.ValueIgnoreEmpty(d.Get("description")),
		"steps":                utils.ValueIgnoreEmpty(buildV2WorkflowSteps(d.Get("steps").([]interface{}))),
		"workspace_id":         d.Get("workspace_id"),
		"data_requirements":    utils.ValueIgnoreEmpty(buildV2WorkflowDataRequirements(d.Get("data_requirements").([]interface{}))),
		"data":                 utils.ValueIgnoreEmpty(buildV2WorkflowData(d.Get("data").([]interface{}))),
		"parameters":           utils.ValueIgnoreEmpty(buildV2WorkflowParameters(d.Get("parameters").([]interface{}))),
		"source_workflow_id":   utils.ValueIgnoreEmpty(d.Get("source_workflow_id")),
		"gallery_subscription": utils.ValueIgnoreEmpty(buildV2WorkflowGallerySubscription(d.Get("gallery_subscription").([]interface{}))),
		"storages":             utils.ValueIgnoreEmpty(buildV2WorkflowStorages(d.Get("storages").([]interface{}))),
		"labels":               utils.ValueIgnoreEmpty(utils.ExpandToStringList(d.Get("labels").([]interface{}))),
		"assets":               utils.ValueIgnoreEmpty(buildV2WorkflowAssets(d.Get("assets").([]interface{}))),
		"sub_graphs":           utils.ValueIgnoreEmpty(buildV2WorkflowSubGraphs(d.Get("sub_graphs").([]interface{}))),
		"extend":               utils.ValueIgnoreEmpty(utils.StringToJson(d.Get("extend").(string))),
		"policy":               utils.ValueIgnoreEmpty(buildV2WorkflowPolicy(d.Get("policy").([]interface{}))),
		"with_subscription":    utils.ValueIgnoreEmpty(d.Get("with_subscription")),
		"smn_switch":           utils.ValueIgnoreEmpty(d.Get("smn_switch")),
		"subscription_id":      utils.ValueIgnoreEmpty(d.Get("subscription_id")),
		"exeml_template_id":    utils.ValueIgnoreEmpty(d.Get("exeml_template_id")),
		"package":              utils.ValueIgnoreEmpty(buildV2WorkflowPackage(d.Get("package").([]interface{}))),
	}
}

func createV2Workflow(client *golangsdk.ServiceClient, d *schema.ResourceData) (interface{}, error) {
	httpUrl := "v2/{project_id}/workflows"
	createPath := client.Endpoint + httpUrl
	createPath = strings.ReplaceAll(createPath, "{project_id}", client.ProjectID)

	opt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json",
		},
		JSONBody: utils.RemoveNil(buildV2WorkflowCreateBodyParams(d)),
	}

	requestResp, err := client.Request("POST", createPath, &opt)
	if err != nil {
		return nil, err
	}

	return utils.FlattenResponse(requestResp)
}

func resourceV2WorkflowCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg    = meta.(*config.Config)
		region = cfg.GetRegion(d)
	)

	client, err := cfg.NewServiceClient("modelarts", region)
	if err != nil {
		return diag.Errorf("error creating ModelArts client: %s", err)
	}

	resp, err := createV2Workflow(client, d)
	if err != nil {
		return diag.Errorf("error creating workflow: %s", err)
	}

	workflowId := utils.PathSearch("workflow_id", resp, "").(string)
	if workflowId == "" {
		return diag.Errorf("unable to find the workflow ID from the API response")
	}
	d.SetId(workflowId)

	return resourceV2WorkflowRead(ctx, d, meta)
}

func GetV2WorkflowById(client *golangsdk.ServiceClient, workflowId string) (interface{}, error) {
	httpUrl := "v2/{project_id}/workflows/{workflow_id}"

	getPath := client.Endpoint + httpUrl
	getPath = strings.ReplaceAll(getPath, "{project_id}", client.ProjectID)
	getPath = strings.ReplaceAll(getPath, "{workflow_id}", workflowId)

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

func buildV2WorkflowQueryParams(d *schema.ResourceData) string {
	queryParams := ""

	if v, ok := d.GetOk("workspace_id"); ok {
		queryParams += fmt.Sprintf("%s&workspace_id=%v", queryParams, v)
	}

	if queryParams != "" {
		queryParams = "?" + queryParams[1:]
	}
	return queryParams
}

func queryV2Workflow(client *golangsdk.ServiceClient, d *schema.ResourceData) (interface{}, error) {
	var (
		httpUrl    = "v2/{project_id}/workflows/{workflow_id}"
		workflowId = d.Id()
	)

	getPath := client.Endpoint + httpUrl
	getPath = strings.ReplaceAll(getPath, "{project_id}", client.ProjectID)
	getPath = strings.ReplaceAll(getPath, "{workflow_id}", workflowId)
	getPath += buildV2WorkflowQueryParams(d)

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

func flattenV2WorkflowStepInputs(inputs []interface{}) []map[string]interface{} {
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

func flattenV2WorkflowStepOutputs(outputs []interface{}) []map[string]interface{} {
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

func flattenV2WorkflowStepConditions(conditions []interface{}) []map[string]interface{} {
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

func flattenV2WorkflowStepPolicy(policy interface{}) []map[string]interface{} {
	if policy == nil {
		return nil
	}

	return []map[string]interface{}{
		{
			"poll_interval_seconds": utils.PathSearch("poll_interval_seconds", policy, nil),
			"max_execution_minutes": utils.PathSearch("max_execution_minutes", policy, nil),
		},
	}
}

func flattenV2WorkflowSteps(steps []interface{}) []map[string]interface{} {
	if len(steps) < 1 {
		return nil
	}

	result := make([]map[string]interface{}, 0, len(steps))
	for _, step := range steps {
		result = append(result, map[string]interface{}{
			"name":         utils.PathSearch("name", step, nil),
			"type":         utils.PathSearch("type", step, nil),
			"inputs":       flattenV2WorkflowStepInputs(utils.PathSearch("inputs", step, make([]interface{}, 0)).([]interface{})),
			"outputs":      flattenV2WorkflowStepOutputs(utils.PathSearch("outputs", step, make([]interface{}, 0)).([]interface{})),
			"title":        utils.PathSearch("title", step, nil),
			"description":  utils.PathSearch("description", step, nil),
			"properties":   utils.JsonToString(utils.PathSearch("properties", step, nil)),
			"depend_steps": utils.PathSearch("depend_steps", step, make([]interface{}, 0)).([]interface{}),
			"conditions": flattenV2WorkflowStepConditions(
				utils.PathSearch("conditions", step, make([]interface{}, 0)).([]interface{})),
			"if_then_steps":   utils.PathSearch("if_then_steps", step, make([]interface{}, 0)).([]interface{}),
			"else_then_steps": utils.PathSearch("else_then_steps", step, make([]interface{}, 0)).([]interface{}),
			"policy":          flattenV2WorkflowStepPolicy(utils.PathSearch("policy", step, nil)),
			"created_at": utils.FormatTimeStampRFC3339(
				utils.ConvertTimeStrToNanoTimestamp(utils.PathSearch("created_at", step, "").(string))/1000, true),
		})
	}
	return result
}

func flattenV2WorkflowDataRequirementConditions(conditions []interface{}) []map[string]interface{} {
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

func flattenV2WorkflowDataRequirements(dataRequirements []interface{}) []map[string]interface{} {
	if len(dataRequirements) < 1 {
		return nil
	}

	result := make([]map[string]interface{}, 0, len(dataRequirements))
	for _, req := range dataRequirements {
		result = append(result, map[string]interface{}{
			"name":       utils.PathSearch("name", req, nil),
			"type":       utils.PathSearch("type", req, nil),
			"conditions": flattenV2WorkflowDataRequirementConditions(utils.PathSearch("conditions", req, make([]interface{}, 0)).([]interface{})),
			"value":      utils.JsonToString(utils.PathSearch("value", req, nil)),
			"used_steps": utils.PathSearch("used_steps", req, make([]interface{}, 0)).([]interface{}),
			"delay":      utils.PathSearch("delay", req, nil),
		})
	}
	return result
}

func flattenV2WorkflowParametersEnum(enumList []interface{}) []string {
	if len(enumList) < 1 {
		return nil
	}

	result := make([]string, 0, len(enumList))
	for _, enum := range enumList {
		result = append(result, utils.JsonToString(enum))
	}
	return result
}

func flattenV2WorkflowParameters(parameters []interface{}) []map[string]interface{} {
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
			"enum":        flattenV2WorkflowParametersEnum(utils.PathSearch("enum", param, make([]interface{}, 0)).([]interface{})),
			"used_steps":  utils.PathSearch("used_steps", param, make([]interface{}, 0)).([]interface{}),
			"format":      utils.PathSearch("format", param, nil),
			"constraint":  utils.JsonToString(utils.PathSearch("constraint", param, nil)),
		})
	}
	return result
}

func flattenV2WorkflowStorages(storages []interface{}) []map[string]interface{} {
	if len(storages) < 1 {
		return nil
	}

	result := make([]map[string]interface{}, 0, len(storages))
	for _, storage := range storages {
		result = append(result, map[string]interface{}{
			"name": utils.PathSearch("name", storage, nil),
			"type": utils.PathSearch("type", storage, nil),
			"path": utils.PathSearch("path", storage, nil),
		})
	}
	return result
}

func flattenV2WorkflowLatestExecution(execution interface{}) []map[string]interface{} {
	if execution == nil {
		return nil
	}

	return []map[string]interface{}{
		{
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

func resourceV2WorkflowRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg        = meta.(*config.Config)
		region     = cfg.GetRegion(d)
		workflowId = d.Id()
	)

	client, err := cfg.NewServiceClient("modelarts", region)
	if err != nil {
		return diag.Errorf("error creating ModelArts client: %s", err)
	}

	resp, err := queryV2Workflow(client, d)
	if err != nil {
		return common.CheckDeletedDiag(d, common.ConvertExpected400ErrInto404Err(err, "error_code", v2WorkflowNotFoundErrCodes...),
			fmt.Sprintf("error retrieving workflow (%s)", workflowId))
	}

	mErr := multierror.Append(
		d.Set("region", region),
		d.Set("name", utils.PathSearch("name", resp, nil)),
		d.Set("description", utils.PathSearch("description", resp, nil)),
		d.Set("steps", flattenV2WorkflowSteps(utils.PathSearch("steps", resp, make([]interface{}, 0)).([]interface{}))),
		d.Set("data_requirements", flattenV2WorkflowDataRequirements(
			utils.PathSearch("data_requirements", resp, make([]interface{}, 0)).([]interface{}))),
		d.Set("parameters", flattenV2WorkflowParameters(utils.PathSearch("parameters", resp, make([]interface{}, 0)).([]interface{}))),
		d.Set("storages", flattenV2WorkflowStorages(utils.PathSearch("storages", resp, make([]interface{}, 0)).([]interface{}))),
		d.Set("labels", utils.PathSearch("labels", resp, make([]interface{}, 0)).([]interface{})),
		d.Set("smn_switch", utils.PathSearch("smn_switch", resp, nil)),
		d.Set("created_at", utils.FormatTimeStampRFC3339(
			utils.ConvertTimeStrToNanoTimestamp(utils.PathSearch("created_at", resp, "").(string))/1000, true)),
		d.Set("user_name", utils.PathSearch("user_name", resp, nil)),
		d.Set("latest_execution", flattenV2WorkflowLatestExecution(utils.PathSearch("latest_execution", resp, nil))),
		d.Set("run_count", utils.PathSearch("run_count", resp, nil)),
		d.Set("param_ready", utils.PathSearch("param_ready", resp, nil)),
		d.Set("source", utils.PathSearch("source", resp, nil)),
		d.Set("last_modified_at", utils.FormatTimeStampRFC3339(
			utils.ConvertTimeStrToNanoTimestamp(utils.PathSearch("last_modified_at", resp, "").(string))/1000, true)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func buildV2WorkflowUpdateBodyParams(d *schema.ResourceData) map[string]interface{} {
	return map[string]interface{}{
		"name":              d.Get("name"),
		"description":       utils.ValueIgnoreEmpty(d.Get("description")),
		"steps":             utils.ValueIgnoreEmpty(buildV2WorkflowSteps(d.Get("steps").([]interface{}))),
		"data_requirements": utils.ValueIgnoreEmpty(buildV2WorkflowDataRequirements(d.Get("data_requirements").([]interface{}))),
		"parameters":        utils.ValueIgnoreEmpty(buildV2WorkflowParameters(d.Get("parameters").([]interface{}))),
		"storages":          utils.ValueIgnoreEmpty(buildV2WorkflowStorages(d.Get("storages").([]interface{}))),
		"labels":            utils.ValueIgnoreEmpty(utils.ExpandToStringList(d.Get("labels").([]interface{}))),
		"smn_switch":        utils.ValueIgnoreEmpty(d.Get("smn_switch")),
	}
}

func updateV2Workflow(client *golangsdk.ServiceClient, d *schema.ResourceData) error {
	httpUrl := "v2/{project_id}/workflows/{workflow_id}"
	updatePath := client.Endpoint + httpUrl
	updatePath = strings.ReplaceAll(updatePath, "{project_id}", client.ProjectID)
	updatePath = strings.ReplaceAll(updatePath, "{workflow_id}", d.Id())

	opt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json",
		},
		JSONBody: utils.RemoveNil(buildV2WorkflowUpdateBodyParams(d)),
	}

	_, err := client.Request("PUT", updatePath, &opt)
	return err
}

func resourceV2WorkflowUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg    = meta.(*config.Config)
		region = cfg.GetRegion(d)
	)

	client, err := cfg.NewServiceClient("modelarts", region)
	if err != nil {
		return diag.Errorf("error creating ModelArts client: %s", err)
	}

	err = updateV2Workflow(client, d)
	if err != nil {
		return diag.Errorf("error updating workflow (%s): %s", d.Id(), err)
	}

	return resourceV2WorkflowRead(ctx, d, meta)
}

func deleteV2Workflow(client *golangsdk.ServiceClient, workflowId string) error {
	httpUrl := "v2/{project_id}/workflows/{workflow_id}"
	deletePath := client.Endpoint + httpUrl
	deletePath = strings.ReplaceAll(deletePath, "{project_id}", client.ProjectID)
	deletePath = strings.ReplaceAll(deletePath, "{workflow_id}", workflowId)

	opt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json",
		},
	}

	_, err := client.Request("DELETE", deletePath, &opt)
	return err
}

func resourceV2WorkflowDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg        = meta.(*config.Config)
		region     = cfg.GetRegion(d)
		workflowId = d.Id()
	)

	client, err := cfg.NewServiceClient("modelarts", region)
	if err != nil {
		return diag.Errorf("error creating ModelArts client: %s", err)
	}

	err = deleteV2Workflow(client, workflowId)
	if err != nil {
		return common.CheckDeletedDiag(d, common.ConvertExpected400ErrInto404Err(err, "error_code", v2WorkflowNotFoundErrCodes...),
			fmt.Sprintf("error deleting workflow (%s)", workflowId))
	}

	return nil
}
