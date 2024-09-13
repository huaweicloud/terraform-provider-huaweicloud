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
  The value can contain no more than 64 characters.
  Changing this parameter will create a new resource.

* `keywords_requests` - (Required, List) Specifies the keywords requests.
  The [KeywordsRequests](#KeywordsAlarmRule_KeywordsRequests) structure is documented below.

* `frequency` - (Required, List) Specifies the alarm frequency configurations.
  The [Frequency](#KeywordsAlarmRule_Frequency) structure is documented below.

* `alarm_level` - (Required, String) Specifies the alarm level.
  The value can be: **INFO**, **MINOR**, **MAJOR** and **CRITICAL**.

* `description` - (Optional, String) Specifies the description of the keywords alarm rule.

* `send_notifications` - (Optional, Bool, ForceNew) Specifies whether to send notifications.
  Changing this parameter will create a new resource.

* `notification_rule` - (Optional, List, ForceNew) Specifies the notification rule.
  The [NotificationRule](#KeywordsAlarmRule_NotificationRule) structure is documented below.
  Changing this parameter will create a new resource.

* `trigger_condition_count` - (Optional, Int) Specifies the count to trigger the alarm.
  Defaults to `1`.

* `trigger_condition_frequency` - (Optional, Int) Specifies the frequency to trigger the alarm.
  Defaults to `1`.

* `send_recovery_notifications` - (Optional, Bool) Specifies whether to send recovery notifications.

* `recovery_frequency` - (Optional, Int) Specifies the frequency to recover the alarm.
  Defaults to `3`.

* `status` - (Optional, String) Specifies the status. The value can be: **RUNNING** and **STOPPING**.
  Defaults to **RUNNING**.

<a name="KeywordsAlarmRule_KeywordsRequests"></a>
The `KeywordsRequests` block supports:

* `keywords` - (Required, String) Specifies the keywords.

* `condition` - (Required, String) Specifies the keywords request condition.
  The value can be: **>=**, **<=**, **<** and **>**.

* `number` - (Required, Int) Specifies the line number.

* `log_stream_id` - (Required, String) Specifies the log stream id.

* `log_group_id` - (Required, String) Specifies the log group id.

* `search_time_range_unit` - (Required, String) Specifies the unit of search time range.
  The value can be: **minute** and **hour**.

* `search_time_range` - (Required, Int) Specifies the search time range.
  + When the `search_time_range_unit` is **minute**, the value ranges from `1` to `60`.
  + When the `search_time_range_unit` is **hour**, the value ranges from `1` to `24`.

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

* `template_name` - (Required, String, ForceNew) Specifies the notification template name.
  Changing this parameter will create a new resource.

* `user_name` - (Required, String, ForceNew) Specifies the user name.
  Changing this parameter will create a new resource.

* `topics` - (Required, List, ForceNew) Specifies the SMN topics.
  The [Topic](#KeywordsAlarmRule_Topic) structure is documented below.
  Changing this parameter will create a new resource.

* `timezone` - (Optional, String, ForceNew) Specifies the timezone.
  Changing this parameter will create a new resource.

* `language` - (Optional, String, ForceNew) Specifies the notification language.
  The value can be **zh-cn** and **en-us**.
  Changing this parameter will create a new resource.

<a name="KeywordsAlarmRule_Topic"></a>
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

The keywords alarm rule can be imported using the `id`, e.g.

```bash
$ terraform import huaweicloud_lts_keywords_alarm_rule.test <id>
```

Note that the imported state may not be identical to your resource definition, due to some attributes missing from the
API response. The missing attributes include: `notification_rule`.
It is generally recommended running `terraform plan` after importing a certificate.
You can then decide if changes should be applied to the certificate, or the resource definition should be updated to
align with the certificate. Also you can ignore changes as below.

```hcl
resource "huaweicloud_lts_keywords_alarm_rule" "test" {
  ...

  lifecycle {
    ignore_changes = [
      notification_rule,
    ]
  }
}
```
