vpc_name            = "tf_test_ges_persistence_vpc"
subnet_name         = "tf_test_ges_persistence_subnet"
security_group_name = "tf_test_ges_persistence_secgroup"
graph_name          = "tf_test_ges_persistence_graph"
keep_backup         = true

graph_tags = {
  env = "production"
}
