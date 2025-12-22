# Create a SWR repository image retention policy

This example provides best practice code for using Terraform to create a SWR repository image retention policy in
HuaweiCloud.

## Prerequisites

* A HuaweiCloud account with SWR permissions
* Terraform installed
* HuaweiCloud access key and secret key (AK/SK)

## Variable Introduction

The following variables need to be configured:

### Authentication Variables

* `region_name` - The region where the SWR service is located
* `access_key` - The access key of the IAM user
* `secret_key` - The secret key of the IAM user

### Resource Variables

#### Required Variables

* `organization_name` - The organization name.
* `repository_name` - The repository name.
* `policy_type` - The policy type.
* `policy_number` - The policy number.

#### Optional Variables

* `category` - The category (default: "linux")
* `tag_selectors` - The configuration of the tag selectors (default: [])

## Usage

* Copy this example script to your `main.tf`.

* Create a `terraform.tfvars` file and fill in the required variables:

  ```hcl
  organization_name           = "tf_test_swr_organization_name"
  repository_name             = "tf_test_swr_repository_name"
  policy_type                 = "date_rule"
  policy_number               = 30
  tag_selectors_configuration = [
    {
      kind    = "label"
      pattern = "1.1"
    }
  ]
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
* All resources will be created in the specified region

## Requirements

| Name | Version |
| ---- | ---- |
| terraform | >= 0.14.0 |
| huaweicloud | >= 1.48.0 |
