---
subcategory: "Simple Message Notification (SMN)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_smn_subscriptions"
description: |-
  Use this data source to get a list of SMN subscriptions.
---

# huaweicloud_smn_subscriptions

Use this data source to get a list of SMN subscriptions.

## Example Usage

```hcl
data "huaweicloud_smn_subscriptions" "test" {}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

* `protocol` - (Optional, String) Specifies the protocol name.
  The enumerated values are **http**, **https**, **sms**, **email**, **functionstage**, **dms**, and **application**.

* `status` - (Optional, String) Specifies the subscription status.
  + **0**: The subscription has not been confirmed.
  + **1**: The subscription has been confirmed.
  + **2**: Confirmation is not required.
  + **3**: The subscription was canceled.
  + **4**: The subscription was deleted.

* `endpoint` - (Optional, String) Specifies the subscription endpoint.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `subscriptions` - The subscription list.

  The [subscriptions](#subscriptions_struct) structure is documented below.

<a name="subscriptions_struct"></a>
The `subscriptions` block supports:

* `status` - The subscription status.

* `filter_polices` - The subscription filter polices.

  The [filter_polices](#subscriptions_filter_polices_struct) structure is documented below.

* `topic_urn` - The topic URN.

* `protocol` - The subscription protocol.

* `subscription_urn` - The subscription URN.

* `owner` - The subscription owner.

* `endpoint` - The subscriptions endpoint.

* `remark` - The subscriptions remark.

<a name="subscriptions_filter_polices_struct"></a>
The `filter_polices` block supports:

* `name` - The filter policy name.

* `string_equals` - The string array for exact match.
