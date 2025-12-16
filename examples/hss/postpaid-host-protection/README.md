# Create an HSS host protection in postPaid charging mode

This example provides best practice code for using Terraform to create an HSS host protection in postPaid charging mode
in HuaweiCloud.

## Prerequisites

* A HuaweiCloud account
* Terraform installed
* HuaweiCloud access key and secret key (AK/SK)

## Variable Introduction

The following variables need to be configured:

### Authentication Variables

* `region_name` - The region where the resources are located
* `access_key`  - The access key of the IAM user
* `secret_key`  - The secret key of the IAM user

### Resource Variables

#### Required Variables

* `host_id` - The host ID for the host protection, Before using host protection, it is necessary to ensure that the
  agent status of the host is online
* `protection_version` - The protection version enabled by the host. Possible values are **hss.version.basic**,
  **hss.version.advanced**, **hss.version.enterprise**, and **hss.version.premium**

#### Optional Variables

* `is_wait_host_available` - Whether to wait for the host agent status to become online (default: false)
* `enterprise_project_id` - The ID of the enterprise project to which the host protection belongs (default: null)

## Usage

* Copy this example script to your `main.tf`.

* Create a `terraform.tfvars` file and fill in the required variables:

  ```hcl
  host_id            = "your_host_id"
  protection_version = "hss.version.enterprise"
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

| Name | Version   |
| ---- |-----------|
| terraform | >= 0.14.0 |
| huaweicloud | >= 1.66.1 |
