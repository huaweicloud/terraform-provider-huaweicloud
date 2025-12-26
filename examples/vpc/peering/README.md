# Create a VPC peering connection

This example provides best practice code for using Terraform to create a VPC peering connection in HuaweiCloud VPC
service. The peering connection is used to connect two VPCs using a flexible configuration approach.

## Prerequisites

* A HuaweiCloud account
* Terraform installed
* HuaweiCloud access key and secret key (AK/SK)

## Variable Introduction

The following variables need to be configured:

### Authentication Variables

* `region_name` - The region where the VPC peering connection is located
* `access_key` - The access key of the IAM user
* `secret_key` - The secret key of the IAM user

### Resource Variables

#### Required Variables

* `vpc_configurations` - The list of VPC configurations for peering connection (exactly 2 VPCs required)
* `peering_connection_name` - The name of the VPC peering connection

#### VPC Configuration Structure

Each VPC configuration in the list should contain:

```hcl
{
  vpc_name             = "your_vpc_name"
  vpc_cidr             = "192.168.0.0/18"
  subnet_name          = "your_subnet_name"
  enterprise_project_id = "your_enterprise_project_id" # optional
}
```

## Usage

* Copy this example script to your `main.tf`.

* Create a `terraform.tfvars` file and fill in the required variables:

  ```hcl
  vpc_configurations = [
    {
      vpc_name    = "source_vpc"
      vpc_cidr    = "192.168.0.0/18"
      subnet_name = "source_subnet"
    },
    {
      vpc_name    = "target_vpc"
      vpc_cidr    = "192.168.128.0/18"
      subnet_name = "target_subnet"
    }
  ]
  peering_connection_name = "your_peering_connection_name"
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
* This example creates exactly 2 VPCs with their respective subnets and establishes a peering connection between them
* The peering connection allows communication between the two VPCs
* Route tables are automatically configured to enable traffic flow between the VPCs
* All resources will be created in the specified region
* Exactly 2 VPC configurations are required for peering connection

## Requirements

| Name | Version |
| ---- | ---- |
| terraform | >= 1.3.0 |
| huaweicloud | >= 1.19.1 |
