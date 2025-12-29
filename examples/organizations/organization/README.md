# Create an Organization

This example provides best practice code for using Terraform to create an organization in HuaweiCloud Organizations service.

## Prerequisites

* A HuaweiCloud account with Organizations permissions
* Terraform installed
* HuaweiCloud access key and secret key (AK/SK)

## Variable Introduction

The following variables need to be configured:

### Authentication Variables

* `region_name` - The region where the Organizations service is located
* `access_key` - The access key of the IAM user
* `secret_key` - The secret key of the IAM user

### Resource Variables

#### Required Variables

#### Optional Variables

* `enabled_policy_types` - The list of Organizations policy types to enable in the Organization Root
* `root_tags` - The key/value to attach to the root.

## Usage

* Copy this example script to your `main.tf`.

* Create a `terraform.tfvars` file and fill in the required variables:

  ```hcl
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
* Creating an organization is a critical operation that affects your entire account structure
* The organization's management account will be the account used to create the organization
* Enabled policy types allow you to apply service control policies (SCPs) and tag policies to organizational units and accounts
* Root tags are applied to the organization's root for easier resource management and identification

## Requirements

| Name | Version |
| ---- | ---- |
| terraform | >= 0.14.0 |
| huaweicloud | >= 1.57.0 |
