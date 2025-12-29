smn_topic_name         = "tf_test_topic"
alarm_rule_name        = "tf_test_alarm_rule"
alarm_rule_description = "Monitor SMN topic events"
alarm_type             = "ALL_INSTANCE"

alarm_rule_conditions = [
  {
    metric_name         = "email_total_count"
    period              = "1"
    filter              = "average"
    comparison_operator = ">="
    value               = "80"
    count               = "3"
    unit                = "count"
    alarm_level         = "3"
  }
]

alarm_rule_resource = [
  {
    name = "topic_id"
  }
]
