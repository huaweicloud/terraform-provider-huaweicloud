# Create a migration task

This example provides best practice code for using Terraform to create a migration task in HuaweiCloud SMS service.

## Prerequisites

* A HuaweiCloud account
* Terraform installed
* HuaweiCloud access key and secret key (AK/SK)
* SMS source server already exists and is registered

## Variable Introduction

The following variables need to be configured:

### Authentication Variables

* `region_name` - The region where the SMS task is located
* `access_key` - The access key of the IAM user
* `secret_key` - The secret key of the IAM user

### Resource Variables

#### Required Variables

* `source_server_name` - The name of the SMS source server
* `server_template_name` - The name of the SMS server template
* `migrate_task_type` - The type of the SMS task (e.g., "MIGRATE_FILE" or "MIGRATE_BLOCK")
* `server_os_type` - The OS type of the server (e.g., "LINUX" or "WINDOWS")

## Usage

* Copy this example script to your `main.tf`.

* Create a `terraform.tfvars` file and fill in the required variables:

  ```hcl
  region_name = "cn-north-4"
  access_key  = "your-access-key"
  secret_key  = "your-secret-key"
  
  source_server_name   = "your_source_server_name"
  server_template_name = "your_server_template_name"
  migrate_task_type    = "MIGRATE_BLOCK"
  server_os_type       = "WINDOWS"
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
* This example creates an SMS server template and migration task
* The SMS source server must already exist and be registered in the SMS service
* The migration task type and OS type are configurable via variables
* All resources will be created in the specified region

## Requirements

| Name | Version |
| ---- | ---- |
| terraform | >= 0.12.0 |
| huaweicloud | >= 1.37.0 |
