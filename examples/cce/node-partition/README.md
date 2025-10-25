# Create a CCE node with partition

This example provides best practice code for using Terraform to create a Cloud Container Engine (CCE) node with
partition in HuaweiCloud, including all necessary networking components, ENI subnet configuration,
and node partition configurations. You can also create a node pool with partition.

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

* `node_name` - The name of the worker node
* `node_password` - The root password to login node

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
* `eni_ipv4_subnet_id` - The ID of the ENI subnet (default: "")
  If not provided, a new ENI subnet will be created
* `eni_subnet_name` - The name of the ENI subnet (default: "")
  Required when `eni_ipv4_subnet_id` is not provided
* `eni_subnet_cidr` - The CIDR block of the ENI subnet (default: "")
  If not provided, will be calculated from VPC CIDR
* `cluster_name` - The name of the CCE cluster (default: "")
* `cluster_flavor_id` - The flavor ID of the CCE cluster (default: "cce.s1.small")
* `cluster_version` - The version of the CCE cluster (default: null, uses latest version)
* `cluster_type` - The type of the CCE cluster (default: "VirtualMachine")
* `container_network_type` - The type of container network (default: "eni")
* `cluster_description` - The description of the cluster (default: "")
* `cluster_tags` - The tags of the cluster (default: {})
* `node_flavor_id` - The flavor ID of the node (default: "")
  If not provided, will be determined by performance type and specifications
* `node_flavor_performance_type` - The performance type of the node (default: "normal")
* `node_flavor_cpu_core_count` - The CPU core count of the node (default: 2)
* `node_flavor_memory_size` - The memory size of the node in GB (default: 4)
* `node_partition` - The name of the partition (default: "")
  If not provided, a new partition will be created
* `partition_name` - The name of the partition (default: "")
  Required when `node_partition` is not provided
* `partition_category` - The category of the partition (default: "IES")  
  Valid values are `Default`, `IES`, and `HomeZone`
* `partition_public_border_group` - The group of the partition (default: "")
  Required when `node_partition` is not provided
* `root_volume_type` - The type of the root volume (default: "SSD")
* `root_volume_size` - The size of the root volume in GB (default: 40)
* `data_volumes_configuration` - The configuration of the data volumes (default: [])
  + `volumetype` - The type of the data volume
  + `size` - The size of the data volume in GB
* `node_pool_name` - The name of the node pool (default: "")
* `node_pool_os_type` - The OS type of the node pool (default: "EulerOS 2.9")
* `node_pool_initial_node_count` - The initial number of nodes in the node pool (default: 1)
* `node_pool_password` - The root password to login node (default: "")
  Required when `node_pool_name` is provided

## Usage

* Copy this example script to your `main.tf`.

* Create a `terraform.tfvars` file and fill in the required variables:

  ```hcl
  vpc_name                      = "tf_test_vpc"
  eni_subnet_name               = "tf_test_eni_subnet"
  cluster_name                  = "tf-test-cluster"
  partition_name                = "center"
  partition_category            = "IES"
  partition_public_border_group = "your_partition_border_group"
  node_name                     = "tf-test-node"
  node_password                 = "your_password"

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

## Other Configuration Options

### Use existing partition to create node

 ```hcl
  vpc_name        = "tf_test_vpc"
  eni_subnet_name = "tf_test_eni_subnet"
  cluster_name    = "tf-test-cluster"
  node_partition  = "center"
  node_name       = "tf-test-node"
  node_password   = "your_password"

  data_volumes_configuration = [
    {
      volumetype = "SSD"
      size       = 100
    }
  ]
```

### Create node pool with partition

If you want to create a node pool with partition, you can use the following variables:

```hcl
  node_pool_name     = "tf-test-node-pool"
  node_pool_password = "your_node_pool_password"
```

## Note

* This example creates a CCE cluster with ENI network type, which requires specific network configuration.
* When using existing partitions, make sure the partition is compatible with your cluster configuration.
* The ENI subnet is automatically created if not specified, but you can also use existing ENI subnets.

## Requirements

| Name | Version |
| ---- | ---- |
| terraform | >= 1.9.0 |
| huaweicloud | >= 1.74.0 |
