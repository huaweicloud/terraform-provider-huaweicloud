# Create an enterprise project

This example provides best practice code for using Terraform to create and manage an enterprise project in
HuaweiCloud Enterprise Project Management Service (EPS).

## Prerequisites

* A HuaweiCloud account with EPS permissions
* Terraform installed
* HuaweiCloud access key and secret key (AK/SK)

## Variable Introduction

The following variables need to be configured:

### Authentication Variables

* `region_name` - The region where the EPS service is located
* `access_key` - The access key of the IAM user
* `secret_key` - The secret key of the IAM user

### Resource Variables

#### Required Variables

* `enterprise_project_name` - The name of the enterprise project

#### Optional Variables

* `enterprise_project_description` - The description of the enterprise project (default: "")
* `enterprise_project_type` - The type of the enterprise project (default: "prod")
* `enterprise_project_enable` - Whether to enable the enterprise project (default: true)
* `skip_disable_on_destroy` - Whether to skip disabling the enterprise project on destroy (default: false)
* `delete_flag` - Whether to delete the enterprise project on destroy (default: true)

## Usage

* Copy this example script to your `main.tf`.

* Create a `terraform.tfvars` file and fill in the required variables:

  ```hcl
  enterprise_project_name        = "tf-test-eps"
  enterprise_project_description = "Terraform EPS basic example"
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

## Destroy Behavior

The `huaweicloud_enterprise_project` resource supports three destroy behaviors:

* **Default** (`skip_disable_on_destroy = false`, `delete_flag = false`): Disable the enterprise project on destroy.
  The project remains in the cloud but is removed from the Terraform state.
* **Skip disable** (`skip_disable_on_destroy = true`): No operation is performed on destroy.
* **Delete** (`delete_flag = true`): Delete the enterprise project on destroy.
  Ensure the project has no associated resources before enabling this option.

## Note

* Make sure to keep your credentials secure and never commit them to version control
* The **poc** type enterprise project does not support disabling operation
* The enterprise project name must be unique in the domain and cannot include any form of the word "default"

## Requirements

| Name | Version |
| ---- | ---- |
| terraform | >= 0.14.0 |
| huaweicloud | >= 1.82.3 |
