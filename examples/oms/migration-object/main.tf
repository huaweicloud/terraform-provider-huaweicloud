
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

# create migration task
resource "huaweicloud_oms_migration_task" "test" {
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

  start_task  = false
  type        = "object"
  description = "test task"

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