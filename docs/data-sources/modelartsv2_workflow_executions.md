---
subcategory: "AI Development Platform (ModelArts)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_modelartsv2_workflow_executions"
description: |-
  Use this data source to get the list of ModelArts workflow executions within HuaweiCloud.
---

# huaweicloud_modelartsv2_workflow_executions

Use this data source to get the list of ModelArts workflow executions within HuaweiCloud.

## Example Usage

### Query all workflow executions by workflow ID

```hcl
variable "workflow_id" {}

data "huaweicloud_modelartsv2_workflow_executions" "test" {
  workflow_id = var.workflow_id
}
```

### Query workflow executions using labels filter

```hcl
variable "workflow_id" {}
variable "workflow_execution_labels" {}

data "huaweicloud_modelartsv2_workflow_executions" "test" {
  workflow_id = var.workflow_id
  labels      = var.workflow_execution_labels
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region where the workflow executions are located.  
  If omitted, the provider-level region will be used.

* `workflow_id` - (Required, String) Specifies the workflow ID of the execution records to be queried.

* `workspace_id` - (Optional, String) Specifies the workspace ID of the execution records to be queried.

* `labels` - (Optional, String) Specifies the labels of the execution records to be queried.

* `status` - (Optional, String) Specifies the status of the execution records to be queried.  
  The valid values are as follows:
  + **init**: Initialization.
  + **wait_inputs**: Waiting for input.
  + **pending**: Pending.
  + **creating**: Creating.
  + **created**: Created.
  + **create_failed**: Create failed.
  + **running**: Running.
  + **stopping**: Stopping.
  + **stopped**: Stopped.
  + **timeout**: Timeout.
  + **completed**: Completed.
  + **failed**: Failed.
  + **hold**: Hold.
  + **skipped**: Skipped.

* `scene_id` - (Optional, String) Specifies the scene ID of the execution records to be queried.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `executions` - The list of the workflow executions that matched filter parameters.  
  The [executions](#modelartsv2_workflow_executions) structure is documented below.

<a name="modelartsv2_workflow_executions"></a>
The `executions` block supports:

* `id` - The ID of the workflow execution.

* `name` - The name of the workflow execution.

* `description` - The description of the workflow execution.

* `status` - The status of the workflow execution.

* `workspace_id` - The workspace ID to which the workflow execution belongs.

* `workflow_id` - The workflow ID.

* `workflow_name` - The workflow name.

* `scene_id` - The custom scene ID.

* `scene_name` - The custom scene name.

* `labels` - The labels of the workflow execution.

* `data_requirements` - The data requirements used by the workflow execution steps.  
  The [data_requirements](#modelartsv2_workflow_executions_data_requirements) structure is documented below.

* `parameters` - The parameters used by the workflow execution steps.  
  The [parameters](#modelartsv2_workflow_executions_parameters) structure is documented below.

* `policies` - The execution policies used by the workflow execution.  
  The [policies](#modelartsv2_workflow_executions_policies) structure is documented below.

* `steps_execution` - The step execution information of the workflow execution.  
  The [steps_execution](#modelartsv2_workflow_executions_steps_execution) structure is documented below.

* `sub_graphs` - The subgraph information of the workflow execution.  
  The [sub_graphs](#modelartsv2_workflow_executions_sub_graphs) structure is documented below.

* `events` - The events of the workflow execution.

* `duration` - The duration of the workflow execution.

* `created_at` - The creation time of the workflow execution, in RFC3339 format.

<a name="modelartsv2_workflow_executions_data_requirements"></a>
The `data_requirements` block supports:

* `name` - The name of the training data.

* `type` - The type of the data source.

* `conditions` - The data constraint conditions.  
  The [conditions](#modelartsv2_workflow_executions_data_requirements_conditions) structure is documented below.

* `value` - The value of the data, in JSON format.

* `used_steps` - The workflow steps that use this data.

* `delay` - Whether the data is delayed.

<a name="modelartsv2_workflow_executions_data_requirements_conditions"></a>
The `conditions` block supports:

* `attribute` - The condition attribute.

* `operator` - The operator of the condition.

* `value` - The value of the condition, in JSON format.

<a name="modelartsv2_workflow_executions_parameters"></a>
The `parameters` block supports:

* `name` - The name of the parameter.

* `type` - The type of the parameter.

* `description` - The description of the parameter.

* `example` - The example of the parameter, in JSON format.

* `delay` - Whether the parameter is delayed.

* `default` - The default value of the parameter, in JSON format.

* `value` - The value of the parameter, in JSON format.

* `enum` - The enumeration values of the parameter, in JSON format.

* `used_steps` - The workflow steps that use this parameter.

* `format` - The data format of the parameter.

* `constraint` - The constraint of the parameter, in JSON format.

<a name="modelartsv2_workflow_executions_policies"></a>
The `policies` block supports:

* `use_cache` - Whether to use cache.

<a name="modelartsv2_workflow_executions_steps_execution"></a>
The `steps_execution` block supports:

* `step_name` - The name of the step.

* `execution_name` - The name of the execution record.

* `name` - The name of the step in this execution.

* `uuid` - The UUID of the step in the execution instance.

* `execution_uuid` - The UUID of the execution instance.

* `created_at` - The creation time of the step execution, in RFC3339 format.

* `updated_at` - The update time of the step execution, in RFC3339 format.

* `duration` - The execution duration of the step, in seconds.

* `type` - The type of the step.
  + **job**: Training.
  + **labeling**: Labeling.
  + **release_dataset**: Dataset release.
  + **model**: Model release.
  + **service**: Service deployment.
  + **mrs_job**: MRS job.
  + **dataset_import**: Dataset import.
  + **create_dataset**: Create dataset.

* `instance_id` - The instance ID of the step execution.

* `status` - The status of the step execution.
  + **init**: Initialization.
  + **wait_inputs**: Waiting for input.
  + **pending**: Pending.
  + **creating**: Creating.
  + **created**: Created.
  + **create_failed**: Create failed.
  + **running**: Running.
  + **stopping**: Stopping.
  + **stopped**: Stopped.
  + **timeout**: Timeout.
  + **completed**: Completed.
  + **failed**: Failed.
  + **hold**: Hold.
  + **skipped**: Skipped.

* `inputs` - The inputs of the step.  
  The [inputs](#modelartsv2_workflow_executions_steps_execution_inputs) structure is documented below.

* `outputs` - The outputs of the step.  
  The [outputs](#modelartsv2_workflow_executions_steps_execution_outputs) structure is documented below.

* `step_uuid` - The UUID of the step.

* `properties` - The properties of the step, in JSON format.

* `events` - The events of the step.

* `error_info` - The error information of the step execution.  
  The [error_info](#modelartsv2_workflow_executions_steps_execution_error_info) structure is documented below.

* `policy` - The execution policy of the step.  
  The [policy](#modelartsv2_workflow_executions_steps_execution_policy) structure is documented below.

* `conditions_execution` - The condition execution of the step.  
  The [conditions_execution](#modelartsv2_workflow_executions_steps_execution_conditions_execution) structure is
  documented below.

* `step_title` - The title of the step.

* `conditions` - The conditions of the step.  
  The [conditions](#modelartsv2_workflow_executions_steps_execution_conditions) structure is documented below.

<a name="modelartsv2_workflow_executions_steps_execution_inputs"></a>
The `inputs` block supports:

* `name` - The name of the input data.

* `type` - The type of the input.
  + **dataset**: Dataset.
  + **obs**: OBS.
  + **data_selector**: Data selector.

* `data` - The input data, in JSON format.

* `value` - The value of the input, in JSON format.

<a name="modelartsv2_workflow_executions_steps_execution_outputs"></a>
The `outputs` block supports:

* `name` - The name of the output data.

* `type` - The type of the output.
  + **obs**: OBS.
  + **model**: AI application meta model.

* `config` - The output configuration, in JSON format.

<a name="modelartsv2_workflow_executions_steps_execution_error_info"></a>
The `error_info` block supports:

* `error_code` - The error code.

* `error_message` - The error message.

<a name="modelartsv2_workflow_executions_steps_execution_policy"></a>
The `policy` block supports:

* `execution_policy` - The execution policy.

* `use_cache` - Whether to use cache.

<a name="modelartsv2_workflow_executions_steps_execution_conditions_execution"></a>
The `conditions_execution` block supports:

* `result` - The execution result, in JSON format.

* `metric_list` - The list of workflow metric information.  
  The [metric_list](#modelartsv2_workflow_executions_steps_execution_conditions_execution_metric_list) structure is
  documented below.

<a name="modelartsv2_workflow_executions_steps_execution_conditions_execution_metric_list"></a>
The `metric_list` block supports:

* `key` - The key of the metric.

* `value` - The value of the metric, in JSON format.

<a name="modelartsv2_workflow_executions_steps_execution_conditions"></a>
The `conditions` block supports:

* `type` - The condition type.
  + **==**
  + **!=**
  + **>**
  + **>=**
  + **<**
  + **<=**
  + **in**
  + **or**

* `left` - The left branch when the condition is true, in JSON format.

* `right` - The right branch when the condition is false, in JSON format.

<a name="modelartsv2_workflow_executions_sub_graphs"></a>
The `sub_graphs` block supports:

* `name` - The name of the subgraph.

* `steps` - The step members of the subgraph.
