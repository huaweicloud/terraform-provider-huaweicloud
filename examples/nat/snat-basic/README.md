# Create a SNAT rule

This example provides best practice code for using Terraform to create a SNAT (Source Network Address Translation) rule
in HuaweiCloud.

## Prerequisites

* A HuaweiCloud account
* Terraform installed
* HuaweiCloud access key and secret key (AK/SK)

## Variables Introduction

The following variables need to be configured:

### Authentication Variables

* `region_name` - The region where the SNAT rule is located
* `access_key` - The access key of the IAM user
* `secret_key` - The secret key of the IAM user

### Resource Variables

#### Required Variables

* `nat_gateway_name` - The name of the NAT gateway
* `gateway_spec` - The specification of the NAT gateway
  The valid values are allowed:
  + **1**: Small, `10,000` concurrent connections
  + **2**: Medium, `50,000` concurrent connections
  + **3**: Large, `200,000` concurrent connections
  + **4**: Extra-large, `1,000,000` concurrent connections
* `eip_bandwidth_name` - The name of the EIP bandwidth
* `ecs_instance_name` - The name of the ECS instance
* `ecs_admin_password` - The password of the ECS instance administrator

#### Optional Variables

* `gateway_description` - The description of the NAT gateway (default: "")
* `eip_bandwidth_size` - The size of the EIP bandwidth in Mbps (default: `5`)
* `eip_bandwidth_share_type` - The share type of the EIP bandwidth (default: "PER")
  Valid values: **PER** (dedicated), **WHOLE** (shared)
* `eip_bandwidth_charge_mode` - The charge mode of the EIP bandwidth (default: "traffic")
  Valid values: **traffic**, **bandwidth**
* `snat_source_type` - The resource type of the SNAT rule (default: `0`)
  Valid values: **0** (VPC scenario), **1** (Direct Connect scenario)
* `snat_cidr` - The CIDR block connected by SNAT rule (DC side, required when snat_source_type is `1`)
  This parameter and `snat_subnet_id` are alternative
* `snat_description` - The description of the SNAT rule (default: "")
* `ecs_flavor_id` - The ID of the ECS instance flavor (optional, will be queried if not provided)
* `ecs_flavor_performance_type` - The performance type of the ECS instance flavor (default: "normal")
* `ecs_flavor_cpu_core_count` - The CPU core count of the ECS instance flavor (default: `2`)
* `ecs_flavor_memory_size` - The memory size of the ECS instance flavor in GB (default: `4`)
* `ecs_image_id` - The ID of the ECS instance image (optional, will be queried if not provided)
* `ecs_image_visibility` - The visibility of the ECS instance image (default: "public")
* `ecs_image_os` - The OS of the ECS instance image (default: "Ubuntu")
* `ecs_security_group_name` - The name of the security group for ECS instance (default: "terraform-test-ecs-sg")
* `ecs_admin_password` - The password of the ECS instance administrator
* `ecs_system_disk_type` - The type of the ECS instance system disk (default: "SSD")
* `ecs_system_disk_size` - The size of the ECS instance system disk in GB (default: `40`)
* `ecs_instance_tags` - The tags of the ECS instance (default: {})

## Usage

* Copy this example script to your `main.tf`.

* Create a `terraform.tfvars` file and fill in the required variables:

  ```hcl
  # NAT Gateway Configuration
  nat_gateway_name   = "terraform-test-nat-gateway"
  gateway_spec       = "1"

  # EIP Configuration
  eip_bandwidth_name = "terraform-test-eip-bandwidth"
  eip_bandwidth_size = 5

  # SNAT Rule Configuration
  snat_source_type   = 0
  snat_description   = "SNAT rule for terraform test"

  # ECS Instance Configuration
  ecs_instance_name  = "terraform-test-ecs"
  ecs_admin_password = "YourPassword123!"
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
* The SNAT rule creation may take a few minutes to complete
* The NAT gateway and SNAT rule must be in the same region
* When using `snat_source_type = 1` (Direct Connect scenario), only `snat_cidr` can be specified, and `snat_subnet_id`
  must be empty
* The `snat_subnet_id` and `snat_cidr` parameters are alternative - you must specify one of them
* All resources will be created in the specified region
* The NAT gateway must be in the ACTIVE status before creating SNAT rules
* The ECS instance is created in the subnet specified by `snat_subnet_id`
* Make sure the subnet has available IP addresses for the ECS instance
* The ECS instance will use the SNAT rule to access the internet, so it doesn't need a public IP address

## Requirements

| Name | Version |
| ---- | ---- |
| terraform | >= 1.1.0 |
| huaweicloud | >= 1.57.0 |
