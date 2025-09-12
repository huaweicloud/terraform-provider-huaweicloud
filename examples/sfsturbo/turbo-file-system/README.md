# Create a SFS Turbo file system

This example provides best practice code for using Terraform to create a SFS Turbo file system in HuaweiCloud.

## Prerequisites

* A HuaweiCloud account
* Terraform installed
* HuaweiCloud access key and secret key (AK/SK)

## Required Variables

The following variables need to be configured:

### Authentication Variables

* `region_name` - The region where the SFS Turbo service is located
* `access_key`  - The access key of the IAM user
* `secret_key`  - The secret key of the IAM user

### Resource Variables

#### Required Variables

* `vpc_name` - The VPC name
* `subnet_name` - The subnet name
* `security_group_name` - The security group name
* `turbo_name` - The SFS Turbo file system name

#### Optional Variables

* `vpc_cidr` - The CIDR block of the VPC (default: "192.168.0.0/16")
* `subnet_cidr` - The CIDR block of the subnet (default: "")
* `subnet_gateway_ip` - The gateway IP of the subnet (default: "")
* `turbo_size` - The capacity of the SFS Turbo file system (default: 1228)
* `turbo_share_proto` - The protocol of the SFS Turbo file system (default: "NFS")
* `turbo_share_type` - The type of the SFS Turbo file system (default: "STANDARD")
* `turbo_hpc_bandwidth` - The bandwidth specification of the SFS Turbo file system, only required when
  `turbo_share_type` is `HPC` (default: "")
* `turbo_backup_id` - The ID of the backup, only required when a file system is created from a backup (default: "")
* `turbo_tags` - The tags of the SFS Turbo file system (default: {})
* `enterprise_project_id` - The enterprise project ID of the SFS Turbo file system (default: null)
* `charging_mode` - The charging mode of SFS Turbo file system (default: "postPaid")
* `period_unit` - The period unit of the SFS Turbo file system, only required when
  `charging_mode` is `prePaid` (default: null)
* `period` - The period of the SFS Turbo file system, only required when `charging_mode` is `prePaid` (default: null)
* `auto_renew` - The auto renew of the SFS Turbo file system (default: "false")

## Usage

* Copy this example script to your `main.tf`.

* Create a `terraform.tfvars` file and fill in the required variables:

  ```hcl
  vpc_name            = "your_vpc_name"
  subnet_name         = "your_subnet_name"
  security_group_name = "your_security_group_name"
  turbo_name          = "your_turbo_name"
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

| Name | Version |
| ---- | ---- |
| terraform | >= 1.1.0 |
| huaweicloud | >= 1.67.0 |
