# Create a Log Transfer to OBS bucket

This example provides best practice code for using Terraform to create a Log Transfer that used to transfers logs to
OBS bucket in HuaweiCloud.

## Prerequisites

* A HuaweiCloud account
* Terraform installed
* HuaweiCloud access key and secret key (AK/SK)
* LTS service enabled in the target region
* OBS service enabled in the target region

## Variable Introduction

The following variables need to be configured:

### Authentication Variables

* `region_name` - The region where the LTS service is located
* `access_key` - The access key of the IAM user
* `secret_key` - The secret key of the IAM user

### Resource Variables

#### Required Variables

* `group_name` - The name of the log group
* `stream_name` - The name of the log stream
* `bucket_name` - The name of the OBS bucket

#### Optional Variables

* `group_log_expiration_days` - The log expiration days of the log group (default: 14)
* `stream_log_expiration_days` - The log expiration days of the log stream (default: null)
  **null** or **-1** means **null** or **-1** indicates that the log expiration days is consistent with the log group
* `transfer_type` - The type of the log transfer (default: "OBS")
* `transfer_mode` - The mode of the log transfer (default: "cycle")
* `transfer_storage_format` - The storage format of the log transfer (default: "JSON")
* `transfer_status` - The status of the log transfer (default: "ENABLE")
* `bucket_dir_prefix_name` - The prefix path of the OBS transfer task
  (default: "LTS-test/%GroupName/%StreamName/%Y/%m/%d/%H/%M")
* `bucket_time_zone` - The time zone of the OBS bucket (default: "UTC")
* `bucket_time_zone_id` - The time zone ID of the OBS bucket (default: "Etc/GMT")

## Usage

* Copy this example script to your `main.tf`.

* Create a `terraform.tfvars` file and fill in the required variables:

  ```hcl
  group_name  = "tf_test_log_group"
  stream_name = "tf_test_log_stream"
  bucket_name = "tf-test-log-transfer-obs-bucket"
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
* The log transfer is dependent on the log group, log stream, and OBS bucket
* Please read the implicit and explicit dependencies in the script carefully
* All resources will be created in the specified region
* The log transfer supports different storage formats: JSON, RAW, etc.
* The OBS bucket will be created with private ACL and force_destroy enabled

## Requirements

| Name | Version |
| ---- | ---- |
| terraform | >= 0.12.0 |
| huaweicloud | >= 1.55.0 |
