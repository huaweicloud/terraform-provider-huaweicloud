# Create a migration task group to migrate objects

This example provides best practice code for using Terraform to create an Object Migration Service (OMS) migration task
group in HuaweiCloud for migrating objects from source OBS bucket to destination OBS bucket.

## Prerequisites

* A HuaweiCloud account
* Terraform installed
* HuaweiCloud access key and secret key (AK/SK)

## Variable Introduction

The following variables need to be configured:

### Authentication Variables

* `region_name` - The region where the OBS bucket is located
* `access_key` - The access key of the IAM user
* `secret_key` - The secret key of the IAM user

### Resource Variables

#### Required Variables

* `bucket_name` - The name of the OBS bucket
* `object_name` - The name of the OBS object to be uploaded
* `object_upload_content` - The content of the OBS object to be uploaded
* `target_bucket_configuration` - The target bucket configuration to be migrated

#### Optional Variables

* `bucket_encryption` - Whether to enable encryption for the OBS bucket (default: true)
* `bucket_encryption_key_id` - The encryption key ID of the OBS bucket (default: "")
* `key_alias` - The alias of the KMS key (default: "")
  The alias of the KMS key (required when `bucket_encryption` is true and `bucket_encryption_key_id` is empty)
* `key_usage` - The usage of the KMS key (default: "ENCRYPT_DECRYPT")
* `bucket_storage_class` - The storage class of the OBS bucket (default: "STANDARD")
* `bucket_acl` - The ACL of the OBS bucket (default: "private")
* `bucket_sse_algorithm` - The SSE algorithm of the OBS bucket (default: "kms")
* `bucket_force_destroy` - Whether to force destroy the OBS bucket (default: true)
* `bucket_tags` - The tags of the OBS bucket (default: {})
* `object_extension_name` - The extension name of the OBS object to be uploaded (default: ".txt")
* `group_action_type` - The action type of the migration task group (default: "stop")
* `group_type` - The type of the migration task group (default: "PREFIX")
* `group_enable_kms` - Whether to enable KMS for the migration task group (default: true)
* `group_migrate_since` - The migrate since of the migration task group (default: null)
* `group_object_overwrite_mode` - The object overwrite mode of the migration task group
  (default: "CRC64_COMPARISON_OVERWRITE")
* `group_consistency_check` - The consistency check of the migration task group (default: "crc64")
* `group_enable_requester_pays` - Whether to enable requester pays for the migration task group (default: true)
* `group_enable_failed_object_recording` - Whether to enable failed object recording for the migration task group
  (default: true)
* `bandwidth_policy_configurations` - The configurations of the bandwidth policy

## Usage

* Copy this example script to your `main.tf`.

* Create a `terraform.tfvars` file and fill in the required variables:

  ```hcl
  bucket_name           = "your_source_obs_bucket_name"
  object_name           = "your_object_name"
  object_upload_content = "Your object content here"
  
  target_bucket_configuration = {
    region     = "your_target_region"
    bucket     = "your_target_bucket_name"
    access_key = "your_target_access_key"
    secret_key = "your_target_secret_key"
  }
  
  bandwidth_policy_configurations = [
    {
      max_bandwidth = 100
      start         = "00:00"
      end           = "23:59"
    }
  ]
  ```

* Initialize Terraform:

  ```bash
  $ terraform init
  ```

* Import the existing resources (optional):

  ```bash
  $ terraform import huaweicloud_kms_key.test[0] xxxxxxxx-xxx-xxx-xxx-xxxxxxxxxxxx
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

## Features

This example demonstrates the following features:

1. **OBS Bucket Creation**: Creates a complete OBS bucket with all necessary components
2. **KMS Encryption**: Enables KMS encryption for enhanced data security
3. **Flexible KMS Key Configuration**: Supports both creating new KMS key and using existing KMS key
4. **Object Upload**: Uploads custom content as an object to the bucket
5. **Storage Class Configuration**: Configurable storage class for cost optimization
6. **Access Control**: Configurable ACL for bucket access management
7. **Tagging Support**: Supports custom tags for resource management
8. **Migration Task Group**: Creates an OMS migration task group for object migration
9. **Bandwidth Control**: Configurable bandwidth policies for migration optimization
10. **Consistency Check**: Configurable consistency check methods for data integrity
11. **Failed Object Recording**: Tracks failed objects during migration process

## Encryption Options

### Option 1: Create New KMS Key

If you don't provide an existing KMS key ID, the example will create a new KMS key with the specified alias:

```hcl
bucket_encryption = true
key_alias         = "your_kms_key_alias"
```

### Option 2: Use Existing KMS Key

If you have an existing KMS key, you can use it directly:

```hcl
bucket_encryption        = true
bucket_encryption_key_id = "your_existing_kms_key_id"
```

### Option 3: Disable Encryption

If you don't need encryption, you can disable it:

```hcl
bucket_encryption = false
```

## Migration Configuration Options

### Target Bucket Configuration

Configure the destination bucket for migration:

```hcl
target_bucket_configuration = {
  region     = "cn-north-4"
  bucket     = "destination-bucket-name"
  access_key = "destination_access_key"
  secret_key = "destination_secret_key"
}
```

### Bandwidth Policy Configuration

Configure bandwidth policies to control migration speed:

```hcl
bandwidth_policy_configurations = [
  {
    max_bandwidth = 100
    start         = "00:00"
    end           = "08:00"
  },
  {
    max_bandwidth = 50
    start         = "08:00"
    end           = "18:00"
  },
  {
    max_bandwidth = 200
    start         = "18:00"
    end           = "23:59"
  }
]
```

### Migration Group Settings

Configure migration group behavior:

```hcl
group_type                           = "PREFIX"
group_enable_kms                     = true
group_object_overwrite_mode          = "CRC64_COMPARISON_OVERWRITE"
group_consistency_check              = "crc64"
group_enable_requester_pays          = true
group_enable_failed_object_recording = true
```

## Storage Classes

The example supports different storage classes for cost optimization:

* `STANDARD` - Standard storage for frequently accessed data (default)
* `WARM` - Infrequent access storage for data accessed less than once per month
* `COLD` - Archive storage for data accessed less than once per year

## Migration Types

The example supports different migration group types:

* `PREFIX` - Migrate objects with specific prefix
* `OBJECT` - Migrate specific objects
* `LIST` - Migrate objects from a list

## Note

* Make sure to keep your credentials secure and never commit them to version control
* The creation of the OBS bucket and object upload is usually instantaneous
* Migration task group creation may take some time depending on the configuration
* This example creates the OBS bucket, optionally a KMS key for encryption, uploads an object with custom content, and
  creates a migration task group
* KMS encryption provides server-side encryption for enhanced data security
* All resources will be created in the specified region
* Bucket names must be globally unique across all HuaweiCloud accounts
* When `bucket_force_destroy` is set to true, the bucket can be destroyed even if it contains objects
* The uploaded object will have content type "application/xml" by default
* The object name will be constructed as `object_name + object_extension_name` if extension is provided
* Migration task groups support bandwidth control to optimize migration performance
* Failed object recording helps track migration issues for troubleshooting

## Requirements

| Name | Version |
| ---- | ---- |
| terraform | >= 1.3.0 |
| huaweicloud | >= 1.64.3 |
