resource "huaweicloud_identityv5_password_policy" "test" {
  maximum_consecutive_identical_chars = var.policy_max_consecutive_identical_chars
  minimum_password_age                = var.policy_min_password_age
  minimum_password_length             = var.policy_min_password_length
  password_reuse_prevention           = var.policy_password_reuse_prevention
  password_not_username_or_invert     = var.policy_password_not_username_or_invert
  password_validity_period            = var.policy_password_validity_period
  password_char_combination           = var.policy_password_char_combination
  allow_user_to_change_password       = var.policy_allow_user_to_change_password
}
