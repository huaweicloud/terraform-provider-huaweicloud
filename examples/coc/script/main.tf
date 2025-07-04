resource "huaweicloud_coc_script" "test" {
  name        = var.coc_script_name
  description = var.coc_script_description
  risk_level  = var.coc_script_risk_level
  version     = var.coc_script_version
  type        = var.coc_script_type

  content = var.coc_script_content

  dynamic "parameters" {
    for_each = var.coc_script_parameters

    content {
      name        = parameters.value.name
      value       = parameters.value.value
      description = parameters.value.description
      sensitive   = parameters.value.sensitive
    }
  }
}
