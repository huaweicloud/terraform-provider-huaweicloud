# Create a CBH basic instance

This example provides best practice code for using Terraform to create a CBH basic instance in HuaweiCloud.

## Prerequisites

* A Huawei Cloud account
* Terraform installed
* Huawei Cloud access key and secret key (AK/SK)

## Required Variables

### Authentication Variables

* `region_name` - The region where the CBH instance is located
* `access_key`  - The access key of the IAM user
* `secret_key`  - The secret key of the IAM user

### Resource Variables

#### Required Variables

* `vpc_name` - The name of the VPC
* `subnet_name` - The name of the subnet
* `security_group_name` - The name of the security group
* `instance_name` - The name of the CBH instance
* `instance_flavor_id` - The flavor ID of the CBH instance
* `instance_flavor_type` - The flavor type of the CBH instance
* `instance_password` - The login password of the CBH instance

#### Optional Variables

* `availability_zone` - The availability zone of the CBH instance (default: "")
* `vpc_cidr` - The CIDR block of the VPC (default: "192.168.0.0/16")
* `subnet_cidr` - The CIDR block of the subnet (default: "")
* `subnet_gateway_ip` - The gateway IP of the subnet (default: "")
* `charging_mode` - The charging mode of the CBH instance (default: "prePaid")
* `period_unit` - The charging period unit of the CBH instance (default: "month")
* `period` - The charging period of the CBH instance (default: 1)
* `auto_renew` - Whether to enable auto-renew for the CBH instance (default: "false")

## Usage

* Copy this example script to your `main.tf`.
* Create a `terraform.tfvars` file and fill in the required variables:

  ```hcl
  vpc_name            = "your_vpc_name"
  subnet_name         = "your_subnet_name"
  security_group_name = "your_security_group_name"
  instance_name       = "your_cbh_instance_name"
  instance_flavor_id  = "your_cbh_instance_flavor_id"
  instance_password   = "your_cbh_instance_password"
  ```

* Initialize Terraform:

  ```bash
  terraform init
  ```

* Review the Terraform plan:

  ```bash
  terraform plan
  ```

* Apply the configuration:

  ```bash
  terraform apply
  ```

* To clean up the resources:

  ```bash
  terraform destroy
  ```

## Notes

* Make sure to keep your credentials secure and never commit them to version control
* It takes about 15 minutes to create a CBH basic instance
* All resources will be created in the specified region

## Requirements

| Name        | Version   |
|-------------|-----------|
| terraform   | >= 1.1.0 |
| huaweicloud | >= 1.64.3 |
