# Configure LTS Logging for a DRS Job

This example provides best practice code for using Terraform to configure LTS logging for
a DRS migration job within HuaweiCloud. It creates the full dependency chain including
VPC, RDS instances, DRS job, and LTS resources.

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

### Network Variables

* `vpc_name` - The VPC name
* `vpc_cidr` - The CIDR block of the VPC (default: "192.168.0.0/16")
* `subnet_name` - The subnet name
* `subnet_cidr` - The CIDR block of the subnet (default: auto-derived from VPC)
* `gateway_ip` - The gateway IP address of the subnet (default: auto-derived from subnet)
* `security_group_name` - The security group name

### RDS Variables

* `source_rds_name` - The name of the source RDS instance
* `dest_rds_name` - The name of the destination RDS instance
* `db_password` - The password for the RDS root user and DRS database connections
* `source_rds_fixed_ip` - The fixed IP address of the source RDS instance
* `dest_rds_fixed_ip` - The fixed IP address of the destination RDS instance
* `rds_flavor` - The flavor of the RDS instances (default: "rds.mysql.x1.large.2.ha")
* `rds_db_type` - The database type for querying RDS flavors (default: "MySQL")
* `rds_db_version` - The database version for querying RDS flavors (default: "5.7")
* `rds_instance_mode` - The instance mode for querying RDS flavors (default: "ha")

### DRS Variables

* `job_name` - The DRS job name
* `description` - The description of the DRS job (default: "")

### LTS Variables

* `lts_group_name` - The name of the LTS log group
* `lts_stream_name` - The name of the LTS log stream
* `lts_ttl_in_days` - The log retention period in days (default: 30)

## Usage

* Copy this example script to your `main.tf`.

* Create a `terraform.tfvars` file and fill in the required variables:

  ```hcl
  vpc_name            = "drs-vpc"
  subnet_name         = "drs-subnet"
  security_group_name = "drs-secgroup"
  source_rds_name     = "drs-source-rds"
  dest_rds_name       = "drs-dest-rds"
  source_rds_fixed_ip = "192.168.0.10"
  dest_rds_fixed_ip   = "192.168.0.11"
  db_password         = "your-strong-password"
  job_name            = "drs-migration-job"
  lts_group_name      = "drs-lts-group"
  lts_stream_name     = "drs-lts-stream"
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
* The `job_id` of the DRS LTS config is referenced from the DRS job resource
* The `log_group_id` and `log_stream_id` are referenced from the LTS resources
* Deleting the DRS LTS config resource will disable the LTS switch for the DRS job
* The DRS job uses `force_destroy = true` to allow deletion even if the job is running
* The DRS LTS config can be imported using its `id`

## Requirements

| Name | Version  |
| ---- |----------|
| terraform | >= 1.1.0 |
| huaweicloud | >= 1.92.0 |
