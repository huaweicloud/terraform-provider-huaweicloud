# Create GES metadata and graph instance

This example provides best practice code for using Terraform to create a Graph Engine Service (GES) metadata
schema and graph instance in HuaweiCloud. It demonstrates the recommended resource orchestration order:
OBS bucket -> metadata -> graph instance.

## Prerequisites

* A HuaweiCloud account
* Terraform installed
* HuaweiCloud access key and secret key (AK/SK)
* GES service enabled in the target region

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
* `obs_bucket_name` - The name of the OBS bucket for storing GES metadata
* `metadata_name` - The name of the GES metadata
* `graph_name` - The name of the GES graph instance

#### Optional Variables

* `vpc_cidr` - The CIDR block of the VPC (default: "192.168.0.0/16")
* `subnet_cidr` - The CIDR block of the subnet (default: "")
* `subnet_gateway_ip` - The gateway IP address of the subnet (default: "")
* `metadata_description` - The description of the GES metadata (default: "Graph schema metadata for data import")
* `graph_size_type_index` - The graph size type index (default: "1", 1 million edges)
* `graph_tags` - The tags of the GES graph instance (default: {})

## Usage

* Copy this example script to your `main.tf`.

* Create a `terraform.tfvars` file and fill in the required variables:

  ```hcl
  vpc_name            = "your_vpc_name"
  subnet_name         = "your_subnet_name"
  security_group_name = "your_security_group_name"
  obs_bucket_name     = "your_obs_bucket_name"
  metadata_name       = "your_metadata_name"
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

1. **OBS Bucket for Metadata Storage**: Creates an OBS bucket to store GES metadata schema files
2. **Metadata Schema Definition**: Defines vertex and edge labels with properties for data import
3. **Graph Instance Creation**: Creates a GES graph instance with VPC networking
4. **Resource Orchestration**: Uses `depends_on` to ensure metadata is ready before graph creation

## Note

* Make sure to keep your credentials secure and never commit them to version control
* The creation of the GES graph instance may take up to 60 minutes
* Only 50 metadata resources can be created per region
* Metadata property format must comply with the
  [GES data format specification](https://support.huaweicloud.com/usermanual-ges/ges_01_0153.html)
* Most graph instance parameters are ForceNew and cannot be changed after creation
* All resources will be created in the specified region

## Requirements

| Name | Version |
| ---- | ---- |
| terraform | >= 1.1.0 |
| huaweicloud | >= 1.50.0 |
