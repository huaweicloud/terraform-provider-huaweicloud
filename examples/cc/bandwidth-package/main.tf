# Create a CC bandwidth package
resource "huaweicloud_cc_bandwidth_package" "test" {
  name                  = var.bandwidth_package_name
  local_area_id         = var.local_area_id
  remote_area_id        = var.remote_area_id
  charge_mode           = var.charge_mode
  billing_mode          = var.billing_mode
  bandwidth             = var.bandwidth
  description           = var.bandwidth_package_description
  enterprise_project_id = var.enterprise_project_id

  tags = var.bandwidth_package_tags
}
