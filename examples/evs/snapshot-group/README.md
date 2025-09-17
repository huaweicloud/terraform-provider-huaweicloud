# Create an EVS Snapshot Group

This example provides best practice code for using Terraform to create an EVS snapshot group in HuaweiCloud.

## Prerequisites

* A HuaweiCloud account
* Terraform installed
* HuaweiCloud access key and secret key (AK/SK)

## Variables Introdution

The following variables need to be configured:

### Authentication Variables

* `region_name` - The region where the resources are located （Currently, only `cn-south-4` is supported）
* `access_key` - The access key of the IAM user
* `secret_key` - The secret key of the IAM user

### Resource Variables

#### Required Variables

* `vpc_name` - The name of the VPC
* `subnet_name` - The name of the subnet
* `secgroup_name` - The name of the security group
* `ecs_instance_name` - The name of the ECS instance

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
* `snapshot_group_name` - The name of the snapshot group (default: "")
* `snapshot_group_description` - The description of the snapshot group (default: "Created by Terraform")
* `enterprise_project_id` - The enterprise project ID for the snapshot group (default: "0")
* `incremental` - Whether to create an incremental snapshot (default: false)
* `tags` - The key/value pairs to associate with the snapshot group (default: { environment = "test",
  created_by = "terraform" })
* `volume_configuration` - The list of volume configurations to attach to the ECS instance (default: [])
  - `name` - The name of the volume
  - `size` - The size of the volume in GB
  - `volume_type` - The type of the volume
  - `device_type` - The device type of the volume

## Usage

* Copy this example script to your `main.tf`.

* Create a `terraform.tfvars` file and fill in the required variables:

  ```hcl
  vpc_name          = "your_vpc_name"
  subnet_name       = "your_subnet_name"
  secgroup_name     = "your_secgroup_name"
  ecs_instance_name = "your_ecs_instance_name"

  volume_configuration = [
    {
      name        = "volume-1"
      size        = 100
      volume_type = "SSD"
      device_type = "SCSI"
    }
  ]
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
* Currently, only `cn-south-4` is supported

## Requirements

| Name | Version |
| ---- | ---- |
| terraform | >= 0.14.0 |
| huaweicloud | >= 1.77.3 |
