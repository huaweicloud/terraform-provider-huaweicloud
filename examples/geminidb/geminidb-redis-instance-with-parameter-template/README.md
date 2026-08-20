# Create a GeminiDB Redis Instance with Parameter Template

This example provides best practice code for using Terraform to create a GeminiDB Redis instance
with a custom parameter template on HuaweiCloud GeminiDB service.

## Prerequisites

* A HuaweiCloud account
* Terraform installed
* HuaweiCloud access key and secret key (AK/SK)

## Architecture

This example creates the following resources:

* **VPC & Subnet** - Network infrastructure for the GeminiDB Redis instance
* **Security Group** - Firewall rules to control access to the Redis instance
* **Parameter Template** - Custom Redis parameter configuration applied to the instance
* **GeminiDB Redis Instance** - The Redis database instance using the parameter template

## Variable Introduction

The following variables need to be configured:

### Authentication Variables

* `region_name` - The region where the GeminiDB Redis instance is located
* `access_key` - The access key of the IAM user
* `secret_key` - The secret key of the IAM user

### Required Variables

* `vpc_name` - The VPC name
* `subnet_name` - The subnet name
* `security_group_name` - The security group name
* `instance_name` - The GeminiDB Redis instance name
* `parameter_template_name` - The name of the parameter template

### Optional Variables

* `vpc_cidr` - The CIDR block of the VPC (default: "192.168.0.0/16")
* `subnet_cidr` - The CIDR block of the subnet (default: "", auto-computed from VPC CIDR)
* `gateway_ip` - The gateway IP address of the subnet (default: "", auto-computed from subnet CIDR)
* `availability_zone` - The availability zone (default: "", auto-selected from data source)
* `datastore_type` - The database type (default: "redis")
* `datastore_version` - The database version (default: "5.0")
* `datastore_storage_engine` - The storage engine (default: "rocksDB")
* `instance_db_port` - The database port (default: 8635)
* `instance_password` - The password for the instance (default: "", auto-generated)
* `instance_mode` - The instance mode (default: "Cluster")
* `instance_ssl_option` - The SSL option (default: "on")
* `instance_flavor_num` - The number of nodes in the instance (default: 3)
* `instance_flavor_size` - The storage size in GB per node (default: 16)
* `instance_flavor_storage` - The storage type (default: "ULTRAHIGH")
* `instance_flavor_spec_code` - The resource specification code (default: "")
* `instance_backup_time_window` - The backup time window in HH:MM-HH:MM format (default: "00:00-01:00")
* `instance_backup_keep_days` - The number of days to retain backups (default: 7)
* `tags` - The key/value pairs to associate with the instance (default: {})
* `parameter_template_description` - The description of the parameter template (default: "")
* `parameter_template_values` - The parameter key-value pairs (default: {})

## Usage

* Copy this example script to your `main.tf`.
* Create a `terraform.tfvars` file and fill in the required variables:

  ```hcl
  vpc_name                = "your_vpc_name"
  subnet_name             = "your_subnet_name"
  security_group_name     = "your_security_group_name"
  instance_name           = "your_geminidb_redis_instance_name"
  parameter_template_name = "your_parameter_template_name"
  ```

* To customize Redis parameters, set `parameter_template_values`:

  ```hcl
  parameter_template_values = {
    max_connections = "10000"
    timeout         = "300"
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
* The creation of the GeminiDB instance takes about 10 minutes
* The parameter template is created before the instance and applied via `configuration_id`
* If `instance_password` is empty, a random password will be generated
* If `availability_zone` is empty, the first available zone will be automatically selected
* If `subnet_cidr` is empty, it will be auto-computed from the VPC CIDR block
* If `gateway_ip` is empty, it will be auto-computed from the subnet CIDR block
* The `datastore` block in both the parameter template and the instance share
  the same `datastore_type` and `datastore_version` variables to ensure consistency

## Requirements

| Name | Version |
| ---- | ---- |
| terraform | >= 1.1.0 |
| huaweicloud | >= 1.92.0 |
