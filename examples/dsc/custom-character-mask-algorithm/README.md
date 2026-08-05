# Create a Custom Character Mask Algorithm

This example provides best practice code for using Terraform to create a custom character mask algorithm in HuaweiCloud
DSC (Data Security Center). The algorithm uses the `PRESNM` algorithm with the overwrite type (`MASK_BY_OVERWRITE`) to
retain specified characters at the beginning and end of the source data and masks the characters in between with a
replacement character.

## Prerequisites

* A HuaweiCloud account with DSC permissions
* A purchased DSC instance
* Terraform installed
* HuaweiCloud access key and secret key (AK/SK)

## Variable Introduction

The following variables need to be configured:

### Authentication Variables

* `region_name` - The region where the DSC mask algorithm is located
* `access_key` - The access key of the IAM user
* `secret_key` - The secret key of the IAM user

### Resource Variables

#### Required Variables

* `mask_algorithm_name` - The name of the mask algorithm

#### Optional Variables

* `mask_algorithm_prefix_length` - The number of characters to retain at the beginning of the data (default: 6)
* `mask_algorithm_suffix_length` - The number of characters to retain at the end of the data (default: 4)
* `mask_algorithm_replacement` - The character used to mask the data (default: "*")

## Usage

* Copy this example script to your `main.tf`.

* Create a `terraform.tfvars` file and fill in the required variables:

  ```hcl
  mask_algorithm_name = "tfmaskalgorithm"
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
* The `mask_algorithm_name` must not be the same as an existing mask algorithm name in the current DSC instance
* Changing the variables in this example updates the existing mask algorithm instead of creating a new one
* This example is fixed to the character overwrite combination (`PRESNM` with `MASK_BY_OVERWRITE`). For other algorithm
  combinations, refer to the `huaweicloud_dsc_mask_algorithm` resource documentation

## Requirements

| Name        | Version   |
|-------------|-----------|
| terraform   | >= 0.14.0 |
| huaweicloud | >= 1.96.0 |
