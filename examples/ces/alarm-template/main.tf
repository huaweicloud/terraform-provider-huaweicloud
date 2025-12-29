resource "huaweicloud_ces_alarm_template" "test" {
  name        = var.alarm_template_name
  description = var.alarm_template_description

  dynamic "policies" {
    for_each = var.alarm_template_policies

    content {
      namespace           = policies.value.namespace
      metric_name         = policies.value.metric_name
      period              = policies.value.period
      filter              = policies.value.filter
      comparison_operator = policies.value.comparison_operator
      count               = policies.value.count
      suppress_duration   = policies.value.suppress_duration
      value               = policies.value.value
      alarm_level         = policies.value.alarm_level
      unit                = policies.value.unit
      dimension_name      = policies.value.dimension_name
      hierarchical_value {
        critical = policies.value.hierarchical_value.critical
        major    = policies.value.hierarchical_value.major
        minor    = policies.value.hierarchical_value.minor
        info     = policies.value.hierarchical_value.info
      }
    }
  }
}
