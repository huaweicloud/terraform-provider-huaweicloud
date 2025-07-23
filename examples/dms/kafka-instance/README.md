# Create a DMS Kafka instance

This example provides best practice code for using Terraform to create a DMS Kafka instance in HuaweiCloud.

## Prerequisites

* A HuaweiCloud account
* Terraform installed
* HuaweiCloud access key and secret key (AK/SK)

## Required Variables

The following variables need to be configured:

### Authentication Variables

* `region_name` - The region where the Kafka instance is located
* `access_key`  - The access key of the IAM user
* `secret_key`  - The secret key of the IAM user

### Resource Variables

#### Required Variables

* `vpc_name` - The VPC name
* `subnet_name` - The subnet name
* `security_group_name` - The security group name
* `instance_name` - The Kafka instance name

#### Optional Variables

* `availability_zones` - The availability zones to which the Kafka instance belongs (default: [])
* `vpc_cidr` - The CIDR block of the VPC (default: "192.168.0.0/16")
* `subnet_cidr` - The CIDR block of the subnet (default: "")
* `subnet_gateway_ip` - The gateway IP of the subnet (default: "")
* `instance_flavor_id` - The flavor ID of the Kafka instance (default: "")
* `instance_flavor_type` - The flavor type of the Kafka instance (default: "cluster")
* `instance_storage_spec_code` - The storage specification code of the Kafka instance (default: "dms.physical.storage.ultra.v2")
* `instance_engine_version` - The engine version of the Kafka instance (default: "2.7")
* `instance_storage_space` - The storage space of the Kafka instance (default: 600)
* `instance_broker_num` - The number of brokers of the Kafka instance (default: 3)
* `instance_ssl_enable` - Whether to enable SSL (default: false)
* `instance_access_user_name` - The access user of the Kafka instance, only required when `instance_ssl_enable`
  is `true` (default: "")
* `instance_access_user_password` - The access password of the Kafka instance, only required when `instance_ssl_enable`
  is `true` (default: "")
* `instance_description` - The description of the Kafka instance (default: "")
* `charging_mode` - The charging mode of the Kafka instance (default: "postPaid")
* `period_unit` - The period unit of the Kafka instance, only required when `charging_mode` is `prePaid` (default: null)
* `period` - The period of the Kafka instance, only required when `charging_mode` is `prePaid` (default: null)
* `auto_renew` - The auto renew of the Kafka instance (default: "false")

## Usage

* Copy this example script to your `main.tf`.

* Create a `terraform.tfvars` file and fill in the required variables:

  ```hcl
  vpc_name            = "your_vpc_name"
  subnet_name         = "your_subnet_name"
  security_group_name = "your_security_group_name"
  instance_name       = "your_kafka_instance"
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
* The creation of the Kafka instance takes about 20 to 50 minutes depending on the flavor and broker number
* All resources will be created in the specified region

## Requirements

| Name | Version |
| ---- | ---- |
| terraform | >= 0.12.0 |
| huaweicloud | >= 1.64.4 |
