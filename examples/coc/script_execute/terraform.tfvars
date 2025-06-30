
# Enterprise Project
enterprise_project_id = "0"

# VPC
vpc_name               = "tf_test_vpc"
vpc_cidr               = "192.168.0.0/16"
subnet_name            = "tf_test_subnet"
security_group_name    = "tf_test_secgroup"
ecs_instance_name      = "tf_test_coc_script_execute"
ecs_instance_user_data = "your_user_data" # Please replace it with the command you actually used to install icagent

# COC
coc_script_name        = "tf_coc_script_demo"
coc_script_description = "Created by terraform script"
coc_script_risk_level  = "LOW"
coc_script_version     = "1.0.0"
coc_script_type        = "SHELL"
coc_script_content     = <<EOF
#! /bin/bash
echo "hello world!"
EOF
coc_script_parameters = [
  {
    name        = "name"
    value       = "world"
    description = "the parameter"
  }
]
coc_script_execute_timeout      = 600
coc_script_execute_execute_user = "root"
coc_script_execute_parameters = [
  {
    name  = "name"
    value = "somebody"
  }
]
