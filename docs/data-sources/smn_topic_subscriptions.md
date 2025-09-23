---
subcategory: "Simple Message Notification (SMN)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_smn_topic_subscriptions"
description: |-
  Use this data source to get a list of SMN topic subscriptions.
---

# huaweicloud_smn_topic_subscriptions

Use this data source to get a list of SMN topic subscriptions.

## Example Usage

```hcl
variable "topic_urn" {}

data "huaweicloud_smn_topic_subscriptions" "test" {
  topic_urn = var.topic_urn
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

* `topic_urn` - (Required, String) Specifies the topic URN.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `subscriptions` - The subscription list.

  The [subscriptions](#subscriptions_struct) structure is documented below.

<a name="subscriptions_struct"></a>
The `subscriptions` block supports:

* `remark` - The subscription remark.

* `status` - The subscription status.

* `filter_polices` - The subscription filter polices.

  The [filter_polices](#subscriptions_filter_polices_struct) structure is documented below.

* `topic_urn` - The topic URN.

* `protocol` - The subscription protocol.

* `subscription_urn` - The subscription URN.

* `owner` - The subscription owner.

* `endpoint` - The subscriptions endpoint.

<a name="subscriptions_filter_polices_struct"></a>
The `filter_polices` block supports:

* `name` - The filter policy name.

* `string_equals` - The string array for exact match.
