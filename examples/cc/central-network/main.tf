# Create a CC central network
resource "huaweicloud_cc_central_network" "test" {
  name                  = var.central_network_name
  description           = var.central_network_description
  enterprise_project_id = var.enterprise_project_id

  tags = var.central_network_tags
}
