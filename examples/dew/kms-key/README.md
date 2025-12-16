# Create a KMS key

This example provides best practice code for using Terraform to create a KMS key in HuaweiCloud.

## Prerequisites

* A HuaweiCloud account
* Terraform installed
* HuaweiCloud access key and secret key (AK/SK)

## Variables Introduction

The following variables need to be configured:

### Authentication Variables

* `region_name` - The region where the KMS key is located
* `access_key` - The access key of the IAM user
* `secret_key` - The secret key of the IAM user

### Resource Variables

#### Required Variables

* `key_name` - The alias name of the KMS key

#### Optional Variables

* `key_algorithm` - The generation algorithm of the KMS key (default: "AES_256")
* `key_usage` - The usage of the KMS key (default: "ENCRYPT_DECRYPT")
* `key_source` - The source of the KMS key (default: "kms")
* `key_description` - The description of the KMS key (default: "")
* `enterprise_project_id` - The ID of the enterprise project to which the KMS key belongs (default: null)
* `key_tags` - The key/value pairs to associate with the KMS key (default: {})
* `key_schedule_time` - The number of days after which the KMS key is scheduled to be deleted (default: "7")
  The valid value rangs from `7` to `1,096`

## Usage

* Copy this example script to your `main.tf`.

* Create a `terraform.tfvars` file and fill in the required variables:

  ```hcl
  key_name = "tf_test_Key"
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

## Requirements

| Name | Version |
| ---- | ---- |
| terraform | >= 0.14.0 |
| huaweicloud | >= 1.64.3 |
