---
subcategory: "Distributed Message Service (DMS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_dms_rocketmq_message_send"
description: |-
  Use this resource to send a message to specified DMS RocketMQ topic within HuaweiCloud.
---

# huaweicloud_dms_rocketmq_message_send

Use this resource to send a message to specified DMS RocketMQ topic within HuaweiCloud.

-> This resource is only a one-time action resource for sending message to specified RocketMQ topic. Deleting this
  resource will not clear the corresponding message record, but will only remove the resource information from the
  tfstate file.

## Example Usage

```hcl
variable "instance_id" {}
variable "topic_name" {}
variable "body" {}
variable "property_list" {
  type = list(object({
    name  = string
    value = string
  }))
}

resource "huaweicloud_dms_rocketmq_message_send" "test" {
  instance_id = var.instance_id
  topic       = var.topic_name
  body        = var.body

  dynamic "property_list" {
    for_each = var.property_list
    content {
      name  = property_list.value.name
      value = property_list.value.value
    }
  }
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region where the message to be sent is located.  
  If omitted, the provider-level region will be used. Changing this parameter will create a new resource.

* `instance_id` - (Required, String, NonUpdatable) Specifies the ID of the RocketMQ instance.

* `topic` - (Required, String, NonUpdatable) Specifies the name of the topic to send the message.

* `body` - (Required, String, NonUpdatable) Specifies the content of the message to be sent.

* `property_list` - (Optional, List, NonUpdatable) Specifies the list of message properties.  
  The [property_list](#rocketmq_message_send_property_list) structure is documented below.

<a name="rocketmq_message_send_property_list"></a>
The `property_list` block supports:

* `name` - (Required, String) Specifies the name of the property.

* `value` - (Required, String) Specifies the value of the property.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.

* `msg_id` - The ID of the message that was sent.

* `queue_id` - The queue ID of the message.

* `queue_offset` - The queue offset of the message.

* `broker_name` - The broker name of the message.
