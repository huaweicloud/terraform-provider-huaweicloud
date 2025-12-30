# Create a resource share instance
resource "huaweicloud_ram_resource_share" "test" {
  name                      = var.resource_share_name
  description               = var.description
  principals                = var.principals
  resource_urns             = var.resource_urns
  permission_ids            = var.permission_ids
  allow_external_principals = var.allow_external_principals
}
