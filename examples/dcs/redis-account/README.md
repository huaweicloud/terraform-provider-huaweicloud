# DCS Redis account management

This example provides best practice code for using Terraform to manage a DCS Redis account.

## Prerequisites

* A HuaweiCloud account
* Terraform installed
* HuaweiCloud access key and secret key (AK/SK)

## Required Variables

The following variables need to be configured:

### Authentication Variables

* `region_name` - The region where the DCS Redis instance is located
* `access_key`  - The access key of the IAM user
* `secret_key`  - The secret key of the IAM user

### Resource Variables

#### Required Variables

* `vpc_name` - The name of the VPC
* `subnet_name` - The name of the subnet
* `instance_name` - The name of the Redis single instance
* `account_name` - The name of the DCS account

#### Optional Variables

* `vpc_cidr` - The CIDR block of the VPC (default: "192.168.0.0/16")
* `subnet_cidr` - The CIDR block of the subnet (default: "")
* `subnet_gateway_ip` - The gateway IP address of the subnet (default: "")
* `availability_zone` - The availability zone to which the Redis single instance belongs (default: "")
* `instance_flavor_id` - The flavor ID of the Redis single instance (default: "")
* `instance_capacity` - The capacity of the Redis instance (default: 4)
* `instance_engine_version` - The engine version of the Redis single instance (default: "5.0")
* `instance_password` - The password for the Redis instance (default: null)
* `account_role` - The role of the account: read or write (default: "read")
* `account_password` - The password of the DCS account (default: null)
* `account_description` - The description of the account (default: "")

## Usage

* Copy this example script to your `main.tf`.

* Create a `terraform.tfvars` file and fill in the required variables:

  ```hcl
  vpc_name      = "your_vpc_name"
  subnet_name   = "your_subnet_name"
  instance_name = "your_instance_name"
  account_name  = "tf_test_account"
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
* `instance_id` and `account_name` are non-updatable, changing them will recreate the resource
* `account_role`, `account_password`, and `description` can be updated in-place
* The `account_role` parameter accepts "read" (read-only) or "write" (read and write)
* The `account_password` is sensitive and will not be displayed in the plan output
* The resource supports import using the format `<instance_id>/<account_id>`
* The create, update, and delete operations each have a default timeout of 30 minutes

## Requirements

| Name        | Version   |
|-------------|-----------|
| terraform   | >= 1.1.0  |
| huaweicloud | >= 1.62.0 |
