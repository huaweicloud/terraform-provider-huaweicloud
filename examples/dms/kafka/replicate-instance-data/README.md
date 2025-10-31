# Create a Smart Connect task to replicate data between two Kafka instances

This example provides best practice code for using Terraform to create a Smart Connect task to replicate data between
two Kafka instances in HuaweiCloud.

## Prerequisites

* A HuaweiCloud account
* Terraform installed
* HuaweiCloud access key and secret key (AK/SK)

## Variable Introduction

The following variables need to be configured:

### Authentication Variables

* `region_name` - The region where the Kafka instance is located
* `access_key` - The access key of the IAM user
* `secret_key` - The secret key of the IAM user

### Resource Variables

#### Required Variables

* `instance_configurations` - The list of configurations for multiple Kafka instances
  - `name` - The instance name
  - `availability_zones` - The availability zones of the instance (default: [])
  - `engine_version` - The engine version (default: "3.x")
  - `flavor_id` - The flavor ID (default: "")
  - `flavor_type` - The flavor type (default: "cluster")
  - `storage_spec_code` - The storage specification code (default: "dms.physical.storage.ultra.v2")
  - `storage_space` - The storage space in GB (default: 600)
  - `broker_num` - The number of brokers (default: 3)
  - `access_user` - The access user name (default: "")
  - `password` - The access user password (default: "")
  - `enabled_mechanisms` - The enabled SASL mechanisms (default: null)
  - `port_protocol` - The port protocol configuration (default: {})
    + `private_plain_enable` - Enable private plaintext access (default: true)
    + `private_sasl_ssl_enable` - Enable private SASL SSL access (default: null)
    + `private_sasl_plaintext_enable` - Enable private SASL plaintext access (default: null)
* `vpc_name` - The name of the VPC
* `subnet_name` - The name of the subnet
* `security_group_name` - The name of the security group
* `task_name` - The name of the Smart Connect task

#### Optional Variables

* `vpc_cidr` - The CIDR block of the VPC (default: "192.168.0.0/16")
* `subnet_cidr` - The CIDR block of the subnet (default: "")
* `subnet_gateway_ip` - The gateway IP of the subnet (default: "")
* `task_topics` - The topics of the Smart Connect task (default: [])
* `topic_name` - The name of the Kafka topic (default: "")
* `topic_partitions` - The number of partitions of the topic (default: 10)
* `topic_replicas` - The number of replicas of the topic (default: 3)
* `topic_aging_time` - The aging time of the topic (default: 72)
* `topic_sync_replication` - The sync replication of the topic (default: false)
* `topic_sync_flushing` - The sync flushing of the topic (default: false)
* `topic_description` - The description of the topic (default: null)
* `topic_configs` - The configs of the topic (default: [])
* `smart_connect_storage_spec_code` - The storage specification code of the Smart Connect (default: null)
* `smart_connect_bandwidth` - The bandwidth of the Smart Connect (default: null)
* `smart_connect_node_count` - The number of nodes of the Smart Connect (default: 2)
* `task_start_later` - The start later of the Smart Connect task (default: false)
* `task_direction` - The direction of the Smart Connect task (default: "two-way")
* `task_replication_factor` - The replication factor of the Smart Connect task (default: 3)
* `task_task_num` - The number of tasks of the Smart Connect task (default: 2)
* `task_provenance_header_enabled` - The provenance header enabled of the Smart Connect task (default: false)
* `task_sync_consumer_offsets_enabled` - The sync consumer offsets enabled of the Smart Connect task (default: false)
* `task_rename_topic_enabled` - The rename topic enabled of the Smart Connect task (default: true)
* `task_consumer_strategy` - The consumer strategy of the Smart Connect task (default: "latest")
* `task_compression_type` - The compression type of the Smart Connect task (default: "none")
* `task_topics_mapping` - The topics mapping of the Smart Connect task (default: [])
* `task_peer_instance_access_user_password` - The access password of the peer Kafka instance (default: null)

## Usage

* Copy this example script to your `main.tf`.

* Create a `terraform.tfvars` file and fill in the required variables:

  ```hcl
  vpc_name                = "tf_test_vpc"
  subnet_name             = "tf_test_subnet"
  security_group_name     = "tf_test_security_group"
  instance_configurations = [
    {
      name = "source-instance"
    },
    {
      name = "peer-instance"
    }
  ]

  topic_name = "tf_test_topic"
  task_name  = "tf_test_task"
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
* At least two instances are required. The Smart Connect task uses instance[0] as current and instance[1] as peer
* Single-node Kafka instances do not support creating Smart Connect tasks for Kafka data replication
* After a Smart Connect task is created, changing the peer instance authentication method, mechanism, or password will
  cause the sync task to fail. You need to delete the current Smart Connect task and recreate a new one

## Requirements

| Name | Version |
| ---- | ---- |
| terraform | >= 1.9.0 |
| huaweicloud | >= 1.77.7 |
