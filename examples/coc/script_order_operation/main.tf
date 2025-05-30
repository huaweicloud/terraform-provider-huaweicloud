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
    description = "the parameter"
  }
}

resource "huaweicloud_coc_script_execute" "test" {
  script_id    = huaweicloud_coc_script.test.id
  instance_id  = var.script_execute_name
  timeout      = 600
  execute_user = "root"

  parameters {
    name  = "name"
    value = "somebody"
  }
  parameters {
    name  = "company"
    value = "HuaweiCloud"
  }
}

resource "huaweicloud_coc_script_order_operation" "test" {
  execute_uuid   = huaweicloud_coc_script_execute.test.id
  batch_index    = 1
  instance_id    = 1
  operation_type = var.operation_type
}
