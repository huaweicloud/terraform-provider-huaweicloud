# Produce messages to DMS Kafka topic

This example provides best practice code for producing messages to a DMS Kafka topic in HuaweiCloud.

## Prerequisites

* A HuaweiCloud account
* Terraform installed
* HuaweiCloud access key and secret key (AK/SK)

## Variables Introduction

The following variables need to be configured:

### Authentication Variables

* `region_name` - The region where the Kafka instance is located
* `access_key` - The access key of the IAM user
* `secret_key` - The secret key of the IAM user

### Resource Variables

#### Required Variables

* `vpc_name` - The VPC name
* `subnet_name` - The subnet name
* `security_group_name` - The security group name
* `instance_name` - The Kafka instance name
* `topic_name` - The name of the topic
* `message_body` - The body of the message to be sent

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
* `instance_access_user_name` - The access user of the Kafka instance (default: "")
* `instance_access_user_password` - The access password of the Kafka instance (default: "")
* `instance_enabled_mechanisms` - The enabled mechanisms of the Kafka instance (default: ["PLAIN"])
* `port_protocol` - The port protocol configuration of the Kafka instance (default: {private_plain_enable = true})
  - `private_plain_enable` - Whether to enable private plaintext access (optional)
  - `private_sasl_ssl_enable` - Whether to enable private SASL SSL access (optional)
  - `private_sasl_plaintext_enable` - Whether to enable private SASL plaintext access (optional)
  - `public_plain_enable` - Whether to enable public plaintext access (optional)
  - `public_sasl_ssl_enable` - Whether to enable public SASL SSL access (optional)
  - `public_sasl_plaintext_enable` - Whether to enable public SASL plaintext access (optional)
  
  -> Private access cannot be disabled. At least one of plaintext access or encrypted access must be enabled.
  
* `topic_partitions` - The number of partitions of the topic (default: 10)
* `topic_replicas` - The number of replicas of the topic (default: 3)
* `topic_aging_time` - The aging time of the topic (default: 72)
* `topic_sync_replication` - The sync replication of the topic (default: false)
* `topic_sync_flushing` - The sync flushing of the topic (default: false)
* `topic_description` - The description of the topic (default: null)
* `topic_configs` - The configs of the topic (default: [])
* `message_properties` - The properties of the message to be sent (default: [])

## Usage

* Copy this example script to your `main.tf`.

* Create a `terraform.tfvars` file and fill in the required variables:

  ```hcl
  vpc_name            = "your_vpc_name"
  subnet_name         = "your_subnet_name"
  security_group_name = "your_security_group_name"
  instance_name       = "your_kafka_instance"
  topic_name          = "your_topic_name"
  message_body        = "Hello Kafka!"
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
* All resources will be created in the specified region

## Requirements

| Name | Version |
| ---- | ---- |
| terraform | >= 1.3.0 |
| huaweicloud | >= 1.77.7 |
