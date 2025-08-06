# Create an ECS instance with EIP association

This example provides best practice code for using Terraform to create an ECS instance in HuaweiCloud ECS service
and associate it with an Elastic IP (EIP) for external access.

## Prerequisites

* A HuaweiCloud account
* Terraform installed
* HuaweiCloud access key and secret key (AK/SK)

## Variable Introduction

The following variables need to be configured:

### Authentication Variables

* `region_name` - The region where the ECS instance is located
* `access_key` - The access key of the IAM user
* `secret_key` - The secret key of the IAM user

### Resource Variables

#### Required Variables

* `vpc_name` - The name of the VPC
* `subnet_name` - The name of the VPC subnet
* `security_group_name` - The name of the security group
* `instance_name` - The name of the ECS instance
* `bandwidth_name` - The name of the EIP bandwidth (required if EIP address is not provided)

#### Optional Variables

* `availability_zone` - The availability zone to which the ECS instance belongs (default: "")
* `instance_flavor_id` - The flavor ID of the ECS instance (default: "")
* `instance_flavor_performance_type` - The performance type of the ECS instance flavor (default: "normal")
* `instance_flavor_cpu_core_count` - The CPU core count of the ECS instance flavor (default: 2)
* `instance_flavor_memory_size` - The memory size of the ECS instance flavor (default: 4)
* `instance_image_id` - The image ID of the ECS instance (default: "")
* `instance_image_visibility` - The visibility of the ECS instance image (default: "public")
* `instance_image_os` - The OS of the ECS instance image (default: "Ubuntu")
* `vpc_cidr` - The CIDR block of the VPC (default: "192.168.0.0/16")
* `subnet_cidr` - The CIDR block of the VPC subnet (default: "")
* `subnet_gateway_ip` - The gateway IP of the VPC subnet (default: "")
* `associate_eip_address` - The EIP address to associate with the ECS instance (default: "")
* `eip_type` - The type of the EIP (default: "5_bgp")
* `bandwidth_size` - The size of the bandwidth (default: 5)
* `bandwidth_share_type` - The share type of the bandwidth (default: "PER")
* `bandwidth_charge_mode` - The charge mode of the bandwidth (default: "traffic")

## Usage

* Copy this example script to your `main.tf`.

* Create a `terraform.tfvars` file and fill in the required variables:

  ```hcl
  vpc_name            = "your_vpc_name"
  subnet_name         = "your_subnet_name"
  security_group_name = "your_security_group_name"
  instance_name       = "your_ecs_instance_name"
  bandwidth_name      = "your_bandwidth_name"
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

## EIP Association Options

### Option 1: Create New EIP

If you don't provide an existing EIP address, the example will create a new EIP with the specified bandwidth configuration:

```hcl
bandwidth_name = "your_bandwidth_name"
bandwidth_size = 5
eip_type       = "5_bgp"
```

### Option 2: Use Existing EIP

If you have an existing EIP, you can associate it directly:

```hcl
associate_eip_address = "your_existing_eip_address"
```

## Note

* Make sure to keep your credentials secure and never commit them to version control
* The creation of the ECS instance takes about 2-3 minutes
* This example creates the ECS instance, VPC, subnet, security group, EIP, and EIP association
* The EIP association enables external access to the ECS instance
* All resources will be created in the specified region

## Requirements

| Name | Version |
| ---- | ---- |
| terraform | >= 1.2.0 |
| huaweicloud | >= 1.57.0 |
