# Query and Kill GaussDB Specified Sessions

This example provides best practice code for using Terraform to query and kill specified sessions of a
GaussDB instance in HuaweiCloud.

## Prerequisites

* A HuaweiCloud account with GaussDB permissions
* A running GaussDB instance
* Terraform installed
* HuaweiCloud access key and secret key (AK/SK)

## Variable Introduction

The following variables need to be configured:

### Authentication Variables

* `region_name` - The region where the GaussDB instance is located
* `access_key` - The access key of the IAM user
* `secret_key` - The secret key of the IAM user

### Resource Variables

#### Required Variables

* `instance_id` - The ID of the GaussDB instance
* `node_id` - The ID of the node to kill sessions (CN or DN primary/standby)
* `component_id` - The ID of the component to kill sessions
* `session_ids` - The list of session IDs to be killed

## Outputs

* `success` - Whether the kill session request is successful
* `session_ids` - The list of successfully killed session IDs

## Usage

* Copy this example script to your `main.tf`.

* Create a `terraform.tfvars` file and fill in the required variables:

  ```hcl
  instance_id  = "your_instance_id"
  node_id      = "your_node_id"
  component_id = "your_component_id"
  session_ids  = ["your_session_ids"]
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

## Note

* Make sure to keep your credentials secure and never commit them to version control
* Only nodes with CN or DN (primary, standby) components are supported
* Deleting this resource only removes it from the Terraform state, it does not undo the kill session operation

## Requirements

| Name | Version   |
| ---- |-----------|
| terraform | >= 0.14.0 |
| huaweicloud | >= 1.95.0 |
