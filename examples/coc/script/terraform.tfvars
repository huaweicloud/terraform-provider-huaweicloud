coc_script_name        = "tf_coc_script"
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
    description = "the first parameter"
  },
  {
    name        = "company"
    value       = "Huawei"
    description = "the second parameter"
    sensitive   = true
  }
]
