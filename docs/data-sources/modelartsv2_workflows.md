---
subcategory: "ModelArts"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_modelartsv2_workflows"
description: |-
  Use this data source to get the list of ModelArts workflows.
---

# huaweicloud_modelartsv2_workflows

Use this data source to get the list of ModelArts workflows.

## Example Usage

### Query all workflows and without any filter

```hcl
data "huaweicloud_modelartsv2_workflows" "test" {}
```

### Query the workflows and using labels filter

```hcl
variable "workflow_labels" {
  type = list(string)
}

data "huaweicloud_modelartsv2_workflows" "test" {
  labels = var.workflow_labels
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region where the workflows are located.  
  If omitted, the provider-level region will be used.

* `name` - (Optional, String) Specifies the name of the workflow to be queried for fuzzy matching.

* `description` - (Optional, String) Specifies the description of the workflow to be queried.

* `labels` - (Optional, List) Specifies the labels of the workflows to be queried.

* `template_id` - (Optional, String) Specifies the template ID of the workflows to be queried.

* `search_type` - (Optional, String) Specifies the search type of the workflows to be queried.  
  The valid values are as follows:
  + **equal**: Exact match.
  + **contain**: Fuzzy match.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `workflows` - The list of the workflows that matched filter parameters.  
  The [workflows](#modelartsv2_workflows) structure is documented below.

<a name="modelartsv2_workflows"></a>
The `workflows` block supports:

* `id` - The ID of the workflow.

* `name` - The name of the workflow.

* `description` - The description of the workflow.

* `steps` - The steps of the workflow.  
  The [steps](#modelartsv2_workflows_steps) structure is documented below.

* `workspace_id` - The workspace ID to which the workflow belongs.

* `data_requirements` - The data requirements of the workflow.  
  The [data_requirements](#modelartsv2_workflows_data_requirements) structure is documented below.

* `data` - The data of the workflow.  
  The [data](#modelartsv2_workflows_data) structure is documented below.

* `parameters` - The parameters of the workflow.  
  The [parameters](#modelartsv2_workflows_parameters) structure is documented below.

* `source_workflow_id` - The source workflow ID for copying.

* `gallery_subscription` - The gallery subscription information of the workflow.  
  The [gallery_subscription](#modelartsv2_workflows_gallery_subscription) structure is documented below.

* `storages` - The unified storage definitions of the workflow.  
  The [storages](#modelartsv2_workflows_storages) structure is documented below.

* `labels` - The labels of the workflow.

* `assets` - The assets bound to the workflow.  
  The [assets](#modelartsv2_workflows_assets) structure is documented below.

* `sub_graphs` - The subgraphs of the workflow.  
  The [sub_graphs](#modelartsv2_workflows_sub_graphs) structure is documented below.

* `extend` - The extended fields of the billing workflow, in JSON format.

* `policy` - The partial running policy of the workflow.  
  The [policy](#modelartsv2_workflows_policy) structure is documented below.

* `with_subscription` - Whether to enable the SMN message subscription of the workflow.

* `smn_switch` - Whether to enable the SMN switch of the workflow.

* `subscription_id` - The SMN message subscription ID of the workflow.

* `exeml_template_id` - The auto learning template ID of the workflow.

* `package` - The billing workflow subscription package information.  
  The [package](#modelartsv2_workflows_package) structure is documented below.

* `created_at` - The creation time of the workflow, in RFC3339 format.

* `user_name` - The user name that created the workflow.

* `latest_execution` - The latest execution information of the workflow.  
  The [latest_execution](#modelartsv2_workflows_latest_execution) structure is documented below.

* `run_count` - The number of times the workflow has been run.

* `param_ready` - Whether all required parameters of the workflow are filled in.

* `source` - The source of the workflow.

* `last_modified_at` - The last modified time of the workflow, in RFC3339 format.

<a name="modelartsv2_workflows_steps"></a>
The `steps` block supports:

* `name` - The name of the workflow step.

* `type` - The type of the workflow step.

* `inputs` - The inputs of the workflow step.  
  The [inputs](#modelartsv2_workflows_steps_inputs) structure is documented below.

* `outputs` - The outputs of the workflow step.  
  The [outputs](#modelartsv2_workflows_steps_outputs) structure is documented below.

* `title` - The title of the workflow step.

* `description` - The description of the workflow step.

* `properties` - The properties of the workflow step, in JSON format.

* `depend_steps` - The dependent steps of the workflow step.

* `conditions` - The execution conditions of the workflow step.  
  The [conditions](#modelartsv2_workflows_steps_conditions) structure is documented below.

* `if_then_steps` - The conditional branch steps of the workflow step.

* `else_then_steps` - The other conditional branch steps of the workflow step.

* `policy` - The execution policy of the workflow step.  
  The [policy](#modelartsv2_workflows_steps_policy) structure is documented below.

* `created_at` - The creation time of the workflow step, in RFC3339 format.

<a name="modelartsv2_workflows_steps_inputs"></a>
The `inputs` block supports:

* `name` - The name of the input data.

* `type` - The type of the input.

* `data` - The input data, in JSON format.

* `value` - The value of the input, in JSON format.

<a name="modelartsv2_workflows_steps_outputs"></a>
The `outputs` block supports:

* `name` - The name of the output data.

* `type` - The type of the output.

* `config` - The output configuration, in JSON format.

<a name="modelartsv2_workflows_steps_conditions"></a>
The `conditions` block supports:

* `type` - The condition type.

* `left` - The left branch when the condition is true, in JSON format.

* `right` - The right branch when the condition is false, in JSON format.

<a name="modelartsv2_workflows_steps_policy"></a>
The `policy` block supports:

* `poll_interval_seconds` - The execution interval of the workflow step, in seconds.

* `max_execution_minutes` - The maximum execution time of the workflow step, in minutes.

<a name="modelartsv2_workflows_data_requirements"></a>
The `data_requirements` block supports:

* `name` - The name of the data requirement.

* `type` - The type of the data source.

* `conditions` - The data constraint conditions.  
  The [conditions](#modelartsv2_workflows_data_requirements_conditions) structure is documented below.

* `value` - The value of the data requirement, in JSON format.

* `used_steps` - The steps that use this data requirement.

* `delay` - Whether the data requirement is delayed.

<a name="modelartsv2_workflows_data_requirements_conditions"></a>
The `conditions` block supports:

* `attribute` - The condition attribute.

* `operator` - The operator of the condition.

* `value` - The value of the condition, in JSON format.

<a name="modelartsv2_workflows_data"></a>
The `data` block supports:

* `name` - The name of the data.

* `type` - The type of the data source.

* `value` - The value of the data, in JSON format.

* `used_steps` - The steps that use this data.

<a name="modelartsv2_workflows_parameters"></a>
The `parameters` block supports:

* `name` - The name of the workflow parameter.

* `type` - The type of the workflow parameter.

* `description` - The description of the workflow parameter.

* `example` - The example of the workflow parameter, in JSON format.

* `delay` - Whether the workflow parameter is delayed.

* `default` - The default value of the workflow parameter, in JSON format.

* `value` - The value of the workflow parameter, in JSON format.

* `enum` - The enumeration items of the workflow parameter, in JSON format.

* `used_steps` - The steps that use this parameter.

* `format` - The data format of the workflow parameter.

* `constraint` - The constraint of the workflow parameter, in JSON format.

<a name="modelartsv2_workflows_gallery_subscription"></a>
The `gallery_subscription` block supports:

* `content_id` - The asset ID of the gallery subscription.

* `version_id` - The version ID of the gallery subscription.

* `expired_at` - The expiration time of the gallery subscription, in RFC3339 format.

<a name="modelartsv2_workflows_storages"></a>
The `storages` block supports:

* `name` - The name of the workflow storage.

* `type` - The type of the workflow storage.

* `path` - The root path of the unified storage.

<a name="modelartsv2_workflows_assets"></a>
The `assets` block supports:

* `name` - The name of the asset.

* `type` - The type of the asset.

* `content_id` - The asset ID.

* `subscription_id` - The subscription ID of the asset.

* `expired_at` - The expiration time of the asset, in RFC3339 format.

<a name="modelartsv2_workflows_sub_graphs"></a>
The `sub_graphs` block supports:

* `name` - The name of the subgraph.

* `steps` - The step members of the subgraph.

<a name="modelartsv2_workflows_policy"></a>
The `policy` block supports:

* `use_scene` - The usage scenario of the workflow policy.

* `scene_id` - The scene ID of the workflow policy.

* `scenes` - The scenes of the workflow policy.  
  The [scenes](#modelartsv2_workflows_policy_scenes) structure is documented below.

<a name="modelartsv2_workflows_policy_scenes"></a>
The `scenes` block supports:

* `id` - The scene ID.

* `name` - The scene name.

* `steps` - The step list of the scene.

<a name="modelartsv2_workflows_package"></a>
The `package` block supports:

* `package_id` - The resource package UUID.

* `pool_id` - The resource pool ID.

* `service_id` - The service ID.

* `workflow_id` - The workflow ID.

* `order` - The subscription information.  
  The [order](#modelartsv2_workflows_package_order) structure is documented below.

* `consume_limit` - The subscription limit.

* `current_consume` - The current subscription consumption.

* `current_date` - The current date.

* `limit_enable` - Whether the limit is enabled.

* `status` - The status of the resource package.

* `created_at` - The creation time of the resource package, in RFC3339 format.

<a name="modelartsv2_workflows_package_order"></a>
The `order` block supports:

* `id` - The subscription ID.

* `sku` - The subscription billing information.  
  The [sku](#modelartsv2_workflows_package_order_sku) structure is documented below.

* `sku_count` - The subscription count.

<a name="modelartsv2_workflows_package_order_sku"></a>
The `sku` block supports:

* `code` - The billing code.

* `period` - The billing period.

* `queries_limit` - The query limit.

* `price` - The price.

<a name="modelartsv2_workflows_latest_execution"></a>
The `latest_execution` block supports:

* `execution_id` - The execution ID of the workflow.

* `created_at` - The creation time of the workflow execution, in RFC3339 format.

* `status` - The status of the workflow execution.

* `running_steps` - The running steps of the workflow execution.

* `current_steps` - The current steps of the workflow execution.

* `duration` - The duration of the workflow execution.
