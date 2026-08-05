# Create a Custom Scan Security Level

This example provides best practice code for using Terraform to create a custom scan security level in HuaweiCloud DSC
(Data Security Center). The security level can be used to classify sensitive data identified by DSC scan rules.

## Prerequisites

* A HuaweiCloud account with DSC permissions
* A purchased DSC instance
* Terraform installed
* HuaweiCloud access key and secret key (AK/SK)

## Variable Introduction

The following variables need to be configured:

### Authentication Variables

* `region_name` - The region where the DSC security level is located
* `access_key` - The access key of the IAM user
* `secret_key` - The secret key of the IAM user

### Resource Variables

#### Required Variables

* `security_level_name` - The name of the security level

#### Optional Variables

* `security_level_color_number` - The color number of the security level displayed on the console (default: 6)
* `security_level_description` - The description of the security level (default: `""`)

## Usage

* Copy this example script to your `main.tf`.

* Create a `terraform.tfvars` file and fill in the required variables:

  ```hcl
  security_level_name = "tfleveltest"
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
* The DSC instance is not created by this example, please purchase it before applying this example
* The `security_level_name` must not be the same as an existing security level name in the current DSC instance
* Changing `security_level_name`, `security_level_color_number` or `security_level_description` updates the existing
  security level instead of creating a new one

## Requirements

| Name        | Version   |
|-------------|-----------|
| terraform   | >= 0.14.0 |
| huaweicloud | >= 1.96.0 |
