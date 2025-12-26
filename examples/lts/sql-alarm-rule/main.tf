resource "huaweicloud_smn_topic" "test" {
  name         = var.topic_name
  display_name = "The display name of topic"
}

resource "huaweicloud_lts_group" "test" {
  group_name  = var.group_name
  ttl_in_days = var.group_log_expiration_days
}

resource "huaweicloud_lts_stream" "test" {
  group_id    = huaweicloud_lts_group.test.id
  stream_name = var.stream_name
  ttl_in_days = var.stream_log_expiration_days
}

data "huaweicloud_lts_notification_templates" "test" {
  count = var.notification_template_name != "" ? 0 : 1

  domain_id = var.domain_id
}

resource "huaweicloud_lts_sql_alarm_rule" "test" {
  name                        = var.alarm_rule_name
  condition_expression        = var.alarm_rule_condition_expression
  alarm_level                 = var.alarm_rule_alarm_level
  send_notifications          = true
  trigger_condition_count     = var.alarm_rule_trigger_condition_count
  trigger_condition_frequency = var.alarm_rule_trigger_condition_frequency
  send_recovery_notifications = var.alarm_rule_send_recovery_notifications
  recovery_frequency          = var.alarm_rule_send_recovery_notifications ? var.alarm_rule_recovery_frequency : null
  notification_frequency      = var.alarm_rule_notification_frequency
  alarm_rule_alias            = var.alarm_rule_alias

  sql_requests {
    title                  = var.alarm_rule_request_title
    sql                    = var.alarm_rule_request_sql
    log_group_id           = huaweicloud_lts_group.test.id
    log_stream_id          = huaweicloud_lts_stream.test.id
    search_time_range_unit = var.alarm_rule_request_search_time_range_unit
    search_time_range      = var.alarm_rule_request_search_time_range
    log_group_name         = huaweicloud_lts_group.test.group_name
    log_stream_name        = huaweicloud_lts_stream.test.stream_name
  }

  frequency {
    type = var.alarm_rule_frequency_type
  }

  notification_save_rule {
    template_name = var.notification_template_name!= "" ? var.notification_template_name : try([for v in data.huaweicloud_lts_notification_templates.test[0].templates[*].name :v if v == "sql_template"][0], null)
    user_name     = var.alarm_rule_notification_user_name
    language      = var.alarm_rule_notification_language

    topics {
      name         = huaweicloud_smn_topic.test.name
      topic_urn    = huaweicloud_smn_topic.test.topic_urn
      display_name = huaweicloud_smn_topic.test.display_name
      push_policy  = huaweicloud_smn_topic.test.push_policy
    }
  }
}
