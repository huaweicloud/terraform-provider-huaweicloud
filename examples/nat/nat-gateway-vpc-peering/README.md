# Using NAT Gateway and VPC Peering for Cross-VPC Internet Access

This example provides best practice code for using Terraform to configure a NAT gateway with VPC peering to enable
cross-VPC internet access in HuaweiCloud.

-> The NAT gateway is created in one VPC and shared with other VPCs through VPC peering connections, allowing
   multiple VPCs to access the internet while reducing costs by sharing EIPs and NAT gateway resources.

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

* `ecs_admin_password` - The administrator password of ECS instances

#### Optional Variables

* `image_name` - The name of the image to use for ECS instances (default: "Ubuntu 22.04 server 64bit")
* `image_visibility` - The visibility of the image (default: "public")
* `image_most_recent` - If more than one result is returned, use the most recent image (default: true)
* `flavor_performance_type` - The performance type of the flavor (default: "normal")
* `flavor_cpu_core_count` - The number of CPU cores for the flavor (default: `2`)
* `flavor_memory_size` - The memory size in GB for the flavor (default: `4`)
* `secgroup_name` - The name of the security group (default: "secgroup-example")
* `security_group_rules` - List of security group rules (default: allows SSH, HTTP, HTTPS, and ICMP)
  Each rule object supports the following attributes:
  + `direction` - The direction of the rule (required)
  + `description` - The description of the rule (required)
  + `ethertype` - The ethertype of the rule (required)
  + `protocol` - The protocol of the rule (required)
  + `ports` - The port range of the rule (optional)
    Format: single port (e.g., "22") or port range (e.g., "80,443")
  + `remote_ip_prefix` - The remote IP prefix (required)
    Format: CIDR notation (e.g., "0.0.0.0/0" for all IPs)
* `vpcs` - List of VPC configurations (default: two VPCs with default names and CIDRs)
  Each VPC object supports the following attributes:
  + `vpc_name` - The name of the VPC (required)
  + `vpc_cidr` - The CIDR block of the VPC (required)
  + `subnet_name` - The name of the subnet (required)
  + `subnet_cidr` - The CIDR block of the subnet (required)
  + `subnet_gateway_ip` - The gateway IP of the subnet (optional, auto-calculated if not provided)
  + `instance_name` - The name of the ECS instance in this VPC (required)
  + `ecs_flavor_id` - The ID of the ECS flavor for this instance (optional, uses the first available flavor from
  the data source if not specified)
  + `ecs_image_id` - The ID of the ECS image for this instance (optional, uses the image from the data source if not specified)
  + `ecs_system_disk_type` - The type of the system disk for this instance (optional, defaults to "SSD" if not specified)
  + `ecs_system_disk_size` - The size of the system disk in GB for this instance (optional, defaults to `40` if not specified)
  + `ecs_admin_password` - The administrator password for this instance (optional,
  uses global `ecs_admin_password` if not specified)
  + `ecs_instance_tags` - The tags for this instance (optional, defaults to empty map `{}` if not specified)
* `peering_connection_name` - The name of the VPC peering connection (default: "peering-connection-example")
* `route_type` - The type of the route (default: "peering")
* `eips` - List of EIP configurations (default: two EIPs for SNAT and DNAT)
  Each EIP object supports the following attributes:
  + `name_suffix` - The suffix for the bandwidth name (required, e.g., "snat", "dnat")
* `nat_gateway_name` - The name of the NAT gateway (default: "nat-gateway-example")
* `bandwidth_size` - The bandwidth size in Mbps (default: `5`)
* `bandwidth_charge_mode` - The charge mode of the bandwidth (default: "traffic")
* `nat_gateway_spec` - The specification of the NAT gateway (default: "1")
* `snat_rules` - List of SNAT rule configurations (default: two SNAT rules)
  Each SNAT rule object supports the following attributes:
  + `eip_index` - The index of the EIP to use (required)
  + `subnet_index` - The index of the subnet (optional, for subnet-based SNAT)
  + `source_type` - The source type (optional, for CIDR-based SNAT, use `1`)
  + `vpc_index` - The index of the VPC to get subnet CIDR (optional, for CIDR-based SNAT)
* `dnat_rules` - List of DNAT rule configurations (default: two DNAT rules)
  Each DNAT rule object supports the following attributes:
  + `eip_index` - The index of the EIP to use (required)
  + `instance_index` - The index of the instance (required, corresponds to vpcs index)
  + `external_service_port` - The external service port (required)
* `dnat_protocol` - The protocol for DNAT rules (default: "tcp")
* `dnat_internal_port` - The internal service port for DNAT rules (default: "22")

## Usage

* Copy this example script to your `main.tf`.

* Create a `terraform.tfvars` file and fill in the required variables:

  ```hcl
  ecs_admin_password = "YourSecurePassword123!"

  security_group_rules = [
    {
      direction        = "ingress"
      description      = "Allow SSH access"
      ethertype        = "IPv4"
      protocol         = "tcp"
      ports            = "22"
      remote_ip_prefix = "1.2.3.4/32"  # Restrict to your IP
    },
    {
      direction        = "ingress"
      description      = "Allow HTTP and HTTPS access"
      ethertype        = "IPv4"
      protocol         = "tcp"
      ports            = "80,443"
      remote_ip_prefix = "0.0.0.0/0"
    },
    {
      direction        = "ingress"
      description      = "Allow ICMP access"
      ethertype        = "IPv4"
      protocol         = "icmp"
      ports            = null
      remote_ip_prefix = "0.0.0.0/0"
    }
  ]

  # VPC configurations
  vpcs = [
    {
      vpc_name          = "my-vpc-one"
      vpc_cidr          = "192.168.0.0/16"
      subnet_name       = "my-subnet-one"
      subnet_cidr       = "192.168.1.0/24"
      subnet_gateway_ip = "192.168.1.1"
      instance_name     = "my-instance-one"
    },
    {
      vpc_name          = "my-vpc-other"
      vpc_cidr          = "10.0.0.0/16"
      subnet_name       = "my-subnet-other"
      subnet_cidr       = "10.0.2.0/24"
      subnet_gateway_ip = "10.0.2.1"
      instance_name     = "my-instance-other"
    }
  ]

  nat_gateway_name = "my-nat-gateway"
  bandwidth_size   = 10
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

* After deployment, you can access the instances via DNAT:

  ```bash
  # Access Instance One via DNAT (using external port from dnat_rules[0])
  ssh -p 8022 root@<eip_dnat_address>

  # Access Instance Other via DNAT (using external port from dnat_rules[1])
  ssh -p 8023 root@<eip_dnat_address>
  ```

* To clean up the resources:

  ```bash
  $ terraform destroy
  ```

## Note

* Make sure to keep your credentials secure and never commit them to version control
* The NAT gateway is created in VPC One and shared with VPC Other via VPC peering
* VPC Other's default route (0.0.0.0/0) is configured to use the peering connection, enabling internet access
  through VPC One's NAT gateway
* Both VPCs share the same EIPs for SNAT and DNAT rules, reducing EIP costs
* For DNAT rules, ensure different external ports are used for each instance (configured via
  `dnat_rules` variable, where each rule specifies `external_service_port`)
* VPC One uses CIDR block 192.168.0.0/16 and VPC Other uses 10.0.0.0/16, which is valid for VPC peering
* Security group rules should be adjusted based on your security requirements
* EIP bandwidth and charge mode can be customized via variables
* NAT Gateway: Charged based on specification and usage
* EIP: Only 2 EIPs are created (one for SNAT, one for DNAT), shared by both VPCs, reducing costs compared
  to using separate EIPs for each VPC
* ECS Instances: Charged based on instance specifications
* VPC Peering: Free for same-region peering
* All resources will be created in the specified region

## Requirements

| Name | Version |
| ---- | ---- |
| terraform | >= 1.3.0 |
| huaweicloud | >= 1.34.0 |
