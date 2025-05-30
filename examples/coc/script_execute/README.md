# Create a COC script and execute it

This example provides best practice code for using Terraform to create a COC script and execute it in HuaweiCloud.

## Prerequisites

* A HuaweiCloud account
* Terraform installed
* HuaweiCloud access key and secret key (AK/SK)

## Required Variables

The following variables need to be configured:

### Authentication Variables

* `region_name` - The region where the COC script is located
* `access_key`  - The access key of the IAM user
* `secret_key`  - The secret key of the IAM user

### Resource Variables

#### Required Variables

* `vpc_name` - The name of the VPC
* `subnet_name` - The name of the subnet
* `security_group_name` - The name of the security group
* `instance_name` - The name of the ECS instance
* `instance_user_data` - The user data for installing UniAgent on the ECS instance
* `script_name` - The name of the script
* `script_description` - The description of the script
* `script_risk_level` - The risk level of the script
* `script_version` - The version of the script
* `script_type` - The type of the script
* `script_content` - The content of the script
* `script_parameters` - The parameter list of the script.
  + `name` - The name of the parameter
  + `value` - The value of the parameter
  + `description` - The description of the parameter
  + `sensitive` - Whether the parameter is sensitive
* `script_execute_timeout` - The maximum time to execute the script in seconds
* `script_execute_execute_user` - The user to execute the script
* `script_execute_parameters` - The parameter list of the script execution.
  + `name` - The name of the parameter
  + `value` - The value of the parameter

#### Optional Variables

* `availability_zone` - The availability zone to which the ECS instance and network belong (default: "")
* `instance_flavor_id` - The flavor ID of the ECS instance (default: "")
* `instance_flavor_performance_type` - The performance type of the ECS instance flavor (default: "normal")
* `instance_flavor_cpu_core_count` - The number of the ECS instance flavor CPU cores (default: 4)
* `instance_flavor_memory_size` - The memory size of the ECS instance flavor (default: 8)
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
  vpc_name                    = "your_vpc_name"
  subnet_name                 = "your_subnet_name"
  security_group_name         = "your_security_group_name"
  instance_name               = "your_instance_name"
  instance_user_data          = "your_user_data"
  script_name                 = "your_script_name"
  script_description          = "your_script_description"
  script_risk_level           = "your_script_risk_level"
  script_version              = "your_script_version"
  script_type                 = "your_script_type"
  script_content              = "your_script_content"
  script_parameters           = "your_script_parameters"
  script_execute_timeout      = "your_script_execute_timeout"
  script_execute_execute_user = "your_script_execute_execute_user"
  script_execute_parameters   = "your_script_execute_parameters"
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
|------|---------|
| terraform | >= 0.12.0 |
| huaweicloud | >= 1.58.0 |
