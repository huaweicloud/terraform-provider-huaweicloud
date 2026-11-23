bandwidth_name  = "tf_test_bandwidth"
alarm_rule_name = "tf_test_alarm_rule"

rule_conditions = [
  {
    period              = 300
    filter              = "average"
    comparison_operator = ">"
    value               = 3
    count               = 5
    unit                = "Mbit/s"
  }
]

interval_alarm_conditions = [
  {
    lower_bound = "3"
    upper_bound = "10"
    operation   = "ADD"
    size        = 1
    limits      = 20
  },
  {
    lower_bound = "10"
    upper_bound = "20"
    operation   = "SET"
    size        = 8
  }
]

bandwidth_policy_name = "tf_test_bandwidthn_policy"
