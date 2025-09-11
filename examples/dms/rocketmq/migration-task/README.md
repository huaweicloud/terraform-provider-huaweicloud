# Create a DMS RocketMQ migration task example

This example provides best practice code for using Terraform to migrate metadata to HuaweiCloud DMS RocketMQ instance.

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
* `migration_task_overwrite` - Whether to overwrite existing configurations with the same name
* `migration_task_name` - The name of the migration task
* `migration_task_type` - The type of the migration task  
  Valid values are **rocketmq** and **rabbitToRocket**

#### Optional Variables

* `vpc_cidr` - The CIDR block of the VPC (default: "192.168.0.0/16")
* `subnet_cidr` - The CIDR block of the subnet (default: "")
* `subnet_gateway_ip` - The gateway IP of the subnet (default: "")
* `availability_zones` - The list of the availability zones to which the RocketMQ instance belongs (default: [])  
  If this parameter is not specified, the availability zone is automatically assigned based on the value
  of `availability_zones_count`
* `instance_flavor_id` - The flavor ID of the RocketMQ instance (default: "")
* `instance_flavor_type` - The type of the RocketMQ instance flavor (default: "cluster.small")
* `availability_zones_count` - The number of availability zones (default: 1)
* `instance_engine_version` - The engine version of the RocketMQ instance (default: "")  
  When `instance_flavor_id` is omitted, `instance_engine_version` is required
* `instance_storage_spec_code` - The storage spec code of the RocketMQ instance (default: "")  
  When `instance_flavor_id` is omitted, `instance_storage_spec_code` is required
* `instance_storage_space` - The storage space of the RocketMQ instance (default: 800)
* `instance_broker_num` - The number of the broker of the RocketMQ instance (default: 0)  
  For `4.8.0` version instance, `instance_broker_num` is required and valid
* `instance_description` - The description of the RocketMQ instance (default: "")
* `instance_tags` - The tags of the RocketMQ instance (default: {})
* `enterprise_project_id` - The enterprise project ID of the RocketMQ instance (default: null)
* `instance_enable_acl` - Whether to enable the ACL of the RocketMQ instance (default: false)
* `instance_tls_mode` - The TLS mode of the RocketMQ instance (default: "SSL")
* `instance_configs` - The configs of the RocketMQ instance (default: [])
  - `name` - The name of the config
  - `value` - The value of the config
* `migration_task_topic_configs` - The topic configuration list of the migration task (default: [])  
  This parameter is required when `migration_task_type` is **rocketmq**
  - `topic_name` - The name of the topic (required)
  - `order` - Whether the topic is ordered (optional)
  - `perm` - The permission of the topic (optional)
  - `read_queue_nums` - The number of read queues (optional)
  - `write_queue_nums` - The number of write queues (optional)
  - `topic_filter_type` - The filter type of the topic (optional)
  - `topic_sys_flag` - The system flag of the topic (optional)
* `migration_task_subscription_groups` - The subscription group list of the migration task (default: [])  
  This parameter is required when `migration_task_type` is **rocketmq**
  - `group_name` - The name of the subscription group (required)
  - `consume_broadcast_enable` - Whether to enable broadcast consumption (optional)
  - `consume_enable` - Whether to enable consumption (optional)
  - `consume_from_min_enable` - Whether to enable consumption from minimum offset (optional)
  - `notify_consumerids_changed_enable` - Whether to enable notification when consumer IDs change (optional)
  - `retry_max_times` - The maximum number of retry times (optional)
  - `retry_queue_num` - The number of retry queues (optional)
  - `which_broker_when_consume_slow` - Which broker to use when consumption is slow (optional)
* `migration_task_vhosts` - The virtual host list of the migration task (default: [])  
  This parameter is required when `migration_task_type` is **rabbitToRocket**
  - `name` - The name of the virtual host (required)
* `migration_task_queues` - The queue list of the migration task (default: [])  
  This parameter is required when `migration_task_type` is **rabbitToRocket**
  - `name` - The name of the queue (optional)
  - `vhost` - The virtual host of the queue (optional)
  - `durable` - Whether the queue is durable (optional)
* `migration_task_exchanges` - The exchange list of the migration task (default: [])  
  This parameter is required when `migration_task_type` is **rabbitToRocket**
  - `name` - The name of the exchange (optional)
  - `vhost` - The virtual host of the exchange (optional)
  - `type` - The type of the exchange (optional)
  - `durable` - Whether the exchange is durable (optional)
* `migration_task_bindings` - The binding list of the migration task (default: [])  
  This parameter is required when `migration_task_type` is **rabbitToRocket**
  - `vhost` - The virtual host of the binding (optional)
  - `source` - The source exchange of the binding (optional)
  - `destination` - The destination queue of the binding (optional)
  - `destination_type` - The type of the destination (optional)
  - `routing_key` - The routing key of the binding (optional)

## Usage

* Copy this example script to your `main.tf`.

* Create a `terraform.tfvars` file and fill in the required variables:

  ```hcl
  vpc_name                 = "your_vpc_name"
  subnet_name              = "your_subnet_name"
  security_group_name      = "your_security_group_name"
  instance_name            = "your_rocketmq_instance"
  migration_task_overwrite = "true"
  migration_task_name      = "your_migration_task"
  migration_task_type      = "your_migration_task_type"
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

## Two types of migration task configuration options

### 1. RocketMQ to RocketMQ Migration

Migrate metadata from other RocketMQ instance to HuaweiCloud RocketMQ instance

```hcl
migration_task_type                = "rocketmq"
migration_task_topic_configs       = "your_topic_configs"
migration_task_subscription_groups = "your_subscription_groups"
```

### 2. RabbitMQ to RocketMQ Migration

Migrate metadata from RabbitMQ to HuaweiCloud RocketMQ instance

```hcl
migration_task_type      = "rabbitToRocket"
migration_task_vhosts    = "your_vhosts"
migration_task_queues    = "your_queues"
migration_task_exchanges = "your_exchanges"
migration_task_bindings  = "your_bindings"
```

## Note

* Make sure to keep your credentials secure and never commit them to version control
* All resources will be created in the specified region

## Requirements

| Name | Version |
| ---- | ---- |
| terraform | >= 1.9.0 |
| huaweicloud | >= 1.77.6 |
