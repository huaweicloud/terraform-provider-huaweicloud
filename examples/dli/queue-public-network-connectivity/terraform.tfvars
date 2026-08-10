vpc_name                   = "tf_test_dli_vpc"
vpc_cidr                   = "192.168.0.0/16"
subnet_name                = "tf_test_dli_subnet"
elastic_resource_pool_name = "tf_test_dli_pool"
elastic_resource_pool_cidr = "172.16.0.0/18"
queue_name                 = "tf_test_dli_queue"
datasource_connection_name = "tf_test_dli_conn"

datasource_connection_routes = [
  {
    name = "tf_test_dli_route"
    cidr = "14.17.72.0/24"
  }
]

eip_bandwidth_name = "tf_test_dli_eip_bw"
nat_gateway_name   = "tf_test_dli_nat"
