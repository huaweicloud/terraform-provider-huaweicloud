# Create a COC script

This example provides best practice code for using Terraform to create a COC script in HuaweiCloud.

## Prerequisites

* A HuaweiCloud account
* Terraform installed
* HuaweiCloud access key and secret key (AK/SK)

## Variable Introduction

The following variables need to be configured:

### Authentication Variables

* `region_name` - The region where the COC script is located
* `access_key`  - The access key of the IAM user
* `secret_key`  - The secret key of the IAM user

### Resource Variables

#### Required Variables

* `coc_script_name` - The name of the script
* `coc_script_description` - The description of the script
* `coc_script_risk_level` - The risk level of the script
* `coc_script_version` - The version of the script
* `coc_script_type` - The type of the script
* `coc_script_content` - The content of the script
* `coc_script_parameters` - The parameter list of the script
  + `name` - The name of the parameter
  + `value` - The value of the parameter
  + `description` - The description of the parameter
  + `sensitive` - Whether the parameter is sensitive

## Usage

* Copy this example script to your `main.tf`.

* Create a `terraform.tfvars` file and fill in the required variables:

  ```hcl
  coc_script_name        = "your_coc_script_name"
  coc_script_description = "your_coc_script_description"
  coc_script_risk_level  = "your_coc_script_risk_level"
  coc_script_version     = "your_coc_script_version"
  coc_script_type        = "your_coc_script_type"
  coc_script_content     = "your_coc_script_content"
  coc_script_parameters  = "your_coc_script_parameters"
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

| Name | Version |
|---|---|
| terraform | >= 1.3.0 |
| huaweicloud | >= 1.58.0 |
