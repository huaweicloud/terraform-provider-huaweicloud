resource "huaweicloud_dcs_custom_template" "test" {
  template_id = var.source_template_id
  source_type = var.source_type
  name        = var.template_name
  description = var.template_description

  dynamic "params" {
    for_each = var.template_params

    content {
      param_name  = params.key
      param_value = params.value
    }
  }
}
