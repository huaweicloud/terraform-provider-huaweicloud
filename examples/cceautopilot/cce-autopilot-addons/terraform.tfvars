# Required variables
vpc_name              = "tf_test_vpc"
subnet_name           = "tf_test_subnet"
cluster_name          = "tf-test-autopilot-cluster"
swr_organization_name = "tf-test-swr-org"
# Optional variables with custom values
cluster_description   = "CCE Autopilot cluster for addon deployment"
cluster_version       = "v1.36"
addon_template_name   = "log-agent"
