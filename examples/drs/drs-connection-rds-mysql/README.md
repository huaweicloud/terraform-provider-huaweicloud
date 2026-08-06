# Create a DRS Connection for RDS MySQL

This example provides best practice code for using Terraform to create a DRS connection
for a MySQL RDS instance within HuaweiCloud.

## Prerequisites

* A HuaweiCloud account
* Terraform installed
* HuaweiCloud access key and secret key (AK/SK)

## Variable Introduction

The following variables need to be configured:

### Authentication Variables

* `region_name` - The region where resources will be created
* `access_key` - The access key of the IAM user
* `secret_key` - The secret key of the IAM user

### Resource Variables

#### Required Variables

* `vpc_name` - The VPC name
* `subnet_name` - The subnet name
* `security_group_name` - The security group name
* `rds_name` - The name of the RDS instance
* `db_password` - The password for the RDS root user and DRS connection
* `connection_name` - The DRS connection name

#### Optional Variables

* `vpc_cidr` - The CIDR block of the VPC (default: "192.168.0.0/16")
* `subnet_cidr` - The CIDR block of the subnet (default: "")
* `gateway_ip` - The gateway IP address of the subnet (default: "")
* `rds_flavor` - The flavor of the RDS instance. If not specified, it will be queried from data source (default: "")
* `rds_fixed_ip` - The fixed IP address of the RDS instance (default: "192.168.0.100")
* `rds_db_type` - The database type for querying RDS flavors (default: "MySQL")
* `rds_db_version` - The database version for querying RDS flavors (default: "5.7")
* `rds_instance_mode` - The instance mode for querying RDS flavors (default: "single")
* `description` - The description of the DRS connection (default: "")
* `db_user` - The database username (default: "root")
* `db_port` - The database port (default: "3306")
* `driver_name` - The driver name of the connection configuration (default: "mysql")

## Usage

* Copy this example script to your `main.tf`.

* Create a `terraform.tfvars` file and fill in the required variables:

  ```hcl
  vpc_name            = "your_vpc_name"
  subnet_name         = "your_subnet_name"
  security_group_name = "your_security_group_name"
  rds_name            = "your_rds_name"
  db_password         = "TestDrs@123"
  connection_name     = "your_drs_connection_name"
  description         = "DRS connection for RDS MySQL"
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
* This example creates the DRS connection, VPC, subnet, security group, and RDS instance
* The `endpoint.0.db_password` attribute is ignored in lifecycle changes since it is not returned by the API
* All resources will be created in the specified region

## Requirements

| Name | Version  |
| ---- |----------|
| terraform | >= 0.14.0 |
| huaweicloud | >= 1.93.0 |
