# Restore GaussDB Client Access Authentication Configuration

This example provides best practice code for using Terraform to restore the client access
authentication configuration of a GaussDB instance in HuaweiCloud.

## Prerequisites

* A HuaweiCloud account with GaussDB permissions
* A running GaussDB instance
* Terraform installed
* HuaweiCloud access key and secret key (AK/SK)

## Variable Introduction

The following variables need to be configured:

### Authentication Variables

* `region_name` - The region where the GaussDB instance is located
* `access_key` - The access key of the IAM user
* `secret_key` - The secret key of the IAM user

### Resource Variables

#### Required Variables

* `instance_id` - The ID of the GaussDB instance

#### Optional Variables

* `hba_history_id` - The client access authentication modification history record ID.
  If empty, it means restoring to the default configuration.

## Usage

* Copy this example script to your `main.tf`.

* Create a `terraform.tfvars` file and fill in the required variables:

  ```hcl
  instance_id = "your_instance_id"
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

## Note

* Make sure to keep your credentials secure and never commit them to version control
* Deleting this resource only removes it from the Terraform state, it does not undo the restore operation
* If `hba_history_id` is not specified, the instance will be restored to the default client access authentication configuration

## Requirements

| Name | Version   |
| ---- |-----------|
| terraform | >= 0.14.0 |
| huaweicloud | >= 1.94.0 |
