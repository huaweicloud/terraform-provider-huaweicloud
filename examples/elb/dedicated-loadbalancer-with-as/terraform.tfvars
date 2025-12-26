vpc_name                       = "tf_test_vpc"
subnet_name                    = "tf_test_subnet"
security_group_name            = "tf_test_security_group"
instance_name                  = "tf_test_instance"
loadbalancer_name              = "tf_test_dedicated_loadbalancer"
loadbalancer_cross_vpc_backend = true
listener_name                  = "tf_test_dedicated_listener"
pool_name                      = "tf_test_dedicated_pool"
configuration_name             = "tf_test_as_configuration"
configuration_user_data        = <<EOT
#! /bin/bash
echo 'root:$6$V6azyeLwcD3CHlpY$BN3VVq18fmCkj66B4zdHLWevqcxlig' | chpasswd -e
EOT

configuration_disks = [
  {
    size        = 40
    volume_type = "SSD"
    disk_type   = "SYS"
  }
]

group_name            = "tf_test_as_group"
policy_name           = "tf_test_as_policy"
alarm_rule_name       = "tf_test_alarm_rule2"
alarm_rule_conditions = [
  {
    alarm_level         = 2
    period              = 300
    filter              = "max"
    comparison_operator = ">"
    value               = 80
    unit                = "%"
    count               = 1
    metric_name         = "cpu_util"
  }
]
