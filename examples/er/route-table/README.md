# Create a route table with ER instance

This example provides best practice code for using Terraform to create a route table in HuaweiCloud ER service.

## Prerequisites

* A HuaweiCloud account with ER permissions
* Terraform installed
* HuaweiCloud access key and secret key (AK/SK)

## Variable Introduction

The following variables need to be configured:

### Authentication Variables

* `region_name` - The region where the ER service is located
* `access_key` - The access key of the IAM user
* `secret_key` - The secret key of the IAM user

### Resource Variables

#### Required Variables

* `vpc_name` - The VPC name.
* `subnet_name` - The subnet name.
* `er_instance_name` - The ER instance name.
* `er_instance_asn` - The ER instance asn.
* `er_route_table_name` - The ER route table name.

#### Optional Variables

* `vpc_cidr` - The CIDR block of the VPC (default: "192.168.0.0/16")
* `subnet_cidr` - The CIDR block of the subnet (default: "")
* `subnet_gateway_ip` - The gateway IP of the subnet (default: "")

## Usage

* Copy this example script to your `main.tf`.

* Create a `terraform.tfvars` file and fill in the required variables:

  ```hcl
  vpc_name            = "tf_test_er_instance_vpc"
  subnet_name         = "tf_test_er_instance_subnet"
  er_instance_name    = "tf_test_er_instance"
  er_instance_asn     = 64512
  er_route_table_name = "tf_test_er_route_table"
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
* All resources will be created in the specified region

## Requirements

| Name | Version |
| ---- | ---- |
| terraform | >= 0.12.0 |
| huaweicloud | >= 1.69.0 |
