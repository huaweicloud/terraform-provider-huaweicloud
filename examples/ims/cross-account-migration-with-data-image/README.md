# Cross-Account Migration with Data Disk Image

This example provides best practice code for using Terraform to migrate business data across accounts within the same
region by creating a data disk image, sharing it, and creating a new data disk from the shared image in HuaweiCloud
IMS service.

## Prerequisites

* Terraform installed
* HuaweiCloud access key and secret key (AK/SK) for both sharer and accepter accounts
* Two different HuaweiCloud accounts (sharer and accepter) in the same region for testing the cross-account migration
  functionality

## Variables Introduction

The following variables need to be configured:

### Authentication Variables

* `region_name` - The region where resources will be created (both sharer and accepter accounts must be in the same
  region)
* `access_key` - The access key of the IAM user in sharer account
* `secret_key` - The secret key of the IAM user in sharer account
* `accepter_access_key` - The access key of the IAM user in accepter account
* `accepter_secret_key` - The secret key of the IAM user in accepter account

### Resource Variables

#### Required Variables

* `vpc_name` - The name of the VPC in sharer account
* `subnet_name` - The name of the VPC subnet in sharer account
* `security_group_name` - The name of the security group in sharer account
* `instance_name` - The name of the ECS instance to be created in sharer account
* `administrator_password` - The password of the administrator for the ECS instance
* `data_volume_name` - The name of the data volume to be created and attached to ECS instance in sharer account
* `data_image_name` - The name of the data disk image to be created
* `accepter_vpc_name` - The name of the VPC in accepter account
* `accepter_subnet_name` - The name of the VPC subnet in accepter account
* `accepter_security_group_name` - The name of the security group in accepter account
* `accepter_instance_name` - The name of the ECS instance to be created in accepter account
* `accepter_data_volume_name` - The name of the data volume to be created from shared image in accepter account

#### Optional Variables

* `instance_flavor_id` - The ID of the ECS instance flavor (default: "")
* `instance_flavor_performance_type` - The performance type of the ECS instance flavor (default: "normal")
* `instance_flavor_cpu_core_count` - The CPU core count of the ECS instance flavor (default: 2)
* `instance_flavor_memory_size` - The memory size of the ECS instance flavor in GB (default: 4)
* `instance_image_id` - The ID of the ECS instance image (default: "")
* `instance_image_visibility` - The visibility of the ECS instance image (default: "public")
* `instance_image_os` - The OS of the ECS instance image (default: "Ubuntu")
* `vpc_cidr` - The CIDR block of the VPC in sharer account (default: "192.168.0.0/16")
* `enterprise_project_id` - The ID of the enterprise project (default: null)
* `subnet_cidr` - The CIDR block of the VPC subnet in sharer account (default: "")
* `subnet_gateway_ip` - The gateway IP of the VPC subnet in sharer account (default: "")
* `data_volume_type` - The type of the data volume (default: "SAS")
* `data_volume_size` - The size of the data volume in GB (default: 10)
* `data_image_description` - The description of the data disk image (default: "")
* `accepter_vpc_cidr` - The CIDR block of the VPC in accepter account (default: "192.168.0.0/16")
* `accepter_instance_flavor_id` - The ID of the ECS instance flavor in accepter account (default: "")
* `accepter_instance_image_id` - The ID of the ECS instance image in accepter account (default: "")
* `accepter_data_volume_type` - The type of the data volume in accepter account (default: "SAS")
* `accepter_data_volume_size` - The size of the data volume in accepter account in GB (default: 20)

## Usage

* Copy this example script to your `main.tf`.

* Create a `terraform.tfvars` file and fill in the required variables:

  ```hcl
  vpc_name                     = "tf_test_data_image_vpc"
  subnet_name                  = "tf_test_data_image_subnet"
  security_group_name          = "tf_test_data_image_sg"
  instance_name                = "tf_test_data_image_ecs"
  administrator_password       = "YourPassword@12!"
  data_volume_name             = "tf_test_data_volume"
  data_image_name              = "tf_test_data_image"
  accepter_vpc_name            = "tf_test_accepter_vpc"
  accepter_subnet_name         = "tf_test_accepter_subnet"
  accepter_security_group_name = "tf_test_accepter_sg"
  accepter_instance_name       = "tf_test_accepter_ecs"
  accepter_data_volume_name    = "tf_test_accepter_data_volume"
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
* Both sharer and accepter accounts must be in the same region for image sharing to work
* The data volume must be attached to an ECS instance before creating a data disk image

## Requirements

| Name | Version |
| ---- | ---- |
| terraform | >= 1.1.0 |
| huaweicloud | >= 1.80.1 |
