vpc_name                  = "tf-test-vpc"
subnet_name               = "tf-test-subnet"
security_group_name       = "tf-test-security-group"
turbo_name                = "tf-test-sfs-turbo"
network_name              = "tf-test-network"
resource_pool_name        = "tf-test-resource-pool"
resource_pool_description = "This is a demo"

resource_pool_resources = [
  {
    count = 1
  }
]
