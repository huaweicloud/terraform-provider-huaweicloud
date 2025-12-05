---
subcategory: "Distributed Message Service (DMS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_dms_kafka_consumer_group"
description: ""
---

# huaweicloud_dms_kafka_consumer_group

Manages DMS Kafka consumer group resources within HuaweiCloud.

## Example Usage

```hcl
variable "instance_id" {}
resource "huaweicloud_dms_kafka_consumer_group" "test" {
  instance_id = var.instance_id
  name        = "consumer_group_test"
  description = "the description of the consumer group"
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used. Changing this parameter will create a new resource.

* `instance_id` - (Required, String, ForceNew) Specifies the ID of the kafka instance.

  Changing this parameter will create a new resource.

* `name` - (Required, String, ForceNew) Specifies the name of the consumer group.

  Changing this parameter will create a new resource.

* `description` - (Optional, String) Specifies the description of the consumer group.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.

* `state` - Indicates the state of the consumer group. This value can be :
  **DEAD**, **EMPTY**, **PreparingRebalance**, **CompletingRebalance**, **Stable**.

* `coordinator_id` - Indicates the coordinator id of the consumer group.

* `created_at` - Indicates the creation time of the consumer group.

## Import

The kafka consumer group can be imported using the kafka `instance_id` and `name` separated by a slash, e.g.

```bash
$ terraform import huaweicloud_dms_kafka_consumer_group.test <instance_id>/<name>
```
