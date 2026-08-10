# Create a GeminiDB Instance with Manual Backup

This example provides best practice code for using Terraform to create a GeminiDB Cassandra instance with a manual
backup that includes database tables in HuaweiCloud GeminiDB service.

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
* `instance_name` - The GeminiDB Cassandra instance name
* `backup_name` - The name of the GeminiDB manual backup
* `instance_backup_time_window` - The backup time window in HH:MM-HH:MM format
* `instance_backup_keep_days` - The number of days to retain backups

#### Optional Variables - GeminiDB Instance

* `vpc_cidr` - The CIDR block of the VPC (default: "192.168.0.0/16")
* `subnet_cidr` - The CIDR block of the subnet (default: "")
* `gateway_ip` - The gateway IP address of the subnet (default: "")
* `instance_password` - The password for the GeminiDB Cassandra instance (default: "")
* `availability_zone` - The availability zone (default: "")
* `instance_mode` - The instance mode (default: "Cluster")
* `instance_ssl_option` - Whether to enable SSL (default: "on")
* `vcpus` - The number of vCPUs of the flavor to query (default: 2)
* `instance_flavor_num` - The node quantity (default: 3)
* `instance_flavor_size` - The disk size in GB (default: 16)
* `instance_flavor_storage` - The disk type (default: "ULTRAHIGH")
* `instance_flavor_spec_code` - The resource specification code (default: "")
* `enterprise_project_id` - The enterprise project ID (default: "0")
* `tags` - The tags of the GeminiDB Cassandra instance (default: {"foo" = "bar", "key" = "value"})

#### Optional Variables - GeminiDB Backup

* `backup_description` - The description of the GeminiDB manual backup (default: "test backup with database tables")
* `backup_database_name` - The database name for the backup database_tables block (default: "test_db")
* `backup_table_names` - The list of table names for the backup database_tables block (default: ["users"])

## Usage

* Copy this example script to your `main.tf`.

* Create a `terraform.tfvars` file and fill in the required variables:

  ```hcl
  vpc_name                    = "your_vpc_name"
  subnet_name                 = "your_subnet_name"
  security_group_name         = "your_security_group_name"
  instance_name               = "your_geminidb_cassandra_instance_name"
  backup_name                 = "your_geminidb_backup_name"
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
* The GeminiDB Cassandra instance uses `datastore.type = "redis"` with `version = "5.0"` and `storage_engine = "rocksDB"`
* The `mode` for a GeminiDB Cassandra instance is `Cluster`
* The instance flavor is automatically queried from `huaweicloud_gaussdb_nosql_flavors` data source
* The `lifecycle.ignore_changes` on the instance ignores `password`, `flavor.0.storage`, and `ssl_option`
* The `lifecycle.ignore_changes` on the backup ignores `database_tables` as it is not returned by the API
* All resources will be created in the specified region

## Requirements

| Name | Version   |
| ---- |-----------|
| terraform | >= 0.14.0 |
| huaweicloud | >= 1.92.0 |
| random | >= 3.0.0  |
