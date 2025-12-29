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

resource "huaweicloud_obs_bucket_object" "test" {
  bucket       = huaweicloud_obs_bucket.test.id
  key          = var.object_extension_name != "" ? format("%s%s", var.object_name, var.object_extension_name) : var.object_name
  content_type = "application/xml"
  content      = var.object_upload_content
}

resource "huaweicloud_obs_bucket_policy" "test" {
  bucket = huaweicloud_obs_bucket.test.id
  policy = <<POLICY
{
  "Statement": [
    {
      "Sid": "",
      "Effect": "Allow",
      "Principal": {"ID": "*"},
      "Action": ["GetObject"],
      "Resource": "${huaweicloud_obs_bucket.test.id}/*"
    }
  ]
}
  POLICY
}

resource "huaweicloud_oms_migration_task" "test" {
  depends_on = [
    huaweicloud_obs_bucket_object.test,
    huaweicloud_obs_bucket_policy.test
  ]

  start_task                     = var.task_is_start
  type                           = var.task_type
  enable_kms                     = var.task_enable_kms
  migrate_since                  = var.task_migrate_since
  object_overwrite_mode          = var.task_object_overwrite_mode
  consistency_check              = var.task_consistency_check
  enable_requester_pays          = var.task_enable_requester_pays
  enable_failed_object_recording = var.task_enable_failed_object_recording

  source_object {
    data_source = "HuaweiCloud"
    region      = huaweicloud_obs_bucket.test.region
    bucket      = huaweicloud_obs_bucket.test.id
    access_key  = var.access_key
    secret_key  = var.secret_key
    object      = [var.object_name]
  }

  destination_object {
    region     = var.target_bucket_configuration.region != "" ? var.target_bucket_configuration.region : huaweicloud_obs_bucket.test.region
    bucket     = var.target_bucket_configuration.bucket
    access_key = var.target_bucket_configuration.access_key != "" ? var.target_bucket_configuration.access_key : var.access_key
    secret_key = var.target_bucket_configuration.secret_key != "" ? var.target_bucket_configuration.secret_key : var.secret_key
  }

  dynamic "bandwidth_policy" {
    for_each = var.bandwidth_policy_configurations

    content {
      max_bandwidth = bandwidth_policy.value["max_bandwidth"]
      start         = bandwidth_policy.value["start"]
      end           = bandwidth_policy.value["end"]
    }
  }
}
