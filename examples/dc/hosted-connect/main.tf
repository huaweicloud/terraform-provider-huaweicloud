# Create a hosted connect
resource "huaweicloud_dc_hosted_connect" "test" {
  name               = var.hosted_connect_name
  description        = var.hosted_connect_description
  bandwidth          = var.bandwidth
  hosting_id         = var.hosting_id
  vlan               = var.vlan
  resource_tenant_id = var.resource_tenant_id
  peer_location      = var.peer_location
}
