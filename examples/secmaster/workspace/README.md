# Create a SecMaster workspace

This example provides best practice code for using Terraform to create a SecMaster workspace in HuaweiCloud.

## Prerequisites

* A HuaweiCloud account
* Terraform installed
* HuaweiCloud access key and secret key (AK/SK)

## Required Variables

The following variables need to be configured:

### Authentication Variables

* `region_name` - The region where the RabbitMQ instance is located
* `access_key`  - The access key of the IAM user
* `secret_key`  - The secret key of the IAM user

### Resource Variables

#### Required Variables

* `workspace_name` - The name of the workspace
* `workspace_project_name` - The name of the project to in which to create the workspace

#### Optional Variables

* `workspace_description` - The description of the workspace (default: "")
* `enterprise_project_id` - The ID of the enterprise project to which the workspace belongs (default: null)
* `workspace_tags` - The key/value pairs to associate with the workspace (default: {})

## Usage

* Copy this example script to your `main.tf`.

* Create a `terraform.tfvars` file and fill in the required variables:

  ```hcl
  workspace_name         = "your_workspace name"
  workspace_project_name = "your_project_name"
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
* The creation of the workspace takes about 10 minutes
* All resources will be created in the specified region

## Requirements

| Name | Version |
| ---- | ---- |
| terraform | >= 0.14.0 |
| huaweicloud | >= 1.76.0 |
