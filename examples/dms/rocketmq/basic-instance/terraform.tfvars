vpc_name             = "tf_test_rocketmq_instance"
subnet_name          = "tf_test_rocketmq_instance"
security_group_name  = "tf_test_rocketmq_instance"
instance_name        = "tf_test_basic"
instance_enable_acl  = true
instance_broker_num  = 2
instance_description = "Created by terraform script"

instance_tags = {
  "owner" = "terraform"
}

instance_configs = [
  {
    name  = "fileReservedTime"
    value = "72"
  }
]
