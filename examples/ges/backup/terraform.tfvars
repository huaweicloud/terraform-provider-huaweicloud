vpc_name              = "tf_test_ges_vpc"
vpc_cidr              = "192.168.0.0/16"
subnet_name           = "tf_test_ges_subnet"
subnet_cidr           = "192.168.0.0/24"
gateway_ip            = "192.168.0.1"
security_group_name   = "tf_test_ges_secgroup"
graph_name            = "tf_test_ges_graph"
graph_size_type_index = "1"
graph_cpu_arch        = "x86_64"
graph_crypt_algorithm = "generalCipher"
graph_enable_https    = false
graph_tags            = {
  key = "val"
  foo = "bar"
}
