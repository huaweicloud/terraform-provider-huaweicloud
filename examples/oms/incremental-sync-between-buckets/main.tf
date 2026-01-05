resource "huaweicloud_obs_bucket" "test" {
  count = length(var.buckets)

  bucket        = var.buckets[count.index].name
  storage_class = lookup(var.buckets[count.index], "storage_class", "STANDARD")
  acl           = lookup(var.buckets[count.index], "acl", "private")
  force_destroy = lookup(var.buckets[count.index], "force_destroy", true)
  tags          = lookup(var.buckets[count.index], "tags", {})
}

locals {
  source_bucket            = [for bucket in var.buckets : bucket if bucket.role == "source"][0]
  destination_bucket       = [for bucket in var.buckets : bucket if bucket.role == "destination"][0]
  source_bucket_index      = index([for bucket in var.buckets : bucket.name], local.source_bucket.name)
  destination_bucket_index = index([for bucket in var.buckets : bucket.name], local.destination_bucket.name)
}

resource "huaweicloud_obs_bucket_object" "test" {
  count = length(var.source_object_configurations)

  bucket       = huaweicloud_obs_bucket.test[local.source_bucket_index].id
  key          = var.source_object_configurations[count.index].key
  content      = var.source_object_configurations[count.index].content
  content_type = lookup(var.source_object_configurations[count.index], "content_type", "text/plain")
}

resource "huaweicloud_oms_migration_sync_task" "test" {
  src_region                = local.source_bucket.region
  src_bucket                = huaweicloud_obs_bucket.test[local.source_bucket_index].id
  src_ak                    = local.source_bucket.access_key
  src_sk                    = local.source_bucket.secret_key
  dst_bucket                = huaweicloud_obs_bucket.test[local.destination_bucket_index].id
  dst_ak                    = local.destination_bucket.access_key
  dst_sk                    = local.destination_bucket.secret_key
  description               = var.sync_task_description
  enable_kms                = var.sync_task_enable_kms
  enable_restore            = var.sync_task_enable_restore
  enable_metadata_migration = var.sync_task_enable_metadata_migration
  consistency_check         = var.sync_task_consistency_check
  action                    = var.sync_task_action

  dynamic "source_cdn" {
    for_each = var.source_cdn_configuration != null ? [var.source_cdn_configuration] : []

    content {
      domain              = source_cdn.value.domain
      protocol            = source_cdn.value.protocol
      authentication_type = lookup(source_cdn.value, "authentication_type", "NONE")
      authentication_key  = lookup(source_cdn.value, "authentication_key", null)
    }
  }

  depends_on = [huaweicloud_obs_bucket_object.test]
}
