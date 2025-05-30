# Create a COC script

This example provides best practice code for using Terraform to create a COC script in HuaweiCloud.

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

* `script_name` - The name of the COC script
* `script_description` - The description of the COC script
* `script_parameters` - The parameter list of the COC script

## Usage

* Copy this example script to your `main.tf`.

* Create a `terraform.tfvars` file and fill in the required variables:

  ```hcl
  access_key         = "your_access_key"
  secret_key         = "your_secret_key"
  region_name        = "your_region"
  script_name        = "example-script"
  script_description = "example-script-description"
  script_parameters = [
    {
      name        = "name"
      value       = "world"
      description = "the parameter"
    },
    {
      name        = "company"
      value       = "Huawei"
      description = "the second parameter"
      sensitive   = true
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
* The administrator password must meet the complexity requirements of HuaweiCloud
* All resources will be created in the specified region
