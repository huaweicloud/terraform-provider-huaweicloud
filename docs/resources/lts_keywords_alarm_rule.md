---
subcategory: "Log Tank Service (LTS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_lts_keywords_alarm_rule"
description: |-
  Manages an LTS keywords alarm rule resource within HuaweiCloud.
---

# huaweicloud_lts_keywords_alarm_rule

Manages an LTS keywords alarm rule resource within HuaweiCloud.

## Example Usage

```hcl
variable "log_group_id" {}
variable "log_stream_id" {}

resource "huaweicloud_lts_keywords_alarm_rule" "test" {
  name        = "terraform_test"
  description = "created by terraform"
  alarm_level = "CRITICAL"
  
  keywords_requests {
    keywords               = "terraform_test_keywords"
    condition              = ">"
    number                 = 100
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

* `name` - (Required, String, ForceNew) Specifies the name of the keywords alarm rule.  
  The maximum lanegth is `128` characters, only Chinese characters, letters, digits, hyphens (-) and underscores (_)
  are allowed.  
  The name cannot start with and end with a hyphen or a underscore.  
  Changing this parameter will create a new resource.

* `keywords_requests` - (Required, List) Specifies the keywords requests.
  The [KeywordsRequests](#KeywordsAlarmRule_KeywordsRequests) structure is documented below.

* `frequency` - (Required, List) Specifies the alarm frequency configurations.
  The [Frequency](#KeywordsAlarmRule_Frequency) structure is documented below.

* `alarm_level` - (Required, String) Specifies the alarm level.  
  The value can be: **INFO**, **MINOR**, **MAJOR** and **CRITICAL**.

* `description` - (Optional, String) Specifies the description of the keywords alarm rule.

* `send_notifications` - (Optional, Bool) Specifies whether to send notifications.  
  Defaults to **false**.
  
* `alarm_action_rule_name` - (Optional, String) Specifies the name of the alarm action rule associated with
  the keyword alarm rule.  
  This parameter is available only when `send_notifications` parameter is set to **true**.

  -> This parameter cannot be used together with `notification_save_rule` parameter.
  
* `notification_save_rule` - (Optional, List) Specifies the notification rule.
  The [NotificationRule](#KeywordsAlarmRule_NotificationRule) structure is documented below.
  This parameter is available only when `send_notifications` parameter is set to **true**.

* `trigger_condition_count` - (Optional, Int) Specifies the count to trigger the alarm.
  Defaults to `1`.

* `trigger_condition_frequency` - (Optional, Int) Specifies the frequency to trigger the alarm.
  Defaults to `1`.

* `send_recovery_notifications` - (Optional, Bool) Specifies whether to send recovery notifications.  
  Defaults to **false**

* `recovery_frequency` - (Optional, Int) Specifies the frequency to recover the alarm.
  Defaults to `3`.

* `alarm_rule_alias` - (Optional, String) Specifies the alias name of the keyword alarm rule.  
  The maximum lanegth is `128` characters, only Chinese characters, letters, digits, hyphens (-) and underscores (_)
  are allowed.  
  The name cannot start with and end with a hyphen or a underscore.

* `notification_frequency` - (Optional, Int) Specifies the notification frequency of the keyword alarm rule,
  in minutes.  
  Defaults to `0`, `0` means immediately notification.  
  This parameter is available only when `send_notifications` parameter is set to **true**.
  The valid values are as follows:
  + **0**
  + **5**
  + **10**
  + **15**
  + **30**
  + **60**
  + **180**
  + **360**

* `status` - (Optional, String) Specifies the status. The value can be: **RUNNING** and **STOPPING**.
  Defaults to **RUNNING**.

<a name="KeywordsAlarmRule_KeywordsRequests"></a>
The `KeywordsRequests` block supports:

* `keywords` - (Required, String) Specifies the keywords.

* `condition` - (Required, String) Specifies the keywords request condition.  
  The value can be: **>=**, **<=**, **<** and **>**.

* `number` - (Required, Int) Specifies the line number.

* `log_stream_id` - (Required, String) Specifies the ID of the log stream.

* `log_group_id` - (Required, String) Specifies the ID of the log group.

* `search_time_range_unit` - (Required, String) Specifies the unit of search time range.  
  The value can be: **minute** and **hour**.

* `search_time_range` - (Required, Int) Specifies the search time range.
  + When the `search_time_range_unit` is **minute**, the value ranges from `1` to `60`.
  + When the `search_time_range_unit` is **hour**, the value ranges from `1` to `24`.

* `log_group_name` - (Optional, String) Specifies the name of the log group.

* `log_stream_name` - (Optional, String) Specifies the name of the log stream.

<a name="KeywordsAlarmRule_Frequency"></a>
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
  + When the `fixed_rate_unit` is **hour**, the value ranges from `1` to `24`

<a name="KeywordsAlarmRule_NotificationRule"></a>
The `NotificationRule` block supports:

* `template_name` - (Required, String) Specifies the notification template name.

* `user_name` - (Required, String) Specifies the user name.

* `topics` - (Required, List) Specifies the SMN topics.
  The [Topic](#KeywordsAlarmRule_Topic) structure is documented below.

* `timezone` - (Optional, String) Specifies the timezone.

* `language` - (Optional, String) Specifies the notification language.  
  The value can be **zh-cn** and **en-us**, defaults to **zh-cn**.

<a name="KeywordsAlarmRule_Topic"></a>
The `NotificationRuleTopic` block supports:

* `name` - (Required, String) Specifies the topic name.

* `topic_urn` - (Required, String) Specifies the topic URN.

* `display_name` - (Optional, String) Specifies the display name.
  This will be shown as the sender of the message.

* `push_policy` - (Optional, Int) Specifies the push policy.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.
  
* `domain_id` - The domain ID.

* `created_at` - The creation time of the alarm rule.

* `updated_at` - The last update time of the alarm rule.

* `condition_expression` - The condition expression of the keyword alarm rule.

## Import

The keywords alarm rule can be imported using the `id`, e.g.

```bash
$ terraform import huaweicloud_lts_keywords_alarm_rule.test <id>
```

Note that the imported state may not be identical to your resource definition, due to some attributes missing from the
API response. The missing attributes include: `notification_save_rule.0.user_name`, `notification_save_rule.0.timezone,`.
It is generally recommended running `terraform plan` after importing a certificate.
You can then decide if changes should be applied to the certificate, or the resource definition should be updated to
align with the certificate. Also you can ignore changes as below.

```hcl
resource "huaweicloud_lts_keywords_alarm_rule" "test" {
  ...

  lifecycle {
    ignore_changes = [
      notification_save_rule.0.user_name, notification_save_rule.0.timezone,
    ]
  }
}
```
