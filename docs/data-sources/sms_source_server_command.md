---
subcategory: "Server Migration Service (SMS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_sms_source_server_command"
description: |-
  Use this data source to obtain commands from the SMS server.
---

# huaweicloud_sms_source_server_command

Use this data source to obtain commands from the SMS server.

## Example Usage

```hcl
variable "server_id" {}

data "huaweicloud_sms_source_server_command" "test" {
  server_id = var.server_id
}
```

## Argument Reference

The following arguments are supported:

* `server_id` - (Required, String) Specifies the ID of the source server that the command is sent to.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `command_name` - Indicates the command name.
  Values can be **START**, **STOP**, **DELETE**, **SYNC** and **SKIP**.

* `command_param` - Indicates the command response parameters.

  The [command_param](#command_param_struct) structure is documented below.

<a name="command_param_struct"></a>
The `command_param` block supports:

* `task_id` - Indicates the task ID.

* `bucket` - Indicates the bucket name.
