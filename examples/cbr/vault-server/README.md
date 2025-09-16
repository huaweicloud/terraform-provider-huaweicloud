# Create a CBR Vault for Server Backup

This example provides best practice code for using Terraform to create a vault for server backup in HuaweiCloud.

## Prerequisites

* A HuaweiCloud account
* Terraform installed (>= 1.9.0)
* HuaweiCloud access key and secret key (AK/SK)

## Variable Introduction

The following variables need to be configured:

### Authentication Variables

* `region_name` - The region where the vault is located
* `access_key` - The access key of the IAM user
* `secret_key` - The secret key of the IAM user

### Resource Variables

#### Required Variables

* `vpc_name` - The name of the VPC
* `subnet_name` - The name of the subnet
* `secgroup_name` - The name of the security group
* `ecs_instance_name` - The name of the ECS instance to be backed up
* `vault_name` - The name of the vault
* `vault_size` - The size of the vault in GB

#### Optional Variables

* `vpc_cidr` - The CIDR block of the VPC (default: "192.168.0.0/16")
* `subnet_cidr` - The CIDR block of the subnet (default: "")
  If not specified, a subnet will be calculated within the VPC CIDR
* `subnet_gateway_ip` - The gateway IP of the subnet (default: "")
  If not specified, a gateway IP will be calculated
* `availability_zone` - The availability zone where resources will be created (default: "")
  If not specified, the first available zone will be used
* `instance_flavor_id` - The flavor ID of the ECS instance (default: "")
  If not specified, a flavor will be selected based on other flavor parameters
* `instance_flavor_performance_type` - The performance type of the ECS flavor (default: "normal")
* `instance_flavor_cpu_core_count` - The number of CPU cores for the ECS flavor (default: 2)
* `instance_flavor_memory_size` - The memory size in GB for the ECS flavor (default: 4)
* `instance_image_id` - The ID of the image for the ECS instance (default: "")
  If not specified, an image will be selected based on other image parameters
* `instance_image_os_type` - The OS type of the image (default: "Ubuntu")
* `key_pair_name` - The name of the key pair for ECS login
* `system_disk_type` - The type of the system disk (default: "SAS")
* `system_disk_size` - The size of the system disk in GB (default: 40)
* `protection_type` - The protection type of the vault (default: "backup")
* `consistent_level` - The consistency level (default: "crash_consistent")
* `auto_bind` - Whether to automatically bind the vault to a policy (default: false)
* `auto_expand` - Whether to automatically expand the vault (default: false)
* `is_multi_az` - Whether the vault is deployed across multiple AZs (default: false)
* `enterprise_project_id` - The enterprise project ID (default: "0")
* `backup_name_prefix` - The prefix for backup names (default: "")
* `exclude_volumes` - Whether to exclude volumes from backup (default: false)
* `enable_policy` - Whether to enable backup policy (default: false)
* `tags` - Tags to apply to resources (default: {environment = "test", terraform = "true"})

## Usage

1. Copy this example script to your `main.tf`.

2. Create a `terraform.tfvars` file and fill in the required variables:

   ```hcl
   # Authentication
   region_name = "cn-north-4"
   access_key  = "your_access_key_here"
   secret_key  = "your_secret_key_here"
   
   # Network Configuration
   vpc_name      = "cbr-test-vpc"
   subnet_name   = "cbr-test-subnet"
   secgroup_name = "cbr-test-sg"
   
   # ECS Configuration
   ecs_instance_name = "cbr-test-ecs"
   key_pair_name     = "your_keypair_name_here"
   
   # Vault Configuration
   vault_name = "cbr-vault-server"
   vault_size = 200
   ```

3. Initialize Terraform:

   ```bash
   $ terraform init
   ```

4. Review the Terraform plan:

   ```bash
   $ terraform plan
   ```

5. Apply the configuration:

   ```bash
   $ terraform apply
   ```

6. To clean up the resources:

   ```bash
   $ terraform destroy
   ```

## Note

* Make sure to keep your credentials secure and never commit them to version control
* All resources will be created in the specified region

## Requirements

| Name | Version |
|------|---------|
| terraform | >= 1.9.0 |
| huaweicloud | >= 1.58.0 |
