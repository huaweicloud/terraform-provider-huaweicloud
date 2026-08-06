# DCS Node Priority Config

This example provides best practice code for using Terraform to configure the slave node
priority weight of a DCS (Distributed Cache Service) Redis HA instance in HuaweiCloud. It
creates a VPC, subnet, Redis HA instance, queries the instance shard information via a data
source, and then sets the slave node priority weight for failover control.

## Prerequisites

* A HuaweiCloud account with DCS permissions
* Terraform installed
* HuaweiCloud access key and secret key (AK/SK)

## Variable Introduction

The following variables need to be configured:

### Authentication Variables

* `region_name` - The region where the DCS instance is located
* `access_key` - The access key of the IAM user
* `secret_key` - The secret key of the IAM user

### Resource Variables

#### Required Variables

* `vpc_name` - The name of the VPC
* `subnet_name` - The name of the subnet
* `instance_name` - The name of the Redis single instance

#### Optional Variables

* `vpc_cidr` - The CIDR block of the VPC (default: "192.168.0.0/16")
* `subnet_cidr` - The CIDR block of the subnet (default: "")
* `subnet_gateway_ip` - The gateway IP address of the subnet (default: "")
* `availability_zone` - The availability zone to which the Redis single instance belongs (default: "")
* `instance_flavor_id` - The flavor ID of the Redis single instance (default: "")
* `instance_capacity` - The capacity of the Redis instance (default: 4)
* `instance_engine_version` - The engine version of the Redis single instance (default: "5.0")
* `instance_password` - The password for the Redis instance (default: null)
* `slave_priority_weight`    - The slave node priority weight, 0-100 (default: 50)

## Usage

* Copy this example script to your `main.tf`.

* Create a `terraform.tfvars` file and fill in the required variables:

  ```hcl
  vpc_name      = "your_vpc_name"
  subnet_name   = "your_subnet_name"
  instance_name = "your_instance_name"
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

## Note

* Make sure to keep your credentials secure and never commit them to version control
* This resource requires a Redis HA instance, as slave nodes only exist in HA architecture
* The `slave_priority_weight` ranges from 0 to 100, where 0 means failover is prohibited
* Deleting this resource only removes it from Terraform state and does not reset the priority configuration on the cloud
  side

## Requirements

| Name        | Version   |
|-------------|-----------|
| terraform   | >= 1.1.0  |
| huaweicloud | >= 1.92.0 |
