# Compare GeminiDB Parameter Templates

This example provides best practice code for using Terraform to create two GeminiDB parameter templates with
different parameter values and compare the differences between them in HuaweiCloud GeminiDB service.

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

* `source_template_name` - The name of the source GeminiDB parameter template
* `target_template_name` - The name of the target GeminiDB parameter template

#### Optional Variables - Source Parameter Template

* `source_template_description` - The description of the source template (default: "source parameter template for comparison")
* `source_parameter_values` - The parameter values map for the source template (default: {"request_timeout_in_ms" = "20000"})

#### Optional Variables - Target Parameter Template

* `target_template_description` - The description of the target template (default: "target parameter template for comparison")
* `target_parameter_values` - The parameter values map for the target template (default: {"request_timeout_in_ms" = "30000"})

#### Optional Variables - Datastore (shared by both templates)

* `datastore_type` - The database type (default: "cassandra")
* `datastore_version` - The database version (default: "3.11")
* `datastore_mode` - The database instance mode (default: "CloudNativeCluster")

## Usage

* Copy this example script to your `main.tf`.

* Create a `terraform.tfvars` file and fill in the required variables:

  ```hcl
  source_template_name = "your_source_parameter_template_name"
  target_template_name = "your_target_parameter_template_name"
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
* Both parameter templates are created for GeminiDB Cassandra with `version = "3.11"` and `mode = "CloudNativeCluster"`
* The source template uses `request_timeout_in_ms = "20000"` and the target uses `request_timeout_in_ms = "30000"`
* The `huaweicloud_geminidb_parameter_template_compare` resource compares the differences between the two templates
* This is a one-time operation and cannot be undone
* The `source_configuration_id` and `target_configuration_id` reference the two parameter templates created above
* The `lifecycle.ignore_changes` on both parameter templates ignores `datastore` as it is not returned by the API after
  import
* The compare resource only has `source_configuration_id` and `target_configuration_id` as required arguments and does
  not need `ignore_changes`
* The comparison result is available in the `differences` attribute
* All resources will be created in the specified region

## Requirements

| Name | Version |
| ---- | ---- |
| terraform | >= 0.14.0 |
| huaweicloud | >= 1.92.0 |
| random | >= 3.0.0 |
