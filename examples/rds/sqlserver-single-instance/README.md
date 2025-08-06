# Create a SQLServer RDS instance

This example provides best practice code for using Terraform to create a SQLServer RDS instance in HuaweiCloud RDS service.

## Prerequisites

* A HuaweiCloud account
* Terraform installed
* HuaweiCloud access key and secret key (AK/SK)

## Variable Introduction

The following variables need to be configured:

### Authentication Variables

* `region_name` - The region where the SQLServer RDS instance is located
* `access_key` - The access key of the IAM user
* `secret_key` - The secret key of the IAM user

### Resource Variables

#### Required Variables

* `vpc_name` - The VPC name
* `subnet_name` - The subnet name
* `security_group_name` - The security group name
* `instance_name` - The name of the RDS instance

#### Optional Variables

* `vpc_cidr` - The CIDR block of the VPC (default: "192.168.0.0/16")
* `availability_zone` - The availability zone to which the RDS instance belongs (default: "")
* `subnet_cidr` - The CIDR block of the subnet (default: "")
* `gateway_ip` - The gateway IP address of the subnet (default: "")
* `instance_flavor_id` - The flavor ID of the RDS instance (default: "")
* `instance_db_type` - The database engine type (default: "SQLServer")
* `instance_db_version` - The database engine version (default: "2019_SE")
* `instance_mode` - The instance mode for the RDS instance flavor (default: "single")
* `instance_flavor_group_type` - The group type for the RDS instance flavor (default: "general")
* `instance_flavor_vcpus` - The CPU core numbers for the RDS instance flavor (default: 2)
* `instance_flavor_memory` - The memory size for the RDS instance flavor (default: 4)
* `instance_db_port` - The database port (default: 1433)
* `instance_password` - The password for the RDS instance (default: "")
* `instance_volume_type` - The storage volume type (default: "CLOUDSSD")
* `instance_volume_size` - The storage volume size in GB (default: 40)
* `instance_backup_time_window` - The backup time window in HH:MM-HH:MM format (default: "03:00-04:00")
* `instance_backup_keep_days` - The number of days to retain backups (default: 7)

## Usage

* Copy this example script to your `main.tf`.

* Create a `terraform.tfvars` file and fill in the required variables:

  ```hcl
  vpc_name            = "your_vpc_name"
  subnet_name         = "your_subnet_name"
  security_group_name = "your_security_group_name"
  instance_name       = "your_sqlserver_instance_name"
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

## Features

This example demonstrates the following features:

1. **SQLServer RDS Instance Creation**: Creates a complete SQLServer RDS instance with all necessary components
2. **Network Configuration**: Sets up VPC, subnet, and security group for the RDS instance
3. **Backup Strategy**: Configures automated backup with customizable time window and retention period
4. **Security**: Creates security group rules to allow database access from VPC

## Note

* Make sure to keep your credentials secure and never commit them to version control
* The creation of the SQLServer RDS instance takes about 5 minutes
* This example creates the SQLServer RDS instance, VPC, subnet, security group, and backup strategy
* All resources will be created in the specified region
* If no password is provided, a random password will be generated

## Requirements

| Name | Version |
| ---- | ---- |
| terraform | >= 0.12.0 |
| huaweicloud | >= 1.40.0 |
| random | >= 3.0.0 |
