# Create a CC bandwidth package example

This example provides best practice code for using Terraform to create a Cloud Connect (CC) bandwidth package in HuaweiCloud.

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

* `bandwidth_package_name` - The name of the bandwidth package
* `local_area_id` - The local area ID
* `remote_area_id` - The remote area ID
* `charge_mode` - Billing option of the bandwidth package
* `billing_mode` - Billing mode of the bandwidth package
* `bandwidth` - Bandwidth in the bandwidth package

#### Optional Variables

* `bandwidth_package_description` - The description about the bandwidth package (default: "Created by Terraform")
* `enterprise_project_id` - ID of the enterprise project that the bandwidth package belongs to (default: "0")
* `bandwidth_package_tags` - The tags of the bandwidth package (default: { "Owner" = "terraform", "Env" = "test" })

## Usage

* Copy this example script to your `main.tf`.

* Create a `terraform.tfvars` file and fill in the required variables:

  ```hcl
  bandwidth_package_name = "your_bandwidth_package_name"
  local_area_id          = "your_local_area_id"
  remote_area_id         = "your_remote_area_id"
  charge_mode            = "your_charge_mode"
  billing_mode           = "your_billing_mode"
  bandwidth              = your_bandwidth
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

| Name | Version     |
| ---- |-------------|
| terraform | >= 1.9.0    |
| huaweicloud | >= 1.78.3 |
