# Create a DRS Migration Job

This example provides best practice code for using Terraform to create a DRS migration job
from one MySQL RDS instance to another within HuaweiCloud.

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
* `source_rds_name` - The name of the source RDS instance
* `dest_rds_name` - The name of the destination RDS instance
* `source_rds_fixed_ip` - The fixed IP address of the source RDS instance
* `dest_rds_fixed_ip` - The fixed IP address of the destination RDS instance
* `db_password` - The password for the RDS root user and DRS database connections
* `job_name` - The DRS job name

#### Optional Variables

* `vpc_cidr` - The CIDR block of the VPC (default: "192.168.0.0/16")
* `subnet_cidr` - The CIDR block of the subnet (default: "")
* `gateway_ip` - The gateway IP address of the subnet (default: "")
* `rds_flavor` - The flavor of the RDS instances. If not specified,
  the system will automatically query available flavors (default: "")
* `rds_db_type` - The database type for querying RDS flavors (default: "MySQL")
* `rds_db_version` - The database version for querying RDS flavors (default: "5.7")
* `rds_instance_mode` - The instance mode for querying RDS flavors (default: "ha")
* `description` - The description of the DRS job (default: "")
* `tags` - The tags of the DRS job (default: {})

## Usage

* Copy this example script to your `main.tf`.

* Create a `terraform.tfvars` file and fill in the required variables:

  ```hcl
  vpc_name            = "your_vpc_name"
  subnet_name         = "your_subnet_name"
  security_group_name = "your_security_group_name"
  source_rds_name     = "your_source_rds_name"
  dest_rds_name       = "your_dest_rds_name"
  source_rds_fixed_ip = "192.168.0.58"
  dest_rds_fixed_ip   = "192.168.0.59"
  db_password         = "TestDrs@123"
  job_name            = "your_drs_job_name"
  tags                = {
    foo = "bar"
  }
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
* The creation of the DRS job takes about 10-15 minutes
* This example creates the DRS migration job, VPC, subnet, security group, source and destination RDS instances

* The `source_db.0.password`, `destination_db.0.password`, `force_destroy`, and `action` attributes are ignored in
  lifecycle changes since they are not returned by the API or may change during runtime
* All resources will be created in the specified region

## Requirements

| Name | Version  |
| ---- |----------|
| terraform | >= 1.1.0 |
| huaweicloud | >= 1.68.0 |
