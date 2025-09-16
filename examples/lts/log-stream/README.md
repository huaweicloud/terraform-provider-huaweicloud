# Create LTS Log Group and Log Stream

This example provides best practice code for using Terraform to create a Log Tank Service (LTS) log group and log
stream in HuaweiCloud.

## Prerequisites

* A HuaweiCloud account
* Terraform installed
* HuaweiCloud access key and secret key (AK/SK)
* LTS service enabled in the target region

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

#### Optional Variables

* `group_log_expiration_days` - The log expiration days of the log group (default: 14)
* `group_tags` - The tags of the log group (default: {})
* `enterprise_project_id` - The enterprise project ID of the log group (default: null)
* `stream_log_expiration_days` - The log expiration days of the log stream (default: null)
  **null** or **-1** means **null** or **-1** indicates that the log expiration days is consistent with the log group
* `stream_tags` - The tags of the log stream (default: {})
* `stream_is_favorite` - Whether to favorite the log stream (default: false)

## Usage

* Copy this example script to your `main.tf`.

* Create a `terraform.tfvars` file and fill in the required variables:

  ```hcl
  group_name  = "tf_test_log_group"
  stream_name = "tf_test_log_stream"
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
* The log stream is dependent on the log group
* Please read the implicit and explicit dependencies in the script carefully
* All resources will be created in the specified region
* The log group and log stream support enterprise project isolation
* Tags can be used for resource categorization and cost allocation

## Requirements

| Name | Version |
| ---- | ---- |
| terraform | >= 0.12.0 |
| huaweicloud | >= 1.73.7 |
