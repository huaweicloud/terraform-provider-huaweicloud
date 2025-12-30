# Create an organizational unit
resource "huaweicloud_rgc_organizational_unit" "test" {
  organizational_unit_name      = var.organizational_unit_name
  parent_organizational_unit_id = var.parent_organizational_unit_id
}

# Account enrollment with blueprint configuration
resource "huaweicloud_rgc_account_enroll" "test" {
  managed_account_id            = var.blueprint_managed_account_id
  parent_organizational_unit_id = var.create_organizational_unit ? huaweicloud_rgc_organizational_unit.test.organizational_unit_id : var.parent_organizational_unit_id

  blueprint {
    blueprint_product_id                    = var.blueprint_product_id
    blueprint_product_version               = var.blueprint_product_version
    variables                               = var.blueprint_variables
    is_blueprint_has_multi_account_resource = var.is_blueprint_has_multi_account_resource
  }
}
