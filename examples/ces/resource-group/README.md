# Create a CES resource group example

This example provides best practice code for using Terraform to create a resource group in HuaweiCloud CES service.

## Prerequisites

* A HuaweiCloud account
* Terraform installed
* HuaweiCloud access key and secret key (AK/SK)

## Variables Introduction

The following variables need to be configured:

### Authentication Variables

* `region_name` - The region where the resources are located
* `access_key` - The access key of the IAM user
* `secret_key` - The secret key of the IAM user

### Resource Variables

#### Required Variables

* `name` - The name of the resource group.

#### Optional Variables

* `type` - The type of the CES resource group.
* `enterprise_project_id` - The enterprise project ID of the resource group.
* `tags` - The key/value to match resources.
* `associated_eps_ids` -  The enterprise project IDs where the resources from.
* `resources` - The list of resources to add into the group.
  + `namespace` - The namespace of the service.
  + `dimensions` - The multiple levels of alarm thresholds.
    - `name` - The dimension name.
    - `value` - The dimension value.

## Usage

* Copy this example script to your `main.tf`.

* Create a `terraform.tfvars` file and fill in the required variables:

  ```hcl
  name = "tf_test_ces_resource_group_name"
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

## Requirements

| Name | Version   |
| ---- |-----------|
| terraform | >= 0.14.0 |
| huaweicloud | >= 1.48.0 |
