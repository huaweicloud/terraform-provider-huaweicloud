# Create TMS Preset Tags

This example provides best practice code for using Terraform to create preset tags (predefine tags) in HuaweiCloud Tag
Management Service (TMS).

## Prerequisites

* A HuaweiCloud account with TMS permissions
* Terraform installed
* HuaweiCloud access key and secret key (AK/SK)

## Variable Introduction

The following variables need to be configured:

### Authentication Variables

* `region_name` - The region where the TMS service is located
* `access_key` - The access key of the IAM user
* `secret_key` - The secret key of the IAM user

### Resource Variables

#### Required Variables

* `preset_tags` - The preset tags
  + `key` - The tag key
  + `value` - The tag key

## Usage

* Copy this example script to your `main.tf`.

* Create a `terraform.tfvars` file and fill in the required variables:

  ```hcl
  preset_tags = [
    {
      key   = "foo"
      value = "bar"
    },
    {
      key   = "owner"
      value = "terraform"
    }
  ]
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

* To clean up the resources (delete preset tags):

  ```bash
  $ terraform destroy
  ```

## Note

* Make sure to keep your credentials secure and never commit them to version control
* Preset tags are created at the account level and are available across all regions
* The maximum number of preset tags per account is limited (check TMS service limits)

## Requirements

| Name | Version |
| ---- | ---- |
| terraform | >= 0.14.0 |
| huaweicloud | >= 1.35.0 |
