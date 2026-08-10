# Create a DLI Flink OpenSource SQL job

This example provides best practice code for using Terraform to create a DLI elastic
resource pool, a general queue, and a Flink OpenSource SQL job in HuaweiCloud.

## Prerequisites

* A HuaweiCloud account
* Terraform installed
* HuaweiCloud access key and secret key (AK/SK)

## Variable Introduction

The following variables need to be configured:

### Authentication Variables

* `region_name` - The region where the DLI resources are located
* `access_key`  - The access key of the IAM user
* `secret_key`  - The secret key of the IAM user

### Resource Variables

#### Required Variables

* `elastic_resource_pool_name` - The name of the DLI elastic resource pool
* `elastic_resource_pool_cidr` - The CIDR block of the elastic resource pool.
* `queue_name` - The name of the DLI general queue used to run the Flink OpenSource SQL job
* `job_name` - The name of the Flink OpenSource SQL job
* `job_sql` - The Flink SQL statement that includes source, query and sink
* `job_flink_version` - The Flink version

#### Optional Variables

* `elastic_resource_pool_description` - The description of the elastic resource pool (default: "")
* `elastic_resource_pool_min_cu` - The minimum number of CUs for the elastic resource pool (default: 16)
* `elastic_resource_pool_max_cu` - The maximum number of CUs for the elastic resource pool (default: 64)
* `enterprise_project_id` - The ID of the enterprise project (default: "")
* `elastic_resource_pool_label` - The label of the elastic resource pool (default: {})
* `queue_cu_count` - The CU count of the DLI queue (default: 16)
* `queue_description` - The description of the DLI queue (default: "")
* `job_execution_agency_urn` - The agency URN authorized to DLI.  
  Available when Flink version is 1.15 or later (default: "")
* `job_cu_number` - The number of CUs selected for the Flink OpenSource SQL job (default: 2)
* `job_parallel_number` - The parallelism of the Flink OpenSource SQL job (default: 1)
* `job_manager_cu_number` - The number of CUs in the JobManager (default: 1)
* `job_tm_cus` - The number of CUs for each Task Manager (default: 1)
* `job_tm_slot_num` - The number of slots in each Task Manager (default: null)
* `job_checkpoint_enabled` - Whether to enable the automatic job snapshot function (default: false)
* `job_checkpoint_mode` - The snapshot mode. (default: "exactly_once")
* `job_checkpoint_interval` - The snapshot interval in seconds (default: 10)
* `job_obs_bucket` - The OBS bucket used to save snapshots or logs (default: "")
* `job_log_enabled` - Whether to enable uploading job logs to the OBS bucket (default: false)
* `job_smn_topic` - The SMN topic used to receive job failure notifications (default: "")
* `job_restart_when_exception` - Whether to enable the function of automatically restarting
  a job upon job exceptions (default: false)
* `job_resume_max_num` - The maximum number of retry times upon exceptions (default: -1)
* `job_idle_state_retention` - The idle state retention in seconds (default: 1)
* `job_runtime_config` - The custom optimization parameters when the Flink job is running (default: {})
* `job_tags` - The key/value pairs to associate with the Flink OpenSource SQL job (default: {})
* `job_description` - The description of the Flink OpenSource SQL job (default: "")

## Usage

* Copy this example script to your `main.tf`.

* Create a `terraform.tfvars` file and fill in the required variables:

  ```hcl
  elastic_resource_pool_name = "your_elastic_resource_pool_name"
  elastic_resource_pool_cidr = "your_elastic_resource_pool_cidr"
  queue_name                 = "your_general_queue_name"
  job_name                   = "your_flink_sql_job_name"
  job_flink_version          = "1.15"
  job_sql                    = <<-EOF
  create table dataGenSource(
    user_id string,
    amount int
  ) with (
    'connector' = 'datagen',
    'rows-per-second' = '1',
    'fields.user_id.kind' = 'random',
    'fields.user_id.length' = '3'
  );

  create table printSink(
    user_id string,
    amount int
  ) with (
    'connector' = 'print'
  );

  insert into printSink select * from dataGenSource;
  EOF
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
* This example only creates a **flink_opensource_sql_job**
* The queue type is fixed to **general**
* All resources will be created in the specified region

## Requirements

| Name | Version |
| ---- | ---- |
| terraform | >= 1.1.0 |
| huaweicloud | >= 1.95.0 |
