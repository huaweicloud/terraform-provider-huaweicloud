---
subcategory: "Distributed Message Service (DMS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_dms_kafka_consumer_groups"
description: |-
  Use this data source to get the list of Kafka consumer groups.
---

# huaweicloud_dms_kafka_consumer_groups

Use this data source to get the list of Kafka consumer groups.

## Example Usage

```hcl
variable "instance_id" {}

data "huaweicloud_dms_kafka_consumer_groups" "test" {
  instance_id = var.instance_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

* `instance_id` - (Required, String) Specifies the instance ID.

* `name` - (Optional, String) Specifies the consumer group name.

* `description` - (Optional, String) Specifies the consumer group description.

* `lag` - (Optional, Int) Specifies the the number of accumulated messages.

* `coordinator_id` - (Optional, Int) Specifies the coordinator ID.

* `state` - (Optional, String) Specifies the consumer group status.
  The value can be:
  + **Dead**: The consumer group has no members or metadata.
  + **Empty**: The consumer group has metadata but has no members.
  + **PreparingRebalance**: The consumer group is to be rebalanced.
  + **CompletingRebalance**: All members have joined the group.
  + **Stable**: Members in the consumer group can consume messages.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `groups` - Indicates the consumer groups.

  The [groups](#groups_struct) structure is documented below.

<a name="groups_struct"></a>
The `groups` block supports:

* `name` - Indicates the consumer group name.

* `description` - Indicates the consumer group description.

* `lag` - Indicates the number of accumulated messages.

* `coordinator_id` - Indicates the coordinator ID.

* `state` - Indicates the consumer group status.

* `created_at` - Indicates the create time.

* `group_message_offsets` - Indicates the group message offsets.
  The [group_message_offsets](#attrblock--groups--group_message_offsets) structure is documented below.

* `members` - Indicates the consumer group members.
  The [members](#attrblock--groups--members) structure is documented below.

* `assignment_strategy` - Indicates the partition assignment strategy.

<a name="attrblock--groups--group_message_offsets"></a>
The `group_message_offsets` block supports:

* `lag` - Indicates the number of accumulated messages.

* `message_current_offset` - Indicates the message current offset.

* `message_log_end_offset` - Indicates the message log end offset.

* `partition` - Indicates the partition.

* `topic` - Indicates the topic name.

<a name="attrblock--groups--members"></a>
The `members` block supports:

* `assignment` - Indicates the details about the partition assigned to the consumer.
  The [assignment](#attrblock--groups--members--assignment) structure is documented below.

* `client_id` - Indicates the client ID.

* `host` - Indicates the consumer address.

* `member_id` - Indicates the member ID.

<a name="attrblock--groups--members--assignment"></a>
The `assignment` block supports:

* `partitions` - Indicates the partitions.

* `topic` - Indicates the topic name.
