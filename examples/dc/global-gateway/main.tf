# Create a global DC gateway
resource "huaweicloud_dc_global_gateway" "test" {
  name                  = var.global_gateway_name
  description           = var.global_gateway_description
  address_family        = var.address_family
  bgp_asn               = var.bgp_asn
  enterprise_project_id = var.enterprise_project_id

  tags = var.global_gateway_tags
}
