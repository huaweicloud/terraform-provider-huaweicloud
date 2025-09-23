---
subcategory: "Distributed Message Service (DMS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_dms_rabbitmq_queue_message_clear"
description: |-
  Manages a DMS RabbitMQ queue message clear resource within HuaweiCloud.
---

# huaweicloud_dms_rabbitmq_queue_message_clear

Manages a DMS RabbitMQ queue message clear resource within HuaweiCloud.

## Example Usage

```hcl
variable "instance_id" {}
variable "vhost" {}
variable "queue" {}

resource "huaweicloud_dms_rabbitmq_queue_message_clear" "test" {
  instance_id = var.instance_id
  vhost       = urlencode(replace(var.vhost, "/", "__F_SLASH__"))
  queue       = urlencode(replace(var.queue, "/", "__F_SLASH__"))
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

* `queue` - (Required, String, ForceNew) Specifies the queue name.
  Changing this creates a new resource.

-> If `vhost` and `queue` has slashes, please change them into **\_\_F_SLASH\_\_**.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.
