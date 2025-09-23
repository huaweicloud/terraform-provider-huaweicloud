---
subcategory: "Distributed Message Service (DMS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_dms_rabbitmq_exchange_associate"
description: |-
  Manages a DMS RabbitMQ exchange association resource within HuaweiCloud.
---

# huaweicloud_dms_rabbitmq_exchange_associate

Manages a DMS RabbitMQ exchange association resource within HuaweiCloud.

## Example Usage

```hcl
variable "instance_id" {}
variable "vhost" {}
variable "exchange" {}
variable "destination" {}
variable "routing_key" {}

resource "huaweicloud_dms_rabbitmq_exchange_associate" "test" {
  instance_id      = var.instance_id
  vhost            = var.vhost
  exchange         = var.exchange
  destination_type = "Queue"
  destination      = var.destination
  routing_key      = var.routing_key
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) The region in which to create the resource.
  If omitted, the provider-level region will be used. Changing this creates a new resource.

* `instance_id` - (Required, String, ForceNew) Specifies the DMS RabbitMQ instance ID.
  Changing this creates a new resource.

* `vhost` - (Required, String, ForceNew) Specifies the vhost name. Changing this creates a new resource.

* `exchange` - (Required, String, ForceNew) Specifies the exchange name. Changing this creates a new resource.

-> If `vhost` and `exchange` has slashes, please change them into **\_\_F_SLASH\_\_**.

* `destination_type` - (Required, String, ForceNew) Specifies the type of the binding target.
  The options are **Exchange** and **Queue**. Changing this creates a new resource.

* `destination` - (Required, String, ForceNew) Specifies the name of a target exchange or queue.
  Changing this creates a new resource.

* `routing_key` - (Optional, String, ForceNew) Specifies the binding key-value. Changing this creates a new resource.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.

* `properties_key` - The URL-translated routing key.

## Import

The RabbitMQ exchange association can be imported using the `instance_id`, `vhost`, `exchange`, `destination_type`,
`destination` and `routing_key` separated by commas.

If `routing_key` is empty e.g.

```bash
$ terraform import huaweicloud_dms_rabbitmq_exchange.test <instance_id>,<vhost>,<exchange>,<destination_type>,<destination>
```

If `routing_key` is specified e.g.

```bash
$ terraform import huaweicloud_dms_rabbitmq_exchange.test <instance_id>,<vhost>,<exchange>,<destination_type>,<destination>,<routing_key>
```
