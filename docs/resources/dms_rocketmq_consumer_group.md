---
subcategory: "Distributed Message Service (DMS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_dms_rocketmq_consumer_group"
description: ""
---

# huaweicloud_dms_rocketmq_consumer_group

Manages DMS RocketMQ consumer group resources within HuaweiCloud.

## Example Usage

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

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used. Changing this parameter will create a new resource.

* `instance_id` - (Required, String, ForceNew) Specifies the ID of the rocketMQ instance.

  Changing this parameter will create a new resource.

* `brokers` - (Required, List, ForceNew) Specifies the list of associated brokers of the consumer group.

  Changing this parameter will create a new resource.

* `name` - (Required, String, ForceNew) Specifies the name of the consumer group.

  Changing this parameter will create a new resource.

* `retry_max_times` - (Required, Int) Specifies the maximum number of retry times.

* `enabled` - (Optional, Bool) Specifies the consumer group is enabled or not. Default to true.

* `broadcast` - (Optional, Bool) Specifies whether to broadcast of the consumer group.

* `description` - (Optional, String) Specifies the description of the consumer group.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.

## Import

The rocketmq consumer group can be imported using the rocketMQ instance ID and group name separated by a slash, e.g.

```
$ terraform import huaweicloud_dms_rocketmq_consumer_group.test 8d3c7938-dc47-4937-a30f-c80de381c5e3/group_1
```
