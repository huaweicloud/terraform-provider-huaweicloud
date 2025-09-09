resource "huaweicloud_obs_bucket" "test" {
  bucket        = var.bucket_name
  acl           = "private"
  force_destroy = true
}

resource "huaweicloud_cts_tracker" "test" {
  enabled        = var.tracker_enabled
  tags           = var.tracker_tags
  delete_tracker = var.is_system_tracker_delete
  bucket_name    = huaweicloud_obs_bucket.test.bucket
  file_prefix    = var.trace_object_prefix
  compress_type  = var.trace_file_compression_type
  lts_enabled    = var.is_lts_enabled
}
