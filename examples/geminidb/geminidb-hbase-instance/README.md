# Create a GeminiDB HBase Instance

This example provides best practice code for using Terraform to create a GeminiDB HBase instance with backup
in HuaweiCloud GeminiDB service.

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
* `instance_name` - The GeminiDB HBase instance name
* `instance_backup_time_window` - The backup time window in HH:MM-HH:MM format
* `instance_backup_keep_days` - The number of days to retain backups

#### Optional Variables

* `vpc_cidr` - The CIDR block of the VPC (default: "192.168.0.0/16")
* `subnet_cidr` - The CIDR block of the subnet (default: "")
* `gateway_ip` - The gateway IP address of the subnet (default: "")
* `instance_password` - The password for the GeminiDB HBase instance (default: "")
* `availability_zone` - The availability zone (default: "")
* `instance_mode` - The instance mode (default: "Cluster")
* `instance_ssl_option` - Whether to enable SSL (default: "on")
* `vcpus` - The number of vCPUs of the flavor to query (default: 2)
* `instance_flavor_num` - The node quantity (default: 3)
* `instance_flavor_size` - The disk size in GB (default: 200)
* `instance_flavor_storage` - The disk type (default: "ULTRAHIGH")
* `instance_flavor_spec_code` - The resource specification code (default: "")
* `enterprise_project_id` - The enterprise project ID (default: "0")
* `maintenance_start_time` - The start time of the maintenance window (default: "02:00")
* `maintenance_end_time` - The end time of the maintenance window (default: "06:00")
* `charging_mode` - The charging mode (default: "prePaid")
* `period_unit` - The charging period unit (default: "month")
* `period` - The charging period (default: 1)
* `auto_renew` - Whether to enable auto-renew (default: "true")
* `tags` - The tags of the GeminiDB HBase instance (default: {"foo" = "bar", "key" = "value"})
* `backup_name` - The name of the GeminiDB backup (default: "tf_test_backup")
* `backup_description` - The description of the GeminiDB backup (default: "Created by Terraform")

## Usage

* Copy this example script to your `main.tf`.

* Create a `terraform.tfvars` file and fill in the required variables:

  ```hcl
  vpc_name                    = "your_vpc_name"
  subnet_name                 = "your_subnet_name"
  security_group_name         = "your_security_group_name"
  instance_name               = "your_geminidb_hbase_instance_name"
  instance_backup_time_window = "03:00-04:00"
  instance_backup_keep_days   = 14
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
* The GeminiDB HBase instance uses `datastore.type = "hbase"` with an empty version string and `storage_engine = "rocksDB"`
* The `mode` for a GeminiDB HBase instance must be `Cluster`
* The instance flavor is automatically queried from `huaweicloud_gaussdb_nosql_flavors` data source
* The `lifecycle.ignore_changes` block ignores attributes that may drift after creation: `password`, `flavor.0.storage`,
  `ssl_option`, `auto_renew`, `period`, and `period_unit`
* All resources will be created in the specified region

## Requirements

| Name | Version |
| ---- | ---- |
| terraform | >= 0.14.0 |
| huaweicloud | >= 1.80.2 |
| random | >= 3.0.0 |
