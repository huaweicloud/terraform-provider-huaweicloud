---
subcategory: "Distributed Message Service (DMS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_dms_kafka_consumer_group_members"
description: |-
  Use this data source to get the member list under the specified consumer group within HuaweiCloud.
---

# huaweicloud_dms_kafka_consumer_group_members

Use this data source to get the member list under the specified consumer group within HuaweiCloud.

## Example Usage

### Query all member list under the specified consumer group

```hcl
variable "instance_id" {}
variable "consumer_group_id" {}

data "huaweicloud_dms_kafka_consumer_group_members" "test" {
  instance_id = var.instance_id
  group       = var.consumer_group_id
}
```

### Query member list by the specified consumer address

```hcl
variable "instance_id" {}
variable "consumer_group_id" {}
variable "consumer_address" {}

data "huaweicloud_dms_kafka_consumer_group_members" "test" {
  instance_id = var.instance_id
  group       = var.consumer_group_id
  host        = var.consumer_address
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region where the consumer group members are located.
  If omitted, the provider-level region will be used.

* `instance_id` - (Required, String) Specifies the ID of the Kafka instance.

* `group` - (Required, String) Specifies the ID of the consumer group.

* `host` - (Optional, String) Specifies the address of the consumer.  
  Fuzzy search is supported.

* `member_id` - (Optional, String) Specifies the ID of the consumer.  
  Fuzzy search is supported.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `members` - The list of consumer group members that match the filter parameters.  
  The [members](#kafka_consumer_group_members_struct) structure is documented below.

<a name="kafka_consumer_group_members_struct"></a>
The `members` block supports:

* `id` - The ID of the consumer.

* `host` - The address of the consumer.

* `client_id` - The ID of the client.
