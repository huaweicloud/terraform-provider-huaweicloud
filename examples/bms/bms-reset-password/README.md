# Reset the password of a BMS instance

This example provides best practice code for using Terraform to reset the password of a BMS instance.
Destroying this module will not clear the corresponding request record, but will only remove the resource informations
from the tf state file.

## Prerequisites

* A HuaweiCloud account
* Terraform installed
* HuaweiCloud access key and secret key (AK/SK)

## Variable Introduction

The following variables need to be configured:

### Authentication Variables

* `region_name` - The region where the BMS instance is located
* `access_key`  - The access key of the IAM user
* `secret_key`  - The secret key of the IAM user

### Resource Variables

#### Required Variables

* `bms_instance_id` - The ID of the BMS instance
* `bms_instance_new_password` - The new password of the BMS instance

#### Optional Variables

* None

## Usage

* Copy this example script to your `main.tf`.

* Create a `terraform.tfvars` file and fill in the required variables:

```hcl
bms_instance_id           = "your_bms_instance_id"
bms_instance_new_password = "your_new_password"
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

## Note

* Make sure to keep your credentials secure and never commit them to version control
* All resources will be created in the specified region

## Requirements

| Name | Version |
| ---- | ---- |
| terraform | >= 0.14.0 |
| huaweicloud | >= 1.82.5 |
