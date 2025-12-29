# Create a Fine-Grained Permission Management for RAM Resource Share Operation

This example provides best practice code for using Terraform to create a RAM (Resource Access Manager) resource share
with fine-grained permission management in HuaweiCloud.

## Prerequisites

* A HuaweiCloud account
* Terraform installed
* HuaweiCloud access key and secret key (AK/SK)
* Resources to be shared (e.g., VPC subnets, DNS zones, etc.)

## Variables Introduction

The following variables need to be configured:

### Authentication Variables

* `region_name` - The region where the RAM resource share is located
* `access_key` - The access key of the IAM user
* `secret_key` - The secret key of the IAM user

### Resource Variables

#### Required Variables

* `resource_share_id` - The ID of the RAM resource share

#### Optional Variables

* `query_resource_type` - The resource type for querying available permissions (default: "")
  Valid values: **vpc:subnets**, **dns:zone**, **dns:resolverRule**, etc.
* `query_permission_type` - The type of the permission to query (default: "ALL")
  Valid values: **RAM_MANAGED**, **CUSTOMER_MANAGED**, **ALL**
* `query_permission_name` - The name of the permission to query (default: "")
* `permission_replace` - Whether to replace existing permissions when associating a new permission (default: false)

## Usage

* Copy this example script to your `main.tf`.

* Create a `terraform.tfvars` file and fill in the required variables:

  ```hcl
  # Query available permissions
  query_resource_type = "vpc:subnets"

  # Resource Share ID
  # Should been replace the real ID of resource share
  resource_share_id = "xxx"
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
* Permission IDs can be obtained from the `huaweicloud_ram_resource_permissions` data source
* When using `replace = true`, be cautious as it will remove all existing permissions
* All resources will be created in the specified region
* The `associated_permissions` attribute in the resource share is computed and reflects the current state
  of all associated permissions
* When using `count` to batch create permissions, if the data source query returns no results,
  no permission resources will be created
* Optional query parameters (`resource_type` and `name`) use conditional expressions to convert
  empty strings to `null` for proper handling of optional filters

## Requirements

| Name | Version |
| ---- | ---- |
| terraform | >= 0.14.0 |
| huaweicloud | >= 1.82.3 |
