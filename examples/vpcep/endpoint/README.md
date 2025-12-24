# Create a VPC endpoint

This example provides best practice code for using Terraform to create a VPC endpoint in HuaweiCloud.

## Prerequisites

* A HuaweiCloud account with VPC endpoint permissions
* Terraform installed
* HuaweiCloud access key and secret key (AK/SK)

## Variable Introduction

The following variables need to be configured:

### Authentication Variables

* `region_name` - The region where the SWR service is located
* `access_key` - The access key of the IAM user
* `secret_key` - The secret key of the IAM user

### Resource Variables

#### Required Variables

* `vpc_name` - The name of the VPC
* `vpc_cidr` - The CIDR block of the VPC
* `subnet_name` - The name of the VPC subnet
* `security_group_name` - The name of the security group
* `instance_name` - The name of the ECS instance
* `endpoint_service_name` - The name of the endpoint service
* `endpoint_service_port_mapping` - The port mapping of the endpoint service

#### Optional Variables

* `instance_flavor_performance_type` - The performance type of the ECS instance flavor (default: "normal")
* `instance_flavor_cpu_core_count` - The CPU core count of the ECS instance flavor (default: 2)
* `instance_flavor_memory_size` - The memory size of the ECS instance flavor (default: 4)
* `instance_image_name` - The name of the ECS instance image (default: "Ubuntu 20.04 server 64bit")
* `instance_image_most_recent` - Whether the instance image is most recent or not (default: true)
* `subnet_cidr` - The CIDR block of the VPC subnet (default: "")
* `subnet_gateway_ip` - The gateway IP of the VPC subnet (default: "")
* `instance_flavor_id` - The ID of the ECS instance flavor (default: "")
* `endpoint_service_type` - The type of the endpoint service (default: "VM")

## Usage

* Copy this example script to your `main.tf`.

* Create a `terraform.tfvars` file and fill in the required variables:

  ```hcl
  vpc_name                      = "tf_test_vpc"
  vpc_cidr                      = "192.168.0.0/16"
  subnet_name                   = "tf_test_subnet"
  security_group_name           = "tf_test_security_group"
  instance_name                 = "tf_test_instance"
  endpoint_service_name         = "tf_test_endpoint_service_name"
  endpoint_service_port_mapping = [
    {
      service_port  = 8080
      terminal_port = 80
    }
  ]
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
| terraform | >= 0.14.0 |
| huaweicloud | >= 1.25.0 |
