# Create a CC global connection bandwidth
resource "huaweicloud_cc_global_connection_bandwidth" "test" {
  name        = var.global_connection_bandwidth_name
  type        = var.bandwidth_type
  bordercross = var.bordercross
  size        = var.bandwidth_size
  charge_mode = var.charge_mode
  description = var.global_connection_bandwidth_description

  tags = var.global_connection_bandwidth_tags
}
