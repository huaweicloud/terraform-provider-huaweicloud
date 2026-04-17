---
subcategory: "Resource Formation (RFS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_rfs_stack_resources"
description: |-
  Use this datasource to get the list of resources managed by a resource stack.
---

# huaweicloud_rfs_stack_resources

Use this datasource to get the list of resources managed by a resource stack.

## Example Usage

```hcl
variable "stack_name" {}

data "huaweicloud_rfs_stack_resources" "test" {
  stack_name = var.stack_name
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

* `stack_name` - (Required, String) Specifies the name of the resource stack.

* `stack_id` - (Optional, String) Specifies the unique ID of the resource stack for strong matching.
  If it does not match the current stack, the API returns an error.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `stack_resources` - The list of resources managed by the stack.

  The [stack_resources](#stack_resources_struct) structure is documented below.

<a name="stack_resources_struct"></a>
The `stack_resources` block supports:

* `physical_resource_id` - The physical ID of the resource.

* `physical_resource_name` - The physical name of the resource.

* `logical_resource_name` - The logical name of the resource defined in the template.

* `logical_resource_type` - The logical resource type defined in the template.

* `index_key` - The index of the resource when **count** or **for_each** is used in the template.

* `resource_status` - The status of the resource.

* `status_message` - A brief error summary when the resource is in a failed state.

* `resource_attributes` - The list of resource attributes. It is available when the stack is in a terminal state.

  The [resource_attributes](#stack_resource_attributes_struct) structure is documented below.

<a name="stack_resource_attributes_struct"></a>
The `resource_attributes` block supports:

* `key` - The attribute key.

* `value` - The attribute value.
