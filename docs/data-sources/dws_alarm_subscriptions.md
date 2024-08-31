---
subcategory: "GaussDB(DWS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_dws_alarm_subscriptions"
description: |-
  Use this data source to query the list of alarm subscriptions within HuaweiCloud.
---

# huaweicloud_dws_alarm_subscriptions

Use this data source to query the list of alarm subscriptions within HuaweiCloud.

## Example Usage

```hcl
data "huaweicloud_dws_alarm_subscriptions" "test" {}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `subscriptions` - The list of the alarm subscriptions.

  The [subscriptions](#subscriptions_struct) structure is documented below.

<a name="subscriptions_struct"></a>
The `subscriptions` block supports:

* `id` - The ID of the alarm subscription.

* `name` - The name of the alarm subscription.

* `enable` - Whether alarm subscription is enabled.
  + **1**: Enabled.
  + **0**: Disabled.

* `alarm_level` - The level of the alarm subscription.
  + **urgent**
  + **important**
  + **minor**
  + **prompt**.

* `notification_target_type` - The type of notification topic corresponding to the alarm subscription.
  + **SMN**

* `notification_target` - The address of notification topic corresponding to the alarm subscription.

* `notification_target_name` - The name of notification topic corresponding to the alarm subscription.

* `time_zone` - The time zone of the alarm subscription.

* `language` - The language of the alarm subscription.
