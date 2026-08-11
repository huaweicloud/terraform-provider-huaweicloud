# Create a DDM instance

This example provides best practice code for using Terraform to create a DDM instance in HuaweiCloud.

## Prerequisites

* A HuaweiCloud account
* Terraform installed
* HuaweiCloud access key and secret key (AK/SK)

## Required Variables

The following variables need to be configured:

### Authentication Variables

* `region_name` - The region where the DDM instance is located
* `access_key`  - The access key of the IAM user
* `secret_key`  - The secret key of the IAM user

### Resource Variables

#### Required Variables

* `vpc_name` - The VPC name
* `subnet_name` - The subnet name
* `security_group_name` - The security group name
* `instance_name` - The DDM instance name

#### Optional Variables

* `availability_zones` - The availability zones to which the DDM instance belongs (default: [])
* `vpc_cidr` - The CIDR block of the VPC (default: "192.168.0.0/16")
* `subnet_cidr` - The CIDR block of the subnet (default: "")
* `subnet_gateway_ip` - The gateway IP of the subnet (default: "")
* `instance_engine_id` - The engine ID of the DDM instance (default: "")
* `instance_flavor_id` - The flavor ID of the DDM instance (default: "")
* `instance_node_num` - The number of nodes in the DDM instance (default: 2)
* `instance_admin_user_name` - The administrator username of the DDM instance (default: "")
* `instance_admin_user_password` - The administrator password of the DDM instance (default: "")
* `instance_parameters` - The parameters of the DDM instance (default: [])
  - `name` - The parameter name
  - `value` - The parameter value
* `charging_mode` - The charging mode of the DDM instance (default: "postPaid")
* `period_unit` - The period unit of the DDM instance, only required when `charging_mode` is `prePaid` (default: null)
* `period` - The period of the DDM instance, only required when `charging_mode` is `prePaid` (default: null)
* `auto_renew` - The auto renew of the DDM instance (default: "false")

## Usage

* Copy this example script to your `main.tf`.

* Create a `terraform.tfvars` file and fill in the required variables:

  ```hcl
  vpc_name            = "your_vpc_name"
  subnet_name         = "your_subnet_name"
  security_group_name = "your_security_group_name"
  instance_name       = "your_instance_name"
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
* The creation of the DDM instance takes about 10 to 20 minutes
* All resources will be created in the specified region

## Requirements

| Name | Version |
| ---- | ---- |
| terraform | >= 1.1.0 |
| huaweicloud | >= 1.66.0 |
