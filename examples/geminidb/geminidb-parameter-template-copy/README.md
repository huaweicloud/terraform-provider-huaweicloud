# Copy a GeminiDB Parameter Template

This example provides best practice code for using Terraform to create a GeminiDB parameter template and copy it
with a new name, description, and parameter values in HuaweiCloud GeminiDB service.

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

* `template_name` - The name of the source GeminiDB parameter template
* `copy_name` - The name of the copied GeminiDB parameter template

#### Optional Variables - Source Parameter Template

* `template_description` - The description of the source parameter template (default: "test configuration")
* `datastore_type` - The database type (default: "cassandra")
* `datastore_version` - The database version (default: "3.11")
* `datastore_mode` - The database instance mode (default: "CloudNativeCluster")
* `parameter_values` - The parameter values map for the source template (default: {"request_timeout_in_ms" = "20000"})

#### Optional Variables - Parameter Template Copy

* `copy_description` - The description of the copied parameter template (default: "test parameter template update")
* `copy_values` - The parameter values map for the copied template (default: {"request_timeout_in_ms" = "10000"})

## Usage

* Copy this example script to your `main.tf`.

* Create a `terraform.tfvars` file and fill in the required variables:

  ```hcl
  template_name = "your_geminidb_parameter_template_name"
  copy_name     = "your_geminidb_parameter_template_copy_name"
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
* The source parameter template is created for GeminiDB Cassandra with `version = "3.11"` and `mode = "CloudNativeCluster"`
* The `huaweicloud_geminidb_parameter_template_copy` resource copies the source template with a new name, description,
  and parameter values
* The `config_id` of the copy resource references the source parameter template created above
* The `lifecycle.ignore_changes` on the source template ignores `datastore` as it is not returned by the API after import
* The `lifecycle.ignore_changes` on the copy resource ignores `config_id` and `values` as they are not returned by the
  API after import
* All resources will be created in the specified region

## Requirements

| Name | Version |
| ---- | ---- |
| terraform | >= 0.14.0 |
| huaweicloud | >= 1.93.0 |
| random | >= 3.0.0 |
