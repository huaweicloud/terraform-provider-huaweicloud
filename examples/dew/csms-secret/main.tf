resource "huaweicloud_csms_secret" "test" {
  name                  = var.secret_name
  secret_text           = var.secret_text
  secret_type           = var.secret_type
  kms_key_id            = var.kms_key_id
  description           = var.secret_description
  enterprise_project_id = var.enterprise_project_id
  tags                  = var.secret_tags
}
