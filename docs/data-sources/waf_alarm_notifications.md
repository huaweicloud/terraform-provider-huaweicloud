---
subcategory: "Web Application Firewall (WAF)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_waf_alarm_notifications"
description: |-
  Use this data source to get the list of alarm notifications.
---

# huaweicloud_waf_alarm_notifications

Use this data source to get the list of alarm notifications.

## Example Usage

```hcl
data "huaweicloud_waf_alarm_notifications" "test" {}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

* `enterprise_project_id` - (Optional, String) Specifies the enterprise project ID to which the alarm notification
  belongs.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `items` - The list of alarm notifications.

  The [items](#items_struct) structure is documented below.

<a name="items_struct"></a>
The `items` block supports:

* `id` - The ID of the alarm notification.

* `name` - The name of the alarm notification.

* `enabled` - Whether to enable the  alarm notification.
  The value can be **true** or  **false**.

* `notice_class` - The type of the  alarm notification.

* `times` - The alarm times of the notification alarm.

* `sendfreq` - The alarm frequency of the alarm notification, in minutes.

* `topic_urn` - The theme URN of the SMN associated with the alarm notification.

* `locale` - The language type.

* `threat` - The type of the event which triggered the alarm.

* `nearly_expired_time` - The advance notification days.

* `enterprise_project_id` - The enterprise project ID to which the alarm notification belongs.

* `is_all_enterprise_project` - Whether all enterprise projects are involved.

* `update_time` - The latest update time of the alarm notification.
