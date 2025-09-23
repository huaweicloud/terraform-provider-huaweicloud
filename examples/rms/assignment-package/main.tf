data "huaweicloud_rms_assignment_package_templates" "test" {
  template_key = var.template_key
}

resource "huaweicloud_rms_assignment_package" "test" {
  name         = var.assignment_package_name
  template_key = data.huaweicloud_rms_assignment_package_templates.test.templates.0.template_key

  dynamic "vars_structure" {
    for_each = data.huaweicloud_rms_assignment_package_templates.test.templates.0.parameters

    content {
      var_key   = vars_structure.value["name"]
      var_value = vars_structure.value["default_value"]
    }
  }
}
