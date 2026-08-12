# Create a DDM account

This example provides best practice code for using Terraform to create a DDM account in HuaweiCloud.

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
* `account_name` - The DDM account name

#### Optional Variables

* `availability_zones` - The availability zones to which the DDM instance belongs (default: [])
* `vpc_cidr` - The CIDR block of the VPC (default: "192.168.0.0/16")
* `subnet_cidr` - The CIDR block of the subnet (default: "")
* `subnet_gateway_ip` - The gateway IP of the subnet (default: "")
* `instance_engine_id` - The engine ID of the DDM instance (default: "")
* `instance_flavor_id` - The flavor ID of the DDM instance (default: "")
* `instance_node_num` - The number of nodes in the DDM instance (default: 2)
* `instance_parameters` - The parameters of the DDM instance (default: [])
  - `name` - The parameter name
  - `value` - The parameter value
* `account_password` - The password of the DDM account (default: "")
* `account_permissions` - The basic permissions of the DDM account (default: ["SELECT"])
  The valid values are `CREATE`, `DROP`, `ALTER`, `INDEX`, `INSERT`, `DELETE`, `UPDATE` and `SELECT`
* `account_description` - The description of the DDM account (default: "")

## Usage

* Copy this example script to your `main.tf`.

* Create a `terraform.tfvars` file and fill in the required variables:

  ```hcl
  vpc_name            = "your_vpc_name"
  subnet_name         = "your_subnet_name"
  security_group_name = "your_security_group_name"
  instance_name       = "your_instance_name"
  account_name        = "your_account_name"
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
