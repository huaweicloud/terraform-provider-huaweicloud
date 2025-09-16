# Create a CBR Vault for SFS Turbo Backup

This example provides best practice code for using Terraform to create a vault for SFS Turbo file system backup in HuaweiCloud.

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
* `turbo_name` - The name of the standard SFS Turbo file system
* `turbo_size` - The size of the standard SFS Turbo file system in GB (minimum: 500GB)
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
* `protection_type` - The protection type of the vault (default: "backup")
* `auto_expand` - Whether to automatically expand the vault (default: false)
* `is_multi_az` - Whether the vault is deployed across multiple AZs (default: false)
* `enterprise_project_id` - The enterprise project ID (default: "0")
* `backup_name_prefix` - The prefix for backup names (default: "")
* `enable_policy` - Whether to enable backup policy (default: false)
* `tags` - Tags to apply to resources (default: {environment = "test", terraform = "true", service = "sfs-turbo"})

## Usage

1. Copy this example script to your `main.tf`.

2. Create a `terraform.tfvars` file and fill in the required variables:

   ```hcl
   # Authentication
   region_name = "your_region_name"
   access_key  = "your_access_key_here"
   secret_key  = "your_secret_key_here"
   
   # Network Configuration
   vpc_name      = "your_vpc_name"
   subnet_name   = "your_subnet_name"
   secgroup_name = "your_secgroup_name"
   
   # SFS Turbo Configuration (Standard)
   turbo_name = "your_turbo_name"
   turbo_size = 500
   
   # Vault Configuration
   vault_name = "your_vault_name"
   vault_size = 1000
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

## Notes

* The minimum size for SFS Turbo is 500GB
* Ensure the selected region supports SFS Turbo service
* Make sure to keep your credentials secure and never commit them to version control
* All resources will be created in the specified region

## Requirements

| Name | Version |
|------|---------|
| terraform | >= 1.9.0 |
| huaweicloud | >= 1.58.0 |
