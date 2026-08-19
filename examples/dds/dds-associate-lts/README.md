# Create a DDS instance associate LTS

This example provides best practice code for using Terraform to create a DDS instance associate LTS in HuaweiCloud.

## Prerequisites

* A HuaweiCloud account
* Terraform installed
* HuaweiCloud access key and secret key (AK/SK)

## Variable Introduction

The following variables need to be configured:

### Authentication Variables

* `region_name` - The region where the DDS instance is located
* `access_key`  - The access key of the IAM user
* `secret_key`  - The secret key of the IAM user

### Resource Variables

#### Required Variables

* `vpc_name` - The VPC name
* `subnet_name` - The subnet name
* `security_group_name` - The security group name
* `group_name` - The LTS log group name
* `stream_name` - The LTS log stream name
* `instance_name` - The DDS instance name

#### Optional Variables

* `availability_zone` - The availability zone to which the DDS instance belongs (default: "")
* `vpc_cidr` - The CIDR block of the VPC (default: "192.168.0.0/16")
* `subnet_cidr` - The CIDR block of the subnet (default: "")
* `subnet_gateway_ip` - The gateway IP of the subnet (default: "")
* `group_log_expiration_days` - The log expiration time of the LTS log group (default: 30)
* `node_spec_code` - The node specification code of the DDS instance (default: "")
* `engine_name` - The DB engine name of the DDS instance (default: "DDS-Community")
* `flavor_vcpus` - The VCPUs of the flavor (default: 2)
* `flavor_memory` - The memory of the flavor (default: 4)
* `node_type` - The type of the DDS instance node (default: "replica")
* `instance_mode` - The mode of the DDS instance (default: "ReplicaSet")
* `database_version` - The database version of the DDS instance (default: "4.0")
* `storage_engine` - The storage engine of the DDS instance (default: "wiredTiger")
* `node_number` - The number of nodes of the DDS instance, the value can be `3`, `5` or `7` (default: 3)
* `node_storage_type` - The storage type of the DDS instance node (default: "ULTRAHIGH")
* `node_size` - The disk size of the node of the DDS instance, The value must be a multiple of `10` (default: 10)
* `node_list` - The node IDs to be deleted of the DDS instance, this parameter is available only
  when you delete a replica set instance nodes (default: null)
* `instance_port` - The database access port of the DDS instance (default: 8635)
* `instance_description` - The description of the DDS instance (default: "")
* `instance_password` - The database access password of the DDS instance (default: "")
* `instance_tags` - The tags of the DDS instance (default: {})
* `charging_mode` - The charging mode of the DDS instance (default: "postPaid")
* `period_unit` - The period unit of the DDS instance, only required when `charging_mode` is `prePaid` (default: null)
* `period` - The period of the DDS instance, only required when `charging_mode` is `prePaid` (default: null)
* `auto_renew` - The auto renew of the DDS instance (default: "false")
* `log_type` - The log type (default: "audit_log")

## Usage

* Copy this example script to your `main.tf`.

* Create a `terraform.tfvars` file and fill in the required variables:

  ```hcl
  vpc_name            = "your_vpc_name"
  subnet_name         = "your_subnet_name"
  security_group_name = "your_security_group_name"
  group_name          = "your_log_group_name"
  stream_name         = "your_log_stream_name"
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
* The creation of the DDS instance takes about 20 to 60 minutes depending on the flavor and node number
* All resources will be created in the specified region

## Requirements

| Name | Version |
| ---- | ---- |
| terraform | >= 1.1.0 |
| huaweicloud | >= 1.88.0 |
