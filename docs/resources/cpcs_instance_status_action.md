---
subcategory: "Data Encryption Workshop (DEW)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_cpcs_instance_status_action"
description: |-
  Manages the status (enable/disable) of a CPCS instance within HuaweiCloud.
---

# huaweicloud_cpcs_instance_status_action

Manages the status (enable/disable) of a CPCS instance within HuaweiCloud.

-> Currently, this resource is valid only in cn-north-9 region. Destroying this resource will not change the status of
  the CPCS instance, but will only remove the resource information from the tfstate file.

## Example Usage

```hcl
variable "instance_id" {}
variable "action" {}

resource "huaweicloud_cpcs_instance_status_action" "test" {
  instance_id = var.instance_id
  action      = var.action
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to manage the CPCS instance.
  If omitted, the provider-level region will be used. Changing this parameter will create a new resource.

* `instance_id` - (Required, String, NonUpdatable) Specifies the ID of the CPCS instance to manage.

* `action` - (Required, String) Specifies the action to perform on the CPCS instance.
  Valid values are **enable** to enable the instance or **disable** to disable it.

  -> The CPCS instance can only be disabled when it is in the enabled state. Similarly, it can only be enabled when it
    is in the disabled state.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID which is the same as the `instance_id`.
