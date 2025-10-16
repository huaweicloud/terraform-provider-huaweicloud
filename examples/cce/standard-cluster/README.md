# Create a CCE Standard Cluster

This example provides best practice code for using Terraform to create a Cloud Container Engine (CCE) standard cluster
in HuaweiCloud with all necessary networking components.

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
* `cluster_description` - The description of the CCE cluster (default: "")
* `cluster_tags` - The tags of the CCE cluster (default: {})

## Usage

* Copy this example script to your `main.tf`.

* Create a `terraform.tfvars` file and fill in the required variables:

  ```hcl
  vpc_name            = "tf_test_vpc"
  subnet_name         = "tf_test_subnet"
  bandwidth_name      = "tf_test_bandwidth"
  bandwidth_size      = 5
  cluster_name        = "tf-test-cluster"
  cluster_description = "Created by terraform script"
  cluster_tags        = {
    owner = "terraform"
  }
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

## Cluster Configuration Options

### Option 1: Create New VPC and Subnet

If you don't provide existing VPC and subnet IDs, the example will create new ones:

```hcl
vpc_name    = "tf_test_vpc"    # Create new VPC and using the default CIDR (defined in variables.tf)
subnet_name = "tf_test_subnet" # Create new subnet
vpc_cidr    = "10.0.0.0/16"
```

### Option 2: Use Existing VPC and Subnet

If you have existing VPC and subnet resources, you can use them directly:

```hcl
vpc_id    = "xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx" # Use existing VPC
subnet_id = "xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx" # Use existing subnet
```

### Option 3: Mixed Configuration

You can also mix existing and new resources:

```hcl
vpc_id      = "xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx"
subnet_name = "tf_test_subnet"
```

## EIP Configuration

### Option 1: Create New EIP

If you don't provide an existing EIP address, the example will create a new one:

```hcl
eip_type              = "5_bgp"
bandwidth_name        = "tf_test_bandwidth"
bandwidth_size        = 10
bandwidth_share_type  = "PER"
bandwidth_charge_mode = "traffic"
```

### Option 2: Use Existing EIP

If you have an existing EIP, you can use it directly:

```hcl
eip_address = "1.2.3.4"
```

## Cluster Flavors

The example supports different cluster flavors for different use cases:

* `cce.s1.small` - Small cluster (default, 50 nodes max)
* `cce.s1.medium` - Medium cluster (200 nodes max)
* `cce.s2.small` - Small cluster with enhanced performance (50 nodes max)
* `cce.s2.medium` - Medium cluster with enhanced performance (200 nodes max)
* `cce.s2.large` - Large cluster with enhanced performance (1000 nodes max)
* `cce.s2.xlarge` - Large cluster with enhanced performance (2000 nodes max)

## Container Network Types

The example supports different container network types:

* `overlay_l2` - Overlay L2 network (default)
* `vpc-router` - VPC router network
* `eni` - ENI network

## Cluster Types

* `VirtualMachine` - Virtual machine cluster (default)
* `ARM64` - ARM64 architecture cluster

## Note

* The creation of the CCE cluster usually takes 10-15 minutes
* This example creates a complete CCE cluster with networking components
* All resources will be created in the specified region
* Cluster names must be unique within the region
* The EIP provides external access to the cluster API server
* Container network configuration affects pod-to-pod communication
* Make sure to have sufficient quota for the resources you plan to create

## Requirements

| Name | Version |
| ---- | ---- |
| terraform | >= 1.9.0 |
| huaweicloud | >= 1.35.0 |
