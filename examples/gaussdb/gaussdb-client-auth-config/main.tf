resource "huaweicloud_gaussdb_client_auth_config" "test" {
  instance_id = var.instance_id
  type        = var.config_type
  database    = var.config_database
  user        = var.config_user
  address     = var.config_address
  method      = var.config_method
}
