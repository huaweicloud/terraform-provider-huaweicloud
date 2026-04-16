---
subcategory: "Distributed Message Service (DMS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_dms_rabbitmq_recycle_instances"
description: |-
  Use this data source to query RabbitMQ recycle bin instance list within HuaweiCloud.
---

# huaweicloud_dms_rabbitmq_recycle_instances

Use this data source to query RabbitMQ recycle bin instance list within HuaweiCloud.

## Example Usage

```hcl
data "huaweicloud_dms_rabbitmq_recycle_instances" "test" {}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region where the recycle bin instances are located.  
  If omitted, the provider-level region will be used.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `retention_days` - The retention days of the recycle bin.

* `default_use_recycle` - Whether the recycle bin is enabled.

* `instances` - The list of recycle bin instances.  
  The [instances](#dms_rabbitmq_recycle_instances_attr) structure is documented below.

<a name="dms_rabbitmq_recycle_instances_attr"></a>
The `instances` block supports:

* `instance_id` - The ID of the instance.

* `name` - The name of the instance.

* `status` - The status of the instance.
  + **ERROR**
  + **RECYCLE**

* `engine` - The engine of the instance.

* `in_recycle_time` - The time when the instance was placed in the recycle bin, in FRC3339 format.

* `save_time` - The time when the instance was saved, in day.

* `auto_delete_time` - The time when the instance was automatically deleted, in FRC3339 format.

* `cost_per_hour` - The cost per hour of the instance.

* `error_message` - The error message.

* `product_id` - The ID of the flavor of the instance.
