---
subcategory: "Distributed Message Service (DMS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_dms_rocketmq_consumers"
description: |-
  Use this data source to get the list of RocketMQ consumers.
---

# huaweicloud_dms_rocketmq_consumers

Use this data source to get the list of RocketMQ consumers.

## Example Usage

```hcl
variable "instance_id" {}
variable "group" {}

data "huaweicloud_dms_rocketmq_consumers" "test" {
  instance_id = var.instance_id
  group       = var.group
  is_detail   = true
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

* `instance_id` - (Required, String) Specifies the instance ID.

* `group` - (Required, String) Specifies the consumer group name.

* `is_detail` - (Optional, Bool) Specifies whether to query the consumer details.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `clients` - Indicates the list of consumer subscription details.

  The [clients](#clients_struct) structure is documented below.

* `online` - Indicates whether the consumer group is online.

* `subscription_consistency` - Indicates whether subscriptions are consistent.

<a name="clients_struct"></a>
The `clients` block supports:

* `subscriptions` - Indicates the subscription list.

  The [subscriptions](#clients_subscriptions_struct) structure is documented below.

* `language` - Indicates the client language.

* `version` - Indicates the client version.

* `client_id` - Indicates the client ID.

* `client_address` - Indicates the client address.

<a name="clients_subscriptions_struct"></a>
The `subscriptions` block supports:

* `topic` - Indicates the name of the subscribed topic.

* `type` - Indicates the subscription type. The value can be **TAG** and **SQL92**.

* `expression` - Indicates the subscription tag.
