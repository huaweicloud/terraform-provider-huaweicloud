vpc_name            = "tf_test_coc_script_execute_vpc"
subnet_name         = "tf_test_coc_script_execute_subnet"
security_group_name = "tf_test_coc_script_execute_secgroup"
instance_name       = "tf_test_coc_script_execute"
instance_user_data  = "your_user_data" # Please replace it with the command you actually used to install icagent
script_name         = "tf_coc_script_execute"
script_description  = "Created by terraform script"
script_risk_level   = "LOW"
script_version      = "1.0.0"
script_type         = "SHELL"
script_content      = <<EOF
#! /bin/bash
echo "hello world!"
EOF
script_parameters = [
  {
    name        = "name"
    value       = "world"
    description = "the parameter"
  }
]
script_execute_timeout      = 600
script_execute_execute_user = "root"
script_execute_parameters = [
  {
    name  = "name"
    value = "somebody"
  }
]
