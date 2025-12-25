# Create a DNAT rule

This example provides best practice code for using Terraform to create a DNAT (Destination Network Address Translation)
rule in HuaweiCloud.

## Prerequisites

- A HuaweiCloud account
- Terraform installed
- HuaweiCloud access key and secret key (AK/SK)

## Variables introduction

The following variables need to be configured:

## Required Variables

### Authentication Variables

- `region_name` - The region where the NAT gateway and DNAT rule are located
- `access_key` - The access key of the IAM user
- `secret_key` - The secret key of the IAM user

### Resource Variables

#### Optional Variables

- `vpc_name` - The VPC name (default: "vpc-dnat-basic")
- `vpc_cidr` - The CIDR block of the VPC (default: "172.16.0.0/16")
- `subnet_name` - The subnet name (default: "subnet-dnat-basic")
- `subnet_cidr` - The CIDR block of the subnet (default: "172.16.10.0/24")
- `subnet_gateway_ip` - The gateway IP address of the subnet (default: ""). If not specified,
  it will be automatically calculated from the VPC CIDR
- `ecs_flavor_id` - The flavor ID of the backend ECS instance (default: ""). If not specified,
  a flavor will be automatically queried using data source based on the following parameters:
  + `ecs_flavor_performance_type` - The performance type of the ECS instance flavor (default: "normal")
  + `ecs_flavor_cpu_core_count` - The CPU core count of the ECS instance flavor (default: 2)
  + `ecs_flavor_memory_size` - The memory size of the ECS instance flavor in GB (default: 4)
- `ecs_image_id` - The image ID of the backend ECS instance (default: ""). If not specified,
  an image will be automatically queried using data source based on the following parameters:
  + `ecs_image_visibility` - The visibility of the ECS instance image (default: "public")
  + `ecs_image_os` - The OS of the ECS instance image (default: "Ubuntu")
- `ecs_system_disk_type` - The system disk type of the ECS instance (default: "SSD")
- `ecs_system_disk_size` - The system disk size of the ECS instance in GB (default: 40)
- `ecs_admin_password` - The administrator password of the ECS instance (default: "", sensitive)
- `ecs_instance_tags` - The tags of the ECS instance (default: {})
- `security_group_name` - The security group name of the backend instance (default: "sg-dnat-backend")
- `ingress_cidr` - The CIDR block that is allowed to access the DNAT service from the Internet (default: "0.0.0.0/0")
- `instance_name` - The name of the backend ECS instance (default: "ecs-dnat-backend")
- `backend_protocol` - The protocol used between the NAT gateway and backend ECS instance (default: "tcp")
  Valid values: **tcp**, **udp**, **any**
- `backend_port` - The port on the backend ECS instance that receives DNAT traffic (default: `22`)
- `frontend_protocol` - The protocol used between the client and NAT gateway (default: "tcp")
  Valid values: **tcp**, **udp**, **any**
- `frontend_port` - The port on the public EIP that clients use to access the DNAT service (default: `22`)
- `eip_bandwidth_name` - The name of the EIP bandwidth (default: "eip-dnat-basic")
- `eip_bandwidth_size` - The size of the EIP bandwidth in Mbit/s (default: `5`)
- `eip_bandwidth_share_type` - The share type of the EIP bandwidth (default: "PER")
- `eip_bandwidth_charge_mode` - The charge mode of the EIP bandwidth (default: "traffic")
- `nat_gateway_name` - The NAT gateway name (default: "nat-gateway-dnat-basic")
- `nat_gateway_description` - The description of the NAT gateway (default: "")
- `nat_gateway_spec` - The specification of the NAT gateway (default: "1")
  Valid values: **1** (Small), **2** (Medium), **3** (Large), **4** (Extra-large)

## Usage

1. Copy the files in this directory to your Terraform working directory.

2. Update `terraform.tfvars` with your own values. You can either specify the IDs directly or
   let Terraform automatically query them using data sources:

   ```hcl
   ecs_flavor_id = "s6.small.1"
   ecs_image_id  = "YOUR_IMAGE_ID"
   ecs_admin_password = "YourPassword123!"
   ```

3. Initialize Terraform:

   ```bash
   terraform init
   ```

4. Review the Terraform plan:

   ```bash
   terraform plan
   ```

5. Apply the configuration:

   ```bash
   terraform apply
   ```

6. After the resources are created successfully, you can connect to the backend ECS through the DNAT rule by
   using the EIP and configured frontend port.

7. To clean up all resources:

   ```bash
   terraform destroy
   ```

### Notes

- Keep your credentials secure and never commit them to version control
- DNAT rules are created on the specified NAT gateway and are bound to the created EIP
- Make sure the security group of the backend ECS allows traffic from the NAT gateway on the configured backend port
- All resources are created in the specified region

## Requirements

| Name | Version |
| ---- | ---- |
| terraform | >= 1.1.0 |
| huaweicloud | >= 1.57.0 |
