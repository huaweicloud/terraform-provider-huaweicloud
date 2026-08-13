# ST.001 Disable
resource "huaweicloud_obs_bucket" "source" {
  bucket        = var.source_bucket_name
  storage_class = var.bucket_storage_class
  acl           = var.bucket_acl
  force_destroy = var.bucket_force_destroy
}

resource "huaweicloud_obs_bucket" "dest" {
  bucket        = var.dest_bucket_name
  storage_class = var.bucket_storage_class
  acl           = var.bucket_acl
  force_destroy = var.bucket_force_destroy
}
# ST.001 Enable

resource "huaweicloud_oms_migration_sync_task" "test" {
  src_cloud_type            = var.source_cloud_type
  src_region                = var.source_region
  src_ak                    = var.source_access_key
  src_sk                    = var.source_secret_key
  src_bucket                = huaweicloud_obs_bucket.source.bucket
  dst_ak                    = var.dest_access_key
  dst_sk                    = var.dest_secret_key
  dst_bucket                = huaweicloud_obs_bucket.dest.bucket
  description               = var.task_description
  consistency_check         = var.consistency_check
  enable_metadata_migration = var.enable_metadata_migration
}
