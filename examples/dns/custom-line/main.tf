# Create a DNS custom line
resource "huaweicloud_dns_custom_line" "test" {
  name        = var.dns_custom_line_name
  ip_segments = var.dns_custom_line_ip_segments
}
