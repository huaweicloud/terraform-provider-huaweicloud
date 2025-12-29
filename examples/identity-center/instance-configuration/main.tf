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

resource "huaweicloud_identitycenter_sso_configuration" "test" {
  instance_id                 = var.is_instance_create ? huaweicloud_identitycenter_instance.test[0].identity_store_id : data.huaweicloud_identitycenter_instance.test[0].identity_store_id
  configuration_type          = var.configuration_type
  mfa_mode                    = var.configuration_mfa_mode
  allowed_mfa_types           = var.configuration_allowed_mfa_types
  no_mfa_signin_behavior      = var.configuration_no_mfa_signin_behavior
  no_password_signin_behavior = var.configuration_no_password_signin_behavior
  max_authentication_age      = var.configuration_max_authentication_age
}
