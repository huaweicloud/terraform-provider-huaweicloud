vpc_name              = "tf_test_vpc"
subnet_name           = "tf_test_subnet"
bandwidth_name        = "tf_test_bandwidth"
bandwidth_size        = 5
cluster_name          = "tf-test-cluster"
node_performance_type = "computingv3"
keypair_name          = "tf_test_keypair"
node_pool_name        = "tf-test-node-pool"
node_pool_tags        = {
  "owner" = "terraform"
}

root_volume_size           = 40
root_volume_type           = "SSD"
data_volumes_configuration = [
  {
    volumetype     = "SSD"
    size           = 100
    count          = 2
    virtual_spaces = [
      {
        name        = "kubernetes"
        size        = "10%"
        lvm_lv_type = "linear"
      },
      {
        name = "runtime"
        size = "90%"
      }
    ]
  },
  {
    volumetype     = "SSD"
    size           = 100
    count          = 1
    virtual_spaces = [
      {
        name        = "user"
        size        = "100%"
        lvm_lv_type = "linear"
        lvm_path    = "/workspace"
      }
    ]
  }
]
