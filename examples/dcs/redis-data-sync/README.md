# Create a DCS Redis Data Synchronization

This example provides best practice code for using Terraform to synchronize data between two DCS Redis instances in
HuaweiCloud.

## Prerequisites

* A HuaweiCloud account
* Terraform installed
* HuaweiCloud access key and secret key (AK/SK)
* Two DCS Redis instances (source and target) or create them using this example

## Variables Introduction

The following variables need to be configured:

### Authentication Variables

* `region_name` - The region where the DCS instances are located
* `access_key` - The access key of the IAM user
* `secret_key` - The secret key of the IAM user

### Required Variables

* `instance_name` - The base name of the DCS Redis instances (will be suffixed with -0 and -1 for source and target)
* `instance_password` - The password of the DCS instances

### Optional Variables

* `vpc_name` - The name of the VPC (default: "dcs-sync-vpc")
* `vpc_cidr` - The CIDR block of the VPC (default: "192.168.0.0/16")
* `subnet_name` - The name of the subnet (default: "dcs-sync-subnet")
* `subnet_cidr` - The CIDR block of the subnet (default: "", auto-calculated if empty)
* `subnet_gateway_ip` - The gateway IP of the subnet (default: "", auto-calculated if empty)
* `security_group_name` - The name of the security group (default: "dcs-sync-sg")
* `instance_cache_mode` - The cache mode of the DCS instances (default: "ha")
* `instance_capacity` - The capacity of the DCS instances in GB (default: `4`)
* `instance_engine_version` - The engine version of the DCS instances (default: "5.0")
* `full_migration_task_name` - The name of the full migration task (default: "full-migration-task")
* `full_migration_task_description` - The description of the full migration task
  (default: "Full data migration from source to target DCS instance")
* `full_migration_resume_mode` - The reconnection mode for full migration (default: "auto")
* `full_migration_bandwidth_limit_mb` - The bandwidth limit for full migration in MB/s (default: "", no limit if empty)
* `enable_incremental_migration` - Whether to enable incremental migration after full migration (default: true)
* `incremental_migration_task_name` - The name of the incremental migration task (default: "incremental-migration-task")
* `incremental_migration_task_description` - The description of the incremental migration task
  (default: "Incremental data migration from source to target DCS instance")
* `incremental_migration_resume_mode` - The reconnection mode for incremental migration (default: "auto")
* `incremental_migration_bandwidth_limit_mb` - The bandwidth limit for incremental migration in MB/s
  (default: "", no limit if empty)

## Usage

* Copy this example script to your `main.tf`.

* Create a `terraform.tfvars` file and fill in the required variables:

  ```hcl
  # DCS Instance Configuration
  instance_name     = "redis-instance"
  instance_password = "YourPassword@123"

  # Full Migration Task Configuration
  full_migration_bandwidth_limit_mb = "100"

  # Incremental Migration Task Configuration
  enable_incremental_migration             = true
  incremental_migration_bandwidth_limit_mb = "50"
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

## Notes

* Make sure to keep your credentials secure and never commit them to version control
* The source and target instances must be in the same VPC or have network connectivity
* Both instances share the same password for migration to work properly
* Full migration may take a long time depending on the amount of data
* Incremental migration will run continuously until manually stopped
* Bandwidth limits can help ensure smooth service running during migration
* The migration tasks have a default timeout of 30 minutes for create/update operations
* All resources will be created in the specified region
* Instance names must be unique within your HuaweiCloud account
  (the actual instance names will be `${instance_name}-0` and `${instance_name}-1`)
* Both instances share the same configuration for simplicity, but you can modify the code
  to use different configurations if needed
* When using auto resume mode, be aware that it may trigger full synchronization if incremental synchronization
  becomes unavailable, which requires more bandwidth
* When destroying resources, the incremental migration task will be stopped first (if it's running), then
  both migration tasks will be deleted, followed by the DCS instances and network resources.

## Requirements

| Name | Version |
| ---- | ---- |
| terraform | >= 0.14.0 |
| huaweicloud | >= 1.77.1 |
