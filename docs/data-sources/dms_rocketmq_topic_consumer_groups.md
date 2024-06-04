---
subcategory: "Distributed Message Service (DMS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_dms_rocketmq_topic_consumer_groups"
description: |-
  Use this data source to get the list of RocketMQ topic consumer groups.
---

# huaweicloud_dms_rocketmq_topic_consumer_groups

Use this data source to get the list of RocketMQ topic consumer groups.

## Example Usage

```hcl
variable "instance_id" {}
variable "topic_name" {}

data "huaweicloud_dms_rocketmq_topic_consumer_groups" "test" {
  instance_id = var.instance_id
  topic_name  = var.topic_name
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

* `instance_id` - (Required, String) Specifies the instance ID.

* `topic_name` - (Required, String) Specifies the topic name.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `groups` - Indicates the consumer group list.
