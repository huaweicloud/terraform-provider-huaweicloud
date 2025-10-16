# Create a CCE Node Pool with Advanced Storage Configuration

This example provides best practice code for using Terraform to create a HuaweiCloud CCE (Cloud Container Engine) node
pool with advanced storage configuration, including dynamic volume management and storage groups.

## Prerequisites

* A HuaweiCloud account
* Terraform installed
* HuaweiCloud access key and secret key (AK/SK)
* An existing CCE cluster or permission to create one

## Variable Introduction

The following variables need to be configured:

### Authentication Variables

* `region_name` - The region where the CCE cluster is located
* `access_key` - The access key of the IAM user
* `secret_key` - The secret key of the IAM user

### Resource Variables

#### Required Variables

* `keypair_name` - The name of the keypair for node access
* `node_pool_name` - The name of the node pool

#### Optional Variables

* `vpc_id` - The ID of the VPC (default: "")
  If not specified, a new VPC will be created
* `subnet_id` - The ID of the subnet (default: "")
  If not specified, a new subnet will be created
* `vpc_name` - The name of the VPC (default: "")
  Required if vpc_id is not provided
* `vpc_cidr` - The CIDR block of the VPC (default: "192.168.0.0/16")
* `subnet_name` - The name of the subnet (default: "")
  Required if subnet_id is not provided
* `subnet_cidr` - The CIDR block of the subnet (default: "")
* `subnet_gateway_ip` - The gateway IP of the subnet (default: "")
* `availability_zone` - The availability zone where the CCE cluster will be created (default: "")
* `eip_address` - The EIP address of the CCE cluster (default: "")
  If not specified, a new EIP will be created
* `eip_type` - The type of the EIP (default: "5_bgp")
* `bandwidth_name` - The name of the bandwidth (default: "")
* `bandwidth_size` - The size of the bandwidth (default: 5)
* `bandwidth_share_type` - The share type of the bandwidth (default: "PER")
* `bandwidth_charge_mode` - The charge mode of the bandwidth (default: "traffic")
* `cluster_name` - The name of the CCE cluster (default: "")
* `cluster_flavor_id` - The flavor ID of the CCE cluster (default: "cce.s1.small")
* `cluster_version` - The version of the CCE cluster (default: null)
* `cluster_type` - The type of the CCE cluster (default: "VirtualMachine")
* `container_network_type` - The type of container network (default: "overlay_l2")
* `node_performance_type` - The performance type of the node (default: "general")
* `node_cpu_core_count` - The CPU core count of the node (default: 4)
* `node_memory_size` - The memory size of the node (default: 8)
* `node_pool_type` - The type of the node pool (default: "vm")
* `node_pool_os_type` - The OS type of the node pool (default: "EulerOS 2.9")
* `node_pool_initial_node_count` - The initial node count of the node pool (default: 2)
* `node_pool_min_node_count` - The minimum node count of the node pool (default: 2)
* `node_pool_max_node_count` - The maximum node count of the node pool (default: 10)
* `node_pool_scale_down_cooldown_time` - The scale down cooldown time of the node pool (default: 10)
* `node_pool_priority` - The priority of the node pool (default: 1)
* `node_pool_tags` - The tags of the node pool (default: {})
* `root_volume_type` - The type of the root volume (default: "SATA")
* `root_volume_size` - The size of the root volume (default: 40)
* `data_volumes_configuration` - The configuration of the data volumes (default: [])

## Usage

1. Copy this example script to your `main.tf`.

2. Create a `terraform.tfvars` file and fill in the required variables:

   ```hcl
   vpc_name              = "tf_test_vpc"
   subnet_name           = "tf_test_subnet"
   bandwidth_name        = "tf_test_bandwidth"
   bandwidth_size        = 5
   cluster_name          = "tf-test-cluster"
   node_performance_type = "computingv3"
   keypair_name          = "tf_test_keypair"
   node_pool_name        = "tf-test-node-pool"
   node_pool_tags        = {
     "owner" = "terraform"
   }

   root_volume_size           = 40
   root_volume_type           = "SSD"
   data_volumes_configuration = [
     {
       volumetype     = "SSD"
       size           = 100
       count          = 2
       virtual_spaces = [
         {
           name        = "kubernetes"
           size        = "10%"
           lvm_lv_type = "linear"
         },
         {
           name = "runtime"
           size = "90%"
         }
       ]
     },
     {
       volumetype     = "SSD"
       size           = 100
       count          = 1
       virtual_spaces = [
         {
           name        = "user"
           size        = "100%"
           lvm_lv_type = "linear"
           lvm_path    = "/workspace"
         }
       ]
     }
   ]
   ```

3. Initialize Terraform:

   ```bash
   $ terraform init
   ```

4. Review the Terraform plan:

   ```bash
   $ terraform plan
   ```

5. Apply the configuration:

   ```bash
   $ terraform apply
   ```

6. To clean up the resources:

   ```bash
   $ terraform destroy
   ```

## Configuration Details

### Node Pool Configuration

The example configures the node pool with the following key parameters:

* **Scaling Configuration**:
  - `initial_node_count`: Initial number of nodes
  - `min_node_count`: Minimum number of nodes
  - `max_node_count`: Maximum number of nodes
  - `scale_down_cooldown_time`: Cooldown time for scale-down operations
  - `priority`: Node pool priority for scheduling

* **Compute Configuration**:
  - Automatic flavor selection based on performance requirements
  - Support for different performance types (general, computingv3, etc.)
  - Flexible CPU and memory configuration
  - Availability zone selection

### Advanced Storage Configuration

The example provides advanced storage management with the following features:

* **Dynamic Volume Management**:
  - Automatic volume flattening based on count configuration
  - Support for multiple volume types (SSD, SATA, etc.)
  - KMS encryption support for data volumes
  - Extended parameters for volume customization

* **Storage Groups and Selectors**:
  - Automatic creation of storage selectors based on volume configuration
  - Support for both CCE-managed and user-managed storage groups
  - Virtual space configuration for logical volume management
  - Flexible LVM configuration options

* **Virtual Spaces Configuration**:
  - `name`: Name of the virtual space
  - `size`: Size allocation (percentage or absolute value)
  - `lvm_lv_type`: LVM logical volume type (linear, striped, etc.)
  - `lvm_path`: Mount path for user volumes
  - `runtime_lv_type`: Runtime logical volume type

## Note

* The node pool will automatically manage node scaling based on configured parameters
* Monitor cluster resource usage after creation to optimize configuration
* Storage groups are automatically created based on volume configuration with virtual spaces

## Requirements

| Name | Version |
|------|---------|
| terraform | >= 1.9.0 |
| huaweicloud | >= 1.57.0 |
