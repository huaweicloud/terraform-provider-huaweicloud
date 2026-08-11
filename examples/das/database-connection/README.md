# Create a DAS database connection with user and shared connection management

This example provides best practice code for using Terraform to manage DAS database connections in HuaweiCloud,
including database instance connection, database user creation, and shared connection configuration.

## Prerequisites

* A HuaweiCloud account
* Terraform installed
* HuaweiCloud access key and secret key (AK/SK)
* An existing RDS instance

## Variable Introduction

The following variables need to be configured:

### Authentication Variables

* `region_name` - The region where the DAS resources are located
* `access_key` - The access key of the IAM user
* `secret_key` - The secret key of the IAM user

### Resource Variables

#### Required Variables

* `connection_instance_id` - The ID of the RDS instance to connect
* `connection_engine_type` - The engine type of the database instance
* `connection_network_type` - The network type of the database instance connection
* `connection_username` - The username for the database instance connection
* `connection_password` - The password for the database instance connection
* `db_user_name` - The name of the database user
* `db_user_password` - The password of the database user
* `shared_user_id` - The IAM user ID to share the connection with
* `shared_user_name` - The IAM user name to share the connection with

#### Optional Variables

* `connection_is_save_password` - Whether to save the password for the database instance
  connection (default: true)
* `connection_port` - The port of the database instance connection (default: null)
* `connection_database_name` - The database name of the database instance connection
  (default: null)
* `connection_sql_record_flag` - Whether SQL recording is enabled (default: null)
* `connection_description` - The description of the database instance connection
  (default: null)
* `connection_node_ids` - The unique identifiers of the instance nodes (default: null)
* `shared_expired_at` - The expiration time of the shared connection, in RFC3339 format
  (default: null)

## Usage

* Copy this example script to your `main.tf`.

* Create a `terraform.tfvars` file and fill in the required variables:

  ```hcl
  connection_instance_id  = "your_rds_instance_id"
  connection_engine_type  = "MySQL"
  connection_network_type = "EIP"
  connection_username     = "your_username"
  connection_password     = "your_password"
  db_user_name            = "your_db_user"
  db_user_password        = "your_db_password"
  shared_user_id          = "your_iam_user_id"
  shared_user_name        = "your_iam_user_name"
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
* The `huaweicloud_das_database_instance_connection` and `huaweicloud_das_database_user`
  resources must be created before the `huaweicloud_das_shared_connection` resource
* All resources will be created in the specified region

## Requirements

| Name | Version |
| ---- || ---- |
| terraform | >= 1.1.0 |
| huaweicloud | >= 1.93.0 |
