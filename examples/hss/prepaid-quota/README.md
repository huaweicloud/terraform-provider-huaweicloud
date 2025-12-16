# Create an HSS quota in prePaid charging mode

This example provides best practice code for using Terraform to create an HSS quota in prePaid charging mode
in HuaweiCloud.

## Prerequisites

* A HuaweiCloud account
* Terraform installed
* HuaweiCloud access key and secret key (AK/SK)

## Variable Introduction

The following variables need to be configured:

### Authentication Variables

* `region_name` - The region where the resources are located
* `access_key`  - The access key of the IAM user
* `secret_key`  - The secret key of the IAM user

### Resource Variables

#### Required Variables

* `quota_version` - The protection quota version
* `period_unit` - The charging period unit of the quota
* `period` - The charging period of the quota

#### Optional Variables

* `is_auto_renew` - Whether auto-renew is enabled (default: false)
* `enterprise_project_id` - The enterprise project ID to which the HSS quota belongs (default: null)
* `quota_tags` - The key/value pairs to associate with the HSS quota (default: null)

## Usage

* Copy this example script to your `main.tf`.

* Create a `terraform.tfvars` file and fill in the required variables:

  ```hcl
  quota_version = "hss.version.enterprise"
  period_unit   = "month"
  period        = 1
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
* All resources will be created in the specified region

## Requirements

| Name | Version   |
| ---- |-----------|
| terraform | >= 0.14.0 |
| huaweicloud | >= 1.66.1 |
