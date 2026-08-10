# Create a GA Access Log

This example provides best practice code for using Terraform to create a GA (Global Accelerator) access log
within HuaweiCloud. The access log sends GA listener access logs to LTS (Log Tank Service).

The example creates the full dependency chain:

* GA Accelerator (with IPV4 and IPV6 IP sets)
* GA Listener (with port ranges)
* LTS Log Group
* LTS Log Stream
* GA Access Log (associated with the listener and LTS resources)

-> Currently, the GA flow log interconnects with LTS only in the **cn-north-4** region.

## Prerequisites

* A HuaweiCloud account
* Terraform installed
* HuaweiCloud access key and secret key (AK/SK)

## Variable Introduction

The following variables need to be configured:

### Authentication Variables

* `region_name` - The region where GA resources will be created
(default: "cn-north-4", as LTS interconnect is only available in this region)
* `access_key` - The access key of the IAM user
* `secret_key` - The secret key of the IAM user

### Resource Variables

#### Required Variables

* `accelerator_name` - The name of the GA accelerator
* `listener_name` - The name of the GA listener
* `lts_group_name` - The name of the LTS log group
* `lts_stream_name` - The name of the LTS log stream

#### Optional Variables

* `accelerator_description` - The description of the GA accelerator (default: "")
* `ip_area` - The area of the IP address, valid values: CM, CT, CU, EU, AP, AF, ME, GE (default: "CM")
* `listener_protocol` - The protocol of the GA listener, valid values: TCP, UDP (default: "TCP")
* `listener_description` - The description of the GA listener (default: "")
* `port_from` - The start port of the listener port range (default: 4000)
* `port_to` - The end port of the listener port range (default: 4200)
* `tags` - The tags of the GA accelerator and listener (default: {})
* `lts_ttl_in_days` - The TTL in days for the LTS log group (default: 30)

## Usage

* Copy this example script to your `main.tf`.

* Create a `terraform.tfvars` file and fill in the required variables:

  ```hcl
  accelerator_name = "ga-accelerator-test"
  listener_name    = "ga-listener-test"
  lts_group_name   = "ga-lts-group"
  lts_stream_name  = "ga-lts-stream"
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
* The `resource_type` is set to "LISTENER" as it is the only supported value
* The GA access log must be created in the **cn-north-4** region where LTS interconnect is available
* The `log_group_id` and `log_stream_id` can be updated after creation
* All resources will be created in the specified region

## Requirements

| Name | Version  |
| ---- |----------|
| terraform | >= 0.14.0 |
| huaweicloud | >= 1.72.1 |
