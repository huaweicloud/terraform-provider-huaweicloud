---
subcategory: "Distributed Message Service (DMS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_dms_kafkav2_smart_connect_tasks"
description: |-
  Use this data source to get the list of DMS kafka smart connect tasks.
---

# huaweicloud_dms_kafkav2_smart_connect_tasks

Use this data source to get the list of DMS kafka smart connect tasks.

## Example Usage

```hcl
var "instance_id" {}

data "huaweicloud_dms_kafkav2_smart_connect_tasks" "test" {
  instance_id = var.instance_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

* `instance_id` - (Required, String) Specifies the kafka instance ID.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `tasks` - Indicates the smart connect task details.

  The [tasks](#tasks_struct) structure is documented below.

<a name="tasks_struct"></a>
The `tasks` block supports:

* `id` - Indicates the task ID.

* `task_name` - Indicates the task name.

* `topics` - Indicates the task topics name list.

* `topics_regex` - Indicates the regular expression of the topic of the task.

* `source_type` - Indicates the source type of the task.

* `destination_type` - Indicates the destination type of the task.

* `source_task` - Indicates the source configuration of the task.

  The [source_task](#tasks_source_task_struct) structure is documented below.

* `destination_task` - Indicates the target configuration of the task.

  The [destination_task](#tasks_destination_task_struct) structure is documented below.

* `status` - Indicates the status of the smart connect task.

* `created_at` - Indicates the creation time of the smart connect task.

<a name="tasks_source_task_struct"></a>
The `source_task` block supports:

* `current_instance_alias` - Indicates the current Kafka instance alias.

* `peer_instance_alias` - Indicates the peer Kafka instance alias.

* `peer_instance_id` - Indicates the peer Kafka instance ID.

* `peer_instance_address` - Indicates the peer Kafka instance address.

* `security_protocol` - Indicates the peer Kafka instance authentication.

* `sasl_mechanism` - Indicates the peer Kafka instance authentication mode.

* `user_name` - Indicates the peer Kafka instance username.

* `direction` - Indicates the sync direction.

* `sync_consumer_offsets_enabled` - Indicates whether to sync the consumption progress.

* `replication_factor` - Indicates the number of replicas.

* `task_num` - Indicates the number of data replication tasks.

* `rename_topic_enabled` - Indicates whether to rename a topic.

* `provenance_header_enabled` - Indicates whether to add the source header.

* `consumer_strategy` - Indicates the start offset.

* `compression_type` - Indicates  the compression algorithm to use for copying messages.

* `topics_mapping` - Indicates the topic mapping.

<a name="tasks_destination_task_struct"></a>
The `destination_task` block supports:

* `consumer_strategy` - Indicates the consumer strategy of the smart connect task.

* `deliver_time_interval` - Indicates the dumping period in seconds.

* `obs_bucket_name` - Indicates the obs bucket name of the smart connect task.

* `partition_format` - Indicates the partiton format of the smart connect task.

* `obs_path` - Indicates the obs path of the smart connect task.

* `destination_file_type` - Indicates the destination file type of the smart connect task.

* `record_delimiter` - Indicates the record delimiter of the smart connect task.

* `store_keys` - Indicates whether to store keys.

* `obs_part_size` - Indicates the size of each file to be uploaded.

* `flush_size` - Indicates the flush size.

* `timezone` - Indicates the time zone.

* `schema_generator_class` - Indicates the schema generator class.

* `partitioner_class` - Indicates the partitioner class.

* `key_converter` - Indicates the key converter.

* `value_converter` - Indicates the value converter.

* `kv_delimiter` - Indicates the kv delimiter.
