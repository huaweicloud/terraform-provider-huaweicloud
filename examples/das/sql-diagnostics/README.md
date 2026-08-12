# Configure DAS SQL diagnostics and optimization

This example provides best practice code for using Terraform to configure SQL diagnostics and optimization
in HuaweiCloud DAS service, including SQL limiting switch, batch SQL switch, and search path configuration.

## Prerequisites

* A HuaweiCloud account
* Terraform installed
* HuaweiCloud access key and secret key (AK/SK)
* An existing RDS instance with DAS database connection

## Variable Introduction

The following variables need to be configured:

### Authentication Variables

* `region_name` - The region where the DAS resources are located
* `access_key` - The access key of the IAM user
* `secret_key` - The secret key of the IAM user

### Resource Variables

#### Required Variables

* `sql_diagnostics_instance_id` - The ID of the database instance
* `sql_limiting_status` - The switch status of the SQL limiting
* `sql_diagnostics_datastore_type` - The database type
* `sql_diagnostics_engine_type` - The engine type of the instances
* `batch_sql_switch_on` - Whether to enable the SQL switch
* `batch_sql_switch_type` - The type of SQL switch to set
* `batch_sql_instance_ids` - The list of instance IDs
* `search_path_connection_id` - The ID of the database connection (DB user ID)
* `search_path_switch_on` - Whether to enable the search path switch

#### Optional Variables

* `batch_sql_retention_hours` - The retention hours of the SQL data (default: null)

## Usage

* Copy this example script to your `main.tf`.

* Create a `terraform.tfvars` file and fill in the required variables:

  ```hcl
  sql_diagnostics_engine_type    = "MySQL"
  batch_sql_switch_on            = true
  batch_sql_switch_type          = "DAS_QUERY"
  batch_sql_instance_ids         = ["your_instance_id_1", "your_instance_id_2"]
  search_path_connection_id      = "your_connection_id"
  search_path_switch_on          = true
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
* The SQL limiting switch and search path switch are one-time action resources
* The `search_path_connection_id` is the DB user ID, not the connection ID
* All resources will be created in the specified region

## Requirements

| Name | Version |
| ---- | ---- |
| terraform | >= 0.14.0 |
| huaweicloud | >= 1.93.0 |
