# Associate LTS log with a GaussDB instance
resource "huaweicloud_lts_group" "test" {
  group_name  = var.lts_group_name
  ttl_in_days = 30
}

resource "huaweicloud_lts_stream" "test" {
  group_id    = huaweicloud_lts_group.test.id
  stream_name = var.lts_stream_name
}

resource "huaweicloud_gaussdb_instance_lts_log_associate" "test" {
  instance_id   = var.gaussdb_instance_id
  log_type      = var.lts_log_type
  lts_group_id  = huaweicloud_lts_group.test.id
  lts_stream_id = huaweicloud_lts_stream.test.id
}
