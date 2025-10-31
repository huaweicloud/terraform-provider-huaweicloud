# Create an ECS Instance with attached network interface

This example provides best practice code for using Terraform to create an ECS instance with an attached
network interface in HuaweiCloud.

## Prerequisites

* A HuaweiCloud account
* Terraform installed
* HuaweiCloud access key and secret key (AK/SK)

## Variable Introduction

The following variables need to be configured:

### Authentication Variables

* `region_name` - The region where the ECS instance is located
* `access_key`  - The access key of the IAM user
* `secret_key`  - The secret key of the IAM user

### Resource Variables

#### Required Variables

* `vpc_name` - The name of the VPC
* `subnet_configurations` - The list of subnet configurations for ECS instance
  - `subnet_name` - The name of the subnet
  - `subnet_cidr` - The CIDR block of the subnet (optional)
  - `subnet_gateway_ip` - The gateway IP address of the subnet (optional)
* `security_group_name` - The name of the security group
* `instance_name` - The name of the ECS instance
* `instance_admin_password` - The login password of the ECS instance

#### Optional Variables

* `availability_zone` - The availability zone to which the ECS instance belongs (default: "")
* `instance_flavor_id` - The flavor ID of the ECS instance (default: "")
* `instance_performance_type` - The performance type of the ECS instance flavor (default: "normal")
* `instance_cpu_core_count` - The number of CPU cores of the ECS instance (default: 2)
* `instance_memory_size` - The memory size in GB of the ECS instance (default: 4)
* `instance_image_id` - The image ID of the ECS instance (default: "")
* `instance_image_visibility` - The visibility of the ECS instance image (default: "public")
* `instance_image_os` - The operating system of the ECS instance image (default: "Ubuntu")
* `vpc_cidr` - The CIDR block of the VPC (default: "192.168.0.0/16")
* `attached_network_id` - The ID of the network to which the ECS instance to be attached (default: "")
* `attached_interface_fixed_ip` - The fixed IP address of the ECS instance to be attached (default: null)
* `attached_security_group_ids` - The list of security group IDs of the ECS instance to be attached (default: null)

## Network Interface Attachment Methods

This example supports two methods for attaching network interfaces to ECS instances:

### Optional 1: Attach to existing network

When you have an existing network that you want to attach to the ECS instance:

Example configuration:

```hcl
subnet_configurations = [
  {
    subnet_name = "tf_test_main_subnet"
  }
]
attached_network_id = "existing-network-id"
```

### Optional 2: Create new subnets and attach to ECS instance

When you want to create new subnets and attach one of them to the ECS instance:

* Configure `subnet_configurations` with exactly `2` subnets:
  - First subnet: Used for the ECS instance's primary network
  - Second subnet: Used for the attached network interface

Example configuration:

```hcl
subnet_configurations = [
  {
    subnet_name = "main-subnet"
  },
  {
    subnet_name = "attached-subnet"
  },
]
```

## Usage

* Copy this example script to your `main.tf`.

* Create a `terraform.tfvars` file and fill in the required variables:

  ```hcl
  vpc_name = "your_vpc_name"
  subnet_configurations = [
    {
      subnet_name       = "your_main_subnet_name"
      subnet_cidr       = "your_main_subnet_cidr"
      subnet_gateway_ip = "your_main_subnet_gateway_ip"
    }
  ]
  security_group_name     = "your_security_group_name"
  instance_name           = "your_instance_name"
  instance_admin_password = "your_instance_password"
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
* The creation of the ECS instance and network interface may take several minutes
* All resources will be created in the specified region
* If no security groups are specified for the attached interface, the default security group will be automatically added

## Requirements

| Name         | Version     |
| ------------ | ----------- |
| terraform | >= 1.9.0   |
| huaweicloud | >= 1.61.0   |
