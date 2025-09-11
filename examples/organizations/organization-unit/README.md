# Create an Organizations organizational unit example

This example provides best practice code for using Terraform to create an organizational unit in HuaweiCloud
Organizations service.

## Prerequisites

* A HuaweiCloud account with Organizations service enabled
* Terraform installed
* HuaweiCloud access key and secret key (AK/SK)

## Variables Introduction

The following variables need to be configured:

### Authentication Variables

* `region_name` - The region where the organizational unit will be created
* `access_key` - The access key of the IAM user
* `secret_key` - The secret key of the IAM user

### Resource Variables

#### Required Variables

* `organizational_unit_name` - The name of the organizational unit.

#### Optional Variables

* `tags` - The key/value to attach to the organizational unit.

## Usage

* Copy this example script to your `main.tf`.

* Create a `terraform.tfvars` file and fill in the required variables:

  ```hcl
  organizational_unit_name = "your_organizational_unit_name"
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

* Make sure Organizations service is enabled in your account before running this example
* Make sure to keep your credentials secure and never commit them to version control
* The organizational unit will be created under the root organizational unit of your organization

## Requirements

| Name | Version |
| ---- | ---- |
| terraform | >= 1.9.0 |
| huaweicloud | >= 1.77.6 |
