# Create a GES persistence graph with backup

This example provides best practice code for using Terraform to create a persistence edition GES graph instance
with manual backup in HuaweiCloud. The persistence edition supports graph backup, incremental backup, and
interactive transactions.

## Prerequisites

* A HuaweiCloud account
* Terraform installed
* HuaweiCloud access key and secret key (AK/SK)
* GES persistence edition enabled in the target region

## Variable Introduction

The following variables need to be configured:

### Authentication Variables

* `region_name` - The region where the GES resources are located
* `access_key` - The access key of the IAM user
* `secret_key` - The secret key of the IAM user

### Resource Variables

#### Required Variables

* `vpc_name` - The name of the VPC
* `subnet_name` - The name of the subnet
* `security_group_name` - The name of the security group
* `graph_name` - The name of the GES graph instance

#### Optional Variables

* `vpc_cidr` - The CIDR block of the VPC (default: "192.168.0.0/16")
* `subnet_cidr` - The CIDR block of the subnet (default: "")
* `subnet_gateway_ip` - The gateway IP address of the subnet (default: "")
* `vertex_id_type` - The vertex ID type for the persistence edition graph (default: "hash")
* `enable_encryption` - Whether to enable data encryption (default: false)
* `kms_key_id` - The KMS key ID for graph data encryption (default: "", required when `enable_encryption` is true)
* `keep_backup` - Whether to retain backups after the graph is deleted (default: true)
* `graph_tags` - The tags of the GES graph instance (default: {})

## Usage

* Copy this example script to your `main.tf`.

* Create a `terraform.tfvars` file and fill in the required variables:

  ```hcl
  vpc_name            = "your_vpc_name"
  subnet_name         = "your_subnet_name"
  security_group_name = "your_security_group_name"
  graph_name          = "your_graph_name"
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

## Features

This example demonstrates the following features:

1. **Persistence Edition Graph**: Creates a database edition graph with `graph_size_type_index = "6"`
2. **Vertex ID Configuration**: Configures vertex ID type (hash recommended when ID length is uncertain)
3. **Manual Backup**: Triggers a manual graph backup after the graph instance is running
4. **Backup Retention**: Retains backups when the graph is deleted via `keep_backup`
5. **Optional KMS Encryption**: Supports DEW KMS encryption for graph data

## Note

* Make sure to keep your credentials secure and never commit them to version control
* The creation of the persistence edition graph may take up to 60 minutes
* Graph backup is supported on persistence edition graphs only
* `huaweicloud_ges_backup` is a one-time create resource; each `terraform apply` creates a new backup
* When `enable_encryption` is true, `kms_key_id` must be set to a valid DEW KMS key ID
* For production environments, set `keep_backup = true` to preserve backups after graph deletion
* All resources will be created in the specified region

## Requirements

| Name | Version |
| ---- | ---- |
| terraform | >= 1.1.0 |
| huaweicloud | >= 1.51.0 |
