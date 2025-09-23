---
subcategory: "Data Encryption Workshop (DEW)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_csms_scheduled_delete_secret_task"
description: |
  Manages a scheduled delete secret task resource within HuaweiCloud
---

# huaweicloud_csms_scheduled_delete_secret_task

Manages a scheduled delete secret task resource within HuaweiCloud.

-> The current resource is a one-time resource, and destroying this resource will not recover the deleted secret task,
but will only remove the resource information from the tfstate file.

## Example Usage

### Create a scheduled delete task

```hcl
variable "secret_name" {}

resource "huaweicloud_csms_scheduled_delete_secret_task" "test" {
  secret_name = var.secret_name
  action      = "create"
}
```

### Cancel a scheduled delete task

```hcl
variable "secret_name" {}

resource "huaweicloud_csms_scheduled_delete_secret_task" "test" {
  secret_name = var.secret_name
  action      = "cancel"
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the CSMS scheduled delete secret task.
  If omitted, the provider-level region will be used.
  Changing this setting will create a new resource.

* `secret_name` - (Required, String, NonUpdatable) Specifies the secret name.

* `action` - (Required, String) Specifies the action name. The valid values are as follows:  
  + **create**: Create a scheduled delete task.
  + **cancel**: Cancel a scheduled delete task.

* `recovery_window_in_days` - (Optional, Int, NonUpdatable) Specifies the recovery window in days.
  When `action` is set to **create**, this feild will take effect.
  The valid value ranges form `7` to `30`. The default value is `30`.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.
