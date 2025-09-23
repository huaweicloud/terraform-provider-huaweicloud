---
subcategory: "Distributed Message Service (DMS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_dms_rabbitmq_queue"
description: |-
  Manages a DMS RabbitMQ queue resource within HuaweiCloud.
---

# huaweicloud_dms_rabbitmq_queue

Manages a DMS RabbitMQ queue resource within HuaweiCloud.

## Example Usage

```hcl
variable "instance_id" {}
variable "vhost" {}
variable "name" {}
variable "exchange" {}
variable "routing_key" {}

resource "huaweicloud_dms_rabbitmq_queue" "test" {
  instance_id             = var.instance_id
  vhost                   = var.vhost
  name                    = var.name
  auto_delete             = false
  durable                 = true
  dead_letter_exchange    = var.exchange
  dead_letter_routing_key = var.routing_key
  message_ttl             = 4
  lazy_mode               = "lazy"
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used.
  Changing this creates a new resource.

* `instance_id` - (Required, String, ForceNew) Specifies the DMS RabbitMQ instance ID.
  Changing this creates a new resource.

* `vhost` - (Required, String, ForceNew) Specifies the vhost name.
  Changing this creates a new resource.

  -> If `vhost` has slashes, please change them into **\_\_F_SLASH\_\_**.

* `name` - (Required, String, ForceNew) Specifies the queue name.
  Changing this creates a new resource.

* `auto_delete` - (Required, Bool, ForceNew) Specifies whether to enable auto delete.
  Changing this creates a new resource.

* `durable` - (Optional, Bool, ForceNew) Specifies whether to enable durable. Defaults to **false**.
  Changing this creates a new resource.

* `dead_letter_exchange` - (Optional, String, ForceNew) Specifies the name of the dead letter exchange.
  It's required when `dead_letter_routing_key` is specified.
  Changing this creates a new resource.

* `dead_letter_routing_key` - (Optional, String, ForceNew) Specifies the routing key of the dead letter exchange.
  Changing this creates a new resource.

* `message_ttl` - (Optional, Int, ForceNew) Specifies how long a message in this queue can be retained.
  Changing this creates a new resource.

* `lazy_mode` - (Optional, String, ForceNew) Specifies the lazy mode. Valid value is **lazy**.
  Changing this creates a new resource.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.

* `consumers` - Indicates the connected consumers.

* `messages` - Indicates the accumulated messages.

* `policy` - Indicates the policy.

* `consumer_details` - Indicates the details of subscribed consumers.
  The [consumer_details](#attrblock--consumer_details) structure is documented below.

* `queue_bindings` - Indicates the bindings to this queue.
  The [queue_bindings](#attrblock--queue_bindings) structure is documented below.

<a name="attrblock--consumer_details"></a>
The `consumer_details` block supports:

* `ack_required` - Indicates whether manual acknowledgement is enabled on the consumer client.

* `channel_details` - Indicates the consumer connections.
  The [channel_details](#attrblock--consumer_details--channel_details) structure is documented below.

* `consumer_tag` - Indicates the consumer tag.

* `prefetch_count` - Indicates the consumer client preset value.

<a name="attrblock--consumer_details--channel_details"></a>
The `channel_details` block supports:

* `connection_name` - Indicates the connection details.

* `name` - Indicates the channel details

* `number` - Indicates the channel quantity.

* `peer_host` - Indicates the IP address of the connected consumer.

* `peer_port` - Indicates the port of the process of the connected consumer.

* `user` - Indicates the consumer username. If ACL is enabled, the real username will be returned, otherwise null will
  be returned.

<a name="attrblock--queue_bindings"></a>
The `queue_bindings` block supports:

* `destination` - Indicates the binding target name.

* `destination_type` - Indicates the binding target type.

* `properties_key` - Indicates the URL-translated routing key.

* `routing_key` - Indicates the binding key-value.

* `source` - Indicates the exchange name.

## Import

The RabbitMQ queue can be imported using the `instance_id`, `vhost` and `name` separated by slashes or commas, but if
`name` contains slashes, the import ID can only be separated by commas e.g.

```bash
$ terraform import huaweicloud_dms_rabbitmq_queue.test <instance_id>/<vhost>/<name>
```

```bash
$ terraform import huaweicloud_dms_rabbitmq_queue.test <instance_id>,<vhost>,<name>
```
