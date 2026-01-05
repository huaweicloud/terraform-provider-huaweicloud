# Incremental sync between two OBS buckets

This example provides best practice code for using Terraform to create an Object Migration Service (OMS) incremental
sync task in HuaweiCloud for synchronizing objects from a source OBS bucket to a destination OBS bucket.

## Prerequisites

* A HuaweiCloud account
* Terraform installed
* HuaweiCloud access key and secret key (AK/SK) for both source and destination buckets
* Source and destination OBS buckets (will be created if they don't exist)

## Variables Introduction

The following variables need to be configured:

### Authentication Variables

* `region_name` - The region where the OMS sync task is located
* `access_key` - The access key for the HuaweiCloud provider
* `secret_key` - The secret key for the HuaweiCloud provider

### Resource Variables

#### Required Variables

* `buckets` - The configurations of OBS buckets for sync task (must include exactly one source and one
  destination bucket)
  Each bucket configuration supports the following attributes:
  - `name` - The name of the OBS bucket (required)
  - `role` - The role of the bucket, must be either "source" or "destination" (required)
  - `region` - The region where the OBS bucket is located (required)
  - `access_key` - The access key for accessing the OBS bucket (required)
  - `secret_key` - The secret key for accessing the OBS bucket (required)
  - `storage_class` - The storage class of the OBS bucket (default: "STANDARD")
  - `acl` - The ACL of the OBS bucket (default: "private")
  - `tags` - The tags of the OBS bucket (default: {})

#### Optional Variables

* `source_object_configurations` - The configurations of objects to be uploaded to the source bucket
  (default: [])
  Each object configuration supports the following attributes:
  - `key` - The object key
  - `content` - The object content
  - `content_type` - The content type of the object (default: "text/plain")
* `sync_task_description` - The description of the OMS migration sync task (default: "")
* `sync_task_enable_kms` - Whether to enable KMS for the OMS migration sync task (default: false)
* `sync_task_enable_restore` - Whether to enable restore for the OMS migration sync task (default: false)
* `sync_task_enable_metadata_migration` - Whether to enable metadata migration for the OMS migration sync
  task (default: true)
* `sync_task_consistency_check` - The consistency check method for the OMS migration sync task
  (default: "crc64")
* `sync_task_action` - The action of the OMS migration sync task (default: "start")
* `source_cdn_configuration` - The CDN configuration for the source bucket (default: null)
  Supports the following attributes:
  - `domain` - The CDN domain name
  - `protocol` - The protocol used
  - `authentication_type` - The authentication type (default: "NONE")
  - `authentication_key` - The authentication key

## Usage

* Copy this example script to your `main.tf`.

* Create a `terraform.tfvars` file and fill in the required variables:

  ```hcl
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
  ```

* Initialize Terraform:

  ```bash
  $ terraform init
  ```

* Review the Terraform plan:

  ```bash
  $ terraform plan
  ```

* Apply the configuration:

  ```bash
  $ terraform apply
  ```

* To clean up the resources:

  ```bash
  $ terraform destroy
  ```

## Note

* Make sure to keep your credentials secure and never commit them to version control
* The creation of the OBS buckets and objects is usually instantaneous
* Sync task creation may take some time depending on the configuration
* This example creates both source and destination OBS buckets, uploads test objects to the source bucket,
  and creates an incremental sync task
* Incremental sync tasks will only synchronize objects that have been added or modified since the last sync
* The sync task will continuously monitor and synchronize changes from source to destination
* Bucket names must be globally unique across all HuaweiCloud accounts
* When `bucket_force_destroy` is set to true, the buckets can be destroyed even if they contain objects
* The sync task supports metadata migration to preserve object metadata during synchronization
* Consistency check using CRC64 ensures data integrity during synchronization
* You can control the sync task by updating the `action` parameter (start/stop)
* Source and destination buckets can be in different regions
* The sync task will automatically handle object overwrite conflicts based on the configured mode

## Requirements

| Name | Version |
| ---- | ---- |
| terraform | >= 1.3.0 |
| huaweicloud | >= 1.61.0 |
