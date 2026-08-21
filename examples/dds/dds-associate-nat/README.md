# Create a DDS instance associate NAT

This example provides best practice code for using Terraform to create a DDS instance associate NAT in HuaweiCloud.

## Prerequisites

* A HuaweiCloud account
* Terraform installed
* HuaweiCloud access key and secret key (AK/SK)

## Variable Introduction

The following variables need to be configured:

### Authentication Variables

* `region_name` - The region where the DDS instance is located
* `access_key`  - The access key of the IAM user
* `secret_key`  - The secret key of the IAM user

### Resource Variables

#### Required Variables

* `vpc_name` - The VPC name
* `subnet_name` - The subnet name
* `security_group_name` - The security group name
* `gateway_name` - The NAT gateway name
* `eip_bandwidth_name` - The EIP bandwidth name
* `instance_name` - The DDS instance name

#### Optional Variables

* `availability_zone` - The availability zone to which the DDS instance belongs (default: "")
* `vpc_cidr` - The CIDR block of the VPC (default: "192.168.0.0/16")
* `subnet_cidr` - The CIDR block of the subnet (default: "")
* `subnet_gateway_ip` - The gateway IP of the subnet (default: "")
* `gateway_spec` - The specification of the NAT gateway (default : "1")
* `eip_type` - The type of the EIP (default: "5_bgp")
* `eip_bandwidth_size` - The size of the EIP bandwidth in Mbit/s (default: 5)
* `eip_bandwidth_charge_mode` - The charge mode of the EIP bandwidth (default: "traffic")
* `instance_mode` - The type of the DDS instance (default: "ReplicaSet")
* `database_type` - The database version type of the DDS instance (default: "DDS-Community")
* `database_version` - The database version of the DDS instance (default: "4.0")
* `storage_engine` - The storage engine of the DDS instance (default: "wiredTiger")
* `node_type` - The type of the DDS instance node (default: "replica")
* `node_number` - The number of nodes of the DDS instance, the value can be `3`, `5` or `7` (default: 3)
* `node_spec_code` - The spec code of the DDS instance node (default: "dds.mongodb.s6.large.2.repset")
* `node_storage_type` - The storage type of the DDS instance node (default: "ULTRAHIGH")
* `node_size` - The disk size of the node of the DDS instance, The value must be a multiple of `10` (default: 10)
* `external_service_port` - The port of the EIP for providing services for external systems (default: 8080)
  The valid value ranges from `1` to `65,535`

## Usage

* Copy this example script to your `main.tf`.

* Create a `terraform.tfvars` file and fill in the required variables:

  ```hcl
  vpc_name            = "your_vpc_name"
  subnet_name         = "your_subnet_name"
  security_group_name = "your_security_group_name"
  gateway_name        = "your_nat_gateway_name"
  eip_bandwidth_name  = "your_eip_bandwidth_name"
  instance_name       = "your_instance_name"
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
* The creation of the DDS instance takes about 20 to 60 minutes depending on the flavor and node number
* All resources will be created in the specified region

## Requirements

| Name | Version |
| ---- | ---- |
| terraform | >= 1.1.0 |
| huaweicloud | >= 1.93.0 |
