---
subcategory: "Cloud Eye (CES)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_ces_resource_group_alarm_template_async_associate"
description: |-
  Manages a CES resource group alarm template async association resource within HuaweiCloud.
---

# huaweicloud_ces_resource_group_alarm_template_async_associate

Manages a CES resource group alarm template async association resource within HuaweiCloud.

## Example Usage

```hcl
variable "group_id" {}
variable "template_id" {}
variable "notification_policy_id" {}

resource "huaweicloud_ces_resource_group_alarm_template_async_associate" "test" {
  group_id                = var.group_id
  template_ids            = [var.template_id]
  notification_enabled    = true
  enterprise_project_id   = "0"
  notification_manner     = "NOTIFICATION_POLICY"
  notification_policy_ids = [var.notification_policy_id]
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used. Changing this parameter will create a new resource.

* `group_id` - (Required, String, NonUpdatable) Specifies the resource group ID.

* `template_ids` - (Required, List) Specifies the list of alert template IDs.

* `notification_enabled` - (Required, Bool) Specifies whether to enable alert notifications.
  The value can be **true** or **false**.

* `alarm_notifications` - (Optional, List) Specifies the list of alert trigger notifications.
The [alarm_notifications](#notifications_struct) structure is documented below.

* `ok_notifications` - (Optional, List) Specifies the list of alert recovery notifications.
The [ok_notifications](#notifications_struct) structure is documented below.

* `notification_begin_time` - (Optional, String) Specifies the time when the alert notification was enabled.

* `notification_end_time` - (Optional, String) Specifies the time when the alert notification was closed.

* `effective_timezone` - (Optional, String) Specifies the time zone.
  Use a format like this: **GMT-08:00**, **GMT+08:00**, or **GMT+0:00**.

* `enterprise_project_id` - (Optional, String) Specifies the enterprise project ID.

* `notification_manner` - (Optional, String) Specifies the notification manner.
  The valid values are as follows:
  + **NOTIFICATION_GROUP**: Notification groups.
  + **TOPIC_SUBSCRIPTION**: Topic subscriptions.
  + **NOTIFICATION_POLICY**: Notification policies.

* `notification_policy_ids` - (Optional, List) Specifies the list of associated notification policy IDs.

<a name="notifications_struct"></a>
The `alarm_notifications` or `ok_notifications` block supports:

* `type` - (Required, String) Specifies the notification type.
  The valid values are as follows:
  + **notification**: SMN notification.
  + **contact**: Cloud account contact.
  + **contactGroup**: Notification group.
  + **autoscaling**: AS notification, only used in Auto Scaling.

* `notification_list` - (Required, List) Specifies the list of recipients to be notified when the alarm status changes.
  + When `type` is **notification**, the value cannot be empty.
  + When `type` is **autoscaling**, the value must be **[]**.
  + When `type` is **autoscaling**, the value must be **[]**.
  + When `notification_enabled` is **true**, either `alarm_notifications` or `ok_notifications` must be non-empty.
  + When both `alarm_notifications` and `ok_notifications` are present, the `notification_list` values must be consistent.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID, which equals to `group_id`.
