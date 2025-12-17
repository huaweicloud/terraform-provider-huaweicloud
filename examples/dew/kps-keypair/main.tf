resource "huaweicloud_kps_keypair" "test" {
  name            = var.keypair_name
  scope           = var.keypair_scope
  user_id         = var.keypair_user_id
  encryption_type = var.keypair_encryption_type
  kms_key_id      = var.kms_key_id
  kms_key_name    = var.kms_key_name
  description     = var.keypair_description
}
