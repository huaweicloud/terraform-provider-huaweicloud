# Create a CCI Network

This example provides best practice code for using Terraform to create a Cloud Container Instance (CCI) network
in HuaweiCloud with all necessary networking components. The example demonstrates how to create a VPC, subnet, security
group, CCI namespace, and CCI network with warm pool configuration.

## Prerequisites

* A HuaweiCloud account
* Terraform installed
* HuaweiCloud access key and secret key (AK/SK)
* CCI service enabled in the target region

## Variable Introduction

The following variables need to be configured:

### Authentication Variables

* `region_name` - The region where the CCI network is located
* `access_key` - The access key of the IAM user
* `secret_key` - The secret key of the IAM user

### Resource Variables

#### Required Variables

* `network_name` - The name of the CCI network
* `namespace_name` - The name of the CCI namespace

#### Optional Variables

* `vpc_name` - The name of the VPC (default: "tf-test-vpc")
* `vpc_cidr` - The CIDR block of the VPC (default: "192.168.0.0/16")
* `subnet_name` - The name of the subnet (default: "tf-test-subnet")
* `subnet_cidr` - The CIDR block of the subnet (default: "192.168.0.0/24")
* `subnet_gateway_ip` - The gateway IP of the subnet (default: "192.168.0.1")
* `security_group_name` - The name of the security group (default: "tf-test-secgroup")
* `warm_pool_size` - The size of the warm pool for the network (default: "10")
* `warm_pool_recycle_interval` - The recycle interval of the warm pool in hours (default: "2")

## Usage

* Copy this example script to your `main.tf`.

* Create a `terraform.tfvars` file and fill in the required variables:

  ```hcl
  network_name   = "tf-test-network"
  namespace_name = "tf-test-namespace"
  ```

* Initialize Terraform:

  ```bash
  terraform init
  ```

* Review the Terraform plan:

  ```bash
  terraform plan
  ```

* Apply the configuration:

  ```bash
  terraform apply
  ```

* To clean up the resources:

  ```bash
  terraform destroy
  ```

## Network Configuration

### Namespace Annotations

The following annotations are automatically inherited from the namespace:

```hcl
annotations = {
  "yangtse.io/project-id" = huaweicloud_cciv2_namespace.test.annotations["tenant.kubernetes.io/project-id"]
  "yangtse.io/domain-id"  = huaweicloud_cciv2_namespace.test.annotations["tenant.kubernetes.io/domain-id"]
}
```

### Warm Pool Configuration

The warm pool is used to pre-allocate resources for faster pod startup. The following annotations are used to configure
the warm pool:

```hcl
annotations = {
  "yangtse.io/warm-pool-size"             = "10"
  "yangtse.io/warm-pool-recycle-interval" = "2"
}
```

## Notes

* Make sure to keep your credentials secure and never commit them to version control
* The network must be associated with at least one subnet
* The warm pool configuration helps reduce pod startup time by pre-allocating resources
* All resources will be created in the specified region
* Network names must be unique within the namespace
* Make sure to have sufficient quota for the resources you plan to create
* The security group is created with default rules deleted for a clean configuration
* The annotations `yangtse.io/project-id` and `yangtse.io/domain-id` are automatically inherited from the namespace
* The warm pool size and recycle interval can be updated by modifying the corresponding variables and
  re-running `terraform apply`

## Requirements

| Name        | Version   |
|-------------|-----------|
| terraform   | >= 0.14.0 |
| huaweicloud | >= 1.73.4 |
