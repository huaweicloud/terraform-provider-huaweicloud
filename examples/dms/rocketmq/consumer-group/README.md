# Create a DMS RocketMQ consumer group

This example provides best practice code for using Terraform to create a DMS RocketMQ consumer group in HuaweiCloud.

## Prerequisites

* A HuaweiCloud account
* Terraform installed
* HuaweiCloud access key and secret key (AK/SK)

## Variables Introdution

The following variables need to be configured:

### Authentication Variables

* `region_name` - The region where the RocketMQ instance is located
* `access_key` - The access key of the IAM user
* `secret_key` - The secret key of the IAM user

### Resource Variables

#### Required Variables

* `vpc_name` - The name of the VPC
* `subnet_name` - The name of the subnet
* `security_group_name` - The name of the security group
* `instance_name` - The name of the RocketMQ instance
* `consumer_group_name` - The name of the consumer group

#### Optional Variables

* `vpc_cidr` - The CIDR block of the VPC (default: "192.168.0.0/16")
* `subnet_cidr` - The CIDR block of the subnet (default: "")
* `subnet_gateway_ip` - The gateway IP of the subnet (default: "")
* `availability_zones` - The list of the availability zones to which the RocketMQ instance belongs (default: [])  
  If this parameter is not specified, the availability zone is automatically assigned based on the value of
  `availability_zones_count`
* `instance_flavor_id` - The flavor ID of the RocketMQ instance (default: "")
* `instance_flavor_type` - The type of the RocketMQ instance flavor (default: "cluster.small")
* `availability_zones_count` - The number of availability zones (default: 1)
* `instance_engine_version` - The engine version of the RocketMQ instance (default: "")
* `instance_storage_spec_code` - The storage spec code of the RocketMQ instance (default: "")
* `instance_storage_space` - The storage space of the RocketMQ instance (default: 800)
* `instance_broker_num` - The number of the broker of the RocketMQ instance (default: 0)  
  This parameter is required and valid when the RocketMQ instance version is `4.8.0`
* `instance_description` - The description of the RocketMQ instance (default: "")
* `instance_tags` - The tags of the RocketMQ instance (default: {})
* `enterprise_project_id` - The enterprise project ID of the RocketMQ instance (default: null)
* `instance_enable_acl` - Whether to enable the ACL of the RocketMQ instance (default: false)
* `instance_tls_mode` - The TLS mode of the RocketMQ instance (default: "SSL")
* `instance_configs` - The configs of the RocketMQ instance (default: [])
  - `name` - The name of the config
  - `value` - The value of the config
* `consumer_group_brokers` - The broker list of the consumer group, it's only valid when the RocketMQ instance
  version is `4.8.0` (default: [])
* `consumer_group_retry_max_times` - The retry max times of the consumer group (default: 16)
* `consumer_group_enabled` - Whether to enable the consumer group (default: true)
* `consumer_group_broadcast` - Whether to enable the broadcast of the consumer group (default: false)
* `consumer_group_description` - The description of the consumer group (default: "")
* `consumer_group_consume_orderly` - Whether to enable the consume orderly of the consumer group, it's only valid
  when the RocketMQ instance version is `5.x` (default: false)

## Usage

* Copy this example script to your `main.tf`.

* Create a `terraform.tfvars` file and fill in the required variables:

  ```hcl
  vpc_name            = "your_vpc_name"
  subnet_name         = "your_subnet_name"
  security_group_name = "your_security_group_name"
  instance_name       = "your_rocketmq_instance"
  consumer_group_name = "your_consumer_group_name"
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

## Configuration Options by Version

This example supports different RocketMQ versions with specific consumer group configuration options:

### For RocketMQ 4.8.0 instance

```hcl
consumer_group_brokers = ["broker-0", "broker-1"]
```

### For RocketMQ 5.x instance

```hcl
consumer_group_consume_orderly = true
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
