---
subcategory: "Application Operations Management (AOM 2.0)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_aom_alarm_notified_histories"
description: |-
  Use this data source to query the notification histories of an alarm event within HuaweiCloud.
---

# huaweicloud_aom_alarm_notified_histories

Use this data source to query the notification histories of an alarm event within HuaweiCloud.

## Example Usage

```hcl
variable "event_sn" {}

data "huaweicloud_aom_alarm_notified_histories" "test" {
  event_sn = var.event_sn
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region where the alarm notification histories are located.  
  If omitted, the provider-level region will be used.

* `event_sn` - (Required, String) Specifies the serial number of the alarm event.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `notified_histories` - The list of notified histories.  
  The [notified_histories](#aom_notified_histories) structure is documented below.

<a name="aom_notified_histories"></a>
The `notified_histories` block supports:

* `event_sn` - The serial number of the alarm event.

* `notifications` - The list of notification results that associated the event.  
  The [notifications](#aom_notified_histories_notifications) structure is documented below.

<a name="aom_notified_histories_notifications"></a>
The `notifications` block supports:

* `action_rule` - The name of the alarm notification rule.

* `notifier_channel` - The notification channel type.

* `smn_channel` - The result detail of the notification.
  The [smn_channel](#aom_notified_histories_notifications_smn_channel) structure is documented below.

<a name="aom_notified_histories_notifications_smn_channel"></a>
The `smn_channel` block supports:

* `sent_time` - The timestamp when the notification was sent.

* `smn_notified_history` - The list of smn notification that associated the event.
    The [smn_notified_history](#aom_notified_histories_notifications_smn_channel_info) structure is documented below.

* `smn_request_id` - The request ID of the notification detail.

* `smn_response_body` - The response body of the notification detail.

* `smn_response_code` - The response code of the notification detail.

* `smn_topic` - The SMN topic used for notification.

<a name="aom_notified_histories_notifications_smn_channel_info"></a>
The `smn_notified_history` block supports:

* `smn_notified_content` - The content of the notification.

* `smn_subscription_status` - The subscription status of the notification.

* `smn_subscription_type` - The subscription type of the notification.
