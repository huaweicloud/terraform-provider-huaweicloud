---
subcategory: "AI Development Platform (ModelArts)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_modelartsv2_workflow"
description: |-
  Use this resource to manage a ModelArts workflow within HuaweiCloud.
---

# huaweicloud_modelartsv2_workflow

Use this resource to manage a ModelArts workflow within HuaweiCloud.

## Example Usage

### Create a workflow with dataset step

```hcl
variable "workflow_name" {}
variable "workflow_description" {}
variable "workflow_step_input_data" {}
variable "train_evaluate_sample_ratio" {}

variable "workflow_data_requirements" {
  type = list(object({
    name       = string
    type       = string
    delay      = optional(bool, false)
    used_steps = list(string)
  }))
}

variable "workflow_parameters" {
  type = list(object({
    name       = string
    type       = string
    default    = optional(string)
    used_steps = list(string)
  }))
}

variable "workflow_storages" {
  type = list(object({
    name        = string
    type        = string
    title       = optional(string)
    description = optional(string)
  }))
}

variable "workflow_policy_scenes" {
  type = list(object({
    id    = optional(string)
    name  = string
    steps = list(string)
  }))
}

resource "huaweicloud_modelartsv2_workflow" "test" {
  name        = var.workflow_name
  description = var.workflow_description

  steps {
    name  = "dataset_step"
    type  = "release_dataset"
    title = "dataset release"

    inputs {
      name = "input_name"
      data = jsonencode(var.workflow_step_input_data)
      type = "dataset"
    }
    outputs {
      name = "output_name"
      type = "dataset"
    }
    properties = jsonencode({
      version_format              = "Default"
      train_evaluate_sample_ratio = var.train_evaluate_sample_ratio
      clear_hard_property         = true
      remove_sample_usage         = true
      label_task_type             = 0
    })

    depend_steps = []
  }

  dynamic "data_requirements" {
    for_each = var.workflow_data_requirements

    content {
      name       = data_requirements.value.name
      type       = data_requirements.value.type
      delay      = data_requirements.value.delay
      used_steps = data_requirements.value.used_steps
    }
  }

  dynamic "parameters" {
    for_each = var.workflow_parameters

    content {
      name       = parameters.value.name
      type       = parameters.value.type
      default    = parameters.value.default
      used_steps = parameters.value.used_steps
    }
  }

  dynamic "storages" {
    for_each = var.workflow_storages

    content {
      name        = storages.value.name
      type        = storages.value.type
      title       = storages.value.title
      description = storages.value.description
    }
  }

  policy {
    dynamic "scenes" {
      for_each = var.workflow_policy_scenes

      content {
        id    = scenes.value.id
        name  = scenes.value.name
        steps = scenes.value.steps
      }
    }
  }
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region where the workflow is located.  
  If omitted, the provider-level region will be used. Changing this creates a new resource.

* `name` - (Required, String) Specifies the name of the workflow.  
  The name consists of `1` to `64` characters, including Chinese and English letters, digits, spaces,
  underscores (_), and hyphens (-), and must start with a Chinese or English letter.

* `description` - (Optional, String) Specifies the description of the workflow.

* `steps` - (Optional, List) Specifies the steps of the workflow.  
  The [steps](#modelartsv2_workflow_steps) structure is documented below.

* `workspace_id` - (Optional, String, NonUpdatable) Specifies the workspace ID to which the workflow belongs.
  Defaults to "0".

* `data_requirements` - (Optional, List) Specifies the data requirements of the workflow.  
  The [data_requirements](#modelartsv2_workflow_data_requirements) structure is documented below.

* `data` - (Optional, List, NonUpdatable) Specifies the data of the workflow.  
  The [data](#modelartsv2_workflow_data) structure is documented below.

* `parameters` - (Optional, List) Specifies the parameters of the workflow.  
  The [parameters](#modelartsv2_workflow_parameters) structure is documented below.

* `source_workflow_id` - (Optional, String, NonUpdatable) Specifies the source workflow ID for copying.

* `gallery_subscription` - (Optional, List, NonUpdatable) Specifies the gallery subscription information of
  the workflow.  
  The [gallery_subscription](#modelartsv2_workflow_gallery_subscription) structure is documented below.

* `storages` - (Optional, List) Specifies the unified storage definitions of the workflow.  
  The [storages](#modelartsv2_workflow_storages) structure is documented below.

* `labels` - (Optional, List) Specifies the labels of the workflow.

* `assets` - (Optional, List, NonUpdatable) Specifies the assets bound to the workflow.  
  The [assets](#modelartsv2_workflow_assets) structure is documented below.

* `sub_graphs` - (Optional, List, NonUpdatable) Specifies the subgraphs of the workflow.  
  The [sub_graphs](#modelartsv2_workflow_sub_graphs) structure is documented below.

* `extend` - (Optional, String, NonUpdatable) Specifies the extended fields of the billing workflow, in JSON format.

* `policy` - (Optional, List) Specifies the partial running policy of the workflow.  
  The [policy](#modelartsv2_workflow_policy) structure is documented below.

* `with_subscription` - (Optional, Bool, NonUpdatable) Specifies whether to enable the SMN message subscription of
  the workflow.  
  Defaults to false.

* `smn_switch` - (Optional, Bool) Specifies whether to enable the SMN switch of the workflow.

* `subscription_id` - (Optional, String, NonUpdatable) Specifies the SMN message subscription ID of the workflow.

* `exeml_template_id` - (Optional, String, NonUpdatable) Specifies the auto learning template ID of the workflow.

* `package` - (Optional, List, NonUpdatable) Specifies the billing workflow subscription package information.  
  The [package](#modelartsv2_workflow_package) structure is documented below.

<a name="modelartsv2_workflow_steps"></a>
The `steps` block supports:

* `name` - (Required, String) Specifies the name of the workflow step.

* `type` - (Optional, String) Specifies the type of the workflow step.  
  The valid values are follows:
  + **job**: The job type of the workflow step.
  + **labeling**: The labeling type of the workflow step.
  + **release_dataset**: The release dataset type of the workflow step.
  + **model**: The model type of the workflow step.
  + **service**: The service type of the workflow step.
  + **mrs_job**: The MRS job type of the workflow step.
  + **dataset_import**: The dataset import type of the workflow step.
  + **create_dataset**: The create dataset type of the workflow step.

* `inputs` - (Optional, List) Specifies the inputs of the workflow step.  
  The [inputs](#modelartsv2_workflow_step_inputs) structure is documented below.

* `outputs` - (Optional, List) Specifies the outputs of the workflow step.  
  The [outputs](#modelartsv2_workflow_step_outputs) structure is documented below.

* `title` - (Optional, String) Specifies the title of the workflow step.

* `description` - (Optional, String) Specifies the description of the workflow step.

* `properties` - (Optional, String) Specifies the properties of the workflow step, in JSON format.

* `depend_steps` - (Optional, List) Specifies the dependent steps of the workflow step.

* `conditions` - (Optional, List) Specifies the execution conditions of the workflow step.  
  The [conditions](#modelartsv2_workflow_step_conditions) structure is documented below.

* `if_then_steps` - (Optional, List) Specifies the conditional branch steps of the workflow step.

* `else_then_steps` - (Optional, List) Specifies the other conditional branch steps of the workflow step.

* `policy` - (Optional, List) Specifies the execution policy of the workflow step.  
  The [step_policy](#modelartsv2_workflow_step_policy) structure is documented below.

<a name="modelartsv2_workflow_step_inputs"></a>
The `inputs` block supports:

* `name` - (Optional, String) Specifies the name of the input data.

* `type` - (Optional, String) Specifies the type of the input.  
  The valid values are as follows:
  + **dataset**
  + **obs**
  + **data_selector**

* `data` - (Optional, String) Specifies the input data, in JSON format.

* `value` - (Optional, String) Specifies the value of the input, in JSON format.

<a name="modelartsv2_workflow_step_outputs"></a>
The `outputs` block supports:

* `name` - (Optional, String) Specifies the name of the output data.

* `type` - (Optional, String) Specifies the type of the output.  
  The valid values are as follows:
  + **obs**: The OBS type.
  + **model**: The model type.

* `config` - (Optional, String) Specifies the output configuration, in JSON format.

<a name="modelartsv2_workflow_step_conditions"></a>
The `conditions` block supports:

* `type` - (Optional, String) Specifies the condition type.  
  The valid values are as follows:
  + **==**: The equal operator.
  + **!=**: The not equal operator.
  + **>**: The greater than operator.
  + **>=**: The greater than or equal operator.
  + **<**: The less than operator.
  + **<=**: The less than or equal operator.
  + **in**: The in operator.
  + **or**: The or operator.

* `left` - (Optional, String) Specifies the left branch when the condition is true, in JSON format.

* `right` - (Optional, String) Specifies the right branch when the condition is false, in JSON format.

<a name="modelartsv2_workflow_step_policy"></a>
The `policy` block supports:

* `poll_interval_seconds` - (Optional, Int) Specifies the execution interval of the workflow step, in seconds.

* `max_execution_minutes` - (Optional, Int) Specifies the maximum execution time of the workflow step, in minutes.

<a name="modelartsv2_workflow_data_requirements"></a>
The `data_requirements` block supports:

* `name` - (Required, String) Specifies the name of the data requirement.

* `type` - (Required, String) Specifies the type of the data source.  
  The valid values are as follows:
  + **dataset**: The dataset type.
  + **obs**: The OBS type.
  + **swr**: The SWR type.
  + **model_list**: The model list type.
  + **label_task**: The label task type.
  + **service**: The service type.

* `conditions` - (Optional, List) Specifies the data constraint conditions.  
  The [conditions](#modelartsv2_workflow_data_requirement_conditions) structure is documented below.

* `value` - (Optional, String) Specifies the value of the data requirement, in JSON format.

* `used_steps` - (Optional, List) Specifies the steps that use this data requirement.

* `delay` - (Optional, Bool) Specifies whether the data requirement is delayed.

<a name="modelartsv2_workflow_data_requirement_conditions"></a>
The `conditions` block supports:

* `attribute` - (Optional, String) Specifies the condition attribute.

* `operator` - (Optional, String) Specifies the operator of the condition.

* `value` - (Optional, String) Specifies the value of the condition, in JSON format.

<a name="modelartsv2_workflow_data"></a>
The `data` block supports:

* `name` - (Optional, String) Specifies the name of the data.

* `type` - (Optional, String) Specifies the type of the data source.  
  The valid values are as follows:
  + **dataset**: The dataset type.
  + **obs**: The OBS type.
  + **swr**: The SWR type.
  + **model**: The model type.
  + **label_task**: The label task type.
  + **service**: The service type.
  + **image**: The image type.

* `value` - (Optional, String) Specifies the value of the data, in JSON format.

* `used_steps` - (Optional, List) Specifies the steps that use this data.

<a name="modelartsv2_workflow_parameters"></a>
The `parameters` block supports:

* `name` - (Optional, String) Specifies the name of the workflow parameter.

* `type` - (Optional, String) Specifies the type of the workflow parameter.  
  The valid value are as follows:
  + **str**: The string type.
  + **int**: The integer type.
  + **bool**: The boolean type.
  + **float**: The float type.

* `description` - (Optional, String) Specifies the description of the workflow parameter.

* `example` - (Optional, String) Specifies the example of the workflow parameter, in JSON format.

* `delay` - (Optional, Bool) Specifies whether the workflow parameter is delayed.

* `default` - (Optional, String) Specifies the default value of the workflow parameter, in JSON format.

* `value` - (Optional, String) Specifies the value of the workflow parameter, in JSON format.

* `enum` - (Optional, List) Specifies the enumeration items of the workflow parameter, in JSON format.

* `used_steps` - (Optional, List) Specifies the steps that use this parameter.

* `format` - (Optional, String) Specifies the data format of the workflow parameter.

* `constraint` - (Optional, String) Specifies the constraint of the workflow parameter, in JSON format.

<a name="modelartsv2_workflow_gallery_subscription"></a>
The `gallery_subscription` block supports:

* `content_id` - (Optional, String, NonUpdatable) Specifies the asset ID of the gallery subscription.

* `version_id` - (Optional, String, NonUpdatable) Specifies the version ID of the gallery subscription.

* `expired_at` - (Optional, String, NonUpdatable) Specifies the expiration time of the gallery subscription,
  in RFC3339 format.

<a name="modelartsv2_workflow_storages"></a>
The `storages` block supports:

* `name` - (Optional, String) Specifies the name of the workflow storage.

* `type` - (Optional, String) Specifies the type of the workflow storage.

* `path` - (Optional, String) Specifies the root path of the unified storage.

<a name="modelartsv2_workflow_assets"></a>
The `assets` block supports:

* `name` - (Optional, String) Specifies the name of the asset.

* `type` - (Optional, String) Specifies the type of the asset.  
  The valid values are as follows:
  + **algorithm**: The algorithm asset.
  + **algorithm2**: The algorithm asset.
  + **model**: The model asset.

* `content_id` - (Optional, String) Specifies the asset ID.

* `subscription_id` - (Optional, String) Specifies the subscription ID of the asset.

* `expired_at` - (Optional, String) Specifies the expiration time of the asset, in RFC3339 format.

<a name="modelartsv2_workflow_sub_graphs"></a>
The `sub_graphs` block supports:

* `name` - (Optional, String) Specifies the name of the subgraph.

* `steps` - (Optional, List) Specifies the step members of the subgraph.

<a name="modelartsv2_workflow_policy"></a>
The `policy` block supports:

* `use_scene` - (Optional, String) Specifies the usage scenario of the workflow policy.

* `scene_id` - (Optional, String) Specifies the scene ID of the workflow policy.

* `scenes` - (Optional, List) Specifies the scenes of the workflow policy.
  The [scenes](#modelartsv2_workflow_policy_scenes) structure is documented below.

<a name="modelartsv2_workflow_policy_scenes"></a>
The `scenes` block supports:

* `id` - (Optional, String) Specifies the scene ID.

* `name` - (Optional, String) Specifies the scene name.

* `steps` - (Optional, List) Specifies the step list of the scene.

<a name="modelartsv2_workflow_package"></a>
The `package` block supports:

* `package_id` - (Optional, String, NonUpdatable) Specifies the resource package UUID.

* `pool_id` - (Optional, String, NonUpdatable) Specifies the resource pool ID.

* `service_id` - (Optional, String, NonUpdatable) Specifies the service ID.

* `workflow_id` - (Optional, String, NonUpdatable) Specifies the workflow ID.

* `order` - (Optional, List, NonUpdatable) Specifies the subscription information.  
  The [order](#modelartsv2_workflow_package_order) structure is documented below.

* `consume_limit` - (Optional, Int, NonUpdatable) Specifies the subscription limit.

* `current_consume` - (Optional, Int, NonUpdatable) Specifies the current subscription consumption.

* `current_date` - (Optional, String, NonUpdatable) Specifies the current date.

* `limit_enable` - (Optional, Bool, NonUpdatable) Specifies whether the limit is enabled.

<a name="modelartsv2_workflow_package_order"></a>
The `order` block supports:

* `id` - (Optional, String, NonUpdatable) Specifies the subscription ID.

* `sku` - (Required, List, NonUpdatable) Specifies the subscription billing information.  
  The [sku](#modelartsv2_workflow_package_order_sku) structure is documented below.

* `sku_count` - (Required, Int, NonUpdatable) Specifies the subscription count.

<a name="modelartsv2_workflow_package_order_sku"></a>
The `sku` block supports:

* `code` - (Optional, String, NonUpdatable) Specifies the billing code.

* `period` - (Optional, Int, NonUpdatable) Specifies the billing period.

* `queries_limit` - (Optional, Int, NonUpdatable) Specifies the query limit.

* `price` - (Optional, Float, NonUpdatable) Specifies the price.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.

* `steps` - The steps of the workflow.  
  The [steps](#modelartsv2_workflow_step_attr) structure is documented below.

* `package` - The billing workflow subscription package information.  
  The [package](#modelartsv2_workflow_package_attr) structure is documented below.

* `created_at` - The creation time of the workflow, in RFC3339 format.

* `user_name` - The user name that created the workflow.

* `latest_execution` - The latest execution information of the workflow.  
  The [latest_execution](#modelartsv2_workflow_latest_execution_attr) structure is documented below.

* `run_count` - The number of times the workflow has been run.

* `param_ready` - Whether all required parameters of the workflow are filled in.

* `source` - The source of the workflow.

* `last_modified_at` - The last modified time of the workflow, in RFC3339 format.

<a name="modelartsv2_workflow_step_attr"></a>
The `steps` block supports:

* `created_at` - The creation time of the workflow step, in RFC3339 format.

<a name="modelartsv2_workflow_package_attr"></a>
The `package` block supports:

* `status` - The status of the resource package.

* `created_at` - The creation time of the resource package, in RFC3339 format.

<a name="modelartsv2_workflow_latest_execution_attr"></a>
The `latest_execution` block supports:

* `execution_id` - The execution ID of the workflow.

* `created_at` - The creation time of the workflow execution, in RFC3339 format.

* `status` - The status of the workflow execution.

* `running_steps` - The running steps of the workflow execution.

* `current_steps` - The current steps of the workflow execution.

* `duration` - The duration of the workflow execution.

## Import

The ModelArts workflow resource can be imported using the 'id', e.g.

```bash
terraform import huaweicloud_modelartsv2_workflow.test <id>
```

Note that the imported state may not be identical to your resource definition, due to some attributes missing from the
API response, security or some other reason. The missing attributes include: `workspace_id`, `data`,`source_workflow_id`,
`gallery_subscription`, `assets`, `sub_graphs`, `extend`, `policy`, `with_subscription`, `subscription_id`,
`exeml_template_id`, `package`. It is generally recommended running `terraform plan` after importing a ModelArts workflow.
You can then decide if changes should be applied to the ModelArts workflow, or the resource definition
should be updated to align with the ModelArts workflow. Also you can ignore changes as below.

```hcl
resource "huaweicloud_modelartsv2_workflow" "test" {
  ...

  lifecycle {
    ignore_changes = [
      workspace_id,
      data,
      source_workflow_id,
      gallery_subscription,
      assets,
      sub_graphs,
      extend,
      policy,
      with_subscription,
      subscription_id,
      exeml_template_id,
      package,
    ]
  }
}
