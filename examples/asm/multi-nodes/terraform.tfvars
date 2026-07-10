# ASM mesh configuration
mesh_name    = "multi-nodes-mesh"
mesh_version = "1.18.7-r7"

# CCE cluster and multiple nodes
cluster_id = "your-cce-cluster-id"
node_ids   = ["your-cce-node-1", "your-cce-node-2"]

namespace_name = "multi-nodes-test"

# Tags for the mesh
tags = {
  foo = "bar"
  key = "value"
}
