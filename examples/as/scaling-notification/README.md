# Create an AS notification

This example provides best practice code for using Terraform to create an Auto Scaling (AS) notification in HuaweiCloud.

## Prerequisites

* A HuaweiCloud account
* Terraform installed
* HuaweiCloud access key and secret key (AK/SK)
* AS service enabled in the target region
* VPC service enabled in the target region
* ECS service enabled in the target region

## Variable Introduction

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
* `keypair_name` - The name of the key pair that is used to access the AS instance
* `topic_name` - The name of the SMN topic
* `configuration_name` - The name of the AS configuration
* `disk_configurations` - The disk configurations for the AS instance (must include exactly one system disk)
  + `disk_type` - The type of the disk (SYS for system disk, DATA for data disk)
  + `volume_type` - The type of the volume (SSD, SATA, SAS)
  + `volume_size` - The size of the volume in GB
* `group_name` - The name of the AS group

#### Optional Variables

* `instance_flavor_id` - The flavor ID of the AS instance (default: "")
* `instance_flavor_performance_type` - The performance type of the AS instance flavor (default: "normal")
* `instance_flavor_cpu_core_count` - The CPU core count of the AS instance flavor (default: 2)
* `instance_flavor_memory_size` - The memory size of the AS instance flavor (default: 4)
* `instance_image_id` - The image ID of the AS instance (default: "")
* `image_visibility` - The visibility of the image (default: "public")
* `image_os` - The image OS type (default: "Ubuntu")
* `vpc_cidr` - The CIDR block of the VPC (default: "192.168.0.0/16")
* `subnet_cidr` - The CIDR block of the subnet (default: "")
* `subnet_gateway_ip` - The gateway IP of the subnet (default: "")
* `keypair_public_key` - The public key of the key pair (default: "")
* `desire_instance_number` - The desired number of scaling instances in the AS group (default: 1)
* `min_instance_number` - The minimum number of scaling instances in the AS group (default: 0)
* `max_instance_number` - The maximum number of scaling instances in the AS group (default: 5)
* `is_delete_publicip` - Whether to delete the public IP address when the AS group is deleted (default: true)
* `is_delete_instances` - Whether to delete the scaling instances when the AS group is deleted (default: true)
* `group_cool_down_time` - The cooling duration in seconds (default: 300)
* `group_description` - The description of the AS group (default: "")
* `enterprise_project_id` - The enterprise project ID (default: null)
* `group_tags` - The tags of the AS group (default: {})
* `notification_events` - The events of the AS notification (default: ["SCALING_UP", "SCALING_DOWN"])
  The events include `SCALING_UP`, `SCALING_UP_FAIL`, `SCALING_DOWN`, `SCALING_DOWN_FAIL`, `SCALING_GROUP_ABNORMAL`.

## Usage

* Copy this example script to your `main.tf`.

* Create a `terraform.tfvars` file and fill in the required variables:

  ```hcl
  vpc_name            = "tf_test_vpc"
  subnet_name         = "tf_test_subnet"
  security_group_name = "tf_test_security_group"
  keypair_name        = "tf_test_keypair"
  topic_name          = "tf_test_topic"
  configuration_name  = "tf_test_configuration"

  disk_configurations = [
    {
      disk_type   = "SYS"
      volume_type = "SSD"
      volume_size = 40
    }
  ]

  group_name = "tf_test_group"
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

## Architecture

This example creates the following resources:

1. **VPC and Subnet** - Network infrastructure for the AS group
2. **Security Group** - Security rules for the instances
3. **Key Pair** - SSH key for instance access
4. **AS Configuration** - Instance configuration template
5. **AS Group** - Auto Scaling group with desired, min, and max instance counts
6. **AS notification** - The AS notification

## Notes

* Make sure to keep your credentials secure and never commit them to version control
* The AS group is dependent on the VPC, subnet, security group, and AS configuration
* The example uses Ubuntu public images by default, but you can specify custom images
* Disk configurations must include exactly one system disk (disk_type = "SYS")

## Requirements

| Name | Version |
| ---- | ---- |
| terraform | >= 1.9.0 |
| huaweicloud | >= 1.57.0 |
