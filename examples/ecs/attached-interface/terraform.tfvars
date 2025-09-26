vpc_name = "tf_test_ecs_instace"
subnet_configurations = [
  {
    subnet_name = "tf_test_main"
  },
  {
    subnet_name = "tf_test_standby"
  },
]
security_group_name     = "tf_test_ecs_instace"
instance_name           = "tf_test_ecs_instace"
instance_admin_password = "YourPassword!"
