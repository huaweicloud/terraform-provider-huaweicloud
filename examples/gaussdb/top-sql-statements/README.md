# Query GaussDB Top SQL Statements

This example provides best practice code for using Terraform to query Top SQL statements of a
GaussDB instance in HuaweiCloud.

## Prerequisites

* A HuaweiCloud account with GaussDB permissions
* A running GaussDB instance (engine version >= 8.200 recommended)
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
* `node_ids` - The list of node IDs to query
* `start_time` - The start time for the query, in 13-digit UNIX timestamp format (milliseconds, UTC)
* `end_time` - The end time for the query, in 13-digit UNIX timestamp format (milliseconds, UTC)

#### Optional Variables

* `support_system` - Whether to display system users (default: `false`)
* `multi_queries` - The list of field aggregation query conditions (default: `[]`)
  + `name` - The query field name. Only `"query"` is supported
  + `condition` - The merge condition between multiple filter conditions. Valid
    values: `"and"`, `"or"`, `"AND"`, `"OR"`
  + `values` - The list of filter query values. Contains 1 to 5 strings
  + `is_fuzzy` - Whether to perform fuzzy query. Valid values: `"true"`, `"false"`. Default: `"true"`

## Outputs

* `top_sql_infos` - The list of Top SQL information, including SQL ID, username, SQL text, call frequency, CPU cost, IO
  cost, etc.

## Usage

* Copy this example script to your `main.tf`.

* Create a `terraform.tfvars` file and fill in the required variables:

  ```hcl
  instance_id = "xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx"
  node_ids    = ["xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx"]
  start_time  = 1750108800000
  end_time    = 1750195200000
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

* Make sure to keep your credentials secure and never commit them to version control
* The `start_time` and `end_time` must be in 13-digit UNIX timestamp format (milliseconds)
* The `db_name` filter only works for engine version 8.200 and above
* It is recommended to set a reasonable time range (e.g., the last 24 hours) to avoid excessive query results
* The `multi_queries` filter supports up to 5 entries, each with up to 5 values

## Requirements

| Name        | Version   |
|-------------|-----------|
| terraform   | >= 1.3.0  |
| huaweicloud | >= 1.82.3 |
