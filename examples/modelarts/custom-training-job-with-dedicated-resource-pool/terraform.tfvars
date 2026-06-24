vpc_name                = "tf_test_vpc"
subnet_name             = "tf_test_subnet"
security_group_name     = "tf_test_security_group"
turbo_name              = "tf_test_sfs_turbo"
network_name            = "tf-test-network"
resource_pool_name      = "tf-test-resource-pool"
training_job_name       = "tf_test_training_job"
training_job_code_dir   = "your_training_code_dir"
training_job_command    = "your_training_command"
resource_pool_flavor_id = "your_resource_pool_flavor_id"

training_job_engine = {
  image_url = "your_swr_image_url"
}
