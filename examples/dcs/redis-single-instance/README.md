# Create a DCS Redis single instance

This example provides best practice code for using Terraform to create a Redis single instance in HuaweiCloud DCS service.

## Prerequisites

* A HuaweiCloud account
* Terraform installed
* HuaweiCloud access key and secret key (AK/SK)

## Required Variables

The following variables need to be configured:

### Authentication Variables

* `region_name` - The region where the DCS Redis single instance is located
* `access_key`  - The access key of the IAM user
* `secret_key`  - The secret key of the IAM user

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
* `instance_capacity` - The capacity of the Redis instance (default: 1)
* `instance_engine_version` - The engine version of the Redis single instance (default: "7.0")
* `enterprise_project_id` - The ID of the enterprise project to which the Redis single instance belongs (default: null)
* `instance_password` - The password for the Redis instance (default: null)
* `charging_mode` - The charging mode of the Redis instance (default: "postPaid")
* `period_unit` - The unit of the period, only used when `charging_mode` is `prePaid` (default: null)
* `period` - The period of the Redis instance, only used when `charging_mode` is `prePaid` (default: null)
* `auto_renew` - Whether auto renew is enabled, only available when `charging_mode` is `prePaid` (default: "false")

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

* To clean up the resources:

  ```bash
  $ terraform destroy
  ```

## Note

* Make sure to keep your credentials secure and never commit them to version control
* The creation of the DCS Redis instance may take several minutes
* All resources will be created in the specified region

## Requirements

| Name | Version |
| ---- | ------- |
| terraform | >= 1.1.0 |
| huaweicloud | >= 1.29.0 |
