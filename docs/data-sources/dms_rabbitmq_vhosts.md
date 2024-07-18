---
subcategory: "Distributed Message Service (DMS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_dms_rabbitmq_vhosts"
description: |-
  Use this data source to get the list of DMS RabbitMQ vhosts.
---

# huaweicloud_dms_rabbitmq_vhosts

Use this data source to get the list of DMS RabbitMQ vhosts.

## Example Usage

```hcl
variable "instance_id" {}

data "huaweicloud_dms_rabbitmq_vhosts" "test" {
  instance_id = var.instance_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

* `instance_id` - (Required, String) Specifies the DMS RabbitMQ instance ID.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `vhosts` - Indicates the vhosts list.

  The [vhosts](#vhosts_struct) structure is documented below.

<a name="vhosts_struct"></a>
The `vhosts` block supports:

* `name` - Indicates the vhost name.

* `tracing` - Indicates whether the message tracing of the vhost is enabled.
