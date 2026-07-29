# Create a Database Account for a GaussDB Instance

This example provides best practice code for using Terraform to create a database account for a
GaussDB instance in HuaweiCloud.

## Prerequisites

* A HuaweiCloud account with GaussDB permissions
* A running GaussDB instance
* Terraform installed
* HuaweiCloud access key and secret key (AK/SK)

## Variable Introduction

The following variables need to be configured:

### Authentication Variables

* `region_name` - The region where the GaussDB instance is located
* `access_key` - The access key of the IAM user
* `secret_key` - The secret key of the IAM user

### Resource Variables

#### Required Variables

* `vpc_name` - The VPC name
* `subnet_name` - The subnet name
* `security_group_name` - The security group name
* `instance_name` - The GaussDB instance name
* `database_account_name` - The name of the database account
* `database_account_password` - The password of the database account

#### Optional Variables

* `vpc_cidr` - The CIDR block of the VPC (default: "192.168.0.0/16")
* `enterprise_project_id` - The ID of the enterprise project (default: "")
* `subnet_cidr` - The CIDR block of the subnet (default: "", auto-calculated from VPC CIDR)
* `gateway_ip` - The gateway IP address of the subnet (default: "", auto-calculated)
* `instance_password` - The password for the GaussDB instance (default: "", auto-generated)
* `instance_volume_type` - The storage volume type (default: "ULTRAHIGH")
* `instance_volume_size` - The storage volume size in GB (default: 40)
* `is_login_only` - Whether the database account supports login only. The valid values are `true` and `false`

## Usage

* Copy this example script to your `main.tf`.

* Create a `terraform.tfvars` file and fill in the required variables:

  ```hcl
  vpc_name                  = "your_vpc_name"
  subnet_name               = "your_subnet_name"
  security_group_name       = "your_security_group_name"
  instance_name             = "your_mysql_instance_name"
  database_account_name     = "your_account_name"
  database_account_password = "your_account_password"
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
* The account name must be 1 to 63 characters, including letters, digits, and underscores, and cannot
  start with `pg` or a digit
* The account name cannot be the same as system user names (`rdsAdmin`, `rdsMetric`, `rdsBackup`,
  `rdsRepl`, and `root`)
* The password must contain 8 to 32 characters and at least three types of the following characters:
  uppercase letters, lowercase letters, digits, and special characters (`~!@#%^*-_=+?`)
* The password cannot be the same as the account name or the reverse of the account name
* The `is_login_only` only supports `true` and `false`

## Requirements

| Name        | Version   |
|-------------|-----------|
| terraform   | >= 0.14.0 |
| huaweicloud | >= 1.94.0 |
