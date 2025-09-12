# Create a RMS assignment package example

This example provides best practice code for using Terraform to create an RMS (Resource Management Service)
assignment package in HuaweiCloud.

## Prerequisites

* A HuaweiCloud account
* Terraform installed
* HuaweiCloud access key and secret key (AK/SK)

## Variables Introduction

The following variables need to be configured:

### Authentication Variables

* `region_name` - The region where the resources are located
* `access_key` - The access key of the IAM user
* `secret_key` - The secret key of the IAM user

### Resource Variables

#### Required Variables

* `assignment_package_name` - The name of the assignment package
* `template_key` - The name of a built-in assignment package template

#### Optional Variables

* `assignment_package_vars` - The parameters of the assignment package
  - `var_key` - The name of a parameter
  - `var_value` - The value of a parameter

## Usage

* Copy this example script to your `main.tf`.

* Create a `terraform.tfvars` file and fill in the required variables:

  ```hcl
  assignment_package_name = "your_assignment_package_name"
  template_key            = "your_template_key"
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
* The assignment package will be created based on the specified template
* All resources will be created in the specified region

## Requirements

| Name | Version |
| ---- | ---- |
| terraform | >= 1.9.0 |
| huaweicloud | >= 1.77.6 |
