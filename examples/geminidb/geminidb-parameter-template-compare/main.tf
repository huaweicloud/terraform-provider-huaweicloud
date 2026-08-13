# Create the GeminiDB parameter template
resource "huaweicloud_geminidb_parameter_template" "test" {
  for_each = var.parameter_template_configuration

  name        = each.value.template_name
  description = each.value.template_description

  datastore {
    type    = var.datastore_type
    version = var.datastore_version
    mode    = var.datastore_mode
  }

  values = each.value.parameter_values

  lifecycle {
    ignore_changes = [
      datastore,
    ]
  }
}

# Compare the differences between the source and target parameter templates
resource "huaweicloud_geminidb_parameter_template_compare" "test" {
  source_configuration_id = huaweicloud_geminidb_parameter_template.test["source_parameter_template"].id
  target_configuration_id = huaweicloud_geminidb_parameter_template.test["target_parameter_template"].id
}
