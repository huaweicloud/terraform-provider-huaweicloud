# Create a DIS stream with auto scaling and data format configuration

This example provides best practice code for using Terraform to create a DIS stream in HuaweiCloud DIS service,
including auto scaling, data format, and compression configuration.

## Prerequisites

* A HuaweiCloud account
* Terraform installed
* HuaweiCloud access key and secret key (AK/SK)

## Variable Introduction

The following variables need to be configured:

### Authentication Variables

* `region_name` - The region where the DIS stream is located
* `access_key` - The access key of the IAM user
* `secret_key` - The secret key of the IAM user

### Resource Variables

#### Required Variables

* `stream_name` - The name of the DIS stream
* `stream_partition_count` - The number of partitions for the DIS stream

#### Optional Variables

* `stream_type` - The type of the DIS stream. Possible values are **COMMON** (normal stream) and **ADVANCED**
  (advanced stream) (default: null)
* `stream_retention_period` - The data retention period in hours. The value ranges from 24 to 72 (default: 24)
* `stream_auto_scale_min_partition_count` - The minimum number of partitions for auto scaling (default: null)
* `stream_auto_scale_max_partition_count` - The maximum number of partitions for auto scaling (default: null)
* `stream_compression_format` - The compression format of the data. Possible values are **zip**, **gzip**,
  **snappy**, **lz4**, and **zstd** (default: null)
* `stream_data_type` - The type of the data. Possible values are **CSV**, **JSON**, and **BLOB** (default: null)
* `stream_csv_delimiter` - The delimiter for CSV data (default: null)
* `stream_data_schema` - The schema of the data in JSON format (default: null)
* `stream_tags` - The key/value pairs to associate with the DIS stream (default: {})

## Usage

* Copy this example script to your `main.tf`.
* Create a `terraform.tfvars` file and fill in the required variables:

  ```hcl
  stream_name            = "your_dis_stream_name"
  stream_partition_count = 2
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
* The `stream_type`, `stream_retention_period`, `stream_auto_scale_min_partition_count`,
  `stream_auto_scale_max_partition_count`, `stream_compression_format`, `stream_data_type`,
  `stream_csv_delimiter`, and `stream_data_schema` parameters cannot be changed after creation
* The `stream_partition_count` parameter can be updated after creation
* All resources will be created in the specified region

## Requirements

| Name | Version   |
| ---- | ---- |
| terraform | >= 0.14.0 |
| huaweicloud | >= 1.30.0 |
