# Create an Identity Center Group

This example provides best practice code for using Terraform to create an Identity Center group in HuaweiCloud.

## Prerequisites

* A HuaweiCloud account with Identity Center permissions
* Terraform installed
* HuaweiCloud access key and secret key (AK/SK)
* An Identity Center instance (or permission to create one)

## Variable Introduction

The following variables need to be configured:

### Authentication Variables

* `region_name` - The region where the Identity Center service is located
* `access_key` - The access key of the IAM user
* `secret_key` - The secret key of the IAM user

### Resource Variables

#### Required Variables

* `group_name` - The name of the Identity Center group

#### Optional Variables

* `is_instance_create` - Whether to create the Identity Center instance (default: `true`)
  + If set to `true`, the example will create a new Identity Center instance
  + If set to `false`, the example will query an existing Identity Center instance
* `is_region_need_register` - Whether to register the region before creating the instance (default: `true`)
  + Only applicable when `is_instance_create` is `true`
  + If set to `true`, the example will register the region before creating the instance
* `instance_store_id_alias` - The alias of the Identity Center instance (default: `""`)
  + Only applicable when `is_instance_create` is `true`
  + If left empty, the instance will be created without an alias
* `group_description` - The description of the Identity Center group (default: `""`)

## Usage

* Copy this example script to your `main.tf`.

* Create a `terraform.tfvars` file and fill in the required variables:

  ```hcl
  group_name        = "your_group_name"
  group_description = "Your group description"
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

## Example Scenarios

### Scenario 1: Create a Group with a New Instance

This is the default scenario. The example will:

1. Register the region (if `is_region_need_register` is `true`)
2. Create a new Identity Center instance
3. Create a group in the instance

Set the following variables:

```hcl
is_instance_create      = true
is_region_need_register = true
group_name              = "my_test_group"
```

### Scenario 2: Create a Group in an Existing Instance

If you already have an Identity Center instance, you can create a group in it:

1. Query the existing Identity Center instance
2. Create a group in the existing instance

Set the following variables:

```hcl
is_instance_create = false
group_name         = "my_test_group"
```

## Note

* Make sure to keep your credentials secure and never commit them to version control
* The Identity Center instance must exist in the specified region before creating groups
* If you need to create a new instance, ensure the region is registered first
* Group names must be unique within an Identity Center instance
* The `identity_store_id` is automatically obtained from the instance (created or queried)
* All resources will be created in the specified region

## Requirements

| Name | Version |
| ---- | ---- |
| terraform | >= 0.14.0 |
| huaweicloud | >= 1.80.1 |
