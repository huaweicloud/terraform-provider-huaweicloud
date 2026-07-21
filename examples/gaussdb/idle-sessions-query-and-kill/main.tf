# Query and kill idle sessions of a GaussDB instance
data "huaweicloud_gaussdb_key_view_nodes_deliver" "test" {
  instance_id = var.instance_id
}

resource "huaweicloud_gaussdb_idle_sessions_query_and_kill" "test" {
  instance_id  = var.instance_id
  node_id      = data.huaweicloud_gaussdb_key_view_nodes_deliver.test.nodes.0.node_id
  component_id = data.huaweicloud_gaussdb_key_view_nodes_deliver.test.nodes.0.component_id
}
