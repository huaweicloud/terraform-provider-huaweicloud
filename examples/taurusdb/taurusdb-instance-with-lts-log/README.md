# Create a TaurusDB Instance with LTS Log

This example provides best practice code for using Terraform to create a TaurusDB instance and associate it with
a LTS (Log Tank Service) log group/stream for centralized log management within HuaweiCloud TaurusDB service.

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
* `instance_name` - The TaurusDB instance name
* `instance_backup_time_window` - The backup time window in HH:MM-HH:MM format
* `instance_backup_keep_days` - The number of days to retain backups
* `lts_group_name` - The name of the LTS log group
* `lts_stream_name` - The name of the LTS log stream

#### Optional Variables

* `vpc_cidr` - The CIDR block of the VPC (default: "192.168.0.0/16")
* `subnet_cidr` - The CIDR block of the subnet (default: "")
* `gateway_ip` - The gateway IP address of the subnet (default: "")
* `availability_zone_mode` - The availability zone mode (default: "multi")
* `master_availability_zone` - The master availability zone (default: "")
* `instance_db_port` - The database port (default: 3306)
* `instance_password` - The password for the TaurusDB instance (default: "")
* `instance_flavor_ref` - The flavor code of the TaurusDB instance (default: "")
* `instance_mode` - The instance mode (default: "Cluster")
* `read_replicas` - The number of read replicas (default: 2)
* `enterprise_project_id` - The enterprise project ID (default: "0")
* `volume_type` - The storage type of the instance (default: "DL6")
* `time_zone` - The time zone of the instance (default: "UTC+08:00")
* `ssl_option` - Whether to enable SSL (default: "true")
* `sql_filter_enabled` - Whether to enable SQL filter (default: true)
* `slow_log_show_original_switch` - Whether to enable slow log show original switch (default: true)
* `table_name_case_sensitivity` - Whether the kernel table name is case sensitive (default: true)
* `multi_tenant_switch` - Whether to enable multi-tenancy switch (default: "true")
* `maintain_begin` - The start time of the maintenance window (default: "02:00")
* `maintain_end` - The end time of the maintenance window (default: "06:00")
* `description` - The description of the TaurusDB instance (default: "")
* `seconds_level_monitoring_enabled` - Whether to enable seconds level monitoring (default: true)
* `seconds_level_monitoring_period` - The seconds level collection period (default: 5)
* `audit_log_enabled` - Whether to enable audit log (default: true)
* `audit_log_keep_days` - The number of days for storing audit logs (default: 7)
* `reserve_audit_logs` - Whether to reserve historical audit logs when SQL audit is disabled (default: "true")
* `tags` - The tags of the TaurusDB instance (default: {})
* `lts_group_ttl_in_days` - The expiration time (in days) of the LTS log group (default: 7)
* `lts_stream_ttl_in_days` - The expiration time (in days) of the LTS log stream (default: 7)
* `lts_stream_is_favorite` - Whether to favorite the log stream (default: false)
* `lts_log_type` - The type of the TaurusDB LTS log (default: "error_log")

## Usage

* Copy this example script to your `main.tf`.

* Create a `terraform.tfvars` file and fill in the required variables:

  ```hcl
  vpc_name                    = "your_vpc_name"
  subnet_name                 = "your_subnet_name"
  security_group_name         = "your_security_group_name"
  instance_name               = "your_taurusdb_instance_name"
  instance_backup_time_window = "02:00-03:00"
  instance_backup_keep_days   = 7
  lts_group_name              = "your_lts_group_name"
  lts_stream_name             = "your_lts_stream_name"
  tags = {
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
* The creation of the TaurusDB instance takes about 15-20 minutes
* This example creates the TaurusDB instance, VPC, subnet, security group, LTS log group, LTS log stream, and
  associates the instance with the LTS log stream for centralized log management
* The instance flavor and availability zones are automatically queried from `huaweicloud_taurusdb_flavors` data source
* The `lts_log_type` variable specifies which type of TaurusDB log to send to LTS. Valid values are `error_log` and
  `slow_log`. Change this value as needed
* All resources will be created in the specified region

## Requirements

| Name | Version   |
| ---- |-----------|
| terraform | >= 0.14.0 |
| huaweicloud | >= 1.95.0 |
| random | >= 3.0.0  |
