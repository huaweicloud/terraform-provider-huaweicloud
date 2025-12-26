vpc_name                       = "tf_test_sdrs_protection_instance_vpc"
subnet_name                    = "tf_test_sdrs_protection_instance_subnet"
security_group_name            = "tf_test_sdrs_protection_instance_secgroup"
ecs_instance_name              = "tf_test_sdrs_protection_instance_ecs_instance"
protection_group_name          = "tf_test_sdrs_protection_instance_group"
protected_instance_name        = "tf_test_sdrs_protected_instance"
delete_target_server           = true
delete_target_eip              = true
protected_instance_description = "Created by terraform script"
protected_instance_tags        = {
  foo = "bar"
  key = "value"
}
