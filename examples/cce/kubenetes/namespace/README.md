# Create a Kubernetes namespace under the CCE environment

This example provides best practice code for using Terraform to create a namespace (via Kubenetes provider) under the
CCE (Cloud Container Engine) environment (The cluster with a node, and enable the public EIP access) in HuaweiCloud.

For the Kubernetes provider authentication, this example demonstrates how to configure the Kubernetes provider to
authenticate with the CCE cluster using cluster certificates.
The authentication is automatically configured using the following CCE cluster certificate information:

* **Cluster CA Certificate**: Used for verifying the cluster's identity
* **Client Certificate**: Used for authenticating the Terraform client
* **Client Key**: Private key for client certificate authentication

The Kubernetes provider configuration automatically extracts and decodes these certificates from the CCE cluster
resource, ensuring secure and seamless authentication without manual certificate management.

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
* `node_name` - The name of the CCE node
* `namespace_name` - The name of the Kubernetes namespace

#### Optional Variables

* `availability_zone` - The availability zone where the resources will be created (default: "")
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
* `eip_address` - The EIP address of the CCE cluster (default: "")
  If not specified, a new EIP will be created
* `eip_type` - The type of the EIP (default: "5_bgp")
* `bandwidth_name` - The name of the bandwidth (default: "")
* `bandwidth_size` - The size of the bandwidth (default: 5)
* `bandwidth_share_type` - The share type of the bandwidth (default: "PER")
* `bandwidth_charge_mode` - The charge mode of the bandwidth (default: "traffic")
* `cluster_flavor_id` - The flavor ID of the CCE cluster (default: "cce.s1.small")
* `cluster_version` - The version of the CCE cluster (default: null)
* `cluster_type` - The type of the CCE cluster (default: "VirtualMachine")
* `container_network_type` - The type of container network (default: "overlay_l2")
* `authentication_mode` - The mode of the CCE cluster (default: "rbac")
* `delete_all_resources_on_terminal` - Whether to delete all resources on terminal (default: true)
* `node_flavor_id` - The flavor ID of the node (default: "")
  If not specified, will be automatically selected based on performance requirements
* `node_performance_type` - The performance type of the node (default: "general")
* `node_cpu_core_count` - The CPU core count of the node (default: 4)
* `node_memory_size` - The memory size of the node (default: 8)
* `root_volume_type` - The type of the root volume (default: "SATA")
* `root_volume_size` - The size of the root volume (default: 40)
* `data_volumes_configuration` - The configuration of the data volumes (default: [])

## Usage

1. Copy this example script to your `main.tf`.

2. Create a `terraform.tfvars` file and fill in the required variables:

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
   namespace_name             = "tf-test-namespace"
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

## Features

This example demonstrates the following features:

1. **Complete Infrastructure Setup**: Creates a complete CCE cluster with all necessary components
2. **Network Configuration**: Automatic VPC and subnet creation with flexible configuration
3. **EIP Management**: Automatic EIP creation and configuration for cluster access
4. **Node Management**: Flexible node configuration with automatic flavor selection
5. **Storage Configuration**: Support for root and data volumes with different types
6. **Kubernetes Integration**: Automatic Kubernetes namespace creation
7. **Security Configuration**: RBAC authentication and keypair-based access
8. **Resource Management**: Configurable resource cleanup on termination

## Configuration Details

### Kubernetes Namespace

The example creates a Kubernetes namespace with the following features:

* **Namespace Management**: Creates a dedicated namespace for application deployment
* **Dependency Management**: Ensures the namespace is created after the node is ready
* **Integration**: Seamless integration with the CCE cluster

### Storage Configuration

The example supports flexible storage configuration:

* **Root Volume**: System volume for the node
* **Data Volumes**: Additional volumes for application data
* **Volume Types**: Support for different volume types (SSD, SATA, etc.)
* **Size Configuration**: Flexible volume size configuration

## Note

* Make sure to keep your credentials secure and never commit them to version control
* The cluster creation process may take several minutes to complete
* The namespace will be created automatically after the node is ready
* All resources will be created in the specified region
* Cluster names must be unique within the region
* When `delete_all_resources_on_terminal` is set to true, all resources will be deleted when the cluster is terminated
* The Kubernetes provider will automatically configure authentication using cluster certificates

## Requirements

| Name | Version |
|------|---------|
| terraform | >= 1.9.0 |
| huaweicloud | >= 1.57.0 |
| kubernetes | >= 1.6.2 |
