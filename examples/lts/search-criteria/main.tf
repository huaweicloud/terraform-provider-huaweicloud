resource "huaweicloud_lts_group" "test" {
  group_name            = var.group_name
  ttl_in_days           = var.group_log_expiration_days
  tags                  = var.group_tags
  enterprise_project_id = var.enterprise_project_id
}

resource "huaweicloud_lts_stream" "test" {
  group_id              = huaweicloud_lts_group.test.id
  stream_name           = var.stream_name
  ttl_in_days           = var.stream_log_expiration_days
  tags                  = var.stream_tags
  enterprise_project_id = var.enterprise_project_id
  is_favorite           = var.stream_is_favorite
}

resource "huaweicloud_lts_search_criteria" "test" {
  log_group_id  = huaweicloud_lts_group.test.id
  log_stream_id = huaweicloud_lts_stream.test.id

  criteria = var.search_criteria
  name     = var.search_criteria_name
  type     = var.search_criteria_type
}
