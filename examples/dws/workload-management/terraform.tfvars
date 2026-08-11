vpc_name                = "tf_test_dws_vpc"
vpc_cidr                = "192.168.0.0/16"
subnet_name             = "tf_test_dws_subnet"
security_group_name     = "tf_test_dws_sg"
cluster_name            = "tf_test_cluster"
cluster_admin_user_name = "dbadmin"
cluster_admin_user_pwd  = "YourPassword@123"
workload_queue_name     = "tf_test_queue"

workload_queue_configurations = [
  {
    resource_name  = "cpu_limit"
    resource_value = "10"
  },
  {
    resource_name  = "memory"
    resource_value = "10"
  },
  {
    resource_name  = "tablespace"
    resource_value = "-1"
  },
  {
    resource_name  = "activestatements"
    resource_value = "-1"
  }
]

user_name                = "tf_test_user"
user_password            = "YourPassword@456"
workload_plan_name       = "tf_test_plan"
workload_plan_stage_name = "tf_test_stage"

workload_plan_stage_configurations = [
  {
    resource_name  = "cpu"
    resource_value = 1
  },
  {
    resource_name  = "cpu_limit"
    resource_value = 0
  },
  {
    resource_name  = "memory"
    resource_value = 10
  },
  {
    resource_name  = "concurrency"
    resource_value = 10
  },
  {
    resource_name  = "shortQueryConcurrencyNum"
    resource_value = -1
  }
]

exception_rule_name = "tf_test_exception_rule"

exception_rule_configurations = [
  {
    key   = "action"
    value = "abort"
  },
  {
    key   = "blocktime"
    value = "300"
  },
  {
    key   = "elapsedtime"
    value = "400"
  },
  {
    key   = "allcputime"
    value = "500"
  },
  {
    key   = "cpuskewpercent"
    value = "60"
  },
  {
    key   = "cpuavgpercent"
    value = "70"
  },
]
