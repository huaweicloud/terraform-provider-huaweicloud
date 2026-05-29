---
subcategory: "AI Development Platform (ModelArts)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_modelartsv2_workflow_execution_action"
description: |-
  Use this resource to execute an action on a ModelArts workflow execution within HuaweiCloud.
---

# huaweicloud_modelartsv2_workflow_execution_action

Use this resource to execute an action on a ModelArts workflow execution within HuaweiCloud.

-> This resource is only a one-time action resource for operating the workflow execution. Deleting this resource will
   not clear the corresponding request record, but will only remove the resource information from the tfstate file.

## Example Usage

### Basic Usage

```hcl
variable "workflow_id" {}
variable "workflow_execution_id" {}

resource "huaweicloud_modelartsv2_workflow_execution_action" "test" {
  workflow_id  = var.workflow_id
  execution_id = var.workflow_execution_id
  action_name  = "rerun"
}
```

### Usage with Data Requirements and Parameters

```hcl
variable "workflow_id" {}
variable "workflow_execution_id" {}

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

variable "workflow_execution_steps" {
  type = list(string)
}

resource "huaweicloud_modelartsv2_workflow_execution_action" "test" {
  workflow_id  = var.workflow_id
  execution_id = var.workflow_execution_id
  action_name  = "rerun"

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

  policies {
    rerun_steps = var.workflow_execution_steps
  }
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region where the workflow execution is located.  
  If omitted, the provider-level region will be used.  
  Changing this creates a new resource.

* `workflow_id` - (Required, String, NonUpdatable) Specifies the workflow ID.

* `execution_id` - (Required, String, NonUpdatable) Specifies the workflow execution ID.

* `action_name` - (Required, String, NonUpdatable) Specifies the action name.  
  The valid values are as follows:
  + **stop**: Stops the workflow execution that is currently running.
  + **rerun**: Reruns the workflow execution from specified steps.

* `data_requirements` - (Optional, List, NonUpdatable) Specifies the data requirements used by the workflow steps.  
  The [data_requirements](#modelarts_workflow_execution_action_data_requirements) structure is documented below.

* `parameters` - (Optional, List, NonUpdatable) Specifies the parameters used by the workflow steps.  
  The [parameters](#modelarts_workflow_execution_action_parameters) structure is documented below.

* `policies` - (Optional, List, NonUpdatable) Specifies the execution policies used by the execution record.  
  The [policies](#modelarts_workflow_execution_action_policies) structure is documented below.

<a name="modelarts_workflow_execution_action_data_requirements"></a>
The `data_requirements` block supports:

* `name` - (Required, String, NonUpdatable) Specifies the name of the training data.

* `type` - (Required, String, NonUpdatable) Specifies the type of the training data source.  
  The valid values are as follows:
  + **dataset**: The dataset type.
  + **obs**: The OBS type.
  + **swr**: The SWR type.
  + **model_list**: The model list type.
  + **label_task**: The label task type.
  + **service**: The service type.

* `conditions` - (Optional, List, NonUpdatable) Specifies the constraint conditions of the data.  
  The [conditions](#modelarts_workflow_execution_action_conditions) structure is documented below.

* `value` - (Optional, String, NonUpdatable) Specifies the value of the data, in JSON format.

* `used_steps` - (Optional, List, NonUpdatable) Specifies the steps that use this data.

* `delay` - (Optional, Bool, NonUpdatable) Specifies whether this is a delayed parameter.

<a name="modelarts_workflow_execution_action_conditions"></a>
The `conditions` block supports:

* `attribute` - (Optional, String, NonUpdatable) Specifies the attribute of the constraint.

* `operator` - (Optional, String, NonUpdatable) Specifies the operator of the constraint.

* `value` - (Optional, String, NonUpdatable) Specifies the value of the constraint, in JSON format.

<a name="modelarts_workflow_execution_action_parameters"></a>
The `parameters` block supports:

* `name` - (Optional, String, NonUpdatable) Specifies the name of the parameter.

* `type` - (Optional, String, NonUpdatable) Specifies the type of the parameter.  
  The valid value are as follows:
  + **str**: The string type.
  + **int**: The integer type.
  + **bool**: The boolean type.
  + **float**: The float type.

* `description` - (Optional, String, NonUpdatable) Specifies the description of the parameter.

* `example` - (Optional, String, NonUpdatable) Specifies the example of the parameter, in JSON format.

* `delay` - (Optional, Bool, NonUpdatable) Specifies whether this is a delayed input parameter.

* `default` - (Optional, String, NonUpdatable) Specifies the default value of the parameter, in JSON format.

* `value` - (Optional, String, NonUpdatable) Specifies the value of the parameter, in JSON format.

* `enum` - (Optional, List, NonUpdatable) Specifies the enum values of the parameter, in JSON format.

* `used_steps` - (Optional, List, NonUpdatable) Specifies the steps that use this parameter.

* `format` - (Optional, String, NonUpdatable) Specifies the format of the parameter data.

* `constraint` - (Optional, String, NonUpdatable) Specifies the constraint of the parameter, in JSON format.

<a name="modelarts_workflow_execution_action_policies"></a>
The `policies` block supports:

* `rerun_steps` - (Optional, List, NonUpdatable) Specifies the steps to rerun.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.
