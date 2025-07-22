vpc_configurations = [
  {
    vpc_name    = "tf_test_source_vpc"
    vpc_cidr    = "192.168.0.0/18"
    subnet_name = "tf_test_source_subnet"
  },
  {
    vpc_name    = "tf_test_target_vpc"
    vpc_cidr    = "192.168.128.0/18"
    subnet_name = "tf_test_target_subnet"
  }
]
peering_connection_name = "tf_test_peering"
