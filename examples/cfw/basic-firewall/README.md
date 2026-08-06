# Create a Cloud Firewall (CFW) instance with EIP protection

This example provides best practice code for using Terraform to create a Cloud Firewall (CFW) instance in HuaweiCloud
CFW service, including EIP auto-protection and manual EIP protection.

## Prerequisites

* A HuaweiCloud account
* Terraform installed
* HuaweiCloud access key and secret key (AK/SK)

## Variable Introduction

The following variables need to be configured:

### Authentication Variables

* `region_name` - The region where the CFW firewall is located
* `access_key` - The access key of the IAM user
* `secret_key` - The secret key of the IAM user

### Resource Variables

#### Required Variables

* `firewall_name` - The CFW firewall name

#### Optional Variables

* `firewall_flavor` - The flavor version of the firewall (default: "Professional")
* `firewall_charging_mode` - The charging mode of the firewall (default: "postPaid")
* `eip_auto_protection_status` - Whether to enable auto-protection for EIPs (default: 1)
* `eip_protection_enabled` - Whether to enable manual EIP protection for specific existing EIPs (default: false)
* `eip_protection_eip_ids` - The list of existing EIP IDs to protect (default: [])
* `firewall_tags` - The key/value pairs to associate with the firewall (default: environment=test, managed_by=terraform)

## Usage

* Copy this example script to your `main.tf`.
* Create a `terraform.tfvars` file and fill in the required variables:

  ```hcl
  region_name   = "cn-north-4"
  access_key    = "your_access_key"
  secret_key    = "your_secret_key"
  firewall_name = "your_firewall_name"
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
* The creation of the CFW firewall takes about 5 minutes
* The `eip_protection` resource requires existing EIPs in your account
* All resources will be created in the specified region

## Requirements

| Name | Version |
| ---- | ---- |
| terraform | >= 0.14.0 |
| huaweicloud | >= 1.88.0 |
