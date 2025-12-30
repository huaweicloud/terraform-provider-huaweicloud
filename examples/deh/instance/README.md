# Create a dedicated host instance

This example provides best practice code for using Terraform to create a instance
in HuaweiCloud DEH (Dedicated Host) service.

## Prerequisites

* A HuaweiCloud account with DEH permissions
* Terraform installed
* HuaweiCloud access key and secret key (AK/SK)

## Variables Introduction

The following variables need to be configured:

### Authentication Variables

* `region_name` - The region where the DEH service is located
* `access_key` - The access key of the IAM user
* `secret_key` - The secret key of the IAM user

### Resource Variables

#### Required Variables

* `instance_name` - The name of the dedicated host instance

#### Optional Variables

* `availability_zone` - The availability zone where the dedicated host will be created (default: "")
* `instance_host_type` - The host type of the dedicated host (default: "")
* `instance_auto_placement` - Whether to enable auto placement for the dedicated host (default: "on")
* `instance_metadata` - The metadata of the dedicated host (default: {})
* `instance_tags` - The tags of the dedicated host (default: {})
* `enterprise_project_id` - The enterprise project ID of the dedicated host (default: null)
* `instance_charging_mode` - The charging mode of the dedicated host (default: "prePaid")
* `instance_period_unit` - The unit of the billing period of the dedicated host (default: "month")
* `instance_period` - The billing period of the dedicated host (default: "1")
* `instance_auto_renew` - Whether to enable auto renew for the dedicated host (default: "false")

## Usage

* Copy this example script to your `main.tf`.

* Create a `terraform.tfvars` file and fill in the required variables:

  ```hcl
  instance_name = "tf_test_instance"
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
* Make sure to have sufficient quota for the dedicated host resources you plan to create.

## Requirements

| Name | Version |
| ---- | ---- |
| terraform | >= 1.1.0 |
| huaweicloud | >= 1.74.0 |
