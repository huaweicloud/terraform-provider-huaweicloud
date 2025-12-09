# Create an APIG instance with Kafka forward plugin

This example provides best practice code for using Terraform to create an API Gateway (APIG) instance with a Kafka forward
plugin in HuaweiCloud. This plugin enables asynchronous message processing by forwarding HTTP API requests to Kafka topics.

## Prerequisites

* A HuaweiCloud account
* HuaweiCloud access key and secret key (AK/SK)

## Variables Introduction

The following variables need to be configured:

### Authentication Variables

* `region_name` - The region where resources will be created
* `access_key` - The access key of the IAM user
* `secret_key` - The secret key of the IAM user

### Resource Variables

#### Required Variables

* `vpc_name` - The name of the VPC
* `subnet_name` - The name of the subnet
* `security_group_name` - The name of the security group
* `instance_name` - The name of the APIG instance
* `plugin_name` - The name of the Kafka forward plugin
* `kafka_instance_name` - The name of the DMS Kafka instance
* `kafka_topic_name` - The name of the Kafka topic to receive messages
* `kafka_broker_list` - The broker list for the Kafka instance. Format: host1:port1,host2:port2
* `kafka_instance_storage_spec_code` - The storage spec code of the DMS Kafka instance
* `kafka_instance_engine_version` - The engine version of the DMS Kafka instance
* `kafka_instance_storage_space` - The storage space of the DMS Kafka instance in GB
* `kafka_instance_broker_num` - The number of brokers for the DMS Kafka instance
* `kafka_instance_user_name` - The access user name for the DMS Kafka instance
* `kafka_instance_user_password` - The access user password for the DMS Kafka instance

#### Optional Variables

* `vpc_cidr` - The CIDR block of the VPC (default: "192.168.0.0/16")
* `subnet_cidr` - The CIDR block of the subnet (default: "", auto-calculated if empty)
* `subnet_gateway_ip` - The gateway IP address of the subnet (default: "", auto-calculated if empty)
* `availability_zones` - The availability zones to which the instance belongs (default: [])
  If not specified, will be automatically allocated based on the number of availability_zones_count
* `availability_zones_count` - The number of availability zones to which the instance belongs (default: 1)
* `instance_edition` - The edition of the APIG instance (default: "BASIC")
* `enterprise_project_id` - The ID of the enterprise project (default: null)
* `plugin_description` - The description of the Kafka forward plugin (default: null)
* `kafka_instance_description` - The description of the DMS Kafka instance (default: "")
* `kafka_instance_flavor_id` - The flavor ID of the DMS Kafka instance (default: "", auto-selected if empty)
* `kafka_instance_flavor_type` - The flavor type of the DMS Kafka instance (default: "cluster")
* `kafka_instance_ssl_enable` - Whether to enable SSL for the DMS Kafka instance (default: false)
* `kafka_charging_mode` - The charging mode. Options: prePaid, postPaid (default: "prePaid")
* `kafka_period_unit` - The period unit. Options: month, year (default: "month")
* `kafka_period` - The period (default: 1)
* `kafka_auto_new` - Whether to enable auto renewal (default: "false")
* `kafka_topic_partitions` - The number of partitions for the Kafka topic (default: 1)
* `kafka_message_key` - The message key extraction strategy (default: ""). Can be a static value or a variable
  expression like "$context.requestId"
* `kafka_max_retry_count` - The maximum number of retry attempts for failed message sends (default: 3)
* `kafka_retry_backoff` - The backoff time in seconds between retries (default: 10)
* `kafka_security_protocol` - The security protocol for Kafka connection (default: "PLAINTEXT"). Options: PLAINTEXT,
  SASL_PLAINTEXT, SASL_SSL, SSL
* `kafka_sasl_mechanisms` - The SASL mechanism for authentication (default: "PLAIN"). Options: PLAIN, SCRAM-SHA-256,
  SCRAM-SHA-512
* `kafka_sasl_username` - The SASL username for authentication (default: "", uses kafka_access_user if empty and
  security_protocol is not PLAINTEXT)
* `kafka_sasl_password` - The SASL password for authentication (default: "", uses kafka_password if empty and
  security_protocol is not PLAINTEXT)
* `kafka_ssl_ca_content` - The SSL CA certificate content (default: "", sensitive)
* `kafka_access_user` - The access user for Kafka authentication (default: "", used when kafka_sasl_username is empty,
  sensitive)
* `kafka_password` - The password for Kafka authentication (default: "", used when kafka_sasl_password is empty,
  sensitive)

## Usage

* Copy this example script to your `main.tf`.

* Create a `terraform.tfvars` file and fill in the required variables. Example:

  ```hcl
  vpc_name                         = "your_vpc_name"
  subnet_name                      = "your_subnet_name"
  security_group_name              = "your_security_group_name"
  instance_name                    = "your_apig_instance_name"
  plugin_name                      = "your_plugin_name"
  kafka_instance_name              = "your_kafka_instance_name"
  kafka_topic_name                 = "your_kafka_topic_name"
  kafka_broker_list                = "broker1:9092,broker2:9092,broker3:9092"
  kafka_instance_storage_spec_code = "dms.physical.storage.high.v2"
  kafka_instance_engine_version    = "2.7"
  kafka_instance_storage_space     = 600
  kafka_instance_broker_num        = 3
  kafka_instance_user_name         = "user"
  kafka_instance_user_password     = "YourPassword123"
  ```

* Initialize Terraform:

  ```bash
  terraform init
  ```

* Review the Terraform plan:

  ```bash
  terraform plan
  ```

* Apply the configuration:

  ```bash
  terraform apply
  ```

* To clean up the resources:

  ```bash
  terraform destroy
  ```

## Note

* Make sure to keep your credentials secure and never commit them to version control.
* All resources will be created in the specified region.
* The APIG instance is created with BASIC edition by default.
* The APIG instance will be deployed in the first available zone if `availability_zones` is not specified.
* The Kafka instance flavor will be automatically selected if `kafka_instance_flavor_id` is not provided.
* The plugin type should be "kafka_log". Please verify the correct plugin type name in the APIG documentation or by
  querying available plugin types before using this example.
* The plugin converts HTTP requests to Kafka messages in JSON format, including request metadata, body, and context
  information.
* For authenticated Kafka connections, you can use `kafka_sasl_username` and `kafka_sasl_password`, or leave them empty
  to use `kafka_access_user` and `kafka_password` when security_protocol is not PLAINTEXT.
* Please be aware of API Gateway and DMS Kafka service quotas in your HuaweiCloud account.

## Requirements

| Name        | Version   |
|-------------|-----------|
| terraform   | >= 1.1.0 |
| huaweicloud | >= 1.77.7 |
