# Build Cross-Region IPv4 Network with G-EIP

This example provides best practice code for using Terraform to build a cross-region IPv4 network with Global Elastic IP
(G-EIP) and ECS in HuaweiCloud EIP service.

The example creates a VPC and ECS, creates a VPC internet gateway, purchases a G-EIP with global internet bandwidth,
and associates the G-EIP with the ECS, global internet gateway and global connection bandwidth.

## Prerequisites

* A HuaweiCloud account with EIP, VPC, and ECS permissions
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
* `security_group_name` - The name of the security group
* `security_group_rule_configurations` - The list of security group rule configurations
  - `direction` - The direction of the rule (default: "ingress")
  - `ethertype` - The IP version (default: "IPv4")
  - `protocol` - The protocol of the rule (default: `null`)
  - `ports` - The ports of the rule (default: `null`)
  - `remote_ip_prefix` - The remote IP prefix (default: "0.0.0.0/0")
* `instance_name` - The name of the ECS instance
* `instance_administrator_password` - The administrator password of the ECS instance
* `internet_gateway_name` - The name of the VPC internet gateway
* `global_eip_name` - The name of the global EIP
* `gc_bandwidth_name` - The name of the global connection bandwidth

#### Optional Variables

* `instance_flavor_id` - The ID of the ECS instance flavor (default: `""`)
* `instance_flavor_performance_type` - The performance type of the ECS instance flavor (default: "normal")
* `instance_flavor_cpu_core_count` - The CPU core count of the ECS instance flavor (default: 2)
* `instance_flavor_memory_size` - The memory size of the ECS instance flavor (default: 4)
* `instance_image_id` - The ID of the ECS instance image (default: `""`)
* `instance_image_visibility` - The visibility of the ECS instance image (default: "public")
* `instance_image_os` - The OS of the ECS instance image (default: "Ubuntu")
* `vpc_cidr` - The CIDR block of the VPC (default: "192.168.0.0/16")
* `enterprise_project_id` - The ID of the enterprise project (default: `""`)
* `subnet_cidr` - The CIDR block of the subnet (default: `""`)
* `subnet_gateway_ip` - The gateway IP address of the subnet (default: `""`)
* `internet_gateway_add_route` - Whether to add the default route to the internet gateway (default: true)
* `global_eip_access_site` - The access site used to filter the global EIP pool (default: "cn-north-beijing")
* `global_eip_ip_version` - The IP version of the global EIP (default: "4")
* `internet_bandwidth_charge_mode` - The charge mode of the global internet bandwidth (default: "95peak_guar")
* `internet_bandwidth_size` - The size of the global internet bandwidth in Mbit/s (default: 300)
* `internet_bandwidth_name` - The name of the global internet bandwidth (default: `null`)
* `internet_bandwidth_ingress_size` - The ingress size of the global internet bandwidth in Mbit/s (default: `null`)
* `internet_bandwidth_tags` - The tags of the global internet bandwidth (default: `null`)
* `global_eip_description` - The description of the global EIP (default: `""`)
* `global_eip_tags` - The tags of the global EIP (default: `null`)
* `gc_bandwidth_charge_mode` - The charge mode of the global connection bandwidth (default: "95")
* `gc_bandwidth_size` - The size of the global connection bandwidth in Mbit/s (default: 100)

## Usage

* Copy this example script to your `main.tf`.

* Create a `terraform.tfvars` file and fill in the required variables:

  ```hcl
  vpc_name            = "tf_test_vpc"
  subnet_name         = "tf_test_subnet"
  security_group_name = "tf_test_security_group"

  security_group_rule_configurations = [
    {
      direction        = "ingress"
      ethertype        = "IPv4"
      protocol         = "icmp"
      remote_ip_prefix = "0.0.0.0/0"
    },
    {
      direction        = "ingress"
      ethertype        = "IPv4"
      protocol         = "tcp"
      ports            = "22,3389"
      remote_ip_prefix = "10.1.0.7/32"
    },
    {
      direction        = "egress"
      ethertype        = "IPv4"
      remote_ip_prefix = "0.0.0.0/0"
    },
  ]

  instance_name                   = "tf_test_ecs"
  internet_gateway_name           = "tf_test_igw"
  global_eip_name                 = "tf_test_geip"
  internet_bandwidth_name         = "tf_test_internet_bandwidth"
  gc_bandwidth_name               = "tf_test_gc_bandwidth"
  instance_administrator_password = "YourPassword@123"
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

* Verify network communication:

  1. Use your local PC to remotely log in to the ECS instance.
  2. Run `ping your_ipv4_public_ip_or_domain` on the ECS instance to verify public network access.

* To clean up the resources:

  ```bash
  $ terraform destroy
  ```

## Notes

* Make sure to keep your credentials secure and never commit them to version control
* The ECS instance, VPC internet gateway, and ECS subnet must be in the same region
* The G-EIP association is dependent on the ECS instance and the VPC internet gateway

## Requirements

| Name | Version |
| ---- | ---- |
| terraform | >= 1.3.0 |
| huaweicloud | >= 1.62.0 |
