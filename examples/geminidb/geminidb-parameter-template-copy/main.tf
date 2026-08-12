# Create a source GeminiDB parameter template with custom parameter values
resource "huaweicloud_geminidb_parameter_template" "test" {
  name        = var.template_name
  description = var.template_description

  datastore {
    type    = var.datastore_type
    version = var.datastore_version
    mode    = var.datastore_mode
  }

  values = var.parameter_values

  lifecycle {
    ignore_changes = [
      datastore,
    ]
  }
}

# Copy the GeminiDB parameter template with updated name, description and values
resource "huaweicloud_geminidb_parameter_template_copy" "test" {
  config_id   = huaweicloud_geminidb_parameter_template.test.id
  name        = var.copy_name
  description = var.copy_description

  values = var.copy_values

  lifecycle {
    ignore_changes = [
      config_id,
      values,
    ]
  }
}
