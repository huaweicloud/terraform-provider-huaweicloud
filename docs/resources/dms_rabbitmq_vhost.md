---
subcategory: "Distributed Message Service (DMS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_dms_rabbitmq_vhost"
description: |-
  Manages a DMS RabbitMQ vhost resource within HuaweiCloud.
---

# huaweicloud_dms_rabbitmq_vhost

Manages a DMS RabbitMQ vhost resource within HuaweiCloud.

## Example Usage

```hcl
variable "instance_id" {}
variable "name" {}

resource "huaweicloud_dms_rabbitmq_vhost" "test" {
  instance_id = var.instance_id
  name        = var.name
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) The region in which to create the resource.
  If omitted, the provider-level region will be used. Changing this creates a new resource.

* `instance_id` - (Required, String, ForceNew) Specifies the DMS RabbitMQ instance ID.
  Changing this creates a new resource.

* `name` - (Required, String, ForceNew) Specifies the vhost name. Changing this creates a new resource.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.

* `tracing` - Indicates whether the message tracing is enabled.

## Import

The RabbitMQ vhost can be imported using the `instance_id` and `name` separated by a slash or a comma, but if `name`
contains slashes, the import ID can only be separated by a comma, e.g.

```bash
$ terraform import huaweicloud_dms_rabbitmq_vhost.test <instance_id>/<name>
```

```bash
$ terraform import huaweicloud_dms_rabbitmq_vhost.test <instance_id>,<name>
```
