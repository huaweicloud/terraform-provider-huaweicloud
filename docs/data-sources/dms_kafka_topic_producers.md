---
subcategory: "Distributed Message Service (DMS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_dms_kafka_topic_producers"
description: |-
  Use this data source to get the list of Kafka topic producers.
---

# huaweicloud_dms_kafka_topic_producers

Use this data source to get the list of Kafka topic producers.

## Example Usage

```hcl
variable "instance_id" {}
variable "topic" {}

data "huaweicloud_dms_kafka_topic_producers" "test" {
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

* `producers` - Indicates the producer list.

  The [producers](#producers_struct) structure is documented below.

<a name="producers_struct"></a>
The `producers` block supports:

* `producer_address` - Indicates the producer address.

* `broker_address` - Indicates the broker address.

* `join_time` - Indicates the time when the broker was connected.
