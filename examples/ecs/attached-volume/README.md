# Create an ECS instance with attached volume

This example provides best practice code for using Terraform to create an ECS instance with an attached data volume in HuaweiCloud.

## Prerequisites

* A HuaweiCloud account
* Terraform installed
* HuaweiCloud access key and secret key (AK/SK)

## Variable Introduction

The following variables need to be configured:

### Authentication Variables

* `region_name` - The region where the ECS instance is located
* `access_key` - The access key of the IAM user
* `secret_key` - The secret key of the IAM user

### Resource Variables

#### Required Variables

* `vpc_name` - The name of the VPC
* `subnet_name` - The name of the subnet
* `instance_name` - The name of the ECS instance
* `instance_admin_password` - The login password of the ECS instance
* `volume_name` - The name of the data volume

#### Optional Variables

* `availability_zone` - The availability zone to which the ECS instance belongs (default: "")
* `instance_flavor_id` - The flavor ID of the ECS instance (default: "")
* `instance_performance_type` - The performance type of the ECS instance flavor (default: "normal")
* `instance_cpu_core_count` - The number of CPU cores of the ECS instance (default: 2)
* `instance_memory_size` - The memory size in GB of the ECS instance (default: 4)
* `instance_image_id` - The image ID of the ECS instance (default: "")
* `instance_image_visibility` - The visibility of the ECS instance image (default: "public")
* `instance_image_os` - The operating system of the ECS instance image (default: "Ubuntu")
* `vpc_cidr` - The CIDR block of the VPC (default: "192.168.0.0/16")
* `subnet_cidr` - The CIDR block of the subnet (default: "")
* `subnet_gateway_ip` - The gateway IP address of the subnet (default: "")
* `security_group_name` - The name of the security group (default: "")
* `security_group_ids` - The list of security group IDs of the ECS instance (default: [])
* `enterprise_project_id` - The ID of the enterprise project (only available for enterprise users, default: null)
* `volume_type` - The type of the data volume (default: "SSD")
* `volume_size` - The size of the data volume in GB (default: 10)
* `volume_iops` - The IOPS(Input/Output Operations Per Second) for the data volume (required if `volume_type` is
  **GPSSD2** or **ESSD2**, default: null)
* `volume_throughput` - The throughput for the data volume (required if `volume_type` is **GPSSD2**,default: null)
* `volume_backup_id` - The backup ID from which to create the disk (default: null)
* `volume_snapshot_id` - The snapshot ID from which to create the disk (default: null)

## Usage

* Copy this example script to your `main.tf`.

* Create a `terraform.tfvars` file and fill in the required variables:

  ```hcl
  vpc_name                = "your_vpc_name"
  subnet_name             = "your_subnet_name"
  security_group_name     = "your_security_group_name"
  instance_name           = "your_instance_name"
  instance_admin_password = "your_instance_password"
  volume_name             = "your_volume_name"
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
* The creation of the ECS instance and volume may take several minutes
* All resources will be created in the specified region

## Requirements

| Name         | Version     |
| ------------ | ----------- |
| terraform | >= 1.9.0   |
| huaweicloud | >= 1.60.1  |
