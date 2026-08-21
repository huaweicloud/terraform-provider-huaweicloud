# GaussDB agency permission policy

This example provides best practice code for using Terraform to manage GaussDB agency permission policies in
HuaweiCloud.

## Prerequisites

* A HuaweiCloud account
* Terraform installed
* HuaweiCloud access key and secret key (AK/SK)

## Required Variables

### Authentication Variables

* `region_name` - The region where the GaussDB agency permission policy is located
* `access_key`  - The access key of the IAM user
* `secret_key`  - The secret key of the IAM user

### Resource Variables

#### Required Variables

* `bind_role_names`   - The permission policies to be bound to the agency
* `unbind_role_names` - The permission policies to be unbound from the agency

#### Optional Variables

* `agency_name` - The agency name, only `RDSAccessProjectResource` is supported (default: "RDSAccessProjectResource")

## Usage

* Copy this example script to your `main.tf`.

* Create a `terraform.tfvars` file and fill in the required variables:

  ```hcl
  bind_role_names   = ["DBS AgencyPolicy", "GaussDB FullAccess"]
  unbind_role_names = ["GaussDB ReadOnlyAccess"]
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

## Note

* The `agency_name` only supports `RDSAccessProjectResource`
* The resource uses the GaussDB (opengauss) service endpoint
* Both `bind_role_names` and `unbind_role_names` are applied in a single PUT API call
* The resource ID is set to the `agency_name` value
* The `agency_name` is non-updatable and will trigger resource replacement if changed
* Deleting this resource only removes it from Terraform state, the agency configuration remains unchanged
* The resource supports import, but `bind_role_names` and `unbind_role_names` are not in the API response
* Use `ImportStateVerifyIgnore` for `bind_role_names` and `unbind_role_names` when importing
* The computed attributes `is_existed` and `roles` are populated from the GET API response

## Requirements

| Name        | Version   |
|-------------|-----------|
| terraform   | >= 0.14.0 |
| huaweicloud | >= 1.96.0 |
