# Create a DC global gateway example

This example provides best practice code for using Terraform to create a DC (Direct Connect) global gateway in HuaweiCloud.

## Prerequisites

* A HuaweiCloud account
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

* `global_gateway_name` - The name of the global gateway
* `bgp_asn` - The BGP ASN of the global DC gateway

#### Optional Variables

* `global_gateway_description` - The description of the global gateway (default: "Created by Terraform")
* `address_family` - The IP address family of the global DC gateway (default: "ipv4")
* `enterprise_project_id` - The enterprise project ID that the global DC gateway belongs to (default: "0")
* `global_gateway_tags` - The tags of the global gateway (default: { "Owner" = "terraform", "Env" = "test" })

## Usage

* Copy this example script to your `main.tf`.

* Create a `terraform.tfvars` file and fill in the required variables:

  ```hcl
  global_gateway_name = "your_global_gateway_name"
  bgp_asn             = your_bgp_asn
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
| terraform | >= 1.9.0 |
| huaweicloud | >= 1.77.6 |
