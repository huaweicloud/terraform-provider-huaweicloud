---
subcategory: "Distributed Message Service (DMS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_dms_rocketmq_consumer_group_topics"
description: |-
  Use this data source to get topic list under the specified consumer group within HuaweiCloud.
---

# huaweicloud_dms_rocketmq_consumer_group_topics

Use this data source to get topic list under the specified consumer group within HuaweiCloud.

## Example Usage

### Query all topics under the specified consumer group

```hcl
var "instance_id" {}
var "group_name" {}

data "huaweicloud_dms_rocketmq_consumer_group_topics" "test" {
  instance_id = var.instance_id
  group       = var.group_name
}
```

### Query the specified topic detail

```hcl
var "instance_id" {}
var "group_name" {}
var "topic_name" {}

data "huaweicloud_dms_rocketmq_consumer_group_topics" "test" {
  instance_id = var.instance_id
  group       = var.group_name
  topic       = var.topic_name
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region where the RocketMQ instance and consumer group are located.  
  If omitted, the provider-level region will be used.

* `instance_id` - (Required, String) Specifies the ID of the RocketMQ instance.

* `group` - (Required, String) Specifies the name of the consumer group.

* `topic` - (Optional, String) Specifies the name of the topic to be queried.  
  If omitted, queries the topic list. If specified, queries the topic detail.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `topics` - The list of topics consumed by the consumer group.

  -> Available if the filter parameter `topic` is **null** or omitted.

* `lag` - The number of consumption accumulations.

* `max_offset` - The total number of messages.

* `consumer_offset` - The number of consumed messages.

* `brokers` - The brokers associated with the topic.  
  The [brokers](#rocketmq_consumer_group_topics_brokers_attr) structure is documented below.

  -> The `lag`, `max_offset`, `consumer_offset` and `brokers` are available only when the
     filter parameter `topic` is specified.

<a name="rocketmq_consumer_group_topics_brokers_attr"></a>
The `brokers` block supports:

  * `broker_name` - The name of the broker.

  * `queues` - The queue details of the associated broker.  
    The [queues](#rocketmq_consumer_group_topics_queues_attr) structure is documented below.

<a name="rocketmq_consumer_group_topics_queues_attr"></a>
The `queues` block supports:

  * `id` - The ID of the queue.

  * `lag` - The number of consumption accumulations.

  * `broker_offset` - The total number of messages.

  * `consumer_offset` - The number of consumed messages.

  * `last_message_time` - The storage time of the latest consumed message, in RFC3339 format.
