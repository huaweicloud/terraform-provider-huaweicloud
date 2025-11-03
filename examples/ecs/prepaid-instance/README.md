# Create an ECS instance in prePaid charging mode

This example provides best practice code for using Terraform to create an ECS instance in prepaid charging mode in HuaweiCloud.

## Prerequisites

* A HuaweiCloud account
* Terraform installed
* HuaweiCloud access key and secret key (AK/SK)

## Variable Introduction

The following variables need to be configured:

### Authentication Variables

* `region_name` - The region where the ECS instance is located
* `access_key`  - The access key of the IAM user
* `secret_key`  - The secret key of the IAM user

### Resource Variables

#### Required Variables

* `vpc_name` - The name of the VPC
* `subnet_name` - The name of the subnet
* `instance_name` - The name of the ECS instance
* `instance_admin_password` - The login password of the ECS instance

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
* `security_group_ids` - The list of security group IDs of the ECS instance (either `security_group_ids` or
  `security_group_name` must be provided, default: [])
* `security_group_name` - The name of the security group (default: "")
* `instance_description` - The description of the ECS instance (default: null)
* `instance_system_disk_type` - The type of the system disk of the ECS instance (default: null)
* `instance_system_disk_size` - The size of the system disk in GB of the ECS instance (default: null)
* `instance_system_disk_iops` - The IOPS of the system disk of the ECS instance (required if
  `instance_system_disk_type` is **GPSSD2** or **ESSD2**, default: null)
* `instance_system_disk_throughput` - The throughput of the system disk of the ECS instance (required if
  `instance_system_disk_type` is **GPSSD2**, default: null)
* `instance_system_disk_dss_pool_id` - The DSS pool ID of the system disk of the ECS instance (default: null)
* `instance_metadata` - The metadata key/value pairs of the ECS instance (default: null)
* `instance_tags` - The key/value pairs to associate with the instance (default: null)
* `enterprise_project_id` - The ID of the enterprise project to which the ECS instance belongs(only available for enterprise
  users, default: null)
* `instance_eip_id` - The EIP ID to associate with the ECS instance (default: null)
* `instance_eip_type` - The EIP type to create and associate with the ECS instance (default: null)
* `instance_bandwidth` - The bandwidth configuration of the ECS instance (default: null)
  - `share_type` - The share type of the bandwidth (required)
  - `id` - The ID of the bandwidth (default: null)
  - `size` - The size of the bandwidth (default: null)
  - `charge_mode` - The charge mode of the bandwidth (default: null)
  - `extend_param` - Additional parameters for the bandwidth (default: null)
* `period_unit` - The unit of the period of the ECS instance (default: "month")
* `period` - The charging period of the ECS instance (default: 1)
* `auto_renew` - Whether to enable auto-renewal of the ECS instance (default: "false")

## Usage

* Copy this example script to your `main.tf`.

* Create a `terraform.tfvars` file and fill in the required variables:

  ```hcl
  vpc_name                = "your_vpc_name"
  subnet_name             = "your_subnet_name"
  instance_name           = "your_instance_name"
  instance_admin_password = "your_admin_password"
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
* The ECS instance will be created in prepaid charging mode with the specified period
* All resources will be created in the specified region
* If no security group IDs are specified, a default security group will be created

## Requirements

| Name | Version |
| ---- | ---- |
| terraform | >= 1.9.0 |
| huaweicloud | >= 1.64.3 |
