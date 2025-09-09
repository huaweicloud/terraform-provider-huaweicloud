resource "huaweicloud_lts_group" "test" {
  group_name  = var.group_name
  ttl_in_days = var.group_log_expiration_days
}

resource "huaweicloud_lts_stream" "test" {
  group_id    = huaweicloud_lts_group.test.id
  stream_name = var.stream_name
  ttl_in_days = var.stream_log_expiration_days
}

resource "huaweicloud_obs_bucket" "test" {
  bucket        = var.bucket_name
  acl           = "private"
  force_destroy = true
}

resource "huaweicloud_lts_transfer" "test" {
  log_group_id = huaweicloud_lts_group.test.id

  log_streams {
    log_stream_id = huaweicloud_lts_stream.test.id
  }

  log_transfer_info {
    log_transfer_type   = var.transfer_type
    log_transfer_mode   = var.transfer_mode
    log_storage_format  = var.transfer_storage_format
    log_transfer_status = var.transfer_status

    log_transfer_detail {
      obs_bucket_name     = huaweicloud_obs_bucket.test.bucket
      obs_period          = 3
      obs_period_unit     = "hour"
      obs_dir_prefix_name = var.bucket_dir_prefix_name
      obs_time_zone       = var.bucket_time_zone
      obs_time_zone_id    = var.bucket_time_zone_id
    }
  }

  depends_on = [
    huaweicloud_obs_bucket.test,
  ]
}
