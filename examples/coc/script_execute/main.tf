resource "huaweicloud_coc_script" "test" {
  name        = var.script_name
  description = var.script_description
  risk_level  = "LOW"
  version     = "1.0.0"
  type        = "SHELL"

  content = <<EOF
#! /bin/bash
echo "hello world!"
EOF


  dynamic "parameters" {
    for_each = var.script_parameters
    content {
      name        = parameters.value.name
      value       = parameters.value.value
      description = parameters.value.description
      sensitive   = parameters.value.sensitive != null ? parameters.value.sensitive : null
    }
  }
}

resource "huaweicloud_coc_script_execute" "test" {
  script_id    = huaweicloud_coc_script.test.id
  instance_id  = var.ecs_instance_id
  timeout      = 600
  execute_user = "root"

  dynamic "parameters" {
    for_each = var.script_execute_parameters
    content {
      name  = parameters.value.name
      value = parameters.value.value
    }
  }
}
