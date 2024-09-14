---
subcategory: "Distributed Message Service (DMS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_dms_rocketmq_consumer_group"
description: ""
---

# huaweicloud_dms_rocketmq_consumer_group

Manages DMS RocketMQ consumer group resources within HuaweiCloud.

## Example Usage

### Create consumer group for 4.8.0 version instance

```hcl
variable "instance_id" {}

resource "huaweicloud_dms_rocketmq_consumer_group" "test" {
  instance_id     = var.instance_id
  name            = "consumer_group_test"
  enabled         = true
  broadcast       = true
  brokers         = ["broker-0","broker-1"]
  retry_max_times = 3
  description     = "the description of the consumer group"
}
```

### Create consumer group for 5.x version instance

```hcl
variable "instance_id" {}

resource "huaweicloud_dms_rocketmq_consumer_group" "test" {
  instance_id     = var.instance_id
  name            = "consumer_group_test"
  enabled         = true
  broadcast       = true
  retry_max_times = 3
  description     = "the description of the consumer group"
  consume_orderly = true
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used. Changing this parameter will create a new resource.

* `instance_id` - (Required, String, ForceNew) Specifies the ID of the rocketMQ instance.

  Changing this parameter will create a new resource.

* `name` - (Required, String, ForceNew) Specifies the name of the consumer group.  
  The valid length is limited from `3` to `64`, only letters, digits, vertical lines (|), percent sign (%), hyphens (-)
  and underscores (_) are allowed.

  Changing this parameter will create a new resource.

* `retry_max_times` - (Required, Int) Specifies the maximum number of retry times.  
  The valid value is range from `1` to `16`.

* `enabled` - (Optional, Bool) Specifies the consumer group is enabled or not. Defaults to true.

* `broadcast` - (Optional, Bool) Specifies whether to broadcast of the consumer group.

* `description` - (Optional, String) Specifies the description of the consumer group.

* `brokers` - (Optional, List, ForceNew) Specifies the list of associated brokers of the consumer group.
  It's only valid when RocketMQ instance version is **4.8.0**.
  Changing this parameter will create a new resource.

* `consume_orderly` - (Optional, Bool) Specifies whether to consume orderly.
  It's only valid when RocketMQ instance version is **5.x**.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.

## Import

The rocketmq consumer group can be imported using the rocketMQ instance ID and group name separated by a slash, e.g.

```bash
$ terraform import huaweicloud_dms_rocketmq_consumer_group.test 8d3c7938-dc47-4937-a30f-c80de381c5e3/group_1
```
