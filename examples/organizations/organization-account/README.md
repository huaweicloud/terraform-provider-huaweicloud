# Create an Organizations Account example

This example provides best practice code for using Terraform to create an Organizations account in HuaweiCloud.

## Prerequisites

* A HuaweiCloud account with Organizations service permissions
* Terraform installed
* HuaweiCloud access key and secret key (AK/SK)

## Variables Introduction

The following variables need to be configured:

### Authentication Variables

* `region_name` - The region where the organization account is located
* `access_key` - The access key of the IAM user
* `secret_key` - The secret key of the IAM user

### Resource Variables

#### Required Variables

* `name` - The name of the organization account
* `email` - The email address of the organization account

#### Optional Variables

* `phone` - The mobile number of the account.
* `agency_name` - The agency name of the account.
* `description` - The description of the account.
* `parent_id` - The ID of the root or organization unit in which you want to create a new account. The default is root ID.
* `tags` - The tags of the account.

## Usage

* Copy this example script to your `main.tf`.

* Create a `terraform.tfvars` file and fill in the required variables:

  ```hcl
  name  = "your_name"
  email = "your_email@example.com"
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
* The account email must be unique and not already associated with another HuaweiCloud account
* The account creation process may take a few minutes to complete
* All resources will be created in the specified region

## Requirements

| Name | Version |
| ---- | ---- |
| terraform | >= 1.9.0 |
| huaweicloud | >= 1.77.6 |
