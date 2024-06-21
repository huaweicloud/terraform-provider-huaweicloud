---
subcategory: "Auto Scaling"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_as_lifecycle_hook_callback"
description: |-
  Manages an AS lifecycle hook callback resource within HuaweiCloud.
---

# huaweicloud_as_lifecycle_hook_callback

Manages an AS lifecycle hook callback resource within HuaweiCloud.

-> 1. Callback action can only be performed when the instance lifecycle hook state is **HANGING**.
<br/>2. The lifecycle hook callback is a one-time action.
<br/>3. Destroying resources does not change the current state of the instance lifecycle hook.

## Example Usage

```hcl
variable "as_group_id" {}
variable "lifecycle_action_result" {}
variable "instance_id" {}
variable "hook_name" {}

resource "huaweicloud_as_lifecycle_hook_callback" "test" {
  scaling_group_id        = var.as_group_id
  lifecycle_action_result = var.lifecycle_action_result
  instance_id             = var.instance_id
  lifecycle_hook_name     = var.hook_name
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the AS lifecycle hook callback.
  If omitted, the provider-level region will be used. Changing this will create a new resource.

* `scaling_group_id` - (Required, String, ForceNew) Specifies the ID of the AS group.
  Changing this will create a new resource.

* `lifecycle_action_result` - (Required, String, ForceNew) Specifies the lifecycle hook callback operation.
  The valid values are as follows:
  + **ABANDON**
  + **CONTINUE**
  + **EXTEND**: Extend the timeout by `1` hour each time.

  Changing this will create a new resource.

* `lifecycle_action_key` - (Optional, String, ForceNew) Specifies the lifecycle hook callback operation token.
  Changing this will create a new resource.

* `instance_id` - (Optional, String, ForceNew) Specifies the instance ID for the lifecycle callback.
  Changing this will create a new resource.

* `lifecycle_hook_name` - (Optional, String, ForceNew) Specifies the lifecycle hook name.
  Changing this will create a new resource.

-> The parameters `instance_id` and `lifecycle_hook_name` must be used together, and they are mutually exclusive with
the parameter `lifecycle_action_key`.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID in UUID format.
