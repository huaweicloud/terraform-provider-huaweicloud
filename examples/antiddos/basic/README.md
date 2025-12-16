# Create a basic configuration for Anti-DDoS

This example provides best practice code for using Terraform to create a basic configuration in HuaweiCloud
Anti-DDoS service.

## Prerequisites

* A HuaweiCloud account
* Terraform installed
* HuaweiCloud access key and secret key (AK/SK)

## Variable Introduction

The following variables need to be configured:

### Authentication Variables

* `region_name` - The region where the Anti-DDoS basic configuration is located
* `access_key`  - The access key of the IAM user
* `secret_key`  - The secret key of the IAM user

### Resource Variables

#### Required Variables

* `vpc_eip_publicip_type` - The EIP type. Possible values are **5_bgp** (dynamic BGP) and **5_sbgp** (static BGP)
* `vpc_eip_bandwidth_share_type` - The bandwidth share type. Possible values are **PER** (dedicated bandwidth) and
  **WHOLE** (shared bandwidth)
* `smn_topic_name` - The name of the topic to be created
* `smn_subscription_endpoint` - The message endpoint
* `smn_subscription_protocol` - The protocol of the message endpoint
* `antiddos_traffic_threshold` - The traffic cleaning threshold in Mbps

#### Optional Variables

* `vpc_eip_bandwidth_name` - The bandwidth name (required when `vpc_eip_bandwidth_share_type` is **PER**) (default: null)
* `vpc_eip_bandwidth_size` - The bandwidth size (required when `vpc_eip_bandwidth_share_type` is **PER**) (default: null)
* `vpc_eip_bandwidth_charge_mode` - The bandwidth charge mode (default: null)
* `smn_topic_display_name` - The topic display name (default: null)
* `smn_subscription_remark` - The remark information (default: null)

## Usage

* Copy this example script to your `main.tf`.

* Create a `terraform.tfvars` file and fill in the required variables:

```hcl
vpc_eip_publicip_type         = "5_bgp"
vpc_eip_bandwidth_share_type  = "PER"
vpc_eip_bandwidth_name        = "test-antiddos-basic-name"
vpc_eip_bandwidth_size        = 5
vpc_eip_bandwidth_charge_mode = "traffic"
smn_topic_name                = "test-antiddos-basic-name"
smn_subscription_endpoint     = "mailtest@gmail.com"
smn_subscription_protocol     = "email"
antiddos_traffic_threshold    = 200
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
| terraform | >= 0.14.0 |
| huaweicloud | >= 1.60.1 |
