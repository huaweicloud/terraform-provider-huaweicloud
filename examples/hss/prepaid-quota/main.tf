resource "huaweicloud_hss_quota" "test" {
  version               = var.quota_version
  period_unit           = var.period_unit
  period                = var.period
  auto_renew            = var.is_auto_renew
  enterprise_project_id = var.enterprise_project_id
  tags                  = var.quota_tags
}
