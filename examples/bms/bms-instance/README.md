# Create a BMS instance

This example provides best practice code for using Terraform to create a BMS instance in HuaweiCloud.

## Prerequisites

* A HuaweiCloud account
* Terraform installed
* HuaweiCloud access key and secret key (AK/SK)

## Variables Introduction

The following variables need to be configured:

### Authentication Variables

* `region_name` - The region where the BMS instance is located
* `access_key`  - The access key of the IAM user
* `secret_key`  - The secret key of the IAM user

### Resource Variables

#### Required Variables

* `vpc_name` - The VPC name
* `subnet_name` - The subnet name
* `security_group_name` - The security group name
* `keypair_name` - The KPS keypair name
* `instance_name` - The BMS instance name
* `instance_user_id` - The BMS instance user ID

#### Optional Variables

* `availability_zone` - The availability zones to which the BMS instance belongs (default: "")  
* `vpc_cidr` - The CIDR block of the VPC (default: "192.168.0.0/16")
* `subnet_cidr` - The CIDR block of the subnet (default: "")
* `subnet_gateway_ip` - The gateway IP of the subnet (default: "")
* `instance_flavor_id` - The flavor ID of the BMS instance (default: "")
* `instance_image_id` - The image ID of the BMS instance (default: "")  
* `enterprise_project_id` - The ID of the enterprise project to which the BMS instance belongs (default: null)
* `instance_tags` - The key/value pairs to associate with the BMS instance (default: {})
* `charging_mode` - The charging mode of the BMS instance (default: "prePaid")
* `period_unit` - The period unit of the BMS instance (default: "month")
* `period` - The period of the BMS instance (default: 1)
* `auto_renew` - The auto renew of the BMS instance (default: "false")

## Usage

* Copy this example script to your `main.tf`.

* Create a `terraform.tfvars` file and fill in the required variables:

  ```hcl
  vpc_name            = "your_vpc_name"
  subnet_name         = "your_subnet_name"
  security_group_name = "your_security_group_name"
  keypair_name        = "your_kps_keypair_name"
  instance_name       = "your_bms_instance"
  instance_user_id    = "your_user_id"
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
* The creation of the BMS instance takes about 30 minutes
* All resources will be created in the specified region

## Requirements

| Name | Version |
| ---- | ---- |
| terraform | >= 0.14.0 |
| huaweicloud | >= 1.43.0 |
