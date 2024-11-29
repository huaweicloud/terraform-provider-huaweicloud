---
subcategory: "Distributed Message Service (DMS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_dms_rocketmq_consumption_verify"
description: |-
  Manages a DMS RocketMQ consumption verify resource within HuaweiCloud.
---

# huaweicloud_dms_rocketmq_consumption_verify

Manages a DMS RocketMQ consumption verify resource within HuaweiCloud.

## Example Usage

```hcl
variable "instance_id" {}
variable "group" {}
variable "topic" {}
variable "client_id" {}
variable "message_id_list" {}

resource "huaweicloud_dms_rocketmq_consumption_verify" "test" {
  instance_id     = var.instance_id
  group           = var.group
  topic           = var.topic
  client_id       = var.client_id
  message_id_list = var.message_id_list
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used.
  Changing this creates a new resource.

* `instance_id` - (Required, String, ForceNew) Specifies the instance ID.
  Changing this creates a new resource.

* `group` - (Optional, String, ForceNew) Specifies the group name.
  Changing this creates a new resource.

* `topic` - (Optional, String, ForceNew) Specifies the topic name.
  Changing this creates a new resource.

* `client_id` - (Optional, String, ForceNew) Specifies the client ID.
  Changing this creates a new resource.

* `message_id_list` - (Optional, List, ForceNew) Specifies the message ID list.
  Changing this creates a new resource.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.

* `resend_results` - Indicates the verify results.
  The [resend_results](#attrblock--resend_results) structure is documented below.

<a name="attrblock--resend_results"></a>
The `resend_results` block supports:

* `error_code` - Indicates the error code.

* `error_message` - Indicates the error message.

* `message_id` - Indicates the message ID.
