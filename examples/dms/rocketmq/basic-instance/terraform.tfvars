vpc_name             = "tf_test_vpc"
subnet_name          = "tf_test_subnet"
security_group_name  = "tf_test_security_group"
instance_name        = "tf_test_instance"
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
