---
subcategory: "Distributed Message Service (DMS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_dms_kafka_topics"
description: |-
  Use this data source to get the list of Kafka instance topics.
---

# huaweicloud_dms_kafka_topics

Use this data source to get the list of Kafka instance topics.

## Example Usage

```hcl
variable "instance_id" {}

data "huaweicloud_dms_kafka_topics" "test" {
  instance_id = var.instance_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

* `instance_id` - (Required, String) Specifies the DMS kafka instance ID.

* `name` - (Optional, String) Specifies the topic name.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `topics` - Indicates the topic list.

  The [topics](#topics_struct) structure is documented below.

* `max_partitions` - Indicates the total number of partitions.

* `remain_partitions` - Indicates the number of remaining partitions.

* `topic_max_partitions` - Indicates the maximum number of partitions in a single topic.

<a name="topics_struct"></a>
The `topics` block supports:

* `name` - Indicates the topic name.

* `partitions` - Indicates the number of topic partitions.

* `replicas` - Indicates the number of replicas.

* `aging_time` - Indicates the aging time in hours.

* `sync_replication` - Indicates whether the synchronous replication is enabled.

* `sync_flushing` - Indicates whether the synchronous flushing is enabled.

* `description` - Indicates the topic description.

* `configs` - Indicates the other topic configurations.

  The [configs](#topics_configs_struct) structure is documented below.

* `policies_only` - Indicates whether this policy is the default policy.

* `type` - Indicates the topic type.

* `created_at` - Indicates the topic create time.

<a name="topics_configs_struct"></a>
The `configs` block supports:

* `name` - Indicates the configuration name.

* `value` - Indicates the configuration value.
