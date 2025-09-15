# Create a CBR vault

This example provides best practice code for using Terraform to create a vault in HuaweiCloud.

## Prerequisites

* A HuaweiCloud account
* Terraform installed
* HuaweiCloud access key and secret key (AK/SK)

## Variable Introduction

The following variables need to be configured:

### Authentication Variables

* `region_name` - The region where the vault is located
* `access_key` - The access key of the IAM user
* `secret_key` - The secret key of the IAM user

### Resource Variables

#### Required Variables

* `volume_type` - The type of the volume
* `volume_size` - The size of the volume, in GB
* `name` - The name of the vault
* `size` - The size of the vault

#### Optional Variables

* `availability_zone` - The availability zone to which the volume belongs (default: "")
  If this parameter is not specified, the available zone will be automatically allocated
* `volume_name` - The name of the volume (default: "")
* `volume_device_type` - The device type of the volume (default: "VBD")
* `type` - The type of the vault (default: "disk")
* `protection_type` - The protection type of the vault (default: "backup")
* `enterprise_project_id` - The enterprise project ID of the vault (default: "0")

## Usage

* Copy this example script to your `main.tf`.

* Create a `terraform.tfvars` file and fill in the required variables:

  ```hcl
  volume_type           = "your_vault_volume_type"
  volume_size           = "your_vault_volume_size"
  name                  = "your_vault_name"
  size                  = "your_vault_size"
  availability_zone     = "your_vault_availability_zone"
  volume_name           = "your_vault_volume_name"
  volume_device_type    = "your_vault_volume_device_type"
  enterprise_project_id = "your_vault_enterprise_project_id"
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
| terraform | >= 0.14.0 |
| huaweicloud | >= 1.61.0 |
