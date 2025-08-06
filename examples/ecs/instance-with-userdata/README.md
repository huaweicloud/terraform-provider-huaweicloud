# Create an ECS Instance with User Data

This example provides best practice code for using Terraform to create an ECS instance with user data in HuaweiCloud.

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
* `subnet_name` - The name of the subnet
* `instance_name` - The name of the ECS instance
* `keypair_name` - The keypair name for the ECS instance
* `instance_user_data` - The user data script for the ECS instance

#### Optional Variables

* `availability_zone` - The availability zone to which the ECS instance belongs (default: "")
* `instance_flavor_id` - The flavor ID of the ECS instance (default: "")
* `instance_performance_type` - The performance type of the ECS instance flavor (default: "normal")
* `instance_cpu_core_count` - The number of the vCPUs in the ECS instance flavor (default: 2)
* `instance_memory_size` - The memory size(GB) in the ECS instance flavor (default: 4)
* `vpc_cidr` - The CIDR block of the VPC (default: "192.168.0.0/16")
* `subnet_cidr` - The CIDR block of the subnet (default: "")
* `subnet_gateway_ip` - The gateway IP address of the subnet (default: "")
* `security_group_ids` - The list of security group IDs for the ECS instance (default: [])
* `security_group_names` - The name of the security groups to which the ECS instance belongs (default: [])
* `instance_image_id` - The image ID of the ECS instance (default: "")
* `instance_image_visibility` - The visibility of the image (default: "public")
* `instance_image_os` - The operating system of the image (default: "Ubuntu")
* `keypair_public_key` - The public key for the keypair (default: null)

-> You must provide either `security_group_ids` or `security_group_names`. If both are provided, `security_group_ids` will
   take precedence.

## Usage

* Copy this example script to your `main.tf`.

* Create a `terraform.tfvars` file and fill in the required variables:

  ```hcl
  vpc_name      = "your_vpc_name"
  subnet_name   = "your_subnet_name"
  instance_name = "your_instance_name"
  keypair_name  = "your_keypair_name"
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
* The admin_pass doesn't work with user_data, use key_pair instead
* All resources will be created in the specified region

## Requirements

| Name | Version |
| ---- | ---- |
| terraform | >= 1.1.0 |
| huaweicloud | >=1.57.0 |
