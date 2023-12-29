
# create a obs source
resource "huaweicloud_obs_bucket" "source" {
  region        = var.source_region
  bucket        = var.source_obs_name
  acl           = "private"
  force_destroy = true
}

# create a obs object as source for migration
resource "huaweicloud_obs_bucket_object" "source" {
  region  = var.source_region
  bucket  = huaweicloud_obs_bucket.source.bucket
  key     = "test.txt"
  content = "test content"
}

# create migration dest
resource "huaweicloud_obs_bucket" "dest" {
  region        = var.dest_region
  bucket        = var.source_dest_name
  acl           = "private"
  force_destroy = true
}

# create migration task group
resource "huaweicloud_oms_migration_task_group" "test" {
  source_object {
    data_source = "HuaweiCloud"
    region      = var.source_region
    bucket      = huaweicloud_obs_bucket.source.bucket
    access_key  = var.source_ak
    secret_key  = var.source_sk
    object      = [""]
  }

  destination_object {
    region     = var.dest_region
    bucket     = huaweicloud_obs_bucket.dest.bucket
    access_key = var.dest_ak
    secret_key = var.dest_sk
  }

  action                         = "stop"
  type                           = "PREFIX"
  enable_kms                     = true
  description                    = "test task group"
  migrate_since                  = "2023-01-02 15:04:05"
  object_overwrite_mode          = "CRC64_COMPARISON_OVERWRITE"
  consistency_check              = "crc64"
  enable_requester_pays          = true
  enable_failed_object_recording = true

  bandwidth_policy {
    max_bandwidth = 1
    start         = "15:00"
    end           = "16:00"
  }

  bandwidth_policy {
    max_bandwidth = 2
    start         = "16:00"
    end           = "17:00"
  }
}