resource "huaweicloud_dds_parameter_template" "test"{
  name             = var.template_name
  parameter_values = var.template_mapping
  node_type        = var.template_node_type
  node_version     = var.database_version
  description      = var.template_description
}
