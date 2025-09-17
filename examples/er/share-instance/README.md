# Create an ER instance to share with other accounts

This example provides best practice code for using Terraform to create an Enterprise Router (ER) instance and share it
with other accounts using Resource Access Manager (RAM) in HuaweiCloud.

## Prerequisites

* A HuaweiCloud account
* Terraform installed
* HuaweiCloud access key and secret key (AK/SK) for both owner and principal accounts
* Two different HuaweiCloud accounts (owner and principal) for testing the sharing functionality

## Variable Introduction

The following variables need to be configured:

### Authentication Variables

* `region_name` - The region where the ER service is located
* `access_key` - The access key of the owner IAM user
* `secret_key` - The secret key of the owner IAM user
* `principal_access_key` - The access key of the principal IAM user
* `principal_secret_key` - The secret key of the principal IAM user
* `owner_account_id` - The account ID of the ER instance sharer
* `principal_account_id` - The account ID of the ER instance accepter

### Resource Variables

#### Required Variables

* `instance_name` - The name of the ER instance
* `principal_vpc_name` - The name of the VPC in the principal account
* `principal_subnet_name` - The name of the subnet in the principal account
* `attachment_name` - The name of the ER attachment

#### Optional Variables

* `instance_asn` - The ASN of the ER instance (default: 64512)
* `instance_description` - The description of the ER instance (default: "The ER instance to share with other accounts")
* `instance_enable_default_propagation` - Whether to enable the default propagation (default: true)
* `instance_enable_default_association` - Whether to enable the default association (default: true)
* `instance_auto_accept_shared_attachments` - Whether to automatically accept the shared attachments (default: false)
* `resource_share_name` - The name of the RAM resource share (default: "resource-share-er")
* `principal_vpc_cidr` - The CIDR block of the VPC (default: "192.168.0.0/16")
* `principal_subnet_cidr` - The CIDR block of the subnet (default: auto-calculated from VPC CIDR)
* `principal_subnet_gateway_ip` - The gateway IP of the subnet (default: auto-calculated from subnet CIDR)

## Usage

* Copy this example script to your `main.tf`.

* Create a `terraform.tfvars` file and fill in the required variables:

  ```hcl
  instance_name         = "your_er_instance_name"
  principal_vpc_name    = "your_principal_vpc_name"
  principal_subnet_name = "your_principal_subnet_name"
  attachment_name       = "your_attachment_name"
  owner_account_id      = "your_owner_account_id"
  principal_account_id  = "your_principal_account_id"
  ```

* Create an `authentication.auto.tfvars` file with your credentials:

  ```hcl
  region_name           = "cn-north-4"
  access_key            = "your_owner_access_key"
  secret_key            = "your_owner_secret_key"
  principal_access_key  = "your_principal_access_key"
  principal_secret_key  = "your_principal_secret_key"
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
* The creation of resources may take several minutes depending on the configuration
* This example creates a complete infrastructure including ER instance, RAM share, VPC, and attachment resources
* The ER instance sharing process involves multiple steps: creation, sharing, acceptance, and attachment
* All resources will be created in the specified region
* The example demonstrates cross-account resource sharing using RAM
* The principal account will automatically accept the shared ER instance and create VPC attachments

## Requirements

| Name        | Version   |
|-------------|-----------|
| terraform   | >= 1.1.0  |
| huaweicloud | >= 1.73.4 |
