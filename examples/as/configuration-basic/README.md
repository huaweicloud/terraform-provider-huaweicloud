# Create a basic AS configuration

This example provides best practice code for using Terraform to create a basic Auto Scaling (AS) configuration in HuaweiCloud.

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

* `security_group_name` - Name of the security group
* `key_pair_name` - Name of the key pair for SSH access
* `public_key` - Public key for the key pair (sensitive)
* `configuration_name` - Name of the AS configuration

## Usage

* Copy this example script to your `main.tf`.

* Create a `terraform.tfvars` file and fill in the required variables:

  ```hcl
  access_key           = "your_access_key"
  secret_key           = "your_secret_key"
  region_name          = "your_region"
  security_group_name  = "example-sg"
  key_pair_name        = "example-keypair"
  public_key           = "your_public_key"
  configuration_name   = "example-as-config"
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
* The administrator password must meet the complexity requirements of HuaweiCloud
* All resources will be created in the specified region
