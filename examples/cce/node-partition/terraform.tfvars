vpc_name                     = "tf_test_vpc"
subnet_name                  = "tf_test_subnet"
eni_subnet_name              = "tf_test_eni_subnet"
cluster_name                 = "tf-test-cluster"
node_partition               = "center"
node_flavor_id               = "c7n.large.2"
node_flavor_performance_type = "computingv3"
node_name                    = "tf-test-node"
node_password                = "your_node_password"
data_volumes_configuration   = [
  {
    volumetype = "SSD"
    size       = 100
  }
]

node_pool_password = "your_node_pool_password"
node_pool_name     = "tf-test-node-pool"
