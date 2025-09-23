# Create a CTS Data Tracker

This example provides best practice code for using Terraform to create a CTS (Cloud Trace Service) data tracker in
HuaweiCloud.

## Prerequisites

* A HuaweiCloud account
* Terraform installed
* HuaweiCloud access key and secret key (AK/SK)
* CTS service enabled in the target region
* OBS service enabled in the target region
* Authorize CTS to use the OBS service as an administrator

## Variable Introduction

The following variables need to be configured:

### Authentication Variables

* `region_name` - The region where the CTS service is located
* `access_key` - The access key of the IAM user
* `secret_key` - The secret key of the IAM user

### Resource Variables

#### Required Variables

* `source_bucket_name` - The name of the OBS bucket for storing trace files
* `transfer_bucket_name` - The name of the OBS bucket for transferring trace files
* `tracker_name` - The name of the system tracker
* `trace_object_prefix` - The prefix of the trace object in the OBS bucket

#### Optional Variables

* `tracker_enabled` - Whether to enable the system tracker (default: true)
* `tracker_tags` - The tags of the system tracker (default: {})
* `trace_file_compression_type` - The compression type of the trace file (default: "gzip")
* `is_lts_enabled` - Whether to enable the trace analysis for LTS service (default: true)

## Usage

* Copy this example script to your `main.tf`.

* Create a `terraform.tfvars` file and fill in the required variables:

  ```hcl
  source_bucket_name   = "tf-test-source-bucket"
  transfer_bucket_name = "tf-test-transfer-bucket"
  tracker_name         = "tf-test-tracker"
  trace_object_prefix  = "tf_test"
  tracker_tags         = {
    "owner" = "terraform"
  }
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

## Notes

* Make sure to keep your credentials secure and never commit them to version control
* The CTS data tracker is dependent on the OBS buckets
* When a trigger is created, a log group named CTS is created and a corresponding log stream is created within it.
  If the relevant resources exist, they are used directly.

## Requirements

| Name | Version |
| ---- | ---- |
| terraform | >= 0.12.0 |
| huaweicloud | >= 1.61.1 |
