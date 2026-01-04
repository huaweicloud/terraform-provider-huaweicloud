---
subcategory: "SecMaster"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_secmaster_search_conditions"
description: |-
  Use this data source to query the list of search conditions.
---

# huaweicloud_secmaster_search_conditions

Use this data source to query the list of search conditions.

## Example Usage

```hcl
variable "workspace_id" {}
variable "pipe_id" {}

data "huaweicloud_secmaster_search_conditions" "test" {
  workspace_id = var.workspace_id
  pipe_id      = var.pipe_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

* `workspace_id` - (Required, String) Specifies the workspace ID.

* `pipe_id` - (Required, String) Specifies the pipe ID.

* `sort_key` - (Optional, String) Specifies sorting field.

* `sort_dir` - (Optional, String) Specifies sorting order.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `records` - The search conditions list.

  The [records](#contions_struct) structure is documented below.

<a name="contions_struct"></a>
The `records` block supports:

* `condition_id` - The search condition ID.

* `condition_name` - The search condition name.
