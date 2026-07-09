# Create EIP on Shared Bandwidth

This example provides best practice code for using Terraform to create an EIP on shared bandwidth in HuaweiCloud EIP
service.

## Prerequisites

* A HuaweiCloud account
* Terraform installed
* HuaweiCloud access key and secret key (AK/SK)

## Variable Introduction

The following variables need to be configured:

### Authentication Variables

* `region_name` - The region where the EIP resources are located
* `access_key` - The access key of the IAM user
* `secret_key` - The secret key of the IAM user

### Resource Variables

#### Required Variables

* `bandwidth_name` - The name of the shared bandwidth

#### Optional Variables

* `enterprise_project_id` - The ID of the enterprise project (default: `""`)
* `bandwidth_size` - The size of the shared bandwidth in Mbit/s (default: 5)
* `bandwidth_charge_mode` - The charge mode of the shared bandwidth (default: "bandwidth")
* `bandwidth_type` - The type of the shared bandwidth (default: "share")
* `bandwidth_public_border_group` - The border group of the public IP (default: "center")
* `eip_type` - The type of the EIP (default: "5_bgp")
* `eip_description` - The description of the EIP (default: `""`)
* `eip_tags` - The tags of the EIP (default: `null`)

## Usage

* Copy this example script to your `main.tf`.

* Create a `terraform.tfvars` file and fill in the required variables:

  ```hcl
  bandwidth_name = "tf_test_shared_bandwidth"
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

## Notes

* Make sure to keep your credentials secure and never commit them to version control
* The shared bandwidth and EIP must be in the same region
* All resources will be created in the specified region

## Requirements

| Name | Version |
| ---- | ---- |
| terraform | >= 1.1.0 |
| huaweicloud | >= 1.92.0 |
