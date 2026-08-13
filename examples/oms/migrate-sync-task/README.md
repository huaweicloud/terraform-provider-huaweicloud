# Create a migration synchronization task

This example provides best practice code for using Terraform to create an OMS migration synchronization task
in HuaweiCloud.

## Prerequisites

* A HuaweiCloud account
* Terraform installed
* HuaweiCloud access key and secret key (AK/SK)

## Variable Introduction

The following variables need to be configured:

### Authentication Variables

* `region_name` - The region where the OMS migration synchronization task is located
* `access_key` - The access key of the IAM user
* `secret_key` - The secret key of the IAM user

### Resource Variables

#### Required Variables

* `source_bucket_name` - The name of the source OBS bucket
* `dest_bucket_name` - The name of the destination OBS bucket
* `source_region` - The region where the source bucket is located
* `source_access_key` - The access key for accessing the source bucket
* `source_secret_key` - The secret key for accessing the source bucket
* `dest_access_key` - The access key for accessing the destination bucket
* `dest_secret_key` - The secret key for accessing the destination bucket

#### Optional Variables

* `bucket_storage_class` - The storage class of the OBS bucket (default: "STANDARD")
* `bucket_acl` - The ACL of the OBS bucket (default: "private")
* `bucket_force_destroy` - Whether to force destroy the OBS bucket (default: true)
* `source_cloud_type` - The source cloud service provider (default: "HuaweiCloud")
* `task_description` - The description of the migration synchronization task (default: "")
* `consistency_check` - The consistency check method (default: "size_last_modified")
* `enable_metadata_migration` - Whether to enable metadata migration (default: false)

## Usage

* Copy the example `main.tf` file to your configuration.

* Create a `terraform.tfvars` file and fill in the required variables:

  ```hcl
  source_bucket_name = "your_source_bucket_name"
  dest_bucket_name   = "your_destination_bucket_name"
  source_region      = "your_source_bucket_region"
  source_access_key  = "your_source_bucket_ak"
  source_secret_key  = "your_source_bucket_sk"
  dest_access_key    = "your_destination_ak"
  dest_secret_key    = "your_destination_sk"
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

* Keep your credentials secure and never commit them to version control.
* All resources are created in the specified region.

## Requirements

| Name | Version |
| ---- | ---- |
| terraform | >= 0.14.0 |
| huaweicloud | >= 1.61.0 |
