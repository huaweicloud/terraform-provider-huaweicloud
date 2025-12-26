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

secret_name = "tf-test-secret"
secret_data = {
  "access.key" = "your_access_key"
  "secret.key" = "your_secret_key"
}

pvc_name              = "tf-test-pvc-obs"
deployment_name       = "tf-test-deployment"
deployment_containers = [
  {
    name  = "container-1"
    image = "nginx:latest"

    volume_mounts = [
      {
        mount_path = "/data"
      }
    ]
  }
]
