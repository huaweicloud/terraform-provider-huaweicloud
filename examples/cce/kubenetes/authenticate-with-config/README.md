# Generate Kubernetes Config File using CCE Cluster Configuration

This example provides best practice code for using Terraform to create a complete CCE (Cloud Container Engine)
environment (The cluster with a node, and enable the public EIP access), and automatically generate a Kubernetes
configuration file (`.kube/config`) for external access and management. The generated configuration file enables
seamless integration with kubectl and other Kubernetes tools.

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

### Kubernetes Configuration File Generation

This example automatically generates a Kubernetes configuration file (`.kube/config`) that contains:

* **Cluster Information**: Complete cluster connection details
* **Authentication Data**: All necessary certificates and keys for secure access
* **Context Configuration**: Pre-configured context for easy kubectl usage
* **Server Information**: External endpoint for cluster access

The generated configuration file enables you to:

* Use `kubectl` commands directly from your local machine
* Access the cluster using standard Kubernetes tools
* Manage resources without additional authentication setup
* Integrate with CI/CD pipelines and automation tools

### Kubernetes Provider Configuration

The example configures the Kubernetes provider to use the generated configuration file:

* **Config Path**: Points to the generated `.kube/config` file
* **Config Context**: Uses the "external" context for external access
* **Automatic Authentication**: No manual certificate management required

## Note

* Make sure to keep your credentials secure and never commit them to version control
* The cluster creation process may take several minutes to complete
* The Kubernetes config file will be generated automatically after cluster creation
* All resources will be created in the specified region
* Cluster names must be unique within the region
* When `delete_all_resources_on_terminal` is set to true, all resources will be deleted when the cluster is terminated
* The generated config file is stored at `.kube/config` and can be used immediately with kubectl
* The config file contains sensitive information and should be protected accordingly

## Requirements

| Name | Version |
|------|---------|
| terraform | >= 1.9.0 |
| huaweicloud | >= 1.57.0 |
| kubernetes | >= 1.6.2 |
