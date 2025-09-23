---
subcategory: "Cloud Eye (CES)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_ces_one_click_alarm"
description: |-
  Manages a CES one-click alarm resource within HuaweiCloud.
---

# huaweicloud_ces_one_click_alarm

Manages a CES one-click alarm resource within HuaweiCloud.

## Example Usage

```hcl
variable "notification_object" {}

resource "huaweicloud_ces_one_click_alarm" "test" {
  one_click_alarm_id = "OBSSystemOneClickAlarm"

  dimension_names {
    metric = ["bucket_name"]
    event  = true
  }

  alarm_notifications {
    type = "notification"

    notification_list = [
      var.notification_object
    ]
  }

  ok_notifications {
    type = "notification"

    notification_list = [
      var.notification_object
    ]
  }

  notification_enabled    = true
  notification_begin_time = "00:00"
  notification_end_time   = "23:59"
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used.
  Changing this creates a new resource.

* `one_click_alarm_id` - (Required, String, NonUpdatable) Specifies the default one-click monitoring ID.
  The value can be queried from the CES one-click alarms data source.

* `dimension_names` - (Required, List, NonUpdatable) Specifies dimensions in metric and event alarm rules that have
  one-click monitoring enabled.

  The [dimension_names](#DimensionNames) structure is documented below.

* `notification_enabled` - (Required, Bool) Specifies whether to enable the alarm notification.

* `alarm_notifications` - (Optional, List) Specifies the action to be triggered by an alarm.
  + If the value of `notification_enabled` is **false**, this parameter should not be set.
  + If the value of `notification_enabled` is **true**, this parameter is required.

  The [alarm_notifications](#Notifications) structure is documented below.

* `ok_notifications` - (Optional, List) Specifies the action to be triggered after an alarm is cleared.
  + If the value of `notification_enabled` is **false**, this parameter should not be set.

  The [ok_notifications](#Notifications) structure is documented below.

* `notification_begin_time` - (Optional, String) Specifies the time when the alarm notification was enabled.

* `notification_end_time` - (Optional, String) Specifies the time when the alarm notification was disabled.

<a name="DimensionNames"></a>
The `dimension_names` block supports:

* `event` - (Optional, Bool, NonUpdatable) Specifies whether to enable the event alarm rules.

* `metric` - (Optional, List, NonUpdatable) Specifies dimensions in metric alarm rules that have one-click monitoring enabled.

<a name="Notifications"></a>
The `alarm_notifications` block or `ok_notifications` block supports:

* `type` - (Required, String) Specifies the notification type.
  The value can be **notification** or **contact**.

* `notification_list` - (Required, List) Specifies the list of objects to be notified if the alarm status changes.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.

* `description` - The supplementary information about one-click monitoring.

* `enabled` - Whether the one-click monitoring is enabled.

* `namespace` - The metric namespace.

## Import

The one-click alarm can be imported using `id`, e.g.

```bash
$ terraform import huaweicloud_ces_one_click_alarm.test <id>
```

Note that the imported state may not be identical to your resource definition, due to some attributes missing from the
API response, security or some other reason.
The missing attributes include: `one_click_alarm_id`, `dimension_names`, `notification_enabled`, `alarm_notifications`,
`ok_notifications`, `notification_begin_time`, `notification_end_time`.
It is generally recommended running `terraform plan` after importing the one-click alarm.
You can then decide if changes should be applied to the one-click alarm, or the resource definition should be updated to
align with the one-click alarm. Also you can ignore changes as below.

```hcl
resource "huaweicloud_ces_one_click_alarm" "test" {
    ...

  lifecycle {
    ignore_changes = [
      one_click_alarm_id, dimension_names, notification_enabled, alarm_notifications,
      ok_notifications, notification_begin_time, notification_end_time
    ]
  }
}
```
