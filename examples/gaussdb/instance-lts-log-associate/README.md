# Associate LTS Log with GaussDB Instance

This example provides best practice code for using Terraform to associate LTS log with a
GaussDB instance in HuaweiCloud.

## Prerequisites

* A HuaweiCloud account with GaussDB and LTS permissions
* A running GaussDB instance
* An LTS log group and log stream
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

* `lts_group_name` - The LTS log group name
* `lts_stream_name` - The LTS log stream name
* `gaussdb_instance_id` - The ID of the GaussDB instance

#### Optional Variables

* `lts_log_type` - The LTS log type. The valid value is `audit_log`

## Usage

* Copy this example script to your `main.tf`.

* Create a `terraform.tfvars` file and fill in the required variables:

  ```hcl
  lts_group_name      = "your_lts_group_name"
  lts_stream_name     = "your_lts_stream_name"
  gaussdb_instance_id = "your_instance_id"
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
* The `log_type` only supports `audit_log`
* The LTS log group and stream can be obtained through the LTS console or API

## Requirements

| Name        | Version   |
|-------------|-----------|
| terraform   | >= 0.14.0 |
| huaweicloud | >= 1.94.0 |
