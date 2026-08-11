# Create a DLI Flink Jar job

This example provides best practice code for using Terraform to create a DLI elastic
resource pool, a general queue, and a Flink Jar job in HuaweiCloud.

## Prerequisites

* A HuaweiCloud account
* Terraform installed
* HuaweiCloud access key and secret key (AK/SK)
* A JAR package uploaded to OBS bucket

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
* `queue_name` - The name of the DLI general queue used to run the Flink Jar job
* `job_name` - The name of the Flink Jar job
* `job_entrypoint` - The OBS path of the JAR file where the job main class is located
* `job_flink_version` - The Flink version

#### Optional Variables

* `elastic_resource_pool_description` - The description of the elastic resource pool (default: "")
* `elastic_resource_pool_min_cu` - The minimum number of CUs for the elastic resource pool (default: 16)
* `elastic_resource_pool_max_cu` - The maximum number of CUs for the elastic resource pool (default: 64)
* `enterprise_project_id` - The ID of the enterprise project (default: "")
* `elastic_resource_pool_label` - The label of the elastic resource pool (default: {})
* `queue_cu_count` - The CU count of the DLI queue (default: 16)
* `queue_description` - The description of the DLI queue (default: "")
* `job_description` - The description of the Flink Jar job (default: "")
* `job_main_class` - The main class of the Flink Jar job (default: "")
* `job_entrypoint_args` - The entrypoint arguments of the Flink Jar job (default: "")
* `job_dependency_jars` - The dependency JAR packages of the Flink Jar job (default: [])
* `job_dependency_files` - The dependency files of the Flink Jar job (default: [])
* `job_execution_agency_urn` - The agency URN authorized to DLI.  
  Available when Flink version is 1.15 or later (default: "")
* `job_feature` - The feature type of the Flink image. Valid values are basic and custom (default: "basic")
* `job_image` - The custom Flink image. Available when feature is custom (default: "")
* `job_cu_num` - The number of CUs selected for the Flink Jar job (default: 2)
* `job_parallel_num` - The parallelism of the Flink Jar job (default: 1)
* `job_obs_bucket` - The OBS bucket used to save logs (default: "")
* `job_log_enabled` - Whether to enable uploading job logs to the OBS bucket (default: false)
* `job_smn_topic` - The SMN topic used to receive job failure notifications (default: "")
* `job_restart_when_exception` - Whether to enable the function of automatically restarting a job upon job
  exceptions (default: false)
* `job_manager_cu_num` - The number of CUs in the JobManager (default: 1)
* `job_tm_cu_num` - The number of CUs for each Task Manager (default: 1)
* `job_tm_slot_num` - The number of slots in each Task Manager (default: null)
* `job_resume_checkpoint` - Whether to resume from the checkpoint upon exceptions (default: false)
* `job_resume_max_num` - The maximum number of retry times upon exceptions (default: -1)
* `job_runtime_config` - The custom optimization parameters when the Flink job is running (default: {})
* `job_checkpoint_path` - The OBS path used to store checkpoints (default: "")
* `job_tags` - The key/value pairs to associate with the Flink Jar job (default: {})
* `job_checkpoint_enabled` - Whether to enable the automatic job snapshot function (default: false)
* `job_checkpoint_mode` - The snapshot mode (default: 1)
* `job_checkpoint_interval` - The snapshot interval in seconds (default: 30)

## Usage

* Copy this example script to your `main.tf`.

* Create a `terraform.tfvars` file and fill in the required variables:

  ```hcl
  elastic_resource_pool_name = "your_elastic_resource_pool_name"
  elastic_resource_pool_cidr = "your_elastic_resource_pool_cidr"
  queue_name                 = "your_general_queue_name"
  job_name                   = "your_flink_jar_job_name"
  job_entrypoint             = "obs://your_bucket_path/your.jar"
  job_flink_version          = "1.15"
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
* The queue type is fixed to **general**
* All resources will be created in the specified region

## Requirements

| Name | Version |
| ---- | ---- |
| terraform | >= 1.1.0 |
| huaweicloud | >= 1.80.5 |
