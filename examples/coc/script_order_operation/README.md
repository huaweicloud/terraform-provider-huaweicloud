# Execute COC script order operation

This example provides best practice code for using Terraform to create and operate a COC script in HuaweiCloud.

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

* `enterprise_project_id` - The ID of the enterprise project
* `vpc_name` - The name of the VPC
* `vpc_cidr` - The CIDR block of the VPC
* `subnet_name` - The name of the subnet
* `security_group_name` - The name of the security group
* `ecs_instance_name` - The name of the ECS instance
* `ecs_instance_user_data` - The user data for installing UniAgent on the ECS instance
* `coc_script_name` - The name of the COC script
* `coc_script_description` - The description of the script
* `coc_script_risk_level` - The risk level of the script
* `coc_script_version` - The version of the script
* `coc_script_type` - The type of the script
* `coc_script_content` - The content of the script
* `coc_script_parameters` - The parameter list of the script
  - `name` - The name of the parameter
  - `value` - The value of the parameter
  - `description` - The description of the parameter
  - `sensitive` - Whether the parameter is sensitive
* `coc_script_execute_timeout` - The maximum time to execute the script in seconds
* `coc_script_execute_user` - The user to execute the script
* `coc_script_execute_parameters` - The parameter list of the script execution
  - `name` - The name of the parameter
  - `value` - The value of the parameter
* `coc_script_order_operation_batch_index` - The batch index for the script order
* `coc_script_order_operation_instance_id` - The instance ID for the script order
* `coc_script_order_operation_type` - The operation type for the script order

## Usage

* Create a working directory and create a `versions.tf` file, the content is as follows:

```hcl
terraform {
  required_providers {
    huaweicloud = {
      source  = "huaweicloud/huaweicloud"
      version = ">=1.75.0"
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

* Make sure to keep your credentials secure and never commit them to version control.
* The ECS instance needs to have UniAgent installed for script execution.
* All resources will be created in the specified region.

## Requirements

| Name | Version |
|------|---------|
| terraform | >= 0.12.0 |
| huaweicloud | >= 1.75.0 |
