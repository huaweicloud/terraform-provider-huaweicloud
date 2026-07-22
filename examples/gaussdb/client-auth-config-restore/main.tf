# Restore the client access authentication configuration of a GaussDB instance

resource "huaweicloud_gaussdb_client_auth_config_restore" "test" {
  instance_id    = var.instance_id
  hba_history_id = var.hba_history_id
}
