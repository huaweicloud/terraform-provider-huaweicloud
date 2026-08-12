# Create a DDM schema

This example provides best practice code for using Terraform to create a DDM schema in HuaweiCloud.

## Prerequisites

* A HuaweiCloud account
* Terraform installed
* HuaweiCloud access key and secret key (AK/SK)

## Required Variables

The following variables need to be configured:

### Authentication Variables

* `region_name` - The region where the DDM schema is located
* `access_key`  - The access key of the IAM user
* `secret_key`  - The secret key of the IAM user

### Resource Variables

#### Required Variables

* `vpc_name` - The name of the VPC
* `subnet_name` - The name of the subnet
* `security_group_name` - The name of the security group
* `rds_instance_name` - The name of the RDS instance
* `ddm_instance_name` - The name of the DDM instance
* `schema_name` - The name of the DDM schema

#### Optional Variables

* `availability_zones` - The availability zones to which the DDM instance belongs (default: [])
* `vpc_cidr` - The CIDR block of the VPC (default: "192.168.0.0/16")
* `subnet_cidr` - The CIDR block of the subnet (default: "")
* `subnet_gateway_ip` - The gateway IP of the subnet (default: "")
* `rds_instance_password` - The password of the RDS instance (default: "")
* `instance_flavor` - The flavor of the RDS instance (default: "")
* `database_type` - The database type of the RDS instance (default: "MySQL")
* `database_version` - The database version of the RDS instance (default: "5.7")
* `instance_mode` - The mode of the RDS instance (default: "single")
* `instance_group_type` - The performance specification (default: "dedicated")
* `instance_flavor_vcpus` - The number of vCPUs for the RDS instance flavor (default: 2)
* `database_port` - The port of the RDS instance (default: 3306)
* `volume_type` - The volume type of the RDS instance (default: "CLOUDSSD")
* `volume_size` - The volume size of the RDS instance (default: 40)
* `instance_engine_id` - The engine ID of the DDM instance (default: "")
* `instance_flavor_id` - The flavor ID of the DDM instance (default: "")
* `instance_node_num` - The number of nodes in the DDM instance (default: 2)
* `instance_parameters` - The parameters of the DDM instance (default: [])
  - `name` - The parameter name
  - `value` - The parameter value
* `schema_shard_mode` - The shard mode of the DDM schema (default: "single")
* `schema_shard_number` - The number of shards in the same working mode (default: 1)

## Usage

* Copy this example script to your `main.tf`.

* Create a `terraform.tfvars` file and fill in the required variables:

  ```hcl
  vpc_name            = "your_vpc_name"
  subnet_name         = "your_subnet_name"
  security_group_name = "your_security_group_name"
  rds_instance_name   = "your_rds_instance_name"
  ddm_instance_name   = "your_ddm_instance_name"
  schema_name         = "your_schema_name"
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
