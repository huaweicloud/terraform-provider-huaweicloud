---
subcategory: "Workspace"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_workspace_user_action"
description: |-
  Use this resource to operate Workspace users within HuaweiCloud.
---

# huaweicloud_workspace_user_action

Use this resource to operate Workspace users within HuaweiCloud.

-> This resource is a one-time action resource for operating user. Deleting this resource will not clear
   the corresponding request record, but will only remove the resource information from the tfstate file.

## Example Usage

```hcl
variable "user_id" {}
variable "operation_type" {}

resource "huaweicloud_workspace_user_action" "test" {
  user_id = var.user_id
  op_type = var.operation_type
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region where the user is located.  
  If omitted, the provider-level region will be used. Changing this parameter will create a new resource.

* `user_id` - (Required, String, NonUpdatable) Specifies the ID of the user to be operated.

* `op_type` - (Required, String, NonUpdatable) Specifies the operation type.  
  The valid values are as follows:
  + **LOCK** - Lock the user.
  + **UNLOCK** - Unlock the user.
  + **RESET_PWD** - Reset the user password.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.
