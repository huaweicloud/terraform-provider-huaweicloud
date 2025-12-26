# Create a RMS policy assignment example

This example provides best practice code for using Terraform to create an RMS (Resource Management Service)
policy assignment in HuaweiCloud to check if ECS instances have required tags.

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

* `vpc_name` - The name of the VPC
* `subnet_name` - The name of the subnet
* `security_group_name` - The name of the security group
* `ecs_instance_name` - The name of the ECS instance
* `availability_zone` - The availability zone where the ECS instance is located
* `policy_assignment_name` - The name of the policy assignment

#### Optional Variables

* `vpc_cidr` - The CIDR block of the VPC (default: "192.168.0.0/16")
* `subnet_cidr` - The CIDR block of the subnet (default: "")
* `subnet_gateway_ip` - The gateway IP of the subnet (default: "")
* `ecs_image_name` - The name of the image used to create the ECS instance (default: "Ubuntu 20.04 server 64bit")
* `ecs_flavor_name` - The flavor name of the ECS instance (default: "s6.small.1")
* `ecs_tags` - The tags of the ECS instance (default: { "Owner" = "terraform", "Env" = "test" })
* `policy_assignment_description` - The description of the policy assignment (default: "Check if ECS instances have
  required tags")
* `policy_definition_id` - The ID of the policy definition
* `policy_assignment_policy_filter` - The configuration used to filter resources
  - `region` - The name of the region to which the filtered resources belong
  - `resource_provider` - The service name to which the filtered resources belong
  - `resource_type` - The resource type of the filtered resources
  - `resource_id` - The resource ID used to filter a specified resource
  - `tag_key` - The tag name used to filter resources
  - `tag_value` - The tag value used to filter resources
* `policy_assignment_parameters` - The rule definition of the policy assignment
* `policy_assignment_tags` - The tags of the policy assignment (default: { "Owner" = "terraform", "Env" = "test" })

## Usage

* Copy this example script to your `main.tf`.

* Create a `terraform.tfvars` file and fill in the required variables:

  ```hcl
  vpc_name              = "your_vpc_name"
  subnet_name           = "your_subnet_name"
  security_group_name   = "your_security_group_name"
  ecs_instance_name     = "your_ecs_instance_name"
  availability_zone     = "your_availability_zone"
  policy_assignment_name = "your_policy_assignment_name"
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
* The policy assignment will check if ECS instances have the required tags specified in `policy_parameters`
* All resources will be created in the specified region

## Requirements

| Name | Version |
| ---- | ---- |
| terraform | >= 1.1.0 |
| huaweicloud | >= 1.71.1 |
