# Route messages by event subscription

This example provides best practice code for using Terraform to create Event Grid (EG) event subscription in
HuaweiCloud for routing messages from OBS service to Kafka instance.

## Prerequisites

* A HuaweiCloud account
* Terraform installed
* HuaweiCloud access key and secret key (AK/SK)

## Variable Introduction

The following variables need to be configured:

### Authentication Variables

* `region_name` - The region where the VPC is located
* `access_key` - The access key of the IAM user
* `secret_key` - The secret key of the IAM user

### Resource Variables

#### Required Variables

* `vpc_name` - The name of the VPC
* `subnet_name` - The name of the subnet
* `security_group_name` - The name of the security group
* `bucket_name` - The name of the OBS bucket
* `instance_name` - The name of the Kafka instance
* `topic_name` - The name of the topic
* `connection_name` - The name of the connection
* `object_name` - The name of the OBS object to be uploaded
* `object_upload_content` - The content of the OBS object to be uploaded

#### Optional Variables

* `vpc_cidr` - The CIDR block of the VPC (default: "172.16.0.0/16")
* `subnet_cidr` - The CIDR block of the subnet (default: "172.16.10.0/24")
* `subnet_gateway` - The gateway IP address of the subnet (default: "172.16.10.1")
* `bucket_acl` - The ACL policy for a bucket (default: "private")
* `availability_zones` - The availability zones to which the Kafka instance belongs (default: [])
* `instance_flavor_id` - The flavor ID of the Kafka instance (default: "kafka.2u4g.cluster.small")
* `instance_flavor_type` - The flavor type of the Kafka instance (default: "cluster")
* `instance_storage_spec_code` - The storage specification code of the Kafka instance
  (default: "dms.physical.storage.high.v2")
* `instance_engine_version` - The engine version of the Kafka instance (default: "3.x")
* `instance_storage_space` - The storage space of the Kafka instance (default: 300)
* `instance_broker_num` - The number of brokers of the Kafka instance (default: 3)
* `instance_ssl_enable` - The SSL enable of the Kafka instance (default: false)
* `instance_description` - The description of the Kafka instance (default: "")
* `instance_security_protocol` - The protocol to use after SASL is enabled (default: "SASL_SSL")
* `charging_mode` - The charging mode of the Kafka instance (default: "postPaid")
* `topic_partitions` - The number of the topic partition (default: 3)
* `subscription_source_values` - The event types to be subscribed from OBS service (default: ["OBS:CloudTrace:ApiCall",
  "OBS:CloudTrace:ObsSDK", "OBS:CloudTrace:ConsoleAction", "OBS:CloudTrace:SystemAction", "OBS:CloudTrace:Others"])
* `connection_type` - The type of the connection (default: "KAFKA")
* `connection_acks` - The number of confirmation signals the prouder needs to receive to consider the message sent
  successfully (default: "1")
* `object_extension_name` - The extension name of the OBS object to be uploaded

## Usage

* Copy this example script to your `main.tf`.

* Create a `terraform.tfvars` file and fill in the required variables:

  ```hcl
  vpc_name              = "your_vpc_name"
  subnet_name           = "your_subnet_name"
  security_group_name   = "your_security_group_name"
  bucket_name           = "your_bucket_name"
  instance_name         = "your_kafka_instance_name"
  topic_name            = "your_topic_name"
  object_name           = "your_object_name"
  object_upload_content = <<EOT
  def main():
      print("Hello, World!")

  if __name__ == "__main__":
      main()
  EOT
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
* The creation of resources may take several minutes depending on the configuration
* This example creates a complete infrastructure including VPC, OBS, Kafka, and EG resources
* The event subscription automatically filters OBS events and routes them to Kafka
* Kafka connection details are automatically retrieved and configured
* The OBS object upload triggers events that are captured by the subscription
* All resources will be created in the specified region

## Requirements

| Name        | Version   |
|-------------|-----------|
| terraform   | >= 1.1.0  |
| huaweicloud | >= 1.68.0 |
| time        | ~> 0.13   |
