# Create a CPTS Project and a Test Task

This example provides best practice code for using Terraform to create a project and
a pressure test task of Cloud Performance Test Service (CPTS) in HuaweiCloud. Both resources
are configuration-type resources, creating them will not generate any fee, the fee is only
generated when a pressure test task is actually started.

## Prerequisites

* A HuaweiCloud account
* Terraform installed
* HuaweiCloud access key and secret key (AK/SK)
* The CPTS service has been activated in the console (activation is free)

## Variable Introduction

The following variables need to be configured:

### Authentication Variables

* `region_name` - The region where the CPTS service is located
* `access_key` - The access key of the IAM user
* `secret_key` - The secret key of the IAM user

### Resource Variables

#### Required Variables

* `cpts_project_name` - The name of the CPTS project
* `cpts_task_name` - The name of the CPTS test task

#### Optional Variables

* `cpts_project_description` - The description of the CPTS project
* `cpts_task_benchmark_concurrency` - The benchmark concurrency of the CPTS test task

## Usage

* Copy this example script to your `main.tf`.

* Create a `terraform.tfvars` file and fill in the required variables:

  ```hcl
  cpts_project_name               = "tf_test_cpts_project"
  cpts_project_description        = "Created by Terraform for CPTS best practice"
  cpts_task_name                  = "tf_test_cpts_task"
  cpts_task_benchmark_concurrency = 200
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
* The IDs of the CPTS project and the test task are both numeric strings, the `project_id`
  argument of the task cannot be updated, changing it will recreate the task
* The task name can contain a maximum of `42` characters, and the project description can
  contain a maximum of `50` characters
* The `benchmark_concurrency` argument is only a reference for the calculation of the number
  of concurrent users (`number of concurrent users` = `benchmark concurrency` * `concurrency
  ratio`), setting it will not start any pressure test
* The `operation` argument is not used in this example, setting it to `enable` will actually
  start the pressure test task, which requires that all test cases have been added to the task
  and may generate fees
* The `cluster_id` argument is not used in this example, it requires an existing CPTS resource
  group, and a shared resource group can be used when the number of concurrent users is less
  than `1,000`
* Both resources support import using their numeric IDs, note that the `operation` argument
  will not be read back after import

## Requirements

| Name | Version |
| ---- | ---- |
| terraform | >= 0.14.0 |
| huaweicloud | >= 1.36.0 |
