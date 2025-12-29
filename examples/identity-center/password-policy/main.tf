data "huaweicloud_identitycenter_instance" "test" {
  count = var.is_instance_create ? 0 : 1
}

resource "huaweicloud_identitycenter_registered_region" "test" {
  count = var.is_instance_create ? var.is_region_need_register ? 1 : 0 : 0

  region_id = var.region_name
}

resource "huaweicloud_identitycenter_instance" "test" {
  count = var.is_instance_create ? 1 : 0

  depends_on = [huaweicloud_identitycenter_registered_region.test]

  alias = var.instance_store_id_alias != "" ? var.instance_store_id_alias : null
}

resource "huaweicloud_identitycenter_password_policy" "test" {
  identity_store_id            = var.is_instance_create ? huaweicloud_identitycenter_instance.test[0].identity_store_id : data.huaweicloud_identitycenter_instance.test[0].identity_store_id
  max_password_age             = var.policy_max_password_age
  minimum_password_length      = var.policy_minimum_password_length
  password_reuse_prevention    = var.policy_password_reuse_prevention ? 1 : null # 1 means the password reuse prevention is enabled
  require_uppercase_characters = var.policy_require_uppercase_characters
  require_lowercase_characters = var.policy_require_lowercase_characters
  require_numbers              = var.policy_require_numbers
  require_symbols              = var.policy_require_symbols
}
