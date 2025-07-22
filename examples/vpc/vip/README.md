# Create a Virtual IP (VIP) and associate with an ECS instance

This example demonstrates how to use Terraform to create a Virtual IP (VIP) in HuaweiCloud VPC and associate it with a
single ECS instance.

## Prerequisites

* A HuaweiCloud account
* Terraform installed
* HuaweiCloud access key and secret key (AK/SK)

## Required Variables

### Authentication Variables

* `region_name` - The region where the ECS instance is located
* `access_key` - The access key of the IAM user
* `secret_key` - The secret key of the IAM user

### Resource Variables

#### Required Variables

* `vpc_name` - The name of the VPC
* `vpc_cidr` - The CIDR block of the VPC
* `subnet_name` - The name of the VPC subnet
* `security_group_name` - The name of the security group
* `instance_name` - The name of the ECS instance
* `administrator_password` - The password of the administrator

#### Optional Variables

* `instance_flavor_id` - The ID of the ECS instance flavor (default: auto-selected)
* `instance_flavor_performance_type` - The performance type of the ECS instance flavor (default: "normal")
* `instance_flavor_cpu_core_count` - The CPU core count of the ECS instance flavor (default: 2)
* `instance_flavor_memory_size` - The memory size of the ECS instance flavor (default: 4)
* `instance_image_id` - The ID of the ECS instance image (default: auto-selected)
* `instance_image_visibility` - The visibility of the ECS instance image (default: "public")
* `instance_image_os` - The OS of the ECS instance image (default: "Ubuntu")
* `subnet_cidr` - The CIDR block of the VPC subnet (default: auto-calculated)
* `subnet_gateway_ip` - The gateway IP of the VPC subnet (default: auto-calculated)

## Usage

1. Copy the files `main.tf`, `variables.tf`, `providers.tf` to your working directory.
2. Create a `terraform.tfvars` file and fill in the required variables:

   ```hcl
   vpc_name               = "tf_test_vpc"
   vpc_cidr               = "192.168.0.0/16"
   subnet_name            = "tf_test_subnet"
   security_group_name    = "tf_test_security_group"
   instance_name          = "tf_test_instance"
   administrator_password = "YourPasswordInput!"
   ```

3. Initialize and apply:

   ```bash
   terraform init
   terraform plan
   terraform apply
   ```

4. Destroy resources:

   ```bash
   terraform destroy
   ```

## Note

* Do not commit your AK/SK or other sensitive information to version control.
* This example creates a VPC, subnet, security group, one ECS instance, a VIP, and associates the VIP with the ECS instance.
* The VIP will be associated with the single ECS instance created.
* All resources will be created in the specified region.
* The subnet CIDR and gateway IP will be auto-calculated if not specified.

## Requirements

| Name        | Version    |
| ----------- | ---------- |
| terraform   | >= 1.3.0   |
| huaweicloud | >= 1.70.1  |
