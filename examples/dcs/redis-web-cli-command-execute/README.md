# DCS Redis Web CLI command execute

This example provides best practice code for using Terraform to execute Redis commands via the Web CLI of a DCS
instance.

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

#### Optional Variables

* `vpc_cidr` - The CIDR block of the VPC (default: "192.168.0.0/16")
* `subnet_cidr` - The CIDR block of the subnet (default: "")
* `subnet_gateway_ip` - The gateway IP address of the subnet (default: "")
* `availability_zone` - The availability zone to which the Redis single instance belongs (default: "")
* `instance_flavor_id` - The flavor ID of the Redis single instance (default: "")
* `instance_capacity` - The capacity of the Redis instance (default: 4)
* `instance_engine_version` - The engine version of the Redis single instance (default: "5.0")
* `instance_password` - The password for the Redis instance (default: null)
* `command`                 - The Redis command to execute (default: "scan 0")
* `database`                - The database number to execute on (default: 0)

## Usage

* Copy this example script to your `main.tf`.

* Create a `terraform.tfvars` file and fill in the required variables:

  ```hcl
  vpc_name      = "your_vpc_name"
  subnet_name   = "your_subnet_name"
  instance_name = "your_instance_name"
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
* The `client_id` is obtained from the `huaweicloud_dcs_login_web_cli` resource, not configured manually
* This resource must be used together with `huaweicloud_dcs_login_web_cli` - login first, then execute commands
* The `command` parameter specifies the Redis command to execute, such as "scan 0", "get key", etc
* The `database` parameter specifies which database to execute the command on
* All parameters are non-updatable, changing any of them will trigger resource replacement
* Deleting this resource only removes it from Terraform state, the command execution is not reversed
* The resource ID is a generated UUID, not associated with any real DCS resource

## Requirements

| Name        | Version   |
|-------------|-----------|
| terraform   | >= 1.1.0  |
| huaweicloud | >= 1.92.0 |
