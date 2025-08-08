# Create a DMS RabbitMQ basic instance

This example provides best practice code for using Terraform to create a DMS RabbitMQ basic instance in HuaweiCloud.

## Prerequisites

* A HuaweiCloud account
* Terraform installed
* HuaweiCloud access key and secret key (AK/SK)

## Required Variables

The following variables need to be configured:

### Authentication Variables

* `region_name` - The region where the RabbitMQ instance is located
* `access_key`  - The access key of the IAM user
* `secret_key`  - The secret key of the IAM user

### Resource Variables

#### Required Variables

* `vpc_name` - The VPC name
* `subnet_name` - The subnet name
* `security_group_name` - The security group name
* `instance_name` - The RabbitMQ instance name
* `instance_access_user_name` - The access user of the RabbitMQ instance
* `instance_password` - The access password of the RabbitMQ instance

#### Optional Variables

* `availability_zones` - The availability zones to which the RabbitMQ instance belongs (default: [])  
  If this parameter is not specified, the availability zone is automatically assigned based on the value of `availability_zone_number`
* `availability_zone_number` - The number of availability zones to which the RabbitMQ instance belongs, and this value
  must be `1` or `3`  (default: 1)  
* `vpc_cidr` - The CIDR block of the VPC (default: "192.168.0.0/16")
* `subnet_cidr` - The CIDR block of the subnet (default: "")
* `subnet_gateway_ip` - The gateway IP of the subnet (default: "")
* `instance_flavor_id` - The flavor ID of the RabbitMQ instance (default: "")
* `instance_flavor_type` - The flavor type of the RabbitMQ instance (default: "cluster")  
  The valid values are as follows:
  - **single** - The RabbitMQ single instance
  - **cluster** - The RabbitMQ cluster instance
* `instance_storage_spec_code` - The storage specification code of the RabbitMQ instance (default: "dms.physical.storage.ultra.v2")
* `instance_engine_version` - The engine version of the RabbitMQ instance (default: "3.12.13")
* `instance_storage_space` - The storage space of the RabbitMQ instance (default: 600)
* `instance_broker_num` - The number of brokers of the RabbitMQ instance (default: 3). For single instance,
   this value can only be `1`.
* `instance_ssl_enable` - Whether to enable SSL for the RabbitMQ instance (default: false)
* `instance_description` - The description of the RabbitMQ instance (default: "")
* `enterprise_project_id` - The enterprise project ID of the RabbitMQ instance (default: null)
* `instance_tags` - The tags of the RabbitMQ instance (default: {})
* `charging_mode` - The charging mode of the RabbitMQ instance (default: "postPaid")
* `period_unit` - The period unit of the RabbitMQ instance (default: null)
* `period` - The period of the RabbitMQ instance (default: null)
* `auto_renew` - The auto renew of the RabbitMQ instance (default: "false")

## Usage

* Copy this example script to your `main.tf`.

* Create a `terraform.tfvars` file and fill in the required variables:

  ```hcl
  vpc_name                  = "your_vpc_name"
  subnet_name               = "your_subnet_name"
  security_group_name       = "your_security_group_name"
  instance_name             = "your_rabbitmq_instance"
  instance_access_user_name = "your_access_user_name"
  instance_password         = "your_access_password"
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
* The creation of the RabbitMQ instance takes about 20 to 50 minutes depending on the flavor and broker number
* All resources will be created in the specified region

## Requirements

| Name | Version |
| ---- | ---- |
| terraform | >= 0.14.0 |
| huaweicloud | >= 1.69.1 |
