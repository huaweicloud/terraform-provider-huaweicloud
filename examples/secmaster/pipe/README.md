# Create a SecMaster Data Pipe

This example provides best practice code for using Terraform to create a SecMaster data pipe
within HuaweiCloud. The data pipe belongs to a dataspace within a workspace.

The example creates the full dependency chain:

* SecMaster Workspace
* SecMaster Dataspace (in the workspace)
* SecMaster Data Pipe (in the dataspace, with index mapping configuration)

## Prerequisites

* A HuaweiCloud account
* Terraform installed
* HuaweiCloud access key and secret key (AK/SK)

## Variable Introduction

The following variables need to be configured:

### Authentication Variables

* `region_name` - The region where the SecMaster resources will be created
* `access_key` - The access key of the IAM user
* `secret_key` - The secret key of the IAM user

### Resource Variables

#### Required Variables

* `workspace_name` - The name of the SecMaster workspace
* `dataspace_name` - The name of the SecMaster dataspace
* `pipe_name` - The name of the data pipe. The name must start with a letter and contain only
  lowercase letters, digits, and asterisks (*)

#### Optional Variables

* `workspace_description` - The description of the SecMaster workspace (default: "Created by Terraform")
* `dataspace_description` - The description of the SecMaster dataspace (default: "Created by Terraform")
* `pipe_description` - The description of the data pipe (default: "Created by Terraform")
* `shards` - The number of partitions for the data pipe, range: 1-64 (default: 3)
* `storage_period` - The data retention period in days, range: 7-180 (default: 30)
* `timestamp_field` - The timestamp field for the data pipe (default: "timestamp")
* `pipe_status` - The status of the pipe, valid values: open, closed (default: "open")

## Usage

* Copy this example script to your `main.tf`.

* Create a `terraform.tfvars` file and fill in the required variables:

  ```hcl
  workspace_name = "tf-workspace-test"
  dataspace_name = "tf-dataspace-test"
  pipe_name      = "tfpipe01"
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
* The `workspace_id` and `dataspace_id` are referenced from the workspace and dataspace resources
* The `project_name` of the workspace is set to the region name
* The `pipe_name` is non-updatable after creation
* The `mapping` parameter is a JSON string that defines the index mapping configuration
* The `status` parameter is set to "open" to enable the pipe
* The data pipe can be imported using `<workspace_id>/<pipe_id>` format

## Requirements

| Name | Version   |
| ---- |-----------|
| terraform | >= 0.14.0 |
| huaweicloud | >= 1.94.0 |
