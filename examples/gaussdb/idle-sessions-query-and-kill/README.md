# Query and Kill GaussDB Idle Sessions

This example provides best practice code for using Terraform to query and kill idle sessions of a
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
* `node_id` - The ID of the node to kill idle sessions (CN or DN primary/standby)
* `component_id` - The ID of the component to kill idle sessions

#### Optional Variables

* `success` - Whether the kill session request is successful

## Usage

* Copy this example script to your `main.tf`.

* Create a `terraform.tfvars` file and fill in the required variables:

  ```hcl
  instance_id  = "your_instance_id"
  node_id      = "your_node_id"
  component_id = "your_component_id"
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
* Use `data.huaweicloud_gaussdb_key_view_nodes_deliver` to query the available node IDs and component IDs of the
  instance

## Requirements

| Name        | Version   |
|-------------|-----------|
| terraform   | >= 0.14.0 |
| huaweicloud | >= 1.95.0 |
