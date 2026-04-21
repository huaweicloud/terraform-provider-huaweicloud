---
subcategory: "Resource Formation (RFS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_rfs_stacks"
description: |-
  Use this datasource to get the list of resource stacks.
---

# huaweicloud_rfs_stacks

Use this datasource to get the list of resource stacks.

## Example Usage

```hcl
data "huaweicloud_rfs_stacks" "test" {}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `stacks` - The list of resource stacks. By default, these stacks are sorted in descending order of the
  generation time.

  The [stacks](#stacks_struct) structure is documented below.

<a name="stacks_struct"></a>
The `stacks` block supports:

* `stack_name` - The name of the resource stack.

* `description` - The description of the resource stack.

* `stack_id` - The unique ID of the resource stack.

* `status` - The status of the resource stack.

* `create_time` - The creation time of a resource stack. It is represented in UTC format (yyyy-mm-ddTHH:MM:SSZ),
  such as **1970-01-01T00:00:00Z**.

* `update_time` - The update time of a resource stack. It is represented in UTC format (yyyy-mm-ddTHH:MM:SSZ),
  such as **1970-01-01T00:00:00Z**.

* `status_message` - A brief error summary when the stack status ends with **FAILED**, for debugging.
