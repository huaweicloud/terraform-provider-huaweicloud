# Buckets configuration
# Note: access_key and secret_key are sensitive and should be kept secure
buckets = [
  {
    name          = "your-source-bucket-name"
    role          = "source"
    region        = "cn-north-4"
    access_key    = "your_source_access_key"
    secret_key    = "your_source_secret_key"
    storage_class = "STANDARD"
    acl           = "private"
    force_destroy = true

    tags = {
      Environment = "test"
      Project     = "sync-demo"
    }
  },
  {
    name          = "your-destination-bucket-name"
    role          = "destination"
    region        = "cn-north-4"
    access_key    = "your_destination_access_key"
    secret_key    = "your_destination_secret_key"
    storage_class = "STANDARD"
    acl           = "private"
    force_destroy = true

    tags = {
      Environment = "test"
      Project     = "sync-demo"
    }
  }
]

# Source objects configuration
source_object_configurations = [
  {
    key          = "test-file-1.txt"
    content      = "This is test file 1 content"
    content_type = "text/plain"
  },
  {
    key          = "test-file-2.txt"
    content      = "This is test file 2 content"
    content_type = "text/plain"
  }
]

# Sync task configuration
sync_task_description               = "Incremental sync task between two OBS buckets"
sync_task_enable_kms                = false
sync_task_enable_restore            = false
sync_task_enable_metadata_migration = true
sync_task_consistency_check         = "crc64"
sync_task_action                    = "start"
