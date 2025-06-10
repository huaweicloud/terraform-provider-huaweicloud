# The values in this file will automatically be used by Terraform if you run `terraform apply` and file name is
# 'terraform.tfvars'.

# Enterprise Project
enterprise_project_id = "0"

# VPC
vpc_name            = "tf_test_vpc"
vpc_cidr            = "192.168.0.0/16"
subnet_name         = "tf_test_subnet"
security_group_name = "tf_test_secgroup"

# WAF
waf_dedicated_instance_name               = "tf_test_waf_dedicated_instance"
waf_dedicated_instance_specification_code = "waf.instance.professional"
waf_policy_name                           = "tf_test_waf_policy"
