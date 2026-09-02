vpc_name            = "tf_test_ges_ha_vpc"
subnet_name         = "tf_test_ges_ha_subnet"
security_group_name = "tf_test_ges_ha_secgroup"
graph_name          = "tf_test_ges_ha_graph"
enable_multi_az     = true
enable_https        = false

graph_tags = {
  env = "production"
}
