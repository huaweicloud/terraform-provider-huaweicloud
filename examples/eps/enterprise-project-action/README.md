# Operate an enterprise project

This example provides best practice code for using Terraform to enable or disable an enterprise project in HuaweiCloud
Enterprise Project Management Service (EPS).

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

* `enterprise_project_id` - The ID of an existing enterprise project
* `enterprise_project_name` - The name of the enterprise project. If specified, a new enterprise project will be created

-> The `enterprise_project_name` must be provided if `enterprise_project_id` is not provided

#### Optional Variables

* `enterprise_project_description` - The description of the enterprise project (default: "")
* `enterprise_project_type` - The type of the enterprise project (default: "prod")
* `enterprise_project_enable` - Whether to enable the enterprise project (default: true)
* `delete_flag` - Whether to delete the enterprise project on destroy (default: true)
* `enterprise_project_action` - The action to perform (default: "disable")

## Usage

* Copy this example script to your `main.tf`.

* Create a `terraform.tfvars` file and fill in the required variables:

  Create a new enterprise project and perform the action:

  ```hcl
  enterprise_project_name = "tf-test-eps"
  ```

  Or perform the action on an existing enterprise project:

  ```hcl
  enterprise_project_id = "your-enterprise-project-id"
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
* `huaweicloud_enterprise_project_action` is a one-time action resource. Deleting it only removes the record from the
  Terraform state and does not revert the action
* The **poc** type enterprise project does not support disabling operation

## Requirements

| Name | Version |
| ---- | ---- |
| terraform | >= 1.9.0 |
| huaweicloud | >= 1.82.3 |
