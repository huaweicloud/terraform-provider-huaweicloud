# Create a DDS parameter template

This example provides best practice code for using Terraform to create a DDS parameter template in HuaweiCloud.

## Prerequisites

* A HuaweiCloud account
* Terraform installed
* HuaweiCloud access key and secret key (AK/SK)

## Variable Introduction

The following variables need to be configured:

### Authentication Variables

* `region_name` - The region where the DDS parameter template is located
* `access_key`  - The access key of the IAM user
* `secret_key`  - The secret key of the IAM user

### Resource Variables

#### Required Variables

* `template_name` - The name of the parameter template
* `template_mapping` - The mapping between parameter names and parameter values

#### Optional Variables

* `template_node_type` - The node type of the parameter template (default: "mongos")
  The value can be `mongos`, `shard`, `config`, `replica`, `readonly`,  `shard_readonly` or `single`
* `database_version` - The database version (default: "4.0")
  The value can be `5.0`, `4.4`, `4.2`, `4.0`, `3.4`
* `template_description` - The description of the DDS parameter template (default: "")

## Usage

* Copy this example script to your `main.tf`.

* Create a `terraform.tfvars` file and fill in the required variables:

  ```hcl
  template_name    = "your_template_name"
  template_mapping = "your_template_mapping"
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
| terraform | >= 0.14.0 |
| huaweicloud | >= 1.48.0 |
