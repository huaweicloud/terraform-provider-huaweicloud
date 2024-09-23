---
subcategory: "Distributed Message Service (DMS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_dms_rabbitmq_exchange"
description: |-
  Manages a DMS RabbitMQ exchange resource within HuaweiCloud.
---

# huaweicloud_dms_rabbitmq_exchange

Manages a DMS RabbitMQ exchange resource within HuaweiCloud.

## Example Usage

```hcl
variable "instance_id" {}
variable "vhost" {}
variable "name" {}

resource "huaweicloud_dms_rabbitmq_exchange" "test" {
  instance_id = var.instance_id
  vhost       = var.vhost
  name        = var.name
  type        = "direct"
  auto_delete = false
  durable     = true
  internal    = false
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) The region in which to create the resource.
  If omitted, the provider-level region will be used. Changing this creates a new resource.

* `instance_id` - (Required, String, ForceNew) Specifies the DMS RabbitMQ instance ID.
  Changing this creates a new resource.

* `vhost` - (Required, String, ForceNew) Specifies the vhost name. Changing this creates a new resource.
  
  -> If `vhost` has slashes, please change them into **\_\_F_SLASH\_\_**.

* `name` - (Required, String, ForceNew) Specifies the exchange name. Changing this creates a new resource.

* `type` - (Required, String, ForceNew) Specifies the exchange type. Valid values are **direct**, **fanout**, **topic**
  and **headers**. Changing this creates a new resource.

* `auto_delete` - (Required, Bool, ForceNew) Specifies whether to enable auto delete. Changing this creates a new resource.

* `durable` - (Optional, Bool, ForceNew) Specifies whether to enable durable. Defaults to **false**.
  Changing this creates a new resource.

* `internal` - (Optional, Bool, ForceNew) Specifies whether the exchange is internal. Defaults to **false**.
  Changing this creates a new resource.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.

* `bindings` - Indicates the exchange bindings.
  The [bindings](#bindings_struct) structure is documented below.

<a name="bindings_struct"></a>
The `bindings` block supports:

* `destination_type` - Indicates the destination type.

* `destination` - Indicates the destination.

* `routing_key` - Indicates the routin key.

* `properties_key` - Indicates the properties key.

## Import

The RabbitMQ exchange can be imported using the `instance_id`, `vhost` and `name` separated by slashes or commas, but if
`name` contains slashes, the import ID can only be separated by commas e.g.

```bash
$ terraform import huaweicloud_dms_rabbitmq_exchange.test <instance_id>/<vhost>/<name>
```

```bash
$ terraform import huaweicloud_dms_rabbitmq_exchange.test <instance_id>,<vhost>,<name>
```
