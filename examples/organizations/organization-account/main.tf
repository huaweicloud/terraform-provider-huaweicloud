resource "huaweicloud_organizations_account" "test" {
  name        = var.name
  email       = var.email
  phone       = var.phone
  agency_name = var.agency_name
  parent_id   = var.parent_id
  description = var.description
  tags        = var.tags
}
