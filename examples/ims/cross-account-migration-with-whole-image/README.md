# Cross-Account Migration with Whole Image

This example provides best practice code for using Terraform to migrate business data across accounts by creating
a whole image, sharing it, and creating a new ECS instance from the shared image in HuaweiCloud IMS service.

## Prerequisites

* Terraform installed
* HuaweiCloud access key and secret key (AK/SK) for both sharer and accepter accounts
* Two different HuaweiCloud accounts (sharer and accepter) for testing the cross-account migration functionality

## Variables Introduction

The following variables need to be configured:

### Authentication Variables

* `region_name` - The region where resources will be created
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
* `vault_name` - The name of the CBR vault in sharer account
* `whole_image_name` - The name of the whole image to be created
* `accepter_project_ids` - The project IDs of accepter account for image sharing
* `accepter_vault_name` - The name of the CBR vault in accepter account
* `accepter_instance_name` - The name of the new ECS instance to be created in accepter account

#### Optional Variables

* `instance_flavor_id` - The ID of the ECS instance flavor (default: "")
* `instance_flavor_performance_type` - The performance type of the ECS instance flavor (default: "normal")
* `instance_flavor_cpu_core_count` - The CPU core count of the ECS instance flavor (default: 2)
* `instance_flavor_memory_size` - The memory size of the ECS instance flavor in GB (default: 4)
* `instance_image_id` - The ID of the ECS instance image (default: "")
* `instance_image_visibility` - The visibility of the ECS instance image (default: "public")
* `instance_image_os` - The OS of the ECS instance image (default: "Ubuntu")
* `vpc_cidr` - The CIDR block of the VPC in sharer account (default: "192.168.0.0/16")
* `subnet_cidr` - The CIDR block of the VPC subnet in sharer account (default: "")
* `subnet_gateway_ip` - The gateway IP of the VPC subnet in sharer account (default: "")
* `instance_data_disks` - The data disks of the ECS instance (default: [])
  + `type` - The type of the data disk
  + `size` - The size of the data disk in GB
* `vault_type` - The type of the CBR vault (default: "server")
* `vault_consistent_level` - The consistent level of the CBR vault (default: "crash_consistent")
* `vault_protection_type` - The protection type of the CBR vault (default: "backup")
* `vault_size` - The size of the CBR vault in GB (default: 200)
* `whole_image_description` - The description of the whole image (default: "")
* `accepter_vault_type` - The type of the CBR vault in accepter account (default: "server")
* `accepter_vault_consistent_level` - The consistent level of the CBR vault in accepter account (default: "crash_consistent")
* `accepter_vault_protection_type` - The protection type of the CBR vault in accepter account (default: "backup")
* `accepter_vault_size` - The size of the CBR vault in accepter account in GB (default: 200)

## Usage

* Copy this example script to your `main.tf`.

* Create a `terraform.tfvars` file and fill in the required variables:

  ```hcl
  accepter_access_key    = "your_accepter_access_key"
  accepter_secret_key    = "your_accepter_secret_key"
  vpc_name               = "tf_test_whole_image_vpc"
  subnet_name            = "tf_test_whole_image_subnet"
  security_group_name    = "tf_test_whole_image_sg"
  instance_name          = "tf_test_whole_image_ecs"
  administrator_password = "YourPassword@12!"
  vault_name             = "tf_test_sharer_vault"
  whole_image_name       = "tf_test_sharer_whole_image"
  accepter_project_ids   = ["your_accepter_project_id"]
  accepter_vault_name    = "tf_test_accepter_vault"
  accepter_instance_name = "tf_test_accepter_instance"
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
* Ensure the CBR vault in accepter account has sufficient capacity to store the whole image
* Only whole images created through CBR (Cloud Backup and Recovery) or whole images created from ECS instances that
  have not generated backups through the old CSBS service support sharing

## Requirements

| Name | Version |
| ---- | ---- |
| terraform | >= 1.1.0 |
| huaweicloud | >= 1.68.0 |
