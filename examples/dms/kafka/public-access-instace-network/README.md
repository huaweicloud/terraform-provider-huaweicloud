# Configure network for DMS Kafka instance public access

This example provides best practice code for configuring network for DMS Kafka instance public access in HuaweiCloud,
including EIP allocation and security group configuration for public network connectivity.

## Prerequisites

* A HuaweiCloud account
* Terraform installed
* HuaweiCloud access key and secret key (AK/SK)

## Variables Introduction

The following variables need to be configured:

### Authentication Variables

* `region_name` - The region where the Kafka instance is located
* `access_key`  - The access key of the IAM user
* `secret_key`  - The secret key of the IAM user

### Resource Variables

#### Required Variables

* `vpc_name` - The VPC name
* `subnet_name` - The subnet name
* `security_group_name` - The security group name
* `instance_name` - The Kafka instance name
* `bandwidth_name` - The name of the bandwidth
* `security_group_rule_remote_ip_prefix` - The remote IP prefix of the security group rule  
  The IP address or address group of the client that is allowed to access the Kafka instance

#### Optional Variables

* `availability_zones` - The availability zones to which the Kafka instance belongs (default: [])
* `vpc_cidr` - The CIDR block of the VPC (default: "192.168.0.0/16")
* `subnet_cidr` - The CIDR block of the subnet (default: "")
* `subnet_gateway_ip` - The gateway IP of the subnet (default: "")
* `security_group_rule_ports` - The ports of the security group rule (default: "9094,9095")
  - **9094**: Plaintext access
  - **9095**: Encrypted access
* `instance_flavor_id` - The flavor ID of the Kafka instance (default: "")
* `instance_flavor_type` - The flavor type of the Kafka instance (default: "cluster")
* `instance_storage_spec_code` - The storage specification code of the Kafka instance (default: "dms.physical.storage.ultra.v2")
* `instance_broker_num` - The number of brokers of the Kafka instance (default: 3)
* `eip_type` - The type of the EIP (default: "5_bgp")
* `bandwidth_size` - The size of the bandwidth (default: 5)
* `bandwidth_share_type` - The share type of the bandwidth (default: "PER")
* `bandwidth_charge_mode` - The charge mode of the bandwidth (default: "traffic")
* `instance_engine_version` - The engine version of the Kafka instance (default: "2.7")
* `instance_storage_space` - The storage space of the Kafka instance (default: 600)
* `instance_description` - The description of the Kafka instance (default: "")
* `instance_access_user_name` - The access user of the Kafka instance (default: null)
* `instance_access_user_password` - The access password of the Kafka instance (default: null)
* `instance_enabled_mechanisms` - The enabled mechanisms of the Kafka instance (default: null)  
  The valid values are **PLAIN** and **SCRAM-SHA-512**
* `instance_public_plain_enable` - Whether to enable public plaintext access (default: true)
* `instance_public_sasl_ssl_enable` - Whether to enable public SASL SSL access (default: false)
* `instance_public_sasl_plaintext_enable` - Whether to enable public SASL plaintext access (default: false)

## Usage

* Copy this example script to your `main.tf`.

* Create a `terraform.tfvars` file and fill in the required variables:

  ```hcl
  vpc_name                             = "your_vpc_name"
  subnet_name                          = "your_subnet_name"
  security_group_name                  = "your_security_group_name"
  security_group_rule_remote_ip_prefix = "your_client_ip_address"
  instance_name                        = "your_kafka_instance"
  bandwidth_name                       = "your_bandwidth_name"
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

## Configuration Options

This example supports three public access configuration options:

1. **Plaintext Access**
   - Uses port `9094`
   - No encryption, suitable for internal networks or testing
   - Fastest performance but least secure

   ```hcl
   instance_public_plain_enable = true
   ```

2. **SASL SSL Access**
   - Uses port `9095`
   - Uses SASL authentication with SSL certificate encryption for data transmission
   - Recommended for production environments

   ```hcl
   instance_public_plain_enable    = false
   instance_access_user_name       = "admin"
   instance_access_user_password   = "YourPassword123!"
   instance_enabled_mechanisms     = ["SCRAM-SHA-512"]
   instance_public_sasl_ssl_enable = true
   ```

3. **SASL Plaintext Access**
   - Uses port `9094`
   - Uses SASL authentication with plaintext data transmission
   - Balance between security and performance

   ```hcl
   instance_public_plain_enable          = false
   instance_access_user_name             = "admin"
   instance_access_user_password         = "YourPassword123!"
   instance_enabled_mechanisms           = ["SCRAM-SHA-512"]
   instance_public_sasl_plaintext_enable = true
   ```

## Note

* Make sure to keep your credentials secure and never commit them to version control
* The creation of the Kafka instance takes about 20 to 50 minutes depending on the flavor and broker number
* All resources will be created in the specified region
* EIPs are allocated for each broker in the Kafka instance to enable public access (the number of binded EIPs must
  match the broker count)
* Security group rules must be configured to allow traffic on Kafka ports (`9094` for plaintext, `9095` for encrypted access)
* Public access protocols can be configured based on your security requirements
* This configuration enables external clients to connect to Kafka brokers over the internet

## Requirements

| Name | Version |
| ---- | ---- |
| terraform | >= 1.1.0 |
| huaweicloud | >= 1.77.7 |
