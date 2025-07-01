---
subcategory: "Server Migration Service (SMS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_sms_source_server_command_result_report"
description: |-
  Manages a SMS source server command result report resource within HuaweiCloud.
---

# huaweicloud_sms_source_server_command_result_report

Manages a SMS source server command result report resource within HuaweiCloud.

~> Deleting source server command result report resource is not supported, it will only be removed from the state.

## Example Usage

```hcl
variable "server_id" {}

resource "huaweicloud_sms_source_server_command_result_report" "test" {
  server_id     = var.server_id
  command_name  = "START"
  result        = "success"
  result_detail = jsonencode({
    "msg": "test"
  })
}
```

## Argument Reference

The following arguments are supported:

* `server_id` - (Required, String, NonUpdatable) Specifies the ID of the source server that the command is sent to.

* `command_name` - (Required, String, NonUpdatable) Specifies the command name.
  Values can be **START**, **STOP**, **DELETE**, **SYNC**, **UPLOAD_LOG** and **RSET_LOG_ACL**.

* `result` - (Required, String, NonUpdatable) Specifies the command execution result.
  Values can be as follows:
  + **success**: The command is executed successfully.
  + **fail**: The command fails to be executed.

* `result_detail` - (Required, String, NonUpdatable) Specifies the command execution results in JSON format.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID, which equals to `server_id`.
