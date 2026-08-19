# Create a GaussDB(DWS) snapshot resource

This example provides best practice code for using Terraform to configure DWS snapshot resource within HuaweiCloud.
The example demonstrates how to create a DWS snapshot resource

## Prerequisites

* A HuaweiCloud account
* Terraform installed
* HuaweiCloud access key and secret key (AK/SK)

## Variable Introduction

The following variables need to be configured:

### Authentication Variables

* `region_name` - The region where the DWS resources are located
* `access_key`  - The access key of the IAM user
* `secret_key`  - The secret key of the IAM user

### Resource Variables

#### Required Variables

* `vpc_name` - The name of the VPC
* `vpc_cidr` - The CIDR block of the VPC
* `subnet_name` - The name of the subnet
* `security_group_name` - The name of the security group
* `cluster_name` - The name of the DWS cluster
* `cluster_admin_user_name` - The administrator username of the DWS cluster
* `cluster_admin_user_pwd` - The administrator password of the DWS cluster
* `snapshot_name` - The name of the snapshot

#### Optional Variables

* `availability_zone` - The availability zone of the DWS cluster (default: "")
* `enterprise_project_id` - The ID of the enterprise project (default: "")
* `subnet_cidr` - The CIDR block of the subnet (default: "")
* `subnet_gateway_ip` - The gateway IP of the subnet (default: "")
* `security_group_delete_default_rules` - Whether to delete the default rules of the security group (default: true)
* `cluster_node_type` - The flavor of the DWS cluster node (default: "")
* `cluster_version` - The version of the DWS cluster (default: "")
* `cluster_vcpus` - The vcpus of the DWS cluster (default: 4)
* `cluster_memory` - The memory of the DWS cluster (default: 32)
* `cluster_datastore_type` - The datastore type of the DWS cluster (default: "dws")
* `cluster_number_of_node` - The number of nodes in the DWS cluster (default: 3)
* `cluster_number_of_cn` - The number of CN nodes in the DWS cluster (default: 3)
* `cluster_volume_type` - The volume type of the DWS cluster (default: "SSD")
* `cluster_volume_capacity` - The volume capacity of the DWS cluster in GB (default: "100")
* `snapshot_description` - The description of the snapshot (default: "")

## Usage

* Copy this example script to your `main.tf`.

* Create a `terraform.tfvars` file and fill in the required variables:

  ```hcl
  vpc_name                = "your_vpc_name"
  vpc_cidr                = "your_vpc_cidr"
  subnet_name             = "your_subnet_name"
  security_group_name     = "your_security_group_name"
  cluster_name            = "your_dws_cluster_name"
  cluster_admin_user_name = "your_admin_user_name"
  cluster_admin_user_pwd  = "your_admin_password"
  snapshot_name           = "your_snapshot_name"
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
* Creating a DWS cluster usually takes about 10 to 30 minutes

## Requirements

| Name | Version |
| ---- | ---- |
| terraform | >= 1.1.0 |
| huaweicloud | >= 1.61.0 |
