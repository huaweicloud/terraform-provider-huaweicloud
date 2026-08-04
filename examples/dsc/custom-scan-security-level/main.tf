resource "huaweicloud_dsc_scan_security_level" "test" {
  security_level_name = var.security_level_name
  color_number        = var.security_level_color_number
  security_level_desc = var.security_level_description
}
