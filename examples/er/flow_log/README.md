# Create flow log with ER instance

This example provides best practice code for using Terraform to create flow log in HuaweiCloud ER service.

## Prerequisites

* A HuaweiCloud account with ER permissions
* Terraform installed
* HuaweiCloud access key and secret key (AK/SK)

## Variable Introduction

The following variables need to be configured:

### Authentication Variables

* `region_name` - The region where the ER service is located
* `access_key` - The access key of the IAM user
* `secret_key` - The secret key of the IAM user

### Resource Variables

#### Required Variables

* `vpc_name` - The VPC name.
* `subnet_name` - The subnet name.
* `er_instance_name` - The ER instance name.
* `er_instance_asn` - The ER instance asn.
* `er_vpc_attachment_name` - The ER VPC attachment name.
* `lts_group_name` - The LTS group name.
* `lts_stream_name` - The LTS stream name.
* `er_flow_log_name` - The ER flow log name.

#### Optional Variables

* `vpc_cidr` - The CIDR block of the VPC (default: "192.168.0.0/16")
* `subnet_cidr` - The CIDR block of the subnet (default: "")
* `subnet_gateway_ip` - The gateway IP of the subnet (default: "")
* `er_vpc_attachment_auto_create_vpc_routes` - Whether to enable auto create VPC routes (default: true)
* `lts_group_ttl_in_days` - The log expiration time(days) (default: 7)
* `er_flow_log_store_type` - The storage type of flow log (default: "LTS")
* `er_flow_log_resource_type` - The resource type to which the logs to be collected (default: "attachment")

## Usage

* Copy this example script to your `main.tf`.

* Create a `terraform.tfvars` file and fill in the required variables:

  ```hcl
  vpc_name               = "tf_test_er_instance_vpc"
  subnet_name            = "tf_test_er_instance_subnet"
  er_instance_name       = "tf_test_er_instance"
  er_instance_asn        = 64512
  er_vpc_attachment_name = "tf_test_er_vpc_attachment"
  lts_group_name         = "tf_test_lts_group"
  lts_stream_name        = "tf_test_lts_stream"
  er_flow_log_name       = "tf_test_er_flow_log"
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
| terraform | >= 0.12.0 |
| huaweicloud | >= 1.73.5 |
