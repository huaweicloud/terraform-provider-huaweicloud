# Create a GeminiDB Cassandra Instance

This example provides best practice code for using Terraform to create a GeminiDB Cassandra instance in HuaweiCloud
GeminiDB service.

## Prerequisites

* A HuaweiCloud account
* Terraform installed
* HuaweiCloud access key and secret key (AK/SK)

## Variable Introduction

The following variables need to be configured:

### Authentication Variables

* `region_name` - The region where the GeminiDB Cassandra instance is located
* `access_key` - The access key of the IAM user
* `secret_key` - The secret key of the IAM user

### Resource Variables

#### Required Variables

* `vpc_name` - The VPC name
* `subnet_name` - The subnet name
* `security_group_name` - The security group name
* `instance_name` - The GeminiDB Cassandra instance name
* `instance_backup_time_window` - The backup time window in HH:MM-HH:MM format
* `instance_backup_keep_days` - The number of days to retain backups
* `backup_name` - The name for instance backups

#### Optional Variables

* `vpc_cidr` - The CIDR block of the VPC (default: "192.168.0.0/16")
* `availability_zone` - The availability zone to which the GeminiDB Cassandra instance belongs (default: "", auto-selected)
* `subnet_cidr` - The CIDR block of the subnet (default: "", auto-calculated from VPC CIDR)
* `gateway_ip` - The gateway IP address of the subnet (default: "", auto-calculated)
* `instance_datastore_version` - The Cassandra database engine version (default: "3.11")
* `instance_mode` - The instance mode (default: "Cluster"). Valid values: Cluster, Single
* `instance_flavor_num` - The number of nodes in the Cassandra cluster (default: 3)
* `instance_flavor_size` - The storage size in GB per node (default: 100)
* `instance_flavor_storage` - The storage type (default: "ULTRAHIGH"). Valid values: ULTRAHIGH, ESSD
* `instance_flavor_spec_code` - The resource specification code (default: "", auto-queried from flavors data source)
* `instance_db_port` - The Cassandra database port (default: 9042)
* `instance_password` - The password for the GeminiDB Cassandra instance (default: "", auto-generated)
* `instance_ssl_option` - The SSL option (default: "on"). Valid values: on, off
* `backup_description` - The description for instance backups (default: "Terraform created backup")
* `charging_mode` - The charging mode (default: "postPaid"). Valid values: postPaid, prePaid
* `period_unit` - The period unit for prePaid mode (default: "month"). Valid values: month, year
* `period` - The period for prePaid mode (default: 1)
* `auto_renew` - Whether to auto renew for prePaid mode (default: "false")
* `tags` - The key/value pairs to associate with the GeminiDB Cassandra instance (default: {})

## Usage

* Copy this example script to your `main.tf`.
* Create a `terraform.tfvars` file and fill in the required variables:

  ```hcl
  region_name                    = "cn-north-4"
  vpc_name                       = "your_vpc_name"
  subnet_name                    = "your_subnet_name"
  security_group_name            = "your_security_group_name"
  instance_name                  = "your_geminidb_cassandra_name"
  instance_backup_time_window    = "03:00-04:00"
  instance_backup_keep_days      = 14
  backup_name                    = "your_backup_name"
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
* The creation of the GeminiDB Cassandra instance takes about 10 minutes
* This example creates the GeminiDB Cassandra instance, VPC, subnet, security group, and backup
* All resources will be created in the specified region
* If `instance_flavor_spec_code` is empty, the flavor will be automatically queried from the
  `huaweicloud_gaussdb_nosql_flavors` data source
* If `instance_password` is empty, a random password will be generated
* The security group opens both the Cassandra native transport port (9042) and the Thrift port (9160)

## Requirements

| Name | Version |
| ---- | ---- |
| terraform | >= 0.14.0 |
| huaweicloud | >= 1.47.0 |
| random | >= 3.0.0 |
