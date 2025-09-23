# Create a custom line

This example provides best practice code for using Terraform to create a custom line in HuaweiCloud DNS service.

## Prerequisites

* A HuaweiCloud account with DNS permissions
* Terraform installed
* HuaweiCloud access key and secret key (AK/SK)

## Variable Introduction

The following variables need to be configured:

### Authentication Variables

* `region_name` - The region where the DNS service is located
* `access_key` - The access key of the IAM user
* `secret_key` - The secret key of the IAM user

### Resource Variables

#### Required Variables

* `dns_custom_line_name` - The custom line name.  
* `dns_custom_line_ip_segments` - The IP address range.  

#### Optional Variables

## Usage

* Copy this example script to your `main.tf`.

* Create a `terraform.tfvars` file and fill in the required variables:

  ```hcl
  dns_custom_line_name        = "your_custom_line_name"
  dns_custom_line_ip_segments = ["100.100.100.102-100.100.100.102", "100.100.100.101-100.100.100.101"]
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

## Requirements

| Name | Version |
| ---- | ---- |
| terraform | >= 0.12.0 |
| huaweicloud | >= 1.73.5 |
