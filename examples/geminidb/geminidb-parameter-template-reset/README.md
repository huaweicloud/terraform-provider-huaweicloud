# Reset a GeminiDB Parameter Template

This example provides best practice code for using Terraform to create a GeminiDB parameter template with custom
parameter values and then reset it to its default values in HuaweiCloud GeminiDB service.

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

* `template_description` - The description of the GeminiDB parameter template (default: "test configuration for reset")
* `datastore_type` - The database type (default: "cassandra")
* `datastore_version` - The database version (default: "3.11")
* `datastore_mode` - The database instance mode (default: "CloudNativeCluster")
* `parameter_values` - The parameter values map (default: {"request_timeout_in_ms" = "20000"})

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
* The `huaweicloud_geminidb_parameter_template_reset` resource resets the parameter template to its default values
* This is a one-time operation and cannot be undone
* The `config_id` of the reset resource references the parameter template created above
* The `lifecycle.ignore_changes` on the parameter template ignores `datastore` as it is not returned by the API after import
* The reset resource only has `config_id` as a required argument and does not need `ignore_changes`
* All resources will be created in the specified region

## Requirements

| Name | Version |
| ---- | ---- |
| terraform | >= 0.14.0 |
| huaweicloud | >= 1.96.1 |
| random | >= 3.0.0 |
