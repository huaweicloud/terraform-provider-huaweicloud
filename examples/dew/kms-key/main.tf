resource "huaweicloud_kms_key" "test" {
  key_alias             = var.key_name
  key_algorithm         = var.key_algorithm
  key_usage             = var.key_usage
  origin                = var.key_source
  key_description       = var.key_description
  enterprise_project_id = var.enterprise_project_id
  tags                  = var.key_tags
  pending_days          = var.key_schedule_time
}
