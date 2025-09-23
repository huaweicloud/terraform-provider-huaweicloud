vpc_name            = "tf_test_vpc"
subnet_name         = "tf_test_subnet"
security_group_name = "tf_test_security_group"
keypair_name        = "tf_test_keypair"
configuration_name  = "tf_test_configuration"

disk_configurations = [
  {
    disk_type   = "SYS"
    volume_type = "SSD"
    volume_size = 40
  }
]

group_name      = "tf_test_group"
topic_name      = "tf_test_topic"
alarm_rule_name = "tf_test_alarm_rule"

rule_conditions = [
  {
    metric_name         = "cpu_util"
    period              = 300
    filter              = "average"
    comparison_operator = ">"
    value               = 80
    count               = 1
  }
]

scaling_up_policy_name   = "tf_test_scaling_up_policy"
scaling_down_policy_name = "tf_test_scaling_down_policy"
