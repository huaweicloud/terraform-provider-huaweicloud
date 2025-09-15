# Create a SDRS drill

This example provides best practice code for using Terraform to create a SDRS drill in HuaweiCloud.

## Prerequisites

* A HuaweiCloud account
* Terraform installed
* HuaweiCloud access key and secret key (AK/SK)

## Variable Introduction

The following variables need to be configured:

### Authentication Variables

* `region_name` - The region where the resources are located
* `access_key`  - The access key of the IAM user
* `secret_key`  - The secret key of the IAM user

### Resource Variables

#### Required Variables

* `vpc_name` - The VPC name
* `subnet_name` - The subnet name
* `security_group_name` - The security group name
* `ecs_instance_name` - The ECS instance name
* `protection_group_name` - The protection group name
* `protected_instance_name` - The protected instance name
* `drill_name` - The drill name

#### Optional Variables

* `availability_zone` - The availability zone to which the ECS instance and network belong (default: "")  
  If this parameter is not specified, the available zone will be automatically allocated
* `instance_flavor_id` - The flavor ID of the ECS instance (default: "")
* `instance_flavor_performance_type` - The performance type of the ECS instance flavor (default: "normal")
* `instance_flavor_cpu_core_count` - The number of the ECS instance flavor CPU cores (default: 2)
* `instance_flavor_memory_size` - The number of the ECS instance flavor memories (default: 4)
* `instance_image_id` - The image ID of the ECS instance (default: "")
* `instance_image_os_type` - The OS type of the ECS instance flavor (default: "Ubuntu")
* `instance_image_visibility` - The visibility of the ECS instance flavor (default: "public")
* `vpc_cidr` - The CIDR block of the VPC (default: "192.168.0.0/16")
* `subnet_cidr` - The CIDR block of the subnet (default: "192.168.0.0/24")
* `subnet_gateway_ip` - The gateway IP of the subnet (default: "")
* `sdrs_domain_name` - The name of an available SDRS domain (default: null)
* `source_availability_zone` - The production site AZ of the protection group (default: "")  
  If this parameter is not specified, the available zone will be automatically allocated.
  Both `source_availability_zone` and `target_availability_zone` must be set, or both must be empty
* `target_availability_zone` - The disaster recovery site AZ of the protection group (default: "")

## Usage

* Copy this example script to your `main.tf`.

* Create a `terraform.tfvars` file and fill in the required variables:

  ```hcl
  vpc_name                = "your_vpc_name"
  subnet_name             = "your_subnet_name"
  security_group_name     = "your_security_group_name"
  ecs_instance_name       = "your_ecs_instance_name"
  protection_group_name   = "your_protection_group_name"
  protected_instance_name = "your_protected_instance_name"
  drill_name              = "your_drill_name"
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
* All resources will be created in the specified region

## Requirements

| Name | Version   |
| ---- |-----------|
| terraform | >= 1.9.0  |
| huaweicloud | >= 1.57.0 |
