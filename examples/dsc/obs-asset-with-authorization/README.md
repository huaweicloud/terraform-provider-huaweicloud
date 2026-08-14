# Add an OBS Asset to DSC with Authorization

This example provides best practice code for using Terraform to authorize the OBS asset access of
Data Security Center (DSC) and add an OBS bucket as a DSC asset in HuaweiCloud. The example first
enables the OBS authorization of DSC, and then adds an OBS bucket as a DSC asset for sensitive
data identification.

## Prerequisites

* A HuaweiCloud account
* Terraform installed
* HuaweiCloud access key and secret key (AK/SK)
* An existing DSC instance

## Variable Introduction

The following variables need to be configured:

### Authentication Variables

* `region_name` - The region where the DSC OBS asset is located
* `access_key` - The access key of the IAM user
* `secret_key` - The secret key of the IAM user

### Resource Variables

#### Required Variables

* `bucket_name` - The name of the OBS bucket to be added as a DSC asset
* `asset_name` - The name of the DSC OBS asset

## Usage

* Copy this example script to your `main.tf`.

* Create a `terraform.tfvars` file and fill in the required variables:

  ```hcl
  bucket_name = "tf-test-dsc-obs-bucket"
  asset_name  = "tf-test-dsc-obs-asset"
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
* A DSC instance is required in advance, and the OBS bucket must be in the same region as the DSC asset
* The `bucket_policy` argument must be consistent with the actual ACL of the OBS bucket, and both
  `bucket_name` and `bucket_policy` cannot be updated, changing them will recreate the asset
* The `name` argument of the asset can be updated in-place, and it must be unique among the added
  OBS assets
* The asset authorization resource is a switch-style resource per asset type, destroying it only
  removes the resource from the state but does not revoke the authorization, and changing the
  `authorization_status` argument toggles the authorization in-place
* The OBS bucket created in this example sets `force_destroy` to `true`, destroying it will delete
  all objects in the bucket

## Requirements

| Name | Version |
| ---- | ---- |
| terraform | >= 0.14.0 |
| huaweicloud | >= 1.96.0 |
