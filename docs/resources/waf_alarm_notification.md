---
subcategory: "Web Application Firewall (WAF)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_waf_alarm_notification"
description: |-
  Manages a WAF alarm notification resource within HuaweiCloud.
---

# huaweicloud_waf_alarm_notification

Manages a WAF alarm notification resource within HuaweiCloud.

-> All WAF resources depend on WAF instances, and the WAF instances need to be purchased before they can be used.

## Example Usage

```hcl
variable "name" {}
variable "topic_urn" {}
variable "notice_class" {}
variable "enterprise_project_id" {}

resource "huaweicloud_waf_alarm_notification" "test" {
  name                  = var.name
  topic_urn             = var.topic_urn
  notice_class          = var.notice_class
  enterprise_project_id = var.enterprise_project_id
  enabled               = true
  sendfreq              = 5
  locale                = "zh-cn"
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used. Changing this parameter will create a new resource.

* `name` - (Required, String) Specifies the name of the alarm notification.

* `topic_urn` - (Required, String) Specifies the topic URN of the SMN.

* `notice_class` - (Required, String) Specifies the type of the alarm notification.

* `enterprise_project_id` - (Required, String, NonUpdatable) Specifies the enterprise project ID.

* `enabled` - (Optional, Bool) Specifies whether to enable the alarm notification. Defaults to **false**.

* `sendfreq` - (Optional, Int) Specifies the time interval, in minutes.
  If the notification type is protection event, an alarm notification is sent when the number of attacks is greater than
  or equal to the threshold within the specified interval.
  The value can be `5`, `15`, `30`, `60`, `120`, `360`, `720`, or `1440`.

  When the notification type is certificate expiration, this parameter indicates the interval for sending an alarm
  notification. The value can be `1` day or `1` week (converted into minutes).

* `locale` - (Optional, String) Specifies the language. The value can be `zh-cn` or `en-us`.

* `times` - (Optional, Int) Specifies the threshold of attack. This parameter is mandatory when notification type is set
  to protection event. WAF will report a notification when the number of attacks reaches the configured threshold.

* `threat` - (Optional, List, NonUpdatable) Specifies the event type.

* `nearly_expired_time` - (Optional, String) Specifies number of days in advance for notification.
  This parameter is mandatory when the notification type is certificate expiration notification.

* `is_all_enterprise_project` - (Optional, Bool) Specifies whether all enterprise projects are involved.
  Defaults to **false**.

* `description` - (Optional, String) Specifies the description of the alarm notification.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID (alarm notification ID).

* `update_time` - The update time, in milliseconds.

## Import

The resource can be imported using the `id` and `enterprise_project_id`, separated by a slash, e.g.

```bash
$ terraform import huaweicloud_waf_alarm_notification.test <id>/<enterprise_project_id>
```
