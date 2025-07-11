---
subcategory: "Workspace"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_workspace_desktop_pool_action"
description: |-
  Use this resource to operate desktop pool within HuaweiCloud.
---

# huaweicloud_workspace_desktop_pool_action

Use this resource to dispatch desktop pool operation within HuaweiCloud.

-> This resource is only a one-time action resource for operate desktop pool. Deleting this resource  
  will not clear the corresponding request record, but will only remove the resource information from the tfstate file.

## Example Usage

```hcl
variable "desktop_pool_id" {}
variable "desktop_pool_op_type" {}
variable "action_type" {}

resource "huaweicloud_workspace_desktop_pool_action" "test" {
  pool_id = var.desktop_pool_id
  op_type = var.desktop_pool_op_type
  type    = var.action_type
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region where the desktop pool is located.  
  If omitted, the provider-level region will be used.
  Changing this parameter will create a new resource.

* `pool_id` - (Required, String, NonUpdatable) Specifies the ID of the desktop pool.

* `op_type` - (Required, String, NonUpdatable) Specifies the type of desktop pool operation.  
  The valid values are as follows:
  + **os-start**
  + **os-stop**
  + **reboot**
  + **hibernate**

* `type` - (Required, String, NonUpdatable) Specifies whether the operation is mandatory.  
  The valid values are as follows:
  + **SOFT**
  + **HARD**

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.

* `status` - The status of operation dispatch task.  
  The valid values are as follows:
  + **SUCCESS**
  + **FAIL**

## Timeouts

This resource provides the following timeouts configuration options:

* `create` - Default is 10 minutes.
