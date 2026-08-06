# Create a DDS cluster instance

This example provides best practice code for using Terraform to create a DDS cluster instance in HuaweiCloud.

## Prerequisites

* A HuaweiCloud account
* Terraform installed
* HuaweiCloud access key and secret key (AK/SK)

## Required Variables

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
* `instance_name` - The DDS instance name
* `instance_flavors` - The list of node flavor configurations for DDS instance
  - `type` - The node type of the DDS instance, the value can be `mongos`, `shard`, or `config`
  - `num` - The node  quantity of the DDS instance
    + If the value of type is `mongos` or `shard`, the value ranges from `2` to `16`
    + If the value of type is `config`, the value can only be `1`
  - `spec_code` - The node specification code of the DDS instance
  - `storage` - The disk type (default: "")
    + The value can be `ULTRAHIGH` or `EXTREMEHIGH`
    + This parameter is valid for the `shard` and `config` nodes
  - `size` - The disk size (default: null)
    + The value must be a multiple of `10`
    + This parameter is valid for the `shard` and `config` nodes
  - `node_list` - The ID list of instance nodes to be deleted (default: null)
    + This parameter is available only when you delete the mongos nodes

#### Optional Variables

* `availability_zone` - The availability zone to which the DDS instance belongs (default: "")
* `vpc_cidr` - The CIDR block of the VPC (default: "192.168.0.0/16")
* `subnet_cidr` - The CIDR block of the subnet (default: "")
* `subnet_gateway_ip` - The gateway IP of the subnet (default: "")
* `database_type` - The database version type of the DDS instance (default: "DDS-Community")
* `database_version` - The database version of the DDS instance, the value can be `4.0`, `4.2`, `4.4` or
  `5.0` (default: "4.0")
* `storage_engine` - The storage engine of the DDS instance (default: "wiredTiger")
  If `database_version` is set to `4.0`, the value is `wiredTiger`.
  If `database_version` is set to `4.2`, `4.4` or `5.0`, the value is `rocksDB`.
* `instance_port` - The database access port of the DDS instance (default: 8635)
* `instance_password` - The database access password of the DDS instance (default: "")
* `instance_description` - The description of the DDS instance (default: "")
* `instance_tags` - The tags of the DDS instance (default: {})
* `charging_mode` - The charging mode of the DDS instance (default: "postPaid")
* `period_unit` - The period unit of the DDS instance, only required when `charging_mode` is `prePaid` (default: null)
* `period` - The period of the DDS instance, only required when `charging_mode` is `prePaid` (default: null)
* `auto_renew` - The auto renew of the DDS instance (default: "false")

## Usage

* Copy this example script to your `main.tf`.

* Create a `terraform.tfvars` file and fill in the required variables:

  ```hcl
  vpc_name            = "your_vpc_name"
  subnet_name         = "your_subnet_name"
  security_group_name = "your_security_group_name"
  instance_name       = "your_instance_name"
  instance_flavors    = [
    {
      type      = "mongos"
      num       = 2
      spec_code = "dds.mongodb.s6.large.2.mongos"
    },
    {
      type      = "shard"
      num       = 2
      spec_code = "dds.mongodb.s6.large.2.shard"
      storage   = "ULTRAHIGH"
      size      = 20
    },
    {
      type      = "config"
      num       = 1
      spec_code = "dds.mongodb.s6.large.2.config"
      storage   = "ULTRAHIGH"
      size      = 20
    }
  ]
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
* The creation of the DDS instance takes about 20 to 50 minutes depending on the flavor and node number
* All resources will be created in the specified region

## Requirements

| Name | Version |
| ---- | ---- |
| terraform | >= 1.3.0 |
| huaweicloud | >= 1.87.0 |
