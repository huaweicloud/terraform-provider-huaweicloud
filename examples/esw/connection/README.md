# Using NAT Gateway and VPC Peering for Cross-VPC Internet Access

This example provides best practice code for using Terraform to create a connection in HuaweiCloud ESW service.

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

* `vpc_name` - The VPC name.
* `subnet_name` - The subnet name.
* `esw_instance_name` - The instance name.
* `esw_instance_ha_mode` - The instance high availability mode.
* `esw_connection_name` - The instance high availability mode.

#### Optional Variables

* `vpc_cidr` - The CIDR block of the VPC (default: "192.168.0.0/16")
* `subnet_cidr` - The CIDR block of the subnet (default: "")
* `subnet_gateway_ip` - The gateway IP of the subnet (default: "")
* `esw_instance_description` - The instance description (default: "").
* `esw_instance_tunnel_ip` - The tunnel IP (default: "").
* `esw_connection_segmentation_id` - The tunnel number for the connection corresponds to the VXLAN network
  identifier (VNI).
* `esw_connection_fixed_ips` - The downlink network port primary and standby IPs.

## Usage

* Copy this example script to your `main.tf`.

* Create a `terraform.tfvars` file and fill in the required variables:

  ```hcl
  vpc_name             = "tf_test_esw_instance"
  subnet_name          = "tf_test_esw_instance"
  esw_instance_name    = "tf_test_esw_instance"
  esw_instance_ha_mode = "ha"
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

## Requirements

| Name | Version   |
| ---- |-----------|
| terraform | >= 0.14.0 |
| huaweicloud | >= 1.82.0 |
