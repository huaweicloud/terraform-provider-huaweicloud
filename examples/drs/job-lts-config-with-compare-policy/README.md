# Configure LTS Log Delivery and Compare Policy for a DRS Job

This example provides best practice code for using Terraform to configure the LTS log delivery
and the data comparison policy of an existing Data Replication Service (DRS) job in HuaweiCloud.
The example first creates an LTS log group and stream, then enables the LTS log delivery of the
DRS job, and finally opens a periodic data comparison policy for the same job.

## Prerequisites

* A HuaweiCloud account
* Terraform installed
* HuaweiCloud access key and secret key (AK/SK)
* An existing DRS job, such as a real-time synchronization job

## Variable Introduction

The following variables need to be configured:

### Authentication Variables

* `region_name` - The region where the DRS job is located
* `access_key` - The access key of the IAM user
* `secret_key` - The secret key of the IAM user

### Resource Variables

#### Required Variables

* `lts_group_name` - The name of the LTS group used to store the DRS job logs
* `lts_stream_name` - The name of the LTS stream used to store the DRS job logs
* `drs_job_id` - The ID of the existing DRS job
* `compare_policy_period` - The comparison period of the compare policy, e.g. `* * 1,3,5` for weekly comparison
* `compare_policy_begin_time` - The start time when the comparison policy takes effect, UTC time in `HH:mm:ss` format
* `compare_policy_end_time` - The end time when the comparison policy takes effect, UTC time in `HH:mm:ss` format

#### Optional Variables

* `compare_policy_compare_type` - The list of comparison types, valid values are `object_comparison`, `lines` and `account`
* `compare_policy_compare_policy` - The comparison policy, valid values are `normal` and `manyToOne`
* `compare_policy_interval_hour` - The comparison interval in hours, required for hourly comparison

## Usage

* Copy this example script to your `main.tf`.

* Create a `terraform.tfvars` file and fill in the required variables:

  ```hcl
  lts_group_name            = "tf-test-drs-lts-group"
  lts_stream_name           = "tf-test-drs-lts-stream"
  drs_job_id                = "your-drs-job-id"
  compare_policy_period     = "* * 1,3,5"
  compare_policy_begin_time = "00:00:00"
  compare_policy_end_time   = "04:00:00"
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
* The `job_id` argument must be an existing DRS job that supports data comparison, and the LTS
  log group and stream must be in the same region as the DRS job
* The `job_id` argument of both DRS resources cannot be updated, changing it will only redirect
  the configuration to another job unless the `enable_force_new` argument is set to `true`
* Destroying the `huaweicloud_drs_lts_config` resource only disables the LTS log delivery of the
  job, and destroying the `huaweicloud_drs_compare_policy` resource only closes the comparison
  policy, the DRS job itself will not be affected
* The IDs of both DRS resources are the `job_id` value, they can be imported by the job ID
* The `begin_time` and `end_time` arguments are UTC time, please convert them according to the
  local time zone

## Requirements

| Name | Version |
| ---- | ---- |
| terraform | >= 0.14.0 |
| huaweicloud | >= 1.92.0 |
