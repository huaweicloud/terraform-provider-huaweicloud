# Create an Identity Center Password Policy

This example provides best practice code for using Terraform to create an Identity Center password policy in HuaweiCloud.

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

#### Password Policy Variables

* `policy_max_password_age` - The max password age of the identity center password policy, unit in days (default: `10`)
  + Valid value range: 1 to 1095
* `policy_minimum_password_length` - The minimum password length of the identity center password policy (default: `10`)
* `policy_password_reuse_prevention` - Whether to prohibit the use of the same password as the previous one (default: `true`)
* `policy_require_uppercase_characters` - Whether to require uppercase characters in passwords (default: `true`)
* `policy_require_lowercase_characters` - Whether to require lowercase characters in passwords (default: `true`)
* `policy_require_numbers` - Whether to require numbers in passwords (default: `true`)
* `policy_require_symbols` - Whether to require symbols in passwords (default: `true`)

## Usage

* Copy this example script to your `main.tf`.

* Create a `terraform.tfvars` file and fill in the required variables:

  ```hcl
  # The policy defines a minimum password length of 8 characters, allows uppercase and lowercase letters, numbers,
  # and special characters, and a password validity period of 30 days.
  policy_max_password_age        = 30
  policy_minimum_password_length = 8
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

### Scenario 1: Create a Password Policy with a New Instance

This is the default scenario. The example will:

1. Register the region (if `is_region_need_register` is `true`)
2. Create a new Identity Center instance
3. Create a password policy in the instance

Set the following variables:

```hcl
is_instance_create             = true
is_region_need_register        = true
policy_max_password_age        = 30
policy_minimum_password_length = 8
```

### Scenario 2: Create a Password Policy in an Existing Instance

If you already have an Identity Center instance, you can create a password policy in it:

1. Query the existing Identity Center instance
2. Create a password policy in the existing instance

Set the following variables:

```hcl
is_instance_create             = false
policy_max_password_age        = 30
policy_minimum_password_length = 8
```

### Scenario 3: Custom Password Policy Configuration

You can customize the password policy according to your security requirements:

```hcl
policy_max_password_age             = 90
policy_minimum_password_length      = 12
policy_password_reuse_prevention    = false
policy_require_uppercase_characters = false
```

## Note

* Make sure to keep your credentials secure and never commit them to version control
* The Identity Center instance must exist in the specified region before creating password policies
* If you need to create a new instance, ensure the region is registered first
* Each Identity Center instance can only have one password policy
* The `identity_store_id` is automatically obtained from the instance (created or queried)
* All resources will be created in the specified region
* The `password_reuse_prevention` parameter accepts `1` (enabled) or `null` (disabled) in the API, but the example uses
  a boolean variable for better usability
* The `max_password_age` must be between 1 and 1095 days

## Requirements

| Name | Version |
| ---- | ---- |
| terraform | >= 0.14.0 |
| huaweicloud | >= 1.80.4 |
