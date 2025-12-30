# Deploy ECS instance on dedicated host

This example provides best practice code for using Terraform to deploy an ECS instance on a dedicated host in
HuaweiCloud.

## Prerequisites

* A HuaweiCloud account with DEH and ECS permissions
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

* `deh_instance_name` - The name of the dedicated host instance
* `vpc_name` - The name of the VPC
* `subnet_name` - The name of the VPC subnet
* `security_group_name` - The name of the security group
* `ecs_instance_name` - The name of the ECS instance
* `ecs_instance_admin_pass` - The password of the administrator

#### Optional Variables

* `availability_zone` - The availability zone where the resources will be created (default: "")
* `deh_instance_host_type` - The host type of the dedicated host (default: "")
* `deh_instance_auto_placement` - Whether to enable auto placement for the dedicated host (default: "on")
* `enterprise_project_id` - The enterprise project ID of the dedicated host (default: null)
* `deh_instance_charging_mode` - The charging mode of the dedicated host (default: "prePaid")
* `deh_instance_period_unit` - The unit of the billing period of the dedicated host (default: "month")
* `deh_instance_period` - The billing period of the dedicated host (default: "1")
* `deh_instance_auto_renew` - Whether to enable auto renew for the dedicated host (default: "false")
* `vpc_cidr` - The CIDR block of the VPC (default: "192.168.0.0/16")
* `subnet_cidr` - The CIDR block of the VPC subnet (default: "")
* `subnet_gateway_ip` - The gateway IP of the VPC subnet (default: "")
* `ecs_instance_image_id` - The ID of the ECS instance image (default: "")
* `ecs_instance_flavor_id` - The ID of the ECS instance flavor (default: "")
* `ecs_instance_image_visibility` - The visibility of the ECS instance image (default: "public")
* `ecs_instance_image_os` - The OS of the ECS instance image (default: "Ubuntu")

## Usage

* Copy this example script to your `main.tf`.

* Create a `terraform.tfvars` file and fill in the required variables:

  ```hcl
  deh_instance_name       = "tf_test_deh_instance"
  vpc_name                = "tf_test_vpc"
  subnet_name             = "tf_test_subnet"
  security_group_name     = "tf_test_security_group"
  ecs_instance_name       = "tf_test_ecs_instance"
  ecs_instance_admin_pass = "YourPassword"
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
* Make sure to have sufficient quota for the dedicated host and ECS instance resources you plan to create
* The ECS instance flavor must match the dedicated host type.
* Currently, only on-demand (postPaid) ECS instances are supported on dedicated hosts

## Requirements

| Name | Version |
| ---- | ---- |
| terraform | >= 1.1.0 |
| huaweicloud | >= 1.74.0 |
