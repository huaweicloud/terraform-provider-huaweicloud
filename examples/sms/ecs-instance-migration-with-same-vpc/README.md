# Migrate ECS instances within the same VPC

This example provides best practice code for using Terraform to migrate ECS instances within the same VPC in
HuaweiCloud SMS service.

## Prerequisites

* A HuaweiCloud account
* Terraform installed
* HuaweiCloud access key and secret key (AK/SK)

## Variables Introduction

The following variables need to be configured:

### Authentication Variables

* `region_name` - The region where the resources are located
* `access_key` - The access key of the IAM user
* `secret_key` - The secret key of the IAM user

### Resource Variables

#### Required Variables

* `vpc_name` - The name of the VPC
* `subnet_name` - The name of the subnet
* `security_group_name` - The name of the security group
* `instance_name` - The name of the ECS instance
* `instance_admin_password` - The login password of the ECS instance
* `destination_instance_name` - The name of the destination ECS instance
* `source_server_name` - The name of the SMS source server
* `source_server_os_version` - The OS version of the SMS source server
* `source_server_agent_version` - The agent version of the SMS source server
* `migrate_task_type` - The type of the SMS migration task

#### Optional Variables

* `availability_zone` - The name of the availability zone to which the resources belong (default: "")
* `vpc_cidr` - The CIDR block of the VPC (default: "172.16.0.0/16")
* `subnet_cidr` - The CIDR block of the subnet (default: "")
* `subnet_gateway_ip` - The gateway IP address of the subnet (default: "")
* `instance_flavor_id` - The flavor ID of the ECS instance (default: "")
* `instance_flavor_performance_type` - The performance type of the instance flavor (default: "normal")
* `instance_flavor_cpu_core_count` - The CPU core count of the instance flavor (default: 2)
* `instance_flavor_memory_size` - The memory size of the instance flavor (default: 4)
* `instance_image_id` - The image ID of the instance (default: "")
* `instance_image_visibility` - The visibility of the instance image (default: "public")
* `instance_image_os` - The OS of the instance image (default: "Ubuntu")
* `instance_system_disk_size` - The size of the ECS instance system disk in GB (default: 40)
* `instance_system_disk_type` - The type of the ECS instance system disk (default: "GPSSD")
* `destination_instance_system_disk_size` - The size of the destination ECS instance system disk in GB (default: 80)
* `destination_instance_system_disk_type` - The type of the destination ECS instance system disk (default: "GPSSD")
* `source_server_firmware` - The firmware of the SMS source server (default: "BIOS")
* `source_server_boot_loader` - The boot loader of the SMS source server (default: "GRUB")
* `source_server_has_rsync` - Whether the SMS source server has rsync (default: true)
* `source_server_paravirtualization` - Whether the SMS source server is paravirtualization (default: true)
* `source_server_cpu_quantity` - The CPU quantity of the SMS source server (default: 2)
* `source_server_memory` - The memory of the SMS source server (default: 4018196480)
* `source_server_disks` - The disks of the SMS source server (default: null)
  + `name` - The name of the disk (required)
  + `device_use` - The device use of the disk (required)
  + `size` - The size of the disk in bytes (required)
  + `used_size` - The used size of the disk in bytes (required)
  + `partition_style` - The partition style of the disk (optional)
  + `relation_name` - The relation name of the disk (optional)
  + `inode_size` - The inode size of the disk (optional)
  + `physical_volumes` - The physical volumes of the disk (required, list)
    * `name` - The name of the physical volume (optional)
    * `device_use` - The device use of the physical volume (optional)
    * `file_system` - The file system of the physical volume (optional)
    * `mount_point` - The mount point of the physical volume (optional)
    * `size` - The size of the physical volume in bytes (optional)
    * `used_size` - The used size of the physical volume in bytes (optional)
* `task_auto_start` - Whether to automatically start the SMS task (default: false)
* `task_action` - The action of the SMS task (default: null)
* `task_target_server_disks` - The disks of the SMS task target server (default: null)
  + `name` - The name of the disk (required)
  + `size` - The size of the disk in bytes (required)
  + `device_type` - The device type of the disk (required)
  + `physical_volumes` - The physical volumes of the disk (optional, list)
    * `name` - The name of the physical volume (required)
    * `size` - The size of the physical volume in bytes (required)
    * `device_type` - The device type of the physical volume (required)
    * `file_system` - The file system of the physical volume (required)
    * `mount_point` - The mount point of the physical volume (required)
    * `volume_index` - The volume index of the physical volume (required)
* `task_configurations` - The configurations of the SMS task (default: [])
  + `config_key` - The key of the configuration (required)
  + `config_value` - The value of the configuration (required)

## Usage

* Copy this example script to your `main.tf`.

* Create a `terraform.tfvars` file and fill in the required variables:

  ```hcl
  vpc_name                    = "tf_test_vpc"
  subnet_name                 = "tf_test_subnet"
  security_group_name         = "tf_test_security_group"
  instance_name               = "tf_test_source_server"
  instance_admin_password     = "YourPassword123!"
  destination_instance_name   = "tf_test_destination_server"
  source_server_name          = "tf_test_source_server"
  source_server_os_version    = "UBUNTU_24_4_64BIT"
  source_server_agent_version = "25.2.0"
  migrate_task_type           = "MIGRATE_FILE"
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
| terraform | >= 0.12.0 |
| huaweicloud | >= 1.37.0 |
