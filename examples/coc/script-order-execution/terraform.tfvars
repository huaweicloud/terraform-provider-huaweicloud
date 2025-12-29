vpc_name            = "tf_test_coc_order_vpc"
subnet_name         = "tf_test_coc_order_subnet"
security_group_name = "tf_test_coc_order_secgroup"
instance_name       = "tf_test_coc_order_instance"
instance_user_data  = "your_user_data" # Please replace it with the command you actually used to install icagent
script_name         = "tf_script_order_demo"
script_description  = "Created by terraform script"
script_risk_level   = "LOW"
script_version      = "1.0.0"
script_type         = "SHELL"
script_content      = <<EOF
#! /bin/bash
echo "hello $${name}!"
sleep 2m
EOF
script_parameters   = [
  {
    name        = "name"
    value       = "world"
    description = "the parameter"
  },
  {
    name        = "company"
    value       = "Terraform"
    description = "the second parameter"
    sensitive   = true
  }
]

script_execute_timeout    = 600
script_execute_user       = "root"
script_execute_parameters = [
  {
    name  = "name"
    value = "somebody"
  }
]

script_order_operation_batch_index = 1
script_order_operation_instance_id = 1
script_order_operation_type        = "CANCEL_ORDER"
