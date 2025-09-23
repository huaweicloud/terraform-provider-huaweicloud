---
subcategory: "Distributed Message Service (DMS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_dms_rabbitmq_exchanges"
description: |-
  Use this data source to get the list of DMS RabbitMQ exchanges.
---

# huaweicloud_dms_rabbitmq_exchanges

Use this data source to get the list of DMS RabbitMQ exchanges.

## Example Usage

```hcl
variable "instance_id" {}
variable "vhost" {}

data "huaweicloud_dms_rabbitmq_exchanges" "test" {
  instance_id = var.instance_id
  vhost       = var.vhost
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

* `instance_id` - (Required, String) Specifies the DMS RabbitMQ instance ID.

* `vhost` - (Required, String) Specifies the vhost name.

  -> If `vhost` has slashes, please change them into **\_\_F_SLASH\_\_**.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `exchanges` - Indicates the list of exchanges.

  The [exchanges](#exchanges_struct) structure is documented below.

<a name="exchanges_struct"></a>
The `exchanges` block supports:

* `name` - Indicates the exchange name.

* `type` - Indicates the exchange type.

* `auto_delete` - Indicates whether the auto delete is enabled.

* `durable` - Indicates whether the durable is enabled.

* `internal` - Indicates whether the exchange is internal.

* `default` - Indicates whether the exchange is default.

* `arguments` - The argument configuration of the exchange, in JSON format.
