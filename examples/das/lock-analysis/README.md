# Configure DAS lock analysis and transaction management

This example provides best practice code for using Terraform to configure lock analysis and transaction
management in HuaweiCloud DAS service, including full dead lock detection, history transaction switch,
and history transaction export.

## Prerequisites

* A HuaweiCloud account
* Terraform installed
* HuaweiCloud access key and secret key (AK/SK)
* An existing RDS MySQL instance with DAS database connection
* An existing OBS bucket for history transaction export

## Variable Introduction

The following variables need to be configured:

### Authentication Variables

* `region_name` - The region where the DAS resources are located
* `access_key` - The access key of the IAM user
* `secret_key` - The secret key of the IAM user

### Resource Variables

#### Required Variables

* `lock_analysis_instance_id` - The ID of the database instance
* `history_transaction_bucket_name` - The OBS bucket name for exporting history transactions

#### Optional Variables

* `full_dead_lock_switch_on` - Whether to enable the full dead lock switch (default: false)
* `full_dead_lock_retention_hours` - The retention hours of the full dead lock data (default: null)
* `history_transaction_status` - The switch status of the history transaction (default: "Enabled")
* `lock_analysis_datastore_type` - The database type (default: "MySQL")
* `history_transaction_start_time` - The start time of the history transactions to export, in RFC3339 format (default: "2000-06-01T00:00:00+08:00")
* `history_transaction_end_time` - The end time of the history transactions to export, in RFC3339 format (default: "2099-06-02T00:00:00+08:00")
* `history_transaction_file_path` - The OBS file directory for the export task (default: null)
* `history_transaction_time_zone` - The time zone for the export task (default: "UTC+8")
* `history_transaction_order_field` - The sort field for the export task (default: "collectTime")
* `history_transaction_order_by` - The sort order for the export task (default: "asc")
* `history_transaction_last_sec_min` - The minimum duration for the export task (default: 0)
* `history_transaction_last_sec_max` - The maximum duration for the export task (default: 100)

## Usage

* Copy this example script to your `main.tf`.

* Create a `terraform.tfvars` file and fill in the required variables:

  ```hcl
  lock_analysis_instance_id       = "your_rds_instance_id"
  full_dead_lock_switch_on        = true
  full_dead_lock_retention_hours  = 24
  history_transaction_status      = "Enabled"
  lock_analysis_datastore_type    = "MySQL"
  history_transaction_bucket_name = "your_obs_bucket"
  history_transaction_start_time  = "2026-08-01T00:00:00Z"
  history_transaction_end_time    = "2026-08-10T23:59:59Z"
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

* Make sure to keep your credentials secure and never commit them to version control
* The `huaweicloud_das_full_dead_lock_switch` only supports MySQL engine type
* The `huaweicloud_das_history_transaction_switch` and `huaweicloud_das_history_transaction_export_task`
  are one-time action resources; deleting them from the configuration will not clear the corresponding
  request records on the server side
* The history transaction export task requires an existing OBS bucket
* All resources will be created in the specified region

## Requirements

| Name | Version |
| ---- | ---- |
| terraform | >= 0.14.0 |
| huaweicloud | >= 1.93.0 |
