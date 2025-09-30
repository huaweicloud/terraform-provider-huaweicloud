# Create a scaling group

This example provides best practice code for using Terraform to create an Auto Scaling group with configuration in
HuaweiCloud.

## Prerequisites

* A Huawei Cloud account
* Terraform installed
* Huawei Cloud access key and secret key (AK/SK)

## Required Variables

### Authentication Variables

* `region_name` - The region where the scaling group is located
* `access_key` - The access key of the IAM user
* `secret_key` - The secret key of the IAM user

### Resource Variables

#### Required Variables

* `security_group_name` - The name of the security group
* `keypair_name` - The name of the key pair
* `configuration_name` - The name of the scaling configuration
* `configuration_metadata` - The metadata for the scaling configuration instances
* `configuration_user_data` - The user data script for scaling configuration instances initialization
* `configuration_disks` - The disk configurations for the scaling configuration instances
  - `size` - The size of the disk in GB
  - `volume_type` - The volume type of the disk
  - `disk_type` - The type of the disk
* `configuration_public_eip_settings` - The public IP settings for the scaling configuration instances
  - `ip_type` - The type of the elastic IP
  - `bandwidth` - The bandwidth configuration
    + `size` - The bandwidth size in Mbit/s
    + `share_type` - The bandwidth share type
    + `charging_mode` - The charging mode
* `scaling_group_name` - The name of the scaling group

#### Optional Variables

* `availability_zone` - The availability zone to which the scaling configuration belongs (default: "")
* `configuration_flavor_id` - The flavor ID of the scaling configuration (default: "")
* `configuration_flavor_performance_type` - The performance type of the scaling configuration (default: "normal")
* `configuration_flavor_cpu_core_count` - The CPU core count of the scaling configuration (default: 2)
* `configuration_flavor_memory_size` - The memory size of the scaling configuration (default: 4)
* `configuration_image_id` - The image ID of the scaling configuration (default: "")
* `configuration_image_visibility` - The visibility of the image (default: "public")
* `configuration_image_os` - The OS of the image (default: "Ubuntu")
* `keypair_public_key` - The public key for SSH access (default: "")
* `scaling_group_vpc_id` - The ID of the VPC (default: "")
* `scaling_group_subnet_id` - The ID of the subnet (default: "")
* `scaling_group_vpc_name` - The name of the VPC (default: "")
* `scaling_group_vpc_cidr` - The CIDR block of the VPC (default: "192.168.0.0/16")
* `scaling_group_subnet_name` - The name of the subnet (default: "")
* `scaling_group_subnet_cidr` - The CIDR block of the subnet (default: "")
* `scaling_group_subnet_gateway_ip` - The gateway IP of the subnet (default: "")
* `scaling_group_desire_instance_number` - The desired number of instances (default: 2)
* `scaling_group_min_instance_number` - The minimum number of instances (default: 0)
* `scaling_group_max_instance_number` - The maximum number of instances (default: 10)
* `is_delete_scaling_group_publicip` - Whether to delete the public IP address of the scaling instances when the scaling
  group is deleted (default: true)
* `is_delete_scaling_group_instances` - Whether to delete the scaling instances when the scaling group is deleted
  (default: "yes")

## Usage

* Copy this example script to your `main.tf`.

* Create a `terraform.tfvars` file and fill in the required variables:

  ```hcl
  security_group_name    = "tf_test_secgroup_demo"
  keypair_name           = "tf_test_keypair_demo"
  configuration_name     = "tf_test_as_configuration"
  configuration_metadata = {
    some_key = "some_value"
  }
  configuration_user_data = <<EOT
  #!/bin/sh
  echo "Hello World! The time is now $(date -R)!" | tee /root/output.txt
  EOT

  configuration_disks = [
    {
      size        = 40
      volume_type = "SSD"
      disk_type   = "SYS"
    }
  ]

  configuration_public_eip_settings = [
    {
      ip_type   = "5_bgp"
      bandwidth = {
        size          = 10
        share_type    = "PER"
        charging_mode = "traffic"
      }
    }
  ]

  scaling_group_vpc_name    = "tf_test_vpc_demo"
  scaling_group_subnet_name = "tf_test_subnet_demo"
  scaling_group_name        = "tf_test_scaling_group_demo"
  ```

* Initialize Terraform:

  ```bash
  terraform init
  ```

* Review the Terraform plan:

  ```bash
  terraform plan
  ```

* Apply the configuration:

  ```bash
  terraform apply
  ```

* To clean up the resources:

  ```bash
  terraform destroy
  ```

## Notes

* This example creates a complete Auto Scaling group including:
  - Security group for instance access control
  - Key pair for SSH access
  - Auto Scaling configuration with instance specifications
  - VPC and subnet (if not provided)
  - Auto Scaling group with scaling policies
* The example uses data sources to automatically select appropriate flavors and images if not specified
* Make sure to keep your credentials secure and never commit them to version control
* All resources will be created in the specified region

## Requirements

| Name        | Version   |
|-------------|-----------|
| terraform   | >= 1.9.0  |
| huaweicloud | >= 1.57.0 |
