---
subcategory: "Distributed Message Service (DMS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_dms_rocketmq_topic"
description: ""
---

# huaweicloud_dms_rocketmq_topic

Manages DMS RocketMQ topic resources within HuaweiCloud.

## Example Usage

### Create a topic for 5.x version instance

```hcl
variable "instance_id" {}

resource "huaweicloud_dms_rocketmq_topic" "test" {
  instance_id  = var.instance_id
  name         = "topic_test"
  message_type = "NORMAL"
}
```

### Create a topic with brokers for 4.8.0 version instance

```hcl
variable "instance_id" {}

resource "huaweicloud_dms_rocketmq_topic" "test" {
  instance_id = var.instance_id
  name        = "topic_test"
  queue_num   = 3
  permission  = "all"

  brokers {
    name = "broker-0"
  }
}
```

### Create a topic with queues for 4.8.0 version instance

```hcl
variable "instance_id" {}

resource "huaweicloud_dms_rocketmq_topic" "test" {
  instance_id = var.instance_id
  name        = "topic_test"
  permission  = "all"

  queues {
    broker    = "broker-0"
    queue_num = 3
  }
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used. Changing this parameter will create a new resource.

* `instance_id` - (Required, String, ForceNew) Specifies the ID of the rocketMQ instance.
  Changing this parameter will create a new resource.

* `name` - (Required, String, ForceNew) Specifies the name of the topic.  
  The valid length is limited from `3` to `64`, only letters, digits, vertical lines (|), percent sign (%), hyphens (-)
  and underscores (_) are allowed.
  Changing this parameter will create a new resource.

* `message_type` - (Optional, String, ForceNew) Specifies the message type of the topic.
  It's only valid when RocketMQ instance version is **5.x**. Valid values are:
  + **NORMAL**: Normal messages.
  + **FIFO**: Ordered messages.
  + **DELAY**: Scheduled messages.
  + **TRANSACTION**: Transactional messages.

  Changing this parameter will create a new resource.

* `brokers` - (Optional, List, ForceNew) Specifies the list of associated brokers of the topic.
  It's only valid when RocketMQ instance version is **4.8.0**.
  Changing this parameter will create a new resource.
  The [brokers](#DmsRocketMQTopic_BrokerRef) structure is documented below.

* `queue_num` - (Optional, Int, ForceNew) Specifies the number of queues.  
  The valid value is range from `1` to `50`. Defaults to `8`.
  It's only valid when RocketMQ instance version is **4.8.0**.
  Changing this parameter will create a new resource.

* `queues` - (Optional, List, ForceNew) Specifies the queues information of the topic.
  It's only valid when RocketMQ instance version is **4.8.0**.
  The [queues](#DmsRocketMQTopic_QueueRef) structure is documented below.
  Changing this parameter will create a new resource.

* `permission` - (Optional, String) Specifies the permissions of the topic.
  Value options: **all**, **sub**, **pub**. Defaults to **all**.
  It's only valid when RocketMQ instance version is **4.8.0**.

* `total_read_queue_num` - (Optional, Int) Specifies the total number of read queues.

* `total_write_queue_num` - (Optional, Int) Specifies the total number of write queues.

<a name="DmsRocketMQTopic_BrokerRef"></a>
The `brokers` block supports:

* `name` - (Optional, String) Specifies the name of the broker.

<a name="DmsRocketMQTopic_QueueRef"></a>
The `queues` block supports:

* `broker` - (Optional, String) Specifies the associated broker.

* `queue_num` - (Optional, Int) Specifies the number of the queues.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.

<a name="DmsRocketMQTopic_BrokerRef"></a>
  The `brokers` block supports:

* `read_queue_num` - Indicates the read queues number of the broker. It's useless when create a topic.

* `write_queue_num` - Indicates the read queues number of the broker. It's useless when create a topic.

## Import

The rocketmq topic can be imported using the rocketMQ instance ID and topic name separated by a slash, e.g.

```bash
$ terraform import huaweicloud_dms_rocketmq_topic.test c8057fe5-23a8-46ef-ad83-c0055b4e0c5c/topic_1
```
