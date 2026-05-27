---
subcategory: "AI Development Platform (ModelArts)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_modelartsv2_workflow_execution"
description: |-
  Use this resource to manage a ModelArts workflow execution resource within HuaweiCloud.
---

# huaweicloud_modelartsv2_workflow_execution

Use this resource to manage a ModelArts workflow execution resource within HuaweiCloud.

## Example Usage

### Basic Usage

```hcl
variable "workflow_id" {}
variable "workflow_execution_name" {}
variable "workflow_execution_description" {}

resource "huaweicloud_modelartsv2_workflow_execution" "test" {
  name        = var.workflow_execution_name
  description = var.workflow_execution_description
  workflow_id = var.workflow_id
}
```

### Usage with Parameters and Data Requirements

```hcl
variable "workflow_id" {}
variable "workflow_execution_name" {}
variable "workflow_execution_description" {}

variable "workflow_execution_data_requirements" {
  type = list(object({
    name       = string
    type       = string
    delay      = optional(bool, false)
    used_steps = list(string)
  }))
}

variable "workflow_execution_parameters" {
  type = list(object({
    name       = string
    type       = string
    default    = optional(string)
    used_steps = list(string)
  }))
}

resource "huaweicloud_modelartsv2_workflow_execution" "test" {
  name        = var.workflow_execution_name
  description = var.workflow_execution_description
  workflow_id = var.workflow_id

  dynamic "parameters" {
    for_each = var.workflow_execution_parameters

    content {
      name       = parameters.value.name
      type       = parameters.value.type
      default    = parameters.value.default
      used_steps = parameters.value.used_steps
    }
  }

  dynamic "data_requirements" {
    for_each = var.workflow_execution_data_requirements

    content {
      name       = data_requirements.value.name
      type       = data_requirements.value.type
      delay      = data_requirements.value.delay
      used_steps = data_requirements.value.used_steps
    }
  }
}
```

### Usage with Policies

```hcl
variable "workflow_id" {}
variable "workflow_execution_name" {}
variable "workflow_execution_description" {}

resource "huaweicloud_modelartsv2_workflow_execution" "test" {
  name        = var.workflow_execution_name
  description = var.workflow_execution_description
  workflow_id = var.workflow_id

  policies {
    use_cache = true
  }
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region where the workflow execution is located.  
  If omitted, the provider-level region will be used. Changing this creates a new resource.

* `name` - (Optional, String, NonUpdatable) Specifies the name of the workflow execution.

* `description` - (Optional, String, NonUpdatable) Specifies the description of the workflow execution.

* `workspace_id` - (Optional, String, NonUpdatable) Specifies the workspace ID to which the workflow execution belongs.

* `workflow_id` - (Optional, String, NonUpdatable) Specifies the workflow ID.

* `workflow_name` - (Optional, String, NonUpdatable) Specifies the workflow name.

* `scene_id` - (Optional, String, NonUpdatable) Specifies the custom scene ID.

* `scene_name` - (Optional, String, NonUpdatable) Specifies the custom scene name.

* `labels` - (Optional, List) Specifies the labels of the workflow execution.

* `data_requirements` - (Optional, List) Specifies the data requirements of the workflow execution.  
  The [data_requirements](#modelarts_workflow_execution_data_requirements) structure is documented below.

* `parameters` - (Optional, List) Specifies the parameters of the workflow execution.  
  The [parameters](#modelarts_workflow_execution_parameters) structure is documented below.

* `policies` - (Optional, List, NonUpdatable) Specifies the policies of the workflow execution.  
  The [policies](#modelarts_workflow_execution_policies) structure is documented below.

<a name="modelarts_workflow_execution_data_requirements"></a>
The `data_requirements` block supports:

* `name` - (Required, String) Specifies the name of the training data.

* `type` - (Required, String) Specifies the type of the training data source.

* `conditions` - (Optional, List) Specifies the constraint conditions of the data.  
  The [conditions](#modelarts_workflow_execution_conditions) structure is documented below.

* `value` - (Optional, String) Specifies the value of the data, in JSON format.

* `used_steps` - (Optional, List) Specifies the steps that use this data.

* `delay` - (Optional, Bool) Specifies whether this is a delayed parameter.

<a name="modelarts_workflow_execution_conditions"></a>
The `conditions` block supports:

* `attribute` - (Optional, String) Specifies the attribute of the constraint.

* `operator` - (Optional, String) Specifies the operator of the constraint.

* `value` - (Optional, String) Specifies the value of the constraint, in JSON format.

<a name="modelarts_workflow_execution_parameters"></a>
The `parameters` block supports:

* `name` - (Optional, String) Specifies the name of the parameter.

* `type` - (Optional, String) Specifies the type of the parameter.

* `description` - (Optional, String) Specifies the description of the parameter.

* `example` - (Optional, String) Specifies the example of the parameter, in JSON format.

* `delay` - (Optional, Bool) Specifies whether this is a delayed input parameter.

* `default` - (Optional, String) Specifies the default value of the parameter, in JSON format.

* `value` - (Optional, String) Specifies the value of the parameter, in JSON format.

* `enum` - (Optional, List) Specifies the enum values of the parameter, in JSON format.

* `used_steps` - (Optional, List) Specifies the steps that use this parameter.

* `format` - (Optional, String) Specifies the format of the parameter data.

* `constraint` - (Optional, String) Specifies the constraint of the parameter, in JSON format.

<a name="modelarts_workflow_execution_policies"></a>
The `policies` block supports:

* `use_cache` - (Optional, Bool, NonUpdatable) Specifies whether to use cache.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.

* `status` - The status of the workflow execution.

* `steps_execution` - The step execution information of the workflow execution.  
  The [steps_execution](#modelarts_workflow_execution_steps_execution) structure is documented below.

* `sub_graphs` - The sub graph information of the workflow execution.  
  The [sub_graphs](#modelarts_workflow_execution_sub_graphs) structure is documented below.

* `events` - The events of the workflow execution.

* `duration` - The duration of the workflow execution.

* `created_at` - The creation time of the workflow execution, in RFC3339 format.

<a name="modelarts_workflow_execution_steps_execution"></a>
The `steps_execution` block supports:

* `step_name` - The name of the step.

* `execution_name` - The name of the execution.

* `name` - The name of the step in this execution.

* `uuid` - The UUID of the step in this execution.

* `execution_uuid` - The UUID of the execution.

* `created_at` - The creation time of the step execution, in RFC3339 format.

* `updated_at` - The update time of the step execution, in RFC3339 format.

* `duration` - The duration of the step execution.

* `type` - The type of the step execution.  
  The valid values are as follows:
  + **job**: Training
  + **labeling**: Labeling
  + **release_dataset**: Dataset release
  + **model**: Model release
  + **service**: Service deployment
  + **mrs_job**: MRS job
  + **dataset_import**: Dataset import
  + **create_dataset**: Create dataset

* `instance_id` - The instance ID of the step execution.

* `status` - The status of the step execution.  
  The valid values are as follows:
  + **init**
  + **wait_inputs**
  + **pending**
  + **creating**
  + **created**
  + **create_failed**
  + **running**
  + **stopping**
  + **stopped**
  + **timeout**
  + **completed**
  + **failed**
  + **hold**
  + **skipped**

* `inputs` - The inputs of the step execution.  
  The [inputs](#modelarts_workflow_execution_inputs) structure is documented below.

* `outputs` - The outputs of the step execution.  
  The [outputs](#modelarts_workflow_execution_outputs) structure is documented below.

* `step_uuid` - The UUID of the step.

* `properties` - The properties of the step, in JSON format.

* `events` - The events of the step execution.

* `error_info` - The error information of the step execution.  
  The [error_info](#modelarts_workflow_execution_error_info) structure is documented below.

* `policy` - The policy of the step execution.  
  The [policy](#modelarts_workflow_execution_step_policy) structure is documented below.

* `conditions_execution` - The conditions execution of the step.  
  The [conditions_execution](#modelarts_workflow_execution_conditions_execution) structure is documented below.

* `step_title` - The title of the step.

* `conditions` - The conditions of the step execution.  
  The [conditions](#modelarts_workflow_execution_step_conditions) structure is documented below.

<a name="modelarts_workflow_execution_inputs"></a>
The `inputs` block supports:

* `name` - The name of the input.

* `type` - The type of the input.  
  The valid values are:
  + **dataset**
  + **obs**
  + **data_selector**

* `data` - The data of the input, in JSON format.

* `value` - The value of the input, in JSON format.

<a name="modelarts_workflow_execution_outputs"></a>
The `outputs` block supports:

* `name` - The name of the output.

* `type` - The type of the output.  
  The valid values are:
  + **obs**: OBS
  + **model**: AI application meta model

* `config` - The config of the output, in JSON format.

<a name="modelarts_workflow_execution_error_info"></a>
The `error_info` block supports:

* `error_code` - The error code.

* `error_message` - The error message.

<a name="modelarts_workflow_execution_step_policy"></a>
The `policy` block supports:

* `execution_policy` - The execution policy.

* `use_cache` - Whether to use cache.

<a name="modelarts_workflow_execution_conditions_execution"></a>
The `conditions_execution` block supports:

* `result` - The result of the condition execution.

* `metric_list` - The metric list of the condition execution.  
  The [metric_list](#modelarts_workflow_execution_metric_list) structure is documented below.

<a name="modelarts_workflow_execution_metric_list"></a>
The `metric_list` block supports:

* `key` - The key of the metric.

* `value` - The value of the metric, in JSON format.

<a name="modelarts_workflow_execution_step_conditions"></a>
The `conditions` block supports:

* `type` - The type of the condition.  
  The valid values are:
  + **==**: Equal
  + **!=**: Not equal
  + **>**: Greater than
  + **>=**: Greater than or equal
  + **<**: Less than
  + **<=**: Less than or equal
  + **in**: In
  + **or**: Or

* `left` - The left value of the condition, in JSON format.

* `right` - The right value of the condition, in JSON format.

<a name="modelarts_workflow_execution_sub_graphs"></a>
The `sub_graphs` block supports:

* `name` - The name of the sub graph.

* `steps` - The steps of the sub graph.

## Import

Workflow executions can be imported using `workflow_id` and `id`, separated by a slash, e.g.

```bash
$ terraform import huaweicloud_modelartsv2_workflow_execution.test <workflow_id>/<id>
```
