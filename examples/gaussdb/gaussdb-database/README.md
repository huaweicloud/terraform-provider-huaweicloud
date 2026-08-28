# GaussDB database

This example provides best practice code for using Terraform to manage a database within a GaussDB OpenGauss instance in
HuaweiCloud.

## Prerequisites

* A HuaweiCloud account
* Terraform installed
* HuaweiCloud access key and secret key (AK/SK)
* An existing GaussDB OpenGauss instance

## Required Variables

### Authentication Variables

* `region_name` - The region where the GaussDB database resource is located
* `access_key`  - The access key of the IAM user
* `secret_key`  - The secret key of the IAM user

### Resource Variables

#### Required Variables

* `vpc_name` - The VPC name
* `subnet_name` - The subnet name
* `security_group_name` - The security group name
* `instance_name` - The GaussDB instance name
* `database_name` - The name of the database (1-63 characters, cannot start with `pg` or a digit, cannot be a template
  database name)

#### Optional Variables

* `vpc_cidr` - The CIDR block of the VPC (default: `"192.168.0.0/16"`)
* `enterprise_project_id` - The ID of the enterprise project (default: `""`)
* `subnet_cidr` - The CIDR block of the subnet (default: `""`, auto-calculated from VPC CIDR)
* `gateway_ip` - The gateway IP address of the subnet (default: `""`, auto-calculated)
* `security_group_rule_ports` - The security group ingress rule ports
  (default: `2379-2380,5000-5001,5432-5532,6000,6500,12016,20050`)
* `instance_password` - The password for the GaussDB instance (default: `""`, auto-generated)
* `instance_volume_type` - The storage volume type (default: `"ULTRAHIGH"`)
* `instance_volume_size` - The storage volume size in GB (default: `40`)
* `character_set` - The database character set, defaults to `UTF8`
* `owner` - The database user, defaults to `root` (must not be a system user: `rdsAdmin`, `rdsMetric`,
  `rdsBackup`, `rdsRepl`)
* `template` - The database template name, can be `template0`
* `lc_collate` - The database collation, defaults to `C`
* `lc_ctype` - The database classification, defaults to `C`

## Usage

* Copy this example script to your `main.tf`.

* Create a `terraform.tfvars` file and fill in the required variables:

  ```hcl
  vpc_name            = "your_vpc_name"
  subnet_name         = "your_subnet_name"
  security_group_name = "your_security_group_name"
  instance_name       = "your_gaussdb_instance_name"
  database_name       = "test_db_name"
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

## Note

* This resource manages a database within a GaussDB OpenGauss instance using the opengauss service endpoint
* The create operation sends a POST request to `/v3/{project_id}/instances/{instance_id}/database` and waits for the
  instance to reach ACTIVE state (default timeout: 90 minutes)
* The delete operation sends a DELETE request and waits for the instance to reach ACTIVE state (default timeout: 90
  minutes)
* The resource ID is formatted as `<instance_id>/<name>`
* All input parameters are non-updatable (ForceNew); changing any parameter will trigger resource replacement
* There is no update operation
* The read operation paginates through the database list (100 per page) to find the database by name
* The `template` and `lc_ctype` parameters are not returned by the API and will not be in the state after import
* The `lc_collate` parameter maps to `collate_set` in the API response
* The computed attributes `size` and `compatibility_type` are populated from the API response
* The resource supports import using `<instance_id>/<name>` format
* When importing, `template` and `lc_ctype` should be added to `lifecycle.ignore_changes` as they are not in the API
  response

## Requirements

| Name        | Version   |
|-------------|-----------|
| terraform   | >= 0.14.0 |
| huaweicloud | >= 1.90.0 |
