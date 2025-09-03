---
subcategory: "Cloud Operations Center (COC)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_coc_document_execution_detail"
description: |-
  Use this data source to get the COC document execution detail.
---

# huaweicloud_coc_document_execution_detail

Use this data source to get the COC document execution detail.

## Example Usage

```hcl
variable "execution_id" {}

data "huaweicloud_coc_document_execution_detail" "test" {
  execution_id = var.execution_id
}
```

## Argument Reference

The following arguments are supported:

* `execution_id` - (Required, String) Specifies the work order ID.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

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

* `tags` - Indicates the list of work order tags.

  The [tags](#data_tags_struct) structure is documented below.

* `target_parameter_name` - Indicates the rate mode to execute the specified parameter name.

* `targets` - Indicates the rate mode to execute the selection element.

  The [targets](#data_targets_struct) structure is documented below.

* `type` - Indicates the work order type.

<a name="data_parameters_struct"></a>
The `parameters` block supports:

* `key` - Indicates the parameter name.

* `value` - Indicates the parameter value.

<a name="data_tags_struct"></a>
The `tags` block supports:

* `key` - Indicates the key of the work order tag.

* `value` - Indicates the value of the work order tag.

<a name="data_targets_struct"></a>
The `targets` block supports:

* `key` - Indicates the rate mode execution type, **InstanceValues**.

* `values` - Indicates the rate mode execution element. Currently, only cmdb resources are supported.
