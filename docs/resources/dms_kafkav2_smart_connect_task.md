---
subcategory: "Distributed Message Service (DMS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_dms_kafkav2_smart_connect_task"
description: ""
---

# huaweicloud_dms_kafkav2_smart_connect_task

Manage DMS kafka smart connect task resource within HuaweiCloud.

## Example Usage

### Create a task to dump Kafka data to OBS

```hcl
variable "instance_id" {}
variable "task_name" {}
variable "topics" {}
variable "access_key" {}
variable "secret_key" {}
variable "obs_bucket_name" {}
variable "obs_path" {}

resource "huaweicloud_dms_kafkav2_smart_connect_task" "test" {
  instance_id      = var.instance_id
  task_name        = var.task_name
  topics           = var.topics
  destination_type = "OBS_SINK"

  destination_task {
    consumer_strategy     = "latest"
    destination_file_type = "TEXT"
    access_key            = var.access_key
    secret_key            = var.secret_key
    obs_bucket_name       = var.obs_bucket_name
    obs_path              = var.obs_path
    partition_format      = "yyyy/MM/dd/HH/mm"
    record_delimiter      = "\n"
    deliver_time_interval = 300
  }
}
```

### Create a Kafka data replication task

```hcl
variable "instance_id" {}
variable "task_name" {}
variable "peer_instance_id" {}

resource "huaweicloud_dms_kafkav2_smart_connect_task" "test" {
  instance_id = var.instance_id
  task_name   = var.task_name
  topics      = ["topic-test"]
  source_type = "KAFKA_REPLICATOR_SOURCE"

  source_task {
    peer_instance_id              = var.peer_instance_id
    direction                     = "push"
    replication_factor            = 3
    task_num                      = 2
    provenance_header_enabled     = false
    sync_consumer_offsets_enabled = false
    rename_topic_enabled          = false
    consumer_strategy             = "latest"
    compression_type              = "none"
    topics_mapping                = ["topic-test:topic-test-mapping"]
  }
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used. Changing this parameter will create a new resource.

* `instance_id` - (Required, String, ForceNew) Specifies the kafka instance ID.
  Changing this parameter will create a new resource.

* `task_name` - (Required, String, ForceNew) Specifies the smart connect task name.
  Changing this parameter will create a new resource.

* `topics` - (Optional, List, ForceNew) Specifies the topic names list of the smart connect task.
  Changing this parameter will create a new resource.

* `topics_regex` - (Optional, String, ForceNew) Specifies the regular expression of topic name for the smart connect
  task. Changing this parameter will create a new resource.

  -> Exactly one of `topics`, `topics_regex` should be specified.

* `start_later` - (Optional, Bool, ForceNew) Specifies whether to start a task later.
  Changing this parameter will create a new resource.

* `source_type` - (Optional, String, ForceNew) Specifies the source type of the smart connect task. Valid values are
  **KAFKA_REPLICATOR_SOURCE** and **NONE**. Changing this parameter will create a new resource.

* `destination_type` - (Optional, String, ForceNew) Specifies the destination type of the smart connect task.
  Valid values are **OBS_SINK** and **NONE**. Changing this parameter will create a new resource.

* `source_task` - (Optional, List, ForceNew) Specifies the source configuration of a smart connect task.
  The [source_task](#dms_source_task) structure is documented below.
  Changing this parameter will create a new resource.

* `destination_task` - (Optional, List, ForceNew) Specifies the destination configuration of a smart connect task.
  The [destination_task](#dms_destination_task) structure is documented below.
  Changing this parameter will create a new resource.

<a name="dms_source_task"></a>
The `source_task` block supports:

* `current_instance_alias` - (Optional, String, ForceNew) Specifies the current Kafka instance alias.
  Changing this parameter will create a new resource.

* `peer_instance_alias` - (Optional, String, ForceNew) Specifies the peer Kafka instance alias.
  Changing this parameter will create a new resource.

* `peer_instance_id` - (Optional, String, ForceNew) Specifies the peer Kafka instance ID.
  Changing this parameter will create a new resource.

* `peer_instance_address` - (Optional, List, ForceNew) Specifies the peer Kafka instance address list.
  Changing this parameter will create a new resource.

  -> Exactly one of `peer_instance_id` and `peer_instance_address` should be specified.

* `security_protocol` - (Optional, String, ForceNew) Specifies the peer Kafka authentication. Valid values are:
  + **SASL_SSL**: SASL_SSL is enabled.
  + **PLAINTEXT**: SASL_SSL is disabled.

  Changing this parameter will create a new resource.

* `sasl_mechanism` - (Optional, String, ForceNew) Specifies the peer Kafka authentication mode.
  Changing this parameter will create a new resource.

* `user_name` - (Optional, String, ForceNew) Specifies the peer Kafka user name.
  It's **required** when `security_protocol` is **SASL_SSL**. Changing this parameter will create a new resource.

* `password` - (Optional, String, ForceNew) Specifies the peer Kafka user password.
  It's **required** when `security_protocol` is **SASL_SSL**. Changing this parameter will create a new resource.

* `direction` - (Optional, String, ForceNew) Specifies the sync direction. Valid values are:
  + **pull**: Copy the data of the peer Kafka instance to the current Kafka instance.
  + **push**: Copy the data of the current Kafka instance to the peer Kafka instance.
  + **two-way**: Copy the data of the Kafka instances at both ends in both directions.

  Changing this parameter will create a new resource.

* `sync_consumer_offsets_enabled` - (Optional, Bool, ForceNew) Specifies whether to sync the consumption progress.
  Changing this parameter will create a new resource.

* `replication_factor` - (Optional, Int, ForceNew) Specifies the number of topic replicas. The value of this parameter
  cannot exceed the number of brokers in the peer instance. Changing this parameter will create a new resource.

* `task_num` - (Optional, Int, ForceNew) Specifies the number of data replication tasks.
  If the `direction` is set to **two-way**, the actual number of tasks will be twice the number of tasks you configure
  here. Changing this parameter will create a new resource.

* `rename_topic_enabled` - (Optional, Bool, ForceNew) Specifies whether to rename the topic. If true, will add the
  alias of the source Kafka instance before the target topic name to form a new name of the target topic.
  Changing this parameter will create a new resource.

* `topics_mapping` - (Optional, List, ForceNew) Specifies the topic mapping string list, which is used to customize
  the target topic name, e.g., topic-sc-1:topic-sc-2. Changing this parameter will create a new resource.

  -> When `rename_topic_enabled` is true, `topics_mapping` can not be specified.

* `provenance_header_enabled` - (Optional, Bool, ForceNew) Specifies whether the message header contains the message
  source. Changing this parameter will create a new resource.

* `consumer_strategy` - (Optional, String, ForceNew) Specifies the start offset. Value options are:
  + **latest**: Read the latest data.
  + **earliest**: Read the earliest data.

  Changing this parameter will create a new resource.

* `compression_type` - (Optional, String, ForceNew) Specifies the compression algorithm to use for copying messages.
  Valid values are **none**, **gzip**, **snappy**, **lz4** and **zstd**.
  Changing this parameter will create a new resource.

<a name="dms_destination_task"></a>
The `destination_task` block supports:

* `access_key` - (Optional, String, ForceNew) Specifies the access key used to access the OBS bucket.
  It's **required** when `destination_type` is **OBS_SINK**.
  Changing this parameter will create a new resource.

* `secret_key` - (Optional, String, ForceNew) Specifies the secret access key used to access the OBS bucket.
  It's **required** when `destination_type` is **OBS_SINK**.
  Changing this parameter will create a new resource.

* `consumer_strategy` - (Optional, String, ForceNew) Specifies the consumer strategy of the smart connect task.
  Value options:
  + **latest**: Read the latest data.
  + **earliest**: Read the earliest data.

  It's **required** when `destination_type` is **OBS_SINK**.
  Changing this parameter will create a new resource.

* `deliver_time_interval` - (Optional, Int, ForceNew) Specifies the deliver time interval of the smart connect task.
  It's **required** when `destination_type` is **OBS_SINK**.
  The value should be between `30` and `900`. Changing this parameter will create a new resource.

* `obs_bucket_name` - (Optional, String, ForceNew) Specifies the obs bucket name of the smart connect task.
  It's **required** when `destination_type` is **OBS_SINK**.
  Changing this parameter will create a new resource.

* `partition_format` - (Optional, String, ForceNew) Specifies the time directory format of the smart connect task.
  Value options: **yyyy**, **yyyy/MM**, **yyyy/MM/dd**, **yyyy/MM/dd/HH**, **yyyy/MM/dd/HH/mm**.
  It's **required** when `destination_type` is **OBS_SINK**.
  Changing this parameter will create a new resource.

* `obs_path` - (Optional, String, ForceNew) Specifies the obs path of the smart connect task.
  Obs path is separated by a slash. Changing this parameter will create a new resource.

* `destination_file_type` - (Optional, String, ForceNew) Specifies the destination file type of the smart connect task.
  Only **TEXT** is supported. Changing this parameter will create a new resource.

* `record_delimiter` - (Optional, String, ForceNew) Specifies the record delimiter of the smart connect task.
  Value options: **,**, **;**, **|**, **\n**.
  Changing this parameter will create a new resource.

* `store_keys` - (Optional, Bool, ForceNew) Specifies whether to dump keys.
  Changing this parameter will create a new resource.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.

* `destination_task` - Indicates the destination configuration of a smart connect task.
  The [destination_task](#attr_destination_task) structure is documented below.

* `status` - Indicates the status of the smart connect task.

* `created_at` - Indicates the creation time of the smart connect task.

<a name="attr_destination_task"></a>
The `destination_task` block supports:

* `obs_part_size` - Indicates the size of each file to be uploaded.

* `flush_size` - Indicates the flush size.

* `timezone` - Indicates the time zone.

* `schema_generator_class` - Indicates the schema generator class.

* `partitioner_class` - Indicates the partitioner class.

* `key_converter` - Indicates the key converter.

* `value_converter` - Indicates the value converter.

* `kv_delimiter` - Indicates the kv delimiter.

## Timeouts

This resource provides the following timeout configuration options:

* `create` - Default is 30 minutes.

## Import

The kafka smart connect task can be imported using the kafka instance `instance_id` and `task_id` separated by a slash,
e.g.

```bash
$ terraform import huaweicloud_dms_kafkav2_smart_connect_task.test <instance_id>/<task_id>
```

Note that the imported state may not be identical to your resource definition, due to some attributes missing from
the API response. The missing attributes include: `start_later`, `destination_task.0.access_key`,
`destination_task.0.secret_key` and `source_task.0.password`. It is generally recommended running `terraform plan`
after importing a kafka smart connect task. You can then decide if changes should be applied to the kafka smart connect
task, or the resource definition should be updated to align with the kafka smart connect task. Also you can ignore
changes as below.

```hcl
resource "huaweicloud_dms_kafkav2_smart_connect_task" "test" {
  ...

  lifecycle {
    ignore_changes = [
      destination_task.0.access_key, destination_task.0.secret_key, source_task.0.password,
    ]
  }
}
```
