# Create a CCE Cluster with Node

This example provides best practice code for using Terraform to create a Cloud Container Engine (CCE) cluster with
worker nodes in HuaweiCloud, including all necessary networking components and node configurations.

## Prerequisites

* A HuaweiCloud account
* Terraform installed
* HuaweiCloud access key and secret key (AK/SK)

## Variable Introduction

The following variables need to be configured:

### Authentication Variables

* `region_name` - The region where the CCE cluster is located
* `access_key` - The access key of the IAM user
* `secret_key` - The secret key of the IAM user

### Resource Variables

#### Required Variables

* `cluster_name` - The name of the CCE cluster
* `keypair_name` - The name of the keypair for node access
* `node_name` - The name of the worker node

#### Optional Variables

* `vpc_id` - The ID of the VPC (default: "")
  If not provided, a new VPC will be created
* `subnet_id` - The ID of the subnet (default: "")
  If not provided, a new subnet will be created
* `vpc_name` - The name of the VPC (default: "")
  Required when `vpc_id` is not provided
* `vpc_cidr` - The CIDR block of the VPC (default: "192.168.0.0/16")
* `subnet_name` - The name of the subnet (default: "")
  Required when `subnet_id` is not provided
* `subnet_cidr` - The CIDR block of the subnet (default: "")
  If not provided, will be calculated from VPC CIDR
* `subnet_gateway_ip` - The gateway IP of the subnet (default: "")
  If not provided, will be calculated from subnet CIDR
* `availability_zone` - The availability zone where the CCE cluster will be created (default: "")
* `eip_address` - The EIP address of the CCE cluster (default: "")
  If not provided, a new EIP will be created
* `eip_type` - The type of the EIP (default: "5_bgp")
* `bandwidth_name` - The name of the bandwidth (default: "")
* `bandwidth_size` - The size of the bandwidth (default: 5)
* `bandwidth_share_type` - The share type of the bandwidth (default: "PER")
* `bandwidth_charge_mode` - The charge mode of the bandwidth (default: "traffic")
* `cluster_flavor_id` - The flavor ID of the CCE cluster (default: "cce.s1.small")
* `cluster_version` - The version of the CCE cluster (default: null, uses latest version)
* `cluster_type` - The type of the CCE cluster (default: "VirtualMachine")
* `container_network_type` - The type of container network (default: "overlay_l2")
* `node_performance_type` - The performance type of the node (default: "general")
* `node_cpu_core_count` - The CPU core count of the node (default: 4)
* `node_memory_size` - The memory size of the node in GB (default: 8)
* `root_volume_type` - The type of the root volume (default: "SATA")
* `root_volume_size` - The size of the root volume in GB (default: 40)
* `data_volumes_configuration` - The configuration of the data volumes (default: [])

## Usage

* Copy this example script to your `main.tf`.

* Create a `terraform.tfvars` file and fill in the required variables:

  ```hcl
  vpc_name                   = "tf_test_vpc"
  subnet_name                = "tf_test_subnet"
  bandwidth_name             = "tf_test_bandwidth"
  bandwidth_size             = 5
  cluster_name               = "tf-test-cluster"
  node_performance_type      = "computingv3"
  keypair_name               = "tf_test_keypair"
  node_name                  = "tf-test-node"
  root_volume_size           = 40
  root_volume_type           = "SSD"
  data_volumes_configuration = [
    {
      volumetype = "SSD"
      size       = 100
    }
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

## Node Configuration Options

### Performance Types

The example supports a performance type `computingv3`, for the other values you can refer to the data source
`huaweicloud_compute_flavors`.

## Requirements

| Name | Version |
| ---- | ---- |
| terraform | >= 1.9.0 |
| huaweicloud | >= 1.57.0 |
