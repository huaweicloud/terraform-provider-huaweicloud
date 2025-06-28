# Create a simple ECS instance

This example provides best practice code for using Terraform to create a basic ECS instance in HuaweiCloud.

## Prerequisites

* A HuaweiCloud account
* Terraform installed
* HuaweiCloud access key and secret key (AK/SK)

## Required Variables

The following variables need to be configured:

### Authentication Variables

* `access_key` - HuaweiCloud access key
* `secret_key` - HuaweiCloud secret key
* `region_name` - The region where resources will be created

### Resource Variables

#### Required Variables

* `vpc_name` - The name of the VPC
* `vpc_cidr` - The CIDR block of the VPC
* `subnet_name` - The name of the VPC subnet
* `security_group_name` - The name of the security group
* `instance_name` - The name of the ECS instance
* `administrator_password` - The password of the administrator

#### Optional Variables

* `instance_flavor_id` - The ID of the ECS instance flavor (default: "")
* `instance_flavor_performance_type` - The performance type of the ECS instance flavor (default: "normal")
* `instance_flavor_cpu_core_count` - The CPU core count of the ECS instance flavor (default: 2)
* `instance_flavor_memory_size` - The memory size of the ECS instance flavor (default: 4)
* `instance_image_id` - The ID of the ECS instance image (default: "")
* `instance_image_visibility` - The visibility of the ECS instance image (default: "public")
* `instance_image_os` - The OS of the ECS instance image (default: "Ubuntu")
* `subnet_cidr` - The CIDR block of the VPC subnet (default: "")
* `subnet_gateway_ip` - The gateway IP of the VPC subnet (default: "")

## Usage

* Copy this example script to your `main.tf`.

* Create a `terraform.tfvars` file and fill in the required variables:

  ```hcl
  access_key  = "your_access_key"
  secret_key  = "your_secret_key"
  region_name = "your_region"
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
* The administrator password must meet the complexity requirements of HuaweiCloud
* All resources will be created in the specified region
