# Create a WAF dedicated instance

This example provides best practice code for using Terraform to create a WAF dedicated instance in HuaweiCloud WAF service.

## Prerequisites

* A HuaweiCloud account
* Terraform installed
* HuaweiCloud access key and secret key (AK/SK)

## Required Variables

The following variables need to be configured:

### Authentication Variables

* `region_name` - The region where the WAF dedicated instance is located
* `access_key` - The access key of the IAM user
* `secret_key` - The secret key of the IAM user

### Resource Variables

#### Required Variables

* `vpc_name` - The VPC name
* `subnet_name` - The subnet name
* `security_group_name` - The security group name
* `dedicated_instance_name` - The WAF dedicated instance name

#### Optional Variables

* `vpc_cidr` - The CIDR block of the VPC (default: "192.168.0.0/16")
* `availability_zone` - The availability zone to which the dedicated instance belongs (default: "")
* `subnet_cidr` - The CIDR block of the subnet (default: "")
* `subnet_gateway_ip` - The gateway IP address of the subnet (default: "")
* `dedicated_instance_flavor_id` - The flavor ID of the dedicated instance (default: "")
* `dedicated_instance_performance_type` - The performance type of the dedicated instance (default: "normal")
* `dedicated_instance_cpu_core_count` - The number of the dedicated instance CPU cores (default: 4)
* `dedicated_instance_memory_size` - The memory size of the dedicated instance (default: 8)
* `dedicated_instance_specification_code` - The specification code of the dedicated instance (default: "waf.instance.professional")

## Usage

* Copy this example script to your `main.tf`.

* Create a `terraform.tfvars` file and fill in the required variables:

  ```hcl
  vpc_name                = "your_vpc_name"
  subnet_name             = "your_subnet_name"
  security_group_name     = "your_security_group_name"
  dedicated_instance_name = "your_waf_instance_name"
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
* The creation of the WAF dedicated instance takes about 5 minutes
* This example only creates the WAF dedicated instance, VPC, subnet and security group
* All resources will be created in the specified region

## Requirements

| Name | Version |
| ---- | ---- |
| terraform | >= 0.12.0 |
| huaweicloud | >= 1.28.0 |
