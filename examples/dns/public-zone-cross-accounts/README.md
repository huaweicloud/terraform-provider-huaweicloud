# Create a public zone with cross-accounts authorization

This example provides best practice code for using Terraform to create a public DNS zone with cross-accounts
authorization in HuaweiCloud DNS service.

## Prerequisites

* Two HuaweiCloud accounts with DNS permissions:
  + **Master account**: The account that will own the sub-domain zone
  + **Target account**: The account that owns the main domain

* Terraform installed
* HuaweiCloud access key and secret key (AK/SK) for both accounts

## Variable Introduction

The following variables need to be configured:

### Authentication Variables

* `region_name` - The region where the DNS service is located
* `access_key` - The access key of the IAM user in the master account
* `secret_key` - The secret key of the IAM user in the master account
* `target_account_access_key` - The access key of the IAM user in the target account
* `target_account_secret_key` - The secret key of the IAM user in the target account

### Resource Variables

#### Required Variables

* `main_domain_name` - The name of the main domain in the target account
* `sub_domain_prefix` - The prefix of the sub-domain to be created in the master account

#### Optional Variables

* `recordset_type` - The type of the recordset used for authorization verification (default: "TXT")
* `recordset_ttl` - The time to live (TTL) of the recordset (default: 300)

## Architecture

This example demonstrates the cross-accounts DNS zone authorization workflow:

1. **Query the main domain** in the target account
2. **Create DNS zone authorization** in the master account for the sub-domain
3. **Create a recordset** in the target account to verify the authorization
4. **Verify the authorization** to complete the authorization process
5. **Create the DNS zone** in the master account for the sub-domain

## Usage

* Copy this example script to your `main.tf`.

* Create a `terraform.tfvars` file and fill in the required variables:

  ```hcl
  target_account_access_key = "access_key_of_target_account"
  target_account_secret_key = "secret_key_of_target_account"
  main_domain_name          = "example.com"
  sub_domain_prefix         = "dev"
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

## Note

* Make sure to keep your credentials secure and never commit them to version control
* The main domain must exist in the target account before running this example
* The provider version must be `1.80.1` or higher to avoid errors during resource deletion
* The authorization verification process may take some time, so a `local-exec` provisioner with a sleep command is included

## Requirements

| Name | Version |
| ---- | ---- |
| terraform | >= 0.14.0 |
| huaweicloud | >= 1.80.1 |
