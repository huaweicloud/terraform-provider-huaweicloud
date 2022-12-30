---
subcategory: "Distributed Message Service (DMS)"
---

# huaweicloud_dms_rocketmq_group_topic

Use this data source to get the list of DMS RocketMQ topics that associated with the consumer group.

## Example Usage

```HCL
variable "instance_id" {}
variable "group" {}
data "huaweicloud_dms_rocketmq_group_topic" "test" {
  instance_id = var.instance_id
  group       = var.group
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the data source.
  If omitted, the provider-level region will be used.

* `instance_id` - (Required, String) Specifies the ID of the RocketMQ instance.

* `group` - (Required, String) Specifies the name of the RocketMQ consumer group.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.

* `topics` - Indicates the list of RocketMQ topics associated with the consumer group.
