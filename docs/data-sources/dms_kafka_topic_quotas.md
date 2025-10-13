---
subcategory: "Distributed Message Service (DMS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_dms_kafka_topic_quotas"
description: |-
  Use this data source to get the list of Kafka topic quotas within HuaweiCloud.
---

# huaweicloud_dms_kafka_topic_quotas

Use this data source to get the list of Kafka topic quotas within HuaweiCloud.

## Example Usage

### Query all topic quotas under the specified instance

```hcl
variable "instance_id" {}

data "huaweicloud_dms_kafka_topic_quotas" "test" {
  instance_id = var.instance_id
}
```

### Query topic quotas by keyword

```hcl
variable "instance_id" {}
variable "topic_name" {}

data "huaweicloud_dms_kafka_topic_quotas" "test" {
  instance_id = var.instance_id
  keyword     = var.topic_name
}

```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region where the topic quotas are located.  
  If omitted, the provider-level region will be used.

* `instance_id` - (Required, String) Specifies the ID of the Kafka instance to which the topic quotas belong.

* `keyword` - (Optional, String) Specifies the keyword of the topic quota to be queried.  
  Fuzzy matching is supported.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `quotas` - All topic quotas that match the filter parameters.  
  The [quotas](#kafka_topic_quotas_attr) structure is documented below.

<a name="kafka_topic_quotas_attr"></a>
The `quotas` block supports:

* `topic` - The name of the topic.

* `producer_byte_rate` - The producer byte rate limit. The unit is B/s.

* `consumer_byte_rate` - The consumer byte rate limit. The unit is B/s.
