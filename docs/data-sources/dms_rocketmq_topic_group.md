---
subcategory: "Distributed Message Service (DMS)"
---

# huaweicloud_dms_rocketmq_topic_group

Use this data source to get the list of DMS RocketMQ consumer groups that associated with the topic.

## Example Usage

```hcl
variable "instance_id" {}
variable "topic" {}

data "huaweicloud_dms_rocketmq_topic_group" "test" {
  instance_id = var.instance_id
  topic       = var.topic
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the data source.
  If omitted, the provider-level region will be used.

* `instance_id` - (Required, String) Specifies the ID of the RocketMQ instance.

* `topic` - (Required, String) Specifies the name of the RocketMQ topic.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.

* `groups` - Indicates the list of RocketMQ consumer groups associated with the topic.
