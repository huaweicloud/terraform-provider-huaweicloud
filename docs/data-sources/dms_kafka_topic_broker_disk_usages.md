---
subcategory: "Distributed Message Service (DMS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_dms_kafka_topic_broker_disk_usages"
description: |-
  Use this data source to query the disk usage of topics on brokers within HuaweiCloud.
---

# huaweicloud_dms_kafka_topic_broker_disk_usages

Use this data source to query the disk usage of topics on brokers within HuaweiCloud.

## Example Usage

### Query all topics on brokers

```hcl
variable "instance_id" {}

data "huaweicloud_dms_kafka_topic_broker_disk_usages" "test" {
  instance_id = var.instance_id
}
```

### Query the top 10 topics on brokers

```hcl
variable "instance_id" {}

data "huaweicloud_dms_kafka_topic_broker_disk_usages" "test" {
  instance_id = var.instance_id
  top         = 10
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region where the topic broker disk usages are located.  
  If omitted, the provider-level region will be used.

* `instance_id` - (Required, String) Specifies the ID of the kafka instance.

* `min_size` - (Optional, String) Specifies the minimum disk size threshold to be queried.  
  The valid values are as follows:
  + **1K**
  + **1M**
  + **1G**

* `top` - (Optional, Int) Specifies the number of top topics to be queried.  
  The valid values ranges from `1` to `1,000`.

* `percentage` - (Optional, String) Specifies the percentage threshold to be queried.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `disk_usages` - The disk usage of topics on brokers that match the filter parameters.  
  The [disk_usages](#kafka_topic_broker_disk_usages_struct) structure is documented below.

<a name="kafka_topic_broker_disk_usages_struct"></a>
The `disk_usages` block supports:

* `broker_name` - The name of the broker.

* `data_disk_size` - The total disk capacity.

* `data_disk_use` - The used disk capacity.

* `data_disk_free` - The free disk capacity.

* `data_disk_use_percentage` - The disk usage percentage.

* `status` - The status of the broker.

* `topics` - The list of topic disk usage.  
  The [topics](#kafka_topic_broker_disk_usages_topics_struct) structure is documented below.

<a name="kafka_topic_broker_disk_usages_topics_struct"></a>
The `topics` block supports:

* `size` - The size of the disk usage.

* `topic_name` - The name of the topic.

* `topic_partition` - The partition of the topic.

* `percentage` - The percentage of the disk usage.
