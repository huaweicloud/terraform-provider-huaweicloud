# Create a VPC Internet Gateway

This example provides best practice code for using Terraform to create a VPC internet gateway in HuaweiCloud EIP
service. The example first creates a VPC and a subnet, and then creates an internet gateway associated with them.
The internet gateway is the prerequisite for global EIP (GEIP) cross-region networking.

## Prerequisites

* A HuaweiCloud account
* Terraform installed
* HuaweiCloud access key and secret key (AK/SK)

## Variable Introduction

The following variables need to be configured:

### Authentication Variables

* `region_name` - The region where the EIP resources are located
* `access_key` - The access key of the IAM user
* `secret_key` - The secret key of the IAM user

### Resource Variables

#### Required Variables

* `vpc_name` - The name of the VPC
* `subnet_name` - The name of the subnet
* `internet_gateway_name` - The name of the VPC internet gateway

#### Optional Variables

* `vpc_cidr` - The CIDR block of the VPC (default: "192.168.0.0/16")
* `subnet_cidr` - The CIDR block of the subnet (default: `""`, derived from the VPC CIDR block)
* `subnet_gateway_ip` - The gateway IP address of the subnet (default: `""`, derived from the subnet CIDR block)
* `internet_gateway_add_route` - Whether to add a default route pointing to the internet gateway (default: true)
* `internet_gateway_ipv6_enabled` - Whether to enable IPv6 for the internet gateway (default: false)

## Usage

* Copy this example script to your `main.tf`.

* Create a `terraform.tfvars` file and fill in the required variables:

  ```hcl
  vpc_name              = "tf_test_vpc"
  subnet_name           = "tf_test_subnet"
  internet_gateway_name = "tf_test_internet_gateway"
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

## Notes

* Make sure to keep your credentials secure and never commit them to version control
* A VPC can only be associated with one internet gateway
* Changing `internet_gateway_add_route` recreates the internet gateway
* The `internet_gateway_ipv6_enabled` cannot be changed from true to false, and the subnet must enable IPv6 before
  setting it to true
* When `internet_gateway_add_route` is true, a default route with the destination 0.0.0.0/0 pointing to the internet
  gateway will be added to the default route table of the VPC

## Requirements

| Name | Version |
| ---- | ---- |
| terraform | >= 1.1.0 |
| huaweicloud | >= 1.62.0 |
