# Create a RAM Cross-Account Resource Sharing Operation

This example provides best practice code for using Terraform to create a RAM (Resource Access Manager) resource share
for cross-account resource in HuaweiCloud.

## Prerequisites

* A HuaweiCloud account
* Terraform installed
* HuaweiCloud access key and secret key (AK/SK)
* Target account IDs or organization IDs for sharing resources

## Variables Introduction

The following variables need to be configured:

### Authentication Variables

* `region_name` - The region where the RAM resource share is located
* `access_key` - The access key of the IAM user
* `secret_key` - The secret key of the IAM user

### Resource Variables

#### Required Variables

* `resource_share_name` - The name of the resource share
* `principals` - The list of one or more principals (account IDs or organization IDs) to share resources with
  The valid values are allowed:
  - **Account ID**: The 32-character account ID
  - **Organization ID**: The organization ID if sharing within an organization

#### Optional Variables

* `description` - The description of the resource share (default: "")
* `permission_ids` - The list of RAM permissions associated with the resource share (default: [])
* `resource_urns` - The list of URNs of one or more resources to be shared (default: []).
  If not specified, URNs will be automatically generated from created resources (VPC, subnet, security group)  
  The general format is: `serviceName:regionID:domainID:resourceType:resourceID`
* `allow_external_principals` - Whether resources can be shared with any accounts outside the organization

## Usage

* Copy this example script to your `main.tf`.

* Create a `terraform.tfvars` file and fill in the required variables:

  ```hcl
  region_name         = "cn-north-4"
  resource_share_name = "cross-account-vpc-share"
  description         = "Share VPC resources with other accounts"

  # Share with multiple accounts
  # Should been replace the real account IDs
  principals = [
    "01234567890123456789012345678901",
    "98765432109876543210987654321098"
  ]

  # The list of URNs of one or more resources to be shared.
  # Should been replace the real URNs
  resource_urns = [
    "vpc:cn-north-4:8f06724e5c6f41f59d3e2f3ad897bb4d:subnet:5de72eeb-7977-4602-8186-8766982d9bcc",
  ]

  # The list of RAM permissions associated with the resource share
  # Should been replace the real permission IDs
  permission_ids = [
    "f5153698-ca8b-4b3c-a839-13ff71f67885"
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

## Notes

* Make sure to keep your credentials secure and never commit them to version control
* If you want to share existing resources instead, specify their URNs in the `resource_urns` variable
and leave resource creation variables as defaults
* The resource share creation is usually completed quickly
* Before deleting a resource share, ensure all associated resources and principals are disassociated
* Resource URNs must be valid and the resources must exist
* Principals must be valid account IDs or organization IDs
* All resources will be created in the specified region
* Resource share names must be unique within your HuaweiCloud account

## Requirements

| Name | Version |
| ---- | ---- |
| terraform | >= 0.14.0 |
| huaweicloud | >= 1.70.0 |
