---
subcategory: "Distributed Message Service (DMS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_dms_rabbitmq_recycle_instance_restore"
description: |-
  Use this resource to restore RabbitMQ instance from recycle bin within HuaweiCloud.
---

# huaweicloud_dms_rabbitmq_recycle_instance_restore

Use this resource to restore RabbitMQ instance from recycle bin within HuaweiCloud.

-> This resource is only a one-time action resource for restoring recycle bin instance. Deleting this resource will
   not clear the corresponding request record, but will only remove the resource information from the tfstate file.

## Example Usage

```hcl
variable "instance_id" {}

resource "huaweicloud_dms_rabbitmq_recycle_instance_restore" "test" {
  instance_id = var.instance_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region where the instance is located to be restored.  
  If omitted, the provider-level region will be used. Changing this creates a new resource.

* `instance_id` - (Required, String, NonUpdatable) Specifies the instance ID to be restored from recycle bin.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.

## Timeouts

This resource provides the following timeouts configuration options:

* `create` - Default is 20 minutes.
