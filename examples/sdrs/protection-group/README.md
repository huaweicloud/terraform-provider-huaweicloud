# Create a SDRS protection group

This example provides best practice code for using Terraform to create a SDRS protection group in HuaweiCloud.

## Prerequisites

* A HuaweiCloud account
* Terraform installed
* HuaweiCloud access key and secret key (AK/SK)

## Variable Introduction

The following variables need to be configured:

### Authentication Variables

* `region_name` - The region where the protection group located
* `access_key`  - The access key of the IAM user
* `secret_key`  - The secret key of the IAM user

### Resource Variables

#### Required Variables

* `vpc_name` - The VPC name
* `protection_group_name` - The protection group name

#### Optional Variables

* `vpc_cidr` - The CIDR block of the VPC (default: "192.168.0.0/16")
* `source_availability_zone` - The production site AZ of the protection group (default: "")  
  If this parameter is not specified, the available zone will be automatically allocated.
  Both `source_availability_zone` and `target_availability_zone` must be set, or both must be empty
* `target_availability_zone` - The disaster recovery site AZ of the protection group (default: "")
* `protection_group_dr_type` - The deployment model (default: null)
* `protection_group_description` - The description of the protection group (default: null)
* `protection_group_enable` - Whether enable the protection group start protecting (default: null)

## Usage

* Copy this example script to your `main.tf`.

* Create a `terraform.tfvars` file and fill in the required variables:

  ```hcl
  vpc_name              = "your_vpc_name"
  protection_group_name = "your_protection_group_name"
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

| Name | Version   |
| ---- |-----------|
| terraform | >= 1.9.0  |
| huaweicloud | >= 1.77.0 |
