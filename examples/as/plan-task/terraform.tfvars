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

group_name = "tf_test_group"
task_name  = "tf_test_task"
