resource "huaweicloud_coc_script" "test" {
  name        = var.script_name
  description = "coc script description"
  risk_level  = "LOW"
  version     = "1.0.0"
  type        = "SHELL"

  content = <<EOF
#! /bin/bash
echo "hello world!"
EOF

  parameters {
    name        = "name"
    value       = "world"
    description = "the first parameter"
  }
  parameters {
    name        = "company"
    value       = "Huawei"
    description = "the second parameter"
    sensitive   = true
  }
}
