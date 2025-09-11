# Create a RMS resource aggregator example

This example provides best practice code for using Terraform to create an RMS (Resource Management Service)
resource aggregator in HuaweiCloud to collect resource data across multiple accounts.

## Prerequisites

* A HuaweiCloud account
* Terraform installed
* HuaweiCloud access key and secret key (AK/SK)

## Variables Introduction

The following variables need to be configured:

### Authentication Variables

* `region_name` - The region where the resources are located
* `access_key` - The access key of the IAM user
* `secret_key` - The secret key of the IAM user

### Resource Variables

#### Required Variables

* `aggregator_name` - The name of the resource aggregator
* `aggregator_type` - The type of the resource aggregator, which can be ACCOUNT or ORGANIZATION

#### Optional Variables

* `account_ids` - The list of source account IDs to be aggregated (only applicable for ACCOUNT type)

## Usage

* Copy this example script to your `main.tf`.

* Create a `terraform.tfvars` file and fill in the required variables:

  ```hcl
  aggregator_name = "your_aggregator_name"
  aggregator_type = "ACCOUNT"
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
* For ACCOUNT type aggregators, you can specify account_ids to aggregate resources from specific accounts
* For ORGANIZATION type aggregators, the aggregator will automatically collect data from all accounts in the organization

## Requirements

| Name | Version |
| ---- | ---- |
| terraform | >= 1.9.0 |
| huaweicloud | >= 1.77.6 |
