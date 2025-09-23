---
subcategory: "Workspace"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_workspace_desktop_pool_notification"
description: |-
  Use this resource to dispatch desktop pool message within HuaweiCloud.
---

# huaweicloud_workspace_desktop_pool_notification

Use this resource to dispatch desktop pool message within HuaweiCloud.

-> This resource is only a one-time action resource for dispatch desktop pool message. Deleting this resource will not clear
   the corresponding request record, but will only remove the resource information from the tfstate file.

## Example Usage

```hcl
variable "pool_id" {}
variable "notifications" {}

resource "huaweicloud_workspace_desktop_pool_notification" "test" {
  pool_id       = var.pool_id
  notifications = var.notifications
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region where the desktop pool is located.  
  If omitted, the provider-level region will be used.
  Changing this parameter will create a new resource.

* `pool_id` - (Required, String, NonUpdatable) Specifies the ID of the desktop pool.
  Changing this will create a new resource.

* `notifications` - (Required, String, NonUpdatable) Specifies the message want to dispatch.
  Changing this will create a new resource.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.

* `status` - The status of notification dispatch task.  
  The valid values are as follows:
  + **SUCCESS**
  + **FAIL**

## Timeouts

This resource provides the following timeouts configuration options:

* `create` - Default is 3 minutes.
