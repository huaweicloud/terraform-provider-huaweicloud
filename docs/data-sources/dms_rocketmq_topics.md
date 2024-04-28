---
subcategory: "Distributed Message Service (DMS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_dms_rocketmq_topics"
description: ""
---

# huaweicloud_dms_rocketmq_topics

Use this data source to get the list of DMS rocketMQ topics.

## Example Usage

```hcl
var "instance_id" {}

data "huaweicloud_dms_rocketmq_topics" "test" {
  instance_id           = var.instance_id
  name                  = "topic1"
  total_read_queue_num  = 3
  total_write_queue_num = 3
  permission            = "all"
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the data source.
  If omitted, the provider-level region will be used.

* `instance_id` - (Required, String) Specifies the ID of the rocketMQ instance.

* `name` - (Optional, String) Specifies the topic name.

* `total_read_queue_num` - (Optional, Int) Specifies the number of total read queue.

* `total_write_queue_num` - (Optional, Int) Specifies the number of total write queue.

* `permission` - (Optional, String) Specifies the permission. Value options: **sub**, **pub** or **all**.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `topics` - The list of topics.
  The [topics](#DMS_rockermq_topics) structure is documented below.

<a name="DMS_rockermq_topics"></a>
The `topics` block supports:

* `name` - Indicates the topic name.

* `total_read_queue_num` - Indicates the number of total read queue.

* `total_write_queue_num` - Indicates the number of total write queue.

* `permission` - Indicates the permission. Value options: **sub**, **pub** or **all**.

* `brokers` - The list of brokers.
  The [brokers](#DMS_rockermq_topic_brokers) structure is documented below.

<a name="DMS_rockermq_topic_brokers"></a>
The `brokers` block supports:

* `broker_name` - Indicates the broker name.

* `read_queue_num` - Indicates the number of read queue.

* `write_queue_num` - Indicates the number of write queue.
