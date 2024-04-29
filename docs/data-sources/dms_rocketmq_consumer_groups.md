---
subcategory: "Distributed Message Service (DMS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_dms_rocketmq_consumer_groups"
description: ""
---

# huaweicloud_dms_rocketmq_consumer_groups

Use this data source to get the list of DMS rocketMQ consumer groups.

## Example Usage

```hcl
var "instance_id" {}
data "huaweicloud_dms_rocketmq_consumer_groups" "test" {
  instance_id     = var.instance_id
  enabled         = false
  broadcast       = true
  name            = "consumergroup002"
  retry_max_times = 16
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used. Changing this parameter will create a new resource.

* `instance_id` - (Required, String) Specifies the ID of the rocketMQ instance.

* `name` - (Optional, String) Specifies the name of the consumer group.

* `retry_max_times` - (Optional, Int) Specifies the maximum number of retry times.

* `enabled` - (Optional, Bool) Specifies the consumer group is enabled or not. Defaults to **true**.

* `broadcast` - (Optional, Bool) Specifies whether to broadcast the consumer group. Defaults to **false**.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `groups` - The list of consumer groups.
  The [groups](#DMS_rockermq_consumer_groups) structure is documented below.

<a name="DMS_rockermq_consumer_groups"></a>
The `groups` block supports:

* `name` - Indicates the name of the consumer group.

* `brokers` - Indicates the list of associated brokers of the consumer group.

* `retry_max_times` - Indicates the maximum number of retry times.

* `enabled` - Indicates the consumer group is enabled or not.

* `broadcast` - Indicates whether to broadcast the consumer group.

* `description` - Indicates the description of the consumer group.
