# Create a CC global connection bandwidth example

This example provides best practice code for using Terraform to create a Cloud Connect (CC) global connection bandwidth
in HuaweiCloud.

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

* `global_connection_bandwidth_name` - The name of the global connection bandwidth
* `bandwidth_type` - The type of the global connection bandwidth
* `bordercross` - Whether the global connection bandwidth crosses borders
* `bandwidth_size` - Bandwidth size of the global connection bandwidth
* `charge_mode` - Billing option of the global connection bandwidth

#### Optional Variables

* `global_connection_bandwidth_description` - The description about the global connection bandwidth (default: "Created
  by Terraform")
* `enterprise_project_id` - ID of the enterprise project that the global connection bandwidth belongs to (default: "0")
* `global_connection_bandwidth_tags` - The tags of the global connection bandwidth (default: { "Owner" = "terraform",
  "Env" = "test" })

## Usage

* Copy this example script to your `main.tf`.

* Create a `terraform.tfvars` file and fill in the required variables:

  ```hcl
  global_connection_bandwidth_name = "your_global_connection_bandwidth_name"
  bandwidth_type                   = "your_bandwidth_type"
  bordercross                      = your_bordercross_value
  bandwidth_size                   = your_bandwidth_size
  charge_mode                      = "your_charge_mode"
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
| ---- | ---- |
| terraform | >= 1.9.0 |
| huaweicloud | >= 1.77.2 |
