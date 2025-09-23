# Execute COC script order operation

This example provides best practice code for using Terraform to create and operate a COC script in HuaweiCloud.

## Prerequisites

* A HuaweiCloud account
* Terraform installed
* HuaweiCloud access key and secret key (AK/SK)

## Variable Introduction

The following variables need to be configured:

### Authentication Variables

* `access_key` - HuaweiCloud access key
* `secret_key` - HuaweiCloud secret key
* `region_name` - The region where resources will be created

### Resource Variables

#### Required Variables

* `vpc_name` - The name of the VPC
* `subnet_name` - The name of the subnet
* `security_group_name` - The name of the security group
* `instance_name` - The name of the ECS instance
* `instance_user_data` - The user data for installing UniAgent on the ECS instance
* `script_name` - The name of the COC script
* `script_description` - The description of the script
* `script_risk_level` - The risk level of the script
* `script_version` - The version of the script
* `script_type` - The type of the script
* `script_content` - The content of the script
* `script_parameters` - The parameter list of the script
  - `name` - The name of the parameter
  - `value` - The value of the parameter
  - `description` - The description of the parameter
  - `sensitive` - Whether the parameter is sensitive
* `script_execute_timeout` - The maximum time to execute the script in seconds
* `script_execute_user` - The user to execute the script
* `script_execute_parameters` - The parameter list of the script execution
  - `name` - The name of the parameter
  - `value` - The value of the parameter
* `script_order_operation_batch_index` - The batch index for the script order
* `script_order_operation_instance_id` - The instance ID for the script order
* `script_order_operation_type` - The operation type for the script order

#### Optional Variables

* `availability_zone` - The availability zone to which the ECS instance and network belong (default: "")
* `instance_flavor_id` - The flavor ID of the ECS instance (default: "")
* `instance_flavor_performance_type` - The performance type of the ECS instance flavor (default: "normal")
* `instance_flavor_cpu_core_count` - The number of the ECS instance flavor CPU cores (default: 2)
* `instance_flavor_memory_size` - The number of the ECS instance flavor memories (default: 4)
* `instance_image_id` - The image ID of the ECS instance (default: "")
* `instance_image_os_type` - The OS type of the ECS instance flavor (default: "Ubuntu")
* `instance_image_visibility` - The visibility of the ECS instance flavor (default: "public")
* `vpc_cidr` - The CIDR block of the VPC (default: "192.168.0.0/16")
* `subnet_cidr` - The CIDR block of the subnet (default: "")
* `subnet_gateway_ip` - The gateway IP of the subnet (default: "")

## Usage

* Copy this example script to your `main.tf`.

* Create a `terraform.tfvars` file and fill in the required variables:

  ```hcl
  vpc_name                           = "your_vpc_name"
  subnet_name                        = "your_subnet_name"
  security_group_name                = "your_security_group_name"
  instance_name                      = "your_instance_name"
  instance_user_data                 = "your_user_data"
  script_name                        = "your_script_name"
  script_description                 = "your_script_description"
  script_risk_level                  = "your_script_risk_level"
  script_version                     = "your_script_version"
  script_type                        = "your_script_type"
  script_content                     = "your_script_content"
  script_parameters                  = "your_script_parameters"
  script_execute_timeout             = "your_script_execute_timeout"
  script_execute_execute_user        = "your_script_execute_execute_user"
  script_execute_parameters          = "your_script_execute_parameters"
  script_order_operation_batch_index = "your_script_order_operation_batch_index"
  script_order_operation_instance_id = "your_script_order_operation_instance_id"
  script_order_operation_type        = "your_script_order_operation_type"
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
* The ECS instance needs to have UniAgent installed for script execution
* All resources will be created in the specified region

## Requirements

| Name | Version |
|------|---------|
| terraform | >= 0.12.0 |
| huaweicloud | >= 1.75.0 |
