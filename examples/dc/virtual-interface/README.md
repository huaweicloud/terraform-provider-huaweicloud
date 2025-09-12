# Create a DC virtual interface example

This example provides best practice code for using Terraform to create a DC (Direct Connect) virtual interface in HuaweiCloud.

## Prerequisites

* A HuaweiCloud account
* Terraform installed
* HuaweiCloud access key and secret key (AK/SK)
* A direct connect connection

## Variables Introduction

The following variables need to be configured:

### Authentication Variables

* `region_name` - The region where the resources are located
* `access_key` - The access key of the IAM user
* `secret_key` - The secret key of the IAM user

### Resource Variables

#### Required Variables

* `vpc_name` - The name of the VPC
* `virtual_gateway_name` - The name of the virtual gateway
* `virtual_interface_name` - The name of the virtual interface
* `direct_connect_id` - The ID of the direct connection associated with the virtual interface
* `vlan` - The VLAN for constom side
* `bandwidth` - The ingress bandwidth size of the virtual interface
* `remote_ep_group` - The CIDR list of remote subnets
* `local_gateway_v4_ip` - The IPv4 address of the virtual interface in cloud side
* `remote_gateway_v4_ip` - The IPv4 address of the virtual interface in client side

#### Optional Variables

* `vpc_cidr` - The CIDR block of the VPC (default: "192.168.0.0/16")
* `virtual_interface_description` - The description of the virtual interface (default: "Created by Terraform")
* `virtual_interface_type` - The type of the virtual interface (default: "private")
* `route_mode` - The route mode of the virtual interface (default: "static")
* `address_family` - The address family type of the virtual interface (default: "ipv4")
* `enable_bfd` - Whether to enable the Bidirectional Forwarding Detection (BFD) function (default: false)
* `enable_nqa` - Whether to enable the Network Quality Analysis (NQA) function (default: false)
* `virtual_interface_tags` - The tags of the virtual interface (default: { "Owner" = "terraform", "Env" = "test" })

## Usage

* Copy this example script to your `main.tf`.

* Create a `terraform.tfvars` file and fill in the required variables:

  ```hcl
  vpc_name               = "your_vpc_name"
  virtual_gateway_name   = "your_virtual_gateway_name"
  virtual_interface_name = "your_virtual_interface_name"
  direct_connect_id      = "your_direct_connect_id"
  vlan                   = your_vlan
  bandwidth              = your_bandwidth
  remote_ep_group        = ["your_remote_ep_group"]
  local_gateway_v4_ip    = "your_local_gateway_v4_ip"
  remote_gateway_v4_ip   = "your_remote_gateway_v4_ip"
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
* The virtual interface will be associated with a virtual gateway, which is connected to a VPC
* All resources will be created in the specified region

## Requirements

| Name | Version |
| ---- | ---- |
| terraform | >= 1.9.0 |
| huaweicloud | >= 1.77.6 |
