# Create a DMS RocketMQ message send example

This example provides best practice code for using Terraform to create a DMS RocketMQ instance and
send a message in HuaweiCloud.

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
* `topic_name` - The name of the topic
* `message_body` - The body of the message to be sent

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
* `topic_brokers` - The broker list of the topic, this parameter is valid only when the RocketMQ instance
  version is `4.8.0` (default: [])
* `topic_message_type` - The message type of the topic, this parameter is valid only when the RocketMQ instance
  version is `5.x` (default: null)
* `topic_queue_num` - The number of the queue of the topic, this parameter is valid only when the RocketMQ
  instance version is `4.8.0` (default: 0)
* `topic_permission` - The permission of the topic, this parameter is valid only when the RocketMQ instance
  version is `4.8.0` (default: null)
* `message_properties` - The property list of the message to be sent (default: [])
  - `name` - The name of the property
  - `value` - The value of the property

## Usage

* Copy this example script to your `main.tf`.

* Create a `terraform.tfvars` file and fill in the required variables:

  ```hcl
  vpc_name            = "your_vpc_name"
  subnet_name         = "your_subnet_name"
  security_group_name = "your_security_group_name"
  instance_name       = "your_rocketmq_instance"
  topic_name          = "your_topic_name"
  message_body        = "your_message_body"
  instance_broker_num = 1
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

## Topic Configuration Options by Version

### For RocketMQ 4.8.0

When using RocketMQ 4.8.0, you can configure the following topic variables:

```hcl
topic_queue_num  = 3
topic_permission = "all"
topic_brokers    = ["broker-0", "broker-1"]
```

### For RocketMQ 5.x

When using RocketMQ `5.x`, you can configure the following topic variables:

```hcl
topic_message_type = "NORMAL"
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
