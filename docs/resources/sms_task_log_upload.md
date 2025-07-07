---
subcategory: "Server Migration Service (SMS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_sms_task_log_upload"
description: |-
  Manages an SMS task log upload resource within HuaweiCloud.
---

# huaweicloud_sms_task_log_upload

Manages an SMS task log upload resource within HuaweiCloud.

~> Deleting task log upload resource is not supported, it will only be removed from the state.

## Example Usage

```hcl
variable "task_id" {}
variable "log_bucket" {}

resource "huaweicloud_sms_task_log_upload" "test" {
  task_id    = var.task_id
  log_bucket = var.log_bucket
}
```

## Argument Reference

The following arguments are supported:

* `task_id` - (Required, String, NonUpdatable) Specifies the migration task ID.

* `log_bucket` - (Required, String, NonUpdatable) Specifies the bucket name.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID, which equals to `task_id`.
