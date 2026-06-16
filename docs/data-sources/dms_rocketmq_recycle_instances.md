---
subcategory: "Distributed Message Service (DMS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_dms_rocketmq_recycle_instances"
description: |-
  Use this data source to query the list of DMS RocketMQ recycle instances within HuaweiCloud.
---

# huaweicloud_dms_rocketmq_recycle_instances

Use this data source to query the list of DMS RocketMQ recycle instances within HuaweiCloud.

## Example Usage

```hcl
data "huaweicloud_dms_rocketmq_recycle_instances" "test" {}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region where the RocketMQ recycle instances are located.  
  If omitted, the provider-level region will be used.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `instances` - The list of the recycle instances.  
  The [instances](#dms_rocketmq_recycle_instances_attr) structure is documented below.

<a name="dms_rocketmq_recycle_instances_attr"></a>
The `instances` block supports:

* `id` - The ID of the instance.

* `name` - The name of the instance.

* `status` - The status of the instance.

* `engine` - The message engine type.

* `in_recycle_time` - The time when the instance was recycled, in RFC3339 format.

* `save_time` - The number of days the instance is retained in the recycle bin.

* `auto_delete_time` - The time when the instance will be automatically deleted, in RFC3339 format.
