# DCS Redis diagnosis task

This example provides best practice code for using Terraform to create a DCS Redis diagnosis task that diagnoses an
existing DCS Redis instance within a specified time range.

## Prerequisites

* A HuaweiCloud account
* Terraform installed
* HuaweiCloud access key and secret key (AK/SK)
* An existing DCS Redis instance

## Required Variables

The following variables need to be configured:

### Authentication Variables

* `region_name` - The region where the DCS instance is located
* `access_key`  - The access key of the IAM user
* `secret_key`  - The secret key of the IAM user

### Resource Variables

#### Required Variables

* `vpc_name` - The name of the VPC
* `subnet_name` - The name of the subnet
* `instance_name` - The name of the Redis single instance

#### Optional Variables

* `instance_id` - The ID of the DCS instance to diagnose
* `begin_time` - The start time of the diagnosis task, in RFC3339 format
* `end_time` - The end time of the diagnosis task, in RFC3339 format

## Usage

* Copy this example script to your `main.tf`.

* Create a `terraform.tfvars` file and fill in the required variables:

  ```hcl
  instance_id = "your_dcs_instance_id"
  begin_time  = "your_begin_time"
  end_time    = "your_end_time"
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
* This example uses an existing DCS instance, no new instance is created
* `begin_time` and `end_time` must be in RFC3339 format, e.g., `2024-01-01T00:00:00Z`
* The time range should be within the instance's lifetime for meaningful diagnosis results
* If `node_ip_list` is not specified, all nodes of the instance will be diagnosed
* The create process first waits for the instance to be in `RUNNING` state, then waits for diagnosis status to change
  from `diagnosing` to `finished`
* All parameters are ForceNew, changing any parameter will recreate the resource
* The `begin_time` and `end_time` only compare year-month-day hour:minute, seconds are ignored
* The resource supports import using `<instance_id>/<report_id>` format
* The create timeout is 30 minutes
* The diagnosis results are available as computed attributes, including `abnormal_item_sum`, `failed_item_sum`, and
  `diagnosis_node_report_list`

## Requirements

| Name        | Version   |
|-------------|-----------|
| terraform   | >= 0.14.0 |
| huaweicloud | >= 1.63.0 |
