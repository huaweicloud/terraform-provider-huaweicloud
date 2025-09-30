---
subcategory: "Distributed Message Service (DMS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_dms_kafka_topic_quota"
description: |-
  Manages a Kafka topic quota resource within HuaweiCloud.
---

# huaweicloud_dms_kafka_topic_quota

Manages a Kafka topic quota resource within HuaweiCloud.

-> This resource is unavailable for single-node instance Kafka.

## Example Usage

```hcl
variable "instance_id" {}
variable "topic_name" {}

resource "huaweicloud_dms_kafka_topic_quota" "test" {
  instance_id        = var.instance_id
  topic              = var.topic_name
  producer_byte_rate = 2097152
  consumer_byte_rate = 3145728
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region where the topic quota is located.  
  If omitted, the provider-level region will be used. Changing this parameter will create a new resource.

* `instance_id` - (Required, String, NonUpdatable) Specifies the ID of the Kafka instance to which the topic quota belongs.

* `topic` - (Required, String, NonUpdatable) Specifies the name of the topic.

* `producer_byte_rate` - (Optional, Int) Specifies the producer rate limit. The unit is byte/s.
  The valid value range `2,097,152` to `1,073,741,824`. `0` means no limit.

* `consumer_byte_rate` - (Optional, Int) Specifies the consumer rate limit. The unit is byte/s.
  The valid value range `2,097,152` to `1,073,741,824`. `0` means no limit.

-> At least one of `producer_byte_rate` and `consumer_byte_rate` must be specified.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID in format `<instance_id>/<topic>`.

## Import

Topic quotas can be imported using its the `id` (consist of `instance_id` and `topic`, separated by a slash (/)), e.g.

```bash
$ terraform import huaweicloud_dms_kafka_topic_quota.test <instance_id>/<topic>
```
