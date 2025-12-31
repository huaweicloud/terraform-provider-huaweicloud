resource "huaweicloud_rgc_account" "test" {
  name                            = var.account_name
  email                           = var.account_email
  phone                           = var.account_phone
  identity_store_user_name        = var.identity_store_user_name
  identity_store_email            = var.identity_store_email
  parent_organizational_unit_name = var.parent_organizational_unit_name
  parent_organizational_unit_id   = var.parent_organizational_unit_id
}
