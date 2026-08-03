# Set Database Role Permission for a GaussDB Instance

This example provides best practice code for using Terraform to set database role permissions for
a GaussDB instance in HuaweiCloud. It creates the complete infrastructure (VPC, subnet, security
group, GaussDB instance, database, database role, schema) and then assigns read-only or read-write
permissions to the database role.

## Prerequisites

* A HuaweiCloud account with GaussDB permissions
* Terraform installed
* HuaweiCloud access key and secret key (AK/SK)

## Variable Introduction

The following variables need to be configured:

### Authentication Variables

* `region_name` - The region where the GaussDB instance is located
* `access_key` - The access key of the IAM user
* `secret_key` - The secret key of the IAM user

### Resource Variables

#### Required Variables

* `vpc_name` - The VPC name
* `subnet_name` - The subnet name
* `security_group_name` - The security group name
* `instance_name` - The GaussDB instance name
* `database_role_name` - The name of the database role
* `db_name` - The database name
* `db_owner` - The database owner
* `schema_name` - The schema name

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
* `database_role_password` - The password of the database role (default: `""`, auto-generated)
* `is_login_only` - Whether the database role supports login only (default: `"false"`)
* `permission_readonly` - Whether the database role permission is read-only (default: `"true"`)

## Usage

* Copy this example script to your `main.tf`.

* Create a `terraform.tfvars` file and fill in the required variables:

  ```hcl
  vpc_name               = "your_vpc_name"
  subnet_name            = "your_subnet_name"
  security_group_name    = "your_security_group_name"
  instance_name          = "your_gaussdb_instance_name"
  database_role_name     = "your_role_name"
  database_role_password = "your_role_password"
  db_name                = "your_db_name"
  db_owner               = "your_db_owner"
  schema_name            = "your_schema_name"
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

### Resource Constraints

* The `instance_id` and `db_name` cannot be modified after creation; changing them will create a new
  resource.
* The `user` block supports only one element (`MaxItems: 1`).
* The `name` field in the `user` block cannot be modified after creation.

### Database Name Constraints

* Template databases (`postgres`, `template0`, `template1`) cannot be used.
* The database must already exist in the GaussDB instance.

### Role Name Constraints

* The role name must be 1 to 63 characters, including letters, digits, and underscores, and cannot
  start with `pg` or a digit.
* The role name cannot be the same as system user names (`rdsAdmin`, `rdsMetric`, `rdsBackup`,
  `rdsRepl`, and `root`).
* The role must already exist in the GaussDB instance.

### Schema Name Constraints

* The schemas `public` and `information_schema` cannot be used.
* The schema must already exist in the GaussDB instance.

### Permission Values

* `true` - Read-only permission.
* `false` - Read and write permission.

### Default Privilege Grantee

* The `default_privilege_grantee` field specifies a database user/role whose permissions are granted
  to the role specified by `name`.
* Leave this field empty if default privilege granting is not needed.

### Provider Behavior

* **Read is a no-op**: The resource does not query the current state from the API. Terraform will
  not detect configuration drift for this resource.
* **Delete is a no-op**: Deleting this resource only removes it from the Terraform state. The
  database role permissions are not revoked from the GaussDB instance. A warning message will be
  displayed during deletion.
* **No import support**: This resource cannot be imported into Terraform state.

### General

* Make sure to keep your credentials secure and never commit them to version control.
* The `depends_on` in the `huaweicloud_gaussdb_instance_role_permission` resource ensures the
  database, database role, and schema are created before the permission is set.

## Requirements

| Name        | Version   |
|-------------|-----------|
| terraform   | >= 0.14.0 |
| huaweicloud | >= 1.95.0 |
