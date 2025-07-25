---
subcategory: "Distributed Cache Service (DCS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_dcs_migration_task_stop"
description: |-
  Manages a DCS migration task stop resource within HuaweiCloud.
---

# huaweicloud_dcs_migration_task_stop

Manages a DCS migration task stop resource within HuaweiCloud.

## Example Usage

```hcl
variable "task_id" {}

resource "huaweicloud_dcs_migration_task_stop" "test"{
  task_id = var.task_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used. Changing this parameter will create a new resource.

* `task_id` - (Required, String, NonUpdatable) Specifies the ID of the migration task.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID which equals the `task_id`.

## Timeouts

This resource provides the following timeouts configuration options:

* `create` - Default is 30 minutes.
