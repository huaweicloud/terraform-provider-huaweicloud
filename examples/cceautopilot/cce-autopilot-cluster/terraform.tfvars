# Required variables
vpc_name            = "tf_test_vpc"
subnet_name         = "tf_test_subnet"
cluster_name        = "tf-test-autopilot-cluster"
# Optional variables (using defaults)
cluster_description = "Created by terraform script"
cluster_tags        = {
  owner = "terraform"
}
