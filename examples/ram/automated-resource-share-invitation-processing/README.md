# Create a Automated Resource Share Invitation Processing Operation

This example provides best practice code for using Terraform to create a RAM (Resource Access Manager) resource
for automatically accept or reject resource share invitations in HuaweiCloud.

## Prerequisites

* A HuaweiCloud account
* Terraform installed (>= 1.9.0)
* HuaweiCloud access key and secret key (AK/SK)
* At least one pending resource share invitation (for testing)

## Variable Introduction

The following variables need to be configured:

### Authentication Variables

* `region_name` - The region where the RAM service is located
* `access_key` - The access key of the IAM user
* `secret_key` - The secret key of the IAM user

### Resource Variables

#### Required Variables

* `action` - The action to perform on invitations: "accept" or "reject" (default: "accept")

#### Optional Variables

* `resource_share_ids` - List of resource share IDs to query invitations for.
  If empty, all pending invitations will be queried (default: [])

## Usage

* Copy this example script to your working directory.

* Create a `terraform.tfvars` file and fill in the required variables:

  ```hcl
  # The list of `huaweicloud_ram_resource_share` IDs
  # This should be replaced with the real ID
  resource_share_ids = [
    "resource-share-id-1",
    "resource-share-id-2"
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

* To clean up the resources (note: this will not undo accepted/rejected invitations):

  ```bash
  $ terraform destroy
  ```

## Note

* Make sure to keep your credentials secure and never commit them to version control
* The `huaweicloud_ram_resource_share_accepter` resource is a one-time action resource. Once created, it cannot be updated
* Only invitations with "pending" status can be accepted or rejected. Expired or already processed invitations will be ignored
* Running the same configuration multiple times will not cause side effects. Already processed invitations will be skipped
* If `resource_share_ids` is provided, only invitations from those specific resource shares will be processed.
  If empty, all pending invitations will be processed
* After accepting an invitation, the resource share association will be automatically created

## Requirements

| Name | Version |
| ---- | ---- |
| terraform | >= 0.14.0 |
| huaweicloud | >= 1.64.3 |
