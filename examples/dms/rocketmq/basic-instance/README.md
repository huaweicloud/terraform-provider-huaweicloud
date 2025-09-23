# Create a DMS RocketMQ basic instance

This example provides best practice code for using Terraform to create a DMS RocketMQ basic instance in HuaweiCloud.

## Prerequisites

* A HuaweiCloud account
* Terraform installed
* HuaweiCloud access key and secret key (AK/SK)

## Variable Introduction

The following variables need to be configured:

### Authentication Variables

* `region_name` - The region where the RocketMQ instance is located
* `access_key`  - The access key of the IAM user
* `secret_key`  - The secret key of the IAM user

### Resource Variables

#### Required Variables

* `vpc_name` - The VPC name
* `subnet_name` - The subnet name
* `security_group_name` - The security group name
* `instance_name` - The RocketMQ instance name

#### Optional Variables

* `availability_zones` - The list of the availability zones to which the RocketMQ instance belongs (default: [])  
  If this parameter is not specified, the availability zone is automatically assigned based on the value of `availability_zones_count`
* `availability_zones_count` - The number of availability zones (default: 1)
* `vpc_cidr` - The CIDR block of the VPC (default: "192.168.0.0/16")
* `subnet_cidr` - The CIDR block of the subnet (default: "")
* `subnet_gateway_ip` - The gateway IP of the subnet (default: "")
* `instance_flavor_id` - The flavor ID of the RocketMQ instance (default: "")
* `instance_flavor_type` - The type of the RocketMQ instance flavor (default: "cluster.small")
* `instance_enable_publicip` - Whether to enable the public IP of the RocketMQ instance (default: false)
* `instance_publicip_id` - The ID of the public IP of the RocketMQ instance, multiple IDs separated by commas (,), only
  required and valid when `instance_enable_publicip` is true (default: "")
* `instance_eips_count` - The number of the public IP of the RocketMQ instance, only required and valid when
  `instance_enable_publicip` is true and `instance_publicip_id` is not specified (default: 0)
* `eip_type` - The type of the EIP, only available when `instance_enable_publicip` is true (default: "5_bgp")
* `bandwidth_name` - The name of the EIP bandwidth, only required and valid when `instance_enable_publicip` is true
  and `instance_eips_count` is greater than 0 (default: "")
* `bandwidth_size` - The size of the EIP bandwidth (default: 5)
* `bandwidth_share_type` - The share type of the EIP bandwidth (default: "PER")
* `bandwidth_charge_mode` - The charge mode of the EIP bandwidth (default: "traffic")
* `instance_engine_version` - The engine version of the RocketMQ instance (default: "")  
  When `instance_flavor_id` is omitted, `instance_engine_version` is required
* `instance_storage_spec_code` - The storage spec code of the RocketMQ instance (default: "")  
  When `instance_flavor_id` is omitted, `instance_storage_spec_code` is required
* `instance_storage_space` - The storage space of the RocketMQ instance (default: 800)
* `instance_broker_num` - The number of the broker of the RocketMQ instance (default: 1)
* `instance_description` - The description of the RocketMQ instance (default: "")
* `instance_tags` - The tags of the RocketMQ instance (default: {})
* `enterprise_project_id` - The enterprise project ID of the RocketMQ instance (required for
  enterprise users, default: null)
* `instance_enable_acl` - Whether to enable the ACL of the RocketMQ instance (default: false)
* `instance_tls_mode` - The TLS mode of the RocketMQ instance (default: "SSL")
* `instance_configs` - The configs of the RocketMQ instance (default: [])
* `charging_mode` - The charging mode of the RocketMQ instance (default: "postPaid")
* `period_unit` - The period type of the RocketMQ instance (default: null)
* `period` - The period of the RocketMQ instance (default: null)
* `auto_renew` - Whether to enable the auto-renew of the RocketMQ instance (default: null)

## Usage

* Copy this example script to your `main.tf`.

* Create a `terraform.tfvars` file and fill in the required variables:

  ```hcl
  vpc_name            = "your_vpc_name"
  subnet_name         = "your_subnet_name"
  security_group_name = "your_security_group_name"
  instance_name       = "your_rocketmq_instance"
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

## Enable public IP options

### Option 1: Create new EIPs

If you don't provide an existing EIP address, the example will create new EIPs with the specified bandwidth configuration:

```hcl
instance_broker_num      = 2
instance_enable_publicip = true
instance_eips_count      = 6
bandwidth_name           = "your_bandwidth_name"
bandwidth_size           = 5
eip_type                 = "5_bgp"
```

### Option 2: Use existing EIPs

If you have existing EIPs, you can associate them directly:

```hcl
instance_enable_publicip = true
instance_publicip_id     = "your_existing_eip_address"
```

## Prepaid instance configuration

For prepaid RocketMQ instance, you must set the following variables:

```hcl
charging_mode = "prePaid"
period_unit   = "month"
period        = 1
```

## Note

* Make sure to keep your credentials secure and never commit them to version control
* The creation of the RocketMQ instance takes about 20 to 50 minutes depending on the flavor and broker number
* All resources will be created in the specified region

## Requirements

| Name | Version |
| ---- | ---- |
| terraform | >= 1.9.0 |
| huaweicloud | >= 1.77.6 |
