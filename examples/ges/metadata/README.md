# Create a GES metadata

This example provides best practice code for using Terraform to create a GES (Graph Engine Service) metadata
in HuaweiCloud, including an OBS bucket for storing the metadata schema file.

## Prerequisites

* A HuaweiCloud account
* Terraform installed
* HuaweiCloud access key and secret key (AK/SK)

## Variable Introduction

The following variables need to be configured:

### Authentication Variables

* `region_name` - The region where the GES metadata is located
* `access_key` - The access key of the IAM user
* `secret_key` - The secret key of the IAM user

### Resource Variables

#### Required Variables

* `bucket_name` - The OBS bucket name for storing GES metadata schema files
* `metadata_name` - The GES metadata name

#### Optional Variables

* `metadata_description` - The description of the GES metadata (default: "This is a demo")
* `metadata_schema_file` - The schema file name in the OBS bucket (default: "schema_demo.xml")
* `metadata_properties` - The properties of the GES metadata label
  (default: [{ dataType = "char", name = "sex", cardinality = "single" }])

## Usage

* Copy this example script to your `main.tf`.

* Create a `terraform.tfvars` file and fill in the required variables:

  ```hcl
  bucket_name          = "your_bucket_name"
  metadata_name        = "your_metadata_name"
  metadata_description = "This is a demo"
  metadata_schema_file = "schema_demo.xml"
  metadata_properties  = [{
      dataType    = "char"
      name        = "sex"
      cardinality = "single"
  }]
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
* The OBS bucket name must be globally unique across all HuaweiCloud accounts
* The `force_destroy` option will delete all objects in the bucket when destroying the resource
* The metadata schema file is automatically generated based on the `ges_metadata` block configuration

## Requirements

| Name | Version |
| ---- | ---- |
| terraform | >= 1.1.0 |
| huaweicloud | >= 1.50.0 |
