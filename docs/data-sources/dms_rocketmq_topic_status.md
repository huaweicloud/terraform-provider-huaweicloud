---
subcategory: "Distributed Message Service (DMS)"
---

# huaweicloud_dms_rocketmq_topic_status

Use this data source to get the RocketMQ topic status.

## Example Usage

```hcl
variable "instance_id" {}
variable "topic" {}

data "huaweicloud_dms_rocketmq_topic_status" "test" {
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

* `brokers` - Indicates the broker list of RocketMQ topic associated with.
  The [Broker](#DmsRocketMQTopicStatus_Broker) structure is documented below.

<a name="DmsRocketMQTopicStatus_Broker"></a>
The `Broker` block supports:

* `queues` - Indicates the queue list owned by the RocketMQ topic.
  The [Queue](#DmsRocketMQTopicStatus_BrokerQueue) structure is documented below.

* `name` - Indicates the name of broker.

<a name="DmsRocketMQTopicStatus_BrokerQueue"></a>
The `BrokerQueue` block supports:

* `id` - Indicates the ID of queue.

* `min_offset` - Indicates the minimum offset of queue.

* `max_offset` - Indicates the maximum offset of queue.

* `last_message_time` - Indicates the time of the last message of queue.
