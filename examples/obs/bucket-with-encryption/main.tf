resource "huaweicloud_kms_key" "test" {
  count = var.bucket_encryption && var.bucket_encryption_key_id == "" ? 1 : 0

  key_alias = var.key_alias
  key_usage = var.key_usage
}

resource "huaweicloud_obs_bucket" "test" {
  bucket        = var.bucket_name
  storage_class = var.bucket_storage_class
  acl           = var.bucket_acl
  encryption    = var.bucket_encryption
  sse_algorithm = var.bucket_encryption ? var.bucket_sse_algorithm : null
  kms_key_id    = var.bucket_encryption ? var.bucket_encryption_key_id != "" ? var.bucket_encryption_key_id : huaweicloud_kms_key.test[0].id : null
  force_destroy = var.bucket_force_destroy
  tags          = var.bucket_tags

  lifecycle {
    ignore_changes = [
      sse_algorithm
    ]
  }
}
