---
subcategory: "Cloud Operations Center (COC)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_coc_document_executions"
description: |-
  Use this data source to get the list of COC document executions.
---

# huaweicloud_coc_document_executions

Use this data source to get the list of COC document executions.

## Example Usage

```hcl
data "huaweicloud_coc_document_executions" "test" {}
```

## Argument Reference

The following arguments are supported:

* `creator` - (Optional, String) Specifies the fuzzy query the creator.

* `start_time` - (Optional, Int) Specifies the start time greater than.

* `end_time` - (Optional, Int) Specifies the end time less than.

* `document_name` - (Optional, String) Specifies the fuzzy query the document name.

* `document_id` - (Optional, String) Specifies the document ID.

* `tags` - (Optional, String) Specifies the tag filtering conditions.

* `exclude_child_executions` - (Optional, Bool) Specifies whether list queries should not return child tickets.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `data` - Indicates the list of document execution.

  The [data](#data_struct) structure is documented below.

<a name="data_struct"></a>
The `data` block supports:

* `execution_id` - Indicates the execution ID.

* `document_name` - Indicates the document name.

* `document_id` - Indicates the document ID.

* `document_version_id` - Indicates the document version ID.

* `document_version` - Indicates the document version.

* `start_time` - Indicates the work order execution start time.

* `end_time` - Indicates the work order execution end time.

* `update_time` - Indicates the work order update time.

* `creator` - Indicates the work order creator.

* `status` - Indicates the work order status.

* `description` - Indicates the work order execution description.

* `parameters` - Indicates the global parameters for work order execution.

  The [parameters](#data_parameters_struct) structure is documented below.

* `sys_tags` - Indicates the list of system tags.

  The [sys_tags](#data_sys_tags_struct) structure is documented below.

* `tags` - Indicates the list of custom tags.

  The [tags](#data_tags_struct) structure is documented below.

* `type` - Indicates the work order type.

* `target_parameter_name` - Indicates the rate mode executes the specified parameters.

* `targets` - Indicates the rate mode executes the specified element.

  The [targets](#data_targets_struct) structure is documented below.

<a name="data_parameters_struct"></a>
The `parameters` block supports:

* `key` - Indicates the parameter name.

* `value` - Indicates the parameter value.

<a name="data_sys_tags_struct"></a>
The `sys_tags` block supports:

* `key` - Indicates the key of the tag.

* `value` - Indicates the value of the tag.

<a name="data_tags_struct"></a>
The `tags` block supports:

* `key` - Indicates the key of the tag.

* `value` - Indicates the value of the tag.

<a name="data_targets_struct"></a>
The `targets` block supports:

* `key` - Indicates the rate mode execution type, **InstanceValues**.

* `values` - Indicates the rate mode execution element, currently only supports cmdb resources.
