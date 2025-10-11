vpc_name                   = "tf_test_vpc"
subnet_name                = "tf_test_subnet"
bandwidth_name             = "tf_test_bandwidth"
bandwidth_size             = 5
cluster_name               = "tf-test-cluster"
node_performance_type      = "computingv3"
keypair_name               = "tf_test_keypair"
node_name                  = "tf-test-node"
root_volume_size           = 40
root_volume_type           = "SSD"
data_volumes_configuration = [
  {
    volumetype = "SSD"
    size       = 100
  }
]
