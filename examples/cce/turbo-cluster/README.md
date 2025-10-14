# Create a CCE Turbo Cluster

This example provides best practice code for using Terraform to create a Cloud Container Engine (CCE) turbo cluster in
HuaweiCloud with ENI (Elastic Network Interface) networking and all necessary components.

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
* `eni_ipv4_subnet_id` - The ID of the ENI subnet (default: "")
  If not provided, a new ENI subnet will be created
* `eni_subnet_name` - The name of the ENI subnet (default: "")
  Required when `eni_ipv4_subnet_id` is not provided
* `eni_subnet_cidr` - The CIDR block of the ENI subnet (default: "")
  If not provided, will be calculated from VPC CIDR
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
* `container_network_type` - The type of container network (default: "eni")
* `cluster_description` - The description of the CCE cluster (default: "")
* `cluster_tags` - The tags of the CCE cluster (default: {})

## Usage

* Copy this example script to your `main.tf`.

* Create a `terraform.tfvars` file and fill in the required variables:

  ```hcl
  vpc_name            = "tf_test_vpc"
  subnet_name         = "tf_test_subnet"
  eni_subnet_name     = "tf_test_eni_subnet"
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

## Turbo Cluster vs Standard Cluster

* **Networking**: Turbo clusters use ENI (Elastic Network Interface) for better network performance
* **Performance**: Higher network throughput and lower latency compared to standard clusters
* **Scalability**: Better support for high-performance workloads
* **Cost**: Generally higher cost due to enhanced networking capabilities

### When to Use Turbo Cluster

* High-performance computing workloads
* Applications requiring low latency
* Microservices with high network traffic
* Real-time data processing
* Machine learning workloads

## ENI Configuration Options

### Option 1: Create New ENI Subnet

If you don't provide an existing ENI subnet ID, the example will create a new one:

```hcl
eni_subnet_name = "tf_test_eni_subnet" # Create new VPC and using the second /20 CIDR
```

### Option 2: Use Existing ENI Subnet

If you have an existing ENI subnet, you can use it directly:

```hcl
eni_ipv4_subnet_id = "xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx"
```

### Option 3: Mixed Configuration

You can also mix existing and new resources:

```hcl
vpc_id          = "xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx" # Use existing VPC
subnet_name     = "tf_test_subnet"                       # Create new regular subnet
eni_subnet_name = "tf_test_eni_subnet"                   # Create new ENI subnet
```

## Network Architecture

The turbo cluster creates a dual-subnet architecture:

1. **Regular Subnet**: For cluster management and control plane
2. **ENI Subnet**: For pod networking with enhanced performance

## Cluster Flavors

The example supports different cluster flavors for different use cases:

* `cce.s1.small` - Small cluster (default, 50 nodes max)
* `cce.s1.medium` - Medium cluster (200 nodes max)
* `cce.s2.small` - Small cluster with enhanced performance (50 nodes max)
* `cce.s2.medium` - Medium cluster with enhanced performance (200 nodes max)
* `cce.s2.large` - Large cluster with enhanced performance (1000 nodes max)
* `cce.s2.xlarge` - Large cluster with enhanced performance (2000 nodes max)

## Performance Considerations

### ENI Benefits

* **Higher Throughput**: Up to 25 Gbps per ENI
* **Lower Latency**: Direct network access without overlay
* **Better Security**: Network isolation at the ENI level
* **Scalability**: Better support for high-density workloads

## Note

* The creation of the CCE turbo cluster usually takes 15-20 minutes
* This example creates a complete CCE turbo cluster with dual-subnet networking
* All resources will be created in the specified region
* Cluster names must be unique within the region
* The EIP provides external access to the cluster API server
* ENI networking provides enhanced performance for container workloads
* Make sure to have sufficient quota for the resources you plan to create

## Requirements

| Name | Version |
| ---- | ---- |
| terraform | >= 1.9.0 |
| huaweicloud | >= 1.50.0 |
