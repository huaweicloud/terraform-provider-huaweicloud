---
subcategory: "Distributed Message Service (DMS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_dms_kafka_topic"
description:  |-
  Manages a DMS kafka topic resource within HuaweiCloud.
---

# huaweicloud_dms_kafka_topic

Manages a DMS kafka topic resource within HuaweiCloud.

## Example Usage

```hcl
variable "kafka_instance_id" {}
variable "topic_name" {}

resource "huaweicloud_dms_kafka_topic" "test" {
  instance_id = var.kafka_instance_id
  name        = var.topic_name
  partitions  = 20
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region where the topic is located.  
  If omitted, the provider-level region will be used. Changing this creates a new resource.

* `instance_id` - (Required, String, NonUpdatable) Specifies the ID of the DMS kafka instance to which the topic belongs.

* `name` - (Required, String, NonUpdatable) Specifies the name of the topic.  
  The name starts with a letter, consists of `3` to `200` characters, and supports only letters, digits, hyphens (-),
  underscores (_) and periods (.).

* `partitions` - (Required, Int) Specifies the partition number.  
  The value ranges from `1` to `200`.
  
  -> Only support to add partitions.

* `new_partition_brokers` - (Optional, List) Specifies the integers list of brokers for new partitions.
  
  -> It's only valid when adding partitions.

* `replicas` - (Optional, Int, NonUpdatable) Specifies the replica number.
  The value ranges from `1` to `3` and defaults to `3`.

* `aging_time` - (Optional, Int) Specifies the aging time in hours.
  The value ranges from `1` to `720` and defaults to `72`.

* `sync_replication` - (Optional, Bool) Specifies whether to enable synchronous replication.  
  Defaults to **false**.

* `sync_flushing` - (Optional, Bool) Specifies whether to enable synchronous flushing.
  Defaults to **false**.

* `description` - (Optional, String) Specifies the description of the topic.

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

* `policies_only` - Whether this policy is the default policy.

* `type` - The topic type.

* `created_at` - The topic create time.

## Timeouts

This resource provides the following timeouts configuration options:

* `create` - Default is 10 minutes.

## Import

DMS kafka topics can be imported using the `instance_id` and `name`, separated by a slash (/), e.g.:

```bash
$ terraform import huaweicloud_dms_kafka_topic.test <instance_id>/<name>
```

Note that the imported state may not be identical to your resource definition, due to some attributes missing from the
API response, security or some other reason. The missing attributes include: `new_partition_brokers`.
It is generally recommended running `terraform plan` after importing the resource.
You can then decide if changes should be applied to the resource, or the resource definition should be updated to
align with the resource. Also you can ignore changes as below.

```hcl
resource "huaweicloud_dms_kafka_topic" "test" {
  ...

  lifecycle {
    ignore_changes = [
      new_partition_brokers,
    ]
  }
}
```
