# Create a GeminiDB DynamoDB-compatible instance

This example provides best practice code for using Terraform to create a GeminiDB DynamoDB-compatible instance in HuaweiCloud.

## Prerequisites

* A HuaweiCloud account
* Terraform installed (>= 0.14.0)
* HuaweiCloud access key and secret key (AK/SK)

## Architecture

The example creates the following resources:

* A VPC and subnet for the GeminiDB DynamoDB instance
* A security group for the GeminiDB DynamoDB instance
* A GeminiDB DynamoDB-compatible instance (`huaweicloud_geminidb_instance` with `datastore.type = "dynamodb"`)
* A GeminiDB manual backup for the instance

## Variable Introduction

The following variables need to be configured:

### Authentication Variables

* `access_key` - HuaweiCloud access key
* `secret_key` - HuaweiCloud secret key
* `region_name` - The region where resources will be created

### Resource Variables

#### Required Variables

* `vpc_name` - The name of the VPC
* `vpc_cidr` - The CIDR block of the VPC
* `subnet_name` - The name of the VPC subnet
* `security_group_name` - The name of the security group
* `instance_name` - The name of the GeminiDB DynamoDB instance

#### Optional Variables

* `instance_password` - The password of the instance, generated randomly if left empty (default: "")
* `availability_zone` - The availability zone, uses the first AZ from data source if empty (default: "")
* `instance_mode` - The instance type (default: "Cluster")
* `instance_ssl_option` - Whether SSL is enabled (default: "on")
* `vcpus` - The number of vCPUs of the flavor to query (default: 2)
* `instance_flavor_num` - The node quantity (default: 3)
* `instance_flavor_size` - The disk size in GB (default: 200)
* `instance_flavor_storage` - The disk type (default: "ULTRAHIGH")
* `instance_flavor_spec_code` - The resource specification code, auto-fetched from data source if empty (default: "")
* `instance_backup_time_window` - The backup time window (default: "03:00-04:00")
* `instance_backup_keep_days` - The number of days to retain backup files (default: 14)
* `enterprise_project_id` - The enterprise project ID (default: "0")
* `maintenance_start_time` - The maintenance start time in UTC (default: "02:00")
* `maintenance_end_time` - The maintenance end time in UTC, 4 hours after start time (default: "06:00")
* `charging_mode` - The charging mode, postPaid or prePaid (default: "prePaid")
* `period_unit` - The charging period unit, month or year (default: "month")
* `period` - The charging period (default: 1)
* `auto_renew` - Whether auto-renew is enabled (default: "true")
* `tags` - The key/value pairs to associate with the instance (default: { foo = "bar", key = "value" })
* `subnet_cidr` - The CIDR block of the VPC subnet (default: "", derived from VPC CIDR)
* `subnet_gateway_ip` - The gateway IP of the VPC subnet (default: "", derived from subnet CIDR)
* `backup_name` - The name of the GeminiDB backup (default: "tf_test_backup")
* `backup_description` - The description of the GeminiDB backup (default: "Created by Terraform")

## Usage

* Copy this example script to your `main.tf`.

* Create a `terraform.tfvars` file and fill in the required variables:

  ```hcl
  access_key  = "your_access_key"
  secret_key  = "your_secret_key"
  region_name = "your_region"
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

* The GeminiDB DynamoDB-compatible instance uses `datastore.type = "dynamodb"` with an empty `version` string and
  `storage_engine = "rocksDB"`.
* The `mode` for a GeminiDB DynamoDB-compatible instance must be `Cluster`.
* The `port` parameter is not supported for GeminiDB DynamoDB-compatible instances (only available for GeminiDB Redis instances).
* The flavor `spec_code` is queried via the `huaweicloud_gaussdb_nosql_flavors` data source using `engine = "cassandra"`
  since the DynamoDB-compatible instance shares the underlying infrastructure.
* The `lifecycle.ignore_changes` block follows the resource MD documentation's Import section, ignoring attributes that
  may drift after creation: `password`, `flavor.0.storage`, `ssl_option`, `auto_renew`, `period`, and `period_unit`.
* Make sure to keep your credentials secure and never commit them to version control.
* All resources will be created in the specified region.

## Requirements

| Name | Version |
| ---- | ---- |
| terraform | >= 0.14.0 |
| huaweicloud | >= 1.80.2 |
| random | >= 3.0.0 |
