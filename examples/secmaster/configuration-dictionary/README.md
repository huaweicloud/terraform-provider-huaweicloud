# Create a SecMaster Configuration Dictionary

This example provides best practice code for using Terraform to create a SecMaster configuration
dictionary within HuaweiCloud.

## Prerequisites

* A HuaweiCloud account
* Terraform installed
* HuaweiCloud access key and secret key (AK/SK)

## Variable Introduction

The following variables need to be configured:

### Authentication Variables

* `region_name` - The region where the SecMaster resources will be created
* `access_key` - The access key of the IAM user
* `secret_key` - The secret key of the IAM user

### Resource Variables

#### Required Variables

* `dict_id` - The dictionary ID
* `dict_key` - The dictionary key
* `dict_code` - The dictionary code
* `dict_val` - The dictionary value

#### Optional Variables

* `language` - The language environment, valid values: zh, en (default: "zh")
* `dict_version` - The version number (default: "1.0.0")
* `dict_pkey` - The parent key of the dictionary (default: "")
* `dict_pcode` - The parent code of the dictionary (default: "")
* `scope` - The domain to which the dictionary belongs (default: "ALERT")
* `description` - The description of the dictionary (default: "")
* `extend_field` - The extension field of the dictionary (default: {})

## Usage

* Copy this example script to your `main.tf`.

* Create a `terraform.tfvars` file and fill in the required variables:

  ```hcl
  dict_id   = "3027"
  dict_key  = "alert_comments"
  dict_code = "Open"
  dict_val  = "Open"
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
* The `is_built_in` is set to `false` to create a user-defined dictionary
* The `dict_id`, `dict_key`, `language`, `version`, `scope` are non-updatable after creation
* The `is_built_in` attribute may drift after import, so it is included in
  `lifecycle.ignore_changes`
* The configuration dictionary can be imported using its `id`

## Requirements

| Name | Version  |
| ---- |----------|
| terraform | >= 0.14.0 |
| huaweicloud | >= 1.94.0 |
