resource "huaweicloud_servicestagev3_configuration_group" "test" {
  name        = var.configuration_group_name
  description = var.configuration_group_description
}

resource "huaweicloud_servicestagev3_configuration" "test" {
  config_group_id = huaweicloud_servicestagev3_configuration_group.test.id
  name            = var.configuration_name
  type            = "properties"
  content         = var.configuration_content
  description     = var.configuration_description

  lifecycle {
    ignore_changes = [
      type
    ]
  }
}
