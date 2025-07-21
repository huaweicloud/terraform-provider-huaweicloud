# Create a basic VPC and subnet

This example provides best practice code for using Terraform to create a basic VPC and subnet in HuaweiCloud VPC service.

## Prerequisites

* A HuaweiCloud account
* Terraform installed
* HuaweiCloud access key and secret key (AK/SK)

## Required Variables

The following variables need to be configured:

### Authentication Variables

* `region_name` - The region where the VPC is located
* `access_key` - The access key of the IAM user
* `secret_key` - The secret key of the IAM user

### Resource Variables

#### Required Variables

* `vpc_name` - The VPC name
* `subnet_name` - The subnet name

#### Optional Variables

* `vpc_cidr` - The CIDR block of the VPC (default: "172.16.0.0/16")
* `enterprise_project_id` - The ID of the enterprise project to which the VPC belongs
* `subnet_cidr` - The CIDR block of the subnet (default: "172.16.10.0/24")
* `subnet_gateway` - The gateway IP address of the subnet (default: "172.16.10.1")
* `primary_dns` - The primary DNS server IP address (default: "100.125.1.250")

## Usage

* Copy this example script to your `main.tf`.

* Create a `terraform.tfvars` file and fill in the required variables:

  ```hcl
  vpc_name    = "your_vpc_name"
  subnet_name = "your_subnet_name"
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
* This example creates a basic VPC and subnet with default configurations
* All resources will be created in the specified region
* The VPC and subnet are fundamental networking components for other cloud resources

## Requirements

| Name | Version |
| ---- | ---- |
| terraform | >= 0.12.0 |
| huaweicloud | >= 1.52.1 |
