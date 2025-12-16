resource "huaweicloud_hss_host_protection" "test" {
  host_id                = var.host_id
  version                = var.protection_version
  charging_mode          = "postPaid"
  is_wait_host_available = var.is_wait_host_available
  enterprise_project_id  = var.enterprise_project_id
}
