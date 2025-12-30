# Configure IAM v5 password policy

This example provides best practice code for using Terraform to configure IAM v5 password policy in HuaweiCloud IAM
service.

## Prerequisites

* A HuaweiCloud account with IAM admin permissions
* Terraform installed
* HuaweiCloud access key and secret key (AK/SK)

## Introduction

This example demonstrates how to configure IAM v5 password policy using Terraform, which helps improve account security
through unified password policy management. The password policy includes:

1. **Password length requirements** - Set the minimum length of passwords
2. **Password complexity requirements** - Set the minimum number of character types that passwords must contain
3. **Consecutive identical characters restriction** - Limit the number of consecutive identical characters in passwords
4. **Password history** - Prevent users from reusing historical passwords
5. **Password validity period** - Set the expiration period for passwords
6. **Password change restrictions** - Set the minimum time interval for password changes
7. **Username restrictions** - Prohibit using username or reversed username as password
8. **User self-service password change** - Control whether users can change their own passwords

## Variables Introduction

The following variables need to be configured:

### Authentication Variables

* `region_name` - The region where the IAM service is located
* `access_key` - The access key of the IAM user
* `secret_key` - The secret key of the IAM user

### Resource Variables

#### Required Variables

* `policy_min_password_length` - The minimum number of characters that a password must contain
* `policy_password_char_combination` - The minimum number of character types that a password must contain. Character
  types include: uppercase letters, lowercase letters, digits, and special characters

#### Optional Variables

* `policy_max_consecutive_identical_chars` - The maximum number of times that a character is allowed to consecutively
  present in a password (default: 0). Value `0` indicates that consecutive identical characters are allowed
* `policy_min_password_age` - The minimum period (minutes) after which users are allowed to make a password change
  (default: 0)
* `policy_password_reuse_prevention` - The password reuse prevention feature indicates the number of historical
  passwords that cannot be reused (default: 3)
* `policy_password_not_username_or_invert` - Whether the password can be the username or the username spelled
  backwards (default: false)
* `policy_password_validity_period` - The password validity period (days) (default: 7). Value `0` indicates that this
  requirement does not apply
* `policy_allow_user_to_change_password` - Whether IAM users are allowed to change their own passwords (default: true)

## Architecture

This example demonstrates the password policy configuration workflow:

1. **Apply password policy** - Apply the policy to the entire account, affecting all IAM users

## Usage

* Copy this example script to your `main.tf`.

* Create a `terraform.tfvars` file and fill in the required variables:

  ```hcl
  policy_max_consecutive_identical_chars = 2
  policy_min_password_age                = 60
  policy_min_password_length             = 8
  policy_password_char_combination       = 2
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

* To clean up the resources (Note: The destroy operation will reset the password policy to default values):

  ```bash
  $ terraform destroy
  ```

## Note

* You *must* have admin privileges to use this resource
* This resource overwrites an existing configuration, make sure one resource per account
* During action `terraform destroy` it sets values the same as defaults for this resource
* Make sure to keep your credentials secure and never commit them to version control
* The password policy applies to the entire account and affects all IAM users
* All parameters have value range restrictions, and Terraform will validate them before applying

## Requirements

| Name | Version |
| ---- | ---- |
| terraform | >= 0.14.0 |
| huaweicloud | = 1.82.3 |
