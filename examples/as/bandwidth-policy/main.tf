resource "huaweicloud_vpc_bandwidth" "test" {
  name                  = var.bandwidth_name
  size                  = var.bandwidth_size
  bandwidth_type        = var.bandwidth_type
  enterprise_project_id = var.enterprise_project_id
  charging_mode         = var.charging_mode
  period_unit           = var.period_unit
  period                = var.period
  auto_renew            = var.auto_renew
}

resource "huaweicloud_ces_alarmrule" "test" {
  alarm_name            = var.alarm_rule_name
  alarm_description     = var.alarm_rule_description
  alarm_enabled         = var.alarm_enabled
  alarm_action_enabled  = var.alarm_action_enabled
  enterprise_project_id = var.enterprise_project_id

  metric {
    namespace   = "SYS.VPC"
    metric_name = "downstream_bandwidth"

    dimensions {
      name  = "bandwidth_id"
      value = huaweicloud_vpc_bandwidth.test.id
    }
  }

  dynamic "condition" {
    for_each = var.rule_conditions

    content {
      period              = condition.value.period
      filter              = condition.value.filter
      comparison_operator = condition.value.comparison_operator
      value               = condition.value.value
      count               = condition.value.count
      unit                = condition.value.unit
      alarm_level         = condition.value.alarm_level
      suppress_duration   = condition.value.suppress_duration
    }
  }

  alarm_actions {
    type              = "autoscaling"
    notification_list = []
  }
}

resource "huaweicloud_as_bandwidth_policy" "test" {
  scaling_policy_name = var.bandwidth_policy_name
  scaling_policy_type = var.bandwidth_policy_type
  bandwidth_id        = huaweicloud_vpc_bandwidth.test.id
  alarm_id            = huaweicloud_ces_alarmrule.test.id
  description         = var.bandwidth_policy_description
  cool_down_time      = var.cool_down_time

  dynamic "interval_alarm_actions" {
    for_each = var.interval_alarm_conditions

    content {
      lower_bound = interval_alarm_actions.value.lower_bound
      upper_bound = interval_alarm_actions.value.upper_bound
      operation   = interval_alarm_actions.value.operation
      size        = interval_alarm_actions.value.size
      limits      = interval_alarm_actions.value.limits
    }
  }
}
