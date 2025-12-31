# Query DEH Quotas

This example provides best practice code for using Terraform to query quotas in HuaweiCloud DEH (Dedicated Host)
service, with support for various filtering scenarios.

## Prerequisites

* A HuaweiCloud account with DEH permissions
* Terraform installed
* HuaweiCloud access key and secret key (AK/SK)

## Variables Introduction

The following variables need to be configured:

### Authentication Variables

* `region_name` - The region where the DEH service is located
* `access_key` - The access key of the IAM user
* `secret_key` - The secret key of the IAM user

### Filter Variables

All filter variables are optional. If not specified, all quotas will be returned.

#### Optional Variables

* `host_type` - The type of the dedicated host to filter quotas (default: `""`)

## Outputs

The example provides five outputs:

* `quotas_with_usage` - The quotas that have been used
  + Useful for identifying which resources are in use
* `quotas_available` - The quotas that have available capacity
  + Useful for identifying resources that can still be created
* `quotas_exhausted` - The quotas that are fully used
  + Useful for identifying resources that need quota increase

## Usage

* Copy this example script to your `main.tf`.

* Create a `terraform.tfvars` file and configure filter variables as needed:

  ```hcl
  # Example: Filter by specific host type
  host_type = "s6"
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

* View the outputs:

  ```bash
  $ terraform output quotas_with_usage
  $ terraform output quotas_available
  $ terraform output quotas_exhausted
  ```

## Example Scenarios

### Scenario 1: Query All Quotas

Query all quotas without any filters:

```hcl
# Leave host_type empty or unset
host_type = ""
```

This will return all quota information for all host types.

### Scenario 2: Filter by Specific Host Type

Query quotas for a specific host type:

```hcl
host_type = "s6"
```

## Note

* Make sure to keep your credentials secure and never commit them to version control

## Requirements

| Name | Version |
| ---- | ---- |
| terraform | >= 1.1.0 |
| huaweicloud | >= 1.82.0 |
