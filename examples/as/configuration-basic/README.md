# Create an Auto Scaling Configuration

This example provides best practice code for using Terraform to create an Auto Scaling configuration in HuaweiCloud.

## Prerequisites

* A HuaweiCloud account
* Terraform installed
* HuaweiCloud access key and secret key (AK/SK)

## Required Variables

The following variables need to be configured:

### Authentication Variables

* `region_name` - The region where the AS configuration is located
* `access_key` - The access key of the IAM user
* `secret_key` - The secret key of the IAM user

### Resource Variables

#### Required Variables

* `security_group_name` - The name of the security group
* `keypair_name` - The name of the key pair
* `configuration_name` - The name of the AS configuration
* `configuration_metadata` - The metadata for the instance
* `configuration_user_data` - The user data script for instance initialization
* `configuration_disks` - The disk configurations for the instance
  - `size` - The size of the disk in GB
  - `volume_type` - The volume type of the disk
  - `disk_type` - The type of the disk
* `configuration_public_ip` - The public IP configuration for the instance
  - `eip` - The EIP configuration
    + `ip_type` - The type of the elastic IP
    + `bandwidth` - The bandwidth configuration
      - `size` - The bandwidth size in Mbit/s
      - `share_type` - The bandwidth share type
      - `charging_mode` - The charging mode

#### Optional Variables

* `availability_zone` - The availability zone to which the AS configuration belongs (default: "")
* `configuration_flavor_id` - The flavor ID of the AS configuration (default: "")
* `configuration_flavor_performance_type` - The performance type of the AS configuration (default: "normal")
* `configuration_flavor_cpu_core_count` - The CPU core count of the AS configuration (default: 2)
* `configuration_flavor_memory_size` - The memory size of the AS configuration (default: 4)
* `configuration_flavor_image_id` - The image ID of the AS configuration (default: "")
* `configuration_flavor_image_visibility` - The visibility of the image (default: "public")
* `configuration_flavor_image_os` - The OS of the image (default: "Ubuntu")
* `keypair_public_key` - The public key for SSH access (default: "")

## Usage

* Copy this example script to your `main.tf`.

* Create a `terraform.tfvars` file and fill in the required variables:

  ```hcl
  security_group_name     = "your_security_group_name"
  keypair_name            = "your_keypair_name"
  configuration_name      = "your_configuration_name"
  configuration_metadata  = "your_configuration_metadata"
  configuration_user_data = "your_configuration_user_data"
  configuration_disks     = "your_configuration_disks"
  configuration_public_ip = "your_configuration_public_ip"
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

* The Auto Scaling configuration will be created with the specified instance configuration.
* Make sure to keep your credentials secure and never commit them to version control.
* All resources will be created in the specified region.

## Requirements

| Name | Version |
|------|---------|
| terraform | >= 0.12.0 |
| huaweicloud | >= 1.57.0 |
