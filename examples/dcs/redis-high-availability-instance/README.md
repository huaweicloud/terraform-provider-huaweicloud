# Create a DCS Redis high availability instance

This example provides best practice code for using Terraform to create a Redis HA, Redis Cluster, Proxy Cluster or
read-write separation instance in HuaweiCloud DCS service.

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
* `instance_name` - The name of the Redis instance

#### Optional Variables

* `vpc_cidr` - The CIDR block of the VPC (default: "192.168.0.0/16")
* `subnet_cidr` - The CIDR block of the subnet (default: "")
* `subnet_gateway_ip` - The gateway IP address of the subnet (default: "")
* `availability_zones` - The availability zones to which the Redis instance belongs (default: [])  
  The first is the primary availability zone, and the second is the standby availability zone
* `instance_flavor_id` - The flavor ID of the Redis instance (default: "")
* `instance_cache_mode` - The cache mode of the Redis instance. (default: "ha")  
  The valid values are as follows:
  - **ha**: Master/Standby
  - **cluster**: Redis Cluster
  - **proxy**: Proxy Cluster (Not supported in Redis 7.0)
  - **ha_rw_split**: Read/Write splitting (Not supported in Redis 7.0)
* `instance_capacity` - The capacity of the Redis instance (default: 4)
* `instance_engine_version` - The engine version of the Redis instance (default: "5.0")
* `enterprise_project_id` - The ID of the enterprise project to which the Redis instance belongs (default: null)
* `instance_password` - The password for the Redis instance (default: null)
* `instance_backup_policy` - The backup policy of the Redis instance (default: null)
  - `backup_type` - The type of the backup (default: "auto")
  - `backup_at` - The time of the backup
  - `begin_at` - The begin time of the backup
  - `save_days` - The number of days to save the backup (default: null)
  - `period_type` - The period type of the backup (default: null)
* `instance_whitelists` - The whitelists of the Redis instance (default: [])
  - `group_name` - The name of the whitelist group
  - `ip_address` - The IP address of the whitelist
* `instance_parameters` - The parameters of the Redis instance (default: [])
  - `id` - The ID of the parameter
  - `name` - The name of the parameter
  - `value` - The value of the parameter
* `instance_tags` - The tags of the Redis instance (default: {})
* `instance_rename_commands` - The rename commands of the Redis instance (default: {})
* `charging_mode` - The charging mode of the Redis instance (default: "postPaid")
* `period_unit` - The unit of the period, only used when `charging_mode` is `prePaid` (default: null)
* `period` - The period of the Redis instance, only used when `charging_mode` is `prePaid` (default: null)
* `auto_renew` - Whether auto renew is enabled (default: "false")

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

* To clean up the resources:

  ```bash
  $ terraform destroy
  ```

## Note

* Make sure to keep your credentials secure and never commit them to version control
* The creation of the DCS Redis instance may take several minutes
* All resources will be created in the specified region

## Requirements

| Name | Version |
| ---- | ------- |
| terraform | >= 1.3.0 |
| huaweicloud | >= 1.57.0 |
