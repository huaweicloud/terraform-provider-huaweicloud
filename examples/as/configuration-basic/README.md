# Create an Auto Scaling Configuration

This example provides best practice code for using Terraform to create an Auto Scaling configuration in HuaweiCloud.

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

* `security_group_name` - The name of the security group
* `kps_key_pair_name` - The name of the key pair
* `kps_public_key` - The public key for SSH access
* `as_configuration_name` - The name of the AS configuration
* `as_metadata` - The metadata for the instance
* `as_user_data` - The user data script for instance initialization
* `as_disks` - The disk configurations for the instance
* `as_public_ip` - The public IP configuration for the instance

## Usage

* Create a working directory and create a `versions.tf` file, the content is as follows:

```hcl
terraform {
  required_providers {
    huaweicloud = {
      source  = "huaweicloud/huaweicloud"
      version = ">= 1.57.0"
    }
  }
}
```

* Copy this example scripts (`main.tf` and `variables.tf`) to your working directory.

* Prepare the authentication (AK/SK and region) and configured in the TF script (versions.tf), also you can using
  environment variables.

```hcl
provider "huaweicloud" {
  region     = var.region_name
  access_key = var.access_key
  secret_key = var.secret_key
}

variable "region_name" {
  type = string
}

variable "access_key" {
  type = string
}

variable "secret_key" {
  type = string
}
```

* Create a `terraform.tfvars` [file](./terraform.tfvars) and fill in the required variables.

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

* The Auto Scaling configuration will be created with the specified instance configuration.
* Make sure to keep your credentials secure and never commit them to version control.
* All resources will be created in the specified region.

## Requirements

| Name | Version |
|------|---------|
| terraform | >= 0.12.0 |
| huaweicloud | >= 1.57.0 |
