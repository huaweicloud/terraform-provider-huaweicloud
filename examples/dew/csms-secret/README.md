# Create a CSMS secret

This example provides best practice code for using Terraform to create a CSMS secret in HuaweiCloud.

## Prerequisites

* A HuaweiCloud account
* Terraform installed
* HuaweiCloud access key and secret key (AK/SK)

## Variables Introduction

The following variables need to be configured:

### Authentication Variables

* `region_name` - The region where the CSMS secret is located
* `access_key`  - The access key of the IAM user
* `secret_key`  - The secret key of the IAM user

### Resource Variables

#### Required Variables

* `secret_name` - The name of the secret
* `secret_text` - The plaintext of a text secret

#### Optional Variables

* `secret_type` - The type of the secret (default: "COMMON")
* `kms_key_id` - The ID of the KMS key used to encrypt the secret (default: "")
* `secret_description` - The description of the secret (default: "")
* `enterprise_project_id` - The ID of the enterprise project to which the secret belongs (default: null)
* `secret_tags` - The key/value pairs to associate with the secret (default: {})

## Usage

* Copy this example script to your `main.tf`.

* Create a `terraform.tfvars` file and fill in the required variables:

  ```hcl
  secret_name = "your_secret_name"
  secret_text = "your_secret_text"
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

## Requirements

| Name | Version |
| ---- | ---- |
| terraform | >= 0.14.0 |
| huaweicloud | >= 1.68.0 |
