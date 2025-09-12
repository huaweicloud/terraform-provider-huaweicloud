# Create a VPN gateway

This example provides best practice code for using Terraform to create an gateway in HuaweiCloud VPN service.

## Prerequisites

* A HuaweiCloud account with VPN permissions
* Terraform installed
* HuaweiCloud access key and secret key (AK/SK)

## Variable Introduction

The following variables need to be configured:

### Authentication Variables

* `region_name` - The region where the VPN service is located
* `access_key` - The access key of the IAM user
* `secret_key` - The secret key of the IAM user

### Resource Variables

#### Required Variables

* `vpc_name` - The VPC name.
* `subnet_name` - The subnet name.
* `bandwidth_name` - The bandwidth name.
* `vpn_gateway_name` - The name of the VPN gateway.

#### Optional Variables

* `vpn_gateway_flavor` - The flavor of the VPN gateway (default: "Professional1")
* `vpn_gateway_attachment_type` - The attachment type (default: "vpc")
* `vpc_cidr` - The CIDR block of the VPC (default: "192.168.0.0/16")
* `subnet_cidr` - The CIDR block of the subnet (default: "")
* `subnet_gateway_ip` - The gateway IP of the subnet (default: "")
* `eip_type` - The EIP type (default: "5_bgp")
* `bandwidth_size` - The bandwidth size (default: 8)
* `bandwidth_share_type` - The bandwidth share type (default: "PER")
* `bandwidth_charge_mode` - The bandwidth charge mode" (default: "traffic")
* `vpn_gateway_delete_eip_on_termination` - Whether to delete the EIP when the VPN gateway is deleted (default: false)

## Usage

* Copy this example script to your `main.tf`.

* Create a `terraform.tfvars` file and fill in the required variables:

  ```hcl
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

| Name | Version |
| ---- | ---- |
| terraform | >= 0.12.0 |
| huaweicloud | >= 1.72.0 |
