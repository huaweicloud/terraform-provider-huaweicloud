resource "huaweicloud_coc_script" "test" {
  name        = var.script_name
  description = var.script_description
  risk_level  = var.script_risk_level
  version     = var.script_version
  type        = var.script_type
  content     = var.script_content

  dynamic "parameters" {
    for_each = var.script_parameters

    content {
      name        = parameters.value.name
      value       = parameters.value.value
      description = parameters.value.description
      sensitive   = parameters.value.sensitive
    }
  }
}
