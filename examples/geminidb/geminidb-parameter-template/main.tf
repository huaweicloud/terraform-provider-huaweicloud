# Create a GeminiDB parameter template with custom parameter values
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
