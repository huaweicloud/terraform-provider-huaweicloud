# Create an RFS Private Module with a Version

This example provides best practice code for using Terraform to create a Resource Formation Service (RFS)
private module and a module version in HuaweiCloud. The example first creates a private module, and then
creates a module version that references the module package stored in OBS.

## Prerequisites

* A HuaweiCloud account
* Terraform installed
* HuaweiCloud access key and secret key (AK/SK)
* A module package in ZIP format uploaded to an OBS bucket, and the package must meet the following requirements:
  - Contains at least one template file whose name ends with `.tf` or `.tf.json`
  - Does not contain files whose names end with `.tfvars`
  - Does not exceed `1` MB in size before and after decompression
  - Does not contain more than `100` files

## Variable Introduction

The following variables need to be configured:

### Authentication Variables

* `region_name` - The region where the RFS private module is located
* `access_key` - The access key of the IAM user
* `secret_key` - The secret key of the IAM user

### Resource Variables

#### Required Variables

* `module_name` - The name of the RFS private module
* `module_version` - The version number of the RFS private module
* `module_uri` - The OBS address of the private module package

#### Optional Variables

* `module_description` - The description of the RFS private module (default: `""`)
* `version_description` - The description of the private module version (default: `""`)

## Usage

* Copy this example script to your `main.tf`.

* Create a `terraform.tfvars` file and fill in the required variables:

  ```hcl
  module_name    = "tf-test-rfs-module"
  module_version = "1.0.0"
  module_uri     = "https://your-bucket.obs.cn-north-4.myhuaweicloud.com/tf-test-rfs-module.zip"
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
* The `module_name` must be unique within the domain and region, only letters, digits, underscores (_) and
  hyphens (-) are allowed, and it must start with a letter
* The `module_uri` must be replaced with the actual OBS address of your module package, and the package is used,
  logged, displayed and stored in plaintext by RFS
* The `module_version`, `module_uri` and `version_description` arguments of the module version are non-updatable,
  changing them will recreate the module version
* The `module_description` argument of the module cannot be updated to an empty value once it is set
* To publish a new module version, add another `huaweicloud_rfs_private_module_version` resource with a new
  `module_version` value instead of modifying the existing one

## Requirements

| Name | Version |
| ---- | ---- |
| terraform | >= 0.14.0 |
| huaweicloud | >= 1.91.0 |
