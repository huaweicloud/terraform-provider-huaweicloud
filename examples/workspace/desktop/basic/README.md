# Create a cloud desktop machine

This example provides best practice code for using Terraform to create a simple cloud desktop machine in HuaweiCloud
Workspace service.

## Prerequisites

* A HuaweiCloud account
* Terraform installed
* HuaweiCloud access key and secret key (AK/SK)

## Required Variables

The following variables need to be configured:

### Authentication Variables

* `region_name` - The region where the Workspace desktop is located
* `access_key` - The access key of the IAM user
* `secret_key` - The secret key of the IAM user

### Resource Variables

#### Required Variables

* `vpc_name` - The VPC name
* `vpc_cidr` - The CIDR block of the VPC
* `subnet_name` - The subnet name
* `security_group_name` - The security group name
* `desktop_user_name` - The user name that the cloud desktop used
* `desktop_user_email` - The email address that the user used
* `cloud_desktop_name` - The cloud desktop name

#### Optional Variables

* `availability_zone` - The availability zone to which the cloud desktop flavor and network belong (default: "")
* `desktop_flavor_id` - The flavor ID of the cloud desktop (default: "")
* `desktop_flavor_os_type` - The OS type of the cloud desktop flavor (default: "Windows")
* `desktop_flavor_cpu_core_number` - The number of the cloud desktop flavor CPU cores (default: 4)
* `desktop_flavor_memory_size` - The number of the cloud desktop flavor memories (default: 8)
* `desktop_image_id` - The specified image ID that the cloud desktop used (default: "")
* `desktop_image_os_type` - The OS type of the cloud desktop image (default: "Windows")
* `desktop_image_visibility` - The visibility of the cloud desktop image (default: "market")
* `subnet_cidr` - The CIDR block of the subnet (default: "")
* `subnet_gateway_ip` - The gateway IP of the subnet (default: "")
* `desktop_user_group_name` - The name of the user group that cloud desktop used (default: "users")
* `desktop_root_volume_type` - The storage type of system disk (default: "SSD")
* `desktop_root_volume_size` - The storage capacity of system disk (default: 100)
* `desktop_data_volumes` - The storage configuration of data disks (default: [{type="SSD", size=100}])

## Usage

* Copy this example script to your `main.tf`.

* Create a `terraform.tfvars` file and fill in the required variables:

  ```hcl
  vpc_name            = "your_vpc_name"
  vpc_cidr            = "192.168.0.0/16"
  subnet_name         = "your_subnet_name"
  security_group_name = "your_security_group_name"
  desktop_user_name   = "your_desktop_user_name"
  desktop_user_email  = "your_desktop_user_email"
  cloud_desktop_name  = "your_cloud_desktop_name"
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
* The cloud desktop machine and the user are dependent on the service configuration
* Please read the implicit and explicit dependencies in the script carefully
* All resources will be created in the specified region

## Requirements

| Name | Version |
| ---- | ---- |
| terraform | >= 0.12.0 |
| huaweicloud | >= 1.67.0 |
