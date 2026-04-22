---
subcategory: "Resource Formation (RFS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_rfs_stack_outputs"
description: |-
  Use this dataSource to get the list of stack outputs.
---

# huaweicloud_rfs_stack_outputs

Use this dataSource to get the list of stack outputs.

## Example Usage

```hcl
variable "stack_name" {}

data "huaweicloud_rfs_stack_outputs" "test" {
  stack_name = var.stack_name
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

* `stack_name` - (Required, String) Specifies the name of the resource stack.

* `stack_id` - (Optional, String) Specifies the ID of the resource stack.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID in UUID format.

* `outputs` - The list of resource stack outputs.

  The [outputs](#outputs_struct) structure is documented below.

<a name="outputs_struct"></a>
The `outputs` block supports:

* `name` - The name of the stack output, as defined in the template.

* `description` - The description of the stack output, as defined in the template.

* `type` - The type of the stack output.

* `value` - The value of the stack output.

* `sensitive` - Whether the stack output is sensitive, as defined in the template.
  If the user defines the output as sensitive in the template, the `value` and `type` of the output in the return body
  will not return the true value, but will return `<sensitive>`.
