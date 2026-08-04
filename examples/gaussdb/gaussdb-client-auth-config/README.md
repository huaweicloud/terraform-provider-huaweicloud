# Configure Client Authentication for a GaussDB Instance

This example provides best practice code for using Terraform to configure client authentication
(HBA config) for a GaussDB instance in HuaweiCloud. It creates the complete infrastructure (VPC,
subnet, security group, GaussDB instance) and then sets up a client authentication rule to control
database access.

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

* `instance_id` - The ID of the GaussDB instance

#### Optional Variables

* `config_type` - The client connection type (default: `"host"`)
* `config_database` - The database name that the record matches (default: `"all"`)
* `config_user` - The database user that the record matches (default: `"root"`)
* `config_address` - The IP address range in CIDR format (default: `"10.10.0.0/16"`)
* `config_method` - The authentication method (default: `"md5"`)

## Usage

* Copy this example script to your `main.tf`.

* Create a `terraform.tfvars` file and fill in the required variables:

  ```hcl
  instance_id = "your_instance_id"
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

### Immutable Parameters

* The `instance_id`, `type`, `database`, `user`, and `address` cannot be modified after creation;
  changing them will create a new resource.
* Only the `method` field can be updated in-place after creation.

### Type Values

* `host` - Allows both SSL and non-SSL connections.
* `hostssl` - Only allows SSL connections.
* `hostnossl` - Only allows non-SSL connections.

### Database and User Values

* The `database` can be `all` or an existing database name.
* The `user` can be `all` or an existing username.

### Address Format

* The `address` must be in CIDR format (e.g., `10.10.0.0/16`).

### Method Values

* `md5` - MD5 password authentication.
* `sha256` - SHA-256 password authentication.
* `sm3` - SM3 password authentication (Chinese cryptography standard).
* `reject` - Rejects the connection unconditionally.
* `cert` - Client SSL certificate authentication.

### Import

* The client auth config can be imported using `instance_id`, `type`, `database`, `user`, and
  `address` separated by colons (`:`):

  ```bash
  $ terraform import huaweicloud_gaussdb_client_auth_config.test <instance_id>:<type>:<database>:<user>:<address>
  ```

### Timeouts

* The default timeout for create, update, and delete operations is 90 minutes each.

### General

* Make sure to keep your credentials secure and never commit them to version control.
* The resource ID is formatted as `<instance_id>:<type>:<database>:<user>:<address>`.

## Requirements

| Name        | Version   |
|-------------|-----------|
| terraform   | >= 0.14.0 |
| huaweicloud | >= 1.91.0 |
