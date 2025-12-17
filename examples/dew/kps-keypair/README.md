# Create a KPS keypair

This example provides best practice code for using Terraform to create a KPS keypair in HuaweiCloud.

## Prerequisites

* A HuaweiCloud account
* Terraform installed
* HuaweiCloud access key and secret key (AK/SK)

## Variables Introduction

The following variables need to be configured:

### Authentication Variables

* `region_name` - The region where the KPS keypair is located
* `access_key`  - The access key of the IAM user
* `secret_key`  - The secret key of the IAM user

### Resource Variables

#### Required Variables

* `keypair_name` - The name of the KPS keypair

#### Optional Variables

* `keypair_scope` - The scope of the KPS keypair (default: "user")
* `keypair_user_id` - The user ID to which the KPS keypair belongs (default: "")
* `keypair_encryption_type` - The encryption mode of the KPS keypair (default: "kms")
  The valid value  can be **default** or **kms**
* `kms_key_id` - The ID of the KMS key (At least one of kms_key_id and kms_key_name must be provided
  when keypair_encryption_type set to **kms**)
* `kms_key_name` - The name of the KMS key (At least one of kms_key_id and kms_key_name must be provided
  when keypair_encryption_type set to **kms**)  
* `keypair_description` - The description of the KPS keypair (default: "")

## Usage

* Copy this example script to your `main.tf`.

* Create a `terraform.tfvars` file and fill in the required variables:

  ```hcl
  keypair_name = "your_keypair_name"
  kms_key_id   = "your_kms_key_id"
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
* The creation of the KPS keypair takes about 5 minutes

## Requirements

| Name | Version |
| ---- | ---- |
| terraform | >= 1.9.0 |
| huaweicloud | >= 1.68.0 |
