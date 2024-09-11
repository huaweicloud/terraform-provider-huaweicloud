---
subcategory: "Distributed Message Service (DMS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_dms_rabbitmq_queues"
description: |-
  Use this data source to get the list of DMS RabbitMQ queues.
---

# huaweicloud_dms_rabbitmq_queues

Use this data source to get the list of DMS RabbitMQ queues.

## Example Usage

```hcl
variable "instance_id" {}
variable "vhost" {}

data "huaweicloud_dms_rabbitmq_queues" "test" {
  instance_id = var.instance_id
  vhost       = var.vhost
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the data source.
  If omitted, the provider-level region will be used.

* `instance_id` - (Required, String) Specifies the DMS RabbitMQ instance ID.

* `vhost` - (Required, String) Specifies the vhost name.

  -> If `vhost` has slashes, please change them into **\_\_F_SLASH\_\_**.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `queues` - Indicates the list of queues.
  The [queues](#attrblock--queues) structure is documented below.

<a name="attrblock--queues"></a>
The `queues` block supports:

* `name` - Indicates the queue name.

* `auto_delete` - Indicates whether the auto delete is enabled.

* `durable` - Indicates whether the durable is enabled.

* `dead_letter_exchange` - Indicates the name of the dead letter exchange.

* `dead_letter_routing_key` - Indicates the routing key of the dead letter exchange.

* `lazy_mode` - Indicates the lazy mode.

* `message_ttl` - Indicates how long a message in this queue can be retained.

* `consumers` - Indicates the connected consumers.

* `messages` - Indicates the accumulated messages.

* `policy` - Indicates the policy.
