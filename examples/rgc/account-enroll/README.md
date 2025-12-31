# Enroll an Account in RGC (Resource Governance Center)

This example provides best practice code for RGC (Resource Governance Center) using Terraform
to enroll an account in HuaweiCloud.

## Prerequisites

* A HuaweiCloud account with RGC service enabled
* Terraform installed
* HuaweiCloud access key and secret key (AK/SK)
* An existing parent organizational unit ID (or permission to create one)
* A managed account ID that needs to be enrolled

## Variables Introduction

The following variables need to be configured:

### Authentication Variables

* `region_name` - The region where the RGC resources will be created
* `access_key` - The access key of the IAM user
* `secret_key` - The secret key of the IAM user

### Resource Variables

#### Required Variables

* `parent_organizational_unit_id` - The ID of the parent organizational unit. This is required for
both account enrollment and OU creation

#### Optional Variables

* `organizational_unit_name` - The name of the organizational unit to be created
* `blueprint_managed_account_id` - The ID of the account to be enrolled with blueprint configuration
* `create_organizational_unit` - Whether to create a new organizational unit (default: `true`)
* `blueprint_product_id` - The ID of the blueprint product
* `blueprint_product_version` - The version of the blueprint product
* `blueprint_variables` - The variables for the blueprint configuration (JSON string format)
* `is_blueprint_has_multi_account_resource` - Whether the blueprint has multi-account resources

## Usage

* Copy this example script to your working directory.

* Create a `terraform.tfvars` file and fill in the required variables:

  ```hcl
  # The ID of the parent organizational unit (Required)
  # This is required for both account enrollment and OU creation
  # Replace with your actual parent organizational unit ID
  parent_organizational_unit_id = "ou-xxxxxxxxxxxxx"

  # The ID of the account to be enrolled with blueprint configuration
  # Replace with your actual managed account ID
  blueprint_managed_account_id            = "account-xxxxxxxxxxxxx"
  # Blueprint product configuration
  # Replace with your actual blueprint product ID and version
  blueprint_product_id                    = "blueprint-xxxxxxxxxxxxx"
  blueprint_product_version               = "1.0.0"
  # Blueprint variables in JSON string format
  # Customize these variables according to your blueprint requirements
  blueprint_variables                     = "{\"environment\":\"production\",\"region\":\"cn-north-4\"}"
  # Whether the blueprint has multi-account resources
  is_blueprint_has_multi_account_resource = false
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

* To clean up the resources (note: this will not undo accepted/rejected invitations):

  ```bash
  $ terraform destroy
  ```

## Notes

* All account enrollment, update, and un-enroll operations are asynchronous and may take up to 6 hours to complete.
  The Terraform provider will automatically wait for the operation to complete, but you should be patient
  during the process
* Ensure the blueprint product ID and version are valid and available in your RGC environment
* Blueprint variables must be provided as a valid JSON string format
* If the blueprint contains multi-account resources, set `is_blueprint_has_multi_account_resource` to `true`
* After enrollment, you can check the account status using the `stage` attribute. Common statuses.
  include `ENROLLED`, `ENROLLING`, and `ENROLL_FAILED`
* Make sure to keep your credentials (access_key and secret_key) secure and never commit them to version control
* Consider using environment variables or a secure secrets management system for sensitive information
* If enrollment fails, check the account status and verify that the parent organizational unit ID is correct
* Ensure the managed account ID exists and is accessible before attempting enrollment
* Verify blueprint product ID and version are correct if using blueprint configuration
* When `create_organizational_unit` is set to `true`, the account enrollment will depend on
  the organizational unit creation. Make sure the parent organizational unit ID is valid

## Requirements

| Name | Version |
| ---- | ---- |
| terraform | >= 0.14.0 |
| huaweicloud | >= 1.80.1 |
