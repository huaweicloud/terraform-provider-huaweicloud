# Create EVS volume and v5 snapshot, and attach the volume to an ECS instance

This example provides best practice code for using Terraform to create EVS volume and v5 snapshot,
and attach the volume to an ECS instance in HuaweiCloud.

## Prerequisites

* A HuaweiCloud account
* Terraform installed
* HuaweiCloud access key and secret key (AK/SK)

## Required Variables

The following variables need to be configured:

### Authentication Variables

* `region_name` - The region where the resources are located
* `access_key`  - The access key of the IAM user
* `secret_key`  - The secret key of the IAM user

### Resource Variables

#### Required Variables

* `vpc_name` - The name of the VPC
* `subnet_name` - The name of the subnet
* `secgroup_name` - The name of the security group
* `ecs_instance_name` - The name of the ECS instance
* `volume_name` - The name of the volume
* `snapshot_name` - The name of the snapshot

#### Optional Variables

* `vpc_cidr` - The CIDR block of the VPC (default: "192.168.0.0/16")
* `instance_flavor_id` - The flavor ID of the ECS instance (default: "")
* `availability_zone` - The availability zone for the resources (default: "")
* `instance_flavor_performance_type` - The performance type of the ECS instance flavor (default: "normal")
* `instance_flavor_cpu_core_count` - The number of CPU cores for the ECS instance flavor (default: 0)
* `instance_flavor_memory_size` - The memory size in GB for the ECS instance flavor (default: 0)
* `instance_image_id` - The ID of the image used to create the ECS instance (default: "")
* `instance_image_os_type` - The OS type of the ECS instance image (default: "Ubuntu")
* `instance_image_visibility` - The visibility of the ECS instance image (default: "public")
* `key_pair_name` - The name of the key pair for ECS login (default: "")
* `system_disk_type` - The type of the system disk (default: "SAS")
* `system_disk_size` - The size of the system disk in GB (default: 40)
* `subnet_cidr` - The CIDR block of the subnet (default: "")
* `subnet_gateway_ip` - The gateway IP of the subnet (default: "")
* `instant_access` - Whether to enable instant access for the snapshot group (default: false)
* `volume_availability_zone` - The availability zone for the volume (default: "")
* `volume_type` - The type of the volume, the value can be **SAS**, **SSD**, **GPSSD**, **ESSD**, **GPSSD2** or
  **ESSD2** (default: "SSD")
* `voulme_size` - The size of the volume (default: 40)
* `volume_description` - The description of the volume (default: "")
* `volume_multiattach` - The volume is shared volume or not (default: false)
* `volume_iops` - The IOPS for the volume, only valid and required when `volume_type` is set to
  **GPSSD2** or **ESSD2** (default: null)
* `volume_throughput` - The throughput for the volume, only valid and required when `volume_type` is set to
  **GPSSD2** (default: null)
* `volume_device_type` - The device type of disk to create (default: "VBD")
* `enterprise_project_id` - The enterprise project ID (default: null)
* `volume_tags` - The tags of the volume (default: {})
* `charging_mode` - The charging mode of the volume (default: "postPaid")
* `period_unit` - The period unit of the volume, only required when `charging_mode` is **prePaid** (default: null)
* `period` - The period of the volume, only required when `charging_mode` is **prePaid** (default: null)
* `auto_renew` - Whether auto-renew is enabled for the volume (default: "false")
* `volume_metadata` - The volume metadata (default: {})
* `snapshot_description` - The description of the snapshot (default: "")
* `snapshot_tags` - The tags of the snapshot (default: {})
* `snapshot_instant_access` - Whether to enable instant access for the snapshot (default: false)
* `snapshot_incremental` - Whether to create an incremental snapshot (default: false)
* `snapshot_metadata` - The snapshot metadata (default: {})

## Usage

* Copy this example script to your `main.tf`.

* Create a `terraform.tfvars` file and fill in the required variables:

  ```hcl
  vpc_name          = "your_vpc_name"
  subnet_name       = "your_subnet_name"
  secgroup_name     = "your_secgroup_name"
  ecs_instance_name = "your_ecs_instance_name"
  volume_name       = "your_volume_name"
  snapshot_name     = "your_snapshot_name"
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

| Name | Version   |
| ---- |-----------|
| terraform | >= 1.1.0  |
| huaweicloud | >= 1.77.3 |
