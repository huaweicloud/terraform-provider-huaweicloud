# Create a GeminiDB Parameter Template

This example provides best practice code for using Terraform to create a GeminiDB parameter template with custom
parameter values in HuaweiCloud GeminiDB service.

## Prerequisites

* A HuaweiCloud account
* Terraform installed
* HuaweiCloud access key and secret key (AK/SK)

## Variable Introduction

The following variables need to be configured:

### Authentication Variables

* `region_name` - The region where resources will be created
* `access_key` - The access key of the IAM user
* `secret_key` - The secret key of the IAM user

### Resource Variables

#### Required Variables

* `template_name` - The name of the GeminiDB parameter template

#### Optional Variables

* `template_description` - The description of the GeminiDB parameter template (default: "test configuration update")
* `datastore_type` - The database type (default: "cassandra")
* `datastore_version` - The database version (default: "3.11")
* `datastore_mode` - The database instance mode (default: "CloudNativeCluster")
* `parameter_values` - The parameter values map (default: {"request_timeout_in_ms" = "10000"})

## Usage

* Copy this example script to your `main.tf`.

* Create a `terraform.tfvars` file and fill in the required variables:

  ```hcl
  template_name = "your_geminidb_parameter_template_name"
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
* The parameter template is created for GeminiDB Cassandra with `version = "3.11"` and `mode = "CloudNativeCluster"`
* The `datastore` block is required when `instance_id` is not specified
* The `values` map specifies custom parameter values, e.g. `request_timeout_in_ms = "10000"`
* The `lifecycle.ignore_changes` ignores `datastore` as it is not returned by the API after import
* All resources will be created in the specified region

## Requirements

| Name | Version |
| ---- | ---- |
| terraform | >= 0.14.0 |
| huaweicloud | >= 1.96.0 |
| random | >= 3.0.0 |
