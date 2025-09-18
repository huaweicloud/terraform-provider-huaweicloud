# ST.001 Disable
resource "huaweicloud_obs_bucket" "source" {
  bucket        = var.source_bucket_name
  acl           = "private"
  force_destroy = true
}

resource "huaweicloud_obs_bucket" "transfer" {
  bucket        = var.transfer_bucket_name
  acl           = "private"
  force_destroy = true
}
# ST.001 Enable

resource "huaweicloud_cts_data_tracker" "test" {
  depends_on = [
    huaweicloud_obs_bucket.source,
    huaweicloud_obs_bucket.transfer,
  ]

  name          = var.tracker_name
  enabled       = var.tracker_enabled
  tags          = var.tracker_tags
  data_bucket   = huaweicloud_obs_bucket.source.bucket
  bucket_name   = huaweicloud_obs_bucket.transfer.bucket
  file_prefix   = var.trace_object_prefix
  compress_type = var.trace_file_compression_type
  lts_enabled   = var.is_lts_enabled
}
