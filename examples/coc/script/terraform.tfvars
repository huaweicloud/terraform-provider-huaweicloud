script_name        = "tf_coc_script"
script_description = "Created by terraform script"
script_risk_level  = "LOW"
script_version     = "1.0.0"
script_type        = "SHELL"
script_content     = <<EOF
#! /bin/bash
echo "hello world!"
EOF
script_parameters = [
  {
    name        = "name"
    value       = "world"
    description = "the first parameter"
  },
  {
    name        = "company"
    value       = "Huawei"
    description = "the second parameter"
    sensitive   = true
  }
]
