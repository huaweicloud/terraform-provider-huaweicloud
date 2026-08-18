# DCS Redis big key analysis

This example provides best practice code for using Terraform to create and manage a DCS Redis big key analysis.

## Prerequisites

* A HuaweiCloud account
* Terraform installed
* HuaweiCloud access key and secret key (AK/SK)

## Required Variables

The following variables need to be configured:

### Authentication Variables

* `region_name` - The region where the DCS Redis instance is located
* `access_key`  - The access key of the IAM user
* `secret_key`  - The secret key of the IAM user

### Resource Variables

#### Required Variables

* `vpc_name` - The name of the VPC
* `subnet_name` - The name of the subnet
* `instance_name` - The name of the Redis single instance

#### Optional Variables

* `vpc_cidr` - The CIDR block of the VPC (default: "192.168.0.0/16")
* `subnet_cidr` - The CIDR block of the subnet (default: "")
* `subnet_gateway_ip` - The gateway IP address of the subnet (default: "")
* `availability_zone` - The availability zone to which the Redis single instance belongs (default: "")
* `instance_flavor_id` - The flavor ID of the Redis single instance (default: "")
* `instance_capacity` - The capacity of the Redis instance (default: 4)
* `instance_engine_version` - The engine version of the Redis single instance (default: "5.0")
* `instance_password` - The password for the Redis instance (default: null)

## Usage

* Copy this example script to your `main.tf`.

* Create a `terraform.tfvars` file and fill in the required variables:

  ```hcl
  vpc_name      = "your_vpc_name"
  subnet_name   = "your_subnet_name"
  instance_name = "your_instance_name"
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

## Note

* Make sure to keep your credentials secure and never commit them to version control
* The only required parameter is `instance_id`, all other attributes are computed
* The resource creates a manual big key analysis (`scan_type` will be "manual")
* The create timeout is 30 minutes, the instance must be in RUNNING state before analysis
* The `status` attribute transitions from `waiting` to `running` to `success`
* The `keys` block contains analysis results: `name`, `type`, `shard`, `db`, `size`, `unit`
* The `shard` field is only populated for cluster instances
* All parameters are ForceNew, changing `instance_id` will recreate the analysis
* The resource supports import using `<instance_id>/<id>`
* Deleting this resource deletes the analysis record, not the big keys themselves

## Requirements

| Name        | Version   |
|-------------|-----------|
| terraform   | >= 1.1.0  |
| huaweicloud | >= 1.92.0 |
