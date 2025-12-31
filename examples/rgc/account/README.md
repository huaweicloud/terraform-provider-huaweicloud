# Create an RGC Account

This example provides best practice code for using Terraform to create a Resource Governance Center (RGC) account
in HuaweiCloud. The example demonstrates how to create an RGC account with organizational unit configuration.

## Prerequisites

* A HuaweiCloud account
* Terraform installed
* HuaweiCloud access key and secret key (AK/SK)
* RGC service enabled in target region
* An existing organizational unit in the organization

## Variable Introduction

The following variables need to be configured:

### Authentication Variables

* `region_name` - The region where RGC account is located
* `access_key` - The access key of IAM user
* `secret_key` - The secret key of IAM user

### Resource Variables

#### Required Variables

* `account_name` - The name of RGC account
* `account_email` - The email of RGC account
* `identity_store_user_name` - The identity store user name of RGC account
* `identity_store_email` - The identity store email of RGC account
* `parent_organizational_unit_name` - The parent organizational unit name of RGC account
* `parent_organizational_unit_id` - The parent organizational unit ID of RGC account

#### Optional Variables

* `account_phone` - The phone number of RGC account (default: "")

## Usage

* Copy this example script to your `main.tf`.

* Create a `terraform.tfvars` file and fill in the required variables:

  ```hcl
  account_name                    = "tf-test-account"
  account_email                   = "tf-test-account@terraform.com"
  identity_store_user_name        = "tf-test-account"
  identity_store_email            = "tf-test-account@terraform.com"
  parent_organizational_unit_name = "your-org-unit-name"
  parent_organizational_unit_id   = "your-org-unit-id"
  ```

* Initialize Terraform:

  ```bash
  terraform init
  ```

* Review the Terraform plan:

  ```bash
  terraform plan
  ```

* Apply the configuration:

  ```bash
  terraform apply
  ```

* To clean up the resources:

  ```bash
  terraform destroy
  ```

## Account Configuration

### Identity Store Configuration

The account is associated with an identity store user:

```hcl
identity_store_user_name = "tf-test-account"
identity_store_email     = "tf-test-account@terraform.com"
```

### Organizational Unit Configuration

The account must be associated with a parent organizational unit:

```hcl
parent_organizational_unit_name = "your-org-unit-name"
parent_organizational_unit_id   = "your-org-unit-id"
```

## Notes

* Make sure to keep your credentials secure and never commit them to version control
* The account must be created within an existing organizational unit
* The parent organizational unit must exist before creating the account
* All resources will be created in the specified region
* Account names must be unique within the organization
* Make sure to have sufficient quota for the resources you plan to create
* The phone number should be a valid phone number format
* The email should be a valid email address format

## Requirements

| Name        | Version   |
|-------------|-----------|
| terraform   | >= 0.14.0 |
| huaweicloud | >= 1.69.0 |
