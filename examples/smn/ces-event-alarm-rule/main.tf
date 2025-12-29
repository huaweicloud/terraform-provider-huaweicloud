resource "huaweicloud_smn_topic" "test" {
  name                  = var.smn_topic_name
  enterprise_project_id = var.enterprise_project_id
}

resource "huaweicloud_ces_alarmrule" "test" {
  alarm_name            = var.alarm_rule_name
  alarm_description     = var.alarm_rule_description
  alarm_action_enabled  = var.alarm_action_enabled
  alarm_enabled         = var.alarm_enabled
  alarm_type            = var.alarm_type
  enterprise_project_id = var.enterprise_project_id

  metric {
    namespace = "SYS.SMN"
  }

  dynamic "condition" {
    for_each = var.alarm_rule_conditions

    content {
      metric_name         = condition.value.metric_name
      period              = condition.value.period
      filter              = condition.value.filter
      comparison_operator = condition.value.comparison_operator
      value               = condition.value.value
      unit                = condition.value.unit
      count               = condition.value.count
      suppress_duration   = condition.value.suppress_duration
      alarm_level         = condition.value.alarm_level
    }
  }

  dynamic "resources" {
    for_each = var.alarm_rule_resource

    content {
      dimensions {
        name  = resources.value.name
        value = resources.value.value
      }
    }
  }

  alarm_actions {
    type = "notification"

    notification_list = [
      huaweicloud_smn_topic.test.topic_urn
    ]
  }

  notification_begin_time = var.alarm_rule_notification_begin_time
  notification_end_time   = var.alarm_rule_notification_end_time
  effective_timezone      = var.alarm_rule_effective_timezone
}
