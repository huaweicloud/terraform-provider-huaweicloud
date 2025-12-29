# Configure Identity Center Instance SSO Configuration

This example provides best practice code for using Terraform to configure an Identity Center instance SSO configuration
in HuaweiCloud.

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

#### SSO Configuration Variables

* `configuration_type` - The type of the identity center instance configuration
  (default: `"APP_AUTHENTICATION_CONFIGURATION"`)
* `configuration_mfa_mode` - The MFA mode of the identity center instance configuration (default: `null`)
* `configuration_allowed_mfa_types` - The allowed MFA types of the identity center instance configuration
  (default: `null`)
* `configuration_no_mfa_signin_behavior` - The behavior when signing in without MFA (default: `null`)
* `configuration_no_password_signin_behavior` - The behavior when signing in without password (default: `null`)
* `configuration_max_authentication_age` - The maximum authentication age in ISO 8601 duration format (default: `null`)
  + Example: `"PT8H"` (8 hours), `"PT12H"` (12 hours), `"P1D"` (1 day)
  + Format: `PT{hours}H` for hours or `P{days}D` for days

## Usage

* Copy this example script to your `main.tf`.

* Create a `terraform.tfvars` file and fill in the required variables:

  ```hcl
  # Example SSO configuration with context-aware MFA
  configuration_mfa_mode                    = "CONTEXT_AWARE"           # Environmental perception
  configuration_allowed_mfa_types           = ["TOTP"]
  configuration_no_mfa_signin_behavior      = "ALLOWED_WITH_ENROLLMENT" # Verification during login
  configuration_no_password_signin_behavior = "BLOCKED"
  configuration_max_authentication_age      = "PT12H"
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

### Scenario 1: Configure SSO with a New Instance

This is the default scenario. The example will:

1. Register the region (if `is_region_need_register` is `true`)
2. Create a new Identity Center instance
3. Configure SSO settings in the instance

Set the following variables:

```hcl
is_instance_create                        = true
is_region_need_register                   = true
configuration_mfa_mode                    = "CONTEXT_AWARE"
configuration_allowed_mfa_types           = ["TOTP"]
configuration_no_mfa_signin_behavior      = "ALLOWED_WITH_ENROLLMENT"
configuration_no_password_signin_behavior = "BLOCKED"
configuration_max_authentication_age      = "PT12H"
```

### Scenario 2: Configure SSO in an Existing Instance

If you already have an Identity Center instance, you can configure SSO in it:

1. Query the existing Identity Center instance
2. Configure SSO settings in the existing instance

Set the following variables:

```hcl
is_instance_create                        = false
configuration_mfa_mode                    = "ALWAYS_ON"
configuration_allowed_mfa_types           = ["TOTP", "WEBAUTHN_SECURITY_KEY"]
configuration_no_mfa_signin_behavior      = "BLOCKED"
configuration_no_password_signin_behavior = "BLOCKED"
configuration_max_authentication_age      = "PT8H"
```

### Scenario 3: Strict Security Configuration

For environments requiring strict security:

```hcl
configuration_mfa_mode                    = "ALWAYS_ON"
configuration_allowed_mfa_types           = ["TOTP"]
configuration_no_mfa_signin_behavior      = "BLOCKED"
configuration_no_password_signin_behavior = "BLOCKED"
configuration_max_authentication_age      = "PT4H"
```

## Note

* Make sure to keep your credentials secure and never commit them to version control
* The Identity Center instance must exist in the specified region before configuring SSO
* If you need to create a new instance, ensure the region is registered first
* Each Identity Center instance can only have one SSO configuration
* The `identity_store_id` is automatically obtained from the instance (created or queried)
* All resources will be created in the specified region
* The `max_authentication_age` must be specified in ISO 8601 duration format

## Requirements

| Name | Version |
| ---- | ---- |
| terraform | >= 0.14.0 |
| huaweicloud | = 1.80.4 |
