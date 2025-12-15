resource "huaweicloud_lts_group" "test" {
  group_name            = var.lts_group_name
  ttl_in_days           = var.lts_ttl_in_days
  enterprise_project_id = var.enterprise_project_id
}

resource "huaweicloud_lts_stream" "test" {
  group_id              = huaweicloud_lts_group.test.id
  stream_name           = var.lts_stream_name
  is_favorite           = var.lts_is_favorite
  enterprise_project_id = var.enterprise_project_id
}

resource "huaweicloud_antiddos_lts_config" "test" {
  lts_group_id          = huaweicloud_lts_group.test.id
  lts_attack_stream_id  = huaweicloud_lts_stream.test.id
  enterprise_project_id = var.enterprise_project_id
}
