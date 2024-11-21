---
subcategory: "Distributed Message Service (DMS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_dms_rocketmq_dead_letter_resend"
description: |-
  Manages a DMS RocketMQ dead letter messages resend resource within HuaweiCloud.
---

# huaweicloud_dms_rocketmq_dead_letter_resend

Manages a DMS RocketMQ dead letter messages resend resource within HuaweiCloud.

## Example Usage

```hcl
variable "instance_id" {}
variable "dead_letter_topic" {}
variable "message_id_list" {}

resource "huaweicloud_dms_rocketmq_dead_letter_resend" "test" {
  instance_id     = var.instance_id
  topic           = var.dead_letter_topic
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
  
* `topic` - (Required, String, ForceNew) Specifies the dead letter topic name.
  Changing this creates a new resource.

* `message_id_list` - (Required, List, ForceNew) Specifies the message ID list.
  Changing this creates a new resource.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.

* `resend_results` - Indicates the resend results.
  The [resend_results](#attrblock--resend_results) structure is documented below.

<a name="attrblock--resend_results"></a>
The `resend_results` block supports:

* `message_id` - Indicates the message ID.

* `error_code` - Indicates the error code.

* `error_message` - Indicates the error message.
