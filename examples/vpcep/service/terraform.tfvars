vpc_name                      = "tf_test_vpc"
vpc_cidr                      = "192.168.0.0/16"
subnet_name                   = "tf_test_subnet"
security_group_name           = "tf_test_security_group"
instance_name                 = "tf_test_instance"
endpoint_service_name         = "tf-test-service"
endpoint_service_port_mapping = [
  {
    service_port  = 8080
    terminal_port = 80
  }
]
