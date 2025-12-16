# Create a SecMaster workflow version

This example provides best practice code for using Terraform to create a SecMaster workflow version in HuaweiCloud.

## Prerequisites

* A HuaweiCloud account
* Terraform installed
* HuaweiCloud access key and secret key (AK/SK)

## Variables Introduction

The following variables need to be configured:

### Authentication Variables

* `region_name` - The region where the workspace is located
* `access_key`  - The access key of the IAM user
* `secret_key`  - The secret key of the IAM user

### Resource Variables

#### Required Variables

* `workflow_name` - The name of the workflow
* `workflow_version_taskflow` - The Base64 encoded of the workflow topology diagram
* `workflow_version_taskconfig` - The parameters configuration of the workflow topology diagram, in JSON format

#### Optional Variables

* `workspace_id` - The ID of the workspace (at least one of workspace_id and workspace_name must be provided)
* `workspace_name` - The name of the workspace (at least one of workspace_id and workspace_name must be provided)
* `workflow_id` - The ID of the workflow (default: "")
* `workflow_version_taskflow_type` - The taskflow type of the workflow (default: "JSON")
* `workflow_version_aop_type` - The aop type of the workflow (default: "NORMAL")
* `workflow_version_description` - The description of the workflow version (default: "")

## Usage

* Copy this example script to your `main.tf`.

* Create a `terraform.tfvars` file and fill in the required variables:

  ```hcl
  workspace_name              = "your_workspace_name"
  workflow_name               = "your_workflow_name"
  workflow_version_taskflow   = "your_workflow_taskflow"
  workflow_version_taskconfig = "your_workflow_taskconfig"
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
* All resources will be created in the specified region

## Requirements

| Name | Version |
| ---- | ---- |
| terraform | >= 1.9.0 |
| huaweicloud | >= 1.78.1 |
