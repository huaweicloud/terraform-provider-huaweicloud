# Configure DAS log analysis

This example provides best practice code for using Terraform to configure log analysis in HuaweiCloud DAS service,
including slow log export and binlog parse.

## Prerequisites

* A HuaweiCloud account
* Terraform installed
* HuaweiCloud access key and secret key (AK/SK)
* An existing RDS MySQL instance with DAS database connection
* An existing OBS bucket for slow log and binlog export

## Variable Introduction

The following variables need to be configured:

### Authentication Variables

* `region_name` - The region where the DAS resources are located
* `access_key` - The access key of the IAM user
* `secret_key` - The secret key of the IAM user

### Resource Variables

#### Required Variables

* `log_analysis_instance_id` - The ID of the database instance
* `slow_log_bucket_name` - The OBS bucket name for exporting slow logs
* `slow_log_start_time` - The start time of the slow logs to export, in RFC3339 format
* `slow_log_end_time` - The end time of the slow logs to export, in RFC3339 format
* `binlog_binlog_type` - The binlog type
* `binlog_file_name` - The binlog file name
* `binlog_export_bucket_name` - The OBS bucket name for exporting binlog parse results

#### Optional Variables

* `slow_log_file_path` - The OBS file directory for the export task (default: null)
* `slow_log_export_type` - The export type for the slow log export task (default: null)
* `slow_log_sort_field` - The sort field for the slow log export task (default: null)
* `slow_log_sort_asc` - Whether to sort in ascending order (default: null)
* `slow_log_time_zone` - The time zone for the slow log export task (default: null)
* `binlog_backup_id` - The backup ID (default: null)
* `binlog_filter_db_names` - The list of database names to filter (default: [])
* `binlog_filter_tb_names` - The list of table names to filter (default: [])
* `binlog_filter_start_time` - The start time of the export range, in RFC3339 format (default: null)
* `binlog_filter_end_time` - The end time of the export range, in RFC3339 format (default: null)
* `binlog_filter_types` - The list of SQL types to filter (default: [])
* `binlog_filter_parse_double_insert` - Whether to export UPDATE statements as two INSERT statements (default: null)

## Usage

* Copy this example script to your `main.tf`.

* Create a `terraform.tfvars` file and fill in the required variables:

  ```hcl
  log_analysis_instance_id          = "your_rds_instance_id"
  slow_log_bucket_name              = "your_obs_bucket"
  slow_log_start_time               = "2026-08-12T00:00:00Z"
  slow_log_end_time                 = "2026-08-13T23:59:59Z"
  binlog_binlog_type                = "mysql"
  binlog_file_name                  = "mysql binlog file name"
  binlog_export_bucket_name         = "your_obs_bucket"
  binlog_filter_start_time          = "2000-06-01T00:00:00+08:00"
  binlog_filter_end_time            = "2099-06-02T00:00:00+08:00"
  binlog_filter_types               = ["insert", "update", "delete", "ddl"]
  binlog_filter_parse_double_insert = true
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
* The `huaweicloud_das_slow_log_export_task` and `huaweicloud_das_binlog_parse_task_export`
  are one-time action resources; deleting them from the configuration will not clear the corresponding
  request records on the server side
* The slow log export task requires the RDS instance to have slow query log enabled
* The binlog parse task requires the RDS instance to have binlog enabled
* The user ID is automatically obtained from the `huaweicloud_das_database_users` data source
* All resources will be created in the specified region

## Requirements

| Name | Version |
| ---- | ---- |
| terraform | >= 1.1.0 |
| huaweicloud | >= 1.94.0 |
