# Create a DC connect gateway
resource "huaweicloud_dc_connect_gateway" "test" {
  name           = var.connect_gateway_name
  description    = var.connect_gateway_description
  address_family = var.address_family
}
