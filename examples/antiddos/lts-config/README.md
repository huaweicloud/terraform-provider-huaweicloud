# Create a LTS configuration for Anti-DDoS

This example provides best practice code for using Terraform to create a LTS configuration in HuaweiCloud
Anti-DDoS service.

## Prerequisites

* A HuaweiCloud account
* Terraform installed
* HuaweiCloud access key and secret key (AK/SK)

## Variable Introduction

The following variables need to be configured:

### Authentication Variables

* `region_name` - The region where the Anti-DDoS LTS configuration is located
* `access_key`  - The access key of the IAM user
* `secret_key`  - The secret key of the IAM user

### Resource Variables

#### Required Variables

* `lts_group_name` - The name of the LTS group
* `lts_stream_name` - The name of the LTS stream
* `lts_ttl_in_days` - The log expiration time(days)

#### Optional Variables

* `lts_is_favorite` - Whether to favorite the log stream (default: false)
* `enterprise_project_id` - The enterprise project ID (default: null)

## Usage

* Copy this example script to your `main.tf`.

* Create a `terraform.tfvars` file and fill in the required variables:

```hcl
lts_group_name  = "test-lts-group-name"
lts_stream_name = "test-lts-stream-name"
lts_ttl_in_days = 7
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
| huaweicloud | >= 1.77.6 |
