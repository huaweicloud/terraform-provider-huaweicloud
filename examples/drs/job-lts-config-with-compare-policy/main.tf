resource "huaweicloud_lts_group" "test" {
  group_name  = var.lts_group_name
  ttl_in_days = 30
}

resource "huaweicloud_lts_stream" "test" {
  group_id    = huaweicloud_lts_group.test.id
  stream_name = var.lts_stream_name
}

resource "huaweicloud_drs_lts_config" "test" {
  job_id        = var.drs_job_id
  log_group_id  = huaweicloud_lts_group.test.id
  log_stream_id = huaweicloud_lts_stream.test.id
}

resource "huaweicloud_drs_compare_policy" "test" {
  job_id         = var.drs_job_id
  period         = var.compare_policy_period
  begin_time     = var.compare_policy_begin_time
  end_time       = var.compare_policy_end_time
  compare_type   = var.compare_policy_compare_type
  compare_policy = var.compare_policy_compare_policy
  interval_hour  = var.compare_policy_interval_hour

  depends_on = [huaweicloud_drs_lts_config.test]
}
