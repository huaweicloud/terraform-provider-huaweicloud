---
subcategory: "Distributed Message Service (DMS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_dms_kafka_topic"
description: ""
---

# huaweicloud_dms_kafka_topic

Manages a DMS kafka topic resource within HuaweiCloud.

## Example Usage

```hcl
variable "kafka_instance_id" {}

resource "huaweicloud_dms_kafka_topic" "topic" {
  instance_id = var.kafka_instance_id
  name        = "topic_1"
  partitions  = 20
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) The region in which to create the DMS kafka topic resource. If omitted, the
  provider-level region will be used. Changing this creates a new resource.

* `instance_id` - (Required, String, ForceNew) Specifies the ID of the DMS kafka instance to which the topic belongs.
  Changing this creates a new resource.

* `name` - (Required, String, ForceNew) Specifies the name of the topic. The name starts with a letter, consists of 3 to
  200 characters, and supports only letters, digits, hyphens (-), underscores (_) and periods (.).
  Changing this creates a new resource.

* `partitions` - (Required, Int) Specifies the partition number. The value ranges from `1` to `200`.
  
  -> Only support to add partitions.

* `new_partition_brokers` - (Optional, List) Specifies the integers list of brokers for new partitions.
  
  -> It's only valid when adding partitions.

* `replicas` - (Optional, Int, ForceNew) Specifies the replica number.
  The value ranges from `1` to `3` and defaults to `3`. Changing this creates a new resource.

* `aging_time` - (Optional, Int) Specifies the aging time in hours.
  The value ranges from `1` to `720` and defaults to `72`.

* `sync_replication` - (Optional, Bool) Whether or not to enable synchronous replication.

* `sync_flushing` - (Optional, Bool) Whether or not to enable synchronous flushing.

* `description` - (Optional, String) Specifies the description of topic.

* `configs` - (Optional, List) Specifies the other topic configurations.

  The [configs](#topics_configs_struct) structure is documented below.

<a name="topics_configs_struct"></a>
The `configs` block supports:

* `name` - (Required, String) Specifies the configuration name.

* `value` - (Required, String) Specifies the configuration value.

  -> When `name` is **max.message.bytes**, `value` ranges from `0` to `10,485,760`.
  When `name` is **message.timestamp.type**, `value` can be **LogAppendTime** and **CreateTime**.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID which equals to the topic name.

* `policies_only` - Indicates whether this policy is the default policy.

* `type` - Indicates the topic type.

* `created_at` - Indicates the topic create time.

## Import

DMS kafka topics can be imported using the kafka instance ID and topic name separated by a slash, e.g.:

```sh
terraform import huaweicloud_dms_kafka_topic.topic c8057fe5-23a8-46ef-ad83-c0055b4e0c5c/topic_1
```
