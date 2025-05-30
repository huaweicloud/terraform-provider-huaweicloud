# Create a simple ECS instance

This example provides best practice code for using Terraform to create a basic ECS instance in HuaweiCloud.

## Prerequisites

* A HuaweiCloud account
* Terraform installed
* HuaweiCloud access key and secret key (AK/SK)

## Required Variables

The following variables need to be configured:

### Authentication Variables

* `access_key` - HuaweiCloud access key
* `secret_key` - HuaweiCloud secret key
* `region_name` - The region where resources will be created

### Resource Variables

* `vpc_name` - Name of the Virtual Private Cloud (VPC)
* `subnet_name` - Name of the subnet within the VPC
* `security_group_name` - Name of the security group
* `instance_name` - Name of the ECS instance
* `administrator_password` - Password for the administrator account

## Usage

* Copy this example script to your `main.tf`.

* Create a `terraform.tfvars` file and fill in the required variables:

  ```hcl
  access_key           = "your_access_key"
  secret_key           = "your_secret_key"
  region_name          = "your_region"
  vpc_name             = "example-vpc"
  subnet_name          = "example-subnet"
  security_group_name  = "example-sg"
  instance_name        = "example-ecs"
  administrator_password = "your_secure_password"
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
* The administrator password must meet the complexity requirements of HuaweiCloud
* All resources will be created in the specified region
