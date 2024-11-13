---
subcategory: "Distributed Message Service (DMS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_dms_kafka_topic_partitions"
description: |-
  Use this data source to get the list of Kafka topic partitions.
---

# huaweicloud_dms_kafka_topic_partitions

Use this data source to get the list of Kafka topic partitions.

## Example Usage

```hcl
variable "instance_id" {}
variable "topic" {}

data "huaweicloud_dms_kafka_topic_partitions" "test" {
  instance_id = var.instance_id
  topic       = var.topic
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

* `instance_id` - (Required, String) Specifies the instance ID.

* `topic` - (Required, String) Specifies the topic name.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `partitions` - Indicates the partitions.
  The [partitions](#attrblock--partitions) structure is documented below.

<a name="attrblock--partitions"></a>
The `partitions` block supports:

* `partition` - Indicates the partition ID.

* `start_offset` - Indicates the start offset.

* `last_offset` - Indicates the last offset.

* `last_update_time` - Indicates the last update time.

* `message_count` - Indicates the message count.
