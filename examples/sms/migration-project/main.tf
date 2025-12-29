resource "huaweicloud_sms_migration_project" "test" {
  name          = var.migration_project_name
  region        = var.migration_project_region
  use_public_ip = var.migration_project_use_public_ip
  exist_server  = var.migration_project_exist_server
  type          = var.migration_project_type
  syncing       = var.migration_project_syncing
}
