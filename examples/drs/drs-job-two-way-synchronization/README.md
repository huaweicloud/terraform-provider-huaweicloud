# Create DRS Jobs for Two-Way Synchronization

This example provides best practice code for using Terraform to create two DRS synchronization jobs
between two MySQL RDS instances within HuaweiCloud, which forms a two-way (bidirectional)
synchronization scenario.

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
* `rds1_name` - The name of the first RDS instance
* `rds2_name` - The name of the second RDS instance
* `db_password` - The password for the RDS root user and DRS database connections
* `forward_job_name` - The name of the forward DRS job
* `reverse_job_name` - The name of the reverse DRS job

#### Optional Variables

* `vpc_cidr` - The CIDR block of the VPC (default: "192.168.0.0/16")
* `subnet_cidr` - The CIDR block of the subnet (default: "")
* `gateway_ip` - The gateway IP address of the subnet (default: "")
* `rds_flavor` - The flavor of the RDS instances (default: "rds.mysql.x1.large.2.ha")
* `rds1_fixed_ip` - The fixed IP address of the first RDS instance (default: "192.168.0.50")
* `rds2_fixed_ip` - The fixed IP address of the second RDS instance (default: "192.168.0.51")
* `db_user` - The database user name used by the RDS instances and DRS jobs (default: "root")
* `db_name` - The name of the database to be synchronized in both directions (default: "test")
* `forward_job_description` - The description of the forward DRS job (default: "")
* `reverse_job_description` - The description of the reverse DRS job (default: "")
* `limit_speed` - The speed limit of the DRS jobs, it is not limited by default (default: null)

## Usage

* Copy this example script to your `main.tf`.

* Create a `terraform.tfvars` file and fill in the required variables:

  ```hcl
  vpc_name            = "your_vpc_name"
  subnet_name         = "your_subnet_name"
  security_group_name = "your_security_group_name"
  rds1_name           = "your_rds_name_1"
  rds2_name           = "your_rds_name_2"
  rds1_fixed_ip       = "192.168.0.50"
  rds2_fixed_ip       = "192.168.0.51"
  db_password         = "TestDrs@123"
  db_name             = "test"
  forward_job_name    = "your_forward_drs_job_name"
  reverse_job_name    = "your_reverse_drs_job_name"
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
* This example creates the VPC, subnet, security group, two RDS MySQL instances, the database to be
  synchronized on both instances, and the DRS jobs with their object selection and start actions
* The two RDS instances are created with the fixed IP addresses `rds1_fixed_ip` and `rds2_fixed_ip`,
  both of them must be in the CIDR block of the subnet
* The `source_db.0.password`, `destination_db.0.password`, and `force_destroy` ttributes are ignored in
  lifecycle changes since they are not returned by the API or may change during runtime
* The `huaweicloud_drs_batch_select_objects` and `huaweicloud_drs_start_job` resources are one-time
  action resources, deleting them will not clear the corresponding request record or stop the jobs,
  but will only remove the resource information from the tf state file
* All resources will be created in the specified region

## Requirements

| Name | Version   |
| ---- |-----------|
| terraform | >= 1.1.0  |
| huaweicloud | >= 1.98.0 |
