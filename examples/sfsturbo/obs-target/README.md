# Create an OBS target

This example provides best practice code for using Terraform to create an OBS target in HuaweiCloud.

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
* `target_bucket_name` - The OBS bucket name
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
* `target_file_path` - The linkage directory name of the OBS target (default: "testDir")
* `target_obs_endpoint` - The domain name of the region where the OBS bucket
  located (default: "obs.cn-north-4.myhuaweicloud.com")
* `target_events` - The type of the data automatically exported to the OBS bucket (default: [])
* `target_prefix` - The prefix to be matched in the storage backend (default: "")
* `target_suffix` - The suffix to be matched in the storage backend (default: "")
* `target_file_mode` - The permissions on the imported file (default: "")
* `target_dir_mode` - The permissions on the imported directory (default: "")
* `target_uid` - The ID of the user who owns the imported object (default: 0)
* `target_gid` - The ID of the user group to which the imported object belongs (default: 0)

## Usage

* Copy this example script to your `main.tf`.

* Create a `terraform.tfvars` file and fill in the required variables:

  ```hcl
  vpc_name            = "your_vpc_name"
  subnet_name         = "your_subnet_name"
  security_group_name = "your_security_group_name"
  target_bucket_name  = "your_target_bucket_name"
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
| huaweicloud | >= 1.75.0 |
