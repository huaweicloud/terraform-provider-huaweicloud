---
subcategory: "Log Tank Service (LTS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_lts_sql_alarm_rule"
description: |-
  Manages an LTS SQL alarm rule resource within HuaweiCloud.
---

# huaweicloud_lts_sql_alarm_rule

Manages an LTS SQL alarm rule resource within HuaweiCloud.

## Example Usage

```hcl
variable "log_group_id" {}
variable "log_stream_id" {}

resource "huaweicloud_lts_sql_alarm_rule" "test" {
  name                 = "terraform_test"
  description          = "created by terraform"
  condition_expression = "t>2"
  alarm_level          = "CRITICAL"

  sql_requests {
    title                  = "terraform_test"
    sql                    = "select count(*) as t"
    log_group_id           = var.log_group_id
    log_stream_id          = var.log_stream_id      
    search_time_range_unit = "minute"
    search_time_range      = 5
  }

  frequency {
    type = "HOURLY"
  }
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used. Changing this parameter will create a new resource.

* `name` - (Required, String, ForceNew) Specifies the name of the SQL alarm rule.
  The value can contain no more than 64 characters.
  Changing this parameter will create a new resource.

* `sql_requests` - (Required, List) Specifies the SQL requests.
  The [SQLRequests](#SQLAlarmRule_SQLRequests) structure is documented below.

* `frequency` - (Required, List) Specifies the alarm frequency configurations.
  The [Frequency](#SQLAlarmRule_Frequency) structure is documented below.

* `condition_expression` - (Required, String) Specifies the condition expression.

* `alarm_level` - (Required, String) Specifies the alarm level.
  The value can be: **INFO**, **MINOR**, **MAJOR** and **CRIRICAL**.

* `description` - (Optional, String) Specifies the description of the SQL alarm rule.

* `send_notifications` - (Optional, Bool, ForceNew) Specifies whether to send notifications.
  Changing this parameter will create a new resource.

* `notification_rule` - (Optional, List, ForceNew) Specifies the notification rule.
  Changing this parameter will create a new resource.
  The [NotificationRule](#SQLAlarmRule_NotificationRule) structure is documented below.

* `trigger_condition_count` - (Optional, Int) Specifies the count to trigger the alarm.
  Defaults to `1`.

* `trigger_condition_frequency` - (Optional, Int) Specifies the frequency to trigger the alarm.
  Defaults to `1`.

* `send_recovery_notifications` - (Optional, Bool) Specifies whether to send recovery notifications.

* `recovery_frequency` - (Optional, Int) Specifies the frequency to recover the alarm.
  Defaults to `3`.

* `status` - (Optional, String) Specifies the status. The value can be: **RUNNING** and **STOPPING**.
  Defaults to **RUNNING**

<a name="SQLAlarmRule_SQLRequests"></a>
The `SQLRequests` block supports:

* `title` - (Required, String) Specifies the SQL request title.

* `sql` - (Required, String) Specifies the SQL.

* `log_stream_id` - (Required, String) Specifies the log stream id.

* `log_group_id` - (Required, String) Specifies the log group id.

* `search_time_range_unit` - (Required, String) Specifies the unit of search time range.
  The value can be: **minute** and **hour**.

* `search_time_range` - (Required, Int) Specifies the search time range.
  + When the `search_time_range_unit` is **minute**, the value ranges from `1` to `60`;
  + When the `search_time_range_unit` is **hour**, the value ranges from `1` to `24`;

* `is_time_range_relative` - (Optional, Bool) Specifies the SQL request is relative to time range.

<a name="SQLAlarmRule_Frequency"></a>
The `Frequency` block supports:

* `type` - (Required, String) Specifies the frequency type.
  The value can be: **CRON**, **HOURLY**, **DAILY**, **WEEKLY** and **FIXED_RATE**.

* `cron_expression` - (Optional, String) Specifies the cron expression.
  This parameter is used when `type` is set to **CRON**.

* `hour_of_day` - (Optional, Int) Specifies the hour of day.
  This parameter is used when `type` is set to **DAILY** or **WEEKLY**.
  The value ranges from `0` to `23`.

* `day_of_week` - (Optional, Int) Specifies the day of week.
  This parameter is used when `type` is set to **WEEKLY**.
  The value ranges from `1` to `7`. `1` means Sunday.

* `fixed_rate_unit` - (Optional, String) Specifies the unit of fixed rate.
  The value can be: **minute** and **hour**.

* `fixed_rate` - (Optional, Int) Specifies the unit fixed rate.
  This parameter is used when `type` is set to **FIXED_RATE**.
  + When the `fixed_rate_unit` is **minute**, the value ranges from `1` to `60`.
  + When the `fixed_rate_unit` is **hour**, the value ranges from `1` to `24`.

<a name="SQLAlarmRule_NotificationRule"></a>
The `NotificationRule` block supports:

* `template_name` - (Required, String, ForceNew) Specifies the notification template name.
  Changing this parameter will create a new resource.

* `language` - (Required, String, ForceNew) Specifies the notification language.
  The value can be **zh-cn** and **en-us**.
  Changing this parameter will create a new resource.

* `user_name` - (Required, String, ForceNew) Specifies the user name.
  Changing this parameter will create a new resource.

* `topics` - (Required, List, ForceNew) Specifies the SMN topics.
  The [Topic](#SQLAlarmRule_Topic) structure is documented below.
  Changing this parameter will create a new resource.

* `timezone` - (Optional, String, ForceNew) Specifies the timezone.
  Changing this parameter will create a new resource.

<a name="SQLAlarmRule_Topic"></a>
The `NotificationRuleTopic` block supports:

* `name` - (Required, String, ForceNew) Specifies the topic name.
  Changing this parameter will create a new resource.

* `topic_urn` - (Required, String, ForceNew) Specifies the topic URN.
  Changing this parameter will create a new resource.

* `display_name` - (Optional, String, ForceNew) Specifies the display name.
  This will be shown as the sender of the message.
  Changing this parameter will create a new resource.

* `push_policy` - (Optional, String, ForceNew) Specifies the push policy.
  Changing this parameter will create a new resource.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.

* `domain_id` - The domain ID.

* `created_at` - The creation time of the alarm rule.

* `updated_at` - The last update time of the alarm rule.

## Import

The sql alarm rule can be imported using the `id`, e.g.

```bash
$ terraform import huaweicloud_lts_sql_alarm_rule.test <id>
```

Note that the imported state may not be identical to your resource definition, due to some attributes missing from the
API response. The missing attributes include: `notification_rule`.
It is generally recommended running `terraform plan` after importing a certificate.
You can then decide if changes should be applied to the certificate, or the resource definition should be updated to
align with the certificate. Also you can ignore changes as below.

```hcl
resource "huaweicloud_lts_sql_alarm_rule" "test" {
  ...

  lifecycle {
    ignore_changes = [
      notification_rule,
    ]
  }
}
```
