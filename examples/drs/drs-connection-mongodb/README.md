# Create a DRS Connection for MongoDB

This example provides best practice code for using Terraform to create a DRS connection
for a self-built MongoDB database with sharding within HuaweiCloud.

## Prerequisites

* A HuaweiCloud account
* Terraform installed
* HuaweiCloud access key and secret key (AK/SK)
* Self-built MongoDB databases accessible from the DRS instance

## Variable Introduction

The following variables need to be configured:

### Authentication Variables

* `region_name` - The region where resources will be created
* `access_key` - The access key of the IAM user
* `secret_key` - The secret key of the IAM user

### Resource Variables

#### Required Variables

* `connection_name` - The DRS connection name
* `db_password` - The password for the MongoDB database user
* `endpoint_ip` - The IP address and port of the primary MongoDB database (e.g. `192.168.0.1:8080`)
* `shard1_ip` - The IP address and port of the first MongoDB shard (e.g. `192.168.0.1:8000`)
* `shard2_ip` - The IP address and port of the second MongoDB shard (e.g. `192.168.0.2:8000`)

#### Optional Variables

* `description` - The description of the DRS connection (default: "")
* `db_user` - The database username (default: "mog")
* `db_name` - The database name (default: "root")
* `driver_name` - The driver name of the connection configuration (default: "mongodb")

## Usage

* Copy this example script to your `main.tf`.

* Create a `terraform.tfvars` file and fill in the required variables:

  ```hcl
  connection_name = "your_drs_mongodb_connection"
  db_password     = "Test@123456"
  endpoint_ip     = "192.168.0.1:8080"
  shard1_ip       = "192.168.0.1:8000"
  shard2_ip       = "192.168.0.2:8000"
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
* This example creates a DRS connection to a self-built MongoDB database with two shards
* The `endpoint.0.db_password`, `endpoint.0.source_sharding.*.db_password`,
  and `endpoint.0.source_sharding.*.endpoint_name` attributes are ignored in lifecycle
  changes since they are not returned by the API
* The MongoDB IP address format is `ip:port`, with multiple addresses separated by commas
* All resources will be created in the specified region

## Requirements

| Name | Version  |
| ---- |----------|
| terraform | >= 0.14.0 |
| huaweicloud | >= 1.93.0 |
