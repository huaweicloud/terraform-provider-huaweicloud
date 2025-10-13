---
subcategory: "Host Security Service (HSS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_hss_asset_assign_task"
description: |-
  Use this data source to query the status of a global asset scan task.
---

# huaweicloud_hss_asset_assign_task

Use this data source to query the status of a global asset scan task.

## Example Usage

```hcl
variable "category" {}

data "huaweicloud_hss_asset_assign_task" "test" {
  category = var.category
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used.

* `category` - (Required, String) Specifies the event type.
  The valid values are as follows:
  + **host**: Indicates server security event.
  + **container**: Indicates container security event.

* `enterprise_project_id` - (Optional, String) Specifies the enterprise project ID.
  This parameter is valid only when the enterprise project is enabled.
  The default value is **0**, indicating the default enterprise project.
  If it is necessary to operate the hosts under all enterprise projects, the value is **all_granted_eps**.
  If you only have permissions for a specific enterprise project, you need set the enterprise project ID. Otherwise,
  the operation may fail due to insufficient permissions.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.

* `exist` - Whether a full scan task exists.
  The value can be **true** or **false**.
