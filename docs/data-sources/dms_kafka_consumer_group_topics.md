---
subcategory: "Distributed Message Service (DMS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_dms_kafka_consumer_group_topics"
description: |-
  Use this data source to get the list of topics consumed by a consumer group within HuaweiCloud.
---

# huaweicloud_dms_kafka_consumer_group_topics

Use this data source to get the list of topics consumed by a consumer group within HuaweiCloud.

## Example Usage

### Query topic list under the specified consumer group

```hcl
variable "instance_id" {}
variable "consumer_group_id" {}

data "huaweicloud_dms_kafka_consumer_group_topics" "test" {
  instance_id = var.instance_id
  group       = var.consumer_group_id
}
```

### Query topic list by the specified topic name

```hcl
variable "instance_id" {}
variable "consumer_group_id" {}
variable "topic_name" {}

data "huaweicloud_dms_kafka_consumer_group_topics" "test" {
  instance_id = var.instance_id
  group       = var.consumer_group_id
  topic       = var.topic_name
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region where the Kafka instance and consumer group are located.  
  If omitted, the provider-level region will be used.

* `instance_id` - (Required, String) Specifies the ID of the Kafka instance.

* `group` - (Required, String) Specifies the ID of the consumer group.

* `topic` - (Optional, String) Specifies the name of the topic to be queried.  
  Fuzzy search is supported.

* `sort_key` - (Optional, String) Specifies the sorting field for the query result.  
  The valid values are as follows:
  + **topic**: Sort by topic name.
  + **partition**: Sort by number of partitions.
  + **messages**: Sort by number of messages (default).

* `sort_dir` - (Optional, String) Specifies the sorting order for the query result.  
  The valid values are as follows:
  + **asc**: Ascending order.
  + **desc**: Descending order (default).

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `topics` - The list of topics that match the filter parameters.  
  The [topics](#kafka_consumer_group_topics_topics_attr) structure is documented below.

<a name="kafka_consumer_group_topics_topics_attr"></a>
The `topics` block supports:

* `topic` - The name of the topic.

* `partitions` - The number of partitions.

* `lag` - The number of message accumulations.
